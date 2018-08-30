package conchapp

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/blockchainworkers/conch/libs/log"
	"github.com/jmoiron/sqlx"
	"math/big"
	"sync"
)

// AccountState means current account's info
type AccountState struct {
	sync.RWMutex
	accounts map[string]*Account
	log      log.Logger
	db       *sqlx.DB
	isDirty  bool
}

// NewAccountState return AccountState inst
func NewAccountState(db *sqlx.DB, log log.Logger) *AccountState {
	return &AccountState{
		accounts: make(map[string]*Account),
		log:      log,
		db:       db,
		isDirty:  false,
	}
}

// LoadAccount get account from cache or db
func (as *AccountState) LoadAccount(address string) (*Account, error) {
	as.RLock()
	defer as.RUnlock()
	if k, ok := as.accounts[address]; ok {
		return k, nil
	}
	// try load from db
	acc, err := as.selectAccount(address)
	if err == nil {
		as.accounts[address] = acc
	}
	return acc, err
}

func (as *AccountState) selectAccount(addr string) (*Account, error) {
	sqlStr := "select amount from funds where address = ?"
	var amount string
	err := as.db.QueryRowx(sqlStr, addr).Scan(&amount)
	if err == sql.ErrNoRows {
		return &Account{Address: addr, Amount: new(big.Int).SetInt64(0)}, nil
	}
	if err != nil {
		return nil, err
	}
	acc := Account{Address: addr, Amount: big.NewInt(0)}
	acc.Amount.SetString(amount, 0)
	return &acc, nil
}

// UpdateAccountCache update account in memory
func (as *AccountState) UpdateAccountCache(acc *Account) {
	as.RLock()
	defer as.RUnlock()
	as.accounts[acc.Address] = acc
	as.isDirty = true
}

//SyncToDisk cache to disk
func (as *AccountState) SyncToDisk() error {
	if !as.isDirty {
		return nil
	}
	if len(as.accounts) == 0 {
		return nil
	}

	sqlStr := "replace into funds(address, amount) values "
	for _, val := range as.accounts {
		sqlStr = sqlStr + fmt.Sprintf(" ('%s', '%s'),", val.Address, val.Amount.String())
	}
	sqlStr = sqlStr[0 : len(sqlStr)-2]
	_, err := as.db.Exec(sqlStr)
	as.accounts = make(map[string]*Account)
	return err
}

//Account ...
type Account struct {
	Address string
	Amount  *big.Int
}

// NewAccount return account inst
func NewAccount(addr, amount string) *Account {
	am := big.NewInt(0)
	am.SetString(amount, 0)
	return &Account{Address: addr, Amount: am}
}

// TxState means current tx's info
type TxState struct {
	sync.RWMutex
	Txs Transactions
	log log.Logger
	db  *sqlx.DB
}

// NewTxState txstate inst
func NewTxState(db *sqlx.DB, log log.Logger) *TxState {
	return &TxState{
		Txs: Transactions{},
		log: log,
		db:  db,
	}
}

// UpdateTx append tx
func (txState *TxState) UpdateTx(tx *Transaction) {
	txState.Lock()
	defer txState.Unlock()
	txState.Txs = txState.Txs.AppendTx(tx)
}

// SyncToDisk write tx to db
func (txState *TxState) SyncToDisk(height int64) (hashRoot string, err error) {
	if txState.Txs.Len() == 0 {
		return txState.Txs.HashRoot(), nil
	}

	sqlStr := "replace into transaction_records(id, sender, receiver, amount, input, expired, time_stamp, nonce, ref_block_num, block_num, sign) values "
	for _, val := range txState.Txs {
		sqlStr = sqlStr + fmt.Sprintf(" ('%s', '%s', '%s', '%s', '%s', '%d', '%d', '%s', '%d', '%d', '%s'),",
			val.TxID(), val.Sender, val.Receiver, val.Value, val.Input, val.ExpiredNum, val.TimeStamp, val.Nonce, val.RefBlockNum, height, val.Sign)
	}
	sqlStr = sqlStr[0 : len(sqlStr)-2]
	_, err = txState.db.Exec(sqlStr)
	// merkle tree
	hashRoot = txState.Txs.HashRoot()

	// new trans for next commit
	txState.RLock()
	txState.Txs = Transactions{}
	txState.RUnlock()
	return
}

// TxRepState means current tx receipt's info
type TxRepState struct {
	sync.RWMutex
	Txreps TransactionReceipts
	log    log.Logger
	db     *sqlx.DB
}

// NewTxRepState tx receipt inst
func NewTxRepState(db *sqlx.DB, log log.Logger) *TxRepState {
	return &TxRepState{
		Txreps: TransactionReceipts{},
		log:    log,
		db:     db,
	}
}

// UpdateTxRep append tx
func (txrSt *TxRepState) UpdateTxRep(tr *TransactionReceipt) {
	txrSt.Lock()
	defer txrSt.Unlock()
	txrSt.Txreps = txrSt.Txreps.AppendTxrp(tr)
}

// SyncToDisk write tx to db
func (txrSt *TxRepState) SyncToDisk(height int64) (hashRoot string, err error) {
	if txrSt.Txreps.Len() == 0 {
		return txrSt.Txreps.HashRoot(), nil
	}
	// id | status | fee | block_num | tx_hash | log

	sqlStr := "replace into transaction_receipts (id, status, fee, block_num, tx_hash, log) values "
	for _, val := range txrSt.Txreps {
		sqlStr = sqlStr + fmt.Sprintf(" ('%s', '%d', '%s', '%d', '%s', '%s'),",
			string(val.Hash()), val.Status, val.Fee.String(), height, val.TxHash, val.Log)
	}
	sqlStr = sqlStr[0 : len(sqlStr)-2]
	_, err = txrSt.db.Exec(sqlStr)
	// merkle tree
	hashRoot = txrSt.Txreps.HashRoot()

	// new trans for next commit
	txrSt.RLock()
	txrSt.Txreps = TransactionReceipts{}
	txrSt.RUnlock()
	return
}

// HeaderState appheader state
type HeaderState struct {
	CurBlockNum  int64  `json:"cur_block_num"`
	CurBlockHash string `json:"cur_block_hash"`
	CurAPPHash   string `json:"cur_app_hash"`
	db           *sqlx.DB
	log          log.Logger
}

// LoadHeaderState from db load header
func (hdSt *HeaderState) LoadHeaderState() error {
	sqlStr := "select content from state where id=1"
	var text string
	err := hdSt.db.QueryRowx(sqlStr).Scan(&text)
	if err == sql.ErrNoRows {
		return nil
	}
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(text), hdSt)
}

// SyncToDisk to db
func (hdSt *HeaderState) SyncToDisk() error {
	dat, err := json.Marshal(hdSt)
	if err != nil {
		return err
	}

	sqlStr := fmt.Sprintf("replace into state (id, content) value ('%d', '%s')", 1, string(dat))
	_, err = hdSt.db.Exec(sqlStr)

	return err
}

// APPState state set
type APPState struct {
	HeadSt   *HeaderState
	AccoutSt *AccountState
	TxSt     *TxState
	TxRepSt  *TxRepState
}

//NewAPPState return app state init db (if db is not exist create the database and tables)
func NewAPPState(db *sqlx.DB, log log.Logger) *APPState {
	return &APPState{
		HeadSt:   &HeaderState{db: db, log: log},
		AccoutSt: NewAccountState(db, log),
		TxSt:     NewTxState(db, log),
		TxRepSt:  NewTxRepState(db, log),
	}
}

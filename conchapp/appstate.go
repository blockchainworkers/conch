package conchapp

import (
	"database/sql"
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
	log      *log.Logger
	db       *sqlx.DB
	isDirty  bool
}

// NewAccountState return AccountState inst
func NewAccountState(db *sqlx.DB, log *log.Logger) *AccountState {
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
	accounts map[string]*Account
	log      *log.Logger
	db       *sqlx.DB
	isDirty  bool
}

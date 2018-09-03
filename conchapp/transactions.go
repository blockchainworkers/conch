package conchapp

import (
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/blockchainworkers/conch/account"
	"github.com/blockchainworkers/conch/crypto"
	"github.com/blockchainworkers/conch/crypto/merkle"
	"github.com/blockchainworkers/conch/crypto/secp256k1"
	"math/big"
	"sort"
)

var nilHash = "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"

// Transactions list of Transaction
type Transactions []*Transaction

//Len return length
func (txs Transactions) Len() int { return len(txs) }

//HashRoot merkle root
func (txs Transactions) HashRoot() string {
	hasers := make([]merkle.Hasher, 0, txs.Len())
	for i := range txs {
		hasers = append(hasers, (txs)[i])
	}
	h := merkle.SimpleHashFromHashers(hasers)
	if h == nil {
		return nilHash
	}
	return hex.EncodeToString(h)
}

// AppendTx append an tx
func (txs Transactions) AppendTx(tx *Transaction) Transactions {
	txs = append(txs, tx)
	return txs
}

// Transaction tx type
type Transaction struct {
	ID          string `json:"id"`
	Sender      string `json:"sender"`
	Receiver    string `json:"receiver"`
	Input       string `json:"input"`
	Sign        string `json:"sign"`
	Value       string `json:"value"`
	TimeStamp   int64  `json:"time"`
	RefBlockNum int64  `json:"ref_block"`
	Nonce       string `json:"nonce"`
	ExpiredNum  int    `json:"expired"`
	Cache       struct {
		code []byte // to all field to encode
		hash []byte // hash(code)
		sign []byte // sign(hash)
		id   string // hex(hash)
	} `json:"-"`
}

// DecodeNewTx decode a new tx
func DecodeNewTx(data []byte) (*Transaction, error) {
	var tx Transaction
	dat, err := base64.StdEncoding.DecodeString(string(data))
	if err != nil {
		return &tx, err
	}

	err = json.Unmarshal(dat, &tx)
	return &tx, err
}

// BuildNewTx create a new tx
func BuildNewTx(sender, receive, input, nonce, value string, time, refBlock, expired int64) *Transaction {
	return &Transaction{
		Sender:      sender,
		Receiver:    receive,
		Input:       input,
		Nonce:       nonce,
		TimeStamp:   time,
		Value:       value,
		RefBlockNum: refBlock,
		ExpiredNum:  int(expired),
	}
}

//SignTx sign tx
func (tx *Transaction) SignTx(privKey crypto.PrivKey) (string, error) {
	sign, err := tx.signCache(privKey)
	if err != nil {
		return "", err
	}
	tx.Sign = hex.EncodeToString(sign)
	return tx.Sign, err
}

// IsValidTx check tx valid or not
func (tx *Transaction) IsValidTx() bool {
	// 1. from sign and content to recover public key
	// 2. from public key to generate address
	// 3. check the address qeual sender or not
	if !tx.CheckArgs() {
		return false
	}

	// todo:: only secp256k1 support recover publickey
	msg := tx.hashCache()
	signature, err := hex.DecodeString(tx.Sign)
	if err != nil {
		return false
	}
	pub, err := secp256k1.RecoverPublicKey(signature, msg)
	if err != nil {
		return false
	}
	return account.PublicKeyToAddress(pub) == tx.Sender
}

// CheckArgs tx's args is vaild
func (tx *Transaction) CheckArgs() bool {
	if !account.CheckAddressValid(tx.Sender) {
		return false
	}
	if tx.Nonce == "" {
		return false
	}
	if tx.RefBlockNum == 0 {
		return false
	}
	if tx.Sign == "" {
		return false
	}
	if _, r := new(big.Int).SetString(tx.Value, 0); !r {
		return false
	}
	return true
}

// TxID return tx's unique hash value
func (tx *Transaction) TxID() string {
	if tx.Cache.id == "" {
		tx.Cache.id = hex.EncodeToString(tx.hashCache())
	}
	tx.ID = tx.Cache.id
	return tx.Cache.id
}

// Hash return tx's unique hash value
func (tx *Transaction) Hash() []byte {
	return tx.hashCache()
}

// FeeCalc calc tx fee
func (tx *Transaction) FeeCalc() *big.Int {
	// if tx.input is empty we only charge an base fee 100 gravel
	// now 1 conch == 10**8 gravel

	if len(tx.Input) < 65 {
		return big.NewInt(100)
	}
	baseCharge := big.NewInt(100)
	inputCharg := big.NewInt(10)
	inputCharg = inputCharg.Mul(inputCharg, big.NewInt(int64(len(tx.Input))))
	return baseCharge.Add(baseCharge, inputCharg)
}

// Serialization json tx and base64
func (tx *Transaction) Serialization() string {
	// if tx.input is empty we only charge an base fee 100 gravel
	// now 1 conch == 10**8 gravel
	dat, _ := json.Marshal(tx)
	return base64.StdEncoding.EncodeToString(dat)
}

// -------priv func ----

func (tx *Transaction) signCache(privKey crypto.PrivKey) ([]byte, error) {
	if tx.Cache.sign != nil {
		return tx.Cache.sign, nil
	}
	sign, err := privKey.(secp256k1.PrivKeySecp256k1).SignCompact(tx.hashCache())
	if err == nil {
		tx.Cache.sign = sign
	}
	return sign, err
}

func (tx *Transaction) hashCache() []byte {
	if tx.Cache.hash != nil {
		return tx.Cache.hash
	}
	tx.Cache.hash = crypto.Sha256(tx.codeCache())
	return tx.Cache.hash
}

func (tx *Transaction) codeCache() []byte {
	if tx.Cache.code != nil {
		return tx.Cache.code
	}
	// encode tx
	dat := tx.FormCode()
	buf := make([]byte, base64.StdEncoding.EncodedLen(len(dat)))
	base64.StdEncoding.Encode(buf, dat)
	tx.Cache.code = buf
	return tx.Cache.code
}

// FormCode organize tx field content then marshal
func (tx *Transaction) FormCode() []byte {
	val := make(map[string]interface{})
	keys := []string{"sender", "receiver", "value", "input", "time", "ref_block", "nonce", "expired"}
	val["sender"] = tx.Sender
	val["receiver"] = tx.Receiver
	val["input"] = tx.Input
	val["time"] = tx.TimeStamp
	val["ref_block"] = tx.RefBlockNum
	val["nonce"] = tx.Nonce
	val["expired"] = tx.ExpiredNum
	val["value"] = tx.Value

	sort.Slice(keys, func(i, j int) bool { return keys[i] < keys[j] })
	dat := ""
	for iter := range keys {
		kv := fmt.Sprintf("%v=%v", keys[iter], val[keys[iter]])
		if dat == "" {
			dat = kv
		} else {
			dat = dat + "&" + kv
		}
		continue
	}
	return []byte(dat)
}

// TransactionReceipt tx receipt
type TransactionReceipt struct {
	Status    int
	Fee       *big.Int
	BlockNum  int64
	TxHash    string
	Log       string
	hashCache []byte
}

// Hash return hash
func (txRep *TransactionReceipt) Hash() []byte {
	if txRep.hashCache != nil {
		return txRep.hashCache
	}
	code := fmt.Sprintf("block_num=%d&fee=%s&status=%d&tx_hash=%s&log=%s",
		txRep.BlockNum, txRep.Fee.String(), txRep.Status, txRep.TxHash, txRep.Log)
	txRep.hashCache = crypto.Sha256([]byte(code))
	return txRep.hashCache
}

// ID return hash
func (txRep *TransactionReceipt) ID() string {
	hash := txRep.Hash()
	return hex.EncodeToString(hash)
}

// TransactionReceipts trxreps
type TransactionReceipts []*TransactionReceipt

//Len return length
func (txrp TransactionReceipts) Len() int { return len(txrp) }

//HashRoot merkle root
func (txrp TransactionReceipts) HashRoot() string {
	hasers := make([]merkle.Hasher, 0, txrp.Len())
	for i := range txrp {
		hasers = append(hasers, txrp[i])
	}
	h := merkle.SimpleHashFromHashers(hasers)
	if h == nil {
		return nilHash
	}
	return hex.EncodeToString(h)
}

// AppendTxrp append an tx receipt
func (txrp TransactionReceipts) AppendTxrp(trp *TransactionReceipt) TransactionReceipts {
	txrp = append(txrp, trp)
	return txrp
}

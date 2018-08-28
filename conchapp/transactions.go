package conchapp

import (
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/blockchainworkers/conch/crypto"
	"github.com/blockchainworkers/conch/crypto/merkle"
	"github.com/blockchainworkers/conch/crypto/secp256k1"
	"math/big"
	"sort"
)

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
	return hex.EncodeToString(h)
}

// AppendTx append an tx
func (txs Transactions) AppendTx(tx *Transaction) Transactions {
	txs = append(txs, tx)
	return txs
}

// Transaction tx type
type Transaction struct {
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
func DecodeNewTx(date []byte) (*Transaction, error) {
	var tx Transaction
	err := json.Unmarshal(date, &tx)
	return &tx, err
}

// BuildNewTx create a new tx
func BuildNewTx(sender, receive, input, nonce string, time, refBlock, expired int64) *Transaction {
	return &Transaction{
		Sender:      sender,
		Receiver:    receive,
		Input:       input,
		Nonce:       nonce,
		TimeStamp:   time,
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
	return pub.Address().String() == tx.Sender
}

// TxID return tx's unique hash value
func (tx *Transaction) TxID() string {
	if tx.Cache.id == "" {
		tx.Cache.id = hex.EncodeToString(tx.hashCache())
	}
	return tx.Cache.id
}

// Hash return tx's unique hash value
func (tx *Transaction) Hash() []byte {
	if tx.Cache.id == "" {
		tx.Cache.id = hex.EncodeToString(tx.hashCache())
	}
	return []byte(tx.Cache.id)
}

// -------priv func ----

func (tx *Transaction) signCache(privKey crypto.PrivKey) ([]byte, error) {
	if tx.Cache.sign != nil {
		return tx.Cache.sign, nil
	}
	sign, err := privKey.Sign(tx.hashCache())
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
	hashCache string
}

// Hash return hash
func (txRep *TransactionReceipt) Hash() string {
	if txRep.hashCache != "" {
		return txRep.hashCache
	}
	code := fmt.Sprintf("block_num=%d&fee=%s&status=%d&tx_hash=%s&log=%s",
		txRep.BlockNum, txRep.Fee.String(), txRep.Status, txRep.TxHash, txRep.Log)
	txRep.hashCache = hex.EncodeToString(crypto.Sha256([]byte(code)))
	return txRep.hashCache
}

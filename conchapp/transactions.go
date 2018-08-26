package conchapp

import (
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/blockchainworkers/conch/crypto"
	"github.com/blockchainworkers/conch/crypto/secp256k1"
	"sort"
)

// Transaction tx type
type Transaction struct {
	Sender   string `json:"sender"`
	Receiver string `json:"receiver"`
	Input    string `json:"input"`
	Sign     string `json:"sign"`
	// Status      bool   `json:"status"`
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
	keys := []string{"sender", "receiver", "input", "time", "ref_block", "nonce", "expired"}
	val["sender"] = tx.Sender
	val["receiver"] = tx.Receiver
	val["input"] = tx.Input
	val["time"] = tx.TimeStamp
	val["ref_block"] = tx.RefBlockNum
	val["nonce"] = tx.Nonce
	val["expired"] = tx.ExpiredNum

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

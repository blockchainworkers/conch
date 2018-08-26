package conchapp

import (
	"encoding/base64"
	"encoding/json"
	"github.com/blockchainworkers/conch/crypto"
)

// Transaction tx type
type Transaction struct {
	Sender      string `json:"sender"`
	Receiver    string `json:"receiver"`
	Input       string `json:"input"`
	Status      bool   `json:"status"`
	TimeStamp   int64  `json:"time"`
	RefBlockNum int64  `json:"ref_block"`
	Nonce       string `json:"nonce"`
	Cache       struct {
		code []byte // to all field to encode
		hash []byte // hash(code)
		sign []byte // sign(hash)
	} `json:"-"`
}

// DecodeNewTx decode a new tx
func DecodeNewTx(date []byte) (*Transaction, error) {
	var tx Transaction
	err := json.Unmarshal(date, &tx)
	return &tx, err
}

//SignTx sign tx
func (tx *Transaction) SignTx(privKey crypto.PrivKey) ([]byte, error) {
	return tx.signCache(privKey)
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
	dat, _ := json.Marshal(tx)
	buf := make([]byte, base64.StdEncoding.EncodedLen(len(dat)))
	base64.StdEncoding.Encode(buf, dat)
	tx.Cache.code = buf
	return tx.Cache.code
}

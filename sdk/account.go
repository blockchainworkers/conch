package sdkdemo

import (
	"encoding/hex"
	// "github.com/blockchainworkers/conch/crypto"
	"github.com/blockchainworkers/conch/crypto/secp256k1"
	"github.com/btcsuite/btcutil/base58"
)

// AddrPrefix ...
var AddrPrefix = "CONCH"

//AddrVersion base58 check version
var addrVersion = 88

// Account 账户
type Account struct {
	Privkey string
	Address string
}

// NewAccount new
func NewAccount() *Account {
	return &Account{Privkey: "", Address: ""}
}

//GenerateAccout create account
func (acc *Account) GenerateAccout() *Account {
	privkey := secp256k1.GenPrivKey()

	acc.Privkey = hex.EncodeToString([]byte(privkey[:]))
	acc.Address = AddrPrefix + base58.CheckEncode(privkey.PubKey().Address(), byte(addrVersion))
	return acc
}

// LoadPrivKey from string to privkey
func (acc *Account) LoadPrivKey(prikey string) ([]byte, error) {
	privHex, err := hex.DecodeString(prikey)
	var priv secp256k1.PrivKeySecp256k1
	if err != nil {
		return nil, err
	}
	copy(priv[:], privHex)
	tmp := [32]byte(priv)
	return tmp[:], nil
}

//PublicKeyToAddress pub to address
func (acc *Account) PublicKeyToAddress(pubkey string) (string, error) {
	pub, err := hex.DecodeString(pubkey)
	if err != nil {
		return "", err
	}
	var key secp256k1.PubKeySecp256k1
	copy(key[:], pub)
	return AddrPrefix + base58.CheckEncode(key.Address(), byte(addrVersion)), nil
}

// // PrivKeyToAddress priv to address
// func PrivKeyToAddress(priv crypto.PrivKey) string {
// 	return AddrPrefix + base58.CheckEncode(priv.PubKey().Address(), AddrVersion)
// }

// CheckAddressValid address is valid or not
func CheckAddressValid(addr string) bool {
	if len(addr)-len(AddrPrefix) < 0 {
		return false
	}
	_, _, err := base58.CheckDecode(addr[len(AddrPrefix):])
	if err != nil {
		return false
	}
	return true
}

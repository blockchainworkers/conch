package account

import (
	"encoding/hex"
	"github.com/blockchainworkers/conch/crypto"
	"github.com/blockchainworkers/conch/crypto/secp256k1"
	"github.com/btcsuite/btcutil/base58"
)

// AddrPrefix ...
var AddrPrefix = "CONCH"

//AddrVersion base58 check version
var AddrVersion byte = 88

//GenerateAccout create account
func GenerateAccout() (string, string) {
	privkey := secp256k1.GenPrivKey()

	return hex.EncodeToString([]byte(privkey[:])), AddrPrefix + base58.CheckEncode(privkey.PubKey().Address(), AddrVersion)
}

// LoadPrivKey from string to privkey
func LoadPrivKey(prikey string) (crypto.PrivKey, error) {
	privHex, err := hex.DecodeString(prikey)
	var priv secp256k1.PrivKeySecp256k1
	if err != nil {
		return priv, err
	}
	copy(priv[:], privHex)
	return priv, nil
}

// PublicKeyToAddress pub to address
func PublicKeyToAddress(pub crypto.PubKey) string {
	return AddrPrefix + base58.CheckEncode(pub.Address(), AddrVersion)
}

// PrivKeyToAddress priv to address
func PrivKeyToAddress(priv crypto.PrivKey) string {
	return AddrPrefix + base58.CheckEncode(priv.PubKey().Address(), AddrVersion)
}

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

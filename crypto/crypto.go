package crypto

import (
	cmn "github.com/blockchainworkers/conch/libs/common"
)

type PrivKey interface {
	Bytes() []byte
	Sign(msg []byte) ([]byte, error)
	PubKey() PubKey
	Equals(PrivKey) bool
}

// An address is a []byte, but hex-encoded even in JSON.
// []byte leaves us the option to change the address length.
// Use an alias so Unmarshal methods (with ptr receivers) are available too.
type Address = cmn.HexBytes

type PubKey interface {
	Address() Address
	Bytes() []byte
	ByteArray() []byte
	VerifyBytes(msg []byte, sig []byte) bool
	Equals(PubKey) bool
}

type Symmetric interface {
	Keygen() []byte
	Encrypt(plaintext []byte, secret []byte) (ciphertext []byte)
	Decrypt(ciphertext []byte, secret []byte) (plaintext []byte, err error)
}

// CreateAddress creates address given the bytes and the nonce 根据账号地址和nonce创建合约地址
func CreateAddress(b Address, nonce string) Address {
	// todo::
	var buf []byte
	buf = append(buf, b...)
	buf = append(buf, []byte(nonce)...)
	// data, _ := rlp.EncodeToBytes([]interface{}{b, nonce})
	return Sha256(buf)
}

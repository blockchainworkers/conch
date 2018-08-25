package node

import (
	amino "github.com/tendermint/go-amino"
	cryptoAmino "github.com/blockchainworkers/conch/crypto/encoding/amino"
)

var cdc = amino.NewCodec()

func init() {
	cryptoAmino.RegisterAmino(cdc)
}

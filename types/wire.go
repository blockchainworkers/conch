package types

import (
	"github.com/tendermint/go-amino"
	"github.com/blockchainworkers/conch/crypto/encoding/amino"
)

var cdc = amino.NewCodec()

func init() {
	RegisterBlockAmino(cdc)
}

func RegisterBlockAmino(cdc *amino.Codec) {
	cryptoAmino.RegisterAmino(cdc)
	RegisterEvidences(cdc)
}

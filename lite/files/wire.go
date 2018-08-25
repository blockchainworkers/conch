package files

import (
	"github.com/tendermint/go-amino"
	"github.com/blockchainworkers/conch/types"
)

var cdc = amino.NewCodec()

func init() {
	types.RegisterBlockAmino(cdc)
}

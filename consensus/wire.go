package consensus

import (
	"github.com/tendermint/go-amino"
	"github.com/blockchainworkers/conch/types"
)

var cdc = amino.NewCodec()

func init() {
	RegisterConsensusMessages(cdc)
	RegisterWALMessages(cdc)
	types.RegisterBlockAmino(cdc)
}

package monitor

import (
	amino "github.com/tendermint/go-amino"
	ctypes "github.com/blockchainworkers/conch/rpc/core/types"
)

var cdc = amino.NewCodec()

func init() {
	ctypes.RegisterAmino(cdc)
}

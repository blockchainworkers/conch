package conn

import (
	"github.com/tendermint/go-amino"
	cryptoAmino "github.com/blockchainworkers/conch/crypto/encoding/amino"
)

var cdc *amino.Codec = amino.NewCodec()

func init() {
	cryptoAmino.RegisterAmino(cdc)
	RegisterPacket(cdc)
}

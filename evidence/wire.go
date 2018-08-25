package evidence

import (
	"github.com/tendermint/go-amino"
	cryptoAmino "github.com/blockchainworkers/conch/crypto/encoding/amino"
	"github.com/blockchainworkers/conch/types"
)

var cdc = amino.NewCodec()

func init() {
	RegisterEvidenceMessages(cdc)
	cryptoAmino.RegisterAmino(cdc)
	types.RegisterEvidences(cdc)
}

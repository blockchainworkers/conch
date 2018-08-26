package conchapp

import (
	"github.com/blockchainworkers/conch/abci/types"
	dbm "github.com/blockchainworkers/conch/libs/db"
	"github.com/blockchainworkers/conch/libs/log"
)

//-----------------------------------------

var _ types.Application = (*ConchApplication)(nil)

// ConchApplication is the Application's impl
type ConchApplication struct {
	// validator set
	ValUpdates []types.Validator
	state      dbm.DB
	logger     log.Logger
}

// NewConchApplication return new instance
func NewConchApplication(dbDir string) *ConchApplication {
	name := "conchapplication"
	db, err := dbm.NewGoLevelDB(name, dbDir)
	if err != nil {
		panic(err)
	}

	return &ConchApplication{
		logger: log.NewNopLogger(),
		state:  db,
	}
}

// SetLogger set logger
func (app *ConchApplication) SetLogger(l log.Logger) {
	app.logger = l
}

// Info impl interface
func (app *ConchApplication) Info(req types.RequestInfo) types.ResponseInfo {
	// todo:: need impl
	var res types.ResponseInfo
	return res
}

// SetOption set option
func (app *ConchApplication) SetOption(req types.RequestSetOption) types.ResponseSetOption {
	// todo::
	var res types.ResponseSetOption
	return res
}

//DeliverTx deleiver an transaction to app
func (app *ConchApplication) DeliverTx(tx []byte) types.ResponseDeliverTx {
	// todo::
	return types.ResponseDeliverTx{Code: types.CodeTypeOK}
}

// CheckTx is called when one tx need be send in mempool
func (app *ConchApplication) CheckTx(tx []byte) types.ResponseCheckTx {
	// todo::
	return types.ResponseCheckTx{Code: types.CodeTypeOK}
}

// Commit will panic if InitChain was not called
func (app *ConchApplication) Commit() types.ResponseCommit {
	// todo::
	return types.ResponseCommit{Data: []byte("88888")}
}

// Query for query info
func (app *ConchApplication) Query(reqQuery types.RequestQuery) types.ResponseQuery {
	return types.ResponseQuery{Code: types.CodeTypeOK}
}

//InitChain Save the validators in the merkle tree
func (app *ConchApplication) InitChain(req types.RequestInitChain) types.ResponseInitChain {
	// for _, v := range req.Validators {
	// 	r := app.updateValidator(v)
	// 	if r.IsErr() {
	// 		app.logger.Error("Error updating validators", "r", r)
	// 	}
	// }
	return types.ResponseInitChain{}
}

//BeginBlock Track the block hash and header information
func (app *ConchApplication) BeginBlock(req types.RequestBeginBlock) types.ResponseBeginBlock {
	//todo:: reset valset changes

	return types.ResponseBeginBlock{}
}

//EndBlock Update the validator set
func (app *ConchApplication) EndBlock(req types.RequestEndBlock) types.ResponseEndBlock {
	// todo::
	return types.ResponseEndBlock{ValidatorUpdates: app.ValUpdates}
}

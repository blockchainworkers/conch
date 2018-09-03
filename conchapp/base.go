package conchapp

import (
	"encoding/hex"
	"github.com/blockchainworkers/conch/abci/types"
	"math/big"
	"os"
	//dbm "github.com/blockchainworkers/conch/libs/db"
	"github.com/blockchainworkers/conch/libs/log"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3" // import sqlite
	"path"
)

//-----------------------------------------

var _ types.Application = (*ConchApplication)(nil)

// ConchApplication is the Application's impl
type ConchApplication struct {
	// validator set
	ValUpdates []types.Validator
	logger     log.Logger
	state      *APPState
}

// NewConchApplication return new instance
func NewConchApplication(dbDir string) *ConchApplication {
	name := "conchapplication"

	db, err := sqlx.Open("sqlite3", path.Join(dbDir, name+".db"))
	if err != nil {
		panic(err)
	}
	if err := initDatabase(db); err != nil {
		panic(err)
	}

	//init logger
	logInst := log.NewTMLogger(log.NewSyncWriter(os.Stdout))
	logInst.With("module", "app")
	// init appsate
	appSt := NewAPPState(db, logInst)
	// init queryhandler
	querys.appSt = appSt
	querys.Init()
	return &ConchApplication{
		logger: logInst,
		state:  appSt,
	}
}

// SetLogger set logger
func (app *ConchApplication) SetLogger(l log.Logger) {
	app.logger = l
}

// Info impl interface
func (app *ConchApplication) Info(req types.RequestInfo) types.ResponseInfo {
	// load state from
	err := app.state.HeadSt.LoadHeaderState()
	if err != nil {
		app.logger.Error("load state from db failed", "err", err.Error())
		panic(err)
	}

	var res types.ResponseInfo
	res.LastBlockAppHash, _ = hex.DecodeString(app.state.HeadSt.CurAPPHash)
	res.LastBlockHeight = app.state.HeadSt.CurBlockNum
	res.Version = "v0.01"
	res.Data = "conch is an virtual currency"
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
	// put tx in cahce then exec them when commit
	// txStr := string(tx)
	// if len(txStr) < 3 {
	// 	return types.ResponseDeliverTx{Code: 1, Info: "invalid tx", Log: "invalid tx"}
	// }
	// txS := txStr[1 : len(txStr)-2]

	//var trans Transaction
	trans, err := DecodeNewTx([]byte(tx))
	if err != nil {
		app.logger.Error("when deliver tx, tx can not be UnMarshal", "tx", string(tx), "err", err.Error())
		return types.ResponseDeliverTx{Code: 1, Info: err.Error(), Log: "deliver tx err in json unmarshal"}
	}

	// try load account amount
	account, err := app.state.AccoutSt.LoadAccount(trans.Sender)
	if err != nil {
		app.logger.Error("when deliver tx, load accout info failed", "err", err.Error())
		return types.ResponseDeliverTx{Code: 1, Info: err.Error(), Log: "deliver tx err in loading account"}
	}

	val, ret := new(big.Int).SetString(trans.Value, 0)
	if !ret {
		app.logger.Error("when deliver tx, tx value can't be string of number")
		return types.ResponseDeliverTx{Code: 1, Info: "tx args err", Log: "deliver tx err in tx value"}
	}

	if account.Amount.Cmp(val.Add(val, trans.FeeCalc())) < 0 {
		app.logger.Error("when deliver tx, account Insufficient balance")
		return types.ResponseDeliverTx{Code: 1, Info: "Insufficient balance", Log: "deliver tx err in account balance"}
	}

	if !trans.IsValidTx() {
		app.logger.Error("when deliver tx, tx is not valid")
		return types.ResponseDeliverTx{Code: 1, Info: "sign not valid", Log: "deliver tx err in tx sign"}
	}

	app.state.TxSt.UpdateTx(trans)
	return types.ResponseDeliverTx{Code: types.CodeTypeOK}
}

// CheckTx is called when one tx need be send in mempool
func (app *ConchApplication) CheckTx(tx []byte) types.ResponseCheckTx {
	// 1. tx should be UnMarshal truct
	// 2. check account amount is enough
	// 3. check sign is right

	// app.logger.Error("check tx", "tx", string(tx))
	// txStr := string(tx)
	// if len(txStr) < 3 {
	// 	return types.ResponseCheckTx{Code: 1, Info: "invalid tx", Log: "invalid tx"}
	// }

	// txS := txStr[1 : len(txStr)-2]
	// println(txStr, txS)
	//var trans Transaction
	trans, err := DecodeNewTx(tx)
	if err != nil {
		app.logger.Error("when checking tx, tx can not be UnMarshal", "tx", string(tx), "err", err.Error())
		return types.ResponseCheckTx{Code: 1, Info: err.Error(), Log: "check tx err in json unmarshal"}
	}

	// try load account amount
	account, err := app.state.AccoutSt.LoadAccount(trans.Sender)
	if err != nil {
		app.logger.Error("when checking tx, load accout info err: ", err.Error())
		return types.ResponseCheckTx{Code: 1, Info: err.Error(), Log: "check tx err in loading account"}
	}

	val, ret := new(big.Int).SetString(trans.Value, 0)
	if !ret {
		app.logger.Error("when checking tx, tx value can't be string of number")
		return types.ResponseCheckTx{Code: 1, Info: "tx args err", Log: "check tx err in tx value"}
	}

	if account.Amount.Cmp(val.Add(val, trans.FeeCalc())) < 0 {
		app.logger.Error("when checking tx, account Insufficient balance")
		return types.ResponseCheckTx{Code: 1, Info: "Insufficient balance", Log: "check tx err in account balance"}
	}

	if !trans.IsValidTx() {
		app.logger.Error("when checking tx, tx is not valid")
		return types.ResponseCheckTx{Code: 1, Info: "sign not valid", Log: "check tx err in tx sign"}
	}

	return types.ResponseCheckTx{Code: types.CodeTypeOK}
}

// Commit will panic if InitChain was not called
func (app *ConchApplication) Commit() types.ResponseCommit {
	// todo:: Commit is very import  exect tx update state
	// when commit be exec it means all transaction in this block should have been delivered
	// we should exec transaction then give an receipt for every tx
	// when all transactions have been exec comit tx's state txreceipt's state
	// account's state to db

	appHash, err := app.state.Commit()
	if err != nil {
		app.logger.Error("commit err: ", err.Error())
	}
	appHashByte, _ := hex.DecodeString(appHash)
	return types.ResponseCommit{Data: appHashByte}
}

// Query for query info
func (app *ConchApplication) Query(reqQuery types.RequestQuery) types.ResponseQuery {

	return querys.Query(reqQuery)
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
	app.state.HeadSt.CurBlockHash = hex.EncodeToString(req.Hash)

	// println("------------valitor--------- ", req.Header.Proposer.Power, hex.EncodeToString(req.Header.Proposer.PubKey.Data), hex.EncodeToString(req.Header.ValidatorsHash))
	//app.ValUpdates
	// for iter := range req.ByzantineValidators {

	// 	println("--------ByzantineValidators----- ", "power: ", req.ByzantineValidators[iter].Validator.Power, " pubkey: ", req.ByzantineValidators[iter].Validator.PubKey.String())
	// }
	return types.ResponseBeginBlock{}
}

//EndBlock Update the validator set
func (app *ConchApplication) EndBlock(req types.RequestEndBlock) types.ResponseEndBlock {
	app.state.HeadSt.CurBlockNum = req.Height

	return types.ResponseEndBlock{}
}

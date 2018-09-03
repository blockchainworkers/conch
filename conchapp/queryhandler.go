package conchapp

import (
	"encoding/json"
	"fmt"
	"github.com/blockchainworkers/conch/abci/types"
	"net/url"
	"strconv"
	"sync"
)

// var queryHandlers = make(map[string]func(reqQuery types.RequestQuery) types.ResponseQuery)

type handerFunc func(appSt *APPState, reqQuery types.RequestQuery) types.ResponseQuery

var querys = &QueryProcesser{
	queryHandlers: make(map[string]handerFunc),
}

// QueryProcesser 查询处理
type QueryProcesser struct {
	sync.Mutex
	appSt         *APPState
	queryHandlers map[string]handerFunc
}

// Query main entry
func (q *QueryProcesser) Query(reqQuery types.RequestQuery) types.ResponseQuery {
	q.Lock()
	handler, ok := q.queryHandlers[reqQuery.Path]
	q.Unlock()

	if !ok {
		return types.ResponseQuery{
			Code:  1,
			Info:  "unsported method",
			Value: []byte(fmt.Sprintf("current app does not support %s method", reqQuery.Path)),
		}
	}
	return handler(q.appSt, reqQuery)
}

//RegistHandler ...
func (q *QueryProcesser) RegistHandler(name string, handler handerFunc) {
	q.Lock()
	q.queryHandlers[name] = handler
	q.Unlock()
}

// Init ...
func (q *QueryProcesser) Init() {
	q.RegistHandler("account_balance", accountBalance)
	q.RegistHandler("account_txs", accountTxs)
}

// -------- related query ------

func accountBalance(appSt *APPState, reqQuery types.RequestQuery) types.ResponseQuery {
	// data="address"
	address := string(reqQuery.Data)
	info, _ := appSt.AccoutSt.LoadAccount(address)
	return types.ResponseQuery{
		Code:  0,
		Value: []byte(info.Amount.String()),
		Info:  "success",
		Log:   info.Amount.String(),
	}
}

func accountTxs(appSt *APPState, reqQuery types.RequestQuery) types.ResponseQuery {
	// data="account=?&page=?"
	values, err := url.ParseQuery(string(reqQuery.Data))
	if err != nil {
		return types.ResponseQuery{
			Code:  1,
			Value: []byte(err.Error()),
			Info:  "args requests err",
		}
	}
	// check args
	account := values.Get("account")
	page := values.Get("page")
	if account == "" {
		return types.ResponseQuery{
			Code:  1,
			Value: []byte(err.Error()),
			Info:  "account can't be empty",
		}
	}
	p, err := strconv.Atoi(page)
	if err != nil || p <= 0 {
		p = 1
	}
	txs := appSt.TxSt.QueryTxsByAccount(account, int64((p-1)*10), int64(10))

	dat, _ := json.Marshal(txs)
	return types.ResponseQuery{
		Code:  0,
		Value: dat,
		Info:  "success",
	}
}

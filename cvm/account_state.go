package cvm

import (
	"fmt"
	//"bufio"
	"bytes"
	"encoding/json"
	"github.com/blockchainworkers/conch/crypto"
	"github.com/blockchainworkers/conch/cvm/common"
	"github.com/blockchainworkers/conch/cvm/types"
	// "github.com/blockchainworkers/conch/rlp"
	"math/big"
	"os"
)

var emptyCodeHash = crypto.Sha256(nil)

type accountObject struct {
	Address      common.Address              `json:"address,omitempty"`
	AddrHash     common.Hash                 `json:"addr_hash,omitempty"` // hash of ethereum address of the account
	ByteCode     []byte                      `json:"byte_code,omitempty"`
	Data         accountData                 `json:"data,omitempty"`
	CacheStorage map[common.Hash]common.Hash `json:"cache_storage,omitempty"` // 用于缓存存储的变量
}

type accountData struct {
	Nonce    uint64      `json:"nonce,omitempty"`
	Balance  *big.Int    `json:"balance,omitempty"`
	Root     common.Hash `json:"root,omitempty"` // merkle root of the storage trie
	CodeHash []byte      `json:"code_hash,omitempty"`
}

// newObject creates a state object.
func newAccountObject(address common.Address, data accountData) *accountObject {
	if data.Balance == nil {
		data.Balance = new(big.Int)
	}
	if data.CodeHash == nil {
		data.CodeHash = emptyCodeHash
	}
	return &accountObject{
		Address:      address,
		AddrHash:     common.BytesToHash(crypto.Sha256(address[:])),
		Data:         data,
		CacheStorage: make(map[common.Hash]common.Hash),
	}
}

//balance--
func (object *accountObject) Balance() *big.Int {
	return object.Data.Balance
}

func (object *accountObject) SubBalance(amount *big.Int) {
	if amount.Sign() == 0 {
		return
	}
	object.Data.Balance = new(big.Int).Sub(object.Balance(), amount)
}

func (object *accountObject) AddBalance(amount *big.Int) {
	if amount.Sign() == 0 {
		return
	}
	object.Data.Balance = new(big.Int).Add(object.Balance(), amount)
}

// nonce--
func (object *accountObject) Nonce() uint64 {
	return object.Data.Nonce
}

func (object *accountObject) SetNonce(nonce uint64) {
	object.Data.Nonce = nonce
}

// code-----

func (object *accountObject) CodeHash() []byte {
	return object.Data.CodeHash
}

func (object *accountObject) Code() []byte {
	return object.ByteCode
}

func (object *accountObject) SetCode(codeHash []byte, code []byte) {
	object.Data.CodeHash = codeHash
	object.ByteCode = code
}

// storage sate-------
func (object *accountObject) GetStorageState(key common.Hash) common.Hash {
	value, exist := object.CacheStorage[key]
	if exist {
		// fmt.Println("exist cache ", " key: ", key, " value: ", value)
		return value
	}
	return common.Hash{}
}

func (object *accountObject) SetStorageState(key, value common.Hash) {
	object.CacheStorage[key] = value
}

func (object *accountObject) Empty() bool {
	return object.Data.Nonce == 0 && object.Data.Balance.Sign() == 0 && bytes.Equal(object.Data.CodeHash, emptyCodeHash)
}

//AccountState 实现vm的StateDB的接口 用于进行测试
type AccountState struct {
	Accounts map[common.Address]*accountObject `json:"accounts,omitempty"`
}

// NewAccountStateDb new instance
func NewAccountStateDb() *AccountState {
	return &AccountState{
		Accounts: make(map[common.Address]*accountObject),
	}
}

func (accSt *AccountState) getAccountObject(addr common.Address) *accountObject {
	if value, exist := accSt.Accounts[addr]; exist {
		return value
	}
	return nil
}

func (accSt *AccountState) setAccountObject(obj *accountObject) {
	accSt.Accounts[obj.Address] = obj
}

// 如果不存在则新创建
func (accSt *AccountState) getOrsetAccountObject(addr common.Address) *accountObject {
	get := accSt.getAccountObject(addr)
	if get != nil {
		return get
	}
	set := newAccountObject(addr, accountData{})
	accSt.setAccountObject(set)
	return set
}

// 实现接口-------

//CreateAccount 创建一个新的合约账户
func (accSt *AccountState) CreateAccount(addr common.Address) {
	if accSt.getAccountObject(addr) != nil {
		return
	}
	obj := newAccountObject(addr, accountData{})
	accSt.setAccountObject(obj)
}

// SubBalance 减去某个账户的余额
func (accSt *AccountState) SubBalance(addr common.Address, amount *big.Int) {
	stateObject := accSt.getOrsetAccountObject(addr)
	if stateObject != nil {
		stateObject.SubBalance(amount)
	}
}

// AddBalance 增加某个账户的余额
func (accSt *AccountState) AddBalance(addr common.Address, amount *big.Int) {
	stateObject := accSt.getOrsetAccountObject(addr)
	if stateObject != nil {
		stateObject.AddBalance(amount)
	}
}

// GetBalance 获取某个账户的余额
func (accSt *AccountState) GetBalance(addr common.Address) *big.Int {
	stateObject := accSt.getOrsetAccountObject(addr)
	if stateObject != nil {
		return stateObject.Balance()
	}
	return new(big.Int).SetInt64(0)
}

//GetNonce 获取nonce
func (accSt *AccountState) GetNonce(addr common.Address) uint64 {
	stateObject := accSt.getAccountObject(addr)
	if stateObject != nil {
		return stateObject.Nonce()
	}
	return 0
}

// SetNonce 设置nonce
func (accSt *AccountState) SetNonce(addr common.Address, nonce uint64) {
	stateObject := accSt.getOrsetAccountObject(addr)
	if stateObject != nil {
		stateObject.SetNonce(nonce)
	}
}

// GetCodeHash 获取代码的hash值
func (accSt *AccountState) GetCodeHash(addr common.Address) common.Hash {
	stateObject := accSt.getAccountObject(addr)
	if stateObject == nil {
		return common.Hash{}
	}
	return common.BytesToHash(stateObject.CodeHash())
}

//GetCode 获取智能合约的代码
func (accSt *AccountState) GetCode(addr common.Address) []byte {
	stateObject := accSt.getAccountObject(addr)
	if stateObject != nil {
		return stateObject.Code()
	}
	return nil
}

//SetCode 设置智能合约的code
func (accSt *AccountState) SetCode(addr common.Address, code []byte) {
	stateObject := accSt.getOrsetAccountObject(addr)
	if stateObject != nil {
		stateObject.SetCode(crypto.Sha256(code), code)
	}
}

// GetCodeSize 获取code的大小
func (accSt *AccountState) GetCodeSize(addr common.Address) int {
	stateObject := accSt.getAccountObject(addr)
	if stateObject == nil {
		return 0
	}
	if stateObject.ByteCode != nil {
		return len(stateObject.ByteCode)
	}
	return 0
}

// AddRefund 暂时先忽略补偿
func (accSt *AccountState) AddRefund(uint64) {
	return
}

//GetRefund ...
func (accSt *AccountState) GetRefund() uint64 {
	return 0
}

// GetState 和SetState 是用于保存合约执行时 存储的变量是否发生变化 evm对变量存储的改变消耗的gas是有区别的
func (accSt *AccountState) GetState(addr common.Address, key common.Hash) common.Hash {
	stateObject := accSt.getAccountObject(addr)
	if stateObject != nil {
		return stateObject.GetStorageState(key)
	}
	return common.Hash{}
}

// SetState 设置变量的状态
func (accSt *AccountState) SetState(addr common.Address, key common.Hash, value common.Hash) {
	stateObject := accSt.getOrsetAccountObject(addr)
	if stateObject != nil {
		fmt.Printf("SetState key: %x value: %s", key, new(big.Int).SetBytes(value[:]).String())
		stateObject.SetStorageState(key, value)
	}
}

// Suicide 暂时禁止自杀
func (accSt *AccountState) Suicide(common.Address) bool {
	return false
}

// HasSuicided ...
func (accSt *AccountState) HasSuicided(common.Address) bool {
	return false
}

// Exist 检查账户是否存在
func (accSt *AccountState) Exist(addr common.Address) bool {
	return accSt.getAccountObject(addr) != nil
}

//Empty 是否是空账户
func (accSt *AccountState) Empty(addr common.Address) bool {
	so := accSt.getAccountObject(addr)
	return so == nil || so.Empty()
}

// RevertToSnapshot ...
func (accSt *AccountState) RevertToSnapshot(int) {

}

// Snapshot ...
func (accSt *AccountState) Snapshot() int {
	return 0
}

// AddLog 添加事件触发日志
func (accSt *AccountState) AddLog(log *types.Log) {
	//fmt.Printf("log: %v", log)
}

// AddPreimage 暂时没搞清楚这个是干嘛用的
func (accSt *AccountState) AddPreimage(common.Hash, []byte) {

}

// ForEachStorage  暂时没发现vm调用这个接口
func (accSt *AccountState) ForEachStorage(common.Address, func(common.Hash, common.Hash) bool) {

}

// Commit 进行持久换存储
func (accSt *AccountState) Commit() error {
	// 将bincode写入文件
	file, err := os.Create("./account_sate.db")
	if err != nil {
		return err
	}
	err = json.NewEncoder(file).Encode(accSt)
	//fmt.Println("len(binCode): ", len(binCode), " code: ", binCode)
	// bufW := bufio.NewWriter(file)
	// bufW.Write(binCode)
	// // bufW.WriteByte('\n')
	// bufW.Flush()
	file.Close()
	return err
}

//TryLoadFromDisk  尝试从磁盘加载AccountState
func TryLoadFromDisk() (*AccountState, error) {
	file, err := os.Open("./account_sate.db")
	if err != nil && os.IsNotExist(err) {
		return NewAccountStateDb(), nil
	}
	if err != nil {
		return nil, err
	}

	// stat, _ := file.Stat()
	// // buf := stat.Size()
	var accStat AccountState

	err = json.NewDecoder(file).Decode(&accStat)
	return &accStat, err
}

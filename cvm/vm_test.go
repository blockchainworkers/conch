package cvm

import (
	"encoding/hex"
	"github.com/blockchainworkers/conch/cvm/common"
	"github.com/blockchainworkers/conch/cvm/vm"
	"math/big"
	"testing"
)

// 进行测试

var normalAddress, _ = hex.DecodeString("123456abc")
var contractAddress, _ = hex.DecodeString("987654321")
var normalAccount = common.BytesToAddress(normalAddress)
var contactAccont = common.BytesToAddress(contractAddress)

var byteCodeStr = "6060604052341561000f57600080fd5b640165a0bc006000819055506101878061002a6000396000f300606060405260043610610041576000357c0100000000000000000000000000000000000000000000000000000000900463ffffffff1680634d9b3d5d14610046575b600080fd5b341561005157600080fd5b6100596100a2565b604051808373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020018281526020019250505060405180910390f35b6000807f08c31d20d5c3a5f2cfe0adf83909e6411f43fe97eb091e15c12f3e5a203e8fde33604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200180602001828103825260058152602001807f66756e636b0000000000000000000000000000000000000000000000000000008152506020019250505060405180910390a13360008081548092919060019003919050559150915090915600a165627a7a7230582033de42289d2d250bb25095884bd901f7e794cce85a62df985896e8d1e681eb7b0029"
var byteCode, _ = hex.DecodeString(byteCodeStr)

var input, _ = hex.DecodeString("4d9b3d5d")

func TestRunVM(t *testing.T) {

	// 创建账户State
	stateDb := NewAccountStateDb()
	// 创建一个普通账户
	stateDb.CreateAccount(normalAccount)
	stateDb.CreateAccount(contactAccont)

	stateDb.AddBalance(normalAccount, big.NewInt(0x878999988776612))
	stateDb.SetCode(contactAccont, byteCode)

	evmCtx := NewEVMContext(normalAccount, 100, 1200000, 1)
	vmenv := vm.NewEVM(evmCtx, stateDb, vm.Config{})

	ret, leftgas, err := vmenv.Call(vm.AccountRef(normalAccount), contactAccont, input, 1000000, big.NewInt(0))
	t.Logf("ret: %v, leftGas: %v, err: %v, len(ret): %v, hexret: %v", ret, leftgas, err, len(ret), hex.EncodeToString(ret))
}

package main

import (
	"encoding/hex"
	"fmt"
	"github.com/blockchainworkers/conch/cvm"
	"github.com/blockchainworkers/conch/cvm/common"
)

var normalAddress, _ = hex.DecodeString("123456abc")
var hellWorldcontractAddress, _ = hex.DecodeString("987654321")
var baseContractAddress, _ = hex.DecodeString("038f160ad632409bfb18582241d9fd88c1a072ba")
var normalAccount = common.BytesToAddress(normalAddress)
var helloWorldcontactAccont = common.BytesToAddress(hellWorldcontractAddress)
var baseContractAccont = common.BytesToAddress(baseContractAddress)

// 基本账户字节码
var baseCodeStr = "608060405260043610610062576000357c0100000000000000000000000000000000000000000000000000000000900463ffffffff1680632b225f29146100675780638afc3605146100f75780638da5cb5b1461010e578063f2fde38b14610165575b600080fd5b34801561007357600080fd5b5061007c6101a8565b6040518080602001828103825283818151815260200191508051906020019080838360005b838110156100bc5780820151818401526020810190506100a1565b50505050905090810190601f1680156100e95780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b34801561010357600080fd5b5061010c6101e5565b005b34801561011a57600080fd5b50610123610227565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b34801561017157600080fd5b506101a6600480360381019080803573ffffffffffffffffffffffffffffffffffffffff16906020019092919050505061024c565b005b60606040805190810160405280601081526020017f42617365436f6e747261637456302e3100000000000000000000000000000000815250905090565b336000806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff160217905550565b6000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b6000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff161415156102a757600080fd5b600073ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff16141515156102e357600080fd5b806000806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055508073ffffffffffffffffffffffffffffffffffffffff166000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e060405160405180910390a3505600a165627a7a723058208c3064096245894122f6bcf5e2ee12e30d4775a3b8dca0b21f10d5a5bc386e8b0029"

// hellworld 账户字节码
var hellCodeStr = "6080604052600436106100615763ffffffff7c01000000000000000000000000000000000000000000000000000000006000350416634d9b3d5d81146100665780637e8800a7146100ab578063c3f82bc3146100c2578063fb1669ca14610165575b600080fd5b34801561007257600080fd5b5061007b61017d565b6040805173ffffffffffffffffffffffffffffffffffffffff909316835260208301919091528051918290030190f35b3480156100b757600080fd5b506100c06101fa565b005b3480156100ce57600080fd5b506100f073ffffffffffffffffffffffffffffffffffffffff6004351661028f565b6040805160208082528351818301528351919283929083019185019080838360005b8381101561012a578181015183820152602001610112565b50505050905090810190601f1680156101575780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b34801561017157600080fd5b506100c0600435610389565b60408051338152602081018290526005818301527f66756e636b0000000000000000000000000000000000000000000000000000006060820152905160009182917f08c31d20d5c3a5f2cfe0adf83909e6411f43fe97eb091e15c12f3e5a203e8fde9181900360800190a150506000805460001981019091553391565b600080526001602090815260647fa6eef7e35abe7026729641147f7915573c7e97b47efa546f5f6e3230263bcb4955604080513381529182018190526008828201527f6f6e6c79746573740000000000000000000000000000000000000000000000006060830152517f08c31d20d5c3a5f2cfe0adf83909e6411f43fe97eb091e15c12f3e5a203e8fde9181900360800190a1565b606060008290508073ffffffffffffffffffffffffffffffffffffffff16632b225f296040518163ffffffff167c0100000000000000000000000000000000000000000000000000000000028152600401600060405180830381600087803b1580156102fa57600080fd5b505af115801561030e573d6000803e3d6000fd5b505050506040513d6000823e601f3d908101601f19168201604052602081101561033757600080fd5b81019080805164010000000081111561034f57600080fd5b8201602081018481111561036257600080fd5b815164010000000081118282018710171561037c57600080fd5b5090979650505050505050565b6000555600a165627a7a72305820c63a859d93a3512b52ccaec75bb9aa146648c41b21c8a0cd0cd2e2c1aede35ed0029"

var helloCode, _ = hex.DecodeString(hellCodeStr)
var baseCode, _ = hex.DecodeString(baseCodeStr)

func updateContract() {
	// 加载账户State
	stateDb, err := cvm.TryLoadFromDisk()
	if err != nil {
		panic(err)
	}
	stateDb.SetCode(helloWorldcontactAccont, helloCode)
	stateDb.SetCode(baseContractAccont, baseCode)
	fmt.Println(stateDb.Commit())
}

var baseContractABIJson = `[
	{
		"constant": true,
		"inputs": [],
		"name": "CurrentVersion",
		"outputs": [
			{
				"name": "",
				"type": "string"
			}
		],
		"payable": false,
		"stateMutability": "pure",
		"type": "function"
	},
	{
		"constant": false,
		"inputs": [],
		"name": "Ownable",
		"outputs": [],
		"payable": false,
		"stateMutability": "nonpayable",
		"type": "function"
	},
	{
		"constant": true,
		"inputs": [],
		"name": "owner",
		"outputs": [
			{
				"name": "",
				"type": "address"
			}
		],
		"payable": false,
		"stateMutability": "view",
		"type": "function"
	},
	{
		"constant": false,
		"inputs": [
			{
				"name": "newOwner",
				"type": "address"
			}
		],
		"name": "transferOwnership",
		"outputs": [],
		"payable": false,
		"stateMutability": "nonpayable",
		"type": "function"
	},
	{
		"anonymous": false,
		"inputs": [
			{
				"indexed": true,
				"name": "previousOwner",
				"type": "address"
			},
			{
				"indexed": true,
				"name": "newOwner",
				"type": "address"
			}
		],
		"name": "OwnershipTransferred",
		"type": "event"
	}
]`

var hellWorldContractABIJson = `[
	{
		"constant": false,
		"inputs": [],
		"name": "getbalance",
		"outputs": [
			{
				"name": "a",
				"type": "address"
			},
			{
				"name": "_b",
				"type": "uint256"
			}
		],
		"payable": false,
		"stateMutability": "nonpayable",
		"type": "function"
	},
	{
		"constant": false,
		"inputs": [],
		"name": "onlytest",
		"outputs": [],
		"payable": false,
		"stateMutability": "nonpayable",
		"type": "function"
	},
	{
		"constant": true,
		"inputs": [
			{
				"name": "contractAddr",
				"type": "address"
			}
		],
		"name": "getVersion",
		"outputs": [
			{
				"name": "",
				"type": "string"
			}
		],
		"payable": false,
		"stateMutability": "view",
		"type": "function"
	},
	{
		"constant": false,
		"inputs": [
			{
				"name": "tmp",
				"type": "uint256"
			}
		],
		"name": "setBalance",
		"outputs": [],
		"payable": false,
		"stateMutability": "nonpayable",
		"type": "function"
	},
	{
		"inputs": [],
		"payable": false,
		"stateMutability": "nonpayable",
		"type": "constructor"
	},
	{
		"anonymous": false,
		"inputs": [
			{
				"indexed": false,
				"name": "",
				"type": "address"
			},
			{
				"indexed": false,
				"name": "",
				"type": "string"
			}
		],
		"name": "Triggle",
		"type": "event"
	}
]`

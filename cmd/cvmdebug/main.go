package main

import (
	"encoding/hex"
	"fmt"
	"github.com/blockchainworkers/conch/cvm"
	"github.com/blockchainworkers/conch/cvm/abi"
	"github.com/blockchainworkers/conch/cvm/common"
	"github.com/blockchainworkers/conch/cvm/vm"
	"math/big"
	"reflect"
	"strings"
)

// 4d9b3d5d : getbalance  7e8800a7: onlytest fb1669ca000000000000000000000000000000000000000000000000000000000000029a: setbalance 666
var input, _ = hex.DecodeString("7e8800a7")

func main() {
	// updateContract()
	// return
	// 创建账户State
	stateDb, err := cvm.TryLoadFromDisk()
	if err != nil {
		panic(err)
	}

	evmCtx := cvm.NewEVMContext(normalAccount, 100, 1200000, 1)
	vmenv := vm.NewEVM(evmCtx, stateDb, vm.Config{})

	ret, leftgas, err := vmenv.Call(vm.AccountRef(normalAccount), helloWorldcontactAccont, input, 1000000, big.NewInt(0))
	fmt.Printf("ret: %v, usedGas: %v, err: %v, len(ret): %v, hexret: %v, ", ret, 1000000-leftgas, err, len(ret), hex.EncodeToString(ret))

	abiObjet, _ := abi.JSON(strings.NewReader(hellWorldContractABIJson))

	// begin, length, _ := lengthPrefixPointsTo(0, ret)
	addr := new(common.Address)

	value := big.NewInt(0) //new(*big.Int)
	restult := []interface{}{addr, &value}
	fmt.Println(abiObjet.Unpack(&restult, "getbalance", ret))
	//fmt.Println(unpackAtomic(&restult, string(ret[begin:begin+length])))
	println(restult[0].(*common.Address).String(), (*restult[1].(**big.Int)).String())
	fmt.Println(stateDb.Commit())

}

func lengthPrefixPointsTo(index int, output []byte) (start int, length int, err error) {
	bigOffsetEnd := big.NewInt(0).SetBytes(output[index : index+32])
	bigOffsetEnd.Add(bigOffsetEnd, big.NewInt(32))
	outputLength := big.NewInt(int64(len(output)))

	if bigOffsetEnd.Cmp(outputLength) > 0 {
		return 0, 0, fmt.Errorf("abi: cannot marshal in to go slice: offset %v would go over slice boundary (len=%v)", bigOffsetEnd, outputLength)
	}

	if bigOffsetEnd.BitLen() > 63 {
		return 0, 0, fmt.Errorf("abi offset larger than int64: %v", bigOffsetEnd)
	}

	offsetEnd := int(bigOffsetEnd.Uint64())
	lengthBig := big.NewInt(0).SetBytes(output[offsetEnd-32 : offsetEnd])

	totalSize := big.NewInt(0)
	totalSize.Add(totalSize, bigOffsetEnd)
	totalSize.Add(totalSize, lengthBig)
	if totalSize.BitLen() > 63 {
		return 0, 0, fmt.Errorf("abi length larger than int64: %v", totalSize)
	}

	if totalSize.Cmp(outputLength) > 0 {
		return 0, 0, fmt.Errorf("abi: cannot marshal in to go type: length insufficient %v require %v", outputLength, totalSize)
	}
	start = int(bigOffsetEnd.Uint64())
	length = int(lengthBig.Uint64())
	return
}

func unpackAtomic(v interface{}, marshalledValues interface{}) error {

	elem := reflect.ValueOf(v).Elem()
	// kind := elem.Kind()
	reflectValue := reflect.ValueOf(marshalledValues)
	return set(elem, reflectValue)
}

func set(dst, src reflect.Value) error {
	dstType := dst.Type()
	srcType := src.Type()
	switch {
	case dstType.AssignableTo(srcType):
		dst.Set(src)
	case dstType.Kind() == reflect.Interface:
		dst.Set(src)
	case dstType.Kind() == reflect.Ptr:
		return set(dst.Elem(), src)
	default:
		return fmt.Errorf("abi: cannot unmarshal %v in to %v", src.Type(), dst.Type())
	}
	return nil
}

package main

import (
	"fmt"
	"github.com/btcsuite/btcutil"
)

func main() {
	priStr := "cV4SVnFobWUcRwa8ZLRfSQuMeN1P6k1ycjvr8MWvMhxF8gGkNyN3"
	// _, privByte, err := Base58CheckDecode(priStr)
	// if err != nil {
	// 	panic(fmt.Errorf("解码私钥失败: %s", err.Error()))
	// }
	// // 从byte到私钥
	// priv, _ := secp256k1.PrivKeyFromBytes(secp256k1.S256(), privByte)
	// fmt.Printf("私钥为: %s, 公钥为: %0x\n", priStr, priv.PubKey().SerializeCompressed())

	wif, err := btcutil.DecodeWIF(priStr)
	if err != nil {
		panic(fmt.Errorf("解码私钥失败: %s", err.Error()))
	}
	fmt.Printf("私钥为: %s, 公钥为: %0x\n", priStr, wif.SerializePubKey())
}

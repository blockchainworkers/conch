package account

import (
	"encoding/base64"
	"encoding/hex"
	"github.com/blockchainworkers/conch/crypto"
	"github.com/btcsuite/btcutil/base58"
	"testing"
)

func TestGenerateAccount(t *testing.T) {
	prik, address := "4d9502a6ec8d978ca27004880715e4417ed4f61130301f1030b48747d1a6df1c", "CONCHcd1sGBDdmQasWZEVfe6x6y3iPij1g67LZJ" //GenerateAccout()

	t.Logf("prikey: %s, addr: %s", prik, address)
	prv, err := LoadPrivKey(prik)

	if err != nil {
		t.Errorf("%s", err.Error())
	}
	t.Log(base64.StdEncoding.EncodeToString(prv.PubKey().ByteArray()))

	addr := "CONCH" + base58.CheckEncode(prv.PubKey().Address(), AddrVersion)
	t.Logf("addr: %s", addr)
	_, _, err = base58.CheckDecode("cd1sGBDdmQasWZEVfe6x6y3iPij1g67LZJ")
	t.Log(err)

	if address != addr {
		t.Errorf("load prikey address err")
	}

	t.Logf("nil hash: %s", hex.EncodeToString(crypto.Sha256(nil)))
}

//667401408d784bdf6c565895e5cac04afe402c09952bd652a3477813bc8c2354, CONCHcd1sGBDdmQasWZEVfe6x6y3iPij1g67LZJ
// 56e0cd27cb67017942776d8359579c76eb0b01168f53237e132b186d7f64754a, addr: CONCHcj7RJN1thrdPxEXd5h2iALJp8Rz6RM2U3k
// 4d9502a6ec8d978ca27004880715e4417ed4f61130301f1030b48747d1a6df1c

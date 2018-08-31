package account

import (
	"encoding/base64"
	"encoding/hex"
	"github.com/blockchainworkers/conch/crypto"
	"github.com/btcsuite/btcutil/base58"
	"testing"
)

func TestGenerateAccount(t *testing.T) {
	prik, address := "449c27e80a3dfdbe1854b83df3ca6484bab0fcba68c82ad3c141341aa46cba0e", "CONCHcd1sGBDdmQasWZEVfe6x6y3iPij1g67LZJ" //GenerateAccout()
	//prik, address := GenerateAccout()
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

// 667401408d784bdf6c565895e5cac04afe402c09952bd652a3477813bc8c2354(ZnQBQI14S99sVliV5crASv5ALAmVK9ZSo0d4E7yMI1Q=),  CONCHcd1sGBDdmQasWZEVfe6x6y3iPij1g67LZJ  pub: A6jY8EODrVWM86iRJOhWx/KUg3m3EogqO7D4lI+ax9mQ
// 56e0cd27cb67017942776d8359579c76eb0b01168f53237e132b186d7f64754a(VuDNJ8tnAXlCd22DWVecdusLARaPUyN+EysYbX9kdUo=), CONCHcj7RJN1thrdPxEXd5h2iALJp8Rz6RM2U3k A7sZJOMYsJ0DL/lN5CX2Zm1/jWmjAF6Wmex1I+xhoWC4
// 4d9502a6ec8d978ca27004880715e4417ed4f61130301f1030b48747d1a6df1c(TZUCpuyNl4yicASIBxXkQX7U9hEwMB8QMLSHR9Gm3xw=),  A7C1pYP/mrQ6Jnp3oQMpAVKpUOnAQjKpLA95e7MbV/eR
// 2377e8542ce25cb01043d038fb864137931f94c4c4bb699b806149c575609e25(I3foVCziXLAQQ9A4+4ZBN5MflMTEu2mbgGFJxXVgniU=), AhIznLQAiqHR7IHeGfK+pUlevJXlOIDa5hzcZayEwJzb
// 449c27e80a3dfdbe1854b83df3ca6484bab0fcba68c82ad3c141341aa46cba0e(RJwn6Ao9/b4YVLg988pkhLqw/LpoyCrTwUE0GqRsug4=), A10lfiB9p7xqte+3AnaJoEnsWKRJO3/JLClq0KweLeku

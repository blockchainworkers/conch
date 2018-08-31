package conchapp

import (
	"encoding/json"
	"github.com/blockchainworkers/conch/account"
	"testing"
	"time"
)

func TestSendTx(t *testing.T) {
	prv, _ := account.LoadPrivKey("667401408d784bdf6c565895e5cac04afe402c09952bd652a3477813bc8c2354")
	tx := BuildNewTx("CONCHcd1sGBDdmQasWZEVfe6x6y3iPij1g67LZJ", "CONCH3oohB7zfYkApUz5SJ3Xy9D5MEDBM", "hello blockchain", "test", "10000", time.Now().Unix(), 100, 100000)
	tx.SignTx(prv)
	dat, _ := json.Marshal(tx)
	t.Log(string(dat), tx.Serialization())
}

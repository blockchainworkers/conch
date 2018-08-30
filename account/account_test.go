package account

import (
	"github.com/btcsuite/btcutil/base58"
	"testing"
)

func TestGenerateAccount(t *testing.T) {
	prik, address := GenerateAccout()

	t.Logf("prikey: %s, addr: %s", prik, address)
	prv, err := LoadPrivKey(prik)

	if err != nil {
		t.Errorf("%s", err.Error())
	}

	addr := "CONCH" + base58.Encode(prv.PubKey().Address())
	t.Logf("addr: %s", addr)

	if address != addr {
		t.Errorf("load prikey address err")
	}

}

//667401408d784bdf6c565895e5cac04afe402c09952bd652a3477813bc8c2354, CONCH2smxiYHvzFTsK3di3Rd6rJCQgrCN

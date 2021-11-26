package bridge

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"
)

var b *Bridge

func init() {
	b = &Bridge{
		url:       "http://cex-bridge-api.test-15.huobiapps.com/cex-bridge-api/v1",
		apiKey:    "kB01A8gsMv2CP4Tctd7XnUaG3iDE9rQJ",
		secretKey: "ho8XcsfygSKTLEvPYOZ9W12k0i4IFJbC",
		cli:       &http.Client{},
	}
}
func TestMd5Sign(t *testing.T) {
	ret := md5SignHex("&method=getAccountList&secret_key=zgUkV1wTZJ3lasYup60KHePS5MQxF4dq")
	t.Logf("%s", ret)
}

func TestGetChainList(t *testing.T) {
	chainList, err := b.GetChainList()
	if err != nil {
		t.Fatalf("getChainList err:%v", err)
	}
	t.Logf("chainList :%v", chainList)
}

func TestGetCurrencyList(t *testing.T) {
	clist, err := b.GetCurrencyList()
	if err != nil {
		t.Fatalf("get c list err:%v", err)
	}
	b, _ := json.Marshal(clist)
	t.Logf("clist:%s", b)
}

func TestAddAccount(t *testing.T) {
	a := &AccountAdd{
		AccounType: 2,
		ChainId:    128,
		Account:    strings.ToLower("0x70d1B6F6536Bb31ecC375eA600089C3631A51Ff5"),
	}
	if b == nil {
		panic("---")
	}
	id, err := b.AddAccount(a)
	if err != nil {
		t.Fatalf("add account err:%v", err)
	}
	t.Logf("accountid:%d", id)
}

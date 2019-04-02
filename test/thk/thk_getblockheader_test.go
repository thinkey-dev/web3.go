package test

import (
	"testing"
	"web3.go/web3"
	"web3.go/web3/providers"
)

func TestThkGetBlockHeader(t *testing.T) {
	var connection = web3.NewWeb3(providers.NewHTTPProvider("thinkey.natapp1.cc", 10, false))
	res, err := connection.Thk.GetBlockHeader("2", "30")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Log("nonce:", res)
}

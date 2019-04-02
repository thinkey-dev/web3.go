package test

import (
	"testing"
	"web3.go/web3"
	"web3.go/web3/providers"
)

func TestThkGetBalance(t *testing.T) {
	var connection = web3.NewWeb3(providers.NewHTTPProvider("thinkey.natapp1.cc", 10, false))
	connection.DefaultAddress = "0x2c7536e3605d9c16a7a3d7b1898e529396a65c23"
	bal, err := connection.Thk.GetBalance(connection.DefaultAddress, "2")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Log("Balance:", bal)
}

func TestThkGetNonce(t *testing.T) {
	var connection = web3.NewWeb3(providers.NewHTTPProvider("thinkey.natapp1.cc", 10, false))
	connection.DefaultAddress = "0x2c7536e3605d9c16a7a3d7b1898e529396a65c23"
	nonce, err := connection.Thk.GetNonce(connection.DefaultAddress, "2")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Log("nonce:", nonce)
}

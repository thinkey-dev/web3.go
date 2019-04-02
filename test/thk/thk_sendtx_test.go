package test

import (
	"github.com/ethereum/go-ethereum/crypto"
	"strconv"
	"testing"
	"web3.go/web3"
	"web3.go/web3/providers"
	"web3.go/web3/thk/util"
)

var (
	key = "4c0883a69102937d6231471b5dbb6204fe5129617082792ae468d01a3f362318"
)

func TestThkSendTx(t *testing.T) {
	var err error
	var connection = web3.NewWeb3(providers.NewHTTPProvider("thinkey.natapp1.cc", 10, false))
	from := "0x2c7536e3605d9c16a7a3d7b1898e529396a65c23"
	to := "0x6ea0fefc17c877c7a4b0f139728ed39dc134a967"
	nonce, err := connection.Thk.GetNonce(from, "2")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	transaction := util.Transaction{
		ChainId: "2", FromChainId: "2", ToChainId: "2", From: from,
		To: to, Value: "2333", Input: "", Nonce: strconv.Itoa(int(nonce)),
	}
	privatekey, err := crypto.HexToECDSA(key)
	err = connection.Thk.SignTransaction(&transaction, privatekey)

	txhash, err := connection.Thk.SendTx(&transaction)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Log("txhash:", txhash)
}

func TestThkCallTx(t *testing.T) {
	var err error
	var connection = web3.NewWeb3(providers.NewHTTPProvider("thinkey.natapp1.cc", 10, false))
	from := "0x2c7536e3605d9c16a7a3d7b1898e529396a65c23"
	to := "0x6ea0fefc17c877c7a4b0f139728ed39dc134a967"
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	transaction := util.Transaction{
		ChainId: "2", FromChainId: "2", ToChainId: "2", From: from,
		To: to, Value: "2333", Input: "", Nonce: "1",
	}
	res, err := connection.Thk.CallTransaction(&transaction)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Log("result", res)
}

func TestThkGetTransactionByHash(t *testing.T) {
	var err error
	var connection = web3.NewWeb3(providers.NewHTTPProvider("thinkey.natapp1.cc", 10, false))
	hash := "0xcb53f1ec9c02053a46de488b63b219217826fd9c4cfb531567d61003664ef653"
	res, err := connection.Thk.GetTransactionByHash("2", hash)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Log("res:", res)
}

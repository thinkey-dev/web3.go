package test

import (
	"encoding/json"
	"io/ioutil"
	"strconv"
	"testing"
	"time"
	"web3.go/common/cryp/crypto"
	"web3.go/web3"
	"web3.go/web3/providers"
	"web3.go/web3/thk/util"
)

func TestThkContract(t *testing.T) {
	content, err := ioutil.ReadFile("../resources/dahan.json")
	type TruffleContract struct {
		Abi      string `json:"abi"`
		Bytecode string `json:"bytecode"`
	}
	var unmarshalResponse TruffleContract
	err = json.Unmarshal(content, &unmarshalResponse)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	var connection = web3.NewWeb3(providers.NewHTTPProvider("test.thinkey.xyz", 10, false))

	bytecode := unmarshalResponse.Bytecode
	contract, err := connection.Thk.NewContract(unmarshalResponse.Abi)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	from := "0x2c7536e3605d9c16a7a3d7b1898e529396a65c23"
	nonce, err := connection.Thk.GetNonce(from, "2")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	transaction := util.Transaction{
		ChainId: "2", FromChainId: "2", ToChainId: "2", From: from,
		To: "", Value: "0", Input: "", Nonce: strconv.Itoa(int(nonce)),
	}
	privatekey, err := crypto.HexToECDSA(key)
	hash, err := contract.Deploy(transaction, bytecode, privatekey, nil)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Log(hash)
	time.Sleep(time.Second * 5)
	recepit, err := connection.Thk.GetTransactionByHash("2", hash)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Log("contract address:", recepit.ContractAddress)
	transaction.To = recepit.ContractAddress
	result, err := contract.Call(transaction, "getObjsNum")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Log("result:", result)
}

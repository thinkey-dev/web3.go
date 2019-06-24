package test

import (
	"encoding/binary"
	"fmt"
	"io"
	"math/big"
	"strconv"
	"testing"
	"web3.go/common/cryp/crypto"
	"web3.go/common/hexutil"
	"web3.go/encoding"
	"web3.go/web3"
	"web3.go/web3/providers"
	"web3.go/web3/thk/util"
)

var (
	key = "4c0883a69102937d6231471b5dbb6204fe5129617082792ae468d01a3f362318"
)

const (
	AddressLength = 20
)

type (
	ChainID uint32
	Height  uint64

	Address [AddressLength]byte

	Addresser interface {
		Address() Address
	}
)

type CashCheck struct {
	FromChain    ChainID  `json:"FromChain"`    // 转出链
	FromAddress  Address  `json:"FromAddr"`     // 转出账户
	Nonce        uint64   `json:"Nonce"`        // 转出账户提交请求时的nonce
	ToChain      ChainID  `json:"ToChain"`      // 目标链
	ToAddress    Address  `json:"ToAddr"`       // 目标账户
	ExpireHeight Height   `json:"ExpireHeight"` // 过期高度，指的是当目标链高度超过（不含）这个值时，这张支票不能被支取，只能退回
	Amount       *big.Int `json:"Amount"`       // 金额
}

// 4字节FromChain + 20字节FromAddress + 8字节Nonce + 4字节ToChain + 20字节ToAddress +
// 8字节ExpireHeight + 1字节len(Amount.Bytes()) + Amount.Bytes()
// 均为BigEndian
func (c *CashCheck) Serialization(w io.Writer) error {
	buf4 := make([]byte, 4)
	buf8 := make([]byte, 8)

	binary.BigEndian.PutUint32(buf4, uint32(c.FromChain))
	_, err := w.Write(buf4)
	if err != nil {
		return err
	}

	_, err = w.Write(c.FromAddress[:])
	if err != nil {
		return err
	}

	binary.BigEndian.PutUint64(buf8, uint64(c.Nonce))
	_, err = w.Write(buf8)
	if err != nil {
		return err
	}

	binary.BigEndian.PutUint32(buf4, uint32(c.ToChain))
	_, err = w.Write(buf4)
	if err != nil {
		return err
	}

	_, err = w.Write(c.ToAddress[:])
	if err != nil {
		return err
	}

	binary.BigEndian.PutUint64(buf8, uint64(c.ExpireHeight))
	_, err = w.Write(buf8)
	if err != nil {
		return err
	}

	buf4 = buf4[:1]
	var mbytes []byte
	if c.Amount != nil {
		mbytes = c.Amount.Bytes()
	}
	buf4[0] = byte(len(mbytes))
	_, err = w.Write(buf4)
	if err != nil {
		return err
	}
	if buf4[0] > 0 {
		_, err = w.Write(mbytes)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *CashCheck) Deserialization(r io.Reader) error {
	buf4 := make([]byte, 4)
	buf8 := make([]byte, 8)

	_, err := r.Read(buf4)
	if err != nil {
		return err
	}
	c.FromChain = ChainID(binary.BigEndian.Uint32(buf4))

	_, err = r.Read(c.FromAddress[:])
	if err != nil {
		return err
	}

	_, err = r.Read(buf8)
	if err != nil {
		return err
	}
	c.Nonce = binary.BigEndian.Uint64(buf8)

	_, err = r.Read(buf4)
	if err != nil {
		return err
	}
	c.ToChain = ChainID(binary.BigEndian.Uint32(buf4))

	_, err = r.Read(c.ToAddress[:])
	if err != nil {
		return err
	}

	_, err = r.Read(buf8)
	if err != nil {
		return err
	}
	c.ExpireHeight = Height(binary.BigEndian.Uint64(buf8))

	buf4 = buf4[:1]
	_, err = r.Read(buf4)
	if err != nil {
		return err
	}
	length := int(buf4[0])

	if length > 0 {
		mbytes := make([]byte, length)
		_, err = r.Read(mbytes)
		if err != nil {
			return err
		}
		c.Amount = new(big.Int)
		c.Amount.SetBytes(mbytes)
	} else {
		c.Amount = big.NewInt(0)
	}

	return nil
}
func BytesToAddress(b []byte) Address {
	var a Address
	a.SetBytes(b)
	return a
}
func (a *Address) SetBytes(b []byte) {
	if len(b) > len(a) {
		b = b[len(b)-AddressLength:]
	}
	copy(a[AddressLength-len(b):], b)
}

func TestThkCashCheck(t *testing.T) {
	var err error
	var connection = web3.NewWeb3(providers.NewHTTPProvider("192.168.1.106:8093", 10, false))
	from := "0x2c7536e3605d9c16a7a3d7b1898e529396a65c23"
	to := "0x0000000000000000000000000000000000020000"

	nonce, err := connection.Thk.GetNonce(from, "2")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	from_str, err := hexutil.Decode("0x2c7536e3605d9c16a7a3d7b1898e529396a65c23")
	to_str, err := hexutil.Decode("0x4fa1c4e6182b6b7f3bca273390cf587b50b47311")
	vcc := &CashCheck{
		FromChain:    2,
		FromAddress:  BytesToAddress(from_str),
		Nonce:        uint64(nonce),
		ToChain:      3,
		ToAddress:    BytesToAddress(to_str),
		ExpireHeight: 279228 + 5000,
		Amount:       big.NewInt(1),
	}
	println(vcc.Nonce)
	intput, err := encoding.Marshal(vcc)
	println(intput)

	str := hexutil.Encode(intput)
	fmt.Println("------------------")
	fmt.Println(str)
	transaction := util.Transaction{
		ChainId: "2", FromChainId: "2", ToChainId: "3", From: from,
		To: to, Value: "0", Input: str, Nonce: strconv.Itoa(int(nonce)),
	}

	privatekey, err := crypto.HexToECDSA(key)
	err = connection.Thk.SignTransaction(&transaction, privatekey)

	txhash, err := connection.Thk.SendTx(&transaction)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Log("txhash:", txhash)
	//0x472a80cd5a8aa4664fcca5f3a4fd72c3ff25681c2511325f4613f04c128966e9
}
func TestThkSaveCashCheck(t *testing.T) {
	var err error
	var connection = web3.NewWeb3(providers.NewHTTPProvider("192.168.1.106:8093", 10, false))
	from := "0x2c7536e3605d9c16a7a3d7b1898e529396a65c23"
	to := "0x0000000000000000000000000000000000030000"

	nonce, err := connection.Thk.GetNonce(from, "3")
	fmt.Println(nonce)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	//
	//from_str,err:=hexutil.Decode("0x2c7536e3605d9c16a7a3d7b1898e529396a65c23")
	//to_str,err:=hexutil.Decode("0x4fa1c4e6182b6b7f3bca273390cf587b50b47311")
	//vcc := &CashCheck{
	//	FromChain:    2,
	//	FromAddress:  BytesToAddress(from_str),
	//	Nonce:        uint64(nonce),
	//	ToChain:      3,
	//	ToAddress:    BytesToAddress(to_str),
	//	ExpireHeight: 33772,
	//	Amount:       big.NewInt(1),
	//}
	//println(vcc.Nonce)
	//intput,err:=encoding.Marshal(vcc)
	//println(intput)
	//
	//str:=hexutil.Encode(intput)
	transaction := util.Transaction{
		ChainId: "3", FromChainId: "3", ToChainId: "3", From: from,
		To: to, Value: "0", Input: "0x95000000022c7536e3605d9c16a7a3d7b1898e529396a65c230000000000000019000000034fa1c4e6182b6b7f3bca273390cf587b50b473110000000000045644010102a253a1c03185eec75a271a89a69740c9e1bcaebcefbce87e06f46b8c470aef98d10b425b93941093a1b0dfbdbf5e039a614e6fc5e077e373c8c706fbd529454ee64e9dcae974df7b346bc200008080940a934080c24afd8081000462f62879bcb53487b2b5a7705622002ceef2792208cd5596957e787d413679bc0767b25d85d3de7ba9aae303b755e4c3cafc8e2b2aace4416c4e3dd7f491ebbd20bd4ff241306fd15e11de6aa74d5f6227a4c43a14e5122481a8ccb39eccfd80863435651c07044738dcd9b17e70ff58fe949482de4c8608df9e4a335e276bf30001049424930080c20000c0a3bf8756d9de4122b253672e4592f1360c5c2212751a65e802268064ae6e80d88100054f1eb21e380dcdb2e9373475592a59af5e3c37823777c1977106c49707d0288eadf6ee384dbb670940db93ecd9bf308698157da029742d81f61c4df7f7cca44caa1ba397c25361529a19d50cc737a83b4e5e39ede0bde35f89fcf8b21a9e1407eca078605c1b0ad6ff4323f7c23307585d3dddd504f96e7a7f722f9802d2a1b7e6c7fa9955f3c8ad04a624c45d01dfb82f3d8148fac1be599ab59932295d372d000110", Nonce: strconv.Itoa(int(nonce)),
	}
	privatekey, err := crypto.HexToECDSA(key)
	err = connection.Thk.SignTransaction(&transaction, privatekey)

	txhash, err := connection.Thk.SendTx(&transaction)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Log("txhash:", txhash)
	//0x920a95dc3af9d6ed801258fc8eeb1455b7e6b35a72d4142c995a27f4f0e78c8d
}

func TestThkGetCommittee(t *testing.T) {
	var err error
	var connection = web3.NewWeb3(providers.NewHTTPProvider("thinkey.natapp1.cc", 10, false))

	res, err := connection.Thk.GetCommittee("2", 100)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Log("res", res)
}

func TestThkSendTx(t *testing.T) {
	var err error
	var connection = web3.NewWeb3(providers.NewHTTPProvider("192.168.1.106:8093", 10, false))
	from := "0x2c7536e3605d9c16a7a3d7b1898e529396a65c23"
	to := "0x6ea0fefc17c877c7a4b0f139728ed39dc134a967"
	nonce, err := connection.Thk.GetNonce(from, "2")
	fmt.Println("nonce:", nonce)
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

//Chainid:
//From:
//To:
//Nonce:   nonce,
// ExpireHeight: expireheight,
// Amount: value.(string),
func TestThkRpcMakeVccProof(t *testing.T) {
	var err error
	var connection = web3.NewWeb3(providers.NewHTTPProvider("192.168.1.106:8093", 10, false))
	from := "0x2c7536e3605d9c16a7a3d7b1898e529396a65c23"
	to := "0x0000000000000000000000000000000000020000"

	nonce, err := connection.Thk.GetNonce(from, "2")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	//stats, _ := connection.Thk.GetStats(2)
	expireHeight := 279228 + 5000

	//fmt.Println(stats.Currentheight)

	fmt.Println(expireHeight)

	transaction := util.Transaction{
		ChainId: "2", FromChainId: "2", ToChainId: "3", From: from,
		To: to, Nonce: strconv.Itoa(int(nonce)), Value: "2333", ExpireHeight: expireHeight,
	}
	input, err := connection.Thk.RpcMakeVccProof(&transaction)
	t.Log("input:", input)

}
func TestThkMakeCCCExistenceProof(t *testing.T) {
	var err error
	var connection = web3.NewWeb3(providers.NewHTTPProvider("192.168.1.106:8093", 10, false))
	from := "0x2c7536e3605d9c16a7a3d7b1898e529396a65c23"
	to := "0x0000000000000000000000000000000000020000"

	nonce, err := connection.Thk.GetNonce(from, "3")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	//stats, _ := connection.Thk.GetStats(3)
	expireHeight := 279228 + 5000

	transaction := util.Transaction{
		ChainId: "3", FromChainId: "3", ToChainId: "2", From: from,
		To: to, Nonce: strconv.Itoa(int(nonce)), Value: "2333", ExpireHeight: expireHeight,
	}
	input, err := connection.Thk.MakeCCCExistenceProof(&transaction)
	t.Log("input:", input)
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
	var connection = web3.NewWeb3(providers.NewHTTPProvider("192.168.1.13:8093", 10, false))
	hash := "0xcb53f1ec9c02053a46de488b63b219217826fd9c4cfb531567d61003664ef653"
	res, err := connection.Thk.GetTransactionByHash("2", hash)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Log("res:", res)
}

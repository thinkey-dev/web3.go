package test

import (
	"encoding/binary"
	"fmt"
	"github.com/go-ethereum/crypto"
	"io"
	"math/big"
	"strconv"
	"testing"
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
		ExpireHeight: 33772,
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
		To: to, Value: "2333", Input: str, Nonce: strconv.Itoa(int(nonce)),
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
func TestThkSaveCashCheck(t *testing.T) {
	var err error
	var connection = web3.NewWeb3(providers.NewHTTPProvider("192.168.1.106:8093", 10, false))
	from := "0x2c7536e3605d9c16a7a3d7b1898e529396a65c23"
	to := "0x0000000000000000000000000000000000030000"

	nonce, err := connection.Thk.GetNonce(from, "3")
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
		To: to, Value: "2333", Input: "0x95000000022c7536e3605d9c16a7a3d7b1898e529396a65c230000000000000012000000034fa1c4e6182b6b7f3bca273390cf587b50b4731100000000000083ec010102a280c0c0b207570b7d18d77cf14dda7b145b624b6829a5c871d3b325c1db03fbc1c8776f93941093a1bbdfe64a6c872bce61e40e6f24a4c5b281d606969773ce1a80daec7c4500ccf59ec200008080940e934080c202d6808100039f6cfc2b0f1741dbe88919515314ff10514b28d80db83cf35e8379c1f0b50df99c6b2c1b0d0b25a008e6c882cc7b415f309965c72ad2b944ac0931048ca31cd57855b2fc1b3086df9936adf68d035acaf19b5f8e775b12a611f44d15eb048c820001049424930080c20000c06c1d20d0284e5d443571a62662b289e788298ee754958406ce345b3702949a7c8100051214b76617a782b0b3752582f65b95d54e560b78192bdd83a530dc7e446a6924f2e9abc3465d31d92b6e9677ffc891f609bdc57bc3d125e589177ad63cb2cac979698713be5c13ec858331f021edfe82c733f815b2cde5f030a4816a8986af3eeca078605c1b0ad6ff4323f7c23307585d3dddd504f96e7a7f722f9802d2a1b7701b444bf94b73018944f48a65b45784ea83896e9c00779eba39d5a855cb9667000110", Nonce: strconv.Itoa(int(nonce)),
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

	stats, _ := connection.Thk.GetStats(2)
	expireHeight := stats.Currentheight + 5000

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

	stats, _ := connection.Thk.GetStats(3)
	expireHeight := stats.Currentheight + 5000

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
	var connection = web3.NewWeb3(providers.NewHTTPProvider("thinkey.natapp1.cc", 10, false))
	hash := "0xcb53f1ec9c02053a46de488b63b219217826fd9c4cfb531567d61003664ef653"
	res, err := connection.Thk.GetTransactionByHash("2", hash)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Log("res:", res)
}

package thk

import (
	"crypto/ecdsa"
	"errors"
	"math/big"
	"strings"
	"web3.go/common/cryp/crypto"
	"web3.go/common/cryp/sha3"
	"web3.go/common/hexutil"
	"web3.go/web3/dto"
	"web3.go/web3/providers"
	"web3.go/web3/thk/util"
)

type Thk struct {
	provider providers.ProviderInterface
}

func NewThk(provider providers.ProviderInterface) *Thk {
	thk := new(Thk)
	thk.provider = provider
	return thk
}

func (thk *Thk) GetBalance(address string, chainId string) (*big.Int, error) {
	params := new(util.GetAccountJson)
	if err := params.FormatParams(address, chainId); err != nil {
		return nil, err
	}
	res := make(map[string]interface{})
	if err := thk.provider.SendRequest(&res, "GetAccount", params); err != nil {
		return nil, err
	}

	if _, ok := res["errMsg"]; ok {
		return nil, errors.New(res["errMsg"].(string))
	}
	ret := big.NewInt(int64(res["balance"].(float64)))

	return ret, nil
}

func (thk *Thk) GetNonce(address string, chainId string) (int64, error) {
	params := new(util.GetAccountJson)
	if err := params.FormatParams(address, chainId); err != nil {
		return 0, err
	}
	res := make(map[string]interface{})
	if err := thk.provider.SendRequest(&res, "GetAccount", params); err != nil {
		return 0, err
	}

	if _, ok := res["errMsg"]; ok {
		return 0, errors.New(res["errMsg"].(string))
	}
	ret := int64(res["nonce"].(float64))

	return ret, nil
}

func (thk *Thk) GetBlockTxs(chainId string, height string, page string, size string) {
	params := new(util.GetBlockTxsJson)
	if err := params.FormatParams(chainId, height, page, size); err != nil {
		return
	}
}

func (thk *Thk) SendTx(transaction *util.Transaction) (string, error) {
	// params := new(util.Transaction)
	// if err := params.FormatParams(transaction); err != nil {
	// 	return err
	// }
	res := new(dto.SendTxResult)
	if err := thk.provider.SendRequest(res, "SendTx", transaction); err != nil {
		return "", err
	}
	if res.ErrMsg != "" {
		err := errors.New(res.ErrMsg)
		return "", err
	}
	return res.TXhash, nil
}

func (thk *Thk) SignTransaction(transaction *util.Transaction, privatekey *ecdsa.PrivateKey) error {
	var toaddr string
	if len(transaction.To) > 2 {
		toaddr = transaction.To[2:]
	}
	var input string
	if len(transaction.Input) > 2 {
		input = transaction.Input[2:]
	}

	str := []string{transaction.ChainId, transaction.From[2:], toaddr, transaction.Nonce, transaction.Value, input}
	p := strings.Join(str, "")
	tmp := sha3.NewKeccak256()
	_, err := tmp.Write([]byte(p))
	if err != nil {
		return err
	}
	hash := tmp.Sum(nil)
	sig, err := crypto.Sign(hash, privatekey)
	if err != nil {
		return err
	}
	transaction.Sig = hexutil.Encode(sig)
	transaction.Pub = hexutil.Encode(crypto.FromECDSAPub(&privatekey.PublicKey))
	return nil
}

func (thk *Thk) CallTransaction(transaction *util.Transaction) (*dto.TxResult, error) {
	res := new(dto.TxResult)
	if err := thk.provider.SendRequest(res, "CallTransaction", transaction); err != nil {
		return nil, err
	}
	if res.ErrMsg != "" {
		err := errors.New(res.ErrMsg)
		return nil, err
	}
	return res, nil
}

func (thk *Thk) GetTransactionByHash(chainId string, hash string) (*dto.TxResult, error) {
	params := new(util.GetTxByHash)
	if err := params.FormatParams(chainId, hash); err != nil {
		return nil, err
	}
	res := new(dto.TxResult)
	if err := thk.provider.SendRequest(res, "GetTransactionByHash", params); err != nil {
		return nil, err
	}
	if res.ErrMsg != "" {
		err := errors.New(res.ErrMsg)
		return nil, err
	}
	return res, nil
}

func (thk *Thk) GetBlockHeader(chainId string, height string) (*dto.GetBlockResult, error) {
	params := new(util.GetBlockHeader)
	if err := params.FormatParams(chainId, height); err != nil {
		return nil, err
	}
	res := new(dto.GetBlockResult)
	if err := thk.provider.SendRequest(res, "GetBlockHeader", params); err != nil {
		return nil, err
	}
	if res.ErrMsg != "" {
		err := errors.New(res.ErrMsg)
		return nil, err
	}
	return res, nil
}
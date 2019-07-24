package thk

import (
	"crypto/ecdsa"
	"errors"
	"fmt"
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

//获取余额11
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

//获取之前交易数
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

//	获取块交易11
func (thk *Thk) GetBlockTxs(chainId string, height string, page string, size string) {
	params := new(util.GetBlockTxsJson)
	if err := params.FormatParams(chainId, height, page, size); err != nil {
		return
	}
}

//11
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

//交易签名
func (thk *Thk) SignTransaction(transaction *util.Transaction, privatekey *ecdsa.PrivateKey) error {
	var toAddr string
	var fromAddr string
	if len(transaction.To) > 2 {
		toAddr = transaction.To[2:]

		toAddr = strings.ToLower(toAddr)
	}

	if len(transaction.From) > 2 {
		fromAddr = transaction.From[2:]

		fromAddr = strings.ToLower(fromAddr)
	}


	var input string
	if len(transaction.Input) > 2 {
		input = transaction.Input[2:]
	}

	str := []string{transaction.ChainId, fromAddr, toAddr, transaction.Nonce, transaction.Value, input}
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

//
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

//通过hash获取交易11
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

//获取块结果11
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

//11
func (thk *Thk) Ping(chainId string) (int64, error) {
	params := new(util.PingJson)
	if err := params.FormatParams(chainId); err != nil {
		return 0, err
	}
	res := make(map[string]interface{})
	if err := thk.provider.SendRequest(&res, "Ping", params); err != nil {
		return 0, err
	}

	if _, ok := res["errMsg"]; ok {
		return 0, errors.New(res["errMsg"].(string))
	}
	ret := int64(res["nonce"].(float64))

	return ret, nil
}

// func (thk *Thk) GetChainInfo(chainId string) {
// 	params := new(util.GetChainInfoJson)
// 	if err := params.FormatParams(chainId); err != nil {
// 		return 0, err
// 	}
// 	res := make(map[string]interface{})
// 	if err := thk.provider.SendRequest(&res, "Ping", params); err != nil {
// 		return 0, err
// 	}
//
// 	if _, ok := res["errMsg"]; ok {
// 		return 0, errors.New(res["errMsg"].(string))
// 	}
// 	ret := int64(res["nonce"].(float64))
//
// 	return ret, nil
// }
//19.5.25 获取链信息11
func (thk *Thk) GetChainInfo(chainIds []int) ([]dto.GetChainInfo, error) {
	params := new(util.GetChainInfoJson)
	if err := params.FormatParams(chainIds); err != nil {
		return nil, err
	}
	res := new(dto.GetChainInfo)
	if err := thk.provider.SendRequest(res, "GetChainInfo", params); err != nil {
		return nil, err
	}
	if res.ErrMsg != "" {
		err := errors.New(res.ErrMsg)
		return nil, err
	}

	res_array := []dto.GetChainInfo{*res}
	return res_array, nil
}

//11
func (thk *Thk) GetStats(chainId int) (gts dto.GetChainStats, err error) {
	params := new(util.GetStatsJson)
	ers := params.FormatParams(chainId)
	if ers != nil {
		fmt.Println(ers)
	}

	res := new(dto.GetChainStats)
	if err := thk.provider.SendRequest(res, "GetStats", params); err != nil {
		return *res, err
	}
	res_array := dto.GetChainStats{ChainId: chainId}

	return res_array, nil

}

//GetTransactions
func (thk *Thk) GetTransactions(chainId, address, startHeight, endHeight string) ([]dto.GetTransactions, error) {
	params := new(util.GetTransactionsJson)
	if err := params.FormatParams(chainId, address, startHeight, endHeight); err != nil {
		return nil, err
	}

	res := new(dto.GetTransactions)
	if err := thk.provider.SendRequest(res, "GetTransactions", params); err != nil {
		return nil, err
	}

	res_array := []dto.GetTransactions{*res}
	return res_array, nil

}

//5.25 获取委员会详情11
func (thk *Thk) GetCommittee(chainId string, epoch int) ([]string, error) {
	params := new(util.GetCommitteeJson)
	if err := params.FormatParams(chainId, epoch); err != nil {
		return nil, err
	}

	res := new(dto.GetCommittee)
	if err := thk.provider.SendRequest(res, "GetCommittee", params); err != nil {
		return nil, err
	}
	if res.ErrMsg != "" {
		err := errors.New(res.ErrMsg)
		return nil, err
	}
	return res.MemberDetails, nil
}

//RpcMakeVccProof 11
func (thk *Thk) RpcMakeVccProof(transaction *util.Transaction) (map[string]interface{}, error) {
	res := new(dto.RpcMakeVccProofJson)
	if err := thk.provider.SendRequest(res, "RpcMakeVccProof", transaction); err != nil {
		return nil, err
	}
	if res.ErrMsg != "" {
		err := errors.New(res.ErrMsg)
		return nil, err
	}
	return res.Proof, nil
}

//MakeCCCExistenceProof  11
func (thk *Thk) MakeCCCExistenceProof(transaction *util.Transaction) (map[string]interface{}, error) {
	res := new(dto.MakeCCCExistenceProofJson)
	if err := thk.provider.SendRequest(res, "MakeCCCExistenceProof", transaction); err != nil {
		return nil, err
	}
	if res.ErrMsg != "" {
		err := errors.New(res.ErrMsg)
		return nil, err
	}
	return res.Proof, nil
}

//GetCCCRelativeTx

package util

type GetAccountJson struct {
	Address string `json:"address"`
	ChainId string `json:"chainId"`
}

type GetBlockTxsJson struct {
	ChainId string `json:"chainId"`
	Height  string `json:"height"`
	Page    string `json:"page"`
	Size    string `json:"size"`
}

type Transaction struct {
	ChainId      string `json:"chainId"`
	FromChainId  string `json:"fromChainId,omitempty"`
	ToChainId    string `json:"toChainId,omitempty"`
	From         string `json:"from"`
	To           string `json:"to"`
	Nonce        string `json:"nonce"`
	Value        string `json:"value"`
	Sig          string `json:"sig,omitempty"`
	Pub          string `json:"pub,omitempty"`
	Input        string `json:"input"`
	ExpireHeight int    `json:"expireheight"`
}

type GetTxByHash struct {
	ChainId string `json:"chainId"`
	Hash    string `json:"hash"`
}

type GetBlockHeader struct {
	ChainId string `json:"chainId"`
	Height  string `json:"height"`
}

type PingJson struct {
	ChainId string `json:"chainId"`
}

type GetChainInfoJson struct {
	ChainId []int `json:"chainId"`
}

type GetStatsJson struct {
	ChainId int `json:"chainId"`
}

type GetTransactionsJson struct {
	ChainId     string `json:"chainId"`
	Address     string `json:"address"`
	StartHeight string `json:"startHeight"`
	EndHeight   string `json:"endHeight"`
}

type GetMultiStatsJson struct {
	ChainId string `json:"chainId"`
}

type GetCommitteeJson struct {
	ChainId string `json:"chainId"`
	Epoch   int    `json:"epoch"`
}
type CompileContractJson struct {
	ChainId  string `json:"chainId"`
	Contract string `json:"contract"`
}

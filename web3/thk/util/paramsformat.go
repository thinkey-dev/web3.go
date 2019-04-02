package util

func (param *GetAccountJson) FormatParams(address string, chainid string) error {
	param.ChainId = chainid
	param.Address = address
	return nil
}

func (param *GetBlockTxsJson) FormatParams(chainId string, height string, page string, size string) error {
	param.ChainId = chainId
	param.Height = height
	param.Page = page
	param.Size = size
	return nil
}

func (param *Transaction) FormatParams(transcation *Transaction) error {
	return nil
}

func (param *GetTxByHash) FormatParams(chainId string, hash string) error {
	param.ChainId = chainId
	param.Hash = hash
	return nil
}

func (param *GetBlockHeader) FormatParams(chainId string, height string) error {
	param.ChainId = chainId
	param.Height = height
	return nil
}

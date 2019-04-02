package utils

import (
	"web3.go/web3/complex/types"
	"web3.go/web3/dto"
	"web3.go/web3/providers"
)

type Utils struct {
	provider providers.ProviderInterface
}

func NewUtils(provider providers.ProviderInterface) *Utils {
	utils := new(Utils)
	utils.provider = provider
	return utils
}

func (utils *Utils) Sha3(data types.ComplexString) (string, error) {

	params := make([]string, 1)
	params[0] = data.ToHex()

	pointer := &dto.RequestResult{}

	err := utils.provider.SendRequest(pointer, "web3_sha3", params)

	if err != nil {
		return "", err
	}

	return pointer.ToString()

}

package net

import (
	"math/big"
	"web3.go/web3/dto"
	"web3.go/web3/providers"
)

type Net struct {
	provider providers.ProviderInterface
}

func NewNet(provider providers.ProviderInterface) *Net {
	net := new(Net)
	net.provider = provider
	return net
}

func (net *Net) IsListening() (bool, error) {

	pointer := &dto.RequestResult{}

	err := net.provider.SendRequest(pointer, "net_listening", nil)

	if err != nil {
		return false, err
	}

	return pointer.ToBoolean()

}

func (net *Net) GetPeerCount() (*big.Int, error) {

	pointer := &dto.RequestResult{}

	err := net.provider.SendRequest(pointer, "net_peerCount", nil)

	if err != nil {
		return nil, err
	}

	return pointer.ToBigInt()

}

func (net *Net) GetVersion() (string, error) {

	pointer := &dto.RequestResult{}

	err := net.provider.SendRequest(pointer, "net_version", nil)

	if err != nil {
		return "", err
	}

	return pointer.ToString()

}

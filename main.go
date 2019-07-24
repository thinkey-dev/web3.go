package main

import (
	"fmt"
	"web3.go/web3"
	"web3.go/web3/providers"
)

func main() {
	var err error
	var connection = web3.NewWeb3(providers.NewHTTPProvider("192.168.1.13:8089", 10, false))
	from := "0x2c7536e3605d9c16a7a3d7b1898e529396a65c23"

	nonce, err := connection.Thk.GetNonce(from, "2")
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(nonce)
	}

}

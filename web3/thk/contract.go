package thk

import (
	"crypto/ecdsa"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"strings"
	"web3.go/common/hexutil"
	"web3.go/web3/complex/types"
	"web3.go/web3/dto"
	"web3.go/web3/thk/abi"
	"web3.go/web3/thk/util"
)

type Contract struct {
	super     *Thk
	abi       abi.ABI
	functions map[string][]string
}

func (thk *Thk) NewContract(abistr string) (*Contract, error) {

	contract := new(Contract)
	var mockInterface interface{}

	err := json.Unmarshal([]byte(abistr), &mockInterface)

	if err != nil {
		return nil, err
	}

	jsonInterface := mockInterface.([]interface{})
	contract.functions = make(map[string][]string)
	for index := 0; index < len(jsonInterface); index++ {
		function := jsonInterface[index].(map[string]interface{})

		if function["type"] == "constructor" || function["type"] == "fallback" {
			function["name"] = function["type"]
		}

		functionName := function["name"].(string)
		contract.functions[functionName] = make([]string, 0)

		if function["inputs"] == nil {
			continue
		}

		inputs := function["inputs"].([]interface{})
		for paramIndex := 0; paramIndex < len(inputs); paramIndex++ {
			params := inputs[paramIndex].(map[string]interface{})
			contract.functions[functionName] = append(contract.functions[functionName], params["type"].(string))
		}

	}
	readerstr := strings.NewReader(abistr)
	Abi, err := abi.JSON(readerstr)
	if err != nil {
		return nil, err
	}
	contract.abi = Abi
	contract.super = thk

	return contract, nil
}

// func (contract *Contract) prepareTransaction(transaction util.Transaction, functionName string, args []interface{}) (util.Transaction, error) {
//
// 	function, ok := contract.functions[functionName]
// 	if !ok {
// 		return transaction, errors.New("Function not finded on passed abi")
// 	}
//
// 	fullFunction := functionName + "("
//
// 	comma := ""
// 	for arg := range function {
// 		fullFunction += comma + function[arg]
// 		comma = ","
// 	}
//
// 	fullFunction += ")"
//
// 	hash := sha3.NewKeccak256()
// 	hash.Write([]byte(fullFunction))
// 	inputbyte := hash.Sum(nil)
// 	input := hexutil.Encode(inputbyte)
// 	var data string
//
// 	// for index := 0; index < len(function); index++ {
// 	// 	currentData, err := contract.getHexValue(function[index], args[index])
// 	//
// 	// 	if err != nil {
// 	// 		return transaction, err
// 	// 	}
// 	//
// 	// 	data += currentData
// 	// }
//
// 	transaction.Input = string(types.ComplexString(input[0:10] + data))
//
// 	return transaction, nil
//
// }

func (contract *Contract) getHexValue(inputType string, value interface{}) (string, error) {

	var data string

	if strings.HasPrefix(inputType, "int") ||
		strings.HasPrefix(inputType, "uint") ||
		strings.HasPrefix(inputType, "fixed") ||
		strings.HasPrefix(inputType, "ufixed") {

		bigVal := value.(*big.Int)

		// Checking that the string actually is the correct inputType
		if strings.Contains(inputType, "128") {
			// 128 bit
			if bigVal.BitLen() > 128 {
				return "", errors.New(fmt.Sprintf("Input type %s not met", inputType))
			}
		} else if strings.Contains(inputType, "256") {
			// 256 bit
			if bigVal.BitLen() > 256 {
				return "", errors.New(fmt.Sprintf("Input type %s not met", inputType))
			}
		}

		data += fmt.Sprintf("%064s", fmt.Sprintf("%x", bigVal.String()))
	}

	if strings.Compare("address", inputType) == 0 {
		data += fmt.Sprintf("%064s", value.(string)[2:])
	}

	if strings.Compare("string", inputType) == 0 {
		data += fmt.Sprintf("%064s", fmt.Sprintf("%x", value.(string)))
	}

	return data, nil

}

func (contract *Contract) Send(transaction util.Transaction, functionName string, privatekey *ecdsa.PrivateKey, args ...interface{}) (string, error) {

	// transaction, err := contract.prepareTransaction(transaction, functionName, args)
	fixedArrStrPack, err := contract.abi.Pack(functionName, args...)
	if err != nil {
		return "", err
	}
	transaction.Input = hexutil.Encode(fixedArrStrPack)
	if err = contract.super.SignTransaction(&transaction, privatekey); err != nil {
		return "", err
	}
	return contract.super.SendTx(&transaction)

}

func (contract *Contract) Deploy(transaction util.Transaction, bytecode string, privatekey *ecdsa.PrivateKey, args ...interface{}) (string, error) {

	// constructor := contract.functions["constructor"]
	constructor := contract.functions["constructor"]

	for index := 0; index < len(constructor); index++ {
		tmpBytes, err := contract.getHexValue(constructor[index], args[index])

		if err != nil {
			return "", err
		}

		bytecode += tmpBytes
	}

	transaction.Input = string(types.ComplexString(bytecode))
	err := contract.super.SignTransaction(&transaction, privatekey)
	if err != nil {
		return "", err
	}
	return contract.super.SendTx(&transaction)

}

func (contract *Contract) Call(transaction util.Transaction, functionName string, args ...interface{}) (*dto.TxResult, error) {

	// transaction, err := contract.prepareTransaction(transaction, functionName, args)
	fixedArrStrPack, err := contract.abi.Pack(functionName, args...)
	if err != nil {
		return nil, err
	}
	transaction.Input = hexutil.Encode(fixedArrStrPack)
	return contract.super.CallTransaction(&transaction)

}

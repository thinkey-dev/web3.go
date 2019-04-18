package abi

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"github.com/Alex-Chris/log/log"
	"github.com/ethereum/go-ethereum/common"
	"math/big"
	"reflect"
	"strings"
	"testing"
)

type ParamJson struct {
	Name  string `json:"name"`  // 参数名称
	Type  string `json:"type"`  // 参数类型
	Value string `json:"value"` // 参数值
}

type MethodJson struct {
	Method string      `json:"method"`           // 方法名称
	Params []ParamJson `json:"params,omitempty"` // 方法参数列表
}

func TestUnpackEventIntoMap(t *testing.T) {
	const abiJSON = `[{"constant":true,"inputs":[],"name":"getObjsIndex","outputs":[{"name":"curOrderIndex","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"name":"_address","type":"address"}],"name":"removeAdmin","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[{"name":"objID","type":"uint256"}],"name":"ifObjectInRange","outputs":[{"name":"","type":"bool"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[{"name":"objID","type":"uint256"}],"name":"getObjById","outputs":[{"components":[{"name":"orderId","type":"string"},{"name":"date","type":"string"},{"name":"orderType","type":"string"},{"name":"company","type":"string"},{"name":"organization","type":"string"},{"name":"department","type":"string"},{"name":"salesman","type":"string"},{"name":"buyer","type":"string"},{"name":"project","type":"string"},{"name":"quantity","type":"uint64"},{"name":"weight","type":"uint64"},{"name":"feeAmount","type":"uint64"},{"name":"orderStatus","type":"uint64"},{"name":"takenQuantity","type":"uint64"},{"name":"takenWeight","type":"uint64"},{"name":"takenOrderId","type":"string"},{"name":"takenDate","type":"string"},{"name":"detail","type":"string"}],"name":"obj","type":"tuple"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"components":[{"name":"orderId","type":"string"},{"name":"date","type":"string"},{"name":"orderType","type":"string"},{"name":"company","type":"string"},{"name":"organization","type":"string"},{"name":"department","type":"string"},{"name":"salesman","type":"string"},{"name":"buyer","type":"string"},{"name":"project","type":"string"},{"name":"quantity","type":"uint64"},{"name":"weight","type":"uint64"},{"name":"feeAmount","type":"uint64"},{"name":"orderStatus","type":"uint64"},{"name":"takenQuantity","type":"uint64"},{"name":"takenWeight","type":"uint64"},{"name":"takenOrderId","type":"string"},{"name":"takenDate","type":"string"},{"name":"detail","type":"string"}],"name":"orderObj","type":"tuple"}],"name":"updateObj","outputs":[{"name":"result","type":"bool"}],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"name":"_newAddress","type":"address"}],"name":"setAdmin","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[{"name":"_userAddress","type":"address"}],"name":"isExitAddress","outputs":[{"name":"isIndeed","type":"bool"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[],"name":"isOrderShare","outputs":[{"name":"","type":"bool"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[{"name":"_orderId","type":"string"}],"name":"getObjByOrderId","outputs":[{"components":[{"name":"orderId","type":"string"},{"name":"date","type":"string"},{"name":"orderType","type":"string"},{"name":"company","type":"string"},{"name":"organization","type":"string"},{"name":"department","type":"string"},{"name":"salesman","type":"string"},{"name":"buyer","type":"string"},{"name":"project","type":"string"},{"name":"quantity","type":"uint64"},{"name":"weight","type":"uint64"},{"name":"feeAmount","type":"uint64"},{"name":"orderStatus","type":"uint64"},{"name":"takenQuantity","type":"uint64"},{"name":"takenWeight","type":"uint64"},{"name":"takenOrderId","type":"string"},{"name":"takenDate","type":"string"},{"name":"detail","type":"string"}],"name":"obj","type":"tuple"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"name":"callAddress","type":"address"}],"name":"setCallerAddress","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"components":[{"name":"orderId","type":"string"},{"name":"date","type":"string"},{"name":"orderType","type":"string"},{"name":"company","type":"string"},{"name":"organization","type":"string"},{"name":"department","type":"string"},{"name":"salesman","type":"string"},{"name":"buyer","type":"string"},{"name":"project","type":"string"},{"name":"quantity","type":"uint64"},{"name":"weight","type":"uint64"},{"name":"feeAmount","type":"uint64"},{"name":"orderStatus","type":"uint64"},{"name":"takenQuantity","type":"uint64"},{"name":"takenWeight","type":"uint64"},{"name":"takenOrderId","type":"string"},{"name":"takenDate","type":"string"},{"name":"detail","type":"string"}],"name":"orderObj","type":"tuple"}],"name":"createObj","outputs":[{"name":"totalNum","type":"uint256"}],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[{"name":"_offset","type":"uint64"},{"name":"_limit","type":"uint64"}],"name":"getAllObjs","outputs":[{"name":"totalNum","type":"uint256"},{"name":"offset","type":"uint64"},{"name":"limit","type":"uint64"},{"components":[{"name":"orderId","type":"string"},{"name":"date","type":"string"},{"name":"orderType","type":"string"},{"name":"company","type":"string"},{"name":"organization","type":"string"},{"name":"department","type":"string"},{"name":"salesman","type":"string"},{"name":"buyer","type":"string"},{"name":"project","type":"string"},{"name":"quantity","type":"uint64"},{"name":"weight","type":"uint64"},{"name":"feeAmount","type":"uint64"},{"name":"orderStatus","type":"uint64"},{"name":"takenQuantity","type":"uint64"},{"name":"takenWeight","type":"uint64"},{"name":"takenOrderId","type":"string"},{"name":"takenDate","type":"string"},{"name":"detail","type":"string"}],"name":"objs","type":"tuple[]"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[],"name":"getObjsNum","outputs":[{"name":"totalNum","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"name":"_orderId","type":"string"}],"name":"deleteObj","outputs":[{"name":"result","type":"bool"}],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[],"name":"caller","outputs":[{"name":"","type":"address"}],"payable":false,"stateMutability":"view","type":"function"},{"inputs":[],"payable":false,"stateMutability":"nonpayable","type":"constructor"},{"anonymous":false,"inputs":[{"indexed":false,"name":"orderID","type":"string"}],"name":"NewObject","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"name":"orderID","type":"string"}],"name":"ModifyObject","type":"event"}]`
	abi, err := JSON(strings.NewReader(abiJSON))
	if err != nil {
		t.Fatal(err)
	}

	const hexdata = ``
	data, err := hex.DecodeString(hexdata)
	if err != nil {
		t.Fatal(err)
	}
	if len(data)%32 == 0 {
		t.Errorf("len(data) is %d, want a non-multiple of 32", len(data))
	}

	receivedMap := map[string]interface{}{}
	expectedReceivedMap := map[string]interface{}{
		"sender": common.HexToAddress("0x376c47978271565f56DEB45495afa69E59c16Ab2"),
		"amount": big.NewInt(1),
		"memo":   []byte{88},
	}
	if err := abi.UnpackIntoMap(receivedMap, "getObjById", data); err != nil {
		t.Error(err)
	}
	if len(receivedMap) != 3 {
		t.Error("unpacked `received` map expected to have length 3")
	}
	if receivedMap["sender"] != expectedReceivedMap["sender"] {
		t.Error("unpacked `received` map does not match expected map")
	}
	if receivedMap["amount"].(*big.Int).Cmp(expectedReceivedMap["amount"].(*big.Int)) != 0 {
		t.Error("unpacked `received` map does not match expected map")
	}
	if !bytes.Equal(receivedMap["memo"].([]byte), expectedReceivedMap["memo"].([]byte)) {
		t.Error("unpacked `received` map does not match expected map")
	}

	receivedAddrMap := map[string]interface{}{}
	if err = abi.UnpackIntoMap(receivedAddrMap, "receivedAddr", data); err != nil {
		t.Error(err)
	}
	if len(receivedAddrMap) != 1 {
		t.Error("unpacked `receivedAddr` map expected to have length 1")
	}
	if receivedAddrMap["sender"] != expectedReceivedMap["sender"] {
		t.Error("unpacked `receivedAddr` map does not match expected map")
	}
}

func TestDecodeAbi(t *testing.T) {

	const definition = `[{"constant":true,"inputs":[],"name":"getObjsIndex","outputs":[{"name":"curOrderIndex","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"name":"_address","type":"address"}],"name":"removeAdmin","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[{"name":"objID","type":"uint256"}],"name":"ifObjectInRange","outputs":[{"name":"","type":"bool"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[{"name":"objID","type":"uint256"}],"name":"getObjById","outputs":[{"components":[{"name":"orderId","type":"string"},{"name":"date","type":"string"},{"name":"orderType","type":"string"},{"name":"company","type":"string"},{"name":"organization","type":"string"},{"name":"department","type":"string"},{"name":"salesman","type":"string"},{"name":"buyer","type":"string"},{"name":"project","type":"string"},{"name":"quantity","type":"uint64"},{"name":"weight","type":"uint64"},{"name":"feeAmount","type":"uint64"},{"name":"orderStatus","type":"uint64"},{"name":"takenQuantity","type":"uint64"},{"name":"takenWeight","type":"uint64"},{"name":"takenOrderId","type":"string"},{"name":"takenDate","type":"string"},{"name":"detail","type":"string"}],"name":"obj","type":"tuple"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"components":[{"name":"orderId","type":"string"},{"name":"date","type":"string"},{"name":"orderType","type":"string"},{"name":"company","type":"string"},{"name":"organization","type":"string"},{"name":"department","type":"string"},{"name":"salesman","type":"string"},{"name":"buyer","type":"string"},{"name":"project","type":"string"},{"name":"quantity","type":"uint64"},{"name":"weight","type":"uint64"},{"name":"feeAmount","type":"uint64"},{"name":"orderStatus","type":"uint64"},{"name":"takenQuantity","type":"uint64"},{"name":"takenWeight","type":"uint64"},{"name":"takenOrderId","type":"string"},{"name":"takenDate","type":"string"},{"name":"detail","type":"string"}],"name":"orderObj","type":"tuple"}],"name":"updateObj","outputs":[{"name":"result","type":"bool"}],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"name":"_newAddress","type":"address"}],"name":"setAdmin","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[{"name":"_userAddress","type":"address"}],"name":"isExitAddress","outputs":[{"name":"isIndeed","type":"bool"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[],"name":"isOrderShare","outputs":[{"name":"","type":"bool"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[{"name":"_orderId","type":"string"}],"name":"getObjByOrderId","outputs":[{"components":[{"name":"orderId","type":"string"},{"name":"date","type":"string"},{"name":"orderType","type":"string"},{"name":"company","type":"string"},{"name":"organization","type":"string"},{"name":"department","type":"string"},{"name":"salesman","type":"string"},{"name":"buyer","type":"string"},{"name":"project","type":"string"},{"name":"quantity","type":"uint64"},{"name":"weight","type":"uint64"},{"name":"feeAmount","type":"uint64"},{"name":"orderStatus","type":"uint64"},{"name":"takenQuantity","type":"uint64"},{"name":"takenWeight","type":"uint64"},{"name":"takenOrderId","type":"string"},{"name":"takenDate","type":"string"},{"name":"detail","type":"string"}],"name":"obj","type":"tuple"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"name":"callAddress","type":"address"}],"name":"setCallerAddress","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"components":[{"name":"orderId","type":"string"},{"name":"date","type":"string"},{"name":"orderType","type":"string"},{"name":"company","type":"string"},{"name":"organization","type":"string"},{"name":"department","type":"string"},{"name":"salesman","type":"string"},{"name":"buyer","type":"string"},{"name":"project","type":"string"},{"name":"quantity","type":"uint64"},{"name":"weight","type":"uint64"},{"name":"feeAmount","type":"uint64"},{"name":"orderStatus","type":"uint64"},{"name":"takenQuantity","type":"uint64"},{"name":"takenWeight","type":"uint64"},{"name":"takenOrderId","type":"string"},{"name":"takenDate","type":"string"},{"name":"detail","type":"string"}],"name":"orderObj","type":"tuple"}],"name":"createObj","outputs":[{"name":"totalNum","type":"uint256"}],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[{"name":"_offset","type":"uint64"},{"name":"_limit","type":"uint64"}],"name":"getAllObjs","outputs":[{"name":"totalNum","type":"uint256"},{"name":"offset","type":"uint64"},{"name":"limit","type":"uint64"},{"components":[{"name":"orderId","type":"string"},{"name":"date","type":"string"},{"name":"orderType","type":"string"},{"name":"company","type":"string"},{"name":"organization","type":"string"},{"name":"department","type":"string"},{"name":"salesman","type":"string"},{"name":"buyer","type":"string"},{"name":"project","type":"string"},{"name":"quantity","type":"uint64"},{"name":"weight","type":"uint64"},{"name":"feeAmount","type":"uint64"},{"name":"orderStatus","type":"uint64"},{"name":"takenQuantity","type":"uint64"},{"name":"takenWeight","type":"uint64"},{"name":"takenOrderId","type":"string"},{"name":"takenDate","type":"string"},{"name":"detail","type":"string"}],"name":"objs","type":"tuple[]"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[],"name":"getObjsNum","outputs":[{"name":"totalNum","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"name":"_orderId","type":"string"}],"name":"deleteObj","outputs":[{"name":"result","type":"bool"}],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[],"name":"caller","outputs":[{"name":"","type":"address"}],"payable":false,"stateMutability":"view","type":"function"},{"inputs":[],"payable":false,"stateMutability":"nonpayable","type":"constructor"},{"anonymous":false,"inputs":[{"indexed":false,"name":"orderID","type":"string"}],"name":"NewObject","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"name":"orderID","type":"string"}],"name":"ModifyObject","type":"event"}]`
	payload := "0xced1c24e00000000000000000000000000000000000000000000000000000000000000200000000000000000000000000000000000000000000000000000000000000240000000000000000000000000000000000000000000000000000000000000028000000000000000000000000000000000000000000000000000000000000002c000000000000000000000000000000000000000000000000000000000000003000000000000000000000000000000000000000000000000000000000000000340000000000000000000000000000000000000000000000000000000000000038000000000000000000000000000000000000000000000000000000000000003c000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000000440000000000000000000000000000000000000000000000000000000000000000a00000000000000000000000000000000000000000000000000000000000000c800000000000000000000000000000000000000000000000000000000000000c8000000000000000000000000000000000000000000000000000000000000000100000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000014000000000000000000000000000000000000000000000000000000000000048000000000000000000000000000000000000000000000000000000000000004c0000000000000000000000000000000000000000000000000000000000000050000000000000000000000000000000000000000000000000000000000000000185848544430303031303130313138313031322d303030323200000000000000000000000000000000000000000000000000000000000000000000000000000013323031382f31302f31322020303a30303a30300000000000000000000000000000000000000000000000000000000000000000000000000000000000000000045848544400000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000001ee5a4a7e6b189e794b5e5ad90e59586e58aa1e69c89e99990e585ace58fb80000000000000000000000000000000000000000000000000000000000000000001ee5a4a7e6b189e794b5e5ad90e59586e58aa1e69c89e99990e585ace58fb80000000000000000000000000000000000000000000000000000000000000000000ce4b89ae58aa1e983a8e997a800000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000006e5bca0e4b8890000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000ce4b88be6b8b8e4b9b0e5aeb60000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000ce4b88be6b8b8e4b9b0e5aeb6000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000185848434b30303031303130313138313232382d30303130380000000000000000000000000000000000000000000000000000000000000000000000000000000f323031382f31302f313220303a3030000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000001345b7b226c6f744e756d626572223a225848544430303031303130313138313031322d3030303031303031222c2273746f7265686f7573654e616d65223a22e6b996e58d97e4b880e58a9be5ba93222c22676f6f644e616d65223a22e89ebae7bab9e992a2222c226d6174657269616c223a22485242333335222c2273706563696669636174696f6e223a2231322a39222c22706c6163654f664f726967696e223a22e99e8de992a2222c227072696365223a3130302c22746178496e636c75646564416d6f756e74223a323030302c22666565223a32303030302c22666565416d6f756e74223a31302c2273616c655175616e74697479223a3130302c2273616c65576569676874223a31302c2274616b656e5175616e74697479223a3230302c2274616b656e576569676874223a3230307d5d000000000000000000000000"
	name := "createObj"

	// 解析Abi格式成为Json格式
	abiDecoder, err := JSON(strings.NewReader(definition))
	if err != nil {
		t.Error(err)
	}

	// 剔除最前面的0x标记
	var decodeString string = ""
	hexFlag := strings.Index(payload, "0x")
	if hexFlag == -1 {
		decodeString = payload
	} else {
		decodeString = payload[2:]
	}

	// 将字符串转换成[]Byte
	decodeBytes, err := hex.DecodeString(decodeString)
	if err != nil {
		t.Error(err)
	}

	// 根据函数的名称，设置函数的输入参数信息
	method, ok := abiDecoder.Methods[name]
	if !ok {
		t.Error(err)
	}

	// 写入获取参数类型
	params := make([]interface{}, 0)
	for i := 0; i < len(method.Inputs); i++ {

		input := method.Inputs[i]

		// 设置参数类型
		var param reflect.Value

		switch input.Type.T {
		case SliceTy:

			var paramType = []byte{}
			param = reflect.New(reflect.TypeOf(paramType))
		case StringTy:

			var paramType = string("")
			param = reflect.New(reflect.TypeOf(paramType))
		case IntTy:

			switch input.Type.Type {

			case reflect.TypeOf(int8(0)):

				var paramType = int8(0)
				param = reflect.New(reflect.TypeOf(paramType))
			case reflect.TypeOf(int16(0)):

				var paramType = int16(0)
				param = reflect.New(reflect.TypeOf(paramType))
			case reflect.TypeOf(int32(0)):

				var paramType = int32(0)
				param = reflect.New(reflect.TypeOf(paramType))
			case reflect.TypeOf(int64(0)):

				var paramType = int64(0)
				param = reflect.New(reflect.TypeOf(paramType))
			case reflect.TypeOf(&big.Int{}):

				var paramType = big.NewInt(1)
				param = reflect.New(reflect.TypeOf(paramType))
			}
		case UintTy:

			switch input.Type.Type {

			case reflect.TypeOf(uint8(0)):

				var paramType = uint8(0)
				param = reflect.New(reflect.TypeOf(paramType))
			case reflect.TypeOf(uint16(0)):

				var paramType = uint16(0)
				param = reflect.New(reflect.TypeOf(paramType))
			case reflect.TypeOf(uint32(0)):

				var paramType = uint32(0)
				param = reflect.New(reflect.TypeOf(paramType))
			case reflect.TypeOf(uint64(0)):

				var paramType = uint64(0)
				param = reflect.New(reflect.TypeOf(paramType))
			case reflect.TypeOf(&big.Int{}):

				var paramType = big.NewInt(1)
				param = reflect.New(reflect.TypeOf(paramType))
			}
		case BoolTy:

			log.Info("DecodeAbi abi.BoolTy")
			var paramType = bool(true)
			param = reflect.New(reflect.TypeOf(paramType))
		case AddressTy:

			log.Info("DecodeAbi abi.AddressTy")
			var paramType = common.Address{}
			param = reflect.New(reflect.TypeOf(paramType))
		case BytesTy:

			log.Info("DecodeAbi abi.BytesTy")
			var paramType = common.Hex2Bytes("0100000000000000000000000000000000000000000000000000000000000000")
			param = reflect.New(reflect.TypeOf(paramType))
		case FunctionTy:
			log.Info("DecodeAbi abi.FunctionTy")
			var paramType = common.Hex2Bytes("0100000000000000000000000000000000000000000000000000000000000000")
			param = reflect.New(reflect.TypeOf(paramType))
		}
		params = append(params, param.Interface())
	}

	// 解码
	if err := abiDecoder.InputUnpack(params, name, decodeBytes[4:]); err != nil {
		log.Error("DecodeAbi ", "err", err)
		t.Error(err)
	}

	// 将返回的信息放入到Json格式中
	json := MethodJson{}
	json.Method = name
	json.Params = make([]ParamJson, 0)

	for i := 0; i < len(params); i++ {

		valueOf := reflect.ValueOf(params[i])
		out := valueOf.Elem().Interface()
		s := fmt.Sprintf("%v", out)

		param := ParamJson{
			Name:  abiDecoder.Methods[name].Inputs[i].Name,
			Type:  abiDecoder.Methods[name].Inputs[i].Type.String(),
			Value: s,
		}
		json.Params = append(json.Params, param)
	}
	log.Info("DecodeAbi ", "Json", json)

	/*for i := 0; i < len(params); i++ {

		valueOf := reflect.ValueOf(params[i])
		out := valueOf.Elem().Interface()
		s := fmt.Sprintf("%v", out)
		outArray = append(outArray, s)
		log.Info("DecodeAbi ", "Param", s)
	}*/

	log.Info(json)
}

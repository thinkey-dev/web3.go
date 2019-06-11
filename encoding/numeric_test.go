package encoding

import (
	"encoding/hex"
	"fmt"
	"math"
	"reflect"
	"testing"
)

type testvalue struct {
	L int
	N bool
}

var (
	u64 = map[uint64]int{0: 0, 0xd4d314ff89a03c6d: 8, 0x8b3: 2, math.MaxUint64: 8}
	u32 = map[uint32]int{0: 0, 0x89a03c6d: 4, 0xa3: 1, math.MaxUint32: 4}
	u16 = map[uint16]int{0: 0, 0x3c6d: 2, 0x4: 1, 0xf00: 2, math.MaxUint16: 2}
	u8  = map[uint8]int{0: 0, 0x19: 1, 0x6: 1, math.MaxUint8: 1}

	i64 = map[int64]testvalue{0: {0, false}, math.MaxInt64: {8, false}, math.MinInt64: {8, true}, 876234434: {4, false}, -8764234434: {5, true}}
	i32 = map[int32]testvalue{0: {0, false}, math.MaxInt32: {4, false}, math.MinInt32: {4, true}, 2234: {2, false}, -5234: {2, true}}
	i16 = map[int16]testvalue{0: {0, false}, math.MaxInt16: {2, false}, math.MinInt16: {2, true}, 443: {2, false}, -4: {1, true}}
	i8  = map[int8]testvalue{0: {0, false}, math.MaxInt8: {1, false}, math.MinInt8: {1, true}, 9: {1, false}, -32: {1, true}}
)

func testuint(value interface{}, t *testing.T, callee reflect.Value, methodName string) {
	vv := reflect.ValueOf(value)
	vkeys := vv.MapKeys()
	for _, v := range vkeys {
		length := int(vv.MapIndex(v).Int())

		b := Numeric.UintToBytes(v.Uint())
		l := 0
		if b != nil {
			l = len(b)
		}
		if l != length {
			t.Error("UintToBytes length error:", hex.EncodeToString(ToBinaryBytes(v.Interface())), "->", hex.EncodeToString(b))
		}
		m := callee.MethodByName(methodName)
		called := m.Call([]reflect.Value{reflect.ValueOf(b)})
		if called[0].Uint() != v.Uint() {
			t.Error(methodName, "codec error:", hex.EncodeToString(b), "->", hex.EncodeToString(ToBinaryBytes(v.Interface())))
		} else {
			t.Log("PASS: ", hex.EncodeToString(ToBinaryBytes(v.Interface())), "->", hex.EncodeToString(b))
		}
	}
}

func testint(value interface{}, t *testing.T, callee reflect.Value, methodName string) {
	vv := reflect.ValueOf(value)
	vkeys := vv.MapKeys()
	for _, v := range vkeys {
		length := int(vv.MapIndex(v).Field(0).Int())
		neg := vv.MapIndex(v).Field(1).Bool()

		n, b := Numeric.IntToBytes(v.Int())

		if n != neg {
			t.Error("IntToBytes neg error:", v, "should", neg, "but", n)
		}

		l := 0
		if b != nil {
			l = len(b)
		}
		if l != length {
			t.Error("IntToBytes length error:", hex.EncodeToString(ToBinaryBytes(v.Interface())), "->", hex.EncodeToString(b))
		}

		m := callee.MethodByName(methodName)
		called := m.Call([]reflect.Value{reflect.ValueOf(b), reflect.ValueOf(n)})
		if called[0].Int() != v.Int() {
			t.Error(methodName, "codec error:", hex.EncodeToString(b), "->", hex.EncodeToString(ToBinaryBytes(v.Interface())))
		} else {
			t.Log("PASS: ", v, "->", vv.MapIndex(v), hex.EncodeToString(ToBinaryBytes(v.Interface())), "->", hex.EncodeToString(b))
		}
	}
}

func TestNumeric_Uint(t *testing.T) {
	nv := reflect.ValueOf(Numeric)
	t.Log("testing uint64")
	testuint(u64, t, nv, "BytesToUint64")

	t.Log("testing uint32")
	testuint(u32, t, nv, "BytesToUint32")

	t.Log("testing uint16")
	testuint(u16, t, nv, "BytesToUint16")

	t.Log("testing uint8")
	testuint(u8, t, nv, "BytesToUint8")
}

func TestNumeric_Int(t *testing.T) {
	nv := reflect.ValueOf(Numeric)
	t.Log("testing int64")
	testint(i64, t, nv, "BytesToInt64")

	t.Log("testing int32")
	testint(i32, t, nv, "BytesToInt32")

	t.Log("testing int16")
	testint(i16, t, nv, "BytesToInt16")

	t.Log("testing int8")
	testint(i8, t, nv, "BytesToInt8")
}

func TestType(t *testing.T) {
	fmt.Println(headerTypeMap)
}

func TestNumeric_NegIntToMinBytes(t *testing.T) {
	var i int64
	i = -0x7FFFFFFFFFFFFFFF
	t.Log(i, hex.EncodeToString(Numeric.NegIntToMinBytes(i)))
	i = -0x7FFFFFFFFFFFFFF
	t.Log(i, hex.EncodeToString(Numeric.NegIntToMinBytes(i)))
	i = -0x7FFFFFFFFFFF
	t.Log(i, hex.EncodeToString(Numeric.NegIntToMinBytes(i)))
	i = -0x7FFFFFFFFF
	t.Log(i, hex.EncodeToString(Numeric.NegIntToMinBytes(i)))
	i = -0x7FFFFFFF
	t.Log(i, hex.EncodeToString(Numeric.NegIntToMinBytes(i)))
	i = -0x7FFFFF
	t.Log(i, hex.EncodeToString(Numeric.NegIntToMinBytes(i)))
	i = -0x7FFF
	t.Log(i, hex.EncodeToString(Numeric.NegIntToMinBytes(i)))
	i = -0x7F
	t.Log(i, hex.EncodeToString(Numeric.NegIntToMinBytes(i)))
	i = -0x80000000000001
	t.Log(i, hex.EncodeToString(Numeric.NegIntToMinBytes(i)))
	i = -0x800000000001
	t.Log(i, hex.EncodeToString(Numeric.NegIntToMinBytes(i)))
	i = -0x8000000001
	t.Log(i, hex.EncodeToString(Numeric.NegIntToMinBytes(i)))
	i = -0x80000001
	t.Log(i, hex.EncodeToString(Numeric.NegIntToMinBytes(i)))
	i = -0x800001
	t.Log(i, hex.EncodeToString(Numeric.NegIntToMinBytes(i)))
	i = -0x8001
	t.Log(i, hex.EncodeToString(Numeric.NegIntToMinBytes(i)))
	i = -0x81
	t.Log(i, hex.EncodeToString(Numeric.NegIntToMinBytes(i)))
}

func TestNumeric_BytesToInt(t *testing.T) {
	bb := [][]byte{{0x0}, {0xff}, {0xff, 0xff}, {0xff, 0x0, 0x0}, {0xff, 0xff, 0xff},
		{0x7f, 0x0, 0x0, 0x0}, {0xff, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0}}
	vv := []int{0, 255, 65535, 16711680, 16777215, 2130706432}
	for i := 0; i < len(bb); i++ {
		j := Numeric.BytesToInt(bb[i])
		if i < len(vv) {
			if j == vv[i] {
				t.Logf("0x%X == %d", bb[i], vv[i])
			} else {
				t.Errorf("0x%X -> %d != %d", bb[i], j, vv[i])
			}
		} else {
			t.Logf("0x%X -> %d", bb[i], j)
		}
	}
}

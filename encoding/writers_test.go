package encoding

import (
	"bytes"
	"fmt"
	"io"
	"math/big"
	"reflect"
	"testing"
)

type encodeTest struct {
	i1    int8
	i2    int16
	i3    int32
	i4    int64
	u1    uint8
	u2    uint16
	u3    uint32
	u4    uint64
	bytes [32]byte
}

func (t *encodeTest) Serialization(w io.Writer) error {
	w.Write(ToBinaryBytes(t.i1))
	w.Write(ToBinaryBytes(t.i2))
	w.Write(ToBinaryBytes(t.i3))
	w.Write(ToBinaryBytes(t.i4))
	w.Write(ToBinaryBytes(t.u1))
	w.Write(ToBinaryBytes(t.u2))
	w.Write(ToBinaryBytes(t.u3))
	w.Write(ToBinaryBytes(t.u4))
	w.Write(t.bytes[:])
	return nil
}

func (t *encodeTest) Deserialization(r io.Reader) error {
	bs := make([]byte, 32)
	r.Read(bs[:1])
	t.i1 = int8(BinaryToInt(bs[:1]))

	r.Read(bs[:2])
	t.i2 = int16(BinaryToInt(bs[:2]))

	r.Read(bs[:4])
	t.i3 = int32(BinaryToInt(bs[:4]))

	r.Read(bs[:8])
	t.i4 = int64(BinaryToInt(bs[:8]))

	r.Read(bs[:1])
	t.u1 = uint8(BinaryToUint(bs[:1]))

	r.Read(bs[:2])
	t.u2 = uint16(BinaryToUint(bs[:2]))

	r.Read(bs[:4])
	t.u3 = uint32(BinaryToUint(bs[:4]))

	r.Read(bs[:8])
	t.u4 = uint64(BinaryToUint(bs[:8]))

	r.Read(bs[0:32])
	copy(t.bytes[:], bs)

	return nil
}

type namedByteType byte

type RawValue []byte

type simplestruct struct {
	A uint
	B string
}

type recstruct struct {
	I     uint
	Child *recstruct `rlp:"nil"`
}

type tailRaw struct {
	A    uint
	Tail []RawValue `rlp:"tail"`
}

type hasIgnoredField struct {
	A uint
	B uint `rtl:"-"`
	C uint
	D uint `json:"-"`
}

type inout struct {
	val           interface{}
	output, error string
}

type mapstruct struct {
	A map[string]int64
	B map[int64]*string
}

var (
	string1 = "string1"
	string2 = "string2"
)

var encTests = []inout{

	{val: float32(111.3), output: ""},
	{val: float64(34343434.333), output: ""},

	{val: &encodeTest{i1: 0x12, i2: 0x3456, i3: 0x567890ab, i4: -1, u1: 0xF1, u2: 0xFFF2, u3: 0xFFFFFFF3,
		u4: 0xFFFFFFFFFFFFFFF4, bytes: [...]byte{'0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '0',
			'1', '2', '3', '4', '5', '6', '7', '8', '9', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9',
			'0', '1'}}, output: "FF"},

	{val: mapstruct{map[string]int64{"key1": 1, "key2": 2}, map[int64]*string{1: &string1, 2: &string2}}, output: "FF"},

	// booleans
	{val: true, output: "01"},
	{val: false, output: "80"},

	// integers
	{val: uint32(0), output: "80"},
	{val: uint32(127), output: "7F"},
	{val: uint32(128), output: "8180"},
	{val: uint32(256), output: "820100"},
	{val: uint32(1024), output: "820400"},
	{val: uint32(0xFFFFFF), output: "83FFFFFF"},
	{val: uint32(0xFFFFFFFF), output: "84FFFFFFFF"},
	{val: uint64(0xFFFFFFFF), output: "84FFFFFFFF"},
	{val: uint64(0xFFFFFFFFFF), output: "85FFFFFFFFFF"},
	{val: uint64(0xFFFFFFFFFFFF), output: "86FFFFFFFFFFFF"},
	{val: uint64(0xFFFFFFFFFFFFFF), output: "87FFFFFFFFFFFFFF"},
	{val: uint64(0xFFFFFFFFFFFFFFFF), output: "88FFFFFFFFFFFFFFFF"},

	{val: int8(0), output: "80"},
	{val: int8(127), output: "7F"},
	{val: int8(-127), output: "7F"},
	{val: int16(128), output: "8180"},
	{val: int16(-128), output: "8180"},
	{val: int32(256), output: "820100"},
	{val: int32(-1024), output: "820400"},
	{val: int32(0xFFFFFF), output: "83FFFFFF"},
	{val: int64(0xFFFFFFFF), output: "84FFFFFFFF"},
	{val: int64(-0xFFFFFFFF), output: "84FFFFFFFF"},
	{val: int64(0xFFFFFFFFFF), output: "85FFFFFFFFFF"},
	{val: int64(-0xFFFFFFFFFF), output: "85FFFFFFFFFF"},
	{val: int64(0xFFFFFFFFFFFF), output: "86FFFFFFFFFFFF"},
	{val: int64(-0xFFFFFFFFFFFF), output: "86FFFFFFFFFFFF"},
	{val: int64(0xFFFFFFFFFFFFFF), output: "87FFFFFFFFFFFFFF"},
	{val: int64(-0xFFFFFFFFFFFFFF), output: "87FFFFFFFFFFFFFF"},
	{val: int64(0x7FFFFFFFFFFFFFFF), output: "88FFFFFFFFFFFFFFFF"},
	{val: int64(-0x7FFFFFFFFFFFFFFF), output: "88FFFFFFFFFFFFFFFF"},

	// big integers (should match uint for small values)
	{val: big.NewInt(0), output: "80"},
	{val: big.NewInt(1), output: "01"},
	{val: big.NewInt(127), output: "7F"},
	{val: big.NewInt(128), output: "8180"},
	{val: big.NewInt(256), output: "820100"},
	{val: big.NewInt(1024), output: "820400"},
	{val: big.NewInt(0xFFFFFF), output: "83FFFFFF"},
	{val: big.NewInt(0xFFFFFFFF), output: "84FFFFFFFF"},
	{val: big.NewInt(0xFFFFFFFFFF), output: "85FFFFFFFFFF"},
	{val: big.NewInt(0xFFFFFFFFFFFF), output: "86FFFFFFFFFFFF"},
	{val: big.NewInt(0xFFFFFFFFFFFFFF), output: "87FFFFFFFFFFFFFF"},
	{
		val:    big.NewInt(0).SetBytes(Unhex("102030405060708090A0B0C0D0E0F2")),
		output: "8F102030405060708090A0B0C0D0E0F2",
	},
	{
		val:    big.NewInt(0).SetBytes(Unhex("0100020003000400050006000700080009000A000B000C000D000E01")),
		output: "9C0100020003000400050006000700080009000A000B000C000D000E01",
	},
	{
		val:    big.NewInt(0).SetBytes(Unhex("010000000000000000000000000000000000000000000000000000000000000000")),
		output: "A1010000000000000000000000000000000000000000000000000000000000000000",
	},
	{
		val:    big.NewInt(0).Sub(big.NewInt(0), big.NewInt(0).SetBytes(Unhex("0100020003000400050006000700080009000A000B000C000D000E01"))),
		output: "9C0100020003000400050006000700080009000A000B000C000D000E01",
	},

	// non-pointer big.Int
	{val: *big.NewInt(0), output: "80"},
	{val: *big.NewInt(0xFFFFFF), output: "83FFFFFF"},

	// negative ints are not supported
	{val: big.NewInt(-1), error: "rlp: cannot encode negative *big.Int"},

	// byte slices, strings
	{val: []byte{}, output: "80"},
	{val: []byte{0x7E}, output: "7E"},
	{val: []byte{0x7F}, output: "7F"},
	{val: []byte{0x80}, output: "8180"},
	{val: []byte{1, 2, 3}, output: "83010203"},

	{val: []namedByteType{1, 2, 3}, output: "83010203"},
	{val: [...]namedByteType{1, 2, 3}, output: "83010203"},

	{val: "", output: "80"},
	{val: "\x7E", output: "7E"},
	{val: "\x7F", output: "7F"},
	{val: "\x80", output: "8180"},
	{val: "dog", output: "83646F67"},
	{
		val:    "Lorem ipsum dolor sit amet, consectetur adipisicing eli",
		output: "B74C6F72656D20697073756D20646F6C6F722073697420616D65742C20636F6E7365637465747572206164697069736963696E6720656C69",
	},
	{
		val:    "Lorem ipsum dolor sit amet, consectetur adipisicing elit",
		output: "B8384C6F72656D20697073756D20646F6C6F722073697420616D65742C20636F6E7365637465747572206164697069736963696E6720656C6974",
	},
	{
		val:    "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Curabitur mauris magna, suscipit sed vehicula non, iaculis faucibus tortor. Proin suscipit ultricies malesuada. Duis tortor elit, dictum quis tristique eu, ultrices at risus. Morbi a est imperdiet mi ullamcorper aliquet suscipit nec lorem. Aenean quis leo mollis, vulputate elit varius, consequat enim. Nulla ultrices turpis justo, et posuere urna consectetur nec. Proin non convallis metus. Donec tempor ipsum in mauris congue sollicitudin. Vestibulum ante ipsum primis in faucibus orci luctus et ultrices posuere cubilia Curae; Suspendisse convallis sem vel massa faucibus, eget lacinia lacus tempor. Nulla quis ultricies purus. Proin auctor rhoncus nibh condimentum mollis. Aliquam consequat enim at metus luctus, a eleifend purus egestas. Curabitur at nibh metus. Nam bibendum, neque at auctor tristique, lorem libero aliquet arcu, non interdum tellus lectus sit amet eros. Cras rhoncus, metus ac ornare cursus, dolor justo ultrices metus, at ullamcorper volutpat",
		output: "B904004C6F72656D20697073756D20646F6C6F722073697420616D65742C20636F6E73656374657475722061646970697363696E6720656C69742E20437572616269747572206D6175726973206D61676E612C20737573636970697420736564207665686963756C61206E6F6E2C20696163756C697320666175636962757320746F72746F722E2050726F696E20737573636970697420756C74726963696573206D616C6573756164612E204475697320746F72746F7220656C69742C2064696374756D2071756973207472697374697175652065752C20756C7472696365732061742072697375732E204D6F72626920612065737420696D70657264696574206D6920756C6C616D636F7270657220616C6971756574207375736369706974206E6563206C6F72656D2E2041656E65616E2071756973206C656F206D6F6C6C69732C2076756C70757461746520656C6974207661726975732C20636F6E73657175617420656E696D2E204E756C6C6120756C74726963657320747572706973206A7573746F2C20657420706F73756572652075726E6120636F6E7365637465747572206E65632E2050726F696E206E6F6E20636F6E76616C6C6973206D657475732E20446F6E65632074656D706F7220697073756D20696E206D617572697320636F6E67756520736F6C6C696369747564696E2E20566573746962756C756D20616E746520697073756D207072696D697320696E206661756369627573206F726369206C756374757320657420756C74726963657320706F737565726520637562696C69612043757261653B2053757370656E646973736520636F6E76616C6C69732073656D2076656C206D617373612066617563696275732C2065676574206C6163696E6961206C616375732074656D706F722E204E756C6C61207175697320756C747269636965732070757275732E2050726F696E20617563746F722072686F6E637573206E69626820636F6E64696D656E74756D206D6F6C6C69732E20416C697175616D20636F6E73657175617420656E696D206174206D65747573206C75637475732C206120656C656966656E6420707572757320656765737461732E20437572616269747572206174206E696268206D657475732E204E616D20626962656E64756D2C206E6571756520617420617563746F72207472697374697175652C206C6F72656D206C696265726F20616C697175657420617263752C206E6F6E20696E74657264756D2074656C6C7573206C65637475732073697420616D65742065726F732E20437261732072686F6E6375732C206D65747573206163206F726E617265206375727375732C20646F6C6F72206A7573746F20756C747269636573206D657475732C20617420756C6C616D636F7270657220766F6C7574706174",
	},

	// slices
	{val: []uint{}, output: "C0"},
	{val: []uint{1, 2, 3}, output: "C3010203"},
	{
		// [ [], [[]], [ [], [[]] ] ]
		val:    []interface{}{[]interface{}{}, []interface{}{[]interface{}{}}, []interface{}{[]interface{}{}, []interface{}{[]interface{}{}}}},
		output: "C7C0C1C0C3C0C1C0",
	},
	{
		val:    []string{"aaa", "bbb", "ccc", "ddd", "eee", "fff", "ggg", "hhh", "iii", "jjj", "kkk", "lll", "mmm", "nnn", "ooo"},
		output: "F83C836161618362626283636363836464648365656583666666836767678368686883696969836A6A6A836B6B6B836C6C6C836D6D6D836E6E6E836F6F6F",
	},
	{
		val:    []interface{}{uint64(1), uint64(0xFFFFFF), []interface{}{[]interface{}{uint64(4), uint64(5), uint64(5)}}, "abc"},
		output: "CE0183FFFFFFC4C304050583616263",
	},
	{
		val: [][]string{
			{"asdf", "qwer", "zxcv"},
			{"asdf", "qwer", "zxcv"},
			{"asdf", "qwer", "zxcv"},
			{"asdf", "qwer", "zxcv"},
			{"asdf", "qwer", "zxcv"},
			{"asdf", "qwer", "zxcv"},
			{"asdf", "qwer", "zxcv"},
			{"asdf", "qwer", "zxcv"},
			{"asdf", "qwer", "zxcv"},
			{"asdf", "qwer", "zxcv"},
			{"asdf", "qwer", "zxcv"},
			{"asdf", "qwer", "zxcv"},
			{"asdf", "qwer", "zxcv"},
			{"asdf", "qwer", "zxcv"},
			{"asdf", "qwer", "zxcv"},
			{"asdf", "qwer", "zxcv"},
			{"asdf", "qwer", "zxcv"},
			{"asdf", "qwer", "zxcv"},
			{"asdf", "qwer", "zxcv"},
			{"asdf", "qwer", "zxcv"},
			{"asdf", "qwer", "zxcv"},
			{"asdf", "qwer", "zxcv"},
			{"asdf", "qwer", "zxcv"},
			{"asdf", "qwer", "zxcv"},
			{"asdf", "qwer", "zxcv"},
			{"asdf", "qwer", "zxcv"},
			{"asdf", "qwer", "zxcv"},
			{"asdf", "qwer", "zxcv"},
			{"asdf", "qwer", "zxcv"},
			{"asdf", "qwer", "zxcv"},
			{"asdf", "qwer", "zxcv"},
			{"asdf", "qwer", "zxcv"},
		},
		output: "F90200CF84617364668471776572847A786376CF84617364668471776572847A786376CF84617364668471776572847A786376CF84617364668471776572847A786376CF84617364668471776572847A786376CF84617364668471776572847A786376CF84617364668471776572847A786376CF84617364668471776572847A786376CF84617364668471776572847A786376CF84617364668471776572847A786376CF84617364668471776572847A786376CF84617364668471776572847A786376CF84617364668471776572847A786376CF84617364668471776572847A786376CF84617364668471776572847A786376CF84617364668471776572847A786376CF84617364668471776572847A786376CF84617364668471776572847A786376CF84617364668471776572847A786376CF84617364668471776572847A786376CF84617364668471776572847A786376CF84617364668471776572847A786376CF84617364668471776572847A786376CF84617364668471776572847A786376CF84617364668471776572847A786376CF84617364668471776572847A786376CF84617364668471776572847A786376CF84617364668471776572847A786376CF84617364668471776572847A786376CF84617364668471776572847A786376CF84617364668471776572847A786376CF84617364668471776572847A786376",
	},

	// RawValue
	{val: RawValue(Unhex("01")), output: "01"},
	{val: RawValue(Unhex("82FFFF")), output: "82FFFF"},
	{val: []RawValue{Unhex("01"), Unhex("02")}, output: "C20102"},

	// structs
	{val: simplestruct{}, output: "C28080"},
	{val: simplestruct{A: 3, B: "foo"}, output: "C50383666F6F"},
	{val: &recstruct{5, nil}, output: "C205C0"},
	{val: &recstruct{5, &recstruct{4, &recstruct{3, nil}}}, output: "C605C404C203C0"},
	{val: &tailRaw{A: 1, Tail: []RawValue{Unhex("02"), Unhex("03")}}, output: "C3010203"},
	{val: &tailRaw{A: 1, Tail: []RawValue{Unhex("02")}}, output: "C20102"},
	{val: &tailRaw{A: 1, Tail: []RawValue{}}, output: "C101"},
	{val: &tailRaw{A: 1, Tail: nil}, output: "C101"},
	{val: &hasIgnoredField{A: 1, B: 0, C: 3, D: 0}, output: "C20103"},

	// nil
	{val: (*uint)(nil), output: "80"},
	{val: (*string)(nil), output: "80"},
	{val: (*[]byte)(nil), output: "80"},
	{val: (*[10]byte)(nil), output: "80"},
	{val: (*big.Int)(nil), output: "80"},
	{val: (*[]string)(nil), output: "C0"},
	{val: (*[10]string)(nil), output: "C0"},
	{val: (*[]interface{})(nil), output: "C0"},
	{val: (*[]struct{ uint })(nil), output: "C0"},
	{val: (*interface{})(nil), output: "C0"},

	// interfaces
	// {val: []io.Reader{reader}, output: "C3C20102"}, // the contained value is a struct
	//
	// Encoder
	// {val: (*testEncoder)(nil), output: "00000000"},
	// {val: &testEncoder{}, output: "00010001000100010001"},
	// {val: &testEncoder{errors.New("test error")}, error: "test error"},
	// // verify that pointer method testEncoder.EncodeRLP is called for
	// // addressable non-pointer values.
	// {val: &struct{ TE testEncoder }{testEncoder{}}, output: "CA00010001000100010001"},
	// {val: &struct{ TE testEncoder }{testEncoder{errors.New("test error")}}, error: "test error"},
	// // verify the error for non-addressable non-pointer Encoder
	// {val: testEncoder{}, error: "rlp: game over: unadressable value of type rlp.testEncoder, EncodeRLP is pointer method"},
	// // verify the special case for []byte
	// {val: []byteEncoder{0, 1, 2, 3, 4}, output: "C5C0C0C0C0C0"},
}

func TestEncode(t *testing.T) {
	buf := new(bytes.Buffer)
	for _, test := range encTests {
		val := reflect.ValueOf(test.val)
		buf.Reset()
		// valueWriter(buf, val)
		Encode(test.val, buf)
		bs := buf.Bytes()
		// fmt.Println(test.val, "->", hex.EncodeToString(buf.Bytes()))

		typ := val.Type()
		nv := reflect.New(typ)
		// vr := NewValueReader(buf, 100)
		nvv := nv.Elem()
		// if err := valueReader(vr, nvv); err != nil {
		// vr := NewValueReader(buf, 256)
		vr := buf
		if err := Decode(vr, nv.Interface()); err != nil {
			t.Error(err)
		}

		fmt.Printf("%v: %#v\n\t%X\n%v: %#v\n", typ, test.val, bs, nvv.Type(), nvv)
		if reflect.DeepEqual(test.val, nvv.Interface()) {
			t.Log(test.val, "check")
		} else {
			t.Error(test.val, "error")
		}
	}
}

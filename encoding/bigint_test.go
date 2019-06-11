package encoding

import (
	"bytes"
	"encoding/hex"
	"math/big"
	"testing"
)

var (
	bigintarray = []*big.Int{
		big.NewInt(0),
		big.NewInt(8767398),
		big.NewInt(0).SetBytes(Unhex("203040506070809010a0b0c0d0e0f0")),
		big.NewInt(0).SetBytes(Unhex("9000000000a0000000000000000011110010101010101010101010")),
		big.NewInt(-8765421),
		big.NewInt(0).Sub(big.NewInt(0), big.NewInt(0).SetBytes(Unhex("378777256599386728cb371923"))),
	}
)

func TestBigIntCodec(t *testing.T) {
	buf := new(bytes.Buffer)
	for _, i := range bigintarray {
		buf.Reset()
		err := EncodeBigInt(i, buf)
		if err != nil {
			t.Error(err)
			continue
		}
		b := buf.Bytes()
		t.Log(hex.EncodeToString(b))
		ni := new(big.Int)
		err = DecodeBigInt(buf, ni)
		if err != nil {
			t.Error(err)
		}
		if i.Cmp(ni) != 0 {
			t.Error(ni, "should be", i)
		} else {
			t.Log(i, "check")
		}
	}
}

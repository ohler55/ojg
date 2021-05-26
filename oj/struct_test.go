// Copyright (c) 2020, Peter Ohler, All rights reserved.

package oj_test

import (
	"fmt"
	"math/big"
	"strings"
	"testing"

	"github.com/ohler55/ojg"
	"github.com/ohler55/ojg/oj"
	"github.com/ohler55/ojg/tt"
)

func TestJSONTagPrimitive(t *testing.T) {
	type Sample struct {
		Yes bool    `json:"yes"`
		No  bool    `json:"no"`
		I   int     `json:"a"`
		I8  int8    `json:"a8"`
		I16 int16   `json:"a16"`
		I32 int32   `json:"a32"`
		I64 int64   `json:"a64"`
		U   uint    `json:"b"`
		U8  uint8   `json:"b8"`
		U16 uint16  `json:"b16"`
		U32 uint32  `json:"b32"`
		U64 uint64  `json:"b64"`
		F32 float32 `json:"f32"`
		F64 float64 `json:"f64"`
		Str string  `json:"z"`
	}
	sample := Sample{
		Yes: true,
		No:  false,
		I:   1,
		I8:  2,
		I16: 3,
		I32: 4,
		I64: 5,
		U:   6,
		U8:  7,
		U16: 8,
		U32: 9,
		U64: 10,
		F32: 11.5,
		F64: 12.5,
		Str: "abc",
	}
	wr := oj.Writer{Options: ojg.Options{UseTags: true}}

	out := wr.MustJSON(&sample)
	tt.Equal(t,
		`{"a":1,"a16":3,"a32":4,"a64":5,"a8":2,"b":6,"b16":8,"b32":9,"b64":10,"b8":7,"f32":11.5,"f64":12.5,"no":false,"yes":true,"z":"abc"}`,
		string(out))
	out = wr.MustJSON(sample)
	tt.Equal(t,
		`{"a":1,"a16":3,"a32":4,"a64":5,"a8":2,"b":6,"b16":8,"b32":9,"b64":10,"b8":7,"f32":11.5,"f64":12.5,"no":false,"yes":true,"z":"abc"}`,
		string(out))

	wr.UseTags = false
	out = wr.MustJSON(&sample)
	tt.Equal(t,
		`{"f32":11.5,"f64":12.5,"i":1,"i16":3,"i32":4,"i64":5,"i8":2,"no":false,"str":"abc","u":6,"u16":8,"u32":9,"u64":10,"u8":7,"yes":true}`,
		string(out))
	out = wr.MustJSON(sample)
	tt.Equal(t,
		`{"f32":11.5,"f64":12.5,"i":1,"i16":3,"i32":4,"i64":5,"i8":2,"no":false,"str":"abc","u":6,"u16":8,"u32":9,"u64":10,"u8":7,"yes":true}`,
		string(out))

	wr.KeyExact = true
	out = wr.MustJSON(&sample)
	tt.Equal(t,
		`{"F32":11.5,"F64":12.5,"I":1,"I16":3,"I32":4,"I64":5,"I8":2,"No":false,"Str":"abc","U":6,"U16":8,"U32":9,"U64":10,"U8":7,"Yes":true}`,
		string(out))
	out = wr.MustJSON(sample)
	tt.Equal(t,
		`{"F32":11.5,"F64":12.5,"I":1,"I16":3,"I32":4,"I64":5,"I8":2,"No":false,"Str":"abc","U":6,"U16":8,"U32":9,"U64":10,"U8":7,"Yes":true}`,
		string(out))
}

func TestJSONTagAsString(t *testing.T) {
	type Sample struct {
		Yes bool    `json:"yes,string"`
		No  bool    `json:"no,string"`
		I   int     `json:"a,string"`
		I8  int8    `json:"a8,string"`
		I16 int16   `json:"a16,string"`
		I32 int32   `json:"a32,string"`
		I64 int64   `json:"a64,string"`
		U   uint    `json:"b,string"`
		U8  uint8   `json:"b8,string"`
		U16 uint16  `json:"b16,string"`
		U32 uint32  `json:"b32,string"`
		U64 uint64  `json:"b64,string"`
		F32 float32 `json:"f32,string"`
		F64 float64 `json:"f64,string"`
		Str string  `json:"z,string"`
	}
	sample := Sample{
		Yes: true,
		No:  false,
		I:   1,
		I8:  2,
		I16: 3,
		I32: 4,
		I64: 5,
		U:   6,
		U8:  7,
		U16: 8,
		U32: 9,
		U64: 10,
		F32: 11.5,
		F64: 12.5,
		Str: "abc",
	}
	wr := oj.Writer{Options: ojg.Options{UseTags: true}}

	out := wr.MustJSON(&sample)
	tt.Equal(t,
		`{"a":"1","a16":"3","a32":"4","a64":"5","a8":"2","b":"6","b16":"8","b32":"9","b64":"10","b8":"7","f32":"11.5","f64":"12.5","no":"false","yes":"true","z":"abc"}`,
		string(out))
	out = wr.MustJSON(sample)
	tt.Equal(t,
		`{"a":"1","a16":"3","a32":"4","a64":"5","a8":"2","b":"6","b16":"8","b32":"9","b64":"10","b8":"7","f32":"11.5","f64":"12.5","no":"false","yes":"true","z":"abc"}`,
		string(out))
}

func TestJSONTagOmitEmpty(t *testing.T) {
	type Sample struct {
		Yes bool    `json:"yes,omitempty"`
		No  bool    `json:"no,omitempty"`
		I   int     `json:"a,omitempty"`
		I8  int8    `json:"a8,omitempty"`
		I16 int16   `json:"a16,omitempty"`
		I32 int32   `json:"a32,omitempty"`
		I64 int64   `json:"a64,omitempty"`
		U   uint    `json:"b,omitempty"`
		U8  uint8   `json:"b8,omitempty"`
		U16 uint16  `json:"b16,omitempty"`
		U32 uint32  `json:"b32,omitempty"`
		U64 uint64  `json:"b64,omitempty"`
		F32 float32 `json:"f32,omitempty"`
		F64 float64 `json:"f64,omitempty"`
		Str string  `json:"z,omitempty"`
	}
	sample := Sample{
		Yes: true,
		No:  false,
		I:   1,
		I8:  2,
		I16: 3,
		I32: 4,
		I64: 5,
		U:   6,
		U8:  7,
		U16: 8,
		U32: 9,
		U64: 10,
		F32: 11.5,
		F64: 12.5,
		Str: "abc",
	}
	wr := oj.Writer{Options: ojg.Options{UseTags: true}}

	out := wr.MustJSON(&sample)
	tt.Equal(t,
		`{"a":1,"a16":3,"a32":4,"a64":5,"a8":2,"b":6,"b16":8,"b32":9,"b64":10,"b8":7,"f32":11.5,"f64":12.5,"yes":true,"z":"abc"}`,
		string(out))

	out = wr.MustJSON(sample)
	tt.Equal(t,
		`{"a":1,"a16":3,"a32":4,"a64":5,"a8":2,"b":6,"b16":8,"b32":9,"b64":10,"b8":7,"f32":11.5,"f64":12.5,"yes":true,"z":"abc"}`,
		string(out))

	out = wr.MustJSON(&Sample{})
	tt.Equal(t, "{}", string(out))

	out = wr.MustJSON(Sample{})
	tt.Equal(t, "{}", string(out))
}

func TestJSONTagOmitEmptyAsString(t *testing.T) {
	type Sample struct {
		Yes bool    `json:"yes,omitempty,string"`
		No  bool    `json:"no,omitempty,string"`
		I   int     `json:"a,omitempty,string"`
		I8  int8    `json:"a8,omitempty,string"`
		I16 int16   `json:"a16,omitempty,string"`
		I32 int32   `json:"a32,omitempty,string"`
		I64 int64   `json:"a64,omitempty,string"`
		U   uint    `json:"b,omitempty,string"`
		U8  uint8   `json:"b8,omitempty,string"`
		U16 uint16  `json:"b16,omitempty,string"`
		U32 uint32  `json:"b32,omitempty,string"`
		U64 uint64  `json:"b64,omitempty,string"`
		F32 float32 `json:"f32,omitempty,string"`
		F64 float64 `json:"f64,omitempty,string"`
		Str string  `json:"z,omitempty,string"`
	}
	sample := Sample{
		Yes: true,
		No:  false,
		I:   1,
		I8:  2,
		I16: 3,
		I32: 4,
		I64: 5,
		U:   6,
		U8:  7,
		U16: 8,
		U32: 9,
		U64: 10,
		F32: 11.5,
		F64: 12.5,
		Str: "abc",
	}
	wr := oj.Writer{Options: ojg.Options{UseTags: true}}

	out := wr.MustJSON(&sample)
	tt.Equal(t,
		`{"a":"1","a16":"3","a32":"4","a64":"5","a8":"2","b":"6","b16":"8","b32":"9","b64":"10","b8":"7","f32":"11.5","f64":"12.5","yes":"true","z":"abc"}`,
		string(out))
	out = wr.MustJSON(sample)
	tt.Equal(t,
		`{"a":"1","a16":"3","a32":"4","a64":"5","a8":"2","b":"6","b16":"8","b32":"9","b64":"10","b8":"7","f32":"11.5","f64":"12.5","yes":"true","z":"abc"}`,
		string(out))

	out = wr.MustJSON(&Sample{})
	tt.Equal(t, "{}", string(out))

	out = wr.MustJSON(Sample{})
	tt.Equal(t, "{}", string(out))
}

func TestJSONTagPtrOmitEmpty(t *testing.T) {
	type Bare struct {
	}
	type Sample struct {
		Ptr    *Bare         `json:"p,omitempty"`
		NilPtr *Bare         `json:"np,omitempty"`
		Slice  []interface{} `json:"s,omitempty"`
		Empty  []interface{} `json:"e,omitempty"`
		Any    interface{}   `json:"a,omitempty"`
		NilAny interface{}   `json:"na,omitempty"`
		Bar    **Bare        `json:"bar"`
	}
	sample := Sample{
		Ptr:    &Bare{},
		NilPtr: nil,
		Slice:  []interface{}{true},
		Empty:  []interface{}{},
		Any:    &Bare{},
		NilAny: nil,
	}
	wr := oj.Writer{Options: ojg.Options{UseTags: true}}

	out := wr.MustJSON(&sample)
	tt.Equal(t, `{"a":{},"bar":null,"p":{},"s":[true]}`, string(out))

	out = wr.MustJSON(sample)
	tt.Equal(t, `{"a":{},"bar":null,"p":{},"s":[true]}`, string(out))

	wr.Indent = 2
	out = wr.MustJSON(&sample)
	tt.Equal(t, `{
  "a": {},
  "bar": null,
  "p": {},
  "s": [
    true
  ]
}`, string(out))
}

func TestJSONTagPtr(t *testing.T) {
	type Bare struct {
	}
	type Sample struct {
		Ptr    *Bare         `json:"p"`
		NilPtr *Bare         `json:"np"`
		Slice  []interface{} `json:"s"`
		Empty  []interface{} `json:"e"`
		Any    interface{}   `json:"a"`
		NilAny interface{}   `json:"na"`
	}
	sample := Sample{
		Ptr:    &Bare{},
		NilPtr: nil,
		Slice:  []interface{}{true},
		Empty:  []interface{}{},
		Any:    &Bare{},
		NilAny: nil,
	}
	wr := oj.Writer{Options: ojg.Options{UseTags: true}}

	out := wr.MustJSON(&sample)
	tt.Equal(t, `{"a":{},"e":[],"na":null,"np":null,"p":{},"s":[true]}`, string(out))

	out = wr.MustJSON(sample)
	tt.Equal(t, `{"a":{},"e":[],"na":null,"np":null,"p":{},"s":[true]}`, string(out))

	wr.OmitNil = true
	wr.UseTags = false
	out = wr.MustJSON(&sample)
	tt.Equal(t, `{"any":{},"empty":[],"ptr":{},"slice":[true]}`, string(out))
}

func TestJSONTagOther(t *testing.T) {
	type Sample struct {
		AsIs int         `json:",omitempty"`
		Dash int         `json:"-,"`
		Skip int         `json:"-"`
		Nil  interface{} `json:"nil"`
		x    int
	}
	sample := Sample{
		AsIs: 1,
		Dash: 2,
		Skip: 3,
		x:    4,
	}
	wr := oj.Writer{Options: ojg.Options{UseTags: true, OmitNil: true, Indent: 2}}

	out := wr.MustJSON(&sample)
	tt.Equal(t, `{
  "-": 2,
  "AsIs": 1
}`, string(out))
}

type Decimal struct {
	value *big.Int
	exp   int32
	fail  bool
}

func (d Decimal) MarshalJSON() ([]byte, error) {
	if d.fail {
		return nil, fmt.Errorf("don't like this one")
	}
	return []byte(fmt.Sprintf(`"%d,%d"`, d.value, d.exp)), nil
}

type TestStruct struct {
	Outer   bool    `json:"outer"`
	Decimal Decimal `json:"decimal"`
}

func TestMarshalStructMarshaler(t *testing.T) {
	tsa := []TestStruct{{
		Outer:   true,
		Decimal: Decimal{value: big.NewInt(5), exp: 2},
	}}
	out, err := oj.Marshal(tsa)
	tt.Nil(t, err)
	tt.Equal(t, true, strings.Contains(string(out), `"decimal":"5,2"`))

	tsa = []TestStruct{{
		Outer:   true,
		Decimal: Decimal{value: big.NewInt(5), exp: 2, fail: true},
	}}
	_, err = oj.Marshal(tsa)
	tt.NotNil(t, err)
}

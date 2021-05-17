// Copyright (c) 2021, Peter Ohler, All rights reserved.

package sen_test

import (
	"fmt"
	"strings"

	"github.com/ohler55/ojg"
	"github.com/ohler55/ojg/pretty"
	"github.com/ohler55/ojg/sen"
)

func ExampleParse() {
	val, err := sen.Parse([]byte("[true false [3 2 1] {a:1 b:2 c:3 d:[x y z []]}]"))
	if err != nil {
		panic(err)
	}
	fmt.Println(pretty.SEN(val, 80.3))

	// Output:
	// [
	//   true
	//   false
	//   [3 2 1]
	//   {a: 1 b: 2 c: 3 d: [x y z []]}
	// ]
}

func ExampleMustParse() {
	val := sen.MustParse([]byte("[true false [3 2 1] {a:1 b:2 c:3 d:[x y z []]}]"))
	fmt.Println(pretty.SEN(val, 80.3))

	// Output:
	// [
	//   true
	//   false
	//   [3 2 1]
	//   {a: 1 b: 2 c: 3 d: [x y z []]}
	// ]
}

func ExampleParseReader() {
	r := strings.NewReader("[true false [3 2 1] {a:1 b:2 c:3 d:[x y z []]}]")
	val, err := sen.ParseReader(r)
	if err != nil {
		panic(err)
	}
	fmt.Println(pretty.SEN(val, 80.3))

	// Output:
	// [
	//   true
	//   false
	//   [3 2 1]
	//   {a: 1 b: 2 c: 3 d: [x y z []]}
	// ]
}

func ExampleMustParseReader() {
	r := strings.NewReader("[true false [3 2 1] {a:1 b:2 c:3 d:[x y z []]}]")
	val := sen.MustParseReader(r)
	fmt.Println(pretty.SEN(val, 80.3))

	// Output:
	// [
	//   true
	//   false
	//   [3 2 1]
	//   {a: 1 b: 2 c: 3 d: [x y z []]}
	// ]
}

func ExampleUnmarshal() {
	type Sample struct {
		X int
		Y string
	}
	var sample Sample
	if err := sen.Unmarshal([]byte("{x: 3 y: why}"), &sample); err != nil {
		panic(err)
	}
	fmt.Printf("sample.X: %d  sample.Y: %s\n", sample.X, sample.Y)

	// Output: sample.X: 3  sample.Y: why
}

func ExampleString() {
	type Sample struct {
		X int
		Y string
	}
	s := sen.String(&Sample{X: 3, Y: "why"})
	fmt.Println(s)

	// Output: {x:3 y:why}
}

func ExampleBytes() {
	type Sample struct {
		X int
		Y string
	}
	b := sen.Bytes(&Sample{X: 3, Y: "why"})
	fmt.Println(string(b))

	// Output: {x:3 y:why}
}

func ExampleWrite() {
	type Sample struct {
		X int
		Y string
	}
	var buf strings.Builder
	if err := sen.Write(&buf, &Sample{X: 3, Y: "why"}); err != nil {
		panic(err)
	}
	fmt.Println(buf.String())

	// Output: {x:3 y:why}
}

func ExampleMustWrite() {
	type Sample struct {
		X int
		Y string
	}
	var buf strings.Builder
	sen.MustWrite(&buf, &Sample{X: 3, Y: "why"})
	fmt.Println(buf.String())

	// Output: {x:3 y:why}
}

func ExampleParser_Parse() {
	p := sen.Parser{}
	// An invalid JSON but valid SEN.
	simple, err := p.Parse([]byte(`{abc: [{"x": {"y": [{b: true}]} z: 7}]}`))
	if err != nil {
		panic(err)
	}
	fmt.Println(sen.String(simple, &ojg.Options{Sort: true}))

	// Output: {abc:[{x:{y:[{b:true}]} z:7}]}
}

func ExampleParser_MustParse() {
	p := sen.Parser{}
	// An invalid JSON but valid SEN.
	simple := p.MustParse([]byte(`{abc: [{"x": {"y": [{b: true}]} z: 7}]}`))
	fmt.Println(sen.String(simple, &ojg.Options{Sort: true}))

	// Output: {abc:[{x:{y:[{b:true}]} z:7}]}
}

func ExampleParser_ParseReader() {
	p := sen.Parser{}
	// An invalid JSON but valid SEN.
	r := strings.NewReader(`{abc: [{"x": {"y": [{b: true}]} z: 7}]}`)
	simple, err := p.ParseReader(r)
	if err != nil {
		panic(err)
	}
	fmt.Println(sen.String(simple, &ojg.Options{Sort: true}))

	// Output: {abc:[{x:{y:[{b:true}]} z:7}]}
}

func ExampleParser_MustParseReader() {
	p := sen.Parser{}
	// An invalid JSON but valid SEN.
	r := strings.NewReader(`{abc: [{"x": {"y": [{b: true}]} z: 7}]}`)
	simple := p.MustParseReader(r)
	fmt.Println(sen.String(simple, &ojg.Options{Sort: true}))

	// Output: {abc:[{x:{y:[{b:true}]} z:7}]}
}

func ExampleParser_Unmarshal() {
	type Sample struct {
		X int
		Y string
	}
	p := sen.Parser{}
	var sample Sample
	if err := p.Unmarshal([]byte("{x: 3 y: why}"), &sample); err != nil {
		panic(err)
	}
	fmt.Printf("sample.X: %d  sample.Y: %s\n", sample.X, sample.Y)

	// Output: sample.X: 3  sample.Y: why
}

func ExampleWriter_SEN() {
	type Sample struct {
		X int
		Y string
	}
	wr := sen.Writer{}
	s := wr.SEN(&Sample{X: 3, Y: "why"})
	fmt.Println(s)

	// Output: {x:3 y:why}
}

func ExampleWriter_MustSEN() {
	type Sample struct {
		X int
		Y string
	}
	wr := sen.Writer{}
	b := wr.MustSEN(&Sample{X: 3, Y: "why"})
	fmt.Println(string(b))

	// Output: {x:3 y:why}
}

func ExampleWriter_Write() {
	type Sample struct {
		X int
		Y string
	}
	wr := sen.Writer{}
	var buf strings.Builder
	if err := wr.Write(&buf, &Sample{X: 3, Y: "why"}); err != nil {
		panic(err)
	}
	fmt.Println(buf.String())

	// Output: {x:3 y:why}
}

func ExampleWriter_MustWrite() {
	type Sample struct {
		X int
		Y string
	}
	wr := sen.Writer{}
	var buf strings.Builder
	wr.MustWrite(&buf, &Sample{X: 3, Y: "why"})
	fmt.Println(buf.String())

	// Output: {x:3 y:why}
}

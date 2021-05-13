// Copyright (c) 2020, Peter Ohler, All rights reserved.

package alt_test

import (
	"fmt"
	"time"

	"github.com/ohler55/ojg"
	"github.com/ohler55/ojg/alt"
	"github.com/ohler55/ojg/gen"
	"github.com/ohler55/ojg/oj"
	"github.com/ohler55/ojg/sen"
)

func ExampleDecompose() {
	type Sample struct {
		Int int
		Str string
	}
	sample := Sample{Int: 3, Str: "three"}
	// Decompose and add a CreateKey to indicate the type with a full path.
	simple := alt.Decompose(&sample, &ojg.Options{CreateKey: "^", FullTypePath: true})

	fmt.Println(oj.JSON(simple, &oj.Options{Sort: true}))

	// Output: {"^":"github.com/ohler55/ojg/alt_test/Sample","int":3,"str":"three"}
}

func ExampleDup() {
	type Sample struct {
		Int int
		Str string
	}
	sample := []interface{}{&Sample{Int: 3, Str: "three"}, 42}
	// Dup creates a deep duplicate of a simple type and decomposes any
	// structs according to the optional options just like alt.Decompose does.
	simple := alt.Decompose(sample, &ojg.Options{CreateKey: "^"})

	fmt.Println(oj.JSON(simple, &ojg.Options{Sort: true}))

	// Output: [{"^":"Sample","int":3,"str":"three"},42]
}

func ExampleRecomposer_Recompose() {
	type Sample struct {
		Int int
		Str string
	}
	// Recomposers are reuseable. Create one and use the default reflect composer (nil).
	r, err := alt.NewRecomposer("^", map[interface{}]alt.RecomposeFunc{&Sample{}: nil})
	if err != nil {
		panic(err)
	}
	var v interface{}
	// Recompose without providing a struct to populate.
	v, err = r.Recompose(map[string]interface{}{"^": "Sample", "int": 3, "str": "three"})
	if err != nil {
		panic(err)
	}
	fmt.Printf("type: %T\n", v)
	if sample, _ := v.(*Sample); sample != nil {
		fmt.Printf("sample: {Int: %d, Str: %q}\n", sample.Int, sample.Str)
	}
	// Output:
	// type: *alt_test.Sample
	// sample: {Int: 3, Str: "three"}
}

func ExampleNewRecomposer() {
	type Sample struct {
		Int int
		Str string
	}
	// Recomposers are reuseable. Create one and use the default reflect composer (nil).
	r, err := alt.NewRecomposer("^", map[interface{}]alt.RecomposeFunc{&Sample{}: nil})
	if err != nil {
		panic(err)
	}
	var v interface{}
	// Recompose without providing a struct to populate.
	v, err = r.Recompose(map[string]interface{}{"^": "Sample", "int": 3, "str": "three"})
	if err != nil {
		panic(err)
	}
	fmt.Printf("type: %T\n", v)
	if sample, _ := v.(*Sample); sample != nil {
		fmt.Printf("sample: {Int: %d, Str: %q}\n", sample.Int, sample.Str)
	}
	// Output:
	// type: *alt_test.Sample
	// sample: {Int: 3, Str: "three"}
}

func ExampleRecompose() {
	type Sample struct {
		Int int
		Str string
	}
	// Simplified sample data or JSON as a map[string]interface{}.
	data := map[string]interface{}{"int": 3, "str": "three"}
	var sample Sample
	// Recompose into the sample struct. Panic on failure.
	v, err := alt.Recompose(data, &sample)
	if err != nil {
		panic(err)
	}
	fmt.Printf("type: %T\n", v)
	fmt.Printf("sample: {Int: %d, Str: %q}\n", sample.Int, sample.Str)

	// Output:
	// type: *alt_test.Sample
	// sample: {Int: 3, Str: "three"}
}

func ExampleMustRecompose() {
	type Sample struct {
		Int int
		Str string
	}
	// Simplified sample data or JSON as a map[string]interface{}.
	data := map[string]interface{}{"int": 3, "str": "three"}
	var sample Sample
	// Recompose into the sample struct. Panic on failure.
	v := alt.MustRecompose(data, &sample)

	fmt.Printf("type: %T\n", v)
	fmt.Printf("sample: {Int: %d, Str: %q}\n", sample.Int, sample.Str)

	// Output:
	// type: *alt_test.Sample
	// sample: {Int: 3, Str: "three"}
}

func ExampleRecomposer_MustRecompose() {
	type Sample struct {
		Int  int
		When time.Time
	}
	// Create a new Recomposer that uses "^" as the create key and register a
	// default reflection recompose function (nil). A time recomposer from an
	// integer is also included in the new recomposer compser options.
	r := alt.MustNewRecomposer("^",
		map[interface{}]alt.RecomposeFunc{&Sample{}: nil},
		map[interface{}]alt.RecomposeAnyFunc{&time.Time{}: func(v interface{}) (interface{}, error) {
			if s, _ := v.(string); 0 < len(s) {
				return time.ParseInLocation(time.RFC3339, s, time.UTC)
			}
			return nil, fmt.Errorf("can not convert a %v to a time.Time", v)
		}})
	// Simplified sample data or JSON as a map[string]interface{} with an
	// included create key using "^" to avoid possible conflicts with other
	// fields in the struct.
	data := map[string]interface{}{"^": "Sample", "int": 3, "when": "2021-02-09T01:02:03Z"}
	v := r.MustRecompose(data)

	if sample, _ := v.(*Sample); sample != nil {
		fmt.Printf("sample: {Int: %d, When: %q}\n", sample.Int, sample.When.Format(time.RFC3339))
	}
	// Output:
	// sample: {Int: 3, When: "2021-02-09T01:02:03Z"}
}

func ExampleMustNewRecomposer() {
	type Sample struct {
		Int  int
		When time.Time
	}
	// Create a new Recomposer that uses "^" as the create key and register a
	// default reflection recompose function (nil). A time recomposer from an
	// integer is also included in the new recomposer compser options.
	r := alt.MustNewRecomposer("^",
		map[interface{}]alt.RecomposeFunc{&Sample{}: nil},
		map[interface{}]alt.RecomposeAnyFunc{&time.Time{}: func(v interface{}) (interface{}, error) {
			if s, _ := v.(string); 0 < len(s) {
				return time.ParseInLocation(time.RFC3339, s, time.UTC)
			}
			return nil, fmt.Errorf("can not convert a %v to a time.Time", v)
		}})
	// Simplified sample data or JSON as a map[string]interface{} with an
	// included create key using "^" to avoid possible conflicts with other
	// fields in the struct.
	data := map[string]interface{}{"^": "Sample", "int": 3, "when": "2021-02-09T01:02:03Z"}
	v := r.MustRecompose(data)

	if sample, _ := v.(*Sample); sample != nil {
		fmt.Printf("sample: {Int: %d, When: %q}\n", sample.Int, sample.When.Format(time.RFC3339))
	}
	// Output:
	// sample: {Int: 3, When: "2021-02-09T01:02:03Z"}
}

func ExampleAlter() {
	src := map[string]interface{}{"a": 1, "b": 4, "c": 9}
	// Alter the src as needed avoiding duplicating when possible.
	val := alt.Alter(src)
	// Modify src should change val since they are the same map.
	src["d"] = 16
	fmt.Println(sen.String(val, &oj.Options{Sort: true}))

	// Output: {a:1 b:4 c:9 d:16}
}

func ExampleGenAlter() {
	m := map[string]interface{}{"a": 1, "b": 4, "c": 9}
	// Convert to a gen.Node.
	node := alt.GenAlter(m)
	fmt.Println(sen.String(node, &oj.Options{Sort: true}))
	obj, _ := node.(gen.Object)
	fmt.Printf("member type: %T\n", obj["b"])

	// Output: {a:1 b:4 c:9}
	// member type: gen.Int
}

func ExampleRecomposer_RegisterComposer() {
	type Sample struct {
		Int  int
		When time.Time
	}
	r := alt.MustNewRecomposer("^", nil)
	err := r.RegisterComposer(&Sample{}, nil)
	if err != nil {
		panic(err)
	}
	err = r.RegisterAnyComposer(time.Time{},
		func(v interface{}) (interface{}, error) {
			if secs, ok := v.(int); ok {
				return time.Unix(int64(secs), 0), nil
			}
			return nil, fmt.Errorf("can not convert a %T to a time.Time", v)
		})
	if err != nil {
		panic(err)
	}
	data := map[string]interface{}{"^": "Sample", "int": 3, "when": 1612872722}
	sample, _ := r.MustRecompose(data).(*Sample)

	fmt.Printf("sample.Int: %d\n", sample.Int)
	fmt.Printf("sample.When: %d\n", sample.When.Unix())

	// Output:
	// sample.Int: 3
	// sample.When: 1612872722
}

func ExampleRecomposer_RegisterAnyComposer() {
	type Sample struct {
		Int  int
		When time.Time
	}
	r := alt.MustNewRecomposer("^", nil)
	err := r.RegisterComposer(&Sample{}, nil)
	if err != nil {
		panic(err)
	}
	err = r.RegisterAnyComposer(time.Time{},
		func(v interface{}) (interface{}, error) {
			if secs, ok := v.(int); ok {
				return time.Unix(int64(secs), 0), nil
			}
			return nil, fmt.Errorf("can not convert a %T to a time.Time", v)
		})
	if err != nil {
		panic(err)
	}
	data := map[string]interface{}{"^": "Sample", "int": 3, "when": 1612872722}
	sample, _ := r.MustRecompose(data).(*Sample)

	fmt.Printf("sample.Int: %d\n", sample.Int)
	fmt.Printf("sample.When: %d\n", sample.When.Unix())

	// Output:
	// sample.Int: 3
	// sample.When: 1612872722
}

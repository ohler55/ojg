// Copyright (c) 2020, Peter Ohler, All rights reserved.

package alt_test

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/ohler55/ojg/alt"
	"github.com/ohler55/ojg/gen"
	"github.com/ohler55/ojg/oj"
)

func ExampleDecompose() {
	type Sample struct {
		Int int
		Str string
	}
	sample := Sample{Int: 3, Str: "three"}
	simple := alt.Decompose(&sample, &alt.Options{CreateKey: "^", FullTypePath: true})

	fmt.Println(oj.JSON(simple, &oj.Options{Sort: true}))

	// Output: {"^":"github.com/ohler55/ojg/alt_test/Sample","int":3,"str":"three"}
}

func ExampleRecomposer_Recompose() {
	type Sample struct {
		Int int
		Str string
	}
	r, err := alt.NewRecomposer("^", map[interface{}]alt.RecomposeFunc{&Sample{}: nil})
	var v interface{}
	if err == nil {
		v, err = r.Recompose(map[string]interface{}{"^": "Sample", "int": 3, "str": "three"})
	}
	if err == nil {
		fmt.Printf("type: %T\n", v)
		if sample, _ := v.(*Sample); sample != nil {
			fmt.Printf("sample: {Int: %d, Str: %q}\n", sample.Int, sample.Str)
		}
	} else {
		fmt.Println(err.Error())
	}
	// Output:
	// type: *alt_test.Sample
	// sample: {Int: 3, Str: "three"}
}

func ExampleMustRecompose() {
	type Sample struct {
		Int int
		Str string
	}
	// Create a Recomposer that

	// Simplified sample data or JSON as a map[string]interface{}.
	data := map[string]interface{}{"int": 3, "str": "three"}
	var sample Sample
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
			if secs, ok := v.(int); ok {
				return time.Unix(int64(secs), 0), nil
			}
			return nil, fmt.Errorf("can not convert a %T to a time.Time", v)
		}})
	// Simplified sample data or JSON as a map[string]interface{} with an
	// included create key using "^" to avoid possible conflicts with other
	// fields in the struct.
	data := map[string]interface{}{"^": "Sample", "int": 3, "when": 1612872722}
	v := r.MustRecompose(data)

	if sample, _ := v.(*Sample); sample != nil {
		fmt.Printf("sample: {Int: %d, When: %q}\n", sample.Int, sample.When.Format(time.RFC3339))
	}
	// Output:
	// sample: {Int: 3, When: "2021-02-09T07:12:02-05:00"}
}

type Genny struct {
	val int
}

func (g *Genny) Generic() gen.Node {
	return gen.Object{"type": gen.String("Genny"), "val": gen.Int(g.val)}
}

func ExampleGenerify() {
	// type Genny struct {
	// 	val int
	// }
	//
	// func (g *Genny) Generic() gen.Node {
	// 	return gen.Object{"type": gen.String("genny"), "val": gen.Int(g.val)}
	// }
	ga := []*Genny{&Genny{val: 3}}
	v := alt.Generify(ga)
	fmt.Println(oj.JSON(v, &oj.Options{Sort: true}))

	// Output: [{"type":"Genny","val":3}]
}

func ExampleAlter() {
	m := map[string]interface{}{"a": 1, "b": 4, "c": 9}
	v := alt.GenAlter(m)
	fmt.Println(oj.JSON(v, &oj.Options{Sort: true}))

	// Output: {"a":1,"b":4,"c":9}
}

type Animal interface {
	Kind() string
}

type Dog struct {
	Size string
}

func (d *Dog) Kind() string {
	return fmt.Sprintf("%s dog", d.Size)
}

type Cat struct {
	Color string
}

func (c *Cat) Kind() string {
	return fmt.Sprintf("%s cat", c.Color)
}

func ExampleDecompose_animal() {
	pets := []Animal{&Dog{Size: "big"}, &Cat{Color: "black"}}

	// First marshal using the go json package.
	pj, err := json.Marshal(pets)
	if err != nil {
		fmt.Printf("error: %s\n", err)
	}
	// Works just fine.
	fmt.Printf("json.Marshal: %s\n", pj)

	// Now try to unmarshall. An error is returned with a list of nils.
	var petsOut []Animal
	err = json.Unmarshal(pj, &petsOut)
	fmt.Printf("error: %s\n", err)
	fmt.Printf("jsom.Unmarshal: %v\n", petsOut)

	// Now try OjG. Decompress and create a JSON []byte slice.
	simple := alt.Decompose(pets, &alt.Options{CreateKey: "^"})
	// Sort the object members in the output for repeatability.
	ps := oj.JSON(simple, &oj.Options{Sort: true})
	fmt.Printf("oj.JSON: %s\n", ps)

	// Create a new Recomposer. This can be use over and over again. Register
	// the types with a nil creation function to let reflection do the work
	// since the styles are exported.
	var r *alt.Recomposer
	if r, err = alt.NewRecomposer("^", map[interface{}]alt.RecomposeFunc{&Dog{}: nil, &Cat{}: nil}); err != nil {
		fmt.Printf("error: %s\n", err)
	}

	// Recompose from the simplified data earlier. The one that matches the JSON.
	var result interface{}
	if result, err = r.Recompose(simple, []Animal{}); err != nil {
		fmt.Printf("error: %s\n", err)
	}
	// Check the results.
	// members.
	pets, _ = result.([]Animal)
	for _, animal := range pets {
		fmt.Printf("  %T - %s\n", animal, animal.Kind())
	}
	// Output:
	// json.Marshal: [{"Size":"big"},{"Color":"black"}]
	// error: json: cannot unmarshal object into Go value of type alt_test.Animal
	// jsom.Unmarshal: [<nil> <nil>]
	// oj.JSON: [{"^":"Dog","size":"big"},{"^":"Cat","color":"black"}]
	//   *alt_test.Dog - big dog
	//   *alt_test.Cat - black cat
}

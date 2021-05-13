// Copyright (c) 2021, Peter Ohler, All rights reserved.

package alt_test

import (
	"fmt"

	"github.com/ohler55/ojg/alt"
	"github.com/ohler55/ojg/oj"
)

// Encode and decode slice of interfaces. Similar behavior is available with
// oj.Unmarshal and sen.Unmarshal.

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

func ExampleRecomposer_Recompose_animals() {
	pets := []Animal{&Dog{Size: "big"}, &Cat{Color: "black"}}

	// Decompose and use a create key to identify the encoded type.
	simple := alt.Decompose(pets, &alt.Options{CreateKey: "^"})
	// Sort the object members in the output for repeatability.
	fmt.Printf("as JSON: %s\n", oj.JSON(simple, &oj.Options{Sort: true}))

	// Create a new Recomposer. This can be use over and over again. Register
	// the types with a nil creation function to let reflection do the work
	// since the types are exported.
	r, err := alt.NewRecomposer("^", map[interface{}]alt.RecomposeFunc{&Dog{}: nil, &Cat{}: nil})
	if err != nil {
		panic(err)
	}
	// Recompose from the simplified data without providing a target which
	// returns a []interface{} populated with the correct types.
	var result interface{}
	if result, err = r.Recompose(simple); err != nil {
		panic(err)
	}
	list, _ := result.([]interface{})
	for _, item := range list {
		animal, _ := item.(Animal)
		fmt.Printf("  %s\n", animal.Kind())
	}
	// Recompose with a target.
	var animals []Animal
	if _, err = r.Recompose(simple, &animals); err != nil {
		panic(err)
	}
	fmt.Println("Recompose into a target struct")
	for _, animal := range animals {
		fmt.Printf("  %T - %s\n", animal, animal.Kind())
	}
	// Output:
	// as JSON: [{"^":"Dog","size":"big"},{"^":"Cat","color":"black"}]
	//   big dog
	//   black cat
	// Recompose into a target struct
	//   *alt_test.Dog - big dog
	//   *alt_test.Cat - black cat
}

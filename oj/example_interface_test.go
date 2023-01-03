// Copyright (c) 2021, Peter Ohler, All rights reserved.

package oj_test

import (
	"fmt"

	"github.com/ohler55/ojg"
	"github.com/ohler55/ojg/alt"
	"github.com/ohler55/ojg/oj"
)

// Encode and decode slice of interfaces.

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

func ExampleUnmarshal_interface() {
	pets := []Animal{&Dog{Size: "big"}, &Cat{Color: "black"}}

	// Encode and use a create key to identify the encoded type.
	b, err := oj.Marshal(pets, &ojg.Options{CreateKey: "^", Sort: true})
	if err != nil {
		panic(err)
	}
	// Sort the object members in the output for repeatability.
	fmt.Printf("as JSON: %s\n", b)

	// Create a new Recomposer. This can be use over and over again. Register
	// the types with a nil creation function to let reflection do the work
	// since the types are exported.
	var r *alt.Recomposer
	if r, err = alt.NewRecomposer("^", map[any]alt.RecomposeFunc{&Dog{}: nil, &Cat{}: nil}); err != nil {
		panic(err)
	}
	var result any
	if err = oj.Unmarshal(b, &result, r); err != nil {
		panic(err)
	}
	list, _ := result.([]any)
	for _, item := range list {
		animal, _ := item.(Animal)
		fmt.Printf("  %s\n", animal.Kind())
	}
	// Unmarshal with a typed target.
	var animals []Animal
	if err = oj.Unmarshal(b, &animals, r); err != nil {
		panic(err)
	}
	fmt.Println("Unmarshal into a target struct")
	for _, animal := range animals {
		fmt.Printf("  %T - %s\n", animal, animal.Kind())
	}

	// Output:
	// as JSON: [{"^":"Dog","size":"big"},{"^":"Cat","color":"black"}]
	//   big dog
	//   black cat
	// Unmarshal into a target struct
	//   *oj_test.Dog - big dog
	//   *oj_test.Cat - black cat
}

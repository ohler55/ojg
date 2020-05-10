// Copyright (c) 2020, Peter Ohler, All rights reserved.

package gen

func Recompose(
	data interface{},
	createKey string,
	builders map[interface{}]func(Node) (interface{}, error)) (interface{}, error) {

	// TBD create builder map that is indexed by reflected names and path+name
	//   hang on to reflection info

	// TBD depth first walk, at leaves, if a map then look for createKey
	//  if key then lookup in builders
	//  if func found then call it
	//  if nil (found but nil) the use reflections

	return nil, nil
}

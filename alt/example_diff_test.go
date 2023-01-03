// Copyright (c) 2021, Peter Ohler, All rights reserved.

package alt_test

import (
	"fmt"
	"sort"
	"strings"

	"github.com/ohler55/ojg/alt"
)

func ExampleDiff() {
	diffs := alt.Diff(
		map[string]any{"x": 1, "y": 2, "z": []any{1, 2, 3}},
		map[string]any{"x": 1, "y": 4, "z": []any{1, 3, 5}},
	)
	sort.Slice(diffs, func(i, j int) bool {
		return 0 < strings.Compare(fmt.Sprintf("%v", diffs[j]), fmt.Sprintf("%v", diffs[i]))
	})
	fmt.Printf("diff: %v\n", diffs)

	// Output: diff: [[y] [z 1] [z 2]]
}

func ExampleCompare() {
	diff := alt.Compare(
		map[string]any{"x": 1, "y": 2, "z": []any{1, 2, 3}},
		map[string]any{"x": 1, "y": 2, "z": []any{1, 3, 5}},
	)
	fmt.Printf("diff: %v\n", diff)

	// Output: diff: [z 1]
}

func ExampleMatch() {
	fingerprint := map[string]any{"x": 1, "z": 3}
	match := alt.Match(
		fingerprint,
		map[string]any{"x": 1, "y": 2, "z": 3},
	)
	fmt.Printf("match: %t\n", match)

	// Output: match: true
}

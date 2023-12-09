// Copyright (c) 2020, Peter Ohler, All rights reserved.

package jp

import "github.com/ohler55/ojg/gen"

// Descent is used as a flag to indicate the path should be displayed in a
// recursive descent representation.
type Descent byte

// Append a fragment string representation of the fragment to the buffer
// then returning the expanded buffer.
func (f Descent) Append(buf []byte, bracket, first bool) []byte {
	if bracket {
		buf = append(buf, "[..]"...)
	} else {
		buf = append(buf, '.')
	}
	return buf
}

func (f Descent) locate(pp Expr, data any, rest Expr, max int) (locs []Expr) {
	switch td := data.(type) {
	case map[string]any:
		if len(rest) == 0 { // last one
			for k := range td {
				locs = locateAppendFrag(locs, pp, Child(k))
				if 0 < max && max <= len(locs) {
					break
				}
			}
		} else {
			// Depth first with rest[1:].
			cp := append(pp, nil) // place holder
			mx := max
			for k, v := range td {
				cp[len(pp)] = Child(k)
				locs = locateContinueFrag(locs, cp, v, rest, max)
				if 0 < max && max <= len(locs) {
					break
				}
				if 0 < max {
					mx = max - len(locs)
				}
				locs = append(locs, f.locate(cp, v, rest, mx)...)
			}
		}
	case []any:
		if len(rest) == 0 { // last one
			for i := range td {
				locs = locateAppendFrag(locs, pp, Nth(i))
				if 0 < max && max <= len(locs) {
					break
				}
			}
		} else {
			cp := append(pp, nil) // place holder
			mx := max
			for i, v := range td {
				cp[len(pp)] = Nth(i)
				locs = locateContinueFrag(locs, cp, v, rest, max)
				if 0 < max && max <= len(locs) {
					break
				}
				if 0 < max {
					mx = max - len(locs)
				}
				locs = append(locs, f.locate(cp, v, rest, mx)...)
			}
		}
	case gen.Object:
		if len(rest) == 0 { // last one
			for k := range td {
				locs = locateAppendFrag(locs, pp, Child(k))
				if 0 < max && max <= len(locs) {
					break
				}
			}
		} else {
			// Depth first with rest[1:].
			cp := append(pp, nil) // place holder
			r2 := rest[1:]
			mx := max
			for k, v := range td {
				cp[len(pp)] = Child(k)
				locs = append(locs, rest[0].locate(cp, v, r2, mx)...)
				if 0 < max {
					if max <= len(locs) {
						break
					}
					mx = max - len(locs)
				}
				locs = append(locs, f.locate(cp, v, rest, mx)...)
			}
		}
	case gen.Array:
		if len(rest) == 0 { // last one
			for i := range td {
				locs = locateAppendFrag(locs, pp, Nth(i))
				if 0 < max && max <= len(locs) {
					break
				}
			}
		} else {
			cp := append(pp, nil) // place holder
			mx := max
			for i, v := range td {
				cp[len(pp)] = Nth(i)
				locs = locateContinueFrag(locs, cp, v, rest, max)
				if 0 < max && max <= len(locs) {
					break
				}
				if 0 < max {
					mx = max - len(locs)
				}
				locs = append(locs, f.locate(cp, v, rest, mx)...)
			}
		}
	case Keyed:
		keys := td.Keys()
		if len(rest) == 0 { // last one
			for _, k := range keys {
				locs = locateAppendFrag(locs, pp, Child(k))
				if 0 < max && max <= len(locs) {
					break
				}
			}
		} else {
			cp := append(pp, nil) // place holder
			mx := max
			for _, k := range keys {
				v, _ := td.ValueForKey(k)
				cp[len(pp)] = Child(k)
				locs = locateContinueFrag(locs, cp, v, rest, max)
				if 0 < max && max <= len(locs) {
					break
				}
				if 0 < max {
					mx = max - len(locs)
				}
				locs = append(locs, f.locate(cp, v, rest, mx)...)
			}
		}
	case Indexed:
		size := td.Size()
		if len(rest) == 0 { // last one
			for i := 0; i < size; i++ {
				locs = locateAppendFrag(locs, pp, Nth(i))
				if 0 < max && max <= len(locs) {
					break
				}
			}
		} else {
			cp := append(pp, nil) // place holder
			mx := max
			for i := 0; i < size; i++ {
				v := td.ValueAtIndex(i)
				cp[len(pp)] = Nth(i)
				locs = locateContinueFrag(locs, cp, v, rest, max)
				if 0 < max && max <= len(locs) {
					break
				}
				if 0 < max {
					mx = max - len(locs)
				}
				locs = append(locs, f.locate(cp, v, rest, mx)...)
			}
		}
	default:
		// TBD
	}
	return
}

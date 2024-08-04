// Copyright (c) 2024, Peter Ohler, All rights reserved.

package jp

// PathMatch returns true if the provided path would match the target
// expression. The path argument is expected to be a normalized path with only
// elements of Root ($), At (@), Child (string), or Nth (int). A Filter
// fragment in the target expression will match any value in path since it
// requires data from a JSON document to be evaluated.
func PathMatch(target, path Expr) bool {
	if 0 < len(path) {
		switch path[0].(type) {
		case Root, At:
			path = path[1:]
		}
	}
	for i, f := range target {
		if len(path) == 0 {
			return false
		}
		switch path[0].(type) {
		case Child, Nth:
		default:
			return false
		}
		switch tf := f.(type) {
		case Root, At:
			if 0 < i { // $ and @ can only be the first fragment
				return false
			}
		case Child, Nth:
			if tf != path[0] {
				return false
			}
			path = path[1:]
		case Bracket:
			// ignore and don't advance path
		case Wildcard:
			path = path[1:]
		case Union:
			var ok bool
			for _, u := range tf {
			check:
				switch tu := u.(type) {
				case string:
					if Child(tu) == path[0] {
						ok = true
						break check
					}
				case int64:
					if Nth(tu) == path[0] {
						ok = true
						break check
					}
				}
			}
			if !ok {
				return false
			}
			path = path[1:]
		case Slice:
			// TBD
		case Filter:
			// Assume a match since there is no data for comparison.
			path = path[1:]
		case Descent:
			// TBD look for match on next target or if no next then return true
		}
	}
	return true
}

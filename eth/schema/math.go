package schema

import "golang.org/x/exp/constraints"

func Abs[T constraints.Integer](a T) T {
	if a < 0 {
		return -a
	}

	return a
}

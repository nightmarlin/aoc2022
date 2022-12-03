package lib

// Reduce performs a sequential functional reduction on the input slice.
func Reduce[S any, Out any](
	slice []S,
	reducer func(prev Out, next S) Out,
	init Out,
) Out {
	o := init

	for i := range slice {
		o = reducer(o, slice[i])
	}

	return o
}

// Map performs a sequential functional map on the input slice.
func Map[T any, U any](slice []T, mapper func(T) U) []U {
	out := make([]U, len(slice))
	for i := range slice {
		out[i] = mapper(slice[i])
	}
	return out
}

func Filter[T any](slice []T, filter func(T) bool) []T {
	res := make([]T, 0)
	for i := range slice {
		if filter(slice[i]) {
			res = append(res, slice[i])
		}
	}
	return res
}

func Any[T any](slice []T, comp func(T) bool) bool {
	for i := range slice {
		if comp(slice[i]) {
			return true
		}
	}
	return false
}

func Keys[K comparable, V any](m map[K]V) []K {
	res := make([]K, 0, len(m))
	for k := range m {
		res = append(res, k)
	}
	return res
}

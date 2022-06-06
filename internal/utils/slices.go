package utils

func SliceIncludes[T comparable](data []T, want T) (found bool) {
	for i := range data {
		if data[i] == want {
			found = true
			break
		}
	}

	return
}

func sliceEmpty[T any](data []T) bool {
	return data == nil || len(data) == 0
}

func SliceFirst[T any](data []T) (first T) {
	if sliceEmpty(data) {
		return
	}

	first = data[0]

	return
}

func SliceLast[T any](data []T) (last T) {
	if sliceEmpty(data) {
		return
	}

	last = data[len(data)-1]

	return
}

// AppendSlices is _not_ an efficient function by any means.
func AppendSlices[T any](slices ...[]T) []T {
	var newSlice = make([]T, 0)

	for i, s := range slices {
		if i == 0 {
			newSlice = s
		} else {
			newSlice2 := make([]T, len(s)+len(newSlice))

			copy(newSlice2, newSlice)
			copy(newSlice2[len(newSlice):], s)

			newSlice = newSlice2
		}
	}

	return newSlice
}

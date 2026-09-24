package changelog

// groupBy groups the elements of collection by the key iteratee returns for each.
func groupBy[T any, U comparable, Slice ~[]T](collection Slice, iteratee func(item T) U) map[U]Slice {
	result := map[U]Slice{}

	for i := range collection {
		key := iteratee(collection[i])

		result[key] = append(result[key], collection[i])
	}

	return result
}

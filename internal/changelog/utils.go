package changelog

import (
	"gitlab.com/l0nax/changelog-go/internal/config"
)

// groupBy returns an object composed of keys generated from the results of running each element of collection through iteratee.
func groupBy[T any, U comparable, Slice ~[]T](collection Slice, iteratee func(item T) U) map[U]Slice {
	result := map[U]Slice{}

	for i := range collection {
		key := iteratee(collection[i])

		result[key] = append(result[key], collection[i])
	}

	return result
}

func resolveChangeTypeID(id string) (config.ChangeType, bool) {
	for _, entry := range config.C.Entry.Types {
		if entry.ID == id {
			return entry, true
		}
	}

	return config.ChangeType{}, false
}

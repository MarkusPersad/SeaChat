package utils

import (
	"slices"
	"strings"
)

func GetSortedID(seq string ,id ...string,) string {
	slices.Sort(id)
	if len(id) == 1 {
		return id[0]
	}
	return strings.Join(id,seq)
}
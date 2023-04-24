package helpers

import "strconv"

func UintToString(n uint) string {
	return strconv.FormatUint(uint64(n), 10)
}

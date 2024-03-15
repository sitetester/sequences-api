package api

import "strconv"

func StrToUint(IDStr string) (uint64, error) {
	return strconv.ParseUint(IDStr, 10, 64)

}

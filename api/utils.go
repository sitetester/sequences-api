package api

import "strconv"

func StrToUint(IDStr string) (uint64, error) {
	return strconv.ParseUint(IDStr, 10, 64)

}

func UintToStr(id uint) string {
	return strconv.FormatUint(uint64(id), 10)
}

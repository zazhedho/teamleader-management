package utils

import (
	"strconv"
	"strings"
)

func ParseInt(val string) (int, error) {
	clean := strings.ReplaceAll(val, ",", "")
	num, err := strconv.Atoi(clean)
	if err != nil {
		return 0, err
	}
	return num, nil
}

package stringsx

import "strconv"

func ParserDefaultInt(numStr string, defaultInt int) int {
	if num, err := strconv.Atoi(numStr); err == nil {
		return num
	}
	return defaultInt
}

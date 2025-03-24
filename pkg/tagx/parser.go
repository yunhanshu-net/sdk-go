package tagx

import "strings"

func ParserKv(tag string) map[string]string {
	mp := make(map[string]string)
	split := strings.Split(tag, ";")
	for _, s := range split {
		vals := strings.Split(s, ":")
		key := vals[0]
		value := vals[1]
		mp[key] = value
	}
	return mp
}

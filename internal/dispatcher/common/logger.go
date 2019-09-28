package common

// RequestResponseHeadersToString
func RequestResponseHeadersToString(headers map[string][]string) string {
	var out string
	for k, v := range headers {
		out += k + ":" + v[0] + "\n "
	}
	return out
}

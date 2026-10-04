package httpgeneric

import "net/url"

func templateURLQuery(s string) string { return url.QueryEscape(s) }

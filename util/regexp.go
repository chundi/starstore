package util

import "regexp"

const (
	REGEXP_EMAIL  string = "^[a-zA-Z0-9_.+-]+@[a-zA-Z0-9-]+\\.[a-zA-Z0-9-.]+$"
	REGEXP_MOBILE string = "^\\+?(\\d{1,4})?[- \\+_]?\\d+$"
)

var (
	REG_EMAIL  *regexp.Regexp
	REG_MOBILE *regexp.Regexp
)

func init() {
	REG_EMAIL = regexp.MustCompile(REGEXP_EMAIL)
	REG_MOBILE = regexp.MustCompile(REGEXP_MOBILE)
}

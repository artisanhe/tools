package validatetpl

import (
	"fmt"
	"regexp"
)

const InvalidNumString = "非法数字字符串"

var numStringRegexp = regexp.MustCompile(`^\d*$`)

func NewValidateNumStringLength(min_len int, max_len int) func(v interface{}) (bool, string) {
	return func(v interface{}) (bool, string) {
		if stringer, ok := v.(fmt.Stringer); ok {
			v = stringer.String()
		}

		if value, ok := v.(string); ok {
			if !numStringRegexp.MatchString(value) {
				return false, InvalidNumString
			}

			str_len := len(value)
			if str_len < min_len || (str_len > max_len && max_len != STRING_UNLIMIT_VALUE) {
				return false, fmt.Sprintf(STRING_LENGHT_NOT_IN_RANGE, min_len, max_len, str_len)
			}
			return true, ""
		}

		return false, InvalidNumString
	}
}

package utils

import "strings"

func Trim(strs ...*string) {
	for _, s := range strs {
		if s != nil {
			*s = strings.TrimSpace(*s)
		}
	}
}

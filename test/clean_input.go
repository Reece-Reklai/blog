package test

import (
	"strings"
	"unicode"
)

func CleanInput(test string) []string {
	for _, value := range test {
		if unicode.IsDigit(value) == true {
			return nil
		}
	}
	convert := []rune(test)
	var slice []string
	var tmp []rune
	if test == "" {
		return nil
	}
	if len(test) == 1 {
		slice = append(slice, strings.ToLower(test))
		return slice
	}
	for index, value := range convert {
		if index == 0 {
			if value != 32 {
				tmp = append(tmp, value)
			}
			continue
		}
		if value == 32 {
			if tmp != nil {
				conversion := string(tmp)
				slice = append(slice, strings.ToLower(conversion))
				tmp = nil
			}
			continue
		}
		if index == len(test)-1 {
			if tmp != nil {
				tmp = append(tmp, value)
				conversion := string(tmp)
				slice = append(slice, strings.ToLower(conversion))
				tmp = nil
			}
		}
		tmp = append(tmp, value)
	}
	return slice
}

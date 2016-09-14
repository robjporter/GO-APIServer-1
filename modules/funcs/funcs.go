package funcs

import (
	"github.com/aymerick/raymond"
)

var (
	Bold = func(content string) raymond.SafeString {
		return raymond.SafeString("<b>" + content + "</b>")
	}
	Italic = func(content string) raymond.SafeString {
		return raymond.SafeString("<i>" + content + "</i>")
	}
	Underline = func(content string) raymond.SafeString {
		return raymond.SafeString("<u>" + content + "</u>")
	}
)

func RemoveDuplicates(a []string) []string { 
	result := []string{} 
	seen := map[string]string{} 
	for _, val := range a { 
		if _, ok := seen[val]; !ok { 
			result = append(result, val) 
			seen[val] = val 
		} 
	} 
	return result 
} 
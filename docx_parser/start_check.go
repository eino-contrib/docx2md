package docx_parser

import (
	"regexp"
)

// Check if the string starts with an Arabic number
func startsWithArabicNumber(s string) bool {
	if len(s) == 0 {
		return false
	}
	return s[0] >= '0' && s[0] <= '9'
}

// Check if the string starts with "第" + an Arabic number
func startsWithDiAndArabicNumber(s string) bool {
	matched, _ := regexp.MatchString(`^第[0-9]+`, s)
	return matched
}

// Check if the string starts with a Chinese number
func startsWithChineseNumber(s string) bool {
	if len(s) == 0 {
		return false
	}
	matched, _ := regexp.MatchString(`[一二三四五六七八九十]`, s)
	return matched
}

// Check if the string starts with "第" + a Chinese number
func startsWithDiAndChineseNumber(s string) bool {
	matched, _ := regexp.MatchString(`^第[一二三四五六七八九十]`, s)
	return matched
}

// Comprehensive check function
func CheckString(s string) bool {
	return startsWithArabicNumber(s) ||
		startsWithDiAndArabicNumber(s) ||
		startsWithChineseNumber(s) ||
		startsWithDiAndChineseNumber(s)
}

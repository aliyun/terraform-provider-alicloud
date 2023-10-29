// This file is auto-generated, don't edit it. Thanks.
/**
 * This is a string module
 */
package client

import (
	"strings"

	"github.com/alibabacloud-go/tea/tea"
)

func Split(raw *string, sep *string, limit *int) (_result []*string) {
	if limit == nil || tea.IntValue(limit) == -1 || tea.IntValue(limit) == 0 {
		return tea.StringSlice(strings.Split(tea.StringValue(raw), tea.StringValue(sep)))
	}
	return tea.StringSlice(strings.SplitN(tea.StringValue(raw), tea.StringValue(sep), tea.IntValue(limit)))
}

func Replace(raw *string, oldStr *string, newStr *string, count *int) (_result *string) {
	if count == nil {
		return tea.String(strings.ReplaceAll(tea.StringValue(raw), tea.StringValue(oldStr), tea.StringValue(newStr)))
	}
	return tea.String(strings.Replace(tea.StringValue(raw), tea.StringValue(oldStr), tea.StringValue(newStr), tea.IntValue(count)))
}

func Contains(s *string, substr *string) (_result *bool) {
	return tea.Bool(strings.Contains(tea.StringValue(s), tea.StringValue(substr)))
}

func Count(s *string, substr *string) (_result *int) {
	return tea.Int(strings.Count(tea.StringValue(s), tea.StringValue(substr)))
}

func HasPrefix(s *string, prefix *string) (_result *bool) {
	return tea.Bool(strings.HasPrefix(tea.StringValue(s), tea.StringValue(prefix)))
}

func HasSuffix(s *string, suffix *string) (_result *bool) {
	return tea.Bool(strings.HasSuffix(tea.StringValue(s), tea.StringValue(suffix)))
}

func Index(s *string, substr *string) (_result *int) {
	return tea.Int(strings.Index(tea.StringValue(s), tea.StringValue(substr)))
}

func ToLower(s *string) (_result *string) {
	return tea.String(strings.ToLower(tea.StringValue(s)))
}

func ToUpper(s *string) (_result *string) {
	return tea.String(strings.ToUpper(tea.StringValue(s)))
}

func SubString(s *string, start, end *int) (_result *string) {
	endIndex := tea.IntValue(end)
	startIndex := tea.IntValue(start)
	str := tea.StringValue(s)
	if endIndex == -1 {
		return tea.String(str[startIndex : len(str)-1])
	}
	return tea.String(str[startIndex:endIndex])
}

func Equals(expect *string, actual *string) (_result *bool) {
	return tea.Bool(tea.StringValue(expect) == tea.StringValue(actual))
}

func ToBytes(s *string, sep *string) []byte {
	return []byte(tea.StringValue(s))
}

func Trim(s *string) *string {
	return tea.String(strings.Trim(tea.StringValue(s), " "))
}

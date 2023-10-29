// This file is auto-generated, don't edit it. Thanks.
/**
 * This is a map module
 */
package client

import "github.com/alibabacloud-go/tea/tea"

func Size(raw map[string]interface{}) (_result *int) {
	var count int
	for range raw {
		count = count + 1
	}
	return tea.Int(count)
}

func KeySet(raw interface{}) (_result []*string) {
	slice := tea.StringSliceValue(_result)
	switch v := raw.(type) {
	case map[string]interface{}:
		for k := range v {
			slice = append(slice, k)
		}
	case map[string]string:
		for k := range v {
			slice = append(slice, k)
		}
	case map[string]*string:
		for k := range v {
			slice = append(slice, k)
		}
	}
	return tea.StringSlice(slice)
}

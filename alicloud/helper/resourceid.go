package helper

import (
	"fmt"
	"strings"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
)

// Resource-ID parsing helpers moved from the root package
// (alicloud/common.go); package alicloud re-exports them as var aliases
// so the flat layout keeps compiling unchanged.

func ParseResourceIds(id string) (parts []string, err error) {
	parts = strings.Split(id, ":")
	return parts, err
}

func ParseResourceId(id string, length int) (parts []string, err error) {
	parts = strings.Split(id, ":")

	if len(parts) != length {
		err = errs.WrapError(fmt.Errorf("Invalid Resource Id %s. Expected parts' length %d, got %d", id, length, len(parts)))
	}
	return parts, err
}

func ParseResourceIdN(id string, length int) (parts []string, err error) {
	parts = strings.SplitN(id, ":", length)

	if len(parts) != length {
		err = errs.WrapError(fmt.Errorf("Invalid Resource Id %s. Expected parts' length %d, got %d", id, length, len(parts)))
	}
	return parts, err
}

func ParseResourceIdWithEscaped(id string, length int) (parts []string, err error) {
	parts = make([]string, 0)
	var currentPart strings.Builder
	for i := 0; i < len(id); i++ {
		if id[i] == '\\' {
			i++
			if i < len(id) {
				currentPart.WriteByte(id[i])
			}
		} else if id[i] == ':' {
			parts = append(parts, currentPart.String())
			currentPart.Reset()
		} else {
			currentPart.WriteByte(id[i])
		}
	}

	if currentPart.Len() > 0 {
		parts = append(parts, currentPart.String())
	}

	if len(parts) != length {
		err = errs.WrapError(fmt.Errorf("Invalid Resource Id %s. Expected parts' length %d, got %d", id, length, len(parts)))
	}
	return parts, err
}

func EscapeColons(s string) string {
	return strings.ReplaceAll(s, ":", "\\:")
}

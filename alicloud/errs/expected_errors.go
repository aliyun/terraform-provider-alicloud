package errs

import (
	"fmt"
	"strings"

	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/errors"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/denverdino/aliyungo/common"
)

// IsExpectedErrors reports whether err matches one of the expected error
// codes. It mirrors the root alicloud.IsExpectedErrors helper so that service
// subpackages (which cannot import the root provider package without a cycle)
// can classify expected errors through the errs package they already import.
// The type-aware branches cover the error kinds this package already depends
// on (ComplexError, tea.SDKError, sdkerrors.ServerError, ProviderError,
// common.Error, oss.ServiceError); anything else falls through to a string
// containment check on err.Error(), matching the root fallback.
func IsExpectedErrors(err error, expectCodes []string) bool {
	if err == nil {
		return false
	}

	if e, ok := err.(*ComplexError); ok {
		return IsExpectedErrors(e.Cause, expectCodes)
	}

	if e, ok := err.(*tea.SDKError); ok {
		for _, code := range expectCodes {
			// The second statement aims to match the tea sdk history bug.
			if *e.Code == code || strings.HasPrefix(code, *e.Code) || strings.Contains(*e.Data, code) {
				return true
			}
		}
		return false
	}

	if e, ok := err.(*errors.ServerError); ok {
		for _, code := range expectCodes {
			if e.ErrorCode() == code || strings.Contains(e.Message(), code) {
				return true
			}
		}
		return false
	}

	if e, ok := err.(*ProviderError); ok {
		for _, code := range expectCodes {
			if e.ErrorCode() == code || strings.Contains(e.Message(), code) {
				return true
			}
		}
		return false
	}

	if e, ok := err.(*common.Error); ok {
		for _, code := range expectCodes {
			if e.Code == code || fmt.Sprint(e.StatusCode) == code || strings.Contains(e.Message, code) {
				return true
			}
		}
		return false
	}

	if e, ok := err.(oss.ServiceError); ok {
		for _, code := range expectCodes {
			if e.Code == code || strings.Contains(e.Message, code) {
				return true
			}
		}
		return false
	}

	for _, code := range expectCodes {
		if strings.Contains(err.Error(), code) {
			return true
		}
	}
	return false
}

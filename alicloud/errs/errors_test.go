package errs_test

import (
	"fmt"
	"testing"

	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
)

func teaSDKError(status int, code, message, data string) error {
	return &tea.SDKError{
		StatusCode: tea.Int(status),
		Code:       tea.String(code),
		Message:    tea.String(message),
		Data:       tea.String(data),
	}
}

func TestUnitIsTypedNotFound(t *testing.T) {
	cases := map[string]struct {
		err  error
		want bool
	}{
		"nil":                    {err: nil, want: false},
		"transport404":           {err: teaSDKError(404, "InvalidParameter", "not found", ""), want: true},
		"bodyCode404":            {err: teaSDKError(200, "404", "not found", ""), want: true},
		"forbiddenWith404InData": {err: teaSDKError(403, "APINotAuthorize", "forbidden", "previous response 404"), want: false},
		"serverError":            {err: teaSDKError(500, "InternalError", "boom", ""), want: false},
		"nonTeaError":            {err: fmt.Errorf("404"), want: false},
		"wrappedTransport404":    {err: fmt.Errorf("calling ARMS: %w", teaSDKError(404, "InvalidParameter", "not found", "")), want: true},
		"wrappedBodyCode404":     {err: fmt.Errorf("calling ARMS: %w", teaSDKError(200, "404", "not found", "")), want: true},
		"wrappedForbidden404":    {err: fmt.Errorf("calling ARMS: %w", teaSDKError(403, "APINotAuthorize", "forbidden", "previous response 404")), want: false},
		"wrappedNonTeaError":     {err: fmt.Errorf("calling ARMS: %w", fmt.Errorf("connection reset")), want: false},
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			if got := errs.IsTypedNotFound(c.err); got != c.want {
				t.Fatalf("IsTypedNotFound() = %v, want %v", got, c.want)
			}
		})
	}
}

func TestUnitIsResponseCodeNotFound(t *testing.T) {
	cases := map[string]struct {
		err  error
		want bool
	}{
		"nil":                    {err: nil, want: false},
		"successBody404":         {err: teaSDKError(200, "404", "not found", ""), want: true},
		"transport404":           {err: teaSDKError(404, "404", "not found", ""), want: false},
		"forbiddenWith404InData": {err: teaSDKError(403, "APINotAuthorize", "forbidden", "previous response 404"), want: false},
		"nonTeaError":            {err: fmt.Errorf("404"), want: false},
		"wrappedSuccessBody404":  {err: fmt.Errorf("calling ARMS: %w", teaSDKError(200, "404", "not found", "")), want: true},
		"wrappedTransport404":    {err: fmt.Errorf("calling ARMS: %w", teaSDKError(404, "404", "not found", "")), want: false},
		"wrappedForbidden404":    {err: fmt.Errorf("calling ARMS: %w", teaSDKError(403, "APINotAuthorize", "forbidden", "previous response 404")), want: false},
		"wrappedNonTeaError":     {err: fmt.Errorf("calling ARMS: %w", fmt.Errorf("connection reset")), want: false},
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			if got := errs.IsResponseCodeNotFound(c.err); got != c.want {
				t.Fatalf("IsResponseCodeNotFound() = %v, want %v", got, c.want)
			}
		})
	}
}

func TestUnitIsProviderNotFound(t *testing.T) {
	providerNotFound := errs.WrapErrorf(errs.NotFoundErr("NotificationPolicy", "1234"), errs.NotFoundMsg, map[string]interface{}{"RequestId": "test"})
	non404WithNotFoundMessage := errs.WrapErrorf(
		teaSDKError(403, "APINotAuthorize", "NotFound while resolving authorization", ""),
		errs.DefaultErrorMsg, "1234", "ListNotificationPolicies", "SDK")
	cases := map[string]struct {
		err  error
		want bool
	}{
		"nil":                                {err: nil, want: false},
		"providerNotFoundWrapper":            {err: providerNotFound, want: true},
		"outerWrapperAroundProviderNotFound": {err: errs.WrapError(providerNotFound), want: true},
		"non404WithNotFoundMessage":          {err: non404WithNotFoundMessage, want: false},
		"rawTyped404":                        {err: teaSDKError(404, "InvalidParameter", "not found", ""), want: false},
		"plainError":                         {err: fmt.Errorf("NotFound"), want: false},
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			if got := errs.IsProviderNotFound(c.err); got != c.want {
				t.Fatalf("IsProviderNotFound() = %v, want %v", got, c.want)
			}
		})
	}
}

func TestUnitIsExpectedErrors(t *testing.T) {
	cases := map[string]struct {
		err         error
		expectCodes []string
		want        bool
	}{
		"nil":                               {err: nil, expectCodes: []string{"404"}, want: false},
		"emptyCodeEmptyDataPermission":      {err: teaSDKError(403, "", "Permission denied", ""), expectCodes: []string{"404"}, want: false},
		"emptyCodeEmptyDataMultipleExpects": {err: teaSDKError(403, "", "Permission denied", ""), expectCodes: []string{"404", "ServiceUnavailable"}, want: false},
		"nonMatchingNonEmptyCode":           {err: teaSDKError(403, "APINotAuthorize", "forbidden", ""), expectCodes: []string{"404"}, want: false},
		"exactCodeMatch":                    {err: teaSDKError(404, "404", "not found", ""), expectCodes: []string{"404"}, want: true},
		"exactCodeAmongMultiple":            {err: teaSDKError(500, "InternalError", "boom", ""), expectCodes: []string{"ServiceUnavailable", "InternalError"}, want: true},
		"reversePrefixHistoricalCompat":     {err: teaSDKError(404, "Invalid", "not found", ""), expectCodes: []string{"InvalidInstanceId.NotFound"}, want: true},
		"dataMatchNonEmptyCode":             {err: teaSDKError(403, "APINotAuthorize", "forbidden", "previous response 404"), expectCodes: []string{"404"}, want: true},
		"dataMatchEmptyCode":                {err: teaSDKError(403, "", "forbidden", "previous response 404"), expectCodes: []string{"404"}, want: true},
		"emptyExpectCodes":                  {err: teaSDKError(404, "404", "not found", ""), expectCodes: []string{}, want: false},
		"complexWrapsTeaExact":              {err: errs.WrapError(teaSDKError(404, "404", "not found", "")), expectCodes: []string{"404"}, want: true},
		"complexWrapsTeaEmptyCode":          {err: errs.WrapError(teaSDKError(403, "", "Permission denied", "")), expectCodes: []string{"404"}, want: false},
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			if got := errs.IsExpectedErrors(c.err, c.expectCodes); got != c.want {
				t.Fatalf("IsExpectedErrors() = %v, want %v", got, c.want)
			}
		})
	}
}

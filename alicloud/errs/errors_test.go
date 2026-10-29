package errs_test

import (
	"testing"

	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
)

// teaSDKError builds a *tea.SDKError the way the live RPC layer surfaces
// failures: Code/Message/Data are the structured fields and StatusCode is
// the HTTP status. NotFoundError reads these directly (tea.IntValue/
// tea.StringValue), so the constructed value exercises the real branch.
func teaSDKError(statusCode int, code, message, data string) *tea.SDKError {
	return &tea.SDKError{
		Code:       tea.String(code),
		Message:    tea.String(message),
		Data:       tea.String(data),
		StatusCode: tea.Int(statusCode),
	}
}

// TestUnitNotFoundErrorTeaSDKErrorByCodeNotMessage pins the fixed tea
// SDKError branch: classification follows HTTP status and the structured
// error Code, never a "NotFound" substring in the human-readable Message.
//
// The first group is genuine resource-missing (must stay true): HTTP 404,
// the SDK-level "404" code, and the product not-found codes that carry a
// NotFound/NotExist token in Code. The second group is the regression
// target — non-missing errors whose Message happens to contain "NotFound";
// these previously cleared live state in Read (d.SetId("")) and faked
// destruction in CheckDestroy, and must now be false.
func TestUnitNotFoundErrorTeaSDKErrorByCodeNotMessage(t *testing.T) {
	cases := []struct {
		name string
		err  error
		want bool
	}{
		// Genuine not-found: HTTP status and Code-level markers.
		{name: "http404-genuine-not-found", err: teaSDKError(404, "NotFound", "The resource not found.", ""), want: true},
		{name: "code-404-sdk-level-marker-on-2xx", err: teaSDKError(200, "404", "code-level not found in response body", ""), want: true},
		{name: "product-code-InvalidInstanceId-NotFound", err: teaSDKError(400, "InvalidInstanceId.NotFound", "The specified instance is not found.", ""), want: true},
		{name: "product-code-ResourceNotFound", err: teaSDKError(404, "ResourceNotFound", "The resource not found.", ""), want: true},
		{name: "product-code-EntityNotExist-User", err: teaSDKError(404, "EntityNotExist.User", "The user does not exist.", ""), want: true},
		{name: "product-code-EntityNotExists-ResourceDirectory", err: teaSDKError(404, "EntityNotExists.ResourceDirectory", "The resource directory does not exist.", ""), want: true},
		{name: "ram-code-Forbidden-InstanceNotFound-403", err: teaSDKError(403, "Forbidden.InstanceNotFound", "The RAM instance is not found.", ""), want: true},
		{name: "product-code-Instance-Notfound-lowercase", err: teaSDKError(400, "Instance.Notfound", "instance not found", ""), want: true},

		// Regression: false positives the old Message-substring match produced.
		// The Message contains "NotFound" but the Code is auth/throttle/service,
		// so the resource is not missing. Compare with the http404 case above.
		{name: "arms-counter-403-APINotAuthorize-NotFound-message", err: teaSDKError(403, "APINotAuthorize", "NotFound while resolving authorization", ""), want: false},
		{name: "403-Forbidden-RAM-Instance-NotFound-denied-message", err: teaSDKError(403, "Forbidden.RAM", "Instance.NotFound denied", ""), want: false},
		{name: "403-access-denied-404-in-metadata", err: teaSDKError(403, "Forbidden.AccessDenied", "access denied", "HttpStatus 404"), want: false},
		{name: "500-InternalError-NotFound-message", err: teaSDKError(500, "InternalError", "NotFound handler failed", ""), want: false},
		{name: "503-ServiceUnavailable-NotFound-message", err: teaSDKError(503, "ServiceUnavailable", "service unavailable NotFound", ""), want: false},
		{name: "429-Throttling-NotFound-message", err: teaSDKError(429, "Throttling", "Throttling NotFound retry later", ""), want: false},
		{name: "empty-code-lowercase-message", err: teaSDKError(400, "", "instance is not found", ""), want: false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := errs.NotFoundError(c.err); got != c.want {
				t.Fatalf("NotFoundError(%v) = %v, want %v", c.err, got, c.want)
			}
		})
	}
}

// TestUnitNotFoundErrorWrapErrorfPassThrough simulates the VPC Read chain:
// DescribeVpcAttribute fails with a *tea.SDKError, the service layer wraps
// it via WrapErrorf(DefaultErrorMsg, ...) into a ComplexError{Cause:
// *tea.SDKError}, and NotFoundError recurses Cause to reach the tea branch.
// The wrap must not obscure classification — the SDK error's Code still
// decides, so the ARMS-style 403 stays false while a genuine 404 stays true.
func TestUnitNotFoundErrorWrapErrorfPassThrough(t *testing.T) {
	wrap := func(cause error) error {
		return errs.WrapErrorf(cause, errs.DefaultErrorMsg, "the-id", "DescribeX", errs.AlibabaCloudSdkGoERROR)
	}
	cases := []struct {
		name string
		err  error
		want bool
	}{
		{name: "wrapped-arms-counter-403-APINotAuthorize", err: wrap(teaSDKError(403, "APINotAuthorize", "NotFound while resolving authorization", "")), want: false},
		{name: "wrapped-403-Forbidden-RAM-NotFound-message", err: wrap(teaSDKError(403, "Forbidden.RAM", "Instance.NotFound denied", "")), want: false},
		{name: "wrapped-500-InternalError-NotFound-message", err: wrap(teaSDKError(500, "InternalError", "NotFound handler failed", "")), want: false},
		{name: "wrapped-http404-genuine-not-found", err: wrap(teaSDKError(404, "NotFound", "The resource not found.", "")), want: true},
		{name: "wrapped-product-code-InvalidInstanceId-NotFound", err: wrap(teaSDKError(400, "InvalidInstanceId.NotFound", "The specified instance is not found.", "")), want: true},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := errs.NotFoundError(c.err); got != c.want {
				t.Fatalf("NotFoundError(%v) = %v, want %v", c.err, got, c.want)
			}
		})
	}
}

// TestUnitNotFoundErrorProviderMarkers covers the provider's own explicit
// not-found markers, which do not flow through the tea SDKError branch and
// must keep working: the ComplexError branch (NotFoundErr+NotFoundMsg build
// a ResourceNotfound-prefixed Err) and the ProviderError branch
// (GetNotFoundErrorFromString sets errorCode=InstanceNotFound).
func TestUnitNotFoundErrorProviderMarkers(t *testing.T) {
	cases := []struct {
		name string
		err  error
		want bool
	}{
		{name: "complex-error-NotFoundMsg-marker", err: errs.WrapErrorf(errs.NotFoundErr("Vpc", "vpc-123"), errs.NotFoundMsg, "response-body"), want: true},
		{name: "provider-error-GetNotFoundErrorFromString", err: errs.GetNotFoundErrorFromString("the instance is not found"), want: true},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := errs.NotFoundError(c.err); got != c.want {
				t.Fatalf("NotFoundError(%v) = %v, want %v", c.err, got, c.want)
			}
		})
	}
}

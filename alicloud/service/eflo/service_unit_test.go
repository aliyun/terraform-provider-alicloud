package eflo

import (
	"fmt"
	"testing"

	"github.com/alibabacloud-go/tea/tea"
)

// TestUnitRetryableErAttachmentError pins the retry policy to the shared transient
// classifier: quoted and unquoted transport failures and throttling retry, while
// terminal API answers — including the 1003 not-found used for absence — do not.
func TestUnitRetryableErAttachmentError(t *testing.T) {
	cases := []struct {
		name string
		err  error
		want bool
	}{
		{name: "nil", err: nil, want: false},
		{name: "quoted network error", err: fmt.Errorf(`Post "https://eflo.cn-hangzhou.aliyuncs.com": EOF`), want: true},
		{name: "unquoted network error", err: fmt.Errorf("Post https://eflo.cn-hangzhou.aliyuncs.com: EOF"), want: true},
		{name: "throttling", err: sdkError("Throttling.User", "Request was denied due to user flow control."), want: true},
		{name: "service unavailable", err: sdkError("ServiceUnavailable", "The request has failed due to a temporary failure of the server."), want: true},
		{name: "not found", err: sdkError("1003", "er attachment not found"), want: false},
		{name: "invalid parameter", err: sdkError("InvalidParameter", "The specified parameter is invalid."), want: false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := retryableErAttachmentError(tc.err); got != tc.want {
				t.Fatalf("retryableErAttachmentError(%v) = %v, want %v", tc.err, got, tc.want)
			}
		})
	}
}

func sdkError(code, message string) error {
	return &tea.SDKError{
		Code:    tea.String(code),
		Message: tea.String(message),
		Data:    tea.String(""),
	}
}

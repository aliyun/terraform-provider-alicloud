package connectivity

import (
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/alibabacloud-go/tea/tea"
)

func TestUnitSplitRpcBudget(t *testing.T) {
	cases := []struct {
		name            string
		budget          time.Duration
		wantConnect     time.Duration
		wantRead        time.Duration
		wantReadMinimum bool
	}{
		{name: "long budget caps connect", budget: 300 * time.Second, wantConnect: 10 * time.Second, wantRead: 290 * time.Second},
		{name: "short budget is split", budget: 5 * time.Second, wantConnect: 2500 * time.Millisecond, wantRead: 2500 * time.Millisecond},
		{name: "very short budget keeps a read window", budget: 2 * time.Millisecond, wantConnect: time.Millisecond, wantRead: time.Millisecond, wantReadMinimum: true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			connect, read := splitRpcBudget(tc.budget)
			if connect != tc.wantConnect || read != tc.wantRead {
				t.Fatalf("split(%s) = (%s, %s), want (%s, %s)", tc.budget, connect, read, tc.wantConnect, tc.wantRead)
			}
			if connect+read > tc.budget {
				t.Fatalf("split(%s) exceeds the operation budget: (%s, %s)", tc.budget, connect, read)
			}
			if read < time.Millisecond {
				t.Fatalf("read budget %s is below the header wait floor", read)
			}
		})
	}
}

func TestUnitSanitizeRpcErrorRedactsSignedUrls(t *testing.T) {
	signed := `Post "https://slb.aliyuncs.com/?Action=SetLoadBalancerHTTPListenerAttribute&AccessKeyId=LTAI-test&Signature=abc123&SignatureNonce=xyz&SecurityToken=token-value": dial tcp 10.0.0.1:443: i/o timeout`
	cause := errors.New(signed)

	sanitized := sanitizeRpcError(cause)
	message := sanitized.Error()
	for _, secret := range []string{"AccessKeyId", "Signature=", "SignatureNonce", "SecurityToken", "LTAI-test", "abc123", "token-value"} {
		if strings.Contains(message, secret) {
			t.Fatalf("the sanitized error still exposes %q: %s", secret, message)
		}
	}
	if !strings.Contains(message, "slb.aliyuncs.com") {
		t.Fatalf("the sanitized error lost the endpoint: %s", message)
	}
	if !errors.Is(sanitized, cause) {
		t.Fatal("the sanitized error no longer wraps the original error")
	}
	// The transport classification still sees a Post request.
	if !strings.HasPrefix(message, "Post \"https://") {
		t.Fatalf("the retry classification prefix was lost: %s", message)
	}
}

func TestUnitSanitizeRpcErrorKeepsSdkCode(t *testing.T) {
	sdkErr := &tea.SDKError{
		Code:    tea.String("IpVersionConflict"),
		Message: tea.String(`code: 400, conflict https://slb.aliyuncs.com/?Signature=abc&AccessKeyId=LTAI-test`),
		Data:    tea.String(`{"url":"https://slb.aliyuncs.com/?Signature=abc&AccessKeyId=LTAI-test"}`),
	}
	sanitized, ok := sanitizeRpcError(sdkErr).(*tea.SDKError)
	if !ok {
		t.Fatalf("the SDK error type was not preserved: %T", sanitizeRpcError(sdkErr))
	}
	if tea.StringValue(sanitized.Code) != "IpVersionConflict" {
		t.Fatalf("the SDK error code changed: %q", tea.StringValue(sanitized.Code))
	}
	if strings.Contains(tea.StringValue(sanitized.Message), "Signature=") || strings.Contains(tea.StringValue(sanitized.Data), "AccessKeyId") {
		t.Fatalf("the SDK error still exposes credentials: %q / %q", tea.StringValue(sanitized.Message), tea.StringValue(sanitized.Data))
	}
}

func TestUnitDecodeRpcResponse(t *testing.T) {
	newResponse := func(status int, body string) *http.Response {
		return &http.Response{StatusCode: status, Body: io.NopCloser(strings.NewReader(body)), Header: http.Header{}}
	}

	t.Run("json object", func(t *testing.T) {
		result, err := decodeRpcResponse(newResponse(200, `{"RequestId":"r-1","AclStatus":"on"}`))
		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}
		if result["AclStatus"] != "on" {
			t.Fatalf("unexpected payload: %v", result)
		}
	})

	t.Run("error status keeps the sdk code", func(t *testing.T) {
		_, err := decodeRpcResponse(newResponse(400, `{"Code":"IpVersionConflict","Message":"conflict","RequestId":"r-1"}`))
		sdkErr, ok := err.(*tea.SDKError)
		if !ok {
			t.Fatalf("expected an SDK error, got %T: %v", err, err)
		}
		if tea.StringValue(sdkErr.Code) != "IpVersionConflict" || tea.IntValue(sdkErr.StatusCode) != 400 {
			t.Fatalf("unexpected SDK error: %+v", sdkErr)
		}
	})

	for _, tc := range []struct {
		name   string
		status int
		body   string
	}{
		{name: "empty body"},
		{name: "json array", status: 200, body: `[1,2,3]`},
		{name: "json scalar", status: 200, body: `"text"`},
		{name: "invalid json", status: 200, body: `{"Code":`},
	} {
		t.Run(tc.name, func(t *testing.T) {
			status := tc.status
			if status == 0 {
				status = 200
			}
			_, err := decodeRpcResponse(newResponse(status, tc.body))
			if !errors.Is(err, ErrUnreadableResponse) {
				t.Fatalf("expected ErrUnreadableResponse, got %v", err)
			}
		})
	}

	t.Run("nil response", func(t *testing.T) {
		if _, err := decodeRpcResponse(nil); !errors.Is(err, ErrUnreadableResponse) {
			t.Fatalf("expected ErrUnreadableResponse, got %v", err)
		}
	})

	t.Run("truncated body", func(t *testing.T) {
		response := &http.Response{StatusCode: 200, Body: io.NopCloser(&failingReader{}), Header: http.Header{}}
		if _, err := decodeRpcResponse(response); !errors.Is(err, ErrUnreadableResponse) {
			t.Fatalf("expected ErrUnreadableResponse, got %v", err)
		}
	})
}

type failingReader struct{}

func (f *failingReader) Read([]byte) (int, error) { return 0, errors.New("unexpected EOF") }

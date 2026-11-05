package connectivity

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
)

func TestUnitRpcTransportTimeouts(t *testing.T) {
	cases := []struct {
		name        string
		budget      time.Duration
		wantConnect time.Duration
		wantRead    time.Duration
	}{
		{name: "long budget caps only the dial", budget: 300 * time.Second, wantConnect: 10 * time.Second, wantRead: 300 * time.Second},
		{name: "short budget keeps the whole header window", budget: 5 * time.Second, wantConnect: 5 * time.Second, wantRead: 5 * time.Second},
		{name: "very short budget keeps a dial and header window", budget: 2 * time.Millisecond, wantConnect: 2 * time.Millisecond, wantRead: 2 * time.Millisecond},
		{name: "sub millisecond budget is rounded up", budget: 500 * time.Microsecond, wantConnect: time.Millisecond, wantRead: time.Millisecond},
		{name: "no deadline keeps the transport defaults", budget: 0, wantConnect: 0, wantRead: 0},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			connect, read := rpcTransportTimeouts(tc.budget)
			if connect != tc.wantConnect || read != tc.wantRead {
				t.Fatalf("timeouts(%s) = (%s, %s), want (%s, %s)", tc.budget, connect, read, tc.wantConnect, tc.wantRead)
			}
		})
	}

	// The header wait is not a fixed share of the budget: a fast dial must leave
	// the whole remaining budget for the response headers.
	if connect, read := rpcTransportTimeouts(10 * time.Second); connect+read <= 10*time.Second {
		t.Fatalf("the header window was pre-split: connect %s, read %s", connect, read)
	}
}

func TestUnitRedactURLMalformedEscape(t *testing.T) {
	// url.Parse rejects the malformed path escape, which used to make the
	// redactor return the complete signed URL unchanged.
	raw := `https://slb.aliyuncs.com/%zz?Action=SetLoadBalancerHTTPListenerAttribute&AccessKeyId=LTAI-test&Signature=abc123&SignatureNonce=xyz&SecurityToken=token-value`
	redacted := redactURL(raw)
	for _, secret := range []string{"Action=", "AccessKeyId", "Signature=", "SignatureNonce", "SecurityToken", "LTAI-test", "abc123", "token-value"} {
		if strings.Contains(redacted, secret) {
			t.Fatalf("the redacted URL still exposes %q: %s", secret, redacted)
		}
	}
	if !strings.Contains(redacted, "slb.aliyuncs.com") {
		t.Fatalf("the redacted URL lost the endpoint: %s", redacted)
	}
}

func TestUnitSanitizeRpcErrorRedactsMalformedSignedUrl(t *testing.T) {
	cause := errors.New(`parse "https://slb.aliyuncs.com/%zz?AccessKeyId=LTAI-test&Signature=abc123&SecurityToken=token-value": invalid URL escape "%zz"`)
	sanitized := sanitizeRpcError(cause)
	message := sanitized.Error()
	for _, secret := range []string{"AccessKeyId", "Signature=", "SecurityToken", "LTAI-test", "abc123", "token-value"} {
		if strings.Contains(message, secret) {
			t.Fatalf("the sanitized error still exposes %q: %s", secret, message)
		}
	}
	if !errors.Is(sanitized, cause) {
		t.Fatal("the sanitized error no longer wraps the original error")
	}
	if helper.NeedRetry(sanitized) && !helper.NeedRetry(cause) {
		t.Fatal("the redaction changed the retry classification")
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
	newResponse := func(status int, body string) *tea.Response {
		return &tea.Response{StatusCode: tea.Int(status), Body: io.NopCloser(strings.NewReader(body)), Headers: map[string]*string{}}
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
		{name: "bad body on a 5xx response", status: 500, body: `<html>gateway</html>`},
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
			if errors.Is(err, ErrRpcRejectedResponse) {
				t.Fatalf("an unreadable 2xx/5xx body must stay unknown, got a rejection: %v", err)
			}
		})
	}

	t.Run("bad body on a 4xx response is an explicit rejection", func(t *testing.T) {
		for _, body := range []string{"", "<html>Bad Request</html>", `{"Code":`} {
			_, err := decodeRpcResponse(newResponse(400, body))
			if !errors.Is(err, ErrRpcRejectedResponse) {
				t.Fatalf("a 4xx rejection with body %q must be explicit, got %v", body, err)
			}
			if errors.Is(err, ErrUnreadableResponse) {
				t.Fatalf("a 4xx rejection must not be reported as an unknown outcome: %v", err)
			}
		}
	})

	// A 404 without a readable body may come from a gateway or proxy that never
	// reached the SLB service, so the status alone must not be reported as a
	// missing resource: the outcome is an explicit rejection with the real status
	// preserved, and no business "not found" code is invented.
	t.Run("bad body on a 404 response is an explicit rejection", func(t *testing.T) {
		_, err := decodeRpcResponse(newResponse(404, ""))
		if errors.Is(err, ErrUnreadableResponse) {
			t.Fatalf("a 404 must keep its status classification, got an unknown outcome: %v", err)
		}
		if !errors.Is(err, ErrRpcRejectedResponse) {
			t.Fatalf("an undecodable 404 must be an explicit rejection, got %v", err)
		}
		if strings.Contains(err.Error(), "status 404") == false {
			t.Fatalf("the real status was lost: %v", err)
		}
		if _, ok := err.(*tea.SDKError); ok {
			t.Fatal("an undecodable 404 must not be turned into a business SDK error")
		}
	})

	// A decodable 4xx body that carries no business code must report an empty
	// code, not the stringified missing value: callers classify the error by the
	// code, and "<nil>" would look like a real one.
	t.Run("missing business code decodes as empty", func(t *testing.T) {
		_, err := decodeRpcResponse(newResponse(403, `{"Forbidden":true}`))
		sdkErr, ok := err.(*tea.SDKError)
		if !ok {
			t.Fatalf("expected a typed business error, got %T: %v", err, err)
		}
		if got := tea.StringValue(sdkErr.Code); got != "" {
			t.Fatalf("a response without a business code reported %q", got)
		}
		if strings.Contains(sdkErr.Error(), "<nil>") {
			t.Fatalf("the decoded error carries a stringified missing value: %s", sdkErr.Error())
		}
	})

	// The control case: a decodable 404 that carries the SLB business code stays
	// usable, so callers can still prove that the entity itself is gone.
	t.Run("decoded 404 keeps the slb business code", func(t *testing.T) {
		_, err := decodeRpcResponse(newResponse(404, `{"Code":"InvalidLoadBalancerId.NotFound","Message":"The specified resource does not exist.","RequestId":"r-1"}`))
		sdkErr, ok := err.(*tea.SDKError)
		if !ok {
			t.Fatalf("expected a typed business error, got %T: %v", err, err)
		}
		if tea.StringValue(sdkErr.Code) != "InvalidLoadBalancerId.NotFound" || tea.IntValue(sdkErr.StatusCode) != 404 {
			t.Fatalf("the business code was lost: %+v", sdkErr)
		}
	})

	t.Run("decoded error does not echo the signed request url", func(t *testing.T) {
		// A gateway error page may quote the request URL. The decoded error is
		// sanitized before it can reach a Terraform diagnostic.
		_, err := decodeRpcResponse(newResponse(400, `{"Code":"InvalidAction","Message":"bad request https://slb.aliyuncs.com/?AccessKeyId=LTAI-test&Signature=abc123&SecurityToken=token-value","RequestId":"r-1"}`))
		if err == nil {
			t.Fatal("expected the 400 response to fail")
		}
		message := sanitizeRpcError(err).Error()
		for _, secret := range []string{"AccessKeyId", "Signature=", "SecurityToken", "LTAI-test", "abc123", "token-value"} {
			if strings.Contains(message, secret) {
				t.Fatalf("the decoded error exposes %q: %s", secret, message)
			}
		}
		if !strings.Contains(message, "InvalidAction") {
			t.Fatalf("the decoded error lost its classification: %s", message)
		}
	})

	t.Run("nil response", func(t *testing.T) {
		if _, err := decodeRpcResponse(nil); !errors.Is(err, ErrUnreadableResponse) {
			t.Fatalf("expected ErrUnreadableResponse, got %v", err)
		}
	})

	t.Run("truncated body", func(t *testing.T) {
		response := &tea.Response{StatusCode: tea.Int(200), Body: io.NopCloser(&failingReader{}), Headers: map[string]*string{}}
		if _, err := decodeRpcResponse(response); !errors.Is(err, ErrUnreadableResponse) {
			t.Fatalf("expected ErrUnreadableResponse, got %v", err)
		}
	})
}

type failingReader struct{}

func (f *failingReader) Read([]byte) (int, error) { return 0, errors.New("unexpected EOF") }

// The debug log boundary receives the provider payloads, so the redaction has to
// reach every string shape the service logs: plain values, nested objects, lists
// and the response headers map. A service may echo only the query string instead
// of the complete URL, or repeat the signed parameters as independent fields,
// which carry the same credentials.
func TestUnitRedactLogValue(t *testing.T) {
	signed := "https://slb.aliyuncs.com/?AccessKeyId=LTAI-test&Signature=synthetic&SecurityToken=token-value"
	queryOnly := "AccessKeyId=LTAI-test&Signature=synthetic&SignatureNonce=nonce-value&SecurityToken=token-value"
	redacted := RedactLogValue(map[string]interface{}{
		"AclStatus":   "on",
		"RequestId":   "req-1",
		"Echo":        signed,
		"QueryEcho":   queryOnly,
		"Nested":      map[string]interface{}{"Message": "request " + signed},
		"List":        []interface{}{signed, "plain"},
		"_headers":    map[string]*string{"location": tea.String(signed), "empty": nil},
		"AccessKeyId": "independent-ak",
		"Signature":   "independent-signature",
		"Auth": map[string]interface{}{
			"SecurityToken": "nested-token",
			"Details": []interface{}{
				map[string]interface{}{"SignatureNonce": "nested-nonce"},
				"plain-detail",
			},
		},
		"securitytoken": []interface{}{"list-token-1", "list-token-2"},
		"HeaderMap":     map[string]*string{"Signature": tea.String("header-signature")},
		// A credential key whose value is an object or a header map must redact the
		// values below it, not only the string shapes.
		"AccessKeySecret": map[string]interface{}{"Value": "object-secret"},
		"signaturenonce":  map[string]*string{"x-nonce": tea.String("header-nonce")},
	})
	payload, err := json.Marshal(redacted)
	if err != nil {
		t.Fatalf("marshalling the redacted payload: %s", err)
	}
	// The parameter names may stay in the log; only their values are secret.
	for _, secret := range []string{
		"LTAI-test", "synthetic", "token-value", "nonce-value",
		"independent-ak", "independent-signature", "nested-token", "nested-nonce",
		"list-token-1", "list-token-2", "header-signature", "object-secret", "header-nonce",
	} {
		if strings.Contains(string(payload), secret) {
			t.Fatalf("the logged payload still exposes %q: %s", secret, payload)
		}
	}
	if strings.Contains(string(payload), "Signature=nonce") || strings.Contains(string(payload), "AccessKeyId=LTAI") {
		t.Fatalf("a signed parameter pair survived the redaction: %s", payload)
	}
	if !strings.Contains(string(payload), "redacted") {
		t.Fatalf("the signed parameter rule did not run: %s", payload)
	}
	// Public metadata and ordinary values stay readable.
	for _, kept := range []string{"slb.aliyuncs.com", "AclStatus", "req-1", "plain", "plain-detail"} {
		if !strings.Contains(string(payload), kept) {
			t.Fatalf("the redaction dropped the usable payload item %q: %s", kept, payload)
		}
	}
	if redacted.(map[string]interface{})["AclStatus"] != "on" {
		t.Fatal("a value without a URL must be preserved")
	}
	if list, ok := redacted.(map[string]interface{})["List"].([]interface{}); !ok || list[1] != "plain" {
		t.Fatalf("a plain list item was changed: %v", redacted.(map[string]interface{})["List"])
	}
}

// A service or an intermediary may echo the signed request as independent fields
// (or as header entries) instead of a URL. Such an echo has no scheme and no
// "key=value" text, so it used to survive the textual URL redaction and reach the
// SDK error and the Terraform diagnostic. The structural rule has to hide every
// credential value while keeping the business code that callers classify by.
func TestUnitSanitizeRpcErrorRedactsStructuredCredentials(t *testing.T) {
	sdkErr := &tea.SDKError{
		Code:       tea.String("InvalidLoadBalancerId.NotFound"),
		StatusCode: tea.Int(404),
		Message:    tea.String(`code: 404, refused request {"Signature":"synthetic-message-signature"}`),
		Data: tea.String(`{
			"AccessKeyId": "LTAI-synthetic-key",
			"Signature": "synthetic-data-signature",
			"SignatureNonce": "synthetic-nonce",
			"SecurityToken": "synthetic-token",
			"nested": {"AccessKeySecret": "synthetic-secret"},
			"list": [{"SecurityToken": "synthetic-list-token"}],
			"RequestId": "req-1",
			"Code": "InvalidLoadBalancerId.NotFound"
		}`),
	}
	sanitized, ok := sanitizeRpcError(sdkErr).(*tea.SDKError)
	if !ok {
		t.Fatalf("the SDK error type was not preserved: %T", sanitizeRpcError(sdkErr))
	}
	if got := tea.StringValue(sanitized.Code); got != "InvalidLoadBalancerId.NotFound" {
		t.Fatalf("the classification code was lost or changed: %q", got)
	}
	if got := tea.IntValue(sanitized.StatusCode); got != 404 {
		t.Fatalf("the status code was lost: %d", got)
	}
	rendered := tea.StringValue(sanitized.Message) + " " + tea.StringValue(sanitized.Data)
	for _, secret := range []string{
		"LTAI-synthetic-key", "synthetic-data-signature", "synthetic-nonce",
		"synthetic-token", "synthetic-secret", "synthetic-list-token",
		"synthetic-message-signature",
	} {
		if strings.Contains(rendered, secret) {
			t.Fatalf("the sanitized SDK error still exposes %q: %s", secret, rendered)
		}
	}
	for _, kept := range []string{"RequestId", "req-1", "InvalidLoadBalancerId.NotFound"} {
		if !strings.Contains(rendered, kept) {
			t.Fatalf("the sanitized SDK error dropped the usable item %q: %s", kept, rendered)
		}
	}
}

// The decoded error is built from a response that may carry the signed request
// both as independent body fields and as response headers. Both boundaries are
// redacted before the error can reach a diagnostic.
func TestUnitDecodeRpcResponseRedactsStructuredEcho(t *testing.T) {
	body := `{
		"Code": "InvalidLoadBalancerId.NotFound",
		"Message": "The specified resource does not exist.",
		"AccessKeyId": "LTAI-synthetic-key",
		"Signature": "synthetic-body-signature",
		"SecurityToken": "synthetic-body-token",
		"RequestId": "req-1"
	}`
	response := &tea.Response{
		StatusCode: tea.Int(404),
		Body:       io.NopCloser(strings.NewReader(body)),
		Headers: map[string]*string{
			"x-acs-signature":      tea.String("synthetic-header-signature"),
			"x-acs-security-token": tea.String("synthetic-header-token"),
			"x-acs-request-id":     tea.String("req-1"),
		},
	}
	_, err := decodeRpcResponse(response)
	if err == nil {
		t.Fatal("expected the 404 response to fail")
	}
	sanitized, ok := sanitizeRpcError(err).(*tea.SDKError)
	if !ok {
		t.Fatalf("expected a typed business error, got %T: %v", sanitizeRpcError(err), err)
	}
	if got := tea.StringValue(sanitized.Code); got != "InvalidLoadBalancerId.NotFound" {
		t.Fatalf("the business code was lost: %q", got)
	}
	rendered := tea.StringValue(sanitized.Data)
	for _, secret := range []string{
		"LTAI-synthetic-key", "synthetic-body-signature", "synthetic-body-token",
		"synthetic-header-signature", "synthetic-header-token",
	} {
		if strings.Contains(rendered, secret) {
			t.Fatalf("the SDK error payload still exposes %q: %s", secret, rendered)
		}
	}
	if !strings.Contains(rendered, "req-1") {
		t.Fatalf("the usable request id was dropped: %s", rendered)
	}
}

// A redirect is a non success receipt. The signed request must not be followed to
// another endpoint (its Referer would carry the credentials), and the answer must
// be classified as an explicit rejection rather than a successful payload.
func TestUnitSignedRpcRequestDoesNotFollowRedirect(t *testing.T) {
	var receivedByTarget int32
	target := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&receivedByTarget, 1)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"AclStatus":"on"}`))
	}))
	defer target.Close()

	var originRequests int32
	origin := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&originRequests, 1)
		if referer := r.Referer(); referer != "" {
			t.Errorf("the resolved request carried a Referer: %s", referer)
		}
		http.Redirect(w, r, target.URL+"/redirected?AccessKeyId=forwarded", http.StatusFound)
	}))
	defer origin.Close()

	originURL, err := url.Parse(origin.URL)
	if err != nil {
		t.Fatalf("parsing the origin URL: %s", err)
	}
	request := tea.NewRequest()
	request.Protocol = tea.String("HTTP")
	request.Method = tea.String("POST")
	request.Pathname = tea.String("/")
	request.Headers = map[string]*string{"host": tea.String(originURL.Host)}
	request.Query = map[string]*string{
		"Action":        tea.String("SetLoadBalancerHTTPListenerAttribute"),
		"AccessKeyId":   tea.String("LTAI-synthetic-key"),
		"Signature":     tea.String("synthetic-signature"),
		"SecurityToken": tea.String("synthetic-token"),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	response, err := executeSignedRpcRequest(ctx, request)
	if err != nil {
		t.Fatalf("executing the signed request: %s", err)
	}
	if got := tea.IntValue(response.StatusCode); got != http.StatusFound {
		t.Fatalf("the redirect status was lost: %d", got)
	}
	if got := atomic.LoadInt32(&originRequests); got != 1 {
		t.Fatalf("the origin received %d requests, want exactly 1", got)
	}
	if got := atomic.LoadInt32(&receivedByTarget); got != 0 {
		t.Fatalf("the redirect target received %d request(s); the signed request must not be followed", got)
	}
	if _, err := decodeRpcResponse(response); !errors.Is(err, ErrRpcRejectedResponse) {
		t.Fatalf("a redirect must be an explicit rejection, got %v", err)
	} else if errors.Is(err, ErrUnreadableResponse) {
		t.Fatalf("a redirect must not be reported as an unknown write: %v", err)
	}
}

// A redirect with a readable JSON body must not be mistaken for a success: the
// body is a redirect receipt, not the answer of the addressed action.
func TestUnitDecodeRpcResponseRedirectWithJsonBodyIsRejected(t *testing.T) {
	response := &tea.Response{
		StatusCode: tea.Int(307),
		Body:       io.NopCloser(strings.NewReader(`{"AclStatus":"on","AclId":"acl-other"}`)),
		Headers:    map[string]*string{},
	}
	result, err := decodeRpcResponse(response)
	if err == nil {
		t.Fatalf("a redirect payload was accepted as a result: %v", result)
	}
	if !errors.Is(err, ErrRpcRejectedResponse) {
		t.Fatalf("a redirect must be an explicit rejection, got %v", err)
	}
	if !strings.Contains(err.Error(), "status 307") {
		t.Fatalf("the redirect status was lost: %v", err)
	}
}

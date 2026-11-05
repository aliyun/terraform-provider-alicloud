package connectivity

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	utilV2 "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/credentials-go/credentials"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
)

// newHangingServer accepts the request and holds the response until the client
// disconnects. Draining the request body first is what lets the server notice
// the disconnect while the handler is still running; without it the connection
// stays active and the handler never sees the cancellation.
func newHangingServer(t *testing.T) (server *httptest.Server, arrived, aborted <-chan struct{}) {
	t.Helper()
	arrivedCh := make(chan struct{})
	abortedCh := make(chan struct{})
	var arrivedOnce, abortedOnce sync.Once
	server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		arrivedOnce.Do(func() { close(arrivedCh) })
		_, _ = io.Copy(io.Discard, r.Body)
		select {
		case <-r.Context().Done():
			abortedOnce.Do(func() { close(abortedCh) })
		case <-time.After(30 * time.Second):
		}
	}))
	t.Cleanup(func() {
		server.CloseClientConnections()
		server.Close()
	})
	return server, arrivedCh, abortedCh
}

func newLoopbackClient(t *testing.T, serverURL string) *AliyunClient {
	t.Helper()
	credential, err := credentials.NewCredential(new(credentials.Config).
		SetType("access_key").
		SetAccessKeyId("access-key").
		SetAccessKeySecret("secret-key"))
	require.NoError(t, err)

	var endpoints sync.Map
	endpoints.Store("ecs", strings.TrimPrefix(serverURL, "http://"))

	return &AliyunClient{
		config: &Config{
			RegionId:             "cn-beijing",
			Protocol:             "HTTP",
			Credential:           credential,
			ClientReadTimeout:    60000,
			ClientConnectTimeout: 60000,
			Endpoints:            &endpoints,
		},
	}
}

func waitForSignal(t *testing.T, signal <-chan struct{}, name string) {
	t.Helper()
	select {
	case <-signal:
	case <-time.After(5 * time.Second):
		t.Fatalf("%s did not happen", name)
	}
}

// A read that is cancelled while its request is in flight must abort that
// request. The numeric timeout alone would leave it running to completion.
func TestUnitDoWithContextAbortsInFlightRequestOnCancel(t *testing.T) {
	server, arrived, aborted := newHangingServer(t)
	client := newLoopbackClient(t, server.URL)

	ctx, cancel := context.WithCancel(context.Background())
	result := make(chan error, 1)
	go func() {
		_, err := client.DoWithContext(ctx, "Ecs", helper.RpcParam("POST", "2014-05-26", "DescribeInstances"), nil,
			map[string]interface{}{"RegionId": "cn-beijing"}, nil, nil, &utilV2.RuntimeOptions{})
		result <- err
	}()

	waitForSignal(t, arrived, "the request reaching the server")
	cancel()

	select {
	case err := <-result:
		assert.Error(t, err, "a cancelled read must fail instead of returning a response")
	case <-time.After(5 * time.Second):
		t.Fatal("the call did not return after the context was cancelled")
	}
	waitForSignal(t, aborted, "the server observing the client disconnect")
}

// The read deadline has to reach the transport: a request that outlives the
// deadline is aborted instead of waiting for the server.
func TestUnitDoWithContextAbortsInFlightRequestOnDeadline(t *testing.T) {
	server, arrived, aborted := newHangingServer(t)
	client := newLoopbackClient(t, server.URL)

	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()

	start := time.Now()
	result := make(chan error, 1)
	go func() {
		_, err := client.DoWithContext(ctx, "Ecs", helper.RpcParam("POST", "2014-05-26", "DescribeInstances"), nil,
			map[string]interface{}{"RegionId": "cn-beijing"}, nil, nil, &utilV2.RuntimeOptions{})
		result <- err
	}()

	waitForSignal(t, arrived, "the request reaching the server")
	select {
	case err := <-result:
		assert.Error(t, err, "an expired read must fail instead of returning a response")
		assert.Less(t, time.Since(start), 4*time.Second, "the transport must stop at the deadline instead of the client read timeout")
	case <-time.After(5 * time.Second):
		t.Fatal("the call did not return when the deadline expired")
	}
	waitForSignal(t, aborted, "the server observing the client disconnect")
}

// newStaticServer answers every request with the given status and body, which
// is what a real HTTP fixture needs to pin the body shapes the SDK accepts.
func newStaticServer(t *testing.T, status int, body string) *httptest.Server {
	t.Helper()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = io.Copy(io.Discard, r.Body)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		if body != "" {
			_, _ = w.Write([]byte(body))
		}
	}))
	t.Cleanup(server.Close)
	return server
}

// The body check itself is pinned for every shape the SDK can hand back: the
// transport wrapper with or without a body, the already decoded body, and the
// non-object body types.
func TestUnitRequireContextJSONBody(t *testing.T) {
	for name, response := range map[string]map[string]interface{}{
		"nil response":            nil,
		"empty response":          {},
		"wrapper without body":    {"statusCode": 200, "headers": map[string]interface{}{}},
		"wrapper with empty body": {"statusCode": 200, "body": map[string]interface{}{}},
		"wrapper with nil body":   {"statusCode": 200, "body": nil},
	} {
		assert.Error(t, requireContextJSONBody(response, "Ecs"), name+" must be rejected")
	}

	for name, response := range map[string]map[string]interface{}{
		"decoded body":          {"RequestId": "mock", "Instances": map[string]interface{}{"Instance": []interface{}{}}},
		"wrapper with object":   {"statusCode": 200, "body": map[string]interface{}{"RequestId": "mock"}},
		"wrapper with a string": {"statusCode": 200, "body": "not-a-json-document"},
	} {
		assert.NoError(t, requireContextJSONBody(response, "Ecs"), name+" must be accepted")
	}
}

// A 2xx response without a JSON object body must fail the read. Handing the
// empty map back would be read as a zero result and replace the state the
// caller already holds.
func TestUnitDoWithContextRejectsAnEmptySuccessBody(t *testing.T) {
	for name, body := range map[string]string{
		"empty body": "",
		"null body":  "null",
	} {
		server := newStaticServer(t, http.StatusOK, body)
		client := newLoopbackClient(t, server.URL)

		_, err := client.DoWithContext(context.Background(), "Ecs", helper.RpcParam("POST", "2014-05-26", "DescribeInstances"), nil,
			map[string]interface{}{"RegionId": "cn-beijing"}, nil, nil, &utilV2.RuntimeOptions{})
		if assert.Error(t, err, "%s must not be returned as an empty result", name) {
			assert.Contains(t, err.Error(), "without a JSON object body")
		}
	}

	// A body that is not JSON at all fails in the SDK decoder; it must fail as
	// well, and it must not be turned into an empty result either.
	server := newStaticServer(t, http.StatusOK, "  \n")
	client := newLoopbackClient(t, server.URL)
	_, err := client.DoWithContext(context.Background(), "Ecs", helper.RpcParam("POST", "2014-05-26", "DescribeInstances"), nil,
		map[string]interface{}{"RegionId": "cn-beijing"}, nil, nil, &utilV2.RuntimeOptions{})
	assert.Error(t, err, "a blank body must not be returned as an empty result")
}

// A legitimate empty result keeps its JSON object: an empty list inside the
// response is a successful read, not a missing body.
func TestUnitDoWithContextKeepsALegitimateEmptyList(t *testing.T) {
	server := newStaticServer(t, http.StatusOK, `{"RequestId":"mock","Instances":{"Instance":[]}}`)
	client := newLoopbackClient(t, server.URL)

	response, err := client.DoWithContext(context.Background(), "Ecs", helper.RpcParam("POST", "2014-05-26", "DescribeInstances"), nil,
		map[string]interface{}{"RegionId": "cn-beijing"}, nil, nil, &utilV2.RuntimeOptions{})
	require.NoError(t, err)
	assert.Equal(t, "mock", response["RequestId"])
	instances, ok := response["Instances"].(map[string]interface{})
	require.True(t, ok, "the nested object must survive the body check")
	assert.Empty(t, instances["Instance"])
}

// A 2xx error-shaped body may carry a numeric StatusCode in any shape the
// transports produce: the legacy path used *int while the context bound path
// decodes JSON numbers as json.Number. Reading it must never panic and must keep
// the status the service reported.
func TestUnitFormatErrorReadsNumericStatusCode(t *testing.T) {
	cases := []struct {
		name  string
		value interface{}
		want  int
	}{
		{"json number", json.Number("500"), 500},
		{"float", float64(503), 503},
		{"numeric string", "429", 429},
		{"pointer", tea.Int(502), 502},
		{"nil value", nil, 200},
		{"fractional", 1.5, 200},
		{"unreadable", "oops", 200},
		{"negative", -1, -1},
	}
	for _, testCase := range cases {
		err := formatError(map[string]interface{}{
			"Code":       "InvalidSomething",
			"Message":    "denied",
			"StatusCode": testCase.value,
		}, nil)
		var sdkErr *tea.SDKError
		require.ErrorAs(t, err, &sdkErr, testCase.name)
		assert.Equal(t, testCase.want, tea.IntValue(sdkErr.StatusCode), testCase.name)
	}
}

// The body check must not replace an SDK error: a failed call keeps its own
// code and message.
func TestUnitDoWithContextKeepsTheSDKError(t *testing.T) {
	server := newStaticServer(t, http.StatusBadRequest, `{"RequestId":"mock","Code":"InvalidInstanceId.NotFound","Message":"The specified InstanceId does not exist."}`)
	client := newLoopbackClient(t, server.URL)

	_, err := client.DoWithContext(context.Background(), "Ecs", helper.RpcParam("POST", "2014-05-26", "DescribeInstances"), nil,
		map[string]interface{}{"RegionId": "cn-beijing"}, nil, nil, &utilV2.RuntimeOptions{})
	require.Error(t, err, "a failed call must stay failed")
	assert.Contains(t, err.Error(), "InvalidInstanceId.NotFound", "the original SDK error must be preserved")
	assert.NotContains(t, err.Error(), "without a JSON object body")
}

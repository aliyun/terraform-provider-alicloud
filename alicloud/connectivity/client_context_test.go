package connectivity

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	utilV2 "github.com/alibabacloud-go/tea-utils/v2/service"
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

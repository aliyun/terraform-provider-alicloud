package connectivity

import (
	"context"
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
// disconnects. It reports arrival and the moment the client went away, so a
// test can prove the request physically reached the server and was then aborted.
func newHangingServer(t *testing.T) (*httptest.Server, <-chan struct{}, <-chan struct{}) {
	t.Helper()
	arrived := make(chan struct{})
	aborted := make(chan struct{})
	var arrivedOnce, abortedOnce sync.Once
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		arrivedOnce.Do(func() { close(arrived) })
		<-r.Context().Done()
		abortedOnce.Do(func() { close(aborted) })
	}))
	return server, arrived, aborted
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
	defer server.Close()
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
	defer server.Close()
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

package connectivity

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/aliyun/credentials-go/credentials"
	"github.com/stretchr/testify/assert"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
)

// TestUnitSlsRequestHonorsClientTimeouts pins the fix in the SLS branch of
// AliyunClient.Do. That branch rebuilds sdkConfig from an empty openapi.Config
// (required for the tmp-ak signature fix), which used to drop the
// provider-level client_connect_timeout / client_read_timeout and left the
// underlying http.Client with no total deadline. Against a server that accepts
// the request but never answers, the call must fail within the configured
// timeout budget instead of hanging until an outer resource.Retry gives up.
func TestUnitSlsRequestHonorsClientTimeouts(t *testing.T) {
	var arrived atomic.Bool
	release := make(chan struct{})
	var releaseOnce sync.Once
	releaseFn := func() { releaseOnce.Do(func() { close(release) }) }
	// Safety valve so a control run without the fix cannot deadlock the
	// binary: the elapsed-time assertion below is what fails in that run.
	go func() {
		time.Sleep(10 * time.Second)
		releaseFn()
	}()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		arrived.Store(true)
		<-release
	}))
	defer server.Close()
	defer releaseFn()

	credential, err := credentials.NewCredential(new(credentials.Config).
		SetType("access_key").
		SetAccessKeyId("access-key").
		SetAccessKeySecret("secret-key"))
	assert.NoError(t, err)

	// Byte-for-byte what provider endpoints { log = ... } stores (scheme-less
	// host, the SDK prepends the protocol), so the endpoint resolution below
	// never leaves the loopback server.
	var endpoints sync.Map
	endpoints.Store("log", strings.TrimPrefix(server.URL, "http://"))

	client := &AliyunClient{
		config: &Config{
			RegionId:             "cn-beijing",
			Protocol:             "HTTP",
			Credential:           credential,
			ClientReadTimeout:    1000,
			ClientConnectTimeout: 1000,
			Endpoints:            &endpoints,
		},
	}

	start := time.Now()
	_, err = client.Do("Sls", helper.RoaParam("GET", "2020-12-30", "GetProject", "/"), map[string]*string{}, nil, nil, nil, false)
	elapsed := time.Since(start)

	assert.True(t, arrived.Load(), "the request should physically reach the slow server")
	assert.Error(t, err, "a stalled response must surface as an error, not block forever")
	assert.Less(t, elapsed, 8*time.Second, "the request must respect connect+read timeout instead of waiting for the 10s safety valve")
	assert.Contains(t, strings.ToLower(err.Error()), "timeout", "the error should be the client timeout, got: %v", err)
}

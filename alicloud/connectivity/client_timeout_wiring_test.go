package connectivity

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk"
	sdkCredentials "github.com/aliyun/alibaba-cloud-sdk-go/sdk/auth/credentials"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	ossv2 "github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
	"github.com/aliyun/credentials-go/credentials"
	"github.com/aliyun/fc-go-sdk"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// The FC SDK takes its whole-request timeout in seconds and treats 0 as
// "no timeout", so an unset provider value must keep the SDK default and a
// positive sub-second value must be rounded up instead of collapsing to 0.
func TestUnitFcClientUsesClientReadTimeout(t *testing.T) {
	credential, err := credentials.NewCredential(new(credentials.Config).
		SetType("access_key").
		SetAccessKeyId("access-key").
		SetAccessKeySecret("secret-key"))
	require.NoError(t, err)

	// fc.WithTimeout also writes the process-wide resty timeout; restore it so
	// other tests in this package keep the SDK default.
	defer fc.WithTimeout(fc.RequestTimeout)(&fc.Client{Connect: &fc.Connection{}})

	var endpoints sync.Map
	endpoints.Store("fc", "fc.unit-test.aliyuncs.com")

	testCases := []struct {
		name      string
		timeoutMs int
		want      uint
	}{
		{name: "provider read timeout is converted to seconds", timeoutMs: 45000, want: 45},
		{name: "sub-second read timeout is rounded up to one second", timeoutMs: 500, want: 1},
		{name: "partial seconds are rounded up", timeoutMs: 1999, want: 2},
		{name: "unset read timeout keeps the sdk default", timeoutMs: 0, want: fc.RequestTimeout},
		{name: "negative read timeout keeps the sdk default", timeoutMs: -1500, want: fc.RequestTimeout},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			client := &AliyunClient{
				accountId: "1234567890",
				config: &Config{
					RegionId:          "cn-beijing",
					Protocol:          "HTTPS",
					Credential:        credential,
					ClientReadTimeout: tc.timeoutMs,
					Endpoints:         &endpoints,
					SignVersion:       &sync.Map{},
				},
			}

			var got uint
			_, err := client.WithFcClient(func(c *fc.Client) (interface{}, error) {
				got = c.Connect.Timeout
				return nil, nil
			})

			assert.NoError(t, err)
			assert.Equal(t, tc.want, got)
		})
	}
}

// WithOssClientV2 used to keep the SDK defaults (5s connect / 10s read-write).
// A stalled response must now fail inside the provider-configured budget
// instead of the SDK default one.
func TestUnitOssV2ClientUsesClientTimeouts(t *testing.T) {
	for _, name := range []string{"HTTP_PROXY", "http_proxy", "HTTPS_PROXY", "https_proxy"} {
		t.Setenv(name, "")
	}

	var arrived atomic.Bool
	release := make(chan struct{})
	var releaseOnce sync.Once
	releaseFn := func() { releaseOnce.Do(func() { close(release) }) }
	// Safety valve so a run without the wiring cannot deadlock the binary:
	// the elapsed-time assertion below is what fails in that run.
	go func() {
		time.Sleep(20 * time.Second)
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
	require.NoError(t, err)

	var endpoints sync.Map
	endpoints.Store("oss", strings.TrimPrefix(server.URL, "http://"))

	client := &AliyunClient{
		accountId: "1234567890",
		config: &Config{
			RegionId:             "cn-beijing",
			Protocol:             "HTTP",
			Credential:           credential,
			ClientReadTimeout:    300,
			ClientConnectTimeout: 300,
			Endpoints:            &endpoints,
			SignVersion:          &sync.Map{},
		},
	}

	start := time.Now()
	_, err = client.WithOssClientV2(func(c *ossv2.Client) (interface{}, error) {
		_, err := c.GetBucketInfo(context.Background(), &ossv2.GetBucketInfoRequest{Bucket: ossv2.Ptr("unit-test-bucket")})
		return nil, err
	})
	elapsed := time.Since(start)

	assert.True(t, arrived.Load(), "the request should physically reach the stalled server")
	assert.Error(t, err, "a stalled response must surface as an error, not block forever")
	assert.Less(t, elapsed, 9*time.Second, "the request must respect the configured read-write timeout instead of the SDK default")
}

func TestUnitOssV2ClientTimeoutsKeepSdkDefaultsWhenUnset(t *testing.T) {
	cfg := withOssV2ClientTimeouts(ossv2.LoadDefaultConfig(), 0, 0)
	assert.Nil(t, cfg.ConnectTimeout, "a zero connect timeout must keep the SDK default")
	assert.Nil(t, cfg.ReadWriteTimeout, "a zero read-write timeout must keep the SDK default")

	cfg = withOssV2ClientTimeouts(ossv2.LoadDefaultConfig(), 1500, 2500)
	if assert.NotNil(t, cfg.ConnectTimeout) {
		assert.Equal(t, 1500*time.Millisecond, *cfg.ConnectTimeout)
	}
	if assert.NotNil(t, cfg.ReadWriteTimeout) {
		assert.Equal(t, 2500*time.Millisecond, *cfg.ReadWriteTimeout)
	}
}

// getSdkConfig is shared by the legacy alibaba-cloud-sdk-go clients. Its
// sdk-level timeout follows the provider client_read_timeout, and a
// non-positive value keeps the 30s fallback so zero-config callers (e.g. the
// sweeper client) do not drop to the SDK's own default. The per-client
// SetReadTimeout calls in the With*Client closures take precedence at request
// time, so this test covers the shared config field only.
func TestUnitLegacySdkConfigUsesClientReadTimeout(t *testing.T) {
	testCases := []struct {
		name      string
		timeoutMs int
		want      time.Duration
	}{
		{name: "provider read timeout is converted to a duration", timeoutMs: 45000, want: 45 * time.Second},
		{name: "sub-second read timeout is preserved", timeoutMs: 1500, want: 1500 * time.Millisecond},
		{name: "unset read timeout keeps the 30s fallback", timeoutMs: 0, want: 30 * time.Second},
		{name: "negative read timeout keeps the 30s fallback", timeoutMs: -1500, want: 30 * time.Second},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			client := &AliyunClient{
				config: &Config{
					RegionId:          "cn-beijing",
					Protocol:          "HTTPS",
					ClientReadTimeout: tc.timeoutMs,
				},
			}

			assert.Equal(t, tc.want, client.getSdkConfig().Timeout)
		})
	}
}

// The public client constructors forward the raw provider timeouts to the
// client-level setters, which the SDK checks before the shared config value.
func TestUnitLegacyVpcClientUsesClientTimeouts(t *testing.T) {
	for _, name := range []string{"HTTP_PROXY", "http_proxy", "HTTPS_PROXY", "https_proxy"} {
		t.Setenv(name, "")
	}

	credential, err := credentials.NewCredential(new(credentials.Config).
		SetType("access_key").
		SetAccessKeyId("access-key").
		SetAccessKeySecret("secret-key"))
	require.NoError(t, err)

	var endpoints sync.Map
	endpoints.Store("vpc", "vpc.unit-test.aliyuncs.com")

	testCases := []struct {
		name         string
		readTimeout  int
		connectTimer int
		wantRead     time.Duration
		wantConnect  time.Duration
	}{
		{name: "provider timeouts reach the client", readTimeout: 45000, connectTimer: 2500, wantRead: 45 * time.Second, wantConnect: 2500 * time.Millisecond},
		{name: "unset timeouts leave the client-level fields unset", readTimeout: 0, connectTimer: 0, wantRead: 0, wantConnect: 0},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			client := &AliyunClient{
				accountId: "1234567890",
				config: &Config{
					RegionId:             "cn-beijing",
					Protocol:             "HTTPS",
					Credential:           credential,
					ClientReadTimeout:    tc.readTimeout,
					ClientConnectTimeout: tc.connectTimer,
					Endpoints:            &endpoints,
					SignVersion:          &sync.Map{},
				},
			}

			var gotRead, gotConnect time.Duration
			_, err := client.WithVpcClient(func(c *vpc.Client) (interface{}, error) {
				gotRead = c.GetReadTimeout()
				gotConnect = c.GetConnectTimeout()
				return nil, nil
			})
			require.NoError(t, err)

			assert.Equal(t, tc.wantRead, gotRead)
			assert.Equal(t, tc.wantConnect, gotConnect)
		})
	}
}

// A legacy client built from getSdkConfig must stop waiting for a stalled
// response once the configured read timeout elapses. The sdk client is built
// directly from the config because most With*Client closures also call
// SetReadTimeout, which would mask the shared config value.
func TestUnitLegacySdkClientUsesClientReadTimeout(t *testing.T) {
	for _, name := range []string{"HTTP_PROXY", "http_proxy", "HTTPS_PROXY", "https_proxy"} {
		t.Setenv(name, "")
	}

	var arrived atomic.Bool
	release := make(chan struct{})
	var releaseOnce sync.Once
	releaseFn := func() { releaseOnce.Do(func() { close(release) }) }
	// Safety valve so a run without the wiring cannot hang the binary: the
	// elapsed-time assertion below is what fails in that run.
	go func() {
		time.Sleep(12 * time.Second)
		releaseFn()
	}()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		arrived.Store(true)
		<-release
	}))
	defer server.Close()
	defer releaseFn()

	client := &AliyunClient{
		config: &Config{
			RegionId:          "cn-beijing",
			Protocol:          "HTTP",
			ClientReadTimeout: 500,
		},
	}
	sdkClient, err := sdk.NewClientWithOptions("cn-beijing", client.getSdkConfig(), sdkCredentials.NewAccessKeyCredential("access-key", "secret-key"))
	require.NoError(t, err)

	request := requests.NewCommonRequest()
	request.Method = http.MethodGet
	request.Scheme = "http"
	request.Domain = strings.TrimPrefix(server.URL, "http://")
	request.Version = "2014-05-26"
	request.ApiName = "DescribeRegions"
	request.Product = "Ecs"

	start := time.Now()
	_, err = sdkClient.ProcessCommonRequest(request)
	elapsed := time.Since(start)

	assert.True(t, arrived.Load(), "the request should physically reach the stalled server")
	assert.Error(t, err, "a stalled response must surface as an error once the configured read timeout elapses")
	assert.Less(t, elapsed, 10*time.Second, "the request must respect the configured read timeout instead of the sdk default")
}

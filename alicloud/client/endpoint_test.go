package client

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/location"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/client/config"
	"github.com/stretchr/testify/assert"
)

func TestUnitCommonLoadEndpointFromEnv(t *testing.T) {
	t.Setenv("ALIBABA_CLOUD_ENDPOINT_ECS", "https://ecs.test.com")
	defer os.Unsetenv("ALIBABA_CLOUD_ENDPOINT_ECS")

	client := &AliyunClient{
		config: &config.Config{
			Endpoints: new(sync.Map),
			RegionId:  "cn-beijing",
		},
	}

	err := client.config.LoadEndpoint("ecs", client.RegionId)
	assert.NoError(t, err)

	val, ok := client.config.Endpoints.Load("ecs")
	assert.True(t, ok)
	assert.Equal(t, "https://ecs.test.com", val)
}

func TestUnitCommonIrregularProductEndpoint(t *testing.T) {
	client := &AliyunClient{
		RegionId: "cn-hangzhou",
		config: &config.Config{
			Endpoints: new(sync.Map),
			RegionId:  "cn-hangzhou",
		},
	}

	testCases := []struct {
		productCode string
		expected    string
	}{
		{"ram", "ram.aliyuncs.com"},
		{"cloudfw", "cloudfw.aliyuncs.com"},
		{"agentloop", "agentloop.cn-hangzhou.aliyuncs.com"},
	}

	for _, tc := range testCases {
		t.Run(tc.productCode, func(t *testing.T) {
			err := client.config.LoadEndpoint(tc.productCode, client.RegionId)
			assert.NoError(t, err)

			val, ok := client.config.Endpoints.Load(tc.productCode)
			assert.True(t, ok)
			assert.Equal(t, tc.expected, val)
		})
	}
}

func TestUnitCommonInternationalRegionEndpoint(t *testing.T) {
	client := &AliyunClient{
		config: &config.Config{
			Endpoints: new(sync.Map),
			RegionId:  "ap-southeast-1",
		},
	}

	testCases := []struct {
		productCode string
		expected    string
	}{
		{"cloudfw", "cloudfw.aliyuncs.com"},
		{"sas", "tds.ap-southeast-1.aliyuncs.com"},
	}

	for _, tc := range testCases {
		t.Run(tc.productCode, func(t *testing.T) {
			err := client.config.LoadEndpoint(tc.productCode, client.RegionId)
			assert.NoError(t, err)

			val, ok := client.config.Endpoints.Load(tc.productCode)
			assert.True(t, ok)
			assert.Equal(t, tc.expected, val)
		})
	}
}

func TestUnitCommonEndpointErrorHandling(t *testing.T) {
	client := NewTestClient(t)

	err := client.config.LoadEndpoint("invalid_product", client.RegionId)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Illegal")
}

func TestUnitSlsEndpoint(t *testing.T) {
	var endpoints sync.Map
	endpoints.Store("log", "cn-shanghai.log.aliyuncs.com")

	client := NewTestClient(t)
	client.config.Region = config.Hangzhou
	client.config.RegionId = "cn-hangzhou"
	client.config.Endpoints = &endpoints

	testCases := []struct {
		productCode    string
		configEndpoint string
		expected       string
	}{
		{"sls", "log", "cn-shanghai.log.aliyuncs.com"},
		{"alb", "alb", "alb.cn-hangzhou.aliyuncs.com"},
	}

	for _, tc := range testCases {
		t.Run(tc.productCode, func(t *testing.T) {
			endpoint, err := client.loadApiEndpoint(tc.productCode)
			assert.NoError(t, err)
			assert.Equal(t, tc.expected, endpoint)

			val, ok := client.config.Endpoints.Load(tc.configEndpoint)
			assert.True(t, ok)
			assert.Equal(t, tc.expected, val)
		})
	}
}

func TestUnitEfloControllerEndpointFromEnv(t *testing.T) {
	// Explicit provider endpoint config and environment variables take
	// precedence over the regularProductEndpoint fallback.
	t.Setenv("ALIBABA_CLOUD_ENDPOINT_EFLO_CONTROLLER", "https://eflo-controller.test.com")
	defer os.Unsetenv("ALIBABA_CLOUD_ENDPOINT_EFLO_CONTROLLER")

	client := &AliyunClient{
		config: &config.Config{
			Endpoints: new(sync.Map),
			RegionId:  "us-southeast-1",
		},
	}

	err := client.config.LoadEndpoint("eflo_controller", client.RegionId)
	assert.NoError(t, err)

	val, ok := client.config.Endpoints.Load("eflo_controller")
	assert.True(t, ok)
	assert.Equal(t, "https://eflo-controller.test.com", val)
}

func TestUnitSlsWithoutEndpoint(t *testing.T) {
	testCases := []struct {
		productCode    string
		configEndpoint string
		RegionId       string
		expected       string
	}{
		{"sls", "", "cn-hangzhou", "cn-hangzhou.log.aliyuncs.com"},
		{"slb", "", "cn-hangzhou", "slb.aliyuncs.com"},
		{"sls", "", "eu-central-1", "eu-central-1.log.aliyuncs.com"},
		{"alb", "", "eu-central-1", "alb.eu-central-1.aliyuncs.com"},
	}

	client := NewTestClient(t)

	for _, tc := range testCases {
		t.Run(tc.productCode, func(t *testing.T) {
			// A fresh endpoint cache keeps the endpoint discovered for another
			// region out of this case.
			client.config.Endpoints = new(sync.Map)
			client.config.Region = config.Region(tc.RegionId)
			client.config.RegionId = tc.RegionId

			endpoint, err := client.loadApiEndpoint(tc.productCode)
			assert.NoError(t, err)
			assert.Equal(t, tc.expected, endpoint)

			val, ok := client.config.Endpoints.Load(tc.productCode)
			assert.True(t, ok)
			assert.Equal(t, tc.expected, val)
		})
	}
}

func TestUnitSlsEndpointCacheReuse(t *testing.T) {
	// Mirror the reported scenario: the location service is unreachable, so the
	// first call falls back to the regular endpoint form and logs a WARN. Every
	// following call must be served from the endpoint cache, without probing the
	// location service again and without logging another WARN.
	var locationProbes int32
	locationProbe := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&locationProbes, 1)
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer locationProbe.Close()

	var logs bytes.Buffer
	previousLogOutput := log.Writer()
	log.SetOutput(&logs)
	defer log.SetOutput(previousLogOutput)

	t.Setenv("ALIBABA_CLOUD_ENDPOINT_SLS", "")
	t.Setenv("SLS_ENDPOINT", "")

	client := &AliyunClient{
		RegionId: "cn-hangzhou",
		config: &config.Config{
			Endpoints: new(sync.Map),
			RegionId:  "cn-hangzhou",
			// HTTP keeps the simulated unreachable-location failure out of the
			// SDK's HTTPS-specific retry branch, which would take minutes.
			Protocol:         "HTTP",
			LocationEndpoint: locationProbe.Listener.Addr().String(),
			AccessKey:        "test-ak",
			SecretKey:        "test-sk",
			Credential:       &mockCredential{err: fmt.Errorf("unit test uses the static access key")},
		},
	}

	first, err := client.loadApiEndpoint("sls")
	assert.NoError(t, err)
	assert.Equal(t, "cn-hangzhou.log.aliyuncs.com", first)
	assert.NotZero(t, atomic.LoadInt32(&locationProbes))

	cached, ok := client.config.Endpoints.Load("sls")
	assert.True(t, ok)
	assert.Equal(t, first, cached)
	_, logCached := client.config.Endpoints.Load("log")
	assert.False(t, logCached)
	assert.Equal(t, 1, strings.Count(logs.String(), "[WARN] loading sls endpoint got an error"))

	// Change the probe source. An un-cached implementation would return the new
	// value here; a cache hit keeps the first endpoint.
	t.Setenv("ALIBABA_CLOUD_ENDPOINT_SLS", "https://sls.second.test.com")
	probesAfterFirstCall := atomic.LoadInt32(&locationProbes)

	second, err := client.loadApiEndpoint("sls")
	assert.NoError(t, err)
	assert.Equal(t, first, second)
	assert.Equal(t, probesAfterFirstCall, atomic.LoadInt32(&locationProbes))
	assert.Equal(t, 1, strings.Count(logs.String(), "[WARN] loading sls endpoint got an error"))

	// Clear the probe source again. With the endpoint cached nothing may probe
	// the location service or log another WARN here; without the cache this call
	// re-runs discovery and breaks both counters, while still returning the same
	// regular endpoint form, so the counters are what distinguish the two paths.
	t.Setenv("ALIBABA_CLOUD_ENDPOINT_SLS", "")

	third, err := client.loadApiEndpoint("sls")
	assert.NoError(t, err)
	assert.Equal(t, first, third)
	assert.Equal(t, probesAfterFirstCall, atomic.LoadInt32(&locationProbes))
	assert.Equal(t, 1, strings.Count(logs.String(), "[WARN] loading sls endpoint got an error"))
}

func TestUnitSlsConfiguredEndpointPriority(t *testing.T) {
	t.Setenv("ALIBABA_CLOUD_ENDPOINT_SLS", "https://sls.env.test.com")

	var endpoints sync.Map
	endpoints.Store("log", "cn-shanghai.log.aliyuncs.com")

	client := &AliyunClient{
		config: &config.Config{
			Endpoints: &endpoints,
			RegionId:  "cn-hangzhou",
		},
	}

	endpoint, err := client.loadApiEndpoint("sls")
	assert.NoError(t, err)
	assert.Equal(t, "cn-shanghai.log.aliyuncs.com", endpoint)

	cached, ok := client.config.Endpoints.Load("log")
	assert.True(t, ok)
	assert.Equal(t, "cn-shanghai.log.aliyuncs.com", cached)
}

func TestUnitProductEndpointCacheReuse(t *testing.T) {
	testCases := []struct {
		productCode string
		envName     string
	}{
		{"ecs", "ALIBABA_CLOUD_ENDPOINT_ECS"},
		{"alb", "ALIBABA_CLOUD_ENDPOINT_ALB"},
	}

	for _, tc := range testCases {
		t.Run(tc.productCode, func(t *testing.T) {
			firstName := fmt.Sprintf("https://%s.first.test.com", tc.productCode)
			secondName := fmt.Sprintf("https://%s.second.test.com", tc.productCode)
			t.Setenv(tc.envName, firstName)

			client := &AliyunClient{
				config: &config.Config{
					Endpoints: new(sync.Map),
					RegionId:  "cn-hangzhou",
				},
			}

			first, err := client.loadApiEndpoint(tc.productCode)
			assert.NoError(t, err)
			assert.Equal(t, firstName, first)

			cached, ok := client.config.Endpoints.Load(tc.productCode)
			assert.True(t, ok)
			assert.Equal(t, firstName, cached)

			t.Setenv(tc.envName, secondName)

			second, err := client.loadApiEndpoint(tc.productCode)
			assert.NoError(t, err)
			assert.Equal(t, first, second)

			cached, ok = client.config.Endpoints.Load(tc.productCode)
			assert.True(t, ok)
			assert.Equal(t, first, cached)
		})
	}
}

// TestUnitSlsEndpointUserConfigTakesPrecedenceOverCache verifies the fix keeps
// the v1.253.0 user-config priority: endpoints.log for sls wins over a
// previously probed endpoint cached under "sls".
func TestUnitSlsEndpointUserConfigTakesPrecedenceOverCache(t *testing.T) {
	var endpoints sync.Map
	endpoints.Store("log", "user.log.example.com")
	endpoints.Store("sls", "cached.sls.example.com")

	client := &AliyunClient{
		config: &config.Config{
			Endpoints: &endpoints,
			RegionId:  "cn-hangzhou",
		},
	}

	ep, err := client.loadApiEndpoint("sls")
	assert.NoError(t, err)
	assert.Equal(t, "user.log.example.com", ep, "user-configured log endpoint must take precedence over the cached sls endpoint")

	val, ok := client.config.Endpoints.Load("log")
	assert.True(t, ok)
	assert.Equal(t, "user.log.example.com", val)

	cached, ok := client.config.Endpoints.Load("sls")
	assert.True(t, ok)
	assert.Equal(t, "cached.sls.example.com", cached, "the cached sls entry must be left untouched")
}

// TestUnitEndpointDiscoveryUsesProviderProxyRules drives uncached NewCommonRequest construction
// through local HTTP stubs. The proxied case must not fall back to a direct request, while NO_PROXY
// must bypass the configured proxy and reach the Location stub directly.
func TestUnitEndpointDiscoveryUsesProviderProxyRules(t *testing.T) {
	const locationResponse = `{"RequestId":"test-request","Success":true,"Endpoints":{"Endpoint":[{"Endpoint":"cn-hangzhou.log.aliyuncs.com","Id":"cn-hangzhou","ServiceCode":"sls","Type":"openAPI"}]}}`

	for _, tc := range []struct {
		name          string
		bypassProxy   bool
		wantProxy     int32
		wantDirect    int32
		wantProxyHost bool
	}{
		{name: "HTTP proxy", wantProxy: 1, wantProxyHost: true},
		{name: "NO_PROXY", bypassProxy: true, wantDirect: 1},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var proxyRequests int32
			var directRequests int32
			proxyHosts := make(chan string, 2)
			writeLocationResponse := func(w http.ResponseWriter) {
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write([]byte(locationResponse))
			}

			locationServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				atomic.AddInt32(&directRequests, 1)
				writeLocationResponse(w)
			}))
			defer locationServer.Close()
			proxyServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				atomic.AddInt32(&proxyRequests, 1)
				proxyHosts <- r.URL.Host
				writeLocationResponse(w)
			}))
			defer proxyServer.Close()

			t.Setenv("HTTP_PROXY", proxyServer.URL)
			t.Setenv("http_proxy", "")
			t.Setenv("HTTPS_PROXY", "")
			t.Setenv("https_proxy", "")
			t.Setenv("NO_PROXY", "")
			t.Setenv("no_proxy", "")
			if tc.bypassProxy {
				t.Setenv("NO_PROXY", locationServer.Listener.Addr().String())
			}

			client := &AliyunClient{
				RegionId: "cn-hangzhou",
				config: &config.Config{
					Endpoints:            new(sync.Map),
					RegionId:             "cn-hangzhou",
					Protocol:             "HTTP",
					LocationEndpoint:     locationServer.Listener.Addr().String(),
					ClientReadTimeout:    1000,
					ClientConnectTimeout: 1000,
					AccessKey:            "test-ak",
					SecretKey:            "test-sk",
					Credential:           &mockCredential{err: fmt.Errorf("unit test uses the static access key")},
				},
			}
			request, err := client.NewCommonRequest("sls", "log", "HTTP", config.ApiVersion("2015-06-01"))
			if assert.NoError(t, err) && assert.NotNil(t, request) {
				assert.Equal(t, "cn-hangzhou.log.aliyuncs.com", request.Domain)
			}
			assert.Equal(t, tc.wantProxy, atomic.LoadInt32(&proxyRequests))
			assert.Equal(t, tc.wantDirect, atomic.LoadInt32(&directRequests))
			if tc.wantProxyHost {
				select {
				case host := <-proxyHosts:
					assert.Equal(t, locationServer.Listener.Addr().String(), host)
				default:
					t.Fatal("proxy did not receive the Location request")
				}
			}
		})
	}
}

func TestUnitEndpointRetryBackoff(t *testing.T) {
	assert.Equal(t, 3*time.Second, endpointRetryBackoff(1))
	assert.Equal(t, 5*time.Second, endpointRetryBackoff(2))
	assert.Equal(t, 5*time.Second, endpointRetryBackoff(3))
}

// TestUnitEndpointRetryWaitHonoursCancellation pins that the discovery backoff stops as soon as the
// operation is cancelled instead of sleeping out the whole pause.
func TestUnitEndpointRetryWaitHonoursCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		time.Sleep(50 * time.Millisecond)
		cancel()
	}()
	started := time.Now()
	err := endpointRetryWait(ctx, 5*time.Second)
	elapsed := time.Since(started)
	assert.ErrorIs(t, err, context.Canceled)
	assert.Less(t, elapsed, time.Second, "the backoff must stop when the operation is cancelled")
}

// TestUnitEndpointRetryWaitHonoursDeadline pins that a near operation deadline caps the backoff and
// reports the deadline. The capped pause ends on the context itself, so the result never depends on
// which timer fires first.
func TestUnitEndpointRetryWaitHonoursDeadline(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Millisecond)
	defer cancel()
	started := time.Now()
	err := endpointRetryWait(ctx, 5*time.Second)
	elapsed := time.Since(started)
	assert.ErrorIs(t, err, context.DeadlineExceeded)
	assert.Less(t, elapsed, time.Second, "the backoff must not outlast the operation deadline")
	assert.GreaterOrEqual(t, elapsed, 50*time.Millisecond)
}

// TestUnitEndpointRetryWaitReturnsAfterBackoff keeps the pause itself in the happy path.
func TestUnitEndpointRetryWaitReturnsAfterBackoff(t *testing.T) {
	started := time.Now()
	err := endpointRetryWait(context.Background(), 20*time.Millisecond)
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, time.Since(started), 20*time.Millisecond)
}

// TestUnitEndpointDiscoveryCancelsDuringBackoff drives the real discovery path with an unreachable
// location service: the first attempt fails with the read timeout so the retry backoff is entered,
// then the operation is cancelled mid-backoff. Without the cancellable wait the loop sleeps its
// full three seconds; with it the call returns promptly and no further attempt is made.
func TestUnitEndpointDiscoveryCancelsDuringBackoff(t *testing.T) {
	var connections int32
	probe := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Hold the request until the client gives up, so the attempt fails with the read timeout.
		<-r.Context().Done()
	}))
	probe.Config.ConnState = func(_ net.Conn, state http.ConnState) {
		if state == http.StateNew {
			atomic.AddInt32(&connections, 1)
		}
	}
	probe.Start()
	defer probe.Close()

	client := &AliyunClient{
		RegionId: "cn-hangzhou",
		config: &config.Config{
			Endpoints:            new(sync.Map),
			RegionId:             "cn-hangzhou",
			Protocol:             "HTTP",
			LocationEndpoint:     probe.Listener.Addr().String(),
			ClientReadTimeout:    100,
			ClientConnectTimeout: 100,
			AccessKey:            "test-ak",
			SecretKey:            "test-sk",
			Credential:           &mockCredential{err: fmt.Errorf("unit test uses the static access key")},
		},
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		time.Sleep(1200 * time.Millisecond)
		cancel()
	}()

	started := time.Now()
	_, err := client.describeEndpointForService(ctx, "sls", client.RegionId)
	elapsed := time.Since(started)
	assert.Error(t, err)
	assert.GreaterOrEqual(t, elapsed, 700*time.Millisecond, "the call must still be blocked when it is cancelled")
	assert.Less(t, elapsed, 3*time.Second, "cancelling during the backoff must return before the pause ends")

	settled := atomic.LoadInt32(&connections)
	assert.NotZero(t, settled)
	time.Sleep(300 * time.Millisecond)
	assert.Equal(t, settled, atomic.LoadInt32(&connections), "no discovery attempt may run after the cancellation")
}

func TestUnitEndpointDiscoveryCancellationStopsSDKRetries(t *testing.T) {
	var requests int32
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if atomic.AddInt32(&requests, 1) == 1 {
			cancel()
			<-r.Context().Done()
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusServiceUnavailable)
		_, _ = w.Write([]byte(`{"Code":"ServiceUnavailable"}`))
	}))
	defer server.Close()

	client := &AliyunClient{
		RegionId: "cn-hangzhou",
		config: &config.Config{
			Endpoints:            new(sync.Map),
			RegionId:             "cn-hangzhou",
			Protocol:             "HTTP",
			LocationEndpoint:     server.Listener.Addr().String(),
			ClientReadTimeout:    5000,
			ClientConnectTimeout: 1000,
			AccessKey:            "test-ak",
			SecretKey:            "test-sk",
			Credential:           &mockCredential{err: fmt.Errorf("unit test uses the static access key")},
		},
	}
	started := time.Now()
	_, err := client.describeEndpointForService(ctx, "sls", client.RegionId)
	assert.ErrorIs(t, err, context.Canceled)
	assert.Less(t, time.Since(started), 2*time.Second)
	assert.Equal(t, int32(1), atomic.LoadInt32(&requests))
	time.Sleep(100 * time.Millisecond)
	assert.Equal(t, int32(1), atomic.LoadInt32(&requests), "the SDK must not retransmit after cancellation")
}

// TestUnitEndpointDiscoveryRecalculatesAttemptTimeout makes the second attempt start after the
// retry backoff. Its read timeout must use the then-remaining operation budget; reusing the first
// attempt's larger timeout would let the SDK call outlive the lifecycle deadline.
func TestUnitEndpointDiscoveryRecalculatesAttemptTimeout(t *testing.T) {
	originalDescribe := describeLocationEndpoints
	t.Cleanup(func() { describeLocationEndpoints = originalDescribe })
	var readTimeouts []time.Duration
	describeLocationEndpoints = func(_ context.Context, _ *location.Client, _ *location.DescribeEndpointsRequest, readTimeout, _ time.Duration, _ *http.Transport, _ *url.URL) (*location.DescribeEndpointsResponse, error) {
		readTimeouts = append(readTimeouts, readTimeout)
		if len(readTimeouts) == 1 {
			return nil, fmt.Errorf(`Post "https://location.aliyuncs.com": i/o timeout`)
		}
		return nil, context.DeadlineExceeded
	}

	client := &AliyunClient{
		RegionId: "cn-hangzhou",
		config: &config.Config{
			Endpoints:            new(sync.Map),
			RegionId:             "cn-hangzhou",
			Protocol:             "HTTP",
			LocationEndpoint:     "location.test",
			ClientReadTimeout:    10000,
			ClientConnectTimeout: 10000,
			AccessKey:            "test-ak",
			SecretKey:            "test-sk",
			Credential:           &mockCredential{err: fmt.Errorf("unit test uses the static access key")},
		},
	}
	ctx, cancel := context.WithTimeout(context.Background(), 4500*time.Millisecond)
	defer cancel()
	_, err := client.describeEndpointForService(ctx, "sls", client.RegionId)
	assert.Error(t, err)
	if !assert.Len(t, readTimeouts, 2, "the test must reach exactly one post-backoff attempt") {
		return
	}
	assert.Greater(t, readTimeouts[0], 4*time.Second)
	assert.Greater(t, readTimeouts[1], 500*time.Millisecond)
	assert.Less(t, readTimeouts[1], 2*time.Second)
	assert.Less(t, readTimeouts[1], readTimeouts[0], "the second attempt must use the reduced remaining budget")
}

// TestUnitEndpointDiscoveryStopsAtDeadline drives the real discovery path with an unreachable
// location service and an operation deadline: the first attempt fails with the read timeout, the
// pause is capped by the deadline and must report it, and no further connection may be opened after
// the call returns. A pause that silently succeeded would sleep the full three seconds and probe
// again, which both breaks the elapsed bound and grows the connection count.
func TestUnitEndpointDiscoveryStopsAtDeadline(t *testing.T) {
	var connections int32
	probe := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		<-r.Context().Done()
	}))
	probe.Config.ConnState = func(_ net.Conn, state http.ConnState) {
		if state == http.StateNew {
			atomic.AddInt32(&connections, 1)
		}
	}
	probe.Start()
	defer probe.Close()

	client := &AliyunClient{
		RegionId: "cn-hangzhou",
		config: &config.Config{
			Endpoints:            new(sync.Map),
			RegionId:             "cn-hangzhou",
			Protocol:             "HTTP",
			LocationEndpoint:     probe.Listener.Addr().String(),
			ClientReadTimeout:    100,
			ClientConnectTimeout: 100,
			AccessKey:            "test-ak",
			SecretKey:            "test-sk",
			Credential:           &mockCredential{err: fmt.Errorf("unit test uses the static access key")},
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 1200*time.Millisecond)
	defer cancel()

	started := time.Now()
	_, err := client.describeEndpointForService(ctx, "sls", client.RegionId)
	elapsed := time.Since(started)
	assert.Error(t, err)
	assert.GreaterOrEqual(t, elapsed, 800*time.Millisecond, "the call must reach the operation deadline")
	assert.Less(t, elapsed, 3*time.Second, "a deadline-capped pause must not sleep the full backoff and probe again")

	settled := atomic.LoadInt32(&connections)
	assert.NotZero(t, settled)
	time.Sleep(300 * time.Millisecond)
	assert.Equal(t, settled, atomic.LoadInt32(&connections), "no discovery attempt may run after the deadline")
}

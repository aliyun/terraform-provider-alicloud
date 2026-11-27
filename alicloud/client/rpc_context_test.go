package client

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/location"
	"github.com/aliyun/credentials-go/credentials"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/client/config"
)

// newContextRPCTestClient builds a client whose foasconsole endpoint is the given local stub, so the
// context-aware RPC path is exercised without real credentials or network access.
func newContextRPCTestClient(t *testing.T, endpoint string) *AliyunClient {
	t.Helper()
	return newContextRPCTestClientWithTimeouts(t, endpoint, 0, 0)
}

// newContextRPCTestClientWithTimeouts is newContextRPCTestClient with explicit client timeouts in
// milliseconds, matching the provider's client_connect_timeout / client_read_timeout wiring.
func newContextRPCTestClientWithTimeouts(t *testing.T, endpoint string, connectTimeoutMs int, readTimeoutMs int) *AliyunClient {
	t.Helper()
	t.Setenv("NO_PROXY", "127.0.0.1")
	t.Setenv("no_proxy", "127.0.0.1")
	credential, err := credentials.NewCredential(new(credentials.Config).SetType("access_key").
		SetAccessKeyId("test-access-key-id").
		SetAccessKeySecret("test-access-key-secret"))
	if err != nil {
		t.Fatalf("build credential failed: %v", err)
	}
	endpoints := new(sync.Map)
	cfg := config.Config{
		AccessKey:            "test-access-key-id",
		SecretKey:            "test-access-key-secret",
		Credential:           credential,
		RegionId:             "cn-hangzhou",
		AccountType:          "test",
		Protocol:             "http",
		SkipRegionValidation: true,
		Endpoints:            endpoints,
		SignVersion:          new(sync.Map),
		ClientConnectTimeout: connectTimeoutMs,
		ClientReadTimeout:    readTimeoutMs,
	}
	client, err := NewClient(&cfg)
	if err != nil {
		t.Fatalf("build client failed: %v", err)
	}
	// The context-aware RPC path resolves the host through the endpoint cache.
	endpoints.Store("foasconsole", strings.TrimPrefix(endpoint, "http://"))
	return client
}

// TestUnitRpcPostWithContextInFlightRequestHonoursBudget pins the HTTP-level half of the operation
// budget: a request that is already in flight has to stop at the remaining deadline or at a lifecycle
// cancellation. The retry loop alone only stops new attempts, so these cases point the client at a
// server that never answers.
func TestUnitRpcPostWithContextInFlightRequestHonoursBudget(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		<-r.Context().Done()
	}))
	t.Cleanup(func() {
		server.CloseClientConnections()
		server.Close()
	})
	client := newContextRPCTestClient(t, server.URL)

	t.Run("remaining deadline bounds the request", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 250*time.Millisecond)
		defer cancel()
		started := time.Now()
		_, err := client.RpcPostWithContext(ctx, "foasconsole", "2019-06-01", "CreateNamespace", nil, map[string]interface{}{})
		if err == nil {
			t.Fatalf("a hung request must not report success")
		}
		if !errors.Is(err, context.DeadlineExceeded) {
			t.Fatalf("error = %v, want context.DeadlineExceeded in its chain", err)
		}
		if elapsed := time.Since(started); elapsed > 3*time.Second {
			t.Fatalf("the hung request returned after %s; the remaining budget was not applied to the HTTP attempt", elapsed)
		}
	})

	t.Run("cancellation interrupts the request", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		go func() {
			time.Sleep(150 * time.Millisecond)
			cancel()
		}()
		started := time.Now()
		_, err := client.RpcPostWithContext(ctx, "foasconsole", "2019-06-01", "CreateNamespace", nil, map[string]interface{}{})
		if !errors.Is(err, context.Canceled) {
			t.Fatalf("error = %v, want context.Canceled", err)
		}
		if elapsed := time.Since(started); elapsed > 3*time.Second {
			t.Fatalf("the cancelled request returned after %s", elapsed)
		}
	})
}

func TestUnitRpcPostWithContextUncachedEndpointPreservesContextError(t *testing.T) {
	originalDescribe := describeLocationEndpoints
	t.Cleanup(func() { describeLocationEndpoints = originalDescribe })

	t.Run("cancellation", func(t *testing.T) {
		client := newContextRPCTestClient(t, "http://127.0.0.1:1")
		client.config.LocationEndpoint = "location.test"
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		describeLocationEndpoints = func(_ context.Context, _ *location.Client, _ *location.DescribeEndpointsRequest, _, _ time.Duration, _ *http.Transport, _ *url.URL) (*location.DescribeEndpointsResponse, error) {
			cancel()
			return nil, ctx.Err()
		}
		_, err := client.RpcPostWithContext(ctx, "vvpuncached", "2019-06-01", "DescribeNamespaces", nil, nil)
		if !errors.Is(err, context.Canceled) {
			t.Fatalf("error = %v, want context.Canceled in its chain", err)
		}
	})

	t.Run("deadline", func(t *testing.T) {
		client := newContextRPCTestClient(t, "http://127.0.0.1:1")
		client.config.LocationEndpoint = "location.test"
		ctx, cancel := context.WithTimeout(context.Background(), 25*time.Millisecond)
		defer cancel()
		describeLocationEndpoints = func(_ context.Context, _ *location.Client, _ *location.DescribeEndpointsRequest, _, _ time.Duration, _ *http.Transport, _ *url.URL) (*location.DescribeEndpointsResponse, error) {
			<-ctx.Done()
			return nil, ctx.Err()
		}
		_, err := client.RpcPostWithContext(ctx, "vvpuncached", "2019-06-01", "DescribeNamespaces", nil, nil)
		if !errors.Is(err, context.DeadlineExceeded) {
			t.Fatalf("error = %v, want context.DeadlineExceeded in its chain", err)
		}
	})
}

// TestUnitRpcPostWithContextCancelStopsFurtherRequests proves that lifecycle cancellation reaches the
// actual HTTP request: the test synchronizes on the request arriving at the server, cancels the
// context, and then waits for the server-side request context to end. It also asserts that no second
// request is issued, because the caller owns the only retry loop.
func TestUnitRpcPostWithContextCancelStopsFurtherRequests(t *testing.T) {
	var mu sync.Mutex
	requests := 0
	started := make(chan struct{}, 1)
	finished := make(chan struct{}, 1)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		requests++
		mu.Unlock()
		started <- struct{}{}
		<-r.Context().Done()
		finished <- struct{}{}
	}))
	t.Cleanup(func() {
		server.CloseClientConnections()
		server.Close()
	})
	client := newContextRPCTestClient(t, server.URL)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	result := make(chan error, 1)
	go func() {
		_, err := client.RpcPostWithContext(ctx, "foasconsole", "2019-06-01", "CreateNamespace", nil, map[string]interface{}{})
		result <- err
	}()
	select {
	case <-started:
	case <-time.After(5 * time.Second):
		t.Fatalf("the request never reached the server")
	}
	cancel()
	select {
	case err := <-result:
		if !errors.Is(err, context.Canceled) {
			t.Fatalf("error = %v, want context.Canceled", err)
		}
	case <-time.After(3 * time.Second):
		t.Fatalf("the cancelled call did not return")
	}
	select {
	case <-finished:
	case <-time.After(5 * time.Second):
		t.Fatalf("cancelling the lifecycle did not end the in-flight HTTP request")
	}
	time.Sleep(300 * time.Millisecond)
	mu.Lock()
	defer mu.Unlock()
	if requests != 1 {
		t.Fatalf("HTTP attempts = %d, want exactly 1 after cancellation (no further retry or write)", requests)
	}
}

// TestUnitRpcPostWithContextResponseBodyShape pins how a context-aware RPC response body is decoded.
// An empty body is a successful call without payload (the range of products that answer 200 with no
// body), while a body that is not a JSON object must fail with an error instead of panicking. Before
// this guard AssertAsMap panicked on nil and on every non-map value.
func TestUnitRpcPostWithContextResponseBodyShape(t *testing.T) {
	cases := []struct {
		name       string
		body       string
		wantErr    bool
		errContain string
	}{
		{name: "empty body is an empty result", body: "", wantErr: false},
		{name: "json null stays an empty result", body: "null", wantErr: false},
		{name: "json array is rejected", body: "[]", wantErr: true, errContain: "instead of a JSON object"},
		{name: "json string is rejected", body: `"ok"`, wantErr: true, errContain: "instead of a JSON object"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(tc.body))
			}))
			t.Cleanup(server.Close)
			client := newContextRPCTestClient(t, server.URL)

			result, err := client.RpcPostWithContext(context.Background(), "foasconsole", "2019-06-01", "DescribeNamespaces", nil, map[string]interface{}{})
			if tc.wantErr {
				if err == nil {
					t.Fatalf("body %q must report an error", tc.body)
				}
				if !strings.Contains(err.Error(), tc.errContain) {
					t.Fatalf("error = %v, want it to contain %q", err, tc.errContain)
				}
				return
			}
			if err != nil {
				t.Fatalf("body %q returned an unexpected error: %v", tc.body, err)
			}
			if result == nil {
				t.Fatalf("body %q returned a nil result map", tc.body)
			}
			if _, ok := result["_headers"]; !ok {
				t.Fatalf("body %q did not produce the response headers entry", tc.body)
			}
		})
	}
}

// TestUnitRpcPostWithContextErrorBodyStatusCode pins the error path for a business error that carries
// a JSON StatusCode. json.Unmarshal decodes the value as float64 because UseNumber is not enabled,
// and the shared formatError used to assert *int there, which panicked instead of returning the SDKError.
func TestUnitRpcPostWithContextErrorBodyStatusCode(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"Code":"InvalidNamespace","Message":"bad namespace","StatusCode":400}`))
	}))
	t.Cleanup(server.Close)
	client := newContextRPCTestClient(t, server.URL)

	_, err := client.RpcPostWithContext(context.Background(), "foasconsole", "2019-06-01", "CreateNamespace", nil, map[string]interface{}{})
	if err == nil {
		t.Fatalf("a business error body must report an error")
	}
	if !strings.Contains(err.Error(), "InvalidNamespace") {
		t.Fatalf("error = %v, want it to carry the product code", err)
	}
	if !strings.Contains(err.Error(), "400") {
		t.Fatalf("error = %v, want it to carry the body StatusCode", err)
	}
}

// TestUnitRpcPostWithContextUsesConfiguredReadTimeout proves the new entry honours the provider's
// client_read_timeout: with a 300ms read timeout a server that never sends a response must surface an
// error long before the 5s operation deadline, instead of waiting for the context budget.
func TestUnitRpcPostWithContextUsesConfiguredReadTimeout(t *testing.T) {
	release := make(chan struct{})
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		select {
		case <-release:
		case <-r.Context().Done():
		}
	}))
	t.Cleanup(func() {
		close(release)
		server.CloseClientConnections()
		server.Close()
	})
	client := newContextRPCTestClientWithTimeouts(t, server.URL, 2000, 300)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	started := time.Now()
	_, err := client.RpcPostWithContext(ctx, "foasconsole", "2019-06-01", "DescribeNamespaces", nil, map[string]interface{}{})
	elapsed := time.Since(started)
	if err == nil {
		t.Fatalf("a hung response must report an error")
	}
	if ctx.Err() != nil {
		t.Fatalf("the 5s operation context ended (%v); the configured read timeout must end the attempt first", ctx.Err())
	}
	if elapsed > 2*time.Second {
		t.Fatalf("the attempt returned after %s; client_read_timeout=300ms was not applied", elapsed)
	}
}

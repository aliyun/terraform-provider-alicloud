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

	rpc "github.com/alibabacloud-go/tea-rpc/client"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/credentials-go/credentials"
)

// newLocalRpcClient points the SLB endpoint of a client at a local test server
// so RpcPostWithContext exercises the real Set/Describe transport path without
// reaching the cloud.
func newLocalRpcClient(t *testing.T, endpoint string) *AliyunClient {
	t.Helper()
	credential, err := credentials.NewCredential(new(credentials.Config).
		SetType("access_key").
		SetAccessKeyId("unit-test-access-key").
		SetAccessKeySecret("unit-test-access-key-secret"))
	if err != nil {
		t.Fatalf("building the credential: %s", err)
	}
	var endpoints sync.Map
	endpoints.Store("slb", endpoint)
	return &AliyunClient{
		RegionId: "cn-beijing",
		config: &Config{
			RegionId:    "cn-beijing",
			Protocol:    "HTTP",
			Credential:  credential,
			Endpoints:   &endpoints,
			SignVersion: &sync.Map{},
		},
		teaSdkConfig: rpc.Config{Protocol: tea.String("HTTP")},
	}
}

func clearProxyEnv(t *testing.T) {
	t.Helper()
	for _, name := range []string{"HTTP_PROXY", "http_proxy", "HTTPS_PROXY", "https_proxy"} {
		t.Setenv(name, "")
	}
}

// The transport used by RpcPostWithContext must abort an in-flight request when
// the lifecycle context is cancelled, not only stop the outer waiter.
func TestUnitContextHTTPClientCancelsInFlightRequest(t *testing.T) {
	requestStarted := make(chan struct{})
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		close(requestStarted)
		<-r.Context().Done()
	}))
	defer server.Close()

	ctx, cancel := context.WithCancel(context.Background())
	client := contextHTTPClient{ctx: ctx, budget: 30 * time.Second}
	request, err := http.NewRequest(http.MethodGet, server.URL, nil)
	if err != nil {
		t.Fatalf("building the request: %s", err)
	}
	go func() {
		<-requestStarted
		cancel()
	}()

	start := time.Now()
	if _, err := client.Call(request, &http.Transport{}); err == nil {
		t.Fatal("expected the cancelled request to fail")
	}
	if elapsed := time.Since(start); elapsed > 5*time.Second {
		t.Fatalf("the cancelled request did not stop promptly: %s", elapsed)
	}
}

// A request is also bounded by the remaining budget even without an explicit
// cancellation.
func TestUnitContextHTTPClientHonoursBudget(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(2 * time.Second)
	}))
	defer server.Close()

	client := contextHTTPClient{ctx: context.Background(), budget: 100 * time.Millisecond}
	request, err := http.NewRequest(http.MethodGet, server.URL, nil)
	if err != nil {
		t.Fatalf("building the request: %s", err)
	}
	start := time.Now()
	if _, err := client.Call(request, &http.Transport{}); err == nil {
		t.Fatal("expected the request to exceed the budget")
	}
	if elapsed := time.Since(start); elapsed > time.Second {
		t.Fatalf("the budget was not enforced: %s", elapsed)
	}
}

// A response that arrives inside the remaining budget still succeeds even when
// the budget is short.
func TestUnitContextHTTPClientDelayedSuccessWithinBudget(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(50 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"ok":true}`))
	}))
	defer server.Close()

	client := contextHTTPClient{ctx: context.Background(), budget: 500 * time.Millisecond}
	request, err := http.NewRequest(http.MethodGet, server.URL, nil)
	if err != nil {
		t.Fatalf("building the request: %s", err)
	}
	response, err := client.Call(request, &http.Transport{})
	if err != nil {
		t.Fatalf("the delayed response should fit the budget: %s", err)
	}
	response.Body.Close()
}

// A fast dial must leave the whole remaining budget for the response headers.
// The header arrives after more than half of the budget, which the previous
// 50/50 split rejected.
func TestUnitContextRpcFastDialDelayedHeader(t *testing.T) {
	clearProxyEnv(t)
	const headerDelay = 700 * time.Millisecond
	var requests atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests.Add(1)
		time.Sleep(headerDelay)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"RequestId":"r-1","AclStatus":"on","ListenerPort":80}`))
	}))
	defer server.Close()

	client := newLocalRpcClient(t, strings.TrimPrefix(server.URL, "http://"))
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	start := time.Now()
	result, err := client.RpcPostWithContext(ctx, "Slb", "2014-05-15", "DescribeLoadBalancerHTTPListenerAttribute", nil, map[string]interface{}{
		"LoadBalancerId": "lb-unit-test",
		"ListenerPort":   80,
	}, false)
	elapsed := time.Since(start)
	if err != nil {
		t.Fatalf("a delayed header inside the budget must succeed: %s", err)
	}
	if got := requests.Load(); got != 1 {
		t.Fatalf("expected exactly one request, got %d", got)
	}
	if result["AclStatus"] != "on" {
		t.Fatalf("unexpected response: %v", result)
	}
	if elapsed >= time.Second {
		t.Fatalf("the request did not respect the total deadline: %s", elapsed)
	}
}

// Cancelling the lifecycle context must abort an in-flight Set or Describe
// request and prevent another attempt.
func TestUnitContextRpcCancellationStopsInFlightWrite(t *testing.T) {
	clearProxyEnv(t)
	for _, action := range []string{"SetLoadBalancerHTTPListenerAttribute", "DescribeLoadBalancerHTTPListenerAttribute"} {
		t.Run(action, func(t *testing.T) {
			arrived := make(chan struct{})
			release := make(chan struct{})
			var once sync.Once
			var requests atomic.Int32
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				requests.Add(1)
				once.Do(func() { close(arrived) })
				// The handler must never outlive the test, otherwise Close waits
				// for the connection instead of reporting the assertion result.
				select {
				case <-r.Context().Done():
				case <-release:
				case <-time.After(3 * time.Second):
				}
			}))
			defer server.Close()
			defer close(release)

			client := newLocalRpcClient(t, strings.TrimPrefix(server.URL, "http://"))
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			go func() {
				<-arrived
				cancel()
			}()

			start := time.Now()
			_, err := client.RpcPostWithContext(ctx, "Slb", "2014-05-15", action, nil, map[string]interface{}{
				"LoadBalancerId": "lb-unit-test",
				"ListenerPort":   80,
				"AclStatus":      "on",
				"AclId":          "acl-unit-test",
				"AclType":        "white",
			}, false)
			if err == nil {
				t.Fatal("expected the cancelled request to fail")
			}
			if elapsed := time.Since(start); elapsed > 2*time.Second {
				t.Fatalf("the cancelled request did not stop promptly: %s", elapsed)
			}
			// Give a runaway retry a chance to show up before asserting the count.
			time.Sleep(100 * time.Millisecond)
			if got := requests.Load(); got != 1 {
				t.Fatalf("a cancelled request must not be retried, got %d requests", got)
			}
		})
	}
}

// A malformed endpoint must fail before tea builds the signed request URL, so
// neither the SDK debug output nor the parse error can echo the credentials.
func TestUnitContextRpcMalformedEndpointIsRedacted(t *testing.T) {
	clearProxyEnv(t)
	var requests atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests.Add(1)
		_, _ = w.Write([]byte(`{"RequestId":"r-1"}`))
	}))
	defer server.Close()

	endpoint := strings.TrimPrefix(server.URL, "http://") + "/%zz"
	client := newLocalRpcClient(t, endpoint)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	_, err := client.RpcPostWithContext(ctx, "Slb", "2014-05-15", "SetLoadBalancerHTTPListenerAttribute", nil, map[string]interface{}{
		"LoadBalancerId": "lb-unit-test",
		"ListenerPort":   80,
		"AclStatus":      "on",
		"AclId":          "acl-unit-test",
		"AclType":        "white",
	}, false)
	if err == nil {
		t.Fatal("expected the malformed endpoint to fail")
	}
	if got := requests.Load(); got != 0 {
		t.Fatalf("no request may be built or sent for a malformed endpoint, got %d", got)
	}
	message := err.Error()
	for _, secret := range []string{"AccessKeyId", "Signature=", "SignatureNonce", "SecurityToken", "unit-test-access-key"} {
		if strings.Contains(message, secret) {
			t.Fatalf("the endpoint error exposes %q: %s", secret, message)
		}
	}
	if !strings.Contains(message, "invalid") || !strings.Contains(message, "endpoint") {
		t.Fatalf("the endpoint error lost its classification: %s", message)
	}
}

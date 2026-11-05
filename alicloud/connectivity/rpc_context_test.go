package connectivity

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/aliyun/credentials-go/credentials"
)

// newContextRPCTestClient builds a client whose foasconsole endpoint is the given local stub, so the
// context-aware RPC path is exercised without real credentials or network access.
func newContextRPCTestClient(t *testing.T, endpoint string) *AliyunClient {
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
	config := Config{
		AccessKey:            "test-access-key-id",
		SecretKey:            "test-access-key-secret",
		Credential:           credential,
		RegionId:             "cn-hangzhou",
		AccountType:          "test",
		Protocol:             "http",
		SkipRegionValidation: true,
		Endpoints:            endpoints,
		SignVersion:          new(sync.Map),
	}
	client, err := config.Client()
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

package connectivity

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

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

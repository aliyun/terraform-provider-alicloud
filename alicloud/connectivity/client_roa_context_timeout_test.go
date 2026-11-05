package connectivity

import (
	"context"
	"testing"
	"time"

	roa "github.com/alibabacloud-go/tea-roa/client"
	"github.com/alibabacloud-go/tea/tea"
)

// TestUnitRoaContextTimeout proves the HTTP budget derived from an operation
// deadline: it is rounded up to whole seconds so the SDK keeps reusing one
// pooled client, it is capped by the configured connect and read timeouts
// together, and an expired deadline is rejected before any request is built.
func TestUnitRoaContextTimeout(t *testing.T) {
	configured := roa.Config{}
	configured.SetConnectTimeout(60000).SetReadTimeout(60000)
	unbounded := roa.Config{}

	cases := map[string]struct {
		deadline time.Time
		config   roa.Config
		want     int
		wantErr  bool
	}{
		"sub-second-budget-rounds-up":   {time.Now().Add(400 * time.Millisecond), unbounded, 1000, false},
		"whole-second-budget":           {time.Now().Add(2 * time.Second), unbounded, 2000, false},
		"below-configured-total":        {time.Now().Add(90 * time.Second), configured, 90000, false},
		"capped-by-configured-total":    {time.Now().Add(5 * time.Minute), configured, 120000, false},
		"no-configured-timeout-is-open": {time.Now().Add(5 * time.Minute), unbounded, 300000, false},
		"expired-deadline":              {time.Now().Add(-time.Second), configured, 0, true},
	}
	for name, testCase := range cases {
		t.Run(name, func(t *testing.T) {
			got, err := roaContextTimeout(testCase.deadline, testCase.config)
			if testCase.wantErr {
				if err == nil {
					t.Fatalf("roaContextTimeout returned %d, want an error for an expired deadline", got)
				}
				return
			}
			if err != nil {
				t.Fatalf("roaContextTimeout returned error: %v", err)
			}
			// The rounding is allowed to add at most one second; the cap must not
			// be exceeded.
			if got < testCase.want || got > testCase.want+1000 {
				t.Fatalf("roaContextTimeout = %d, want %d (plus at most one second of rounding)", got, testCase.want)
			}
		})
	}
}

// TestUnitRoaContextTimeouts proves the split of the request budget into the
// connect and read timeouts: the connect share comes from the configured value,
// is never larger than the budget, and the two shares add up to the budget.
func TestUnitRoaContextTimeouts(t *testing.T) {
	configured := roa.Config{}
	configured.SetConnectTimeout(60000).SetReadTimeout(60000)
	unbounded := roa.Config{}

	cases := map[string]struct {
		deadline    time.Time
		config      roa.Config
		wantConnect int
		wantRead    int
		wantErr     bool
	}{
		"connect-kept-below-budget":   {time.Now().Add(5 * time.Minute), configured, 60000, 60000, false},
		"connect-trimmed-to-budget":   {time.Now().Add(90 * time.Second), configured, 60000, 30000, false},
		"connect-capped-by-budget":    {time.Now().Add(30 * time.Second), configured, 30000, 0, false},
		"no-configured-connect":       {time.Now().Add(400 * time.Millisecond), unbounded, 0, 1000, false},
		"expired-deadline-is-refused": {time.Now().Add(-time.Second), configured, 0, 0, true},
	}
	for name, testCase := range cases {
		t.Run(name, func(t *testing.T) {
			connect, read, err := roaContextTimeouts(testCase.deadline, testCase.config)
			if testCase.wantErr {
				if err == nil {
					t.Fatalf("roaContextTimeouts returned %d/%d, want an error", connect, read)
				}
				return
			}
			if err != nil {
				t.Fatalf("roaContextTimeouts returned error: %v", err)
			}
			if connect != testCase.wantConnect || read < testCase.wantRead || read > testCase.wantRead+1000 {
				t.Fatalf("roaContextTimeouts = %d/%d, want %d/%d (plus at most one second of rounding)", connect, read, testCase.wantConnect, testCase.wantRead)
			}
			if connect < 0 || read < 0 {
				t.Fatalf("negative timeout shares: connect=%d read=%d", connect, read)
			}
		})
	}
}

// TestUnitRoaContextRuntime proves a bounded runtime is built only when the
// context carries a deadline, and that it keeps the automatic retry flag.
func TestUnitRoaContextRuntime(t *testing.T) {
	unbounded := roa.Config{}
	runtime, err := roaContextRuntime(context.Background(), unbounded, true)
	if err != nil {
		t.Fatalf("roaContextRuntime returned error: %v", err)
	}
	if runtime != nil {
		t.Fatalf("roaContextRuntime built a bounded runtime for a context without a deadline")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	runtime, err = roaContextRuntime(ctx, unbounded, false)
	if err != nil {
		t.Fatalf("roaContextRuntime returned error: %v", err)
	}
	if runtime == nil {
		t.Fatal("roaContextRuntime did not bound a context with a deadline")
	}
	if got := tea.IntValue(runtime.ReadTimeout); got < 3000 || got > 4000 {
		t.Fatalf("read timeout = %d, want the three-second budget rounded up to a whole second", got)
	}
	if got := tea.IntValue(runtime.ConnectTimeout); got != 0 {
		t.Fatalf("connect timeout = %d, want the unbounded client value", got)
	}
	if got := tea.BoolValue(runtime.Autoretry); got {
		t.Fatalf("autoretry = %v, want the caller's false to be kept", got)
	}
}

// TestUnitRoaWithContextRejectsExpiredDeadline proves an expired operation
// stops before any request is issued.
func TestUnitRoaWithContextRejectsExpiredDeadline(t *testing.T) {
	client := &AliyunClient{}
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(-time.Second))
	defer cancel()

	if _, err := client.RoaWithContext(ctx, "GET", "Cms", "2024-03-30", "/subscriptions", nil, nil, nil, false); err == nil {
		t.Fatal("RoaWithContext accepted a request whose operation deadline had already passed")
	}
}

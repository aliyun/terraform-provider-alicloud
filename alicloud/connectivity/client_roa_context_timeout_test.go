package connectivity

import (
	"context"
	"testing"
	"time"

	roa "github.com/alibabacloud-go/tea-roa/client"
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

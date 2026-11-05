package rds

import (
	"strings"
	"testing"
)

// The page collector is an unexported helper, so this boundary case lives in
// an internal test; the external tests exercise it end to end through the
// mock HTTP client.
func TestUnitRdsDatabasesCollectPagesStopsAtLimit(t *testing.T) {
	calls := 0
	_, err := databasesCollectPages(3, 100, func(pageNumber int) ([]map[string]interface{}, error) {
		calls++
		items := make([]map[string]interface{}, 100)
		for i := range items {
			items[i] = map[string]interface{}{"DBName": "db"}
		}
		return items, nil
	})
	if err == nil {
		t.Fatal("expected the collector to reject a pagination loop that never shrinks")
	}
	if calls != 3 {
		t.Fatalf("expected the collector to stop after 3 pages, got %d", calls)
	}
	if !strings.Contains(err.Error(), "pagination did not terminate") {
		t.Fatalf("unexpected error: %v", err)
	}
}

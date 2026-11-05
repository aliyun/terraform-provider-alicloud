package rocketmq

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"testing"
)

func rocketmqTestPage(prefix string, count int) []map[string]interface{} {
	items := make([]map[string]interface{}, 0, count)
	for i := 0; i < count; i++ {
		items = append(items, map[string]interface{}{"instanceId": fmt.Sprintf("%s-%d", prefix, i)})
	}
	return items
}

func TestUnitRocketmqCollectPages(t *testing.T) {
	key := func(item map[string]interface{}) string { return rocketmqString(item, "instanceId") }

	t.Run("short page ends the walk", func(t *testing.T) {
		calls := 0
		items, err := rocketmqCollectPages("/instances", key, func(pageNumber int) ([]map[string]interface{}, error) {
			calls++
			return rocketmqTestPage("a", 2), nil
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if calls != 1 {
			t.Fatalf("expected 1 request, got %d", calls)
		}
		if len(items) != 2 {
			t.Fatalf("expected 2 items, got %d", len(items))
		}
	})

	t.Run("exact page boundary reads one more page", func(t *testing.T) {
		calls := 0
		items, err := rocketmqCollectPages("/instances", key, func(pageNumber int) ([]map[string]interface{}, error) {
			calls++
			if pageNumber == 1 {
				return rocketmqTestPage("a", rocketmqDataSourcePageSize), nil
			}
			return rocketmqTestPage("b", 1), nil
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if calls != 2 {
			t.Fatalf("expected 2 requests, got %d", calls)
		}
		if len(items) != rocketmqDataSourcePageSize+1 {
			t.Fatalf("expected %d items, got %d", rocketmqDataSourcePageSize+1, len(items))
		}
	})

	t.Run("duplicates across pages are dropped without ending the walk early", func(t *testing.T) {
		calls := 0
		items, err := rocketmqCollectPages("/instances", key, func(pageNumber int) ([]map[string]interface{}, error) {
			calls++
			switch pageNumber {
			case 1:
				return rocketmqTestPage("a", rocketmqDataSourcePageSize), nil
			case 2:
				// A full page that repeats 99 known rows and adds one new row.
				// The walk must continue because the raw page is full, even
				// though only one row is new.
				page := rocketmqTestPage("a", rocketmqDataSourcePageSize-1)
				return append(page, map[string]interface{}{"instanceId": "new-2"}), nil
			case 3:
				return []map[string]interface{}{{"instanceId": "new-3"}}, nil
			}
			t.Fatalf("unexpected page %d", pageNumber)
			return nil, nil
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if calls != 3 {
			t.Fatalf("expected 3 requests, got %d", calls)
		}
		if len(items) != rocketmqDataSourcePageSize+2 {
			t.Fatalf("expected %d items after de-duplication, got %d", rocketmqDataSourcePageSize+2, len(items))
		}
	})

	t.Run("fetch error is propagated", func(t *testing.T) {
		_, err := rocketmqCollectPages("/instances", key, func(pageNumber int) ([]map[string]interface{}, error) {
			return nil, fmt.Errorf("boom")
		})
		if err == nil || !strings.Contains(err.Error(), "boom") {
			t.Fatalf("expected the fetch error, got %v", err)
		}
	})

	t.Run("a page that never shrinks is an error, not a partial result", func(t *testing.T) {
		_, err := rocketmqCollectPages("/instances", key, func(pageNumber int) ([]map[string]interface{}, error) {
			return rocketmqTestPage(fmt.Sprintf("p%d", pageNumber), rocketmqDataSourcePageSize), nil
		})
		if err == nil || !strings.Contains(err.Error(), "did not terminate") {
			t.Fatalf("expected a non-termination error, got %v", err)
		}
	})
}

func TestUnitRocketmqListItems(t *testing.T) {
	t.Run("data.list is parsed", func(t *testing.T) {
		response := map[string]interface{}{
			"success": true,
			"data": map[string]interface{}{
				"list": []interface{}{
					map[string]interface{}{"instanceId": "rmq-1"},
					"not-an-object",
					map[string]interface{}{"instanceId": "rmq-2"},
				},
			},
		}
		items, err := rocketmqListItems(response)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(items) != 2 {
			t.Fatalf("expected 2 items, got %d", len(items))
		}
	})

	t.Run("missing list is an empty result", func(t *testing.T) {
		items, err := rocketmqListItems(map[string]interface{}{"success": true, "data": map[string]interface{}{}})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(items) != 0 {
			t.Fatalf("expected no items, got %d", len(items))
		}
	})

	t.Run("unsuccessful response is an error", func(t *testing.T) {
		_, err := rocketmqListItems(map[string]interface{}{
			"success": false,
			"code":    "InvalidParameter",
			"message": "bad request",
		})
		if err == nil || !strings.Contains(err.Error(), "InvalidParameter") {
			t.Fatalf("expected an error naming the code, got %v", err)
		}
	})
}

func TestUnitRocketmqInt64(t *testing.T) {
	cases := []struct {
		name  string
		value interface{}
		want  int64
	}{
		{"float64", float64(20), 20},
		{"json number", json.Number("1500"), 1500},
		{"int", 7, 7},
		{"int64", int64(9), 9},
		{"missing", nil, 0},
		{"unknown type", "20", 0},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			item := map[string]interface{}{"count": tc.value}
			if tc.value == nil {
				item = map[string]interface{}{}
			}
			got := rocketmqInt64(item, "count")
			if got.ValueInt64() != tc.want {
				t.Fatalf("expected %d, got %d", tc.want, got.ValueInt64())
			}
		})
	}
}

func TestUnitRocketmqTags(t *testing.T) {
	if got := rocketmqTagsRequestValue(map[string]string{"b": "2", "a": "1"}); got != `[{"key":"a","value":"1"},{"key":"b","value":"2"}]` {
		t.Fatalf("unexpected request value: %s", got)
	}

	item := map[string]interface{}{
		"tags": []interface{}{
			map[string]interface{}{"key": "a", "value": "1"},
			map[string]interface{}{"key": "b", "value": "2"},
		},
	}
	if !rocketmqTagsMatch(item, map[string]string{"a": "1"}) {
		t.Fatal("expected the configured tag to match")
	}
	if rocketmqTagsMatch(item, map[string]string{"a": "2"}) {
		t.Fatal("expected a different tag value to be rejected")
	}
	if rocketmqTagsMatch(item, map[string]string{"c": "1"}) {
		t.Fatal("expected a missing tag to be rejected")
	}
	if !rocketmqTagsMatch(map[string]interface{}{}, nil) {
		t.Fatal("expected an empty filter to match everything")
	}

	tags := rocketmqTagsFromAPI(item["tags"])
	if len(tags) != 2 {
		t.Fatalf("expected 2 tags, got %d", len(tags))
	}
	if !reflect.DeepEqual([]string{tags[0].Key.ValueString(), tags[1].Key.ValueString()}, []string{"a", "b"}) {
		t.Fatalf("expected tags sorted by key, got %v", tags)
	}
}

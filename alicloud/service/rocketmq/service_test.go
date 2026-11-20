package rocketmq

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
)

func rocketmqTestPage(prefix string, count, pageNumber, totalCount int) rocketmqPage {
	items := make([]map[string]interface{}, 0, count)
	for i := 0; i < count; i++ {
		items = append(items, map[string]interface{}{"instanceId": fmt.Sprintf("%s-%d", prefix, i)})
	}
	return rocketmqPage{
		Items:      items,
		PageNumber: int64(pageNumber),
		PageSize:   int64(helper.PageSizeXLarge),
		TotalCount: int64(totalCount),
	}
}

func TestUnitRocketmqCollectPages(t *testing.T) {
	key := func(item map[string]interface{}) string { return rocketmqString(item, "instanceId") }

	t.Run("short page ends the walk", func(t *testing.T) {
		calls := 0
		items, err := rocketmqCollectPages("/instances", key, func(pageNumber int) (rocketmqPage, error) {
			calls++
			return rocketmqTestPage("a", 2, pageNumber, 2), nil
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if calls != 1 || len(items) != 2 {
			t.Fatalf("expected one request and two items, got calls=%d items=%d", calls, len(items))
		}
	})

	t.Run("exact page boundary reads one more page", func(t *testing.T) {
		calls := 0
		items, err := rocketmqCollectPages("/instances", key, func(pageNumber int) (rocketmqPage, error) {
			calls++
			if pageNumber == 1 {
				return rocketmqTestPage("a", helper.PageSizeXLarge, pageNumber, helper.PageSizeXLarge+1), nil
			}
			return rocketmqTestPage("b", 1, pageNumber, helper.PageSizeXLarge+1), nil
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if calls != 2 || len(items) != helper.PageSizeXLarge+1 {
			t.Fatalf("expected two requests and %d items, got calls=%d items=%d", helper.PageSizeXLarge+1, calls, len(items))
		}
	})

	t.Run("actual page size continues while total count has more rows", func(t *testing.T) {
		calls := 0
		items, err := rocketmqCollectPages("/instances", key, func(pageNumber int) (rocketmqPage, error) {
			calls++
			switch pageNumber {
			case 1:
				return rocketmqPage{Items: rocketmqTestPage("a", 10, 1, 20).Items, PageNumber: 1, PageSize: 10, TotalCount: 20}, nil
			case 2:
				return rocketmqPage{Items: rocketmqTestPage("b", 10, 2, 20).Items, PageNumber: 2, PageSize: 10, TotalCount: 20}, nil
			case 3:
				return rocketmqPage{Items: []map[string]interface{}{}, PageNumber: 3, PageSize: 10, TotalCount: 20}, nil
			}
			t.Fatalf("unexpected page %d", pageNumber)
			return rocketmqPage{}, nil
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if calls != 3 || len(items) != 20 {
			t.Fatalf("expected all 20 rows from two actual pages, got calls=%d items=%d", calls, len(items))
		}
	})

	t.Run("short page before total count is an error", func(t *testing.T) {
		items, err := rocketmqCollectPages("/instances", key, func(pageNumber int) (rocketmqPage, error) {
			return rocketmqPage{Items: rocketmqTestPage("a", 10, 1, 20).Items, PageNumber: 1, PageSize: helper.PageSizeXLarge, TotalCount: 20}, nil
		})
		if err == nil || !strings.Contains(err.Error(), "10 unique of 20") || len(items) != 0 {
			t.Fatalf("expected an incomplete-result error, got items=%d err=%v", len(items), err)
		}
	})

	t.Run("duplicate rows cannot satisfy total count", func(t *testing.T) {
		calls := 0
		items, err := rocketmqCollectPages("/instances", key, func(pageNumber int) (rocketmqPage, error) {
			calls++
			switch pageNumber {
			case 1, 2:
				return rocketmqPage{Items: rocketmqTestPage("a", 10, 1, 20).Items, PageNumber: int64(pageNumber), PageSize: 10, TotalCount: 20}, nil
			case 3:
				return rocketmqPage{Items: []map[string]interface{}{}, PageNumber: 3, PageSize: 10, TotalCount: 20}, nil
			}
			t.Fatalf("unexpected page %d", pageNumber)
			return rocketmqPage{}, nil
		})
		if err == nil || !strings.Contains(err.Error(), "10 unique of 20") || len(items) != 0 {
			t.Fatalf("expected duplicates to fail unique completeness, got calls=%d items=%d err=%v", calls, len(items), err)
		}
	})

	t.Run("duplicates across pages are dropped without ending the walk early", func(t *testing.T) {
		calls := 0
		items, err := rocketmqCollectPages("/instances", key, func(pageNumber int) (rocketmqPage, error) {
			calls++
			switch pageNumber {
			case 1:
				return rocketmqTestPage("a", helper.PageSizeXLarge, pageNumber, helper.PageSizeXLarge+2), nil
			case 2:
				page := rocketmqTestPage("a", helper.PageSizeXLarge-1, pageNumber, helper.PageSizeXLarge+2)
				page.Items = append(page.Items, map[string]interface{}{"instanceId": "new-2"})
				return page, nil
			case 3:
				return rocketmqPage{Items: []map[string]interface{}{{"instanceId": "new-3"}}, PageNumber: 3, PageSize: helper.PageSizeXLarge, TotalCount: helper.PageSizeXLarge + 2}, nil
			}
			t.Fatalf("unexpected page %d", pageNumber)
			return rocketmqPage{}, nil
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if calls != 3 || len(items) != helper.PageSizeXLarge+2 {
			t.Fatalf("expected %d unique items, got calls=%d items=%d", helper.PageSizeXLarge+2, calls, len(items))
		}
	})

	t.Run("empty primary key is an error naming its page", func(t *testing.T) {
		calls := 0
		items, err := rocketmqCollectPages("/instances", key, func(pageNumber int) (rocketmqPage, error) {
			calls++
			if pageNumber == 1 {
				return rocketmqTestPage("a", helper.PageSizeXLarge, pageNumber, helper.PageSizeXLarge+1), nil
			}
			return rocketmqPage{Items: []map[string]interface{}{{"instanceName": "missing-primary-key"}}, PageNumber: 2, PageSize: helper.PageSizeXLarge, TotalCount: helper.PageSizeXLarge + 1}, nil
		})
		if err == nil || !strings.Contains(err.Error(), "empty primary key") || !strings.Contains(err.Error(), "page 2") {
			t.Fatalf("expected an empty-primary-key error naming page 2, got items=%d err=%v", len(items), err)
		}
		if calls != 2 {
			t.Fatalf("expected 2 requests, got %d", calls)
		}
	})

	t.Run("fetch error is propagated", func(t *testing.T) {
		_, err := rocketmqCollectPages("/instances", key, func(pageNumber int) (rocketmqPage, error) {
			return rocketmqPage{}, fmt.Errorf("boom")
		})
		if err == nil || !strings.Contains(err.Error(), "boom") {
			t.Fatalf("expected the fetch error, got %v", err)
		}
	})

	t.Run("response page number must match the request", func(t *testing.T) {
		_, err := rocketmqCollectPages("/instances", key, func(pageNumber int) (rocketmqPage, error) {
			return rocketmqTestPage("a", 1, pageNumber+1, 1), nil
		})
		if err == nil || !strings.Contains(err.Error(), "pagination is inconsistent") {
			t.Fatalf("expected a page-number mismatch, got %v", err)
		}
	})

	t.Run("a page that never shrinks is an error, not a partial result", func(t *testing.T) {
		_, err := rocketmqCollectPages("/instances", key, func(pageNumber int) (rocketmqPage, error) {
			return rocketmqTestPage(fmt.Sprintf("p%d", pageNumber), helper.PageSizeXLarge, pageNumber, 0), nil
		})
		if err == nil || !strings.Contains(err.Error(), "did not terminate") {
			t.Fatalf("expected a non-termination error, got %v", err)
		}
	})

	t.Run("a malformed full page is an error, not a short final page", func(t *testing.T) {
		page := make([]interface{}, helper.PageSizeXLarge)
		for i := range page {
			page[i] = map[string]interface{}{"instanceId": fmt.Sprintf("id-%d", i)}
		}
		page[len(page)-1] = "not-an-object"
		response := map[string]interface{}{"success": true, "data": map[string]interface{}{
			"pageNumber": 1, "pageSize": helper.PageSizeXLarge, "totalCount": helper.PageSizeXLarge, "list": page,
		}}

		items, err := rocketmqCollectPages("/instances", key, func(pageNumber int) (rocketmqPage, error) {
			return rocketmqListPage(response)
		})
		if err == nil || !strings.Contains(err.Error(), fmt.Sprintf("$.data.list[%d]", len(page)-1)) {
			t.Fatalf("expected the malformed element to fail the walk, got items=%d err=%v", len(items), err)
		}
	})
}

func TestUnitRocketmqListPage(t *testing.T) {
	t.Run("data.list and pagination metadata are parsed", func(t *testing.T) {
		response := map[string]interface{}{
			"success": true,
			"data": map[string]interface{}{
				"pageNumber": 1,
				"pageSize":   2,
				"totalCount": 2,
				"list": []interface{}{
					map[string]interface{}{"instanceId": "rmq-1"},
					map[string]interface{}{"instanceId": "rmq-2"},
				},
			},
		}
		page, err := rocketmqListPage(response)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(page.Items) != 2 || page.PageNumber != 1 || page.PageSize != 2 || page.TotalCount != 2 {
			t.Fatalf("unexpected page: %+v", page)
		}
	})

	t.Run("non-object element is an error naming its index", func(t *testing.T) {
		response := map[string]interface{}{
			"success": true,
			"data": map[string]interface{}{
				"pageNumber": 1,
				"pageSize":   2,
				"totalCount": 2,
				"list": []interface{}{
					map[string]interface{}{"instanceId": "rmq-1"},
					"not-an-object",
				},
			},
		}
		_, err := rocketmqListPage(response)
		if err == nil || !strings.Contains(err.Error(), "$.data.list[1]") {
			t.Fatalf("expected an element error naming index 1, got %v", err)
		}
	})

	t.Run("empty list is an empty result", func(t *testing.T) {
		page, err := rocketmqListPage(map[string]interface{}{"success": true, "data": map[string]interface{}{
			"pageNumber": 1, "pageSize": 10, "totalCount": 0, "list": []interface{}{},
		}})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(page.Items) != 0 {
			t.Fatalf("expected no items, got %d", len(page.Items))
		}
	})

	t.Run("null list is an error", func(t *testing.T) {
		_, err := rocketmqListPage(map[string]interface{}{"success": true, "data": map[string]interface{}{
			"pageNumber": 1, "pageSize": 10, "totalCount": 0, "list": nil,
		}})
		if err == nil || !strings.Contains(err.Error(), "null") {
			t.Fatalf("expected a null-list error, got %v", err)
		}
	})

	t.Run("missing list with a zero count is an empty result", func(t *testing.T) {
		page, err := rocketmqListPage(map[string]interface{}{
			"success": true,
			"data":    map[string]interface{}{"pageNumber": float64(1), "pageSize": float64(10), "totalCount": float64(0)},
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(page.Items) != 0 {
			t.Fatalf("expected no items, got %d", len(page.Items))
		}
	})

	t.Run("missing list with rows reported is an error", func(t *testing.T) {
		_, err := rocketmqListPage(map[string]interface{}{
			"success": true,
			"data":    map[string]interface{}{"pageNumber": 1, "pageSize": 10, "totalCount": float64(3)},
		})
		if err == nil || !strings.Contains(err.Error(), "$.data.list") {
			t.Fatalf("expected a parse error naming $.data.list, got %v", err)
		}
	})

	t.Run("missing pagination metadata is an error", func(t *testing.T) {
		_, err := rocketmqListPage(map[string]interface{}{"success": true, "data": map[string]interface{}{"list": []interface{}{}}})
		if err == nil || !strings.Contains(err.Error(), "$.data.pageNumber") {
			t.Fatalf("expected a pageNumber metadata error, got %v", err)
		}
	})

	t.Run("non-integer pagination metadata is an error", func(t *testing.T) {
		_, err := rocketmqListPage(map[string]interface{}{"success": true, "data": map[string]interface{}{
			"pageNumber": 1, "pageSize": 10, "totalCount": 1.5, "list": []interface{}{},
		}})
		if err == nil || !strings.Contains(err.Error(), "$.data.totalCount") {
			t.Fatalf("expected a totalCount metadata error, got %v", err)
		}
	})

	t.Run("list under the wrong type is an error", func(t *testing.T) {
		_, err := rocketmqListPage(map[string]interface{}{"success": true, "data": map[string]interface{}{
			"pageNumber": 1, "pageSize": 10, "totalCount": 0, "list": "not-an-array",
		}})
		if err == nil || !strings.Contains(err.Error(), "instead of an array") {
			t.Fatalf("expected a type error naming the list, got %v", err)
		}
	})

	t.Run("empty response body is an error", func(t *testing.T) {
		_, err := rocketmqListPage(nil)
		if err == nil || !strings.Contains(err.Error(), "empty response body") {
			t.Fatalf("expected an empty-body error, got %v", err)
		}
	})

	t.Run("unsuccessful response is an error", func(t *testing.T) {
		_, err := rocketmqListPage(map[string]interface{}{
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
	valid := []struct {
		name  string
		value interface{}
		want  int64
	}{
		{"float64", float64(20), 20},
		{"json number", json.Number("1500"), 1500},
		{"int", 7, 7},
		{"int64", int64(9), 9},
		{"zero", json.Number("0"), 0},
	}
	for _, tc := range valid {
		t.Run(tc.name, func(t *testing.T) {
			got, err := rocketmqInt64(map[string]interface{}{"count": tc.value}, "count")
			if err != nil || got.ValueInt64() != tc.want {
				t.Fatalf("expected %d, got %d with err=%v", tc.want, got.ValueInt64(), err)
			}
		})
	}

	invalid := []struct {
		name  string
		item  map[string]interface{}
		field string
	}{
		{"missing", map[string]interface{}{}, "count"},
		{"null", map[string]interface{}{"count": nil}, "count"},
		{"string", map[string]interface{}{"count": "20"}, "count"},
		{"fractional", map[string]interface{}{"count": 1.5}, "count"},
		{"invalid json number", map[string]interface{}{"count": json.Number("1.5")}, "count"},
	}
	for _, tc := range invalid {
		t.Run(tc.name, func(t *testing.T) {
			got, err := rocketmqInt64(tc.item, tc.field)
			if err == nil || !strings.Contains(err.Error(), "count") || !got.IsNull() {
				t.Fatalf("expected a field-specific error and null value, got value=%v err=%v", got, err)
			}
		})
	}
}

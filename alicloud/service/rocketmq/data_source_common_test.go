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

	t.Run("empty list is an empty result", func(t *testing.T) {
		items, err := rocketmqListItems(map[string]interface{}{"success": true, "data": map[string]interface{}{"list": []interface{}{}}})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(items) != 0 {
			t.Fatalf("expected no items, got %d", len(items))
		}
	})

	t.Run("null list is an empty result", func(t *testing.T) {
		items, err := rocketmqListItems(map[string]interface{}{"success": true, "data": map[string]interface{}{"list": nil}})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(items) != 0 {
			t.Fatalf("expected no items, got %d", len(items))
		}
	})

	t.Run("missing list with a zero count is an empty result", func(t *testing.T) {
		// ListInstances answers an empty result without a list key at all.
		items, err := rocketmqListItems(map[string]interface{}{
			"success": true,
			"data":    map[string]interface{}{"pageNumber": float64(1), "pageSize": float64(10), "totalCount": float64(0)},
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(items) != 0 {
			t.Fatalf("expected no items, got %d", len(items))
		}
	})

	t.Run("missing list with rows reported is an error", func(t *testing.T) {
		_, err := rocketmqListItems(map[string]interface{}{
			"success": true,
			"data":    map[string]interface{}{"totalCount": float64(3)},
		})
		if err == nil || !strings.Contains(err.Error(), "$.data.list") {
			t.Fatalf("expected a parse error naming $.data.list, got %v", err)
		}
	})

	t.Run("missing list without a count is an error", func(t *testing.T) {
		_, err := rocketmqListItems(map[string]interface{}{"success": true, "data": map[string]interface{}{}})
		if err == nil || !strings.Contains(err.Error(), "$.data.list") {
			t.Fatalf("expected a parse error naming $.data.list, got %v", err)
		}
	})

	t.Run("list under the wrong type is an error", func(t *testing.T) {
		_, err := rocketmqListItems(map[string]interface{}{"success": true, "data": map[string]interface{}{"list": "not-an-array"}})
		if err == nil || !strings.Contains(err.Error(), "instead of an array") {
			t.Fatalf("expected a type error naming the list, got %v", err)
		}
	})

	t.Run("empty response body is an error", func(t *testing.T) {
		_, err := rocketmqListItems(nil)
		if err == nil || !strings.Contains(err.Error(), "empty response body") {
			t.Fatalf("expected an empty-body error, got %v", err)
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

func TestUnitRocketmqInstanceMapping(t *testing.T) {
	item := map[string]interface{}{
		"instanceId":      "rmq-1",
		"instanceName":    "name-1",
		"status":          "RUNNING",
		"regionId":        "cn-hangzhou",
		"seriesCode":      "ultimate",
		"subSeriesCode":   "cluster_ha",
		"paymentType":     "PayAsYouGo",
		"commodityCode":   "ons_rmqsub_public_cn",
		"serviceCode":     "rmq",
		"resourceGroupId": "rg-1",
		"userId":          "1234567890",
		"remark":          "example",
		"topicCount":      json.Number("3"),
		"groupCount":      json.Number("2"),
		"createTime":      "2026-01-01 00:00:00",
		"updateTime":      "2026-01-02 00:00:00",
		"startTime":       "2026-01-01 00:00:00",
		"expireTime":      "2026-02-01 00:00:00",
		"releaseTime":     "2026-02-02 00:00:00",
		"tags": []interface{}{
			map[string]interface{}{"key": "Created", "value": "TF"},
		},
		"productInfo": map[string]interface{}{
			"traceOn":      true,
			"capacityType": "provisioned",
		},
	}

	got := instanceFromAPI(item)
	checks := []struct {
		name string
		got  string
		want string
	}{
		{"instance_id", got.InstanceId.ValueString(), "rmq-1"},
		{"user_id", got.UserId.ValueString(), "1234567890"},
		{"commodity_code", got.CommodityCode.ValueString(), "ons_rmqsub_public_cn"},
		{"remark", got.Remark.ValueString(), "example"},
		{"update_time", got.UpdateTime.ValueString(), "2026-01-02 00:00:00"},
	}
	for _, c := range checks {
		if c.got != c.want {
			t.Fatalf("%s: expected %q, got %q", c.name, c.want, c.got)
		}
	}
	if got.TopicCount.ValueInt64() != 3 || got.GroupCount.ValueInt64() != 2 {
		t.Fatalf("expected numeric counts 3/2, got %d/%d", got.TopicCount.ValueInt64(), got.GroupCount.ValueInt64())
	}
	if len(got.ProductInfo) != 1 {
		t.Fatalf("expected one product_info element, got %d", len(got.ProductInfo))
	}
	if !got.ProductInfo[0].TraceOn.ValueBool() || got.ProductInfo[0].CapacityType.ValueString() != "provisioned" {
		t.Fatalf("unexpected product_info mapping: %+v", got.ProductInfo[0])
	}
	if len(got.Tags) != 1 || got.Tags[0].Key.ValueString() != "Created" {
		t.Fatalf("unexpected tags mapping: %+v", got.Tags)
	}
}

func TestUnitRocketmqConsumerGroupMapping(t *testing.T) {
	item := map[string]interface{}{
		"consumerGroupId": "GID-1",
		"status":          "RUNNING",
		"remark":          "example",
		"instanceId":      "rmq-1",
		"regionId":        "cn-hangzhou",
		"topicName":       "lite-topic",
		"createTime":      "2026-01-01 00:00:00",
		"updateTime":      "2026-01-02 00:00:00",
		"maxReceiveTps":   json.Number("1000"),
		"messageModel":    "LITE_SELECTIVE",
	}

	got := consumerGroupFromAPI(item)
	if got.MessageModel.ValueString() != "LITE_SELECTIVE" {
		t.Fatalf("expected message_model to be mapped from messageModel, got %q", got.MessageModel.ValueString())
	}
	if got.TopicName.ValueString() != "lite-topic" {
		t.Fatalf("expected topic_name to be mapped from topicName, got %q", got.TopicName.ValueString())
	}
	if got.MaxReceiveTps.ValueInt64() != 1000 {
		t.Fatalf("expected max_receive_tps 1000, got %d", got.MaxReceiveTps.ValueInt64())
	}

	// A consumer group that is not a lite group has no messageModel field; the
	// attribute stays in the schema with an empty value.
	normal := consumerGroupFromAPI(map[string]interface{}{"consumerGroupId": "GID-2"})
	if normal.MessageModel.ValueString() != "" {
		t.Fatalf("expected an empty message_model for a normal group, got %q", normal.MessageModel.ValueString())
	}
}

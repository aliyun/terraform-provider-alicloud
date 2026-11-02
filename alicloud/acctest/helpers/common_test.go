package helpers

import (
	"testing"
)

func TestUnitResourceDescribeMethod(t *testing.T) {
	testCases := []struct {
		name       string
		resourceId string
		want       string
		wantErr    bool
	}{
		{"single word", "alicloud_instance.foo", "DescribeInstance", false},
		{"multi underscore", "alicloud_cs_kubernetes_addon.bar", "DescribeCsKubernetesAddon", false},
		{"data source", "data.alicloud_dns_records.record", "DescribeDnsRecords", false},
		{"no alicloud_ prefix", "instance.foo", "", true},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			got, err := getResourceDescribeMethod(testCase.resourceId)
			if testCase.wantErr {
				if err == nil {
					t.Fatalf("getResourceDescribeMethod(%q) expected error, got %q", testCase.resourceId, got)
				}
				return
			}
			if err != nil {
				t.Fatalf("getResourceDescribeMethod(%q) unexpected error: %s", testCase.resourceId, err)
			}
			if got != testCase.want {
				t.Fatalf("getResourceDescribeMethod(%q) = %q, want %q", testCase.resourceId, got, testCase.want)
			}
		})
	}
}

func TestUnitResourceConfigUpdateSentinels(t *testing.T) {
	t.Run("REMOVEKEY on absent key", func(t *testing.T) {
		b := &resourceConfig{attributeMap: map[string]interface{}{"keep": "v"}}
		b.configUpdate(map[string]interface{}{"absent": REMOVEKEY})
		if _, ok := b.attributeMap["absent"]; ok {
			t.Fatalf("REMOVEKEY on absent key must delete without leaking the sentinel, got attributeMap %v", b.attributeMap)
		}
		if got := b.attributeMap["keep"]; got != "v" {
			t.Fatalf("unrelated key must survive, attributeMap %v", b.attributeMap)
		}
	})

	t.Run("CLEARMAP replaces with empty map", func(t *testing.T) {
		b := &resourceConfig{attributeMap: map[string]interface{}{"m": map[string]interface{}{"k": "v"}}}
		b.configUpdate(map[string]interface{}{"m": CLEARMAP})
		got, ok := b.attributeMap["m"].(map[string]interface{})
		if !ok || len(got) != 0 {
			t.Fatalf("CLEARMAP must set an empty map, got %#v", b.attributeMap["m"])
		}
	})

	t.Run("CLEARLIST replaces with empty list", func(t *testing.T) {
		b := &resourceConfig{attributeMap: map[string]interface{}{"l": []interface{}{"a"}}}
		b.configUpdate(map[string]interface{}{"l": CLEARLIST})
		got, ok := b.attributeMap["l"].([]interface{})
		if !ok || len(got) != 0 {
			t.Fatalf("CLEARLIST must set an empty list, got %#v", b.attributeMap["l"])
		}
	})

	t.Run("overwrite keeps other keys", func(t *testing.T) {
		b := &resourceConfig{attributeMap: map[string]interface{}{"a": "1"}}
		b.configUpdate(map[string]interface{}{"a": "2", "b": "3"})
		if b.attributeMap["a"] != "2" || b.attributeMap["b"] != "3" {
			t.Fatalf("configUpdate must upsert, got attributeMap %v", b.attributeMap)
		}
	})
}

func TestUnitResourceAttrUpdateCheckMapPair(t *testing.T) {
	original := map[string]string{"a": "1", "b": "2"}
	ra := ResourceAttrInit("alicloud_instance.foo", original)

	changeMap := map[string]string{"b": REMOVEKEY, "c": "3", ForceSleep: "0"}
	ra.updateCheckMapPair(changeMap)

	if _, ok := ra.checkMap["b"]; ok {
		t.Fatalf("REMOVEKEY must drop the key, got checkMap %v", ra.checkMap)
	}
	if ra.checkMap["a"] != "1" || ra.checkMap["c"] != "3" {
		t.Fatalf("updateCheckMapPair must keep other keys and add new ones, got checkMap %v", ra.checkMap)
	}
	if _, ok := changeMap[ForceSleep]; ok {
		t.Fatalf("valid ForceSleep interval must be consumed, got changeMap %v", changeMap)
	}
	if original["b"] != "2" {
		t.Fatalf("updateCheckMapPair must copy before mutating, original map changed: %v", original)
	}

	ra.updateCheckMapPair(map[string]string{ForceSleep: "not-a-number"})
	if len(ra.checkMap) != 3 {
		t.Fatalf("invalid ForceSleep interval must be ignored, got checkMap %v", ra.checkMap)
	}
}

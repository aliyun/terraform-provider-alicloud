package cr

import (
	"context"
	"testing"
)

func TestCRRepoStateUpgradeV0(t *testing.T) {
	cases := []struct {
		name     string
		input    map[string]interface{}
		expected []interface{}
	}{
		{
			name: "map with data",
			input: map[string]interface{}{
				"id":        "test-ns/test-repo",
				"namespace": "test-ns",
				"name":      "test-repo",
				"domain_list": map[string]interface{}{
					"vpc":      "registry-vpc.cn-hangzhou.cr.aliyuncs.com",
					"public":   "registry.cn-hangzhou.cr.aliyuncs.com",
					"internal": "registry-internal.cn-hangzhou.cr.aliyuncs.com",
				},
			},
			expected: []interface{}{
				map[string]interface{}{
					"vpc":      "registry-vpc.cn-hangzhou.cr.aliyuncs.com",
					"public":   "registry.cn-hangzhou.cr.aliyuncs.com",
					"internal": "registry-internal.cn-hangzhou.cr.aliyuncs.com",
				},
			},
		},
		{
			name: "empty map",
			input: map[string]interface{}{
				"id":          "test-ns/test-repo",
				"domain_list": map[string]interface{}{},
			},
			expected: []interface{}{},
		},
		{
			name: "nil value",
			input: map[string]interface{}{
				"id":          "test-ns/test-repo",
				"domain_list": nil,
			},
			expected: nil,
		},
		{
			name: "field not present",
			input: map[string]interface{}{
				"id": "test-ns/test-repo",
			},
			expected: nil,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := resourceAlicloudCRRepoStateUpgradeV0(context.Background(), tc.input, nil)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			got := result["domain_list"]
			if tc.expected == nil {
				if got != nil {
					t.Errorf("expected nil, got %v", got)
				}
				return
			}
			gotList, ok := got.([]interface{})
			if !ok {
				t.Fatalf("expected []interface{}, got %T", got)
			}
			if len(gotList) != len(tc.expected) {
				t.Fatalf("expected length %d, got %d", len(tc.expected), len(gotList))
			}
			for i, item := range gotList {
				gotMap := item.(map[string]interface{})
				expMap := tc.expected[i].(map[string]interface{})
				for k, v := range expMap {
					if gotMap[k] != v {
						t.Errorf("item[%d].%s: expected %v, got %v", i, k, v, gotMap[k])
					}
				}
			}
		})
	}
}

func TestCRRepoSchemaVersionV0ToV1(t *testing.T) {
	r := resourceAlicloudCRRepo()
	if r.SchemaVersion != 1 {
		t.Errorf("expected SchemaVersion 1, got %d", r.SchemaVersion)
	}
	if len(r.StateUpgraders) != 1 {
		t.Fatalf("expected 1 StateUpgrader, got %d", len(r.StateUpgraders))
	}
	if r.StateUpgraders[0].Version != 0 {
		t.Errorf("expected StateUpgrader version 0, got %d", r.StateUpgraders[0].Version)
	}
}

package cs

import (
	"context"
	"testing"
)

func TestUnitCSEdgeKubernetesStateUpgradeV0(t *testing.T) {
	fields := []string{"runtime", "certificate_authority", "connections"}
	cases := []struct {
		name     string
		input    map[string]interface{}
		expected map[string][]interface{}
	}{
		{
			name: "map with data",
			input: map[string]interface{}{
				"id": "c-edge-123",
				"runtime": map[string]interface{}{
					"name":    "containerd",
					"version": "1.6.20",
				},
				"certificate_authority": map[string]interface{}{
					"cluster_cert": "cluster-cert",
					"client_cert":  "client-cert",
					"client_key":   "client-key",
				},
				"connections": map[string]interface{}{
					"api_server_internet": "https://1.2.3.4:6443",
					"api_server_intranet": "https://10.0.0.1:6443",
					"master_public_ip":    "1.2.3.4",
					"service_domain":      "*.c-edge-123.cs.local",
				},
			},
			expected: map[string][]interface{}{
				"runtime": {
					map[string]interface{}{
						"name":    "containerd",
						"version": "1.6.20",
					},
				},
				"certificate_authority": {
					map[string]interface{}{
						"cluster_cert": "cluster-cert",
						"client_cert":  "client-cert",
						"client_key":   "client-key",
					},
				},
				"connections": {
					map[string]interface{}{
						"api_server_internet": "https://1.2.3.4:6443",
						"api_server_intranet": "https://10.0.0.1:6443",
						"master_public_ip":    "1.2.3.4",
						"service_domain":      "*.c-edge-123.cs.local",
					},
				},
			},
		},
		{
			name: "empty map",
			input: map[string]interface{}{
				"id":                    "c-edge-456",
				"runtime":               map[string]interface{}{},
				"certificate_authority": map[string]interface{}{},
				"connections":           map[string]interface{}{},
			},
			expected: map[string][]interface{}{
				"runtime":               {},
				"certificate_authority": {},
				"connections":           {},
			},
		},
		{
			name: "nil value",
			input: map[string]interface{}{
				"id":                    "c-edge-789",
				"runtime":               nil,
				"certificate_authority": nil,
				"connections":           nil,
			},
			expected: nil,
		},
		{
			name: "field not present",
			input: map[string]interface{}{
				"id": "c-edge-000",
			},
			expected: nil,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := resourceCSEdgeKubernetesStateUpgradeV0(context.Background(), tc.input, nil)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			for _, field := range fields {
				got := result[field]
				if tc.expected == nil {
					if got != nil {
						t.Errorf("%s: expected nil, got %v", field, got)
					}
					continue
				}
				exp := tc.expected[field]
				gotList, ok := got.([]interface{})
				if !ok {
					t.Fatalf("%s: expected []interface{}, got %T", field, got)
				}
				if len(gotList) != len(exp) {
					t.Fatalf("%s: expected length %d, got %d", field, len(exp), len(gotList))
				}
				for i, item := range gotList {
					gotMap := item.(map[string]interface{})
					expMap := exp[i].(map[string]interface{})
					for k, v := range expMap {
						if gotMap[k] != v {
							t.Errorf("%s[%d].%s: expected %v, got %v", field, i, k, v, gotMap[k])
						}
					}
				}
			}
		})
	}
}

func TestUnitCSEdgeKubernetesSchemaVersionV0ToV1(t *testing.T) {
	r := ResourceCSEdgeKubernetes()
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

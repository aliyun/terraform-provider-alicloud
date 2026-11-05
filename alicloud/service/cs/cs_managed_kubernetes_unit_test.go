package cs

import (
	"context"
	"testing"
)

func TestUnit_versionCompare(t *testing.T) {
	tests := []struct {
		name           string
		oldVersion     string
		newVersion     string
		expectedResult int
		expectError    bool
	}{
		// basic semver
		{name: "equal versions", oldVersion: "1.2.3", newVersion: "1.2.3", expectedResult: 0},
		{name: "new version newer patch", oldVersion: "1.2.3", newVersion: "1.2.4", expectedResult: 1},
		{name: "new version older patch", oldVersion: "1.2.4", newVersion: "1.2.3", expectedResult: -1},
		{name: "new version newer minor", oldVersion: "1.2.3", newVersion: "1.3.0", expectedResult: 1},
		{name: "new version older minor", oldVersion: "1.3.0", newVersion: "1.2.3", expectedResult: -1},
		{name: "new version newer major", oldVersion: "1.2.3", newVersion: "2.0.0", expectedResult: 1},
		{name: "new version older major", oldVersion: "2.0.0", newVersion: "1.2.3", expectedResult: -1},

		// v prefix
		{name: "v prefix both", oldVersion: "v1.2.3", newVersion: "v1.2.4", expectedResult: 1},
		{name: "v prefix old only", oldVersion: "v1.2.3", newVersion: "1.2.4", expectedResult: 1},
		{name: "v prefix new only", oldVersion: "1.2.3", newVersion: "v1.2.4", expectedResult: 1},
		{name: "v prefix equal", oldVersion: "v1.2.3", newVersion: "v1.2.3", expectedResult: 0},

		// pre-release suffix (per semver: release > pre-release, e.g. 1.2.3 > 1.2.3-alpha.1)
		{name: "release vs pre-release", oldVersion: "1.2.3-alpha.1", newVersion: "1.2.3", expectedResult: 1},
		{name: "pre-release vs release", oldVersion: "1.2.3", newVersion: "1.2.3-alpha.1", expectedResult: -1},
		{name: "pre-release equal", oldVersion: "1.2.3-alpha.1", newVersion: "1.2.3-alpha.1", expectedResult: 0},
		{name: "pre-release numeric diff", oldVersion: "1.2.3-alpha.1", newVersion: "1.2.3-alpha.2", expectedResult: 1},
		{name: "pre-release string diff", oldVersion: "1.2.3-alpha.1", newVersion: "1.2.3-beta.1", expectedResult: 1},
		{name: "pre-release with aliyun", oldVersion: "1.9.3-aliyun.1", newVersion: "1.9.7-aliyun.2", expectedResult: 1},
		{name: "pre-release newer major version", oldVersion: "1.9.3-aliyun.1", newVersion: "2.0.0-aliyun.1", expectedResult: 1},

		// apsara format
		{name: "apsara same main version, suffix newer", oldVersion: "v1.31.0-apsara.6.11.6.ad796663", newVersion: "v1.31.0-apsara.6.11.7.17d202a9", expectedResult: 1},
		{name: "apsara same main version, suffix older", oldVersion: "v1.31.0-apsara.6.11.7.17d202a9", newVersion: "v1.31.0-apsara.6.11.6.ad796663", expectedResult: -1},
		{name: "apsara same numeric suffix, diff commit hash", oldVersion: "v1.31.0-apsara.6.11.6.ad796663", newVersion: "v1.31.0-apsara.6.11.6.bf123456", expectedResult: 1},
		{name: "apsara different main version", oldVersion: "v1.30.0-apsara.6.11.6.ad796663", newVersion: "v1.31.0-apsara.6.11.6.ad796663", expectedResult: 1},

		// empty versions
		{name: "both empty", oldVersion: "", newVersion: "", expectedResult: 0},
		{name: "old empty", oldVersion: "", newVersion: "1.2.3", expectedResult: 1},
		{name: "new empty", oldVersion: "1.2.3", newVersion: "", expectedResult: -1},

		// invalid format
		{name: "invalid old version", oldVersion: "latest", newVersion: "1.2.3", expectError: true},
		{name: "invalid new version", oldVersion: "1.2.3", newVersion: "latest", expectError: true},
		{name: "both invalid", oldVersion: "abc", newVersion: "xyz", expectError: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := versionCompare(tt.oldVersion, tt.newVersion)
			if tt.expectError {
				if err == nil {
					t.Errorf("expected error but got nil, result=%d", result)
				}
				return
			}
			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}
			if result != tt.expectedResult {
				t.Errorf("versionCompare(%q, %q) = %d, want %d", tt.oldVersion, tt.newVersion, result, tt.expectedResult)
			}
		})
	}
}

func TestUnitCSManagedKubernetesStateUpgradeV0(t *testing.T) {
	fields := []string{"certificate_authority", "connections"}
	cases := []struct {
		name     string
		input    map[string]interface{}
		expected map[string][]interface{}
	}{
		{
			name: "map with data",
			input: map[string]interface{}{
				"id": "c-managed-123",
				"certificate_authority": map[string]interface{}{
					"cluster_cert": "cert",
					"client_cert":  "client",
					"client_key":   "key",
				},
				"connections": map[string]interface{}{
					"api_server_internet": "https://3.4.5.6:6443",
					"api_server_intranet": "https://10.0.0.2:6443",
					"master_public_ip":    "3.4.5.6",
					"service_domain":      "*.c-managed-123.cs.local",
				},
			},
			expected: map[string][]interface{}{
				"certificate_authority": {
					map[string]interface{}{
						"cluster_cert": "cert",
						"client_cert":  "client",
						"client_key":   "key",
					},
				},
				"connections": {
					map[string]interface{}{
						"api_server_internet": "https://3.4.5.6:6443",
						"api_server_intranet": "https://10.0.0.2:6443",
						"master_public_ip":    "3.4.5.6",
						"service_domain":      "*.c-managed-123.cs.local",
					},
				},
			},
		},
		{
			name: "empty map",
			input: map[string]interface{}{
				"id":                    "c-managed-456",
				"certificate_authority": map[string]interface{}{},
				"connections":           map[string]interface{}{},
			},
			expected: map[string][]interface{}{
				"certificate_authority": {},
				"connections":           {},
			},
		},
		{
			name: "nil value",
			input: map[string]interface{}{
				"id":                    "c-managed-789",
				"certificate_authority": nil,
				"connections":           nil,
			},
			expected: nil,
		},
		{
			name: "field not present",
			input: map[string]interface{}{
				"id": "c-managed-000",
			},
			expected: nil,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := resourceCSManagedKubernetesStateUpgradeV0(context.Background(), tc.input, nil)
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

func TestUnitCSManagedKubernetesSchemaVersionV0ToV1(t *testing.T) {
	r := ResourceCSManagedKubernetes()
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

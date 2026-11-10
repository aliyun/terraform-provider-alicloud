package cs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnitCSAutoRepairPolicyDiffSuppress(t *testing.T) {
	tests := []struct {
		name     string
		oldValue string
		newValue string
		want     bool
	}{
		{name: "identical values", oldValue: `{"a":1}`, newValue: `{"a":1}`, want: true},
		{name: "empty and empty json object", oldValue: "{}", newValue: "", want: true},
		{name: "null and empty json object", oldValue: "{}", newValue: "null", want: true},
		{name: "same object with different key order", oldValue: `{"a":1,"b":2}`, newValue: `{"b":2,"a":1}`, want: true},
		{name: "different values", oldValue: `{"a":1}`, newValue: `{"a":2}`, want: false},
		{name: "invalid json falls back to a diff", oldValue: `{"a":1}`, newValue: `not-json`, want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, autoRepairPolicyConfigDiffSuppress("", tt.oldValue, tt.newValue, nil))
		})
	}
}

func TestUnitCSAutoRepairPolicyParseId(t *testing.T) {
	tests := []struct {
		name        string
		id          string
		wantErr     bool
		wantCluster string
		wantPolicy  string
	}{
		{name: "valid id", id: "c-1:r-1", wantCluster: "c-1", wantPolicy: "r-1"},
		{name: "missing policy part", id: "c-1", wantErr: true},
		{name: "empty policy part", id: "c-1:", wantErr: true},
		{name: "extra part", id: "c-1:r-1:extra", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clusterId, policyId, err := parseCSAutoRepairPolicyId(tt.id)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.wantCluster, clusterId)
			assert.Equal(t, tt.wantPolicy, policyId)
		})
	}
}

func TestUnitCSAutoRepairPolicyValidateRules(t *testing.T) {
	newRules := func(procedures ...map[string]interface{}) []interface{} {
		list := make([]interface{}, 0, len(procedures))
		for _, procedure := range procedures {
			list = append(list, procedure)
		}
		return []interface{}{map[string]interface{}{"repair_procedure": list}}
	}

	tests := []struct {
		name    string
		rules   interface{}
		wantErr string
	}{
		{
			name:  "nil rules",
			rules: nil,
		},
		{
			name:  "quarantine before evict",
			rules: newRules(map[string]interface{}{"name": "QuarantineNode"}, map[string]interface{}{"name": "Evict", "config": `{"evict_timeout":300}`}),
		},
		{
			name:  "quarantine, evict, repair and unquarantine in order",
			rules: newRules(map[string]interface{}{"name": "QuarantineNode"}, map[string]interface{}{"name": "Evict"}, map[string]interface{}{"name": "Repair"}, map[string]interface{}{"name": "UnquarantineNode"}),
		},
		{
			name:  "repair only does not depend on quarantine",
			rules: newRules(map[string]interface{}{"name": "Repair"}),
		},
		{
			name:    "evict without quarantine",
			rules:   newRules(map[string]interface{}{"name": "Evict"}),
			wantErr: "repair procedure Evict requires QuarantineNode to be defined before it",
		},
		{
			name:    "evict before quarantine",
			rules:   newRules(map[string]interface{}{"name": "Evict"}, map[string]interface{}{"name": "QuarantineNode"}),
			wantErr: "repair procedure Evict requires QuarantineNode to be defined before it",
		},
		{
			name:    "unquarantine without quarantine",
			rules:   newRules(map[string]interface{}{"name": "Repair"}, map[string]interface{}{"name": "UnquarantineNode"}),
			wantErr: "repair procedure UnquarantineNode requires QuarantineNode to be defined before it",
		},
		{
			name:    "duplicated procedure",
			rules:   newRules(map[string]interface{}{"name": "QuarantineNode"}, map[string]interface{}{"name": "QuarantineNode"}),
			wantErr: "repair procedure QuarantineNode is duplicated",
		},
		{
			name:    "invalid config json",
			rules:   newRules(map[string]interface{}{"name": "QuarantineNode"}, map[string]interface{}{"name": "Evict", "config": `evict_timeout=300`}, map[string]interface{}{"name": "Repair"}),
			wantErr: "repair_procedure.1.config must be a valid JSON object",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateCSAutoRepairPolicyRules(tt.rules)
			if tt.wantErr == "" {
				assert.NoError(t, err)
				return
			}
			if assert.Error(t, err) {
				assert.Contains(t, err.Error(), tt.wantErr)
			}
		})
	}
}

func TestUnitCSAutoRepairPolicyMatchResourceSubType(t *testing.T) {
	tests := []struct {
		name string
		raw  string
		want string
		ok   bool
	}{
		{name: "single value", raw: "ess", want: "ess", ok: true},
		{name: "set contains value", raw: "ess|lingjun", want: "lingjun", ok: true},
		{name: "set without value", raw: "lingjun", want: "ess", ok: false},
		{name: "empty value", raw: "", want: "ess", ok: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.ok, matchCSAutoRepairResourceSubType(tt.raw, tt.want))
		})
	}
}

func TestUnitCSAutoRepairRateLimit(t *testing.T) {
	tests := []struct {
		name    string
		value   interface{}
		wantErr bool
	}{
		{name: "empty value is accepted", value: "", wantErr: false},
		{name: "node count", value: "2", wantErr: false},
		{name: "upper bound node count", value: "100000", wantErr: false},
		{name: "zero node count", value: "0", wantErr: true},
		{name: "negative node count", value: "-1", wantErr: true},
		{name: "not a number", value: "two", wantErr: true},
		{name: "percentage", value: "50%", wantErr: false},
		{name: "percentage out of range", value: "101%", wantErr: true},
		{name: "zero percentage", value: "0%", wantErr: true},
		{name: "non string value", value: 2, wantErr: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, errors := validateCSAutoRepairRateLimit(tt.value, "max_parallel_repairing_nodes")
			assert.Equal(t, tt.wantErr, len(errors) > 0)
		})
	}
}

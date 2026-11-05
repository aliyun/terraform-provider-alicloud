package alikafka

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
)

func TestUnitCommonAlikafkaInstanceConfigDiffSuppressFunc(t *testing.T) {
	testCases := []struct {
		name        string
		old         string
		new         string
		expected    bool
		description string
	}{
		{
			name:        "New_Empty",
			old:         `{"key":"value"}`,
			new:         "",
			expected:    true,
			description: "Empty new value should suppress diff",
		},
		{
			name:        "Old_Empty",
			old:         "",
			new:         `{"key":"value"}`,
			expected:    false,
			description: "Empty old value should not suppress diff",
		},
		{
			name:        "Old_Invalid_JSON",
			old:         "invalid json",
			new:         `{"key":"value"}`,
			expected:    false,
			description: "Invalid old JSON should not suppress diff",
		},
		{
			name:        "New_Invalid_JSON",
			old:         `{"key":"value"}`,
			new:         "invalid json",
			expected:    false,
			description: "Invalid new JSON should not suppress diff",
		},
		{
			name:        "Same_Values",
			old:         `{"key":"value","other":"data"}`,
			new:         `{"key":"value"}`,
			expected:    true,
			description: "New is subset of old with same values should suppress diff",
		},
		{
			name:        "Different_Values",
			old:         `{"key":"old_value"}`,
			new:         `{"key":"new_value"}`,
			expected:    false,
			description: "Different values for same key should not suppress diff",
		},
		{
			name:        "New_Key_Not_In_Old",
			old:         `{"other":"value"}`,
			new:         `{"key":"value"}`,
			expected:    true,
			description: "New key not in old map should suppress diff",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			d := schema.TestResourceDataRaw(t, map[string]*schema.Schema{}, map[string]interface{}{})
			result := alikafkaInstanceConfigDiffSuppressFunc("key", tc.old, tc.new, d)
			assert.Equal(t, tc.expected, result, tc.description)
		})
	}
}

func TestUnitAlikafkaTopicConfigsDiffSuppressFunc(t *testing.T) {
	// The service always reports the full effective topic config, so `old` (state) is a
	// superset of `new` (the keys declared in HCL) whenever the instance injects defaults.
	serverInjectedConfig := `{"cloud.native.topic.type":"true","replication-factor":"3","max.message.bytes":"1048576","retention.hours":"72"}`

	testCases := []struct {
		name        string
		old         string
		new         string
		expected    bool
		description string
	}{
		{
			name:        "ServerInjectedKeysIgnored",
			old:         serverInjectedConfig,
			new:         `{"max.message.bytes":"1048576","retention.hours":"72"}`,
			expected:    true,
			description: "keys injected by the service must not produce a diff",
		},
		{
			name:        "DeclaredKeyChanged",
			old:         serverInjectedConfig,
			new:         `{"max.message.bytes":"10485760","retention.hours":"72"}`,
			expected:    false,
			description: "a real change on a declared key must still produce a diff",
		},
		{
			name:        "NumberOnStateSide",
			old:         `{"cloud.native.topic.type":"true","replication-factor":3,"retention.hours":72}`,
			new:         `{"retention.hours":"72"}`,
			expected:    true,
			description: "numeric values reported by the service compare equal to their string form",
		},
		{
			name:        "NumberOnConfigSide",
			old:         serverInjectedConfig,
			new:         `{"max.message.bytes":1048576}`,
			expected:    true,
			description: "jsonencode with a numeric value compares equal to the string form",
		},
		{
			name:        "BoolAndStringBoolAreEqual",
			old:         `{"cloud.native.topic.type":"true","preallocate":true}`,
			new:         `{"preallocate":"true"}`,
			expected:    true,
			description: "boolean values compare equal to their string form",
		},
		{
			name:        "RemovedKeyIsNotADiff",
			old:         serverInjectedConfig,
			new:         `{"retention.hours":"72"}`,
			expected:    true,
			description: "dropping a previously declared key is suppressed because the update API only merges keys",
		},
		{
			name:        "NewlyDeclaredKeyIsADiff",
			old:         `{"cloud.native.topic.type":"true","replication-factor":"3"}`,
			new:         `{"retention.hours":"72"}`,
			expected:    false,
			description: "a key the service does not report yet must produce a diff so it gets applied",
		},
		{
			name:        "EmptyConfigDeclared",
			old:         serverInjectedConfig,
			new:         `{}`,
			expected:    true,
			description: "declaring an empty object enforces nothing",
		},
		{
			name:        "EmptyNew",
			old:         serverInjectedConfig,
			new:         "",
			expected:    true,
			description: "an unset attribute keeps the value read back from the service",
		},
		{
			name:        "EmptyOld",
			old:         "",
			new:         `{"retention.hours":"72"}`,
			expected:    false,
			description: "nothing in state yet, so the declared config must be applied",
		},
		{
			name:        "InvalidNewJson",
			old:         serverInjectedConfig,
			new:         `{"retention.hours":`,
			expected:    false,
			description: "malformed config must not be suppressed",
		},
		{
			name:        "InvalidOldJson",
			old:         `not-json`,
			new:         `{"retention.hours":"72"}`,
			expected:    false,
			description: "malformed state must not be suppressed",
		},
		{
			name:        "TrailingDataInState",
			old:         `{"retention.hours":"72"}garbage`,
			new:         `{"retention.hours":"72"}`,
			expected:    false,
			description: "a document with data after the object must not be suppressed",
		},
		{
			name:        "TrailingDataInConfig",
			old:         `{"retention.hours":"72"}`,
			new:         `{"retention.hours":"72"}{"max.message.bytes":"1048576"}`,
			expected:    false,
			description: "a declared document with data after the object must not be suppressed",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			d := schema.TestResourceDataRaw(t, map[string]*schema.Schema{}, map[string]interface{}{})
			result := alikafkaTopicConfigsDiffSuppressFunc("configs", tc.old, tc.new, d)
			assert.Equal(t, tc.expected, result, tc.description)
		})
	}
}

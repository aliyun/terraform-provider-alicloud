// lintignore: S013
package smartag

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
)

func TestUnitCommonSagDnatEntryTypeDiffSuppressFunc(t *testing.T) {
	testCases := []struct {
		name        string
		entryType   string
		expected    bool
		description string
	}{
		{
			name:        "Intranet_Type",
			entryType:   "Intranet",
			expected:    false,
			description: "Intranet type should not suppress diff",
		},
		{
			name:        "Internet_Type",
			entryType:   "Internet",
			expected:    true,
			description: "Internet type should suppress diff",
		},
		{
			name:        "Empty_Type",
			entryType:   "",
			expected:    true,
			description: "Empty type should suppress diff",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			d := schema.TestResourceDataRaw(t, map[string]*schema.Schema{
				"type": {Type: schema.TypeString},
			}, map[string]interface{}{
				"type": tc.entryType,
			})
			result := sagDnatEntryTypeDiffSuppressFunc("key", "old", "new", d)
			assert.Equal(t, tc.expected, result, tc.description)
		})
	}
}

func TestUnitCommonSagClientUserPasswordSuppressFunc(t *testing.T) {
	testCases := []struct {
		name        string
		userName    string
		expected    bool
		description string
	}{
		{
			name:        "With_UserName",
			userName:    "alice",
			expected:    false,
			description: "Non-empty user_name should not suppress diff",
		},
		{
			name:        "Empty_UserName",
			userName:    "",
			expected:    true,
			description: "Empty user_name should suppress diff",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			d := schema.TestResourceDataRaw(t, map[string]*schema.Schema{
				"user_name": {Type: schema.TypeString},
			}, map[string]interface{}{
				"user_name": tc.userName,
			})
			result := sagClientUserPasswordSuppressFunc("key", "old", "new", d)
			assert.Equal(t, tc.expected, result, tc.description)
		})
	}
}

// lintignore: S013
package expressconnect

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
)

func TestUnitCommonRouterInterfaceAcceptsideDiffSuppressFunc(t *testing.T) {
	testCases := []struct {
		name        string
		role        string
		expected    bool
		description string
	}{
		{
			name:        "AcceptingSide",
			role:        "AcceptingSide",
			expected:    true,
			description: "AcceptingSide role should suppress diff",
		},
		{
			name:        "InitiatingSide",
			role:        "InitiatingSide",
			expected:    false,
			description: "InitiatingSide role should not suppress diff",
		},
		{
			name:        "Empty_Role",
			role:        "",
			expected:    false,
			description: "Empty role should not suppress diff",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			d := schema.TestResourceDataRaw(t, map[string]*schema.Schema{
				"role": {Type: schema.TypeString},
			}, map[string]interface{}{
				"role": tc.role,
			})
			result := routerInterfaceAcceptsideDiffSuppressFunc("key", "old", "new", d)
			assert.Equal(t, tc.expected, result, tc.description)
		})
	}
}

func TestUnitCommonRouterInterfaceVBRTypeDiffSuppressFunc(t *testing.T) {
	testCases := []struct {
		name        string
		role        string
		routerType  string
		expected    bool
		description string
	}{
		{
			name:        "AcceptingSide",
			role:        "AcceptingSide",
			expected:    true,
			description: "AcceptingSide role should suppress diff",
		},
		{
			name:        "InitiatingSide_VRouter",
			role:        "InitiatingSide",
			routerType:  "VRouter",
			expected:    true,
			description: "VRouter type should suppress diff",
		},
		{
			name:        "InitiatingSide_VBR",
			role:        "InitiatingSide",
			routerType:  "VBR",
			expected:    false,
			description: "InitiatingSide with VBR should not suppress diff",
		},
		{
			name:        "Empty_Role_VBR",
			role:        "",
			routerType:  "VBR",
			expected:    false,
			description: "Empty role with VBR should not suppress diff",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			d := schema.TestResourceDataRaw(t, map[string]*schema.Schema{
				"role":        {Type: schema.TypeString},
				"router_type": {Type: schema.TypeString},
			}, map[string]interface{}{
				"role":        tc.role,
				"router_type": tc.routerType,
			})
			result := routerInterfaceVBRTypeDiffSuppressFunc("key", "old", "new", d)
			assert.Equal(t, tc.expected, result, tc.description)
		})
	}
}

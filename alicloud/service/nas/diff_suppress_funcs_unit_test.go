// lintignore: S013
package nas

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
)

func TestUnitNasFileSystemVpcDiffSuppressFunc(t *testing.T) {
	testCases := []struct {
		name           string
		fileSystemType string
		expected       bool
		description    string
	}{
		{
			name:           "Standard_Suppress",
			fileSystemType: "standard",
			expected:       true,
			description:    "standard file system should suppress vpc_id/vswitch_id diff (API does not return them)",
		},
		{
			name:           "Extreme_Suppress",
			fileSystemType: "extreme",
			expected:       true,
			description:    "extreme file system should suppress vpc_id/vswitch_id diff (API does not return them)",
		},
		{
			name:           "Cpfs_NoSuppress",
			fileSystemType: "cpfs",
			expected:       false,
			description:    "cpfs file system should not suppress diff (API returns VpcId/VSwitchId)",
		},
		{
			name:           "Cpfsse_NoSuppress",
			fileSystemType: "cpfsse",
			expected:       false,
			description:    "cpfsse file system should not suppress diff",
		},
		{
			name:           "Empty_NoSuppress",
			fileSystemType: "",
			expected:       false,
			description:    "unknown file system type should not suppress diff",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			d := schema.TestResourceDataRaw(t, map[string]*schema.Schema{
				"file_system_type": {Type: schema.TypeString},
				"vpc_id":           {Type: schema.TypeString},
				"vswitch_id":       {Type: schema.TypeString},
			}, map[string]interface{}{
				"file_system_type": tc.fileSystemType,
				"vpc_id":           "vpc-123456",
				"vswitch_id":       "vsw-123456",
			})

			result := nasFileSystemVpcDiffSuppressFunc("vpc_id", "", "vpc-123456", d)
			assert.Equal(t, tc.expected, result, tc.description)
		})
	}
}

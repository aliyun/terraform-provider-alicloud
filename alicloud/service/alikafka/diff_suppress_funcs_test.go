// lintignore: S013
package alikafka

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
)

func TestUnitCommonKmsDiffSuppressFunc(t *testing.T) {
	testCases := []struct {
		name            string
		password        string
		accountPassword string
		expected        bool
		description     string
	}{
		{
			name:        "With_Password",
			password:    "mypassword",
			expected:    true,
			description: "With password should suppress diff",
		},
		{
			name:            "With_AccountPassword",
			accountPassword: "myaccountpassword",
			expected:        true,
			description:     "With account_password should suppress diff",
		},
		{
			name:        "No_Password",
			expected:    false,
			description: "Without any password should not suppress diff",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			data := map[string]interface{}{}
			if tc.password != "" {
				data["password"] = tc.password
			}
			if tc.accountPassword != "" {
				data["account_password"] = tc.accountPassword
			}
			d := schema.TestResourceDataRaw(t, map[string]*schema.Schema{
				"password":         {Type: schema.TypeString},
				"account_password": {Type: schema.TypeString},
			}, data)
			result := kmsDiffSuppressFunc("key", "old", "new", d)
			assert.Equal(t, tc.expected, result, tc.description)
		})
	}
}

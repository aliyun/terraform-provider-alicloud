// lintignore: S013
package ram

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnitCommonRamSAMLProviderDiffSuppressFunc(t *testing.T) {
	testCases := []struct {
		name        string
		old         string
		new         string
		expected    bool
		description string
	}{
		{
			name:        "Same_Content",
			old:         "<saml>content</saml>",
			new:         "<saml>content</saml>",
			expected:    true,
			description: "Identical strings should suppress diff",
		},
		{
			name:        "Different_Content",
			old:         "<saml>old</saml>",
			new:         "<saml>new</saml>",
			expected:    false,
			description: "Different strings should not suppress diff",
		},
		{
			name:        "Same_With_Newlines",
			old:         "<saml>\ncontent\n</saml>",
			new:         "<saml>content</saml>",
			expected:    true,
			description: "Strings differing only by newlines should suppress diff",
		},
		{
			name:        "Same_With_Whitespace",
			old:         "  content  ",
			new:         "content",
			expected:    true,
			description: "Strings differing only by surrounding whitespace should suppress diff",
		},
		{
			name:        "Empty_Strings",
			old:         "",
			new:         "",
			expected:    true,
			description: "Empty strings should suppress diff",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := ramSAMLProviderDiffSuppressFunc(tc.old, tc.new)
			assert.Equal(t, tc.expected, result, tc.description)
		})
	}
}

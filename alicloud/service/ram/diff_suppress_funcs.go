package ram

import (
	"strings"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
)

// ramSAMLProviderDiffSuppressFunc migrated from the root alicloud package.
// Two SAML metadata documents compare equal when newlines are stripped and the
// surrounding whitespace is trimmed, so re-formatted but semantically identical
// metadata does not produce a diff.
func ramSAMLProviderDiffSuppressFunc(old, new string) bool {
	if helper.Trim(strings.Replace(old, "\n", "", -1)) != helper.Trim(strings.Replace(new, "\n", "", -1)) {
		return false
	}
	return true
}

package validators

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// PolarDBPostPaidDiffSuppressFunc suppresses period diffs for PolarDB
// PostPaid (pay-as-you-go) applications: the period attribute only applies to
// PrePaid, so changes to it are meaningless for PostPaid and must not drive a
// diff. Previously polardbPostPaidDiffSuppressFunc in the flat layout; it is
// product-flavoured but a schema suppressor, so it lives in the validators
// leaf package and the flat layout keeps a behavior-preserving alias.
func PolarDBPostPaidDiffSuppressFunc(k, old, new string, d *schema.ResourceData) bool {
	if d.Get("pay_type").(string) == "PrePaid" {
		return false
	}
	return true
}

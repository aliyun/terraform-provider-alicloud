package polardb

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// postPaidDiffSuppressFunc suppresses period diffs for PostPaid applications:
// period only applies to PrePaid, so a change to it is meaningless otherwise.
func postPaidDiffSuppressFunc(k, old, new string, d *schema.ResourceData) bool {
	return d.Get("pay_type").(string) != string(PrePaid)
}

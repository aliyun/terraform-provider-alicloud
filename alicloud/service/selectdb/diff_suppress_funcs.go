package selectdb

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

// postPaidDiffSuppressFunc suppresses plan-time diffs for payment-related
// attributes that only apply to Subscription instances. PayAsYouGo instances
// leave these fields empty, so any diff is noise.
func postPaidDiffSuppressFunc(k, old, new string, d *schema.ResourceData) bool {
	if d.Get("payment_type").(string) == "Subscription" {
		return false
	}
	return true
}

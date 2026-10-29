// Package natgateway. Local mirrors of root alicloud package symbols that the
// migrated NAT Gateway source depends on. The originals live in the root
// alicloud package (common.go / diff_suppress_funcs.go) and cannot be imported
// from a service subpackage without an import cycle, so the subset used here is
// kept with identical semantics.
package natgateway

import (
	"strings"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/convert"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// TimeType mirrors alicloud.TimeType for charge-period constants.
type TimeType string

const (
	Hour  = TimeType("Hour")
	Day   = TimeType("Day")
	Week  = TimeType("Week")
	Month = TimeType("Month")
	Year  = TimeType("Year")
)

// COLON_SEPARATED mirrors alicloud.COLON_SEPARATED for composite resource IDs.
const COLON_SEPARATED = ":"

// convertMaptoJsonString mirrors the alicloud package-level alias of
// convert.ConvertMaptoJsonString used by the migrated source.
var convertMaptoJsonString = convert.ConvertMaptoJsonString

// PostPaidDiffSuppressFunc mirrors alicloud.PostPaidDiffSuppressFunc for the
// nat_gateway period field.
func PostPaidDiffSuppressFunc(k, old, new string, d *schema.ResourceData) bool {
	// payment_type is the instance_charge_type's replacement.
	// If both instance_charge_type and payment_type are "", it means hiding a default "PostPaid"
	if v, ok := d.GetOk("instance_charge_type"); ok && strings.ToLower(v.(string)) == "prepaid" {
		return false
	}
	if v, ok := d.GetOk("payment_type"); ok && v.(string) == "Subscription" {
		return false
	}
	return true
}

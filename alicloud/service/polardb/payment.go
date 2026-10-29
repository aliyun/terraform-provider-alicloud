package polardb

import (
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// PayType/TimeType and getChargeType are cross-product symbols owned by other
// flat-layout products. Per the stock-method-disposition standard a product
// package collects the logic it needs instead of importing another product
// package, and the flat copies keep their own definition; these are copies
// limited to the values CreateApplication/DeleteApplication use.
type PayType string

const (
	PrePaid  = PayType("PrePaid")
	PostPaid = PayType("PostPaid")
	Prepaid  = PayType("Prepaid")
	Postpaid = PayType("Postpaid")
)

type TimeType string

const (
	Month = TimeType("Month")
	Year  = TimeType("Year")
)

// getChargeType normalizes a POP PayType response value to the Terraform
// pay_type representation: Postpaid stays PostPaid, anything else is reported
// as PrePaid.
func getChargeType(paymentType string) string {
	if strings.ToLower(paymentType) == strings.ToLower(string(PostPaid)) {
		return string(PostPaid)
	}
	return string(PrePaid)
}

// postPaidDiffSuppressFunc suppresses period diffs for PostPaid applications:
// period only applies to PrePaid, so a change to it is meaningless otherwise.
func postPaidDiffSuppressFunc(k, old, new string, d *schema.ResourceData) bool {
	return d.Get("pay_type").(string) != string(PrePaid)
}

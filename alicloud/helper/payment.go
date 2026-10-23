package helper

import "strings"

// PayType and TimeType, with their enumerated values, were previously defined
// in the flat alicloud package. They are generic (used by many products'
// resources) and had no leaf-package home after the #10858 shared-tool sink,
// so migrating any resource that depends on them into a service package
// requires them to be importable from here. The flat package keeps behavior-
// preserving aliases (type PayType = helper.PayType etc.).
type PayType string

const (
	PrePaid    = PayType("PrePaid")
	PostPaid   = PayType("PostPaid")
	Prepaid    = PayType("Prepaid")
	Postpaid   = PayType("Postpaid")
	Serverless = PayType("Serverless")
)

type TimeType string

const (
	Hour  = TimeType("Hour")
	Day   = TimeType("Day")
	Week  = TimeType("Week")
	Month = TimeType("Month")
	Year  = TimeType("Year")
)

// GetChargeType normalizes a POP PayType response value to the Terraform
// pay_type attribute representation: PostPaid stays PostPaid, anything else
// is reported as PrePaid.
func GetChargeType(paymentType string) string {
	if strings.ToLower(paymentType) == strings.ToLower(string(PostPaid)) {
		return string(PostPaid)
	}
	return string(PrePaid)
}

package alicloud

import "strings"

// getChargeType normalizes a payment type string to the canonical PrePaid /
// PostPaid charge type. It is a shared root-package helper used by adb and
// polardb resources/data-sources that have not yet migrated to their own
// service subpackages. Retained in the root package for those unmigrated
// cross-product callers after the elasticsearch service moved to
// alicloud/service/elasticsearch, which carries its own package-private
// copy for in-product use.
func getChargeType(paymentType string) string {
	if strings.ToLower(paymentType) == strings.ToLower(string(PostPaid)) {
		return string(PostPaid)
	} else {
		return string(PrePaid)
	}
}

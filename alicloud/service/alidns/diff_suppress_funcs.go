package alidns

import (
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// dnsValueDiffSuppressFunc and dnsPriorityDiffSuppressFunc are local 1:1 copies
// of the DNS-specific definitions in alicloud/diff_suppress_funcs.go. The root
// copies are retained for the root-package test consumers in
// alicloud/diff_suppress_funcs_test.go; alidns cannot import the root package.
func dnsValueDiffSuppressFunc(k, old, new string, d *schema.ResourceData) bool {
	switch d.Get("type") {
	case "NS", "MX", "CNAME", "SRV":
		new = strings.TrimSuffix(strings.TrimSpace(new), ".")
	}
	return old == new
}

func dnsPriorityDiffSuppressFunc(k, old, new string, d *schema.ResourceData) bool {
	return d.Get("type").(string) != "MX"
}

// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudCloudFirewallTlsInspectCaCertificateDataSource(t *testing.T) {
	testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
	rand := acctest.RandIntRange(1000000, 9999999)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCloudFirewallTlsInspectCaCertificateSourceConfig(rand, map[string]string{}),
		fakeConfig: testAccCheckAlicloudCloudFirewallTlsInspectCaCertificateSourceConfig(rand, map[string]string{
			"ids": `["fake"]`,
		}),
	}

	CloudFirewallTlsInspectCaCertificateCheckInfo.dataSourceTestCheck(t, rand, idsConf)
}

var existCloudFirewallTlsInspectCaCertificateMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"certificates.#":            "2",
		"certificates.0.ca_cert_id": CHECKSET,
	}
}

var fakeCloudFirewallTlsInspectCaCertificateMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"certificates.#": "0",
	}
}

var CloudFirewallTlsInspectCaCertificateCheckInfo = dataSourceAttr{
	resourceId:   "data.alicloud_cloud_firewall_tls_inspect_ca_certificates.default",
	existMapFunc: existCloudFirewallTlsInspectCaCertificateMapFunc,
	fakeMapFunc:  fakeCloudFirewallTlsInspectCaCertificateMapFunc,
}

func testAccCheckAlicloudCloudFirewallTlsInspectCaCertificateSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
	default = "tf-testAccCloudFirewallTlsInspectCaCertificate%d"
}

data "alicloud_cloud_firewall_tls_inspect_ca_certificates" "default" {
%s
}
`, rand, strings.Join(pairs, "\n   "))
	return config
}

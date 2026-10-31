// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package cloudfirewall_test

import (
	"fmt"
	"strings"
	"testing"

	acctest "github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	connectivity "github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAlicloudCloudFirewallTlsInspectCaCertificateDataSource(t *testing.T) {
	acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
	rand := tfacctest.RandIntRange(1000000, 9999999)

	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudCloudFirewallTlsInspectCaCertificateSourceConfig(rand, map[string]string{}),
		FakeConfig: testAccCheckAlicloudCloudFirewallTlsInspectCaCertificateSourceConfig(rand, map[string]string{
			"ids": `["fake"]`,
		}),
	}

	CloudFirewallTlsInspectCaCertificateCheckInfo.DataSourceTestCheck(t, rand, idsConf)
}

var existCloudFirewallTlsInspectCaCertificateMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"certificates.#":            acctest.REGEXMATCH + "^[1-9]\\d*$",
		"certificates.0.ca_cert_id": acctest.CHECKSET,
	}
}

var fakeCloudFirewallTlsInspectCaCertificateMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"certificates.#": "0",
	}
}

var CloudFirewallTlsInspectCaCertificateCheckInfo = acctest.DataSourceAttr{
	ResourceId:   "data.alicloud_cloud_firewall_tls_inspect_ca_certificates.default",
	ExistMapFunc: existCloudFirewallTlsInspectCaCertificateMapFunc,
	FakeMapFunc:  fakeCloudFirewallTlsInspectCaCertificateMapFunc,
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

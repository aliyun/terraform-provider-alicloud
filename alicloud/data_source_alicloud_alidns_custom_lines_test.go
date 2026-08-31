package alicloud

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudAlidnsCustomLinesDataSource(t *testing.T) {
	rand := acctest.RandIntRange(10000, 99999)
	checkoutSupportedRegions(t, true, connectivity.AlidnsSupportRegions)
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudAlidnsCustomLinesDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_alidns_custom_line.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudAlidnsCustomLinesDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_alidns_custom_line.default.id}_fake"]`,
		}),
	}
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudAlidnsCustomLinesDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_alidns_custom_line.default.custom_line_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudAlidnsCustomLinesDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_alidns_custom_line.default.custom_line_name}_fake"`,
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudAlidnsCustomLinesDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_alidns_custom_line.default.id}"]`,
			"name_regex": `"${alicloud_alidns_custom_line.default.custom_line_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudAlidnsCustomLinesDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_alidns_custom_line.default.id}_fake"]`,
			"name_regex": `"${alicloud_alidns_custom_line.default.custom_line_name}_fake"`,
		}),
	}
	var existAlicloudAlidnsCustomLinesDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                              "1",
			"names.#":                            "1",
			"lines.#":                            "1",
			"lines.0.code":                       CHECKSET,
			"lines.0.id":                         CHECKSET,
			"lines.0.custom_line_id":             CHECKSET,
			"lines.0.custom_line_name":           fmt.Sprintf("tf-testAcc%d", rand),
			"lines.0.domain_name":                CHECKSET,
			"lines.0.ip_segment_list.#":          "1",
			"lines.0.ip_segment_list.0.start_ip": "192.0.2.123",
			"lines.0.ip_segment_list.0.end_ip":   "192.0.2.125",
		}
	}
	var fakeAlicloudAlidnsCustomLinesDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"names.#": "0",
		}
	}
	var alicloudAlidnsCustomLinesCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_alidns_custom_lines.default",
		existMapFunc: existAlicloudAlidnsCustomLinesDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlicloudAlidnsCustomLinesDataSourceNameMapFunc,
	}

	preCheck := func() {
		testAccPreCheck(t)
		testAccPreCheckWithEnvVariable(t, "ALICLOUD_ICP_DOMAIN_NAME")
	}
	alicloudAlidnsCustomLinesCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, nameRegexConf, allConf)
}
func testAccCheckAlicloudAlidnsCustomLinesDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`

variable "name" {	
	default = "tf-testAcc%d"
}

variable "domain_name" {
  default = "%s"
}

resource "alicloud_alidns_custom_line" "default" {
  custom_line_name = var.name
  domain_name = var.domain_name
  ip_segment_list {
		start_ip =  "192.0.2.123"
		end_ip   =  "192.0.2.125"
	}
}

data "alicloud_alidns_custom_lines" "default" {	
	enable_details = true
    domain_name    = var.domain_name
	%s	
}
`, rand, os.Getenv("ALICLOUD_ICP_DOMAIN_NAME"), strings.Join(pairs, " \n "))
	return config
}

package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudRdcOrganizationDataSource(t *testing.T) {
	checkoutSupportedRegions(t, true, connectivity.RDCupportRegions)
	rand := acctest.RandInt()
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudRdcOrganizationDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_rdc_organization.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudRdcOrganizationDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_rdc_organization.default.id}_fake"]`,
		}),
	}

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudRdcOrganizationDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_rdc_organization.default.id}"]`,
			"name_regex": `"${alicloud_rdc_organization.default.organization_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudRdcOrganizationDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_rdc_organization.default.id}_fake"]`,
			"name_regex": `"${alicloud_rdc_organization.default.organization_name}_fake"`,
		}),
	}

	var existAlicloudRdcOrganizationDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":           "1",
			"names.#":         "1",
			"organizations.#": "1",
		}
	}
	var fakeAlicloudRdcOrganizationDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":           "0",
			"names.#":         "0",
			"organizations.#": "0",
		}
	}

	var AlicloudRdcOrganizationCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_rdc_organizations.default",
		existMapFunc: existAlicloudRdcOrganizationDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlicloudRdcOrganizationDataSourceNameMapFunc,
	}
	AlicloudRdcOrganizationCheckInfo.dataSourceTestCheck(t, rand, idsConf, nameRegexConf)
}
func testAccCheckAlicloudRdcOrganizationDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`

variable "name" {	
	default = "tf-testAccOrganizations-%d"
}

resource "alicloud_rdc_organization" "default"{
  organization_name = var.name
  source =            var.name
}

data "alicloud_rdc_organizations" "default" {	
	%s
}
`, rand, strings.Join(pairs, " \n "))
	return config
}

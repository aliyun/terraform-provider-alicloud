package devopsrdc_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

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

var AlicloudRdcOrganizationCheckInfo = acctest.DataSourceAttr{
	ResourceId:   "data.alicloud_rdc_organizations.default",
	ExistMapFunc: existAlicloudRdcOrganizationDataSourceNameMapFunc,
	FakeMapFunc:  fakeAlicloudRdcOrganizationDataSourceNameMapFunc,
}

// data source alicloud_rdc_organizations has been deprecated from version 1.238.0
func SkipTestAccAlicloudRdcOrganizationDataSource(t *testing.T) {
	acctest.CheckoutSupportedRegions(t, true, connectivity.RDCupportRegions)
	rand := tfacctest.RandInt()
	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudRdcOrganizationDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_rdc_organization.default.id}"]`,
		}),
		FakeConfig: testAccCheckAlicloudRdcOrganizationDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_rdc_organization.default.id}_fake"]`,
		}),
	}

	nameRegexConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudRdcOrganizationDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_rdc_organization.default.id}"]`,
			"name_regex": `"${alicloud_rdc_organization.default.organization_name}"`,
		}),
		FakeConfig: testAccCheckAlicloudRdcOrganizationDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_rdc_organization.default.id}_fake"]`,
			"name_regex": `"${alicloud_rdc_organization.default.organization_name}_fake"`,
		}),
	}

	AlicloudRdcOrganizationCheckInfo.DataSourceTestCheck(t, rand, idsConf, nameRegexConf)
}

// TestAccAliCloudRdcOrganizationsDataSource is the representative acceptance test for
// the migrated alicloud_rdc_organizations data source in the devopsrdc subpackage.
func TestAccAliCloudRdcOrganizationsDataSource(t *testing.T) {
	acctest.CheckoutSupportedRegions(t, true, connectivity.RDCupportRegions)
	rand := tfacctest.RandInt()
	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudRdcOrganizationDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_rdc_organization.default.id}"]`,
		}),
		FakeConfig: testAccCheckAlicloudRdcOrganizationDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_rdc_organization.default.id}_fake"]`,
		}),
	}

	nameRegexConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudRdcOrganizationDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_rdc_organization.default.id}"]`,
			"name_regex": `"${alicloud_rdc_organization.default.organization_name}"`,
		}),
		FakeConfig: testAccCheckAlicloudRdcOrganizationDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_rdc_organization.default.id}_fake"]`,
			"name_regex": `"${alicloud_rdc_organization.default.organization_name}_fake"`,
		}),
	}

	AlicloudRdcOrganizationCheckInfo.DataSourceTestCheck(t, rand, idsConf, nameRegexConf)
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

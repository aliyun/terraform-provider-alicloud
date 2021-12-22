package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

// Skip this testcase because of the account cannot purchase non-internal products.
func SkipTestAccAlicloudCenBandwidthLimitsDataSource(t *testing.T) {
	rand := acctest.RandIntRange(1000000, 99999999)
	idConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCenBandwidthLimitsDataSourceConfig(rand, map[string]string{
			"instance_ids": `["${alicloud_cen_bandwidth_limit.default.instance_id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudCenBandwidthLimitsDataSourceConfig(rand, map[string]string{
			"instance_ids": `["${alicloud_cen_bandwidth_limit.default.instance_id}-fake"]`,
		}),
	}

	steps := idConf.buildDataSourceSteps(t, &cenBandwidthLimitsCheckInfo, rand)

	// multi provideris
	var providers []*schema.Provider
	providerFactories := map[string]terraform.ResourceProviderFactory{
		"alicloud": func() (terraform.ResourceProvider, error) {
			p := Provider()
			providers = append(providers, p.(*schema.Provider))
			return p, nil
		},
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithAccountSiteType(t, DomesticSite)
			testAccPreCheckWithRegions(t, true, connectivity.CenNoSkipRegions)
		},

		// module name
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckCenBandwidthLimitDestroyWithProviders(&providers),
		Steps:             steps,
	})
}

func testAccCheckAlicloudCenBandwidthLimitsDataSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`

variable "name" {
	  default = "tf-testAcc%sCenBandwidthLimitsDataSource-%d"
	}

provider "alicloud" {
  alias = "bj"
  region = "cn-beijing"
}

provider "alicloud" {
  alias = "us"
  region = "us-west-1"
}

data "alicloud_vpcs" "default" {
    provider = "alicloud.bj"
	name_regex = "default-NODELETING"
}

data "alicloud_vpcs" "default1" {
	provider = "alicloud.us"
	name_regex = "default-NODELETING"
}

resource "alicloud_cen_instance" "default" {
     name = "${var.name}"
     description = "tf-testAccTerraform01"
}

resource "alicloud_cen_bandwidth_package" "default" {
    name = "${var.name}"
    bandwidth = 20
    geographic_region_ids = [
		"China",
		"North-America"]
}

resource "alicloud_cen_bandwidth_package_attachment" "default" {
    instance_id = "${alicloud_cen_instance.default.id}"
    bandwidth_package_id = "${alicloud_cen_bandwidth_package.default.id}"
}

resource "alicloud_cen_instance_attachment" "default" {
    instance_id = "${alicloud_cen_instance.default.id}"
    child_instance_id = data.alicloud_vpcs.default.ids.0
    child_instance_type = "VPC""
    child_instance_region_id = "cn-beijing"
}

resource "alicloud_cen_instance_attachment" "default1" {
    instance_id = "${alicloud_cen_instance.default.id}"
    child_instance_id = data.alicloud_vpcs.default1.ids.0
    child_instance_type = "VPC""
    child_instance_region_id = "us-west-1"
}

resource "alicloud_cen_bandwidth_limit" "default" {
     instance_id = "${alicloud_cen_instance.default.id}"
     region_ids = ["cn-beijing",
                   "us-west-1"]
     bandwidth_limit = 15
     depends_on = [
        "alicloud_cen_bandwidth_package_attachment.default",
        "alicloud_cen_instance_attachment.default",
        "alicloud_cen_instance_attachment.default1"]
}

data "alicloud_cen_bandwidth_limits" "default" {
	%s
}
`, defaultRegionToTest, rand, strings.Join(pairs, "\n  "))
	return config
}

var existCenBandwidthLimitsMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"instance_ids.#":              "1",
		"limits.#":                    "1",
		"limits.0.local_region_id":    "cn-beijing",
		"limits.0.opposite_region_id": "us-west-1",
		"limits.0.status":             "Active",
		"limits.0.bandwidth_limit":    "15",
	}
}

var fakeCenBandwidthLimitsMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"limits.#": "0",
	}
}

var cenBandwidthLimitsCheckInfo = dataSourceAttr{
	resourceId:   "data.alicloud_cen_bandwidth_limits.default",
	existMapFunc: existCenBandwidthLimitsMapFunc,
	fakeMapFunc:  fakeCenBandwidthLimitsMapFunc,
}

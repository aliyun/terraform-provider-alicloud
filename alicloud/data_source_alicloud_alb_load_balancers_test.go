package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudALBLoadBalancersDataSource(t *testing.T) {
	rand := acctest.RandInt()
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudAlbLoadBalancersDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_alb_load_balancer.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudAlbLoadBalancersDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_alb_load_balancer.default.id}_fake"]`,
		}),
	}
	addressTypeConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudAlbLoadBalancersDataSourceName(rand, map[string]string{
			"ids":          `["${alicloud_alb_load_balancer.default.id}"]`,
			"address_type": `"Internet"`,
		}),
		fakeConfig: testAccCheckAlicloudAlbLoadBalancersDataSourceName(rand, map[string]string{
			"ids":          `["${alicloud_alb_load_balancer.default.id}_fake"]`,
			"address_type": `"Intranet"`,
		}),
	}
	loadBalancerBussinessstatus := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudAlbLoadBalancersDataSourceName(rand, map[string]string{
			"ids":                           `["${alicloud_alb_load_balancer.default.id}"]`,
			"load_balancer_business_status": `"Normal"`,
		}),
		fakeConfig: testAccCheckAlicloudAlbLoadBalancersDataSourceName(rand, map[string]string{
			"ids":                           `["${alicloud_alb_load_balancer.default.id}_fake"]`,
			"load_balancer_business_status": `"Abnormal"`,
		}),
	}
	vpcIdsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudAlbLoadBalancersDataSourceName(rand, map[string]string{
			"ids":     `["${alicloud_alb_load_balancer.default.id}"]`,
			"vpc_ids": `["${alicloud_alb_load_balancer.default.vpc_id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudAlbLoadBalancersDataSourceName(rand, map[string]string{
			"ids":     `["${alicloud_alb_load_balancer.default.id}_fake"]`,
			"vpc_ids": `["${alicloud_alb_load_balancer.default.vpc_id}"]`,
		}),
	}
	vpcIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudAlbLoadBalancersDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_alb_load_balancer.default.id}"]`,
			"vpc_id": `"${alicloud_alb_load_balancer.default.vpc_id}"`,
		}),
		fakeConfig: testAccCheckAlicloudAlbLoadBalancersDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_alb_load_balancer.default.id}_fake"]`,
			"vpc_id": `"${alicloud_alb_load_balancer.default.vpc_id}"`,
		}),
	}
	loadBalancerIdsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudAlbLoadBalancersDataSourceName(rand, map[string]string{
			"load_balancer_ids": `["${alicloud_alb_load_balancer.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudAlbLoadBalancersDataSourceName(rand, map[string]string{
			"load_balancer_ids": `["${alicloud_alb_load_balancer.default.id}_fake"]`,
		}),
	}
	resourceGroupIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudAlbLoadBalancersDataSourceName(rand, map[string]string{
			"ids":               `["${alicloud_alb_load_balancer.default.id}"]`,
			"resource_group_id": `"${alicloud_alb_load_balancer.default.resource_group_id}"`,
		}),
		fakeConfig: testAccCheckAlicloudAlbLoadBalancersDataSourceName(rand, map[string]string{
			"ids":               `["${alicloud_alb_load_balancer.default.id}_fake"]`,
			"resource_group_id": `"${alicloud_alb_load_balancer.default.resource_group_id}_fake"`,
		}),
	}
	tagsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudAlbLoadBalancersDataSourceName(rand, map[string]string{
			"ids":  `["${alicloud_alb_load_balancer.default.id}"]`,
			"tags": `{Created = "TF"}`,
		}),
		fakeConfig: testAccCheckAlicloudAlbLoadBalancersDataSourceName(rand, map[string]string{
			"ids":  `["${alicloud_alb_load_balancer.default.id}_fake"]`,
			"tags": `{Created = "TF_fake"}`,
		}),
	}
	loadBalancerNameConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudAlbLoadBalancersDataSourceName(rand, map[string]string{
			"load_balancer_name": `"${alicloud_alb_load_balancer.default.load_balancer_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudAlbLoadBalancersDataSourceName(rand, map[string]string{
			"load_balancer_name": `"${alicloud_alb_load_balancer.default.load_balancer_name}_fake"`,
		}),
	}
	statusConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudAlbLoadBalancersDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_alb_load_balancer.default.id}"]`,
			"status": `"Active"`,
		}),
		fakeConfig: testAccCheckAlicloudAlbLoadBalancersDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_alb_load_balancer.default.id}"]`,
			"status": `"Configuring"`,
		}),
	}
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudAlbLoadBalancersDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_alb_load_balancer.default.load_balancer_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudAlbLoadBalancersDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_alb_load_balancer.default.load_balancer_name}_fake"`,
		}),
	}
	zoneIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudAlbLoadBalancersDataSourceName(rand, map[string]string{
			"ids":     `["${alicloud_alb_load_balancer.default.id}"]`,
			"zone_id": `"${data.alicloud_alb_zones.default.zones.0.id}"`,
		}),
		fakeConfig: testAccCheckAlicloudAlbLoadBalancersDataSourceName(rand, map[string]string{
			"ids":     `["${alicloud_alb_load_balancer.default.id}_fake"]`,
			"zone_id": `"${data.alicloud_alb_zones.default.zones.0.id}"`,
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudAlbLoadBalancersDataSourceName(rand, map[string]string{
			"ids":                            `["${alicloud_alb_load_balancer.default.id}"]`,
			"name_regex":                     `"${alicloud_alb_load_balancer.default.load_balancer_name}"`,
			"vpc_id":                         `"${alicloud_alb_load_balancer.default.vpc_id}"`,
			"zone_id":                        `"${data.alicloud_alb_zones.default.zones.0.id}"`,
			"load_balancer_name":             `"${alicloud_alb_load_balancer.default.load_balancer_name}"`,
			"status":                         `"Active"`,
			"tags":                           `{Created = "TF"}`,
			"resource_group_id":              `"${alicloud_alb_load_balancer.default.resource_group_id}"`,
			"load_balancer_ids":              `["${alicloud_alb_load_balancer.default.id}"]`,
			"load_balancer_business_status":  `"Normal"`,
			"load_balancer_bussiness_status": `"Normal"`,
			"address_type":                   `"Internet"`,
		}),
		fakeConfig: testAccCheckAlicloudAlbLoadBalancersDataSourceName(rand, map[string]string{
			"ids":                            `["${alicloud_alb_load_balancer.default.id}_fake"]`,
			"name_regex":                     `"${alicloud_alb_load_balancer.default.load_balancer_name}_fake"`,
			"vpc_id":                         `"${alicloud_alb_load_balancer.default.vpc_id}"`,
			"zone_id":                        `"${data.alicloud_alb_zones.default.zones.0.id}"`,
			"load_balancer_name":             `"${alicloud_alb_load_balancer.default.load_balancer_name}_fake"`,
			"status":                         `"Configuring"`,
			"tags":                           `{Created = "TF_fake"}`,
			"resource_group_id":              `"${alicloud_alb_load_balancer.default.resource_group_id}_fake"`,
			"load_balancer_ids":              `["${alicloud_alb_load_balancer.default.id}_fake"]`,
			"load_balancer_business_status":  `"Abnormal"`,
			"load_balancer_bussiness_status": `"Abnormal"`,
			"address_type":                   `"Intranet"`,
		}),
	}
	var existAlicloudAlbLoadBalancersDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                              "1",
			"names.#":                            "1",
			"balancers.#":                        "1",
			"balancers.0.load_balancer_name":     fmt.Sprintf("tf-testAccLoadBalancer-%d", rand),
			"balancers.0.address_type":           "Internet",
			"balancers.0.address_allocated_mode": "Fixed",
			"balancers.0.load_balancer_edition":  "Basic",
		}
	}
	var fakeAlicloudAlbLoadBalancersDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":       "0",
			"names.#":     "0",
			"balancers.#": "0",
		}
	}
	var AlicloudAlbLoadBalancersCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_alb_load_balancers.default",
		existMapFunc: existAlicloudAlbLoadBalancersDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlicloudAlbLoadBalancersDataSourceNameMapFunc,
	}
	preCheck := func() {
		testAccPreCheckWithRegions(t, true, connectivity.AlbSupportRegions)
	}
	AlicloudAlbLoadBalancersCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, addressTypeConf, loadBalancerBussinessstatus, vpcIdsConf, vpcIdConf, loadBalancerIdsConf, statusConf, resourceGroupIdConf, tagsConf, loadBalancerNameConf, nameRegexConf, zoneIdConf, allConf)
}
func testAccCheckAlicloudAlbLoadBalancersDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`

variable "name" {	
	default = "tf-testAccLoadBalancer-%d"
}

data "alicloud_alb_zones" "default"{}

data "alicloud_vpcs" "default" {
 name_regex = "default-NODELETING"
}
data "alicloud_vswitches" "default_1" {
  vpc_id = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_alb_zones.default.zones.0.id
}
resource "alicloud_vswitch" "vswitch_1" {
  count             = length(data.alicloud_vswitches.default_1.ids) > 0 ? 0 : 1
  vpc_id            = data.alicloud_vpcs.default.ids.0
  cidr_block        = cidrsubnet(data.alicloud_vpcs.default.vpcs[0].cidr_block, 8, 2)
  zone_id =  data.alicloud_alb_zones.default.zones.0.id
  vswitch_name              = var.name
}

data "alicloud_vswitches" "default_2" {
  vpc_id = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_alb_zones.default.zones.1.id
}
resource "alicloud_vswitch" "vswitch_2" {
  count             = length(data.alicloud_vswitches.default_2.ids) > 0 ? 0 : 1
  vpc_id            = data.alicloud_vpcs.default.ids.0
  cidr_block        = cidrsubnet(data.alicloud_vpcs.default.vpcs[0].cidr_block, 8, 4)
  zone_id = data.alicloud_alb_zones.default.zones.1.id
  vswitch_name              = var.name
}
data "alicloud_resource_manager_resource_groups" "default" {}

resource "alicloud_alb_load_balancer" "default" {
  vpc_id =                data.alicloud_vpcs.default.ids.0
  address_type =        "Internet"
  address_allocated_mode = "Fixed"
  load_balancer_name =    var.name
  load_balancer_edition = "Basic"
  resource_group_id = data.alicloud_resource_manager_resource_groups.default.groups.0.id
  load_balancer_billing_config {
    pay_type = 	"PayAsYouGo"
  }
  tags = {
		Created = "TF"
  }
  zone_mappings{
		vswitch_id =  length(data.alicloud_vswitches.default_1.ids) > 0 ? data.alicloud_vswitches.default_1.ids[0] : concat(alicloud_vswitch.vswitch_1.*.id, [""])[0]
		zone_id =  data.alicloud_alb_zones.default.zones.0.id
	}
  zone_mappings{
		vswitch_id = length(data.alicloud_vswitches.default_2.ids) > 0 ? data.alicloud_vswitches.default_2.ids[0] : concat(alicloud_vswitch.vswitch_2.*.id, [""])[0]
		zone_id =   data.alicloud_alb_zones.default.zones.1.id
	}
  modification_protection_config{
	status = "NonProtection"
  }
}

data "alicloud_alb_load_balancers" "default" {	
	%s
}
`, rand, strings.Join(pairs, " \n "))
	return config
}

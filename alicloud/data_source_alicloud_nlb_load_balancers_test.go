package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudNlbLoadBalancersDataSource(t *testing.T) {
	rand := acctest.RandInt()
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudNlbLoadBalancersDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_nlb_load_balancer.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudNlbLoadBalancersDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_nlb_load_balancer.default.id}_fake"]`,
		}),
	}
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudNlbLoadBalancersDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_nlb_load_balancer.default.load_balancer_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudNlbLoadBalancersDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_nlb_load_balancer.default.load_balancer_name}_fake"`,
		}),
	}
	addressIpVersionConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudNlbLoadBalancersDataSourceName(rand, map[string]string{
			"address_ip_version": `"${alicloud_nlb_load_balancer.default.address_ip_version}"`,
			"ids":                `["${alicloud_nlb_load_balancer.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudNlbLoadBalancersDataSourceName(rand, map[string]string{
			"address_ip_version": `"DualStack"`,
			"ids":                `["${alicloud_nlb_load_balancer.default.id}"]`,
		}),
	}
	addressTypeConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudNlbLoadBalancersDataSourceName(rand, map[string]string{
			"address_type": `"${alicloud_nlb_load_balancer.default.address_type}"`,
			"ids":          `["${alicloud_nlb_load_balancer.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudNlbLoadBalancersDataSourceName(rand, map[string]string{
			"address_type": `"Intranet"`,
			"ids":          `["${alicloud_nlb_load_balancer.default.id}"]`,
		}),
	}
	dnsNameConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudNlbLoadBalancersDataSourceName(rand, map[string]string{
			"dns_name": `"${alicloud_nlb_load_balancer.default.dns_name}"`,
			"ids":      `["${alicloud_nlb_load_balancer.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudNlbLoadBalancersDataSourceName(rand, map[string]string{
			"dns_name": `"${alicloud_nlb_load_balancer.default.dns_name}_fake"`,
			"ids":      `["${alicloud_nlb_load_balancer.default.id}"]`,
		}),
	}
	loadBalancerBusinessStatusConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudNlbLoadBalancersDataSourceName(rand, map[string]string{
			"load_balancer_business_status": `"${alicloud_nlb_load_balancer.default.load_balancer_business_status}"`,
			"ids":                           `["${alicloud_nlb_load_balancer.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudNlbLoadBalancersDataSourceName(rand, map[string]string{
			"load_balancer_business_status": `"Abnormal"`,
			"ids":                           `["${alicloud_nlb_load_balancer.default.id}"]`,
		}),
	}
	resourceGroupIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudNlbLoadBalancersDataSourceName(rand, map[string]string{
			"resource_group_id": `"${alicloud_nlb_load_balancer.default.resource_group_id}"`,
			"ids":               `["${alicloud_nlb_load_balancer.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudNlbLoadBalancersDataSourceName(rand, map[string]string{
			"resource_group_id": `"${alicloud_nlb_load_balancer.default.resource_group_id}_fake"`,
			"ids":               `["${alicloud_nlb_load_balancer.default.id}"]`,
		}),
	}
	vpcIdsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudNlbLoadBalancersDataSourceName(rand, map[string]string{
			"vpc_ids": `["${alicloud_nlb_load_balancer.default.vpc_id}"]`,
			"ids":     `["${alicloud_nlb_load_balancer.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudNlbLoadBalancersDataSourceName(rand, map[string]string{
			"vpc_ids": `["${alicloud_nlb_load_balancer.default.vpc_id}_fake"]`,
			"ids":     `["${alicloud_nlb_load_balancer.default.id}"]`,
		}),
	}
	statusConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudNlbLoadBalancersDataSourceName(rand, map[string]string{
			"status": `"${alicloud_nlb_load_balancer.default.status}"`,
			"ids":    `["${alicloud_nlb_load_balancer.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudNlbLoadBalancersDataSourceName(rand, map[string]string{
			"status": `"Deleted"`,
			"ids":    `["${alicloud_nlb_load_balancer.default.id}"]`,
		}),
	}
	loadBalancerNamesConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudNlbLoadBalancersDataSourceName(rand, map[string]string{
			"load_balancer_names": `["${alicloud_nlb_load_balancer.default.load_balancer_name}"]`,
			"ids":                 `["${alicloud_nlb_load_balancer.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudNlbLoadBalancersDataSourceName(rand, map[string]string{
			"load_balancer_names": `["${alicloud_nlb_load_balancer.default.load_balancer_name}_fake"]`,
			"ids":                 `["${alicloud_nlb_load_balancer.default.id}"]`,
		}),
	}
	zoneIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudNlbLoadBalancersDataSourceName(rand, map[string]string{
			"zone_id": `"${local.zone_id_1}"`,
			"ids":     `["${alicloud_nlb_load_balancer.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudNlbLoadBalancersDataSourceName(rand, map[string]string{
			"zone_id": `"${local.zone_id_1}_fake"`,
			"ids":     `["${alicloud_nlb_load_balancer.default.id}"]`,
		}),
	}
	tagsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudNlbLoadBalancersDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_nlb_load_balancer.default.id}"]`,
			"tags": `{ 
						"Created" = "tfTestAcc0"
						"For"     =  "Tftestacc 0"
					}`,
		}),
		fakeConfig: testAccCheckAlicloudNlbLoadBalancersDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_nlb_load_balancer.default.id}"]`,
			"tags": `{ 
						"Created" = "tfTestAcc0_fake"
						"For"     =  "Tftestacc0fake"
					}`,
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudNlbLoadBalancersDataSourceName(rand, map[string]string{
			"name_regex":                    `"${alicloud_nlb_load_balancer.default.load_balancer_name}"`,
			"load_balancer_business_status": `"${alicloud_nlb_load_balancer.default.load_balancer_business_status}"`,
			"tags": `{ 
						"Created" = "tfTestAcc0"
						"For"     =  "Tftestacc 0"
					}`,
			"dns_name":            `"${alicloud_nlb_load_balancer.default.dns_name}"`,
			"resource_group_id":   `"${alicloud_nlb_load_balancer.default.resource_group_id}"`,
			"status":              `"${alicloud_nlb_load_balancer.default.status}"`,
			"ids":                 `["${alicloud_nlb_load_balancer.default.id}"]`,
			"address_ip_version":  `"${alicloud_nlb_load_balancer.default.address_ip_version}"`,
			"address_type":        `"${alicloud_nlb_load_balancer.default.address_type}"`,
			"vpc_ids":             `["${alicloud_nlb_load_balancer.default.vpc_id}"]`,
			"load_balancer_names": `["${alicloud_nlb_load_balancer.default.load_balancer_name}"]`,
			"zone_id":             `"${local.zone_id_1}"`,
		}),
		fakeConfig: testAccCheckAlicloudNlbLoadBalancersDataSourceName(rand, map[string]string{
			"status": `"Deleted"`,
			"tags": `{ 
						"Created" = "tfTestAcc0_fake"
						"For"     =  "Tftestacc0fake"
					}`,
			"ids":                           `["${alicloud_nlb_load_balancer.default.id}_fake"]`,
			"name_regex":                    `"${alicloud_nlb_load_balancer.default.load_balancer_name}_fake"`,
			"address_type":                  `"Intranet"`,
			"load_balancer_business_status": `"Abnormal"`,
			"resource_group_id":             `"${alicloud_nlb_load_balancer.default.resource_group_id}_fake"`,
			"address_ip_version":            `"DualStack"`,
			"dns_name":                      `"${alicloud_nlb_load_balancer.default.dns_name}_fake"`,
			"vpc_ids":                       `["${alicloud_nlb_load_balancer.default.vpc_id}_fake"]`,
			"load_balancer_names":           `["${alicloud_nlb_load_balancer.default.load_balancer_name}_fake"]`,
			"zone_id":                       `"${local.zone_id_1}_fake"`,
		}),
	}
	var existAlicloudNlbLoadBalancersDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                                            "1",
			"names.#":                                          "1",
			"balancers.#":                                      "1",
			"balancers.0.address_ip_version":                   "Ipv4",
			"balancers.0.address_type":                         "Internet",
			"balancers.0.cross_zone_enabled":                   "true",
			"balancers.0.load_balancer_name":                   CHECKSET,
			"balancers.0.load_balancer_type":                   "Network",
			"balancers.0.resource_group_id":                    CHECKSET,
			"balancers.0.tags.%":                               "2",
			"balancers.0.tags.Created":                         "tfTestAcc0",
			"balancers.0.tags.For":                             "Tftestacc 0",
			"balancers.0.vpc_id":                               CHECKSET,
			"balancers.0.dns_name":                             CHECKSET,
			"balancers.0.zone_mappings.#":                      "2",
			"balancers.0.zone_mappings.0.allocation_id":        CHECKSET,
			"balancers.0.zone_mappings.0.eni_id":               CHECKSET,
			"balancers.0.zone_mappings.0.ipv6_address":         "",
			"balancers.0.zone_mappings.0.private_ipv4_address": CHECKSET,
			"balancers.0.zone_mappings.0.public_ipv4_address":  CHECKSET,
			"balancers.0.zone_mappings.0.vswitch_id":           CHECKSET,
			"balancers.0.zone_mappings.0.zone_id":              CHECKSET,
			"balancers.0.zone_mappings.1.allocation_id":        CHECKSET,
			"balancers.0.zone_mappings.1.eni_id":               CHECKSET,
			"balancers.0.zone_mappings.1.ipv6_address":         "",
			"balancers.0.zone_mappings.1.private_ipv4_address": CHECKSET,
			"balancers.0.zone_mappings.1.public_ipv4_address":  CHECKSET,
			"balancers.0.zone_mappings.1.vswitch_id":           CHECKSET,
			"balancers.0.zone_mappings.1.zone_id":              CHECKSET,
			"balancers.0.id":                                   CHECKSET,
			"balancers.0.bandwidth_package_id":                 "",
			"balancers.0.create_time":                          CHECKSET,
			"balancers.0.ipv6_address_type":                    CHECKSET,
			"balancers.0.load_balancer_business_status":        CHECKSET,
			"balancers.0.load_balancer_id":                     CHECKSET,
			"balancers.0.security_group_ids.#":                 "0",
			"balancers.0.status":                               CHECKSET,
			"balancers.0.operation_locks.#":                    "0",
		}
	}
	var fakeAlicloudNlbLoadBalancersDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":       "0",
			"names.#":     "0",
			"balancers.#": "0",
		}
	}
	var AlicloudNlbLoadBalancersCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_nlb_load_balancers.default",
		existMapFunc: existAlicloudNlbLoadBalancersDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlicloudNlbLoadBalancersDataSourceNameMapFunc,
	}

	preCheck := func() {
		testAccPreCheck(t)
	}
	AlicloudNlbLoadBalancersCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, nameRegexConf, addressIpVersionConf, addressTypeConf, dnsNameConf, loadBalancerBusinessStatusConf, resourceGroupIdConf, resourceGroupIdConf, vpcIdsConf, statusConf, loadBalancerNamesConf, zoneIdConf, tagsConf, allConf)
}
func testAccCheckAlicloudNlbLoadBalancersDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`

variable "name" {
	default = "tf-testAccLoadBalancer-%d"
}

data "alicloud_nlb_zones" "default" {}

data "alicloud_vpcs" "default" {
    name_regex = "^default-NODELETING$"
}

data "alicloud_resource_manager_resource_groups" "default" {}

data "alicloud_vswitches" "default_1" {
  vpc_id  = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_nlb_zones.default.zones.0.id
}
data "alicloud_vswitches" "default_2" {
  vpc_id  = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_nlb_zones.default.zones.1.id
}
locals {
  zone_id_1    = data.alicloud_nlb_zones.default.zones.0.id
  vswitch_id_1 = data.alicloud_vswitches.default_1.ids[0]
  zone_id_2    = data.alicloud_nlb_zones.default.zones.1.id
  vswitch_id_2 = data.alicloud_vswitches.default_2.ids[0]
}

resource "alicloud_nlb_load_balancer" "default" {
  load_balancer_name = var.name
  resource_group_id  = data.alicloud_resource_manager_resource_groups.default.ids.0
  load_balancer_type = "Network"
  address_type       = "Internet"
  address_ip_version = "Ipv4"
  tags               = {
    Created = "tfTestAcc0"
    For     = "Tftestacc 0"
  }
  vpc_id = data.alicloud_vpcs.default.ids.0
  zone_mappings {
    vswitch_id = local.vswitch_id_1
    zone_id    = local.zone_id_1
  }
  zone_mappings {
    vswitch_id = local.vswitch_id_2
    zone_id    = local.zone_id_2
  }
}

data "alicloud_nlb_load_balancers" "default" {
	%s
}
`, rand, strings.Join(pairs, " \n "))
	return config
}

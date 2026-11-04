package nlb_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/nlb"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAliCloudNlbLoadBalancer_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_nlb_load_balancer.default"
	acctest.CheckoutSupportedRegions(t, true, connectivity.NLBSupportRegions)
	ra := acctest.ResourceAttrInit(resourceId, AlicloudNLBLoadBalancerMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return nlb.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeNlbLoadBalancer")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%snlbloadbalancer%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudNLBLoadBalancerBasicDependence0)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"load_balancer_name":             "${var.name}",
					"resource_group_id":              "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"load_balancer_type":             "Network",
					"address_type":                   "Internet",
					"address_ip_version":             "DualStack",
					"payment_type":                   "PayAsYouGo",
					"cps":                            "300",
					"vpc_id":                         "${alicloud_vpc.default.id}",
					"deletion_protection_enabled":    "false",
					"modification_protection_status": "NonProtection",
					"tags": map[string]string{
						"Created": "tfTestAcc0",
						"For":     "Tftestacc 0",
					},
					"zone_mappings": []map[string]interface{}{
						{
							"vswitch_id":           "${local.vswitch_id_1}",
							"zone_id":              "${local.zone_id_1}",
							"ipv6_address":         "${cidrhost(alicloud_vswitch.default.0.ipv6_cidr_block, 128)}",
							"private_ipv4_address": "${cidrhost(alicloud_vswitch.default.0.cidr_block, 100)}",
							"ipv6_local_addresses": []string{
								"${cidrhost(alicloud_vswitch.default.0.ipv6_cidr_block, 32)}",
								"${cidrhost(alicloud_vswitch.default.0.ipv6_cidr_block, 64)}",
							},
							"ipv4_local_addresses": []string{
								"${cidrhost(alicloud_vswitch.default.0.cidr_block, 32)}",
								"${cidrhost(alicloud_vswitch.default.0.cidr_block, 64)}",
							},
						},
						{
							"vswitch_id": "${local.vswitch_id_2}",
							"zone_id":    "${local.zone_id_2}",
						},
					},
					"cross_zone_enabled": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"load_balancer_name":             name,
						"resource_group_id":              acctest.CHECKSET,
						"load_balancer_type":             "Network",
						"address_type":                   "Internet",
						"address_ip_version":             "DualStack",
						"vpc_id":                         acctest.CHECKSET,
						"deletion_protection_enabled":    "false",
						"modification_protection_status": "NonProtection",
						"tags.%":                         "2",
						"tags.Created":                   "tfTestAcc0",
						"tags.For":                       "Tftestacc 0",
						"zone_mappings.#":                "2",
						"cross_zone_enabled":             "false",
						"cps":                            "300",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"cps": "400",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cps": "400",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"ipv6_address_type": "Internet",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ipv6_address_type": "Internet",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"zone_mappings": []map[string]interface{}{
						{
							"vswitch_id":           "${local.vswitch_id_1}",
							"zone_id":              "${local.zone_id_1}",
							"allocation_id":        "${alicloud_eip.zone_b.id}",
							"private_ipv4_address": "${cidrhost(alicloud_vswitch.default.0.cidr_block, 90)}",
						},
						{
							"vswitch_id": "${local.vswitch_id_2}",
							"zone_id":    "${local.zone_id_2}",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"deletion_protection_enabled":    "true",
					"modification_protection_status": "ConsoleProtection",
					"deletion_protection_reason":     "tf-open",
					"modification_protection_reason": "tf-open",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"deletion_protection_enabled":    "true",
						"modification_protection_status": "ConsoleProtection",
						"deletion_protection_reason":     "tf-open",
						"modification_protection_reason": "tf-open",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"deletion_protection_reason":     "tf-open-update",
					"modification_protection_reason": "tf-open-update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"deletion_protection_reason":     "tf-open-update",
						"modification_protection_reason": "tf-open-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"deletion_protection_enabled":    "false",
					"modification_protection_status": "NonProtection",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"deletion_protection_enabled":    "false",
						"modification_protection_status": "NonProtection",
						"deletion_protection_reason":     "",
						"modification_protection_reason": "",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAliCloudNlbLoadBalancer_basic1(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_nlb_load_balancer.default"
	acctest.CheckoutSupportedRegions(t, true, connectivity.NLBSupportRegions)
	ra := acctest.ResourceAttrInit(resourceId, AlicloudNLBLoadBalancerMap1)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return nlb.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeNlbLoadBalancer")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%snlbloadbalancer%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudNLBLoadBalancerBasicDependence1)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"load_balancer_name": "${var.name}",
					"address_type":       "Internet",
					"vpc_id":             "${alicloud_vpc.default.id}",
					"zone_mappings": []map[string]interface{}{
						{
							"vswitch_id": "${local.vswitch_id_1}",
							"zone_id":    "${local.zone_id_1}",
						},
						{
							"vswitch_id": "${local.vswitch_id_2}",
							"zone_id":    "${local.zone_id_2}",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"load_balancer_name": name,
						"address_type":       "Internet",
						"vpc_id":             acctest.CHECKSET,
						"zone_mappings.#":    "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"load_balancer_name": "${var.name}_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"load_balancer_name": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "tfTestAcc0",
						"For":     "Tftestacc 0",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "tfTestAcc0",
						"tags.For":     "Tftestacc 0",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"zone_mappings": []map[string]interface{}{
						{
							"vswitch_id": "${local.vswitch_id_1}",
							"zone_id":    "${local.zone_id_1}",
						},
						{
							"vswitch_id": "${local.vswitch_id_2}",
							"zone_id":    "${local.zone_id_2}",
						},
						{
							"vswitch_id": "${local.vswitch_id_3}",
							"zone_id":    "${local.zone_id_3}",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"zone_mappings.#": "3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"bandwidth_package_id": "${alicloud_common_bandwidth_package.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bandwidth_package_id": acctest.CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"bandwidth_package_id": acctest.REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bandwidth_package_id": acctest.REMOVEKEY,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"address_type": "Intranet",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"address_type": "Intranet",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"cross_zone_enabled": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cross_zone_enabled": "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"load_balancer_name": "${var.name}",
					"zone_mappings": []map[string]interface{}{
						{
							"vswitch_id": "${local.vswitch_id_1}",
							"zone_id":    "${local.zone_id_1}",
						},
						{
							"vswitch_id": "${local.vswitch_id_2}",
							"zone_id":    "${local.zone_id_2}",
						},
					},
					"cross_zone_enabled": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"load_balancer_name": name,
						"zone_mappings.#":    "2",
						"cross_zone_enabled": "true",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

var AlicloudNLBLoadBalancerMap0 = map[string]string{
	"cross_zone_enabled":            acctest.CHECKSET,
	"load_balancer_type":            acctest.CHECKSET,
	"status":                        acctest.CHECKSET,
	"address_ip_version":            acctest.CHECKSET,
	"load_balancer_name":            acctest.CHECKSET,
	"vpc_id":                        acctest.CHECKSET,
	"zone_mappings.#":               acctest.CHECKSET,
	"address_type":                  acctest.CHECKSET,
	"resource_group_id":             acctest.CHECKSET,
	"dns_name":                      acctest.CHECKSET,
	"load_balancer_business_status": acctest.CHECKSET,
	"tags.%":                        acctest.CHECKSET,
}

func AlicloudNLBLoadBalancerBasicDependence0(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}

data "alicloud_nlb_zones" "default" {
}

data "alicloud_resource_manager_resource_groups" "default" {
}

resource "alicloud_eip" "zone_a" {
  bandwidth            = "10"
  internet_charge_type = "PayByTraffic"
}

resource "alicloud_eip" "zone_b" {
  bandwidth            = "10"
  internet_charge_type = "PayByTraffic"
}

resource "alicloud_vpc" "default" {
  name        = var.name
  cidr_block  = "172.16.0.0/16"
  enable_ipv6 = true
}

resource "alicloud_vpc_ipv6_gateway" "default" {
  description       = var.name
  ipv6_gateway_name = var.name
  vpc_id            = alicloud_vpc.default.id
}

resource "alicloud_vswitch" "default" {
  count                = length(data.alicloud_nlb_zones.default.zones)
  vpc_id               = alicloud_vpc.default.id
  cidr_block           = cidrsubnet(alicloud_vpc.default.cidr_block, 3, count.index)
  zone_id              = data.alicloud_nlb_zones.default.zones[count.index].id
  ipv6_cidr_block_mask = count.index + 60
  enable_ipv6          = true
}

locals {
  zone_id_1    = data.alicloud_nlb_zones.default.zones.0.id
  vswitch_id_1 = alicloud_vswitch.default.0.id
  zone_id_2    = data.alicloud_nlb_zones.default.zones.1.id
  vswitch_id_2 = alicloud_vswitch.default.1.id
  zone_id_3    = data.alicloud_nlb_zones.default.zones.2.id
  vswitch_id_3 = alicloud_vswitch.default.2.id
}

resource "alicloud_common_bandwidth_package" "default" {
  bandwidth            = 2
  internet_charge_type = "PayByBandwidth"
  name                 = "${var.name}"
  description          = "${var.name}_description"
}
`, name)
}

func AlicloudNLBLoadBalancerBasicDependence1(name string) string {
	return strings.ReplaceAll(AlicloudNLBLoadBalancerBasicDependence0(name), "172.16.", "172.17.")
}

var AlicloudNLBLoadBalancerMap1 = map[string]string{
	"address_type":       acctest.CHECKSET,
	"resource_group_id":  acctest.CHECKSET,
	"status":             acctest.CHECKSET,
	"cross_zone_enabled": acctest.CHECKSET,
	"load_balancer_type": acctest.CHECKSET,
	"vpc_id":             acctest.CHECKSET,
	"zone_mappings.#":    acctest.CHECKSET,
	"address_ip_version": acctest.CHECKSET,
	"load_balancer_name": acctest.CHECKSET,
}

// lintignore: R001

func TestAccAliCloudNlbLoadBalancer_basic3678(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_nlb_load_balancer.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudNlbLoadBalancerMap3678)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return nlb.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeNlbLoadBalancer")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%snlbloadbalancer%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudNlbLoadBalancerBasicDependence3678)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, connectivity.NLBSupportRegions)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"zone_mappings": []map[string]interface{}{
						{
							"vswitch_id": "${alicloud_vswitch.vsj.id}",
							"zone_id":    "${alicloud_vswitch.vsj.zone_id}",
						},
						{
							"vswitch_id": "${alicloud_vswitch.vsk.id}",
							"zone_id":    "${alicloud_vswitch.vsk.zone_id}",
						},
						{
							"vswitch_id": "${alicloud_vswitch.vsg.id}",
							"zone_id":    "${alicloud_vswitch.vsg.zone_id}",
						},
					},
					"address_type":       "Intranet",
					"vpc_id":             "${alicloud_vpc.vpc.id}",
					"load_balancer_name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"zone_mappings.#":    "3",
						"address_type":       "Intranet",
						"vpc_id":             acctest.CHECKSET,
						"load_balancer_name": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id": acctest.CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"security_group_ids": []string{},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"security_group_ids.#": "0",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"deletion_protection_config": []map[string]interface{}{
						{
							"enabled": "true",
							"reason":  "tf-open-update",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"modification_protection_config": []map[string]interface{}{
						{
							"status": "ConsoleProtection",
							"reason": "tf-open-update",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"modification_protection_config": []map[string]interface{}{
						{
							"status": "NonProtection",
						},
					},
					"deletion_protection_config": []map[string]interface{}{
						{
							"enabled": "false",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"cross_zone_enabled": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cross_zone_enabled": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"zone_mappings": []map[string]interface{}{
						{
							"vswitch_id": "${alicloud_vswitch.vsj.id}",
							"zone_id":    "${alicloud_vswitch.vsj.zone_id}",
						},
						{
							"vswitch_id": "${alicloud_vswitch.vsk.id}",
							"zone_id":    "${alicloud_vswitch.vsk.zone_id}",
						},
						{
							"vswitch_id": "${alicloud_vswitch.vsg.id}",
							"zone_id":    "${alicloud_vswitch.vsg.zone_id}",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"zone_mappings.#": "3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"address_type": "Intranet",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"address_type": "Intranet",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"load_balancer_name": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"load_balancer_name": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"load_balancer_name": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"load_balancer_name": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id": acctest.CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"security_group_ids": []string{},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"security_group_ids.#": "0",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"deletion_protection_config": []map[string]interface{}{
						{
							"enabled": "true",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"modification_protection_config": []map[string]interface{}{
						{
							"status": "ConsoleProtection",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"security_group_ids": []string{},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"security_group_ids.#": "0",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"deletion_protection_config": []map[string]interface{}{
						{},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"cross_zone_enabled": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cross_zone_enabled": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"load_balancer_name": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"load_balancer_name": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.1}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id": acctest.CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"security_group_ids": []string{},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"security_group_ids.#": "0",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"deletion_protection_config": []map[string]interface{}{
						{},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"modification_protection_config": []map[string]interface{}{
						{
							"status": "NonProtection",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"load_balancer_name": name + "_update",
					"zone_mappings": []map[string]interface{}{
						{
							"vswitch_id": "${alicloud_vswitch.vsj.id}",
							"zone_id":    "${alicloud_vswitch.vsj.zone_id}",
						},
						{
							"vswitch_id": "${alicloud_vswitch.vsk.id}",
							"zone_id":    "${alicloud_vswitch.vsk.zone_id}",
						},
						{
							"vswitch_id": "${alicloud_vswitch.vsg.id}",
							"zone_id":    "${alicloud_vswitch.vsg.zone_id}",
						},
					},
					"address_type":       "Intranet",
					"address_ip_version": "Ipv4",
					"load_balancer_type": "Network",
					"vpc_id":             "${alicloud_vpc.vpc.id}",
					"resource_group_id":  "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"security_group_ids": []string{},
					"deletion_protection_config": []map[string]interface{}{
						{
							"enabled": "false",
						},
					},
					"modification_protection_config": []map[string]interface{}{
						{
							"status": "NonProtection",
						},
					},
					"cross_zone_enabled": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"load_balancer_name":   name + "_update",
						"zone_mappings.#":      "3",
						"address_type":         "Intranet",
						"address_ip_version":   "Ipv4",
						"load_balancer_type":   "Network",
						"vpc_id":               acctest.CHECKSET,
						"resource_group_id":    acctest.CHECKSET,
						"security_group_ids.#": "0",
						"cross_zone_enabled":   "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": acctest.REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": acctest.REMOVEKEY,
						"tags.For":     acctest.REMOVEKEY,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

var AlicloudNlbLoadBalancerMap3678 = map[string]string{
	"load_balancer_type": "Network",
	"status":             acctest.CHECKSET,
	"create_time":        acctest.CHECKSET,
}

func AlicloudNlbLoadBalancerBasicDependence3678(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

data "alicloud_nlb_zones" "default" {
}

resource "alicloud_vpc" "vpc" {
  vpc_name = var.name

  cidr_block = "192.168.0.0/16"
}

resource "alicloud_vswitch" "vsj" {
  vpc_id       = alicloud_vpc.vpc.id
  zone_id      = data.alicloud_nlb_zones.default.zones.0.id
  cidr_block   = "192.168.10.0/24"
  vswitch_name = var.name

}

resource "alicloud_vswitch" "vsk" {
  vpc_id       = alicloud_vpc.vpc.id
  zone_id      = data.alicloud_nlb_zones.default.zones.1.id
  cidr_block   = "192.168.20.0/24"
  vswitch_name = var.name

}

resource "alicloud_security_group" "defaultLkkjal" {
  vpc_id = alicloud_vpc.vpc.id
  name   = var.name

}

resource "alicloud_security_group" "defaultmlAdy7" {
  vpc_id = alicloud_vpc.vpc.id
  name   = var.name

}

resource "alicloud_security_group" "defaultCr6BU3" {
  vpc_id = alicloud_vpc.vpc.id
  name   = var.name

}

data "alicloud_resource_manager_resource_groups" "default" {
}

resource "alicloud_vswitch" "vsg" {
  vpc_id       = alicloud_vpc.vpc.id
  zone_id      = data.alicloud_nlb_zones.default.zones.2.id
  cidr_block   = "192.168.30.0/24"
  vswitch_name = var.name

}`, name)
}

// Case 3862

func TestAccAliCloudNlbLoadBalancer_basic3862(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_nlb_load_balancer.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudNlbLoadBalancerMap3862)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return nlb.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeNlbLoadBalancer")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%snlbloadbalancer%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudNlbLoadBalancerBasicDependence3862)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"zone_mappings": []map[string]interface{}{
						{
							"zone_id":    "${alicloud_vswitch.defaultVSwitch.zone_id}",
							"vswitch_id": "${alicloud_vswitch.defaultVSwitch.id}",
						},
						{
							"vswitch_id": "${alicloud_vswitch.defaultkR35um.id}",
							"zone_id":    "${alicloud_vswitch.defaultkR35um.zone_id}",
						},
					},
					"address_type":       "Internet",
					"vpc_id":             "${alicloud_vpc.defaultvVpc.id}",
					"load_balancer_name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"zone_mappings.#":    "2",
						"address_type":       "Internet",
						"vpc_id":             acctest.CHECKSET,
						"load_balancer_name": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"cross_zone_enabled": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cross_zone_enabled": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"zone_mappings": []map[string]interface{}{
						{
							"zone_id":    "${alicloud_vswitch.defaultVSwitch.zone_id}",
							"vswitch_id": "${alicloud_vswitch.defaultVSwitch.id}",
						},
						{
							"vswitch_id": "${alicloud_vswitch.defaultkR35um.id}",
							"zone_id":    "${alicloud_vswitch.defaultkR35um.zone_id}",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"zone_mappings.#": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"address_type": "Internet",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"address_type": "Internet",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"load_balancer_name": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"load_balancer_name": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"load_balancer_name": name + "_update",
					"zone_mappings": []map[string]interface{}{
						{
							"zone_id":    "${alicloud_vswitch.defaultVSwitch.zone_id}",
							"vswitch_id": "${alicloud_vswitch.defaultVSwitch.id}",
						},
						{
							"vswitch_id": "${alicloud_vswitch.defaultkR35um.id}",
							"zone_id":    "${alicloud_vswitch.defaultkR35um.zone_id}",
						},
					},
					"address_type":         "Internet",
					"cross_zone_enabled":   "true",
					"vpc_id":               "${alicloud_vpc.defaultvVpc.id}",
					"bandwidth_package_id": "${alicloud_common_bandwidth_package.cbwp.id}",
					"load_balancer_type":   "Network",
					"address_ip_version":   "Ipv4",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"load_balancer_name":   name + "_update",
						"zone_mappings.#":      "2",
						"address_type":         "Internet",
						"cross_zone_enabled":   "true",
						"vpc_id":               acctest.CHECKSET,
						"bandwidth_package_id": acctest.CHECKSET,
						"load_balancer_type":   "Network",
						"address_ip_version":   "Ipv4",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": acctest.REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": acctest.REMOVEKEY,
						"tags.For":     acctest.REMOVEKEY,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

var AlicloudNlbLoadBalancerMap3862 = map[string]string{
	"load_balancer_type": "Network",
	"status":             acctest.CHECKSET,
	"create_time":        acctest.CHECKSET,
}

func AlicloudNlbLoadBalancerBasicDependence3862(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

data "alicloud_nlb_zones" "default" {
}

resource "alicloud_vpc" "defaultvVpc" {
  description = "test"
  cidr_block  = "10.0.0.0/8"
  enable_ipv6 = true
  vpc_name    = var.name

}

resource "alicloud_vswitch" "defaultVSwitch" {
  vpc_id       = alicloud_vpc.defaultvVpc.id
  cidr_block   = "10.0.1.0/24"
  zone_id      = data.alicloud_nlb_zones.default.zones.0.id
  vswitch_name = var.name

}

resource "alicloud_vswitch" "defaultkR35um" {
  description  = "test"
  zone_id      = data.alicloud_nlb_zones.default.zones.1.id
  vpc_id       = alicloud_vpc.defaultvVpc.id
  cidr_block   = "10.0.2.0/24"
  vswitch_name = var.name

  ipv6_cidr_block_mask = "8"
}

resource "alicloud_common_bandwidth_package" "cbwp" {
  bandwidth            = "1000"
  internet_charge_type = "PayByBandwidth"
}`, name)
}

// Case 3678  twin

func TestAccAliCloudNlbLoadBalancer_basic3678_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_nlb_load_balancer.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudNlbLoadBalancerMap3678)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return nlb.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeNlbLoadBalancer")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%snlbloadbalancer%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudNlbLoadBalancerBasicDependence3678Twin)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, connectivity.NLBSupportRegions)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"load_balancer_name": name,
					"zone_mappings": []map[string]interface{}{
						{
							"vswitch_id": "${alicloud_vswitch.vsj.id}",
							"zone_id":    "${alicloud_vswitch.vsj.zone_id}",
						},
						{
							"vswitch_id": "${alicloud_vswitch.vsk.id}",
							"zone_id":    "${alicloud_vswitch.vsk.zone_id}",
						},
						{
							"vswitch_id": "${alicloud_vswitch.vsg.id}",
							"zone_id":    "${alicloud_vswitch.vsg.zone_id}",
						},
					},
					"address_type":       "Intranet",
					"address_ip_version": "Ipv4",
					"load_balancer_type": "Network",
					"vpc_id":             "${alicloud_vpc.vpc.id}",
					"resource_group_id":  "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"security_group_ids": []string{
						"${alicloud_security_group.defaultLkkjal.id}"},
					"deletion_protection_config": []map[string]interface{}{
						{
							"enabled": "false",
							"reason":  "",
						},
					},
					"modification_protection_config": []map[string]interface{}{
						{
							"status": "NonProtection",
							"reason": "",
						},
					},
					"cross_zone_enabled": "true",
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"load_balancer_name":   name,
						"zone_mappings.#":      "3",
						"address_type":         "Intranet",
						"address_ip_version":   "Ipv4",
						"load_balancer_type":   "Network",
						"vpc_id":               acctest.CHECKSET,
						"resource_group_id":    acctest.CHECKSET,
						"security_group_ids.#": "1",
						"cross_zone_enabled":   "true",
						"tags.%":               "2",
						"tags.Created":         "TF",
						"tags.For":             "Test",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

// Case 3862  twin

func TestAccAliCloudNlbLoadBalancer_basic3862_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_nlb_load_balancer.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudNlbLoadBalancerMap3862)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return nlb.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeNlbLoadBalancer")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%snlbloadbalancer%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudNlbLoadBalancerBasicDependence3862Twin)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"load_balancer_name": name,
					"zone_mappings": []map[string]interface{}{
						{
							"zone_id":    "${alicloud_vswitch.defaultVSwitch.zone_id}",
							"vswitch_id": "${alicloud_vswitch.defaultVSwitch.id}",
						},
						{
							"vswitch_id": "${alicloud_vswitch.defaultkR35um.id}",
							"zone_id":    "${alicloud_vswitch.defaultkR35um.zone_id}",
						},
					},
					"address_type":         "Internet",
					"cross_zone_enabled":   "true",
					"vpc_id":               "${alicloud_vpc.defaultvVpc.id}",
					"bandwidth_package_id": "${alicloud_common_bandwidth_package.cbwp.id}",
					"load_balancer_type":   "Network",
					"address_ip_version":   "Ipv4",
					"ipv6_address_type":    "Intranet",
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"load_balancer_name":   name,
						"zone_mappings.#":      "2",
						"address_type":         "Internet",
						"cross_zone_enabled":   "true",
						"vpc_id":               acctest.CHECKSET,
						"bandwidth_package_id": acctest.CHECKSET,
						"load_balancer_type":   "Network",
						"address_ip_version":   "Ipv4",
						"ipv6_address_type":    "Intranet",
						"tags.%":               "2",
						"tags.Created":         "TF",
						"tags.For":             "Test",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

// Test Nlb LoadBalancer. <<< Resource test cases, automatically generated.

func AlicloudNlbLoadBalancerBasicDependence3678Twin(name string) string {
	return strings.ReplaceAll(AlicloudNlbLoadBalancerBasicDependence3678(name), "192.168.", "172.18.")
}

func AlicloudNlbLoadBalancerBasicDependence3862Twin(name string) string {
	return strings.NewReplacer(
		"10.0.0.0/8", "10.1.0.0/16",
		"10.0.1.0/24", "10.1.1.0/24",
		"10.0.2.0/24", "10.1.2.0/24",
	).Replace(AlicloudNlbLoadBalancerBasicDependence3862(name))
}

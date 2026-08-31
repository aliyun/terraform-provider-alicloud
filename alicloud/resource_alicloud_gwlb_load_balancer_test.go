package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test Gwlb LoadBalancer. >>> Resource test cases, automatically generated.
// Case LoadBalancer Test 7994
func TestAccAliCloudGwlbLoadBalancer_basic7994(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_gwlb_load_balancer.default"
	ra := resourceAttrInit(resourceId, AlicloudGwlbLoadBalancerMap7994)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &GwlbServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeGwlbLoadBalancer")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sgwlbloadbalancer%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudGwlbLoadBalancerBasicDependence7994)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-wulanchabu"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"vpc_id":             "${alicloud_vpc.defaulti9Axhl.id}",
					"load_balancer_name": name,
					"zone_mappings": []map[string]interface{}{
						{
							"vswitch_id": "${alicloud_vswitch.default9NaKmL.id}",
							"zone_id":    "${var.zone_id1}",
						},
					},
					"address_ip_version": "Ipv4",
					"resource_group_id":  "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"dry_run":            "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"vpc_id":             CHECKSET,
						"load_balancer_name": name,
						"zone_mappings.#":    "1",
						"address_ip_version": "Ipv4",
						"resource_group_id":  CHECKSET,
						"dry_run":            "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"load_balancer_name": name + "_update",
					"zone_mappings": []map[string]interface{}{
						{
							"vswitch_id": "${alicloud_vswitch.defaultH4pKT4.id}",
							"zone_id":    "${var.zone_id2}",
						},
					},
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.1}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"load_balancer_name": name + "_update",
						"zone_mappings.#":    "1",
						"resource_group_id":  CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id": CHECKSET,
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
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"dry_run"},
			},
		},
	})
}

var AlicloudGwlbLoadBalancerMap7994 = map[string]string{
	"status":      CHECKSET,
	"create_time": CHECKSET,
}

func AlicloudGwlbLoadBalancerBasicDependence7994(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

variable "region_id" {
  default = "cn-wulanchabu"
}

variable "zone_id2" {
  default = "cn-wulanchabu-c"
}

variable "zone_id1" {
  default = "cn-wulanchabu-b"
}

data "alicloud_resource_manager_resource_groups" "default" {}

resource "alicloud_vpc" "defaulti9Axhl" {
  cidr_block = "10.0.0.0/8"
  vpc_name   = var.name
}

resource "alicloud_vswitch" "default9NaKmL" {
  vpc_id       = alicloud_vpc.defaulti9Axhl.id
  zone_id      = var.zone_id1
  cidr_block   = "10.0.0.0/24"
  vswitch_name = format("%%s1", var.name)
}

resource "alicloud_vswitch" "defaultH4pKT4" {
  vpc_id       = alicloud_vpc.defaulti9Axhl.id
  zone_id      = var.zone_id2
  cidr_block   = "10.0.1.0/24"
  vswitch_name = format("%%s2", var.name)
}


`, name)
}

// Case LoadBalancer Test_依赖资源 8565
func TestAccAliCloudGwlbLoadBalancer_basic8565(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_gwlb_load_balancer.default"
	ra := resourceAttrInit(resourceId, AlicloudGwlbLoadBalancerMap8565)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &GwlbServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeGwlbLoadBalancer")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sgwlbloadbalancer%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudGwlbLoadBalancerBasicDependence8565)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-wulanchabu"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"vpc_id":             "${alicloud_vpc.defaulti9Axhl.id}",
					"load_balancer_name": name,
					"zone_mappings": []map[string]interface{}{
						{
							"vswitch_id": "${alicloud_vswitch.default9NaKmL.id}",
							"zone_id":    "${var.zone_id1}",
						},
					},
					"address_ip_version": "Ipv4",
					"resource_group_id":  "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"dry_run":            "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"vpc_id":             CHECKSET,
						"load_balancer_name": name,
						"zone_mappings.#":    "1",
						"address_ip_version": "Ipv4",
						"resource_group_id":  CHECKSET,
						"dry_run":            "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"load_balancer_name": name + "_update",
					"zone_mappings": []map[string]interface{}{
						{
							"vswitch_id": "${alicloud_vswitch.defaultH4pKT4.id}",
							"zone_id":    "${var.zone_id2}",
						},
					},
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.1}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"load_balancer_name": name + "_update",
						"zone_mappings.#":    "1",
						"resource_group_id":  CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id": CHECKSET,
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
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"dry_run"},
			},
		},
	})
}

var AlicloudGwlbLoadBalancerMap8565 = map[string]string{
	"status":      CHECKSET,
	"create_time": CHECKSET,
}

func AlicloudGwlbLoadBalancerBasicDependence8565(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

variable "region_id" {
  default = "cn-wulanchabu"
}

variable "zone_id2" {
  default = "cn-wulanchabu-c"
}

variable "zone_id1" {
  default = "cn-wulanchabu-b"
}

data "alicloud_resource_manager_resource_groups" "default" {}

resource "alicloud_vpc" "defaulti9Axhl" {
  cidr_block = "10.0.0.0/8"
  vpc_name   = var.name
}

resource "alicloud_vswitch" "default9NaKmL" {
  vpc_id       = alicloud_vpc.defaulti9Axhl.id
  zone_id      = var.zone_id1
  cidr_block   = "10.0.0.0/24"
  vswitch_name = format("%%s1", var.name)
}

resource "alicloud_vswitch" "defaultH4pKT4" {
  vpc_id       = alicloud_vpc.defaulti9Axhl.id
  zone_id      = var.zone_id2
  cidr_block   = "10.0.1.0/24"
  vswitch_name = format("%%s2", var.name)
}


`, name)
}

// Test Gwlb LoadBalancer. <<< Resource test cases, automatically generated.

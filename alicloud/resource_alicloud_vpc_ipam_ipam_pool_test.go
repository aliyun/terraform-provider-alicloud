package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test VpcIpam IpamPool. >>> Resource test cases, automatically generated.
// Case test_parent_ipam_pool 8374
func TestAccAliCloudVpcIpamIpamPool_basic8374(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_vpc_ipam_ipam_pool.default"
	ra := resourceAttrInit(resourceId, AlicloudVpcIpamIpamPoolMap8374)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &VpcIpamServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeVpcIpamIpamPool")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%svpcipamipampool%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudVpcIpamIpamPoolBasicDependence8374)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"ipam_scope_id":       "${alicloud_vpc_ipam_ipam.defaultIpam.private_default_scope_id}",
					"pool_region_id":      "${alicloud_vpc_ipam_ipam_pool.parentIpamPool.pool_region_id}",
					"ipam_pool_name":      name,
					"source_ipam_pool_id": "${alicloud_vpc_ipam_ipam_pool.parentIpamPool.id}",
					"ip_version":          "IPv4",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ipam_scope_id":       CHECKSET,
						"pool_region_id":      CHECKSET,
						"ipam_pool_name":      name,
						"source_ipam_pool_id": CHECKSET,
						"ip_version":          "IPv4",
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
				ImportStateVerifyIgnore: []string{"clear_allocation_default_cidr_mask", "region_id"},
			},
		},
	})
}

var AlicloudVpcIpamIpamPoolMap8374 = map[string]string{
	"status":      CHECKSET,
	"create_time": CHECKSET,
	"region_id":   CHECKSET,
}

func AlicloudVpcIpamIpamPoolBasicDependence8374(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_vpc_ipam_ipam" "defaultIpam" {
  operating_region_list = ["cn-hangzhou"]
}

resource "alicloud_vpc_ipam_ipam_pool" "parentIpamPool" {
  ipam_scope_id  = alicloud_vpc_ipam_ipam.defaultIpam.private_default_scope_id
  ipam_pool_name = format("%%s1", var.name)
  pool_region_id = alicloud_vpc_ipam_ipam.defaultIpam.region_id
}


`, name)
}

// Case test_ipam_pool 8026
func TestAccAliCloudVpcIpamIpamPool_basic8026(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_vpc_ipam_ipam_pool.default"
	ra := resourceAttrInit(resourceId, AlicloudVpcIpamIpamPoolMap8026)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &VpcIpamServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeVpcIpamIpamPool")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%svpcipamipampool%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudVpcIpamIpamPoolBasicDependence8026)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"ipam_scope_id":                "${alicloud_vpc_ipam_ipam.defaultIpam.private_default_scope_id}",
					"ipam_pool_description":        "This is my ipam pool.",
					"ipam_pool_name":               name,
					"ip_version":                   "IPv4",
					"allocation_default_cidr_mask": "20",
					"allocation_min_cidr_mask":     "16",
					"allocation_max_cidr_mask":     "24",
					"pool_region_id":               "cn-hangzhou",
					"auto_import":                  "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ipam_scope_id":                CHECKSET,
						"ipam_pool_description":        "This is my ipam pool.",
						"ipam_pool_name":               name,
						"ip_version":                   "IPv4",
						"allocation_default_cidr_mask": "20",
						"allocation_min_cidr_mask":     "16",
						"allocation_max_cidr_mask":     "24",
						"pool_region_id":               "cn-hangzhou",
						"auto_import":                  "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"ipam_pool_description":              "This is my new ipam pool description.",
					"ipam_pool_name":                     name + "_update",
					"allocation_default_cidr_mask":       "24",
					"allocation_min_cidr_mask":           "12",
					"allocation_max_cidr_mask":           "26",
					"auto_import":                        "false",
					"clear_allocation_default_cidr_mask": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ipam_pool_description":              "This is my new ipam pool description.",
						"ipam_pool_name":                     name + "_update",
						"allocation_default_cidr_mask":       "24",
						"allocation_min_cidr_mask":           "12",
						"allocation_max_cidr_mask":           "26",
						"auto_import":                        "false",
						"clear_allocation_default_cidr_mask": "false",
					}),
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
				ImportStateVerifyIgnore: []string{"clear_allocation_default_cidr_mask", "region_id"},
			},
		},
	})
}

var AlicloudVpcIpamIpamPoolMap8026 = map[string]string{
	"status":      CHECKSET,
	"create_time": CHECKSET,
	"region_id":   CHECKSET,
}

func AlicloudVpcIpamIpamPoolBasicDependence8026(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_vpc_ipam_ipam" "defaultIpam" {
  operating_region_list = ["cn-hangzhou"]
  ipam_name             = var.name
}


`, name)
}

// Test VpcIpam IpamPool. <<< Resource test cases, automatically generated.

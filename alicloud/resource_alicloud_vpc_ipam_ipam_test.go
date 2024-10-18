package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test VpcIpam Ipam. >>> Resource test cases, automatically generated.
// Case test_ipam_1.1.2正式 8005
func TestAccAliCloudVpcIpamIpam_basic8005(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_vpc_ipam_ipam.default"
	ra := resourceAttrInit(resourceId, AlicloudVpcIpamIpamMap8005)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &VpcIpamServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeVpcIpamIpam")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%svpcipamipam%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudVpcIpamIpamBasicDependence8005)
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
					"ipam_description":  "This is my first Ipam.",
					"ipam_name":         name,
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"operating_region_list": []string{
						"cn-hangzhou"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ipam_description":        "This is my first Ipam.",
						"ipam_name":               name,
						"resource_group_id":       CHECKSET,
						"operating_region_list.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"ipam_description":  "This is my new ipam.",
					"ipam_name":         name + "_update",
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.1}",
					"operating_region_list": []string{
						"cn-hangzhou", "cn-beijing", "cn-qingdao"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ipam_description":        "This is my new ipam.",
						"ipam_name":               name + "_update",
						"resource_group_id":       CHECKSET,
						"operating_region_list.#": "3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"operating_region_list": []string{
						"cn-hangzhou"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id":       CHECKSET,
						"operating_region_list.#": "1",
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

var AlicloudVpcIpamIpamMap8005 = map[string]string{
	"status":      CHECKSET,
	"create_time": CHECKSET,
}

func AlicloudVpcIpamIpamBasicDependence8005(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_resource_manager_resource_groups" "default" {}

`, name)
}

// Case test_ipam_副本_副本 7857
func TestAccAliCloudVpcIpamIpam_basic7857(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_vpc_ipam_ipam.default"
	ra := resourceAttrInit(resourceId, AlicloudVpcIpamIpamMap7857)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &VpcIpamServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeVpcIpamIpam")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%svpcipamipam%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudVpcIpamIpamBasicDependence7857)
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
					"ipam_description":  "This is my first Ipam.",
					"ipam_name":         name,
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"operating_region_list": []string{
						"cn-hangzhou"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ipam_description":        "This is my first Ipam.",
						"ipam_name":               name,
						"resource_group_id":       CHECKSET,
						"operating_region_list.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"ipam_description":  "This is my new ipam.",
					"ipam_name":         name + "_update",
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.1}",
					"operating_region_list": []string{
						"cn-hangzhou", "cn-beijing", "cn-qingdao"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ipam_description":        "This is my new ipam.",
						"ipam_name":               name + "_update",
						"resource_group_id":       CHECKSET,
						"operating_region_list.#": "3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"operating_region_list": []string{
						"cn-hangzhou"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id":       CHECKSET,
						"operating_region_list.#": "1",
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

var AlicloudVpcIpamIpamMap7857 = map[string]string{
	"status":      CHECKSET,
	"create_time": CHECKSET,
}

func AlicloudVpcIpamIpamBasicDependence7857(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_resource_manager_resource_groups" "default" {}


`, name)
}

// Case test_ipam_副本 7856
func TestAccAliCloudVpcIpamIpam_basic7856(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_vpc_ipam_ipam.default"
	ra := resourceAttrInit(resourceId, AlicloudVpcIpamIpamMap7856)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &VpcIpamServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeVpcIpamIpam")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%svpcipamipam%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudVpcIpamIpamBasicDependence7856)
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
					"ipam_description":  "This is my first Ipam.",
					"ipam_name":         name,
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"operating_region_list": []string{
						"cn-hangzhou"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ipam_description":        "This is my first Ipam.",
						"ipam_name":               name,
						"resource_group_id":       CHECKSET,
						"operating_region_list.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"ipam_description":  "This is my new ipam.",
					"ipam_name":         name + "_update",
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.1}",
					"operating_region_list": []string{
						"cn-hangzhou", "cn-beijing", "cn-qingdao"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ipam_description":        "This is my new ipam.",
						"ipam_name":               name + "_update",
						"resource_group_id":       CHECKSET,
						"operating_region_list.#": "3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"operating_region_list": []string{
						"cn-hangzhou"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id":       CHECKSET,
						"operating_region_list.#": "1",
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

var AlicloudVpcIpamIpamMap7856 = map[string]string{
	"status":      CHECKSET,
	"create_time": CHECKSET,
}

func AlicloudVpcIpamIpamBasicDependence7856(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_resource_manager_resource_groups" "default" {}


`, name)
}

// Case test_ipam 7530
func TestAccAliCloudVpcIpamIpam_basic7530(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_vpc_ipam_ipam.default"
	ra := resourceAttrInit(resourceId, AlicloudVpcIpamIpamMap7530)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &VpcIpamServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeVpcIpamIpam")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%svpcipamipam%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudVpcIpamIpamBasicDependence7530)
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
					"ipam_description":  "This is my first Ipam.",
					"ipam_name":         name,
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"operating_region_list": []string{
						"cn-hangzhou"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ipam_description":        "This is my first Ipam.",
						"ipam_name":               name,
						"resource_group_id":       CHECKSET,
						"operating_region_list.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"ipam_description":  "This is my new ipam.",
					"ipam_name":         name + "_update",
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.1}",
					"operating_region_list": []string{
						"cn-hangzhou", "cn-beijing", "cn-qingdao"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ipam_description":        "This is my new ipam.",
						"ipam_name":               name + "_update",
						"resource_group_id":       CHECKSET,
						"operating_region_list.#": "3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"operating_region_list": []string{
						"cn-hangzhou"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id":       CHECKSET,
						"operating_region_list.#": "1",
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

var AlicloudVpcIpamIpamMap7530 = map[string]string{
	"status":      CHECKSET,
	"create_time": CHECKSET,
}

func AlicloudVpcIpamIpamBasicDependence7530(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_resource_manager_resource_groups" "default" {}


`, name)
}

// Test VpcIpam Ipam. <<< Resource test cases, automatically generated.

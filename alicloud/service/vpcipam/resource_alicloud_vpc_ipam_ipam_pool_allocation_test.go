package vpcipam_test

import (
	"fmt"
	"testing"

	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/vpcipam"
)

// Test VpcIpam IpamPoolAllocation. >>> Resource test cases, automatically generated.
// Case test_ipam_pool_allocation 7885
func TestAccAliCloudVpcIpamIpamPoolAllocation_basic7885(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_vpc_ipam_ipam_pool_allocation.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudVpcIpamIpamPoolAllocationMap7885)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return vpcipam.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeVpcIpamIpamPoolAllocation")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%svpcipamipampoolallocation%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudVpcIpamIpamPoolAllocationBasicDependence7885)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"ipam_pool_allocation_description": "init alloc desc",
					"ipam_pool_allocation_name":        name,
					"cidr":                             "10.0.0.0/20",
					"ipam_pool_id":                     "${alicloud_vpc_ipam_ipam_pool_cidr.defaultIpamPoolCidr.ipam_pool_id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ipam_pool_allocation_description": "init alloc desc",
						"ipam_pool_allocation_name":        name,
						"cidr":                             "10.0.0.0/20",
						"ipam_pool_id":                     acctest.CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"ipam_pool_allocation_description": "update alloc desc",
					"ipam_pool_allocation_name":        name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ipam_pool_allocation_description": "update alloc desc",
						"ipam_pool_allocation_name":        name + "_update",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"cidr_mask"},
			},
		},
	})
}

var AlicloudVpcIpamIpamPoolAllocationMap7885 = map[string]string{
	"status":      acctest.CHECKSET,
	"create_time": acctest.CHECKSET,
	"region_id":   acctest.CHECKSET,
}

func AlicloudVpcIpamIpamPoolAllocationBasicDependence7885(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_vpc_ipam_ipam" "defaultIpam" {
  operating_region_list = ["cn-hangzhou"]
}

resource "alicloud_vpc_ipam_ipam_pool" "defaultIpamPool" {
  ipam_scope_id  = alicloud_vpc_ipam_ipam.defaultIpam.private_default_scope_id
  pool_region_id = "cn-hangzhou"
}

resource "alicloud_vpc_ipam_ipam_pool_cidr" "defaultIpamPoolCidr" {
  cidr         = "10.0.0.0/8"
  ipam_pool_id = alicloud_vpc_ipam_ipam_pool.defaultIpamPool.id
}


`, name)
}

// Case test_ipam_pool_allocation_creating_by_cidrmask 9268
func TestAccAliCloudVpcIpamIpamPoolAllocation_basic9268(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_vpc_ipam_ipam_pool_allocation.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudVpcIpamIpamPoolAllocationMap9268)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return vpcipam.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeVpcIpamIpamPoolAllocation")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%svpcipamipampoolallocation%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudVpcIpamIpamPoolAllocationBasicDependence9268)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"ipam_pool_allocation_description": "init alloc desc",
					"ipam_pool_allocation_name":        name,
					"ipam_pool_id":                     "${alicloud_vpc_ipam_ipam_pool_cidr.defaultIpamPoolCidr.ipam_pool_id}",
					"cidr_mask":                        "20",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ipam_pool_allocation_description": "init alloc desc",
						"ipam_pool_allocation_name":        name,
						"ipam_pool_id":                     acctest.CHECKSET,
						"cidr_mask":                        "20",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"ipam_pool_allocation_description": "update alloc desc",
					"ipam_pool_allocation_name":        name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ipam_pool_allocation_description": "update alloc desc",
						"ipam_pool_allocation_name":        name + "_update",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"cidr_mask"},
			},
		},
	})
}

var AlicloudVpcIpamIpamPoolAllocationMap9268 = map[string]string{
	"status":      acctest.CHECKSET,
	"create_time": acctest.CHECKSET,
	"region_id":   acctest.CHECKSET,
}

func AlicloudVpcIpamIpamPoolAllocationBasicDependence9268(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_vpc_ipam_ipam" "defaultIpam" {
  operating_region_list = ["cn-hangzhou"]
}

resource "alicloud_vpc_ipam_ipam_pool" "defaultIpamPool" {
  ipam_scope_id  = alicloud_vpc_ipam_ipam.defaultIpam.private_default_scope_id
  pool_region_id = "cn-hangzhou"
}

resource "alicloud_vpc_ipam_ipam_pool_cidr" "defaultIpamPoolCidr" {
  cidr         = "10.0.0.0/8"
  ipam_pool_id = alicloud_vpc_ipam_ipam_pool.defaultIpamPool.id
}


`, name)
}

// Test VpcIpam IpamPoolAllocation. <<< Resource test cases, automatically generated.

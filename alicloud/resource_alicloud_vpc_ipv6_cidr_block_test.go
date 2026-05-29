// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// Test Vpc Ipv6CidrBlock. >>> Resource test cases, automatically generated.
// Case 从IPAM地址池为VPC添加IPv6网段_ula_create_by_cidrblock 12802
func TestAccAliCloudVPCIpv6CidrBlock_basic12802(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_vpc_ipv6_cidr_block.default"
	ra := resourceAttrInit(resourceId, AlicloudVpcIpv6CidrBlockMap12802)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &VpcServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeVpcIpv6CidrBlock")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccvpc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudVpcIpv6CidrBlockBasicDependence12802)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"ipv6_ipam_pool_id": "${alicloud_vpc_ipam_ipam_pool_cidr.defaultIpv6PoolCidr.ipam_pool_id}",
					"vpc_id":            "${alicloud_vpc.defaultVpc.id}",
					"ipv6_cidr_block":   "fd03:d00:a000::/60",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ipv6_ipam_pool_id": CHECKSET,
						"vpc_id":            CHECKSET,
						"ipv6_cidr_block":   "fd03:d00:a000::/60",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"ipv6_cidr_mask", "ipv6_ipam_pool_id"},
			},
		},
	})
}

var AlicloudVpcIpv6CidrBlockMap12802 = map[string]string{}

func AlicloudVpcIpv6CidrBlockBasicDependence12802(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_vpc_ipam_ipam" "defaultIpam" {
  operating_region_list = ["cn-hangzhou", "cn-beijing"]
}

resource "alicloud_vpc_ipam_ipam_pool" "defaultIpv6Pool" {
  ipam_scope_id  = alicloud_vpc_ipam_ipam.defaultIpam.private_default_scope_id
  pool_region_id = alicloud_vpc_ipam_ipam.defaultIpam.region_id
  ip_version     = "IPv6"
}

resource "alicloud_vpc_ipam_ipam_pool_cidr" "defaultIpv6PoolCidr" {
  ipam_pool_id = alicloud_vpc_ipam_ipam_pool.defaultIpv6Pool.id
  cidr         = "fd03:d00:a000::/48"
}

resource "alicloud_vpc" "defaultVpc" {
  cidr_block = "10.0.0.0/8"
  vpc_name   = "test-ipv6-cidr-block"

  lifecycle {
    ignore_changes = [enable_ipv6, ipv6_cidr_block, ipv6_cidr_blocks]
  }
}


`, name)
}

// Case 从IPAM地址池为VPC添加IPv6网段_byoip 12801
func TestAccAliCloudVPCIpv6CidrBlock_basic12801(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_vpc_ipv6_cidr_block.default"
	ra := resourceAttrInit(resourceId, AlicloudVpcIpv6CidrBlockMap12801)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &VpcServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeVpcIpv6CidrBlock")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccvpc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudVpcIpv6CidrBlockBasicDependence12801)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"ipv6_ipam_pool_id": "${alicloud_vpc_ipam_ipam_pool_cidr.defaultIpv6PoolCidr.ipam_pool_id}",
					"ipv6_cidr_mask":    "60",
					"vpc_id":            "${alicloud_vpc.defaultVpc.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ipv6_ipam_pool_id": CHECKSET,
						"ipv6_cidr_mask":    "60",
						"vpc_id":            CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"ipv6_cidr_mask", "ipv6_ipam_pool_id"},
			},
		},
	})
}

var AlicloudVpcIpv6CidrBlockMap12801 = map[string]string{}

func AlicloudVpcIpv6CidrBlockBasicDependence12801(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_vpc_ipam_ipam" "defaultIpam" {
  operating_region_list = ["cn-hangzhou", "cn-beijing"]
}

resource "alicloud_vpc_ipam_ipam_pool" "defaultIpv6Pool" {
  ipam_scope_id  = alicloud_vpc_ipam_ipam.defaultIpam.private_default_scope_id
  pool_region_id = alicloud_vpc_ipam_ipam.defaultIpam.region_id
  ip_version     = "IPv6"
}

resource "alicloud_vpc_ipam_ipam_pool_cidr" "defaultIpv6PoolCidr" {
  ipam_pool_id = alicloud_vpc_ipam_ipam_pool.defaultIpv6Pool.id
  cidr         = "2a03:1b40:7c2e:9f00::/58"
}

resource "alicloud_vpc" "defaultVpc" {
  cidr_block = "10.0.0.0/8"
  vpc_name   = "test-ipv6-cidr-block"

  lifecycle {
    ignore_changes = [enable_ipv6, ipv6_cidr_block, ipv6_cidr_blocks]
  }
}


`, name)
}

// Test Vpc Ipv6CidrBlock. <<< Resource test cases, automatically generated.

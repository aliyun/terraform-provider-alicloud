// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test Lindorm PublicNetwork. >>> Resource test cases, automatically generated.
// Case PublicNetwork用例_线上 10782
func TestAccAliCloudLindormPublicNetwork_basic10782(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_lindorm_public_network.default"
	ra := resourceAttrInit(resourceId, AlicloudLindormPublicNetworkMap10782)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &LindormServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeLindormPublicNetwork")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacclindorm%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudLindormPublicNetworkBasicDependence10782)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-shanghai"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_id":           "${alicloud_lindorm_instance.defaultQpsLKr.id}",
					"enable_public_network": "1",
					"engine_type":           "lindorm",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_id":           CHECKSET,
						"enable_public_network": "1",
						"engine_type":           "lindorm",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"enable_public_network", "engine_type"},
			},
		},
	})
}

var AlicloudLindormPublicNetworkMap10782 = map[string]string{
	"status": CHECKSET,
}

func AlicloudLindormPublicNetworkBasicDependence10782(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

variable "zone_id" {
  default = "cn-shanghai-f"
}

variable "region_id" {
  default = "cn-shanghai"
}

resource "alicloud_vpc" "defaultX7MgJO" {
  description = var.name
  cidr_block  = "10.0.0.0/8"
  vpc_name    = "amp-test-shanghai"
}

resource "alicloud_vswitch" "default45mCzM" {
  description = var.name
  vpc_id      = alicloud_vpc.defaultX7MgJO.id
  zone_id     = var.zone_id
  cidr_block  = "10.0.0.0/24"
}

resource "alicloud_lindorm_instance" "defaultQpsLKr" {
  payment_type               = "PayAsYouGo"
  table_engine_node_count    = "2"
  instance_storage           = "80"
  zone_id                    = var.zone_id
  vswitch_id                 = alicloud_vswitch.default45mCzM.id
  disk_category              = "cloud_efficiency"
  table_engine_specification = "lindorm.g.xlarge"
  instance_name              = "tf-test"
  vpc_id                     = alicloud_vpc.defaultX7MgJO.id
}


`, name)
}

// Case PublicNetwork用例1 10758
func TestAccAliCloudLindormPublicNetwork_basic10758(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_lindorm_public_network.default"
	ra := resourceAttrInit(resourceId, AlicloudLindormPublicNetworkMap10758)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &LindormServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeLindormPublicNetwork")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacclindorm%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudLindormPublicNetworkBasicDependence10758)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-shanghai"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_id":           "${alicloud_lindorm_instance.defaultQpsLKr.id}",
					"enable_public_network": "1",
					"engine_type":           "lindorm",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_id":           CHECKSET,
						"enable_public_network": "1",
						"engine_type":           "lindorm",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"enable_public_network", "engine_type"},
			},
		},
	})
}

var AlicloudLindormPublicNetworkMap10758 = map[string]string{
	"status": CHECKSET,
}

func AlicloudLindormPublicNetworkBasicDependence10758(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

variable "zone_id" {
  default = "cn-shanghai-f"
}

variable "region_id" {
  default = "cn-shanghai"
}

resource "alicloud_vpc" "defaultX7MgJO" {
  description = var.name
  cidr_block  = "10.0.0.0/8"
  vpc_name    = "amp-test-shanghai"
}

resource "alicloud_vswitch" "default45mCzM" {
  description = var.name
  vpc_id      = alicloud_vpc.defaultX7MgJO.id
  zone_id     = var.zone_id
  cidr_block  = "10.0.0.0/24"
}

resource "alicloud_lindorm_instance" "defaultQpsLKr" {
  payment_type               = "PayAsYouGo"
  table_engine_node_count    = "2"
  instance_storage           = "80"
  zone_id                    = var.zone_id
  vswitch_id                 = alicloud_vswitch.default45mCzM.id
  disk_category              = "cloud_efficiency"
  table_engine_specification = "lindorm.g.xlarge"
  instance_name              = "tf-test"
  vpc_id                     = alicloud_vpc.defaultX7MgJO.id
}


`, name)
}

// Test Lindorm PublicNetwork. <<< Resource test cases, automatically generated.

package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Case 1
func TestAccAliCloudDiskDiskReplicaPair_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ebs_disk_replica_pair.default"
	ra := resourceAttrInit(resourceId, AlicloudDiskDiskReplicaPairMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EbsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEbsDiskReplicaPair")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sDiskDiskReplicaPair%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudDiskDiskReplicaPairBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"destination_disk_id":   "${alicloud_ecs_disk.default.id}",
					"destination_region_id": "cn-hangzhou",
					"destination_zone_id":   "cn-hangzhou-h",
					"source_zone_id":        "cn-hangzhou-g",
					"disk_id":               "${alicloud_ecs_disk.defaultone.id}",
					"description":           name,
					"pair_name":             name,
					"payment_type":          "POSTPAY",
					"resource_group_id":     "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"one_shot":              "true",
					"reverse_replicate":     "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"destination_disk_id":   CHECKSET,
						"destination_region_id": "cn-hangzhou",
						"destination_zone_id":   "cn-hangzhou-h",
						"source_zone_id":        "cn-hangzhou-g",
						"disk_id":               CHECKSET,
						"description":           name,
						"pair_name":             name,
						"payment_type":          CHECKSET,
						"resource_group_id":     CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.1}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"pair_name": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"pair_name": name + "_update",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"one_shot", "reverse_replicate"},
			},
		},
	})
}

var AlicloudDiskDiskReplicaPairMap = map[string]string{}

func AlicloudDiskDiskReplicaPairBasicDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_resource_manager_resource_groups" "default" {}

resource "alicloud_ecs_disk" "default" {
	zone_id = "cn-hangzhou-h"
	category = "cloud_essd"
	delete_auto_snapshot = "true"
	delete_with_instance = "true"
	description = "Test For Terraform"
	disk_name = var.name
	enable_auto_snapshot = "true"
	encrypted = "true"
	size = "500"
	lifecycle {
		ignore_changes = [tags]
	}
}

resource "alicloud_ecs_disk" "defaultone" {
	zone_id = "cn-hangzhou-g"
	category = "cloud_essd"
	delete_auto_snapshot = "true"
	delete_with_instance = "true"
	description = "Test For Terraform"
	disk_name = var.name
	enable_auto_snapshot = "true"
	encrypted = "true"
	size = "500"
	lifecycle {
		ignore_changes = [tags]
	}
}

`, name)
}

// Test Ebs DiskReplicaPair. >>> Resource test cases, automatically generated.
// Case 帆过测试2024年03月07日-有预付费bandwidth-102400 6274
func TestAccAliCloudEbsDiskReplicaPair_basic6274(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ebs_disk_replica_pair.default"
	ra := resourceAttrInit(resourceId, AlicloudEbsDiskReplicaPairMap6274)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EbsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEbsDiskReplicaPair")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccebs%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudEbsDiskReplicaPairBasicDependence6274)
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
					"destination_disk_id":    "${alicloud_ecs_disk.defaultHU3EL9.id}",
					"destination_region_id":  "${var.disk-region}",
					"destination_zone_id":    "${alicloud_ecs_disk.defaultHU3EL9.zone_id}",
					"source_zone_id":         "${alicloud_ecs_disk.defaultoPIuxo.zone_id}",
					"disk_id":                "${alicloud_ecs_disk.defaultoPIuxo.id}",
					"description":            "ccapi-test",
					"payment_type":           "Subscription",
					"disk_replica_pair_name": name,
					"rpo":                    "900",
					"period":                 "1",
					"bandwidth":              "102400",
					"period_unit":            "Month",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"destination_disk_id":    CHECKSET,
						"destination_region_id":  CHECKSET,
						"destination_zone_id":    CHECKSET,
						"source_zone_id":         CHECKSET,
						"disk_id":                CHECKSET,
						"description":            "ccapi-test",
						"payment_type":           "Subscription",
						"disk_replica_pair_name": name,
						"rpo":                    "900",
						"period":                 "1",
						"bandwidth":              "102400",
						"period_unit":            "Month",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"disk_replica_pair_name": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"disk_replica_pair_name": name + "_update",
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
				ImportStateVerifyIgnore: []string{"one_shot", "period", "period_unit", "reverse_replicate"},
			},
		},
	})
}

var AlicloudEbsDiskReplicaPairMap6274 = map[string]string{
	"status":      CHECKSET,
	"create_time": CHECKSET,
	"region_id":   CHECKSET,
}

func AlicloudEbsDiskReplicaPairBasicDependence6274(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

variable "disk-region" {
  default = "cn-hangzhou"
}

resource "alicloud_ecs_disk" "defaultoPIuxo" {
  category  = "cloud_essd"
  zone_id   = "cn-hangzhou-g"
  encrypted = false
  size      = "20"
  disk_name = "fg-terr"
  lifecycle {
	ignore_changes = [tags]
  }
}

resource "alicloud_ecs_disk" "defaultHU3EL9" {
  category  = "cloud_essd"
  zone_id   = "cn-hangzhou-h"
  encrypted = false
  size      = "20"
  disk_name = "fg-terr"
  lifecycle {
	ignore_changes = [tags]
  }
}


`, name)
}

// Test Ebs DiskReplicaPair. <<< Resource test cases, automatically generated.

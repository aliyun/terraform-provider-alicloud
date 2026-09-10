package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAliCloudECSInstanceSet_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ecs_instance_set.default"
	ra := resourceAttrInit(resourceId, AliCloudECSInstanceSetMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EcsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEcsInstanceSet")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%secs%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudECSInstanceSetBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.TestSalveRegions)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"amount":                     "10",
					"image_id":                   "${data.alicloud_images.default.images[0].id}",
					"instance_type":              "${data.alicloud_instance_types.default.instance_types[0].id}",
					"instance_name":              "${var.name}",
					"security_group_ids":         []string{"${alicloud_security_group.default.id}"},
					"internet_charge_type":       "PayByTraffic",
					"internet_max_bandwidth_out": "10",
					"zone_id":                    "${alicloud_vswitch.default.zone_id}",
					"instance_charge_type":       "PostPaid",
					"system_disk_category":       "cloud_efficiency",
					"vswitch_id":                 "${alicloud_vswitch.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"image_id":                   CHECKSET,
						"instance_type":              CHECKSET,
						"instance_name":              CHECKSET,
						"security_group_ids.#":       "1",
						"internet_charge_type":       "PayByTraffic",
						"internet_max_bandwidth_out": "10",
						"zone_id":                    CHECKSET,
						"instance_charge_type":       "PostPaid",
						"system_disk_category":       "cloud_efficiency",
						"vswitch_id":                 CHECKSET,
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
		},
	})
}

func TestAccAliCloudECSInstanceSet_basic1(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ecs_instance_set.default"
	ra := resourceAttrInit(resourceId, AliCloudECSInstanceSetMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EcsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEcsInstanceSet")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%secs%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudECSInstanceSetBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.TestSalveRegions)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"amount":                     "10",
					"image_id":                   "${data.alicloud_images.default.images[0].id}",
					"instance_type":              "${data.alicloud_instance_types.default.instance_types[0].id}",
					"instance_name":              "${var.name}",
					"security_group_ids":         []string{"${alicloud_security_group.default.id}"},
					"internet_charge_type":       "PayByTraffic",
					"internet_max_bandwidth_out": "10",
					"zone_id":                    "${alicloud_vswitch.default.zone_id}",
					"instance_charge_type":       "PostPaid",
					"system_disk_category":       "cloud_efficiency",
					"vswitch_id":                 "${alicloud_vswitch.default.id}",
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"image_id":                   CHECKSET,
						"instance_type":              CHECKSET,
						"instance_name":              CHECKSET,
						"security_group_ids.#":       "1",
						"internet_charge_type":       "PayByTraffic",
						"internet_max_bandwidth_out": "10",
						"zone_id":                    CHECKSET,
						"instance_charge_type":       "PostPaid",
						"system_disk_category":       "cloud_efficiency",
						"vswitch_id":                 CHECKSET,
						"tags.%":                     "2",
						"tags.Created":               "TF",
						"tags.For":                   "Test",
					}),
				),
			},
		},
	})
}

var AliCloudECSInstanceSetMap0 = map[string]string{}

func AliCloudECSInstanceSetBasicDependence0(name string) string {
	return fmt.Sprintf(` 
	variable "name" {
  		default = "%s"
	}

	data "alicloud_images" "default" {
  		name_regex  = "^ubuntu_[0-9]+_[0-9]+_x64*"
  		most_recent = true
  		owners      = "system"
	}

	data "alicloud_instance_types" "default" {
  		image_id = data.alicloud_images.default.images.0.id
	}

	resource "alicloud_vpc" "default" {
  		cidr_block = "192.168.0.0/16"
  		vpc_name   = var.name
	}

	resource "alicloud_vswitch" "default" {
  		vpc_id       = alicloud_vpc.default.id
  		cidr_block   = cidrsubnet(alicloud_vpc.default.cidr_block, 8, 8)
  		zone_id      = data.alicloud_instance_types.default.instance_types.0.availability_zones.0
  		vswitch_name = var.name
	}

	resource "alicloud_security_group" "default" {
  		name   = var.name
  		vpc_id = alicloud_vpc.default.id
	}
`, name)
}

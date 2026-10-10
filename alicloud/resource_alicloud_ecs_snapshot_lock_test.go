package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAliCloudECSSnapshotLock_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ecs_snapshot_lock.default"
	ra := resourceAttrInit(resourceId, AliCloudEcsSnapshotLockMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EcsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEcsSnapshotLock")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%secssnapshotlock%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudEcsSnapshotLockBasicDependence0)
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
					"snapshot_lock_id": "${alicloud_ecs_snapshot.default.id}",
					"cool_off_period":  "72",
					"lock_duration":    "73",
					"lock_mode":        "compliance",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"snapshot_lock_id": CHECKSET,
						"cool_off_period":  "72",
						"lock_duration":    "73",
						"lock_mode":        "compliance",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"snapshot_lock_id": "${alicloud_ecs_snapshot.default.id}",
					"cool_off_period":  "48",
					"lock_duration":    "74",
					"lock_mode":        "compliance",
					"dry_run":          "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cool_off_period": "48",
						"lock_duration":   "74",
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

var AliCloudEcsSnapshotLockMap0 = map[string]string{
	"lock_status": CHECKSET,
}

func AliCloudEcsSnapshotLockBasicDependence0(name string) string {
	return fmt.Sprintf(`
	variable "name" {
		default = "%s"
	}

	data "alicloud_resource_manager_resource_groups" "default" {
		status = "OK"
	}

	data "alicloud_zones" "default" {
		available_disk_category     = "cloud_essd"
		available_instance_type     = data.alicloud_instance_types.default.instance_types.0.id
		available_resource_creation = "Instance"
	}

	data "alicloud_images" "default" {
		most_recent = true
		owners      = "system"
		architecture = "x86_64"
		os_type      = "linux"
	}

	data "alicloud_instance_types" "default" {
		instance_type_family = "ecs.g7"
		sorted_by            = "CPU"
		image_id             = data.alicloud_images.default.images.0.id
		system_disk_category = "cloud_essd"
	}

	resource "alicloud_vpc" "default" {
		vpc_name   = var.name
		cidr_block = "192.168.0.0/16"
	}

	resource "alicloud_vswitch" "default" {
		vswitch_name = var.name
		vpc_id       = alicloud_vpc.default.id
		cidr_block   = "192.168.192.0/24"
		zone_id      = data.alicloud_zones.default.zones.0.id
	}

	resource "alicloud_security_group" "default" {
		security_group_name = var.name
		vpc_id              = alicloud_vpc.default.id
	}

	resource "alicloud_instance" "default" {
		image_id                   = data.alicloud_images.default.images.0.id
		instance_type              = data.alicloud_instance_types.default.instance_types.0.id
		security_groups            = alicloud_security_group.default.*.id
		internet_charge_type       = "PayByTraffic"
		internet_max_bandwidth_out = "10"
		availability_zone          = data.alicloud_zones.default.zones.0.id
		instance_charge_type       = "PostPaid"
		system_disk_category       = "cloud_essd"
		system_disk_encrypted      = true
		vswitch_id                 = alicloud_vswitch.default.id
		instance_name              = var.name
		data_disks {
			category = "cloud_essd"
			encrypted = true
			size     = 20
		}
	}

	resource "alicloud_ecs_disk" "default" {
		disk_name = var.name
		zone_id   = data.alicloud_zones.default.zones.0.id
		category  = "cloud_essd"
		encrypted = true
		size      = 20
	}

	resource "alicloud_ecs_disk_attachment" "default" {
		disk_id     = alicloud_ecs_disk.default.id
		instance_id = alicloud_instance.default.id
	}

	resource "alicloud_ecs_snapshot" "default" {
		disk_id       = alicloud_ecs_disk_attachment.default.disk_id
		snapshot_name = var.name
		wait_until    = "available"
	}
`, name)
}

package ecs_test

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/ecs"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccAliCloudECSSnapshotLock_basic0(t *testing.T) {
	resourceId := "alicloud_ecs_snapshot_lock.default"
	rand := tfacctest.RandIntRange(10000, 99999)
	region := os.Getenv("ALICLOUD_REGION")
	if region == "" {
		region = "cn-beijing"
	}
	name := fmt.Sprintf("tf-testacc%secssnapshotlock%d", region, rand)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckEcsSnapshotLockDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccAliCloudECSSnapshotLockBasic0DryRun(name),
				ExpectError: regexp.MustCompile("DryRunOperation"),
			},
			{
				Config: testAccAliCloudECSSnapshotLockBasic0(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceId, "snapshot_lock_id"),
					resource.TestCheckResourceAttr(resourceId, "cool_off_period", "72"),
					resource.TestCheckResourceAttr(resourceId, "lock_duration", "73"),
					resource.TestCheckResourceAttr(resourceId, "lock_mode", "compliance"),
					// DescribeLockedSnapshots omits the three phase-dependent time
					// fields while the lock is still in its cool-off period.
					resource.TestCheckResourceAttrSet(resourceId, "lock_creation_time"),
					resource.TestCheckResourceAttrSet(resourceId, "lock_status"),
				),
			},
			{
				Config: testAccAliCloudECSSnapshotLockBasic0Update(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "cool_off_period", "48"),
					resource.TestCheckResourceAttr(resourceId, "lock_duration", "74"),
					resource.TestCheckResourceAttr(resourceId, "dry_run", "false"),
				),
			},
			{
				Config: testAccAliCloudECSSnapshotLockBasic0RemoveDryRun(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "cool_off_period", "48"),
					resource.TestCheckResourceAttr(resourceId, "lock_duration", "74"),
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

func testAccCheckEcsSnapshotLockDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_ecs_snapshot_lock" {
			continue
		}
		client := acctest.Provider.Meta().(*connectivity.AliyunClient)
		_, err := ecs.NewService(client).DescribeSnapshotLock(rs.Primary.ID)
		if err != nil {
			if errs.NotFoundError(err) {
				continue
			}
			return fmt.Errorf("describing snapshot lock %s failed: %s", rs.Primary.ID, err)
		}
		return fmt.Errorf("snapshot lock %s still exists", rs.Primary.ID)
	}
	return nil
}

func testAccAliCloudECSSnapshotLockBasic0(name string) string {
	return fmt.Sprintf(`%s

	resource "alicloud_ecs_snapshot_lock" "default" {
		snapshot_lock_id = alicloud_ecs_snapshot.default.id
		cool_off_period  = 72
		lock_duration    = 73
		lock_mode        = "compliance"
	}
`, testAccAliCloudECSSnapshotLockBasicDependence0(name))
}

func testAccAliCloudECSSnapshotLockBasic0DryRun(name string) string {
	return fmt.Sprintf(`%s

	resource "alicloud_ecs_snapshot_lock" "default" {
		snapshot_lock_id = alicloud_ecs_snapshot.default.id
		cool_off_period  = 72
		lock_duration    = 73
		lock_mode        = "compliance"
		dry_run          = true
	}
`, testAccAliCloudECSSnapshotLockBasicDependence0(name))
}

func testAccAliCloudECSSnapshotLockBasic0Update(name string) string {
	return fmt.Sprintf(`%s

	resource "alicloud_ecs_snapshot_lock" "default" {
		snapshot_lock_id = alicloud_ecs_snapshot.default.id
		cool_off_period  = 48
		lock_duration    = 74
		lock_mode        = "compliance"
		dry_run          = false
	}
`, testAccAliCloudECSSnapshotLockBasicDependence0(name))
}

func testAccAliCloudECSSnapshotLockBasic0RemoveDryRun(name string) string {
	return fmt.Sprintf(`%s

	resource "alicloud_ecs_snapshot_lock" "default" {
		snapshot_lock_id = alicloud_ecs_snapshot.default.id
		cool_off_period  = 48
		lock_duration    = 74
		lock_mode        = "compliance"
	}
`, testAccAliCloudECSSnapshotLockBasicDependence0(name))
}

func testAccAliCloudECSSnapshotLockBasicDependence0(name string) string {
	return fmt.Sprintf(`
	variable "name" {
		default = "%s"
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

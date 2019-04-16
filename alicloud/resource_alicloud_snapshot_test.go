package alicloud

import (
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func init() {
	resource.AddTestSweepers("alicloud_snapshot", &resource.Sweeper{
		Name: "alicloud_snapshot",
		F:    testSweepSnapshots,
	})
}

func testSweepSnapshots(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}

	var snapshots []ecs.Snapshot
	req := ecs.CreateDescribeSnapshotsRequest()
	req.RegionId = client.RegionId
	req.PageSize = requests.NewInteger(PageSizeLarge)
	req.PageNumber = requests.NewInteger(1)
	for {
		raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.DescribeSnapshots(req)
		})
		if err != nil {
			return fmt.Errorf("Error retrieving snapshots: %s", err)
		}
		resp, _ := raw.(*ecs.DescribeSnapshotsResponse)
		if resp == nil || len(resp.Snapshots.Snapshot) < 1 {
			break
		}
		snapshots = append(snapshots, resp.Snapshots.Snapshot...)

		if len(resp.Snapshots.Snapshot) < PageSizeLarge {
			break
		}

		if page, err := getNextpageNumber(req.PageNumber); err != nil {
			return err
		} else {
			req.PageNumber = page
		}
	}

	sweeped := false
	for _, v := range snapshots {
		name := v.SnapshotName
		id := v.SnapshotId
		skip := true
		for _, prefix := range prefixes {
			if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
				skip = false
				break
			}
		}
		if skip {
			log.Printf("[INFO] Skipping snapshot: %s (%s)", name, id)
			continue
		}
		sweeped = true
		log.Printf("[INFO] Deleting snapshot: %s (%s)", name, id)
		req := ecs.CreateDeleteSnapshotRequest()
		req.SnapshotId = id
		_, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.DeleteSnapshot(req)
		})
		if err != nil {
			log.Printf("[ERROR] Failed to delete snapshot(%s (%s)): %s", name, id, err)
		}
	}

	if sweeped {
		time.Sleep(30 * time.Second)
	}
	return nil
}

func TestAccAlicloudSnapshot_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "alicloud_snapshot.snapshot",

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSnapshotDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccSnapshotConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSnapshotExists("alicloud_snapshot.snapshot"),
					resource.TestCheckResourceAttrSet("alicloud_snapshot.snapshot", "disk_id"),
					resource.TestCheckResourceAttr("alicloud_snapshot.snapshot", "name", "tf-testAcc-snapshot"),
					resource.TestCheckResourceAttr("alicloud_snapshot.snapshot", "description", "TF Test"),
					resource.TestCheckResourceAttr("alicloud_snapshot.snapshot", "tags.%", "1"),
					resource.TestCheckResourceAttr("alicloud_snapshot.snapshot", "tags.version", "1.0"),
				),
			},
			resource.TestStep{
				Config: testAccSnapshotConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSnapshotExists("alicloud_snapshot.snapshot"),
					resource.TestCheckResourceAttrSet("alicloud_snapshot.snapshot", "disk_id"),
					resource.TestCheckResourceAttr("alicloud_snapshot.snapshot", "name", "tf-testAcc-snapshot"),
					resource.TestCheckResourceAttr("alicloud_snapshot.snapshot", "description", "TF Test"),
					resource.TestCheckResourceAttr("alicloud_snapshot.snapshot", "tags.%", "1"),
					resource.TestCheckResourceAttr("alicloud_snapshot.snapshot", "tags.version", "1.1"),
				),
			},
		},
	})
}

func testAccCheckSnapshotDestroy(s *terraform.State) error {

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_snapshot" {
			continue
		}

		client := testAccProvider.Meta().(*connectivity.AliyunClient)
		ecsService := EcsService{client}

		_, err := ecsService.DescribeSnapshotById(rs.Primary.ID)

		if err != nil {
			if NotFoundError(err) {
				continue
			}
			return fmt.Errorf("Describing snapshot (%s) failed while destoring, error: %#v.", rs.Primary.ID, err)
		}
		return fmt.Errorf("Error ECS Snapshot (%s) still exist", rs.Primary.ID)
	}

	return nil
}

func testAccCheckSnapshotExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Snapshot ID is set")
		}

		client := testAccProvider.Meta().(*connectivity.AliyunClient)
		ecsService := EcsService{client}
		_, err := ecsService.DescribeSnapshotById(rs.Primary.ID)

		if err != nil {
			return fmt.Errorf("While checking disk existing, describing disk got an error: %#v.", err)
		}

		return nil
	}
}

const testAccSnapshotConfig = `
data "alicloud_instance_types" "instance_type" {
  availability_zone = "${data.alicloud_zones.zone.zones.0.id}"
}

resource "alicloud_vpc" "vpc" {
  name = "tf-testAcc-vpc"
  cidr_block = "192.168.0.0/16"
}

data "alicloud_zones" "zone" {
}

resource "alicloud_vswitch" "vswitch" {
  name = "tf-testAcc-vswitch"
  cidr_block = "192.168.0.0/24"
  availability_zone = "${data.alicloud_zones.zone.zones.0.id}"
  vpc_id = "${alicloud_vpc.vpc.id}"
}

resource "alicloud_security_group" "group" {
  name        = "tf-testACC-group"
  description = "New security group"
  vpc_id = "${alicloud_vpc.vpc.id}"
}

resource "alicloud_disk" "disk" {
  name = "disk"
  availability_zone = "${alicloud_instance.instance.availability_zone}"
  category          = "cloud_efficiency"
  size              = "20"
}

data "alicloud_images" "sys_images" {
  owners = "system"
}

resource "alicloud_instance" "instance" {
  instance_name   = "tf-testAcc-instance"
  host_name       = "tf-testAcc"
  image_id        = "${data.alicloud_images.sys_images.images.0.id}"
  instance_type   = "${data.alicloud_instance_types.instance_type.instance_types.0.id}"
  security_groups = ["${alicloud_security_group.group.id}"]
  vswitch_id      = "${alicloud_vswitch.vswitch.id}"
}

resource "alicloud_disk_attachment" "instance-attachment" {
  disk_id     = "${alicloud_disk.disk.id}"
  instance_id = "${alicloud_instance.instance.id}"
}

resource "alicloud_snapshot" "snapshot" {
  disk_id = "${alicloud_disk_attachment.instance-attachment.disk_id}"
  name = "tf-testAcc-snapshot"
  description = "TF Test"
  tags = {
    version = "1.0"
  }
}
`

const testAccSnapshotConfigUpdate = `
data "alicloud_instance_types" "instance_type" {
  availability_zone = "${data.alicloud_zones.zone.zones.0.id}"
}

resource "alicloud_vpc" "vpc" {
  name = "tf-testAcc-vpc"
  cidr_block = "192.168.0.0/16"
}

data "alicloud_zones" "zone" {
}

resource "alicloud_vswitch" "vswitch" {
  name = "tf-testAcc-vswitch"
  cidr_block = "192.168.0.0/24"
  availability_zone = "${data.alicloud_zones.zone.zones.0.id}"
  vpc_id = "${alicloud_vpc.vpc.id}"
}

resource "alicloud_security_group" "group" {
  name        = "tf-testACC-group"
  description = "New security group"
  vpc_id = "${alicloud_vpc.vpc.id}"
}

resource "alicloud_disk" "disk" {
  name = "disk"
  availability_zone = "${alicloud_instance.instance.availability_zone}"
  category          = "cloud_efficiency"
  size              = "20"
}

data "alicloud_images" "sys_images" {
  owners = "system"
}

resource "alicloud_instance" "instance" {
  instance_name   = "tf-testAcc-instance"
  host_name       = "tf-testAcc"
  image_id        = "${data.alicloud_images.sys_images.images.0.id}"
  instance_type   = "${data.alicloud_instance_types.instance_type.instance_types.0.id}"
  security_groups = ["${alicloud_security_group.group.id}"]
  vswitch_id      = "${alicloud_vswitch.vswitch.id}"
}

resource "alicloud_disk_attachment" "instance-attachment" {
  disk_id     = "${alicloud_disk.disk.id}"
  instance_id = "${alicloud_instance.instance.id}"
}

resource "alicloud_snapshot" "snapshot" {
  disk_id = "${alicloud_disk_attachment.instance-attachment.disk_id}"
  name = "tf-testAcc-snapshot"
  description = "TF Test"
  tags = {
    version = "1.1"
  }
}
`

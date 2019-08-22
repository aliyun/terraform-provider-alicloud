package alicloud

import (
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform/helper/acctest"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/hashicorp/terraform/helper/resource"
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
		return WrapError(err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}

	var snapshots []ecs.Snapshot
	request := ecs.CreateDescribeSnapshotsRequest()
	request.RegionId = client.RegionId
	request.PageSize = requests.NewInteger(PageSizeLarge)
	request.PageNumber = requests.NewInteger(1)
	for {
		raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.DescribeSnapshots(request)
		})
		if err != nil {
			return WrapError(err)
		}
		response, _ := raw.(*ecs.DescribeSnapshotsResponse)
		if len(response.Snapshots.Snapshot) < 1 {
			break
		}
		snapshots = append(snapshots, response.Snapshots.Snapshot...)

		if len(response.Snapshots.Snapshot) < PageSizeLarge {
			break
		}

		if page, err := getNextpageNumber(request.PageNumber); err != nil {
			return err
		} else {
			request.PageNumber = page
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

func TestAccAlicloudSnapshotBasic(t *testing.T) {

	var v *ecs.Snapshot
	resourceId := "alicloud_snapshot.default"
	rand := acctest.RandInt()
	name := fmt.Sprintf("tf-testAccSnapshotBasic%d", rand)
	ra := resourceAttrInit(resourceId, map[string]string{
		"disk_id":      CHECKSET,
		"name":         name,
		"description":  name,
		"tags.%":       "1",
		"tags.version": "1.0",
	})

	serviceFunc := func() interface{} {
		return &EcsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)

	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceSnapshotConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: resourceId,

		Providers:    testAccProviders,
		CheckDestroy: rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"disk_id":     "${alicloud_disk_attachment.default.0.disk_id}",
					"name":        "${var.name}",
					"description": "${var.name}",
					"tags": map[string]string{
						"version": "1.0",
					},
				}),
				Check: resource.ComposeTestCheckFunc(),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"version": "1.0",
						"Tag2":    "Tag2",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":    "2",
						"tags.Tag2": "Tag2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"version": "1.0",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":    "1",
						"tags.Tag2": REMOVEKEY,
					}),
				),
			},
		},
	})
}

func TestAccAlicloudSnapshotMulti(t *testing.T) {

	var v *ecs.Snapshot
	resourceId := "alicloud_snapshot.default.1"
	rand := acctest.RandInt()
	name := fmt.Sprintf("tf-testAccSnapshotMulti%d", rand)
	ra := resourceAttrInit(resourceId, map[string]string{
		"disk_id":      CHECKSET,
		"name":         name,
		"description":  name,
		"tags.%":       "1",
		"tags.version": "1.0",
	})

	serviceFunc := func() interface{} {
		return &EcsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)

	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceSnapshotConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: resourceId,

		Providers:    testAccProviders,
		CheckDestroy: rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"count":       "2",
					"disk_id":     "${element(alicloud_disk_attachment.default.*.disk_id,count.index)}",
					"name":        "${var.name}",
					"description": "${var.name}",
					"tags": map[string]string{
						"version": "1.0",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
		},
	})
}

func resourceSnapshotConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}


data "alicloud_instance_types" "default" {
  	cpu_core_count    = 1
	memory_size       = 2
}

resource "alicloud_vpc" "default" {
  name = "${var.name}"
  cidr_block = "192.168.0.0/16"
}


resource "alicloud_vswitch" "default" {
  name = "${var.name}"
  cidr_block = "192.168.0.0/24"
  availability_zone = "${data.alicloud_instance_types.default.instance_types.0.availability_zones.0}"
  vpc_id = "${alicloud_vpc.default.id}"
}

resource "alicloud_security_group" "default" {
  name = "${var.name}"
  description = "New security group"
  vpc_id = "${alicloud_vpc.default.id}"
}

resource "alicloud_disk" "default" {
  count = "2"
  name = "${var.name}"
  availability_zone = "${data.alicloud_instance_types.default.instance_types.0.availability_zones.0}"
  category          = "cloud_efficiency"
  size              = "20"
}

data "alicloud_images" "default" {
  owners = "system"
}

resource "alicloud_instance" "default" {
  availability_zone = "${data.alicloud_instance_types.default.instance_types.0.availability_zones.0}"
  instance_name   = "${var.name}"
  host_name       = "tf-testAcc"
  image_id        = "${data.alicloud_images.default.images.0.id}"
  instance_type   = "${data.alicloud_instance_types.default.instance_types.0.id}"
  security_groups = ["${alicloud_security_group.default.id}"]
  vswitch_id      = "${alicloud_vswitch.default.id}"
}

resource "alicloud_disk_attachment" "default" {
  count = "2"
  disk_id     = "${element(alicloud_disk.default.*.id,count.index)}"
  instance_id = "${alicloud_instance.default.id}"
}

`, name)
}

package alicloud

import (
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func init() {
	resource.AddTestSweepers("alicloud_disk", &resource.Sweeper{
		Name: "alicloud_disk",
		F:    testSweepDisks,
	})
}

func testSweepDisks(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
		"tf_test_",
		"tf-test-",
		"testAcc",
	}

	var disks []ecs.Disk
	req := ecs.CreateDescribeDisksRequest()
	req.RegionId = client.RegionId
	req.PageSize = requests.NewInteger(PageSizeLarge)
	req.PageNumber = requests.NewInteger(1)
	for {
		raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.DescribeDisks(req)
		})
		if err != nil {
			return WrapError(err)
		}
		resp, _ := raw.(*ecs.DescribeDisksResponse)
		if resp == nil || len(resp.Disks.Disk) < 1 {
			break
		}
		disks = append(disks, resp.Disks.Disk...)

		if len(resp.Disks.Disk) < PageSizeLarge {
			break
		}

		page, err := getNextpageNumber(req.PageNumber)
		if err != nil {
			return err
		}
		req.PageNumber = page
	}

	for _, v := range disks {
		name := v.DiskName
		id := v.DiskId
		skip := true
		for _, prefix := range prefixes {
			if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
				skip = false
				break
			}
		}
		if skip {
			log.Printf("[INFO] Skipping Disk: %s (%s)", name, id)
			continue
		}
		log.Printf("[INFO] Deleting Disk: %s (%s)", name, id)
		req := ecs.CreateDeleteDiskRequest()
		req.DiskId = id
		_, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.DeleteDisk(req)
		})
		if err != nil {
			log.Printf("[ERROR] Failed to delete Disk (%s (%s)): %s", name, id, err)
		}
	}
	return nil
}

func testAccCheckDiskExists(n string, disk *ecs.Disk) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return WrapError(fmt.Errorf("Not found: %s", n))
		}

		if rs.Primary.ID == "" {
			return WrapError(fmt.Errorf("No Disk ID is set"))
		}

		client := testAccProvider.Meta().(*connectivity.AliyunClient)
		ecsService := EcsService{client}

		d, err := ecsService.DescribeDisk(rs.Primary.ID)

		if err != nil {
			return WrapError(err)
		}

		*disk = d
		return nil
	}
}

func testAccCheckDiskDestroy(s *terraform.State) error {

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_disk" {
			continue
		}

		// Try to find the Disk
		client := testAccProvider.Meta().(*connectivity.AliyunClient)
		ecsService := EcsService{client}

		_, err := ecsService.DescribeDisk(rs.Primary.ID)

		if err != nil {
			if NotFoundError(err) {
				continue
			}
			return WrapError(err)
		}
	}

	return nil
}

func TestAccAlicloudDisk_basic(t *testing.T) {
	var v ecs.Disk
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "alicloud_disk.foo",

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDiskDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDiskConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDiskExists("alicloud_disk.foo", &v),
					resource.TestCheckResourceAttrSet("alicloud_disk.foo", "availability_zone"),
					resource.TestCheckResourceAttr("alicloud_disk.foo", "size", "50"),
					resource.TestCheckResourceAttr("alicloud_disk.foo", "name", ""),
					resource.TestCheckResourceAttr("alicloud_disk.foo", "description", ""),
					resource.TestCheckResourceAttr("alicloud_disk.foo", "category", "cloud_efficiency"),
					resource.TestCheckResourceAttr("alicloud_disk.foo", "snapshot_id", ""),
					resource.TestCheckResourceAttr("alicloud_disk.foo", "encrypted", "false"),
					resource.TestCheckNoResourceAttr("alicloud_disk.foo", "tags"),
					resource.TestCheckResourceAttr("alicloud_disk.foo", "status", string(Available)),
				),
			},
			{
				Config: testAccDiskConfig_size,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDiskExists("alicloud_disk.foo", &v),
					resource.TestCheckResourceAttrSet("alicloud_disk.foo", "availability_zone"),
					resource.TestCheckResourceAttr("alicloud_disk.foo", "size", "70"),
					resource.TestCheckResourceAttr("alicloud_disk.foo", "name", ""),
					resource.TestCheckResourceAttr("alicloud_disk.foo", "description", ""),
					resource.TestCheckResourceAttr("alicloud_disk.foo", "category", "cloud_efficiency"),
					resource.TestCheckResourceAttr("alicloud_disk.foo", "snapshot_id", ""),
					resource.TestCheckResourceAttr("alicloud_disk.foo", "encrypted", "false"),
					resource.TestCheckNoResourceAttr("alicloud_disk.foo", "tags"),
					resource.TestCheckResourceAttr("alicloud_disk.foo", "status", string(Available)),
				),
			},
			{
				Config: testAccDiskConfig_name,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDiskExists("alicloud_disk.foo", &v),
					resource.TestCheckResourceAttrSet("alicloud_disk.foo", "availability_zone"),
					resource.TestCheckResourceAttr("alicloud_disk.foo", "size", "70"),
					resource.TestCheckResourceAttr("alicloud_disk.foo", "name", "tf-testAccDiskConfig"),
					resource.TestCheckResourceAttr("alicloud_disk.foo", "description", ""),
					resource.TestCheckResourceAttr("alicloud_disk.foo", "category", "cloud_efficiency"),
					resource.TestCheckResourceAttr("alicloud_disk.foo", "snapshot_id", ""),
					resource.TestCheckResourceAttr("alicloud_disk.foo", "encrypted", "false"),
					resource.TestCheckNoResourceAttr("alicloud_disk.foo", "tags"),
					resource.TestCheckResourceAttr("alicloud_disk.foo", "status", string(Available)),
				),
			},
			{
				Config: testAccDiskConfig_description,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDiskExists("alicloud_disk.foo", &v),
					resource.TestCheckResourceAttrSet("alicloud_disk.foo", "availability_zone"),
					resource.TestCheckResourceAttr("alicloud_disk.foo", "size", "70"),
					resource.TestCheckResourceAttr("alicloud_disk.foo", "name", "tf-testAccDiskConfig"),
					resource.TestCheckResourceAttr("alicloud_disk.foo", "description", "tf-testAccDiskConfig_description"),
					resource.TestCheckResourceAttr("alicloud_disk.foo", "category", "cloud_efficiency"),
					resource.TestCheckResourceAttr("alicloud_disk.foo", "snapshot_id", ""),
					resource.TestCheckResourceAttr("alicloud_disk.foo", "encrypted", "false"),
					resource.TestCheckNoResourceAttr("alicloud_disk.foo", "tags"),
					resource.TestCheckResourceAttr("alicloud_disk.foo", "status", string(Available)),
				),
			},
			{
				Config: testAccDiskConfig_tags,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDiskExists("alicloud_disk.foo", &v),
					resource.TestCheckResourceAttrSet("alicloud_disk.foo", "availability_zone"),
					resource.TestCheckResourceAttr("alicloud_disk.foo", "size", "70"),
					resource.TestCheckResourceAttr("alicloud_disk.foo", "name", "tf-testAccDiskConfig"),
					resource.TestCheckResourceAttr("alicloud_disk.foo", "description", "tf-testAccDiskConfig_description"),
					resource.TestCheckResourceAttr("alicloud_disk.foo", "category", "cloud_efficiency"),
					resource.TestCheckResourceAttr("alicloud_disk.foo", "snapshot_id", ""),
					resource.TestCheckResourceAttr("alicloud_disk.foo", "encrypted", "false"),
					resource.TestCheckResourceAttr("alicloud_disk.foo", "tags.%", "3"),
					resource.TestCheckResourceAttr("alicloud_disk.foo", "status", string(Available)),
				),
			},
			{
				Config: testAccDiskConfig_all,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDiskExists("alicloud_disk.foo", &v),
					resource.TestCheckResourceAttrSet("alicloud_disk.foo", "availability_zone"),
					resource.TestCheckResourceAttr("alicloud_disk.foo", "size", "70"),
					resource.TestCheckResourceAttr("alicloud_disk.foo", "name", "tf-testAccDiskConfig_all"),
					resource.TestCheckResourceAttr("alicloud_disk.foo", "description", "nothing"),
					resource.TestCheckResourceAttr("alicloud_disk.foo", "category", "cloud_efficiency"),
					resource.TestCheckResourceAttr("alicloud_disk.foo", "snapshot_id", ""),
					resource.TestCheckResourceAttr("alicloud_disk.foo", "encrypted", "false"),
					resource.TestCheckResourceAttr("alicloud_disk.foo", "tags.%", "0"),
					resource.TestCheckResourceAttr("alicloud_disk.foo", "status", string(Available)),
				),
			},
		},
	})

}

func TestAccAlicloudDisk_multi(t *testing.T) {
	var v ecs.Disk
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "alicloud_disk.foo.4",

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDiskDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDiskConfig_multi,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDiskExists("alicloud_disk.foo.4", &v),
					resource.TestCheckResourceAttrSet("alicloud_disk.foo.4", "availability_zone"),
					resource.TestCheckResourceAttr("alicloud_disk.foo.4", "size", "70"),
					resource.TestCheckResourceAttr("alicloud_disk.foo.4", "name", "tf-testAccDiskConfig_multi"),
					resource.TestCheckResourceAttr("alicloud_disk.foo.4", "description", "nothing"),
					resource.TestCheckResourceAttr("alicloud_disk.foo.4", "category", "cloud_efficiency"),
					resource.TestCheckResourceAttr("alicloud_disk.foo.4", "snapshot_id", ""),
					resource.TestCheckResourceAttr("alicloud_disk.foo.4", "encrypted", "false"),
					resource.TestCheckNoResourceAttr("alicloud_disk.foo.4", "tags.%"),
					resource.TestCheckResourceAttr("alicloud_disk.foo.4", "status", string(Available)),
				),
			},
		},
	})

}

const testAccDiskConfig_basic = `
data "alicloud_zones" "az" {
	"available_resource_creation"= "VSwitch"
}

resource "alicloud_disk" "foo" {
	availability_zone = "${data.alicloud_zones.az.zones.0.id}"
  	size = "50"
}
`

const testAccDiskConfig_size = `
data "alicloud_zones" "az" {
	"available_resource_creation"= "VSwitch"
}


resource "alicloud_disk" "foo" {
	availability_zone = "${data.alicloud_zones.az.zones.0.id}"
  	size = "70"
}
`
const testAccDiskConfig_name = `
data "alicloud_zones" "az" {
	"available_resource_creation"= "VSwitch"
}


variable "name" {
	default = "tf-testAccDiskConfig"
}

resource "alicloud_disk" "foo" {
	availability_zone = "${data.alicloud_zones.az.zones.0.id}"
  	size = "70"
	name = "${var.name}"
}
`

const testAccDiskConfig_description = `
data "alicloud_zones" "az" {
	"available_resource_creation"= "VSwitch"
}


variable "name" {
	default = "tf-testAccDiskConfig"
}

resource "alicloud_disk" "foo" {
	availability_zone = "${data.alicloud_zones.az.zones.0.id}"
  	size = "70"
	name = "${var.name}"
	description = "${var.name}_description"
}
`

const testAccDiskConfig_tags = `
data "alicloud_zones" "az" {
	"available_resource_creation"= "VSwitch"
}


variable "name" {
	default = "tf-testAccDiskConfig"
}

resource "alicloud_disk" "foo" {
	availability_zone = "${data.alicloud_zones.az.zones.0.id}"
  	size = "70"
	name = "${var.name}"
	description = "${var.name}_description"
	category = "cloud_efficiency"
	encrypted = "false"
	tags = {
		name1 = "name1"
		name2 = "name2"
		name3 = "name3"
			}
}
`

const testAccDiskConfig_all = `
data "alicloud_zones" "az" {
	"available_resource_creation"= "VSwitch"
}

variable "name" {
	default = "tf-testAccDiskConfig"
}

resource "alicloud_disk" "foo" {
	availability_zone = "${data.alicloud_zones.az.zones.0.id}"
  	size = "70"
	name = "${var.name}_all"
	description = "nothing"
	category = "cloud_efficiency"
	encrypted = "false"
}
`

const testAccDiskConfig_multi = `
data "alicloud_zones" "az" {
	"available_resource_creation"= "VSwitch"
}

variable "name" {
	default = "tf-testAccDiskConfig"
}

resource "alicloud_disk" "foo" {
	count = "5"
	availability_zone = "${data.alicloud_zones.az.zones.0.id}"
  	size = "70"
	name = "${var.name}_multi"
	description = "nothing"
	category = "cloud_efficiency"
	encrypted = "false"
}
`

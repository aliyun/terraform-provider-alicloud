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
			return fmt.Errorf("Error retrieving Disks: %s", err)
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

func TestAccAlicloudDisk_Resize(t *testing.T) {
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
				Config: testAccDiskConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDiskExists("alicloud_disk.foo", &v),
					resource.TestCheckResourceAttr("alicloud_disk.foo", "category", "cloud_efficiency"),
					resource.TestCheckResourceAttr("alicloud_disk.foo", "size", "30"),
					resource.TestCheckResourceAttr("alicloud_disk.foo", "encrypted", "false"),
					resource.TestCheckResourceAttrSet("alicloud_disk.foo", "availability_zone"),
					resource.TestCheckResourceAttr("alicloud_disk.foo", "description", "Hello ecs disk."),
					resource.TestCheckResourceAttr("alicloud_disk.foo", "name", "tf-testAccDiskConfig"),
					resource.TestCheckResourceAttr("alicloud_disk.foo", "snapshot_id", ""),
					resource.TestCheckResourceAttr("alicloud_disk.foo", "status", string(Available)),
					resource.TestCheckNoResourceAttr("alicloud_disk.foo", "tags"),
				),
			},
			{
				Config: testAccDiskConfigResize,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDiskExists("alicloud_disk.foo", &v),
					resource.TestCheckResourceAttr("alicloud_disk.foo", "category", "cloud_efficiency"),
					resource.TestCheckResourceAttr("alicloud_disk.foo", "size", "40"),
					resource.TestCheckResourceAttr("alicloud_disk.foo", "encrypted", "false"),
					resource.TestCheckResourceAttrSet("alicloud_disk.foo", "availability_zone"),
					resource.TestCheckResourceAttr("alicloud_disk.foo", "description", "Hello ecs disk."),
					resource.TestCheckResourceAttr("alicloud_disk.foo", "name", "tf-testAccDiskConfig"),
					resource.TestCheckResourceAttr("alicloud_disk.foo", "snapshot_id", ""),
					resource.TestCheckResourceAttr("alicloud_disk.foo", "status", string(Available)),
					resource.TestCheckNoResourceAttr("alicloud_disk.foo", "tags"),
				),
			},
		},
	})
}
func TestAccAlicloudDisk_Multi(t *testing.T) {
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
				Config: testAccDiskConfigMulti,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDiskExists("alicloud_disk.foo.4", &v),
					resource.TestCheckResourceAttr("alicloud_disk.foo.4", "name", "tf-testAccDiskConfigMulti-4"),
					testAccCheckDiskExists("alicloud_disk.foo.3", &v),
					resource.TestCheckResourceAttr("alicloud_disk.foo.3", "name", "tf-testAccDiskConfigMulti-3"),
					testAccCheckDiskExists("alicloud_disk.foo.2", &v),
					resource.TestCheckResourceAttr("alicloud_disk.foo.2", "name", "tf-testAccDiskConfigMulti-2"),
					testAccCheckDiskExists("alicloud_disk.foo.1", &v),
					resource.TestCheckResourceAttr("alicloud_disk.foo.1", "name", "tf-testAccDiskConfigMulti-1"),
					testAccCheckDiskExists("alicloud_disk.foo.0", &v),
					resource.TestCheckResourceAttr("alicloud_disk.foo.0", "name", "tf-testAccDiskConfigMulti-0"),
				),
			},
		},
	})
}

func TestAccAlicloudDisk_withTags(t *testing.T) {
	var v ecs.Disk

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		//module name
		IDRefreshName: "alicloud_disk.bar",

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDiskDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDiskConfigWithTags,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDiskExists("alicloud_disk.bar", &v),
					resource.TestCheckResourceAttr("alicloud_disk.bar", "tags.%", "6"),
					resource.TestCheckResourceAttr("alicloud_disk.bar", "tags.Name", "TerraformTest"),
					resource.TestCheckResourceAttr("alicloud_disk.bar", "category", "cloud_efficiency"),
					resource.TestCheckResourceAttr("alicloud_disk.bar", "size", "20"),
					resource.TestCheckResourceAttr("alicloud_disk.bar", "encrypted", "false"),
					resource.TestCheckResourceAttrSet("alicloud_disk.bar", "availability_zone"),
					resource.TestCheckResourceAttr("alicloud_disk.bar", "description", ""),
					resource.TestCheckResourceAttr("alicloud_disk.bar", "name", "tf-testAccDiskConfigWithTags"),
					resource.TestCheckResourceAttr("alicloud_disk.bar", "snapshot_id", ""),
					resource.TestCheckResourceAttr("alicloud_disk.bar", "status", string(Available)),
				),
			},
			{
				Config: testAccDiskConfigWithTags1,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDiskExists("alicloud_disk.bar", &v),
					resource.TestCheckResourceAttr("alicloud_disk.bar", "tags.%", "7"),
					resource.TestCheckResourceAttr("alicloud_disk.bar", "tags.Name6", "hahah"),
					resource.TestCheckResourceAttr("alicloud_disk.bar", "category", "cloud_efficiency"),
					resource.TestCheckResourceAttr("alicloud_disk.bar", "size", "20"),
					resource.TestCheckResourceAttr("alicloud_disk.bar", "encrypted", "false"),
					resource.TestCheckResourceAttrSet("alicloud_disk.bar", "availability_zone"),
					resource.TestCheckResourceAttr("alicloud_disk.bar", "description", ""),
					resource.TestCheckResourceAttr("alicloud_disk.bar", "name", "tf-testAccDiskConfigWithTags"),
					resource.TestCheckResourceAttr("alicloud_disk.bar", "snapshot_id", ""),
					resource.TestCheckResourceAttr("alicloud_disk.bar", "status", string(Available)),
				),
			},
			{
				Config: testAccDiskConfigWithTags2,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDiskExists("alicloud_disk.bar", &v),
					resource.TestCheckResourceAttr("alicloud_disk.bar", "tags.%", "7"),
					resource.TestCheckResourceAttr("alicloud_disk.bar", "tags.Name7", "please"),
					resource.TestCheckResourceAttr("alicloud_disk.bar", "category", "cloud_efficiency"),
					resource.TestCheckResourceAttr("alicloud_disk.bar", "size", "20"),
					resource.TestCheckResourceAttr("alicloud_disk.bar", "encrypted", "false"),
					resource.TestCheckResourceAttrSet("alicloud_disk.bar", "availability_zone"),
					resource.TestCheckResourceAttr("alicloud_disk.bar", "description", ""),
					resource.TestCheckResourceAttr("alicloud_disk.bar", "name", "tf-testAccDiskConfigWithTags"),
					resource.TestCheckResourceAttr("alicloud_disk.bar", "snapshot_id", ""),
					resource.TestCheckResourceAttr("alicloud_disk.bar", "status", string(Available)),
				),
			},
			{
				Config: testAccDiskConfigWithTags3,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDiskExists("alicloud_disk.bar", &v),
					resource.TestCheckResourceAttr("alicloud_disk.bar", "tags.%", "6"),
					resource.TestCheckResourceAttr("alicloud_disk.bar", "tags.Name", "TerraformTest"),
					resource.TestCheckResourceAttr("alicloud_disk.bar", "category", "cloud_efficiency"),
					resource.TestCheckResourceAttr("alicloud_disk.bar", "size", "20"),
					resource.TestCheckResourceAttr("alicloud_disk.bar", "encrypted", "false"),
					resource.TestCheckResourceAttrSet("alicloud_disk.bar", "availability_zone"),
					resource.TestCheckResourceAttr("alicloud_disk.bar", "description", ""),
					resource.TestCheckResourceAttr("alicloud_disk.bar", "name", "tf-testAccDiskConfigWithTags"),
					resource.TestCheckResourceAttr("alicloud_disk.bar", "snapshot_id", ""),
					resource.TestCheckResourceAttr("alicloud_disk.bar", "status", string(Available)),
				),
			},
		},
	})
}

func TestAccAlicloudDisk_encrypted(t *testing.T) {
	var v ecs.Disk

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "alicloud_disk.encrypted",

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDiskDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDiskConfigEncrypted,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDiskExists("alicloud_disk.encrypted", &v),
					resource.TestCheckResourceAttr("alicloud_disk.encrypted", "category", "cloud_efficiency"),
					resource.TestCheckResourceAttr("alicloud_disk.encrypted", "size", "30"),
					resource.TestCheckResourceAttr("alicloud_disk.encrypted", "encrypted", "true"),
					resource.TestCheckResourceAttrSet("alicloud_disk.encrypted", "availability_zone"),
					resource.TestCheckResourceAttr("alicloud_disk.encrypted", "description", "Hello ecs disk."),
					resource.TestCheckResourceAttr("alicloud_disk.encrypted", "name", "tf-testAccDiskConfigEncrypted"),
					resource.TestCheckResourceAttr("alicloud_disk.encrypted", "snapshot_id", ""),
					resource.TestCheckResourceAttr("alicloud_disk.encrypted", "status", string(Available)),
					resource.TestCheckNoResourceAttr("alicloud_disk.encrypted", "tags"),
				),
			},
		},
	})
}
func TestAccAlicloudDisk_ModifyDiskAttribute(t *testing.T) {
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
				Config: testAccDiskConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDiskExists("alicloud_disk.foo", &v),
					resource.TestCheckResourceAttr("alicloud_disk.foo", "category", "cloud_efficiency"),
					resource.TestCheckResourceAttr("alicloud_disk.foo", "size", "30"),
					resource.TestCheckResourceAttr("alicloud_disk.foo", "encrypted", "false"),
					resource.TestCheckResourceAttrSet("alicloud_disk.foo", "availability_zone"),
					resource.TestCheckResourceAttr("alicloud_disk.foo", "description", "Hello ecs disk."),
					resource.TestCheckResourceAttr("alicloud_disk.foo", "name", "tf-testAccDiskConfig"),
					resource.TestCheckResourceAttr("alicloud_disk.foo", "snapshot_id", ""),
					resource.TestCheckResourceAttr("alicloud_disk.foo", "status", string(Available)),
					resource.TestCheckNoResourceAttr("alicloud_disk.foo", "tags"),
				),
			},
			{
				Config: testAccDiskConfigRename,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDiskExists("alicloud_disk.foo", &v),
					resource.TestCheckResourceAttr("alicloud_disk.foo", "category", "cloud_efficiency"),
					resource.TestCheckResourceAttr("alicloud_disk.foo", "size", "30"),
					resource.TestCheckResourceAttr("alicloud_disk.foo", "encrypted", "false"),
					resource.TestCheckResourceAttrSet("alicloud_disk.foo", "availability_zone"),
					resource.TestCheckResourceAttr("alicloud_disk.foo", "description", "Hello ecs disk."),
					resource.TestCheckResourceAttr("alicloud_disk.foo", "name", "tf-testAccDiskConfigRedescription"),
					resource.TestCheckResourceAttr("alicloud_disk.foo", "snapshot_id", ""),
					resource.TestCheckResourceAttr("alicloud_disk.foo", "status", string(Available)),
					resource.TestCheckNoResourceAttr("alicloud_disk.foo", "tags"),
				),
			},
			{
				Config: testAccDiskConfigRedescription,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDiskExists("alicloud_disk.foo", &v),
					resource.TestCheckResourceAttr("alicloud_disk.foo", "category", "cloud_efficiency"),
					resource.TestCheckResourceAttr("alicloud_disk.foo", "size", "30"),
					resource.TestCheckResourceAttr("alicloud_disk.foo", "encrypted", "false"),
					resource.TestCheckResourceAttrSet("alicloud_disk.foo", "availability_zone"),
					resource.TestCheckResourceAttr("alicloud_disk.foo", "description", "Hello MrX."),
					resource.TestCheckResourceAttr("alicloud_disk.foo", "name", "tf-testAccDiskConfigRedescription"),
					resource.TestCheckResourceAttr("alicloud_disk.foo", "snapshot_id", ""),
					resource.TestCheckResourceAttr("alicloud_disk.foo", "status", string(Available)),
					resource.TestCheckNoResourceAttr("alicloud_disk.foo", "tags"),
				),
			},
		},
	})
}

func testAccCheckDiskExists(n string, disk *ecs.Disk) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Disk ID is set")
		}

		client := testAccProvider.Meta().(*connectivity.AliyunClient)
		ecsService := EcsService{client}

		d, err := ecsService.DescribeDiskById("", rs.Primary.ID)

		if err != nil {
			return fmt.Errorf("While checking disk existing, describing disk got an error: %#v.", err)
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

		d, err := ecsService.DescribeDiskById("", rs.Primary.ID)

		if err != nil {
			if NotFoundError(err) {
				continue
			}
			return fmt.Errorf("While checking disk destroy, describing disk got an error: %#v.", err)
		}

		if d.DiskId != "" {
			return fmt.Errorf("Error ECS Disk still exist")
		}
	}

	return nil
}

const testAccDiskConfig = `
data "alicloud_zones" "default" {
	"available_disk_category"= "cloud_efficiency"
}
variable "name" {
	default = "tf-testAccDiskConfig"
}
resource "alicloud_disk" "foo" {
	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
	name = "${var.name}"
	description = "Hello ecs disk."
	category = "cloud_efficiency"
  	size = "30"
}
`
const testAccDiskConfigResize = `
data "alicloud_zones" "default" {
	"available_disk_category"= "cloud_efficiency"
}
variable "name" {
	default = "tf-testAccDiskConfig"
}
resource "alicloud_disk" "foo" {
	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
	name = "${var.name}"
	description = "Hello ecs disk."
	category = "cloud_efficiency"
	size = "40"
}
`
const testAccDiskConfigWithTags = `
data "alicloud_zones" "default" {
	"available_disk_category"= "cloud_efficiency"
}
variable "name" {
	default = "tf-testAccDiskConfigWithTags"
}
resource "alicloud_disk" "bar" {
	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
	category = "cloud_efficiency"
	name = "${var.name}"
	size = "20"
	tags {
	    Name = "TerraformTest"
	    Name1 = "TerraformTest"
	    Name2 = "TerraformTest"
	    Name3 = "TerraformTest"
	    Name4 = "TerraformTest"
	    Name5 = "TerraformTest"
	}
}
`
const testAccDiskConfigWithTags1 = `
data "alicloud_zones" "default" {
	"available_disk_category"= "cloud_efficiency"
}
variable "name" {
	default = "tf-testAccDiskConfigWithTags"
}
resource "alicloud_disk" "bar" {
	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
	category = "cloud_efficiency"
	name = "${var.name}"
	size = "20"
	tags {
	    Name = "TerraformTest"
	    Name1 = "TerraformTest"
	    Name2 = "TerraformTest"
	    Name3 = "TerraformTest"
		Name4 = "TerraformTest"
		Name5 = "TerraformTest"
	    Name6 = "hahah"
	}
}
`
const testAccDiskConfigWithTags2 = `
data "alicloud_zones" "default" {
	"available_disk_category"= "cloud_efficiency"
}
variable "name" {
	default = "tf-testAccDiskConfigWithTags"
}
resource "alicloud_disk" "bar" {
	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
	category = "cloud_efficiency"
	name = "${var.name}"
	size = "20"
	tags {
	    Name = "TerraformTest"
	    Name1 = "TerraformTest"
	    Name2 = "TerraformTest"
	    Name3 = "TerraformTest"
	    Name4 = "TerraformTest"
		Name5 = "TerraformTest"
		Name7 = "please"
	}
}
`
const testAccDiskConfigWithTags3 = `
data "alicloud_zones" "default" {
	"available_disk_category"= "cloud_efficiency"
}
variable "name" {
	default = "tf-testAccDiskConfigWithTags"
}
resource "alicloud_disk" "bar" {
	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
	category = "cloud_efficiency"
	name = "${var.name}"
	size = "20"
	tags {
	    Name = "TerraformTest"
	    Name1 = "TerraformTest"
	    Name2 = "TerraformTest"
	    Name3 = "TerraformTest"
	    Name4 = "TerraformTest"
	    Name5 = "TerraformTest"
	}
}
`
const testAccDiskConfigEncrypted = `
data "alicloud_zones" "default" {
	"available_disk_category"= "cloud_efficiency"
}
variable "name" {
	default = "tf-testAccDiskConfigEncrypted"
}
resource "alicloud_disk" "encrypted" {
	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
	name = "${var.name}"
	description = "Hello ecs disk."
	category = "cloud_efficiency"
  	size = "30"
	encrypted = true
}
`
const testAccDiskConfigRename = `
data "alicloud_zones" "default" {
	"available_disk_category"= "cloud_efficiency"
}
variable "name" {
	default = "tf-testAccDiskConfigRedescription"
}
resource "alicloud_disk" "foo" {	
	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
	name = "${var.name}"
	description = "Hello ecs disk."
	category = "cloud_efficiency"
  	size = "30"
}
`
const testAccDiskConfigRedescription = `
data "alicloud_zones" "default" {
	"available_disk_category"= "cloud_efficiency"
}
variable "name" {
	default = "tf-testAccDiskConfigRedescription"
}
resource "alicloud_disk" "foo" {
	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
	name = "${var.name}"
	description = "Hello MrX."
	category = "cloud_efficiency"
  	size = "30"
}
`

const testAccDiskConfigMulti = `
data "alicloud_zones" "default" {
	"available_disk_category"= "cloud_efficiency"
}
variable "name" {
	default = "tf-testAccDiskConfigMulti"
}
resource "alicloud_disk" "foo" {
	count = 5
	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
	name = "${var.name}-${count.index}"
	description = "Hello MrX."
	category = "cloud_efficiency"
  	size = "30"
}
`

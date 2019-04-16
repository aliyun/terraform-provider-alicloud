package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/nas"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func TestAccAlicloudNas_MountTarget_basic(t *testing.T) {
	var mt nas.DescribeMountTargetsMountTarget1
	rand1 := acctest.RandIntRange(10000, 499999)
	rand2 := acctest.RandIntRange(500000, 999999)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, connectivity.NasClassicSupportedRegions)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMountTargetDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccNasMountTargetConfig(rand1, rand2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMountTargetExists("alicloud_nas_mount_target.foo", &mt),
					resource.TestCheckResourceAttr(
						"alicloud_nas_mount_target.foo", "access_group_name", fmt.Sprintf("tf-testAccNasConfig-%d", rand1)),
					resource.TestCheckResourceAttr("alicloud_nas_mount_target.foo", "vswitch_id", ""),
					resource.TestCheckResourceAttr("alicloud_nas_mount_target.foo", "status", "Active"),
					resource.TestCheckResourceAttrSet("alicloud_nas_mount_target.foo", "file_system_id"),
				),
			},
		},
	})
}

func TestAccAlicloudNas_MountTarget_Vpc_basic(t *testing.T) {
	var mt nas.DescribeMountTargetsMountTarget1
	rand1 := acctest.RandIntRange(10000, 499999)
	rand2 := acctest.RandIntRange(500000, 999999)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMountTargetDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccNasMountTargetVpcConfig(rand1, rand2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMountTargetExists("alicloud_nas_mount_target.foo", &mt),
					resource.TestCheckResourceAttr(
						"alicloud_nas_mount_target.foo", "access_group_name", fmt.Sprintf("tf-testAccNasConfig-%d", rand1)),
					resource.TestCheckResourceAttrSet("alicloud_nas_mount_target.foo", "vswitch_id"),
					resource.TestCheckResourceAttr("alicloud_nas_mount_target.foo", "status", "Active"),
					resource.TestCheckResourceAttrSet("alicloud_nas_mount_target.foo", "file_system_id"),
				),
			},
		},
	})
}

func TestAccAlicloudNas_MountTarget_update(t *testing.T) {
	var mt nas.DescribeMountTargetsMountTarget1
	rand1 := acctest.RandIntRange(10000, 499999)
	rand2 := acctest.RandIntRange(500000, 999999)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: "alicloud_nas_mount_target.foo",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckMountTargetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNasMountTargetVpcConfig(rand1, rand2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMountTargetExists("alicloud_nas_mount_target.foo", &mt),
					resource.TestCheckResourceAttr(
						"alicloud_nas_mount_target.foo", "access_group_name", fmt.Sprintf("tf-testAccNasConfig-%d", rand1)),
					resource.TestCheckResourceAttrSet("alicloud_nas_mount_target.foo", "vswitch_id"),
					resource.TestCheckResourceAttr("alicloud_nas_mount_target.foo", "status", "Active"),
					resource.TestCheckResourceAttrSet("alicloud_nas_mount_target.foo", "file_system_id"),
				),
			},
			{
				Config: testAccNasMountTargetConfigUpdateAccessGroup(rand1, rand2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMountTargetExists("alicloud_nas_mount_target.foo", &mt),
					resource.TestCheckResourceAttr(
						"alicloud_nas_mount_target.foo", "status", "Active"),
					resource.TestCheckResourceAttrSet("alicloud_nas_mount_target.foo", "vswitch_id"),
					resource.TestCheckResourceAttrSet("alicloud_nas_mount_target.foo", "file_system_id"),
					resource.TestCheckResourceAttr("alicloud_nas_mount_target.foo", "access_group_name", fmt.Sprintf("tf-testAccNasConfig-2-%d", rand2)),
				),
			},
			{
				Config: testAccNasMountTargetConfigUpdateStatus(rand1, rand2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMountTargetExists("alicloud_nas_mount_target.foo", &mt),
					resource.TestCheckResourceAttr(
						"alicloud_nas_mount_target.foo", "status", "Inactive"),
					resource.TestCheckResourceAttrSet("alicloud_nas_mount_target.foo", "vswitch_id"),
					resource.TestCheckResourceAttrSet("alicloud_nas_mount_target.foo", "file_system_id"),
					resource.TestCheckResourceAttr("alicloud_nas_mount_target.foo", "access_group_name", fmt.Sprintf("tf-testAccNasConfig-2-%d", rand2)),
				),
			},
			{
				Config: testAccNasMountTargetConfigUpdateAll(rand1, rand2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMountTargetExists("alicloud_nas_mount_target.foo", &mt),
					resource.TestCheckResourceAttr(
						"alicloud_nas_mount_target.foo", "status", "Active"),
					resource.TestCheckResourceAttrSet("alicloud_nas_mount_target.foo", "vswitch_id"),
					resource.TestCheckResourceAttrSet("alicloud_nas_mount_target.foo", "file_system_id"),
					resource.TestCheckResourceAttr("alicloud_nas_mount_target.foo", "access_group_name", fmt.Sprintf("tf-testAccNasConfig-%d", rand1)),
				),
			},
		},
	})
}

func testAccCheckMountTargetExists(n string, nas *nas.DescribeMountTargetsMountTarget1) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return WrapError(fmt.Errorf("Not found: %s", n))
		}

		if rs.Primary.ID == "" {
			return WrapError(fmt.Errorf("No NAS ID is set"))
		}
		client := testAccProvider.Meta().(*connectivity.AliyunClient)
		nasService := NasService{client}
		instance, err := nasService.DescribeNasMountTarget(rs.Primary.ID)
		if err != nil {
			return WrapError(err)
		}

		*nas = instance
		return nil
	}
}

func testAccCheckMountTargetDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.AliyunClient)
	nasService := NasService{client}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_nas_mount_target" {
			continue
		}
		instance, err := nasService.DescribeNasMountTarget(rs.Primary.ID)

		if err != nil {
			if NotFoundError(err) {
				continue
			}
			return WrapError(err)
		}
		return WrapError(fmt.Errorf("NAS %s still exist", instance.MountTargetDomain))
	}
	return nil
}

func testAccNasMountTargetConfig(rand1 int, rand2 int) string {
	return fmt.Sprintf(`
	resource "alicloud_nas_file_system" "foo" {
			protocol_type = "NFS"
			storage_type = "Performance"
			description = "tf-testAccNasConfigFs"
	}
	resource "alicloud_nas_access_group" "foo" {
			name = "tf-testAccNasConfig-%d"
			type = "Classic"
			description = "tf-testAccNasConfig"
	}
	resource "alicloud_nas_access_group" "bar" {
			name = "tf-testAccNasConfig-2-%d"
			type = "Classic"
			description = "tf-testAccNasConfig-2"
	}
	resource "alicloud_nas_mount_target" "foo" {
			file_system_id = "${alicloud_nas_file_system.foo.id}"
			access_group_name = "${alicloud_nas_access_group.foo.id}"
	}`, rand1, rand2)
}

func testAccNasMountTargetVpcConfig(rand1 int, rand2 int) string {
	return fmt.Sprintf(`
	variable "name" {
			default = "tf-testAccVswitch"
	}
	data "alicloud_zones" "default" {
			"available_resource_creation"= "VSwitch"
	}
	resource "alicloud_vpc" "foo" {
			name = "${var.name}"
			cidr_block = "172.16.0.0/12"
	}
	resource "alicloud_vswitch" "foo" {
			vpc_id = "${alicloud_vpc.foo.id}"
			cidr_block = "172.16.0.0/24"
			availability_zone = "${data.alicloud_zones.default.zones.0.id}"
			name = "${var.name}-1"
	}
	resource "alicloud_nas_file_system" "foo" {
        	        protocol_type = "NFS"
                	storage_type = "Performance"
			description = "tf-testAccNasConfigFs"
	}
	resource "alicloud_nas_access_group" "foo" {
        	        name = "tf-testAccNasConfig-%d"
                	type = "Vpc"
	                description = "tf-testAccNasConfig"
	}
	resource "alicloud_nas_access_group" "bar" {
        	        name = "tf-testAccNasConfig-2-%d"
                	type = "Vpc"
	                description = "tf-testAccNasConfig-2"
	}
	resource "alicloud_nas_mount_target" "foo" {
        	        file_system_id = "${alicloud_nas_file_system.foo.id}"
                	access_group_name = "${alicloud_nas_access_group.foo.id}"
			vswitch_id = "${alicloud_vswitch.foo.id}"               
	}
`, rand1, rand2)
}

func testAccNasMountTargetConfigUpdateAccessGroup(rand1 int, rand2 int) string {
	return fmt.Sprintf(`
	variable "name" {
        	        default = "tf-testAccVswitch"
	}
	data "alicloud_zones" "default" {
        	        "available_resource_creation"= "VSwitch"
	}
	resource "alicloud_vpc" "foo" {
        	        name = "${var.name}"
                	cidr_block = "172.16.0.0/12"
	}
	resource "alicloud_vswitch" "foo" {
        	        vpc_id = "${alicloud_vpc.foo.id}"
                	cidr_block = "172.16.0.0/24"
	                availability_zone = "${data.alicloud_zones.default.zones.0.id}"
        	        name = "${var.name}-1"
	}
	resource "alicloud_nas_file_system" "foo" {
        	        protocol_type = "NFS"
                	storage_type = "Performance"
	                description = "tf-testAccNasConfigFs"
	}
	resource "alicloud_nas_access_group" "foo" {
        	        name = "tf-testAccNasConfig-%d"
                	type = "Vpc"
	                description = "tf-testAccNasConfig"
	}
	resource "alicloud_nas_access_group" "bar" {
        	        name = "tf-testAccNasConfig-2-%d"
                	type = "Vpc"
	                description = "tf-testAccNasConfig-2"
	}
	resource "alicloud_nas_mount_target" "foo" {
			file_system_id = "${alicloud_nas_file_system.foo.id}"
			access_group_name = "${alicloud_nas_access_group.bar.id}"
			vswitch_id = "${alicloud_vswitch.foo.id}"
	}`, rand1, rand2)
}

func testAccNasMountTargetConfigUpdateStatus(rand1 int, rand2 int) string {
	return fmt.Sprintf(`
	variable "name" {
        	        default = "tf-testAccVswitch"
	}
	data "alicloud_zones" "default" {
        	        "available_resource_creation"= "VSwitch"
	}
	resource "alicloud_vpc" "foo" {
        	        name = "${var.name}"
                	cidr_block = "172.16.0.0/12"
	}
	resource "alicloud_vswitch" "foo" {
        	        vpc_id = "${alicloud_vpc.foo.id}"
                	cidr_block = "172.16.0.0/24"
	                availability_zone = "${data.alicloud_zones.default.zones.0.id}"
        	        name = "${var.name}-1"
	}
	resource "alicloud_nas_file_system" "foo" {
        	        protocol_type = "NFS"
                	storage_type = "Performance"
	                description = "tf-testAccNasConfigFs"
	}
	resource "alicloud_nas_access_group" "foo" {
        	        name = "tf-testAccNasConfig-%d"
                	type = "Vpc"
	                description = "tf-testAccNasConfig"
	}
	resource "alicloud_nas_access_group" "bar" {
        	        name = "tf-testAccNasConfig-2-%d"
                	type = "Vpc"
	                description = "tf-testAccNasConfig-2"
	}
	resource "alicloud_nas_mount_target" "foo" {
			file_system_id = "${alicloud_nas_file_system.foo.id}"
                	access_group_name = "${alicloud_nas_access_group.bar.id}"
	                status = "Inactive"
			vswitch_id = "${alicloud_vswitch.foo.id}"
	}`, rand1, rand2)
}
func testAccNasMountTargetConfigUpdateAll(rand1 int, rand2 int) string {
	return fmt.Sprintf(`
	variable "name" {
        	        default = "tf-testAccVswitch"
	}
	data "alicloud_zones" "default" {
        	        "available_resource_creation"= "VSwitch"
	}
	resource "alicloud_vpc" "foo" {
        	        name = "${var.name}"
                	cidr_block = "172.16.0.0/12"
	}
	resource "alicloud_vswitch" "foo" {
        	        vpc_id = "${alicloud_vpc.foo.id}"
                	cidr_block = "172.16.0.0/24"
	                availability_zone = "${data.alicloud_zones.default.zones.0.id}"
        	        name = "${var.name}-1"
	}
	resource "alicloud_nas_file_system" "foo" {
        	        protocol_type = "NFS"
                	storage_type = "Performance"
	                description = "tf-testAccNasConfigFs"
	}
	resource "alicloud_nas_access_group" "foo" {
        	        name = "tf-testAccNasConfig-%d"
	                type = "Vpc"
        	        description = "tf-testAccNasConfig"
	}
	resource "alicloud_nas_access_group" "bar" {
        	        name = "tf-testAccNasConfig-2-%d"
	                type = "Vpc"
        	        description = "tf-testAccNasConfig-2"
	}
	resource "alicloud_nas_mount_target" "foo" {
        	        file_system_id = "${alicloud_nas_file_system.foo.id}"
                	access_group_name = "${alicloud_nas_access_group.foo.id}"
	                status = "Active"
			vswitch_id = "${alicloud_vswitch.foo.id}"
	}`, rand1, rand2)
}

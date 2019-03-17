package alicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func TestAccAlicloudNasMountTargetDataSourceAccessGroupName(t *testing.T) {
	rand := acctest.RandIntRange(10000, 999999)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, connectivity.NasClassicSupportedRegions)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudMountTargetDataSourceAccessGroupName(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_nas_mount_targets.mt"),
					resource.TestCheckResourceAttrSet("data.alicloud_nas_mount_targets.mt", "targets.0.mount_target_domain"),
					resource.TestCheckResourceAttr("data.alicloud_nas_mount_targets.mt", "targets.0.type", "Classic"),
					resource.TestCheckResourceAttr("data.alicloud_nas_mount_targets.mt", "targets.0.vpc_id", ""),
					resource.TestCheckResourceAttr("data.alicloud_nas_mount_targets.mt", "targets.0.vswitch_id", ""),
					resource.TestCheckResourceAttr("data.alicloud_nas_mount_targets.mt", "targets.0.access_group_name", fmt.Sprintf("tf-testAccNasConfig-%d", rand)),
					resource.TestCheckResourceAttr("data.alicloud_nas_mount_targets.mt", "ids.#", "1"),
					resource.TestCheckResourceAttrSet("data.alicloud_nas_mount_targets.mt", "ids.0"),
				),
			},
			{
				Config: testAccCheckAlicloudMountTargetDataSourceAccessGroupNameEmpty(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_nas_mount_targets.mt"),
					resource.TestCheckResourceAttr("data.alicloud_nas_mount_targets.mt", "targets.#", "0"),
					resource.TestCheckResourceAttr("data.alicloud_nas_mount_targets.mt", "ids.#", "0"),
				),
			},
		},
	})
}

func TestAccAlicloudNasMountTargetDataSourceType(t *testing.T) {
	rand := acctest.RandIntRange(10000, 999999)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudMountTargetDataSourceType(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_nas_mount_targets.mt"),
					resource.TestCheckResourceAttr("data.alicloud_nas_mount_targets.mt", "targets.#", "1"),
					resource.TestCheckResourceAttrSet("data.alicloud_nas_mount_targets.mt", "targets.0.mount_target_domain"),
					resource.TestCheckResourceAttr("data.alicloud_nas_mount_targets.mt", "targets.0.type", "Vpc"),
					resource.TestCheckResourceAttrSet("data.alicloud_nas_mount_targets.mt", "targets.0.vpc_id"),
					resource.TestCheckResourceAttrSet("data.alicloud_nas_mount_targets.mt", "targets.0.vswitch_id"),
					resource.TestCheckResourceAttr("data.alicloud_nas_mount_targets.mt", "targets.0.access_group_name", fmt.Sprintf("tf-testAccNasConfig-%d", rand)),
					resource.TestCheckResourceAttr("data.alicloud_nas_mount_targets.mt", "ids.#", "1"),
					resource.TestCheckResourceAttrSet("data.alicloud_nas_mount_targets.mt", "ids.0"),
				),
			},
			{
				Config: testAccCheckAlicloudMountTargetDataSourceTypeEmpty(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_nas_mount_targets.mt"),
					resource.TestCheckResourceAttr("data.alicloud_nas_mount_targets.mt", "targets.#", "0"),
					resource.TestCheckResourceAttr("data.alicloud_nas_mount_targets.mt", "ids.#", "0"),
				),
			},
		},
	})
}

func TestAccAlicloudNasMountTargetDataSourceMountTargetDomain(t *testing.T) {
	rand := acctest.RandIntRange(10000, 999999)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudMountTargetDataSourceMountTargetDomain(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_nas_mount_targets.mt"),
					resource.TestCheckResourceAttrSet("data.alicloud_nas_mount_targets.mt", "targets.0.mount_target_domain"),
					resource.TestCheckResourceAttr("data.alicloud_nas_mount_targets.mt", "targets.0.type", "Vpc"),
					resource.TestCheckResourceAttrSet("data.alicloud_nas_mount_targets.mt", "targets.0.vpc_id"),
					resource.TestCheckResourceAttrSet("data.alicloud_nas_mount_targets.mt", "targets.0.vswitch_id"),
					resource.TestCheckResourceAttr("data.alicloud_nas_mount_targets.mt", "targets.0.access_group_name", fmt.Sprintf("tf-testAccNasConfig-%d", rand)),
					resource.TestCheckResourceAttr("data.alicloud_nas_mount_targets.mt", "ids.#", "1"),
					resource.TestCheckResourceAttrSet("data.alicloud_nas_mount_targets.mt", "ids.0"),
				),
			},
			{
				Config: testAccCheckAlicloudMountTargetDataSourceMountTargetDomainEmpty(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_nas_mount_targets.mt"),
					resource.TestCheckResourceAttr("data.alicloud_nas_mount_targets.mt", "targets.#", "0"),
					resource.TestCheckResourceAttr("data.alicloud_nas_mount_targets.mt", "ids.#", "0"),
				),
			},
		},
	})
}

func TestAccAlicloudNasMountTargetDataSourceVpcId(t *testing.T) {
	rand := acctest.RandIntRange(10000, 999999)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudMountTargetDataSourceVpcId(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_nas_mount_targets.mt"),
					resource.TestCheckResourceAttrSet("data.alicloud_nas_mount_targets.mt", "targets.0.mount_target_domain"),
					resource.TestCheckResourceAttr("data.alicloud_nas_mount_targets.mt", "targets.0.type", "Vpc"),
					resource.TestCheckResourceAttrSet("data.alicloud_nas_mount_targets.mt", "targets.0.vpc_id"),
					resource.TestCheckResourceAttrSet("data.alicloud_nas_mount_targets.mt", "targets.0.vswitch_id"),
					resource.TestCheckResourceAttr("data.alicloud_nas_mount_targets.mt", "targets.0.access_group_name", fmt.Sprintf("tf-testAccNasConfig-%d", rand)),
					resource.TestCheckResourceAttr("data.alicloud_nas_mount_targets.mt", "ids.#", "1"),
					resource.TestCheckResourceAttrSet("data.alicloud_nas_mount_targets.mt", "ids.0"),
				),
			},
			{
				Config: testAccCheckAlicloudMountTargetDataSourceVpcIdEmpty(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_nas_mount_targets.mt"),
					resource.TestCheckResourceAttr("data.alicloud_nas_mount_targets.mt", "targets.#", "0"),
					resource.TestCheckResourceAttr("data.alicloud_nas_mount_targets.mt", "ids.#", "0"),
				),
			},
		},
	})
}

func TestAccAlicloudNasMountTargetDataSourceVSwitchId(t *testing.T) {
	rand := acctest.RandIntRange(10000, 999999)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudMountTargetDataSourceVSwitchId(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_nas_mount_targets.mt"),
					resource.TestCheckResourceAttrSet("data.alicloud_nas_mount_targets.mt", "targets.0.mount_target_domain"),
					resource.TestCheckResourceAttr("data.alicloud_nas_mount_targets.mt", "targets.0.type", "Vpc"),
					resource.TestCheckResourceAttrSet("data.alicloud_nas_mount_targets.mt", "targets.0.vpc_id"),
					resource.TestCheckResourceAttrSet("data.alicloud_nas_mount_targets.mt", "targets.0.vswitch_id"),
					resource.TestCheckResourceAttr("data.alicloud_nas_mount_targets.mt", "targets.0.access_group_name", fmt.Sprintf("tf-testAccNasConfig-%d", rand)),
					resource.TestCheckResourceAttr("data.alicloud_nas_mount_targets.mt", "ids.#", "1"),
					resource.TestCheckResourceAttrSet("data.alicloud_nas_mount_targets.mt", "ids.0"),
				),
			},
			{
				Config: testAccCheckAlicloudMountTargetDataSourceVSwitchIdEmpty(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_nas_mount_targets.mt"),
					resource.TestCheckResourceAttr("data.alicloud_nas_mount_targets.mt", "targets.#", "0"),
					resource.TestCheckResourceAttr("data.alicloud_nas_mount_targets.mt", "ids.#", "0"),
				),
			},
		},
	})
}

func TestAccAlicloudNasMountTargetDataSourceAll(t *testing.T) {
	rand := acctest.RandIntRange(10000, 999999)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudMountTargetDataSourceAll(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_nas_mount_targets.mt"),
					resource.TestCheckResourceAttrSet("data.alicloud_nas_mount_targets.mt", "targets.0.mount_target_domain"),
					resource.TestCheckResourceAttr("data.alicloud_nas_mount_targets.mt", "targets.0.type", "Vpc"),
					resource.TestCheckResourceAttrSet("data.alicloud_nas_mount_targets.mt", "targets.0.vpc_id"),
					resource.TestCheckResourceAttrSet("data.alicloud_nas_mount_targets.mt", "targets.0.vswitch_id"),
					resource.TestCheckResourceAttr("data.alicloud_nas_mount_targets.mt", "targets.0.access_group_name", fmt.Sprintf("tf-testAccNasConfig-%d", rand)),
					resource.TestCheckResourceAttr("data.alicloud_nas_mount_targets.mt", "ids.#", "1"),
					resource.TestCheckResourceAttrSet("data.alicloud_nas_mount_targets.mt", "ids.0"),
				),
			},
			{
				Config: testAccCheckAlicloudMountTargetDataSourceAllEmpty(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_nas_mount_targets.mt"),
					resource.TestCheckResourceAttr("data.alicloud_nas_mount_targets.mt", "targets.#", "0"),
					resource.TestCheckResourceAttr("data.alicloud_nas_mount_targets.mt", "ids.#", "0"),
				),
			},
		},
	})
}

func testAccCheckAlicloudMountTargetDataSourceAccessGroupName(rand int) string {
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
	resource "alicloud_nas_mount_target" "foo" {
			file_system_id = "${alicloud_nas_file_system.foo.id}"
			access_group_name = "${alicloud_nas_access_group.foo.id}"
	}
	data "alicloud_nas_mount_targets" "mt" {
			file_system_id = "${alicloud_nas_mount_target.foo.file_system_id}"
			access_group_name = "${alicloud_nas_access_group.foo.id}"
	}`, rand)
}

func testAccCheckAlicloudMountTargetDataSourceAccessGroupNameEmpty(rand int) string {
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
	resource "alicloud_nas_mount_target" "foo" {
			file_system_id = "${alicloud_nas_file_system.foo.id}"
			access_group_name = "${alicloud_nas_access_group.foo.id}"
	}
	data "alicloud_nas_mount_targets" "mt" {
			file_system_id = "${alicloud_nas_mount_target.foo.file_system_id}"
			access_group_name = "${alicloud_nas_access_group.foo.id}-fake"
	}`, rand)
}

func testAccCheckAlicloudMountTargetDataSourceType(rand int) string {
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
	resource "alicloud_nas_mount_target" "foo" {
			file_system_id = "${alicloud_nas_file_system.foo.id}"
			access_group_name = "${alicloud_nas_access_group.foo.id}"
			vswitch_id = "${alicloud_vswitch.foo.id}"               
	}
	data "alicloud_nas_mount_targets" "mt" {
			file_system_id = "${alicloud_nas_mount_target.foo.file_system_id}"
			type = "Vpc"	
	}`, rand)
}

func testAccCheckAlicloudMountTargetDataSourceTypeEmpty(rand int) string {
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
	resource "alicloud_nas_mount_target" "foo" {
			file_system_id = "${alicloud_nas_file_system.foo.id}"
			access_group_name = "${alicloud_nas_access_group.foo.id}"
			vswitch_id = "${alicloud_vswitch.foo.id}"
	}
	data "alicloud_nas_mount_targets" "mt" {
			file_system_id = "${alicloud_nas_mount_target.foo.file_system_id}"
			type = "Classic"
	}`, rand)
}

func testAccCheckAlicloudMountTargetDataSourceMountTargetDomain(rand int) string {
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
	resource "alicloud_nas_mount_target" "foo" {
			file_system_id = "${alicloud_nas_file_system.foo.id}"
			access_group_name = "${alicloud_nas_access_group.foo.id}"
			vswitch_id = "${alicloud_vswitch.foo.id}"               
	}
	data "alicloud_nas_mount_targets" "mt" {
                        file_system_id = "${alicloud_nas_mount_target.foo.file_system_id}"
			mount_target_domain = "${alicloud_nas_mount_target.foo.id}"
	}`, rand)
}

func testAccCheckAlicloudMountTargetDataSourceMountTargetDomainEmpty(rand int) string {
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
	resource "alicloud_nas_mount_target" "foo" {
			file_system_id = "${alicloud_nas_file_system.foo.id}"
			access_group_name = "${alicloud_nas_access_group.foo.id}"
			vswitch_id = "${alicloud_vswitch.foo.id}"
	}
	data "alicloud_nas_mount_targets" "mt" {
			file_system_id = "${alicloud_nas_mount_target.foo.file_system_id}"
			mount_target_domain = "${alicloud_nas_mount_target.foo.id}-fake"
	}`, rand)
}

func testAccCheckAlicloudMountTargetDataSourceVpcId(rand int) string {
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
	resource "alicloud_nas_mount_target" "foo" {
			file_system_id = "${alicloud_nas_file_system.foo.id}"
			access_group_name = "${alicloud_nas_access_group.foo.id}"
			vswitch_id = "${alicloud_vswitch.foo.id}"               
	}
	data "alicloud_nas_mount_targets" "mt" {
			file_system_id = "${alicloud_nas_mount_target.foo.file_system_id}"
			vpc_id = "${alicloud_vpc.foo.id}"
	}`, rand)
}

func testAccCheckAlicloudMountTargetDataSourceVpcIdEmpty(rand int) string {
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
	resource "alicloud_nas_mount_target" "foo" {
			file_system_id = "${alicloud_nas_file_system.foo.id}"
			access_group_name = "${alicloud_nas_access_group.foo.id}"
			vswitch_id = "${alicloud_vswitch.foo.id}"
	}
	data "alicloud_nas_mount_targets" "mt" {
			file_system_id = "${alicloud_nas_mount_target.foo.file_system_id}"
			vpc_id = "vpc_123445"
	}`, rand)
}

func testAccCheckAlicloudMountTargetDataSourceVSwitchId(rand int) string {
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
	resource "alicloud_nas_mount_target" "foo" {
			file_system_id = "${alicloud_nas_file_system.foo.id}"
			access_group_name = "${alicloud_nas_access_group.foo.id}"
			vswitch_id = "${alicloud_vswitch.foo.id}"               
	}
	data "alicloud_nas_mount_targets" "mt" {
			file_system_id = "${alicloud_nas_mount_target.foo.file_system_id}"
			vswitch_id = "${alicloud_nas_mount_target.foo.vswitch_id}"
	}`, rand)
}

func testAccCheckAlicloudMountTargetDataSourceVSwitchIdEmpty(rand int) string {
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
	resource "alicloud_nas_mount_target" "foo" {
			file_system_id = "${alicloud_nas_file_system.foo.id}"
			access_group_name = "${alicloud_nas_access_group.foo.id}"
			vswitch_id = "${alicloud_vswitch.foo.id}"
	}
	data "alicloud_nas_mount_targets" "mt" {
			file_system_id = "${alicloud_nas_mount_target.foo.file_system_id}"
			vswitch_id = "vsw-123456"
	}`, rand)
}

func testAccCheckAlicloudMountTargetDataSourceAll(rand int) string {
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
        resource "alicloud_nas_mount_target" "foo" {
                        file_system_id = "${alicloud_nas_file_system.foo.id}"
                        access_group_name = "${alicloud_nas_access_group.foo.id}"
                        vswitch_id = "${alicloud_vswitch.foo.id}"               
        }
        data "alicloud_nas_mount_targets" "mt" {
                        file_system_id = "${alicloud_nas_mount_target.foo.file_system_id}"
                        vswitch_id = "${alicloud_nas_mount_target.foo.vswitch_id}"
			access_group_name = "${alicloud_nas_mount_target.foo.access_group_name}"
			type = "Vpc"
			vpc_id = "${alicloud_vpc.foo.id}"
			mount_target_domain = "${alicloud_nas_mount_target.foo.id}"
        }`, rand)
}

func testAccCheckAlicloudMountTargetDataSourceAllEmpty(rand int) string {
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
        resource "alicloud_nas_mount_target" "foo" {
                        file_system_id = "${alicloud_nas_file_system.foo.id}"
                        access_group_name = "${alicloud_nas_access_group.foo.id}"
                        vswitch_id = "${alicloud_vswitch.foo.id}"               
        }
        data "alicloud_nas_mount_targets" "mt" {
                        file_system_id = "${alicloud_nas_mount_target.foo.file_system_id}"
                        vswitch_id = "vsw-123456"
			vpc_id = "vpc_123445"
                        access_group_name = "tf-testAccNasConfig"
                        type = "Classic"
                        mount_target_domain = "${alicloud_nas_mount_target.foo.id}"
        }`, rand)
}

package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/nas"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func TestAccAlicloudNas_MountTarget_update(t *testing.T) {
	var v nas.DescribeMountTargetsMountTarget1
	rand1 := acctest.RandIntRange(10000, 499999)
	rand2 := acctest.RandIntRange(500000, 999999)
	resourceID := "alicloud_nas_mount_target.default"
	ra := resourceAttrInit(resourceID, map[string]string{})
	serviceFunc := func() interface{} {
		return &NasService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceID, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, false, connectivity.NasNoSupportedRegions)
		},
		IDRefreshName: resourceID,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckMountTargetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNasMountTargetVpcConfig(rand1, rand2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"vswitch_id":        CHECKSET,
						"file_system_id":    CHECKSET,
						"status":            "Active",
						"access_group_name": fmt.Sprintf("tf-testAccNasConfig-%d", rand1),
					}),
				),
			},
			{
				ResourceName:      resourceID,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccNasMountTargetConfigUpdateAccessGroup(rand1, rand2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"vswitch_id":        CHECKSET,
						"file_system_id":    CHECKSET,
						"status":            "Active",
						"access_group_name": fmt.Sprintf("tf-testAccNasConfig-2-%d", rand2),
					}),
				),
			},
			{
				Config: testAccNasMountTargetConfigUpdateStatus(rand1, rand2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"vswitch_id":        CHECKSET,
						"file_system_id":    CHECKSET,
						"status":            "Inactive",
						"access_group_name": fmt.Sprintf("tf-testAccNasConfig-2-%d", rand2),
					}),
				),
			},
			{
				Config: testAccNasMountTargetConfigUpdateAll(rand1, rand2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"vswitch_id":        CHECKSET,
						"file_system_id":    CHECKSET,
						"status":            "Active",
						"access_group_name": fmt.Sprintf("tf-testAccNasConfig-%d", rand1),
					}),
				),
			},
		},
	})
}

func TestAccAlicloudNas_MountTarget_updateT(t *testing.T) {
	var v nas.DescribeMountTargetsMountTarget1
	rand1 := acctest.RandIntRange(10000, 499999)
	rand2 := acctest.RandIntRange(500000, 999999)
	resourceID := "alicloud_nas_mount_target.default"
	ra := resourceAttrInit(resourceID, map[string]string{})
	serviceFunc := func() interface{} {
		return &NasService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceID, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceID,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckMountTargetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNasMountTargetVpcConfigT(rand1, rand2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"vswitch_id":        CHECKSET,
						"file_system_id":    CHECKSET,
						"status":            "Active",
						"access_group_name": fmt.Sprintf("tf-testAccNasConfig-%d", rand1),
					}),
				),
			},
			{
				ResourceName:      resourceID,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccNasMountTargetConfigUpdateAccessGroupT(rand1, rand2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"vswitch_id":        CHECKSET,
						"file_system_id":    CHECKSET,
						"status":            "Active",
						"access_group_name": fmt.Sprintf("tf-testAccNasConfig-2-%d", rand2),
					}),
				),
			},
			{
				Config: testAccNasMountTargetConfigUpdateStatusT(rand1, rand2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"vswitch_id":        CHECKSET,
						"file_system_id":    CHECKSET,
						"status":            "Inactive",
						"access_group_name": fmt.Sprintf("tf-testAccNasConfig-2-%d", rand2),
					}),
				),
			},
			{
				Config: testAccNasMountTargetConfigUpdateAllT(rand1, rand2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"vswitch_id":        CHECKSET,
						"file_system_id":    CHECKSET,
						"status":            "Active",
						"access_group_name": fmt.Sprintf("tf-testAccNasConfig-%d", rand1),
					}),
				),
			},
		},
	})
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

func testAccNasMountTargetVpcConfig(rand1 int, rand2 int) string {
	return fmt.Sprintf(`
	variable "name" {
		default = "tf-testAccVswitch"
	}
	data "alicloud_zones" "default" {
		available_resource_creation = "VSwitch"
	}
	resource "alicloud_vpc" "default" {
		name = "${var.name}"
		cidr_block = "172.16.0.0/12"
	}
	resource "alicloud_vswitch" "default" {
		vpc_id = "${alicloud_vpc.default.id}"
		cidr_block = "172.16.0.0/24"
		availability_zone = "${data.alicloud_zones.default.zones.0.id}"
		name = "${var.name}-1"
	}
	variable "storage_type" {
  		default = "Performance"
	}
	data "alicloud_nas_protocols" "default" {
        	type = "${var.storage_type}"
	}
	resource "alicloud_nas_file_system" "default" {
	        protocol_type = "${data.alicloud_nas_protocols.default.protocols.0}"
        	storage_type = "${var.storage_type}"
	        description = "tf-testAccNasConfigUpdateName"
	}
	resource "alicloud_nas_access_group" "default" {
                name = "tf-testAccNasConfig-%d"
              	type = "Vpc"
	        description = "tf-testAccNasConfig"
	}
	resource "alicloud_nas_access_group" "bar" {
        	        name = "tf-testAccNasConfig-2-%d"
                	type = "Vpc"
	                description = "tf-testAccNasConfig-2"
	}
	resource "alicloud_nas_mount_target" "default" {
		file_system_id = "${alicloud_nas_file_system.default.id}"
		access_group_name = "${alicloud_nas_access_group.default.id}"
		vswitch_id = "${alicloud_vswitch.default.id}"               
	}
`, rand1, rand2)
}

func testAccNasMountTargetConfigUpdateAccessGroup(rand1 int, rand2 int) string {
	return fmt.Sprintf(`
	variable "name" {
                default = "tf-testAccVswitch"
	}
	data "alicloud_zones" "default" {
		available_resource_creation = "VSwitch"
	}
	resource "alicloud_vpc" "default" {
                name = "${var.name}"
               	cidr_block = "172.16.0.0/12"
	}
	resource "alicloud_vswitch" "default" {
                vpc_id = "${alicloud_vpc.default.id}"
               	cidr_block = "172.16.0.0/24"
	        availability_zone = "${data.alicloud_zones.default.zones.0.id}"
                name = "${var.name}-1"
	}
	variable "storage_type" {
  		default = "Performance"
	}
	data "alicloud_nas_protocols" "default" {
        	type = "${var.storage_type}"
	}
	resource "alicloud_nas_file_system" "default" {
		protocol_type = "${data.alicloud_nas_protocols.default.protocols.0}"
	        storage_type = "${var.storage_type}"
	        description = "tf-testAccNasConfigUpdateName"
	}
	resource "alicloud_nas_access_group" "default" {
                name = "tf-testAccNasConfig-%d"
               	type = "Vpc"
	        description = "tf-testAccNasConfig"
	}
	resource "alicloud_nas_access_group" "bar" {
        	        name = "tf-testAccNasConfig-2-%d"
                	type = "Vpc"
	                description = "tf-testAccNasConfig-2"
	}
	resource "alicloud_nas_mount_target" "default" {
		file_system_id = "${alicloud_nas_file_system.default.id}"
		access_group_name = "${alicloud_nas_access_group.bar.id}"
		vswitch_id = "${alicloud_vswitch.default.id}"
	}`, rand1, rand2)
}

func testAccNasMountTargetConfigUpdateStatus(rand1 int, rand2 int) string {
	return fmt.Sprintf(`
	variable "name" {
                default = "tf-testAccVswitch"
	}
	data "alicloud_zones" "default" {
		available_resource_creation = "VSwitch"
	}
	resource "alicloud_vpc" "default" {
                name = "${var.name}"
               	cidr_block = "172.16.0.0/12"
	}
	resource "alicloud_vswitch" "default" {
                vpc_id = "${alicloud_vpc.default.id}"
               	cidr_block = "172.16.0.0/24"
	        availability_zone = "${data.alicloud_zones.default.zones.0.id}"
                name = "${var.name}-1"
	}
	variable "storage_type" {
  		default = "Performance"
	}
	data "alicloud_nas_protocols" "default" {
        	type = "${var.storage_type}"
	}
	resource "alicloud_nas_file_system" "default" {
	        protocol_type = "${data.alicloud_nas_protocols.default.protocols.0}"
        	storage_type = "${var.storage_type}"
	        description = "tf-testAccNasConfigUpdateName"
	}
	resource "alicloud_nas_access_group" "default" {
                name = "tf-testAccNasConfig-%d"
               	type = "Vpc"
	        description = "tf-testAccNasConfig"
	}
	resource "alicloud_nas_access_group" "bar" {
                name = "tf-testAccNasConfig-2-%d"
               	type = "Vpc"
	        description = "tf-testAccNasConfig-2"
	}
	resource "alicloud_nas_mount_target" "default" {
		file_system_id = "${alicloud_nas_file_system.default.id}"
               	access_group_name = "${alicloud_nas_access_group.bar.id}"
	        status = "Inactive"
		vswitch_id = "${alicloud_vswitch.default.id}"
	}`, rand1, rand2)
}

func testAccNasMountTargetConfigUpdateAll(rand1 int, rand2 int) string {
	return fmt.Sprintf(`
	variable "name" {
                default = "tf-testAccVswitch"
	}
	data "alicloud_zones" "default" {
		available_resource_creation = "VSwitch"
	}
	resource "alicloud_vpc" "default" {
                name = "${var.name}"
               	cidr_block = "172.16.0.0/12"
	}
	resource "alicloud_vswitch" "default" {
                vpc_id = "${alicloud_vpc.default.id}"
               	cidr_block = "172.16.0.0/24"
	        availability_zone = "${data.alicloud_zones.default.zones.0.id}"
                name = "${var.name}-1"
	}
	variable "storage_type" {
  		default = "Performance"
	}
	data "alicloud_nas_protocols" "default" {
        	type = "${var.storage_type}"
	}
	resource "alicloud_nas_file_system" "default" {
	        protocol_type = "${data.alicloud_nas_protocols.default.protocols.0}"
        	storage_type = "${var.storage_type}"
	        description = "tf-testAccNasConfigUpdateName"
	}
	resource "alicloud_nas_access_group" "default" {
                name = "tf-testAccNasConfig-%d"
               	type = "Vpc"
	        description = "tf-testAccNasConfig"
	}
	resource "alicloud_nas_access_group" "bar" {
                name = "tf-testAccNasConfig-2-%d"
	        type = "Vpc"
                description = "tf-testAccNasConfig-2"
	}
	resource "alicloud_nas_mount_target" "default" {
                file_system_id = "${alicloud_nas_file_system.default.id}"
               	access_group_name = "${alicloud_nas_access_group.default.id}"
	        status = "Active"
		vswitch_id = "${alicloud_vswitch.default.id}"
	}`, rand1, rand2)
}

func testAccNasMountTargetVpcConfigT(rand1 int, rand2 int) string {
	return fmt.Sprintf(`
	variable "name" {
		default = "tf-testAccVswitch"
	}
	data "alicloud_zones" "default" {
		available_resource_creation = "VSwitch"
	}
	resource "alicloud_vpc" "default" {
		name = "${var.name}"
		cidr_block = "172.16.0.0/12"
	}
	resource "alicloud_vswitch" "default" {
		vpc_id = "${alicloud_vpc.default.id}"
		cidr_block = "172.16.0.0/24"
		availability_zone = "${data.alicloud_zones.default.zones.0.id}"
		name = "${var.name}-1"
	}
	variable "storage_type" {
  		default = "Capacity"
	}
	data "alicloud_nas_protocols" "default" {
	        type = "${var.storage_type}"
	}
	resource "alicloud_nas_file_system" "default" {
        	protocol_type = "${data.alicloud_nas_protocols.default.protocols.0}"
	        storage_type = "${var.storage_type}"
        	description = "tf-testAccNasConfigUpdateName"
	}
	resource "alicloud_nas_access_group" "default" {
        	name = "tf-testAccNasConfig-%d"
                type = "Vpc"
	        description = "tf-testAccNasConfig"
	}
	resource "alicloud_nas_access_group" "bar" {
		name = "tf-testAccNasConfig-2-%d"
		type = "Vpc"
		description = "tf-testAccNasConfig-2"
	}
	resource "alicloud_nas_mount_target" "default" {
		file_system_id = "${alicloud_nas_file_system.default.id}"
		access_group_name = "${alicloud_nas_access_group.default.id}"
		vswitch_id = "${alicloud_vswitch.default.id}"               
	}
`, rand1, rand2)
}

func testAccNasMountTargetConfigUpdateAccessGroupT(rand1 int, rand2 int) string {
	return fmt.Sprintf(`
	variable "name" {
        	default = "tf-testAccVswitch"
	}
	data "alicloud_zones" "default" {
		available_resource_creation = "VSwitch"
	}
	resource "alicloud_vpc" "default" {
                name = "${var.name}"
               	cidr_block = "172.16.0.0/12"
	}
	resource "alicloud_vswitch" "default" {
                vpc_id = "${alicloud_vpc.default.id}"
               	cidr_block = "172.16.0.0/24"
	        availability_zone = "${data.alicloud_zones.default.zones.0.id}"
                name = "${var.name}-1"
	}
	variable "storage_type" {
  		default = "Capacity"
	}
	data "alicloud_nas_protocols" "default" {
        	type = "${var.storage_type}"
	}
	resource "alicloud_nas_file_system" "default" {
	        protocol_type = "${data.alicloud_nas_protocols.default.protocols.0}"
	        storage_type = "${var.storage_type}"
	        description = "tf-testAccNasConfigUpdateName"
	}
	resource "alicloud_nas_access_group" "default" {
                name = "tf-testAccNasConfig-%d"
               	type = "Vpc"
	        description = "tf-testAccNasConfig"
	}
	resource "alicloud_nas_access_group" "bar" {
                name = "tf-testAccNasConfig-2-%d"
               	type = "Vpc"
	        description = "tf-testAccNasConfig-2"
	}
	resource "alicloud_nas_mount_target" "default" {
		file_system_id = "${alicloud_nas_file_system.default.id}"
		access_group_name = "${alicloud_nas_access_group.bar.id}"
		vswitch_id = "${alicloud_vswitch.default.id}"
	}`, rand1, rand2)
}

func testAccNasMountTargetConfigUpdateStatusT(rand1 int, rand2 int) string {
	return fmt.Sprintf(`
	variable "name" {
                default = "tf-testAccVswitch"
	}
	data "alicloud_zones" "default" {
		available_resource_creation = "VSwitch"
	}
	resource "alicloud_vpc" "default" {
                name = "${var.name}"
               	cidr_block = "172.16.0.0/12"
	}
	resource "alicloud_vswitch" "default" {
                vpc_id = "${alicloud_vpc.default.id}"
               	cidr_block = "172.16.0.0/24"
	        availability_zone = "${data.alicloud_zones.default.zones.0.id}"
                name = "${var.name}-1"
	}
	variable "storage_type" {
  		default = "Capacity"
	}
	data "alicloud_nas_protocols" "default" {
        	type = "${var.storage_type}"
	}
	resource "alicloud_nas_file_system" "default" {
	        protocol_type = "${data.alicloud_nas_protocols.default.protocols.0}"
        	storage_type = "${var.storage_type}"
	        description = "tf-testAccNasConfigUpdateName"
	}
	resource "alicloud_nas_access_group" "default" {
                name = "tf-testAccNasConfig-%d"
               	type = "Vpc"
	        description = "tf-testAccNasConfig"
	}
	resource "alicloud_nas_access_group" "bar" {
                name = "tf-testAccNasConfig-2-%d"
               	type = "Vpc"
	        description = "tf-testAccNasConfig-2"
	}
	resource "alicloud_nas_mount_target" "default" {
		file_system_id = "${alicloud_nas_file_system.default.id}"
               	access_group_name = "${alicloud_nas_access_group.bar.id}"
	        status = "Inactive"
		vswitch_id = "${alicloud_vswitch.default.id}"
	}`, rand1, rand2)
}

func testAccNasMountTargetConfigUpdateAllT(rand1 int, rand2 int) string {
	return fmt.Sprintf(`
	variable "name" {
                default = "tf-testAccVswitch"
	}
	data "alicloud_zones" "default" {
		available_resource_creation = "VSwitch"
	}
	resource "alicloud_vpc" "default" {
                name = "${var.name}"
               	cidr_block = "172.16.0.0/12"
	}
	resource "alicloud_vswitch" "default" {
                vpc_id = "${alicloud_vpc.default.id}"
               	cidr_block = "172.16.0.0/24"
	        availability_zone = "${data.alicloud_zones.default.zones.0.id}"
                name = "${var.name}-1"
	}
	variable "storage_type" {
  		default = "Capacity"
	}
	data "alicloud_nas_protocols" "default" {
        	type = "${var.storage_type}"
	}
	resource "alicloud_nas_file_system" "default" {
	        protocol_type = "${data.alicloud_nas_protocols.default.protocols.0}"
        	storage_type = "${var.storage_type}"
	        description = "tf-testAccNasConfigUpdateName"
	}
	resource "alicloud_nas_access_group" "default" {
                name = "tf-testAccNasConfig-%d"
	        type = "Vpc"
                description = "tf-testAccNasConfig"
	}
	resource "alicloud_nas_access_group" "bar" {
                name = "tf-testAccNasConfig-2-%d"
	        type = "Vpc"
                description = "tf-testAccNasConfig-2"
	}
	resource "alicloud_nas_mount_target" "default" {
                file_system_id = "${alicloud_nas_file_system.default.id}"
               	access_group_name = "${alicloud_nas_access_group.default.id}"
	        status = "Active"
		vswitch_id = "${alicloud_vswitch.default.id}"
	}`, rand1, rand2)
}

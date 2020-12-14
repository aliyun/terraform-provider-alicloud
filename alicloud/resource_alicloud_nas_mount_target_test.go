package alicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/nas"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccAlicloudNas_MountTarget_update(t *testing.T) {
	var v nas.MountTarget
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
			testAccPreCheckWithNoDefaultVpc(t)
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
				ResourceName:            resourceID,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"security_group_id"},
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
	var v nas.MountTarget
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
			testAccPreCheckWithNoDefaultVpc(t)
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
				ResourceName:            resourceID,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"security_group_id"},
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
	variable "storage_type" {
  		default = "Performance"
	}
	data "alicloud_nas_protocols" "default" {
        	type = "${var.storage_type}"
	}
	data "alicloud_vpcs" "default" {
			is_default = true
	}
	resource "alicloud_security_group" "default" {
		  name = var.name
		  vpc_id = "${data.alicloud_vpcs.default.vpcs.0.id}"
	}
	resource "alicloud_nas_file_system" "default" {
	        protocol_type = "${data.alicloud_nas_protocols.default.protocols.0}"
        	storage_type = "${var.storage_type}"
	        description = "tf-testAccNasConfigUpdateName"
	}
	resource "alicloud_nas_access_group" "default" {
                access_group_name = "tf-testAccNasConfig-%d"
              	access_group_type = "Vpc"
	        	description = "tf-testAccNasConfig"
	}
	resource "alicloud_nas_access_group" "bar" {
				access_group_name = "tf-testAccNasConfig-2-%d"
				access_group_type = "Vpc"
				description = "tf-testAccNasConfig-2"
	}
	resource "alicloud_nas_mount_target" "default" {
		file_system_id = "${alicloud_nas_file_system.default.id}"
		access_group_name = "${alicloud_nas_access_group.default.access_group_name}"
		vswitch_id = "${data.alicloud_vpcs.default.vpcs.0.vswitch_ids.0}"
  		security_group_id = "${alicloud_security_group.default.id}"
	}
`, rand1, rand2)
}

func testAccNasMountTargetConfigUpdateAccessGroup(rand1 int, rand2 int) string {
	return fmt.Sprintf(`
	variable "name" {
                default = "tf-testAccVswitch"
	}

	data "alicloud_vpcs" "default" {
			is_default = true
	}
	resource "alicloud_security_group" "default" {
		  name = var.name
		  vpc_id = "${data.alicloud_vpcs.default.vpcs.0.id}"
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
                access_group_name = "tf-testAccNasConfig-%d"
               	access_group_type = "Vpc"
	        	description = "tf-testAccNasConfig"
	}
	resource "alicloud_nas_access_group" "bar" {
        	        access_group_name = "tf-testAccNasConfig-2-%d"
                	access_group_type = "Vpc"
	                description = "tf-testAccNasConfig-2"
	}
	resource "alicloud_nas_mount_target" "default" {
		file_system_id = "${alicloud_nas_file_system.default.id}"
		access_group_name = "${alicloud_nas_access_group.bar.access_group_name}"
		vswitch_id = "${data.alicloud_vpcs.default.vpcs.0.vswitch_ids.0}"
  		security_group_id = "${alicloud_security_group.default.id}"
	}`, rand1, rand2)
}

func testAccNasMountTargetConfigUpdateStatus(rand1 int, rand2 int) string {
	return fmt.Sprintf(`
	variable "name" {
                default = "tf-testAccVswitch"
	}
		data "alicloud_vpcs" "default" {
			is_default = true
	}
	resource "alicloud_security_group" "default" {
		  name = var.name
		  vpc_id = "${data.alicloud_vpcs.default.vpcs.0.id}"
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
                access_group_name = "tf-testAccNasConfig-%d"
               	access_group_type = "Vpc"
	        	description = "tf-testAccNasConfig"
	}
	resource "alicloud_nas_access_group" "bar" {
                access_group_name = "tf-testAccNasConfig-2-%d"
               	access_group_type = "Vpc"
	        	description = "tf-testAccNasConfig-2"
	}
	resource "alicloud_nas_mount_target" "default" {
			file_system_id = "${alicloud_nas_file_system.default.id}"
			access_group_name = "${alicloud_nas_access_group.bar.access_group_name}"
	        status = "Inactive"
			vswitch_id = "${data.alicloud_vpcs.default.vpcs.0.vswitch_ids.0}"
  			security_group_id = "${alicloud_security_group.default.id}"
	}`, rand1, rand2)
}

func testAccNasMountTargetConfigUpdateAll(rand1 int, rand2 int) string {
	return fmt.Sprintf(`
	variable "name" {
                default = "tf-testAccVswitch"
	}
		data "alicloud_vpcs" "default" {
			is_default = true
	}
	resource "alicloud_security_group" "default" {
		  name = var.name
		  vpc_id = "${data.alicloud_vpcs.default.vpcs.0.id}"
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
                access_group_name = "tf-testAccNasConfig-%d"
               	access_group_type = "Vpc"
	        	description = "tf-testAccNasConfig"
	}
	resource "alicloud_nas_access_group" "bar" {
                access_group_name = "tf-testAccNasConfig-2-%d"
	        	access_group_type = "Vpc"
                description = "tf-testAccNasConfig-2"
	}
	resource "alicloud_nas_mount_target" "default" {
			file_system_id = "${alicloud_nas_file_system.default.id}"
			access_group_name = "${alicloud_nas_access_group.default.access_group_name}"
	        status = "Active"
			vswitch_id = "${data.alicloud_vpcs.default.vpcs.0.vswitch_ids.0}"
  			security_group_id = "${alicloud_security_group.default.id}"
	}`, rand1, rand2)
}

func testAccNasMountTargetVpcConfigT(rand1 int, rand2 int) string {
	return fmt.Sprintf(`
	variable "name" {
		default = "tf-testAccVswitch"
	}
		data "alicloud_vpcs" "default" {
			is_default = true
	}
	resource "alicloud_security_group" "default" {
		  name = var.name
		  vpc_id = "${data.alicloud_vpcs.default.vpcs.0.id}"
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
        	access_group_name = "tf-testAccNasConfig-%d"
			access_group_type = "Vpc"
	        description = "tf-testAccNasConfig"
	}
	resource "alicloud_nas_access_group" "bar" {
		access_group_name = "tf-testAccNasConfig-2-%d"
		access_group_type = "Vpc"
		description = "tf-testAccNasConfig-2"
	}
	resource "alicloud_nas_mount_target" "default" {
		file_system_id = "${alicloud_nas_file_system.default.id}"
		access_group_name = "${alicloud_nas_access_group.default.access_group_name}"
			vswitch_id = "${data.alicloud_vpcs.default.vpcs.0.vswitch_ids.0}"
  			security_group_id = "${alicloud_security_group.default.id}"              
	}
`, rand1, rand2)
}

func testAccNasMountTargetConfigUpdateAccessGroupT(rand1 int, rand2 int) string {
	return fmt.Sprintf(`
	variable "name" {
        	default = "tf-testAccVswitch"
	}
		data "alicloud_vpcs" "default" {
			is_default = true
	}
	resource "alicloud_security_group" "default" {
		  name = var.name
		  vpc_id = "${data.alicloud_vpcs.default.vpcs.0.id}"
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
                access_group_name = "tf-testAccNasConfig-%d"
               	access_group_type = "Vpc"
	        	description = "tf-testAccNasConfig"
	}
	resource "alicloud_nas_access_group" "bar" {
                access_group_name = "tf-testAccNasConfig-2-%d"
               	access_group_type = "Vpc"
	        	description = "tf-testAccNasConfig-2"
	}
	resource "alicloud_nas_mount_target" "default" {
		file_system_id = "${alicloud_nas_file_system.default.id}"
		access_group_name = "${alicloud_nas_access_group.bar.access_group_name}"
			vswitch_id = "${data.alicloud_vpcs.default.vpcs.0.vswitch_ids.0}"
  			security_group_id = "${alicloud_security_group.default.id}"
	}`, rand1, rand2)
}

func testAccNasMountTargetConfigUpdateStatusT(rand1 int, rand2 int) string {
	return fmt.Sprintf(`
	variable "name" {
                default = "tf-testAccVswitch"
	}
		data "alicloud_vpcs" "default" {
			is_default = true
	}
	resource "alicloud_security_group" "default" {
		  name = var.name
		  vpc_id = "${data.alicloud_vpcs.default.vpcs.0.id}"
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
                access_group_name = "tf-testAccNasConfig-%d"
               	access_group_type = "Vpc"
	        	description = "tf-testAccNasConfig"
	}
	resource "alicloud_nas_access_group" "bar" {
                access_group_name = "tf-testAccNasConfig-2-%d"
               	access_group_type = "Vpc"
	        	description = "tf-testAccNasConfig-2"
	}
	resource "alicloud_nas_mount_target" "default" {
		file_system_id = "${alicloud_nas_file_system.default.id}"
               	access_group_name = "${alicloud_nas_access_group.bar.access_group_name}"
	        status = "Inactive"
			vswitch_id = "${data.alicloud_vpcs.default.vpcs.0.vswitch_ids.0}"
  			security_group_id = "${alicloud_security_group.default.id}"
	}`, rand1, rand2)
}

func testAccNasMountTargetConfigUpdateAllT(rand1 int, rand2 int) string {
	return fmt.Sprintf(`
	variable "name" {
                default = "tf-testAccVswitch"
	}
		data "alicloud_vpcs" "default" {
			is_default = true
	}
	resource "alicloud_security_group" "default" {
		  name = var.name
		  vpc_id = "${data.alicloud_vpcs.default.vpcs.0.id}"
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
                access_group_name = "tf-testAccNasConfig-%d"
	        	access_group_type = "Vpc"
                description = "tf-testAccNasConfig"
	}
	resource "alicloud_nas_access_group" "bar" {
                access_group_name = "tf-testAccNasConfig-2-%d"
	        	access_group_type = "Vpc"
                description = "tf-testAccNasConfig-2"
	}
	resource "alicloud_nas_mount_target" "default" {
                file_system_id = "${alicloud_nas_file_system.default.id}"
               	access_group_name = "${alicloud_nas_access_group.default.access_group_name}"
	        status = "Active"
			vswitch_id = "${data.alicloud_vpcs.default.vpcs.0.vswitch_ids.0}"
  			security_group_id = "${alicloud_security_group.default.id}"
	}`, rand1, rand2)
}

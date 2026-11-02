package nas_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/nas"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAliCloudNASMountTarget_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_nas_mount_target.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudNasMountTarget0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return nas.NewNasService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeNasMountTarget")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%sAlicloudNasMountTarget%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudNasMountTargetBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},

		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"access_group_name": "${alicloud_nas_access_group.example.access_group_name}",
					"file_system_id":    "${alicloud_nas_file_system.example.id}",
					"vswitch_id":        "${alicloud_vswitch.main.id}",
					"security_group_id": "${alicloud_security_group.example.id}",
					"status":            "Active",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"access_group_name": name,
						"file_system_id":    acctest.CHECKSET,
						"vswitch_id":        acctest.CHECKSET,
						"security_group_id": acctest.CHECKSET,
						"status":            "Active",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"security_group_id"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"access_group_name": "${alicloud_nas_access_group.example1.access_group_name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"access_group_name": name + "change",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "Inactive",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "Inactive",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"access_group_name": "${alicloud_nas_access_group.example.access_group_name}",
					"file_system_id":    "${alicloud_nas_file_system.example.id}",
					"vswitch_id":        "${alicloud_vswitch.main.id}",
					"security_group_id": "${alicloud_security_group.example.id}",
					"status":            "Active",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"access_group_name": name,
						"file_system_id":    acctest.CHECKSET,
						"vswitch_id":        acctest.CHECKSET,
						"security_group_id": acctest.CHECKSET,
						"status":            "Active",
					}),
				),
			},
		},
	})
}

func TestAccAliCloudNasExtremeMountTarget_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_nas_mount_target.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudNasMountTarget1)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return nas.NewNasService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeNasMountTarget")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%sAlicloudNasMountTarget%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudNasMountTargetBasicDependence1)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},

		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"access_group_name": "${alicloud_nas_access_group.example.access_group_name}",
					"file_system_id":    "${alicloud_nas_file_system.example.id}",
					"vswitch_id":        "${alicloud_vswitch.main.id}",
					"vpc_id":            "${alicloud_vpc.main.id}",
					"network_type":      "${alicloud_nas_access_group.example.access_group_type}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"access_group_name":   name,
						"file_system_id":      acctest.CHECKSET,
						"vswitch_id":          acctest.CHECKSET,
						"vpc_id":              acctest.CHECKSET,
						"mount_target_domain": acctest.CHECKSET,
						"status":              acctest.CHECKSET,
						"network_type":        "Vpc",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"access_group_name": "${alicloud_nas_access_group.example1.access_group_name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"access_group_name": name + "change",
					}),
				),
			},
		},
	})
}

var AlicloudNasMountTarget0 = map[string]string{}

func AlicloudNasMountTargetBasicDependence0(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}

variable "name1" {
	default = "%schange"
}

data "alicloud_nas_protocols" "example" {
	type = "Performance"
}

data "alicloud_nas_zones" "default" {
	file_system_type = "standard"
}
	
locals {
	count_size = length(data.alicloud_nas_zones.default.zones)
	zone_id    = data.alicloud_nas_zones.default.zones[local.count_size - 1].zone_id
}

resource "alicloud_vpc" "main" {
	vpc_name   = "terraform-example"
	cidr_block = "172.17.3.0/24"
}

resource "alicloud_vswitch" "main" {
  	vswitch_name = alicloud_vpc.main.vpc_name
  	cidr_block   = alicloud_vpc.main.cidr_block
  	vpc_id       = alicloud_vpc.main.id
  	zone_id      = local.zone_id
}

resource "alicloud_security_group" "example" {
	name = var.name
	vpc_id = "${alicloud_vpc.main.id}"
}

resource "alicloud_nas_file_system" "example" {
	protocol_type = "${data.alicloud_nas_protocols.example.protocols.0}"
	storage_type = "Performance"
}

resource "alicloud_nas_access_group" "example" {
	access_group_name = "${var.name}"
	access_group_type = "Vpc"
}

resource "alicloud_nas_access_group" "example1" {
	access_group_name = "${var.name1}"
	access_group_type = "Vpc"
}

resource "alicloud_nas_mount_target" "example" {
	file_system_id    = "${alicloud_nas_file_system.example.id}"
	access_group_name = "${alicloud_nas_access_group.example.access_group_name}"
	vswitch_id        = "${alicloud_vswitch.main.id}"
	security_group_id = "${alicloud_security_group.example.id}"
}
`, name, name)
}

var AlicloudNasMountTarget1 = map[string]string{}

func AlicloudNasMountTargetBasicDependence1(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}

variable "name1" {
	default = "%schange"
}

data "alicloud_nas_zones" "default" {
	file_system_type = "extreme"
}

locals {
	count_size = length(data.alicloud_nas_zones.default.zones)
	zone_id    = data.alicloud_nas_zones.default.zones[local.count_size - 1].zone_id
}

resource "alicloud_vpc" "main" {
  vpc_name   = "terraform-example"
  cidr_block = "172.17.3.0/24"
}

resource "alicloud_vswitch" "main" {
  vswitch_name = alicloud_vpc.main.vpc_name
  cidr_block   = alicloud_vpc.main.cidr_block
  vpc_id       = alicloud_vpc.main.id
  zone_id      = local.zone_id
}

resource "alicloud_nas_file_system" "example" {
  file_system_type = "extreme"
  protocol_type    = "NFS"
  zone_id          = data.alicloud_nas_zones.default.zones[local.count_size - 1].zone_id
  storage_type     = "advance"
  capacity         = 100
}

resource "alicloud_nas_access_group" "example" {
	access_group_name = "${var.name}"
	access_group_type = "Vpc"
	file_system_type  = "extreme"
}

resource "alicloud_nas_access_group" "example1" {
	access_group_name = "${var.name1}"
	access_group_type = "Vpc"
	file_system_type  = "extreme"
}
`, name, name)
}

// lintignore: R001
func TestAccAliCloudNASMountTargettest_extreme_MountTarget(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_nas_mount_target.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudNASMountTargettest_extreme_MountTargetMap)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return nas.NewNasServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeNasMountTarget")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sNASMountTarget%d", acctest.DefaultRegionToTest, rand)

	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudNASMountTargettest_extreme_MountTargetBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"vpc_id":            "${alicloud_vpc.createEVpc.id}",
					"network_type":      "Vpc",
					"vswitch_id":        "${alicloud_vswitch.CreateEVswitch1.id}",
					"security_group_id": "${alicloud_security_group.CreateSecurityGroup.id}",
					"access_group_name": "${alicloud_nas_access_group.create_eaccess_group.access_group_name}",
					"file_system_id":    "${alicloud_nas_file_system.create_extreme_file_system.id}",
					"dual_stack":        true,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"dual_stack", "security_group_id"},
			},
		},
	})
}

var AliCloudNASMountTargettest_extreme_MountTargetMap = map[string]string{
	"id": acctest.CHECKSET,
}

func AliCloudNASMountTargettest_extreme_MountTargetBasicDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

variable "var_zone_id1" {
  type = object({
    zone_id1 = string
  })
  default = {
    zone_id1 = "cn-hangzhou-g"
  }
}

resource "alicloud_vpc" "createEVpc" {
  cidr_block  = "192.168.0.0/16"
  vpc_name    = "nas-teste0308-vpc"
  enable_ipv6 = true
}

resource "alicloud_nas_file_system" "create_extreme_file_system" {
  description      = "挂载点E0307测试资源"
  storage_type     = "standard"
  zone_id          = var.var_zone_id1.zone_id1
  encrypt_type     = 0
  capacity         = 100
  protocol_type    = "NFS"
  file_system_type = "extreme"
}

resource "alicloud_security_group" "CreateSecurityGroup" {
  security_group_name = "testMountTarget"
  vpc_id              = alicloud_vpc.createEVpc.id
}

resource "alicloud_nas_access_group" "create_eaccess_group" {
  access_group_type = "Vpc"
  description       = "挂载点创建测试"
  access_group_name = "ExtremeMountTarget"
  file_system_type  = "extreme"
}

resource "alicloud_vswitch" "CreateEVswitch1" {
  ipv6_cidr_block_mask = 38
  vpc_id               = alicloud_vpc.createEVpc.id
  zone_id              = var.var_zone_id1.zone_id1
  cidr_block           = "192.168.0.0/24"
  vswitch_name         = "nas-teste0307-vsw1sdw"
  enable_ipv6          = true
}

`, name)
}

// Case test_MountTarget
func TestAccAliCloudNASMountTargettest_MountTarget(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_nas_mount_target.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudNASMountTargettest_MountTargetMap)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return nas.NewNasServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeNasMountTarget")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sNASMountTarget%d", acctest.DefaultRegionToTest, rand)

	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudNASMountTargettest_MountTargetBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"status":            "Active",
					"vpc_id":            "${alicloud_vpc.createVpc.id}",
					"network_type":      "Vpc",
					"access_group_name": "${alicloud_nas_access_group.create_access_group.access_group_name}",
					"vswitch_id":        "${alicloud_vswitch.CreateVswitch1.id}",
					"file_system_id":    "${alicloud_nas_file_system.create_file_system.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"access_group_name": "${var.access_group_default.access_group_name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "Inactive",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

var AliCloudNASMountTargettest_MountTargetMap = map[string]string{
	"id": acctest.CHECKSET,
}

func AliCloudNASMountTargettest_MountTargetBasicDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

variable "var_region_id" {
  type = object({
    region_id = string
  })
  default = {
    region_id = "cn-hangzhou"
  }
}

variable "var_zone_id1" {
  type = object({
    zone_id1 = string
  })
  default = {
    zone_id1 = "cn-hangzhou-g"
  }
}

variable "access_group_default" {
  type = object({
    access_group_name = string
  })
  default = {
    access_group_name = "DEFAULT_VPC_GROUP_NAME"
  }
}

resource "alicloud_vpc" "createVpc" {
  cidr_block = "172.16.0.0/12"
  vpc_name   = "nas-test0307-vpc"
}

resource "alicloud_nas_file_system" "create_file_system" {
  description      = "挂载点0307测试资源"
  storage_type     = "Performance"
  zone_id          = var.var_zone_id1.zone_id1
  encrypt_type     = 0
  protocol_type    = "NFS"
  file_system_type = "standard"
}

resource "alicloud_vswitch" "CreateVswitch1" {
  vpc_id       = alicloud_vpc.createVpc.id
  zone_id      = var.var_zone_id1.zone_id1
  cidr_block   = "172.16.0.0/24"
  vswitch_name = "nas-test0307-vsw1sdw"
}

resource "alicloud_nas_access_group" "create_access_group" {
  access_group_type = "Vpc"
  description       = "挂载点创建测试"
  access_group_name = "StandardMountTarget"
  file_system_type  = "standard"
}

`, name)
}

// Test NAS MountTarget. <<< Resource test cases, automatically generated.

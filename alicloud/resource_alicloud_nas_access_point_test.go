package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test Nas AccessPoint. >>> Resource test cases, automatically generated.
// Case 通用型接入点不传入 RootDirectory 6611
func TestAccAliCloudNasAccessPoint_basic6611(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_nas_access_point.default"
	ra := resourceAttrInit(resourceId, AliCloudNasAccessPointMap6611)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &NasServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeNasAccessPoint")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%snasaccesspoint%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudNasAccessPointBasicDependence6611)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"vpc_id":         "${alicloud_vpc.defaultkyVC70.id}",
					"access_group":   "${alicloud_nas_access_group.defaultBbc7ev.access_group_name}",
					"vswitch_id":     "${alicloud_vswitch.defaultoZAPmO.id}",
					"file_system_id": "${alicloud_nas_file_system.defaultVtUpDh.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"vpc_id":         CHECKSET,
						"access_group":   CHECKSET,
						"vswitch_id":     CHECKSET,
						"file_system_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"access_group": "${alicloud_nas_access_group.default6mnIjY.access_group_name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"access_group": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"enabled_ram": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"enabled_ram": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"enabled_ram": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"enabled_ram": "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"access_point_name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"access_point_name": name,
					}),
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

var AliCloudNasAccessPointMap6611 = map[string]string{
	"status":          CHECKSET,
	"create_time":     CHECKSET,
	"access_point_id": CHECKSET,
	"region_id":       CHECKSET,
}

func AliCloudNasAccessPointBasicDependence6611(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

variable "region_id" {
  default = "cn-hangzhou"
}

variable "azone" {
  default = "cn-hangzhou-g"
}

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

resource "alicloud_vpc" "defaultkyVC70" {
  cidr_block  = "172.16.0.0/12"
  description = "接入点测试noRootDirectory"
}

resource "alicloud_vswitch" "defaultoZAPmO" {
  vpc_id     = alicloud_vpc.defaultkyVC70.id
  zone_id    = data.alicloud_zones.default.zones.0.id
  cidr_block = "172.16.0.0/24"
}

resource "alicloud_nas_access_group" "defaultBbc7ev" {
  access_group_type = "Vpc"
  access_group_name = var.name
  file_system_type  = "standard"
}

resource "alicloud_nas_access_group" "default6mnIjY" {
  access_group_type = "Vpc"
  access_group_name = "${var.name}-update"
  file_system_type  = "standard"
}

resource "alicloud_nas_file_system" "defaultVtUpDh" {
  storage_type     = "Performance"
  zone_id          = var.azone
  encrypt_type     = "0"
  protocol_type    = "NFS"
  file_system_type = "standard"
  description      = "AccessPointnoRootDirectory"
}


`, name)
}

// Case 通用型接入点 6611  twin
func TestAccAliCloudNasAccessPoint_basic6611_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_nas_access_point.default"
	ra := resourceAttrInit(resourceId, AliCloudNasAccessPointMap6611)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &NasServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeNasAccessPoint")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%snasaccesspoint%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudNasAccessPointBasicDependence6611)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"vpc_id":            "${alicloud_vpc.defaultkyVC70.id}",
					"access_group":      "${alicloud_nas_access_group.defaultBbc7ev.access_group_name}",
					"vswitch_id":        "${alicloud_vswitch.defaultoZAPmO.id}",
					"enabled_ram":       "true",
					"file_system_id":    "${alicloud_nas_file_system.defaultVtUpDh.id}",
					"root_path":         "/terraform",
					"access_point_name": name,
					"posix_user": []map[string]interface{}{
						{
							"posix_group_id": "123",
							"posix_user_id":  "123",
						},
					},
					"root_path_permission": []map[string]interface{}{
						{
							"owner_group_id": "1",
							"owner_user_id":  "1",
							"permission":     "0777",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"vpc_id":                 CHECKSET,
						"access_group":           CHECKSET,
						"vswitch_id":             CHECKSET,
						"enabled_ram":            "true",
						"file_system_id":         CHECKSET,
						"root_path":              "/terraform",
						"access_point_name":      name,
						"posix_user.#":           "1",
						"root_path_permission.#": "1",
					}),
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

// Test Nas AccessPoint. <<< Resource test cases, automatically generated.

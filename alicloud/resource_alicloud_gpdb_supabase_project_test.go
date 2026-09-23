// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test Gpdb SupabaseProject. >>> Resource test cases, automatically generated.
// Case supabase资源测试 11921
func TestAccAliCloudGpdbSupabaseProject_basic11921(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_gpdb_supabase_project.default"
	ra := resourceAttrInit(resourceId, AliCloudGpdbSupabaseProjectMap11921)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &GpdbServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeGpdbSupabaseProject")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccgpdb%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudGpdbSupabaseProjectBasicDependence11921)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"account_password": "YourPassword123!",
					"project_name":     name,
					"project_spec":     "1C2G",
					"security_ip_list": []string{
						"127.0.0.1"},
					"zone_id":    "cn-hangzhou-j",
					"vpc_id":     "${data.alicloud_vpcs.default.ids.0}",
					"vswitch_id": "${data.alicloud_vswitches.default.ids.0}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"project_name":       name,
						"project_spec":       "1C2G",
						"security_ip_list.#": "1",
						"zone_id":            CHECKSET,
						"vpc_id":             CHECKSET,
						"vswitch_id":         CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"security_ip_list": []string{
						"127.0.0.1", "0.0.0.0/0", "140.205.11.0/24", "140.205.11.11"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"security_ip_list.#": "4",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"account_password": "YourPassword123!update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"account_password"},
			},
		},
	})
}

func TestAccAliCloudGpdbSupabaseProject_basic11921_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_gpdb_supabase_project.default"
	ra := resourceAttrInit(resourceId, AliCloudGpdbSupabaseProjectMap11921)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &GpdbServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeGpdbSupabaseProject")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccgpdb%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudGpdbSupabaseProjectBasicDependence11921)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"account_password": "YourPassword123!",
					"project_name":     name,
					"project_spec":     "1C2G",
					"security_ip_list": []string{
						"127.0.0.1", "0.0.0.0/0", "140.205.11.0/24", "140.205.11.11"},
					"zone_id":                "cn-hangzhou-j",
					"vpc_id":                 "${data.alicloud_vpcs.default.ids.0}",
					"vswitch_id":             "${data.alicloud_vswitches.default.ids.0}",
					"disk_performance_level": "PL0",
					"storage_size":           "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"project_name":           name,
						"project_spec":           "1C2G",
						"security_ip_list.#":     "4",
						"zone_id":                CHECKSET,
						"vpc_id":                 CHECKSET,
						"vswitch_id":             CHECKSET,
						"disk_performance_level": "PL0",
						"storage_size":           "1",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"account_password"},
			},
		},
	})
}

var AliCloudGpdbSupabaseProjectMap11921 = map[string]string{
	"disk_performance_level": CHECKSET,
	"storage_size":           CHECKSET,
	"create_time":            CHECKSET,
	"region_id":              CHECKSET,
	"status":                 CHECKSET,
}

func AliCloudGpdbSupabaseProjectBasicDependence11921(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

	data "alicloud_vpcs" "default" {
  		name_regex = "^default-NODELETING$"
	}

	data "alicloud_vswitches" "default" {
  		vpc_id  = data.alicloud_vpcs.default.ids.0
  		zone_id = "cn-hangzhou-j"
	}
`, name)
}

// Test Gpdb SupabaseProject. <<< Resource test cases, automatically generated.

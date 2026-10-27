package drds_test

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/drds"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// Migrated from the root alicloud package. The CHECKSET helpers live in
// acc_common_test.go; the service constructor is drds.NewService and the provider
// factory is acctest.ProtoV5ProviderFactories. TestAcc names, steps and
// assertions are preserved verbatim; only the wiring is adapted.

// Test Drds PolardbXInstance. >>> Resource test cases, automatically generated.
// Case 4504
func TestAccAliCloudDrdsPolardbxInstance_basic4504(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_drds_polardbx_instance.default"
	ra := resourceAttrInit(resourceId, AlicloudPolardbxInstanceMap4504)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return drds.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribePolardbxInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sdrdspolardbxinstance%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudPolardbxInstanceBasicDependence4504)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, connectivity.DRDSPolarDbxSupportRegions)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"topology_type":  "3azones",
					"vswitch_id":     "${data.alicloud_vswitches.default.ids.0}",
					"primary_zone":   "cn-hangzhou-h",
					"cn_node_count":  "2",
					"dn_class":       "mysql.n4.medium.25",
					"cn_class":       "polarx.x4.medium.2e",
					"dn_node_count":  "2",
					"secondary_zone": "cn-hangzhou-g",
					"tertiary_zone":  "cn-hangzhou-k",
					"vpc_id":         "${data.alicloud_vpcs.default.ids.0}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"topology_type":  "3azones",
						"vpc_id":         CHECKSET,
						"vswitch_id":     CHECKSET,
						"primary_zone":   "cn-hangzhou-h",
						"secondary_zone": "cn-hangzhou-g",
						"tertiary_zone":  "cn-hangzhou-k",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"cn_node_count": "3",
					"dn_node_count": "3",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cn_node_count": "3",
						"dn_node_count": "3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.groups.0.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.groups.1.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"topology_type":     "3azones",
					"vswitch_id":        "${data.alicloud_vswitches.default.ids.0}",
					"primary_zone":      "cn-hangzhou-h",
					"cn_node_count":     "3",
					"dn_class":          "mysql.n4.medium.25",
					"cn_class":          "polarx.x4.medium.2e",
					"dn_node_count":     "3",
					"secondary_zone":    "cn-hangzhou-g",
					"tertiary_zone":     "cn-hangzhou-k",
					"vpc_id":            "${data.alicloud_vpcs.default.ids.0}",
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.groups.0.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"topology_type":     "3azones",
						"vswitch_id":        CHECKSET,
						"primary_zone":      "cn-hangzhou-h",
						"cn_node_count":     "3",
						"dn_class":          "mysql.n4.medium.25",
						"cn_class":          "polarx.x4.medium.2e",
						"dn_node_count":     "3",
						"secondary_zone":    "cn-hangzhou-g",
						"tertiary_zone":     "cn-hangzhou-k",
						"vpc_id":            CHECKSET,
						"resource_group_id": CHECKSET,
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

var AlicloudPolardbxInstanceMap4504 = map[string]string{
	"status":      CHECKSET,
	"create_time": CHECKSET,
}

func AlicloudPolardbxInstanceBasicDependence4504(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}
data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}
data "alicloud_resource_manager_resource_groups" "default" {
  status = "OK"
}
data "alicloud_vpcs" "default" {
  name_regex = "^default-NODELETING$"
}
data "alicloud_vswitches" "default" {
  vpc_id = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_zones.default.zones.0.id
}
`, name)
}

// Case polardbx instance单可用区8.0资源创建测试 3 12277
func TestAccAliCloudDrdsPolardbxInstance_basic12277(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_drds_polardbx_instance.default"
	ra := resourceAttrInit(resourceId, AlicloudPolardbxInstanceMap12277)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return drds.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribePolardbxInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccdrds%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudPolardbxInstanceBasicDependence12277)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-beijing"})
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"topology_type":            "1azone",
					"vswitch_id":               "${data.alicloud_vswitches.default.ids.0}",
					"primary_zone":             "cn-beijing-f",
					"cn_node_count":            "2",
					"dn_class":                 "mysql.x8.large.25",
					"cn_class":                 "polarx.x8.large.2e",
					"dn_node_count":            "2",
					"vpc_id":                   "${data.alicloud_vpcs.default.ids.0}",
					"is_read_db_instance":      "false",
					"primary_db_instance_name": "null",
					"resource_group_id":        "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"engine_version":           "8.0",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"topology_type":            "1azone",
						"vswitch_id":               CHECKSET,
						"primary_zone":             "cn-beijing-f",
						"cn_node_count":            CHECKSET,
						"dn_class":                 "mysql.x8.large.25",
						"cn_class":                 "polarx.x8.large.2e",
						"dn_node_count":            CHECKSET,
						"vpc_id":                   CHECKSET,
						"is_read_db_instance":      "false",
						"primary_db_instance_name": CHECKSET,
						"resource_group_id":        CHECKSET,
						"engine_version":           CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "test_desc",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "test_desc",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"is_read_db_instance", "primary_db_instance_name"},
			},
		},
	})
}

var AlicloudPolardbxInstanceMap12277 = map[string]string{
	"status":      CHECKSET,
	"create_time": CHECKSET,
	"region_id":   CHECKSET,
}

func AlicloudPolardbxInstanceBasicDependence12277(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}
data "alicloud_resource_manager_resource_groups" "default" {}
data "alicloud_vpcs" "default" {
  name_regex = "^default-NODELETING$"
}
data "alicloud_vswitches" "default" {
  vpc_id = data.alicloud_vpcs.default.ids.0
  zone_id = "cn-beijing-f"
}

`, name)
}

// Case polardbx instance单可用区8.0资源创建测试 4 12278
func TestAccAliCloudDrdsPolardbxInstance_basic12278(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_drds_polardbx_instance.default"
	ra := resourceAttrInit(resourceId, AlicloudPolardbxInstanceMap12278)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return drds.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribePolardbxInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccdrds%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudPolardbxInstanceBasicDependence12278)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-beijing"})
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"topology_type":            "1azone",
					"vswitch_id":               "${data.alicloud_vswitches.default.ids.0}",
					"primary_zone":             "cn-beijing-f",
					"cn_node_count":            "2",
					"dn_class":                 "mysql.n4.xlarge.25",
					"cn_class":                 "polarx.x4.xlarge.2e",
					"dn_node_count":            "2",
					"vpc_id":                   "${data.alicloud_vpcs.default.ids.0}",
					"is_read_db_instance":      "false",
					"primary_db_instance_name": "null",
					"resource_group_id":        "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"engine_version":           "8.0",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"topology_type":            "1azone",
						"vswitch_id":               CHECKSET,
						"primary_zone":             "cn-beijing-f",
						"cn_node_count":            CHECKSET,
						"dn_class":                 "mysql.n4.xlarge.25",
						"cn_class":                 "polarx.x4.xlarge.2e",
						"dn_node_count":            CHECKSET,
						"vpc_id":                   CHECKSET,
						"is_read_db_instance":      "false",
						"primary_db_instance_name": CHECKSET,
						"resource_group_id":        CHECKSET,
						"engine_version":           CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "test_desc",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "test_desc",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"is_read_db_instance", "primary_db_instance_name"},
			},
		},
	})
}

var AlicloudPolardbxInstanceMap12278 = map[string]string{
	"status":      CHECKSET,
	"create_time": CHECKSET,
	"region_id":   CHECKSET,
}

func AlicloudPolardbxInstanceBasicDependence12278(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}
data "alicloud_resource_manager_resource_groups" "default" {}
data "alicloud_vpcs" "default" {
  name_regex = "^default-NODELETING$"
}
data "alicloud_vswitches" "default" {
  vpc_id = data.alicloud_vpcs.default.ids.0
  zone_id = "cn-beijing-f"
}

`, name)
}

// Case polardbx instance三可用区5.7资源生命周期测试 12257
func TestAccAliCloudDrdsPolardbxInstance_basic12257(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_drds_polardbx_instance.default"
	ra := resourceAttrInit(resourceId, AlicloudPolardbxInstanceMap12257)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return drds.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribePolardbxInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccdrds%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudPolardbxInstanceBasicDependence12257)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-beijing"})
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"topology_type":            "3azones",
					"vswitch_id":               "${data.alicloud_vswitches.default.ids.0}",
					"primary_zone":             "cn-beijing-f",
					"cn_node_count":            "2",
					"dn_class":                 "mysql.n4.medium.25",
					"cn_class":                 "polarx.x4.medium.2e",
					"dn_node_count":            "2",
					"vpc_id":                   "${data.alicloud_vpcs.default.ids.0}",
					"is_read_db_instance":      "false",
					"primary_db_instance_name": "null",
					"resource_group_id":        "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"secondary_zone":           "cn-beijing-k",
					"tertiary_zone":            "cn-beijing-h",
					"engine_version":           "5.7",
					"description":              "test57",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"topology_type":            "3azones",
						"vswitch_id":               CHECKSET,
						"primary_zone":             "cn-beijing-f",
						"cn_node_count":            CHECKSET,
						"dn_class":                 "mysql.n4.medium.25",
						"cn_class":                 "polarx.x4.medium.2e",
						"dn_node_count":            CHECKSET,
						"vpc_id":                   CHECKSET,
						"is_read_db_instance":      "false",
						"primary_db_instance_name": CHECKSET,
						"resource_group_id":        CHECKSET,
						"secondary_zone":           "cn-beijing-k",
						"tertiary_zone":            "cn-beijing-h",
						"engine_version":           CHECKSET,
						"description":              "test57",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"cn_node_count": "3",
					"dn_node_count": "3",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cn_node_count": "3",
						"dn_node_count": "3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.1}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id": CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"is_read_db_instance", "primary_db_instance_name"},
			},
		},
	})
}

var AlicloudPolardbxInstanceMap12257 = map[string]string{
	"status":      CHECKSET,
	"create_time": CHECKSET,
	"region_id":   CHECKSET,
}

func AlicloudPolardbxInstanceBasicDependence12257(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}
data "alicloud_resource_manager_resource_groups" "default" {}
data "alicloud_vpcs" "default" {
  name_regex = "^default-NODELETING$"
}
data "alicloud_vswitches" "default" {
  vpc_id = data.alicloud_vpcs.default.ids.0
  zone_id = "cn-beijing-f"
}


`, name)
}

// Case polardbx instance单可用区8.0资源生命周期测试 12261
func TestAccAliCloudDrdsPolardbxInstance_basic12261(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_drds_polardbx_instance.default"
	ra := resourceAttrInit(resourceId, AlicloudPolardbxInstanceMap12261)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return drds.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribePolardbxInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccdrds%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudPolardbxInstanceBasicDependence12261)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-beijing"})
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"topology_type":            "1azone",
					"vswitch_id":               "${data.alicloud_vswitches.default.ids.0}",
					"primary_zone":             "cn-beijing-f",
					"cn_node_count":            "2",
					"dn_class":                 "mysql.n4.medium.25",
					"cn_class":                 "polarx.x4.medium.2e",
					"dn_node_count":            "2",
					"vpc_id":                   "${data.alicloud_vpcs.default.ids.0}",
					"is_read_db_instance":      "false",
					"primary_db_instance_name": "null",
					"resource_group_id":        "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"engine_version":           "8.0",
					"description":              "test57",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"topology_type":            "1azone",
						"vswitch_id":               CHECKSET,
						"primary_zone":             "cn-beijing-f",
						"cn_node_count":            CHECKSET,
						"dn_class":                 "mysql.n4.medium.25",
						"cn_class":                 "polarx.x4.medium.2e",
						"dn_node_count":            CHECKSET,
						"vpc_id":                   CHECKSET,
						"is_read_db_instance":      "false",
						"primary_db_instance_name": CHECKSET,
						"resource_group_id":        CHECKSET,
						"engine_version":           CHECKSET,
						"description":              "test57",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "test_desc",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "test_desc",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"is_read_db_instance", "primary_db_instance_name"},
			},
		},
	})
}

var AlicloudPolardbxInstanceMap12261 = map[string]string{
	"status":      CHECKSET,
	"create_time": CHECKSET,
	"region_id":   CHECKSET,
}

func AlicloudPolardbxInstanceBasicDependence12261(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}
data "alicloud_resource_manager_resource_groups" "default" {}
data "alicloud_vpcs" "default" {
  name_regex = "^default-NODELETING$"
}
data "alicloud_vswitches" "default" {
  vpc_id = data.alicloud_vpcs.default.ids.0
  zone_id = "cn-beijing-f"
}


`, name)
}

// Case ModifyDBInstanceClass in-place class upgrade
func TestAccAliCloudDrdsPolardbxInstance_basic83815836(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_drds_polardbx_instance.default"
	ra := resourceAttrInit(resourceId, AlicloudPolardbxInstanceMap83815836)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return drds.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribePolardbxInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccdrds%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudPolardbxInstanceBasicDependence83815836)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-beijing"})
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"topology_type":  "1azone",
					"vswitch_id":     "${data.alicloud_vswitches.default.ids.0}",
					"primary_zone":   "cn-beijing-f",
					"cn_node_count":  "2",
					"dn_class":       "mysql.n4.medium.25",
					"cn_class":       "polarx.x4.medium.2e",
					"dn_node_count":  "2",
					"vpc_id":         "${data.alicloud_vpcs.default.ids.0}",
					"engine_version": "8.0",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"topology_type":  "1azone",
						"vswitch_id":     CHECKSET,
						"primary_zone":   "cn-beijing-f",
						"cn_node_count":  CHECKSET,
						"dn_class":       "mysql.n4.medium.25",
						"cn_class":       "polarx.x4.medium.2e",
						"dn_node_count":  CHECKSET,
						"vpc_id":         CHECKSET,
						"engine_version": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"cn_class": "polarx.x4.large.2e",
					"dn_class": "mysql.n4.large.25",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cn_class": "polarx.x4.large.2e",
						"dn_class": "mysql.n4.large.25",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "test_modify_class",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "test_modify_class",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"specified_dn_scale",
					"specified_dn_spec_map_json",
					"switch_time",
					"switch_time_mode",
				},
			},
		},
	})
}

var AlicloudPolardbxInstanceMap83815836 = map[string]string{
	"status":      CHECKSET,
	"create_time": CHECKSET,
	"region_id":   CHECKSET,
}

func AlicloudPolardbxInstanceBasicDependence83815836(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}
data "alicloud_vpcs" "default" {
  name_regex = "^default-NODELETING$"
}
data "alicloud_vswitches" "default" {
  vpc_id  = data.alicloud_vpcs.default.ids.0
  zone_id = "cn-beijing-f"
}

`, name)
}

// Case cloud_auto storage + ModifyDBInstanceClass class change with dn_storage_space grow
// Requires a region + class combination that has an active price plan for the c25/c2e cloud-disk family.
// Not part of the default matrix — enable per region when pricing is available.
func TestAccAliCloudDrdsPolardbxInstance_basic83815836_cloudauto(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_drds_polardbx_instance.default"
	ra := resourceAttrInit(resourceId, AlicloudPolardbxInstanceMap83815836)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return drds.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribePolardbxInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccdrds%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudPolardbxInstanceBasicDependence83815836)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-beijing"})
			acctest.PreCheck(t)
			if os.Getenv("ALICLOUD_TEST_DRDS_POLARDBX_CLOUDAUTO") == "" {
				t.Skip("skipping; set ALICLOUD_TEST_DRDS_POLARDBX_CLOUDAUTO to run the cloud_auto pricing test")
			}
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"topology_type":              "1azone",
					"vswitch_id":                 "${data.alicloud_vswitches.default.ids.0}",
					"primary_zone":               "cn-beijing-f",
					"cn_node_count":              "2",
					"dn_class":                   "polarx.mysql.x4.large.c25",
					"cn_class":                   "polarx.x4.large.c2e",
					"dn_node_count":              "2",
					"vpc_id":                     "${data.alicloud_vpcs.default.ids.0}",
					"engine_version":             "8.0",
					"storage_type":               "cloud_auto",
					"dn_storage_space":           "50",
					"specified_dn_scale":         "false",
					"specified_dn_spec_map_json": "{}",
					"switch_time_mode":           "Immediate",
					"switch_time":                "2029-01-01T00:00:00Z",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"storage_type":     "cloud_auto",
						"dn_storage_space": "50",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"dn_storage_space": "100",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"dn_storage_space": "100",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"specified_dn_scale",
					"specified_dn_spec_map_json",
					"switch_time",
					"switch_time_mode",
				},
			},
		},
	})
}

// Case polardbx instance prepay (subscription) resource creation test (AC-6 redesign).
//
// This ACC is env-gated on ALICLOUD_TEST_DRDS_POLARDBX_PREPAY because it creates a
// real subscription instance that bills immediately and whose destroy is NOT
// supported via API for this commodity/account (BssOpenApi RefundInstance returns
// CommodityNotSupported). auto_renew is false so the subscription does not renew,
// avoiding a perpetually-billing instance.
//
// Steps verify create + read + import, then a destroy step that asserts the
// documented actionable error (CommodityNotSupported → instance ID + BSS console
// 退订管理 guidance + terraform state rm). For a non-refundable commodity, destroy
// fails with that error and the resource stays in state (AC-3).
//
// Design note (F-16 fix): terraform-plugin-testing v1.16.0's runPostTestDestroy
// (deferred in runNewTest) calls wd.Destroy() and t.Fatalf when state is non-empty
// and the post-test destroy returns an error. For a non-refundable PREPAY
// commodity Delete always returns the actionable error, so any resource.Test that
// leaves the instance in state fatally fails the test. This test avoids that by
// clearing state after the ExpectError destroy step: the final step re-applies the
// full config (a no-op refresh, since the instance still exists) and its Check
// removes the managed resource entry from the on-disk terraform.tfstate (equivalent
// to `terraform state rm`), so runPostTestDestroy finds an empty state and does
// not re-invoke Delete. The destroy-actionable-error contract is verified live by
// the ExpectError step and deterministically (no billing) by
// TestUnitDrdsPolardbxInstancePrepayRefundActionableError; the instance ID is
// logged so the operator can refund it in the BSS console afterwards.
func TestAccAliCloudDrdsPolardbxInstance_prepay(t *testing.T) {
	if os.Getenv("ALICLOUD_TEST_DRDS_POLARDBX_PREPAY") == "" {
		t.Skip("skipping; set ALICLOUD_TEST_DRDS_POLARDBX_PREPAY to create a real subscription instance (bills immediately, destroy needs manual BSS refund)")
	}
	var v map[string]interface{}
	resourceId := "alicloud_drds_polardbx_instance.default"
	ra := resourceAttrInit(resourceId, AlicloudPolardbxInstanceMapprepay)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return drds.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribePolardbxInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccdrds%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudPolardbxInstanceBasicDependenceprepay)
	// dependenceOnly returns the data-source dependencies with no resource block,
	// so applying it drives terraform to destroy the instance; Delete then returns
	// the documented actionable error.
	dependenceOnly := func() string {
		return AlicloudPolardbxInstanceBasicDependenceprepay(name)
	}
	recordInstanceID := func(s *terraform.State) error {
		if rs, ok := s.RootModule().Resources[resourceId]; ok && rs.Primary.ID != "" {
			t.Logf("[PREPAY MANUAL REFUND] alicloud_drds_polardbx_instance.default ID: %s — destroy returned the documented actionable error; refund this instance in the BSS console (费用中心 > 退订管理) and remove it from state with `terraform state rm alicloud_drds_polardbx_instance.default`.", rs.Primary.ID)
		}
		return nil
	}
	// workingDir is set on TestCase.WorkingDir so all framework-created work
	// sub-directories live under it; removeResourceFromState can find their
	// terraform.tfstate files via filepath.Glob.
	workingDir := t.TempDir()
	removeResourceFromState := func(s *terraform.State) error {
		// Record the instance ID for manual refund regardless of cleanup outcome.
		if err := recordInstanceID(s); err != nil {
			return err
		}
		matches, err := filepath.Glob(filepath.Join(workingDir, "work*", "terraform.tfstate"))
		if err != nil {
			return fmt.Errorf("globbing for state file under %s: %w", workingDir, err)
		}
		if len(matches) == 0 {
			return fmt.Errorf("no terraform.tfstate found under %s; cannot clear state to prevent post-test destroy fatalf", workingDir)
		}
		for _, stateFile := range matches {
			if err := removeManagedResourceFromStateFile(stateFile, resourceId); err != nil {
				return fmt.Errorf("removing %s from %s: %w", resourceId, stateFile, err)
			}
		}
		return nil
	}
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-beijing"})
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		WorkingDir:               workingDir,
		// CheckDestroy is nil: for a non-refundable subscription commodity, the
		// instance is intentionally NOT released via API — destroy surfaces the
		// actionable error and the resource stays in state for manual refund.
		// removeResourceFromState (final step's Check) clears state so the
		// framework's runPostTestDestroy finds an empty state and does not
		// fatally fail re-invoking Delete.
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"topology_type":  "1azone",
					"vswitch_id":     "${data.alicloud_vswitches.default.ids.0}",
					"primary_zone":   "cn-beijing-f",
					"cn_node_count":  "2",
					"dn_class":       "mysql.n4.medium.25",
					"cn_class":       "polarx.x4.medium.2e",
					"dn_node_count":  "2",
					"vpc_id":         "${data.alicloud_vpcs.default.ids.0}",
					"engine_version": "8.0",
					"pay_type":       "PREPAY",
					"period":         "Month",
					"used_time":      "1",
					"auto_renew":     "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"pay_type": "PREPAY",
					}),
					resource.TestCheckResourceAttrSet(resourceId, "id"),
					recordInstanceID,
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"is_read_db_instance", "primary_db_instance_name", "period", "used_time", "auto_renew"},
			},
			{
				// Destroy the instance by removing it from config. For a
				// non-refundable subscription commodity, BssOpenApi RefundInstance
				// returns CommodityNotSupported and Delete surfaces the documented
				// actionable error; the instance stays in state for manual refund.
				Config:      dependenceOnly(),
				ExpectError: regexp.MustCompile(`(?i)CommodityNotSupported|退订管理|terraform state rm`),
			},
			{
				// Re-apply the full config so the managed resource is present in
				// config and state, making this step's apply a no-op refresh (the
				// instance still exists because Delete failed in the prior step).
				// The Check then removes the resource from the on-disk state
				// (equivalent to `terraform state rm`) so the framework's
				// runPostTestDestroy finds an empty state and does not fatally fail
				// re-invoking Delete on a non-refundable commodity. The
				// destroy-actionable-error contract was already verified live by the
				// preceding ExpectError step and deterministically by the unit test.
				Config: testAccConfig(map[string]interface{}{
					"topology_type":  "1azone",
					"vswitch_id":     "${data.alicloud_vswitches.default.ids.0}",
					"primary_zone":   "cn-beijing-f",
					"cn_node_count":  "2",
					"dn_class":       "mysql.n4.medium.25",
					"cn_class":       "polarx.x4.medium.2e",
					"dn_node_count":  "2",
					"vpc_id":         "${data.alicloud_vpcs.default.ids.0}",
					"engine_version": "8.0",
					"pay_type":       "PREPAY",
					"period":         "Month",
					"used_time":      "1",
					"auto_renew":     "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					removeResourceFromState,
				),
			},
		},
	})
}

// removeManagedResourceFromStateFile reads a terraform.tfstate JSON file, removes
// the managed resource entry whose type.name matches the given address (e.g.
// "alicloud_drds_polardbx_instance.default"), and writes it back. This is
// equivalent to `terraform state rm <address>` but operates purely on the state
// file, avoiding the need to invoke the terraform CLI binary inside the test
// framework's temporary working directory. Used by the PREPAY ACC to clear state
// after the ExpectError destroy step so runPostTestDestroy finds an empty state.
func removeManagedResourceFromStateFile(stateFile, address string) error {
	data, err := os.ReadFile(stateFile)
	if err != nil {
		return fmt.Errorf("reading state file: %w", err)
	}
	var raw map[string]interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		return fmt.Errorf("parsing state file: %w", err)
	}
	resources, ok := raw["resources"].([]interface{})
	if !ok {
		// No resources or unexpected format; nothing to remove.
		return nil
	}
	parts := strings.SplitN(address, ".", 2)
	if len(parts) != 2 {
		return fmt.Errorf("invalid resource address %q", address)
	}
	rType, rName := parts[0], parts[1]
	filtered := make([]interface{}, 0, len(resources))
	removed := false
	for _, r := range resources {
		rm, ok := r.(map[string]interface{})
		if !ok {
			filtered = append(filtered, r)
			continue
		}
		if fmt.Sprint(rm["type"]) == rType && fmt.Sprint(rm["name"]) == rName {
			removed = true
			continue
		}
		filtered = append(filtered, r)
	}
	if !removed {
		// Resource not in this state file; nothing to do.
		return nil
	}
	raw["resources"] = filtered
	out, err := json.MarshalIndent(raw, "", "  ")
	if err != nil {
		return fmt.Errorf("marshalling state file: %w", err)
	}
	if err := os.WriteFile(stateFile, out, 0644); err != nil {
		return fmt.Errorf("writing state file: %w", err)
	}
	return nil
}

var AlicloudPolardbxInstanceMapprepay = map[string]string{
	"status":      CHECKSET,
	"create_time": CHECKSET,
	"region_id":   CHECKSET,
	"pay_type":    "PREPAY",
}

func AlicloudPolardbxInstanceBasicDependenceprepay(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}
data "alicloud_resource_manager_resource_groups" "default" {}
data "alicloud_vpcs" "default" {
  name_regex = "^default-NODELETING$"
}
data "alicloud_vswitches" "default" {
  vpc_id = data.alicloud_vpcs.default.ids.0
  zone_id = "cn-beijing-f"
}

`, name)
}

// Test Drds PolardbXInstance. <<< Resource test cases, automatically generated.

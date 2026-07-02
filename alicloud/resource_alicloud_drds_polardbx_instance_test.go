package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test Drds PolardbXInstance. >>> Resource test cases, automatically generated.
// Case 4504
func TestAccAliCloudDrdsPolardbxInstance_basic4504(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_drds_polardbx_instance.default"
	ra := resourceAttrInit(resourceId, AlicloudDrdsPolardbxInstanceMap4504)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &DrdsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDrdsPolardbxInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sdrdspolardbxinstance%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudDrdsPolardbxInstanceBasicDependence4504)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.DRDSPolarDbxSupportRegions)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
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

var AlicloudDrdsPolardbxInstanceMap4504 = map[string]string{
	"status":      CHECKSET,
	"create_time": CHECKSET,
}

func AlicloudDrdsPolardbxInstanceBasicDependence4504(name string) string {
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
	ra := resourceAttrInit(resourceId, AlicloudDrdsPolardbxInstanceMap12277)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &DrdsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDrdsPolardbxInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccdrds%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudDrdsPolardbxInstanceBasicDependence12277)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-beijing"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
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

var AlicloudDrdsPolardbxInstanceMap12277 = map[string]string{
	"status":      CHECKSET,
	"create_time": CHECKSET,
	"region_id":   CHECKSET,
}

func AlicloudDrdsPolardbxInstanceBasicDependence12277(name string) string {
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
	ra := resourceAttrInit(resourceId, AlicloudDrdsPolardbxInstanceMap12278)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &DrdsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDrdsPolardbxInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccdrds%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudDrdsPolardbxInstanceBasicDependence12278)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-beijing"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
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

var AlicloudDrdsPolardbxInstanceMap12278 = map[string]string{
	"status":      CHECKSET,
	"create_time": CHECKSET,
	"region_id":   CHECKSET,
}

func AlicloudDrdsPolardbxInstanceBasicDependence12278(name string) string {
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
	ra := resourceAttrInit(resourceId, AlicloudDrdsPolardbxInstanceMap12257)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &DrdsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDrdsPolardbxInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccdrds%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudDrdsPolardbxInstanceBasicDependence12257)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-beijing"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
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

var AlicloudDrdsPolardbxInstanceMap12257 = map[string]string{
	"status":      CHECKSET,
	"create_time": CHECKSET,
	"region_id":   CHECKSET,
}

func AlicloudDrdsPolardbxInstanceBasicDependence12257(name string) string {
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
	ra := resourceAttrInit(resourceId, AlicloudDrdsPolardbxInstanceMap12261)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &DrdsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDrdsPolardbxInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccdrds%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudDrdsPolardbxInstanceBasicDependence12261)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-beijing"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
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

var AlicloudDrdsPolardbxInstanceMap12261 = map[string]string{
	"status":      CHECKSET,
	"create_time": CHECKSET,
	"region_id":   CHECKSET,
}

func AlicloudDrdsPolardbxInstanceBasicDependence12261(name string) string {
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
	ra := resourceAttrInit(resourceId, AlicloudDrdsPolardbxInstanceMap83815836)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &DrdsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDrdsPolardbxInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccdrds%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudDrdsPolardbxInstanceBasicDependence83815836)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-beijing"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
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

var AlicloudDrdsPolardbxInstanceMap83815836 = map[string]string{
	"status":      CHECKSET,
	"create_time": CHECKSET,
	"region_id":   CHECKSET,
}

func AlicloudDrdsPolardbxInstanceBasicDependence83815836(name string) string {
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
	ra := resourceAttrInit(resourceId, AlicloudDrdsPolardbxInstanceMap83815836)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &DrdsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDrdsPolardbxInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccdrds%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudDrdsPolardbxInstanceBasicDependence83815836)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-beijing"})
			testAccPreCheck(t)
			testAccPreCheckWithEnvVariable(t, "ALICLOUD_TEST_DRDS_POLARDBX_CLOUDAUTO")
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
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

// Test Drds PolardbxInstance. <<< Resource test cases, automatically generated.

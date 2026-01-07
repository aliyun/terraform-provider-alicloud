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
					"vswitch_id":     "${alicloud_vswitch.defaultV9mMOX.id}",
					"primary_zone":   "cn-hangzhou-h",
					"cn_node_count":  "2",
					"dn_class":       "mysql.n4.medium.25",
					"cn_class":       "polarx.x4.medium.2e",
					"dn_node_count":  "2",
					"secondary_zone": "cn-hangzhou-g",
					"tertiary_zone":  "cn-hangzhou-k",
					"vpc_id":         "${alicloud_vpc.defaultI3SPrf.id}",
					"engine_version": "8.0",
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
					"vswitch_id":        "${alicloud_vswitch.defaultV9mMOX.id}",
					"primary_zone":      "cn-hangzhou-h",
					"cn_node_count":     "3",
					"dn_class":          "mysql.n4.medium.25",
					"cn_class":          "polarx.x4.medium.2e",
					"dn_node_count":     "3",
					"secondary_zone":    "cn-hangzhou-g",
					"tertiary_zone":     "cn-hangzhou-k",
					"vpc_id":            "${alicloud_vpc.defaultI3SPrf.id}",
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
resource "alicloud_vpc" "defaultI3SPrf" {
  vpc_name = var.name
}
resource "alicloud_vswitch" "defaultV9mMOX" {
  vpc_id       = alicloud_vpc.defaultI3SPrf.id
  zone_id      = data.alicloud_zones.default.zones.0.id
  cidr_block   = "172.16.0.0/24"
  vswitch_name = var.name
}
`, name)
}

// Case 4497
func TestAccAliCloudDrdsPolardbxInstance_basic4497(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_drds_polardbx_instance.default"
	ra := resourceAttrInit(resourceId, AlicloudDrdsPolardbxInstanceMap4497)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &DrdsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDrdsPolardbxInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sdrdspolardbxinstance%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudDrdsPolardbxInstanceBasicDependence4497)
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
					"vswitch_id":     "${alicloud_vswitch.defaultV9mMOX.id}",
					"primary_zone":   "cn-hangzhou-h",
					"cn_node_count":  "3",
					"dn_class":       "mysql.n4.medium.25",
					"cn_class":       "polarx.x4.medium.2e",
					"dn_node_count":  "3",
					"secondary_zone": "cn-hangzhou-g",
					"tertiary_zone":  "cn-hangzhou-k",
					"vpc_id":         "${alicloud_vpc.defaultI3SPrf.id}",
					"engine_version": "8.0",
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
					"vswitch_id":        "${alicloud_vswitch.defaultV9mMOX.id}",
					"primary_zone":      "cn-hangzhou-h",
					"cn_node_count":     "3",
					"dn_class":          "mysql.n4.medium.25",
					"cn_class":          "polarx.x4.medium.2e",
					"dn_node_count":     "3",
					"secondary_zone":    "cn-hangzhou-g",
					"tertiary_zone":     "cn-hangzhou-k",
					"vpc_id":            "${alicloud_vpc.defaultI3SPrf.id}",
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

var AlicloudDrdsPolardbxInstanceMap4497 = map[string]string{
	"status":      CHECKSET,
	"create_time": CHECKSET,
}

func AlicloudDrdsPolardbxInstanceBasicDependence4497(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}
data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}
data "alicloud_resource_manager_resource_groups" "default" {}
resource "alicloud_vpc" "defaultI3SPrf" {
  vpc_name = var.name
}
resource "alicloud_vswitch" "defaultV9mMOX" {
  vpc_id       = alicloud_vpc.defaultI3SPrf.id
  zone_id      = data.alicloud_zones.default.zones.0.id
  cidr_block   = "172.16.0.0/24"
  vswitch_name = var.name
}
`, name)
}

// Case 4504  twin
func TestAccAliCloudDrdsPolardbxInstance_basic4504_twin(t *testing.T) {
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
					"vswitch_id":     "${alicloud_vswitch.defaultV9mMOX.id}",
					"primary_zone":   "cn-hangzhou-h",
					"cn_node_count":  "3",
					"dn_class":       "mysql.n4.medium.25",
					"cn_class":       "polarx.x4.medium.2e",
					"dn_node_count":  "3",
					"secondary_zone": "cn-hangzhou-g",
					"tertiary_zone":  "cn-hangzhou-k",
					"vpc_id":         "${alicloud_vpc.defaultI3SPrf.id}",
					"engine_version": "8.0",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"topology_type":  "3azones",
						"vswitch_id":     CHECKSET,
						"primary_zone":   "cn-hangzhou-h",
						"cn_node_count":  "3",
						"dn_class":       "mysql.n4.medium.25",
						"cn_class":       "polarx.x4.medium.2e",
						"dn_node_count":  "3",
						"secondary_zone": "cn-hangzhou-g",
						"tertiary_zone":  "cn-hangzhou-k",
						"vpc_id":         CHECKSET,
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

// Case 4497  twin
func TestAccAliCloudDrdsPolardbxInstance_basic4497_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_drds_polardbx_instance.default"
	ra := resourceAttrInit(resourceId, AlicloudDrdsPolardbxInstanceMap4497)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &DrdsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDrdsPolardbxInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sdrdspolardbxinstance%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudDrdsPolardbxInstanceBasicDependence4497)
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
					"vswitch_id":     "${alicloud_vswitch.defaultV9mMOX.id}",
					"primary_zone":   "cn-hangzhou-h",
					"cn_node_count":  "3",
					"dn_class":       "mysql.n4.medium.25",
					"cn_class":       "polarx.x4.medium.2e",
					"dn_node_count":  "3",
					"secondary_zone": "cn-hangzhou-g",
					"tertiary_zone":  "cn-hangzhou-k",
					"vpc_id":         "${alicloud_vpc.defaultI3SPrf.id}",
					"engine_version": "8.0",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"topology_type":  "3azones",
						"vswitch_id":     CHECKSET,
						"primary_zone":   "cn-hangzhou-h",
						"cn_node_count":  "3",
						"dn_class":       "mysql.n4.medium.25",
						"cn_class":       "polarx.x4.medium.2e",
						"dn_node_count":  "3",
						"secondary_zone": "cn-hangzhou-g",
						"tertiary_zone":  "cn-hangzhou-k",
						"vpc_id":         CHECKSET,
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

func TestAccAliCloudDrdsPolardbxInstance_basic4497_single(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_drds_polardbx_instance.default"
	ra := resourceAttrInit(resourceId, AlicloudDrdsPolardbxInstanceMap4497)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &DrdsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDrdsPolardbxInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sdrdspolardbxinstance%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudDrdsPolardbxInstanceBasicDependence4497)
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
					"topology_type": "1azone",
					"vswitch_id":    "${alicloud_vswitch.defaultV9mMOX.id}",
					"primary_zone":  "cn-hangzhou-h",
					"cn_node_count": "2",
					"dn_class":      "mysql.n4.medium.25",
					"cn_class":      "polarx.x4.medium.2e",
					"dn_node_count": "2",
					"vpc_id":        "${alicloud_vpc.defaultI3SPrf.id}",
					"engine_version": "8.0",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"topology_type": "1azone",
						"vswitch_id":    CHECKSET,
						"primary_zone":  "cn-hangzhou-h",
						"cn_node_count": "2",
						"dn_class":      "mysql.n4.medium.25",
						"cn_class":      "polarx.x4.medium.2e",
						"dn_node_count": "2",
						"vpc_id":        CHECKSET,
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

// Test Drds PolardbXInstance. <<< Resource test cases, automatically generated.
// Test Drds PolardbxInstance. >>> Resource test cases, automatically generated.
// Case polardbx instance单可用区8.0资源创建测试 7 12282
func TestAccAliCloudDrdsPolardbxInstance_basic12282(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_drds_polardbx_instance.default"
	ra := resourceAttrInit(resourceId, AlicloudDrdsPolardbxInstanceMap12282)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &DrdsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDrdsPolardbxInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccdrds%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudDrdsPolardbxInstanceBasicDependence12282)
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
					"topology_type":            "1azone",
					"vswitch_id":               "${alicloud_vswitch.default0amZnE.id}",
					"primary_zone":             "cn-beijing-i",
					"cn_node_count":            "2",
					"dn_class":                 "mysql.x8.2xlarge.25",
					"cn_class":                 "polarx.x8.2xlarge.2e",
					"dn_node_count":            "2",
					"vpc_id":                   "${alicloud_vpc.defaultlKripL.id}",
					"is_read_db_instance":      "false",
					"primary_db_instance_name": "null",
					"resource_group_id":        "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"engine_version":           "8.0",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"topology_type":            "1azone",
						"vswitch_id":               CHECKSET,
						"primary_zone":             "cn-beijing-i",
						"cn_node_count":            CHECKSET,
						"dn_class":                 "mysql.x8.2xlarge.25",
						"cn_class":                 "polarx.x8.2xlarge.2e",
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

var AlicloudDrdsPolardbxInstanceMap12282 = map[string]string{
	"status":      CHECKSET,
	"create_time": CHECKSET,
	"region_id":   CHECKSET,
}

func AlicloudDrdsPolardbxInstanceBasicDependence12282(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_resource_manager_resource_groups" "default" {}

resource "alicloud_vpc" "defaultlKripL" {
  vpc_name = "terraform-1azone-test-1"
}

resource "alicloud_vswitch" "default0amZnE" {
  vpc_id       = alicloud_vpc.defaultlKripL.id
  zone_id      = "cn-beijing-f"
  cidr_block   = "172.16.0.0/24"
  vswitch_name = "terraform-1azone-test-1"
}


`, name)
}

// Case polardbx instance单可用区8.0资源创建测试 5 12279
func TestAccAliCloudDrdsPolardbxInstance_basic12279(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_drds_polardbx_instance.default"
	ra := resourceAttrInit(resourceId, AlicloudDrdsPolardbxInstanceMap12279)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &DrdsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDrdsPolardbxInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccdrds%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudDrdsPolardbxInstanceBasicDependence12279)
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
					"topology_type":            "1azone",
					"vswitch_id":               "${alicloud_vswitch.default0amZnE.id}",
					"primary_zone":             "cn-beijing-f",
					"cn_node_count":            "2",
					"dn_class":                 "mysql.x8.xlarge.25",
					"cn_class":                 "polarx.x8.xlarge.2e",
					"dn_node_count":            "2",
					"vpc_id":                   "${alicloud_vpc.defaultlKripL.id}",
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
						"dn_class":                 "mysql.x8.xlarge.25",
						"cn_class":                 "polarx.x8.xlarge.2e",
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

var AlicloudDrdsPolardbxInstanceMap12279 = map[string]string{
	"status":      CHECKSET,
	"create_time": CHECKSET,
	"region_id":   CHECKSET,
}

func AlicloudDrdsPolardbxInstanceBasicDependence12279(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_resource_manager_resource_groups" "default" {}

resource "alicloud_vpc" "defaultlKripL" {
  vpc_name = "terraform-1azone-test-1"
}

resource "alicloud_vswitch" "default0amZnE" {
  vpc_id       = alicloud_vpc.defaultlKripL.id
  zone_id      = "cn-beijing-f"
  cidr_block   = "172.16.0.0/24"
  vswitch_name = "terraform-1azone-test-1"
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
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"topology_type":            "1azone",
					"vswitch_id":               "${alicloud_vswitch.default0amZnE.id}",
					"primary_zone":             "cn-beijing-f",
					"cn_node_count":            "2",
					"dn_class":                 "mysql.x8.large.25",
					"cn_class":                 "polarx.x8.large.2e",
					"dn_node_count":            "2",
					"vpc_id":                   "${alicloud_vpc.defaultlKripL.id}",
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

data "alicloud_resource_manager_resource_groups" "default" {}

resource "alicloud_vpc" "defaultlKripL" {
  vpc_name = "terraform-1azone-test-1"
}

resource "alicloud_vswitch" "default0amZnE" {
  vpc_id       = alicloud_vpc.defaultlKripL.id
  zone_id      = "cn-beijing-f"
  cidr_block   = "172.16.0.0/24"
  vswitch_name = "terraform-1azone-test-1"
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
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"topology_type":            "1azone",
					"vswitch_id":               "${alicloud_vswitch.default0amZnE.id}",
					"primary_zone":             "cn-beijing-f",
					"cn_node_count":            "2",
					"dn_class":                 "mysql.n4.xlarge.25",
					"cn_class":                 "polarx.x4.xlarge.2e",
					"dn_node_count":            "2",
					"vpc_id":                   "${alicloud_vpc.defaultlKripL.id}",
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

data "alicloud_resource_manager_resource_groups" "default" {}

resource "alicloud_vpc" "defaultlKripL" {
  vpc_name = "terraform-1azone-test-1"
}

resource "alicloud_vswitch" "default0amZnE" {
  vpc_id       = alicloud_vpc.defaultlKripL.id
  zone_id      = "cn-beijing-f"
  cidr_block   = "172.16.0.0/24"
  vswitch_name = "terraform-1azone-test-1"
}


`, name)
}

// Case polardbx instance单可用区8.0资源创建测试 9 12284
func TestAccAliCloudDrdsPolardbxInstance_basic12284(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_drds_polardbx_instance.default"
	ra := resourceAttrInit(resourceId, AlicloudDrdsPolardbxInstanceMap12284)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &DrdsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDrdsPolardbxInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccdrds%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudDrdsPolardbxInstanceBasicDependence12284)
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
					"topology_type":            "1azone",
					"vswitch_id":               "${alicloud_vswitch.default0amZnE.id}",
					"primary_zone":             "cn-beijing-i",
					"cn_node_count":            "2",
					"dn_class":                 "mysql.x8.4xlarge.25",
					"cn_class":                 "polarx.x8.4xlarge.2e",
					"dn_node_count":            "2",
					"vpc_id":                   "${alicloud_vpc.defaultlKripL.id}",
					"is_read_db_instance":      "false",
					"primary_db_instance_name": "null",
					"resource_group_id":        "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"engine_version":           "8.0",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"topology_type":            "1azone",
						"vswitch_id":               CHECKSET,
						"primary_zone":             "cn-beijing-i",
						"cn_node_count":            CHECKSET,
						"dn_class":                 "mysql.x8.4xlarge.25",
						"cn_class":                 "polarx.x8.4xlarge.2e",
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

var AlicloudDrdsPolardbxInstanceMap12284 = map[string]string{
	"status":      CHECKSET,
	"create_time": CHECKSET,
	"region_id":   CHECKSET,
}

func AlicloudDrdsPolardbxInstanceBasicDependence12284(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_resource_manager_resource_groups" "default" {}

resource "alicloud_vpc" "defaultlKripL" {
  vpc_name = "terraform-1azone-test-1"
}

resource "alicloud_vswitch" "default0amZnE" {
  vpc_id       = alicloud_vpc.defaultlKripL.id
  zone_id      = "cn-beijing-f"
  cidr_block   = "172.16.0.0/24"
  vswitch_name = "terraform-1azone-test-1"
}


`, name)
}

// Case polardbx instance单可用区8.0资源创建测试 6 12280
func TestAccAliCloudDrdsPolardbxInstance_basic12280(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_drds_polardbx_instance.default"
	ra := resourceAttrInit(resourceId, AlicloudDrdsPolardbxInstanceMap12280)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &DrdsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDrdsPolardbxInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccdrds%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudDrdsPolardbxInstanceBasicDependence12280)
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
					"topology_type":            "1azone",
					"vswitch_id":               "${alicloud_vswitch.default0amZnE.id}",
					"primary_zone":             "cn-beijing-i",
					"cn_node_count":            "2",
					"dn_class":                 "mysql.n4.2xlarge.25",
					"cn_class":                 "polarx.x4.2xlarge.2e",
					"dn_node_count":            "2",
					"vpc_id":                   "${alicloud_vpc.defaultlKripL.id}",
					"is_read_db_instance":      "false",
					"primary_db_instance_name": "null",
					"resource_group_id":        "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"engine_version":           "8.0",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"topology_type":            "1azone",
						"vswitch_id":               CHECKSET,
						"primary_zone":             "cn-beijing-i",
						"cn_node_count":            CHECKSET,
						"dn_class":                 "mysql.n4.2xlarge.25",
						"cn_class":                 "polarx.x4.2xlarge.2e",
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

var AlicloudDrdsPolardbxInstanceMap12280 = map[string]string{
	"status":      CHECKSET,
	"create_time": CHECKSET,
	"region_id":   CHECKSET,
}

func AlicloudDrdsPolardbxInstanceBasicDependence12280(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_resource_manager_resource_groups" "default" {}

resource "alicloud_vpc" "defaultlKripL" {
  vpc_name = "terraform-1azone-test-1"
}

resource "alicloud_vswitch" "default0amZnE" {
  vpc_id       = alicloud_vpc.defaultlKripL.id
  zone_id      = "cn-beijing-f"
  cidr_block   = "172.16.0.0/24"
  vswitch_name = "terraform-1azone-test-1"
}


`, name)
}

// Case polardbx instance单可用区8.0资源创建测试 8 12283
func TestAccAliCloudDrdsPolardbxInstance_basic12283(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_drds_polardbx_instance.default"
	ra := resourceAttrInit(resourceId, AlicloudDrdsPolardbxInstanceMap12283)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &DrdsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDrdsPolardbxInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccdrds%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudDrdsPolardbxInstanceBasicDependence12283)
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
					"topology_type":            "1azone",
					"vswitch_id":               "${alicloud_vswitch.default0amZnE.id}",
					"primary_zone":             "cn-beijing-i",
					"cn_node_count":            "2",
					"dn_class":                 "mysql.x4.4xlarge.25",
					"cn_class":                 "polarx.x4.4xlarge.2e",
					"dn_node_count":            "2",
					"vpc_id":                   "${alicloud_vpc.defaultlKripL.id}",
					"is_read_db_instance":      "false",
					"primary_db_instance_name": "null",
					"resource_group_id":        "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"engine_version":           "8.0",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"topology_type":            "1azone",
						"vswitch_id":               CHECKSET,
						"primary_zone":             "cn-beijing-i",
						"cn_node_count":            CHECKSET,
						"dn_class":                 "mysql.x4.4xlarge.25",
						"cn_class":                 "polarx.x4.4xlarge.2e",
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

var AlicloudDrdsPolardbxInstanceMap12283 = map[string]string{
	"status":      CHECKSET,
	"create_time": CHECKSET,
	"region_id":   CHECKSET,
}

func AlicloudDrdsPolardbxInstanceBasicDependence12283(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_resource_manager_resource_groups" "default" {}

resource "alicloud_vpc" "defaultlKripL" {
  vpc_name = "terraform-1azone-test-1"
}

resource "alicloud_vswitch" "default0amZnE" {
  vpc_id       = alicloud_vpc.defaultlKripL.id
  zone_id      = "cn-beijing-f"
  cidr_block   = "172.16.0.0/24"
  vswitch_name = "terraform-1azone-test-1"
}


`, name)
}

// Case polardbx instance单可用区8.0资源创建测试 11 12290
func TestAccAliCloudDrdsPolardbxInstance_basic12290(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_drds_polardbx_instance.default"
	ra := resourceAttrInit(resourceId, AlicloudDrdsPolardbxInstanceMap12290)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &DrdsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDrdsPolardbxInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccdrds%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudDrdsPolardbxInstanceBasicDependence12290)
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
					"topology_type":            "1azone",
					"vswitch_id":               "${alicloud_vswitch.default0amZnE.id}",
					"primary_zone":             "cn-hangzhou-j",
					"cn_node_count":            "2",
					"dn_class":                 "mysql.st.12xlarge.25",
					"cn_class":                 "polarx.st.12xlarge.2e",
					"dn_node_count":            "2",
					"vpc_id":                   "${alicloud_vpc.defaultlKripL.id}",
					"is_read_db_instance":      "false",
					"primary_db_instance_name": "null",
					"resource_group_id":        "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"engine_version":           "8.0",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"topology_type":            "1azone",
						"vswitch_id":               CHECKSET,
						"primary_zone":             "cn-hangzhou-j",
						"cn_node_count":            CHECKSET,
						"dn_class":                 "mysql.st.12xlarge.25",
						"cn_class":                 "polarx.st.12xlarge.2e",
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

var AlicloudDrdsPolardbxInstanceMap12290 = map[string]string{
	"status":      CHECKSET,
	"create_time": CHECKSET,
	"region_id":   CHECKSET,
}

func AlicloudDrdsPolardbxInstanceBasicDependence12290(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_resource_manager_resource_groups" "default" {}

resource "alicloud_vpc" "defaultlKripL" {
  vpc_name   = "杭州polarx测试VPC可用区I"
  cidr_block = "172.16.0.0/12"
}

resource "alicloud_vswitch" "default0amZnE" {
  vpc_id       = alicloud_vpc.defaultlKripL.id
  zone_id      = "cn-hangzhou-k"
  cidr_block   = "172.31.0.0/16"
  vswitch_name = "terraform-hangzhou"
  description  = "tf-hangzhou"
}


`, name)
}

// Case polardbx instance单可用区8.0资源创建测试 10 12285
func TestAccAliCloudDrdsPolardbxInstance_basic12285(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_drds_polardbx_instance.default"
	ra := resourceAttrInit(resourceId, AlicloudDrdsPolardbxInstanceMap12285)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &DrdsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDrdsPolardbxInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccdrds%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudDrdsPolardbxInstanceBasicDependence12285)
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
					"topology_type":            "1azone",
					"vswitch_id":               "${alicloud_vswitch.default0amZnE.id}",
					"primary_zone":             "cn-hangzhou-k",
					"cn_node_count":            "2",
					"dn_class":                 "mysql.st.8xlarge.25",
					"cn_class":                 "polarx.st.8xlarge.2e",
					"dn_node_count":            "2",
					"vpc_id":                   "${alicloud_vpc.defaultlKripL.id}",
					"is_read_db_instance":      "false",
					"primary_db_instance_name": "null",
					"resource_group_id":        "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"engine_version":           "8.0",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"topology_type":            "1azone",
						"vswitch_id":               CHECKSET,
						"primary_zone":             "cn-hangzhou-k",
						"cn_node_count":            CHECKSET,
						"dn_class":                 "mysql.st.8xlarge.25",
						"cn_class":                 "polarx.st.8xlarge.2e",
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

var AlicloudDrdsPolardbxInstanceMap12285 = map[string]string{
	"status":      CHECKSET,
	"create_time": CHECKSET,
	"region_id":   CHECKSET,
}

func AlicloudDrdsPolardbxInstanceBasicDependence12285(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_resource_manager_resource_groups" "default" {}

resource "alicloud_vpc" "defaultlKripL" {
  vpc_name   = "杭州polarx测试VPC可用区I"
  cidr_block = "172.16.0.0/12"
}

resource "alicloud_vswitch" "default0amZnE" {
  vpc_id       = alicloud_vpc.defaultlKripL.id
  zone_id      = "cn-hangzhou-k"
  cidr_block   = "172.31.0.0/16"
  vswitch_name = "terraform-hangzhou"
  description  = "tf-hangzhou"
}


`, name)
}

// Case polardbx instance单可用区8.0资源创建测试 2 12276
func TestAccAliCloudDrdsPolardbxInstance_basic12276(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_drds_polardbx_instance.default"
	ra := resourceAttrInit(resourceId, AlicloudDrdsPolardbxInstanceMap12276)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &DrdsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDrdsPolardbxInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccdrds%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudDrdsPolardbxInstanceBasicDependence12276)
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
					"topology_type":            "1azone",
					"vswitch_id":               "${alicloud_vswitch.default0amZnE.id}",
					"primary_zone":             "cn-beijing-f",
					"cn_node_count":            "2",
					"dn_class":                 "mysql.n4.large.25",
					"cn_class":                 "polarx.x4.large.2e",
					"dn_node_count":            "2",
					"vpc_id":                   "${alicloud_vpc.defaultlKripL.id}",
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
						"dn_class":                 "mysql.n4.large.25",
						"cn_class":                 "polarx.x4.large.2e",
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

var AlicloudDrdsPolardbxInstanceMap12276 = map[string]string{
	"status":      CHECKSET,
	"create_time": CHECKSET,
	"region_id":   CHECKSET,
}

func AlicloudDrdsPolardbxInstanceBasicDependence12276(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_resource_manager_resource_groups" "default" {}

resource "alicloud_vpc" "defaultlKripL" {
  vpc_name = "terraform-1azone-test-1"
}

resource "alicloud_vswitch" "default0amZnE" {
  vpc_id       = alicloud_vpc.defaultlKripL.id
  zone_id      = "cn-beijing-f"
  cidr_block   = "172.16.0.0/24"
  vswitch_name = "terraform-1azone-test-1"
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
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"topology_type":            "3azones",
					"vswitch_id":               "${alicloud_vswitch.default0amZnE.id}",
					"primary_zone":             "cn-beijing-f",
					"cn_node_count":            "2",
					"dn_class":                 "mysql.n4.medium.25",
					"cn_class":                 "polarx.x4.medium.2e",
					"dn_node_count":            "2",
					"vpc_id":                   "${alicloud_vpc.defaultlKripL.id}",
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

data "alicloud_resource_manager_resource_groups" "default" {}

resource "alicloud_vpc" "defaultlKripL" {
  vpc_name = "terraform-1azone-test-1"
}

resource "alicloud_vswitch" "default0amZnE" {
  vpc_id       = alicloud_vpc.defaultlKripL.id
  zone_id      = "cn-beijing-f"
  cidr_block   = "172.16.0.0/24"
  vswitch_name = "terraform-1azone-test-1"
}

resource "alicloud_vpc" "defaultXKFrZe" {
  vpc_name = "terraform-1azone-test-1"
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
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"topology_type":            "1azone",
					"vswitch_id":               "${alicloud_vswitch.default0amZnE.id}",
					"primary_zone":             "cn-beijing-f",
					"cn_node_count":            "2",
					"dn_class":                 "mysql.n4.medium.25",
					"cn_class":                 "polarx.x4.medium.2e",
					"dn_node_count":            "2",
					"vpc_id":                   "${alicloud_vpc.defaultlKripL.id}",
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

data "alicloud_resource_manager_resource_groups" "default" {}

resource "alicloud_vpc" "defaultlKripL" {
  vpc_name = "terraform-1azone-test-1"
}

resource "alicloud_vswitch" "default0amZnE" {
  vpc_id       = alicloud_vpc.defaultlKripL.id
  zone_id      = "cn-beijing-f"
  cidr_block   = "172.16.0.0/24"
  vswitch_name = "terraform-1azone-test-1"
}

resource "alicloud_vpc" "defaultXKFrZe" {
  vpc_name = "terraform-1azone-test-1"
}


`, name)
}

// Test Drds PolardbxInstance. <<< Resource test cases, automatically generated.

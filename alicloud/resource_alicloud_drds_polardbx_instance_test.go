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

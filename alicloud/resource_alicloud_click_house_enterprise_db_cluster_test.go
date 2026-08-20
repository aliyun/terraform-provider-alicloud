package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

// Test ClickHouse EnterpriseDBCluster. >>> Resource test cases, automatically generated.
// Case CK企业版-基本资源-多可用区1-线上 10560
func TestAccAliCloudClickHouseEnterpriseDBCluster_basic10560(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_click_house_enterprise_db_cluster.default"
	ra := resourceAttrInit(resourceId, AlicloudClickHouseEnterpriseDBClusterMap10560)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ClickHouseServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeClickHouseEnterpriseDBCluster")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccclickhouse%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudClickHouseEnterpriseDBClusterBasicDependence10560)
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
					"zone_id":    "${var.zone_id_i}",
					"vpc_id":     "${alicloud_vpc.defaultktKLuM.id}",
					"scale_min":  "8",
					"scale_max":  "16",
					"vswitch_id": "${alicloud_vswitch.defaultTQWN3k.id}",
					"multi_zones": []map[string]interface{}{
						{
							"vswitch_ids": []string{
								"${alicloud_vswitch.defaultTQWN3k.id}"},
							"zone_id": "${var.zone_id_i}",
						},
						{
							"vswitch_ids": []string{
								"${alicloud_vswitch.defaultylyLu8.id}"},
							"zone_id": "${var.zone_id_l}",
						},
						{
							"vswitch_ids": []string{
								"${alicloud_vswitch.defaultRNbPh8.id}"},
							"zone_id": "${var.zone_id_k}",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"zone_id":       CHECKSET,
						"vpc_id":        CHECKSET,
						"scale_min":     CHECKSET,
						"scale_max":     CHECKSET,
						"vswitch_id":    CHECKSET,
						"multi_zones.#": "3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"scale_min": "32",
					"scale_max": "64",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"scale_min": CHECKSET,
						"scale_max": CHECKSET,
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

var AlicloudClickHouseEnterpriseDBClusterMap10560 = map[string]string{
	"status":      CHECKSET,
	"create_time": CHECKSET,
	"region_id":   CHECKSET,
}

func AlicloudClickHouseEnterpriseDBClusterBasicDependence10560(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

variable "vsw_ip_range_i" {
  default = "172.16.1.0/24"
}

variable "region_id" {
  default = "cn-beijing"
}

variable "vpc_ip_range" {
  default = "172.16.0.0/12"
}

variable "vsw_ip_range_k" {
  default = "172.16.3.0/24"
}

variable "vsw_ip_range_l" {
  default = "172.16.2.0/24"
}

variable "zone_id_i" {
  default = "cn-beijing-i"
}

variable "zone_id_l" {
  default = "cn-beijing-l"
}

variable "zone_id_k" {
  default = "cn-beijing-k"
}

resource "alicloud_vpc" "defaultktKLuM" {
  cidr_block = var.vpc_ip_range
}

resource "alicloud_vswitch" "defaultTQWN3k" {
  vpc_id     = alicloud_vpc.defaultktKLuM.id
  zone_id    = var.zone_id_i
  cidr_block = var.vsw_ip_range_i
}

resource "alicloud_vswitch" "defaultylyLu8" {
  vpc_id     = alicloud_vpc.defaultktKLuM.id
  zone_id    = var.zone_id_l
  cidr_block = var.vsw_ip_range_l
}

resource "alicloud_vswitch" "defaultRNbPh8" {
  vpc_id     = alicloud_vpc.defaultktKLuM.id
  zone_id    = var.zone_id_k
  cidr_block = var.vsw_ip_range_k
}


`, name)
}

// Case CK企业版-基本资源-实例Id_网络_CCU_单可用区 10226
func TestAccAliCloudClickHouseEnterpriseDBCluster_basic10226(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_click_house_enterprise_db_cluster.default"
	ra := resourceAttrInit(resourceId, AlicloudClickHouseEnterpriseDBClusterMap10226)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ClickHouseServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeClickHouseEnterpriseDBCluster")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccclickhouse%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudClickHouseEnterpriseDBClusterBasicDependence10226)
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
					"zone_id":    "${var.zone_id_i}",
					"vpc_id":     "${alicloud_vpc.defaultktKLuM.id}",
					"scale_min":  "8",
					"scale_max":  "16",
					"vswitch_id": "${alicloud_vswitch.defaultTQWN3k.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"zone_id":    CHECKSET,
						"vpc_id":     CHECKSET,
						"scale_min":  CHECKSET,
						"scale_max":  CHECKSET,
						"vswitch_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"scale_min": "32",
					"scale_max": "64",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"scale_min": CHECKSET,
						"scale_max": CHECKSET,
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

var AlicloudClickHouseEnterpriseDBClusterMap10226 = map[string]string{
	"status":      CHECKSET,
	"create_time": CHECKSET,
	"region_id":   CHECKSET,
}

func AlicloudClickHouseEnterpriseDBClusterBasicDependence10226(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

variable "vsw_ip_range_i" {
  default = "172.16.1.0/24"
}

variable "region_id" {
  default = "cn-beijing"
}

variable "vpc_ip_range" {
  default = "172.16.0.0/12"
}

variable "zone_id_i" {
  default = "cn-beijing-i"
}

resource "alicloud_vpc" "defaultktKLuM" {
  cidr_block = var.vpc_ip_range
}

resource "alicloud_vswitch" "defaultTQWN3k" {
  vpc_id     = alicloud_vpc.defaultktKLuM.id
  zone_id    = var.zone_id_i
  cidr_block = var.vsw_ip_range_i
}


`, name)
}

// Test ClickHouse EnterpriseDBCluster. <<< Resource test cases, automatically generated.

// Test ClickHouse EnterpriseDbCluster. >>> Resource test cases, automatically generated.
// Case 线上-CK企业版-基本资源-多可用区1-多属性支持2-iam测试账号 12419
func TestAccAliCloudClickHouseEnterpriseDbCluster_basic12419(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_click_house_enterprise_db_cluster.default"
	ra := resourceAttrInit(resourceId, AlicloudClickHouseEnterpriseDbClusterMap12419)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ClickHouseServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeClickHouseEnterpriseDbCluster")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccclickhouse%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudClickHouseEnterpriseDbClusterBasicDependence12419)
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
					"zone_id":    "${var.zone_id_i}",
					"vpc_id":     "${alicloud_vpc.defaultktKLuM.id}",
					"scale_min":  "8",
					"scale_max":  "8",
					"vswitch_id": "${alicloud_vswitch.defaultTQWN3k.id}",
					"multi_zones": []map[string]interface{}{
						{
							"vswitch_ids": []string{
								"${alicloud_vswitch.defaultTQWN3k.id}"},
							"zone_id": "${var.zone_id_i}",
						},
						{
							"vswitch_ids": []string{
								"${alicloud_vswitch.defaultylyLu8.id}"},
							"zone_id": "${var.zone_id_l}",
						},
						{
							"vswitch_ids": []string{
								"${alicloud_vswitch.defaultRNbPh8.id}"},
							"zone_id": "${var.zone_id_k}",
						},
					},
					"node_scale_min":    "4",
					"node_scale_max":    "4",
					"node_count":        "2",
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"description":       "test-create",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"zone_id":           CHECKSET,
						"vpc_id":            CHECKSET,
						"scale_min":         CHECKSET,
						"scale_max":         CHECKSET,
						"vswitch_id":        CHECKSET,
						"multi_zones.#":     "3",
						"node_scale_min":    CHECKSET,
						"node_scale_max":    CHECKSET,
						"node_count":        CHECKSET,
						"resource_group_id": CHECKSET,
						"description":       "test-create",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.1}",
					"description":       "test-update-1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id": CHECKSET,
						"description":       "test-update-1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"scale_min":      "24",
					"scale_max":      "24",
					"node_scale_min": "8",
					"node_scale_max": "8",
					"node_count":     "3",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"scale_min":      CHECKSET,
						"scale_max":      CHECKSET,
						"node_scale_min": "8",
						"node_scale_max": "8",
						"node_count":     "3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "test-update-2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "test-update-2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
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

var AlicloudClickHouseEnterpriseDbClusterMap12419 = map[string]string{
	"engine_minor_version":  CHECKSET,
	"category":              CHECKSET,
	"instance_network_type": CHECKSET,
	"endpoints.#":           CHECKSET,
	"storage_quota":         CHECKSET,
	"computing_group_ids.#": CHECKSET,
	"status":                CHECKSET,
	"storage_type":          CHECKSET,
	"create_time":           CHECKSET,
	"storage_size":          CHECKSET,
	"charge_type":           CHECKSET,
	"region_id":             CHECKSET,
}

func AlicloudClickHouseEnterpriseDbClusterBasicDependence12419(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

variable "vsw__ip_range_i" {
  default = "172.16.9.0/24"
}

variable "region_id" {
  default = "cn-beijing"
}

variable "vpc__ip_range" {
  default = "172.16.0.0/12"
}

variable "vsw__ip_range_k" {
  default = "172.16.10.0/24"
}

variable "vsw__ip_range_l" {
  default = "172.16.11.0/24"
}

variable "resource_group_name_2" {
  default = "test-resource-group-10"
}

variable "zone_id_i" {
  default = "cn-beijing-i"
}

variable "zone_id_l" {
  default = "cn-beijing-l"
}

variable "zone_id_k" {
  default = "cn-beijing-k"
}

variable "resource_group_name" {
  default = "test-resource-group-11"
}

data "alicloud_resource_manager_resource_groups" "default" {}

resource "alicloud_vpc" "defaultktKLuM" {
  cidr_block = var.vpc__ip_range
}

resource "alicloud_vswitch" "defaultTQWN3k" {
  vpc_id     = alicloud_vpc.defaultktKLuM.id
  zone_id    = var.zone_id_i
  cidr_block = var.vsw__ip_range_i
}

resource "alicloud_vswitch" "defaultylyLu8" {
  vpc_id     = alicloud_vpc.defaultktKLuM.id
  zone_id    = var.zone_id_l
  cidr_block = var.vsw__ip_range_l
}

resource "alicloud_vswitch" "defaultRNbPh8" {
  vpc_id     = alicloud_vpc.defaultktKLuM.id
  zone_id    = var.zone_id_k
  cidr_block = var.vsw__ip_range_k
}


`, name)
}

// Case CK企业版-基本资源-多可用区1-线上 10560
func TestAccAliCloudClickHouseEnterpriseDbCluster_basic10560(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_click_house_enterprise_db_cluster.default"
	ra := resourceAttrInit(resourceId, AlicloudClickHouseEnterpriseDbClusterMap10560)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ClickHouseServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeClickHouseEnterpriseDbCluster")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccclickhouse%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudClickHouseEnterpriseDbClusterBasicDependence10560)
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
					"zone_id":    "${var.zone_id_i}",
					"vpc_id":     "${alicloud_vpc.defaultktKLuM.id}",
					"scale_min":  "8",
					"scale_max":  "16",
					"vswitch_id": "${alicloud_vswitch.defaultTQWN3k.id}",
					"multi_zones": []map[string]interface{}{
						{
							"vswitch_ids": []string{
								"${alicloud_vswitch.defaultTQWN3k.id}"},
							"zone_id": "${var.zone_id_i}",
						},
						{
							"vswitch_ids": []string{
								"${alicloud_vswitch.defaultylyLu8.id}"},
							"zone_id": "${var.zone_id_l}",
						},
						{
							"vswitch_ids": []string{
								"${alicloud_vswitch.defaultRNbPh8.id}"},
							"zone_id": "${var.zone_id_k}",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"zone_id":       CHECKSET,
						"vpc_id":        CHECKSET,
						"scale_min":     CHECKSET,
						"scale_max":     CHECKSET,
						"vswitch_id":    CHECKSET,
						"multi_zones.#": "3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"scale_min": "32",
					"scale_max": "64",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"scale_min": CHECKSET,
						"scale_max": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
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

var AlicloudClickHouseEnterpriseDbClusterMap10560 = map[string]string{
	"engine_minor_version":  CHECKSET,
	"category":              CHECKSET,
	"instance_network_type": CHECKSET,
	"endpoints.#":           CHECKSET,
	"storage_quota":         CHECKSET,
	"computing_group_ids.#": CHECKSET,
	"status":                CHECKSET,
	"storage_type":          CHECKSET,
	"create_time":           CHECKSET,
	"storage_size":          CHECKSET,
	"charge_type":           CHECKSET,
	"region_id":             CHECKSET,
}

func AlicloudClickHouseEnterpriseDbClusterBasicDependence10560(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

variable "vsw__ip_range_i" {
  default = "172.16.1.0/24"
}

variable "region_id" {
  default = "cn-beijing"
}

variable "vpc__ip_range" {
  default = "172.16.0.0/12"
}

variable "vsw__ip_range_k" {
  default = "172.16.3.0/24"
}

variable "vsw__ip_range_l" {
  default = "172.16.2.0/24"
}

variable "zone_id_i" {
  default = "cn-beijing-i"
}

variable "zone_id_l" {
  default = "cn-beijing-l"
}

variable "zone_id_k" {
  default = "cn-beijing-k"
}

resource "alicloud_vpc" "defaultktKLuM" {
  cidr_block = var.vpc__ip_range
}

resource "alicloud_vswitch" "defaultTQWN3k" {
  vpc_id     = alicloud_vpc.defaultktKLuM.id
  zone_id    = var.zone_id_i
  cidr_block = var.vsw__ip_range_i
}

resource "alicloud_vswitch" "defaultylyLu8" {
  vpc_id     = alicloud_vpc.defaultktKLuM.id
  zone_id    = var.zone_id_l
  cidr_block = var.vsw__ip_range_l
}

resource "alicloud_vswitch" "defaultRNbPh8" {
  vpc_id     = alicloud_vpc.defaultktKLuM.id
  zone_id    = var.zone_id_k
  cidr_block = var.vsw__ip_range_k
}


`, name)
}

// Case CK企业版-基本资源-实例Id_网络_CCU_单可用区 10226
func TestAccAliCloudClickHouseEnterpriseDbCluster_basic10226(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_click_house_enterprise_db_cluster.default"
	ra := resourceAttrInit(resourceId, AlicloudClickHouseEnterpriseDbClusterMap10226)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ClickHouseServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeClickHouseEnterpriseDbCluster")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccclickhouse%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudClickHouseEnterpriseDbClusterBasicDependence10226)
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
					"zone_id":    "${var.zone_id_i}",
					"vpc_id":     "${alicloud_vpc.defaultktKLuM.id}",
					"scale_min":  "8",
					"scale_max":  "16",
					"vswitch_id": "${alicloud_vswitch.defaultTQWN3k.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"zone_id":    CHECKSET,
						"vpc_id":     CHECKSET,
						"scale_min":  CHECKSET,
						"scale_max":  CHECKSET,
						"vswitch_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"scale_min": "32",
					"scale_max": "64",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"scale_min": CHECKSET,
						"scale_max": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
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

var AlicloudClickHouseEnterpriseDbClusterMap10226 = map[string]string{
	"engine_minor_version":  CHECKSET,
	"category":              CHECKSET,
	"instance_network_type": CHECKSET,
	"endpoints.#":           CHECKSET,
	"storage_quota":         CHECKSET,
	"computing_group_ids.#": CHECKSET,
	"status":                CHECKSET,
	"storage_type":          CHECKSET,
	"create_time":           CHECKSET,
	"storage_size":          CHECKSET,
	"charge_type":           CHECKSET,
	"region_id":             CHECKSET,
}

func AlicloudClickHouseEnterpriseDbClusterBasicDependence10226(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

variable "vsw__ip_range_i" {
  default = "172.16.1.0/24"
}

variable "region_id" {
  default = "cn-beijing"
}

variable "vpc__ip_range" {
  default = "172.16.0.0/12"
}

variable "zone_id_i" {
  default = "cn-beijing-i"
}

resource "alicloud_vpc" "defaultktKLuM" {
  cidr_block = var.vpc__ip_range
}

resource "alicloud_vswitch" "defaultTQWN3k" {
  vpc_id     = alicloud_vpc.defaultktKLuM.id
  zone_id    = var.zone_id_i
  cidr_block = var.vsw__ip_range_i
}


`, name)
}

// Test ClickHouse EnterpriseDbCluster. <<< Resource test cases, automatically generated.

// TestAccAliCloudClickHouseEnterpriseDBClusterReorderMultiZonesPrimaryFirst
// locks the multi-zone ordering helper. multi_zones is a schema.TypeSet, whose
// Set.List() iteration order is hash-based and does not preserve the HCL
// declaration order. For multi-zone deployments the top-level zone_id is not
// forwarded to CreateDBInstance, but ClickHouse still treats MultiZone[0] as
// the primary zone, so when the top-level zone_id is set the helper normalizes
// the matching entry to index 0 to keep the server-selected primary zone
// aligned with the configuration and avoid a state drift / ForceNew
// replacement loop across plans.
func TestAccAliCloudClickHouseEnterpriseDBClusterReorderMultiZonesPrimaryFirst(t *testing.T) {
	primaryZone := "cn-beijing-i"

	// Case 1: primary zone is not at index 0 (simulates the hash-sorted
	// schema.TypeSet order). The helper must move it to index 0 and preserve
	// the relative order of the remaining entries.
	input := []interface{}{
		map[string]interface{}{"ZoneId": "cn-beijing-l", "VSwitchIds": []interface{}{"vsw-l"}},
		map[string]interface{}{"ZoneId": "cn-beijing-k", "VSwitchIds": []interface{}{"vsw-k"}},
		map[string]interface{}{"ZoneId": primaryZone, "VSwitchIds": []interface{}{"vsw-i"}},
	}
	got, err := reorderMultiZonesPrimaryFirst(input, primaryZone)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(got))
	}
	if got[0].(map[string]interface{})["ZoneId"] != primaryZone {
		t.Fatalf("expected primary zone %q at index 0, got %v", primaryZone, got[0].(map[string]interface{})["ZoneId"])
	}
	if got[1].(map[string]interface{})["ZoneId"] != "cn-beijing-l" ||
		got[2].(map[string]interface{})["ZoneId"] != "cn-beijing-k" {
		t.Fatalf("expected remaining order [cn-beijing-l, cn-beijing-k], got %v / %v", got[1], got[2])
	}

	// Case 2: primary zone already at index 0 -> array returned unchanged.
	input2 := []interface{}{
		map[string]interface{}{"ZoneId": primaryZone, "VSwitchIds": []interface{}{"vsw-i"}},
		map[string]interface{}{"ZoneId": "cn-beijing-l", "VSwitchIds": []interface{}{"vsw-l"}},
	}
	got2, err := reorderMultiZonesPrimaryFirst(input2, primaryZone)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got2[0].(map[string]interface{})["ZoneId"] != primaryZone {
		t.Fatalf("expected primary zone at index 0, got %v", got2[0])
	}

	// Case 3: empty primary zone id -> unchanged (contract not enforced client-side).
	input3 := []interface{}{
		map[string]interface{}{"ZoneId": "cn-beijing-l", "VSwitchIds": []interface{}{"vsw-l"}},
	}
	got3, err := reorderMultiZonesPrimaryFirst(input3, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got3) != 1 || got3[0].(map[string]interface{})["ZoneId"] != "cn-beijing-l" {
		t.Fatalf("expected unchanged array for empty primary zone, got %v", got3)
	}

	// Case 4: primary zone set but absent from multi_zones -> explicit error.
	input4 := []interface{}{
		map[string]interface{}{"ZoneId": "cn-beijing-l", "VSwitchIds": []interface{}{"vsw-l"}},
	}
	if _, err := reorderMultiZonesPrimaryFirst(input4, primaryZone); err == nil {
		t.Fatal("expected error when primary zone is absent from multi_zones, got nil")
	}

	// Case 5: end-to-end through a real schema.TypeSet to confirm the helper
	// corrects the hash-based Set.List() order regardless of HCL declaration
	// order. This reproduces the original bug shape: a non-primary zone sorted
	// to MultiZone[0] must be moved back so CreateDBInstance accepts the request.
	multiZoneSchema := resourceAliCloudClickHouseEnterpriseDbCluster().Schema["multi_zones"]
	set := schema.NewSet(schema.HashResource(multiZoneSchema.Elem.(*schema.Resource)), []interface{}{
		map[string]interface{}{
			"zone_id":     "cn-beijing-l",
			"vswitch_ids": schema.NewSet(schema.HashString, []interface{}{"vsw-l"}),
		},
		map[string]interface{}{
			"zone_id":     "cn-beijing-k",
			"vswitch_ids": schema.NewSet(schema.HashString, []interface{}{"vsw-k"}),
		},
		map[string]interface{}{
			"zone_id":     primaryZone,
			"vswitch_ids": schema.NewSet(schema.HashString, []interface{}{"vsw-i"}),
		},
	})
	multiZoneMapsArray := make([]interface{}, 0)
	for _, dataLoop1 := range convertToInterfaceArray(set) {
		dataLoop1Tmp := dataLoop1.(map[string]interface{})
		dataLoop1Map := make(map[string]interface{})
		dataLoop1Map["VSwitchIds"] = convertToInterfaceArray(dataLoop1Tmp["vswitch_ids"])
		dataLoop1Map["ZoneId"] = dataLoop1Tmp["zone_id"]
		multiZoneMapsArray = append(multiZoneMapsArray, dataLoop1Map)
	}
	got5, err := reorderMultiZonesPrimaryFirst(multiZoneMapsArray, primaryZone)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got5) != 3 {
		t.Fatalf("expected 3 entries from Set, got %d", len(got5))
	}
	if got5[0].(map[string]interface{})["ZoneId"] != primaryZone {
		t.Fatalf("expected primary zone %q at index 0 after Set normalization, got %v", primaryZone, got5[0].(map[string]interface{})["ZoneId"])
	}
}

// TestAccAliCloudClickHouseEnterpriseDBCluster_multiZonesPrimaryZoneStable
// locks the fix end-to-end: a multi-zone deployment where the primary zone
// (top-level zone_id) is declared in multi_zones but not as the first block.
// multi_zones is a schema.TypeSet, whose iteration order is hash-based. The
// fix has two parts: (1) the top-level zone_id/vswitch_id are not forwarded
// to CreateDBInstance for multi-zone deployments, so the API no longer rejects
// the request with InvalidZoneId.InconsistentWithMultiZone; (2) when the
// top-level zone_id is set, the matching multi_zones entry is normalized to
// MultiZone[0] so the server-selected primary zone aligns with the
// configuration and the state stays stable across plans. The top-level
// zone_id/vswitch_id remain in the config to keep state stable and to exercise
// the primary-zone normalization end-to-end.
func TestAccAliCloudClickHouseEnterpriseDBCluster_multiZonesPrimaryZoneStable(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_click_house_enterprise_db_cluster.default"
	ra := resourceAttrInit(resourceId, AlicloudClickHouseEnterpriseDBClusterMapPrimaryZoneStable)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ClickHouseServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeClickHouseEnterpriseDbCluster")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccclickhouse%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudClickHouseEnterpriseDbClusterBasicDependence10560)
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
					"zone_id":    "${var.zone_id_i}",
					"vpc_id":     "${alicloud_vpc.defaultktKLuM.id}",
					"scale_min":  "8",
					"scale_max":  "16",
					"vswitch_id": "${alicloud_vswitch.defaultTQWN3k.id}",
					// The primary zone (zone_id_i) is intentionally declared as
					// the LAST multi_zones block. Because multi_zones is a
					// TypeSet, its serialization order is hash-based; the fix
					// must move the primary zone to MultiZone[0] so that
					// CreateDBInstance accepts the request.
					"multi_zones": []map[string]interface{}{
						{
							"vswitch_ids": []string{
								"${alicloud_vswitch.defaultylyLu8.id}"},
							"zone_id": "${var.zone_id_l}",
						},
						{
							"vswitch_ids": []string{
								"${alicloud_vswitch.defaultRNbPh8.id}"},
							"zone_id": "${var.zone_id_k}",
						},
						{
							"vswitch_ids": []string{
								"${alicloud_vswitch.defaultTQWN3k.id}"},
							"zone_id": "${var.zone_id_i}",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"zone_id":       CHECKSET,
						"vpc_id":        CHECKSET,
						"scale_min":     CHECKSET,
						"scale_max":     CHECKSET,
						"vswitch_id":    CHECKSET,
						"multi_zones.#": "3",
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

var AlicloudClickHouseEnterpriseDBClusterMapPrimaryZoneStable = map[string]string{
	"status":      CHECKSET,
	"create_time": CHECKSET,
	"region_id":   CHECKSET,
}

package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
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
	"maintain_end_time":     CHECKSET,
	"object_store_size":     CHECKSET,
	"disabled_ports":        CHECKSET,
	"storage_quota":         CHECKSET,
	"computing_group_ids.#": CHECKSET,
	"maintain_start_time":   CHECKSET,
	"status":                CHECKSET,
	"engine_version":        CHECKSET,
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
	"maintain_end_time":     CHECKSET,
	"object_store_size":     CHECKSET,
	"disabled_ports":        CHECKSET,
	"storage_quota":         CHECKSET,
	"computing_group_ids.#": CHECKSET,
	"maintain_start_time":   CHECKSET,
	"status":                CHECKSET,
	"engine_version":        CHECKSET,
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
	"maintain_end_time":     CHECKSET,
	"object_store_size":     CHECKSET,
	"disabled_ports":        CHECKSET,
	"storage_quota":         CHECKSET,
	"computing_group_ids.#": CHECKSET,
	"maintain_start_time":   CHECKSET,
	"status":                CHECKSET,
	"engine_version":        CHECKSET,
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

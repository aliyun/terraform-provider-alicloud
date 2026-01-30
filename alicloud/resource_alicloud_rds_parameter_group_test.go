package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudRdsParameterGroup_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_rds_parameter_group.default"
	ra := resourceAttrInit(resourceId, AlicloudRdsParameterGroupMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &RdsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeRdsParameterGroup")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf_testAccAlicloudRdsParameterGroup%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudRdsParameterGroupBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"engine":         "mysql",
					"engine_version": `5.7`,
					"param_detail": []map[string]interface{}{
						{
							"param_name":  "back_log",
							"param_value": `3000`,
						},
						{
							"param_name":  "wait_timeout",
							"param_value": `86400`,
						},
					},
					"parameter_group_name": "${var.name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"engine":               "mysql",
						"engine_version":       "5.7",
						"param_detail.#":       "2",
						"parameter_group_name": name,
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
					"param_detail": []map[string]interface{}{
						{
							"param_name":  "back_log",
							"param_value": `4000`,
						},
						{
							"param_name":  "wait_timeout",
							"param_value": `86460`,
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"param_detail.#": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"parameter_group_desc": "update_test",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"parameter_group_desc": "update_test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"parameter_group_name": name + "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"parameter_group_name": name + "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"param_detail": []map[string]interface{}{
						{
							"param_name":  "back_log",
							"param_value": `3000`,
						},
						{
							"param_name":  "wait_timeout",
							"param_value": `86400`,
						},
					},
					"parameter_group_desc": "test",
					"parameter_group_name": "${var.name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"param_detail.#":       "2",
						"parameter_group_desc": "test",
						"parameter_group_name": name,
					}),
				),
			},
		},
	})
}

func TestAccAlicloudRdsParameterGroup_basic1(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_rds_parameter_group.default"
	ra := resourceAttrInit(resourceId, AlicloudRdsParameterGroupMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &RdsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeRdsParameterGroup")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf_testAccAlicloudRdsParameterGroup%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudRdsParameterGroupBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"engine":         "mysql",
					"engine_version": `5.7`,
					"param_detail": []map[string]interface{}{
						{
							"param_name":  "back_log",
							"param_value": `3000`,
						},
						{
							"param_name":  "wait_timeout",
							"param_value": `86400`,
						},
					},
					"parameter_group_name": "${var.name}",
					"parameter_group_desc": "update_test",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"engine":               "mysql",
						"engine_version":       "5.7",
						"param_detail.#":       "2",
						"parameter_group_name": name,
						"parameter_group_desc": "update_test",
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

var AlicloudRdsParameterGroupMap0 = map[string]string{}

func AlicloudRdsParameterGroupBasicDependence0(name string) string {
	return fmt.Sprintf(`
variable "name" {
			default = "%s"
		}
`, name)
}

func TestAccAlicloudRdsParameterGroupPostgreSQL(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_rds_parameter_group.default"
	ra := resourceAttrInit(resourceId, AlicloudRdsParameterGroupMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &RdsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeRdsParameterGroup")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf_testAccAlicloudRdsParameterGroup%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudRdsParameterGroupBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"engine":         "PostgreSQL",
					"engine_version": "11.0",
					"param_detail": []map[string]interface{}{
						{
							"param_name":  "enable_sort",
							"param_value": "off",
						},
						{
							"param_name":  "geqo_seed",
							"param_value": "0.1",
						},
					},
					"parameter_group_name": "${var.name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"engine":               "PostgreSQL",
						"engine_version":       "11.0",
						"param_detail.#":       "2",
						"parameter_group_name": name,
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
					"param_detail": []map[string]interface{}{
						{
							"param_name":  "enable_sort",
							"param_value": "on",
						},
						{
							"param_name":  "geqo_seed",
							"param_value": "0.2",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"param_detail.#": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"parameter_group_desc": "update_test",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"parameter_group_desc": "update_test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"parameter_group_name": name + "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"parameter_group_name": name + "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"param_detail": []map[string]interface{}{
						{
							"param_name":  "enable_sort",
							"param_value": "off",
						},
						{
							"param_name":  "geqo_seed",
							"param_value": "0.1",
						},
					},
					"parameter_group_desc": "test",
					"parameter_group_name": "${var.name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"param_detail.#":       "2",
						"parameter_group_desc": "test",
						"parameter_group_name": name,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"modify_mode", "resource_group_id", "parameter_detail", "param_detail"},
			},
		},
	})
}

// Test Rds ParameterGroup. >>> Resource test cases, automatically generated.
// Case RDS/参数模板测试 12250
func TestAccAliCloudRdsParameterGroup_basic12250(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_rds_parameter_group.default"
	ra := resourceAttrInit(resourceId, AlicloudRdsParameterGroupMap12250)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &RdsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeRdsParameterGroup")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccrds%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudRdsParameterGroupBasicDependence12250)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"param_detail": []map[string]interface{}{
						{
							"param_name":  "loose_performance_schema_max_index_stat",
							"param_value": "1000",
						},
					},
					"engine_version":       "8.0",
					"parameter_group_name": name,
					"parameter_group_desc": "test001",
					"engine":               "mysql",
					"resource_group_id":    "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"param_detail.#":       "1",
						"engine_version":       CHECKSET,
						"parameter_group_name": name,
						"parameter_group_desc": "test001",
						"engine":               "mysql",
						"resource_group_id":    CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"param_detail": []map[string]interface{}{
						{
							"param_name":  "loose_performance_schema_max_index_stat",
							"param_value": "2000",
						},
						{
							"param_name":  "bulk_insert_buffer_size",
							"param_value": "30",
						},
					},
					"parameter_group_name": name + "_update",
					"parameter_group_desc": "test02",
					"resource_group_id":    "${data.alicloud_resource_manager_resource_groups.default.ids.1}",
					"modify_mode":          "Individual",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"param_detail.#":       "2",
						"parameter_group_name": name + "_update",
						"parameter_group_desc": "test02",
						"resource_group_id":    CHECKSET,
						"modify_mode":          "Individual",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"modify_mode", "resource_group_id"},
			},
		},
	})
}

var AlicloudRdsParameterGroupMap12250 = map[string]string{
	"create_time": CHECKSET,
}

func AlicloudRdsParameterGroupBasicDependence12250(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

variable "test_region_id" {
  default = "cn-beijing"
}

variable "test_zone_id" {
  default = "cn-beijing-h"
}

data "alicloud_resource_manager_resource_groups" "default" {}


`, name)
}

// Case RDS/参数模板测试 12067
func TestAccAliCloudRdsParameterGroup_basic12067(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_rds_parameter_group.default"
	ra := resourceAttrInit(resourceId, AlicloudRdsParameterGroupMap12067)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &RdsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeRdsParameterGroup")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccrds%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudRdsParameterGroupBasicDependence12067)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"param_detail": []map[string]interface{}{
						{
							"param_name":  "loose_performance_schema_max_index_stat",
							"param_value": "1000",
						},
					},
					"engine_version":       "8.0",
					"parameter_group_name": name,
					"parameter_group_desc": "test001",
					"engine":               "mysql",
					"resource_group_id":    "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"param_detail.#":       "1",
						"engine_version":       CHECKSET,
						"parameter_group_name": name,
						"parameter_group_desc": "test001",
						"engine":               "mysql",
						"resource_group_id":    CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"param_detail": []map[string]interface{}{
						{
							"param_name":  "loose_performance_schema_max_index_stat",
							"param_value": "2000",
						},
						{
							"param_name":  "bulk_insert_buffer_size",
							"param_value": "30",
						},
					},
					"parameter_group_name": name + "_update",
					"parameter_group_desc": "test02",
					"resource_group_id":    "${data.alicloud_resource_manager_resource_groups.default.ids.1}",
					"modify_mode":          "Individual",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"parameter_detail.#":   "2",
						"parameter_group_name": name + "_update",
						"parameter_group_desc": "test02",
						"resource_group_id":    CHECKSET,
						"modify_mode":          "Individual",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"modify_mode", "resource_group_id"},
			},
		},
	})
}

var AlicloudRdsParameterGroupMap12067 = map[string]string{
	"create_time": CHECKSET,
}

func AlicloudRdsParameterGroupBasicDependence12067(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

variable "test_region_id" {
  default = "cn-beijing"
}

variable "test_zone_id" {
  default = "cn-beijing-h"
}

data "alicloud_resource_manager_resource_groups" "default" {}


`, name)
}

// Test Rds ParameterGroup. <<< Resource test cases, automatically generated.

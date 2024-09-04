package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test Gpdb JdbcDataSource. >>> Resource test cases, automatically generated.
// Case 创建JDBC数据源_资源依赖_NEW_ACCOUNT 7592
func TestAccAliCloudGpdbJdbcDataSource_basic7592(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_gpdb_jdbc_data_source.default"
	ra := resourceAttrInit(resourceId, AlicloudGpdbJdbcDataSourceMap7592)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &GpdbServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeGpdbJdbcDataSource")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1, 9999)
	name := fmt.Sprintf("tf_testgpdb%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudGpdbJdbcDataSourceBasicDependence7592)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-beijing"})
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"jdbc_password":           "test_005",
					"data_source_name":        "${alicloud_gpdb_external_data_service.defaultRXkfKL.service_name}",
					"data_source_type":        "mysql",
					"jdbc_user_name":          "test_005",
					"jdbc_connection_string":  "jdbc:mysql://rm-2ze327yr44c61183c.mysql.rds.aliyuncs.com:3306/test_005",
					"data_source_description": "mytest",
					"db_instance_id":          "${alicloud_gpdb_instance.defaulttuqTmM.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"jdbc_password":           "test_005",
						"data_source_name":        CHECKSET,
						"data_source_type":        "mysql",
						"jdbc_user_name":          "test_005",
						"jdbc_connection_string":  "jdbc:mysql://rm-2ze327yr44c61183c.mysql.rds.aliyuncs.com:3306/test_005",
						"data_source_description": "mytest",
						"db_instance_id":          CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"jdbc_password":           "test_006",
					"data_source_type":        "postgresql",
					"jdbc_user_name":          "test_006",
					"jdbc_connection_string":  "jdbc:mysql://rm-2ze327yr44c61182c.mysql.rds.aliyuncs.com:3306/test_006",
					"data_source_description": "test2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"jdbc_password":           "test_006",
						"data_source_type":        "postgresql",
						"jdbc_user_name":          "test_006",
						"jdbc_connection_string":  "jdbc:mysql://rm-2ze327yr44c61182c.mysql.rds.aliyuncs.com:3306/test_006",
						"data_source_description": "test2",
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

var AlicloudGpdbJdbcDataSourceMap7592 = map[string]string{
	"status":         CHECKSET,
	"create_time":    CHECKSET,
	"data_source_id": CHECKSET,
}

func AlicloudGpdbJdbcDataSourceBasicDependence7592(name string) string {
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
  zone_id = "cn-beijing-h"
}

resource "alicloud_gpdb_instance" "defaulttuqTmM" {
  instance_spec              = "2C8G"
  description                = "创建依赖的Local实例"
  seg_node_num               = "2"
  seg_storage_type           = "cloud_essd"
  instance_network_type      = "VPC"
  payment_type               = "PayAsYouGo"
  ssl_enabled                = "0"
  engine_version             = "6.0"
  zone_id                    = "cn-beijing-h"
  vswitch_id                 = data.alicloud_vswitches.default.ids[0]
  storage_size               = "50"
  master_cu                  = "4"
  vpc_id                     = data.alicloud_vpcs.default.ids.0
  db_instance_mode           = "StorageElastic"
  engine                     = "gpdb"
  db_instance_category       = "Basic"
}

resource "alicloud_gpdb_account" "defaultsF2fsl" {
  db_instance_id   = alicloud_gpdb_instance.defaulttuqTmM.id
  account_name     = format("%%s1", var.name)
  account_password = "test_005"
}

resource "alicloud_gpdb_account" "default8txVNo" {
  account_name     = format("%%s2", var.name)
  account_password = "test_006"
  db_instance_id   = alicloud_gpdb_instance.defaulttuqTmM.id
}

resource "alicloud_gpdb_external_data_service" "defaultRXkfKL" {
  service_name        = var.name
  db_instance_id      = alicloud_gpdb_instance.defaulttuqTmM.id
  service_description = "mytest"
  service_spec        = "8"
}


`, name)
}

// Case 创建JDBC数据源_资源依赖_case 6965
func TestAccAliCloudGpdbJdbcDataSource_basic6965(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_gpdb_jdbc_data_source.default"
	ra := resourceAttrInit(resourceId, AlicloudGpdbJdbcDataSourceMap6965)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &GpdbServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeGpdbJdbcDataSource")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1, 9999)
	name := fmt.Sprintf("tf_testgpdb%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudGpdbJdbcDataSourceBasicDependence6965)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-beijing"})
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"jdbc_password":           "test_001",
					"data_source_name":        "${alicloud_gpdb_external_data_service.defaultRXkfKL.service_name}",
					"data_source_type":        "mysql",
					"jdbc_user_name":          "test_001",
					"jdbc_connection_string":  "jdbc:mysql://rm-2ze327yr44c61183c.mysql.rds.aliyuncs.com:3306/test_001",
					"data_source_description": "mytest",
					"db_instance_id":          "${alicloud_gpdb_instance.defaulttuqTmM.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"jdbc_password":           "test_001",
						"data_source_name":        CHECKSET,
						"data_source_type":        "mysql",
						"jdbc_user_name":          "test_001",
						"jdbc_connection_string":  "jdbc:mysql://rm-2ze327yr44c61183c.mysql.rds.aliyuncs.com:3306/test_001",
						"data_source_description": "mytest",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"data_source_type":        "postgresql",
					"jdbc_connection_string":  "jdbc:mysql://rm-2ze327yr44c61182c.mysql.rds.aliyuncs.com:3306/test_001",
					"data_source_description": "test2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"data_source_type":        "postgresql",
						"jdbc_connection_string":  "jdbc:mysql://rm-2ze327yr44c61182c.mysql.rds.aliyuncs.com:3306/test_001",
						"data_source_description": "test2",
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

var AlicloudGpdbJdbcDataSourceMap6965 = map[string]string{
	"status":         CHECKSET,
	"create_time":    CHECKSET,
	"data_source_id": CHECKSET,
}

func AlicloudGpdbJdbcDataSourceBasicDependence6965(name string) string {
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
  zone_id = "cn-beijing-h"
}

resource "alicloud_gpdb_instance" "defaulttuqTmM" {
  instance_spec              = "2C8G"
  description                = "创建依赖的Local实例"
  seg_node_num               = "2"
  seg_storage_type           = "cloud_essd"
  instance_network_type      = "VPC"
  payment_type               = "PayAsYouGo"
  ssl_enabled                = "0"
  engine_version             = "6.0"
  zone_id                    = "cn-beijing-h"
  vswitch_id                 = data.alicloud_vswitches.default.ids[0]
  storage_size               = "50"
  master_cu                  = "4"
  vpc_id                     = data.alicloud_vpcs.default.ids.0
  db_instance_mode           = "StorageElastic"
  engine                     = "gpdb"
  db_instance_category       = "Basic"
}

resource "alicloud_gpdb_account" "defaultsk1eaS" {
  account_description = "test_001"
  db_instance_id      = alicloud_gpdb_instance.defaulttuqTmM.id
  account_name        = format("%%s1", var.name)
  account_password    = "test_001"
}

resource "alicloud_gpdb_external_data_service" "defaultRXkfKL" {
  service_name        = var.name
  db_instance_id      = alicloud_gpdb_instance.defaulttuqTmM.id
  service_description = "mytest"
  service_spec        = "8"
}


`, name)
}

// Test Gpdb JdbcDataSource. <<< Resource test cases, automatically generated.

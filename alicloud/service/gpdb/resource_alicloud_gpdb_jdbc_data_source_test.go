package gpdb_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/gpdb"

	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Test Gpdb JdbcDataSource. >>> Resource test cases, automatically generated.
// Case 创建JDBC数据源_资源依赖_NEW_ACCOUNT 7592
func TestAccAliCloudGpdbJdbcDataSource_basic7592(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_gpdb_jdbc_data_source.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudGpdbJdbcDataSourceMap7592)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return gpdb.NewServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeJdbcDataSource")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccgpdb%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudGpdbJdbcDataSourceBasicDependence7592)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-beijing"})
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
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
						"data_source_name":        acctest.CHECKSET,
						"data_source_type":        "mysql",
						"jdbc_user_name":          "test_005",
						"jdbc_connection_string":  "jdbc:mysql://rm-2ze327yr44c61183c.mysql.rds.aliyuncs.com:3306/test_005",
						"data_source_description": "mytest",
						"db_instance_id":          acctest.CHECKSET,
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
	"status":         acctest.CHECKSET,
	"create_time":    acctest.CHECKSET,
	"data_source_id": acctest.CHECKSET,
}

func AlicloudGpdbJdbcDataSourceBasicDependence7592(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_gpdb_zones" "default" {}

resource "alicloud_vpc" "default4Mf0nY" {
  cidr_block = "192.168.0.0/16"
}

resource "alicloud_vswitch" "defaultwSAVpf" {
  vpc_id     = alicloud_vpc.default4Mf0nY.id
  zone_id    = data.alicloud_gpdb_zones.default.zones.0.id
  cidr_block = "192.168.1.0/24"
}

resource "alicloud_gpdb_instance" "defaulttuqTmM" {
  instance_spec              = "2C8G"
  description                = "创建依赖的Local实例"
  seg_node_num               = "2"
  seg_storage_type           = "cloud_essd"
  instance_network_type      = "VPC"
  db_instance_category       = "Basic"
  security_ip_list           = ["127.0.0.1"]
  payment_type               = "PayAsYouGo"
  ssl_enabled                = "0"
  engine_version             = "6.0"
  zone_id                    = data.alicloud_gpdb_zones.default.zones.0.id
  vswitch_id                 = alicloud_vswitch.defaultwSAVpf.id
  storage_size               = "50"
  master_cu                  = "4"
  vpc_id                     = alicloud_vpc.default4Mf0nY.id
  db_instance_mode           = "StorageElastic"
  engine                     = "gpdb"
}

resource "alicloud_gpdb_account" "defaultsF2fsl" {
  db_instance_id   = alicloud_gpdb_instance.defaulttuqTmM.id
  account_name     = format("%%s1", var.name)
  account_password = "test_005"
  account_type     = "Normal"
}

resource "alicloud_gpdb_account" "default8txVNo" {
  account_name     = format("%%s2", var.name)
  account_password = "test_006"
  db_instance_id   = alicloud_gpdb_instance.defaulttuqTmM.id
  account_type     = "Normal"
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
	ra := acctest.ResourceAttrInit(resourceId, AlicloudGpdbJdbcDataSourceMap6965)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return gpdb.NewServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeJdbcDataSource")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccgpdb%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudGpdbJdbcDataSourceBasicDependence6965)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-beijing"})
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
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
						"data_source_name":        acctest.CHECKSET,
						"data_source_type":        "mysql",
						"jdbc_user_name":          "test_001",
						"jdbc_connection_string":  "jdbc:mysql://rm-2ze327yr44c61183c.mysql.rds.aliyuncs.com:3306/test_001",
						"data_source_description": "mytest",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"jdbc_password":           "test_002",
					"data_source_type":        "postgresql",
					"jdbc_connection_string":  "jdbc:mysql://rm-2ze327yr44c61182c.mysql.rds.aliyuncs.com:3306/test_002",
					"jdbc_user_name":          "test_002",
					"data_source_description": "test2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"jdbc_password":           "test_002",
						"data_source_type":        "postgresql",
						"jdbc_connection_string":  "jdbc:mysql://rm-2ze327yr44c61182c.mysql.rds.aliyuncs.com:3306/test_002",
						"jdbc_user_name":          "test_002",
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
	"status":         acctest.CHECKSET,
	"create_time":    acctest.CHECKSET,
	"data_source_id": acctest.CHECKSET,
}

func AlicloudGpdbJdbcDataSourceBasicDependence6965(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_gpdb_zones" "default" {}

resource "alicloud_vpc" "default4Mf0nY" {
  cidr_block = "192.168.0.0/16"
}

resource "alicloud_vswitch" "defaultwSAVpf" {
  vpc_id     = alicloud_vpc.default4Mf0nY.id
  zone_id    = data.alicloud_gpdb_zones.default.zones.0.id
  cidr_block = "192.168.1.0/24"
}

resource "alicloud_gpdb_instance" "defaulttuqTmM" {
  instance_spec              = "2C8G"
  description                = "创建依赖的Local实例"
  seg_node_num               = "2"
  seg_storage_type           = "cloud_essd"
  instance_network_type      = "VPC"
  db_instance_category       = "Basic"
  security_ip_list           = ["127.0.0.1"]
  payment_type               = "PayAsYouGo"
  ssl_enabled                = "0"
  engine_version             = "6.0"
  zone_id                    = data.alicloud_gpdb_zones.default.zones.0.id
  vswitch_id                 = alicloud_vswitch.defaultwSAVpf.id
  storage_size               = "50"
  master_cu                  = "4"
  vpc_id                     = alicloud_vpc.default4Mf0nY.id
  db_instance_mode           = "StorageElastic"
  engine                     = "gpdb"
}

resource "alicloud_gpdb_account" "defaultsk1eaS" {
  account_description = "test_001"
  db_instance_id      = alicloud_gpdb_instance.defaulttuqTmM.id
  account_name        = format("%%s1", var.name)
  account_password    = "test_001"
  account_type        = "Normal"
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

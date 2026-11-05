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

func TestAccAliCloudGpdbRemoteADBDataSource_basic6853(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_gpdb_remote_adb_data_source.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudGpdbRemoteADBDataSourceMap6853)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return gpdb.NewServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeRemoteADBDataSource")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sgpdbremoteadbdatasource%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudGpdbRemoteADBDataSourceBasicDependence6853)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"remote_database":       "test_001",
					"manager_user_name":     "test_001",
					"user_name":             "test_001",
					"remote_db_instance_id": "${alicloud_gpdb_account.defaultwXePof.db_instance_id}",
					"local_database":        "test_001",
					"data_source_name":      "mytest",
					"user_password":         "test_001",
					"manager_user_password": "test_001",
					"local_db_instance_id":  "${alicloud_gpdb_instance.defaultEtEzMF.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"remote_database":       "test_001",
						"manager_user_name":     "test_001",
						"user_name":             "test_001",
						"remote_db_instance_id": acctest.CHECKSET,
						"local_database":        "test_001",
						"user_password":         "test_001",
						"manager_user_password": "test_001",
						"local_db_instance_id":  acctest.CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"data_source_name": "mytest",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"data_source_name": "mytest",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"data_source_name": "testDataSource2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"data_source_name": "testDataSource2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"user_name":     "test_002",
					"user_password": "test_002",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"user_name":     "test_002",
						"user_password": "test_002",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"remote_database":       "test_001",
					"manager_user_name":     "test_001",
					"user_name":             "test_001",
					"remote_db_instance_id": "${alicloud_gpdb_account.defaultwXePof.db_instance_id}",
					"local_database":        "test_001",
					"data_source_name":      "mytest",
					"user_password":         "test_001",
					"manager_user_password": "test_001",
					"local_db_instance_id":  "${alicloud_gpdb_instance.defaultEtEzMF.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"remote_database":       "test_001",
						"manager_user_name":     "test_001",
						"user_name":             "test_001",
						"remote_db_instance_id": acctest.CHECKSET,
						"local_database":        "test_001",
						"data_source_name":      "mytest",
						"user_password":         "test_001",
						"manager_user_password": "test_001",
						"local_db_instance_id":  acctest.CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"manager_user_password", "user_password"},
			},
		},
	})
}

var AlicloudGpdbRemoteADBDataSourceMap6853 = map[string]string{
	"status":                    acctest.CHECKSET,
	"remote_adb_data_source_id": acctest.CHECKSET,
}

func AlicloudGpdbRemoteADBDataSourceBasicDependence6853(name string) string {
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

resource "alicloud_gpdb_instance" "defaultEtEzMF" {
  instance_spec              = "2C8G"
  description                = "创建依赖的Local实例"
  seg_node_num               = "2"
  seg_storage_type           = "cloud_essd"
  instance_network_type      = "VPC"
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
  db_instance_category       = "Basic"
}

resource "alicloud_gpdb_instance" "defaultEY7t9t" {
  instance_spec              = "2C8G"
  description                = "创建远端依赖的实例"
  seg_node_num               = "2"
  seg_storage_type           = "cloud_essd"
  instance_network_type      = "VPC"
  db_instance_category       = "Basic"
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

resource "alicloud_gpdb_account" "default26qpEo" {
  account_description = "test_001"
  db_instance_id      = alicloud_gpdb_instance.defaultEtEzMF.id
  account_name        = "test_001"
  account_password    = "test_001"
}

resource "alicloud_gpdb_account" "defaultwXePof" {
  account_description = "test_001"
  db_instance_id      = alicloud_gpdb_instance.defaultEY7t9t.id
  account_name        = "test_001"
  account_password    = "test_001"
}

resource "alicloud_gpdb_account" "defaultTQoLwv" {
  account_description = "test_002"
  db_instance_id      = alicloud_gpdb_instance.defaultEY7t9t.id
  account_name        = "test_002"
  account_password    = "test_002"
  account_type        = "Normal"
}

`, name)
}

// Case adb2adb测试用例_yb_test 6853  raw
func TestAccAliCloudGpdbRemoteADBDataSource_basic6853_raw(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_gpdb_remote_adb_data_source.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudGpdbRemoteADBDataSourceMap6853)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return gpdb.NewServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeRemoteADBDataSource")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sgpdbremoteadbdatasource%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudGpdbRemoteADBDataSourceBasicDependence6853)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"eu-central-1"})
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"remote_database":       "test_001",
					"manager_user_name":     "test_001",
					"user_name":             "test_001",
					"remote_db_instance_id": "${alicloud_gpdb_instance.defaultEtEzMF.id}",
					"local_database":        "test_001",
					"data_source_name":      name,
					"user_password":         "test_001",
					"manager_user_password": "test_001",
					"local_db_instance_id":  "${alicloud_gpdb_instance.defaultEtEzMF.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"remote_database":       "test_001",
						"manager_user_name":     "test_001",
						"user_name":             "test_001",
						"remote_db_instance_id": acctest.CHECKSET,
						"local_database":        "test_001",
						"data_source_name":      name,
						"user_password":         "test_001",
						"manager_user_password": "test_001",
						"local_db_instance_id":  acctest.CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"data_source_name": "testDataSource2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"data_source_name": "testDataSource2",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"manager_user_password", "user_password"},
			},
		},
	})
}

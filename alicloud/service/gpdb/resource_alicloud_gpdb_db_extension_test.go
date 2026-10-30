// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
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

// Test Gpdb DbExtension. >>> Resource test cases, automatically generated.
// Case 插件 - 跳过更新版本测试 7888
func TestAccAliCloudGpdbDbExtension_basic7888(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_gpdb_db_extension.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudGpdbDbExtensionMap7888)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return gpdb.NewServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeDbExtension")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccgpdb%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudGpdbDbExtensionBasicDependence7888)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-beijing"})
			acctest.PreCheck(t)
		},
		Providers:    acctest.SDKProviders,
		CheckDestroy: rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"extension_name":    "uuid-ossp",
					"db_instance_id":    "${alicloud_gpdb_instance.defaultqsmpIy.id}",
					"database_name":     "${alicloud_gpdb_database.defaultPPmRVa.database_name}",
					"is_latest_version": true,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"extension_name":    "uuid-ossp",
						"db_instance_id":    acctest.CHECKSET,
						"database_name":     acctest.CHECKSET,
						"is_latest_version": "true",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
				// CreateExtensions always installs the newest available version, so a
				// freshly created extension already reports IsLatestVersion = true and
				// there is no API to install an older one. The false -> true transition
				// this attribute drives therefore cannot be produced in a test.
				ImportStateVerifyIgnore: []string{"is_latest_version"},
			},
		},
	})
}

var AlicloudGpdbDbExtensionMap7888 = map[string]string{
	"status":                  acctest.CHECKSET,
	"description":             acctest.CHECKSET,
	"is_install_need_restart": acctest.CHECKSET,
	"latest_version":          acctest.CHECKSET,
	"current_version":         acctest.CHECKSET,
	"extension_id":            acctest.CHECKSET,
}

func AlicloudGpdbDbExtensionBasicDependence7888(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_vpc" "defaultiYyNGW" {
  cidr_block = "192.168.0.0/16"
}

resource "alicloud_vswitch" "defaultPlruct" {
  vpc_id     = alicloud_vpc.defaultiYyNGW.id
  zone_id    = "cn-beijing-h"
  cidr_block = "192.168.1.0/24"
}

resource "alicloud_gpdb_instance" "defaultqsmpIy" {
  instance_spec         = "2C8G"
  seg_node_num          = "2"
  seg_storage_type      = "cloud_essd"
  instance_network_type = "VPC"
  db_instance_category  = "Basic"
  payment_type          = "PayAsYouGo"
  ssl_enabled           = "0"
  engine_version        = "6.0"
  engine                = "gpdb"
  zone_id               = "cn-beijing-h"
  vswitch_id            = alicloud_vswitch.defaultPlruct.id
  storage_size          = "50"
  master_cu             = "4"
  vpc_id                = alicloud_vpc.defaultiYyNGW.id
  db_instance_mode      = "StorageElastic"
}

resource "alicloud_gpdb_account" "defaultOwner" {
  account_name        = "tf_example"
  account_password    = "Example1234"
  account_description = "tf_example"
  db_instance_id      = alicloud_gpdb_instance.defaultqsmpIy.id
}

resource "alicloud_gpdb_database" "defaultPPmRVa" {
  owner              = alicloud_gpdb_account.defaultOwner.account_name
  database_name      = "seagull"
  db_instance_id     = alicloud_gpdb_instance.defaultqsmpIy.id
  character_set_name = "UTF8"
  collate            = "en_US.utf8"
  ctype              = "en_US.utf8"
}


`, name)
}

// Test Gpdb DbExtension. <<< Resource test cases, automatically generated.

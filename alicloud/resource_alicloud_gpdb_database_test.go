package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test Gpdb Database. >>> Resource test cases, automatically generated.
// Case 测试 7868
func TestAccAliCloudGpdbDatabase_basic7868(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_gpdb_database.default"
	ra := resourceAttrInit(resourceId, AlicloudGpdbDatabaseMap7868)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &GpdbServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeGpdbDatabase")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sgpdbdatabase%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudGpdbDatabaseBasicDependence7868)
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
					"character_set_name": "UTF8",
					"owner":              "adbpgadmin",
					"description":        "go-to-the-docks-for-french-fries",
					"database_name":      "seagull",
					"collate":            "en_US.utf8",
					"ctype":              "en_US.utf8",
					"db_instance_id":     "${alicloud_gpdb_instance.defaultTC08a9.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"character_set_name": "UTF8",
						"owner":              "adbpgadmin",
						"description":        "go-to-the-docks-for-french-fries",
						"database_name":      "seagull",
						"collate":            "en_US.utf8",
						"ctype":              "en_US.utf8",
						"db_instance_id":     CHECKSET,
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

var AlicloudGpdbDatabaseMap7868 = map[string]string{}

func AlicloudGpdbDatabaseBasicDependence7868(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_vpc" "default35OkxY" {
  cidr_block = "192.168.0.0/16"
}

resource "alicloud_vswitch" "defaultl8haQ3" {
  vpc_id     = alicloud_vpc.default35OkxY.id
  zone_id    = "cn-beijing-h"
  cidr_block = "192.168.1.0/24"
}

resource "alicloud_gpdb_instance" "defaultTC08a9" {
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
  vswitch_id            = alicloud_vswitch.defaultl8haQ3.id
  storage_size          = "50"
  master_cu             = "4"
  vpc_id                = alicloud_vpc.default35OkxY.id
  db_instance_mode      = "StorageElastic"
}


`, name)
}

// Test Gpdb Database. <<< Resource test cases, automatically generated.

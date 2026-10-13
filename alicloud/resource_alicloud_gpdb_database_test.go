package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
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
			testAccPreCheckWithRegions(t, true, connectivity.GPDBDBInstancePlanSupportRegions)
			testAccPreCheck(t)
		},
		IDRefreshName:     resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"character_set_name": "UTF8",
					"owner":              "${alicloud_gpdb_account.default.account_name}",
					"description":        "go-to-the-docks-for-french-fries",
					"database_name":      "seagull",
					"collate":            "en_US.utf8",
					"ctype":              "en_US.utf8",
					"db_instance_id":     "${alicloud_gpdb_instance.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"character_set_name": "UTF8",
						"owner":              CHECKSET,
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

data "alicloud_gpdb_zones" "default" {
}

data "alicloud_vpcs" "default" {
  name_regex = "^default-NODELETING$"
}

data "alicloud_vswitches" "default" {
  vpc_id  = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_gpdb_zones.default.ids.0
}

resource "alicloud_gpdb_instance" "default" {
  db_instance_category  = "HighAvailability"
  db_instance_class     = "gpdb.group.segsdx1"
  db_instance_mode      = "StorageElastic"
  description           = var.name
  engine                = "gpdb"
  engine_version        = "6.0"
  zone_id               = data.alicloud_gpdb_zones.default.ids.0
  instance_network_type = "VPC"
  instance_spec         = "2C16G"
  payment_type          = "PayAsYouGo"
  seg_storage_type      = "cloud_essd"
  seg_node_num          = 4
  storage_size          = 50
  vpc_id                = data.alicloud_vpcs.default.ids.0
  vswitch_id            = data.alicloud_vswitches.default.ids.0
}

resource "alicloud_gpdb_account" "default" {
  account_name        = "tfte123456"
  db_instance_id      = alicloud_gpdb_instance.default.id
  account_password    = "Example1234"
  account_description = "tf_example"
  account_type        = "Normal"
}
`, name)
}

// Test Gpdb Database. <<< Resource test cases, automatically generated.

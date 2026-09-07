// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test PolarDb Extension. >>> Resource test cases, automatically generated.
// Case 已重置副本 11821
func TestAccAliCloudPolarDbExtension_basic11821(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_polar_db_extension.default"
	ra := resourceAttrInit(resourceId, AlicloudPolarDbExtensionMap11821)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &PolarDbServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribePolarDbExtension")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1, 999)
	name := fmt.Sprintf("tfacc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudPolarDbExtensionBasicDependence11821)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"extension_name": "postgres_fdw",
					"db_cluster_id":  "${alicloud_polardb_cluster.resource_DBCluster_test_01.id}",
					"db_name":        "${alicloud_polardb_database.resource_Database_test.db_name}",
					"account_name":   "${alicloud_polardb_account.resource_Account_test.account_name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"extension_name": "postgres_fdw",
						"db_cluster_id":  CHECKSET,
						"db_name":        CHECKSET,
						"account_name":   CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"installed_version": "1.1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
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

var AlicloudPolarDbExtensionMap11821 = map[string]string{}

func AlicloudPolarDbExtensionBasicDependence11821(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_polardb_node_classes" "default" {
  db_type    = "PostgreSQL"
  pay_type   = "PostPaid"
  category   = "Normal"
}

resource "alicloud_vpc" "default" {
  vpc_name   = "terraform-example"
  cidr_block = "172.16.0.0/16"
}

resource "alicloud_vswitch" "default" {
  vpc_id       = alicloud_vpc.default.id
  cidr_block   = "172.16.0.0/24"
  zone_id      = data.alicloud_polardb_node_classes.default.classes[0].zone_id
  vswitch_name = "terraform-example"
}

resource "alicloud_polardb_cluster" "resource_DBCluster_test_01" {
  default_time_zone = "SYSTEM"
  creation_category = "Normal"
  storage_type      = "PSL5"
  db_version        = "14"
  pay_type          = "PostPaid"
  db_node_class     = "polar.pg.x4.medium"
  db_type           = "PostgreSQL"
  creation_option   = "Normal"
  vswitch_id        = alicloud_vswitch.default.id
}

resource "alicloud_polardb_account" "resource_Account_test" {
  db_cluster_id    = alicloud_polardb_cluster.resource_DBCluster_test_01.id
  account_name     = "nzh"
  account_password = "Ali123456"
  account_type     = "Super"
}

resource "alicloud_polardb_database" "resource_Database_test" {
  character_set_name = "utf8"
  db_cluster_id      = alicloud_polardb_cluster.resource_DBCluster_test_01.id
  db_name            = "nzh"
  account_name       = alicloud_polardb_account.resource_Account_test.account_name
  db_description     = var.name
}


`, name)
}

// Test PolarDb Extension. <<< Resource test cases, automatically generated.

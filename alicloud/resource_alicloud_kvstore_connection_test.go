package alicloud

import (
	"fmt"
	"testing"

	r_kvstore "github.com/aliyun/alibaba-cloud-sdk-go/services/r-kvstore"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudKVStoreConnection_basic(t *testing.T) {
	var v r_kvstore.InstanceNetInfo
	resourceId := "alicloud_kvstore_connection.default"
	ra := resourceAttrInit(resourceId, RedisConnectionMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &R_kvstoreService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeKvstoreConnection")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccKvstoreConnection%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, KvstoreConnectionBasicdependence)
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
					"connection_string_prefix": "allocatetest",
					"instance_id":              "${alicloud_kvstore_instance.default.id}",
					"port":                     "6370",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_id":              CHECKSET,
						"port":                     "6370",
						"connection_string_prefix": "allocatetest",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"connection_string_prefix"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"connection_string_prefix": "allocatetestupdate",
					"port":                     "6371",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"connection_string_prefix": "allocatetestupdate",
						"port":                     "6371",
					}),
				),
			},
		},
	})
}

var RedisConnectionMap = map[string]string{}

func KvstoreConnectionBasicdependence(name string) string {
	return fmt.Sprintf(`
	data "alicloud_kvstore_zones" "default"{
		instance_charge_type = "PostPaid"
	}
	data "alicloud_vpcs" "default" {
	  name_regex = "default-NODELETING"
	}
	data "alicloud_vswitches" "default" {
	  zone_id = data.alicloud_kvstore_zones.default.zones[length(data.alicloud_kvstore_zones.default.ids) - 1].id
	  vpc_id = data.alicloud_vpcs.default.ids.0
	}
	data "alicloud_resource_manager_resource_groups" "default" {
	}
	resource "alicloud_kvstore_instance" "default" {
		db_instance_name = "%s"
  		vswitch_id = data.alicloud_vswitches.default.ids.0
		instance_type = "Redis"
		engine_version = "4.0"
		tags = {
			Created = "TF",
			For = "update test",
		}
		resource_group_id = data.alicloud_resource_manager_resource_groups.default.ids.1
		zone_id = data.alicloud_kvstore_zones.default.zones[length(data.alicloud_kvstore_zones.default.ids) - 1].id
		instance_class="redis.master.large.default"
	}
	`, name)
}

package rkvstore_test

import (
	"fmt"
	"testing"

	r_kvstore "github.com/aliyun/alibaba-cloud-sdk-go/services/r-kvstore"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/rkvstore"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAliCloudKvstoreConnection_basic(t *testing.T) {
	var v r_kvstore.InstanceNetInfo
	resourceId := "alicloud_kvstore_connection.default"
	ra := acctest.ResourceAttrInit(resourceId, RedisConnectionMap)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return rkvstore.NewR_kvstoreService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeKvstoreConnection")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccKvstoreConnection%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, KvstoreConnectionBasicdependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"connection_string_prefix": "allocatetest",
					"instance_id":              "${alicloud_kvstore_instance.default.id}",
					"port":                     "6370",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_id":              acctest.CHECKSET,
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
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"connection_string_prefix": "allocatetestupdate",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"port": "6371",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"port": "6371",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"connection_string_prefix": "allocatetestup",
					"port":                     "6372",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"connection_string_prefix": "allocatetestup",
						"port":                     "6372",
					}),
				),
			},
		},
	})
}

var RedisConnectionMap = map[string]string{
	"connection_string": acctest.CHECKSET,
}

func KvstoreConnectionBasicdependence(name string) string {
	return fmt.Sprintf(`
	data "alicloud_kvstore_zones" "default"{
		instance_charge_type = "PostPaid"
	}
	data "alicloud_vpcs" "default" {
	  name_regex = "^default-NODELETING$"
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

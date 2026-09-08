package alicloud

import (
	"fmt"
	"os"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test Kvcachestore KvCacheStore. >>> Resource test cases.

// The HPN cluster ID in cn-shanghai-cloudspe, confirmed by the service team.
const testAccKvcachestoreCloudspeHpnZone = "b1"

// KVCacheStore is currently available only in whitelist regions; pin the test
// region to cn-shanghai-cloudspe unless ALICLOUD_REGION is set explicitly.
// NOTE: the pin must happen before testAccPreCheck, which defaults an empty
// ALICLOUD_REGION to cn-beijing. cn-shanghai-cloudspe is not in the provider's
// static region list, so region validation is skipped for these tests.
func testAccPreCheckKvcachestoreRegion(t *testing.T) {
	if v := os.Getenv("ALICLOUD_REGION"); v == "" {
		os.Setenv("ALICLOUD_REGION", "cn-shanghai-cloudspe")
		os.Setenv("ALICLOUD_SKIP_REGION_VALIDATION", "true")
	}
	testAccPreCheck(t)
}

var AliCloudKvcachestoreKvCacheStoreMap = map[string]string{
	"id":           CHECKSET,
	"create_time":  CHECKSET,
	"status":       CHECKSET,
	"payment_type": CHECKSET,
}

func TestAccAliCloudKvcachestoreKvCacheStore_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_kvcachestore_kv_cache_store.default"
	ra := resourceAttrInit(resourceId, AliCloudKvcachestoreKvCacheStoreMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &KvcachestoreServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeKvcachestoreKvCacheStore")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacckvcs%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudKvcachestoreKvCacheStoreBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckKvcachestoreRegion(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"name":              name,
					"description":       name,
					"zone_id":           "cn-shanghai-cloudspe-b",
					"hpn_zone":          testAccKvcachestoreCloudspeHpnZone,
					"capacity":          "307200",
					"payment_type":      "POSTPAY",
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.groups.0.id}",
					"tag": []map[string]interface{}{
						{
							"tag_key":   "env",
							"tag_value": "acceptance",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":              name,
						"description":       name,
						"zone_id":           "cn-shanghai-cloudspe-b",
						"hpn_zone":          testAccKvcachestoreCloudspeHpnZone,
						"capacity":          "307200",
						"payment_type":      "POSTPAY",
						"resource_group_id": CHECKSET,
						"tag.#":             "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"name":        name + "-update",
					"description": name + "-update",
					"capacity":    "614400",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":        name + "-update",
						"description": name + "-update",
						"capacity":    "614400",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tag": []map[string]interface{}{
						{
							"tag_key":   "env_updated",
							"tag_value": "acceptance-update",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tag.#": "1",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"resource_group_id", "tag"},
			},
		},
	})
}

func AliCloudKvcachestoreKvCacheStoreBasicDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_resource_manager_resource_groups" "default" {
  status = "OK"
}
`, name)
}

// Test Kvcachestore KvCacheStore. <<< Resource test cases.

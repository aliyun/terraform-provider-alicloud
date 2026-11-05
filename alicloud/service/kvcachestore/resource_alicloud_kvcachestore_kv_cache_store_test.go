package kvcachestore_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/kvcachestore"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
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
	acctest.PreCheck(t)
}

var AliCloudKvcachestoreKvCacheStoreMap = map[string]string{
	"id":           acctest.CHECKSET,
	"create_time":  acctest.CHECKSET,
	"status":       acctest.CHECKSET,
	"payment_type": acctest.CHECKSET,
}

func TestAccAliCloudKvcachestoreKvCacheStore_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_kvcachestore_kv_cache_store.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudKvcachestoreKvCacheStoreMap)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return kvcachestore.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeKvCacheStore")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacckvcs%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudKvcachestoreKvCacheStoreBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckKvcachestoreRegion(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"name":         name,
					"description":  name,
					"zone_id":      "cn-shanghai-cloudspe-b",
					"hpn_zone":     testAccKvcachestoreCloudspeHpnZone,
					"capacity":     "307200",
					"payment_type": "POSTPAY",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":         name,
						"description":  name,
						"zone_id":      "cn-shanghai-cloudspe-b",
						"hpn_zone":     testAccKvcachestoreCloudspeHpnZone,
						"capacity":     "307200",
						"payment_type": "POSTPAY",
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
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func AliCloudKvcachestoreKvCacheStoreBasicDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}
`, name)
}

// Test Kvcachestore KvCacheStore. <<< Resource test cases.

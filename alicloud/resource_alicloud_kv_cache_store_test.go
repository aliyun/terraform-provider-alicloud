package alicloud

import (
	"fmt"
	"os"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// KVCacheStore requires a minimum capacity of 307200 GiB (300 TiB) and expands
// in 300 TiB steps — an enterprise-grade resource that needs special account
// quota. The acceptance test is gated behind ALICLOUD_KVCACHESTORE_CAPACITY so
// CI skips it by default; set the env to a valid capacity (e.g. 307200) when
// the account has quota to provision a KVCacheStore instance.
func testAccPreCheckKvCacheStore(t *testing.T) {
	testAccPreCheck(t)
	if v := os.Getenv("ALICLOUD_KVCACHESTORE_CAPACITY"); v == "" {
		t.Skip("Skipping KVCacheStore acceptance test: minimum 300 TiB (307200 GiB) capacity requires special account quota. Set ALICLOUD_KVCACHESTORE_CAPACITY env to enable.")
	}
}

func TestAccAliCloudKvCacheStore_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_kv_cache_store.default"
	ra := resourceAttrInit(resourceId, AlicloudKvCacheStoreMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &KvcachestoreServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeKvcachestoreKVCacheStore")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc-kvcs-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudKvCacheStoreBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckKvCacheStore(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"capacity":          "307200",
					"zone_id":           "cn-beijing-i",
					"hpn_zone":          "default",
					"name":              name,
					"description":       "terraform-test-description",
					"payment_type":      "POSTPAY",
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"tags": map[string]string{
						"Created": "Terraform",
						"Env":     "test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"capacity":          "307200",
						"zone_id":           "cn-beijing-i",
						"hpn_zone":          "default",
						"name":              name,
						"description":       "terraform-test-description",
						"payment_type":      "POSTPAY",
						"resource_group_id": CHECKSET,
						"tags.Created":      "Terraform",
						"tags.Env":          "test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"capacity":          "614400",
					"name":              name + "-update",
					"description":       "terraform-test-description-updated",
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.1}",
					"tags": map[string]string{
						"Created": "Terraform",
						"Env":     "prod",
						"Team":    "iac",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"capacity":          "614400",
						"name":              name + "-update",
						"description":       "terraform-test-description-updated",
						"resource_group_id": CHECKSET,
						"tags.Created":      "Terraform",
						"tags.Env":          "prod",
						"tags.Team":         "iac",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"tags"},
			},
		},
	})
}

var AlicloudKvCacheStoreMap0 = map[string]string{
	"kvcs_id": CHECKSET,
}

func AlicloudKvCacheStoreBasicDependence0(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

data "alicloud_resource_manager_resource_groups" "default" {
}
`, name)
}

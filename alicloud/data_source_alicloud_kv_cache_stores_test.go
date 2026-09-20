package alicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudKvCacheStoresDataSource_basic(t *testing.T) {
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc-kvcs-ds-%d", rand)
	zoneId := "cn-beijing-i"
	hpnZone := "default"

	testAccConfig := fmt.Sprintf(`
variable "name" {
  default = "%s"
}

resource "alicloud_kv_cache_store" "default" {
  capacity     = 307200
  zone_id      = "%s"
  hpn_zone     = "%s"
  name         = var.name
  description  = "terraform-test-ds-description"
  payment_type = "POSTPAY"
}

data "alicloud_kv_cache_stores" "default" {
  ids = ["${alicloud_kv_cache_store.default.id}"]
}
`, name, zoneId, hpnZone)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckKvCacheStore(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.alicloud_kv_cache_stores.default", "stores.#"),
					resource.TestCheckResourceAttrSet("data.alicloud_kv_cache_stores.default", "stores.0.kvcs_id"),
				),
			},
		},
	})
}

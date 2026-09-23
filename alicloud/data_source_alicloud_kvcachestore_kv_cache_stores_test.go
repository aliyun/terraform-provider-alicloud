package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAliCloudKvcachestoreKvCacheStoresDataSource(t *testing.T) {
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacckvcs%d", rand)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudKvcachestoreKvCacheStoresDataSourceConfig(name, map[string]string{
			"ids": `[alicloud_kvcachestore_kv_cache_store.default.id]`,
		}),
		fakeConfig: testAccCheckAliCloudKvcachestoreKvCacheStoresDataSourceConfig(name, map[string]string{
			"ids": `["kvcs-fake00000"]`,
		}),
	}

	kvcsIdsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudKvcachestoreKvCacheStoresDataSourceConfig(name, map[string]string{
			"kvcs_ids": `"${alicloud_kvcachestore_kv_cache_store.default.id}"`,
		}),
		fakeConfig: testAccCheckAliCloudKvcachestoreKvCacheStoresDataSourceConfig(name, map[string]string{
			"kvcs_ids": `"${alicloud_kvcachestore_kv_cache_store.default.id}_fake"`,
		}),
	}

	nameConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudKvcachestoreKvCacheStoresDataSourceConfig(name, map[string]string{
			"name": `"${alicloud_kvcachestore_kv_cache_store.default.name}"`,
		}),
		fakeConfig: testAccCheckAliCloudKvcachestoreKvCacheStoresDataSourceConfig(name, map[string]string{
			"name": `"${alicloud_kvcachestore_kv_cache_store.default.name}_fake"`,
		}),
	}

	statusConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudKvcachestoreKvCacheStoresDataSourceConfig(name, map[string]string{
			"ids":    `[alicloud_kvcachestore_kv_cache_store.default.id]`,
			"status": `"${alicloud_kvcachestore_kv_cache_store.default.status}"`,
		}),
		fakeConfig: testAccCheckAliCloudKvcachestoreKvCacheStoresDataSourceConfig(name, map[string]string{
			"ids":    `[alicloud_kvcachestore_kv_cache_store.default.id]`,
			"status": `"Deleting"`,
		}),
	}

	zoneIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudKvcachestoreKvCacheStoresDataSourceConfig(name, map[string]string{
			"ids":     `[alicloud_kvcachestore_kv_cache_store.default.id]`,
			"zone_id": `"${alicloud_kvcachestore_kv_cache_store.default.zone_id}"`,
		}),
		fakeConfig: testAccCheckAliCloudKvcachestoreKvCacheStoresDataSourceConfig(name, map[string]string{
			"ids":     `[alicloud_kvcachestore_kv_cache_store.default.id]`,
			"zone_id": `"cn-shanghai-cloudspe-z"`,
		}),
	}

	detailsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudKvcachestoreKvCacheStoresDataSourceConfig(name, map[string]string{
			"ids":            `[alicloud_kvcachestore_kv_cache_store.default.id]`,
			"enable_details": `true`,
		}),
		fakeConfig: testAccCheckAliCloudKvcachestoreKvCacheStoresDataSourceConfig(name, map[string]string{
			"ids":            `["kvcs-fake00000"]`,
			"enable_details": `true`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudKvcachestoreKvCacheStoresDataSourceConfig(name, map[string]string{
			"ids":            `[alicloud_kvcachestore_kv_cache_store.default.id]`,
			"name":           `"${alicloud_kvcachestore_kv_cache_store.default.name}"`,
			"zone_id":        `"${alicloud_kvcachestore_kv_cache_store.default.zone_id}"`,
			"enable_details": `true`,
		}),
		fakeConfig: testAccCheckAliCloudKvcachestoreKvCacheStoresDataSourceConfig(name, map[string]string{
			"ids":            `["kvcs-fake00000"]`,
			"name":           `"${alicloud_kvcachestore_kv_cache_store.default.name}_fake"`,
			"zone_id":        `"cn-shanghai-cloudspe-z"`,
			"enable_details": `true`,
		}),
	}

	var existAliCloudKvcachestoreKvCacheStoresDataSourceMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                 "1",
			"stores.#":              "1",
			"stores.0.id":           CHECKSET,
			"stores.0.kvcs_id":      CHECKSET,
			"stores.0.name":         CHECKSET,
			"stores.0.zone_id":      CHECKSET,
			"stores.0.hpn_zone":     CHECKSET,
			"stores.0.capacity":     CHECKSET,
			"stores.0.payment_type": CHECKSET,
			"stores.0.region_id":    CHECKSET,
			"stores.0.description":  CHECKSET,
			"stores.0.create_time":  CHECKSET,
		}
	}

	var fakeAliCloudKvcachestoreKvCacheStoresDataSourceMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":    "0",
			"stores.#": "0",
		}
	}

	aliCloudKvcachestoreKvCacheStoresCheckInfo := dataSourceAttr{
		resourceId:   "data.alicloud_kvcachestore_kv_cache_stores.default",
		existMapFunc: existAliCloudKvcachestoreKvCacheStoresDataSourceMapFunc,
		fakeMapFunc:  fakeAliCloudKvcachestoreKvCacheStoresDataSourceMapFunc,
	}

	preCheck := func() {
		testAccPreCheckKvcachestoreRegion(t)
	}
	aliCloudKvcachestoreKvCacheStoresCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, kvcsIdsConf, nameConf, statusConf, zoneIdConf, detailsConf, allConf)
}

func testAccCheckAliCloudKvcachestoreKvCacheStoresDataSourceConfig(name string, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

resource "alicloud_kvcachestore_kv_cache_store" "default" {
  name         = var.name
  description  = var.name
  zone_id      = "cn-shanghai-cloudspe-b"
  hpn_zone     = "%s"
  capacity     = 307200
  payment_type = "POSTPAY"
}

data "alicloud_kvcachestore_kv_cache_stores" "default" {
  %s
}
`, name, testAccKvcachestoreCloudspeHpnZone, strings.Join(pairs, "\n  "))
}

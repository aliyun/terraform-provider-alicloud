package kvcachestore_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAliCloudKvcachestoreKvCacheStoresDataSource(t *testing.T) {
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacckvcs%d", rand)

	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAliCloudKvcachestoreKvCacheStoresDataSourceConfig(name, map[string]string{
			"ids": `[alicloud_kvcachestore_kv_cache_store.default.id]`,
		}),
		FakeConfig: testAccCheckAliCloudKvcachestoreKvCacheStoresDataSourceConfig(name, map[string]string{
			"ids": `["kvcs-fake00000"]`,
		}),
	}

	kvcsIdsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAliCloudKvcachestoreKvCacheStoresDataSourceConfig(name, map[string]string{
			"kvcs_ids": `"${alicloud_kvcachestore_kv_cache_store.default.id}"`,
		}),
		FakeConfig: testAccCheckAliCloudKvcachestoreKvCacheStoresDataSourceConfig(name, map[string]string{
			"kvcs_ids": `"${alicloud_kvcachestore_kv_cache_store.default.id}_fake"`,
		}),
	}

	nameConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAliCloudKvcachestoreKvCacheStoresDataSourceConfig(name, map[string]string{
			"name": `"${alicloud_kvcachestore_kv_cache_store.default.name}"`,
		}),
		FakeConfig: testAccCheckAliCloudKvcachestoreKvCacheStoresDataSourceConfig(name, map[string]string{
			"name": `"${alicloud_kvcachestore_kv_cache_store.default.name}_fake"`,
		}),
	}

	statusConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAliCloudKvcachestoreKvCacheStoresDataSourceConfig(name, map[string]string{
			"ids":    `[alicloud_kvcachestore_kv_cache_store.default.id]`,
			"status": `"${alicloud_kvcachestore_kv_cache_store.default.status}"`,
		}),
		FakeConfig: testAccCheckAliCloudKvcachestoreKvCacheStoresDataSourceConfig(name, map[string]string{
			"ids":    `[alicloud_kvcachestore_kv_cache_store.default.id]`,
			"status": `"Deleting"`,
		}),
	}

	zoneIdConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAliCloudKvcachestoreKvCacheStoresDataSourceConfig(name, map[string]string{
			"ids":     `[alicloud_kvcachestore_kv_cache_store.default.id]`,
			"zone_id": `"${alicloud_kvcachestore_kv_cache_store.default.zone_id}"`,
		}),
		FakeConfig: testAccCheckAliCloudKvcachestoreKvCacheStoresDataSourceConfig(name, map[string]string{
			"ids":     `[alicloud_kvcachestore_kv_cache_store.default.id]`,
			"zone_id": `"cn-shanghai-cloudspe-z"`,
		}),
	}

	detailsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAliCloudKvcachestoreKvCacheStoresDataSourceConfig(name, map[string]string{
			"ids":            `[alicloud_kvcachestore_kv_cache_store.default.id]`,
			"enable_details": `true`,
		}),
		FakeConfig: testAccCheckAliCloudKvcachestoreKvCacheStoresDataSourceConfig(name, map[string]string{
			"ids":            `["kvcs-fake00000"]`,
			"enable_details": `true`,
		}),
	}

	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAliCloudKvcachestoreKvCacheStoresDataSourceConfig(name, map[string]string{
			"ids":            `[alicloud_kvcachestore_kv_cache_store.default.id]`,
			"name":           `"${alicloud_kvcachestore_kv_cache_store.default.name}"`,
			"zone_id":        `"${alicloud_kvcachestore_kv_cache_store.default.zone_id}"`,
			"enable_details": `true`,
		}),
		FakeConfig: testAccCheckAliCloudKvcachestoreKvCacheStoresDataSourceConfig(name, map[string]string{
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
			"stores.0.id":           acctest.CHECKSET,
			"stores.0.kvcs_id":      acctest.CHECKSET,
			"stores.0.name":         acctest.CHECKSET,
			"stores.0.zone_id":      acctest.CHECKSET,
			"stores.0.hpn_zone":     acctest.CHECKSET,
			"stores.0.capacity":     acctest.CHECKSET,
			"stores.0.payment_type": acctest.CHECKSET,
			"stores.0.region_id":    acctest.CHECKSET,
			"stores.0.description":  acctest.CHECKSET,
			"stores.0.create_time":  acctest.CHECKSET,
		}
	}

	var fakeAliCloudKvcachestoreKvCacheStoresDataSourceMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":    "0",
			"stores.#": "0",
		}
	}

	aliCloudKvcachestoreKvCacheStoresCheckInfo := acctest.DataSourceAttr{
		ResourceId:   "data.alicloud_kvcachestore_kv_cache_stores.default",
		ExistMapFunc: existAliCloudKvcachestoreKvCacheStoresDataSourceMapFunc,
		FakeMapFunc:  fakeAliCloudKvcachestoreKvCacheStoresDataSourceMapFunc,
	}

	preCheck := func() {
		testAccPreCheckKvcachestoreRegion(t)
	}
	aliCloudKvcachestoreKvCacheStoresCheckInfo.DataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, kvcsIdsConf, nameConf, statusConf, zoneIdConf, detailsConf, allConf)
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

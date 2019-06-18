package alicloud

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/acctest"
	"strings"
	"testing"
)

func TestAccAlicloudKVStoreInstanceClasses(t *testing.T) {
	rand := acctest.RandInt()
	EngineVersionConf := dataSourceTestAccConfig{
		existConfig: generateConfig(map[string]string{
			"engine":         `"Redis"`,
			"engine_version": `"5.0"`,
		}),
		fakeConfig: generateConfig(map[string]string{
			"engine":         `"Redis"`,
			"engine_version": `"4.9"`,
		}),
	}

	ChargeTypeConfPostpaid := dataSourceTestAccConfig{
		existConfig: generateConfig(map[string]string{
			"instance_charge_type": `"PostPaid"`,
		}),
	}
	ArchitectureStandard := dataSourceTestAccConfig{
		existConfig: generateConfig(map[string]string{
			"architecture": `"standard"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: generateConfig(map[string]string{
			"instance_charge_type": `"PostPaid"`,
			"engine":               `"Redis"`,
			"engine_version":       `"5.0"`,
			"architecture":         `"standard"`,
			"performance_type":     `"standard_performance_type"`,
			"storage_type":         `"inmemory"`,
			"node_type":            `"double"`,
			"package_type":         `"standard"`,
		}),
		fakeConfig: generateConfig(map[string]string{
			"instance_charge_type": `"PostPaid"`,
			"engine":               `"Fake"`,
			"engine_version":       `"5.6"`,
			"architecture":         `"standard"`,
			"performance_type":     `"standard_performance_type"`,
			"storage_type":         `"inmemory"`,
			"node_type":            `"double"`,
			"package_type":         `"standard"`,
		}),
	}

	var existKVStoreInstanceMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"instance_classes.#": CHECKSET,
		}
	}

	var fakeKVStoreInstanceMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"instance_classes.#": "0",
		}
	}

	var KVStoreInstanceCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_kvstore_instance_classes.resources",
		existMapFunc: existKVStoreInstanceMapFunc,
		fakeMapFunc:  fakeKVStoreInstanceMapFunc,
	}

	KVStoreInstanceCheckInfo.dataSourceTestCheck(t, rand, EngineVersionConf, ChargeTypeConfPostpaid,
		ArchitectureStandard, allConf)
}

func generateConfig(attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
data "alicloud_zones" "resources" {
	available_resource_creation= "KVStore"
}
data "alicloud_kvstore_instance_classes" "resources" {
	"zone_id" = "${data.alicloud_zones.resources.zones.0.id}"
	%s
}
`, strings.Join(pairs, "\n  "))
	return config
}

package alicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
)

func TestAccAlicloudKVStoreInstanceClasses(t *testing.T) {
	rand := acctest.RandInt()
	resourceId := "data.alicloud_kvstore_instance_classes.default"

	testAccConfig := dataSourceTestAccConfigFunc(resourceId, "KVStore", kvstoreConfigHeader)

	EngineVersionConfRedis := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"zone_id":        "${data.alicloud_zones.resources.zones.0.id}",
			"engine":         "Redis",
			"engine_version": "5.0",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"zone_id":        "${data.alicloud_zones.resources.zones.0.id}",
			"engine":         "Redis",
			"engine_version": "4.9",
		}),
	}

	// At present, there are some limitation for sorted
	//prePaidSortedByConfRedis := dataSourceTestAccConfig{
	//	existConfig: testAccConfig(map[string]interface{}{
	//		"zone_id":              "${data.alicloud_zones.resources.zones.0.id}",
	//		"engine":               "Redis",
	//		"engine_version":       "5.0",
	//		"instance_charge_type": "PrePaid",
	//		"sorted_by":            "Price",
	//	}),
	//	existChangMap: map[string]string{
	//		"classes.#":                CHECKSET,
	//		"classes.0.instance_class": CHECKSET,
	//		"classes.0.price":          CHECKSET,
	//	},
	//}
	//
	//postPaidSortedByConfRedis := dataSourceTestAccConfig{
	//	existConfig: testAccConfig(map[string]interface{}{
	//		"zone_id":              "${data.alicloud_zones.resources.zones.0.id}",
	//		"engine":               "Redis",
	//		"engine_version":       "5.0",
	//		"instance_charge_type": "PostPaid",
	//		"sorted_by":            "Price",
	//	}),
	//	existChangMap: map[string]string{
	//		"classes.#":                CHECKSET,
	//		"classes.0.instance_class": CHECKSET,
	//		"classes.0.price":          CHECKSET,
	//	},
	//}
	EngineVersionConfMemcache := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"zone_id": "${data.alicloud_zones.resources.zones.0.id}",
			"engine":  "Memcache",
		}),
	}

	ChargeTypeConfPostpaid := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"zone_id":              "${data.alicloud_zones.resources.zones.0.id}",
			"instance_charge_type": "PostPaid",
		}),
	}
	seriesTypeEnhancedPerformanceType := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"zone_id":     "${data.alicloud_zones.resources.zones.0.id}",
			"series_type": "enhanced_performance_type",
		}),
	}
	editionTypeCommunity := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"zone_id":      "${data.alicloud_zones.resources.zones.0.id}",
			"edition_type": "Community",
		}),
	}
	shardNumber8 := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"zone_id":      "${data.alicloud_zones.resources.zones.0.id}",
			"shard_number": "8",
		}),
	}
	ArchitectureStandard := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"zone_id":      "${data.alicloud_zones.resources.zones.0.id}",
			"architecture": "standard",
		}),
	}
	ArchitectureCluster := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"zone_id":      "${data.alicloud_zones.resources.zones.0.id}",
			"architecture": "cluster",
		}),
	}
	// Not all of zone support rwsplit
	//ArchitectureRwsplit := dataSourceTestAccConfig{
	//	existConfig: testAccConfig(map[string]interface{}{
	//		"zone_id":      "${data.alicloud_zones.resources.zones.0.id}",
	//		"architecture": "rwsplit",
	//	}),
	//}
	NodeType := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"zone_id":   "${data.alicloud_zones.resources.zones.0.id}",
			"node_type": "double",
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"zone_id":              "${data.alicloud_zones.resources.zones.0.id}",
			"instance_charge_type": "PostPaid",
			"engine":               "Redis",
			"engine_version":       "5.0",
			"architecture":         "standard",
			"series_type":          "enhanced_performance_type",
			"edition_type":         "Community",
			"node_type":            "double",
			"shard_number":         "1",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"zone_id":              "${data.alicloud_zones.resources.zones.0.id}",
			"instance_charge_type": "PostPaid",
			"engine":               "Redis",
			"engine_version":       "5.6",
			"architecture":         "standard",
			"series_type":          "enhanced_performance_type",
			"edition_type":         "Community",
			"node_type":            "double",
			"shard_number":         "1",
		}),
	}

	var existKVStoreInstanceMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"instance_classes.#": CHECKSET,
			"instance_classes.0": CHECKSET,
		}
	}

	var fakeKVStoreInstanceMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"instance_classes.#": "0",
			"classes.#":          "0",
		}
	}

	var KVStoreInstanceCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existKVStoreInstanceMapFunc,
		fakeMapFunc:  fakeKVStoreInstanceMapFunc,
	}

	// At present, the datasource does not support memcache
	//KVStoreInstanceCheckInfo.dataSourceTestCheck(t, rand, EngineVersionConfRedis, EngineVersionConfMemcache,
	//	ChargeTypeConfPostpaid, PerformanceTypeStandardPerformanceType, PerformanceTypeEnhancePerformanceType,
	//	StorageTypeInmemory, PackageTypeStandard, PackageTypeCustomized, ArchitectureStandard, ArchitectureCluster,
	//	ArchitectureRwsplit, NodeTypeDouble, NodeTypeSingle, NodeTypeReadone, NodeTypeReadthree, NodeTypeReadfive,
	//	ArchitectureStandard, allConf)
	KVStoreInstanceCheckInfo.dataSourceTestCheck(t, rand, EngineVersionConfRedis,
		//prePaidSortedByConfRedis, postPaidSortedByConfRedis
		ChargeTypeConfPostpaid, seriesTypeEnhancedPerformanceType, editionTypeCommunity,
		shardNumber8, ArchitectureStandard, ArchitectureCluster,
		NodeType, allConf, EngineVersionConfMemcache)
}

func kvstoreConfigHeader(name string) string {
	return fmt.Sprintf(`
data "alicloud_zones" "resources" {
	available_resource_creation= "%s"
}
`, name)
}

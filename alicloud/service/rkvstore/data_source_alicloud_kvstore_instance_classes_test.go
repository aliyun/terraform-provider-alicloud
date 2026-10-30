package rkvstore_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAliCloudKvstoreInstanceClasses(t *testing.T) {
	rand := tfacctest.RandInt()
	resourceId := "data.alicloud_kvstore_instance_classes.default"

	testAccConfig := acctest.DataSourceTestAccConfigFunc(resourceId, "KVStore", kvstoreConfigHeader)

	EngineVersionConfRedis := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"zone_id":        "${data.alicloud_zones.resources.zones.0.id}",
			"engine":         "Redis",
			"engine_version": "5.0",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"zone_id":        "${data.alicloud_zones.resources.zones.0.id}",
			"engine":         "Redis",
			"engine_version": "4.9",
		}),
	}

	// At present, there are some limitation for sorted
	//prePaidSortedByConfRedis := acctest.DataSourceTestAccConfig{
	//	ExistConfig: testAccConfig(map[string]interface{}{
	//		"zone_id":              "${data.alicloud_zones.resources.zones.0.id}",
	//		"engine":               "Redis",
	//		"engine_version":       "5.0",
	//		"instance_charge_type": "PrePaid",
	//		"sorted_by":            "Price",
	//	}),
	//	ExistChangMap: map[string]string{
	//		"classes.#":                acctest.CHECKSET,
	//		"classes.0.instance_class": acctest.CHECKSET,
	//		"classes.0.price":          acctest.CHECKSET,
	//	},
	//}
	//
	//postPaidSortedByConfRedis := acctest.DataSourceTestAccConfig{
	//	ExistConfig: testAccConfig(map[string]interface{}{
	//		"zone_id":              "${data.alicloud_zones.resources.zones.0.id}",
	//		"engine":               "Redis",
	//		"engine_version":       "5.0",
	//		"instance_charge_type": "PostPaid",
	//		"sorted_by":            "Price",
	//	}),
	//	ExistChangMap: map[string]string{
	//		"classes.#":                acctest.CHECKSET,
	//		"classes.0.instance_class": acctest.CHECKSET,
	//		"classes.0.price":          acctest.CHECKSET,
	//	},
	//}
	EngineVersionConfMemcache := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"zone_id": "${data.alicloud_zones.resources.zones.0.id}",
			"engine":  "Memcache",
		}),
	}

	ChargeTypeConfPostpaid := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"zone_id":              "${data.alicloud_zones.resources.zones.0.id}",
			"instance_charge_type": "PostPaid",
		}),
	}
	seriesTypeEnhancedPerformanceType := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"zone_id":     "${data.alicloud_zones.resources.zones.0.id}",
			"series_type": "enhanced_performance_type",
		}),
	}
	editionTypeCommunity := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"zone_id":      "${data.alicloud_zones.resources.zones.0.id}",
			"edition_type": "Community",
		}),
	}
	shardNumber8 := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"zone_id":      "${data.alicloud_zones.resources.zones.0.id}",
			"shard_number": "8",
		}),
	}
	ArchitectureStandard := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"zone_id":      "${data.alicloud_zones.resources.zones.0.id}",
			"architecture": "standard",
		}),
	}
	ArchitectureCluster := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"zone_id":      "${data.alicloud_zones.resources.zones.0.id}",
			"architecture": "cluster",
		}),
	}
	// Not all of zone support rwsplit
	//ArchitectureRwsplit := acctest.DataSourceTestAccConfig{
	//	ExistConfig: testAccConfig(map[string]interface{}{
	//		"zone_id":      "${data.alicloud_zones.resources.zones.0.id}",
	//		"architecture": "rwsplit",
	//	}),
	//}
	NodeType := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"zone_id":   "${data.alicloud_zones.resources.zones.0.id}",
			"node_type": "double",
		}),
	}

	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
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
		FakeConfig: testAccConfig(map[string]interface{}{
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
			"instance_classes.#": acctest.CHECKSET,
			"instance_classes.0": acctest.CHECKSET,
		}
	}

	var fakeKVStoreInstanceMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"instance_classes.#": "0",
			"classes.#":          "0",
		}
	}

	var KVStoreInstanceCheckInfo = acctest.DataSourceAttr{
		ResourceId:   resourceId,
		ExistMapFunc: existKVStoreInstanceMapFunc,
		FakeMapFunc:  fakeKVStoreInstanceMapFunc,
	}

	// At present, the datasource does not support memcache
	//KVStoreInstanceCheckInfo.DataSourceTestCheck(t, rand, EngineVersionConfRedis, EngineVersionConfMemcache,
	//	ChargeTypeConfPostpaid, PerformanceTypeStandardPerformanceType, PerformanceTypeEnhancePerformanceType,
	//	StorageTypeInmemory, PackageTypeStandard, PackageTypeCustomized, ArchitectureStandard, ArchitectureCluster,
	//	ArchitectureRwsplit, NodeTypeDouble, NodeTypeSingle, NodeTypeReadone, NodeTypeReadthree, NodeTypeReadfive,
	//	ArchitectureStandard, allConf)
	KVStoreInstanceCheckInfo.DataSourceTestCheck(t, rand, EngineVersionConfRedis,
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

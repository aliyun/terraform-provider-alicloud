package rkvstore_test

import (
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
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
	editionTypeEnterprise := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"zone_id":      "${data.alicloud_zones.resources.zones.0.id}",
			"product_type": "Local",
			"edition_type": "Enterprise",
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
	architectureNonCluster := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"zone_id":      "${data.alicloud_zones.resources.zones.0.id}",
			"product_type": "Tair_rdb",
			"architecture": "non_cluster",
		}),
	}
	editionTypeEnhancedPerformance := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"zone_id":      "${data.alicloud_zones.resources.zones.0.id}",
			"product_type": "Tair_rdb",
			"edition_type": "EnhancedPerformance",
		}),
	}
	productTypeOnECS := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"zone_id":      "${data.alicloud_zones.resources.zones.0.id}",
			"product_type": "OnECS",
			"architecture": "non_cluster",
		}),
	}
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
	KVStoreInstanceCheckInfo.DataSourceTestCheckWithPreCheck(t, rand, func() {
		acctest.PreCheck(t)
		acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
	}, EngineVersionConfRedis,
		//prePaidSortedByConfRedis, postPaidSortedByConfRedis
		ChargeTypeConfPostpaid, seriesTypeEnhancedPerformanceType, editionTypeCommunity, editionTypeEnterprise,
		shardNumber8, ArchitectureStandard, ArchitectureCluster,
		architectureNonCluster, editionTypeEnhancedPerformance, productTypeOnECS,
		NodeType, allConf, EngineVersionConfMemcache)
}

// TestAccAliCloudKvstoreInstanceClassesFilterIdentity checks that the architecture and
// edition filters partition the unfiltered result instead of merely returning a
// non-empty list: the filtered counts add up to the unfiltered count and sibling values
// yield no classes at all. The equalities hold while Tair_rdb exposes exactly the
// non_cluster/cluster architectures and the Community/EnhancedPerformance editions,
// which is the case in the pinned zone.
func TestAccAliCloudKvstoreInstanceClassesFilterIdentity(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		ProviderFactories: acctest.SDKProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
data "alicloud_zones" "resources" {
  available_resource_creation = "KVStore"
}

data "alicloud_kvstore_instance_classes" "all" {
  zone_id      = data.alicloud_zones.resources.zones.0.id
  product_type = "Tair_rdb"
}

data "alicloud_kvstore_instance_classes" "non_cluster" {
  zone_id      = data.alicloud_zones.resources.zones.0.id
  product_type = "Tair_rdb"
  architecture = "non_cluster"
}

data "alicloud_kvstore_instance_classes" "cluster" {
  zone_id      = data.alicloud_zones.resources.zones.0.id
  product_type = "Tair_rdb"
  architecture = "cluster"
}

data "alicloud_kvstore_instance_classes" "community" {
  zone_id      = data.alicloud_zones.resources.zones.0.id
  product_type = "Tair_rdb"
  edition_type = "Community"
}

data "alicloud_kvstore_instance_classes" "enhanced" {
  zone_id      = data.alicloud_zones.resources.zones.0.id
  product_type = "Tair_rdb"
  edition_type = "EnhancedPerformance"
}

data "alicloud_kvstore_instance_classes" "rwsplit" {
  zone_id      = data.alicloud_zones.resources.zones.0.id
  product_type = "Tair_rdb"
  architecture = "rwsplit"
}

data "alicloud_kvstore_instance_classes" "enterprise" {
  zone_id      = data.alicloud_zones.resources.zones.0.id
  product_type = "Tair_rdb"
  edition_type = "Enterprise"
}

data "alicloud_kvstore_instance_classes" "standard" {
  zone_id      = data.alicloud_zones.resources.zones.0.id
  product_type = "Tair_rdb"
  architecture = "standard"
}

data "alicloud_kvstore_instance_classes" "on_ecs" {
  zone_id      = data.alicloud_zones.resources.zones.0.id
  product_type = "OnECS"
  architecture = "non_cluster"
}
`,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKvstoreInstanceClassesFilterIdentity(),
					testAccCheckKvstoreInstanceClassesContains("data.alicloud_kvstore_instance_classes.non_cluster", "tair.rdb."),
					testAccCheckKvstoreInstanceClassesContains("data.alicloud_kvstore_instance_classes.enhanced", "tair.rdb."),
					testAccCheckKvstoreInstanceClassesContains("data.alicloud_kvstore_instance_classes.on_ecs", "redis.shard."),
				),
			},
		},
	})
}

// TestAccAliCloudKvstoreInstanceClassesInternational runs the cloud-native query in an
// international region, where the original issue was reported. It changes the process
// wide region, so it must be selected by name rather than run as part of the package.
func TestAccAliCloudKvstoreInstanceClassesInternational(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{connectivity.APSouthEast1})
		},
		ProviderFactories: acctest.SDKProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
data "alicloud_zones" "resources" {
  available_resource_creation = "KVStore"
}

data "alicloud_kvstore_instance_classes" "default" {
  zone_id      = data.alicloud_zones.resources.zones.0.id
  product_type = "Tair_rdb"
  architecture = "non_cluster"
}
`,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKvstoreInstanceClassesContains("data.alicloud_kvstore_instance_classes.default", "tair.rdb."),
				),
			},
		},
	})
}

func testAccKvstoreInstanceClasses(id string, s *terraform.State) (int, []string, error) {
	rs, ok := s.RootModule().Resources[id]
	if !ok || rs.Primary == nil {
		return 0, nil, fmt.Errorf("data source %s not found in state", id)
	}
	count, err := strconv.Atoi(rs.Primary.Attributes["instance_classes.#"])
	if err != nil {
		return 0, nil, fmt.Errorf("read instance_classes.# of %s: %s", id, err)
	}
	classes := make([]string, 0, count)
	for i := 0; i < count; i++ {
		classes = append(classes, rs.Primary.Attributes[fmt.Sprintf("instance_classes.%d", i)])
	}
	return count, classes, nil
}

func testAccCheckKvstoreInstanceClassesContains(id, prefix string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		count, classes, err := testAccKvstoreInstanceClasses(id, s)
		if err != nil {
			return err
		}
		if count == 0 {
			return fmt.Errorf("%s returned no instance classes", id)
		}
		for _, class := range classes {
			if strings.HasPrefix(class, prefix) {
				return nil
			}
		}
		return fmt.Errorf("%s returned no class starting with %q, got %v", id, prefix, classes)
	}
}

func testAccCheckKvstoreInstanceClassesFilterIdentity() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		counts := map[string]int{}
		for _, name := range []string{"all", "non_cluster", "cluster", "community", "enhanced", "rwsplit", "enterprise", "standard"} {
			id := "data.alicloud_kvstore_instance_classes." + name
			count, _, err := testAccKvstoreInstanceClasses(id, s)
			if err != nil {
				return err
			}
			counts[name] = count
		}

		for _, name := range []string{"all", "non_cluster", "cluster", "community", "enhanced"} {
			if counts[name] == 0 {
				return fmt.Errorf("instance_classes of %s is empty", name)
			}
		}
		if counts["non_cluster"]+counts["cluster"] != counts["all"] {
			return fmt.Errorf("architecture filter did not partition the result: non_cluster=%d + cluster=%d != all=%d",
				counts["non_cluster"], counts["cluster"], counts["all"])
		}
		if counts["community"]+counts["enhanced"] != counts["all"] {
			return fmt.Errorf("edition filter did not partition the result: community=%d + enhanced=%d != all=%d",
				counts["community"], counts["enhanced"], counts["all"])
		}
		if counts["rwsplit"] != 0 {
			return fmt.Errorf("architecture=rwsplit returned %d classes for Tair_rdb, want 0", counts["rwsplit"])
		}
		if counts["enterprise"] != 0 {
			return fmt.Errorf("edition_type=Enterprise returned %d classes for Tair_rdb, want 0", counts["enterprise"])
		}
		if counts["standard"] != 0 {
			return fmt.Errorf("architecture=standard returned %d classes for Tair_rdb, want 0", counts["standard"])
		}
		return nil
	}
}

func kvstoreConfigHeader(name string) string {
	return fmt.Sprintf(`
data "alicloud_zones" "resources" {
	available_resource_creation= "%s"
}
`, name)
}

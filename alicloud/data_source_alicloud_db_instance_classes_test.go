package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
)

func TestAccAlicloudDBInstanceClasses_base(t *testing.T) {
	rand := acctest.RandInt()
	EngineVersionConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudDBInstanceClassesDataSourceConfig(map[string]string{
			"engine":         `"MySQL"`,
			"engine_version": `"5.6"`,
		}),
		fakeConfig: testAccCheckAlicloudDBInstanceClassesDataSourceConfig(map[string]string{
			"engine":         `"MySQL"`,
			"engine_version": `"3.0"`,
		}),
	}
	// At present, there are some limitation for sorted
	//prePaidSortedByConf := dataSourceTestAccConfig{
	//	existConfig: testAccCheckAlicloudDBInstanceClassesDataSourceConfig(map[string]string{
	//		"engine":               `"MySQL"`,
	//		"engine_version":       `"5.6"`,
	//		"instance_charge_type": `"PrePaid"`,
	//		"sorted_by":            `"Price"`,
	//	}),
	//	existChangMap: map[string]string{
	//		"instance_classes.0.price": CHECKSET,
	//	},
	//}
	//
	//postPaidSortedByConf := dataSourceTestAccConfig{
	//	existConfig: testAccCheckAlicloudDBInstanceClassesDataSourceConfig(map[string]string{
	//		"engine":               `"MySQL"`,
	//		"engine_version":       `"5.6"`,
	//		"instance_charge_type": `"PostPaid"`,
	//		"sorted_by":            `"Price"`,
	//	}),
	//	existChangMap: map[string]string{
	//		"instance_classes.0.price": CHECKSET,
	//	},
	//}

	ChargeTypeConf_Prepaid := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudDBInstanceClassesDataSourceConfig(map[string]string{
			"instance_charge_type": `"PrePaid"`,
		}),
	}
	ChargeTypeConf_Postpaid := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudDBInstanceClassesDataSourceConfig(map[string]string{
			"instance_charge_type": `"PostPaid"`,
		}),
	}
	StorageTypeConf_local_ssd := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudDBInstanceClassesDataSourceConfig(map[string]string{
			"storage_type": `"local_ssd"`,
		}),
	}

	StorageTypeConf_cloud_ssd := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudDBInstanceClassesDataSourceConfig(map[string]string{
			"storage_type": `"cloud_ssd"`,
		}),
	}
	multiZoneConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudDBInstanceClassesDataSourceConfig(map[string]string{
			"multi_zone": `"true"`,
		}),
	}
	falseMultiZoneConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudDBInstanceClassesDataSourceConfig(map[string]string{
			"multi_zone": `"false"`,
		}),
	}
	CategoryConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudDBInstanceClassesDataSourceConfig(map[string]string{
			"category": `"Basic"`,
		}),
		fakeConfig: testAccCheckAlicloudDBInstanceClassesDataSourceConfig(map[string]string{
			"category": `"fake"`,
		}),
	}
	DBInstanceClassConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudDBInstanceClassesDataSourceConfig(map[string]string{
			"db_instance_class": `"mysql.n2.large.1"`,
		}),
		fakeConfig: testAccCheckAlicloudDBInstanceClassesDataSourceConfig(map[string]string{
			"db_instance_class": `"fake"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudDBInstanceClassesDataSourceConfig(map[string]string{
			"instance_charge_type": `"PostPaid"`,
			"engine":               `"MySQL"`,
			"engine_version":       `"5.6"`,
			//"sorted_by":            `"Price"`,
		}),
		//existChangMap: map[string]string{
		//	"instance_classes.0.price": CHECKSET,
		//},
		fakeConfig: testAccCheckAlicloudDBInstanceClassesDataSourceConfig(map[string]string{
			"instance_charge_type": `"PostPaid"`,
			"engine":               `"Fake"`,
			"engine_version":       `"5.6"`,
			//"sorted_by":            `"Price"`,
		}),
	}

	var existDBInstanceMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"instance_classes.#":                CHECKSET,
			"instance_classes.0.instance_class": CHECKSET,
			//"instance_classes.0.price":                     "",
			"instance_classes.0.storage_range.min":         CHECKSET,
			"instance_classes.0.storage_range.max":         CHECKSET,
			"instance_classes.0.storage_range.step":        CHECKSET,
			"instance_classes.0.zone_ids.0.id":             CHECKSET,
			"instance_classes.0.zone_ids.0.sub_zone_ids.0": CHECKSET,
		}
	}

	var fakeDBInstanceMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"instance_classes.#": "0",
		}
	}

	var DBInstanceCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_db_instance_classes.default",
		existMapFunc: existDBInstanceMapFunc,
		fakeMapFunc:  fakeDBInstanceMapFunc,
	}

	//DBInstanceCheckInfo.dataSourceTestCheck(t, rand, EngineVersionConf, prePaidSortedByConf, postPaidSortedByConf,
	//	ChargeTypeConf_Prepaid, ChargeTypeConf_Postpaid, CategoryConf, DBInstanceClassConf, multiZoneConf, falseMultiZoneConf, StorageTypeConf_local_ssd, StorageTypeConf_cloud_ssd, allConf)
	DBInstanceCheckInfo.dataSourceTestCheck(t, rand, EngineVersionConf,
		ChargeTypeConf_Prepaid, ChargeTypeConf_Postpaid, CategoryConf, DBInstanceClassConf, multiZoneConf, falseMultiZoneConf, StorageTypeConf_local_ssd, StorageTypeConf_cloud_ssd, allConf)
}

func testAccCheckAlicloudDBInstanceClassesDataSourceConfig(attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
data "alicloud_zones" "default" {
  available_resource_creation= "Rds"
}
data "alicloud_db_instance_classes" "default" {
  %s
}
`, strings.Join(pairs, "\n  "))
	return config
}

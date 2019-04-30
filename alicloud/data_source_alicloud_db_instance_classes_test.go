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

	CategoryConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudDBInstanceClassesDataSourceConfig(map[string]string{
			"category": `"Basic"`,
		}),
		fakeConfig: testAccCheckAlicloudDBInstanceClassesDataSourceConfig(map[string]string{
			"category": `"fake"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudDBInstanceClassesDataSourceConfig(map[string]string{
			"instance_charge_type": `"PostPaid"`,
			"zone_id":              `"${data.alicloud_zones.resources.zones.0.id}"`,
			"engine":               `"MySQL"`,
			"engine_version":       `"5.6"`,
		}),
		fakeConfig: testAccCheckAlicloudDBInstanceClassesDataSourceConfig(map[string]string{
			"instance_charge_type": `"PostPaid"`,
			"zone_id":              `"${data.alicloud_zones.resources.zones.0.id}"`,
			"engine":               `"Fake"`,
			"engine_version":       `"5.6"`,
		}),
	}

	var existDBInstanceMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"instance_classes.#":                    CHECKSET,
			"instance_classes.0.zone_id":            CHECKSET,
			"instance_classes.0.instance_class":     CHECKSET,
			"instance_classes.0.storage_range.min":  CHECKSET,
			"instance_classes.0.storage_range.max":  CHECKSET,
			"instance_classes.0.storage_range.step": CHECKSET,
		}
	}

	var fakeDBInstanceMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"instance_classes.#": "0",
		}
	}

	var DBInstanceCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_db_instance_classes.resources",
		existMapFunc: existDBInstanceMapFunc,
		fakeMapFunc:  fakeDBInstanceMapFunc,
	}

	DBInstanceCheckInfo.dataSourceTestCheck(t, rand, EngineVersionConf, ChargeTypeConf_Prepaid, ChargeTypeConf_Postpaid, CategoryConf, StorageTypeConf_local_ssd, StorageTypeConf_cloud_ssd, allConf)
}

func testAccCheckAlicloudDBInstanceClassesDataSourceConfig(attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
data "alicloud_zones" "resources" {
  available_resource_creation= "Rds"
}
data "alicloud_db_instance_classes" "resources" {
  %s
}
`, strings.Join(pairs, "\n  "))
	return config
}

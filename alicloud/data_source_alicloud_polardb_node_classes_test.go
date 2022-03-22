package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudPolarDBNodeClasses(t *testing.T) {
	checkoutSupportedRegions(t, true, connectivity.PolarDBSupportRegions)
	rand := acctest.RandInt()
	resourceId := "data.alicloud_polardb_node_classes.default"

	testAccConfig := dataSourceTestAccConfigFunc(resourceId, "PolarDB", polardbConfigHeader)

	PayTypeConfPrepaid := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"pay_type": "PrePaid",
		}),
	}
	PayTypeConfPostpaid := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"pay_type": "PostPaid",
		}),
	}

	EngineVersionConfMysql := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"pay_type":   "PostPaid",
			"db_type":    "Mysql",
			"db_version": "5.6",
		}),
	}

	EngineVersionConfpgsql := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"pay_type":   "PostPaid",
			"db_type":    "PostgreSQL",
			"db_version": "11",
		}),
	}

	DBNodeClassConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"pay_type":      "PostPaid",
			"db_node_class": "polar.pg.x4.large",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"pay_type":      "PostPaid",
			"db_node_class": "fake",
		}),
	}

	RegionIdConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"pay_type":  "PostPaid",
			"region_id": fmt.Sprintf("%s", defaultRegionToTest),
		}),
	}

	// TODO: There is an api bug. It should reopen after the bug has been fixed.
	//ZoneIdConf := dataSourceTestAccConfig{
	//	existConfig: testAccConfig(map[string]interface{}{
	//		"pay_type": "PostPaid",
	//		"zone_id":  "${data.alicloud_polardb_zones.resources.zones.0.id}",
	//	}),
	//}

	var existPolardbInstanceMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"classes.#": CHECKSET,
			"classes.0.supported_engines.0.available_resources.0.db_node_class": CHECKSET,
			"classes.0.supported_engines.0.engine":                              CHECKSET,
			"classes.0.zone_id":                                                 CHECKSET,
		}
	}

	var fakePolardbInstanceMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"classes.#": "0",
		}
	}

	var PolardbInstanceCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existPolardbInstanceMapFunc,
		fakeMapFunc:  fakePolardbInstanceMapFunc,
	}

	PolardbInstanceCheckInfo.dataSourceTestCheck(t, rand, PayTypeConfPrepaid,
		PayTypeConfPostpaid, EngineVersionConfMysql, EngineVersionConfpgsql,
		DBNodeClassConf, RegionIdConf)
}

func polardbConfigHeader(name string) string {
	return fmt.Sprintf(`
data "alicloud_polardb_zones" "resources" {}
`)
}

package alicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

var (
	existHBaseInstanceTypesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                                 CHECKSET,
			"types.#":                               CHECKSET,
			"master_instance_types.#":               CHECKSET,
			"master_instance_types.0.instance_type": CHECKSET,
			"master_instance_types.0.cpu_size":      CHECKSET,
			"master_instance_types.0.mem_size":      CHECKSET,
			"core_instance_types.#":                 CHECKSET,
			"core_instance_types.0.instance_type":   CHECKSET,
			"core_instance_types.0.cpu_size":        CHECKSET,
			"core_instance_types.0.mem_size":        CHECKSET,
		}
	}

	fakeHBaseInstanceTypesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"master_instance_types.#": "0",
			"core_instance_types.#":   "0",
		}
	}

	hbaseInstanceTypesCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_hbase_instance_types.default",
		existMapFunc: existHBaseInstanceTypesMapFunc,
		fakeMapFunc:  fakeHBaseInstanceTypesMapFunc,
	}
)

func TestAccAlicloudHBaseInstanceTypesDataSource(t *testing.T) {
	rand := acctest.RandInt()
	resourceID := "data.alicloud_hbase_instance_types.default"

	testAccConfig := dataSourceTestAccConfigFunc(resourceID,
		fmt.Sprintf("tf-testacc%sbase-instance-instance_types%v.abc", defaultRegionToTest, rand),
		dataSourceHBaseInstanceTypesConfigDependence)

	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"charge_type": "Postpaid",
			"region_id":   "cn-shanghai",
			"zone_id":     "cn-shanghai-g",
			"engine":      "hbaseue",
			"version":     "2.0",
		}),
	}

	oneConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"charge_type":   "Postpaid",
			"region_id":     "cn-shanghai",
			"zone_id":       "cn-shanghai-g",
			"engine":        "hbaseue",
			"version":       "2.0",
			"instance_type": "hbase.sn1.8xlarge",
			"disk_type":     "cloud_ssd",
		}),
	}

	hbaseInstanceTypesCheckInfo.dataSourceTestCheck(t, rand, oneConf, allConf)
}

func dataSourceHBaseInstanceTypesConfigDependence(name string) string {
	return ``
}

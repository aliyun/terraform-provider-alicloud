package hbase_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

var (
	existHBaseInstanceTypesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                                 acctest.CHECKSET,
			"types.#":                               acctest.CHECKSET,
			"master_instance_types.#":               acctest.CHECKSET,
			"master_instance_types.0.instance_type": acctest.CHECKSET,
			"master_instance_types.0.cpu_size":      acctest.CHECKSET,
			"master_instance_types.0.mem_size":      acctest.CHECKSET,
			"core_instance_types.#":                 acctest.CHECKSET,
			"core_instance_types.0.instance_type":   acctest.CHECKSET,
			"core_instance_types.0.cpu_size":        acctest.CHECKSET,
			"core_instance_types.0.mem_size":        acctest.CHECKSET,
		}
	}

	fakeHBaseInstanceTypesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"master_instance_types.#": "0",
			"core_instance_types.#":   "0",
		}
	}

	hbaseInstanceTypesCheckInfo = acctest.DataSourceAttr{
		ResourceId:   "data.alicloud_hbase_instance_types.default",
		ExistMapFunc: existHBaseInstanceTypesMapFunc,
		FakeMapFunc:  fakeHBaseInstanceTypesMapFunc,
	}
)

func TestAccAlicloudHBaseInstanceTypesDataSource(t *testing.T) {
	rand := tfacctest.RandInt()
	resourceID := "data.alicloud_hbase_instance_types.default"

	testAccConfig := acctest.DataSourceTestAccConfigFunc(resourceID,
		fmt.Sprintf("tf-testacc%sbase-instance-instance_types%v.abc", acctest.DefaultRegionToTest, rand),
		dataSourceHBaseInstanceTypesConfigDependence)

	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"charge_type": "Postpaid",
			"region_id":   "cn-shanghai",
			"zone_id":     "cn-shanghai-g",
			"engine":      "hbaseue",
			"version":     "2.0",
		}),
	}

	oneConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"charge_type":   "Postpaid",
			"region_id":     "cn-shanghai",
			"zone_id":       "cn-shanghai-g",
			"engine":        "hbaseue",
			"version":       "2.0",
			"instance_type": "hbase.sn1.8xlarge",
			"disk_type":     "cloud_ssd",
		}),
	}

	hbaseInstanceTypesCheckInfo.DataSourceTestCheck(t, rand, oneConf, allConf)
}

func dataSourceHBaseInstanceTypesConfigDependence(name string) string {
	return ``
}

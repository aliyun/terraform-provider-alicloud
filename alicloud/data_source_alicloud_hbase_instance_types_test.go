package alicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

var (
	existHBaseInstanceTypesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"types.#":          CHECKSET,
			"types.0.cpu_size": CHECKSET,
			"types.0.mem_size": CHECKSET,
			"types.0.value":    CHECKSET,
		}
	}

	fakeHBaseInstanceTypesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"types.#": "0",
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
		fmt.Sprintf("tf-testacc%sbase-instance-types%v.abc", defaultRegionToTest, rand),
		dataSourceHBaseInstanceTypesConfigDependence)

	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"instance_type": "hbase.sn1.large",
		}),
	}

	oneConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"instance_type": "hbase.sn1.large",
		}),
	}

	hbaseInstanceTypesCheckInfo.dataSourceTestCheck(t, rand, oneConf, allConf)
}

func dataSourceHBaseInstanceTypesConfigDependence(name string) string {
	return ``
}

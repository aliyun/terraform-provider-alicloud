package hbase_test

import (
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAlicloudHBaseZonesDataSource_basic(t *testing.T) {
	rand := tfacctest.RandInt()
	resourceId := "data.alicloud_hbase_zones.default"

	var existHBaseZonesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   acctest.CHECKSET,
			"zones.#": acctest.CHECKSET,
		}
	}

	var fakeHBaseZonesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"zones.#": "0",
		}
	}

	var HBaseZonesCheckInfo = acctest.DataSourceAttr{
		ResourceId:   resourceId,
		ExistMapFunc: existHBaseZonesMapFunc,
		FakeMapFunc:  fakeHBaseZonesMapFunc,
	}

	// The flat-suite original passed zero configs, which left the TestCase with
	// no Steps and could never actually exercise the data source. Provide a
	// data source block so the zones query and its ids/zones outputs are
	// validated end-to-end.
	zonesConf := acctest.DataSourceTestAccConfig{
		ExistConfig: `data "alicloud_hbase_zones" "default" {}`,
	}

	HBaseZonesCheckInfo.DataSourceTestCheck(t, rand, zonesConf)
}

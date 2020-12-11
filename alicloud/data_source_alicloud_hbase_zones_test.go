package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudHBaseZonesDataSource_basic(t *testing.T) {
	rand := acctest.RandInt()
	resourceId := "data.alicloud_hbase_zones.default"

	var existHBaseZonesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   CHECKSET,
			"zones.#": CHECKSET,
		}
	}

	var fakeHBaseZonesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"zones.#": "0",
		}
	}

	var HBaseZonesCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existHBaseZonesMapFunc,
		fakeMapFunc:  fakeHBaseZonesMapFunc,
	}

	HBaseZonesCheckInfo.dataSourceTestCheck(t, rand)
}

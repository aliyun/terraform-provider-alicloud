package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudFCZonesDataSource_basic(t *testing.T) {
	rand := acctest.RandInt()
	resourceId := "data.alicloud_fc_zones.default"

	var fcZonesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":      CHECKSET,
			"ids.0":      CHECKSET,
			"zones.#":    CHECKSET,
			"zones.0.id": CHECKSET,
		}
	}

	var fakeFCZonesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"zones.#": "0",
		}
	}

	var fcZonesCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: fcZonesMapFunc,
		fakeMapFunc:  fakeFCZonesMapFunc,
	}

	fcZonesCheckInfo.dataSourceTestCheck(t, rand)
}

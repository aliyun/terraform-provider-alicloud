package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudEnhancedNatAvailableZones_basic(t *testing.T) {
	rand := acctest.RandInt()
	resourceId := "data.alicloud_enhanced_nat_available_zones.default"

	var existEnhancedNatAvailableZonesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   CHECKSET,
			"zones.#": CHECKSET,
		}
	}

	var fakeEnhancedNatAvailableZonesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"zones.#": "0",
		}
	}

	var EnhancedNatZonesCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existEnhancedNatAvailableZonesMapFunc,
		fakeMapFunc:  fakeEnhancedNatAvailableZonesMapFunc,
	}

	EnhancedNatZonesCheckInfo.dataSourceTestCheck(t, rand)
}

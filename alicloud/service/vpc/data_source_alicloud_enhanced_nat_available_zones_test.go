package vpc_test

import (
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAliCloudEnhancedNatAvailableZones_basic(t *testing.T) {
	rand := tfacctest.RandInt()
	resourceId := "data.alicloud_enhanced_nat_available_zones.default"

	var existEnhancedNatAvailableZonesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   acctest.CHECKSET,
			"zones.#": acctest.CHECKSET,
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

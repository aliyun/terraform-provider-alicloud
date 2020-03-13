package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
)

func TestAccAlicloudSlbZonesDataSource_basic(t *testing.T) {
	rand := acctest.RandInt()
	resourceId := "data.alicloud_slb_zones.default"

	testAccConfig := dataSourceTestAccConfigFunc(resourceId, "", dataSourceslbZonesConfigDependence)

	addressTypeConfig := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"available_slb_address_type": "Vpc",
		}),
	}

	ipVersionConfig := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"available_slb_address_ip_version": "ipv4",
		}),
	}

	allConfig := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"available_slb_address_type":       "Vpc",
			"available_slb_address_ip_version": "ipv4",
		}),
	}

	var existSlbZonesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                        CHECKSET,
			"zones.#":                      CHECKSET,
			"zones.0.slb_slave_zone_ids.#": CHECKSET,
		}
	}

	var fakeSlbZonesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"zones.#": "0",
		}
	}

	var slbZonesCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existSlbZonesMapFunc,
		fakeMapFunc:  fakeSlbZonesMapFunc,
	}

	slbZonesCheckInfo.dataSourceTestCheck(t, rand, addressTypeConfig, ipVersionConfig, allConfig)
}

func dataSourceslbZonesConfigDependence(name string) string {
	return ""
}

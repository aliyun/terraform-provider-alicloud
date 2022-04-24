package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudSLBZonesDataSource_basic(t *testing.T) {
	rand := acctest.RandInt()
	resourceId := "data.alicloud_slb_zones.default"

	testAccConfig := dataSourceTestAccConfigFunc(resourceId, "", dataSourceslbZonesConfigDependence)

	addressTypeConfig := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"available_slb_address_type": "vpc",
		}),
	}

	ipVersionConfig := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"available_slb_address_ip_version": "ipv4",
		}),
	}

	masterZoneIdConfig := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"master_zone_id": "${data.alicloud_zones.default.ids.0}",
		}),
	}

	slaveZoneIdConfig := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"slave_zone_id": "${data.alicloud_zones.default.ids.0}",
		}),
	}

	allConfig := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"available_slb_address_type":       "vpc",
			"available_slb_address_ip_version": "ipv4",
		}),
	}

	var existSlbZonesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                         CHECKSET,
			"zones.#":                       CHECKSET,
			"zones.0.master_zone_id":        CHECKSET,
			"zones.0.slave_zone_id":         CHECKSET,
			"zones.0.supported_resources.#": CHECKSET,
			"zones.0.supported_resources.0.address_type":       CHECKSET,
			"zones.0.supported_resources.0.address_ip_version": CHECKSET,
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

	slbZonesCheckInfo.dataSourceTestCheck(t, rand, addressTypeConfig, ipVersionConfig, masterZoneIdConfig, slaveZoneIdConfig, allConfig)
}

func dataSourceslbZonesConfigDependence(name string) string {
	return `
		data "alicloud_zones" "default" {
			available_resource_creation = "Slb"
		}
	`
}

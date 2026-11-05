package slb_test

import (
	"testing"

	acctest "github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAlicloudSLBZonesDataSource_basic(t *testing.T) {
	rand := tfacctest.RandInt()
	resourceId := "data.alicloud_slb_zones.default"

	testAccConfig := acctest.DataSourceTestAccConfigFunc(resourceId, "", dataSourceslbZonesConfigDependence)

	addressTypeConfig := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"available_slb_address_type": "vpc",
		}),
	}

	ipVersionConfig := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"available_slb_address_ip_version": "ipv4",
		}),
	}

	masterZoneIdConfig := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"master_zone_id": "${data.alicloud_zones.default.ids.0}",
		}),
	}

	slaveZoneIdConfig := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"slave_zone_id": "${data.alicloud_zones.default.ids.0}",
		}),
	}

	allConfig := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
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

	var slbZonesCheckInfo = acctest.DataSourceAttr{
		ResourceId:   resourceId,
		ExistMapFunc: existSlbZonesMapFunc,
		FakeMapFunc:  fakeSlbZonesMapFunc,
	}

	slbZonesCheckInfo.DataSourceTestCheck(t, rand, addressTypeConfig, ipVersionConfig, masterZoneIdConfig, slaveZoneIdConfig, allConfig)
}

func dataSourceslbZonesConfigDependence(name string) string {
	return `
		data "alicloud_zones" "default" {
			available_resource_creation = "Slb"
		}
	`
}

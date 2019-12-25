package alicloud

import (
	"fmt"
	"testing"

	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"

	"github.com/hashicorp/terraform/helper/acctest"
)

func TestAccAlicloudDRDSInstancesDataSource(t *testing.T) {
	rand := acctest.RandIntRange(1000000, 9999999)
	resourceId := "data.alicloud_drds_instances.default"

	testAccConfig := dataSourceTestAccConfigFunc(resourceId,
		fmt.Sprintf("tf-testAcc%sDRDSInstancesDataSource-%d", defaultRegionToTest, rand),
		dataSourceDRDSInstancesConfigDependence)

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex": alicloud_drds_instance.default.description,
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alicloud_drds_instance.default.description}-fake",
		}),
	}

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids": []string{alicloud_drds_instance.default.id},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_drds_instance.default.id}-fake"},
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex": alicloud_drds_instance.default.description,
			"ids":        []string{alicloud_drds_instance.default.id},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alicloud_drds_instance.default.description}-fake",
			"ids":        []string{"${alicloud_drds_instance.default.id}-fake"},
		}),
	}

	var existDRDSInstancesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                    "1",
			"descriptions.#":           "1",
			"ids.0":                    CHECKSET,
			"descriptions.0":           fmt.Sprintf("tf-testAcc%sDRDSInstancesDataSource-%d", defaultRegionToTest, rand),
			"instances.#":              "1",
			"instances.0.description":  fmt.Sprintf("tf-testAcc%sDRDSInstancesDataSource-%d", defaultRegionToTest, rand),
			"instances.0.type":         "PRIVATE",
			"instances.0.zone_id":      CHECKSET,
			"instances.0.id":           CHECKSET,
			"instances.0.network_type": "VPC",
			"instances.0.create_time":  CHECKSET,
		}
	}

	var fakeDRDSInstancesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":          "0",
			"descriptions.#": "0",
			"instances.#":    "0",
		}
	}

	var drdsInstancesCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existDRDSInstancesMapFunc,
		fakeMapFunc:  fakeDRDSInstancesMapFunc,
	}

	preCheck := func() {
		testAccPreCheckWithRegions(t, true, connectivity.DrdsSupportedRegions)
		testAccPreCheckWithAccountSiteType(t, DomesticSite)
	}

	drdsInstancesCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, nameRegexConf, idsConf, allConf)
}

func dataSourceDRDSInstancesConfigDependence(name string) string {
	return fmt.Sprintf(`
 	data "alicloud_zones" "default" {
		available_resource_creation = "VSwitch"
	}
 	variable "name" {
		default = "%s"
	}
 	resource "alicloud_vpc" "default" {
		name = var.name
		cidr_block = "172.16.0.0/12"
	}
 	resource "alicloud_vswitch" "default" {
 		vpc_id = alicloud_vpc.default.id
 		cidr_block = "172.16.0.0/21"
 		availability_zone = data.alicloud_zones.default.zones.0.id
 		name = var.name
	}

 	resource "alicloud_drds_instance" "default" {
  		description = var.name
  		zone_id = data.alicloud_zones.default.zones.0.id
  		instance_series = "drds.sn1.4c8g"
  		instance_charge_type = "PostPaid"
		vswitch_id = alicloud_vswitch.default.id
  		specification = "drds.sn1.4c8g.8C16G"
}
 `, name)
}

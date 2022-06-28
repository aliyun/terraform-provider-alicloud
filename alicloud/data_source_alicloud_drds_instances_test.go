package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudDRDSInstancesDataSource(t *testing.T) {
	rand := acctest.RandIntRange(0, 9999)
	resourceId := "data.alicloud_drds_instances.default"
	name := fmt.Sprintf("tf-testAcc%s%d", defaultRegionToTest, rand)
	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceDRDSInstancesConfigDependence)

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alicloud_drds_instance.default.description}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alicloud_drds_instance.default.description}-fake",
		}),
	}

	descriptionRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"description_regex": "${alicloud_drds_instance.default.description}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"description_regex": "${alicloud_drds_instance.default.description}-fake",
		}),
	}

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_drds_instance.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_drds_instance.default.id}-fake"},
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex":        "${alicloud_drds_instance.default.description}",
			"description_regex": "${alicloud_drds_instance.default.description}",
			"ids":               []string{"${alicloud_drds_instance.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex":        "${alicloud_drds_instance.default.description}-fake",
			"description_regex": "${alicloud_drds_instance.default.description}-fake",
			"ids":               []string{"${alicloud_drds_instance.default.id}-fake"},
		}),
	}

	var existDRDSInstancesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                    "1",
			"descriptions.#":           "1",
			"ids.0":                    CHECKSET,
			"instances.#":              "1",
			"instances.0.description":  name,
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
	}

	drdsInstancesCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, nameRegexConf, descriptionRegexConf, idsConf, allConf)
}

func dataSourceDRDSInstancesConfigDependence(name string) string {
	return fmt.Sprintf(`
 	data "alicloud_zones" "default" {
		available_resource_creation = "VSwitch"
	}
 	variable "name" {
		default = "%s"
	}
	data "alicloud_vpcs" "default"	{
        name_regex = "default-NODELETING"
	}
	data "alicloud_vswitches" "default" {
	  vpc_id = "${data.alicloud_vpcs.default.ids.0}"
      zone_id = data.alicloud_zones.default.ids.0
	}
 	resource "alicloud_drds_instance" "default" {
  		description = "${var.name}"
  		zone_id = "${data.alicloud_vswitches.default.vswitches.0.zone_id}"
  		instance_series = "drds.sn1.4c8g"
  		instance_charge_type = "PostPaid"
		vswitch_id = "${data.alicloud_vswitches.default.vswitches.0.id}"
  		specification = "drds.sn1.4c8g.8C16G"
}
 `, name)
}

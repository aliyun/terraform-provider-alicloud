package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudEaisInstancesDataSource(t *testing.T) {

	rand := acctest.RandInt()
	resourceId := "data.alicloud_eais_instances.default"
	name := fmt.Sprintf("tf-testacc-eaisInstance%v", rand)

	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceEaisInstancesConfigDependence)

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alicloud_eais_instance.default.instance_name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": "fake_tf-testacc*",
		}),
	}

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_eais_instance.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_eais_instance.default.id}_fake"},
		}),
	}

	instanceTypeConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":           []string{"${alicloud_eais_instance.default.id}"},
			"instance_type": "eais.ei-a6.medium",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":           []string{"${alicloud_eais_instance.default.id}_fake"},
			"instance_type": "eais.ei-a6.2xlarge",
		}),
	}

	statusConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":    []string{"${alicloud_eais_instance.default.id}"},
			"status": "Available",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":    []string{"${alicloud_eais_instance.default.id}_fake"},
			"status": "Unavailable",
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex":    "${alicloud_eais_instance.default.instance_name}",
			"ids":           []string{"${alicloud_eais_instance.default.id}"},
			"status":        "Available",
			"instance_type": "eais.ei-a6.medium",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex":    "${alicloud_eais_instance.default.instance_name}",
			"ids":           []string{"${alicloud_eais_instance.default.id}"},
			"status":        "Unavailable",
			"instance_type": "eais.ei-a6.2xlarge",
		}),
	}
	var existEaisInstancesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                     "1",
			"names.#":                   "1",
			"instances.#":               "1",
			"instances.0.instance_name": name,
		}
	}

	var fakeEaisInstancesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"instances.#": "0",
			"names.#":     "0",
			"ids.#":       "0",
		}
	}

	var EaisInstancesCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existEaisInstancesMapFunc,
		fakeMapFunc:  fakeEaisInstancesMapFunc,
	}
	preCheck := func() {
		testAccPreCheck(t)
		testAccPreCheckWithRegions(t, true, connectivity.EAISSystemSupportRegions)
	}
	EaisInstancesCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, nameRegexConf, idsConf, instanceTypeConf, statusConf, allConf)
}

func dataSourceEaisInstancesConfigDependence(name string) string {
	return fmt.Sprintf(`
		variable "name" {
		 default = "%v"
		}
		data "alicloud_zones" "default" {
			available_resource_creation = "VSwitch"
		}
		data "alicloud_vpcs" "default"{
			name_regex = "default-NODELETING"
		}
		data "alicloud_vswitches" "default" {
		  vpc_id  = data.alicloud_vpcs.default.ids.0
          zone_id = data.alicloud_zones.default.ids.0
		}
		
		resource "alicloud_security_group" "default" {
		  name        = var.name
		  description = "tf test"
		  vpc_id      = data.alicloud_vpcs.default.ids.0
		}
		resource "alicloud_eais_instance" "default" {
		  instance_type     = "eais.ei-a6.medium"
		  instance_name     = var.name
		  security_group_id = alicloud_security_group.default.id
		  vswitch_id        = data.alicloud_vswitches.default.ids.0
		}
		`, name)
}

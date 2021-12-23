package alicloud

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudGaEndpointGroupsDataSource(t *testing.T) {
	rand := acctest.RandInt()
	resourceId := "data.alicloud_ga_endpoint_groups.default"
	name := fmt.Sprintf("tf-testEndpointGroups_datasource-%d", rand)
	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceGaEndpointGroupsConfigDependence)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"accelerator_id": "${data.alicloud_ga_accelerators.default.ids.0}",
			"ids":            []string{"${alicloud_ga_endpoint_group.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"accelerator_id": "${data.alicloud_ga_accelerators.default.ids.0}",
			"ids":            []string{"${alicloud_ga_endpoint_group.default.id}_fake"},
		}),
	}

	statusConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"accelerator_id": "${data.alicloud_ga_accelerators.default.ids.0}",
			"ids":            []string{"${alicloud_ga_endpoint_group.default.id}"},
			"status":         "active",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"accelerator_id": "${data.alicloud_ga_accelerators.default.ids.0}",
			"ids":            []string{"${alicloud_ga_endpoint_group.default.id}_fake"},
			"status":         "creating",
		}),
	}

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"accelerator_id": "${data.alicloud_ga_accelerators.default.ids.0}",
			"name_regex":     "${alicloud_ga_endpoint_group.default.name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"accelerator_id": "${data.alicloud_ga_accelerators.default.ids.0}",
			"name_regex":     "${alicloud_ga_endpoint_group.default.name}_fake",
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"accelerator_id": "${data.alicloud_ga_accelerators.default.ids.0}",
			"ids":            []string{"${alicloud_ga_endpoint_group.default.id}"},
			"name_regex":     "${alicloud_ga_endpoint_group.default.name}",
			"status":         "active",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"accelerator_id": "${data.alicloud_ga_accelerators.default.ids.0}",
			"ids":            []string{"${alicloud_ga_endpoint_group.default.id}_fake"},
			"name_regex":     "${alicloud_ga_endpoint_group.default.name}_fake",
			"status":         "creating",
		}),
	}

	var existMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                                  "1",
			"names.#":                                "1",
			"groups.#":                               CHECKSET,
			"groups.0.description":                   name + "_desc",
			"groups.0.endpoint_configurations.#":     "1",
			"groups.0.id":                            CHECKSET,
			"groups.0.endpoint_group_id":             CHECKSET,
			"groups.0.endpoint_group_region":         "cn-hangzhou",
			"groups.0.health_check_interval_seconds": "3",
			"groups.0.health_check_path":             "/healthcheck",
			"groups.0.health_check_port":             "9999",
			"groups.0.health_check_protocol":         "http",
			"groups.0.listener_id":                   CHECKSET,
			"groups.0.name":                          name,
			"groups.0.port_overrides.#":              "1",
			"groups.0.status":                        "active",
			"groups.0.threshold_count":               "4",
			"groups.0.traffic_percentage":            "20",
		}
	}

	var fakeMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"groups.#": "0",
			"ids.#":    "0",
		}
	}

	var CheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existMapFunc,
		fakeMapFunc:  fakeMapFunc,
	}
	preCheck := func() {}

	CheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, statusConf, nameRegexConf, allConf)
}

func dataSourceGaEndpointGroupsConfigDependence(name string) string {
	return fmt.Sprintf(`data "alicloud_ga_accelerators" "default"{
  
}
resource "alicloud_ga_endpoint_group" "default" {
  accelerator_id="${data.alicloud_ga_accelerators.default.ids.0}"
  endpoint_configurations{
    endpoint=alicloud_eip_address.example.ip_address
    type="PublicIp"
    weight="20"
  }
  description="%s_desc"
  name="%s"
  threshold_count=4
  endpoint_group_region="%s"
  health_check_interval_seconds="3"
  health_check_path="/healthcheck"
  health_check_port="9999"
  health_check_protocol="http"
  port_overrides{
    endpoint_port="10"
    listener_port="60"
  }
  traffic_percentage=20
  listener_id="${alicloud_ga_listener.default.id}"
}
resource "alicloud_eip_address" "example" {
  bandwidth            = "10"
  internet_charge_type = "PayByBandwidth"
}
resource "alicloud_ga_listener" "default" {
  port_ranges{
    from_port="60"
    to_port="70"
  }
  accelerator_id="${data.alicloud_ga_accelerators.default.ids.0}"
  client_affinity="SOURCE_IP"
  protocol="UDP"
  name="%s"
}`, name, name, os.Getenv("ALICLOUD_REGION"), name)
}

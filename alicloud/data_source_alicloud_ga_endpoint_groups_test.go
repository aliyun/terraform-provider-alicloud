package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudGaEndpointGroupsDataSource(t *testing.T) {
	rand := acctest.RandInt()
	resourceId := "data.alicloud_ga_endpoint_groups.default"
	name := fmt.Sprintf("tf-testEndpointGroups_datasource-%d", rand)
	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceGaEndpointGroupsConfigDependence)
	checkoutSupportedRegions(t, true, connectivity.GaSupportRegions)
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"accelerator_id": "${alicloud_ga_endpoint_group.default.accelerator_id}",
			"ids":            []string{"${alicloud_ga_endpoint_group.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"accelerator_id": "${alicloud_ga_endpoint_group.default.accelerator_id}",
			"ids":            []string{"${alicloud_ga_endpoint_group.default.id}_fake"},
		}),
	}

	statusConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"accelerator_id": "${alicloud_ga_endpoint_group.default.accelerator_id}",
			"ids":            []string{"${alicloud_ga_endpoint_group.default.id}"},
			"status":         "active",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"accelerator_id": "${alicloud_ga_endpoint_group.default.accelerator_id}",
			"ids":            []string{"${alicloud_ga_endpoint_group.default.id}_fake"},
			"status":         "creating",
		}),
	}

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"accelerator_id": "${alicloud_ga_endpoint_group.default.accelerator_id}",
			"name_regex":     "${alicloud_ga_endpoint_group.default.name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"accelerator_id": "${alicloud_ga_endpoint_group.default.accelerator_id}",
			"name_regex":     "${alicloud_ga_endpoint_group.default.name}_fake",
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"accelerator_id": "${alicloud_ga_endpoint_group.default.accelerator_id}",
			"ids":            []string{"${alicloud_ga_endpoint_group.default.id}"},
			"name_regex":     "${alicloud_ga_endpoint_group.default.name}",
			"status":         "active",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"accelerator_id": "${alicloud_ga_endpoint_group.default.accelerator_id}",
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
			"groups.0.description":                   name,
			"groups.0.endpoint_configurations.#":     "1",
			"groups.0.id":                            CHECKSET,
			"groups.0.endpoint_group_id":             CHECKSET,
			"groups.0.endpoint_group_region":         defaultRegionToTest,
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
	preCheck := func() {
		testAccPreCheck(t)

	}

	CheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, statusConf, nameRegexConf, allConf)
}

func dataSourceGaEndpointGroupsConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default  = "%s"
}

data "alicloud_ga_accelerators" "default" {
  status = "active"
}

resource "alicloud_ga_bandwidth_package" "default" {
   	bandwidth              =  100
  	type                   = "Basic"
  	bandwidth_type         = "Basic"
	payment_type           = "PayAsYouGo"
  	billing_type           = "PayBy95"
	ratio       = 30
	bandwidth_package_name = var.name
    auto_pay               = true
    auto_use_coupon        = true
}

resource "alicloud_ga_bandwidth_package_attachment" "default" {
	// Please run resource ga_accelerator test case to ensure this account has at least one accelerator before run this case.
	accelerator_id = data.alicloud_ga_accelerators.default.ids.0
	bandwidth_package_id = alicloud_ga_bandwidth_package.default.id
}

resource "alicloud_ga_listener" "default" {
  port_ranges{
    from_port="60"
    to_port="70"
  }
  accelerator_id=alicloud_ga_bandwidth_package_attachment.default.accelerator_id
  client_affinity="SOURCE_IP"
  protocol="UDP"
  name=var.name
}

resource "alicloud_eip_address" "default" {
  bandwidth            = "10"
  internet_charge_type = "PayByBandwidth"
  address_name = var.name
}

resource "alicloud_ga_endpoint_group" "default" {
  accelerator_id=alicloud_ga_listener.default.accelerator_id
  endpoint_configurations{
    endpoint=alicloud_eip_address.default.ip_address
    type="PublicIp"
    weight="20"
  }
  description=var.name
  name=var.name
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
  listener_id=alicloud_ga_listener.default.id
}
`, name, defaultRegionToTest)
}

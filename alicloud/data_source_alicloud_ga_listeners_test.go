package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudGaListenersDataSource(t *testing.T) {
	rand := acctest.RandInt()
	checkoutSupportedRegions(t, true, connectivity.GaSupportRegions)
	resourceId := "data.alicloud_ga_listeners.default"
	name := fmt.Sprintf("tf-testListeners_datasource-%d", rand)
	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceGaListenersConfigDependence)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"accelerator_id": "${alicloud_ga_listener.default.accelerator_id}",
			"ids":            []string{"${alicloud_ga_listener.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"accelerator_id": "${alicloud_ga_listener.default.accelerator_id}",
			"ids":            []string{"${alicloud_ga_listener.default.id}_fake"},
		}),
	}

	statusConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"accelerator_id": "${alicloud_ga_listener.default.accelerator_id}",
			"ids":            []string{"${alicloud_ga_listener.default.id}"},
			"status":         "active",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"accelerator_id": "${alicloud_ga_listener.default.accelerator_id}",
			"ids":            []string{"${alicloud_ga_listener.default.id}_fake"},
			"status":         "creating",
		}),
	}

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"accelerator_id": "${alicloud_ga_listener.default.accelerator_id}",
			"name_regex":     "${alicloud_ga_listener.default.name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"accelerator_id": "${alicloud_ga_listener.default.accelerator_id}",
			"name_regex":     "${alicloud_ga_listener.default.name}_fake",
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"accelerator_id": "${alicloud_ga_listener.default.accelerator_id}",
			"ids":            []string{"${alicloud_ga_listener.default.id}"},
			"name_regex":     "${alicloud_ga_listener.default.name}",
			"status":         "active",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"accelerator_id": "${alicloud_ga_listener.default.accelerator_id}",
			"ids":            []string{"${alicloud_ga_listener.default.id}_fake"},
			"name_regex":     "${alicloud_ga_listener.default.name}_fake",
			"status":         "creating",
		}),
	}

	var existMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                       "1",
			"names.#":                     "1",
			"listeners.#":                 CHECKSET,
			"listeners.0.certificates.#":  "0",
			"listeners.0.client_affinity": "NONE",
			"listeners.0.description":     "create_description",
			"listeners.0.id":              CHECKSET,
			"listeners.0.name":            fmt.Sprintf("tf-testListeners_datasource-%d", rand),
			"listeners.0.port_ranges.#":   "1",
			"listeners.0.protocol":        "TCP",
			"listeners.0.status":          "active",
		}
	}

	var fakeMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"accelerators.#": "0",
			"ids.#":          "0",
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

func dataSourceGaListenersConfigDependence(name string) string {
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
resource "alicloud_ga_listener" "default"{
  port_ranges{
    from_port = "80"
    to_port   = "90"
  }
  accelerator_id = alicloud_ga_bandwidth_package_attachment.default.accelerator_id
  name           =var.name
  description    ="create_description"
}`, name)
}

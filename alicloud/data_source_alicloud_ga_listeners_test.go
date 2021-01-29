package alicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudGaListenersDataSource(t *testing.T) {
	rand := acctest.RandInt()
	resourceId := "data.alicloud_ga_listeners.default"
	name := fmt.Sprintf("tf-testListeners_datasource-%d", rand)
	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceGaListenersConfigDependence)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"accelerator_id": "${data.alicloud_ga_accelerators.default.ids.0}",
			"ids":            []string{"${alicloud_ga_listener.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"accelerator_id": "${data.alicloud_ga_accelerators.default.ids.0}",
			"ids":            []string{"${alicloud_ga_listener.default.id}_fake"},
		}),
	}

	statusConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"accelerator_id": "${data.alicloud_ga_accelerators.default.ids.0}",
			"ids":            []string{"${alicloud_ga_listener.default.id}"},
			"status":         "active",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"accelerator_id": "${data.alicloud_ga_accelerators.default.ids.0}",
			"ids":            []string{"${alicloud_ga_listener.default.id}_fake"},
			"status":         "creating",
		}),
	}

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"accelerator_id": "${data.alicloud_ga_accelerators.default.ids.0}",
			"name_regex":     "${alicloud_ga_listener.default.name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"accelerator_id": "${data.alicloud_ga_accelerators.default.ids.0}",
			"name_regex":     "${alicloud_ga_listener.default.name}_fake",
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"accelerator_id": "${data.alicloud_ga_accelerators.default.ids.0}",
			"ids":            []string{"${alicloud_ga_listener.default.id}"},
			"name_regex":     "${alicloud_ga_listener.default.name}",
			"status":         "active",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"accelerator_id": "${data.alicloud_ga_accelerators.default.ids.0}",
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
	return fmt.Sprintf(`data "alicloud_ga_accelerators" "default"{
}
resource "alicloud_ga_listener" "default"{
  port_ranges{
    from_port = "80"
    to_port   = "90"
  }
  accelerator_id = "${data.alicloud_ga_accelerators.default.ids.0}"
  name           ="%s"
  description    ="create_description"
}`, name)
}

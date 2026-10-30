package ga_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAliCloudGaAcceleratorsDataSource(t *testing.T) {
	rand := tfacctest.RandInt()
	acctest.CheckoutSupportedRegions(t, true, connectivity.GaSupportRegions)
	resourceId := "data.alicloud_ga_accelerators.default"
	name := fmt.Sprintf("tf-testAccelerators_datasource-%d", rand)
	testAccConfig := acctest.DataSourceTestAccConfigFunc(resourceId, name, dataSourceGaAcceleratorsConfigDependence)

	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_ga_accelerator.default.id}"},
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_ga_accelerator.default.id}_fake"},
		}),
	}

	nameRegexConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alicloud_ga_accelerator.default.accelerator_name}",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alicloud_ga_accelerator.default.accelerator_name}_fake",
		}),
	}

	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"ids":        []string{"${alicloud_ga_accelerator.default.id}"},
			"name_regex": "${alicloud_ga_accelerator.default.accelerator_name}",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"ids":        []string{"${alicloud_ga_accelerator.default.id}_fake"},
			"name_regex": "${alicloud_ga_accelerator.default.accelerator_name}_fake",
		}),
	}

	var existMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                                           "1",
			"names.#":                                         "1",
			"accelerators.#":                                  acctest.CHECKSET,
			"accelerators.0.id":                               acctest.CHECKSET,
			"accelerators.0.accelerator_name":                 name,
			"accelerators.0.cen_id":                           "",
			"accelerators.0.ddos_id":                          "",
			"accelerators.0.description":                      name,
			"accelerators.0.dns_name":                         acctest.CHECKSET,
			"accelerators.0.expired_time":                     acctest.CHECKSET,
			"accelerators.0.payment_type":                     "PREPAY",
			"accelerators.0.second_dns_name":                  "",
			"accelerators.0.spec":                             "1",
			"accelerators.0.status":                           "active",
			"accelerators.0.cross_domain_bandwidth_package.#": "0",
			"accelerators.0.basic_bandwidth_package.#":        "0",
		}
	}

	var fakeMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"accelerators.#": "0",
			"ids.#":          "0",
		}
	}

	var CheckInfo = acctest.DataSourceAttr{
		ResourceId:   resourceId,
		ExistMapFunc: existMapFunc,
		FakeMapFunc:  fakeMapFunc,
	}

	CheckInfo.DataSourceTestCheck(t, rand, idsConf, nameRegexConf, allConf)
}

func TestAccAliCloudGaAcceleratorsDataSource_withAnycast(t *testing.T) {
	rand := tfacctest.RandInt()
	acctest.CheckoutSupportedRegions(t, true, connectivity.GaSupportRegions)
	resourceId := "data.alicloud_ga_accelerators.default"
	name := fmt.Sprintf("tf-testAccelerators-anycast-%d", rand)
	testAccConfig := acctest.DataSourceTestAccConfigFunc(resourceId, name, dataSourceGaAcceleratorsConfigDependenceWithAnycast)

	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_ga_accelerator.default.id}"},
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_ga_accelerator.default.id}_fake"},
		}),
	}

	var existMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                                      "1",
			"accelerators.#":                             acctest.CHECKSET,
			"accelerators.0.id":                          acctest.CHECKSET,
			"accelerators.0.payment_type":                "POSTPAY",
			"accelerators.0.bandwidth_billing_type":      "CDT",
			"accelerators.0.bandwidth":                   "200",
			"accelerators.0.ip_set_config.#":             "1",
			"accelerators.0.ip_set_config.0.access_mode": "Anycast",
		}
	}

	var fakeMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"accelerators.#": "0",
			"ids.#":          "0",
		}
	}

	var CheckInfo = acctest.DataSourceAttr{
		ResourceId:   resourceId,
		ExistMapFunc: existMapFunc,
		FakeMapFunc:  fakeMapFunc,
	}

	CheckInfo.DataSourceTestCheck(t, rand, idsConf)
}

func dataSourceGaAcceleratorsConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}
resource "alicloud_ga_accelerator" "default" {
  duration         = 1
  spec             = "1"
  accelerator_name = var.name
  auto_use_coupon  = true
  description      = var.name
}
`, name)
}

func dataSourceGaAcceleratorsConfigDependenceWithAnycast(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}
resource "alicloud_ga_accelerator" "default" {
  bandwidth_billing_type = "CDT"
  payment_type           = "PayAsYouGo"
  bandwidth              = 200
  ip_set_config {
    access_mode = "Anycast"
  }
}
`, name)
}

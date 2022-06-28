package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudGaBandwidthPackagesDataSource(t *testing.T) {
	rand := acctest.RandInt()
	checkoutSupportedRegions(t, true, connectivity.GaSupportRegions)
	resourceId := "data.alicloud_ga_bandwidth_packages.default"
	name := fmt.Sprintf("tf-testBandwidthPackages_datasource-%d", rand)
	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceGaBandwidthPackagesConfigDependence)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":            []string{"${alicloud_ga_bandwidth_package.default.id}"},
			"enable_details": "true",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_ga_bandwidth_package.default.id}_fake"},
		}),
	}
	statusConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":            []string{"${alicloud_ga_bandwidth_package.default.id}"},
			"status":         "active",
			"enable_details": "true",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":    []string{"${alicloud_ga_bandwidth_package.default.id}_fake"},
			"status": "init",
		}),
	}
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex":     "${alicloud_ga_bandwidth_package.default.bandwidth_package_name}",
			"enable_details": "true",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alicloud_ga_bandwidth_package.default.bandwidth_package_name}_fake",
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":            []string{"${alicloud_ga_bandwidth_package.default.id}"},
			"name_regex":     "${alicloud_ga_bandwidth_package.default.bandwidth_package_name}",
			"status":         "active",
			"enable_details": "true",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":        []string{"${alicloud_ga_bandwidth_package.default.id}_fake"},
			"name_regex": "${alicloud_ga_bandwidth_package.default.bandwidth_package_name}_fake",
			"status":     "init",
		}),
	}

	var existMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                                "1",
			"names.#":                              "1",
			"packages.#":                           CHECKSET,
			"packages.0.bandwidth":                 "100",
			"packages.0.id":                        CHECKSET,
			"packages.0.bandwidth_package_id":      CHECKSET,
			"packages.0.bandwidth_package_name":    name,
			"packages.0.bandwidth_type":            "Basic",
			"packages.0.cbn_geographic_region_ida": "",
			"packages.0.cbn_geographic_region_idb": "",
			"packages.0.description":               "",
			"packages.0.expired_time":              CHECKSET,
			"packages.0.payment_type":              "PayAsYouGo",
			"packages.0.status":                    "active",
			"packages.0.type":                      "Basic",
		}
	}

	var fakeMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"packages.#": "0",
			"ids.#":      "0",
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

func dataSourceGaBandwidthPackagesConfigDependence(name string) string {
	return fmt.Sprintf(`
resource "alicloud_ga_bandwidth_package" "default" {
   	bandwidth              =  100
  	type                   = "Basic"
  	bandwidth_type         = "Basic"
	payment_type           = "PayAsYouGo"
  	billing_type           = "PayBy95"
	ratio       = 30
	bandwidth_package_name = "%s"
}`, name)
}

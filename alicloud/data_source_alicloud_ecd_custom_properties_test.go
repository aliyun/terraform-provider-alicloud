package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudEcdCustomPropertiesDataSource(t *testing.T) {
	rand := acctest.RandInt()
	checkoutSupportedRegions(t, true, connectivity.EcdUserSupportRegions)
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEcdCustomPropertiesDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_ecd_custom_property.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudEcdCustomPropertiesDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_ecd_custom_property.default.id}_fake"]`,
		}),
	}
	var existAlicloudEcdCustomPropertiesDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                          "1",
			"properties.#":                   "1",
			"properties.0.property_key":      fmt.Sprintf("tf-testAccCustomProperty-%d", rand),
			"properties.0.property_values.#": "1",
			"properties.0.property_values.0.property_value_id": CHECKSET,
			"properties.0.property_values.0.property_value":    fmt.Sprintf("tf-testAccCustomProperty-%d", rand),
			"properties.0.id":                 CHECKSET,
			"properties.0.custom_property_id": CHECKSET,
		}
	}
	var fakeAlicloudEcdCustomPropertiesDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#": "0",
		}
	}
	var alicloudEcdCustomPropertiesCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_ecd_custom_properties.default",
		existMapFunc: existAlicloudEcdCustomPropertiesDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlicloudEcdCustomPropertiesDataSourceNameMapFunc,
	}

	preCheck := func() {
		testAccPreCheck(t)
	}
	alicloudEcdCustomPropertiesCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf)
}
func testAccCheckAlicloudEcdCustomPropertiesDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`

variable "name" {	
	default = "tf-testAccCustomProperty-%d"
}

resource "alicloud_ecd_custom_property" "default" {
	property_key = var.name
	property_values {
		property_value = var.name
	}
}

data "alicloud_ecd_custom_properties" "default" {	
	%s
}
`, rand, strings.Join(pairs, " \n "))
	return config
}

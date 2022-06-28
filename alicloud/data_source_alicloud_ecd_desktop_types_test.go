package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudECDDesktopTypesDataSource(t *testing.T) {
	rand := acctest.RandIntRange(10000, 99999)
	checkoutSupportedRegions(t, true, connectivity.EcdSupportRegions)
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEcdDesktopTypesDataSourceName(rand, map[string]string{
			"instance_type_family": `"eds.hf"`,
			"cpu_count":            `"4"`,
			"memory_size":          `"8192"`,
			"gpu_count":            `"0"`,
		}),
		fakeConfig: "",
	}
	var existAlicloudEcdDesktopTypesDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   CHECKSET,
			"types.#": CHECKSET,
		}
	}
	var fakeAlicloudEcdDesktopTypesDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"types.#": "0",
		}
	}
	var alicloudEcdDesktopTypesBusesCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_ecd_desktop_types.default",
		existMapFunc: existAlicloudEcdDesktopTypesDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlicloudEcdDesktopTypesDataSourceNameMapFunc,
	}

	alicloudEcdDesktopTypesBusesCheckInfo.dataSourceTestCheck(t, rand, idsConf)
}

func testAccCheckAlicloudEcdDesktopTypesDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`

variable "name" {
  default = "tf-testaccdesktoptypes%d"
}
data "alicloud_ecd_desktop_types" "default" {	
	%s
}
`, rand, strings.Join(pairs, " \n "))
	return config
}

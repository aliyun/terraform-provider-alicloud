package ecd_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAlicloudECDDesktopTypesDataSource(t *testing.T) {
	rand := tfacctest.RandIntRange(10000, 99999)
	acctest.CheckoutSupportedRegions(t, true, connectivity.EcdSupportRegions)
	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudEcdDesktopTypesDataSourceName(rand, map[string]string{
			"instance_type_family": `"eds.hf"`,
			"cpu_count":            `"4"`,
			"memory_size":          `"8192"`,
			"gpu_count":            `"0"`,
		}),
		FakeConfig: "",
	}
	var existAlicloudEcdDesktopTypesDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   acctest.CHECKSET,
			"types.#": acctest.CHECKSET,
		}
	}
	var fakeAlicloudEcdDesktopTypesDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"types.#": "0",
		}
	}
	var alicloudEcdDesktopTypesBusesCheckInfo = acctest.DataSourceAttr{
		ResourceId:   "data.alicloud_ecd_desktop_types.default",
		ExistMapFunc: existAlicloudEcdDesktopTypesDataSourceNameMapFunc,
		FakeMapFunc:  fakeAlicloudEcdDesktopTypesDataSourceNameMapFunc,
	}

	alicloudEcdDesktopTypesBusesCheckInfo.DataSourceTestCheck(t, rand, idsConf)
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

package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudCloudControlProductDataSource(t *testing.T) {
	testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
	rand := acctest.RandIntRange(1000000, 9999999)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCloudControlProductSourceConfig(rand, map[string]string{
			"ids": `["Live"]`,
		}),
		fakeConfig: testAccCheckAlicloudCloudControlProductSourceConfig(rand, map[string]string{
			"ids": `["Live_fake"]`,
		}),
	}

	CloudControlProductCheckInfo.dataSourceTestCheck(t, rand, idsConf)
}

var existCloudControlProductMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"products.#":              "1",
		"products.0.product_name": CHECKSET,
		"products.0.product_code": CHECKSET,
	}
}

var fakeCloudControlProductMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"products.#": "0",
	}
}

var CloudControlProductCheckInfo = dataSourceAttr{
	resourceId:   "data.alicloud_cloud_control_products.default",
	existMapFunc: existCloudControlProductMapFunc,
	fakeMapFunc:  fakeCloudControlProductMapFunc,
}

func testAccCheckAlicloudCloudControlProductSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
	default = "tf-testAccCloudControlProduct%d"
}


data "alicloud_cloud_control_products" "default" {
%s
}
`, rand, strings.Join(pairs, "\n   "))
	return config
}

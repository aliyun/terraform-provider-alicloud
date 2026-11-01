package apigateway_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"

	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAlicloudApiGatewayBackendsDataSource(t *testing.T) {
	rand := tfacctest.RandInt()
	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudApiGatewayBackendsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_api_gateway_backend.default.id}"]`,
		}),
		FakeConfig: testAccCheckAlicloudApiGatewayBackendsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_api_gateway_backend.default.id}_fake"]`,
		}),
	}
	nameRegexConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudApiGatewayBackendsDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_api_gateway_backend.default.backend_name}"`,
		}),
		FakeConfig: testAccCheckAlicloudApiGatewayBackendsDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_api_gateway_backend.default.backend_name}_fake"`,
		}),
	}
	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudApiGatewayBackendsDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_api_gateway_backend.default.id}"]`,
			"name_regex": `"${alicloud_api_gateway_backend.default.backend_name}"`,
		}),
		FakeConfig: testAccCheckAlicloudApiGatewayBackendsDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_api_gateway_backend.default.id}_fake"]`,
			"name_regex": `"${alicloud_api_gateway_backend.default.backend_name}_fake"`,
		}),
	}
	var existAlicloudApiGatewayBackendsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                    "1",
			"names.#":                  "1",
			"backends.#":               "1",
			"backends.0.id":            acctest.CHECKSET,
			"backends.0.backend_id":    acctest.CHECKSET,
			"backends.0.backend_type":  "HTTP",
			"backends.0.backend_name":  fmt.Sprintf("tf-testAccBackend-%d", rand),
			"backends.0.create_time":   acctest.CHECKSET,
			"backends.0.description":   fmt.Sprintf("tf-testAccBackend-%d", rand),
			"backends.0.modified_time": acctest.CHECKSET,
		}
	}
	var fakeAlicloudApiGatewayBackendsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"names.#": "0",
		}
	}
	var alicloudApiGatewayBackendsCheckInfo = acctest.DataSourceAttr{
		ResourceId:   "data.alicloud_api_gateway_backends.default",
		ExistMapFunc: existAlicloudApiGatewayBackendsDataSourceNameMapFunc,
		FakeMapFunc:  fakeAlicloudApiGatewayBackendsDataSourceNameMapFunc,
	}

	preCheck := func() {
		acctest.PreCheck(t)
	}
	alicloudApiGatewayBackendsCheckInfo.DataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, nameRegexConf, allConf)
}
func testAccCheckAlicloudApiGatewayBackendsDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`

variable "name" {	
	default = "tf-testAccBackend-%d"
}

resource "alicloud_api_gateway_backend" "default" {
 backend_name = var.name
 description  = var.name
 backend_type = "HTTP"
}

data "alicloud_api_gateway_backends" "default" {	
	%s
}
`, rand, strings.Join(pairs, " \n "))
	return config
}

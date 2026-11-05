package apig_test

import (
	"fmt"
	"strings"
	"testing"

	acctest "github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAliCloudApigAiModelProvidersDataSource(t *testing.T) {
	acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
	rand := tfacctest.RandIntRange(10000, 99999)

	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAliCloudApigAiModelProvidersSourceConfig(rand, map[string]string{
			"gateway_id": `"${alicloud_apig_ai_model_provider.default.gateway_id}"`,
			"ids":        `["${alicloud_apig_ai_model_provider.default.id}"]`,
		}),
		FakeConfig: testAccCheckAliCloudApigAiModelProvidersSourceConfig(rand, map[string]string{
			"gateway_id": `"${alicloud_apig_ai_model_provider.default.gateway_id}"`,
			"ids":        `["${alicloud_apig_ai_model_provider.default.id}_fake"]`,
		}),
	}

	ApigAiModelProvidersCheckInfo.DataSourceTestCheck(t, rand, idsConf)
}

var existApigAiModelProvidersMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"ids.#":                         "1",
		"providers.#":                   "1",
		"providers.0.id":                acctest.CHECKSET,
		"providers.0.model_provider_id": acctest.CHECKSET,
		"providers.0.gateway_id":        acctest.CHECKSET,
		"providers.0.model_provider":    "openai",
		"providers.0.display_name":      fmt.Sprintf("tfaccapigaimp%d", rand),
		"providers.0.source":            acctest.CHECKSET,
		"providers.0.update_time":       acctest.CHECKSET,
	}
}

var fakeApigAiModelProvidersMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"ids.#":       "0",
		"providers.#": "0",
	}
}

var ApigAiModelProvidersCheckInfo = acctest.DataSourceAttr{
	ResourceId:   "data.alicloud_apig_ai_model_providers.default",
	ExistMapFunc: existApigAiModelProvidersMapFunc,
	FakeMapFunc:  fakeApigAiModelProvidersMapFunc,
}

func testAccCheckAliCloudApigAiModelProvidersSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	name := fmt.Sprintf("tfaccapigaimp%d", rand)
	return AlicloudApigAiModelProviderBasicDependence(name) + fmt.Sprintf(`
resource "alicloud_apig_ai_model_provider" "default" {
  gateway_id     = alicloud_apig_gateway.default.id
  model_provider = "openai"
  display_name   = "%s"
}

data "alicloud_apig_ai_model_providers" "default" {
  %s
}
`, name, strings.Join(pairs, "\n  "))
}

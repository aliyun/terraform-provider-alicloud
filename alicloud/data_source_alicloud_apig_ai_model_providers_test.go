package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAliCloudApigAiModelProvidersDataSource(t *testing.T) {
	testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
	rand := acctest.RandIntRange(10000, 99999)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudApigAiModelProvidersSourceConfig(rand, map[string]string{
			"gateway_id": `"${alicloud_apig_ai_model_provider.default.gateway_id}"`,
			"ids":        `["${alicloud_apig_ai_model_provider.default.id}"]`,
		}),
		fakeConfig: testAccCheckAliCloudApigAiModelProvidersSourceConfig(rand, map[string]string{
			"gateway_id": `"${alicloud_apig_ai_model_provider.default.gateway_id}"`,
			"ids":        `["${alicloud_apig_ai_model_provider.default.id}_fake"]`,
		}),
	}

	ApigAiModelProvidersCheckInfo.dataSourceTestCheck(t, rand, idsConf)
}

var existApigAiModelProvidersMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"ids.#":                          "1",
		"providers.#":                    "1",
		"providers.0.id":                 CHECKSET,
		"providers.0.model_provider_id":  CHECKSET,
		"providers.0.gateway_id":         CHECKSET,
		"providers.0.model_provider":     "openai",
		"providers.0.display_name":       fmt.Sprintf("tfaccapigaimp%d", rand),
		"providers.0.source":             CHECKSET,
		"providers.0.update_time":        CHECKSET,
	}
}

var fakeApigAiModelProvidersMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"ids.#":       "0",
		"providers.#": "0",
	}
}

var ApigAiModelProvidersCheckInfo = dataSourceAttr{
	resourceId:   "data.alicloud_apig_ai_model_providers.default",
	existMapFunc: existApigAiModelProvidersMapFunc,
	fakeMapFunc:  fakeApigAiModelProvidersMapFunc,
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

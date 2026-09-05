package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAliCloudApigPoliciesDataSource(t *testing.T) {
	testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
	rand := acctest.RandIntRange(10000, 99999)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudApigPoliciesSourceConfig(rand, map[string]string{
			"ids":            `["${alicloud_apig_policy.default.id}"]`,
			"enable_details": `true`,
		}),
		fakeConfig: testAccCheckAliCloudApigPoliciesSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_apig_policy.default.id}_fake"]`,
		}),
	}

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudApigPoliciesSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_apig_policy.default.policy_name}"`,
		}),
		fakeConfig: testAccCheckAliCloudApigPoliciesSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_apig_policy.default.policy_name}_fake"`,
		}),
	}

	attachConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudApigPoliciesSourceConfig(rand, map[string]string{
			"gateway_id":           `"${alicloud_apig_policy.default.gateway_id}"`,
			"environment_id":       `"${alicloud_apig_policy.default.environment_id}"`,
			"attach_resource_type": `"GatewayRoute"`,
			"attach_resource_ids":  `"${alicloud_apig_route.policy_route.route_id}"`,
			"enable_details":       `true`,
		}),
		fakeConfig: testAccCheckAliCloudApigPoliciesSourceConfig(rand, map[string]string{
			"gateway_id":           `"${alicloud_apig_policy.default.gateway_id}"`,
			"attach_resource_type": `"GatewayRoute"`,
			"attach_resource_ids":  `"${alicloud_apig_route.policy_route.route_id}_fake"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudApigPoliciesSourceConfig(rand, map[string]string{
			"ids":                  `["${alicloud_apig_policy.default.id}"]`,
			"name_regex":           `"${alicloud_apig_policy.default.policy_name}"`,
			"gateway_id":           `"${alicloud_apig_policy.default.gateway_id}"`,
			"environment_id":       `"${alicloud_apig_policy.default.environment_id}"`,
			"attach_resource_type": `"GatewayRoute"`,
			"attach_resource_ids":  `"${alicloud_apig_route.policy_route.route_id}"`,
			"enable_details":       `true`,
			"output_file":          `"/tmp/policies.json"`,
		}),
		fakeConfig: testAccCheckAliCloudApigPoliciesSourceConfig(rand, map[string]string{
			"ids":        `["${alicloud_apig_policy.default.id}_fake"]`,
			"name_regex": `"${alicloud_apig_policy.default.policy_name}_fake"`,
		}),
	}

	ApigPoliciesCheckInfo.dataSourceTestCheck(t, rand, idsConf, nameRegexConf, attachConf, allConf)
}

var existApigPoliciesMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"ids.#":                            "1",
		"policies.#":                       "1",
		"policies.0.id":                    CHECKSET,
		"policies.0.policy_id":             CHECKSET,
		"policies.0.policy_name":           fmt.Sprintf("tfaccapig%d", rand),
		"policies.0.policy_class_name":     "RateLimit",
		"policies.0.policy_config":         CHECKSET,
		"policies.0.gateway_id":            CHECKSET,
		"policies.0.environment_id":        CHECKSET,
		"policies.0.attach_resource_type":  "GatewayRoute",
		"policies.0.attach_resource_ids.#": "1",
		"policies.0.policy_attachment_id":  CHECKSET,
	}
}

var fakeApigPoliciesMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"ids.#":      "0",
		"policies.#": "0",
	}
}

var ApigPoliciesCheckInfo = dataSourceAttr{
	resourceId:   "data.alicloud_apig_policies.default",
	existMapFunc: existApigPoliciesMapFunc,
	fakeMapFunc:  fakeApigPoliciesMapFunc,
}

func testAccCheckAliCloudApigPoliciesSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	name := fmt.Sprintf("tfaccapig%d", rand)
	return AlicloudApigPolicyBasicDependence13010(name) + fmt.Sprintf(`
resource "alicloud_apig_policy" "default" {
  attach_resource_ids  = [alicloud_apig_route.policy_route.route_id]
  policy_config        = "{\"unit\":\"Second\",\"limit\":100,\"responseStatusCode\":429}"
  policy_name          = "%s"
  gateway_id           = alicloud_apig_gateway.policy_gateway.id
  attach_resource_type = "GatewayRoute"
  policy_class_name    = "RateLimit"
  environment_id       = alicloud_apig_gateway.policy_gateway.environments.0.environment_id
}

data "alicloud_apig_policies" "default" {
  %s
}
`, name, strings.Join(pairs, "\n  "))
}

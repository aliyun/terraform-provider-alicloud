package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAliCloudApigPolicyClassesDataSource(t *testing.T) {
	testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
	rand := acctest.RandIntRange(10000, 99999)

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudApigPolicyClassesSourceConfig(rand, map[string]string{}),
		fakeConfig: testAccCheckAlicloudApigPolicyClassesSourceConfig(rand, map[string]string{
			"name_regex": fmt.Sprintf(`"tfacc-no-such-policy-class-%d"`, rand),
		}),
	}

	typeConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudApigPolicyClassesSourceConfig(rand, map[string]string{
			"type": `"Auth"`,
		}),
		fakeConfig: testAccCheckAlicloudApigPolicyClassesSourceConfig(rand, map[string]string{
			"type":       `"Auth"`,
			"name_regex": fmt.Sprintf(`"tfacc-no-such-policy-class-%d"`, rand),
		}),
	}

	ApigPolicyClassesCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, func() {
		testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
	}, allConf, typeConf)
}

var existApigPolicyClassesMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"classes.#":                             CHECKSET,
		"classes.0.id":                          CHECKSET,
		"classes.0.policy_class_id":             CHECKSET,
		"classes.0.policy_class_name":           CHECKSET,
		"classes.0.alias":                       CHECKSET,
		"classes.0.version":                     CHECKSET,
		"classes.0.description":                 CHECKSET,
		"classes.0.type":                        CHECKSET,
		"classes.0.direction":                   CHECKSET,
		"classes.0.attachable_resource_types.#": CHECKSET,
		"classes.0.execute_stage":               CHECKSET,
		"classes.0.execute_priority":            CHECKSET,
		"classes.0.enable_log":                  CHECKSET,
		"classes.0.config_example":              CHECKSET,
		"classes.0.deprecated":                  CHECKSET,
		"ids.#":                                 CHECKSET,
		"names.#":                               CHECKSET,
	}
}

var fakeApigPolicyClassesMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"classes.#": "0",
		"ids.#":     "0",
		"names.#":   "0",
	}
}

var ApigPolicyClassesCheckInfo = dataSourceAttr{
	resourceId:   "data.alicloud_apig_policy_classes.default",
	existMapFunc: existApigPolicyClassesMapFunc,
	fakeMapFunc:  fakeApigPolicyClassesMapFunc,
}

func testAccCheckAlicloudApigPolicyClassesSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	return fmt.Sprintf(`
data "alicloud_apig_policy_classes" "default" {
  %s
}
`, strings.Join(pairs, "\n  "))
}

package alicloud

import (
	"fmt"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudMaxComputeProjectDataSource(t *testing.T) {
	rand := acctest.RandIntRange(1000000, 9999999)
	checkoutSupportedRegions(t, true, connectivity.MaxComputeProjectSupportRegions)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudMaxComputeProjectSourceConfig(rand, map[string]string{
			"ids":        `["${alicloud_maxcompute_project.default.id}"]`,
			"name_regex": `"${var.name}"`,
		}),
		fakeConfig: testAccCheckAlicloudMaxComputeProjectSourceConfig(rand, map[string]string{
			"ids":        `["${alicloud_maxcompute_project.default.id}"]`,
			"name_regex": `"${var.name}_fake"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudMaxComputeProjectSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_maxcompute_project.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudMaxComputeProjectSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_maxcompute_project.default.id}_fake"]`,
		}),
	}

	MaxComputeProjectCheckInfo.dataSourceTestCheck(t, rand, idsConf, allConf)
}

var existMaxComputeProjectMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"projects.#":                       "1",
		"projects.0.id":                    CHECKSET,
		"projects.0.comment":               CHECKSET,
		"projects.0.default_quota":         CHECKSET,
		"projects.0.owner":                 CHECKSET,
		"projects.0.project_name":          CHECKSET,
		"projects.0.properties.#":          "1",
		"projects.0.security_properties.#": "1",
		"projects.0.status":                CHECKSET,
		"projects.0.type":                  CHECKSET,
	}
}

var fakeMaxComputeProjectMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"projects.#": "0",
	}
}

var MaxComputeProjectCheckInfo = dataSourceAttr{
	resourceId:   "data.alicloud_maxcompute_projects.default",
	existMapFunc: existMaxComputeProjectMapFunc,
	fakeMapFunc:  fakeMaxComputeProjectMapFunc,
}

func testAccCheckAlicloudMaxComputeProjectSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
	default = "tf_testaccmp%d"
}

resource "alicloud_maxcompute_project" "default" {
  default_quota = "默认后付费Quota"
  project_name  = var.name
  comment       = var.name
  product_type  = "PayAsYouGo"
}

data "alicloud_maxcompute_projects" "default" {
%s
}
`, rand, strings.Join(pairs, "\n   "))
	return config
}

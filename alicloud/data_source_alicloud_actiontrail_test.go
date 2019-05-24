package alicloud

import (
	"fmt"
	"testing"

	"strings"

	"github.com/hashicorp/terraform/helper/acctest"
)

func TestAccAlicloudActiontrailDataSource(t *testing.T) {
	rand := acctest.RandIntRange(1000, 9999)
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudActiontrailDataSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_actiontrail.default.name}"`,
		}),
		fakeConfig: testAccCheckAlicloudActiontrailDataSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_actiontrail.default.name}_fake"`,
		}),
	}

	var existActiontrailMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"actiontrails.#":                    "1",
			"names.#":                           "1",
			"names.0":                           fmt.Sprintf("tf-testacc-actiontrail-%v", rand),
			"actiontrails.0.name":               fmt.Sprintf("tf-testacc-actiontrail-%v", rand),
			"actiontrails.0.event_rw":           "Write",
			"actiontrails.0.oss_bucket_name":    fmt.Sprintf("tf-testacc-actiontrail-%v", rand),
			"actiontrails.0.role_name":          fmt.Sprintf("tf-testacc-actiontrail-%v", rand),
			"actiontrails.0.oss_key_prefix":     "at-product-account-audit-B",
			"actiontrails.0.sls_project_arn":    "",
			"actiontrails.0.sls_write_role_arn": "",
		}
	}

	var fakeActiontrailMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"actiontrails.#": "0",
			"names.#":        "0",
		}
	}

	var actiontrailCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_actiontrails.default",
		existMapFunc: existActiontrailMapFunc,
		fakeMapFunc:  fakeActiontrailMapFunc,
	}

	actiontrailCheckInfo.dataSourceTestCheck(t, rand, nameRegexConf)

}

func testAccCheckAlicloudActiontrailDataSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`
variable "name" {
    default = "tf-testacc-actiontrail-%v"
}

resource "alicloud_ram_role" "default" {
	  name = "${var.name}"
	  services = ["actiontrail.aliyuncs.com", "oss.aliyuncs.com"]
	  description = "this is a test"
	  force = "true"
}

resource "alicloud_oss_bucket" "default" {
    bucket  = "${var.name}"
}

resource "alicloud_ram_policy" "default" {
	  name = "${var.name}"
	  statement = [
	    {
	      effect = "Allow"
	      action = ["*"]
	      resource = [
		"acs:oss:*:*:${alicloud_oss_bucket.default.id}",
		"acs:oss:*:*:${alicloud_oss_bucket.default.id}"]
	    }]
	  description = "this is a policy test"
	  force = true
	}

resource "alicloud_ram_role_policy_attachment" "default" {
	policy_name = "${alicloud_ram_policy.default.name}"
	role_name = "${alicloud_ram_role.default.name}"
	policy_type = "${alicloud_ram_policy.default.type}"
}
	
resource "alicloud_actiontrail" "default" {
	name = "${var.name}"
	event_rw = "Write"
	oss_bucket_name = "${alicloud_oss_bucket.default.id}"
	role_name = "${alicloud_ram_role_policy_attachment.default.role_name}"
	oss_key_prefix = "at-product-account-audit-B"
}

data "alicloud_actiontrails" "default"{
	  %s
}
`, rand, strings.Join(pairs, "\n  "))
	return config
}

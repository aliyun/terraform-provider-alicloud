package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
)

func TestAccAlicloudEnsSecurityGroupsDataSource(t *testing.T) {
	rand := acctest.RandIntRange(10000, 99999)

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEnsSecurityGroupsSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_ens_security_group.default.security_group_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudEnsSecurityGroupsSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_ens_security_group.default.security_group_name}_fake"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEnsSecurityGroupsSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_ens_security_group.default.security_group_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudEnsSecurityGroupsSourceConfig(rand, map[string]string{
			"name_regex": `"TestAccAlicloudEnsSecurityGroupsDataSource_fake"`,
		}),
	}

	EnsSecurityGroupCheckInfo.dataSourceTestCheck(t, rand, nameRegexConf, allConf)
}

func testAccCheckAlicloudEnsSecurityGroupsSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
    default = "tf-testAccEnsSecurityGroupsDataSource%d"
}

resource "alicloud_ens_security_group" "default" {
  security_group_name = var.name
  description         = var.name
  permissions {
    direction      = "ingress"
    ip_protocol    = "TCP"
    port_range     = "80/80"
    policy         = "Accept"
    priority       = 1
    source_cidr_ip = "0.0.0.0/0"
    dest_cidr_ip   = "0.0.0.0/0"
  }
}

data "alicloud_ens_security_groups" "default" {
  enable_details = true
%s
}
`, rand, strings.Join(pairs, "\n   "))
	return config
}

var existEnsSecurityGroupMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"groups.#":                     "1",
		"groups.0.security_group_name": fmt.Sprintf("tf-testAccEnsSecurityGroupsDataSource%d", rand),
		"groups.0.permissions.#":       "1",
	}
}

var fakeEnsSecurityGroupMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"groups.#": "0",
	}
}

var EnsSecurityGroupCheckInfo = dataSourceAttr{
	resourceId:   "data.alicloud_ens_security_groups.default",
	existMapFunc: existEnsSecurityGroupMapFunc,
	fakeMapFunc:  fakeEnsSecurityGroupMapFunc,
}

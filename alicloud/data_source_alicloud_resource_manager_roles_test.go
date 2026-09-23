package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudResourceManagerRolesDataSource(t *testing.T) {
	rand := acctest.RandIntRange(1000000, 9999999)
	resourceId := "data.alicloud_resource_manager_roles.example"

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudResourceManagerRolesSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_resource_manager_role.example.role_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudResourceManagerRolesSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_resource_manager_role.example.role_name}_fake"`,
		}),
	}

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudResourceManagerRolesSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_resource_manager_role.example.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudResourceManagerRolesSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_resource_manager_role.example.id}_fake"]`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudResourceManagerRolesSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_resource_manager_role.example.role_name}"`,
			"ids":        `["${alicloud_resource_manager_role.example.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudResourceManagerRolesSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_resource_manager_role.example.role_name}_fake"`,
			"ids":        `["${alicloud_resource_manager_role.example.id}"]`,
		}),
	}

	var existResourceManagerRolesRecordsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"roles.#":                             "1",
			"names.#":                             "1",
			"ids.#":                               "1",
			"roles.0.id":                          fmt.Sprintf("tf-testAccRole-%d", rand),
			"roles.0.role_id":                     CHECKSET,
			"roles.0.role_name":                   fmt.Sprintf("tf-testAccRole-%d", rand),
			"roles.0.arn":                         CHECKSET,
			"roles.0.assume_role_policy_document": CHECKSET,
			"roles.0.description":                 "tf_test",
			"roles.0.max_session_duration":        CHECKSET,
			"roles.0.update_date":                 CHECKSET,
		}
	}

	var fakeResourceManagerRolesRecordsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"roles.#": "0",
			"ids.#":   "0",
			"names.#": "0",
		}
	}

	var rolesRecordsCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existResourceManagerRolesRecordsMapFunc,
		fakeMapFunc:  fakeResourceManagerRolesRecordsMapFunc,
	}

	rolesRecordsCheckInfo.dataSourceTestCheck(t, rand, nameRegexConf, idsConf, allConf)

}

func testAccCheckAlicloudResourceManagerRolesSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
data "alicloud_account" "example" {}

resource "alicloud_resource_manager_role" "example"{
	  role_name = "tf-testAccRole-%d"
	  assume_role_policy_document = <<EOF
		{
			"Statement": [{
				"Action": "sts:AssumeRole",
				"Effect": "Allow",
				"Principal": {
					"RAM":["acs:ram::${data.alicloud_account.example.id}:root"]
				}
			}],
			"Version": "1"
		}
	EOF
 	description = "tf_test"
}

data "alicloud_resource_manager_roles" "example"{
	enable_details = true
%s
}
`, rand, strings.Join(pairs, "\n   "))
	return config
}

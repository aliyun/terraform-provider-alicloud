package resourcemanager_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAlicloudResourceManagerRolesDataSource(t *testing.T) {
	rand := tfacctest.RandIntRange(1000000, 9999999)
	resourceId := "data.alicloud_resource_manager_roles.example"

	nameRegexConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudResourceManagerRolesSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_resource_manager_role.example.role_name}"`,
		}),
		FakeConfig: testAccCheckAlicloudResourceManagerRolesSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_resource_manager_role.example.role_name}_fake"`,
		}),
	}

	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudResourceManagerRolesSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_resource_manager_role.example.id}"]`,
		}),
		FakeConfig: testAccCheckAlicloudResourceManagerRolesSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_resource_manager_role.example.id}_fake"]`,
		}),
	}

	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudResourceManagerRolesSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_resource_manager_role.example.role_name}"`,
			"ids":        `["${alicloud_resource_manager_role.example.id}"]`,
		}),
		FakeConfig: testAccCheckAlicloudResourceManagerRolesSourceConfig(rand, map[string]string{
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
			"roles.0.role_id":                     acctest.CHECKSET,
			"roles.0.role_name":                   fmt.Sprintf("tf-testAccRole-%d", rand),
			"roles.0.arn":                         acctest.CHECKSET,
			"roles.0.assume_role_policy_document": acctest.CHECKSET,
			"roles.0.description":                 "tf_test",
			"roles.0.max_session_duration":        acctest.CHECKSET,
			"roles.0.update_date":                 acctest.CHECKSET,
		}
	}

	var fakeResourceManagerRolesRecordsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"roles.#": "0",
			"ids.#":   "0",
			"names.#": "0",
		}
	}

	var rolesRecordsCheckInfo = acctest.DataSourceAttr{
		ResourceId:   resourceId,
		ExistMapFunc: existResourceManagerRolesRecordsMapFunc,
		FakeMapFunc:  fakeResourceManagerRolesRecordsMapFunc,
	}

	rolesRecordsCheckInfo.DataSourceTestCheck(t, rand, nameRegexConf, idsConf, allConf)

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

package alicloud

import (
	"strings"
	"testing"

	"fmt"

	"github.com/hashicorp/terraform/helper/acctest"
)

func TestAccAlicloudRamRolesDataSource(t *testing.T) {
	rand := acctest.RandIntRange(1000000, 99999999)

	policyTypeConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudRamRolesDataSourceConfig(rand, map[string]string{
			"policy_name": `"${alicloud_ram_role_policy_attachment.default.policy_name}"`,
			"policy_type": `"Custom"`,
		}),
		/*fakeConfig: testAccCheckAlicloudRamRolesDataSourceConfig(rand, map[string]string{
			"policy_name": `"${alicloud_ram_role_policy_attachment.default.policy_name}"`,
			"policy_type": `"System"`,
		}),*/
	}
	policyNameConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudRamRolesDataSourceConfig(rand, map[string]string{
			"policy_name": `"${alicloud_ram_role_policy_attachment.default.policy_name}"`,
		}),
		/*fakeConfig: testAccCheckAlicloudRamRolesDataSourceConfig(rand, map[string]string{
			"policy_name": `"${alicloud_ram_role_policy_attachment.default.policy_name}_fake"`,
		}),*/
	}

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudRamRolesDataSourceConfig(rand, map[string]string{
			"policy_name": `"${alicloud_ram_role_policy_attachment.default.policy_name}"`,
			"name_regex":  `"${alicloud_ram_role.default.name}"`,
		}),
		/*fakeConfig: testAccCheckAlicloudRamRolesDataSourceConfig(rand, map[string]string{
			"policy_name": `"${alicloud_ram_role_policy_attachment.default.policy_name}"`,
			"name_regex":  `"${alicloud_ram_role.default.name}_fake"`,
		}),*/
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudRamRolesDataSourceConfig(rand, map[string]string{
			"policy_type": "Custom",
			"policy_name": `"${alicloud_ram_role_policy_attachment.default.policy_name}"`,
			"name_regex":  `"${alicloud_ram_role.default.name}"`,
		}),
		/*fakeConfig: testAccCheckAlicloudRamRolesDataSourceConfig(rand, map[string]string{
			"policy_type": `"Custom"`,
			"policy_name": `"${alicloud_ram_role_policy_attachment.default.policy_name}_fake"`,
			"name_regex":  `"${alicloud_ram_role.default.name}_fake"`,
		}),*/
	}

	ramRolesCheckInfo.dataSourceTestCheck(t, rand, policyTypeConf, policyNameConf, nameRegexConf, allConf)

}

func testAccCheckAlicloudRamRolesDataSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`
variable "name" {
	  default = "tf-testAccRamRolesDataSourceForPolicy-%d"
	}
	resource "alicloud_ram_policy" "default" {
	  name = "${var.name}"
	  statement = [
	    {
	      effect = "Deny"
	      action = [
		"oss:ListObjects",
		"oss:ListObjects"]
	      resource = [
		"acs:oss:*:*:mybucket",
		"acs:oss:*:*:mybucket/*"]
	    }]
	  description = "this is a policy test"
	  force = true
	}

	resource "alicloud_ram_role" "default" {
	  name = "${var.name}"
	  services = ["apigateway.aliyuncs.com", "ecs.aliyuncs.com"]
	  description = "this is a test"
	  force = true
	}

	resource "alicloud_ram_role_policy_attachment" "default" {
	  policy_name = "${alicloud_ram_policy.default.name}"
	  role_name = "${alicloud_ram_role.default.name}"
	  policy_type = "${alicloud_ram_policy.default.type}"
	}

data "alicloud_ram_roles" "default" {
	%s
}`, rand, strings.Join(pairs, "\n  "))
	return config
}

var existRamRolesMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"ids.#":               "1",
		"names.#":             "1",
		"roles.#":             "1",
		"roles.0.id":          CHECKSET,
		"roles.0.name":        fmt.Sprintf("tf-testAccRamRolesDataSourceForPolicy-%d", rand),
		"roles.0.arn":         CHECKSET,
		"roles.0.description": "",
		"roles.0.document":    CHECKSET,
	}
}

var fakeRamRolesMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"ids.#":   "0",
		"names.#": "0",
		"roles.#": "0",
	}
}

var ramRolesCheckInfo = dataSourceAttr{
	resourceId:   "data.alicloud_ram_roles.default",
	existMapFunc: existRamRolesMapFunc,
	fakeMapFunc:  fakeRamRolesMapFunc,
}

/*const testAccCheckAlicloudRamRolesDataSourceForAllConfig = `
data "alicloud_ram_roles" "role" {
    policy_type="System"
}`

func TestAccAlicloudRamRolesDataSource_for_policy(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudRamRolesDataSourceForPolicyConfig(acctest.RandIntRange(1000000, 99999999)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_ram_roles.role"),
					resource.TestCheckResourceAttr("data.alicloud_ram_roles.role", "roles.#", "1"),
					resource.TestMatchResourceAttr("data.alicloud_ram_roles.role", "roles.0.name",
						regexp.MustCompile("^tf-testAccRamRolesDataSourceForPolicy-*")),
				),
			},
		},
	})
}

func TestAccAlicloudRamRolesDataSource_for_all(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudRamRolesDataSourceForAllConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_ram_roles.role"),
				),
			},
		},
	})
}

func TestAccAlicloudRamRolesDataSource_role_name_regex(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudRamRolesDataSourceRoleNameRegexConfig(acctest.RandIntRange(1000000, 99999999)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_ram_roles.role"),
					resource.TestCheckResourceAttr("data.alicloud_ram_roles.role", "roles.#", "1"),
				),
			},
		},
	})
}

func TestAccAlicloudRamRolesDataSource_empty(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudRamRolesDataSourceEmpty,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_ram_roles.role"),
					resource.TestCheckResourceAttr("data.alicloud_ram_roles.role", "roles.#", "0"),
					resource.TestCheckNoResourceAttr("data.alicloud_ram_roles.role", "roles.0.id"),
					resource.TestCheckNoResourceAttr("data.alicloud_ram_roles.role", "roles.0.name"),
					resource.TestCheckNoResourceAttr("data.alicloud_ram_roles.role", "roles.0.arn"),
					resource.TestCheckNoResourceAttr("data.alicloud_ram_roles.role", "roles.0.description"),
					resource.TestCheckNoResourceAttr("data.alicloud_ram_roles.role", "roles.0.assume_role_policy_document"),
					resource.TestCheckNoResourceAttr("data.alicloud_ram_roles.role", "roles.0.document"),
					resource.TestCheckNoResourceAttr("data.alicloud_ram_roles.role", "roles.0.create_date"),
					resource.TestCheckNoResourceAttr("data.alicloud_ram_roles.role", "roles.0.update_date"),
				),
			},
		},
	})
}

func testAccCheckAlicloudRamRolesDataSourceForPolicyConfig(rand int) string {
	return fmt.Sprintf(`
	variable "name" {
	  default = "tf-testAccRamRolesDataSourceForPolicy-%d"
	}
	resource "alicloud_ram_policy" "policy" {
	  name = "${var.name}"
	  statement = [
	    {
	      effect = "Deny"
	      action = [
		"oss:ListObjects",
		"oss:ListObjects"]
	      resource = [
		"acs:oss:*:*:mybucket",
		"acs:oss:*:*:mybucket/*"]
	    }]
	  description = "this is a policy test"
	  force = true
	}

	resource "alicloud_ram_role" "role" {
	  name = "${var.name}"
	  services = ["apigateway.aliyuncs.com", "ecs.aliyuncs.com"]
	  description = "this is a test"
	  force = true
	}

	resource "alicloud_ram_role_policy_attachment" "attach" {
	  policy_name = "${alicloud_ram_policy.policy.name}"
	  role_name = "${alicloud_ram_role.role.name}"
	  policy_type = "${alicloud_ram_policy.policy.type}"
	}

	data "alicloud_ram_roles" "role" {
	  policy_name = "${alicloud_ram_role_policy_attachment.attach.policy_name}"
	  policy_type = "Custom"
	}`, rand)
}

func testAccCheckAlicloudRamRolesDataSourceRoleNameRegexConfig(rand int) string {
	return fmt.Sprintf(`
	resource "alicloud_ram_role" "role" {
	  name = "tf-testAccRamRolesDataSourceRoleNameRegex-%d"
	  services = ["apigateway.aliyuncs.com", "ecs.aliyuncs.com"]
	  description = "this is a test"
	  force = true
	}
	data "alicloud_ram_roles" "role" {
	  name_regex = "${alicloud_ram_role.role.name}"
	}`, rand)
}

const testAccCheckAlicloudRamRolesDataSourceEmpty = `
data "alicloud_ram_roles" "role" {
	name_regex = "^tf-testacc-fake-name"
}`
*/

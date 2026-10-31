package ram_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAliCloudRamRolesDataSource_basic0(t *testing.T) {
	rand := tfacctest.RandIntRange(1000, 9999)
	resourceId := "data.alicloud_ram_roles.default"
	name := fmt.Sprintf("tfaccram%d", rand)
	testAccConfig := acctest.DataSourceTestAccConfigFunc(resourceId, name, dataSourceRamRolesConfig0)

	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_ram_role.default.role_id}"},
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_ram_role.default.role_id}_fake"},
		}),
	}

	nameRegexConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alicloud_ram_role_policy_attachment.default.role_name}",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alicloud_ram_role_policy_attachment.default.role_name}_fake",
		}),
	}

	policyNameConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"policy_name": "${alicloud_ram_role_policy_attachment.default.policy_name}",
		}),
	}

	policyTypeConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"policy_name": "${alicloud_ram_role_policy_attachment.default.policy_name}",
			"policy_type": "${alicloud_ram_role_policy_attachment.default.policy_type}",
		}),
	}

	tagsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"tags": map[string]string{
				"Created": "TF",
				"For":     "Role",
			},
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"tags": map[string]string{
				"Created": "TF",
				"For":     "Role_Fake",
			},
		}),
	}

	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"ids":         []string{"${alicloud_ram_role.default.role_id}"},
			"name_regex":  "${alicloud_ram_role_policy_attachment.default.role_name}",
			"policy_name": "${alicloud_ram_role_policy_attachment.default.policy_name}",
			"policy_type": "${alicloud_ram_role_policy_attachment.default.policy_type}",
			"tags": map[string]string{
				"Created": "TF",
				"For":     "Role",
			},
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"ids":         []string{"${alicloud_ram_role.default.role_id}_fake"},
			"name_regex":  "${alicloud_ram_role_policy_attachment.default.role_name}_fake",
			"policy_name": "${alicloud_ram_role_policy_attachment.default.policy_name}",
			"policy_type": "${alicloud_ram_role_policy_attachment.default.policy_type}",
			"tags": map[string]string{
				"Created": "TF",
				"For":     "Role_Fake",
			},
		}),
	}

	var existAliCloudRamRolesDefaultDataSourceMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                               "1",
			"names.#":                             "1",
			"roles.#":                             "1",
			"roles.0.id":                          acctest.CHECKSET,
			"roles.0.name":                        acctest.CHECKSET,
			"roles.0.assume_role_policy_document": acctest.CHECKSET,
			"roles.0.document":                    acctest.CHECKSET,
			"roles.0.description":                 acctest.CHECKSET,
			"roles.0.tags.%":                      "2",
			"roles.0.tags.Created":                "TF",
			"roles.0.tags.For":                    "Role",
			"roles.0.arn":                         acctest.CHECKSET,
			"roles.0.create_date":                 acctest.CHECKSET,
			"roles.0.update_date":                 acctest.CHECKSET,
		}
	}

	var fakeAliCloudRamRolesDefaultDataSourceMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"names.#": "0",
			"roles.#": "0",
		}
	}

	var aliCloudRamRolesCheckInfo = acctest.DataSourceAttr{
		ResourceId:   "data.alicloud_ram_roles.default",
		ExistMapFunc: existAliCloudRamRolesDefaultDataSourceMapFunc,
		FakeMapFunc:  fakeAliCloudRamRolesDefaultDataSourceMapFunc,
	}

	preCheck := func() {
		acctest.PreCheck(t)
	}

	aliCloudRamRolesCheckInfo.DataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, nameRegexConf, policyNameConf, policyTypeConf, tagsConf, allConf)
}

func TestAccAliCloudRamRolesDataSource_basic1(t *testing.T) {
	rand := tfacctest.RandIntRange(1000, 9999)
	resourceId := "data.alicloud_ram_roles.default"
	name := fmt.Sprintf("tfaccram%d", rand)
	testAccConfig := acctest.DataSourceTestAccConfigFunc(resourceId, name, dataSourceRamRolesConfig1)

	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_ram_role.default.role_id}"},
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_ram_role.default.role_id}_fake"},
		}),
	}

	nameRegexConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alicloud_ram_role_policy_attachment.default.role_name}",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alicloud_ram_role_policy_attachment.default.role_name}_fake",
		}),
	}

	policyNameAndpolicyTypeConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"policy_name": "${alicloud_ram_role_policy_attachment.default.policy_name}",
			"policy_type": "${alicloud_ram_role_policy_attachment.default.policy_type}",
		}),
	}

	tagsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"tags": map[string]string{
				"Created": "TF",
				"For":     "Role",
			},
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"tags": map[string]string{
				"Created": "TF",
				"For":     "Role_Fake",
			},
		}),
	}

	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"ids":         []string{"${alicloud_ram_role.default.role_id}"},
			"name_regex":  "${alicloud_ram_role_policy_attachment.default.role_name}",
			"policy_name": "${alicloud_ram_role_policy_attachment.default.policy_name}",
			"policy_type": "${alicloud_ram_role_policy_attachment.default.policy_type}",
			"tags": map[string]string{
				"Created": "TF",
				"For":     "Role",
			},
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"ids":         []string{"${alicloud_ram_role.default.role_id}_fake"},
			"name_regex":  "${alicloud_ram_role_policy_attachment.default.role_name}_fake",
			"policy_name": "${alicloud_ram_role_policy_attachment.default.policy_name}",
			"policy_type": "${alicloud_ram_role_policy_attachment.default.policy_type}",
			"tags": map[string]string{
				"Created": "TF",
				"For":     "Role_Fake",
			},
		}),
	}

	var existAliCloudRamRolesDefaultDataSourceMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                               "1",
			"names.#":                             "1",
			"roles.#":                             "1",
			"roles.0.id":                          acctest.CHECKSET,
			"roles.0.name":                        acctest.CHECKSET,
			"roles.0.assume_role_policy_document": acctest.CHECKSET,
			"roles.0.document":                    acctest.CHECKSET,
			"roles.0.description":                 acctest.CHECKSET,
			"roles.0.tags.%":                      "2",
			"roles.0.tags.Created":                "TF",
			"roles.0.tags.For":                    "Role",
			"roles.0.arn":                         acctest.CHECKSET,
			"roles.0.create_date":                 acctest.CHECKSET,
			"roles.0.update_date":                 acctest.CHECKSET,
		}
	}

	var fakeAliCloudRamRolesDefaultDataSourceMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"names.#": "0",
			"roles.#": "0",
		}
	}

	var aliCloudRamRolesCheckInfo = acctest.DataSourceAttr{
		ResourceId:   "data.alicloud_ram_roles.default",
		ExistMapFunc: existAliCloudRamRolesDefaultDataSourceMapFunc,
		FakeMapFunc:  fakeAliCloudRamRolesDefaultDataSourceMapFunc,
	}

	preCheck := func() {
		acctest.PreCheck(t)
	}

	aliCloudRamRolesCheckInfo.DataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, nameRegexConf, policyNameAndpolicyTypeConf, tagsConf, allConf)
}

func dataSourceRamRolesConfig0(name string) string {
	return fmt.Sprintf(`
	variable "name" {
  		default = "%s"
	}

	resource "alicloud_ram_role" "default" {
  		role_name                   = var.name
  		description                 = var.name
  		force                       = true
  		assume_role_policy_document = <<EOF
  		{
    		"Statement": [
      			{
        			"Action": "sts:AssumeRole",
        			"Effect": "Allow",
        			"Principal": {
          				"Service": [
            				"ecs.aliyuncs.com"
          				]
        			}
      			}
    		],
    		"Version": "1"
  		}
  		EOF
  		tags = {
    		Created = "TF"
    		For     = "Role"
  		}
	}

	resource "alicloud_ram_role_policy_attachment" "default" {
  		role_name   = alicloud_ram_role.default.role_name
  		policy_name = "AliyunRDSGADReadOnlyAccess"
  		policy_type = "System"
	}
`, name)
}

func dataSourceRamRolesConfig1(name string) string {
	return fmt.Sprintf(`
	variable "name" {
  		default = "%s"
	}

	resource "alicloud_ram_role" "default" {
  		role_name                   = var.name
  		description                 = var.name
  		force                       = true
  		assume_role_policy_document = <<EOF
  		{
    		"Statement": [
      			{
        			"Action": "sts:AssumeRole",
        			"Effect": "Allow",
        			"Principal": {
          				"Service": [
            				"ecs.aliyuncs.com"
          				]
        			}
      			}
    		],
    		"Version": "1"
  		}
  		EOF
  		tags = {
    		Created = "TF"
    		For     = "Role"
  		}
	}

	resource "alicloud_ram_policy" "default" {
  		policy_name     = var.name
  		description     = var.name
  		force           = true
  		policy_document = <<EOF
		{
			"Statement": [
				{
					"Effect": "Allow",
					"Action": "*",
					"Resource": "*"
				}
			],
			"Version": "1"
		}
		EOF
	}

	resource "alicloud_ram_role_policy_attachment" "default" {
  		role_name   = alicloud_ram_role.default.role_name
  		policy_name = alicloud_ram_policy.default.policy_name
  		policy_type = alicloud_ram_policy.default.type
	}
`, name)
}

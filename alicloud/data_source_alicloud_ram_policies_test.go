package alicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAliCloudRamPoliciesDataSource_basic0(t *testing.T) {
	rand := acctest.RandIntRange(1000, 9999)
	resourceId := "data.alicloud_ram_policies.default"
	name := fmt.Sprintf("tfaccram%d", rand)
	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceRamPoliciesConfig0)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_ram_role_policy_attachment.default.policy_name}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_ram_role_policy_attachment.default.policy_name}_fake"},
		}),
	}

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alicloud_ram_role_policy_attachment.default.policy_name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alicloud_ram_role_policy_attachment.default.policy_name}_fake",
		}),
	}

	typeConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":  []string{"${alicloud_ram_role_policy_attachment.default.policy_name}"},
			"type": "System",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":  []string{"${alicloud_ram_role_policy_attachment.default.policy_name}"},
			"type": "Custom",
		}),
	}

	userNameConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"user_name": "${alicloud_ram_user_policy_attachment.default.user_name}",
		}),
	}

	groupNameConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"group_name": "${alicloud_ram_group_policy_attachment.default.group_name}",
		}),
	}

	roleNameConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"role_name": "${alicloud_ram_role_policy_attachment.default.role_name}",
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":        []string{"${alicloud_ram_role_policy_attachment.default.policy_name}"},
			"name_regex": "${alicloud_ram_role_policy_attachment.default.policy_name}",
			"type":       "System",
			"user_name":  "${alicloud_ram_user_policy_attachment.default.user_name}",
			"group_name": "${alicloud_ram_group_policy_attachment.default.group_name}",
			"role_name":  "${alicloud_ram_role_policy_attachment.default.role_name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":        []string{"${alicloud_ram_role_policy_attachment.default.policy_name}_fake"},
			"name_regex": "${alicloud_ram_role_policy_attachment.default.policy_name}_fake",
			"type":       "Custom",
			"user_name":  "${alicloud_ram_user_policy_attachment.default.user_name}",
			"group_name": "${alicloud_ram_group_policy_attachment.default.group_name}",
			"role_name":  "${alicloud_ram_role_policy_attachment.default.role_name}",
		}),
	}

	var existAliCloudRamPoliciesDefaultDataSourceMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                       "1",
			"names.#":                     "1",
			"policies.#":                  "1",
			"policies.0.id":               CHECKSET,
			"policies.0.policy_name":      CHECKSET,
			"policies.0.name":             CHECKSET,
			"policies.0.type":             CHECKSET,
			"policies.0.description":      CHECKSET,
			"policies.0.tags.%":           "0",
			"policies.0.default_version":  CHECKSET,
			"policies.0.attachment_count": CHECKSET,
			"policies.0.policy_document":  CHECKSET,
			"policies.0.document":         CHECKSET,
			"policies.0.version_id":       CHECKSET,
			"policies.0.create_date":      CHECKSET,
			"policies.0.update_date":      CHECKSET,
		}
	}

	var fakeAliCloudRamPoliciesDefaultDataSourceMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":      "0",
			"names.#":    "0",
			"policies.#": "0",
		}
	}

	var aliCloudRamPoliciesCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_ram_policies.default",
		existMapFunc: existAliCloudRamPoliciesDefaultDataSourceMapFunc,
		fakeMapFunc:  fakeAliCloudRamPoliciesDefaultDataSourceMapFunc,
	}

	preCheck := func() {
		testAccPreCheck(t)
	}

	aliCloudRamPoliciesCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, nameRegexConf, typeConf, userNameConf, groupNameConf, roleNameConf, allConf)
}

func TestAccAliCloudRamPoliciesDataSource_basic1(t *testing.T) {
	rand := acctest.RandIntRange(1000, 9999)
	resourceId := "data.alicloud_ram_policies.default"
	name := fmt.Sprintf("tfaccram%d", rand)
	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceRamPoliciesConfig0)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":            []string{"${alicloud_ram_role_policy_attachment.default.policy_name}"},
			"enable_details": "true",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":            []string{"${alicloud_ram_role_policy_attachment.default.policy_name}"},
			"enable_details": "false",
		}),
	}

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex":     "${alicloud_ram_role_policy_attachment.default.policy_name}",
			"enable_details": "true",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex":     "${alicloud_ram_role_policy_attachment.default.policy_name}",
			"enable_details": "false",
		}),
	}

	typeConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":            []string{"${alicloud_ram_role_policy_attachment.default.policy_name}"},
			"type":           "System",
			"enable_details": "true",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":            []string{"${alicloud_ram_role_policy_attachment.default.policy_name}"},
			"type":           "System",
			"enable_details": "false",
		}),
	}

	userNameConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"user_name":      "${alicloud_ram_user_policy_attachment.default.user_name}",
			"enable_details": "true",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"user_name":      "${alicloud_ram_user_policy_attachment.default.user_name}",
			"enable_details": "false",
		}),
	}

	groupNameConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"group_name":     "${alicloud_ram_group_policy_attachment.default.group_name}",
			"enable_details": "true",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"group_name":     "${alicloud_ram_group_policy_attachment.default.group_name}",
			"enable_details": "false",
		}),
	}

	roleNameConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"role_name":      "${alicloud_ram_role_policy_attachment.default.role_name}",
			"enable_details": "true",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"role_name":      "${alicloud_ram_role_policy_attachment.default.role_name}",
			"enable_details": "false",
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":            []string{"${alicloud_ram_role_policy_attachment.default.policy_name}"},
			"name_regex":     "${alicloud_ram_role_policy_attachment.default.policy_name}",
			"type":           "System",
			"user_name":      "${alicloud_ram_user_policy_attachment.default.user_name}",
			"group_name":     "${alicloud_ram_group_policy_attachment.default.group_name}",
			"role_name":      "${alicloud_ram_role_policy_attachment.default.role_name}",
			"enable_details": "true",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":            []string{"${alicloud_ram_role_policy_attachment.default.policy_name}"},
			"name_regex":     "${alicloud_ram_role_policy_attachment.default.policy_name}",
			"type":           "System",
			"user_name":      "${alicloud_ram_user_policy_attachment.default.user_name}",
			"group_name":     "${alicloud_ram_group_policy_attachment.default.group_name}",
			"role_name":      "${alicloud_ram_role_policy_attachment.default.role_name}",
			"enable_details": "false",
		}),
	}

	var existAliCloudRamPoliciesDefaultDataSourceMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                       "1",
			"names.#":                     "1",
			"policies.#":                  "1",
			"policies.0.id":               CHECKSET,
			"policies.0.policy_name":      CHECKSET,
			"policies.0.name":             CHECKSET,
			"policies.0.type":             CHECKSET,
			"policies.0.description":      CHECKSET,
			"policies.0.tags.%":           "0",
			"policies.0.default_version":  CHECKSET,
			"policies.0.attachment_count": CHECKSET,
			"policies.0.policy_document":  CHECKSET,
			"policies.0.document":         CHECKSET,
			"policies.0.version_id":       CHECKSET,
			"policies.0.create_date":      CHECKSET,
			"policies.0.update_date":      CHECKSET,
		}
	}

	var fakeAliCloudRamPoliciesDefaultDataSourceMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                       "1",
			"names.#":                     "1",
			"policies.#":                  "1",
			"policies.0.id":               CHECKSET,
			"policies.0.policy_name":      CHECKSET,
			"policies.0.name":             CHECKSET,
			"policies.0.type":             CHECKSET,
			"policies.0.description":      CHECKSET,
			"policies.0.tags.%":           "0",
			"policies.0.default_version":  CHECKSET,
			"policies.0.attachment_count": CHECKSET,
			"policies.0.policy_document":  "",
			"policies.0.document":         "",
			"policies.0.version_id":       "",
			"policies.0.create_date":      CHECKSET,
			"policies.0.update_date":      CHECKSET,
		}
	}

	var aliCloudRamPoliciesCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_ram_policies.default",
		existMapFunc: existAliCloudRamPoliciesDefaultDataSourceMapFunc,
		fakeMapFunc:  fakeAliCloudRamPoliciesDefaultDataSourceMapFunc,
	}

	preCheck := func() {
		testAccPreCheck(t)
	}

	aliCloudRamPoliciesCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, nameRegexConf, typeConf, userNameConf, groupNameConf, roleNameConf, allConf)
}

func TestAccAliCloudRamPoliciesDataSource_basic2(t *testing.T) {
	rand := acctest.RandIntRange(1000, 9999)
	resourceId := "data.alicloud_ram_policies.default"
	name := fmt.Sprintf("tfaccram%d", rand)
	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceRamPoliciesConfig1)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_ram_role_policy_attachment.default.policy_name}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_ram_role_policy_attachment.default.policy_name}_fake"},
		}),
	}

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alicloud_ram_role_policy_attachment.default.policy_name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alicloud_ram_role_policy_attachment.default.policy_name}_fake",
		}),
	}

	typeConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":  []string{"${alicloud_ram_role_policy_attachment.default.policy_name}"},
			"type": "Custom",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":  []string{"${alicloud_ram_role_policy_attachment.default.policy_name}"},
			"type": "System",
		}),
	}

	userNameConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"user_name": "${alicloud_ram_user_policy_attachment.default.user_name}",
		}),
	}

	groupNameConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"group_name": "${alicloud_ram_group_policy_attachment.default.group_name}",
		}),
	}

	roleNameConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"role_name": "${alicloud_ram_role_policy_attachment.default.role_name}",
		}),
	}

	tagsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"tags": map[string]string{
				"Created": "TF",
				"For":     "Policy",
			},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"tags": map[string]string{
				"Created": "TF",
				"For":     "Policy_Fake",
			},
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":        []string{"${alicloud_ram_role_policy_attachment.default.policy_name}"},
			"name_regex": "${alicloud_ram_role_policy_attachment.default.policy_name}",
			"type":       "Custom",
			"user_name":  "${alicloud_ram_user_policy_attachment.default.user_name}",
			"group_name": "${alicloud_ram_group_policy_attachment.default.group_name}",
			"role_name":  "${alicloud_ram_role_policy_attachment.default.role_name}",
			"tags": map[string]string{
				"Created": "TF",
				"For":     "Policy",
			},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":        []string{"${alicloud_ram_role_policy_attachment.default.policy_name}_fake"},
			"name_regex": "${alicloud_ram_role_policy_attachment.default.policy_name}_fake",
			"type":       "Custom",
			"user_name":  "${alicloud_ram_user_policy_attachment.default.user_name}",
			"group_name": "${alicloud_ram_group_policy_attachment.default.group_name}",
			"role_name":  "${alicloud_ram_role_policy_attachment.default.role_name}",
			"tags": map[string]string{
				"Created": "TF",
				"For":     "Policy_Fake",
			},
		}),
	}

	var existAliCloudRamPoliciesDefaultDataSourceMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                       "1",
			"names.#":                     "1",
			"policies.#":                  "1",
			"policies.0.id":               CHECKSET,
			"policies.0.policy_name":      CHECKSET,
			"policies.0.name":             CHECKSET,
			"policies.0.type":             CHECKSET,
			"policies.0.description":      CHECKSET,
			"policies.0.tags.%":           "2",
			"policies.0.tags.Created":     "TF",
			"policies.0.tags.For":         "Policy",
			"policies.0.default_version":  CHECKSET,
			"policies.0.attachment_count": CHECKSET,
			"policies.0.policy_document":  CHECKSET,
			"policies.0.document":         CHECKSET,
			"policies.0.version_id":       CHECKSET,
			"policies.0.create_date":      CHECKSET,
			"policies.0.update_date":      CHECKSET,
		}
	}

	var fakeAliCloudRamPoliciesDefaultDataSourceMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":      "0",
			"names.#":    "0",
			"policies.#": "0",
		}
	}

	var aliCloudRamPoliciesCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_ram_policies.default",
		existMapFunc: existAliCloudRamPoliciesDefaultDataSourceMapFunc,
		fakeMapFunc:  fakeAliCloudRamPoliciesDefaultDataSourceMapFunc,
	}

	preCheck := func() {
		testAccPreCheck(t)
	}

	aliCloudRamPoliciesCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, nameRegexConf, typeConf, userNameConf, groupNameConf, roleNameConf, tagsConf, allConf)
}

func TestAccAliCloudRamPoliciesDataSource_basic3(t *testing.T) {
	rand := acctest.RandIntRange(1000, 9999)
	resourceId := "data.alicloud_ram_policies.default"
	name := fmt.Sprintf("tfaccram%d", rand)
	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceRamPoliciesConfig1)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":            []string{"${alicloud_ram_role_policy_attachment.default.policy_name}"},
			"enable_details": "true",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":            []string{"${alicloud_ram_role_policy_attachment.default.policy_name}"},
			"enable_details": "false",
		}),
	}

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex":     "${alicloud_ram_role_policy_attachment.default.policy_name}",
			"enable_details": "true",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex":     "${alicloud_ram_role_policy_attachment.default.policy_name}",
			"enable_details": "false",
		}),
	}

	typeConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":            []string{"${alicloud_ram_role_policy_attachment.default.policy_name}"},
			"type":           "Custom",
			"enable_details": "true",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":            []string{"${alicloud_ram_role_policy_attachment.default.policy_name}"},
			"type":           "Custom",
			"enable_details": "false",
		}),
	}

	userNameConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"user_name":      "${alicloud_ram_user_policy_attachment.default.user_name}",
			"enable_details": "true",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"user_name":      "${alicloud_ram_user_policy_attachment.default.user_name}",
			"enable_details": "false",
		}),
	}

	groupNameConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"group_name":     "${alicloud_ram_group_policy_attachment.default.group_name}",
			"enable_details": "true",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"group_name":     "${alicloud_ram_group_policy_attachment.default.group_name}",
			"enable_details": "false",
		}),
	}

	roleNameConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"role_name":      "${alicloud_ram_role_policy_attachment.default.role_name}",
			"enable_details": "true",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"role_name":      "${alicloud_ram_role_policy_attachment.default.role_name}",
			"enable_details": "false",
		}),
	}

	tagsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"tags": map[string]string{
				"Created": "TF",
				"For":     "Policy",
			},
			"enable_details": "true",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"tags": map[string]string{
				"Created": "TF",
				"For":     "Policy",
			},
			"enable_details": "false",
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":        []string{"${alicloud_ram_role_policy_attachment.default.policy_name}"},
			"name_regex": "${alicloud_ram_role_policy_attachment.default.policy_name}",
			"type":       "Custom",
			"user_name":  "${alicloud_ram_user_policy_attachment.default.user_name}",
			"group_name": "${alicloud_ram_group_policy_attachment.default.group_name}",
			"role_name":  "${alicloud_ram_role_policy_attachment.default.role_name}",
			"tags": map[string]string{
				"Created": "TF",
				"For":     "Policy",
			},
			"enable_details": "true",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":        []string{"${alicloud_ram_role_policy_attachment.default.policy_name}"},
			"name_regex": "${alicloud_ram_role_policy_attachment.default.policy_name}",
			"type":       "Custom",
			"user_name":  "${alicloud_ram_user_policy_attachment.default.user_name}",
			"group_name": "${alicloud_ram_group_policy_attachment.default.group_name}",
			"role_name":  "${alicloud_ram_role_policy_attachment.default.role_name}",
			"tags": map[string]string{
				"Created": "TF",
				"For":     "Policy",
			},
			"enable_details": "false",
		}),
	}

	var existAliCloudRamPoliciesDefaultDataSourceMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                       "1",
			"names.#":                     "1",
			"policies.#":                  "1",
			"policies.0.id":               CHECKSET,
			"policies.0.policy_name":      CHECKSET,
			"policies.0.name":             CHECKSET,
			"policies.0.type":             CHECKSET,
			"policies.0.description":      CHECKSET,
			"policies.0.tags.%":           "2",
			"policies.0.tags.Created":     "TF",
			"policies.0.tags.For":         "Policy",
			"policies.0.default_version":  CHECKSET,
			"policies.0.attachment_count": CHECKSET,
			"policies.0.policy_document":  CHECKSET,
			"policies.0.document":         CHECKSET,
			"policies.0.version_id":       CHECKSET,
			"policies.0.create_date":      CHECKSET,
			"policies.0.update_date":      CHECKSET,
		}
	}

	var fakeAliCloudRamPoliciesDefaultDataSourceMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                       "1",
			"names.#":                     "1",
			"policies.#":                  "1",
			"policies.0.id":               CHECKSET,
			"policies.0.policy_name":      CHECKSET,
			"policies.0.name":             CHECKSET,
			"policies.0.type":             CHECKSET,
			"policies.0.description":      CHECKSET,
			"policies.0.tags.%":           "2",
			"policies.0.tags.Created":     "TF",
			"policies.0.tags.For":         "Policy",
			"policies.0.default_version":  CHECKSET,
			"policies.0.attachment_count": CHECKSET,
			"policies.0.policy_document":  "",
			"policies.0.document":         "",
			"policies.0.version_id":       "",
			"policies.0.create_date":      CHECKSET,
			"policies.0.update_date":      CHECKSET,
		}
	}

	var aliCloudRamPoliciesCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_ram_policies.default",
		existMapFunc: existAliCloudRamPoliciesDefaultDataSourceMapFunc,
		fakeMapFunc:  fakeAliCloudRamPoliciesDefaultDataSourceMapFunc,
	}

	preCheck := func() {
		testAccPreCheck(t)
	}

	aliCloudRamPoliciesCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, nameRegexConf, typeConf, userNameConf, groupNameConf, roleNameConf, tagsConf, allConf)
}

func dataSourceRamPoliciesConfig0(name string) string {
	return fmt.Sprintf(`
	variable "name" {
  		default = "%s"
	}

	resource "alicloud_ram_user" "default" {
  		name = "${var.name}-user"
	}

	resource "alicloud_ram_group" "default" {
  		group_name = "${var.name}-group"
	}

	resource "alicloud_ram_role" "default" {
  		role_name                   = "${var.name}-role"
  		description                 = "${var.name}-role"
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

	resource "alicloud_ram_user_policy_attachment" "default" {
  		user_name   = alicloud_ram_user.default.name
  		policy_name = "AliyunGlobalAccelerationReadOnlyAccess"
  		policy_type = "System"
	}

	resource "alicloud_ram_group_policy_attachment" "default" {
  		group_name  = alicloud_ram_group.default.group_name
  		policy_name = alicloud_ram_user_policy_attachment.default.policy_name
  		policy_type = "System"
	}

	resource "alicloud_ram_role_policy_attachment" "default" {
  		role_name   = alicloud_ram_role.default.role_name
  		policy_name = alicloud_ram_group_policy_attachment.default.policy_name
  		policy_type = "System"
	}
`, name)
}

func dataSourceRamPoliciesConfig1(name string) string {
	return fmt.Sprintf(`
	variable "name" {
  		default = "%s"
	}

	resource "alicloud_ram_user" "default" {
  		name = "${var.name}-user"
	}

	resource "alicloud_ram_group" "default" {
  		group_name = "${var.name}-group"
	}

	resource "alicloud_ram_role" "default" {
  		role_name                   = "${var.name}-role"
  		description                 = "${var.name}-role"
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
  		policy_name     = "${var.name}-policy"
  		description     = "${var.name}-policy"
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
  		tags = {
    		Created = "TF"
    		For     = "Policy"
  		}
	}

	resource "alicloud_ram_user_policy_attachment" "default" {
  		user_name   = alicloud_ram_user.default.name
  		policy_name = alicloud_ram_policy.default.policy_name
  		policy_type = alicloud_ram_policy.default.type
	}

	resource "alicloud_ram_group_policy_attachment" "default" {
  		group_name  = alicloud_ram_group.default.group_name
  		policy_name = alicloud_ram_user_policy_attachment.default.policy_name
  		policy_type = alicloud_ram_user_policy_attachment.default.policy_type
	}

	resource "alicloud_ram_role_policy_attachment" "default" {
  		role_name   = alicloud_ram_role.default.role_name
  		policy_name = alicloud_ram_group_policy_attachment.default.policy_name
  		policy_type = alicloud_ram_group_policy_attachment.default.policy_type
	}
`, name)
}

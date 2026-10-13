package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAliCloudIacServiceModuleDataSource(t *testing.T) {
	rand := acctest.RandIntRange(1000000, 9999999)
	projectId, groupId, groupIdEmpty, emptyProjectId := prepareIacServiceModuleGroupFixtures(t)
	name := fmt.Sprintf("tf-testacc%siacmodule%d", defaultRegionToTest, rand)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudIacServiceModuleSourceConfig(rand, groupId, projectId, map[string]string{
			"ids": `["${alicloud_iac_service_module.default.id}"]`,
		}),
		fakeConfig: testAccCheckAliCloudIacServiceModuleSourceConfig(rand, groupId, projectId, map[string]string{
			"ids": `["${alicloud_iac_service_module.default.id}_fake"]`,
		}),
	}

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudIacServiceModuleSourceConfig(rand, groupId, projectId, map[string]string{
			"name_regex": fmt.Sprintf("\"^%s$\"", name),
		}),
		fakeConfig: testAccCheckAliCloudIacServiceModuleSourceConfig(rand, groupId, projectId, map[string]string{
			"name_regex": fmt.Sprintf("\"^%s_fake$\"", name),
		}),
	}

	moduleNameConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudIacServiceModuleSourceConfig(rand, groupId, projectId, map[string]string{
			"module_name": `"${alicloud_iac_service_module.default.module_name}"`,
		}),
		fakeConfig: testAccCheckAliCloudIacServiceModuleSourceConfig(rand, groupId, projectId, map[string]string{
			"module_name": `"${alicloud_iac_service_module.default.module_name}_fake"`,
		}),
	}

	groupIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudIacServiceModuleSourceConfig(rand, groupId, projectId, map[string]string{
			"group_id": fmt.Sprintf("%q", groupId),
		}),
		fakeConfig: testAccCheckAliCloudIacServiceModuleSourceConfig(rand, groupId, projectId, map[string]string{
			"group_id": fmt.Sprintf("%q", groupIdEmpty),
		}),
	}

	projectIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudIacServiceModuleSourceConfig(rand, groupId, projectId, map[string]string{
			"project_id": fmt.Sprintf("%q", projectId),
		}),
		fakeConfig: testAccCheckAliCloudIacServiceModuleSourceConfig(rand, groupId, projectId, map[string]string{
			"project_id": fmt.Sprintf("%q", emptyProjectId),
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudIacServiceModuleSourceConfig(rand, groupId, projectId, map[string]string{
			"ids":         `["${alicloud_iac_service_module.default.id}"]`,
			"name_regex":  fmt.Sprintf("\"^%s$\"", name),
			"module_name": `"${alicloud_iac_service_module.default.module_name}"`,
			"group_id":    fmt.Sprintf("%q", groupId),
			"project_id":  fmt.Sprintf("%q", projectId),
		}),
		fakeConfig: testAccCheckAliCloudIacServiceModuleSourceConfig(rand, groupId, projectId, map[string]string{
			"ids":         `["${alicloud_iac_service_module.default.id}_fake"]`,
			"name_regex":  fmt.Sprintf("\"^%s_fake$\"", name),
			"module_name": `"${alicloud_iac_service_module.default.module_name}_fake"`,
			"group_id":    fmt.Sprintf("%q", groupIdEmpty),
			"project_id":  fmt.Sprintf("%q", emptyProjectId),
		}),
	}

	IacServiceModuleCheckInfo.dataSourceTestCheck(t, rand, idsConf, nameRegexConf, moduleNameConf, groupIdConf, projectIdConf, allConf)
}

var existIacServiceModuleMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"modules.#":                  "1",
		"modules.0.module_id":        CHECKSET,
		"modules.0.module_name":      CHECKSET,
		"modules.0.description":      CHECKSET,
		"modules.0.source":           CHECKSET,
		"modules.0.status":           CHECKSET,
		"modules.0.create_time":      CHECKSET,
		"modules.0.group_info.#":     "1",
		"modules.0.tags.#":           "1",
		"modules.0.tags.0.tag_key":   CHECKSET,
		"modules.0.tags.0.tag_value": CHECKSET,
	}
}

var fakeIacServiceModuleMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"modules.#": "0",
	}
}

var IacServiceModuleCheckInfo = dataSourceAttr{
	resourceId:   "data.alicloud_iac_service_modules.default",
	existMapFunc: existIacServiceModuleMapFunc,
	fakeMapFunc:  fakeIacServiceModuleMapFunc,
}

func testAccCheckAliCloudIacServiceModuleSourceConfig(rand int, groupId, projectId string, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	name := fmt.Sprintf("tf-testacc%siacmodule%d", defaultRegionToTest, rand)
	config := fmt.Sprintf(`
variable "name" {
	default = "%s"
}

resource "alicloud_iac_service_module" "default" {
  module_name      = var.name
  source           = "Registry"
  source_path      = "alibaba/security-group:2.4.1"
  version_strategy = "Manual"
  description      = var.name
  group_info {
    group_id   = "%s"
    project_id = "%s"
  }
  tags {
    tag_key   = "Created"
    tag_value = "TF"
  }
}

data "alicloud_iac_service_modules" "default" {
%s
}
`, name, groupId, projectId, strings.Join(pairs, "\n   "))
	return config
}

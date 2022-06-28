package alicloud

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudECSSecurityGroupsDataSourceBasic(t *testing.T) {
	rand := acctest.RandIntRange(1000, 9999)

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudSecurityGroupsDataSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_security_group.default.name}"`,
		}),
		fakeConfig: testAccCheckAlicloudSecurityGroupsDataSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_security_group.default.name}_fake"`,
		}),
	}

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudSecurityGroupsDataSourceConfig(rand, map[string]string{
			"ids": `[ "${alicloud_security_group.default.id}" ]`,
		}),
		fakeConfig: testAccCheckAlicloudSecurityGroupsDataSourceConfig(rand, map[string]string{
			"ids": `[ "${alicloud_security_group.default.id}_fake" ]`,
		}),
	}

	vpcIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudSecurityGroupsDataSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_security_group.default.name}"`,
			"vpc_id":     `"${alicloud_security_group.default.vpc_id}"`,
		}),
		fakeConfig: testAccCheckAlicloudSecurityGroupsDataSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_security_group.default.name}"`,
			"vpc_id":     `"${alicloud_security_group.default.vpc_id}_fake"`,
		}),
	}

	tagsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudSecurityGroupsDataSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_security_group.default.name}"`,
			"tags": `{
                         from = "datasource"
                         usage1 = "test"
                         usage2 = "test"
                        }`,
		}),
		fakeConfig: testAccCheckAlicloudSecurityGroupsDataSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_security_group.default.name}"`,
			"tags": `{
                         from = "datasource"
                         usage1 = "test"
                         usage2 = "test_fake"
                        }`,
		}),
	}

	resourceGroupIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudSecurityGroupsDataSourceConfig(rand, map[string]string{
			"name_regex":        `"${alicloud_security_group.default.name}"`,
			"resource_group_id": fmt.Sprintf(`"%s"`, os.Getenv("ALICLOUD_RESOURCE_GROUP_ID")),
		}),
		fakeConfig: testAccCheckAlicloudSecurityGroupsDataSourceConfig(rand, map[string]string{
			"name_regex":        `"${alicloud_security_group.default.name}"`,
			"resource_group_id": fmt.Sprintf(`"%s_fake"`, os.Getenv("ALICLOUD_RESOURCE_GROUP_ID")),
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudSecurityGroupsDataSourceConfig(rand, map[string]string{
			"name_regex":        `"${alicloud_security_group.default.name}"`,
			"ids":               `[ "${alicloud_security_group.default.id}" ]`,
			"vpc_id":            `"${alicloud_security_group.default.vpc_id}"`,
			"resource_group_id": fmt.Sprintf(`"%s"`, os.Getenv("ALICLOUD_RESOURCE_GROUP_ID")),
			"tags": `{
                         from = "datasource"
                         usage1 = "test"
                         usage2 = "test"
                        }`,
		}),
		fakeConfig: testAccCheckAlicloudSecurityGroupsDataSourceConfig(rand, map[string]string{
			"name_regex":        `"${alicloud_security_group.default.name}_fake"`,
			"ids":               `[ "${alicloud_security_group.default.id}" ]`,
			"vpc_id":            `"${alicloud_security_group.default.vpc_id}"`,
			"resource_group_id": fmt.Sprintf(`"%s"`, os.Getenv("ALICLOUD_RESOURCE_GROUP_ID")),
			"tags": `{
                         from = "datasource"
                         usage1 = "test"
                         usage2 = "test"
                        }`,
		}),
	}

	securityGroupsCheckInfo.dataSourceTestCheck(t, rand, nameRegexConf, idsConf, vpcIdConf, tagsConf, resourceGroupIdConf, allConf)
}

func testAccCheckAlicloudSecurityGroupsDataSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`
variable "name" {
	default = "tf-testAccCheckAlicloudSecurityGroupsDataSourceConfig%d"
}

data "alicloud_vpcs" "default" {
	name_regex = "default-NODELETING"
}

resource "alicloud_security_group" "default" {
  name        = "${var.name}"
  description = "test security group"
  vpc_id      = data.alicloud_vpcs.default.ids.0
  resource_group_id = "%s"
  tags = {
		from = "datasource"
		usage1 = "test"
		usage2 = "test"
  }
}

data "alicloud_security_groups" "default" {
  %s
}`, rand, os.Getenv("ALICLOUD_RESOURCE_GROUP_ID"), strings.Join(pairs, "\n  "))
	return config
}

var existSecurityGroupsMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"ids.#":                        "1",
		"names.#":                      "1",
		"groups.#":                     "1",
		"total_count":                  CHECKSET,
		"groups.0.vpc_id":              CHECKSET,
		"groups.0.resource_group_id":   os.Getenv("ALICLOUD_RESOURCE_GROUP_ID"),
		"groups.0.security_group_type": "normal",
		"groups.0.name":                fmt.Sprintf("tf-testAccCheckAlicloudSecurityGroupsDataSourceConfig%d", rand),
		"groups.0.tags.from":           "datasource",
		"groups.0.tags.usage1":         "test",
		"groups.0.tags.usage2":         "test",
		"groups.0.inner_access":        "true",
		"groups.0.creation_time":       CHECKSET,
		"groups.0.description":         "test security group",
		"groups.0.id":                  CHECKSET,
	}
}

var fakeSecurityGroupsMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"ids.#":    "0",
		"names.#":  "0",
		"groups.#": "0",
	}
}

var securityGroupsCheckInfo = dataSourceAttr{
	resourceId:   "data.alicloud_security_groups.default",
	existMapFunc: existSecurityGroupsMapFunc,
	fakeMapFunc:  fakeSecurityGroupsMapFunc,
}

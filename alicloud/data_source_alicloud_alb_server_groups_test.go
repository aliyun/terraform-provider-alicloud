package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudAlbServerGroupsDataSource(t *testing.T) {
	rand := acctest.RandInt()

	serverGroupIdsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudAlbServerGroupsDataSourceName(rand, map[string]string{
			"server_group_ids": `["${alicloud_alb_server_group.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudAlbServerGroupsDataSourceName(rand, map[string]string{
			"server_group_ids": `["${alicloud_alb_server_group.default.id}_fake"]`,
		}),
	}
	resourceGroupIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudAlbServerGroupsDataSourceName(rand, map[string]string{
			"ids":               `["${alicloud_alb_server_group.default.id}"]`,
			"resource_group_id": `"${alicloud_alb_server_group.default.resource_group_id}"`,
		}),
		fakeConfig: testAccCheckAlicloudAlbServerGroupsDataSourceName(rand, map[string]string{
			"ids":               `["${alicloud_alb_server_group.default.id}_fake"]`,
			"resource_group_id": `"${alicloud_alb_server_group.default.resource_group_id}_fake"`,
		}),
	}
	tagsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudAlbServerGroupsDataSourceName(rand, map[string]string{
			"ids":  `["${alicloud_alb_server_group.default.id}"]`,
			"tags": `{Created = "TF"}`,
		}),
		fakeConfig: testAccCheckAlicloudAlbServerGroupsDataSourceName(rand, map[string]string{
			"ids":  `["${alicloud_alb_server_group.default.id}_fake"]`,
			"tags": `{Created = "TF_fake"}`,
		}),
	}
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudAlbServerGroupsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_alb_server_group.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudAlbServerGroupsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_alb_server_group.default.id}_fake"]`,
		}),
	}
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudAlbServerGroupsDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_alb_server_group.default.server_group_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudAlbServerGroupsDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_alb_server_group.default.server_group_name}_fake"`,
		}),
	}
	vpcIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudAlbServerGroupsDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_alb_server_group.default.id}"]`,
			"vpc_id": `"${alicloud_alb_server_group.default.vpc_id}"`,
		}),
		fakeConfig: testAccCheckAlicloudAlbServerGroupsDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_alb_server_group.default.id}_fake"]`,
			"vpc_id": `"${alicloud_alb_server_group.default.vpc_id}"`,
		}),
	}
	serverGroupNameConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudAlbServerGroupsDataSourceName(rand, map[string]string{
			"server_group_name": `"${alicloud_alb_server_group.default.server_group_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudAlbServerGroupsDataSourceName(rand, map[string]string{
			"server_group_name": `"${alicloud_alb_server_group.default.server_group_name}_fake"`,
		}),
	}
	statusConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudAlbServerGroupsDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_alb_server_group.default.id}"]`,
			"status": `"Available"`,
		}),
		fakeConfig: testAccCheckAlicloudAlbServerGroupsDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_alb_server_group.default.id}"]`,
			"status": `"Configuring"`,
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudAlbServerGroupsDataSourceName(rand, map[string]string{
			"ids":               `["${alicloud_alb_server_group.default.id}"]`,
			"name_regex":        `"${alicloud_alb_server_group.default.server_group_name}"`,
			"vpc_id":            `"${alicloud_alb_server_group.default.vpc_id}"`,
			"server_group_name": `"${alicloud_alb_server_group.default.server_group_name}"`,
			"status":            `"Available"`,
			"tags":              `{Created = "TF"}`,
			"resource_group_id": `"${alicloud_alb_server_group.default.resource_group_id}"`,
			"server_group_ids":  `["${alicloud_alb_server_group.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudAlbServerGroupsDataSourceName(rand, map[string]string{
			"ids":               `["${alicloud_alb_server_group.default.id}_fake"]`,
			"name_regex":        `"${alicloud_alb_server_group.default.server_group_name}_fake"`,
			"vpc_id":            `"${alicloud_alb_server_group.default.vpc_id}"`,
			"server_group_name": `"${alicloud_alb_server_group.default.server_group_name}_fake"`,
			"status":            `"Configuring"`,
			"tags":              `{Created = "TF_fake"}`,
			"resource_group_id": `"${alicloud_alb_server_group.default.resource_group_id}_fake"`,
			"server_group_ids":  `["${alicloud_alb_server_group.default.id}_fake"]`,
		}),
	}
	var existDataAlicloudAlbServerGroupsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                      "1",
			"groups.#":                   "1",
			"groups.0.server_group_name": fmt.Sprintf("tf-testAccAlbServerGroup-%d", rand),
		}
	}
	var fakeDataAlicloudAlbServerGroupsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":    "0",
			"groups.#": "0",
		}
	}
	var alicloudAlbServerGroupCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_alb_server_groups.default",
		existMapFunc: existDataAlicloudAlbServerGroupsDataSourceNameMapFunc,
		fakeMapFunc:  fakeDataAlicloudAlbServerGroupsDataSourceNameMapFunc,
	}
	preCheck := func() {
		testAccPreCheckWithRegions(t, true, connectivity.AlbSupportRegions)
	}
	alicloudAlbServerGroupCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, nameRegexConf, vpcIdConf, statusConf, serverGroupNameConf, tagsConf, serverGroupIdsConf, resourceGroupIdConf, allConf)
}
func testAccCheckAlicloudAlbServerGroupsDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`

variable "name" {	
	default = "tf-testAccAlbServerGroup-%d"
}

data "alicloud_vpcs" "default" {
  name_regex = "default-NODELETING"
}

data "alicloud_resource_manager_resource_groups" "default" {}

resource "alicloud_alb_server_group" "default" {
	protocol = "HTTP"
	vpc_id = data.alicloud_vpcs.default.vpcs.0.id
	server_group_name = var.name
    resource_group_id = data.alicloud_resource_manager_resource_groups.default.groups.0.id
	health_check_config {
       health_check_enabled = "false"
	}
	sticky_session_config {
       sticky_session_enabled = "false"
	}
	tags = {
		Created = "TF"
	}
}

data "alicloud_alb_server_groups" "default" {	
	%s
}
`, rand, strings.Join(pairs, " \n "))
	return config
}

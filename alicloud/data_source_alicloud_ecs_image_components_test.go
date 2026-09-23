package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudECSImageComponentsDataSource(t *testing.T) {
	rand := acctest.RandInt()
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEcsImageComponentsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_ecs_image_component.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudEcsImageComponentsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_ecs_image_component.default.id}_fake"]`,
		}),
	}
	imageComponentNameConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEcsImageComponentsDataSourceName(rand, map[string]string{
			"ids":                  `["${alicloud_ecs_image_component.default.id}"]`,
			"image_component_name": `"${alicloud_ecs_image_component.default.image_component_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudEcsImageComponentsDataSourceName(rand, map[string]string{
			"ids":                  `["${alicloud_ecs_image_component.default.id}"]`,
			"image_component_name": `"${alicloud_ecs_image_component.default.image_component_name}_fake"`,
		}),
	}
	ownerConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEcsImageComponentsDataSourceName(rand, map[string]string{
			"ids":   `["${alicloud_ecs_image_component.default.id}"]`,
			"owner": `"SELF"`,
		}),
		fakeConfig: testAccCheckAlicloudEcsImageComponentsDataSourceName(rand, map[string]string{
			"ids":   `["${alicloud_ecs_image_component.default.id}"]`,
			"owner": `"ALIYUN"`,
		}),
	}
	resourceGroupIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEcsImageComponentsDataSourceName(rand, map[string]string{
			"ids":               `["${alicloud_ecs_image_component.default.id}"]`,
			"resource_group_id": `"${alicloud_ecs_image_component.default.resource_group_id}"`,
		}),
		fakeConfig: testAccCheckAlicloudEcsImageComponentsDataSourceName(rand, map[string]string{
			"ids":               `["${alicloud_ecs_image_component.default.id}"]`,
			"resource_group_id": `"${alicloud_ecs_image_component.default.resource_group_id}_fake"`,
		}),
	}
	tagsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEcsImageComponentsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_ecs_image_component.default.id}"]`,
			"tags": `{
				Created = "TF"
		}`,
		}),
		fakeConfig: testAccCheckAlicloudEcsImageComponentsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_ecs_image_component.default.id}"]`,
			"tags": `{
				Created = "TF-fake"
			}`,
		}),
	}
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEcsImageComponentsDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_ecs_image_component.default.image_component_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudEcsImageComponentsDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_ecs_image_component.default.image_component_name}_fake"`,
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEcsImageComponentsDataSourceName(rand, map[string]string{
			"ids":                  `["${alicloud_ecs_image_component.default.id}"]`,
			"image_component_name": `"${alicloud_ecs_image_component.default.image_component_name}"`,
			"name_regex":           `"${alicloud_ecs_image_component.default.image_component_name}"`,
			"owner":                `"SELF"`,
			"resource_group_id":    `"${alicloud_ecs_image_component.default.resource_group_id}"`,
			"tags": `{
				Created = "TF"
		}`,
		}),
		fakeConfig: testAccCheckAlicloudEcsImageComponentsDataSourceName(rand, map[string]string{
			"ids":                  `["${alicloud_ecs_image_component.default.id}_fake"]`,
			"image_component_name": `"${alicloud_ecs_image_component.default.image_component_name}_fake"`,
			"name_regex":           `"${alicloud_ecs_image_component.default.image_component_name}_fake"`,
			"owner":                `"ALIYUN"`,
			"resource_group_id":    `"${alicloud_ecs_image_component.default.resource_group_id}_fake"`,
			"tags": `{
				Created = "TF-fake"
			}`,
		}),
	}
	var existAlicloudEcsImageComponentsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                             "1",
			"names.#":                           "1",
			"components.#":                      "1",
			"components.0.component_type":       "Build",
			"components.0.create_time":          CHECKSET,
			"components.0.content":              "RUN yum update -y",
			"components.0.description":          fmt.Sprintf("tf-testAccImageComponent-%d", rand),
			"components.0.image_component_name": fmt.Sprintf("tf-testAccImageComponent-%d", rand),
			"components.0.resource_group_id":    CHECKSET,
			"components.0.image_component_id":   CHECKSET,
			"components.0.id":                   CHECKSET,
			"components.0.owner":                "SELF",
			"components.0.system_type":          "Linux",
			"components.0.tags.%":               "1",
			"components.0.tags.Created":         "TF",
		}
	}
	var fakeAlicloudEcsImageComponentsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"names.#": "0",
		}
	}
	var alicloudEcsImageComponentsCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_ecs_image_components.default",
		existMapFunc: existAlicloudEcsImageComponentsDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlicloudEcsImageComponentsDataSourceNameMapFunc,
	}

	preCheck := func() {
		testAccPreCheck(t)
	}
	alicloudEcsImageComponentsCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, imageComponentNameConf, ownerConf, resourceGroupIdConf, tagsConf, nameRegexConf, allConf)
}
func testAccCheckAlicloudEcsImageComponentsDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`

variable "name" {	
	default = "tf-testAccImageComponent-%d"
}

data "alicloud_resource_manager_resource_groups" "default" {
	name_regex = "default"
}

resource "alicloud_ecs_image_component" "default" {
	component_type = "Build"
	content = "RUN yum update -y"
	description = var.name
	image_component_name = var.name
	resource_group_id = data.alicloud_resource_manager_resource_groups.default.groups.0.id
	system_type = "Linux"
	tags = {
		Created = "TF"
	}
}

data "alicloud_ecs_image_components" "default" {	
	%s
}
`, rand, strings.Join(pairs, " \n "))
	return config
}

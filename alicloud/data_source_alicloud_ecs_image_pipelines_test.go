package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudECSImagePipelinesDataSource(t *testing.T) {
	rand := acctest.RandInt()
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudECSImagePipelinesDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_ecs_image_pipeline.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudECSImagePipelinesDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_ecs_image_pipeline.default.id}_fake"]`,
		}),
	}
	nameConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudECSImagePipelinesDataSourceName(rand, map[string]string{
			"ids":  `["${alicloud_ecs_image_pipeline.default.id}"]`,
			"name": `"${alicloud_ecs_image_pipeline.default.name}"`,
		}),
		fakeConfig: testAccCheckAlicloudECSImagePipelinesDataSourceName(rand, map[string]string{
			"ids":  `["${alicloud_ecs_image_pipeline.default.id}"]`,
			"name": `"${alicloud_ecs_image_pipeline.default.name}_fake"`,
		}),
	}
	resourceGroupIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudECSImagePipelinesDataSourceName(rand, map[string]string{
			"ids":               `["${alicloud_ecs_image_pipeline.default.id}"]`,
			"resource_group_id": `"${alicloud_ecs_image_pipeline.default.resource_group_id}"`,
		}),
		fakeConfig: testAccCheckAlicloudECSImagePipelinesDataSourceName(rand, map[string]string{
			"ids":               `["${alicloud_ecs_image_pipeline.default.id}"]`,
			"resource_group_id": `"${alicloud_ecs_image_pipeline.default.resource_group_id}_fake"`,
		}),
	}
	tagsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudECSImagePipelinesDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_ecs_image_pipeline.default.id}"]`,
			"tags": `{
				"Created" = "TF"
				"For" = "Acceptance-test"
		}`,
		}),
		fakeConfig: testAccCheckAlicloudECSImagePipelinesDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_ecs_image_pipeline.default.id}"]`,
			"tags": `{
				"Created" = "TF-fake"
				"For" = "Acceptance-test"
		}`,
		}),
	}
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudECSImagePipelinesDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_ecs_image_pipeline.default.name}"`,
		}),
		fakeConfig: testAccCheckAlicloudECSImagePipelinesDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_ecs_image_pipeline.default.name}_fake"`,
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudECSImagePipelinesDataSourceName(rand, map[string]string{
			"ids":               `["${alicloud_ecs_image_pipeline.default.id}"]`,
			"name":              `"${alicloud_ecs_image_pipeline.default.name}"`,
			"name_regex":        `"${alicloud_ecs_image_pipeline.default.name}"`,
			"resource_group_id": `"${alicloud_ecs_image_pipeline.default.resource_group_id}"`,
			"tags": `{
				"Created" = "TF"
				"For" = "Acceptance-test"
		}`,
		}),
		fakeConfig: testAccCheckAlicloudECSImagePipelinesDataSourceName(rand, map[string]string{
			"ids":               `["${alicloud_ecs_image_pipeline.default.id}_fake"]`,
			"name":              `"${alicloud_ecs_image_pipeline.default.name}_fake"`,
			"name_regex":        `"${alicloud_ecs_image_pipeline.default.name}_fake"`,
			"resource_group_id": `"${alicloud_ecs_image_pipeline.default.resource_group_id}_fake"`,
			"tags": `{
				"Created" = "TF-fake"
				"For" = "Acceptance-test"
		}`,
		}),
	}
	var existAlicloudEcsImagePipelinesDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                                  "1",
			"names.#":                                "1",
			"pipelines.#":                            "1",
			"pipelines.0.base_image":                 CHECKSET,
			"pipelines.0.base_image_type":            "IMAGE",
			"pipelines.0.build_content":              CHECKSET,
			"pipelines.0.delete_instance_on_failure": "false",
			"pipelines.0.description":                fmt.Sprintf("tf-testAccImagePipeline-%d", rand),
			"pipelines.0.image_name":                 fmt.Sprintf("tf-testAccImagePipeline-%d", rand),
			"pipelines.0.name":                       fmt.Sprintf("tf-testAccImagePipeline-%d", rand),
			"pipelines.0.instance_type":              CHECKSET,
			"pipelines.0.internet_max_bandwidth_out": "20",
			"pipelines.0.system_disk_size":           "40",
			"pipelines.0.to_region_id.#":             "2",
			"pipelines.0.add_account.#":              "0",
			"pipelines.0.creation_time":              CHECKSET,
			"pipelines.0.vswitch_id":                 CHECKSET,
			"pipelines.0.id":                         CHECKSET,
			"pipelines.0.image_pipeline_id":          CHECKSET,
			"pipelines.0.resource_group_id":          CHECKSET,
			"pipelines.0.tags.%":                     "2",
			"pipelines.0.tags.Created":               "TF",
			"pipelines.0.tags.For":                   "Acceptance-test",
		}
	}
	var fakeAlicloudEcsImagePipelinesDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"names.#": "0",
		}
	}
	var alicloudEcsImagePipelinesCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_ecs_image_pipelines.default",
		existMapFunc: existAlicloudEcsImagePipelinesDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlicloudEcsImagePipelinesDataSourceNameMapFunc,
	}

	preCheck := func() {
		testAccPreCheck(t)
	}
	alicloudEcsImagePipelinesCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, nameConf, resourceGroupIdConf, tagsConf, nameRegexConf, allConf)
}
func testAccCheckAlicloudECSImagePipelinesDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`

variable "name" {	
	default = "tf-testAccImagePipeline-%d"
}
data "alicloud_resource_manager_resource_groups" "default" {
  name_regex = "default"
}
data "alicloud_zones" "default" {}
data "alicloud_vpcs" "default" {
	name_regex = "default-NODELETING"
}
data "alicloud_vswitches" "default" {
	vpc_id = data.alicloud_vpcs.default.ids.0
	zone_id      = data.alicloud_zones.default.zones.0.id
}
data "alicloud_images" "default" {
  name_regex = "^ubuntu_[0-9]+_[0-9]+_x64*"
  most_recent = true
  owners      = "system"
}
data "alicloud_instance_types" "default" {
  image_id = data.alicloud_images.default.ids.0
}
resource "alicloud_ecs_image_pipeline" "default" {
	base_image = data.alicloud_images.default.ids.0
	base_image_type = "IMAGE"
	build_content = "RUN yum update -y"
	delete_instance_on_failure = false
	image_name = var.name
	name = var.name
	description = var.name
	instance_type = data.alicloud_instance_types.default.ids.0
	internet_max_bandwidth_out = 20
	system_disk_size = 40
	to_region_id = ["cn-qingdao", "cn-zhangjiakou"]
	vswitch_id = data.alicloud_vswitches.default.ids.0
	resource_group_id = data.alicloud_resource_manager_resource_groups.default.groups.0.id
	tags = {
		Created = "TF"
		For ="Acceptance-test"
	}
}

data "alicloud_ecs_image_pipelines" "default" {	
	%s
}
`, rand, strings.Join(pairs, " \n "))
	return config
}

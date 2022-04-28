package alicloud

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"strings"
	"testing"
)

func TestAccAlicloudEcsImagePipelineExecutionsDataSource(t *testing.T) {
	rand := acctest.RandInt()
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEcsImagePipelineExecutionsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_ecs_image_pipeline_execution.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudEcsImagePipelineExecutionsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_ecs_image_pipeline_execution.default.id}_fake"]`,
		}),
	}
	statusConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEcsImagePipelineExecutionsDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_ecs_image_pipeline_execution.default.id}"]`,
			"status": `"${alicloud_ecs_image_pipeline_execution.default.status}"`,
		}),
		fakeConfig: testAccCheckAlicloudEcsImagePipelineExecutionsDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_ecs_image_pipeline_execution.default.id}"]`,
			"status": `"BUILDING"`,
		}),
	}
	var existAlicloudEcsImagePipelineExecutionsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                                    "1",
			"executions.#":                             "1",
			"executions.0.image_pipeline_id":           CHECKSET,
			"executions.0.status":                      CHECKSET,
			"executions.0.id":                          CHECKSET,
			"executions.0.create_time":                 CHECKSET,
			"executions.0.image_id":                    "",
			"executions.0.resource_group_id":           "",
			"executions.0.message":                     CHECKSET,
			"executions.0.modified_time":               CHECKSET,
			"executions.0.image_pipeline_execution_id": CHECKSET,
		}
	}
	var fakeAlicloudEcsImagePipelineExecutionsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#": "0",
		}
	}
	var alicloudEcsImagePipelineExecutionsCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_ecs_image_pipeline_executions.default",
		existMapFunc: existAlicloudEcsImagePipelineExecutionsDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlicloudEcsImagePipelineExecutionsDataSourceNameMapFunc,
	}

	preCheck := func() {
		testAccPreCheck(t)
	}
	alicloudEcsImagePipelineExecutionsCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, statusConf)
}
func testAccCheckAlicloudEcsImagePipelineExecutionsDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`
variable "name" {	
	default = "tf-testAccImagePipelineExecution-%d"
}
data "alicloud_resource_manager_resource_groups" "default" {}
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

resource "alicloud_ecs_image_pipeline_execution" "default" {
	image_pipeline_id = alicloud_ecs_image_pipeline.default.id
}

data "alicloud_ecs_image_pipeline_executions" "default" {	
	image_pipeline_id = alicloud_ecs_image_pipeline.default.id
	%s
}
`, rand, strings.Join(pairs, " \n "))
	return config
}

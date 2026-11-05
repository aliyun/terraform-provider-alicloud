package ecs_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAliCloudECSImagePipelineExecutionsDataSource(t *testing.T) {
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%secsimagepipelineexecutionds%d", acctest.DefaultRegionToTest, rand)
	resourceId := "data.alicloud_ecs_image_pipeline_executions.default"
	baseConfig := AlicloudEcsImagePipelineExecutionBasicDependence8237(name)
	executionId := "${alicloud_ecs_image_pipeline_execution.default.id}"
	pipelineId := "${alicloud_ecs_image_pipeline.pipelineExection-pipeline.id}"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: baseConfig + imagePipelineExecutionsDataBlock(fmt.Sprintf("  execution_id = %s\n", executionId)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "executions.#", "1"),
					resource.TestCheckResourceAttr(resourceId, "ids.#", "1"),
					resource.TestCheckResourceAttr(resourceId, "total_count", "1"),
					resource.TestCheckResourceAttrPair(resourceId, "ids.0", "alicloud_ecs_image_pipeline_execution.default", "id"),
					resource.TestCheckResourceAttrPair(resourceId, "executions.0.execution_id", "alicloud_ecs_image_pipeline_execution.default", "id"),
					resource.TestCheckResourceAttrPair(resourceId, "executions.0.image_pipeline_id", "alicloud_ecs_image_pipeline.pipelineExection-pipeline", "id"),
					resource.TestMatchResourceAttr(resourceId, "executions.0.status", regexp.MustCompile(`^[A-Z_]+$`)),
					resource.TestCheckResourceAttrSet(resourceId, "executions.0.creation_time"),
				),
			},
			{
				Config: baseConfig + imagePipelineExecutionsDataBlock(fmt.Sprintf("  image_pipeline_id = %s\n", pipelineId)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "executions.#", "1"),
					resource.TestCheckResourceAttrPair(resourceId, "executions.0.execution_id", "alicloud_ecs_image_pipeline_execution.default", "id"),
				),
			},
			{
				Config: baseConfig + imagePipelineExecutionsDataBlock(fmt.Sprintf("  execution_id = %s\n  ids           = [%s]\n", executionId, executionId)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "executions.#", "1"),
					resource.TestCheckResourceAttrPair(resourceId, "ids.0", "alicloud_ecs_image_pipeline_execution.default", "id"),
				),
			},
			{
				Config: baseConfig + imagePipelineExecutionsDataBlock(fmt.Sprintf("  execution_id = %s\n  ids           = [\"exec-does-not-exist\"]\n", executionId)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "executions.#", "0"),
					resource.TestCheckResourceAttr(resourceId, "ids.#", "0"),
					resource.TestCheckResourceAttr(resourceId, "total_count", "0"),
				),
			},
			{
				Config: baseConfig + imagePipelineExecutionsDataBlock(fmt.Sprintf("  execution_id = %s\n  status       = \"CANCELLED\"\n", executionId)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "executions.#", "0"),
					resource.TestCheckResourceAttr(resourceId, "total_count", "0"),
				),
			},
			{
				Config: baseConfig + imagePipelineExecutionsDataBlock(fmt.Sprintf("  execution_id = %s\n  tags = {\n    env = \"tf-not-matching\"\n  }\n", executionId)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "executions.#", "0"),
				),
			},
			{
				Config: baseConfig + imagePipelineExecutionsDataBlock(fmt.Sprintf("  execution_id = %s\n  max_results  = 10\n", executionId)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "executions.#", "1"),
				),
			},
			{
				Config:      baseConfig + imagePipelineExecutionsDataBlock(fmt.Sprintf("  execution_id = %s\n  status       = \"SUCCESS,NOT_A_STATUS\"\n", executionId)),
				ExpectError: regexp.MustCompile(`Invalid Status|not a valid image pipeline execution status`),
			},
		},
	})
}

func imagePipelineExecutionsDataBlock(arguments string) string {
	return fmt.Sprintf(`
resource "alicloud_ecs_image_pipeline_execution" "default" {
  image_pipeline_id = alicloud_ecs_image_pipeline.pipelineExection-pipeline.id
}

data "alicloud_ecs_image_pipeline_executions" "default" {
%s}
`, arguments)
}

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

// A real image build is slow and its execution can outlive the data source
// assertion, which makes the shared image pipeline undeletable during teardown
// (InvalidImagePipeline.HasRunningExecution). The pipeline itself is created as
// a real fixture here; the non-empty execution mapping is covered by unit tests.
func TestAccAliCloudECSImagePipelineExecutionsDataSource(t *testing.T) {
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%secsimagepipelineexecutionds%d", acctest.DefaultRegionToTest, rand)
	resourceId := "data.alicloud_ecs_image_pipeline_executions.default"
	baseConfig := AlicloudEcsImagePipelineExecutionBasicDependence8237(name)
	pipelineFilter := "  image_pipeline_id = alicloud_ecs_image_pipeline.pipelineExection-pipeline.id\n"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: baseConfig + imagePipelineExecutionsDataBlock(pipelineFilter),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "executions.#", "0"),
					resource.TestCheckResourceAttr(resourceId, "ids.#", "0"),
					resource.TestCheckResourceAttr(resourceId, "total_count", "0"),
				),
			},
			{
				Config: baseConfig + imagePipelineExecutionsDataBlock(
					pipelineFilter+"  ids = [\"exec-does-not-exist\"]\n"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "executions.#", "0"),
					resource.TestCheckResourceAttr(resourceId, "ids.#", "0"),
					resource.TestCheckResourceAttr(resourceId, "total_count", "0"),
				),
			},
			{
				Config: baseConfig + imagePipelineExecutionsDataBlock(
					pipelineFilter+"  status      = \"CANCELLED\"\n  max_results = 10\n  tags = {\n    env = \"tf-not-matching\"\n  }\n"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "executions.#", "0"),
					resource.TestCheckResourceAttr(resourceId, "total_count", "0"),
				),
			},
			{
				Config:      baseConfig + imagePipelineExecutionsDataBlock(""),
				ExpectError: regexp.MustCompile(`At least one of|Invalid Attribute Combination`),
			},
			{
				Config:      baseConfig + imagePipelineExecutionsDataBlock(pipelineFilter+"  status = \"SUCCESS,NOT_A_STATUS\"\n"),
				ExpectError: regexp.MustCompile(`Invalid Status|not a valid image pipeline execution status`),
			},
			{
				// The destroy phase re-reads the last configuration, so a valid
				// step must follow the expected validation failures.
				Config: baseConfig + imagePipelineExecutionsDataBlock(pipelineFilter),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "executions.#", "0"),
				),
			},
		},
	})
}

func imagePipelineExecutionsDataBlock(arguments string) string {
	return fmt.Sprintf(`
data "alicloud_ecs_image_pipeline_executions" "default" {
%s}
`, arguments)
}

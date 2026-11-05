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

// The fixture starts a real execution so the query is proven against a real
// execution record instead of an empty result. The execution is cancelled as
// soon as the non-empty assertions have run, which also lets the shared image
// pipeline be deleted during teardown.
func TestAccAliCloudECSImagePipelineExecutionsDataSource(t *testing.T) {
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%secsimagepipelineexecutionds%d", acctest.DefaultRegionToTest, rand)
	resourceId := "data.alicloud_ecs_image_pipeline_executions.default"
	executionId := "alicloud_ecs_image_pipeline_execution.default"
	pipelineId := "alicloud_ecs_image_pipeline.default"
	statusPattern := regexp.MustCompile(`^(PREPARING|REPAIRING|BUILDING|TESTING|DISTRIBUTING|RELEASING|SUCCESS|PARTITION_SUCCESS|FAILED|TEST_FAILED|CANCELLING|CANCELLED)$`)
	executionConfig := "  image_pipeline_id = alicloud_ecs_image_pipeline.default.id\n"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: imagePipelineExecutionsConfig(name, executionConfig, "  execution_id = alicloud_ecs_image_pipeline_execution.default.id\n"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "executions.#", "1"),
					resource.TestCheckResourceAttr(resourceId, "ids.#", "1"),
					resource.TestCheckResourceAttrPair(resourceId, "ids.0", executionId, "id"),
					resource.TestCheckResourceAttrPair(resourceId, "executions.0.execution_id", executionId, "id"),
					resource.TestCheckResourceAttrPair(resourceId, "executions.0.image_pipeline_id", pipelineId, "id"),
					resource.TestCheckResourceAttrSet(resourceId, "executions.0.status"),
					resource.TestMatchResourceAttr(resourceId, "executions.0.status", statusPattern),
					resource.TestCheckResourceAttr(resourceId, "total_count", "1"),
				),
			},
			{
				Config: imagePipelineExecutionsConfig(name, executionConfig, `  execution_id = alicloud_ecs_image_pipeline_execution.default.id
  ids          = [alicloud_ecs_image_pipeline_execution.default.id]
`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "executions.#", "1"),
					resource.TestCheckResourceAttrPair(resourceId, "ids.0", executionId, "id"),
					resource.TestCheckResourceAttr(resourceId, "total_count", "1"),
				),
			},
			{
				Config: imagePipelineExecutionsConfig(name, executionConfig, `  execution_id = alicloud_ecs_image_pipeline_execution.default.id
  ids          = ["exec-does-not-exist"]
`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "executions.#", "0"),
					resource.TestCheckResourceAttr(resourceId, "ids.#", "0"),
					resource.TestCheckResourceAttr(resourceId, "total_count", "0"),
				),
			},
			{
				Config:      imagePipelineExecutionsConfig(name, executionConfig, "  execution_id = alicloud_ecs_image_pipeline_execution.default.id\n  status       = \"SUCCESS,,FAILED\"\n"),
				ExpectError: regexp.MustCompile(`Every comma separated status must be non-empty|Invalid Status`),
			},
			{
				Config:      imagePipelineExecutionsConfig(name, executionConfig, "  execution_id = alicloud_ecs_image_pipeline_execution.default.id\n  status       = \",\"\n"),
				ExpectError: regexp.MustCompile(`Every comma separated status must be non-empty|Invalid Status`),
			},
			{
				// Cancelling the execution after the non-empty assertions keeps the
				// fixture deletable, and the cancelled status is a deterministic
				// non-empty value for the follow-up read.
				Config: imagePipelineExecutionsConfig(name, executionConfig+"  status            = \"CANCELLED\"\n", `  execution_id = alicloud_ecs_image_pipeline_execution.default.id
  status       = "CANCELLED"
`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "executions.#", "1"),
					resource.TestCheckResourceAttr(resourceId, "executions.0.status", "CANCELLED"),
					resource.TestCheckResourceAttrPair(resourceId, "executions.0.execution_id", executionId, "id"),
				),
			},
			{
				Config: imagePipelineExecutionsConfig(name, executionConfig, "  execution_id = alicloud_ecs_image_pipeline_execution.default.id\n"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "executions.#", "1"),
				),
			},
		},
	})
}

func imagePipelineExecutionsConfig(name, executionArguments, dataArguments string) string {
	return fmt.Sprintf(`
variable "name" {
  default = %q
}

data "alicloud_instance_types" "default" {
  instance_type_family = "ecs.g6"
  cpu_core_count       = 2
}

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
  available_instance_type     = data.alicloud_instance_types.default.instance_types.0.id
}

data "alicloud_images" "default" {
  name_regex   = "^ubuntu"
  most_recent  = true
  owners       = "system"
  architecture = "x86_64"
}

resource "alicloud_vpc" "default" {
  vpc_name   = var.name
  cidr_block = "172.16.0.0/16"
}

resource "alicloud_vswitch" "default" {
  vpc_id       = alicloud_vpc.default.id
  cidr_block   = "172.16.0.0/24"
  zone_id      = data.alicloud_zones.default.zones.0.id
  vswitch_name = var.name
}

resource "alicloud_ecs_image_pipeline" "default" {
  base_image                 = data.alicloud_images.default.images.0.id
  base_image_type            = "IMAGE"
  vswitch_id                 = alicloud_vswitch.default.id
  instance_type              = data.alicloud_instance_types.default.instance_types.0.id
  build_content              = "RUN echo terraform-acceptance-test"
  system_disk_size           = 40
  internet_max_bandwidth_out = 5
  delete_instance_on_failure = true
  name                       = var.name
}

resource "alicloud_ecs_image_pipeline_execution" "default" {
%s}

data "alicloud_ecs_image_pipeline_executions" "default" {
%s}
`, name, executionArguments, dataArguments)
}

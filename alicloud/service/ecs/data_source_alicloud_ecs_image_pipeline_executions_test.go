package ecs_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// The account used for acceptance tests cannot start an image build: every
// StartImagePipelineExecution attempt stops while preparing the transition
// instance with ErrorCode NO_AVAILABLE_INSTANCE_TYPE. The fixture therefore
// keeps a real image pipeline and exercises the live query, the complete
// filter request and the empty-result semantics; the non-empty object mapping
// and the pagination logic are covered by the data source unit tests.
func TestAccAliCloudECSImagePipelineExecutionsDataSource(t *testing.T) {
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%secsimagepipelineexecutionds%d", acctest.DefaultRegionToTest, rand)
	resourceId := "data.alicloud_ecs_image_pipeline_executions.default"
	pipelineFilter := "  image_pipeline_id = alicloud_ecs_image_pipeline.default.id\n"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: imagePipelineExecutionsConfig(name, pipelineFilter),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "executions.#", "0"),
					resource.TestCheckResourceAttr(resourceId, "ids.#", "0"),
					resource.TestCheckResourceAttr(resourceId, "total_count", "0"),
				),
			},
			{
				Config: imagePipelineExecutionsConfig(name, pipelineFilter+"  ids = [\"exec-does-not-exist\"]\n"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "executions.#", "0"),
					resource.TestCheckResourceAttr(resourceId, "ids.#", "0"),
					resource.TestCheckResourceAttr(resourceId, "total_count", "0"),
				),
			},
			{
				Config: imagePipelineExecutionsConfig(name, pipelineFilter+"  status      = \"CANCELLED\"\n  max_results = 10\n  tags = {\n    env = \"tf-not-matching\"\n  }\n"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "executions.#", "0"),
					resource.TestCheckResourceAttr(resourceId, "total_count", "0"),
				),
			},
			{
				Config:      imagePipelineExecutionsConfig(name, ""),
				ExpectError: regexp.MustCompile(`At least one of|Invalid Attribute Combination`),
			},
			{
				Config:      imagePipelineExecutionsConfig(name, pipelineFilter+"  status = \"SUCCESS,NOT_A_STATUS\"\n"),
				ExpectError: regexp.MustCompile(`Invalid Status|not a valid image pipeline execution status`),
			},
			{
				Config:      imagePipelineExecutionsConfig(name, pipelineFilter+"  status = \",\"\n"),
				ExpectError: regexp.MustCompile(`Every comma separated status must be non-empty|Invalid Status`),
			},
			{
				Config:      imagePipelineExecutionsConfig(name, pipelineFilter+"  status = \"SUCCESS,,FAILED\"\n"),
				ExpectError: regexp.MustCompile(`Every comma separated status must be non-empty|Invalid Status`),
			},
			{
				// The destroy phase re-reads the last configuration, so a valid
				// step must follow the expected validation failures.
				Config: imagePipelineExecutionsConfig(name, pipelineFilter),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "executions.#", "0"),
				),
			},
		},
	})
}

func imagePipelineExecutionsConfig(name, dataArguments string) string {
	return fmt.Sprintf(`
variable "name" {
  default = %q
}

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
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
  build_content              = "RUN echo terraform-acceptance-test"
  system_disk_size           = 40
  internet_max_bandwidth_out = 5
  delete_instance_on_failure = true
  name                       = var.name
}

data "alicloud_ecs_image_pipeline_executions" "default" {
%s}
`, name, dataArguments)
}

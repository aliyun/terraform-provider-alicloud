package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test Ecs ImagePipelineExecution. >>> Resource test cases, automatically generated.
// Case ImagePipelineExecution-status 8237
func TestAccAliCloudEcsImagePipelineExecution_basic8237(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ecs_image_pipeline_execution.default"
	ra := resourceAttrInit(resourceId, AlicloudEcsImagePipelineExecutionMap8237)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EcsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEcsImagePipelineExecution")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%secsimagepipelineexecution%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudEcsImagePipelineExecutionBasicDependence8237)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"image_pipeline_id": "${alicloud_ecs_image_pipeline.pipelineExection-pipeline.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"image_pipeline_id": CHECKSET,
						"status":            CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "CANCELLED",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "CANCELLED",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

var AlicloudEcsImagePipelineExecutionMap8237 = map[string]string{
	"create_time": CHECKSET,
}

func AlicloudEcsImagePipelineExecutionBasicDependence8237(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_vpc" "pipelineExecution-vpc" {
  description = "test-pipeline"
  enable_ipv6 = true
  vpc_name    = var.name
}

resource "alicloud_vswitch" "vs" {
  description  = "pipelineExecution-start"
  vpc_id       = alicloud_vpc.pipelineExecution-vpc.id
  cidr_block   = "172.16.0.0/24"
  vswitch_name = format("%%s1", var.name)
  zone_id      = "cn-hangzhou-i"
}

resource "alicloud_ecs_image_pipeline" "pipelineExection-pipeline" {
  base_image_type            = "IMAGE"
  description                = "test"
  system_disk_size           = "40"
  vswitch_id                 = alicloud_vswitch.vs.id
  add_account                = ["1284387915995949"]
  image_name                 = "test-image-pipeline"
  delete_instance_on_failure = true
  internet_max_bandwidth_out = "5"
  to_region_id               = ["cn-beijing"]
  base_image                 = "aliyun_3_x64_20G_dengbao_alibase_20240819.vhd"
  build_content              = "COMPONENT ic-bp122acttbs2sxdyq2ky"
}


`, name)
}

// Case ImagePipelineExecution-start 8232
func TestAccAliCloudEcsImagePipelineExecution_basic8232(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ecs_image_pipeline_execution.default"
	ra := resourceAttrInit(resourceId, AlicloudEcsImagePipelineExecutionMap8232)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EcsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEcsImagePipelineExecution")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%secsimagepipelineexecution%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudEcsImagePipelineExecutionBasicDependence8232)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"image_pipeline_id": "${alicloud_ecs_image_pipeline.pipelineExection-pipeline.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"image_pipeline_id": CHECKSET,
						"status":            CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "CANCELLED",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "CANCELLED",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

var AlicloudEcsImagePipelineExecutionMap8232 = map[string]string{
	"create_time": CHECKSET,
}

func AlicloudEcsImagePipelineExecutionBasicDependence8232(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_vpc" "pipelineExecution-vpc" {
  description = "test-pipeline"
  enable_ipv6 = true
  vpc_name    = var.name
}

resource "alicloud_vswitch" "vs" {
  description  = "pipelineExecution-start-test"
  vpc_id       = alicloud_vpc.pipelineExecution-vpc.id
  cidr_block   = "172.16.0.0/24"
  vswitch_name = format("%%s1", var.name)
  zone_id      = "cn-hangzhou-i"
}

resource "alicloud_ecs_image_pipeline" "pipelineExection-pipeline" {
  base_image_type            = "IMAGE"
  description                = "test"
  system_disk_size           = "40"
  vswitch_id                 = alicloud_vswitch.vs.id
  add_account                = ["1284387915995949"]
  image_name                 = "test-image-pipeline"
  delete_instance_on_failure = true
  internet_max_bandwidth_out = "5"
  to_region_id               = ["cn-hangzhou"]
  base_image                 = "aliyun_3_x64_20G_dengbao_alibase_20240819.vhd"
  build_content              = "COMPONENT ic-bp122acttbs2sxdyq2ky"
}


`, name)
}

// Test Ecs ImagePipelineExecution. <<< Resource test cases, automatically generated.

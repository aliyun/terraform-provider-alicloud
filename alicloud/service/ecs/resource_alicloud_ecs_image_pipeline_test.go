package ecs_test

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/ecs"

	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAliCloudECSImagePipeline_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ecs_image_pipeline.default"
	ra := resourceAttrInit(resourceId, AlicloudECSImagePipelineMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return ecs.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeEcsImagePipeline")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%secsimagepipeline%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudECSImagePipelineBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"base_image_type": "IMAGE",
					"base_image":      "${data.alicloud_images.default.ids.0}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"base_image_type": "IMAGE",
						"base_image":      CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"build_content"},
			},
		},
	})
}
func TestAccAliCloudECSImagePipeline_basic1(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ecs_image_pipeline.default"
	ra := resourceAttrInit(resourceId, AlicloudECSImagePipelineMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return ecs.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeEcsImagePipeline")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%secsimagepipeline%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudECSImagePipelineBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"base_image_type":            "IMAGE",
					"name":                       "${var.name}",
					"description":                "${var.name}",
					"system_disk_size":           "40",
					"vswitch_id":                 "${data.alicloud_vswitches.default.ids.0}",
					"image_name":                 "${var.name}",
					"delete_instance_on_failure": "true",
					"internet_max_bandwidth_out": "20",
					"to_region_id":               []string{"cn-qingdao", "cn-zhangjiakou"},
					"add_account":                []string{"${data.alicloud_account.default.id}"},
					"instance_type":              "${data.alicloud_instance_types.default.ids.0}",
					"base_image":                 "${data.alicloud_images.default.ids.0}",
					"build_content":              "RUN yum update -y",
					"resource_group_id":          "${data.alicloud_resource_manager_resource_groups.default.groups.0.id}",
					"tags":                       map[string]interface{}{"Created": "TF", "Foo": "Bar"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"base_image_type":            "IMAGE",
						"description":                name,
						"system_disk_size":           "40",
						"vswitch_id":                 CHECKSET,
						"image_name":                 name,
						"delete_instance_on_failure": "true",
						"internet_max_bandwidth_out": "20",
						"to_region_id.#":             "2",
						"add_account.#":              "1",
						"instance_type":              CHECKSET,
						"base_image":                 CHECKSET,
						"resource_group_id":          CHECKSET,
						"tags.%":                     "2",
						"tags.Created":               "TF",
						"tags.Foo":                   "Bar",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"build_content"},
			},
		},
	})
}

var AlicloudECSImagePipelineMap0 = map[string]string{}

func AlicloudECSImagePipelineBasicDependence0(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}
data "alicloud_resource_manager_resource_groups" "default" {
  name_regex = "default"
}
data "alicloud_zones" "default" {}
data "alicloud_vpcs" "default" {
	name_regex = "^default-NODELETING$"
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
data "alicloud_account" "default" {}
`, name)
}

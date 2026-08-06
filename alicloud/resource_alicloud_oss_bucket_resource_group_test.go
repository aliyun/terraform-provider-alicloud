package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAliCloudOssBucketResourceGroupBasic(t *testing.T) {
	var v string

	resourceId := "alicloud_oss_bucket_resource_group.default"
	ra := resourceAttrInit(resourceId, ossBucketResourceGroupMap)

	serviceFunc := func() interface{} {
		return &OssService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)

	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc-bucket-resourcegroup-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceOssBucketResourceGroupDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"bucket":            "${local.bucket_src}",
					"resource_group_id": "${local.resource_group_id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bucket":            name + "-t-1",
						"resource_group_id": CHECKSET,
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"bucket":            "${local.bucket_src}",
					"resource_group_id": "${local.resource_group_id1}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bucket":            name + "-t-1",
						"resource_group_id": CHECKSET,
					}),
				),
			},			
		},
	})
}

func resourceOssBucketResourceGroupDependence(name string) string {
	return fmt.Sprintf(`

	variable "name" {
		default = "%s-t"
	}

	data "alicloud_resource_manager_resource_groups" "default" {
		name_regex = "default"
	}

	resource "alicloud_oss_bucket" "bucket_src" {
		bucket = "${var.name}-1"
	}

	locals {
		bucket_src = alicloud_oss_bucket.bucket_src.id
		resource_group_id  = data.alicloud_resource_manager_resource_groups.default.groups.0.id
		resource_group_id1 = data.alicloud_resource_manager_resource_groups.default.groups.1.id
	}
`, name)
}

var ossBucketResourceGroupMap = map[string]string{}

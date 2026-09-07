package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAliCloudImageImport_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_image_import.default"
	ra := resourceAttrInit(resourceId, AliCloudImageImportMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EcsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeImageById")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000, 9999)
	name := fmt.Sprintf("tf-testacc%simageimport%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudImageImportBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"image_name": name,
					"disk_device_mapping": []map[string]interface{}{
						{
							"oss_bucket": "${alicloud_oss_bucket.default.id}",
							"oss_object": "${alicloud_oss_bucket_object.default.id}",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"image_name":            name,
						"disk_device_mapping.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"image_name": name + "-update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"image_name": name + "-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": name,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"license_type"},
			},
		},
	})

}

func TestAccAliCloudImageImport_basic0_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_image_import.default"
	ra := resourceAttrInit(resourceId, AliCloudImageImportMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EcsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeImageById")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000, 9999)
	name := fmt.Sprintf("tf-testacc%simageimport%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudImageImportBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"architecture": "i386",
					"os_type":      "linux",
					"platform":     "Aliyun",
					"boot_mode":    "UEFI",
					"license_type": "Auto",
					"image_name":   name,
					"description":  name,
					"disk_device_mapping": []map[string]interface{}{
						{
							"format":          "RAW",
							"oss_bucket":      "${alicloud_oss_bucket.default.id}",
							"oss_object":      "${alicloud_oss_bucket_object.default.id}",
							"device":          "/dev/xvda",
							"disk_image_size": "6",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"architecture":          "i386",
						"os_type":               "linux",
						"platform":              "Aliyun",
						"boot_mode":             "UEFI",
						"image_name":            name,
						"description":           name,
						"disk_device_mapping.#": "1",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"license_type"},
			},
		},
	})

}

var AliCloudImageImportMap0 = map[string]string{
	"platform":   CHECKSET,
	"boot_mode":  CHECKSET,
	"image_name": CHECKSET,
}

func AliCloudImageImportBasicDependence0(name string) string {
	return fmt.Sprintf(`
	variable "name" {
  		default = "%s"
	}

	resource "alicloud_oss_bucket" "default" {
  		bucket = var.name
	}

	resource "alicloud_oss_bucket_object" "default" {
  		bucket  = alicloud_oss_bucket.default.id
  		key     = "fc/hello.zip"
  		content = <<EOF
		# -*- coding: utf-8 -*-
		def handler(event, context):
		print "hello world"
		return 'hello world'
		EOF
	}
`, name)
}

package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
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
		IDRefreshName:     resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
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
		IDRefreshName:     resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
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

// imageImportIdTracker records the image id in one step and asserts it stays
// the same afterwards, so in-place updates (e.g. features.nvme_support) can
// be verified not to replace the image.
func imageImportIdTracker(resourceId string) (record, assertKept resource.TestCheckFunc) {
	imageId := ""
	readId := func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceId]
		if !ok {
			return "", fmt.Errorf("not found: %s", resourceId)
		}
		return rs.Primary.ID, nil
	}
	record = func(s *terraform.State) error {
		id, err := readId(s)
		if err != nil {
			return err
		}
		imageId = id
		return nil
	}
	assertKept = func(s *terraform.State) error {
		id, err := readId(s)
		if err != nil {
			return err
		}
		if id != imageId {
			return fmt.Errorf("the image was replaced, but updating features.nvme_support should modify the image in place and keep the image id")
		}
		return nil
	}
	return record, assertKept
}

func TestAccAliCloudImageImport_features(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_image_import.default"
	recordImageId, assertImageIdKept := imageImportIdTracker(resourceId)
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
		IDRefreshName:     resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"image_name": name,
					"features": []map[string]interface{}{
						{
							"nvme_support": "supported",
						},
					},
					"disk_device_mapping": []map[string]interface{}{
						{
							"oss_bucket": "${alicloud_oss_bucket.default.id}",
							"oss_object": "${alicloud_oss_bucket_object.default.id}",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					recordImageId,
					testAccCheck(map[string]string{
						"image_name":              name,
						"features.#":              "1",
						"features.0.nvme_support": "supported",
						"disk_device_mapping.#":   "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"image_name": name,
					"features": []map[string]interface{}{
						{
							"nvme_support": "unsupported",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					assertImageIdKept,
					testAccCheck(map[string]string{
						"image_name":              name,
						"features.0.nvme_support": "unsupported",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"image_name": name,
					"features": []map[string]interface{}{
						{
							"nvme_support": "supported",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					assertImageIdKept,
					testAccCheck(map[string]string{
						"image_name":              name,
						"features.0.nvme_support": "supported",
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

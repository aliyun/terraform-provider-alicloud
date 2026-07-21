package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test Fcv3 LayerVersion. >>> Resource test cases, automatically generated.
// Case TestLayer_ZIP 6990
func TestAccAliCloudFcv3LayerVersion_basic6990(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_fcv3_layer_version.default"
	ra := resourceAttrInit(resourceId, AlicloudFcv3LayerVersionMap6990)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &Fcv3ServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeFcv3LayerVersion")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sfcv3layerversion%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudFcv3LayerVersionBasicDependence6990)
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
					"description": "luoni.fz",
					//"compatible_runtime": []string{
					//	"custom.debian10", "python3.10", "python3.9"},
					"layer_name": "FC3LayerResouceTest_ZIP_2024SepWed",
					"public":     "true",
					"license":    "Apache2.0",
					"acl":        "0",
					"code": []map[string]interface{}{
						{
							"zip_file": "UEsDBBQACAAIAAAAAAAAAAAAAAAAAAAAAAAIAAAAaW5kZXgucHmEkEFKxEAQRfd9ig9ZTCJOooIwDMwNXLqXnnQlaalUhU5lRj2KZ/FOXkESGR114bJ/P/7jV4b1xRq1hijtFpM1682cuNgPmgysbRulPT0fRxXnMtwrSPyeCdYRokSLnuMLJTTkbUqEvDMbxm1VdcRD6Tk+T1LW2ldB66knsYdA5iNX17ebm6tN2VnPhcswMPmREPuBacb+CiapLarAj9gT6/H97dVlCNScY3mtYvRkxdZlwDKDEnanPWVLdrdkeXEGlFEazVdfPVHaVeHc3N15CUwppwOJXeK7HshAB8NuOU7J6sP4SRXuH/EvbUfMiqMmDqv5M5FNSfAj/wgAAP//UEsHCPl//NYAAQAArwEAAFBLAQIUABQACAAIAAAAAAD5f/zWAAEAAK8BAAAIAAAAAAAAAAAAAAAAAAAAAABpbmRleC5weVBLBQYAAAAAAQABADYAAAA2AQAAAAA=",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "luoni.fz",
						//"compatible_runtime.#": "3",
						"layer_name": CHECKSET,
						"license":    "Apache2.0",
						"acl":        "0",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"acl": "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"acl": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"acl": "0",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"acl": "0",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"code", "public"},
			},
		},
	})
}

var AlicloudFcv3LayerVersionMap6990 = map[string]string{
	"create_time": CHECKSET,
	"version":     CHECKSET,
}

func AlicloudFcv3LayerVersionBasicDependence6990(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}

// Case TestLayer_OSS 7069
func TestAccAliCloudFcv3LayerVersion_basic7069(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_fcv3_layer_version.default"
	ra := resourceAttrInit(resourceId, AlicloudFcv3LayerVersionMap7069)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &Fcv3ServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeFcv3LayerVersion")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sfcv3layerversion%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudFcv3LayerVersionBasicDependence7069)
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
					"description": "luoni.fz",
					"compatible_runtime": []string{
						"custom.debian10", "python3.10", "python3.9"},
					"layer_name": "FC3LayerResouceTest_OSS_2024SepWed",
					"license":    "Apache2.0",
					"acl":        "0",
					"code": []map[string]interface{}{
						{
							"oss_bucket_name": "${alicloud_oss_bucket.default.bucket}",
							"oss_object_name": "${alicloud_oss_bucket_object.default.id}",
							"checksum":        "4270285996107335518",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":          "luoni.fz",
						"compatible_runtime.#": "3",
						"layer_name":           CHECKSET,
						"license":              "Apache2.0",
						"acl":                  "0",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"acl": "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"acl": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"acl": "0",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"acl": "0",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"code"},
			},
		},
	})
}

var AlicloudFcv3LayerVersionMap7069 = map[string]string{
	"create_time": CHECKSET,
	"version":     CHECKSET,
}

func AlicloudFcv3LayerVersionBasicDependence7069(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_oss_bucket" "default" {
  bucket = var.name
}

resource "alicloud_oss_bucket_object" "default" {
  bucket  = alicloud_oss_bucket.default.bucket
  key     = "fc3Py39.zip"
  content = "print('hello')"
}

`, name)
}

// Case TestLayer_OSS_Online 7337
func TestAccAliCloudFcv3LayerVersion_basic7337(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_fcv3_layer_version.default"
	ra := resourceAttrInit(resourceId, AlicloudFcv3LayerVersionMap7337)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &Fcv3ServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeFcv3LayerVersion")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sfcv3layerversion%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudFcv3LayerVersionBasicDependence7337)
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
					"description": "luoni.fz",
					"compatible_runtime": []string{
						"custom.debian10", "python3.10", "python3.9"},
					"layer_name": "FC3LayerResouceTest_OSS_2024SepWed",
					"license":    "Apache2.0",
					"acl":        "0",
					"code": []map[string]interface{}{
						{
							"oss_bucket_name": "${alicloud_oss_bucket.default.bucket}",
							"oss_object_name": "${alicloud_oss_bucket_object.default.id}",
							"checksum":        "4270285996107335518",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":          "luoni.fz",
						"compatible_runtime.#": "3",
						"layer_name":           CHECKSET,
						"license":              "Apache2.0",
						"acl":                  "0",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"acl": "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"acl": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"acl": "0",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"acl": "0",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"code"},
			},
		},
	})
}

var AlicloudFcv3LayerVersionMap7337 = map[string]string{
	"create_time": CHECKSET,
	"version":     CHECKSET,
}

func AlicloudFcv3LayerVersionBasicDependence7337(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_oss_bucket" "default" {
  bucket = var.name
}

resource "alicloud_oss_bucket_object" "default" {
  bucket  = alicloud_oss_bucket.default.bucket
  key     = "fc3Py39.zip"
  content = "print('hello')"
}
`, name)
}

// Case TestLayer_ZIP_Online 7336
func TestAccAliCloudFcv3LayerVersion_basic7336(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_fcv3_layer_version.default"
	ra := resourceAttrInit(resourceId, AlicloudFcv3LayerVersionMap7336)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &Fcv3ServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeFcv3LayerVersion")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sfcv3layerversion%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudFcv3LayerVersionBasicDependence7336)
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
					"description": "luoni.fz",
					"compatible_runtime": []string{
						"custom.debian10", "python3.10", "python3.9"},
					"layer_name": "FC3LayerResouceTest_ZIP_2024SepWed",
					"license":    "Apache2.0",
					"acl":        "0",
					"code": []map[string]interface{}{
						{
							"zip_file": "UEsDBBQACAAIAAAAAAAAAAAAAAAAAAAAAAAIAAAAaW5kZXgucHmEkEFKxEAQRfd9ig9ZTCJOooIwDMwNXLqXnnQlaalUhU5lRj2KZ/FOXkESGR114bJ/P/7jV4b1xRq1hijtFpM1682cuNgPmgysbRulPT0fRxXnMtwrSPyeCdYRokSLnuMLJTTkbUqEvDMbxm1VdcRD6Tk+T1LW2ldB66knsYdA5iNX17ebm6tN2VnPhcswMPmREPuBacb+CiapLarAj9gT6/H97dVlCNScY3mtYvRkxdZlwDKDEnanPWVLdrdkeXEGlFEazVdfPVHaVeHc3N15CUwppwOJXeK7HshAB8NuOU7J6sP4SRXuH/EvbUfMiqMmDqv5M5FNSfAj/wgAAP//UEsHCPl//NYAAQAArwEAAFBLAQIUABQACAAIAAAAAAD5f/zWAAEAAK8BAAAIAAAAAAAAAAAAAAAAAAAAAABpbmRleC5weVBLBQYAAAAAAQABADYAAAA2AQAAAAA=",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":          "luoni.fz",
						"compatible_runtime.#": "3",
						"layer_name":           CHECKSET,
						"license":              "Apache2.0",
						"acl":                  "0",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"acl": "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"acl": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"acl": "0",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"acl": "0",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"code"},
			},
		},
	})
}

var AlicloudFcv3LayerVersionMap7336 = map[string]string{
	"create_time": CHECKSET,
	"version":     CHECKSET,
}

func AlicloudFcv3LayerVersionBasicDependence7336(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}

// Test Fcv3 LayerVersion. <<< Resource test cases, automatically generated.

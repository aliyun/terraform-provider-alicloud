package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test Aligreen ImageLib. >>> Resource test cases, automatically generated.
// Case 图片库_副本1721974637924 7332
func TestAccAliCloudAligreenImageLib_basic7332(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_aligreen_image_lib.default"
	ra := resourceAttrInit(resourceId, AlicloudAligreenImageLibMap7332)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &AligreenServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeAligreenImageLib")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf_testacc_aligreenimagelib%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudAligreenImageLibBasicDependence7332)
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
					"category":       "BLACK",
					"scene":          "PORN",
					"image_lib_name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"category":       "BLACK",
						"scene":          "PORN",
						"image_lib_name": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"image_lib_name": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"image_lib_name": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"category":       "BLACK",
					"enable":         "true",
					"scene":          "PORN",
					"image_lib_name": name + "_update",
					"biz_types": []string{
						"test_007"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"category":       "BLACK",
						"enable":         "true",
						"scene":          "PORN",
						"image_lib_name": name + "_update",
						"biz_types.#":    "1",
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

var AlicloudAligreenImageLibMap7332 = map[string]string{}

func AlicloudAligreenImageLibBasicDependence7332(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_aligreen_biz_type" "defaultUalunB" {
  biz_type_name = var.name
}


`, name)
}

// Case 图片库 7317
func TestAccAliCloudAligreenImageLib_basic7317(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_aligreen_image_lib.default"
	ra := resourceAttrInit(resourceId, AlicloudAligreenImageLibMap7317)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &AligreenServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeAligreenImageLib")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf_testacc_aligreenimagelib%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudAligreenImageLibBasicDependence7317)
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
					"category":       "BLACK",
					"scene":          "PORN",
					"image_lib_name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"category":       "BLACK",
						"scene":          "PORN",
						"image_lib_name": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"image_lib_name": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"image_lib_name": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"category":       "BLACK",
					"enable":         "true",
					"scene":          "PORN",
					"image_lib_name": name + "_update",
					"biz_types": []string{
						"test_007"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"category":       "BLACK",
						"enable":         "true",
						"scene":          "PORN",
						"image_lib_name": name + "_update",
						"biz_types.#":    "1",
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

var AlicloudAligreenImageLibMap7317 = map[string]string{}

func AlicloudAligreenImageLibBasicDependence7317(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_aligreen_biz_type" "defaultUalunB" {
 biz_type_name = var.name
}


`, name)
}

// Case 图片库_副本1721974637924 7332  twin
func TestAccAliCloudAligreenImageLib_basic7332_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_aligreen_image_lib.default"
	ra := resourceAttrInit(resourceId, AlicloudAligreenImageLibMap7332)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &AligreenServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeAligreenImageLib")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf_testacc_aligreenimagelib%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudAligreenImageLibBasicDependence7332)
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
					"category":       "BLACK",
					"enable":         "true",
					"scene":          "PORN",
					"image_lib_name": name,
					"biz_types": []string{
						"test_007"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"category":       "BLACK",
						"enable":         "true",
						"scene":          "PORN",
						"image_lib_name": name,
						"biz_types.#":    "1",
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

// Case 图片库 7317  twin
func TestAccAliCloudAligreenImageLib_basic7317_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_aligreen_image_lib.default"
	ra := resourceAttrInit(resourceId, AlicloudAligreenImageLibMap7317)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &AligreenServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeAligreenImageLib")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf_testacc_aligreenimagelib%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudAligreenImageLibBasicDependence7317)
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
					"category":       "BLACK",
					"enable":         "true",
					"scene":          "PORN",
					"image_lib_name": name,
					"biz_types": []string{
						"test_007"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"category":       "BLACK",
						"enable":         "true",
						"scene":          "PORN",
						"image_lib_name": name,
						"biz_types.#":    "1",
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

// Case 图片库_副本1721974637924 7332  raw
func TestAccAliCloudAligreenImageLib_basic7332_raw(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_aligreen_image_lib.default"
	ra := resourceAttrInit(resourceId, AlicloudAligreenImageLibMap7332)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &AligreenServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeAligreenImageLib")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf_testacc_aligreenimagelib%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudAligreenImageLibBasicDependence7332)
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
					"category":       "BLACK",
					"enable":         "true",
					"scene":          "PORN",
					"image_lib_name": name,
					"biz_types": []string{
						"${alicloud_aligreen_biz_type.defaultUalunB.biz_type_name}"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"category":       "BLACK",
						"enable":         "true",
						"scene":          "PORN",
						"image_lib_name": name,
						"biz_types.#":    "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"image_lib_name": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"image_lib_name": name + "_update",
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

// Case 图片库 7317  raw
func TestAccAliCloudAligreenImageLib_basic7317_raw(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_aligreen_image_lib.default"
	ra := resourceAttrInit(resourceId, AlicloudAligreenImageLibMap7317)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &AligreenServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeAligreenImageLib")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf_testacc_aligreenimagelib%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudAligreenImageLibBasicDependence7317)
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
					"category":       "BLACK",
					"enable":         "true",
					"scene":          "PORN",
					"image_lib_name": name,
					"biz_types": []string{
						"test_007"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"category":       "BLACK",
						"enable":         "true",
						"scene":          "PORN",
						"image_lib_name": name,
						"biz_types.#":    "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"image_lib_name": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"image_lib_name": name + "_update",
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

// Test Aligreen ImageLib. <<< Resource test cases, automatically generated.

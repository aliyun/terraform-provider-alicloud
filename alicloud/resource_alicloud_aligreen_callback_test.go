package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test Aligreen Callback. >>> Resource test cases, automatically generated.
// Case 回调 7327
func TestAccAliCloudAligreenCallback_basic7327(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_aligreen_callback.default"
	ra := resourceAttrInit(resourceId, AlicloudAligreenCallbackMap7327)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &AligreenServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeAligreenCallback")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%saligreencallback%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudAligreenCallbackBasicDependence7327)
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
					"callback_url":  "https://www.aliyun.com",
					"callback_name": name,
					"callback_types": []string{
						"machineScan", "selfAudit", "test"},
					"callback_suggestions": []string{
						"block", "review", "pass"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"callback_url":           "https://www.aliyun.com",
						"callback_name":          name,
						"callback_types.#":       "3",
						"callback_suggestions.#": "3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"callback_url": "https://www.aliyun.com/",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"callback_url": "https://www.aliyun.com/",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"callback_name": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"callback_name": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"callback_types": []string{
						"machineScan"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"callback_types.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"callback_suggestions": []string{
						"review"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"callback_suggestions.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"callback_url":  "https://www.aliyun.com",
					"crypt_type":    "0",
					"callback_name": name + "_update",
					"callback_types": []string{
						"machineScan", "selfAudit", "test"},
					"callback_suggestions": []string{
						"block", "review", "pass"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"callback_url":           "https://www.aliyun.com",
						"crypt_type":             "0",
						"callback_name":          name + "_update",
						"callback_types.#":       "3",
						"callback_suggestions.#": "3",
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

var AlicloudAligreenCallbackMap7327 = map[string]string{
	"create_time": CHECKSET,
}

func AlicloudAligreenCallbackBasicDependence7327(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}

// Case 回调 7260
func TestAccAliCloudAligreenCallback_basic7260(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_aligreen_callback.default"
	ra := resourceAttrInit(resourceId, AlicloudAligreenCallbackMap7260)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &AligreenServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeAligreenCallback")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%saligreencallback%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudAligreenCallbackBasicDependence7260)
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
					"callback_url":  "https://www.aliyun.com",
					"callback_name": name,
					"callback_types": []string{
						"machineScan", "selfAudit", "test"},
					"callback_suggestions": []string{
						"block", "review", "pass"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"callback_url":           "https://www.aliyun.com",
						"callback_name":          name,
						"callback_types.#":       "3",
						"callback_suggestions.#": "3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"callback_url": "https://www.aliyun.com/",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"callback_url": "https://www.aliyun.com/",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"callback_name": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"callback_name": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"callback_types": []string{
						"machineScan"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"callback_types.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"callback_suggestions": []string{
						"review"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"callback_suggestions.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"callback_url":  "https://www.aliyun.com",
					"crypt_type":    "0",
					"callback_name": name + "_update",
					"callback_types": []string{
						"machineScan", "selfAudit", "test"},
					"callback_suggestions": []string{
						"block", "review", "pass"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"callback_url":           "https://www.aliyun.com",
						"crypt_type":             "0",
						"callback_name":          name + "_update",
						"callback_types.#":       "3",
						"callback_suggestions.#": "3",
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

var AlicloudAligreenCallbackMap7260 = map[string]string{
	"create_time": CHECKSET,
}

func AlicloudAligreenCallbackBasicDependence7260(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}

// Case 回调 7327  twin
func TestAccAliCloudAligreenCallback_basic7327_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_aligreen_callback.default"
	ra := resourceAttrInit(resourceId, AlicloudAligreenCallbackMap7327)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &AligreenServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeAligreenCallback")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%saligreencallback%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudAligreenCallbackBasicDependence7327)
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
					"callback_url":  "https://www.aliyun.com",
					"crypt_type":    "0",
					"callback_name": name,
					"callback_types": []string{
						"machineScan", "selfAudit", "test"},
					"callback_suggestions": []string{
						"block", "review", "pass"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"callback_url":           "https://www.aliyun.com",
						"crypt_type":             "0",
						"callback_name":          name,
						"callback_types.#":       "3",
						"callback_suggestions.#": "3",
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

// Case 回调 7260  twin
func TestAccAliCloudAligreenCallback_basic7260_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_aligreen_callback.default"
	ra := resourceAttrInit(resourceId, AlicloudAligreenCallbackMap7260)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &AligreenServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeAligreenCallback")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%saligreencallback%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudAligreenCallbackBasicDependence7260)
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
					"callback_url":  "https://www.aliyun.com",
					"crypt_type":    "0",
					"callback_name": name,
					"callback_types": []string{
						"machineScan", "selfAudit", "test"},
					"callback_suggestions": []string{
						"block", "review", "pass"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"callback_url":           "https://www.aliyun.com",
						"crypt_type":             "0",
						"callback_name":          name,
						"callback_types.#":       "3",
						"callback_suggestions.#": "3",
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

// Case 回调 7327  raw
func TestAccAliCloudAligreenCallback_basic7327_raw(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_aligreen_callback.default"
	ra := resourceAttrInit(resourceId, AlicloudAligreenCallbackMap7327)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &AligreenServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeAligreenCallback")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%saligreencallback%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudAligreenCallbackBasicDependence7327)
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
					"callback_url":  "https://www.aliyun.com",
					"crypt_type":    "0",
					"callback_name": name,
					"callback_types": []string{
						"machineScan", "selfAudit", "test"},
					"callback_suggestions": []string{
						"block", "review", "pass"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"callback_url":           "https://www.aliyun.com",
						"crypt_type":             "0",
						"callback_name":          name,
						"callback_types.#":       "3",
						"callback_suggestions.#": "3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"callback_url":  "https://www.aliyun.com/",
					"crypt_type":    "1",
					"callback_name": name + "_update",
					"callback_types": []string{
						"machineScan"},
					"callback_suggestions": []string{
						"review"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"callback_url":           "https://www.aliyun.com/",
						"crypt_type":             "1",
						"callback_name":          name + "_update",
						"callback_types.#":       "1",
						"callback_suggestions.#": "1",
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

// Case 回调 7260  raw
func TestAccAliCloudAligreenCallback_basic7260_raw(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_aligreen_callback.default"
	ra := resourceAttrInit(resourceId, AlicloudAligreenCallbackMap7260)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &AligreenServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeAligreenCallback")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%saligreencallback%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudAligreenCallbackBasicDependence7260)
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
					"callback_url":  "https://www.aliyun.com",
					"crypt_type":    "0",
					"callback_name": name,
					"callback_types": []string{
						"machineScan", "selfAudit", "test"},
					"callback_suggestions": []string{
						"block", "review", "pass"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"callback_url":           "https://www.aliyun.com",
						"crypt_type":             "0",
						"callback_name":          name,
						"callback_types.#":       "3",
						"callback_suggestions.#": "3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"callback_url":  "https://www.aliyun.com/",
					"crypt_type":    "1",
					"callback_name": name + "_update",
					"callback_types": []string{
						"machineScan"},
					"callback_suggestions": []string{
						"review"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"callback_url":           "https://www.aliyun.com/",
						"crypt_type":             "1",
						"callback_name":          name + "_update",
						"callback_types.#":       "1",
						"callback_suggestions.#": "1",
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

// Test Aligreen Callback. <<< Resource test cases, automatically generated.

package aligreen_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	aligreen "github.com/aliyun/terraform-provider-alicloud/alicloud/service/aligreen"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Test Aligreen Callback. >>> Resource test cases, automatically generated.
// Case 回调 7327
func TestAccAliCloudAligreenCallback_basic7327(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_aligreen_callback.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudAligreenCallbackMap7327)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return aligreen.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeAligreenCallback")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%saligreencallback%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudAligreenCallbackBasicDependence7327)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
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
	"create_time": acctest.CHECKSET,
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
	ra := acctest.ResourceAttrInit(resourceId, AlicloudAligreenCallbackMap7260)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return aligreen.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeAligreenCallback")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%saligreencallback%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudAligreenCallbackBasicDependence7260)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
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
	"create_time": acctest.CHECKSET,
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
	ra := acctest.ResourceAttrInit(resourceId, AlicloudAligreenCallbackMap7327)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return aligreen.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeAligreenCallback")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%saligreencallback%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudAligreenCallbackBasicDependence7327)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
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
	ra := acctest.ResourceAttrInit(resourceId, AlicloudAligreenCallbackMap7260)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return aligreen.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeAligreenCallback")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%saligreencallback%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudAligreenCallbackBasicDependence7260)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
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
	ra := acctest.ResourceAttrInit(resourceId, AlicloudAligreenCallbackMap7327)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return aligreen.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeAligreenCallback")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%saligreencallback%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudAligreenCallbackBasicDependence7327)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
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
	ra := acctest.ResourceAttrInit(resourceId, AlicloudAligreenCallbackMap7260)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return aligreen.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeAligreenCallback")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%saligreencallback%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudAligreenCallbackBasicDependence7260)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
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

// Case 回调 collection-order coverage
// Covers the configurable multi-member TypeList attributes callback_types and
// callback_suggestions: apply the original member order, then a convergent apply
// with each list's members reordered (one list per step, so the checker's
// permutation comparison sees every other field unchanged).
func TestAccAliCloudAligreenCallback_collectionOrder(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_aligreen_callback.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudAligreenCallbackMap7327)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return aligreen.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeAligreenCallback")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%saligreencallback%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudAligreenCallbackBasicDependence7327)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
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
					"callback_types": []string{
						"test", "selfAudit", "machineScan"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"callback_types.#":       "3",
						"callback_suggestions.#": "3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"callback_suggestions": []string{
						"pass", "review", "block"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"callback_types.#":       "3",
						"callback_suggestions.#": "3",
					}),
				),
			},
		},
	})
}

// Test Aligreen Callback. <<< Resource test cases, automatically generated.

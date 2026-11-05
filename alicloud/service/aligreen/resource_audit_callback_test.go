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

// Test Aligreen AuditCallback. >>> Resource test cases, automatically generated.
// Case 机审回调通知 7328
func TestAccAliCloudAligreenAuditCallback_basic7328(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_aligreen_audit_callback.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudAligreenAuditCallbackMap7328)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return aligreen.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeAligreenAuditCallback")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%saligreenauditcallback%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudAligreenAuditCallbackBasicDependence7328)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"audit_callback_name": name,
					"crypt_type":          "SM3",
					"url":                 "https://www.aliyun.com/",
					"callback_types": []string{
						"aliyunAudit", "selfAduit", "example"},
					"callback_suggestions": []string{
						"block", "review", "pass"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"audit_callback_name":    name,
						"url":                    "https://www.aliyun.com/",
						"callback_types.#":       "3",
						"callback_suggestions.#": "3",
						"crypt_type":             "SM3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"crypt_type": "SM3",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"crypt_type": "SM3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"crypt_type": "SHA256",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"crypt_type": "SHA256",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"url": "https://www.aliyun.com",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"url": "https://www.aliyun.com",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"callback_types": []string{
						"selfAudit"},
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
						"block"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"callback_suggestions.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"crypt_type":          "SM3",
					"audit_callback_name": name + "_update",
					"url":                 "https://www.aliyun.com/",
					"callback_types": []string{
						"aliyunAudit", "selfAduit", "example"},
					"callback_suggestions": []string{
						"block", "review", "pass"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"crypt_type":             "SM3",
						"audit_callback_name":    name + "_update",
						"url":                    "https://www.aliyun.com/",
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

var AlicloudAligreenAuditCallbackMap7328 = map[string]string{}

func AlicloudAligreenAuditCallbackBasicDependence7328(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}

// Case 机审回调通知 7324
func TestAccAliCloudAligreenAuditCallback_basic7324(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_aligreen_audit_callback.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudAligreenAuditCallbackMap7324)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return aligreen.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeAligreenAuditCallback")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%saligreenauditcallback%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudAligreenAuditCallbackBasicDependence7324)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"audit_callback_name": name,
					"url":                 "https://www.aliyun.com/",
					"crypt_type":          "SM3",
					"callback_types": []string{
						"aliyunAudit", "selfAduit", "example"},
					"callback_suggestions": []string{
						"block", "review", "pass"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"audit_callback_name":    name,
						"url":                    "https://www.aliyun.com/",
						"crypt_type":             "SM3",
						"callback_types.#":       "3",
						"callback_suggestions.#": "3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"crypt_type": "SM3",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"crypt_type": "SM3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"crypt_type": "SHA256",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"crypt_type": "SHA256",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"url": "https://www.aliyun.com",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"url": "https://www.aliyun.com",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"callback_types": []string{
						"selfAudit"},
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
						"block"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"callback_suggestions.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"crypt_type":          "SM3",
					"audit_callback_name": name + "_update",
					"url":                 "https://www.aliyun.com/",
					"callback_types": []string{
						"aliyunAudit", "selfAduit", "example"},
					"callback_suggestions": []string{
						"block", "review", "pass"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"crypt_type":             "SM3",
						"audit_callback_name":    name + "_update",
						"url":                    "https://www.aliyun.com/",
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

var AlicloudAligreenAuditCallbackMap7324 = map[string]string{}

func AlicloudAligreenAuditCallbackBasicDependence7324(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}

// Case 机审回调通知 7328  twin
func TestAccAliCloudAligreenAuditCallback_basic7328_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_aligreen_audit_callback.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudAligreenAuditCallbackMap7328)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return aligreen.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeAligreenAuditCallback")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%saligreenauditcallback%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudAligreenAuditCallbackBasicDependence7328)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"crypt_type":          "SM3",
					"audit_callback_name": name,
					"url":                 "https://www.aliyun.com/",
					"callback_types": []string{
						"aliyunAudit", "selfAduit", "example"},
					"callback_suggestions": []string{
						"block", "review", "pass"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"crypt_type":             "SM3",
						"audit_callback_name":    name,
						"url":                    "https://www.aliyun.com/",
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

// Case 机审回调通知 7324  twin
func TestAccAliCloudAligreenAuditCallback_basic7324_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_aligreen_audit_callback.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudAligreenAuditCallbackMap7324)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return aligreen.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeAligreenAuditCallback")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%saligreenauditcallback%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudAligreenAuditCallbackBasicDependence7324)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"crypt_type":          "SM3",
					"audit_callback_name": name,
					"url":                 "https://www.aliyun.com/",
					"callback_types": []string{
						"aliyunAudit", "selfAduit", "example"},
					"callback_suggestions": []string{
						"block", "review", "pass"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"crypt_type":             "SM3",
						"audit_callback_name":    name,
						"url":                    "https://www.aliyun.com/",
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

// Case 机审回调通知 7328  raw
func TestAccAliCloudAligreenAuditCallback_basic7328_raw(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_aligreen_audit_callback.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudAligreenAuditCallbackMap7328)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return aligreen.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeAligreenAuditCallback")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%saligreenauditcallback%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudAligreenAuditCallbackBasicDependence7328)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"crypt_type":          "SM3",
					"audit_callback_name": name,
					"url":                 "https://www.aliyun.com/",
					"callback_types": []string{
						"aliyunAudit", "selfAduit", "example"},
					"callback_suggestions": []string{
						"block", "review", "pass"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"crypt_type":             "SM3",
						"audit_callback_name":    name,
						"url":                    "https://www.aliyun.com/",
						"callback_types.#":       "3",
						"callback_suggestions.#": "3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"crypt_type": "SHA256",
					"url":        "https://www.aliyun.com",
					"callback_types": []string{
						"selfAudit"},
					"callback_suggestions": []string{
						"block"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"crypt_type":             "SHA256",
						"url":                    "https://www.aliyun.com",
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

// Case 机审回调通知 7324  raw
func TestAccAliCloudAligreenAuditCallback_basic7324_raw(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_aligreen_audit_callback.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudAligreenAuditCallbackMap7324)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return aligreen.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeAligreenAuditCallback")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%saligreenauditcallback%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudAligreenAuditCallbackBasicDependence7324)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"crypt_type":          "SM3",
					"audit_callback_name": name,
					"url":                 "https://www.aliyun.com/",
					"callback_types": []string{
						"aliyunAudit", "selfAduit", "example"},
					"callback_suggestions": []string{
						"block", "review", "pass"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"crypt_type":             "SM3",
						"audit_callback_name":    name,
						"url":                    "https://www.aliyun.com/",
						"callback_types.#":       "3",
						"callback_suggestions.#": "3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"crypt_type": "SHA256",
					"url":        "https://www.aliyun.com",
					"callback_types": []string{
						"selfAudit"},
					"callback_suggestions": []string{
						"block"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"crypt_type":             "SHA256",
						"url":                    "https://www.aliyun.com",
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

// Case 机审回调通知 collection-order coverage
// Covers the configurable multi-member TypeList attributes callback_types and
// callback_suggestions: apply the original member order, then a convergent apply
// with each list's members reordered (one list per step, so the checker's
// permutation comparison sees every other field unchanged).
func TestAccAliCloudAligreenAuditCallback_collectionOrder(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_aligreen_audit_callback.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudAligreenAuditCallbackMap7328)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return aligreen.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeAligreenAuditCallback")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%saligreenauditcallback%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudAligreenAuditCallbackBasicDependence7328)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"audit_callback_name": name,
					"url":                 "https://www.aliyun.com/",
					"crypt_type":          "SM3",
					"callback_types": []string{
						"aliyunAudit", "selfAduit", "example"},
					"callback_suggestions": []string{
						"block", "review", "pass"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"audit_callback_name":    name,
						"url":                    "https://www.aliyun.com/",
						"crypt_type":             "SM3",
						"callback_types.#":       "3",
						"callback_suggestions.#": "3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"callback_types": []string{
						"example", "selfAduit", "aliyunAudit"},
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

// Test Aligreen AuditCallback. <<< Resource test cases, automatically generated.

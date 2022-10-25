package alicloud

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	roa "github.com/alibabacloud-go/tea-roa/client"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/stretchr/testify/assert"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudSAEGreyTagRoute_basic0(t *testing.T) {
	var v map[string]interface{}
	checkoutSupportedRegions(t, true, connectivity.SaeSupportRegions)
	resourceId := "alicloud_sae_grey_tag_route.default"
	ra := resourceAttrInit(resourceId, AlicloudSAEGreyTagRouteMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &SaeService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeSaeGreyTagRoute")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1, 100)
	name := fmt.Sprintf("tftestacc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudSAEGreyTagRouteBasicDependence0)
	resource.Test(t, resource.TestCase{
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"description": name,
					"app_id":      "${alicloud_sae_application.default.id}",
					"sc_rules": []map[string]interface{}{
						{
							"items": []map[string]interface{}{
								{
									"type":     "param",
									"name":     "tftest",
									"operator": "rawvalue",
									"value":    "test",
									"cond":     "==",
								},
								{
									"type":     "param",
									"name":     "tftest",
									"operator": "rawvalue",
									"value":    "test1",
									"cond":     "!=",
								},
							},
							"path":      "/tf/test",
							"condition": "AND",
						},
					},
					"dubbo_rules": []map[string]interface{}{
						{
							"items": []map[string]interface{}{
								{
									"cond":     "==",
									"expr":     ".key1",
									"index":    "1",
									"operator": "rawvalue",
									"value":    "value1",
								},
								{
									"cond":     "==",
									"expr":     ".key2",
									"index":    "0",
									"operator": "rawvalue",
									"value":    "value2",
								},
							},
							"condition":    "OR",
							"group":        "DUBBO",
							"method_name":  "test",
							"service_name": "com.test.service",
							"version":      "1.0.0",
						},
					},
					"grey_tag_route_name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":         name,
						"sc_rules.#":          "1",
						"dubbo_rules.#":       "1",
						"grey_tag_route_name": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": name + "update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": name + "update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"sc_rules": []map[string]interface{}{
						{
							"items": []map[string]interface{}{
								{
									"type":     "param",
									"name":     "tftest",
									"operator": "rawvalue",
									"value":    "test3",
									"cond":     "==",
								},
								{
									"type":     "param",
									"name":     "tftest",
									"operator": "rawvalue",
									"value":    "test2",
									"cond":     "!=",
								},
							},
							"path":      "/tf/test1",
							"condition": "OR",
						},
					},
					"dubbo_rules": []map[string]interface{}{
						{
							"items": []map[string]interface{}{
								{
									"cond":     "!=",
									"expr":     ".key1",
									"index":    "1",
									"operator": "rawvalue",
									"value":    "value3",
								},
								{
									"cond":     "!=",
									"expr":     ".key2",
									"index":    "0",
									"operator": "rawvalue",
									"value":    "value3",
								},
							},
							"condition":    "AND",
							"group":        "DUBBO",
							"method_name":  "test",
							"service_name": "com.test.service",
							"version":      "1.2.0",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"sc_rules.#":    "1",
						"dubbo_rules.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": name,
					"sc_rules": []map[string]interface{}{
						{
							"items": []map[string]interface{}{
								{
									"type":     "param",
									"name":     "tftest",
									"operator": "rawvalue",
									"value":    "test",
									"cond":     "==",
								},
								{
									"type":     "param",
									"name":     "tftest",
									"operator": "rawvalue",
									"value":    "test1",
									"cond":     "!=",
								},
							},
							"path":      "/tf/test",
							"condition": "AND",
						},
					},
					"dubbo_rules": []map[string]interface{}{
						{
							"items": []map[string]interface{}{
								{
									"cond":     "==",
									"expr":     ".key1",
									"index":    "1",
									"operator": "rawvalue",
									"value":    "value1",
								},
								{
									"cond":     "==",
									"expr":     ".key2",
									"index":    "0",
									"operator": "rawvalue",
									"value":    "value2",
								},
							},
							"condition":    "OR",
							"group":        "DUBBO",
							"method_name":  "test",
							"service_name": "com.test.service",
							"version":      "1.0.0",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":   name,
						"sc_rules.#":    "1",
						"dubbo_rules.#": "1",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{""},
			},
		},
	})
}

var AlicloudSAEGreyTagRouteMap0 = map[string]string{}

func AlicloudSAEGreyTagRouteBasicDependence0(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}
data "alicloud_vpcs" "default" {
  name_regex = "default-NODELETING"
}

data "alicloud_vswitches" "default" {
  vpc_id = data.alicloud_vpcs.default.ids.0
}

resource "alicloud_sae_namespace" "default" {
  namespace_description = var.name
  namespace_id          = join(":",["%s",var.name])
  namespace_name        = var.name
}

resource "alicloud_sae_application" "default" {
  app_description = var.name
  app_name        = var.name
  namespace_id    = alicloud_sae_namespace.default.namespace_id
  image_url       = "registry-vpc.cn-hangzhou.aliyuncs.com/lxepoo/apache-php5"
  package_type    = "Image"
  jdk             = "Open JDK 8"
  vswitch_id      = data.alicloud_vswitches.default.ids.0
  vpc_id          = data.alicloud_vpcs.default.ids.0
  timezone        = "Asia/Shanghai"
  replicas        = "5"
  cpu             = "500"
  memory          = "2048"
}
`, name, defaultRegionToTest)
}

func TestUnitAlicloudSAEGreyTagRoute(t *testing.T) {
	p := Provider().(*schema.Provider).ResourcesMap
	d, _ := schema.InternalMap(p["alicloud_sae_grey_tag_route"].Schema).Data(nil, nil)
	dCreate, _ := schema.InternalMap(p["alicloud_sae_grey_tag_route"].Schema).Data(nil, nil)
	dCreate.MarkNewResource()
	for key, value := range map[string]interface{}{
		"description": "description",
		"app_id":      "app_id",
		"sc_rules": []map[string]interface{}{
			{
				"items": []map[string]interface{}{
					{
						"type":     "param",
						"name":     "name",
						"operator": "rawvalue",
						"value":    "value",
						"cond":     "==",
					},
				},
				"path":      "/tf/test",
				"condition": "AND",
			},
		},
		"dubbo_rules": []map[string]interface{}{
			{
				"items": []map[string]interface{}{
					{
						"cond":     "==",
						"expr":     "expr",
						"index":    1,
						"operator": "rawvalue",
						"value":    "value",
					},
				},
				"condition":    "OR",
				"group":        "DUBBO",
				"method_name":  "test",
				"service_name": "com.test.service",
				"version":      "1.0.0",
			},
		},
		"grey_tag_route_name": "grey_tag_route_name",
	} {
		err := dCreate.Set(key, value)
		assert.Nil(t, err)
		err = d.Set(key, value)
		assert.Nil(t, err)
	}
	region := os.Getenv("ALICLOUD_REGION")
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		t.Skipf("Skipping the test case with err: %s", err)
		t.Skipped()
	}
	rawClient = rawClient.(*connectivity.AliyunClient)
	ReadMockResponse := map[string]interface{}{
		"body": map[string]interface{}{
			"Data": map[string]interface{}{
				"GreyTagRouteId": "grey_tag_route_id",
				"AppId":          "app_id",
				"Name":           "rule-grey_tag_route_name",
				"Description":    "description",
				"ScRules": []interface{}{
					map[string]interface{}{
						"items": []interface{}{
							map[string]interface{}{
								"type":     "param",
								"name":     "name",
								"operator": "rawvalue",
								"value":    "value",
								"cond":     "==",
							},
						},
						"path":      "/tf/test",
						"condition": "AND",
					},
				},
				"DubboRules": []interface{}{
					map[string]interface{}{
						"items": []interface{}{
							map[string]interface{}{
								"cond":     "==",
								"expr":     "expr",
								"index":    1,
								"operator": "rawvalue",
								"value":    "value",
							},
						},
						"condition":   "OR",
						"group":       "DUBBO",
						"methodName":  "test",
						"serviceName": "com.test.service",
						"version":     "1.0.0",
					},
				},
			},
		},
	}

	responseMock := map[string]func(errorCode string) (map[string]interface{}, error){
		"RetryError": func(errorCode string) (map[string]interface{}, error) {
			return nil, &tea.SDKError{
				Code:    String(errorCode),
				Data:    String(errorCode),
				Message: String(errorCode),
			}
		},
		"NotFoundError": func(errorCode string) (map[string]interface{}, error) {
			return nil, GetNotFoundErrorFromString(GetNotFoundMessage("alicloud_sae_grey_tag_route", "grey_tag_route_id"))
		},
		"NoRetryError": func(errorCode string) (map[string]interface{}, error) {
			return nil, &tea.SDKError{
				Code:    String(errorCode),
				Data:    String(errorCode),
				Message: String(errorCode),
			}
		},
		"CreateNormal": func(errorCode string) (map[string]interface{}, error) {
			result := ReadMockResponse
			return result, nil
		},
		"UpdateNormal": func(errorCode string) (map[string]interface{}, error) {
			result := ReadMockResponse
			return result, nil
		},
		"UpdateStatusNormal": func(errorCode string) (map[string]interface{}, error) {
			result := ReadMockResponse
			return result, nil
		},
		"DeleteNormal": func(errorCode string) (map[string]interface{}, error) {
			result := ReadMockResponse
			return result, nil
		},
		"ReadNormal": func(errorCode string) (map[string]interface{}, error) {
			result := ReadMockResponse
			return result, nil
		},
	}
	// Create
	t.Run("CreateClientAbnormal", func(t *testing.T) {
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewServerlessClient", func(_ *connectivity.AliyunClient) (*roa.Client, error) {
			return nil, &tea.SDKError{
				Code:    String("loadEndpoint error"),
				Data:    String("loadEndpoint error"),
				Message: String("loadEndpoint error"),
			}
		})
		err := resourceAlicloudSaeGreyTagRouteCreate(d, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
	})
	t.Run("CreateAbnormal", func(t *testing.T) {
		retryFlag := true
		noRetryFlag := true
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&roa.Client{}), "DoRequest", func(_ *roa.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]*string, _ map[string]*string, _ interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if retryFlag {
				retryFlag = false
				return responseMock["RetryError"]("Throttling")
			} else if noRetryFlag {
				noRetryFlag = false
				return responseMock["NoRetryError"]("NonRetryableError")
			}
			return responseMock["CreateNormal"]("")
		})
		err := resourceAlicloudSaeGreyTagRouteCreate(d, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
	})
	t.Run("CreateNormal", func(t *testing.T) {
		retryFlag := false
		noRetryFlag := false
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&roa.Client{}), "DoRequest", func(_ *roa.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]*string, _ map[string]*string, _ interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if retryFlag {
				retryFlag = false
				return responseMock["RetryError"]("Throttling")
			} else if noRetryFlag {
				noRetryFlag = false
				return responseMock["NoRetryError"]("NonRetryableError")
			}
			return responseMock["CreateNormal"]("")
		})
		err := resourceAlicloudSaeGreyTagRouteCreate(dCreate, rawClient)
		patches.Reset()
		assert.Nil(t, err)
	})
	t.Run("CreateRetryableError", func(t *testing.T) {
		retryFlag := false
		noRetryFlag := true
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&roa.Client{}), "DoRequest", func(_ *roa.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]*string, _ map[string]*string, _ interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if retryFlag {
				return responseMock["RetryError"]("Throttling")
			} else if noRetryFlag {
				noRetryFlag = false
				return responseMock["NoRetryError"]("NonRetryableError")
			}
			return responseMock["CreateNormal"]("")
		})
		err := resourceAlicloudSaeGreyTagRouteCreate(d, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
	})

	// Set ID for Update and Delete Method
	d.SetId("grey_tag_route_id")

	// Update
	t.Run("UpdateModifySaeGreyTagRouteAttributeAbnormal", func(t *testing.T) {
		diff := terraform.NewInstanceDiff()
		for _, key := range []string{"description"} {
			switch p["alicloud_sae_grey_tag_route"].Schema[key].Type {
			case schema.TypeString:
				diff.SetAttribute(key, &terraform.ResourceAttrDiff{Old: d.Get(key).(string), New: d.Get(key).(string) + "_update"})
			case schema.TypeBool:
				diff.SetAttribute(key, &terraform.ResourceAttrDiff{Old: strconv.FormatBool(d.Get(key).(bool)), New: strconv.FormatBool(true)})
			case schema.TypeInt:
				diff.SetAttribute(key, &terraform.ResourceAttrDiff{Old: strconv.Itoa(d.Get(key).(int)), New: strconv.Itoa(1200)})
			case schema.TypeMap:
				diff.SetAttribute("tags.%", &terraform.ResourceAttrDiff{Old: "0", New: "2"})
				diff.SetAttribute("tags.For", &terraform.ResourceAttrDiff{Old: "", New: "Test"})
				diff.SetAttribute("tags.Created", &terraform.ResourceAttrDiff{Old: "", New: "TF"})
			}
		}
		resourceData1, _ := schema.InternalMap(p["alicloud_sae_grey_tag_route"].Schema).Data(nil, diff)
		resourceData1.SetId(d.Id())
		retryFlag := true
		noRetryFlag := true
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&roa.Client{}), "DoRequest", func(_ *roa.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]*string, _ map[string]*string, _ interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if retryFlag {
				retryFlag = false
				return responseMock["RetryError"]("Throttling")
			} else if noRetryFlag {
				noRetryFlag = false
				return responseMock["NoRetryError"]("NonRetryableError")
			}
			return responseMock["UpdateNormal"]("")
		})
		err := resourceAlicloudSaeGreyTagRouteUpdate(resourceData1, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
	})
	t.Run("UpdateModifySaeGreyTagRouteAttributeNoRetryErrorAbnormal", func(t *testing.T) {
		diff := terraform.NewInstanceDiff()
		for _, key := range []string{"description"} {
			switch p["alicloud_sae_grey_tag_route"].Schema[key].Type {
			case schema.TypeString:
				diff.SetAttribute(key, &terraform.ResourceAttrDiff{Old: d.Get(key).(string), New: d.Get(key).(string) + "_update"})
			case schema.TypeBool:
				diff.SetAttribute(key, &terraform.ResourceAttrDiff{Old: strconv.FormatBool(d.Get(key).(bool)), New: strconv.FormatBool(true)})
			case schema.TypeInt:
				diff.SetAttribute(key, &terraform.ResourceAttrDiff{Old: strconv.Itoa(d.Get(key).(int)), New: strconv.Itoa(1200)})
			case schema.TypeMap:
				diff.SetAttribute("tags.%", &terraform.ResourceAttrDiff{Old: "0", New: "2"})
				diff.SetAttribute("tags.For", &terraform.ResourceAttrDiff{Old: "", New: "Test"})
				diff.SetAttribute("tags.Created", &terraform.ResourceAttrDiff{Old: "", New: "TF"})
			}
		}
		resourceData, _ := schema.InternalMap(p["alicloud_sae_grey_tag_route"].Schema).Data(nil, diff)
		resourceData.SetId(d.Id())
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&roa.Client{}), "DoRequest", func(_ *roa.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]*string, _ map[string]*string, _ interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			return responseMock["NoRetryError"]("NoRetryError")
		})
		err := resourceAlicloudSaeGreyTagRouteUpdate(resourceData, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
	})
	t.Run("UpdateModifySaeGreyTagRouteAttributeNormal", func(t *testing.T) {
		diff := terraform.NewInstanceDiff()
		for _, key := range []string{"description", "dubbo_rules", "sc_rules"} {
			switch p["alicloud_sae_grey_tag_route"].Schema[key].Type {
			case schema.TypeString:
				diff.SetAttribute(key, &terraform.ResourceAttrDiff{Old: d.Get(key).(string), New: d.Get(key).(string) + "_update"})
			case schema.TypeBool:
				diff.SetAttribute(key, &terraform.ResourceAttrDiff{Old: strconv.FormatBool(d.Get(key).(bool)), New: strconv.FormatBool(true)})
			case schema.TypeInt:
				diff.SetAttribute(key, &terraform.ResourceAttrDiff{Old: strconv.Itoa(d.Get(key).(int)), New: strconv.Itoa(1200)})
			case schema.TypeMap:
				diff.SetAttribute("tags.%", &terraform.ResourceAttrDiff{Old: "0", New: "2"})
				diff.SetAttribute("tags.For", &terraform.ResourceAttrDiff{Old: "", New: "Test"})
				diff.SetAttribute("tags.Created", &terraform.ResourceAttrDiff{Old: "", New: "TF"})
			case schema.TypeSet:
				diff.SetAttribute(fmt.Sprintf("%s.#", key), &terraform.ResourceAttrDiff{Old: "1", New: "1"})
				diff.SetAttribute(fmt.Sprintf("%s.0.items.#", key), &terraform.ResourceAttrDiff{Old: "1", New: "1"})
				diff.SetAttribute(fmt.Sprintf("%s.0.items.0.value", key), &terraform.ResourceAttrDiff{Old: "value1", New: "value2"})
				diff.SetAttribute(fmt.Sprintf("%s.0.items.0.operator", key), &terraform.ResourceAttrDiff{Old: "rawvalue", New: "list"})
				diff.SetAttribute(fmt.Sprintf("%s.0.condition", key), &terraform.ResourceAttrDiff{Old: "AND", New: "OR"})
			}
		}
		resourceData1, _ := schema.InternalMap(p["alicloud_sae_grey_tag_route"].Schema).Data(nil, diff)
		resourceData1.SetId(d.Id())
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&roa.Client{}), "DoRequest", func(_ *roa.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]*string, _ map[string]*string, _ interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			return responseMock["UpdateNormal"]("")
		})
		err := resourceAlicloudSaeGreyTagRouteUpdate(resourceData1, rawClient)
		patches.Reset()
		assert.Nil(t, err)
	})

	// Delete
	t.Run("DeleteClientAbnormal", func(t *testing.T) {
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewServerlessClient", func(_ *connectivity.AliyunClient) (*roa.Client, error) {
			return nil, &tea.SDKError{
				Code:    String("loadEndpoint error"),
				Data:    String("loadEndpoint error"),
				Message: String("loadEndpoint error"),
			}
		})
		err := resourceAlicloudSaeGreyTagRouteDelete(d, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
	})
	t.Run("DeleteMockAbnormal", func(t *testing.T) {
		retryFlag := true
		noRetryFlag := true
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&roa.Client{}), "DoRequest", func(_ *roa.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]*string, _ map[string]*string, _ interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if retryFlag {
				retryFlag = false
				// retry until the timeout comes
				return responseMock["RetryError"]("Throttling")
			} else if noRetryFlag {
				noRetryFlag = false
				return responseMock["NoRetryError"]("NonRetryableError")
			}
			return responseMock["DeleteNormal"]("")
		})
		err := resourceAlicloudSaeGreyTagRouteDelete(d, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
	})
	t.Run("DeleteMockNormal", func(t *testing.T) {
		retryFlag := false
		noRetryFlag := false
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&roa.Client{}), "DoRequest", func(_ *roa.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]*string, _ map[string]*string, _ interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if retryFlag {
				return responseMock["RetryError"]("Throttling")
			} else if noRetryFlag {
				return responseMock["NoRetryError"]("NonRetryableError")
			}
			return responseMock["DeleteNormal"]("")
		})
		err := resourceAlicloudSaeGreyTagRouteDelete(d, rawClient)
		patches.Reset()
		assert.Nil(t, err)
	})
	t.Run("DeleteNonRetryableError", func(t *testing.T) {
		retryFlag := false
		noRetryFlag := true
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&roa.Client{}), "DoRequest", func(_ *roa.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]*string, _ map[string]*string, _ interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if retryFlag {
				return responseMock["RetryError"]("Throttling")
			} else if noRetryFlag {
				noRetryFlag = false
				return responseMock["NoRetryError"]("NonRetryableError")
			}
			return responseMock["DeleteNormal"]("")
		})
		err := resourceAlicloudSaeGreyTagRouteDelete(d, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
	})

	//Read
	t.Run("ReadDescribeVpcSaeGreyTagRouteNotFound", func(t *testing.T) {
		patchRequest := gomonkey.ApplyMethod(reflect.TypeOf(&roa.Client{}), "DoRequest", func(_ *roa.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]*string, _ map[string]*string, _ interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			NotFoundFlag := true
			noRetryFlag := false
			if NotFoundFlag {
				return responseMock["NotFoundError"]("ResourceNotfound")
			} else if noRetryFlag {
				return responseMock["NoRetryError"]("NoRetryError")
			}
			return responseMock["ReadNormal"]("")
		})
		err := resourceAlicloudSaeGreyTagRouteRead(d, rawClient)
		patchRequest.Reset()
		assert.Nil(t, err)
	})
	t.Run("ReadDescribeSaeGreyTagRouteAbnormal", func(t *testing.T) {
		patcheDorequest := gomonkey.ApplyMethod(reflect.TypeOf(&roa.Client{}), "DoRequest", func(_ *roa.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]*string, _ map[string]*string, _ interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			return responseMock["ReadNormal"]("")
		})
		err := resourceAlicloudSaeGreyTagRouteRead(d, rawClient)
		patcheDorequest.Reset()
		assert.Nil(t, err)
	})
}

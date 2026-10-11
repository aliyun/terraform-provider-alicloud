package alicloud

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/alibabacloud-go/tea-rpc/client"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/stretchr/testify/assert"
)

func TestAccAliCloudOOSServiceSetting_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_oos_service_setting.default"
	checkoutSupportedRegions(t, true, connectivity.OOSSupportRegions)
	ra := resourceAttrInit(resourceId, AlicloudOOSServiceSettingMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &OosService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeOosServiceSetting")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%soos%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudOOSServiceSettingBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"delivery_oss_enabled":      "true",
					"delivery_oss_key_prefix":   "path1/",
					"delivery_oss_bucket_name":  "${alicloud_oss_bucket.default.0.bucket}",
					"delivery_sls_enabled":      "true",
					"delivery_sls_project_name": "${alicloud_log_project.default.0.name}",
					"service_access_rd_enabled": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"delivery_oss_enabled":      "true",
						"delivery_oss_key_prefix":   "path1/",
						"delivery_oss_bucket_name":  CHECKSET,
						"delivery_sls_enabled":      "true",
						"delivery_sls_project_name": CHECKSET,
						"service_access_rd_enabled": "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"service_access_rd_enabled": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"service_access_rd_enabled": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"delivery_oss_key_prefix": "path2/",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"delivery_oss_key_prefix": "path2/",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"delivery_oss_bucket_name": "${alicloud_oss_bucket.default.1.bucket}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"delivery_oss_bucket_name": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"delivery_oss_enabled": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"delivery_oss_enabled": "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"delivery_sls_project_name": "${alicloud_log_project.default.1.name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"delivery_sls_project_name": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"delivery_sls_enabled": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"delivery_sls_enabled": "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"delivery_oss_enabled":      "true",
					"delivery_oss_key_prefix":   "path1/",
					"delivery_oss_bucket_name":  "${alicloud_oss_bucket.default.0.bucket}",
					"delivery_sls_enabled":      "true",
					"delivery_sls_project_name": "${alicloud_log_project.default.0.name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"delivery_oss_enabled":      "true",
						"delivery_oss_key_prefix":   "path1/",
						"delivery_oss_bucket_name":  CHECKSET,
						"delivery_sls_enabled":      "true",
						"delivery_sls_project_name": CHECKSET,
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// rd_folder_ids is an ordered TypeList. TypeList Order Coverage requires two
// distinct valid members and three adjacent steps: apply order A, then reorder
// to B with PlanOnly+ExpectNonEmptyPlan (the SDK must report a diff), then apply
// B normally and verify a post-apply empty plan. Members are real resource
// directory folder IDs produced by alicloud_resource_manager_folder.
func TestAccAliCloudOOSServiceSetting_memberOrder(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_oos_service_setting.default"
	checkoutSupportedRegions(t, true, connectivity.OOSSupportRegions)
	ra := resourceAttrInit(resourceId, AlicloudOOSServiceSettingMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &OosService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeOosServiceSetting")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%soos%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudOOSServiceSettingMemberOrderDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					// The OOS SetServiceSettings API rejects RdFolderIds unless
					// ServiceAccessRdEnabled=true (RdServiceNotEnabled), so RD
					// folder management must enable RD service access first/in the
					// same request.
					"service_access_rd_enabled": "true",
					"rd_folder_ids": []string{
						"${alicloud_resource_manager_folder.first.id}",
						"${alicloud_resource_manager_folder.second.id}",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"rd_folder_ids.#": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"service_access_rd_enabled": "true",
					"rd_folder_ids": []string{
						"${alicloud_resource_manager_folder.second.id}",
						"${alicloud_resource_manager_folder.first.id}",
					},
				}),
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"service_access_rd_enabled": "true",
					"rd_folder_ids": []string{
						"${alicloud_resource_manager_folder.second.id}",
						"${alicloud_resource_manager_folder.first.id}",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"rd_folder_ids.#": "2",
					}),
				),
			},
		},
	})
}

func AlicloudOOSServiceSettingMemberOrderDependence0(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

resource "alicloud_resource_manager_folder" "first" {
  folder_name = join("", [var.name, "-first"])
}

resource "alicloud_resource_manager_folder" "second" {
  folder_name = join("", [var.name, "-second"])
}
`, name)
}

var AlicloudOOSServiceSettingMap0 = map[string]string{}

func AlicloudOOSServiceSettingBasicDependence0(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

resource "alicloud_oss_bucket" "default" {
  count  = 2
  bucket = join("", [var.name, count.index])
  acl    = "private"
}

resource "alicloud_log_project" "default" {
  count = 2
  name  = join("", [var.name, count.index])
}
`, name)
}

// lintignore: R001
func TestUnitAliCloudOOSServiceSetting(t *testing.T) {
	p := Provider().(*schema.Provider).ResourcesMap
	d, _ := schema.InternalMap(p["alicloud_oos_service_setting"].Schema).Data(nil, nil)
	dCreate, _ := schema.InternalMap(p["alicloud_oos_service_setting"].Schema).Data(nil, nil)
	dCreate.MarkNewResource()
	for key, value := range map[string]interface{}{
		"delivery_oss_enabled":      true,
		"delivery_oss_key_prefix":   "path1/",
		"delivery_oss_bucket_name":  "delivery_oss_bucket_name",
		"delivery_sls_enabled":      true,
		"delivery_sls_project_name": "delivery_sls_project_name",
		"service_access_rd_enabled": false,
		"rd_folder_ids":             []interface{}{"rd-xxx"},
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
		"ServiceSettings": map[string]interface{}{
			"DeliveryOssBucketName":  "delivery_oss_bucket_name",
			"DeliveryOssEnabled":     true,
			"DeliveryOssKeyPrefix":   "delivery_oss_key_prefix",
			"DeliverySlsEnabled":     true,
			"DeliverySlsProjectName": "delivery_sls_project_name",
			"RdcEnterpriseId":        "rdc_enterprise_id",
			"ServiceAccessRdEnabled": true,
			"RdFolders":              "rd-123",
		},
	}

	responseMock := map[string]func(errorCode string) (map[string]interface{}, error){
		"RetryError": func(errorCode string) (map[string]interface{}, error) {
			return nil, &tea.SDKError{
				Code:       String(errorCode),
				Data:       String(errorCode),
				Message:    String(errorCode),
				StatusCode: tea.Int(400),
			}
		},
		"NotFoundError": func(errorCode string) (map[string]interface{}, error) {
			return nil, GetNotFoundErrorFromString(GetNotFoundMessage("alicloud_oos_service_setting", fmt.Sprint(strings.Trim(uuid.New().String(), "-"))))
		},
		"NoRetryError": func(errorCode string) (map[string]interface{}, error) {
			return nil, &tea.SDKError{
				Code:       String(errorCode),
				Data:       String(errorCode),
				Message:    String(errorCode),
				StatusCode: tea.Int(400),
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
		// This resource dispatches via AliyunClient.RpcPost (there is no dedicated
		// NewOosClient builder), so the client/request abnormal path is simulated by
		// making RpcPost return a non-retryable error; Create then wraps and returns it
		// before the trailing Read, so no real network call is made.
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "RpcPost", func(_ *connectivity.AliyunClient, _ string, _ string, _ string, _ map[string]interface{}, _ map[string]interface{}, _ bool) (map[string]interface{}, error) {
			return nil, &tea.SDKError{
				Code:       String("loadEndpoint error"),
				Data:       String("loadEndpoint error"),
				Message:    String("loadEndpoint error"),
				StatusCode: tea.Int(400),
			}
		})
		err := resourceAlicloudOosServiceSettingCreate(d, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
	})
	t.Run("CreateAbnormal", func(t *testing.T) {
		retryFlag := true
		noRetryFlag := true
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if retryFlag {
				retryFlag = false
				return responseMock["RetryError"]("Throttling")
			} else if noRetryFlag {
				noRetryFlag = false
				return responseMock["NoRetryError"]("NonRetryableError")
			}
			return responseMock["CreateNormal"]("")
		})
		err := resourceAlicloudOosServiceSettingCreate(d, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
	})
	t.Run("CreateNormal", func(t *testing.T) {
		// SetServiceSettings is dispatched via RpcPost(nil, request): on the wire the
		// request is the DoRequest BODY (2nd map), not the query (1st map). The trailing
		// Read uses RpcGet(request, nil) so its query carries RegionId; the Set body is
		// distinguished by carrying ServiceAccessRdEnabled/RdFolderIds.
		var setBody map[string]interface{}
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, body map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if _, ok := body["ServiceAccessRdEnabled"]; ok {
				setBody = body
			}
			return responseMock["CreateNormal"]("")
		})
		err := resourceAlicloudOosServiceSettingCreate(dCreate, rawClient)
		patches.Reset()
		assert.Nil(t, err)
		assert.NotNil(t, setBody, "SetServiceSettings request body not captured")
		// service_access_rd_enabled=false must still be transmitted (GetOkExists on bool).
		assert.Equal(t, false, setBody["ServiceAccessRdEnabled"])
		assert.Equal(t, `["rd-xxx"]`, setBody["RdFolderIds"])
	})

	// Set ID for Update and Delete Method
	d.SetId(fmt.Sprint(strings.Trim(uuid.New().String(), "-")))
	// Update
	t.Run("UpdateClientAbnormal", func(t *testing.T) {
		// Update with no changes skips SetServiceSettings and falls through to Read,
		// which dispatches via AliyunClient.RpcGet (GetServiceSettings). Mocking RpcGet
		// to fail makes Read return an error, which Update propagates; no real network.
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "RpcGet", func(_ *connectivity.AliyunClient, _ string, _ string, _ string, _ map[string]interface{}, _ map[string]interface{}) (map[string]interface{}, error) {
			return nil, &tea.SDKError{
				Code:       String("loadEndpoint error"),
				Data:       String("loadEndpoint error"),
				Message:    String("loadEndpoint error"),
				StatusCode: tea.Int(400),
			}
		})

		err := resourceAlicloudOosServiceSettingUpdate(d, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
	})

	t.Run("UpdateSetServiceSettingsAbnormal", func(t *testing.T) {
		diff := terraform.NewInstanceDiff()
		for _, key := range []string{"delivery_oss_enabled", "delivery_oss_bucket_name", "delivery_oss_key_prefix", "delivery_sls_enabled", "delivery_sls_project_name", "service_access_rd_enabled"} {
			switch p["alicloud_oos_service_setting"].Schema[key].Type {
			case schema.TypeString:
				diff.SetAttribute(key, &terraform.ResourceAttrDiff{Old: d.Get(key).(string), New: d.Get(key).(string) + "_update"})
			case schema.TypeBool:
				diff.SetAttribute(key, &terraform.ResourceAttrDiff{Old: strconv.FormatBool(d.Get(key).(bool)), New: strconv.FormatBool(true)})
			case schema.TypeInt:
				diff.SetAttribute(key, &terraform.ResourceAttrDiff{Old: strconv.Itoa(d.Get(key).(int)), New: strconv.Itoa(3)})
			case schema.TypeMap:
				diff.SetAttribute("tags.%", &terraform.ResourceAttrDiff{Old: "0", New: "2"})
				diff.SetAttribute("tags.For", &terraform.ResourceAttrDiff{Old: "", New: "Test"})
				diff.SetAttribute("tags.Created", &terraform.ResourceAttrDiff{Old: "", New: "TF"})
			}
		}
		resourceData1, _ := schema.InternalMap(p["alicloud_oos_service_setting"].Schema).Data(nil, diff)
		resourceData1.SetId(d.Id())
		retryFlag := true
		noRetryFlag := true
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if retryFlag {
				retryFlag = false
				return responseMock["RetryError"]("Throttling")
			} else if noRetryFlag {
				noRetryFlag = false
				return responseMock["NoRetryError"]("NonRetryableError")
			}
			return responseMock["UpdateNormal"]("")
		})
		err := resourceAlicloudOosServiceSettingUpdate(resourceData1, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
	})

	t.Run("UpdateSetServiceSettingsNormal", func(t *testing.T) {
		diff := terraform.NewInstanceDiff()
		for _, key := range []string{"delivery_oss_enabled", "delivery_oss_bucket_name", "delivery_oss_key_prefix", "delivery_sls_enabled", "delivery_sls_project_name", "service_access_rd_enabled"} {
			switch p["alicloud_oos_service_setting"].Schema[key].Type {
			case schema.TypeString:
				diff.SetAttribute(key, &terraform.ResourceAttrDiff{Old: d.Get(key).(string), New: d.Get(key).(string) + "_update"})
			case schema.TypeBool:
				diff.SetAttribute(key, &terraform.ResourceAttrDiff{Old: strconv.FormatBool(d.Get(key).(bool)), New: strconv.FormatBool(true)})
			case schema.TypeInt:
				diff.SetAttribute(key, &terraform.ResourceAttrDiff{Old: strconv.Itoa(d.Get(key).(int)), New: strconv.Itoa(3)})
			case schema.TypeMap:
				diff.SetAttribute("tags.%", &terraform.ResourceAttrDiff{Old: "0", New: "2"})
				diff.SetAttribute("tags.For", &terraform.ResourceAttrDiff{Old: "", New: "Test"})
				diff.SetAttribute("tags.Created", &terraform.ResourceAttrDiff{Old: "", New: "TF"})
			}
		}
		resourceData1, _ := schema.InternalMap(p["alicloud_oos_service_setting"].Schema).Data(nil, diff)
		resourceData1.SetId(d.Id())
		var setBody map[string]interface{}
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, body map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if _, ok := body["ServiceAccessRdEnabled"]; ok {
				setBody = body
			}
			return responseMock["UpdateNormal"]("")
		})
		err := resourceAlicloudOosServiceSettingUpdate(resourceData1, rawClient)
		patches.Reset()
		assert.Nil(t, err)
		// The update request body must transmit the changed service_access_rd_enabled value.
		assert.NotNil(t, setBody)
		assert.Equal(t, true, setBody["ServiceAccessRdEnabled"])
	})

	// Update: rd_folder_ids must send an explicit empty JSON array on clear and the
	// populated JSON array when non-empty (SetServiceSettings is a set/replace API).
	t.Run("UpdateSetServiceSettingsClearRdFolderIds", func(t *testing.T) {
		diff := terraform.NewInstanceDiff()
		diff.SetAttribute("rd_folder_ids.#", &terraform.ResourceAttrDiff{Old: "1", New: "0"})
		resourceData1, _ := schema.InternalMap(p["alicloud_oos_service_setting"].Schema).Data(nil, diff)
		resourceData1.SetId(d.Id())
		var captured []map[string]interface{}
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, body map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			captured = append(captured, body)
			return responseMock["UpdateNormal"]("")
		})
		err := resourceAlicloudOosServiceSettingUpdate(resourceData1, rawClient)
		patches.Reset()
		assert.Nil(t, err)
		var setReq map[string]interface{}
		for _, r := range captured {
			if _, ok := r["RdFolderIds"]; ok {
				setReq = r
			}
		}
		assert.NotNil(t, setReq, "SetServiceSettings request carrying RdFolderIds not captured")
		assert.Equal(t, "[]", setReq["RdFolderIds"])
	})

	t.Run("UpdateSetServiceSettingsRdFolderIdsNonEmpty", func(t *testing.T) {
		diff := terraform.NewInstanceDiff()
		// A count change (0 -> 1) unambiguously triggers HasChange for the list, and the
		// Set below establishes the new single member so d.Get reads ["rd-new"].
		diff.SetAttribute("rd_folder_ids.#", &terraform.ResourceAttrDiff{Old: "0", New: "1"})
		resourceData1, _ := schema.InternalMap(p["alicloud_oos_service_setting"].Schema).Data(nil, diff)
		resourceData1.SetId(d.Id())
		assert.Nil(t, resourceData1.Set("rd_folder_ids", []interface{}{"rd-new"}))
		var captured []map[string]interface{}
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, body map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			captured = append(captured, body)
			return responseMock["UpdateNormal"]("")
		})
		err := resourceAlicloudOosServiceSettingUpdate(resourceData1, rawClient)
		patches.Reset()
		assert.Nil(t, err)
		var setReq map[string]interface{}
		for _, r := range captured {
			if _, ok := r["RdFolderIds"]; ok {
				setReq = r
			}
		}
		assert.NotNil(t, setReq, "SetServiceSettings request carrying RdFolderIds not captured")
		assert.Equal(t, `["rd-new"]`, setReq["RdFolderIds"])
	})

	// Delete
	t.Run("DeleteNormal", func(t *testing.T) {
		err := resourceAlicloudOosServiceSettingDelete(d, rawClient)
		assert.Nil(t, err)
	})

	//Read
	t.Run("ReadDescribeOosServiceSettingNotFound", func(t *testing.T) {
		patcheDorequest := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			NotFoundFlag := true
			noRetryFlag := false
			if NotFoundFlag {
				return responseMock["NotFoundError"]("ResourceNotfound")
			} else if noRetryFlag {
				return responseMock["NoRetryError"]("NoRetryError")
			}
			return responseMock["ReadNormal"]("")
		})
		err := resourceAlicloudOosServiceSettingRead(d, rawClient)
		patcheDorequest.Reset()
		assert.Nil(t, err)
	})

	t.Run("ReadDescribeOosServiceSettingAbnormal", func(t *testing.T) {
		patcheDorequest := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			retryFlag := false
			noRetryFlag := true
			if retryFlag {
				return responseMock["RetryError"]("Throttling")
			} else if noRetryFlag {
				return responseMock["NoRetryError"]("NonRetryableError")
			}
			return responseMock["ReadNormal"]("")
		})
		err := resourceAlicloudOosServiceSettingRead(d, rawClient)
		patcheDorequest.Reset()
		assert.NotNil(t, err)
	})

	// Read backfill: rd_folder_ids is read from the asymmetric RdFolders string field.
	t.Run("ReadRdFoldersJsonArray", func(t *testing.T) {
		mock := map[string]interface{}{
			"ServiceSettings": map[string]interface{}{
				"RdFolders": `["rd-a","rd-b"]`,
			},
		}
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			return mock, nil
		})
		err := resourceAlicloudOosServiceSettingRead(d, rawClient)
		patches.Reset()
		assert.Nil(t, err)
		assert.Equal(t, []interface{}{"rd-a", "rd-b"}, d.Get("rd_folder_ids"))
	})

	t.Run("ReadRdFoldersEmpty", func(t *testing.T) {
		mock := map[string]interface{}{
			"ServiceSettings": map[string]interface{}{
				"RdFolders": "",
			},
		}
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			return mock, nil
		})
		err := resourceAlicloudOosServiceSettingRead(d, rawClient)
		patches.Reset()
		assert.Nil(t, err)
		assert.Equal(t, []interface{}{}, d.Get("rd_folder_ids"))
	})
}

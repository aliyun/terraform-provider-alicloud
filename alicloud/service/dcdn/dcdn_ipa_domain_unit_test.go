package dcdn

import (
	"os"
	"reflect"
	"strconv"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/alibabacloud-go/tea-rpc/client"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/stretchr/testify/assert"
)

// lintignore: R001
func TestUnitAlicloudDCDNIpaDomain(t *testing.T) {
	p := unitTestProviderResources()
	d, _ := schema.InternalMap(p["alicloud_dcdn_ipa_domain"].Schema).Data(nil, nil)
	dCreate, _ := schema.InternalMap(p["alicloud_dcdn_ipa_domain"].Schema).Data(nil, nil)
	dCreate.MarkNewResource()
	for key, value := range map[string]interface{}{
		"domain_name": "domain_name",
		"sources": []map[string]interface{}{
			{
				"content":  "content",
				"port":     8888,
				"priority": "priority",
				"type":     "type",
				"weight":   10,
			},
		},
		"scope":             "scope",
		"resource_group_id": "resource_group_id",
	} {
		err := dCreate.Set(key, value)
		assert.Nil(t, err)
		err = d.Set(key, value)
		assert.Nil(t, err)
	}
	region := os.Getenv("ALICLOUD_REGION")
	rawClient, err := unitTestSharedClientForRegion(region)
	if err != nil {
		t.Skipf("Skipping the test case with err: %s", err)
		t.Skipped()
	}
	rawClient = rawClient.(*connectivity.AliyunClient)
	ReadMockResponse := map[string]interface{}{
		"DomainDetail": map[string]interface{}{
			"Cname":           "domain_name",
			"ResourceGroupId": "resource_group_id",
			"DomainStatus":    "online",
			"DomainName":      "domain_name",
			"Sources": map[string]interface{}{
				"Source": []interface{}{
					map[string]interface{}{
						"Enabled":  "online",
						"Port":     8888,
						"Type":     "type",
						"Content":  "content",
						"Priority": "priority",
						"Weight":   "10",
					},
				},
			},
			"GmtModified": "gmtModified",
			"GmtCreated":  "gmtCreated",
		},
	}

	responseMock := map[string]func(errorCode string) (map[string]interface{}, error){
		"RetryError": func(errorCode string) (map[string]interface{}, error) {
			return nil, &tea.SDKError{
				Code:       helper.String(errorCode),
				Data:       helper.String(errorCode),
				Message:    helper.String(errorCode),
				StatusCode: tea.Int(400),
			}
		},
		"errs.NotFoundError": func(errorCode string) (map[string]interface{}, error) {
			return nil, errs.GetNotFoundErrorFromString(errs.GetNotFoundMessage("alicloud_dcdn_ipa_domain", "xxx_id"))
		},
		"NoRetryError": func(errorCode string) (map[string]interface{}, error) {
			return nil, &tea.SDKError{
				Code:       helper.String(errorCode),
				Data:       helper.String(errorCode),
				Message:    helper.String(errorCode),
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
		"UpdateOfflineStatusNormal": func(errorCode string) (map[string]interface{}, error) {
			result := map[string]interface{}{
				"DomainDetail": map[string]interface{}{
					"Cname":           "domain_name",
					"ResourceGroupId": "resource_group_id",
					"DomainStatus":    "offline",
					"DomainName":      "domain_name",
					"Sources": map[string]interface{}{
						"Source": []interface{}{
							map[string]interface{}{
								"Enabled":  "offline",
								"Port":     8888,
								"Type":     "type",
								"Content":  "content",
								"Priority": "priority",
								"Weight":   "10",
							},
						},
					},
					"GmtModified": "gmtModified",
					"GmtCreated":  "gmtCreated",
				},
			}
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
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewDcdnClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
			return nil, &tea.SDKError{
				Code:       helper.String("loadEndpoint error"),
				Data:       helper.String("loadEndpoint error"),
				Message:    helper.String("loadEndpoint error"),
				StatusCode: tea.Int(400),
			}
		})
		err := resourceAlicloudDcdnIpaDomainCreate(d, rawClient)
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
		err := resourceAlicloudDcdnIpaDomainCreate(d, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
	})
	t.Run("CreateNormal", func(t *testing.T) {
		retryFlag := false
		noRetryFlag := false
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
		err := resourceAlicloudDcdnIpaDomainCreate(dCreate, rawClient)
		patches.Reset()
		assert.Nil(t, err)
	})
	t.Run("CreateRetryableError", func(t *testing.T) {
		retryFlag := false
		noRetryFlag := true
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if retryFlag {
				return responseMock["RetryError"]("Throttling")
			} else if noRetryFlag {
				noRetryFlag = false
				return responseMock["NoRetryError"]("NonRetryableError")
			}
			return responseMock["CreateNormal"]("")
		})
		err := resourceAlicloudDcdnIpaDomainCreate(d, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
	})
	t.Run("CreateNoRetryErrorError", func(t *testing.T) {
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			return responseMock["CreateNormal"]("")
		})
		patchDescribe := gomonkey.ApplyMethod(reflect.TypeOf(&DcdnService{}), "DescribeDcdnIpaDomain", func(*DcdnService, string) (map[string]interface{}, error) {
			return responseMock["NoRetryError"]("NoRetryError")
		})
		err := resourceAlicloudDcdnIpaDomainCreate(d, rawClient)
		patches.Reset()
		patchDescribe.Reset()
		assert.NotNil(t, err)
	})

	// Set ID for Update and Delete Method
	d.SetId("domain_name")

	// Update
	t.Run("UpdateClientAbnormal", func(t *testing.T) {
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewDcdnClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
			return nil, &tea.SDKError{
				Code:       helper.String("loadEndpoint error"),
				Data:       helper.String("loadEndpoint error"),
				Message:    helper.String("loadEndpoint error"),
				StatusCode: tea.Int(400),
			}
		})

		err := resourceAlicloudDcdnIpaDomainUpdate(d, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
	})
	t.Run("UpdateModifyAttributeAbnormal", func(t *testing.T) {
		diff := terraform.NewInstanceDiff()
		for _, key := range []string{"resource_group_id"} {
			switch p["alicloud_dcdn_ipa_domain"].Schema[key].Type {
			case schema.TypeString:
				diff.Attributes[key] = &terraform.ResourceAttrDiff{Old: d.Get(key).(string), New: d.Get(key).(string) + "_update"}
			case schema.TypeBool:
				diff.Attributes[key] = &terraform.ResourceAttrDiff{Old: strconv.FormatBool(d.Get(key).(bool)), New: strconv.FormatBool(true)}
			case schema.TypeInt:
				diff.Attributes[key] = &terraform.ResourceAttrDiff{Old: strconv.Itoa(d.Get(key).(int)), New: strconv.Itoa(1200)}
			case schema.TypeMap:
				diff.Attributes["tags.%"] = &terraform.ResourceAttrDiff{Old: "0", New: "2"}
				diff.Attributes["tags.For"] = &terraform.ResourceAttrDiff{Old: "", New: "Test"}
				diff.Attributes["tags.Created"] = &terraform.ResourceAttrDiff{Old: "", New: "TF"}
			}
		}
		resourceData1, _ := schema.InternalMap(p["alicloud_dcdn_ipa_domain"].Schema).Data(nil, diff)
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
		err := resourceAlicloudDcdnIpaDomainUpdate(resourceData1, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
	})
	t.Run("UpdateModifyAttributeNoRetryErrorAbnormal", func(t *testing.T) {
		diff := terraform.NewInstanceDiff()
		for _, key := range []string{"resource_group_id"} {
			switch p["alicloud_dcdn_ipa_domain"].Schema[key].Type {
			case schema.TypeString:
				diff.Attributes[key] = &terraform.ResourceAttrDiff{Old: d.Get(key).(string), New: d.Get(key).(string) + "_update"}
			case schema.TypeBool:
				diff.Attributes[key] = &terraform.ResourceAttrDiff{Old: strconv.FormatBool(d.Get(key).(bool)), New: strconv.FormatBool(true)}
			case schema.TypeInt:
				diff.Attributes[key] = &terraform.ResourceAttrDiff{Old: strconv.Itoa(d.Get(key).(int)), New: strconv.Itoa(1200)}
			case schema.TypeMap:
				diff.Attributes["tags.%"] = &terraform.ResourceAttrDiff{Old: "0", New: "2"}
				diff.Attributes["tags.For"] = &terraform.ResourceAttrDiff{Old: "", New: "Test"}
				diff.Attributes["tags.Created"] = &terraform.ResourceAttrDiff{Old: "", New: "TF"}
			}
		}
		resourceData, _ := schema.InternalMap(p["alicloud_dcdn_ipa_domain"].Schema).Data(nil, diff)
		resourceData.SetId(d.Id())
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			return responseMock["UpdateNormal"]("")
		})
		patcheDescribe := gomonkey.ApplyMethod(reflect.TypeOf(&DcdnService{}), "DescribeDcdnIpaDomain", func(*DcdnService, string) (map[string]interface{}, error) {
			return responseMock["NoRetryError"]("NoRetryError")
		})
		err := resourceAlicloudDcdnIpaDomainUpdate(resourceData, rawClient)
		patches.Reset()
		patcheDescribe.Reset()
		assert.NotNil(t, err)
	})
	t.Run("UpdateModifyAttributeNormal", func(t *testing.T) {
		diff := terraform.NewInstanceDiff()
		for _, key := range []string{"resource_group_id", "sources"} {
			switch p["alicloud_dcdn_ipa_domain"].Schema[key].Type {
			case schema.TypeString:
				diff.Attributes[key] = &terraform.ResourceAttrDiff{Old: d.Get(key).(string), New: d.Get(key).(string) + "_update"}
			case schema.TypeBool:
				diff.Attributes[key] = &terraform.ResourceAttrDiff{Old: strconv.FormatBool(d.Get(key).(bool)), New: strconv.FormatBool(true)}
			case schema.TypeInt:
				diff.Attributes[key] = &terraform.ResourceAttrDiff{Old: strconv.Itoa(d.Get(key).(int)), New: strconv.Itoa(1200)}
			case schema.TypeSet:
				diff.Attributes["sources.0.content"] = &terraform.ResourceAttrDiff{Old: "content", New: "content_update"}
				diff.Attributes["sources.0.port"] = &terraform.ResourceAttrDiff{Old: "8888", New: "8888"}
				diff.Attributes["sources.0.priority"] = &terraform.ResourceAttrDiff{Old: "priority", New: "priority"}
				diff.Attributes["sources.0.type"] = &terraform.ResourceAttrDiff{Old: "type", New: "type"}
				diff.Attributes["sources.0.weight"] = &terraform.ResourceAttrDiff{Old: "10", New: "9"}
			}
		}
		resourceData1, _ := schema.InternalMap(p["alicloud_dcdn_ipa_domain"].Schema).Data(nil, diff)
		resourceData1.SetId(d.Id())
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			return responseMock["UpdateNormal"]("")
		})
		err := resourceAlicloudDcdnIpaDomainUpdate(resourceData1, rawClient)
		patches.Reset()
		assert.Nil(t, err)
	})
	t.Run("UpdateModifyStatusOfflineAttributeNotFoundErrorAbnormal", func(t *testing.T) {
		diff := terraform.NewInstanceDiff()
		diff.Attributes["status"] = &terraform.ResourceAttrDiff{Old: "online", New: "offline"}
		resourceData1, _ := schema.InternalMap(p["alicloud_dcdn_ipa_domain"].Schema).Data(nil, diff)
		resourceData1.SetId(d.Id())
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			return responseMock["errs.NotFoundError"]("ResourceNotfound")
		})
		err := resourceAlicloudDcdnIpaDomainUpdate(resourceData1, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
	})
	t.Run("UpdateModifyStatusOfflineAttributeAbnormal", func(t *testing.T) {
		diff := terraform.NewInstanceDiff()
		diff.Attributes["status"] = &terraform.ResourceAttrDiff{Old: "online", New: "offline"}
		resourceData1, _ := schema.InternalMap(p["alicloud_dcdn_ipa_domain"].Schema).Data(nil, diff)
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
		patchDescribe := gomonkey.ApplyMethod(reflect.TypeOf(&DcdnService{}), "DescribeDcdnIpaDomain", func(*DcdnService, string) (map[string]interface{}, error) {
			object := map[string]interface{}{
				"Cname":           "domain_name",
				"ResourceGroupId": "resource_group_id",
				"DomainStatus":    "online",
				"DomainName":      "domain_name",
				"Sources": map[string]interface{}{
					"Source": []interface{}{
						map[string]interface{}{
							"Enabled":  "online",
							"Port":     8888,
							"Type":     "type",
							"Content":  "content",
							"Priority": "priority",
							"Weight":   "10",
						},
					},
				},
				"GmtModified": "gmtModified",
				"GmtCreated":  "gmtCreated",
			}
			return object, nil
		})
		err := resourceAlicloudDcdnIpaDomainUpdate(resourceData1, rawClient)
		patches.Reset()
		patchDescribe.Reset()
		assert.NotNil(t, err)
	})
	t.Run("UpdateModifyStatusOnlineAttributeAbnormal", func(t *testing.T) {
		diff := terraform.NewInstanceDiff()
		diff.Attributes["status"] = &terraform.ResourceAttrDiff{Old: "offline", New: "online"}

		resourceData1, _ := schema.InternalMap(p["alicloud_dcdn_ipa_domain"].Schema).Data(nil, diff)
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
		patchDescribe := gomonkey.ApplyMethod(reflect.TypeOf(&DcdnService{}), "DescribeDcdnIpaDomain", func(*DcdnService, string) (map[string]interface{}, error) {
			object := map[string]interface{}{
				"Cname":           "domain_name",
				"ResourceGroupId": "resource_group_id",
				"DomainStatus":    "offline",
				"DomainName":      "domain_name",
				"Sources": map[string]interface{}{
					"Source": []interface{}{
						map[string]interface{}{
							"Enabled":  "offline",
							"Port":     8888,
							"Type":     "type",
							"Content":  "content",
							"Priority": "priority",
							"Weight":   "10",
						},
					},
				},
				"GmtModified": "gmtModified",
				"GmtCreated":  "gmtCreated",
			}
			return object, nil
		})
		err := resourceAlicloudDcdnIpaDomainUpdate(resourceData1, rawClient)
		patches.Reset()
		patchDescribe.Reset()
		assert.NotNil(t, err)
	})
	t.Run("UpdateModifyStatusOfflineAttributeNormal", func(t *testing.T) {
		diff := terraform.NewInstanceDiff()
		diff.Attributes["status"] = &terraform.ResourceAttrDiff{Old: "online", New: "offline"}
		resourceData1, _ := schema.InternalMap(p["alicloud_dcdn_ipa_domain"].Schema).Data(nil, diff)
		resourceData1.SetId(d.Id())
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			return responseMock["UpdateNormal"]("")
		})
		patchDescribe := gomonkey.ApplyMethod(reflect.TypeOf(&DcdnService{}), "DcdnIpaDomainStateRefreshFunc", func(*DcdnService, string, []string) retry.StateRefreshFunc {
			return func() (interface{}, string, error) {
				object := map[string]interface{}{
					"Cname":           "domain_name",
					"ResourceGroupId": "resource_group_id",
					"DomainStatus":    "offline",
					"DomainName":      "domain_name",
					"Sources": map[string]interface{}{
						"Source": []interface{}{
							map[string]interface{}{
								"Enabled":  "offline",
								"Port":     8888,
								"Type":     "type",
								"Content":  "content",
								"Priority": "priority",
								"Weight":   "10",
							},
						},
					},
					"GmtModified": "gmtModified",
					"GmtCreated":  "gmtCreated",
				}
				return object, "offline", nil
			}
		})
		err := resourceAlicloudDcdnIpaDomainUpdate(resourceData1, rawClient)
		patches.Reset()
		patchDescribe.Reset()
		assert.Nil(t, err)
	})
	t.Run("UpdateModifyStatusOnlineAttributeNormal", func(t *testing.T) {
		diff := terraform.NewInstanceDiff()
		diff.Attributes["status"] = &terraform.ResourceAttrDiff{Old: "offline", New: "online"}
		resourceData1, _ := schema.InternalMap(p["alicloud_dcdn_ipa_domain"].Schema).Data(nil, diff)
		resourceData1.SetId(d.Id())
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			return responseMock["UpdateOfflineStatusNormal"]("")
		})
		patchDescribe := gomonkey.ApplyMethod(reflect.TypeOf(&DcdnService{}), "DcdnIpaDomainStateRefreshFunc", func(*DcdnService, string, []string) retry.StateRefreshFunc {
			return func() (interface{}, string, error) {
				object := map[string]interface{}{
					"Cname":           "domain_name",
					"ResourceGroupId": "resource_group_id",
					"DomainStatus":    "online",
					"DomainName":      "domain_name",
					"Sources": map[string]interface{}{
						"Source": []interface{}{
							map[string]interface{}{
								"Enabled":  "online",
								"Port":     8888,
								"Type":     "type",
								"Content":  "content",
								"Priority": "priority",
								"Weight":   "10",
							},
						},
					},
					"GmtModified": "gmtModified",
					"GmtCreated":  "gmtCreated",
				}
				return object, "online", nil
			}
		})
		err := resourceAlicloudDcdnIpaDomainUpdate(resourceData1, rawClient)
		patches.Reset()
		patchDescribe.Reset()
		assert.Nil(t, err)
	})

	// Delete
	t.Run("DeleteClientAbnormal", func(t *testing.T) {
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewDcdnClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
			return nil, &tea.SDKError{
				Code:       helper.String("loadEndpoint error"),
				Data:       helper.String("loadEndpoint error"),
				Message:    helper.String("loadEndpoint error"),
				StatusCode: tea.Int(400),
			}
		})
		err := resourceAlicloudDcdnIpaDomainDelete(d, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
	})
	t.Run("DeleteMockAbnormal", func(t *testing.T) {
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
			return responseMock["DeleteNormal"]("")
		})
		err := resourceAlicloudDcdnIpaDomainDelete(d, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
	})
	t.Run("DeleteMockNormal", func(t *testing.T) {
		retryFlag := false
		noRetryFlag := false
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if retryFlag {
				return responseMock["RetryError"]("Throttling")
			} else if noRetryFlag {
				return responseMock["NoRetryError"]("NonRetryableError")
			}
			return responseMock["DeleteNormal"]("")
		})
		patchDescribe := gomonkey.ApplyMethod(reflect.TypeOf(&DcdnService{}), "DescribeDcdnIpaDomain", func(*DcdnService, string) (map[string]interface{}, error) {
			return responseMock["NoRetryError"]("NoRetryError")
		})
		err := resourceAlicloudDcdnIpaDomainDelete(d, rawClient)
		patches.Reset()
		patchDescribe.Reset()
		assert.NotNil(t, err)
	})
	t.Run("DeleteNonRetryableError", func(t *testing.T) {
		retryFlag := false
		noRetryFlag := true
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if retryFlag {
				return responseMock["RetryError"]("Throttling")
			} else if noRetryFlag {
				noRetryFlag = false
				return responseMock["NoRetryError"]("NonRetryableError")
			}
			return responseMock["DeleteNormal"]("")
		})
		err := resourceAlicloudDcdnIpaDomainDelete(d, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
	})

	//Read
	t.Run("ReadDescribeNotFound", func(t *testing.T) {
		patchRequest := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			NotFoundFlag := true
			noRetryFlag := false
			if NotFoundFlag {
				return responseMock["errs.NotFoundError"]("ResourceNotfound")
			} else if noRetryFlag {
				return responseMock["NoRetryError"]("NoRetryError")
			}
			return responseMock["ReadNormal"]("")
		})
		err := resourceAlicloudDcdnIpaDomainRead(d, rawClient)
		patchRequest.Reset()
		assert.Nil(t, err)
	})
	t.Run("ReadDescribeAbnormal", func(t *testing.T) {
		patcheDorequest := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			return responseMock["ReadNormal"]("")
		})
		err := resourceAlicloudDcdnIpaDomainRead(d, rawClient)
		patcheDorequest.Reset()
		assert.Nil(t, err)
	})
}

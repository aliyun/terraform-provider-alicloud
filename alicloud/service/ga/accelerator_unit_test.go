package ga

import (
	"context"
	"fmt"
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
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/stretchr/testify/assert"
	"github.com/zclconf/go-cty/cty"
)

// White-box (internal package ga) TestUnit cases moved out of the
// external ga_test file so they can call the unexported resource CRUD
// funcs. Shared local helpers (unitTestProviderResources /
// newInstanceDiff / unitTestSharedClientForRegion) live in
// unit_test_helpers_test.go; see that file for why acctest cannot be
// imported here.

// lintignore: R001
func TestUnitAliCloudGaAccelerator(t *testing.T) {
	p := unitTestProviderResources()
	d, _ := schema.InternalMap(p["alicloud_ga_accelerator"].Schema).Data(nil, nil)
	dCreate, _ := schema.InternalMap(p["alicloud_ga_accelerator"].Schema).Data(nil, nil)
	dCreate.MarkNewResource()
	for key, value := range map[string]interface{}{
		"auto_use_coupon":     true,
		"duration":            9,
		"spec":                "1",
		"accelerator_name":    "accelerator_name",
		"auto_renew_duration": 1,
		"description":         "description",
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
		//DescribeGaAccelerator
		"Name":        "accelerator_name",
		"State":       "active",
		"Description": "description",
		"Spec":        "1",

		//DescribeAcceleratorAutoRenewAttribute
		"AutoRenewDuration": 1,
		"RenewalStatus":     false,
	}

	responseMock := map[string]func(errorCode string) (map[string]interface{}, error){
		"RetryError": func(errorCode string) (map[string]interface{}, error) {
			return nil, &tea.SDKError{
				Code:       tea.String(errorCode),
				Data:       tea.String(errorCode),
				Message:    tea.String(errorCode),
				StatusCode: tea.Int(400),
			}
		},
		"NotFoundError": func(errorCode string) (map[string]interface{}, error) {
			return nil, errs.GetNotFoundErrorFromString(errs.GetNotFoundMessage("alicloud_ga_accelerator", "MockId"))
		},
		"NoRetryError": func(errorCode string) (map[string]interface{}, error) {
			return nil, &tea.SDKError{
				Code:       tea.String(errorCode),
				Data:       tea.String(errorCode),
				Message:    tea.String(errorCode),
				StatusCode: tea.Int(400),
			}
		},
		"CreateNormal": func(errorCode string) (map[string]interface{}, error) {
			result := ReadMockResponse
			result["AcceleratorId"] = "MockId"
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
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewGaplusClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
			return nil, &tea.SDKError{
				Code:       tea.String("loadEndpoint error"),
				Data:       tea.String("loadEndpoint error"),
				Message:    tea.String("loadEndpoint error"),
				StatusCode: tea.Int(400),
			}
		})
		err := resourceAliCloudGaAcceleratorCreate(d, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
	})
	t.Run("CreateAbnormal", func(t *testing.T) {
		d.Set("duration", 24)
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
			return responseMock["Normal"]("")
		})
		err := resourceAliCloudGaAcceleratorCreate(d, rawClient)
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
		err := resourceAliCloudGaAcceleratorCreate(dCreate, rawClient)
		patches.Reset()
		assert.Nil(t, err)
	})

	// Set ID for Update and Delete Method
	d.SetId("MockId")
	// Update
	t.Run("UpdateClientAbnormal", func(t *testing.T) {
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewGaplusClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
			return nil, &tea.SDKError{
				Code:       tea.String("loadEndpoint error"),
				Data:       tea.String("loadEndpoint error"),
				Message:    tea.String("loadEndpoint error"),
				StatusCode: tea.Int(400),
			}
		})
		err := resourceAliCloudGaAcceleratorUpdate(d, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
	})
	t.Run("UpdateAcceleratorAutoRenewAttributeAbnormal", func(t *testing.T) {
		diff := terraform.NewInstanceDiff()
		for _, key := range []string{"auto_renew_duration", "renewal_status"} {
			switch p["alicloud_ga_accelerator"].Schema[key].Type {
			case schema.TypeString:
				diff.Attributes[key] = &terraform.ResourceAttrDiff{Old: d.Get(key).(string), New: d.Get(key).(string) + "_update"}
			case schema.TypeBool:
				diff.Attributes[key] = &terraform.ResourceAttrDiff{Old: strconv.FormatBool(d.Get(key).(bool)), New: strconv.FormatBool(true)}
			case schema.TypeInt:
				diff.Attributes[key] = &terraform.ResourceAttrDiff{Old: strconv.Itoa(d.Get(key).(int)), New: strconv.Itoa(3)}
			case schema.TypeMap:
				diff.Attributes["tags.%"] = &terraform.ResourceAttrDiff{Old: "0", New: "2"}
				diff.Attributes["tags.For"] = &terraform.ResourceAttrDiff{Old: "", New: "Test"}
				diff.Attributes["tags.Created"] = &terraform.ResourceAttrDiff{Old: "", New: "TF"}
			}
		}
		resourceData1, _ := schema.InternalMap(p["alicloud_ga_accelerator"].Schema).Data(nil, diff)
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
			return responseMock["Normal"]("")
		})
		err := resourceAliCloudGaAcceleratorUpdate(resourceData1, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
	})
	t.Run("UpdateAcceleratorAutoRenewAttributeNormal", func(t *testing.T) {
		diff := terraform.NewInstanceDiff()
		for _, key := range []string{"auto_renew_duration", "renewal_status"} {
			switch p["alicloud_ga_accelerator"].Schema[key].Type {
			case schema.TypeString:
				diff.Attributes[key] = &terraform.ResourceAttrDiff{Old: d.Get(key).(string), New: d.Get(key).(string) + "_update"}
			case schema.TypeBool:
				diff.Attributes[key] = &terraform.ResourceAttrDiff{Old: strconv.FormatBool(d.Get(key).(bool)), New: strconv.FormatBool(true)}
			case schema.TypeInt:
				diff.Attributes[key] = &terraform.ResourceAttrDiff{Old: strconv.Itoa(d.Get(key).(int)), New: strconv.Itoa(3)}
			case schema.TypeMap:
				diff.Attributes["tags.%"] = &terraform.ResourceAttrDiff{Old: "0", New: "2"}
				diff.Attributes["tags.For"] = &terraform.ResourceAttrDiff{Old: "", New: "Test"}
				diff.Attributes["tags.Created"] = &terraform.ResourceAttrDiff{Old: "", New: "TF"}
			}
		}
		resourceData1, _ := schema.InternalMap(p["alicloud_ga_accelerator"].Schema).Data(nil, diff)
		resourceData1.SetId(d.Id())
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
			return responseMock["UpdateNormal"]("")
		})
		err := resourceAliCloudGaAcceleratorUpdate(resourceData1, rawClient)
		patches.Reset()
		assert.Nil(t, err)
	})
	t.Run("UpdateAcceleratorAbnormal", func(t *testing.T) {
		diff := terraform.NewInstanceDiff()
		for _, key := range []string{"accelerator_name", "description", "spec", "auto_use_coupon"} {
			switch p["alicloud_ga_accelerator"].Schema[key].Type {
			case schema.TypeString:
				diff.Attributes[key] = &terraform.ResourceAttrDiff{Old: d.Get(key).(string), New: d.Get(key).(string) + "_update"}
			case schema.TypeBool:
				diff.Attributes[key] = &terraform.ResourceAttrDiff{Old: strconv.FormatBool(d.Get(key).(bool)), New: strconv.FormatBool(true)}
			case schema.TypeInt:
				diff.Attributes[key] = &terraform.ResourceAttrDiff{Old: strconv.Itoa(d.Get(key).(int)), New: strconv.Itoa(3)}
			case schema.TypeMap:
				diff.Attributes["tags.%"] = &terraform.ResourceAttrDiff{Old: "0", New: "2"}
				diff.Attributes["tags.For"] = &terraform.ResourceAttrDiff{Old: "", New: "Test"}
				diff.Attributes["tags.Created"] = &terraform.ResourceAttrDiff{Old: "", New: "TF"}
			}
		}
		resourceData1, _ := schema.InternalMap(p["alicloud_ga_accelerator"].Schema).Data(nil, diff)
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
			return responseMock["Normal"]("")
		})
		err := resourceAliCloudGaAcceleratorUpdate(resourceData1, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
	})
	t.Run("UpdateAcceleratorNormal", func(t *testing.T) {
		diff := terraform.NewInstanceDiff()
		for _, key := range []string{"accelerator_name", "description", "spec", "auto_use_coupon"} {
			switch p["alicloud_ga_accelerator"].Schema[key].Type {
			case schema.TypeString:
				diff.Attributes[key] = &terraform.ResourceAttrDiff{Old: d.Get(key).(string), New: d.Get(key).(string) + "_update"}
			case schema.TypeBool:
				diff.Attributes[key] = &terraform.ResourceAttrDiff{Old: strconv.FormatBool(d.Get(key).(bool)), New: strconv.FormatBool(true)}
			case schema.TypeInt:
				diff.Attributes[key] = &terraform.ResourceAttrDiff{Old: strconv.Itoa(d.Get(key).(int)), New: strconv.Itoa(3)}
			case schema.TypeMap:
				diff.Attributes["tags.%"] = &terraform.ResourceAttrDiff{Old: "0", New: "2"}
				diff.Attributes["tags.For"] = &terraform.ResourceAttrDiff{Old: "", New: "Test"}
				diff.Attributes["tags.Created"] = &terraform.ResourceAttrDiff{Old: "", New: "TF"}
			}
		}
		resourceData1, _ := schema.InternalMap(p["alicloud_ga_accelerator"].Schema).Data(nil, diff)
		resourceData1.SetId(d.Id())
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
			return responseMock["UpdateNormal"]("")
		})
		err := resourceAliCloudGaAcceleratorUpdate(resourceData1, rawClient)
		patches.Reset()
		assert.Nil(t, err)
	})

	// Delete
	t.Run("DeleteMockNormal", func(t *testing.T) {
		retryFlag := false
		noRetryFlag := false
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if retryFlag {
				retryFlag = false
				return responseMock["RetryError"]("RetryError")
			} else if noRetryFlag {
				noRetryFlag = false
				return responseMock["NoRetryError"]("NonRetryableError")
			}
			return responseMock["DeleteNormal"]("")
		})
		err := resourceAliCloudGaAcceleratorDelete(d, rawClient)
		patches.Reset()
		assert.Nil(t, err)
	})

	//Read
	t.Run("ReadDescribeGaAcceleratorNotFound", func(t *testing.T) {
		patchDoRequest := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			NotFoundFlag := true
			noRetryFlag := false
			if NotFoundFlag {
				return responseMock["NotFoundError"]("ResourceNotfound")
			} else if noRetryFlag {
				return responseMock["NoRetryError"]("NoRetryError")
			}
			return responseMock["ReadNormal"]("")
		})
		err := resourceAliCloudGaAcceleratorRead(d, rawClient)
		patchDoRequest.Reset()
		assert.Nil(t, err)
	})
	t.Run("ReadDescribeCrChartNamespaceUnexpectedError", func(t *testing.T) {
		patchDoRequest := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			NotFoundFlag := false
			noRetryFlag := true
			if NotFoundFlag {
				return responseMock["NotFoundError"]("CHART_REPO_NOT_EXIST")
			} else if noRetryFlag {
				return responseMock["NoRetryError"]("NoRetryError")
			}
			return responseMock["ReadNormal"]("")
		})
		err := resourceAliCloudGaAcceleratorRead(d, rawClient)
		patchDoRequest.Reset()
		assert.NotNil(t, err)
	})
}

// lintignore: R001
func TestUnitAliCloudGaAcceleratorDdosConfigs(t *testing.T) {
	add, remove := ddosConfigDiff(
		[]map[string]string{{"ddos_id": "ddos-a", "ddos_region_id": "cn-hangzhou"}},
		[]map[string]string{{"ddos_id": "ddos-b", "ddos_region_id": "ap-southeast-1"}},
	)
	assert.Len(t, add, 1)
	assert.Equal(t, "ddos-b", add[0]["ddos_id"])
	assert.Len(t, remove, 1)
	assert.Equal(t, "ddos-a", remove[0]["ddos_id"])

	request := buildDdosConfigRequest(map[string]interface{}{}, add)
	assert.Equal(t, "ddos-b", request["DdosConfigList.1.DdosId"])
	assert.Equal(t, "ap-southeast-1", request["DdosConfigList.1.DdosRegionId"])
	request = buildDdosConfigRequest(map[string]interface{}{}, remove)
	assert.Equal(t, "ddos-a", request["DdosConfigList.1.DdosId"])
	assert.Equal(t, "cn-hangzhou", request["DdosConfigList.1.DdosRegionId"])
}

// lintignore: R001
func TestUnitAliCloudGaAcceleratorDdosUpdateRequest(t *testing.T) {
	region := os.Getenv("ALICLOUD_REGION")
	rawClient, err := unitTestSharedClientForRegion(region)
	if err != nil {
		t.Skipf("Skipping the test case with err: %s", err)
	}

	state := &terraform.InstanceState{ID: "MockId", Attributes: map[string]string{"ddos_config_list.#": "0"}}
	diff, err := newInstanceDiff("alicloud_ga_accelerator", map[string]interface{}{}, map[string]interface{}{
		"spec":     "2",
		"auto_pay": false,
		"ddos_config_list": []interface{}{
			map[string]interface{}{
				"ddos_id":        "ddoscoo-cn-test",
				"ddos_region_id": "cn-hangzhou",
			},
		},
	}, state)
	assert.Nil(t, err)

	resourceData, err := schema.InternalMap(unitTestProviderResources()["alicloud_ga_accelerator"].Schema).Data(state, diff)
	assert.Nil(t, err)
	resourceData.SetId("MockId")

	requests := make(map[string]map[string]interface{})
	patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, query map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
		requests[*action] = query
		return map[string]interface{}{
			"Name":              "accelerator_name",
			"State":             "active",
			"Description":       "description",
			"Spec":              "2",
			"DdosId":            "ddoscoo-cn-test",
			"RenewalStatus":     "Normal",
			"AutoRenewDuration": 1,
		}, nil
	})
	err = resourceAliCloudGaAcceleratorUpdate(resourceData, rawClient)
	patches.Reset()
	assert.Nil(t, err)

	attachReq, ok := requests["AttachDdosToAccelerator"]
	assert.True(t, ok)
	assert.Equal(t, "ddoscoo-cn-test", attachReq["DdosConfigList.1.DdosId"])
	assert.Equal(t, "cn-hangzhou", attachReq["DdosConfigList.1.DdosRegionId"])

	updateReq, ok := requests["UpdateAccelerator"]
	assert.True(t, ok)
	assert.Equal(t, false, updateReq["AutoPay"])
}

// lintignore: R001
func TestUnitAliCloudGaAcceleratorCreateRequest(t *testing.T) {
	region := os.Getenv("ALICLOUD_REGION")
	rawClient, err := unitTestSharedClientForRegion(region)
	if err != nil {
		t.Skipf("Skipping the test case with err: %s", err)
	}

	d, _ := schema.InternalMap(unitTestProviderResources()["alicloud_ga_accelerator"].Schema).Data(nil, nil)
	d.MarkNewResource()
	assert.Nil(t, d.Set("spec", "1"))
	assert.Nil(t, d.Set("auto_renew", true))
	assert.Nil(t, d.Set("accelerator_name", "terraform-example-ga-accelerator"))

	requests := make(map[string]map[string]interface{})
	patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, query map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
		requests[*action] = query
		return map[string]interface{}{
			"AcceleratorId":     "MockId",
			"Name":              "accelerator_name",
			"State":             "active",
			"Spec":              "1",
			"RenewalStatus":     "AutoRenewal",
			"AutoRenewDuration": 1,
		}, nil
	})
	err = resourceAliCloudGaAcceleratorCreate(d, rawClient)
	patches.Reset()
	assert.Nil(t, err)

	createReq, ok := requests["CreateAccelerator"]
	assert.True(t, ok)
	assert.Equal(t, true, createReq["AutoRenew"])
	assert.Equal(t, true, createReq["AutoPay"])
	assert.Equal(t, "terraform-example-ga-accelerator", createReq["Name"])
}

// lintignore: R001
func TestUnitAliCloudGaAcceleratorUserDefineBandwidthGuard(t *testing.T) {
	d, _ := schema.InternalMap(unitTestProviderResources()["alicloud_ga_accelerator"].Schema).Data(nil, nil)
	d.MarkNewResource()
	assert.Nil(t, d.Set("bandwidth", 200))
	assert.Nil(t, d.Set("ip_set_config", []interface{}{
		map[string]interface{}{"access_mode": "UserDefine"},
	}))

	err := resourceAliCloudGaAcceleratorCreate(d, &connectivity.AliyunClient{})
	assert.NotNil(t, err)
	assert.Contains(t, fmt.Sprint(err), "bandwidth is valid only when access_mode is set to Anycast")
}

// lintignore: R001
func TestUnitAliCloudGaAcceleratorDdosDetachRequest(t *testing.T) {
	region := os.Getenv("ALICLOUD_REGION")
	rawClient, err := unitTestSharedClientForRegion(region)
	if err != nil {
		t.Skipf("Skipping the test case with err: %s", err)
	}

	state := &terraform.InstanceState{ID: "MockId"}
	diff, err := newInstanceDiff("alicloud_ga_accelerator", map[string]interface{}{
		"ddos_config_list": []interface{}{
			map[string]interface{}{
				"ddos_id":        "ddoscoo-cn-test",
				"ddos_region_id": "cn-hangzhou",
			},
		},
	}, map[string]interface{}{
		"ddos_config_list": []interface{}{},
	}, state)
	assert.Nil(t, err)

	resourceData, err := schema.InternalMap(unitTestProviderResources()["alicloud_ga_accelerator"].Schema).Data(state, diff)
	assert.Nil(t, err)
	resourceData.SetId("MockId")

	requests := make(map[string]map[string]interface{})
	patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, query map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
		requests[*action] = query
		return map[string]interface{}{
			"Name":              "accelerator_name",
			"State":             "active",
			"Description":       "description",
			"Spec":              "1",
			"DdosId":            "ddoscoo-cn-test",
			"RenewalStatus":     "Normal",
			"AutoRenewDuration": 1,
		}, nil
	})
	err = resourceAliCloudGaAcceleratorUpdate(resourceData, rawClient)
	patches.Reset()
	assert.Nil(t, err)

	detachReq, ok := requests["DetachDdosFromAccelerator"]
	assert.True(t, ok)
	assert.Equal(t, "ddoscoo-cn-test", detachReq["DdosConfigList.1.DdosId"])
	assert.Equal(t, "cn-hangzhou", detachReq["DdosConfigList.1.DdosRegionId"])
}

// TestUnitAliCloudGaAcceleratorPlanValidation covers the plan-time rejections
// through the public Resource.Diff entry with a real RawConfig: incomplete DDoS
// entries, writes to the deprecated read-only aliases and the mutually
// exclusive renewal arguments. It also proves the plan never creates a cloud
// client.
// lintignore: R001
func TestUnitAliCloudGaAcceleratorPlanValidation(t *testing.T) {
	accelerator := ResourceAliCloudGaAccelerator()
	clientCalls := 0
	patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewGaplusClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		clientCalls++
		return nil, fmt.Errorf("unexpected cloud client creation during plan")
	})
	defer patches.Reset()

	testCases := []struct {
		name    string
		raw     map[string]interface{}
		wantErr string
	}{
		{
			name:    "ddos_config_list entry missing ddos_id",
			raw:     map[string]interface{}{"spec": "1", "ddos_config_list": []interface{}{map[string]interface{}{"ddos_region_id": "cn-hangzhou"}}},
			wantErr: "ddos_id and ddos_region_id are required in ddos_config_list",
		},
		{
			name:    "ddos_config_list entry missing ddos_region_id",
			raw:     map[string]interface{}{"spec": "1", "ddos_config_list": []interface{}{map[string]interface{}{"ddos_id": "ddoscoo-cn-test"}}},
			wantErr: "ddos_id and ddos_region_id are required in ddos_config_list",
		},
		{
			name:    "explicit ddos_id alias is rejected",
			raw:     map[string]interface{}{"spec": "1", "ddos_id": "ddoscoo-cn-test"},
			wantErr: "`ddos_id` is read-only and deprecated",
		},
		{
			name:    "explicit ddos_region_id alias is rejected",
			raw:     map[string]interface{}{"spec": "1", "ddos_region_id": "cn-hangzhou"},
			wantErr: "`ddos_region_id` is read-only and deprecated",
		},
		{
			name:    "auto_renew and renewal_status are mutually exclusive",
			raw:     map[string]interface{}{"spec": "1", "auto_renew": true, "renewal_status": "AutoRenewal"},
			wantErr: "`auto_renew` conflicts with `renewal_status`",
		},
		{
			name: "valid ddos entry passes plan validation",
			raw:  map[string]interface{}{"spec": "1", "ddos_config_list": []interface{}{map[string]interface{}{"ddos_id": "ddoscoo-cn-test", "ddos_region_id": "cn-hangzhou"}}},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			state := &terraform.InstanceState{RawConfig: acceleratorRawConfig(testCase.raw)}
			_, err := accelerator.Diff(context.Background(), state, terraform.NewResourceConfigRaw(testCase.raw), &connectivity.AliyunClient{})
			if testCase.wantErr == "" {
				assert.NoError(t, err)
				return
			}
			if assert.Error(t, err) {
				assert.Contains(t, err.Error(), testCase.wantErr)
			}
		})
	}

	assert.Equal(t, 0, clientCalls, "plan validation must not create a cloud client")
}

// TestUnitAliCloudGaAcceleratorDdosConfigListPreWriteGuard proves the
// Create/Update entry points reject an incomplete ddos_config_list entry before
// any cloud request is created, including when other fields would change.
// lintignore: R001
func TestUnitAliCloudGaAcceleratorDdosConfigListPreWriteGuard(t *testing.T) {
	clientCalls := 0
	patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewGaplusClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		clientCalls++
		return nil, fmt.Errorf("unexpected cloud client creation")
	})
	defer patches.Reset()

	t.Run("create rejects incomplete ddos entry before any request", func(t *testing.T) {
		d, _ := schema.InternalMap(unitTestProviderResources()["alicloud_ga_accelerator"].Schema).Data(nil, nil)
		d.MarkNewResource()
		assert.Nil(t, d.Set("spec", "1"))
		assert.Nil(t, d.Set("ddos_config_list", []interface{}{
			map[string]interface{}{"ddos_id": "ddoscoo-cn-test"},
		}))

		err := resourceAliCloudGaAcceleratorCreate(d, &connectivity.AliyunClient{})
		if assert.Error(t, err) {
			assert.Contains(t, err.Error(), "ddos_id and ddos_region_id are required in ddos_config_list")
		}
	})

	t.Run("update rejects incomplete ddos entry before other writes", func(t *testing.T) {
		state := &terraform.InstanceState{ID: "MockId", Attributes: map[string]string{"spec": "1"}}
		diff, err := newInstanceDiff("alicloud_ga_accelerator", map[string]interface{}{
			"spec": "1",
		}, map[string]interface{}{
			"spec": "2",
			"ddos_config_list": []interface{}{
				map[string]interface{}{"ddos_region_id": "cn-hangzhou"},
			},
		}, state)
		assert.Nil(t, err)

		resourceData, err := schema.InternalMap(unitTestProviderResources()["alicloud_ga_accelerator"].Schema).Data(state, diff)
		assert.Nil(t, err)
		resourceData.SetId("MockId")

		err = resourceAliCloudGaAcceleratorUpdate(resourceData, &connectivity.AliyunClient{})
		if assert.Error(t, err) {
			assert.Contains(t, err.Error(), "ddos_id and ddos_region_id are required in ddos_config_list")
		}
	})

	assert.Equal(t, 0, clientCalls, "write guards must reject before any cloud request")
}

// acceleratorRawConfig builds the plan-time RawConfig value the SDK exposes to
// CustomizeDiff. Resource.Diff only carries a RawConfig through the state, so
// the negative cases above pass an InstanceState that holds the same config.
func acceleratorRawConfig(raw map[string]interface{}) cty.Value {
	attributeTypes := ResourceAliCloudGaAccelerator().CoreConfigSchema().ImpliedType().AttributeTypes()
	values := make(map[string]cty.Value, len(attributeTypes))
	for name, attributeType := range attributeTypes {
		rawValue, ok := raw[name]
		if !ok {
			values[name] = cty.NullVal(attributeType)
			continue
		}
		values[name] = acceleratorRawConfigAttribute(rawValue, attributeType)
	}
	return cty.ObjectVal(values)
}

func acceleratorRawConfigAttribute(rawValue interface{}, attributeType cty.Type) cty.Value {
	switch {
	case attributeType.Equals(cty.String):
		value, _ := rawValue.(string)
		return cty.StringVal(value)
	case attributeType.Equals(cty.Number):
		value, _ := rawValue.(int)
		return cty.NumberIntVal(int64(value))
	case attributeType.Equals(cty.Bool):
		value, _ := rawValue.(bool)
		return cty.BoolVal(value)
	case attributeType.IsSetType() || attributeType.IsListType():
		elementType := attributeType.ElementType()
		items, _ := rawValue.([]interface{})
		elements := make([]cty.Value, 0, len(items))
		for _, item := range items {
			itemMap, _ := item.(map[string]interface{})
			elementAttributeTypes := elementType.AttributeTypes()
			elementValues := make(map[string]cty.Value, len(elementAttributeTypes))
			for name, elementAttributeType := range elementAttributeTypes {
				elementRawValue, ok := itemMap[name]
				if !ok {
					elementValues[name] = cty.NullVal(elementAttributeType)
					continue
				}
				elementValues[name] = acceleratorRawConfigAttribute(elementRawValue, elementAttributeType)
			}
			elements = append(elements, cty.ObjectVal(elementValues))
		}
		if attributeType.IsSetType() {
			if len(elements) == 0 {
				return cty.SetValEmpty(elementType)
			}
			return cty.SetVal(elements)
		}
		if len(elements) == 0 {
			return cty.ListValEmpty(elementType)
		}
		return cty.ListVal(elements)
	default:
		return cty.NullVal(attributeType)
	}
}

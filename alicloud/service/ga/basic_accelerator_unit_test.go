package ga

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"reflect"
	"strings"
	"sync"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/alibabacloud-go/tea-rpc/client"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/alibabacloud-go/tea/tea"
	credentials "github.com/aliyun/credentials-go/credentials"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
)

// White-box (internal package ga) TestUnit cases moved out of the
// external ga_test file so they can call the unexported resource CRUD
// funcs. Shared local helpers (unitTestProviderResources /
// newInstanceDiff / unitTestSharedClientForRegion) live in
// unit_test_helpers_test.go; see that file for why acctest cannot be
// imported here.

// lintignore: R001
func TestUnitAliCloudGaBasicAccelerator(t *testing.T) {
	p := unitTestProviderResources()
	dInit, _ := schema.InternalMap(p["alicloud_ga_basic_accelerator"].Schema).Data(nil, nil)
	dExisted, _ := schema.InternalMap(p["alicloud_ga_basic_accelerator"].Schema).Data(nil, nil)
	dInit.MarkNewResource()
	attributes := map[string]interface{}{
		"duration":               1,
		"pricing_cycle":          "CreateGaBasicAcceleratorValue",
		"basic_accelerator_name": "CreateGaBasicAcceleratorValue",
		"description":            "CreateGaBasicAcceleratorValue",
		"bandwidth_billing_type": "CreateGaBasicAcceleratorValue",
		"auto_pay":               true,
		"auto_use_coupon":        "true",
		"auto_renew":             false,
		"auto_renew_duration":    1,
	}
	for key, value := range attributes {
		err := dInit.Set(key, value)
		assert.Nil(t, err)
		err = dExisted.Set(key, value)
		assert.Nil(t, err)
		if err != nil {
			log.Printf("[ERROR] the field %s setting error", key)
		}
	}
	region := os.Getenv("ALICLOUD_REGION")
	rawClient, err := unitTestSharedClientForRegion(region)
	if err != nil {
		t.Skipf("Skipping the test case with err: %s", err)
		t.Skipped()
	}

	rawClient = rawClient.(*connectivity.AliyunClient)
	ReadMockResponse := map[string]interface{}{
		// GetBasicAccelerator
		"Name":                 "CreateGaBasicAcceleratorValue",
		"Description":          "CreateGaBasicAcceleratorValue",
		"BandwidthBillingType": "CreateGaBasicAcceleratorValue",
		"State":                "active",
	}
	CreateMockResponse := map[string]interface{}{
		"AcceleratorId": "CreateGaBasicAcceleratorValue",
	}
	failedResponseMock := func(errorCode string) (map[string]interface{}, error) {
		return nil, &tea.SDKError{
			Code:       tea.String(errorCode),
			Data:       tea.String(errorCode),
			Message:    tea.String(errorCode),
			StatusCode: tea.Int(400),
		}
	}
	notFoundResponseMock := func(errorCode string) (map[string]interface{}, error) {
		return nil, errs.GetNotFoundErrorFromString(errs.GetNotFoundMessage("alicloud_ga_basic_accelerator", errorCode))
	}
	successResponseMock := func(operationMockResponse map[string]interface{}) (map[string]interface{}, error) {
		if len(operationMockResponse) > 0 {
			helper.MapMerge(ReadMockResponse, operationMockResponse)
		}
		return ReadMockResponse, nil
	}
	// Create
	patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewGaplusClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:       tea.String("loadEndpoint error"),
			Data:       tea.String("loadEndpoint error"),
			Message:    tea.String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = resourceAliCloudGaBasicAcceleratorCreate(dInit, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	ReadMockResponseDiff := map[string]interface{}{}
	errorCodes := []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1 // a counter used to cover retry scenario; the same below
		patches = gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "CreateBasicAccelerator" {
				switch errorCode {
				case "NonRetryableError":
					return failedResponseMock(errorCode)
				default:
					retryIndex++
					if retryIndex >= len(errorCodes)-1 {
						successResponseMock(ReadMockResponseDiff)
						return CreateMockResponse, nil
					}
					return failedResponseMock(errorCodes[retryIndex])
				}
			}
			return ReadMockResponse, nil
		})
		err := resourceAliCloudGaBasicAcceleratorCreate(dInit, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_ga_basic_accelerator"].Schema).Data(dInit.State(), nil)
			for key, value := range attributes {
				_ = dCompare.Set(key, value)
			}
			assert.Equal(t, dCompare.State().Attributes, dInit.State().Attributes)
		}
		if retryIndex >= len(errorCodes)-1 {
			break
		}
	}

	// Update
	patches = gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewGaplusClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:       tea.String("loadEndpoint error"),
			Data:       tea.String("loadEndpoint error"),
			Message:    tea.String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = resourceAliCloudGaBasicAcceleratorUpdate(dExisted, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	attributesDiff := map[string]interface{}{
		"basic_accelerator_name": "PutGaBasicAcceleratorValue",
		"description":            "PutGaBasicAcceleratorValue",
	}
	diff, err := newInstanceDiff("alicloud_ga_basic_accelerator", attributes, attributesDiff, dInit.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_ga_basic_accelerator"].Schema).Data(dInit.State(), diff)
	ReadMockResponseDiff = map[string]interface{}{
		// GetBasicAccelerator Response
		"Name":        "PutGaBasicAcceleratorValue",
		"Description": "PutGaBasicAcceleratorValue",
		"State":       "active",
	}
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches = gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "UpdateBasicAccelerator" {
				switch errorCode {
				case "NonRetryableError":
					return failedResponseMock(errorCode)
				default:
					retryIndex++
					if retryIndex >= len(errorCodes)-1 {
						return successResponseMock(ReadMockResponseDiff)
					}
					return failedResponseMock(errorCodes[retryIndex])
				}
			}
			return ReadMockResponse, nil
		})
		err := resourceAliCloudGaBasicAcceleratorUpdate(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_ga_basic_accelerator"].Schema).Data(dExisted.State(), nil)
			for key, value := range attributes {
				_ = dCompare.Set(key, value)
			}
			assert.Equal(t, dCompare.State().Attributes, dExisted.State().Attributes)
		}
		if retryIndex >= len(errorCodes)-1 {
			break
		}
	}

	// Read
	diff, err = newInstanceDiff("alicloud_ga_basic_accelerator", attributes, attributesDiff, dInit.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_ga_basic_accelerator"].Schema).Data(dInit.State(), diff)
	errorCodes = []string{"NonRetryableError", "Throttling", "nil", "{}"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches = gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "GetBasicAccelerator" {
				switch errorCode {
				case "{}":
					return notFoundResponseMock(errorCode)
				case "NonRetryableError":
					return failedResponseMock(errorCode)
				default:
					retryIndex++
					if errorCodes[retryIndex] == "nil" {
						return ReadMockResponse, nil
					}
					return failedResponseMock(errorCodes[retryIndex])
				}
			}
			return ReadMockResponse, nil
		})
		err := resourceAliCloudGaBasicAcceleratorRead(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		case "{}":
			assert.Nil(t, err)
		}
	}

	// Delete
	patches = gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewGaplusClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:       tea.String("loadEndpoint error"),
			Data:       tea.String("loadEndpoint error"),
			Message:    tea.String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = resourceAliCloudGaBasicAcceleratorDelete(dExisted, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	attributesDiff = map[string]interface{}{}
	diff, err = newInstanceDiff("alicloud_ga_basic_accelerator", attributes, attributesDiff, dInit.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_ga_basic_accelerator"].Schema).Data(dInit.State(), diff)
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "DeleteBasicAccelerator" {
				switch errorCode {
				case "NonRetryableError":
					return failedResponseMock(errorCode)
				default:
					retryIndex++
					if errorCodes[retryIndex] == "nil" {
						ReadMockResponse = map[string]interface{}{}
						return ReadMockResponse, nil
					}
					return failedResponseMock(errorCodes[retryIndex])
				}
			}
			return ReadMockResponse, nil
		})
		err := resourceAliCloudGaBasicAcceleratorDelete(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		case "nil":
			assert.Nil(t, err)
		}
	}
}

// TestUnitAliCloudGaBasicAcceleratorPaymentMapping pins payment/coupon mappings at a local RPC boundary.
func TestUnitAliCloudGaBasicAcceleratorPaymentMapping(t *testing.T) {
	p := unitTestProviderResources()
	d, _ := schema.InternalMap(p["alicloud_ga_basic_accelerator"].Schema).Data(nil, nil)
	d.MarkNewResource()
	assert.Nil(t, d.Set("duration", 1))
	assert.Nil(t, d.Set("pricing_cycle", "Month"))
	assert.Nil(t, d.Set("basic_accelerator_name", "tf-example-value"))
	assert.Nil(t, d.Set("description", "tf-example-value"))
	assert.Nil(t, d.Set("bandwidth_billing_type", "BandwidthPackage"))
	assert.Nil(t, d.Set("auto_pay", true))
	assert.Nil(t, d.Set("auto_use_coupon", "true"))
	assert.Nil(t, d.Set("payment_type", "PayAsYouGo"))

	readResponse := map[string]interface{}{
		"Name":                 "tf-example-value",
		"Description":          "tf-example-value",
		"BandwidthBillingType": "BandwidthPackage",
		"InstanceChargeType":   "PREPAY",
		"State":                "active",
	}
	var requestMutex sync.Mutex
	var createQuery, createBody url.Values
	createCaptured := false
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			t.Error(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		var response map[string]interface{}
		switch r.Form.Get("Action") {
		case "CreateBasicAccelerator":
			requestMutex.Lock()
			createQuery = r.URL.Query()
			createBody = r.PostForm
			createCaptured = true
			requestMutex.Unlock()
			response = map[string]interface{}{"AcceleratorId": "ga-basic-accelerator-id"}
		case "GetBasicAccelerator":
			response = readResponse
		case "ListTagResources":
			response = map[string]interface{}{"TagResources": map[string]interface{}{"TagResource": []interface{}{}}}
		default:
			t.Errorf("unexpected action: %s", r.Form.Get("Action"))
			w.WriteHeader(http.StatusBadRequest)
			response = map[string]interface{}{"Code": "UnexpectedAction"}
		}
		if err := json.NewEncoder(w).Encode(response); err != nil {
			t.Error(err)
		}
	}))
	t.Cleanup(server.Close)

	credential, err := credentials.NewCredential(new(credentials.Config).
		SetType("access_key").SetAccessKeyId("test-key").SetAccessKeySecret("test-secret"))
	if err != nil {
		t.Fatal(err)
	}
	endpoints := new(sync.Map)
	endpoint := strings.TrimPrefix(server.URL, "http://")
	t.Setenv("NO_PROXY", endpoint)
	clientConfig := &connectivity.Config{
		AccessKey: "test-key", SecretKey: "test-secret", Credential: credential,
		RegionId: "cn-hangzhou", AccountType: "test", Protocol: "http",
		Endpoints: endpoints, SignVersion: new(sync.Map), SkipRegionValidation: true,
	}
	rawClient, err := connectivity.NewClient(clientConfig)
	if err != nil {
		t.Fatalf("building the local client failed: %s", err)
	}
	endpoints.Store("ga", endpoint)

	assert.Nil(t, resourceAliCloudGaBasicAcceleratorCreate(d, rawClient))
	requestMutex.Lock()
	captured := createCaptured
	query := createQuery
	body := createBody
	requestMutex.Unlock()
	assert.True(t, captured, "CreateBasicAccelerator request was not captured")
	assert.Equal(t, "POSTPAY", body.Get("ChargeType"))
	assert.Equal(t, "true", body.Get("AutoUseCoupon"))
	assert.Equal(t, "true", body.Get("AutoPay"))
	for _, key := range []string{"ChargeType", "AutoUseCoupon", "AutoPay"} {
		assert.NotContains(t, query, key)
	}

	assert.Nil(t, resourceAliCloudGaBasicAcceleratorRead(d, rawClient))
	assert.Equal(t, "Subscription", d.Get("payment_type"))
}

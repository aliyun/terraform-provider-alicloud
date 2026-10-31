package alicloud

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"
	"sync"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/alibabacloud-go/tea-rpc/client"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/credentials-go/credentials"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/assert"
)

func TestAccAliCloudSDDPInstance_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_sddp_instance.default"
	ra := resourceAttrInit(resourceId, AlicloudSDDPInstanceMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &SddpService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeSddpInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%ssddpinstance%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudSDDPInstanceBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithTime(t, []int{30})
		},
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"dataphin_count": "CreateInstance",
					"logistics":      "yes",
					"oss_size":       "3000",
					"renew_period":   "1",
					"renewal_status": "CreateInstance",
					"payment_type":   "Subscription",
					"sddp_version":   "version_company",
					"sdc":            "3",
					"sd_cbool":       "no",
					"period":         "1",
					"ud_cbool":       "yes",
					"udc":            "2000",
					"dataphin":       "yes",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"payment_type": "Subscription",
						"sddp_version": "version_company",
						"sd_cbool":     "no",
						"sdc":          "3",
						"ud_cbool":     "yes",
						"udc":          "2000",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"dataphin_count": "ModifyInstance",
					"logistics":      "no",
					"oss_size":       "3001",
					"renew_period":   "2",
					"renewal_status": "AutoRenewal",
					"sddp_version":   "version_audit",
					"modify_type":    "Upgrade",
					"sdc":            "8",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"sdc": "8",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"modify_type": "Upgrade",
					"sd_cbool":    "no",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"sd_cbool": "no",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"modify_type": "Upgrade",
					"udc":         "3000",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"udc": "3000",
					}),
				),
			},

			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"sd_cbool", "period", "ud_cbool", "udc", "dataphin", "modify_type"},
			},
		},
	})
}

var AlicloudSDDPInstanceMap0 = map[string]string{}

func AlicloudSDDPInstanceBasicDependence0(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}
`, name)
}

// lintignore: R001
func TestUnitAlicloudSDDPInstance(t *testing.T) {
	p := Provider().ResourcesMap
	dInit, _ := schema.InternalMap(p["alicloud_sddp_instance"].Schema).Data(nil, nil)
	dExisted, _ := schema.InternalMap(p["alicloud_sddp_instance"].Schema).Data(nil, nil)
	dInit.MarkNewResource()
	attributes := map[string]interface{}{
		"payment_type":   "Subscription",
		"sddp_version":   "CreateInstance",
		"sdc":            "3",
		"sd_cbool":       "CreateInstance",
		"period":         1,
		"ud_cbool":       "CreateInstance",
		"udc":            "2000",
		"dataphin":       "CreateInstance",
		"renewal_status": "CreateInstance",
		"dataphin_count": "CreateInstance",
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
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		t.Skipf("Skipping the test case with err: %s", err)
		t.Skipped()
	}
	rawClient = rawClient.(*connectivity.AliyunClient)
	ReadMockResponse := map[string]interface{}{
		"UserStatus": map[string]interface{}{
			"ChargeType":     "PREPAY",
			"Version":        "CreateInstance",
			"InstanceNum":    3,
			"Authed":         "CreateInstance",
			"RenewStatus":    "CreateInstance",
			"OdpsSet":        "CreateInstance",
			"OssBucketSet":   "CreateInstance",
			"OssSize":        2000,
			"RdsSet":         "CreateInstance",
			"InstanceStatus": 1,
			"RemainDays":     "CreateInstance",
		},
		"Code": "Success",
		"Data": map[string]interface{}{
			"InstanceId": "CreateInstance",
		},
	}
	CreateMockResponse := map[string]interface{}{
		"Code": "Success",
		"Data": map[string]interface{}{
			"InstanceId": "CreateInstance",
		},
	}
	failedResponseMock := func(errorCode string) (map[string]interface{}, error) {
		return nil, &tea.SDKError{
			Code:       String(errorCode),
			Data:       String(errorCode),
			Message:    String(errorCode),
			StatusCode: tea.Int(400),
		}
	}
	notFoundResponseMock := func(errorCode string) (map[string]interface{}, error) {
		return nil, GetNotFoundErrorFromString(GetNotFoundMessage("alicloud_sddp_instance", errorCode))
	}
	successResponseMock := func(operationMockResponse map[string]interface{}) (map[string]interface{}, error) {
		if len(operationMockResponse) > 0 {
			mapMerge(ReadMockResponse, operationMockResponse)
		}
		return ReadMockResponse, nil
	}

	// Create
	patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewBssopenapiClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:       String("loadEndpoint error"),
			Data:       String("loadEndpoint error"),
			Message:    String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = resourceAlicloudSddpInstanceCreate(dInit, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	ReadMockResponseDiff := map[string]interface{}{
		"Code": "Success",
		"Data": map[string]interface{}{
			"InstanceId": "CreateInstance",
		},
	}
	errorCodes := []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1 // a counter used to cover retry scenario; the same below
		patches = gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "CreateInstance" {
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
		err := resourceAlicloudSddpInstanceCreate(dInit, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_sddp_instance"].Schema).Data(dInit.State(), nil)
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
	patches = gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewBssopenapiClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:       String("loadEndpoint error"),
			Data:       String("loadEndpoint error"),
			Message:    String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = resourceAlicloudSddpInstanceUpdate(dExisted, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	attributesDiff := map[string]interface{}{
		"sdc":            "4",
		"oss_size":       "3000",
		"sd_cbool":       "ModifyInstance",
		"ud_cbool":       "ModifyInstance",
		"sddp_version":   "ModifyInstance",
		"udc":            "ModifyInstance",
		"dataphin":       "ModifyInstance",
		"dataphin_count": "ModifyInstance",
		"modify_type":    "ModifyInstance",
	}
	diff, err := newInstanceDiff("alicloud_sddp_instance", attributes, attributesDiff, dInit.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_sddp_instance"].Schema).Data(dInit.State(), diff)
	ReadMockResponseDiff = map[string]interface{}{
		"UserStatus": map[string]interface{}{
			"Version":        "ModifyInstance",
			"InstanceNum":    4,
			"Authed":         "ModifyInstance",
			"OdpsSet":        "ModifyInstance",
			"OssBucketSet":   "ModifyInstance",
			"OssSize":        3000,
			"InstanceStatus": 1,
			"Code":           "Success",
		},
	}
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches = gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "ModifyInstance" {
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
		err := resourceAlicloudSddpInstanceUpdate(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_sddp_instance"].Schema).Data(dExisted.State(), nil)
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
	attributesDiff = map[string]interface{}{}
	diff, err = newInstanceDiff("alicloud_sddp_instance", attributes, attributesDiff, dInit.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_sddp_instance"].Schema).Data(dInit.State(), diff)
	errorCodes = []string{"NonRetryableError", "Throttling", "nil", "{}"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches = gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "DescribeUserStatus" {
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
		err := resourceAlicloudSddpInstanceRead(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		case "{}":
			assert.Nil(t, err)
		}
	}

	// Delete
	err = resourceAlicloudSddpInstanceDelete(dExisted, rawClient)
	assert.Nil(t, err)
}

// TestUnitAlicloudSDDPInstanceRenewPeriodCreateContract verifies that Create assembles
// the real BssOpenApi CreateInstance request body correctly for renewal semantics,
// using a normally-initialized client with static fake credentials and a fail-closed
// DoRequest interception. No real network, AK, or payment is involved.
//
// Cases:
//  1. AutoRenewal + renew_period=1  -> body.RenewalStatus=AutoRenewal, body.RenewPeriod=1
//  2. AutoRenewal without renew_period -> rejected before any RPC (error names renew_period)
//  3. ManualRenewal + renew_period=2 -> body.RenewalStatus=ManualRenewal, body.RenewPeriod=2
//  4. ManualRenewal without renew_period -> accepted, body carries no RenewPeriod
//
// On the pre-fix source (which read the non-existent "renewal_duration" schema key),
// cases 1, 2 and 3 fail; after the fix switches to "renew_period" they pass,
// preserving fix-before-fail / fix-after-pass regression semantics.
//
// lintignore: R001
func TestUnitAlicloudSDDPInstanceRenewPeriodCreateContract(t *testing.T) {
	p := Provider().ResourcesMap
	resourceSchema := p["alicloud_sddp_instance"].Schema

	// Build a normally-initialized AliyunClient with static fake credentials and
	// pre-stored non-routable local endpoints. SkipRegionValidation avoids the
	// region network check; AccountType="Domestic" short-circuits getAccountType's
	// QueryAvailableInstances probe; pre-stored endpoints avoid loadEndpoint
	// discovery. Construction is fully offline.
	endpoints := &sync.Map{}
	endpoints.Store("bssopenapi", "127.0.0.1:0") // SDDP Create via BssOpenApi
	endpoints.Store("sddp", "127.0.0.1:0")       // Read via Sddp/DescribeUserStatus
	conf := connectivity.Config{
		Region:               connectivity.Region("cn-hangzhou"),
		RegionId:             "cn-hangzhou",
		AccessKey:            "FakeAccessKeyID",
		SecretKey:            "FakeAccessKeySecret",
		Protocol:             "HTTPS",
		Endpoints:            endpoints,
		SignVersion:          &sync.Map{},
		SkipRegionValidation: true,
		AccountId:            "0000000000000000",
		AccountType:          "Domestic",
	}
	credConfig := new(credentials.Config).
		SetType("access_key").
		SetAccessKeyId("FakeAccessKeyID").
		SetAccessKeySecret("FakeAccessKeySecret")
	cred, err := credentials.NewCredential(credConfig)
	assert.Nil(t, err)
	conf.Credential = cred
	rawClient, err := conf.Client()
	assert.Nil(t, err)

	readMock := map[string]interface{}{
		"Code": "Success",
		"UserStatus": map[string]interface{}{
			"ChargeType":     "PREPAY",
			"Version":        "version_dlp",
			"InstanceNum":    3,
			"Authed":         "CreateInstance",
			"RenewStatus":    "AutoRenewal",
			"OdpsSet":        "CreateInstance",
			"OssBucketSet":   "CreateInstance",
			"OssSize":        2000,
			"RdsSet":         "CreateInstance",
			"InstanceStatus": 1,
			"RemainDays":     "CreateInstance",
		},
	}

	// baseAttributes are the non-renewal required Create inputs. version_dlp does
	// not require dataphin, keeping the fixture minimal. Each case layers renewal
	// fields on top of these.
	baseAttributes := map[string]interface{}{
		"payment_type": "Subscription",
		"period":       1,
		"sddp_version": "version_dlp",
		"sdc":          "3",
		"sd_cbool":     "no",
		"ud_cbool":     "yes",
		"udc":          "2000",
	}

	newResourceData := func(t *testing.T, extra map[string]interface{}) *schema.ResourceData {
		d, _ := schema.InternalMap(resourceSchema).Data(nil, nil)
		d.MarkNewResource()
		for k, v := range baseAttributes {
			assert.Nil(t, d.Set(k, v))
		}
		for k, v := range extra {
			assert.Nil(t, d.Set(k, v))
		}
		return d
	}

	createInstanceMock := func(instanceId string) map[string]interface{} {
		return map[string]interface{}{
			"Code": "Success",
			"Data": map[string]interface{}{
				"InstanceId": instanceId,
			},
		}
	}

	// failClosedMock intercepts every DoRequest: CreateInstance returns a success
	// mock (and exposes the request body), DescribeUserStatus returns the Read
	// mock, and any other action fails closed so no unexpected external call can
	// ever escape. Because DoRequest is a large leaf method it is not inlined, so
	// the gomonkey patch reliably intercepts real tea-rpc calls.
	failClosedMock := func(capturedBody *map[string]interface{}) func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, body map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
		return func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, body map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			switch *action {
			case "CreateInstance":
				if capturedBody != nil {
					*capturedBody = body
				}
				return createInstanceMock("mock-sddp-instance"), nil
			case "DescribeUserStatus":
				return readMock, nil
			default:
				return nil, fmt.Errorf("fail-closed: unexpected DoRequest action %q", *action)
			}
		}
	}

	// Case 1: AutoRenewal + renew_period=1 -> Create succeeds, body carries
	// RenewalStatus=AutoRenewal and RenewPeriod=1.
	t.Run("AutoRenewalWithRenewPeriodAssemblesBody", func(t *testing.T) {
		var capturedBody map[string]interface{}
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", failClosedMock(&capturedBody))
		defer patches.Reset()

		d := newResourceData(t, map[string]interface{}{
			"renewal_status": "AutoRenewal",
			"renew_period":   1,
		})
		err := resourceAlicloudSddpInstanceCreate(d, rawClient)
		assert.Nil(t, err, "AutoRenewal+renew_period Create should succeed")
		assert.NotNil(t, capturedBody, "CreateInstance DoRequest should be invoked")
		assert.Equal(t, "AutoRenewal", fmt.Sprint(capturedBody["RenewalStatus"]))
		assert.Equal(t, "1", fmt.Sprint(capturedBody["RenewPeriod"]))
	})

	// Case 2: AutoRenewal without renew_period -> rejected before any RPC, error
	// names renew_period (not the non-existent renewal_duration).
	t.Run("AutoRenewalMissingRenewPeriodRejectedPreRPC", func(t *testing.T) {
		var doRequestCalls int
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			doRequestCalls++
			return nil, fmt.Errorf("fail-closed: DoRequest must not be called for missing renew_period")
		})
		defer patches.Reset()

		d := newResourceData(t, map[string]interface{}{
			"renewal_status": "AutoRenewal",
		})
		err := resourceAlicloudSddpInstanceCreate(d, rawClient)
		assert.NotNil(t, err, "AutoRenewal without renew_period must be rejected")
		assert.Equal(t, 0, doRequestCalls, "rejection must occur before any RPC")
		assert.True(t, strings.Contains(err.Error(), "renew_period"), "error must name renew_period, got: %v", err)
		assert.False(t, strings.Contains(err.Error(), "renewal_duration"), "error must not mention the non-existent renewal_duration key, got: %v", err)
	})

	// Case 3: ManualRenewal + renew_period=2 -> Create succeeds, body carries
	// RenewalStatus=ManualRenewal and RenewPeriod=2.
	t.Run("ManualRenewalWithRenewPeriodAssemblesBody", func(t *testing.T) {
		var capturedBody map[string]interface{}
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", failClosedMock(&capturedBody))
		defer patches.Reset()

		d := newResourceData(t, map[string]interface{}{
			"renewal_status": "ManualRenewal",
			"renew_period":   2,
		})
		err := resourceAlicloudSddpInstanceCreate(d, rawClient)
		assert.Nil(t, err, "ManualRenewal+renew_period Create should succeed")
		assert.NotNil(t, capturedBody, "CreateInstance DoRequest should be invoked")
		assert.Equal(t, "ManualRenewal", fmt.Sprint(capturedBody["RenewalStatus"]))
		assert.Equal(t, "2", fmt.Sprint(capturedBody["RenewPeriod"]))
	})

	// Case 4: ManualRenewal without renew_period -> accepted (not rejected),
	// body carries RenewalStatus=ManualRenewal but no RenewPeriod key.
	t.Run("ManualRenewalWithoutRenewPeriodAccepted", func(t *testing.T) {
		var capturedBody map[string]interface{}
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", failClosedMock(&capturedBody))
		defer patches.Reset()

		d := newResourceData(t, map[string]interface{}{
			"renewal_status": "ManualRenewal",
		})
		err := resourceAlicloudSddpInstanceCreate(d, rawClient)
		assert.Nil(t, err, "ManualRenewal without renew_period should be accepted")
		assert.NotNil(t, capturedBody, "CreateInstance DoRequest should be invoked")
		assert.Equal(t, "ManualRenewal", fmt.Sprint(capturedBody["RenewalStatus"]))
		_, hasRenewPeriod := capturedBody["RenewPeriod"]
		assert.False(t, hasRenewPeriod, "body must not carry RenewPeriod when renew_period unset")
	})
}

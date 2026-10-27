package drds

import (
	"fmt"
	"reflect"
	"sync"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	client "github.com/alibabacloud-go/tea-rpc/client"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/alibabacloud-go/tea/tea"
	credential "github.com/aliyun/credentials-go/credentials"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
)

// TestUnitDrdsPolardbxInstancePrepayRefundActionableError verifies the AC-3
// destroy contract for a subscription (PREPAY) PolarDB-X instance whose billing
// commodity does not support API unsubscription: BssOpenApi RefundInstance
// returns CommodityNotSupported, and Delete must surface an actionable error
// naming the instance ID and pointing to the BSS console unsubscribe page, while
// retaining the resource in state (Id stays set). This is deterministic and uses
// no live billing — DoRequest is mocked at the tea-rpc seam.
//
// It is the "verify destroy returns documented actionable error" evidence for
// AC-6: the live PREPAY ACC (env-gated) covers create/read/import, while this
// unit test covers the unsupported-refund destroy path without leaving a paid
// instance behind.
func TestUnitDrdsPolardbxInstancePrepayRefundActionableError(t *testing.T) {
	// buildClient constructs a dummy local *connectivity.AliyunClient whose
	// AccountType, polardbx and BSS endpoints are controlled by the caller, with
	// no network calls during construction. Mirrors the proven buildClient in the
	// QueryAvailableInstance unit test.
	buildClient := func() *connectivity.AliyunClient {
		endpoints := &sync.Map{}
		endpoints.Store("polardbx", "polardbx.dummy.aliyuncs.com")
		endpoints.Store("bssopenapi", connectivity.BssOpenAPIEndpointDomestic)
		credConfig := &credential.Config{
			Type:            tea.String("access_key"),
			AccessKeyId:     tea.String("dummy-access-key"),
			AccessKeySecret: tea.String("dummy-secret-key"),
		}
		cred, err := credential.NewCredential(credConfig)
		assert.NoError(t, err)
		cfg := &connectivity.Config{
			Region:               connectivity.Beijing,
			RegionId:             "cn-beijing",
			AccessKey:            "dummy-access-key",
			SecretKey:            "dummy-secret-key",
			AccountType:          "Domestic",
			Protocol:             "HTTPS",
			SkipRegionValidation: true,
			Endpoints:            endpoints,
			SignVersion:          &sync.Map{},
			Credential:           cred,
		}
		aliyunClient, err := cfg.Client()
		assert.NoError(t, err)
		return aliyunClient
	}

	instanceID := "pxc-test-actionable-instance-id"

	describeResponse := map[string]interface{}{
		"Success": true,
		"Code":    "200",
		"DBInstance": map[string]interface{}{
			"DBInstanceName": instanceID,
			"PayType":        "Prepaid",
			"Status":         "Running",
			"RegionId":       "cn-beijing",
		},
	}

	billingResponse := map[string]interface{}{
		"Success": true,
		"Code":    "200",
		"Data": map[string]interface{}{
			"InstanceList": []interface{}{
				map[string]interface{}{
					"InstanceID":       instanceID,
					"ProductCode":      "drds",
					"ProductType":      "drds_polarxpre_public_cn",
					"SubscriptionType": "Subscription",
				},
			},
		},
	}

	// refundError is a tea SDK error carrying the CommodityNotSupported code, the
	// shape RefundInstance returns for a billing commodity that does not support
	// API unsubscription. It is built with the exact field set that
	// connectivity.formatError (alicloud/connectivity/client.go) uses to construct
	// a real BSS SDKError — statusCode/code/message/data — so the Code, Message and
	// Data *string fields are all non-nil. errs.IsExpectedErrors dereferences *e.Data
	// on every non-matching code comparison (errors.go: "strings.Contains(*e.Data,
	// code)"); a mock SDKError that only sets "code"/"message" leaves Data nil and
	// panics there. Mirroring formatError keeps the mock faithful to the real error
	// shape and avoids the nil-pointer dereference.
	refundResponseBody := map[string]interface{}{
		"Code":    "CommodityNotSupported",
		"Message": "The specified commodity is not supported.",
		"Success": false,
	}
	refundError := tea.NewSDKError(map[string]interface{}{
		"statusCode": 400,
		"code":       tea.ToString(refundResponseBody["Code"]),
		"message":    tea.ToString(refundResponseBody["Message"]),
		"data":       refundResponseBody,
	})

	// Patch DoRequest at the tea-rpc/client.Client seam — the proven-stable point
	//RpcPost/RpcPostWithEndpoint both reach, so a single patch dispatches by action.
	// refundRequest captures the RefundInstance request body so the destroy unit test
	// can assert the resolved billing identity (ProductCode/ProductType from
	// QueryAvailableInstanceByInstanceID), the immediate-release flag and the
	// ClientToken are forwarded — not a hard-coded service name.
	var refundRequest map[string]interface{}
	patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest",
		func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string,
			_ map[string]interface{}, body map[string]interface{},
			_ *util.RuntimeOptions) (map[string]interface{}, error) {
			switch *action {
			case "DescribeDBInstanceAttribute":
				return describeResponse, nil
			case "QueryAvailableInstances":
				return billingResponse, nil
			case "RefundInstance":
				refundRequest = body
				return nil, refundError
			default:
				return nil, fmt.Errorf("unexpected action in unit test: %s", *action)
			}
		})
	defer patches.Reset()

	client := buildClient()

	d := schema.TestResourceDataRaw(t, ResourcePolardbxInstance().Schema, map[string]interface{}{})
	d.SetId(instanceID)
	assert.Equal(t, instanceID, d.Id())

	err := resourcePolardbxInstanceDelete(d, client)
	assert.NotNil(t, err, "Delete must fail for a non-refundable subscription instance")

	// RefundInstance must carry the billing identity resolved from
	// QueryAvailableInstanceByInstanceID (ProductCode/ProductType), the instance
	// ID, the immediate-release flag and a ClientToken — not a hard-coded service
	// name like "polardbx". This is the F-12 request-parameter contract.
	assert.NotNil(t, refundRequest, "RefundInstance request body must be captured")
	assert.Equal(t, instanceID, fmt.Sprint(refundRequest["InstanceId"]), "RefundInstance must carry the instance ID")
	assert.Equal(t, "1", fmt.Sprint(refundRequest["ImmediatelyRelease"]), "RefundInstance must request immediate release")
	assert.Equal(t, "drds", fmt.Sprint(refundRequest["ProductCode"]), "RefundInstance must carry the resolved billing ProductCode")
	assert.Equal(t, "drds_polarxpre_public_cn", fmt.Sprint(refundRequest["ProductType"]), "RefundInstance must carry the resolved billing ProductType")
	assert.NotEmpty(t, fmt.Sprint(refundRequest["ClientToken"]), "RefundInstance must carry a ClientToken")

	msg := err.Error()
	// The actionable error names the instance ID, points to the BSS console
	// unsubscribe page (费用中心 > 退订管理) and the manual state-removal command.
	assert.Contains(t, msg, instanceID, "error must name the instance ID")
	assert.Contains(t, msg, "退订管理", "error must direct to the BSS console unsubscribe page")
	assert.Contains(t, msg, "terraform state rm", "error must give the manual state-removal command")
	// The manual state-removal address is a <name> placeholder, not the
	// description text or a live resource address — F-11.
	assert.Contains(t, msg, "alicloud_drds_polardbx_instance.<name>", "error must use the <name> placeholder address")
	// State is retained: Delete returning an error keeps the resource in state.
	assert.Equal(t, instanceID, d.Id(), "resource Id must be retained when refund is unsupported")

	// The refund failure must not be misreported as NotFound: it is a real,
	// actionable billing rejection.
	assert.False(t, errs.NotFoundError(err), "refund-unsupported error must not be NotFound")
}

// TestUnitDrdsPolardbxInstancePrepayRefundSuccessFalse verifies the AC-3 destroy
// contract for a subscription (PREPAY) PolarDB-X instance whose RefundInstance
// call returns a 200 with Success=false and no SDK error: BSS accepted the
// request but rejected the refund (e.g. ineligible SKU or missing
// bss:RefundBatchRemainRefund permission). Delete must surface an executable
// error and retain the resource in state (Id stays set) instead of silently
// forgetting a still-billing instance. It is the Success=false complement to
// TestUnitDrdsPolardbxInstancePrepayRefundActionableError and a deterministic,
// no-live-billing unit test (DoRequest mocked at the tea-rpc seam).
func TestUnitDrdsPolardbxInstancePrepayRefundSuccessFalse(t *testing.T) {
	buildClient := func() *connectivity.AliyunClient {
		endpoints := &sync.Map{}
		endpoints.Store("polardbx", "polardbx.dummy.aliyuncs.com")
		endpoints.Store("bssopenapi", connectivity.BssOpenAPIEndpointDomestic)
		credConfig := &credential.Config{
			Type:            tea.String("access_key"),
			AccessKeyId:     tea.String("dummy-access-key"),
			AccessKeySecret: tea.String("dummy-secret-key"),
		}
		cred, err := credential.NewCredential(credConfig)
		assert.NoError(t, err)
		cfg := &connectivity.Config{
			Region:               connectivity.Beijing,
			RegionId:             "cn-beijing",
			AccessKey:            "dummy-access-key",
			SecretKey:            "dummy-secret-key",
			AccountType:          "Domestic",
			Protocol:             "HTTPS",
			SkipRegionValidation: true,
			Endpoints:            endpoints,
			SignVersion:          &sync.Map{},
			Credential:           cred,
		}
		aliyunClient, err := cfg.Client()
		assert.NoError(t, err)
		return aliyunClient
	}

	instanceID := "pxc-test-successfalse-instance-id"

	describeResponse := map[string]interface{}{
		"Success": true,
		"Code":    "200",
		"DBInstance": map[string]interface{}{
			"DBInstanceName": instanceID,
			"PayType":        "Prepaid",
			"Status":         "Running",
			"RegionId":       "cn-beijing",
		},
	}

	billingResponse := map[string]interface{}{
		"Success": true,
		"Code":    "200",
		"Data": map[string]interface{}{
			"InstanceList": []interface{}{
				map[string]interface{}{
					"InstanceID":       instanceID,
					"ProductCode":      "drds",
					"ProductType":      "drds_polarxpre_public_cn",
					"SubscriptionType": "Subscription",
				},
			},
		},
	}

	// RefundInstance returns a 200 with Success=false and no SDK error: the SDK
	// call succeeds (nil error) but the refund is rejected at the business level.
	refundResponse := map[string]interface{}{
		"Success": false,
		"Code":    "NotApplicable",
		"Message": "The instance is not eligible for refund.",
	}

	// refundRequest captures the RefundInstance request body to assert the resolved
	// billing identity is forwarded (mirrors the CommodityNotSupported test).
	var refundRequest map[string]interface{}
	patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest",
		func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string,
			_ map[string]interface{}, body map[string]interface{},
			_ *util.RuntimeOptions) (map[string]interface{}, error) {
			switch *action {
			case "DescribeDBInstanceAttribute":
				return describeResponse, nil
			case "QueryAvailableInstances":
				return billingResponse, nil
			case "RefundInstance":
				refundRequest = body
				return refundResponse, nil
			default:
				return nil, fmt.Errorf("unexpected action in unit test: %s", *action)
			}
		})
	defer patches.Reset()

	aliyunClient := buildClient()

	d := schema.TestResourceDataRaw(t, ResourcePolardbxInstance().Schema, map[string]interface{}{})
	d.SetId(instanceID)
	assert.Equal(t, instanceID, d.Id())

	err := resourcePolardbxInstanceDelete(d, aliyunClient)
	assert.NotNil(t, err, "Delete must fail when RefundInstance returns Success=false")

	// The resolved billing identity is still forwarded to RefundInstance.
	assert.NotNil(t, refundRequest, "RefundInstance request body must be captured")
	assert.Equal(t, instanceID, fmt.Sprint(refundRequest["InstanceId"]))
	assert.Equal(t, "1", fmt.Sprint(refundRequest["ImmediatelyRelease"]))
	assert.Equal(t, "drds", fmt.Sprint(refundRequest["ProductCode"]))
	assert.Equal(t, "drds_polarxpre_public_cn", fmt.Sprint(refundRequest["ProductType"]))
	assert.NotEmpty(t, fmt.Sprint(refundRequest["ClientToken"]))

	// State is retained: a Success=false refund must not be misread as success
	// and must not silently forget the still-billing instance.
	assert.Equal(t, instanceID, d.Id(), "resource Id must be retained when refund is rejected with Success=false")
	// The Success=false rejection is a real error, not NotFound.
	assert.False(t, errs.NotFoundError(err), "Success=false rejection must not be NotFound")
}

// TestUnitDrdsQueryAvailableInstanceByInstanceID verifies the helper that
// resolves a subscription instance's billing identity (ProductCode/ProductType)
// solely from its InstanceID via the domestic BSS QueryAvailableInstances API.
// Ported from the root package when alicloud_drds_polardbx_instance was
// migrated; now exercised on drds.Service.
//
// The dummy local AliyunClient is built through connectivity.Config.Client() with a
// dummy access_key credential, a non-empty AccountType (so Client() never probes the
// real account type), SkipRegionValidation=true and pre-stored BSS endpoint (so no
// endpoint discovery happens). DoRequest is mocked at the tea-rpc/client.Client
// boundary — the proven-stable seam — so no real network call is made.
func TestUnitDrdsQueryAvailableInstanceByInstanceID(t *testing.T) {
	buildClient := func(accountType string) *connectivity.AliyunClient {
		endpoints := &sync.Map{}
		if accountType == "International" {
			endpoints.Store("bssopenapi", connectivity.BssOpenAPIEndpointInternational)
		} else {
			endpoints.Store("bssopenapi", connectivity.BssOpenAPIEndpointDomestic)
		}
		credConfig := &credential.Config{
			Type:            tea.String("access_key"),
			AccessKeyId:     tea.String("dummy-access-key"),
			AccessKeySecret: tea.String("dummy-secret-key"),
		}
		cred, err := credential.NewCredential(credConfig)
		assert.NoError(t, err)
		cfg := &connectivity.Config{
			Region:               connectivity.Beijing,
			RegionId:             "cn-beijing",
			AccessKey:            "dummy-access-key",
			SecretKey:            "dummy-secret-key",
			AccountType:          accountType,
			Protocol:             "HTTPS",
			SkipRegionValidation: true,
			Endpoints:            endpoints,
			SignVersion:          &sync.Map{},
			Credential:           cred,
		}
		aliyunClient, err := cfg.Client()
		assert.NoError(t, err)
		return aliyunClient
	}

	instanceID := "pxc-test-instance-id"

	normalResponse := func(id string) map[string]interface{} {
		return map[string]interface{}{
			"Success": true,
			"Code":    "200",
			"Data": map[string]interface{}{
				"InstanceList": []interface{}{
					map[string]interface{}{
						"InstanceID":       id,
						"ProductCode":      "drds",
						"ProductType":      "drds_polarxpre_public_cn",
						"SubscriptionType": "Subscription",
					},
				},
			},
		}
	}

	emptyResponse := map[string]interface{}{
		"Success": true,
		"Code":    "200",
		"Data": map[string]interface{}{
			"InstanceList": []interface{}{},
		},
	}

	patchDoRequest := func(t *testing.T, respond func(action string, body map[string]interface{}) (map[string]interface{}, error)) (*gomonkey.Patches, *string, *map[string]interface{}, *int) {
		var capturedAction string
		var capturedBody map[string]interface{}
		callCount := 0
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest",
			func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string,
				_ map[string]interface{}, body map[string]interface{},
				_ *util.RuntimeOptions) (map[string]interface{}, error) {
				capturedAction = *action
				capturedBody = body
				callCount++
				return respond(*action, body)
			})
		return patches, &capturedAction, &capturedBody, &callCount
	}

	// Domestic account: normal resolution by InstanceIDs only.
	t.Run("DomesticNormal", func(t *testing.T) {
		aliyunClient := buildClient("Domestic")
		assert.False(t, aliyunClient.IsInternationalAccount())
		svc := NewService(aliyunClient)
		patches, capturedAction, capturedBody, _ := patchDoRequest(t, func(action string, _ map[string]interface{}) (map[string]interface{}, error) {
			if action != "QueryAvailableInstances" {
				return nil, fmt.Errorf("unexpected action: %s", action)
			}
			return normalResponse(instanceID), nil
		})
		defer patches.Reset()
		obj, err := svc.QueryAvailableInstanceByInstanceID(instanceID)
		assert.Nil(t, err)
		assert.Equal(t, "QueryAvailableInstances", *capturedAction)
		assert.Equal(t, instanceID, (*capturedBody)["InstanceIDs"])
		// The fix queries solely by InstanceIDs: no ProductCode/ProductType filter.
		_, hasProductCode := (*capturedBody)["ProductCode"]
		assert.False(t, hasProductCode)
		_, hasProductType := (*capturedBody)["ProductType"]
		assert.False(t, hasProductType)
		assert.Equal(t, instanceID, obj["InstanceID"])
		assert.Equal(t, "drds", obj["ProductCode"])
	})

	// Domestic account: empty list → NotFound, with no international-endpoint fallback.
	t.Run("DomesticEmptyListNotFound", func(t *testing.T) {
		aliyunClient := buildClient("Domestic")
		assert.False(t, aliyunClient.IsInternationalAccount())
		svc := NewService(aliyunClient)
		patches, _, _, callCount := patchDoRequest(t, func(action string, _ map[string]interface{}) (map[string]interface{}, error) {
			if action != "QueryAvailableInstances" {
				return nil, fmt.Errorf("unexpected action: %s", action)
			}
			return emptyResponse, nil
		})
		defer patches.Reset()
		_, err := svc.QueryAvailableInstanceByInstanceID(instanceID)
		assert.NotNil(t, err)
		assert.True(t, errs.NotFoundError(err))
		// Exactly one call: an empty domestic list surfaces NotFound instead of
		// blindly falling back to the international endpoint.
		assert.Equal(t, 1, *callCount)
	})

	// Domestic account: returned InstanceID does not match the requested one → NotFound.
	t.Run("DomesticInstanceIDMismatch", func(t *testing.T) {
		aliyunClient := buildClient("Domestic")
		svc := NewService(aliyunClient)
		patches, _, _, _ := patchDoRequest(t, func(action string, _ map[string]interface{}) (map[string]interface{}, error) {
			if action != "QueryAvailableInstances" {
				return nil, fmt.Errorf("unexpected action: %s", action)
			}
			return normalResponse("other-instance-id"), nil
		})
		defer patches.Reset()
		_, err := svc.QueryAvailableInstanceByInstanceID(instanceID)
		assert.NotNil(t, err)
		assert.True(t, errs.NotFoundError(err))
	})

	// International account: IsInternationalAccount() is true and the BSS endpoint is
	// pre-stored as international; resolution still works by InstanceIDs only.
	t.Run("InternationalNormal", func(t *testing.T) {
		aliyunClient := buildClient("International")
		assert.True(t, aliyunClient.IsInternationalAccount())
		svc := NewService(aliyunClient)
		patches, capturedAction, capturedBody, _ := patchDoRequest(t, func(action string, _ map[string]interface{}) (map[string]interface{}, error) {
			if action != "QueryAvailableInstances" {
				return nil, fmt.Errorf("unexpected action: %s", action)
			}
			return normalResponse(instanceID), nil
		})
		defer patches.Reset()
		obj, err := svc.QueryAvailableInstanceByInstanceID(instanceID)
		assert.Nil(t, err)
		assert.Equal(t, "QueryAvailableInstances", *capturedAction)
		assert.Equal(t, instanceID, (*capturedBody)["InstanceIDs"])
		assert.Equal(t, instanceID, obj["InstanceID"])
	})

	// Domestic account: BSS returns Success=false (e.g. permission denied) → the
	// helper must surface a non-NotFound error instead of reading an empty list.
	t.Run("DomesticSuccessFalseRejected", func(t *testing.T) {
		aliyunClient := buildClient("Domestic")
		svc := NewService(aliyunClient)
		rejectedResponse := map[string]interface{}{
			"Success": false,
			"Code":    "NotApplicable",
			"Message": "permission denied",
		}
		patches, _, _, _ := patchDoRequest(t, func(action string, _ map[string]interface{}) (map[string]interface{}, error) {
			if action != "QueryAvailableInstances" {
				return nil, fmt.Errorf("unexpected action: %s", action)
			}
			return rejectedResponse, nil
		})
		defer patches.Reset()
		_, err := svc.QueryAvailableInstanceByInstanceID(instanceID)
		assert.NotNil(t, err)
		assert.False(t, errs.NotFoundError(err))
	})

	// Domestic account: BSS response omits Success entirely → treated as a failure
	// (defensive: a nil Success is never read as success).
	t.Run("DomesticSuccessMissingRejected", func(t *testing.T) {
		aliyunClient := buildClient("Domestic")
		svc := NewService(aliyunClient)
		noSuccessResponse := map[string]interface{}{
			"Data": map[string]interface{}{
				"InstanceList": []interface{}{
					map[string]interface{}{
						"InstanceID":  instanceID,
						"ProductCode": "drds",
					},
				},
			},
		}
		patches, _, _, _ := patchDoRequest(t, func(action string, _ map[string]interface{}) (map[string]interface{}, error) {
			if action != "QueryAvailableInstances" {
				return nil, fmt.Errorf("unexpected action: %s", action)
			}
			return noSuccessResponse, nil
		})
		defer patches.Reset()
		_, err := svc.QueryAvailableInstanceByInstanceID(instanceID)
		assert.NotNil(t, err)
		assert.False(t, errs.NotFoundError(err))
	})
}

// TestUnitDrdsPolardbxInstancePrepayRefundNotFoundDescribeFails verifies the
// F-18 destroy contract for a subscription (PREPAY) PolarDB-X instance whose
// BssOpenApi RefundInstance returns ResourceNotExists (the BSS refund record is
// missing). The preceding DescribePolardbxInstance already proved a Prepaid
// instance exists, so a refund-system NotFound does NOT prove the product is gone.
// Delete must fall through to the product-side release confirmation: when that
// confirmation cannot establish the product is gone — here the post-refund
// Describe query fails with a non-NotFound error (the "product persists or query
// fails" case) — destroy must surface an executable error and retain the resource
// in state (Id stays set) instead of silently returning nil. A plain (non-tea)
// error is used for the poll so NeedRetry is false and DescribePolardbxInstance
// returns immediately rather than retrying, keeping the test deterministic.
// Deterministic, no live billing — DoRequest is mocked at the tea-rpc seam.
func TestUnitDrdsPolardbxInstancePrepayRefundNotFoundDescribeFails(t *testing.T) {
	buildClient := func() *connectivity.AliyunClient {
		endpoints := &sync.Map{}
		endpoints.Store("polardbx", "polardbx.dummy.aliyuncs.com")
		endpoints.Store("bssopenapi", connectivity.BssOpenAPIEndpointDomestic)
		credConfig := &credential.Config{
			Type:            tea.String("access_key"),
			AccessKeyId:     tea.String("dummy-access-key"),
			AccessKeySecret: tea.String("dummy-secret-key"),
		}
		cred, err := credential.NewCredential(credConfig)
		assert.NoError(t, err)
		cfg := &connectivity.Config{
			Region:               connectivity.Beijing,
			RegionId:             "cn-beijing",
			AccessKey:            "dummy-access-key",
			SecretKey:            "dummy-secret-key",
			AccountType:          "Domestic",
			Protocol:             "HTTPS",
			SkipRegionValidation: true,
			Endpoints:            endpoints,
			SignVersion:          &sync.Map{},
			Credential:           cred,
		}
		aliyunClient, err := cfg.Client()
		assert.NoError(t, err)
		return aliyunClient
	}

	instanceID := "pxc-test-refundnotfound-persists"

	// Entry Describe returns the running Prepaid instance so the PREPAY refund
	// path is taken; the post-refund release poll returns a non-NotFound error
	// (a Describe query failure) instead, so the product cannot be confirmed gone.
	describeResponse := map[string]interface{}{
		"Success": true,
		"Code":    "200",
		"DBInstance": map[string]interface{}{
			"DBInstanceName": instanceID,
			"PayType":        "Prepaid",
			"Status":         "Running",
			"RegionId":       "cn-beijing",
		},
	}

	billingResponse := map[string]interface{}{
		"Success": true,
		"Code":    "200",
		"Data": map[string]interface{}{
			"InstanceList": []interface{}{
				map[string]interface{}{
					"InstanceID":       instanceID,
					"ProductCode":      "drds",
					"ProductType":      "drds_polarxpre_public_cn",
					"SubscriptionType": "Subscription",
				},
			},
		},
	}

	// RefundInstance returns the BSS refund-system NotFound (ResourceNotExists):
	// the refund record is missing even though the product still exists.
	refundResponseBody := map[string]interface{}{
		"Code":    "ResourceNotExists",
		"Message": "The resource does not exist.",
		"Success": false,
	}
	refundError := tea.NewSDKError(map[string]interface{}{
		"statusCode": 404,
		"code":       tea.ToString(refundResponseBody["Code"]),
		"message":    tea.ToString(refundResponseBody["Message"]),
		"data":       refundResponseBody,
	})

	var refundRequest map[string]interface{}
	describeCallCount := 0
	patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest",
		func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string,
			_ map[string]interface{}, body map[string]interface{},
			_ *util.RuntimeOptions) (map[string]interface{}, error) {
			switch *action {
			case "DescribeDBInstanceAttribute":
				describeCallCount++
				if describeCallCount == 1 {
					// Entry confirmation: a Prepaid instance exists, so the
					// PREPAY refund path is taken.
					return describeResponse, nil
				}
				// Post-refund release poll: the Describe query fails with a
				// non-NotFound error. The product cannot be confirmed gone, so
				// destroy must surface an error and keep the resource in state.
				// A plain error keeps NeedRetry false so DescribePolardbxInstance
				// returns immediately instead of retrying.
				return nil, fmt.Errorf("DescribeDBInstanceAttribute for %q failed: simulated product-side query error", instanceID)
			case "QueryAvailableInstances":
				return billingResponse, nil
			case "RefundInstance":
				refundRequest = body
				return nil, refundError
			default:
				return nil, fmt.Errorf("unexpected action in unit test: %s", *action)
			}
		})
	defer patches.Reset()

	aliyunClient := buildClient()

	d := schema.TestResourceDataRaw(t, ResourcePolardbxInstance().Schema, map[string]interface{}{})
	d.SetId(instanceID)
	assert.Equal(t, instanceID, d.Id())

	err := resourcePolardbxInstanceDelete(d, aliyunClient)
	assert.NotNil(t, err, "Delete must fail when the product still exists after a refund-system NotFound")

	// The resolved billing identity is still forwarded to RefundInstance.
	assert.NotNil(t, refundRequest, "RefundInstance request body must be captured")
	assert.Equal(t, instanceID, fmt.Sprint(refundRequest["InstanceId"]))
	assert.Equal(t, "1", fmt.Sprint(refundRequest["ImmediatelyRelease"]))
	assert.Equal(t, "drds", fmt.Sprint(refundRequest["ProductCode"]))
	assert.Equal(t, "drds_polarxpre_public_cn", fmt.Sprint(refundRequest["ProductType"]))

	// The release confirmation must poll DescribePolardbxInstance after the
	// refund-system NotFound instead of returning nil immediately (entry + poll).
	assert.GreaterOrEqual(t, describeCallCount, 2, "release confirmation must poll DescribePolardbxInstance after refund NotFound")
	// State retained: the product could not be confirmed gone, so the resource
	// is never silently forgotten.
	assert.Equal(t, instanceID, d.Id(), "resource Id must be retained when the product cannot be confirmed gone after refund NotFound")
	// The error must not be NotFound: the product is not confirmed gone, so
	// destroy must not masquerade as idempotent success.
	assert.False(t, errs.NotFoundError(err), "describe-fails destroy error must not be NotFound")
}

// TestUnitDrdsPolardbxInstancePrepayRefundNotFoundProductGone verifies the F-18
// destroy contract for a subscription (PREPAY) PolarDB-X instance whose
// BssOpenApi RefundInstance returns ResourceNotExists AND the product is
// genuinely gone (the subsequent DescribePolardbxInstance reports NotFound,
// confirmed by the product side). In this case alone destroy may succeed: the
// refund-system NotFound is corroborated by a product-side NotFound, so Delete
// returns nil. Deterministic, no live billing — DoRequest is mocked at the
// tea-rpc seam.
func TestUnitDrdsPolardbxInstancePrepayRefundNotFoundProductGone(t *testing.T) {
	buildClient := func() *connectivity.AliyunClient {
		endpoints := &sync.Map{}
		endpoints.Store("polardbx", "polardbx.dummy.aliyuncs.com")
		endpoints.Store("bssopenapi", connectivity.BssOpenAPIEndpointDomestic)
		credConfig := &credential.Config{
			Type:            tea.String("access_key"),
			AccessKeyId:     tea.String("dummy-access-key"),
			AccessKeySecret: tea.String("dummy-secret-key"),
		}
		cred, err := credential.NewCredential(credConfig)
		assert.NoError(t, err)
		cfg := &connectivity.Config{
			Region:               connectivity.Beijing,
			RegionId:             "cn-beijing",
			AccessKey:            "dummy-access-key",
			SecretKey:            "dummy-secret-key",
			AccountType:          "Domestic",
			Protocol:             "HTTPS",
			SkipRegionValidation: true,
			Endpoints:            endpoints,
			SignVersion:          &sync.Map{},
			Credential:           cred,
		}
		aliyunClient, err := cfg.Client()
		assert.NoError(t, err)
		return aliyunClient
	}

	instanceID := "pxc-test-refundnotfound-gone"

	// Entry Describe returns the running Prepaid instance so the PREPAY refund
	// path is taken. The subsequent release-poll Describe returns
	// DBInstance.NotFound so the product side corroborates the refund-system
	// NotFound.
	describeResponse := map[string]interface{}{
		"Success": true,
		"Code":    "200",
		"DBInstance": map[string]interface{}{
			"DBInstanceName": instanceID,
			"PayType":        "Prepaid",
			"Status":         "Running",
			"RegionId":       "cn-beijing",
		},
	}

	describeNotFoundBody := map[string]interface{}{
		"Code":    "DBInstance.NotFound",
		"Message": "The specified DBInstance does not exist.",
		"Success": false,
	}
	describeNotFoundErr := tea.NewSDKError(map[string]interface{}{
		"statusCode": 404,
		"code":       tea.ToString(describeNotFoundBody["Code"]),
		"message":    tea.ToString(describeNotFoundBody["Message"]),
		"data":       describeNotFoundBody,
	})

	billingResponse := map[string]interface{}{
		"Success": true,
		"Code":    "200",
		"Data": map[string]interface{}{
			"InstanceList": []interface{}{
				map[string]interface{}{
					"InstanceID":       instanceID,
					"ProductCode":      "drds",
					"ProductType":      "drds_polarxpre_public_cn",
					"SubscriptionType": "Subscription",
				},
			},
		},
	}

	refundResponseBody := map[string]interface{}{
		"Code":    "ResourceNotExists",
		"Message": "The resource does not exist.",
		"Success": false,
	}
	refundError := tea.NewSDKError(map[string]interface{}{
		"statusCode": 404,
		"code":       tea.ToString(refundResponseBody["Code"]),
		"message":    tea.ToString(refundResponseBody["Message"]),
		"data":       refundResponseBody,
	})

	var refundRequest map[string]interface{}
	describeCallCount := 0
	patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest",
		func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string,
			_ map[string]interface{}, body map[string]interface{},
			_ *util.RuntimeOptions) (map[string]interface{}, error) {
			switch *action {
			case "DescribeDBInstanceAttribute":
				describeCallCount++
				// First call (entry) confirms the Prepaid instance; later calls
				// (release poll) report product-side NotFound.
				if describeCallCount == 1 {
					return describeResponse, nil
				}
				return nil, describeNotFoundErr
			case "QueryAvailableInstances":
				return billingResponse, nil
			case "RefundInstance":
				refundRequest = body
				return nil, refundError
			default:
				return nil, fmt.Errorf("unexpected action in unit test: %s", *action)
			}
		})
	defer patches.Reset()

	aliyunClient := buildClient()

	d := schema.TestResourceDataRaw(t, ResourcePolardbxInstance().Schema, map[string]interface{}{})
	d.SetId(instanceID)
	assert.Equal(t, instanceID, d.Id())

	err := resourcePolardbxInstanceDelete(d, aliyunClient)
	assert.Nil(t, err, "Delete must succeed when the product is genuinely gone (product-side NotFound corroborates refund-system NotFound)")

	// The resolved billing identity is still forwarded to RefundInstance.
	assert.NotNil(t, refundRequest, "RefundInstance request body must be captured")
	assert.Equal(t, instanceID, fmt.Sprint(refundRequest["InstanceId"]))
	assert.Equal(t, "1", fmt.Sprint(refundRequest["ImmediatelyRelease"]))
	assert.Equal(t, "drds", fmt.Sprint(refundRequest["ProductCode"]))
	assert.Equal(t, "drds_polarxpre_public_cn", fmt.Sprint(refundRequest["ProductType"]))

	// The release confirmation polled DescribePolardbxInstance (entry + poll).
	assert.GreaterOrEqual(t, describeCallCount, 2, "release confirmation must poll DescribePolardbxInstance before succeeding")
}

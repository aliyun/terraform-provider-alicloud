package cr

import (
	"reflect"
	"sync"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/alibabacloud-go/tea-rpc/client"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/credentials-go/credentials"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
)

// TestUnitCrEeInstanceIntlEmptyListSelfHeal covers the site-misjudgement
// self-heal branch added to DescribeInstanceQueryAvailableInstances
// (service_cr_v2.go) for Bug 88792900: an international account whose
// getAccountType() probe cleared AccountType (so IsInternationalAccount() ==
// false) queries BssOpenApi QueryAvailableInstances with
// ProductType=acr_ee_public_cn and receives HTTP 200 Success with an empty
// InstanceList. Before the fix that empty list was misread as NotFound,
// blanking payment_type (Required + ForceNew) and triggering an unwanted
// destroy/rebuild on plan. The self-heal retries once with
// acr_ee_public_intl + the international BssOpenApi endpoint and returns the
// located instance.
//
// A real TestAcc cannot reproduce this branch: it requires an international
// account whose AccountType was cleared by the cfw_pre_intl probe AND a
// CN-product QueryAvailableInstances response that is Success+empty — a
// gateway/account edge case not constructable in the shared acceptance
// account. The gomonkey unit test exercises the branch deterministically.
//
// Mocking notes:
//   - AliyunClient.RpcPostWithEndpoint is a one-line tail call into
//     rpcRequest, so the CI toolchain inlines it and a gomonkey patch on that
//     method is bypassed. The patchable layer is the tea-rpc SDK client's
//     DoRequest — the same layer used by this package's
//     endpoint_acl_policy / chart unit tests — so the mock is applied there.
//     With the minimal AliyunClient (dummy credentials, non-empty AccountType
//     so Client() skips the getAccountType network probe) the real
//     rpcRequest still resolves the BssOpenApi endpoint and builds the SDK
//     client, then the patched DoRequest returns the canned responses.
//   - IsInternationalAccount() is likewise too small to patch reliably, so
//     the CN vs INTL path is selected by setting Config.AccountType ("test"
//     => false, "International" => true) instead of patching the method.
func TestUnitCrEeInstanceIntlEmptyListSelfHeal(t *testing.T) {
	newClient := func(accountType string) *connectivity.AliyunClient {
		credential, err := credentials.NewCredential(new(credentials.Config).
			SetType("access_key").
			SetAccessKeyId("test-key").
			SetAccessKeySecret("test-secret"))
		if err != nil {
			t.Fatalf("failed to create credential: %s", err)
		}
		config := &connectivity.Config{AccessKey: "test-key", SecretKey: "test-secret", Credential: credential,
			RegionId: "cn-hangzhou", AccountType: accountType, Protocol: "http", Endpoints: new(sync.Map),
			SignVersion: new(sync.Map), SkipRegionValidation: true}
		c, err := config.Client()
		if err != nil {
			t.Fatalf("failed to create client: %s", err)
		}
		return c
	}

	p := unitTestProvider().ResourcesMap
	d, _ := schema.InternalMap(p["alicloud_cr_ee_instance"].Schema).Data(nil, nil)
	d.SetId("cri-xmafmag1gmkxi9iq")

	// CN-product response: HTTP 200 Success with an empty InstanceList — the
	// site-misjudgement symptom that previously caused the false NotFound.
	cnEmptyResponse := map[string]interface{}{
		"Data": map[string]interface{}{
			"InstanceList": []interface{}{},
		},
	}
	// INTL-product response: the international retry locates the instance
	// carrying the billing fields Read maps to payment_type / renewal_status /
	// end_time; returning them prevents the ForceNew rebuild.
	intlFoundResponse := map[string]interface{}{
		"Data": map[string]interface{}{
			"InstanceList": []interface{}{
				map[string]interface{}{
					"SubscriptionType": "Subscription",
					"RenewStatus":      "ManualRenewal",
					"EndTime":          "2027-11-11T00:00:00Z",
					"Region":           "eu-central-1",
					"RenewalDuration":  1,
					"CreateTime":       "2026-11-11T00:00:00Z",
				},
			},
		},
	}

	// SelfHealLocatesInstanceViaIntlRetry: the CN query returns Success+empty,
	// the self-heal retries with acr_ee_public_intl + the international
	// endpoint, and the INTL query locates the instance. The method must
	// return the instance (with SubscriptionType) instead of NotFound, so Read
	// keeps payment_type and does not ForceNew.
	t.Run("SelfHealLocatesInstanceViaIntlRetry", func(t *testing.T) {
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest",
			func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string,
				_ map[string]interface{}, body map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
				if body["ProductType"] == "acr_ee_public_cn" {
					return cnEmptyResponse, nil
				}
				return intlFoundResponse, nil
			})
		defer patches.Reset()
		svc := NewCrServiceV2(newClient("test"))
		result, err := svc.DescribeInstanceQueryAvailableInstances(d)
		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "Subscription", result["SubscriptionType"])
		assert.Equal(t, "ManualRenewal", result["RenewStatus"])
	})

	// SelfHealStillEmptyReturnsNotFound: both CN and INTL return Success+empty.
	// The self-heal fires but does not locate the instance; the method must
	// fall through to NotFound (the pre-fix tolerance Read relies on), not
	// panic or blank state silently.
	t.Run("SelfHealStillEmptyReturnsNotFound", func(t *testing.T) {
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest",
			func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string,
				_ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
				return cnEmptyResponse, nil
			})
		defer patches.Reset()
		svc := NewCrServiceV2(newClient("test"))
		_, err := svc.DescribeInstanceQueryAvailableInstances(d)
		assert.NotNil(t, err)
		assert.True(t, errs.NotFoundError(err))
	})

	// IntlAccountNoSelfHeal: when IsInternationalAccount() is true,
	// ProductType is acr_ee_public_intl up front; the self-heal condition
	// (ProductType == acr_ee_public_cn) is not met, so an empty list goes
	// straight to NotFound — correct for a properly-detected international
	// account whose instance is genuinely absent.
	t.Run("IntlAccountNoSelfHeal", func(t *testing.T) {
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest",
			func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string,
				_ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
				return cnEmptyResponse, nil
			})
		defer patches.Reset()
		svc := NewCrServiceV2(newClient("International"))
		_, err := svc.DescribeInstanceQueryAvailableInstances(d)
		assert.NotNil(t, err)
		assert.True(t, errs.NotFoundError(err))
	})
}

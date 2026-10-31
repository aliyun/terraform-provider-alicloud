package bssopenapi_test

import (
	"fmt"
	"reflect"
	"sync"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/alibabacloud-go/tea-rpc/client"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/credentials-go/credentials"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/bssopenapi"
	"github.com/stretchr/testify/assert"
)

func newBssOpenApiTestClient(t *testing.T) *connectivity.AliyunClient {
	t.Helper()
	endpoints := new(sync.Map)
	endpoints.Store("bssopenapi", "bssopenapi.aliyuncs.com")
	credentialConfig := new(credentials.Config).SetType("access_key").SetAccessKeyId("test-access-key").SetAccessKeySecret("test-secret-key")
	testCredential, credErr := credentials.NewCredential(credentialConfig)
	assert.Nil(t, credErr)
	providerConfig := connectivity.Config{
		Region:               connectivity.Region("cn-hangzhou"),
		RegionId:             "cn-hangzhou",
		Protocol:             "HTTPS",
		AccessKey:            "test-access-key",
		SecretKey:            "test-secret-key",
		AccountType:          "Domestic",
		SkipRegionValidation: true,
		Endpoints:            endpoints,
		Credential:           testCredential,
	}
	aliyunClient, err := providerConfig.Client()
	assert.Nil(t, err)
	return aliyunClient
}

// queryResponse builds the BSS response shape the shared helper parses.
func queryResponse(ids ...string) map[string]interface{} {
	list := make([]interface{}, 0, len(ids))
	for _, id := range ids {
		list = append(list, map[string]interface{}{
			"InstanceID":  id,
			"ProductCode": "dts",
			"RenewStatus": "AutoRenewal",
		})
	}
	return map[string]interface{}{
		"Data": map[string]interface{}{
			"InstanceList": list,
		},
	}
}

// runQueryAvailableInstances runs the shared helper against a mocked RPC layer
// and returns the captured request bodies in call order.
func runQueryAvailableInstances(t *testing.T, id, instanceRegion, productCode, productCodeIntl string, respond func(call int, body map[string]interface{}) (map[string]interface{}, error)) (map[string]interface{}, []map[string]interface{}, error) {
	t.Helper()
	aliyunClient := newBssOpenApiTestClient(t)
	calls := make([]map[string]interface{}, 0, 2)
	patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, body map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
		assert.Equal(t, "QueryAvailableInstances", *action)
		captured := make(map[string]interface{}, len(body))
		for key, value := range body {
			captured[key] = value
		}
		calls = append(calls, captured)
		return respond(len(calls), captured)
	})
	defer patches.Reset()

	object, err := bssopenapi.NewService(aliyunClient).QueryAvailableInstances(id, instanceRegion, productCode, productCodeIntl)
	return object, calls, err
}

// TestUnitBssOpenApiQueryAvailableInstances pins the wire contract of the
// shared BSS lookup: the request carries a scalar InstanceIDs filter and the
// product code, the region argument is not forwarded, the response identity
// follows the requested instance, and a missing or mismatched target is
// reported as NotFound instead of another instance's record.
func TestUnitBssOpenApiQueryAvailableInstances(t *testing.T) {
	t.Run("requests the target instance without a region filter", func(t *testing.T) {
		object, calls, err := runQueryAvailableInstances(t, "dts-target-a", "cn-hongkong", "dts", "dts_intl", func(_ int, body map[string]interface{}) (map[string]interface{}, error) {
			return queryResponse(fmt.Sprint(body["InstanceIDs"])), nil
		})
		assert.Nil(t, err)
		assert.Len(t, calls, 1)
		assert.IsType(t, "", calls[0]["InstanceIDs"])
		assert.Equal(t, "dts-target-a", calls[0]["InstanceIDs"])
		assert.Equal(t, "dts", calls[0]["ProductCode"])
		assert.NotContains(t, calls[0], "RegionId")
		assert.NotContains(t, calls[0], "Region")
		for key, value := range calls[0] {
			assert.NotEqual(t, "cn-hongkong", value, "region must not be sent, found in %s", key)
		}
		assert.Equal(t, "dts-target-a", object["InstanceID"])
		assert.Equal(t, "AutoRenewal", object["RenewStatus"])
	})

	t.Run("identity follows the requested instance", func(t *testing.T) {
		objectA, callsA, errA := runQueryAvailableInstances(t, "dts-target-a", "", "dts", "dts_intl", func(_ int, body map[string]interface{}) (map[string]interface{}, error) {
			return queryResponse(fmt.Sprint(body["InstanceIDs"])), nil
		})
		assert.Nil(t, errA)
		assert.Len(t, callsA, 1)
		assert.Equal(t, "dts-target-a", objectA["InstanceID"])

		objectB, callsB, errB := runQueryAvailableInstances(t, "dts-target-b", "", "dts", "dts_intl", func(_ int, body map[string]interface{}) (map[string]interface{}, error) {
			return queryResponse(fmt.Sprint(body["InstanceIDs"])), nil
		})
		assert.Nil(t, errB)
		assert.Len(t, callsB, 1)
		assert.Equal(t, "dts-target-b", objectB["InstanceID"])
	})

	t.Run("rejects a response whose first instance is another object", func(t *testing.T) {
		object, calls, err := runQueryAvailableInstances(t, "dts-target-a", "", "dts", "dts_intl", func(_ int, _ map[string]interface{}) (map[string]interface{}, error) {
			return queryResponse("dts-other"), nil
		})
		assert.Len(t, calls, 1)
		assert.NotNil(t, err)
		assert.True(t, errs.NotFoundError(err))
		assert.Nil(t, object)
	})

	t.Run("falls back to the international endpoint and reports the missing target", func(t *testing.T) {
		object, calls, err := runQueryAvailableInstances(t, "dts-target-a", "", "dts", "dts_intl", func(_ int, _ map[string]interface{}) (map[string]interface{}, error) {
			return queryResponse(), nil
		})
		assert.NotNil(t, err)
		assert.True(t, errs.NotFoundError(err))
		assert.Nil(t, object)
		assert.Len(t, calls, 2)
		assert.Equal(t, "dts", calls[0]["ProductCode"])
		assert.Equal(t, "dts_intl", calls[1]["ProductCode"])
		assert.Equal(t, "dts-target-a", calls[1]["InstanceIDs"])
		assert.NotContains(t, calls[1], "RegionId")
	})
}

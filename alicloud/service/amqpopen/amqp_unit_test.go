package amqpopen

import (
	"testing"

	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
)

// TestUnitAliCloudAmqpInstanceEncryption exercises the encrypted_instance /
// kms_key_id schema flags, the serverless_charge_type validator and the
// addAmqpInstanceEncryptionCreateRequest builder. It migrated from the root
// alicloud package together with the instance resource; the assertions and
// coverage are preserved unchanged. It lives in the internal amqpopen test
// package because addAmqpInstanceEncryptionCreateRequest is unexported.
func TestUnitAliCloudAmqpInstanceEncryption(t *testing.T) {
	resourceSchema := ResourceAmqpInstance()
	assert.True(t, resourceSchema.Schema["encrypted_instance"].Optional)
	assert.True(t, resourceSchema.Schema["encrypted_instance"].Computed)
	assert.True(t, resourceSchema.Schema["encrypted_instance"].ForceNew)
	assert.True(t, resourceSchema.Schema["kms_key_id"].Optional)
	assert.False(t, resourceSchema.Schema["kms_key_id"].Computed)
	assert.True(t, resourceSchema.Schema["kms_key_id"].ForceNew)
	assert.False(t, resourceSchema.Schema["edition"].ForceNew)
	serverlessChargeTypeValidate := resourceSchema.Schema["serverless_charge_type"].ValidateFunc
	assert.NotNil(t, serverlessChargeTypeValidate)
	_, validationErrors := serverlessChargeTypeValidate("onDemand", "serverless_charge_type")
	assert.Empty(t, validationErrors)
	_, validationErrors = serverlessChargeTypeValidate("provisioned", "serverless_charge_type")
	assert.Empty(t, validationErrors)
	_, validationErrors = serverlessChargeTypeValidate("invalid", "serverless_charge_type")
	assert.NotEmpty(t, validationErrors)

	newData := func(values map[string]interface{}) *schema.ResourceData {
		return schema.TestResourceDataRaw(t, resourceSchema.Schema, values)
	}

	for name, values := range map[string]map[string]interface{}{
		"vip": {
			"payment_type":       "Subscription",
			"instance_type":      "vip",
			"encrypted_instance": true,
			"kms_key_id":         "MockKmsKeyId",
		},
		"serverless provisioned dedicated": {
			"payment_type":           "PayAsYouGo",
			"serverless_charge_type": "provisioned",
			"edition":                "dedicated",
			"encrypted_instance":     true,
			"kms_key_id":             "MockKmsKeyId",
		},
	} {
		t.Run("valid "+name, func(t *testing.T) {
			request := make(map[string]interface{})
			assert.NoError(t, addAmqpInstanceEncryptionCreateRequest(newData(values), request))
			assert.Equal(t, true, request["EncryptedInstance"])
			assert.Equal(t, "MockKmsKeyId", request["KmsKeyId"])
		})
	}

	t.Run("encryption arguments omitted", func(t *testing.T) {
		request := make(map[string]interface{})
		assert.NoError(t, addAmqpInstanceEncryptionCreateRequest(newData(map[string]interface{}{
			"payment_type":  "PayAsYouGo",
			"instance_type": "enterprise",
			"edition":       "shared",
		}), request))
		assert.NotContains(t, request, "EncryptedInstance")
		assert.NotContains(t, request, "KmsKeyId")
	})

	t.Run("encryption explicitly disabled", func(t *testing.T) {
		request := make(map[string]interface{})
		assert.NoError(t, addAmqpInstanceEncryptionCreateRequest(newData(map[string]interface{}{
			"payment_type":       "PayAsYouGo",
			"instance_type":      "enterprise",
			"edition":            "shared",
			"encrypted_instance": false,
		}), request))
		assert.Equal(t, false, request["EncryptedInstance"])
		assert.NotContains(t, request, "KmsKeyId")
	})

	for name, testCase := range map[string]struct {
		values map[string]interface{}
		err    string
	}{
		"missing key": {
			values: map[string]interface{}{
				"payment_type":       "Subscription",
				"instance_type":      "vip",
				"encrypted_instance": true,
			},
			err: "kms_key_id must be configured when encrypted_instance is true",
		},
		"blank key": {
			values: map[string]interface{}{
				"payment_type":       "Subscription",
				"instance_type":      "vip",
				"encrypted_instance": true,
				"kms_key_id":         "  ",
			},
			err: "kms_key_id must be configured when encrypted_instance is true",
		},
		"key without encryption": {
			values: map[string]interface{}{
				"payment_type":       "Subscription",
				"instance_type":      "vip",
				"encrypted_instance": false,
				"kms_key_id":         "MockKmsKeyId",
			},
			err: "encrypted_instance must be true when kms_key_id is configured",
		},
	} {
		t.Run("invalid "+name, func(t *testing.T) {
			assert.EqualError(t, addAmqpInstanceEncryptionCreateRequest(newData(testCase.values), make(map[string]interface{})), testCase.err)
		})
	}

	for name, values := range map[string]map[string]interface{}{
		"subscription enterprise dedicated": {
			"payment_type":       "Subscription",
			"instance_type":      "enterprise",
			"edition":            "dedicated",
			"encrypted_instance": true,
			"kms_key_id":         "MockKmsKeyId",
		},
		"pay as you go on demand dedicated": {
			"payment_type":           "PayAsYouGo",
			"serverless_charge_type": "onDemand",
			"edition":                "dedicated",
			"encrypted_instance":     true,
			"kms_key_id":             "MockKmsKeyId",
		},
		"pay as you go provisioned shared": {
			"payment_type":           "PayAsYouGo",
			"serverless_charge_type": "provisioned",
			"edition":                "shared",
			"encrypted_instance":     true,
			"kms_key_id":             "MockKmsKeyId",
		},
		"subscription provisioned dedicated": {
			"payment_type":           "Subscription",
			"serverless_charge_type": "provisioned",
			"edition":                "dedicated",
			"encrypted_instance":     true,
			"kms_key_id":             "MockKmsKeyId",
		},
	} {
		t.Run("passes through "+name, func(t *testing.T) {
			request := make(map[string]interface{})
			assert.NoError(t, addAmqpInstanceEncryptionCreateRequest(newData(values), request))
			assert.Equal(t, true, request["EncryptedInstance"])
			assert.Equal(t, "MockKmsKeyId", request["KmsKeyId"])
		})
	}
}

// TestUnitAliCloudAmqpVirtualHostIsInstanceNotReadyError covers the
// isAmqpInstanceNotReadyError classifier that narrows the "instance still
// initializing" InvalidParameter error to a readiness-retry signal. It
// migrated from the root alicloud package together with the virtual_host
// resource; the assertions and coverage are preserved unchanged. It lives in
// the internal amqpopen test package because isAmqpInstanceNotReadyError is
// unexported.
func TestUnitAliCloudAmqpVirtualHostIsInstanceNotReadyError(t *testing.T) {
	testCases := []struct {
		name     string
		err      error
		expected bool
	}{
		{
			name: "instance still initializing",
			err: &tea.SDKError{
				Code:       helper.String("InvalidParameter"),
				Message:    helper.String("code: 400, the internal error! request id: MockRequestId"),
				StatusCode: tea.Int(400),
			},
			expected: true,
		},
		{
			name: "other InvalidParameter failures still fail fast",
			err: &tea.SDKError{
				Code:       helper.String("InvalidParameter"),
				Message:    helper.String("code: 400, the VirtualHost is invalid. request id: MockRequestId"),
				StatusCode: tea.Int(400),
			},
			expected: false,
		},
		{
			name: "same message under another code is not the readiness signal",
			err: &tea.SDKError{
				Code:       helper.String("InvalidInstanceId.NotFound"),
				Message:    helper.String("code: 400, the internal error! request id: MockRequestId"),
				StatusCode: tea.Int(400),
			},
			expected: false,
		},
		{
			name:     "non sdk error",
			err:      errs.Error("the internal error!"),
			expected: false,
		},
		{
			name:     "nil error",
			err:      nil,
			expected: false,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			assert.Equal(t, testCase.expected, isAmqpInstanceNotReadyError(testCase.err))
		})
	}
}

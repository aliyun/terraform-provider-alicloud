package polardb

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

// KmsService is a minimal local replica of the root-package KmsService
// (alicloud/service_alicloud_kms.go) carrying only the Decrypt method that the
// PolarDB account resource uses for kms_encrypted_password. The product
// package does not import other product packages or the root package; the
// original root definition is retained for unmigrated consumers. Behavior is
// unchanged.
type KmsService struct {
	client *connectivity.AliyunClient
}

func (s *KmsService) Decrypt(ciphertextBlob string, encryptionContext map[string]interface{}) (plaintext string, err error) {
	context, err := json.Marshal(encryptionContext)
	if err != nil {
		return plaintext, WrapError(err)
	}

	var response map[string]interface{}
	client := s.client
	action := "Decrypt"
	request := map[string]interface{}{
		"RegionId":          s.client.RegionId,
		"CiphertextBlob":    ciphertextBlob,
		"EncryptionContext": string(context[:]),
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("Kms", "2016-01-20", action, nil, request, false)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		return plaintext, WrapErrorf(err, DefaultErrorMsg, context, action, AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.Plaintext", response)
	if err != nil {
		return plaintext, WrapErrorf(err, FailedGetAttributeMsg, context, "$.Plaintext", response)
	}

	return fmt.Sprint(v), err
}

package alikafka

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

// KmsService is a local product-package copy of the root alicloud.KmsService
// (service_alicloud_kms.go) limited to Decrypt — the only KMS operation the
// alikafka SASL user resource needs (resource_sasl_user.go). The root type
// cannot be imported without a reverse dependency on the root package; this
// copy preserves the original behavior 1:1 (stock-method-disposition:
// necessary local copy, no behavior change). Root bare aliases are qualified
// to the errs/logs/helper leaf packages they already resolve to.
type KmsService struct {
	client *connectivity.AliyunClient
}

// Decrypt mirrors alicloud.KmsService.Decrypt (service_alicloud_kms.go:124).
// Identical control flow, error wrapping and retry semantics.
func (s *KmsService) Decrypt(ciphertextBlob string, encryptionContext map[string]interface{}) (plaintext string, err error) {
	context, err := json.Marshal(encryptionContext)
	if err != nil {
		return plaintext, errs.WrapError(err)
	}

	var response map[string]interface{}
	client := s.client
	action := "Decrypt"
	request := map[string]interface{}{
		"RegionId":          s.client.RegionId,
		"CiphertextBlob":    ciphertextBlob,
		"EncryptionContext": string(context[:]),
	}
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("Kms", "2016-01-20", action, nil, request, false)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)
	if err != nil {
		return plaintext, errs.WrapErrorf(err, errs.DefaultErrorMsg, context, action, errs.AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.Plaintext", response)
	if err != nil {
		return plaintext, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, context, "$.Plaintext", response)
	}

	return fmt.Sprint(v), err
}

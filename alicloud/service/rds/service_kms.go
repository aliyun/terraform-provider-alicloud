package rds

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

// KmsService is the package-rds copy of the root alicloud/service_alicloud_kms.go
// type. RDS only calls Decrypt, so only that method is copied here. It is a
// cross-product service type; per stock-method-disposition rule #4 the root
// keeps its full type and RDS keeps this behavior-preserving local copy (with
// only the methods RDS uses) rather than importing the root package. The body
// is verbatim from the root; WrapError/NeedRetry/incrementalWait/addDebug and
// the error-code constants resolve via aliases.go.

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

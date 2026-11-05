package cr

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

// decryptKmsPassword decrypts a KMS-encrypted password via the KMS Decrypt API.
// It mirrors KmsService.Decrypt from the root alicloud package (service_alicloud_kms.go)
// verbatim so the cr subpackage keeps the original production behavior without
// importing the kms product service, whose KmsService has not yet been migrated
// to a service subpackage. The cr product only needs the Decrypt operation, so
// the cross-product dependency is localized to this single helper.
func decryptKmsPassword(client *connectivity.AliyunClient, ciphertextBlob string, encryptionContext map[string]interface{}) (string, error) {
	context, err := json.Marshal(encryptionContext)
	if err != nil {
		return "", errs.WrapError(err)
	}

	var response map[string]interface{}
	action := "Decrypt"
	request := map[string]interface{}{
		"RegionId":          client.RegionId,
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
		return "", errs.WrapErrorf(err, errs.DefaultErrorMsg, context, action, errs.AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.Plaintext", response)
	if err != nil {
		return "", errs.WrapErrorf(err, errs.FailedGetAttributeMsg, context, "$.Plaintext", response)
	}

	return fmt.Sprint(v), nil
}

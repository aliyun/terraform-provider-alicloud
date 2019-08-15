package alicloud

import (
	"log"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/kms"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

type KmsService struct {
	client *connectivity.AliyunClient
}

func (k *KmsService) DescribeKmsKey(id string) (*kms.DescribeKeyResponse, error) {
	request := kms.CreateDescribeKeyRequest()
	request.RegionId = k.client.RegionId
	request.KeyId = id

	raw, err := k.client.WithKmsClient(func(kmsClient *kms.Client) (interface{}, error) {
		return kmsClient.DescribeKey(request)
	})
	if err != nil {
		return nil, WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	key, _ := raw.(*kms.DescribeKeyResponse)
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	if key == nil || key.KeyMetadata.KeyId != id {
		return nil, WrapErrorf(Error(GetNotFoundMessage("Kms_key", id)), NotFoundMsg, ProviderERROR)
	}
	if KeyState(key.KeyMetadata.KeyState) == PendingDeletion {
		log.Printf("[WARN] Removing KMS key %s because it's already gone", id)
		return key, WrapErrorf(Error(GetNotFoundMessage("Kms_key", id)), NotFoundMsg, ProviderERROR)
	}
	return key, nil
}

func (k *KmsService) WaitForKmsKey(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := k.DescribeKmsKey(id)
		if err != nil {
			if NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return WrapError(err)
			}
		}
		if object.KeyMetadata.KeyId == id && status != Deleted {
			return nil
		}
		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, id, status, ProviderERROR)
		}
	}
}

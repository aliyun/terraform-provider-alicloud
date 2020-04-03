package alicloud

import (
	"encoding/json"
	"log"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/kms"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

type KmsService struct {
	client *connectivity.AliyunClient
}

func (s *KmsService) DescribeKmsKey(id string) (*kms.DescribeKeyResponse, error) {
	key := &kms.DescribeKeyResponse{}
	request := kms.CreateDescribeKeyRequest()
	request.RegionId = s.client.RegionId
	request.KeyId = id

	raw, err := s.client.WithKmsClient(func(kmsClient *kms.Client) (interface{}, error) {
		return kmsClient.DescribeKey(request)
	})
	if err != nil {
		return key, WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	key, _ = raw.(*kms.DescribeKeyResponse)
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	if key == nil || key.KeyMetadata.KeyId != id {
		return key, WrapErrorf(Error(GetNotFoundMessage("Kms_key", id)), NotFoundMsg, ProviderERROR)
	}
	if KeyState(key.KeyMetadata.KeyState) == PendingDeletion {
		log.Printf("[WARN] Removing KMS key %s because it's already gone", id)
		return key, WrapErrorf(Error(GetNotFoundMessage("Kms_key", id)), NotFoundMsg, ProviderERROR)
	}
	return key, nil
}

func (s *KmsService) WaitForKmsKey(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribeKmsKey(id)
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

func (s *KmsService) Decrypt(ciphertextBlob string, encryptionContext map[string]interface{}) (*kms.DecryptResponse, error) {
	context, err := json.Marshal(encryptionContext)
	if err != nil {
		return nil, WrapError(err)
	}
	request := kms.CreateDecryptRequest()
	request.RegionId = s.client.RegionId
	request.CiphertextBlob = ciphertextBlob
	request.EncryptionContext = string(context[:])
	raw, err := s.client.WithKmsClient(func(kmsClient *kms.Client) (interface{}, error) {
		return kmsClient.Decrypt(request)
	})
	if err != nil {
		return nil, WrapErrorf(err, DefaultErrorMsg, context, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	response, _ := raw.(*kms.DecryptResponse)
	return response, err
}

func (s *KmsService) DescribeKmsSecret(id string) (object kms.DescribeSecretResponse, err error) {
	request := kms.CreateDescribeSecretRequest()
	request.RegionId = s.client.RegionId

	request.SecretName = id
	request.FetchTags = "true"

	raw, err := s.client.WithKmsClient(func(kmsClient *kms.Client) (interface{}, error) {
		return kmsClient.DescribeSecret(request)
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"Forbidden.ResourceNotFound"}) {
			err = WrapErrorf(Error(GetNotFoundMessage("KmsSecret", id)), NotFoundMsg, ProviderERROR)
			return
		}
		err = WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
		return
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*kms.DescribeSecretResponse)
	return *response, nil
}

func (s *KmsService) GetSecretValue(id string) (object kms.GetSecretValueResponse, err error) {
	request := kms.CreateGetSecretValueRequest()
	request.RegionId = s.client.RegionId

	request.SecretName = id

	raw, err := s.client.WithKmsClient(func(kmsClient *kms.Client) (interface{}, error) {
		return kmsClient.GetSecretValue(request)
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"Forbidden.ResourceNotFound"}) {
			err = WrapErrorf(Error(GetNotFoundMessage("kmssecret", id)), NotFoundMsg, ProviderERROR)
			return
		}
		err = WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
		return
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*kms.GetSecretValueResponse)
	return *response, nil
}

func (s *KmsService) setResourceTags(d *schema.ResourceData, resourceType string) error {
	oldItems, newItems := d.GetChange("tags")
	added := make([]JsonTag, 0)
	for key, value := range newItems.(map[string]interface{}) {
		added = append(added, JsonTag{
			TagKey:   key,
			TagValue: value.(string),
		})
	}
	removed := make([]string, 0)
	for key, _ := range oldItems.(map[string]interface{}) {
		removed = append(removed, key)
	}
	if len(removed) > 0 {
		request := kms.CreateUntagResourceRequest()
		request.RegionId = s.client.RegionId
		if resourceType == "kmskey" {
			request.KeyId = d.Id()
		}
		if resourceType == "kmssecret" {
			request.SecretName = d.Id()
		}
		remove, err := json.Marshal(removed)
		if err != nil {
			return WrapError(err)
		}
		request.TagKeys = string(remove)
		raw, err := s.client.WithKmsClient(func(kmsClient *kms.Client) (interface{}, error) {
			return kmsClient.UntagResource(request)
		})
		addDebug(request.GetActionName(), raw)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
	}
	if len(added) > 0 {
		request := kms.CreateTagResourceRequest()
		request.RegionId = s.client.RegionId
		if resourceType == "kmskey" {
			request.KeyId = d.Id()
		}
		if resourceType == "kmssecret" {
			request.SecretName = d.Id()
		}
		add, err := json.Marshal(added)
		if err != nil {
			return WrapError(err)
		}
		request.Tags = string(add)
		raw, err := s.client.WithKmsClient(func(kmsClient *kms.Client) (interface{}, error) {
			return kmsClient.TagResource(request)
		})
		addDebug(request.GetActionName(), raw)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
	}
	return nil
}

func (s *KmsService) DescribeKmsAlias(id string) (object kms.KeyMetadata, err error) {
	request := kms.CreateDescribeKeyRequest()
	request.RegionId = s.client.RegionId

	request.KeyId = id

	raw, err := s.client.WithKmsClient(func(kmsClient *kms.Client) (interface{}, error) {
		return kmsClient.DescribeKey(request)
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"Forbidden.AliasNotFound", "Forbidden.KeyNotFound"}) {
			err = WrapErrorf(Error(GetNotFoundMessage("KmsAlias", id)), NotFoundMsg, ProviderERROR)
			return
		}
		err = WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
		return
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*kms.DescribeKeyResponse)
	return response.KeyMetadata, nil
}

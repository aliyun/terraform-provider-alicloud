package alicloud

import (
	"encoding/json"
	"log"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/kms"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

type KmsService struct {
	client *connectivity.AliyunClient
}

func (s *KmsService) DescribeKmsKey(id string) (object kms.KeyMetadata, err error) {
	request := kms.CreateDescribeKeyRequest()
	request.RegionId = s.client.RegionId

	request.KeyId = id

	raw, err := s.client.WithKmsClient(func(kmsClient *kms.Client) (interface{}, error) {
		return kmsClient.DescribeKey(request)
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"Forbidden.AliasNotFound", "Forbidden.KeyNotFound"}) {
			err = WrapErrorf(Error(GetNotFoundMessage("KmsKey", id)), NotFoundMsg, ProviderERROR)
			return
		}
		err = WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
		return
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*kms.DescribeKeyResponse)
	if response.KeyMetadata.KeyState == "PendingDeletion" {
		log.Printf("[WARN] Removing KmsKey  %s because it's already gone", id)
		return response.KeyMetadata, WrapErrorf(Error(GetNotFoundMessage("KmsKey", id)), NotFoundMsg, ProviderERROR)
	}
	return response.KeyMetadata, nil
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
			return object, err
		}
		err = WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
		return object, err
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
			err = WrapErrorf(Error(GetNotFoundMessage("KmsSecret", id)), NotFoundMsg, ProviderERROR)
			return object, err
		}
		err = WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
		return object, err
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*kms.GetSecretValueResponse)
	return *response, nil
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

func (s *KmsService) DescribeKmsKeyVersion(id string) (object kms.DescribeKeyVersionResponse, err error) {
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		err = WrapError(err)
		return
	}
	request := kms.CreateDescribeKeyVersionRequest()
	request.RegionId = s.client.RegionId
	request.KeyId = parts[0]
	request.KeyVersionId = parts[1]

	raw, err := s.client.WithKmsClient(func(kmsClient *kms.Client) (interface{}, error) {
		return kmsClient.DescribeKeyVersion(request)
	})
	if err != nil {
		err = WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
		return
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*kms.DescribeKeyVersionResponse)
	return *response, nil
}

func (s *KmsService) SetResourceTags(d *schema.ResourceData, resourceType string) error {
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
	// 对系统 Tag 进行过滤
	removedTagKeys := make([]string, 0)
	for _, v := range removed {
		if !ignoredTags(v, "") {
			removedTagKeys = append(removedTagKeys, v)
		}
	}
	if len(removedTagKeys) > 0 {
		request := kms.CreateUntagResourceRequest()
		request.RegionId = s.client.RegionId
		if resourceType == "key" {
			request.KeyId = d.Id()
		}
		if resourceType == "secret" {
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
		if resourceType == "key" {
			request.KeyId = d.Id()
		}
		if resourceType == "secret" {
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

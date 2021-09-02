package alicloud

import (
	"fmt"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

type CloudssoService struct {
	client *connectivity.AliyunClient
}

func (s *CloudssoService) GetDirectory(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewCloudssoClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "GetDirectory"
	request := map[string]interface{}{
		"DirectoryId": id,
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2021-05-15"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		if IsExpectedErrors(err, []string{"EntityNotExists.Directory"}) {
			return object, WrapErrorf(Error(GetNotFoundMessage("CloudSSO:Directory", id)), NotFoundMsg, ProviderERROR, fmt.Sprint(response["RequestId"]))
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.Directory", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Directory", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}
func (s *CloudssoService) GetMFAAuthenticationStatus(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewCloudssoClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "GetMFAAuthenticationStatus"
	request := map[string]interface{}{
		"DirectoryId": id,
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2021-05-15"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		if IsExpectedErrors(err, []string{"EntityNotExists.Directory"}) {
			return object, WrapErrorf(Error(GetNotFoundMessage("CloudSSO:Directory", id)), NotFoundMsg, ProviderERROR, fmt.Sprint(response["RequestId"]))
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}
func (s *CloudssoService) GetSCIMSynchronizationStatus(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewCloudssoClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "GetSCIMSynchronizationStatus"
	request := map[string]interface{}{
		"DirectoryId": id,
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2021-05-15"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		if IsExpectedErrors(err, []string{"EntityNotExists.Directory"}) {
			return object, WrapErrorf(Error(GetNotFoundMessage("CloudSSO:Directory", id)), NotFoundMsg, ProviderERROR, fmt.Sprint(response["RequestId"]))
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}
func (s *CloudssoService) GetDirectoryTasks(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewCloudssoClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "ListTasks"
	request := map[string]interface{}{
		"DirectoryId": id,
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2021-05-15"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		if IsExpectedErrors(err, []string{"EntityNotExists.Directory"}) {
			return object, WrapErrorf(Error(GetNotFoundMessage("CloudSSO:Directory", id)), NotFoundMsg, ProviderERROR, fmt.Sprint(response["RequestId"]))
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Tasks", response)
	}

	object = v.(map[string]interface{})
	return object, nil
}
func (s *CloudssoService) GetExternalSAMLIdentityProvider(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewCloudssoClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "GetExternalSAMLIdentityProvider"
	request := map[string]interface{}{
		"DirectoryId": id,
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2021-05-15"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		if IsExpectedErrors(err, []string{"EntityNotExists.Directory"}) {
			return object, WrapErrorf(Error(GetNotFoundMessage("CloudSSO:Directory", id)), NotFoundMsg, ProviderERROR, fmt.Sprint(response["RequestId"]))
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.SAMLIdentityProviderConfiguration", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.SAMLIdentityProviderConfiguration", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}
func (s *CloudssoService) DescribeCloudSsoDirectory(id string) (object map[string]interface{}, err error) {

	getDirectoryObject, err := s.GetDirectory(id)
	getMFAAuthenticationStatusObject, err := s.GetMFAAuthenticationStatus(id)
	if err != nil {
		return nil, err
	}
	getSCIMSynchronizationStatusObject, err := s.GetSCIMSynchronizationStatus(id)
	if err != nil {
		return nil, err
	}
	getExternalSAMLIdentityProviderObject, err := s.GetExternalSAMLIdentityProvider(id)
	if err != nil {
		return nil, err
	}
	getDirectoryObject["MFAAuthenticationStatus"] = getMFAAuthenticationStatusObject["MFAAuthenticationStatus"]
	getDirectoryObject["SCIMSynchronizationStatus"] = getSCIMSynchronizationStatusObject["SCIMSynchronizationStatus"]
	getDirectoryObject["SAMLIdentityProviderConfiguration"] = getExternalSAMLIdentityProviderObject

	return getDirectoryObject, nil
}

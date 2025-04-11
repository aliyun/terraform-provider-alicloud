package alicloud

import (
	"fmt"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

type CloudSSOServiceV2 struct {
	client *connectivity.AliyunClient
}

// DescribeCloudSSODirectory <<< Encapsulated get interface for CloudSSO Directory.

func (s *CloudSSOServiceV2) DescribeCloudSSODirectory(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["DirectoryId"] = id

	action := "GetDirectorySAMLServiceProviderInfo"

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(1*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("cloudsso", "2021-05-15", action, query, request, true)

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
			return object, WrapErrorf(NotFoundErr("Directory", id), NotFoundMsg, response)
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$.SAMLServiceProvider", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.SAMLServiceProvider", response)
	}

	return v.(map[string]interface{}), nil
}
func (s *CloudSSOServiceV2) DescribeDirectoryGetPasswordPolicy(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["DirectoryId"] = id

	action := "GetPasswordPolicy"

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(1*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("cloudsso", "2021-05-15", action, query, request, true)

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
			return object, WrapErrorf(NotFoundErr("Directory", id), NotFoundMsg, response)
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$.PasswordPolicy", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.PasswordPolicy", response)
	}

	return v.(map[string]interface{}), nil
}
func (s *CloudSSOServiceV2) DescribeDirectoryGetExternalSAMLIdentityProvider(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["DirectoryId"] = id

	action := "GetExternalSAMLIdentityProvider"

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(1*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("cloudsso", "2021-05-15", action, query, request, true)

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
			return object, WrapErrorf(NotFoundErr("Directory", id), NotFoundMsg, response)
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$.SAMLIdentityProviderConfiguration", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.SAMLIdentityProviderConfiguration", response)
	}

	return v.(map[string]interface{}), nil
}
func (s *CloudSSOServiceV2) DescribeDirectoryGetDirectory(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["DirectoryId"] = id

	action := "GetDirectory"

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(1*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("cloudsso", "2021-05-15", action, query, request, true)

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
			return object, WrapErrorf(NotFoundErr("Directory", id), NotFoundMsg, response)
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$.Directory", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Directory", response)
	}

	return v.(map[string]interface{}), nil
}
func (s *CloudSSOServiceV2) DescribeDirectoryGetUserProvisioningConfiguration(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["DirectoryId"] = id

	action := "GetUserProvisioningConfiguration"

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(1*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("cloudsso", "2021-05-15", action, query, request, true)

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
			return object, WrapErrorf(NotFoundErr("Directory", id), NotFoundMsg, response)
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$.UserProvisioningConfiguration", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.UserProvisioningConfiguration", response)
	}

	return v.(map[string]interface{}), nil
}
func (s *CloudSSOServiceV2) DescribeDirectoryGetMFAAuthenticationSettingInfo(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["DirectoryId"] = id

	action := "GetMFAAuthenticationSettingInfo"

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(1*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("cloudsso", "2021-05-15", action, query, request, true)

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
			return object, WrapErrorf(NotFoundErr("Directory", id), NotFoundMsg, response)
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$.MFAAuthenticationSettingInfo", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.MFAAuthenticationSettingInfo", response)
	}

	return v.(map[string]interface{}), nil
}
func (s *CloudSSOServiceV2) DescribeDirectoryGetLoginPreference(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["DirectoryId"] = id

	action := "GetLoginPreference"

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(1*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("cloudsso", "2021-05-15", action, query, request, true)

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
			return object, WrapErrorf(NotFoundErr("Directory", id), NotFoundMsg, response)
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$.LoginPreference", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.LoginPreference", response)
	}

	return v.(map[string]interface{}), nil
}
func (s *CloudSSOServiceV2) DescribeDirectoryGetSCIMSynchronizationStatus(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["DirectoryId"] = id

	action := "GetSCIMSynchronizationStatus"

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(1*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("cloudsso", "2021-05-15", action, query, request, true)

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
			return object, WrapErrorf(NotFoundErr("Directory", id), NotFoundMsg, response)
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}

	return response, nil
}
func (s *CloudSSOServiceV2) DescribeDirectoryGetMFAAuthenticationStatus(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["DirectoryId"] = id

	action := "GetMFAAuthenticationStatus"

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(1*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("cloudsso", "2021-05-15", action, query, request, true)

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
			return object, WrapErrorf(NotFoundErr("Directory", id), NotFoundMsg, response)
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}

	return response, nil
}
func (s *CloudSSOServiceV2) DescribeDirectoryGetDirectoryGlobalAccessStatus(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["DirectoryId"] = id

	action := "GetDirectoryGlobalAccessStatus"

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(1*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("cloudsso", "2021-05-15", action, query, request, true)

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
			return object, WrapErrorf(NotFoundErr("Directory", id), NotFoundMsg, response)
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}

	return response, nil
}

func (s *CloudSSOServiceV2) CloudSSODirectoryStateRefreshFunc(id string, field string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeCloudSSODirectory(id)
		if err != nil {
			if NotFoundError(err) {
				return object, "", nil
			}
			return nil, "", WrapError(err)
		}

		v, err := jsonpath.Get(field, object)
		currentStatus := fmt.Sprint(v)

		if strings.HasPrefix(field, "#") {
			v, _ := jsonpath.Get(strings.TrimPrefix(field, "#"), object)
			if v != nil {
				currentStatus = "#CHECKSET"
			}
		}

		for _, failState := range failStates {
			if currentStatus == failState {
				return object, currentStatus, WrapError(Error(FailedToReachTargetStatus, currentStatus))
			}
		}
		return object, currentStatus, nil
	}
}

// DescribeCloudSSODirectory >>> Encapsulated.

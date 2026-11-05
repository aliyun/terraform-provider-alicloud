package cloudsso

import (
	"fmt"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	errs "github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	helper "github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	logs "github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

type CloudssoService struct {
	client *connectivity.AliyunClient
}

func (s *CloudssoService) GetDirectory(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "GetDirectory"
	request := map[string]interface{}{
		"DirectoryId": id,
	}
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("cloudsso", "2021-05-15", action, nil, request, true)
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
		if errs.IsExpectedErrors(err, []string{"EntityNotExists.Directory"}) {
			return object, errs.WrapErrorf(errs.NotFoundErr("CloudSSO:Directory", id), errs.NotFoundMsg, errs.ProviderERROR, fmt.Sprint(response["RequestId"]))
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.Directory", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.Directory", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}
func (s *CloudssoService) GetMFAAuthenticationStatus(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "GetMFAAuthenticationStatus"
	request := map[string]interface{}{
		"DirectoryId": id,
	}
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("cloudsso", "2021-05-15", action, nil, request, true)
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
		if errs.IsExpectedErrors(err, []string{"EntityNotExists.Directory"}) {
			return object, errs.WrapErrorf(errs.NotFoundErr("CloudSSO:Directory", id), errs.NotFoundMsg, errs.ProviderERROR, fmt.Sprint(response["RequestId"]))
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}
func (s *CloudssoService) GetSCIMSynchronizationStatus(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "GetSCIMSynchronizationStatus"
	request := map[string]interface{}{
		"DirectoryId": id,
	}
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("cloudsso", "2021-05-15", action, nil, request, true)
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
		if errs.IsExpectedErrors(err, []string{"EntityNotExists.Directory"}) {
			return object, errs.WrapErrorf(errs.NotFoundErr("CloudSSO:Directory", id), errs.NotFoundMsg, errs.ProviderERROR, fmt.Sprint(response["RequestId"]))
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}
func (s *CloudssoService) GetDirectoryTasks(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "ListTasks"
	request := map[string]interface{}{
		"DirectoryId": id,
	}
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("cloudsso", "2021-05-15", action, nil, request, true)
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
		if errs.IsExpectedErrors(err, []string{"EntityNotExists.Directory"}) {
			return object, errs.WrapErrorf(errs.NotFoundErr("CloudSSO:Directory", id), errs.NotFoundMsg, errs.ProviderERROR, fmt.Sprint(response["RequestId"]))
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.Tasks", response)
	}

	object = v.(map[string]interface{})
	return object, nil
}
func (s *CloudssoService) GetExternalSAMLIdentityProvider(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "GetExternalSAMLIdentityProvider"
	request := map[string]interface{}{
		"DirectoryId": id,
	}
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("cloudsso", "2021-05-15", action, nil, request, true)
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
		if errs.IsExpectedErrors(err, []string{"EntityNotExists.Directory"}) {
			return object, errs.WrapErrorf(errs.NotFoundErr("CloudSSO:Directory", id), errs.NotFoundMsg, errs.ProviderERROR, fmt.Sprint(response["RequestId"]))
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.SAMLIdentityProviderConfiguration", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.SAMLIdentityProviderConfiguration", response)
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

func (s *CloudssoService) DescribeCloudSsoScimServerCredential(id string) (object map[string]interface{}, err error) {
	client := s.client
	var response map[string]interface{}
	action := "ListSCIMServerCredentials"

	parts, err := helper.ParseResourceId(id, 2)
	if err != nil {
		err = errs.WrapError(err)
		return
	}

	request := map[string]interface{}{
		"DirectoryId": parts[0],
	}

	idExist := false
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("cloudsso", "2021-05-15", action, nil, request, true)
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
		if errs.IsExpectedErrors(err, []string{"EntityNotExists.Directory"}) {
			return object, errs.WrapErrorf(errs.NotFoundErr("CloudSSO:SCIMServerCredential", id), errs.NotFoundMsg, errs.ProviderERROR, fmt.Sprint(response["RequestId"]))
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}

	resp, err := jsonpath.Get("$.SCIMServerCredentials", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.SCIMServerCredentials", response)
	}

	if v, ok := resp.([]interface{}); !ok || len(v) < 1 {
		return object, errs.WrapErrorf(errs.NotFoundErr("CloudSSO:SCIMServerCredential", id), errs.NotFoundWithResponse, response)
	}

	for _, v := range resp.([]interface{}) {
		if fmt.Sprint(v.(map[string]interface{})["CredentialId"]) == parts[1] {
			idExist = true
			return v.(map[string]interface{}), nil
		}
	}

	if !idExist {
		return object, errs.WrapErrorf(errs.NotFoundErr("CloudSSO:SCIMServerCredential", id), errs.NotFoundWithResponse, response)
	}

	return object, nil
}

func (s *CloudssoService) DescribeCloudSsoGroup(id string) (object map[string]interface{}, err error) {
	client := s.client
	parts, err := helper.ParseResourceId(id, 2)
	if err != nil {
		err = errs.WrapError(err)
		return
	}
	request := map[string]interface{}{
		"DirectoryId": parts[0],
		"GroupId":     parts[1],
	}

	var response map[string]interface{}
	action := "GetGroup"
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		resp, err := client.RpcPost("cloudsso", "2021-05-15", action, nil, request, true)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		response = resp
		logs.AddDebug(action, resp, request)
		return nil
	})
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"EntityNotExists.Group"}) {
			return object, errs.WrapErrorf(errs.NotFoundErr("CloudSSO:Group", id), errs.NotFoundMsg, errs.ProviderERROR, fmt.Sprint(response["RequestId"]))
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.Group", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.Group", response)
	}
	return v.(map[string]interface{}), nil
}

func (s *CloudssoService) DescribeCloudSsoUser(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "GetUser"
	parts, err := helper.ParseResourceId(id, 2)
	if err != nil {
		err = errs.WrapError(err)
		return
	}
	request := map[string]interface{}{
		"DirectoryId": parts[0],
		"UserId":      parts[1],
	}
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("cloudsso", "2021-05-15", action, nil, request, true)
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
		if errs.IsExpectedErrors(err, []string{"EntityNotExists.User"}) {
			return object, errs.WrapErrorf(errs.NotFoundErr("CloudSSO:User", id), errs.NotFoundMsg, errs.ProviderERROR, fmt.Sprint(response["RequestId"]))
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.User", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.User", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}

func (s *CloudssoService) ListMFADevicesForUser(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "ListMFADevicesForUser"
	parts, err := helper.ParseResourceId(id, 2)
	if err != nil {
		err = errs.WrapError(err)
		return
	}
	request := map[string]interface{}{
		"DirectoryId": parts[0],
		"UserId":      parts[1],
	}
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("cloudsso", "2021-05-15", action, nil, request, true)
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
		if errs.IsExpectedErrors(err, []string{"EntityNotExists.Directory", "EntityNotExists.User"}) {
			return object, errs.WrapErrorf(errs.NotFoundErr("CloudSSO:User", id), errs.NotFoundMsg, errs.ProviderERROR, fmt.Sprint(response["RequestId"]))
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.MFADevices", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}

func (s *CloudssoService) DescribeCloudSsoAccessConfiguration(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	action := "GetAccessConfiguration"

	client := s.client

	parts, err := helper.ParseResourceId(id, 2)
	if err != nil {
		return nil, errs.WrapError(err)
	}

	request := map[string]interface{}{
		"DirectoryId":           parts[0],
		"AccessConfigurationId": parts[1],
	}

	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("cloudsso", "2021-05-15", action, nil, request, true)
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
		if errs.IsExpectedErrors(err, []string{"EntityNotExists.AccessConfiguration"}) {
			return object, errs.WrapErrorf(errs.NotFoundErr("CloudSSO:AccessConfiguration", id), errs.NotFoundWithResponse, response)
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$.AccessConfiguration", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.AccessConfiguration", response)
	}

	object = v.(map[string]interface{})

	return object, nil
}

func (s *CloudssoService) ListPermissionPoliciesInAccessConfiguration(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	action := "ListPermissionPoliciesInAccessConfiguration"

	client := s.client

	parts, err := helper.ParseResourceId(id, 2)
	if err != nil {
		return nil, errs.WrapError(err)
	}

	request := map[string]interface{}{
		"DirectoryId":           parts[0],
		"AccessConfigurationId": parts[1],
	}

	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("cloudsso", "2021-05-15", action, nil, request, true)
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
		if errs.IsExpectedErrors(err, []string{"EntityNotExists.AccessConfiguration"}) {
			return object, errs.WrapErrorf(errs.NotFoundErr("CloudSSO:AccessConfiguration", id), errs.NotFoundWithResponse, response)
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}

	return response, nil
}

func (s *CloudssoService) DescribeCloudSsoUserAttachment(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "ListGroupMembers"
	parts, err := helper.ParseResourceId(id, 3)
	if err != nil {
		err = errs.WrapError(err)
		return
	}
	request := map[string]interface{}{
		"DirectoryId": parts[0],
		"GroupId":     parts[1],
		"MaxResults":  helper.PageSizeLarge,
	}
	idExist := false
	for {
		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(5*time.Minute, func() *retry.RetryError {

			response, err = client.RpcPost("cloudsso", "2021-05-15", action, nil, request, true)
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
			if errs.IsExpectedErrors(err, []string{"EntityNotExists.Directory", "EntityNotExists.Group"}) {
				return object, errs.WrapErrorf(errs.NotFoundErr("CloudSSO:UserAttachment", id), errs.NotFoundMsg, errs.ProviderERROR, fmt.Sprint(response["RequestId"]))
			}
			return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
		}
		v, err := jsonpath.Get("$.GroupMembers", response)
		if err != nil {
			return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.GroupMembers", response)
		}
		if len(v.([]interface{})) < 1 {
			return object, errs.WrapErrorf(errs.NotFoundErr("CloudSSO", id), errs.NotFoundWithResponse, response)
		}
		for _, v := range v.([]interface{}) {
			if fmt.Sprint(v.(map[string]interface{})["UserId"]) == parts[2] {
				idExist = true
				return v.(map[string]interface{}), nil
			}
		}

		if nextToken, ok := response["NextToken"].(string); ok && nextToken != "" {
			request["NextToken"] = nextToken
		} else {
			break
		}
	}
	if !idExist {
		return object, errs.WrapErrorf(errs.NotFoundErr("CloudSSO", id), errs.NotFoundWithResponse, response)
	}
	return
}

func (s *CloudssoService) DescribeCloudSsoAccessAssignment(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "ListAccessAssignments"
	parts, err := helper.ParseResourceId(id, 6)
	if err != nil {
		err = errs.WrapError(err)
		return
	}
	request := map[string]interface{}{
		"DirectoryId":           parts[0],
		"AccessConfigurationId": parts[1],
		"TargetType":            parts[2],
		"TargetId":              parts[3],
		"PrincipalType":         parts[4],
		"PrincipalId":           parts[5],
	}
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("cloudsso", "2021-05-15", action, nil, request, true)
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
		if errs.IsExpectedErrors(err, []string{"EntityNotExists.Directory"}) {
			return object, errs.WrapErrorf(errs.NotFoundErr("CloudSSO:AccessAttachment", id), errs.NotFoundMsg, errs.ProviderERROR, fmt.Sprint(response["RequestId"]))
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.AccessAssignments", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.AccessAssignments", response)
	}
	if len(v.([]interface{})) < 1 {
		return object, errs.WrapErrorf(errs.NotFoundErr("CloudSSO", id), errs.NotFoundWithResponse, response)
	} else {
		if fmt.Sprint(v.([]interface{})[0].(map[string]interface{})["PrincipalId"]) != parts[5] {
			return object, errs.WrapErrorf(errs.NotFoundErr("CloudSSO", id), errs.NotFoundWithResponse, response)
		}
	}
	object = v.([]interface{})[0].(map[string]interface{})
	return object, nil
}

func (s *CloudssoService) GetTaskStatus(directoryId, taskId string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "GetTaskStatus"
	request := map[string]interface{}{
		"DirectoryId": directoryId,
		"TaskId":      taskId,
	}
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("cloudsso", "2021-05-15", action, nil, request, true)
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
		if errs.IsExpectedErrors(err, []string{"EntityNotExists.Task"}) {
			return object, errs.WrapErrorf(errs.NotFoundErr("CloudSSO", fmt.Sprint(directoryId, ":", taskId)), errs.NotFoundMsg, errs.ProviderERROR, fmt.Sprint(response["RequestId"]))
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, fmt.Sprint(directoryId, ":", taskId), action, errs.AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$.TaskStatus", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, fmt.Sprint(directoryId, ":", taskId), "$.TaskStatus", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}

func (s *CloudssoService) CloudssoServiceAccessAssignmentStateRefreshFunc(directoryId, taskId string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.GetTaskStatus(directoryId, taskId)
		if err != nil {
			if errs.NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", errs.WrapError(err)
		}

		for _, failState := range failStates {
			if fmt.Sprint(object["Status"]) == failState {
				return object, fmt.Sprint(object["Status"]), errs.WrapError(errs.Error(errs.FailedToReachTargetStatus, fmt.Sprintf(directoryId, ":", taskId)))
			}
		}
		return object, fmt.Sprint(object["Status"]), nil
	}
}

func (s *CloudssoService) DescribeCloudSsoAccessConfigurationProvisioning(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "ListAccessConfigurationProvisionings"
	parts, err := helper.ParseResourceId(id, 4)
	if err != nil {
		err = errs.WrapError(err)
		return
	}
	request := map[string]interface{}{
		"AccessConfigurationId": parts[1],
		"DirectoryId":           parts[0],
		"TargetId":              parts[3],
		"TargetType":            parts[2],
	}
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("cloudsso", "2021-05-15", action, nil, request, true)
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
		if errs.IsExpectedErrors(err, []string{"EntityNotExists.AccessConfigurationProvisioning"}) {
			return object, errs.WrapErrorf(errs.NotFoundErr("CloudSSO:AccessConfigurationProvisioning", id), errs.NotFoundMsg, errs.ProviderERROR, fmt.Sprint(response["RequestId"]))
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.AccessConfigurationProvisionings", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.AccessConfigurationProvisionings", response)
	}
	if len(v.([]interface{})) < 1 {
		return object, errs.WrapErrorf(errs.NotFoundErr("CloudSSO", id), errs.NotFoundWithResponse, response)
	} else {
		if fmt.Sprint(v.([]interface{})[0].(map[string]interface{})["TargetId"]) != parts[3] {
			return object, errs.WrapErrorf(errs.NotFoundErr("CloudSSO", id), errs.NotFoundWithResponse, response)
		}
	}
	object = v.([]interface{})[0].(map[string]interface{})
	return object, nil
}

func (s *CloudssoService) CloudssoServiceAccessConfigurationProvisioningStateRefreshFunc(directoryId, taskId string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.GetTaskStatus(directoryId, taskId)
		if err != nil {
			if errs.NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", errs.WrapError(err)
		}

		for _, failState := range failStates {
			if fmt.Sprint(object["Status"]) == failState {
				return object, fmt.Sprint(object["Status"]), errs.WrapError(errs.Error(errs.FailedToReachTargetStatus, fmt.Sprintf(directoryId, ":", taskId)))
			}
		}
		return object, fmt.Sprint(object["Status"]), nil
	}
}

func (s *CloudssoService) CloudssoServicAccessConfigurationProvisioning(directoryId, accessConfigurationId, targetType, targetId string) (err error) {
	var response map[string]interface{}
	action := "ProvisionAccessConfiguration"
	request := make(map[string]interface{})
	client := s.client

	request["DirectoryId"] = directoryId
	request["AccessConfigurationId"] = accessConfigurationId
	request["TargetType"] = targetType
	request["TargetId"] = targetId

	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(10*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("cloudsso", "2021-05-15", action, nil, request, false)
		if err != nil {
			if errs.IsExpectedErrors(err, []string{"OperationConflict.Task"}) || helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)
	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_cloud_sso_access_configuration_provisioning", action, errs.AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.Tasks", response)
	if err != nil || len(v.([]interface{})) < 1 {
		return errs.WrapErrorf(err, errs.IdMsg, fmt.Sprint(directoryId, ":", accessConfigurationId, ":", targetType, ":", targetId))
	}
	response = v.([]interface{})[0].(map[string]interface{})
	_, err = s.GetTaskStatus(fmt.Sprint(request["DirectoryId"]), fmt.Sprint(response["TaskId"]))
	if err != nil {
		if errs.NotFoundError(err) {
			return nil
		}
		return errs.WrapError(err)
	}
	stateConf := helper.BuildStateConf([]string{}, []string{"Success"}, 10*time.Minute, 5*time.Second, s.CloudssoServiceAccessConfigurationProvisioningStateRefreshFunc(fmt.Sprint(request["DirectoryId"]), fmt.Sprint(response["TaskId"]), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errs.WrapErrorf(err, errs.IdMsg, fmt.Sprint(directoryId, ":", accessConfigurationId, ":", targetType, ":", targetId))
	}
	return nil
}

func (s *CloudssoService) DescribeCloudSsoAccessConfigurationProvisionings(directoryId, accessConfigurationId string) (objects []map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client

	action := "ListAccessConfigurationProvisionings"
	request := map[string]interface{}{
		"AccessConfigurationId": accessConfigurationId,
		"DirectoryId":           directoryId,
		"MaxResults":            helper.PageSizeMedium,
	}

	for {
		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(5*time.Minute, func() *retry.RetryError {
			response, err = client.RpcPost("cloudsso", "2021-05-15", action, nil, request, true)
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
			return objects, errs.WrapErrorf(err, errs.DataDefaultErrorMsg, "alicloud_cloud_sso_access_configuration", action, errs.AlibabaCloudSdkGoERROR)
		}
		resp, err := jsonpath.Get("$.AccessConfigurationProvisionings", response)
		if err != nil {
			return objects, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, action, "$.AccessConfigurationProvisionings", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			objects = append(objects, item)
		}

		if nextToken, ok := response["NextToken"].(string); ok && nextToken != "" {
			request["NextToken"] = nextToken
		} else {
			break
		}
	}

	return objects, nil
}

// NewCloudssoService constructs the CloudssoService (cloudsso 2021-05-15 v1
// helpers) from a connectivity client. Used by the SDKv2 factories in this
// package and by the external test package.
func NewCloudssoService(client *connectivity.AliyunClient) *CloudssoService {
	return &CloudssoService{client: client}
}

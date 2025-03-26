package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

type ResourcemanagerService struct {
	client *connectivity.AliyunClient
}

func (s *ResourcemanagerService) DescribeResourceManagerRole(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "GetRole"
	request := map[string]interface{}{
		"RegionId": s.client.RegionId,
		"RoleName": id,
	}
	response, err = client.RpcPost("ResourceManager", "2020-03-31", action, nil, request, true)
	if err != nil {
		if IsExpectedErrors(err, []string{"EntityNotExist.Role"}) {
			err = WrapErrorf(NotFoundErr("ResourceManagerRole", id), NotFoundMsg, ProviderERROR)
			return object, err
		}
		err = WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
		return object, err
	}
	addDebug(action, response, request)
	v, err := jsonpath.Get("$.Role", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Role", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}

func (s *ResourcemanagerService) DescribeResourceManagerResourceGroup(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	action := "GetResourceGroup"

	client := s.client

	request := map[string]interface{}{
		"ResourceGroupId": id,
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("ResourceManager", "2020-03-31", action, nil, request, true)
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
		if IsExpectedErrors(err, []string{"EntityNotExists.ResourceGroup"}) {
			return object, WrapErrorf(NotFoundErr("ResourceManager:ResourceGroup", id), NotFoundWithResponse, response)
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$.ResourceGroup", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.ResourceGroup", response)
	}

	if status, ok := v.(map[string]interface{})["Status"].(string); ok && status == "PendingDelete" {
		log.Printf("[WARN] Removing ResourceManagerResourceGroup  %s because it's already gone", id)
		return v.(map[string]interface{}), WrapErrorf(NotFoundErr("ResourceManager:ResourceGroup", id), NotFoundWithResponse, response)
	}

	object = v.(map[string]interface{})

	return object, nil
}

func (s *ResourcemanagerService) ResourceManagerResourceGroupStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeResourceManagerResourceGroup(id)
		if err != nil {
			if NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}

		for _, failState := range failStates {
			if object["Status"].(string) == failState {
				return object, object["Status"].(string), WrapError(Error(FailedToReachTargetStatus, object["Status"].(string)))
			}
		}

		return object, object["Status"].(string), nil
	}
}

func (s *ResourcemanagerService) DescribeResourceManagerFolder(id string) (object map[string]interface{}, err error) {
	client := s.client
	request := map[string]interface{}{
		"FolderId": id,
	}

	var response map[string]interface{}
	action := "GetFolder"
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		resp, err := client.RpcPost("ResourceManager", "2020-03-31", action, nil, request, false)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		response = resp
		addDebug(action, resp, request)
		return nil
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"EntityNotExists.Folder", "EntityNotExists.ResourceDirectory"}) {
			err = WrapErrorf(NotFoundErr("ResourceManagerFolder", id), NotFoundMsg, ProviderERROR)
			return object, err
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.Folder", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Folder", response)
	}
	return v.(map[string]interface{}), nil
}

func (s *ResourcemanagerService) DescribeResourceManagerHandshake(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "GetHandshake"
	request := map[string]interface{}{
		"RegionId":    s.client.RegionId,
		"HandshakeId": id,
	}
	response, err = client.RpcPost("ResourceManager", "2020-03-31", action, nil, request, true)
	if err != nil {
		if IsExpectedErrors(err, []string{"EntityNotExists.Handshake"}) {
			err = WrapErrorf(NotFoundErr("ResourceManagerHandshake", id), NotFoundMsg, ProviderERROR)
			return object, err
		}
		err = WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
		return object, err
	}
	addDebug(action, response, request)
	v, err := jsonpath.Get("$.Handshake", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Handshake", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}

func (s *ResourcemanagerService) GetPolicyVersion(id string, d *schema.ResourceData) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "GetPolicyVersion"
	request := map[string]interface{}{
		"RegionId":   s.client.RegionId,
		"PolicyName": id,
		"PolicyType": "Custom",
	}
	if v, ok := d.GetOk("default_version"); ok {
		request["VersionId"] = v.(string)
	}
	response, err = client.RpcPost("ResourceManager", "2020-03-31", action, nil, request, true)
	if err != nil {
		if IsExpectedErrors(err, []string{"EntityNotExist.Policy", "EntityNotExist.Policy.Version"}) {
			err = WrapErrorf(NotFoundErr("ResourceManagerPolicy", id), NotFoundMsg, ProviderERROR)
			return object, err
		}
		err = WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
		return object, err
	}
	addDebug(action, response, request)
	v, err := jsonpath.Get("$.PolicyVersion", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.PolicyVersion", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}

func (s *ResourcemanagerService) DescribeResourceManagerPolicy(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "GetPolicy"
	request := map[string]interface{}{
		"RegionId":   s.client.RegionId,
		"PolicyName": id,
		"PolicyType": "Custom",
	}
	response, err = client.RpcPost("ResourceManager", "2020-03-31", action, nil, request, true)
	if err != nil {
		if IsExpectedErrors(err, []string{"EntityNotExist.Policy"}) {
			err = WrapErrorf(NotFoundErr("ResourceManagerPolicy", id), NotFoundMsg, ProviderERROR)
			return object, err
		}
		err = WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
		return object, err
	}
	addDebug(action, response, request)
	v, err := jsonpath.Get("$.Policy", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Policy", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}

func (s *ResourcemanagerService) DescribeResourceManagerAccount(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "GetAccount"
	request := map[string]interface{}{
		"RegionId":  s.client.RegionId,
		"AccountId": id,
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("ResourceManager", "2020-03-31", action, nil, request, true)
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
		if IsExpectedErrors(err, []string{"EntityNotExists.Account", "EntityNotExists.ResourceDirectory"}) {
			return object, WrapErrorf(NotFoundErr("ResourceManager:Account", id), NotFoundMsg, ProviderERROR, fmt.Sprint(response["RequestId"]))
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.Account", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Account", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}

func (s *ResourcemanagerService) GetAccountDeletionStatus(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "GetAccountDeletionStatus"
	request := map[string]interface{}{
		"AccountId": id,
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("ResourceManager", "2020-03-31", action, nil, request, true)
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
		if IsExpectedErrors(err, []string{"EntityNotExists.Account", "EntityNotExists.ResourceDirectory"}) {
			return object, WrapErrorf(NotFoundErr("ResourceManager:Account", id), NotFoundMsg, ProviderERROR, fmt.Sprint(response["RequestId"]))
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.RdAccountDeletionStatus", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.RdAccountDeletionStatus", response)
	}
	object = v.(map[string]interface{})
	v, _ = jsonpath.Get("$.RequestId", response)
	object["RequestId"] = v
	return object, nil
}

func (s *ResourcemanagerService) AccountDeletionStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.GetAccountDeletionStatus(id)
		if err != nil {
			if NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}

		for _, failState := range failStates {

			if fmt.Sprint(object["Status"]) == failState {
				return object, fmt.Sprint(object["Status"]), WrapError(Error(FailedToReachTargetStatusWithRequestId, fmt.Sprint(object["Status"]), fmt.Sprint(object["RequestId"])))
			}
		}
		return object, fmt.Sprint(object["Status"]), nil
	}
}

func (s *ResourcemanagerService) DescribeResourceManagerResourceDirectory(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	action := "GetResourceDirectory"

	client := s.client

	request := map[string]interface{}{}
	response, err = client.RpcPost("ResourceManager", "2020-03-31", action, nil, request, true)
	if err != nil {
		if IsExpectedErrors(err, []string{"ResourceDirectoryNotInUse"}) {
			err = WrapErrorf(NotFoundErr("ResourceManagerResourceDirectory", id), NotFoundMsg, ProviderERROR)
			return object, err
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	addDebug(action, response, request)

	v, err := jsonpath.Get("$.ResourceDirectory", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.ResourceDirectory", response)
	}

	object = v.(map[string]interface{})

	return object, nil
}

func (s *ResourcemanagerService) DescribeResourceManagerPolicyVersion(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "GetPolicyVersion"
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		err = WrapError(err)
		return
	}
	request := map[string]interface{}{
		"RegionId":   s.client.RegionId,
		"PolicyName": parts[0],
		"VersionId":  parts[1],
		"PolicyType": "Custom",
	}
	response, err = client.RpcPost("ResourceManager", "2020-03-31", action, nil, request, true)
	if err != nil {
		if IsExpectedErrors(err, []string{"EntityNotExist.Policy", "EntityNotExist.Policy.Version"}) {
			err = WrapErrorf(NotFoundErr("ResourceManagerPolicyVersion", id), NotFoundMsg, ProviderERROR)
			return object, err
		}
		err = WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
		return object, err
	}
	addDebug(action, response, request)
	v, err := jsonpath.Get("$.PolicyVersion", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.PolicyVersion", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}

func (s *ResourcemanagerService) DescribeResourceManagerPolicyAttachment(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "ListPolicyAttachments"
	parts, err := ParseResourceId(id, 5)
	if err != nil {
		err = WrapError(err)
		return
	}
	request := map[string]interface{}{
		"RegionId":        s.client.RegionId,
		"PolicyName":      parts[0],
		"PolicyType":      parts[1],
		"PrincipalName":   parts[2],
		"PrincipalType":   parts[3],
		"ResourceGroupId": parts[4],
	}
	response, err = client.RpcPost("ResourceManager", "2020-03-31", action, nil, request, true)
	if err != nil {
		if IsExpectedErrors(err, []string{"EntityNotExist.Policy", "EntityNotExists.ResourceGroup"}) {
			err = WrapErrorf(NotFoundErr("ResourceManagerPolicyAttachment", id), NotFoundMsg, ProviderERROR)
			return object, err
		}
		err = WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
		return object, err
	}
	addDebug(action, response, request)
	v, err := jsonpath.Get("$.PolicyAttachments.PolicyAttachment", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.PolicyAttachments.PolicyAttachment", response)
	}
	if len(v.([]interface{})) < 1 {
		return object, WrapErrorf(NotFoundErr("ResourceManager", id), NotFoundWithResponse, response)
	}
	object = v.([]interface{})[0].(map[string]interface{})
	return object, nil
}

func (s *ResourcemanagerService) DescribeResourceManagerControlPolicy(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "GetControlPolicy"
	request := map[string]interface{}{
		"RegionId": s.client.RegionId,
		"PolicyId": id,
	}
	response, err = client.RpcPost("ResourceManager", "2020-03-31", action, nil, request, true)
	if err != nil {
		if IsExpectedErrors(err, []string{"EntityNotExists.ControlPolicy"}) {
			err = WrapErrorf(NotFoundErr("ResourceManagerControlPolicy", id), NotFoundMsg, ProviderERROR)
			return object, err
		}
		err = WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
		return object, err
	}
	addDebug(action, response, request)
	v, err := jsonpath.Get("$.ControlPolicy", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.ControlPolicy", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}

func (s *ResourcemanagerService) DescribeResourceManagerControlPolicyAttachment(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	action := "ListControlPolicyAttachmentsForTarget"

	client := s.client

	parts, err := ParseResourceId(id, 2)
	if err != nil {
		return nil, WrapError(err)
	}

	request := map[string]interface{}{
		"TargetId": parts[1],
	}

	idExist := false
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("ResourceManager", "2020-03-31", action, nil, request, true)
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
		if IsExpectedErrors(err, []string{"EntityNotExists.Target"}) {
			return object, WrapErrorf(NotFoundErr("ResourceManager:ControlPolicyAttachment", id), NotFoundWithResponse, response)
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}

	resp, err := jsonpath.Get("$.ControlPolicyAttachments.ControlPolicyAttachment", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.ControlPolicyAttachments.ControlPolicyAttachment", response)
	}

	if v, ok := resp.([]interface{}); !ok || len(v) < 1 {
		return object, WrapErrorf(NotFoundErr("ResourceManager:ControlPolicyAttachment", id), NotFoundWithResponse, response)
	}

	for _, v := range resp.([]interface{}) {
		if fmt.Sprint(v.(map[string]interface{})["PolicyId"]) == parts[0] {
			idExist = true
			return v.(map[string]interface{}), nil
		}
	}

	if !idExist {
		return object, WrapErrorf(NotFoundErr("ResourceManager:ControlPolicyAttachment", id), NotFoundWithResponse, response)
	}

	return object, nil
}

func (s *ResourcemanagerService) ResourceManagerResourceDirectoryStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeResourceManagerResourceDirectory(id)
		if err != nil {
			if NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}

		for _, failState := range failStates {
			if object["ScpStatus"].(string) == failState {
				return object, object["ScpStatus"].(string), WrapError(Error(FailedToReachTargetStatus, object["ScpStatus"].(string)))
			}
		}
		return object, object["ScpStatus"].(string), nil
	}
}

func (s *ResourcemanagerService) GetPayerForAccount(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "GetPayerForAccount"
	request := map[string]interface{}{
		"RegionId":  s.client.RegionId,
		"AccountId": id,
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("ResourceManager", "2020-03-31", action, nil, request, true)
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
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}

func (s *ResourcemanagerService) ListTagResources(id string, resourceType string) (object interface{}, err error) {
	client := s.client
	action := "ListTagResources"
	request := map[string]interface{}{
		"ResourceType": resourceType,
		"ResourceId.1": id,
	}
	tags := make([]interface{}, 0)
	var response map[string]interface{}

	for {
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			var err error
			response, err = client.RpcPost("ResourceManager", "2020-03-31", action, nil, request, false)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(action, response, request)
			v, err := jsonpath.Get("$.TagResources", response)
			if err != nil {
				return resource.NonRetryableError(WrapErrorf(err, FailedGetAttributeMsg, id, "$.TagResources", response))
			}
			if v != nil {
				tags = append(tags, v.([]interface{})...)
			}
			return nil
		})
		if err != nil {
			err = WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
			return
		}
		if response["NextToken"] == nil {
			break
		}
		request["NextToken"] = response["NextToken"]
	}

	return tags, nil
}

func (s *ResourcemanagerService) SetResourceTags(d *schema.ResourceData, resourceType string) error {

	if d.HasChange("tags") {
		added, removed := parsingTags(d)
		client := s.client
		removedTagKeys := make([]string, 0)
		for _, v := range removed {
			if !ignoredTags(v, "") {
				removedTagKeys = append(removedTagKeys, v)
			}
		}
		if len(removedTagKeys) > 0 {
			action := "UnTagResources"
			request := map[string]interface{}{
				"ResourceType": resourceType,
				"ResourceId.1": d.Id(),
			}
			for i, key := range removedTagKeys {
				request[fmt.Sprintf("TagKey.%d", i+1)] = key
			}
			wait := incrementalWait(2*time.Second, 1*time.Second)
			err := resource.Retry(10*time.Minute, func() *resource.RetryError {
				response, err := client.RpcPost("ResourceManager", "2020-03-31", action, nil, request, false)
				if err != nil {
					if NeedRetry(err) {
						wait()
						return resource.RetryableError(err)

					}
					return resource.NonRetryableError(err)
				}
				addDebug(action, response, request)
				return nil
			})
			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
			}
		}
		if len(added) > 0 {
			action := "TagResources"
			request := map[string]interface{}{
				"ResourceType": resourceType,
				"ResourceId.1": d.Id(),
			}
			count := 1
			for key, value := range added {
				request[fmt.Sprintf("Tag.%d.Key", count)] = key
				request[fmt.Sprintf("Tag.%d.Value", count)] = value
				count++
			}

			wait := incrementalWait(2*time.Second, 1*time.Second)
			err := resource.Retry(10*time.Minute, func() *resource.RetryError {
				response, err := client.RpcPost("ResourceManager", "2020-03-31", action, nil, request, false)
				if err != nil {
					if NeedRetry(err) {
						wait()
						return resource.RetryableError(err)

					}
					return resource.NonRetryableError(err)
				}
				addDebug(action, response, request)
				return nil
			})
			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
			}
		}
		d.SetPartial("tags")
	}
	return nil
}

func (s *ResourcemanagerService) DescribeResourceManagerAccountDeletionCheckTask(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "GetAccountDeletionCheckResult"
	request := map[string]interface{}{
		"AccountId": id,
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("ResourceManager", "2020-03-31", action, nil, request, true)
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
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.AccountDeletionCheckResultInfo", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.AccountDeletionCheckResultInfo", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}

func (s *ResourcemanagerService) ResourceManagerAccountDeletionCheckTaskStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeResourceManagerAccountDeletionCheckTask(id)
		if err != nil {
			if NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}

		for _, failState := range failStates {
			if fmt.Sprint(object["Status"]) == failState {
				return object, fmt.Sprint(object["Status"]), WrapError(Error(FailedToReachTargetStatus, fmt.Sprint(object["Status"])))
			}
		}

		return object, fmt.Sprint(object["Status"]), nil
	}
}

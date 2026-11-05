package resourcemanager

import (
	"fmt"
	"log"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/tags"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type ResourcemanagerService struct {
	client *connectivity.AliyunClient
}

func NewResourcemanagerService(client *connectivity.AliyunClient) *ResourcemanagerService {
	return &ResourcemanagerService{client: client}
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
		if errs.IsExpectedErrors(err, []string{"EntityNotExist.Role"}) {
			err = errs.WrapErrorf(errs.NotFoundErr("ResourceManagerRole", id), errs.NotFoundMsg, errs.ProviderERROR)
			return object, err
		}
		err = errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
		return object, err
	}
	logs.AddDebug(action, response, request)
	v, err := jsonpath.Get("$.Role", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.Role", response)
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

	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("ResourceManager", "2020-03-31", action, nil, request, true)
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
		if errs.IsExpectedErrors(err, []string{"EntityNotExists.ResourceGroup"}) {
			return object, errs.WrapErrorf(errs.NotFoundErr("ResourceManager:ResourceGroup", id), errs.NotFoundWithResponse, response)
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$.ResourceGroup", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.ResourceGroup", response)
	}

	if status, ok := v.(map[string]interface{})["Status"].(string); ok && status == "PendingDelete" {
		log.Printf("[WARN] Removing ResourceManagerResourceGroup  %s because it's already gone", id)
		return v.(map[string]interface{}), errs.WrapErrorf(errs.NotFoundErr("ResourceManager:ResourceGroup", id), errs.NotFoundWithResponse, response)
	}

	object = v.(map[string]interface{})

	return object, nil
}

func (s *ResourcemanagerService) ResourceManagerResourceGroupStateRefreshFunc(id string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeResourceManagerResourceGroup(id)
		if err != nil {
			if errs.NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", errs.WrapError(err)
		}

		for _, failState := range failStates {
			if object["Status"].(string) == failState {
				return object, object["Status"].(string), errs.WrapError(errs.Error(errs.FailedToReachTargetStatus, object["Status"].(string)))
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
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		resp, err := client.RpcPost("ResourceManager", "2020-03-31", action, nil, request, false)
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
		if errs.IsExpectedErrors(err, []string{"EntityNotExists.Folder", "EntityNotExists.ResourceDirectory"}) {
			err = errs.WrapErrorf(errs.NotFoundErr("ResourceManagerFolder", id), errs.NotFoundMsg, errs.ProviderERROR)
			return object, err
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.Folder", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.Folder", response)
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
		if errs.IsExpectedErrors(err, []string{"EntityNotExists.Handshake"}) {
			err = errs.WrapErrorf(errs.NotFoundErr("ResourceManagerHandshake", id), errs.NotFoundMsg, errs.ProviderERROR)
			return object, err
		}
		err = errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
		return object, err
	}
	logs.AddDebug(action, response, request)
	v, err := jsonpath.Get("$.Handshake", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.Handshake", response)
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
		if errs.IsExpectedErrors(err, []string{"EntityNotExist.Policy", "EntityNotExist.Policy.Version"}) {
			err = errs.WrapErrorf(errs.NotFoundErr("ResourceManagerPolicy", id), errs.NotFoundMsg, errs.ProviderERROR)
			return object, err
		}
		err = errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
		return object, err
	}
	logs.AddDebug(action, response, request)
	v, err := jsonpath.Get("$.PolicyVersion", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.PolicyVersion", response)
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
		if errs.IsExpectedErrors(err, []string{"EntityNotExist.Policy"}) {
			err = errs.WrapErrorf(errs.NotFoundErr("ResourceManagerPolicy", id), errs.NotFoundMsg, errs.ProviderERROR)
			return object, err
		}
		err = errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
		return object, err
	}
	logs.AddDebug(action, response, request)
	v, err := jsonpath.Get("$.Policy", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.Policy", response)
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
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("ResourceManager", "2020-03-31", action, nil, request, true)
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
		if errs.IsExpectedErrors(err, []string{"EntityNotExists.Account", "EntityNotExists.ResourceDirectory"}) {
			return object, errs.WrapErrorf(errs.NotFoundErr("ResourceManager:Account", id), errs.NotFoundMsg, errs.ProviderERROR, fmt.Sprint(response["RequestId"]))
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.Account", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.Account", response)
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
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("ResourceManager", "2020-03-31", action, nil, request, true)
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
		if errs.IsExpectedErrors(err, []string{"EntityNotExists.Account", "EntityNotExists.ResourceDirectory"}) {
			return object, errs.WrapErrorf(errs.NotFoundErr("ResourceManager:Account", id), errs.NotFoundMsg, errs.ProviderERROR, fmt.Sprint(response["RequestId"]))
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.RdAccountDeletionStatus", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.RdAccountDeletionStatus", response)
	}
	object = v.(map[string]interface{})
	v, _ = jsonpath.Get("$.RequestId", response)
	object["RequestId"] = v
	return object, nil
}

func (s *ResourcemanagerService) AccountDeletionStateRefreshFunc(id string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.GetAccountDeletionStatus(id)
		if err != nil {
			if errs.NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", errs.WrapError(err)
		}

		for _, failState := range failStates {

			if fmt.Sprint(object["Status"]) == failState {
				return object, fmt.Sprint(object["Status"]), errs.WrapError(errs.Error(errs.FailedToReachTargetStatusWithRequestId, fmt.Sprint(object["Status"]), fmt.Sprint(object["RequestId"])))
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
		if errs.IsExpectedErrors(err, []string{"ResourceDirectoryNotInUse"}) {
			err = errs.WrapErrorf(errs.NotFoundErr("ResourceManagerResourceDirectory", id), errs.NotFoundMsg, errs.ProviderERROR)
			return object, err
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}
	logs.AddDebug(action, response, request)

	v, err := jsonpath.Get("$.ResourceDirectory", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.ResourceDirectory", response)
	}

	object = v.(map[string]interface{})

	return object, nil
}

func (s *ResourcemanagerService) DescribeResourceManagerPolicyVersion(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "GetPolicyVersion"
	parts, err := helper.ParseResourceId(id, 2)
	if err != nil {
		err = errs.WrapError(err)
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
		if errs.IsExpectedErrors(err, []string{"EntityNotExist.Policy", "EntityNotExist.Policy.Version"}) {
			err = errs.WrapErrorf(errs.NotFoundErr("ResourceManagerPolicyVersion", id), errs.NotFoundMsg, errs.ProviderERROR)
			return object, err
		}
		err = errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
		return object, err
	}
	logs.AddDebug(action, response, request)
	v, err := jsonpath.Get("$.PolicyVersion", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.PolicyVersion", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}

func (s *ResourcemanagerService) DescribeResourceManagerPolicyAttachment(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "ListPolicyAttachments"
	parts, err := helper.ParseResourceId(id, 5)
	if err != nil {
		err = errs.WrapError(err)
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
		if errs.IsExpectedErrors(err, []string{"EntityNotExist.Policy", "EntityNotExists.ResourceGroup"}) {
			err = errs.WrapErrorf(errs.NotFoundErr("ResourceManagerPolicyAttachment", id), errs.NotFoundMsg, errs.ProviderERROR)
			return object, err
		}
		err = errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
		return object, err
	}
	logs.AddDebug(action, response, request)
	v, err := jsonpath.Get("$.PolicyAttachments.PolicyAttachment", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.PolicyAttachments.PolicyAttachment", response)
	}
	if len(v.([]interface{})) < 1 {
		return object, errs.WrapErrorf(errs.NotFoundErr("ResourceManager", id), errs.NotFoundWithResponse, response)
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
		if errs.IsExpectedErrors(err, []string{"EntityNotExists.ControlPolicy"}) {
			err = errs.WrapErrorf(errs.NotFoundErr("ResourceManagerControlPolicy", id), errs.NotFoundMsg, errs.ProviderERROR)
			return object, err
		}
		err = errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
		return object, err
	}
	logs.AddDebug(action, response, request)
	v, err := jsonpath.Get("$.ControlPolicy", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.ControlPolicy", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}

func (s *ResourcemanagerService) DescribeResourceManagerControlPolicyAttachment(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	action := "ListControlPolicyAttachmentsForTarget"

	client := s.client

	parts, err := helper.ParseResourceId(id, 2)
	if err != nil {
		return nil, errs.WrapError(err)
	}

	request := map[string]interface{}{
		"TargetId": parts[1],
	}

	idExist := false
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("ResourceManager", "2020-03-31", action, nil, request, true)
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
		if errs.IsExpectedErrors(err, []string{"EntityNotExists.Target"}) {
			return object, errs.WrapErrorf(errs.NotFoundErr("ResourceManager:ControlPolicyAttachment", id), errs.NotFoundWithResponse, response)
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}

	resp, err := jsonpath.Get("$.ControlPolicyAttachments.ControlPolicyAttachment", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.ControlPolicyAttachments.ControlPolicyAttachment", response)
	}

	if v, ok := resp.([]interface{}); !ok || len(v) < 1 {
		return object, errs.WrapErrorf(errs.NotFoundErr("ResourceManager:ControlPolicyAttachment", id), errs.NotFoundWithResponse, response)
	}

	for _, v := range resp.([]interface{}) {
		if fmt.Sprint(v.(map[string]interface{})["PolicyId"]) == parts[0] {
			idExist = true
			return v.(map[string]interface{}), nil
		}
	}

	if !idExist {
		return object, errs.WrapErrorf(errs.NotFoundErr("ResourceManager:ControlPolicyAttachment", id), errs.NotFoundWithResponse, response)
	}

	return object, nil
}

func (s *ResourcemanagerService) ResourceManagerResourceDirectoryStateRefreshFunc(id string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeResourceManagerResourceDirectory(id)
		if err != nil {
			if errs.NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", errs.WrapError(err)
		}

		for _, failState := range failStates {
			if object["ScpStatus"].(string) == failState {
				return object, object["ScpStatus"].(string), errs.WrapError(errs.Error(errs.FailedToReachTargetStatus, object["ScpStatus"].(string)))
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
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("ResourceManager", "2020-03-31", action, nil, request, true)
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
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$", response)
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
		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(5*time.Minute, func() *retry.RetryError {
			var err error
			response, err = client.RpcPost("ResourceManager", "2020-03-31", action, nil, request, false)
			if err != nil {
				if helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			logs.AddDebug(action, response, request)
			v, err := jsonpath.Get("$.TagResources", response)
			if err != nil {
				return retry.NonRetryableError(errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.TagResources", response))
			}
			if v != nil {
				tags = append(tags, v.([]interface{})...)
			}
			return nil
		})
		if err != nil {
			err = errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
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
		added, removed := tags.ParsingTags(d)
		client := s.client
		removedTagKeys := make([]string, 0)
		for _, v := range removed {
			if !tags.IgnoredTags(v, "") {
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
			wait := helper.IncrementalWait(2*time.Second, 1*time.Second)
			err := retry.Retry(10*time.Minute, func() *retry.RetryError {
				response, err := client.RpcPost("ResourceManager", "2020-03-31", action, nil, request, false)
				if err != nil {
					if helper.NeedRetry(err) {
						wait()
						return retry.RetryableError(err)

					}
					return retry.NonRetryableError(err)
				}
				logs.AddDebug(action, response, request)
				return nil
			})
			if err != nil {
				return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
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

			wait := helper.IncrementalWait(2*time.Second, 1*time.Second)
			err := retry.Retry(10*time.Minute, func() *retry.RetryError {
				response, err := client.RpcPost("ResourceManager", "2020-03-31", action, nil, request, false)
				if err != nil {
					if helper.NeedRetry(err) {
						wait()
						return retry.RetryableError(err)

					}
					return retry.NonRetryableError(err)
				}
				logs.AddDebug(action, response, request)
				return nil
			})
			if err != nil {
				return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
			}
		}
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
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("ResourceManager", "2020-03-31", action, nil, request, true)
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
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.AccountDeletionCheckResultInfo", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.AccountDeletionCheckResultInfo", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}

func (s *ResourcemanagerService) ResourceManagerAccountDeletionCheckTaskStateRefreshFunc(id string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeResourceManagerAccountDeletionCheckTask(id)
		if err != nil {
			if errs.NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", errs.WrapError(err)
		}

		for _, failState := range failStates {
			if fmt.Sprint(object["Status"]) == failState {
				return object, fmt.Sprint(object["Status"]), errs.WrapError(errs.Error(errs.FailedToReachTargetStatus, fmt.Sprint(object["Status"])))
			}
		}

		return object, fmt.Sprint(object["Status"]), nil
	}
}

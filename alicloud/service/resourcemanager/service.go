package resourcemanager

import (
	"fmt"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ram"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

type ResourceManagerService struct {
	client *connectivity.AliyunClient
}

func NewResourceManagerService(client *connectivity.AliyunClient) *ResourceManagerService {
	return &ResourceManagerService{client: client}
}

func (s *ResourceManagerService) DescribeResourceManagerDelegatedAdministrator(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "ListDelegatedAdministrators"
	parts, err := helper.ParseResourceId(id, 2)
	if err != nil {
		err = errs.WrapError(err)
		return
	}
	request := map[string]interface{}{
		"ServicePrincipal": parts[1],
		"PageNumber":       1,
		"PageSize":         helper.PageSizeLarge,
	}
	idExist := false
	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
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
		v, err := jsonpath.Get("$.Accounts.Account", response)
		if err != nil {
			return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.Accounts.Account", response)
		}
		if len(v.([]interface{})) < 1 {
			return object, errs.WrapErrorf(errs.NotFoundErr("ResourceManager", id), errs.NotFoundWithResponse, response)
		}
		for _, v := range v.([]interface{}) {
			if fmt.Sprint(v.(map[string]interface{})["AccountId"]) == parts[0] {
				idExist = true
				return v.(map[string]interface{}), nil
			}
		}
		if len(v.([]interface{})) < request["PageSize"].(int) {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}
	if !idExist {
		return object, errs.WrapErrorf(errs.NotFoundErr("ResourceManager", id), errs.NotFoundWithResponse, response)
	}
	return
}

func (s *ResourceManagerService) GetServiceLinkedRoleDeletionStatus(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "GetServiceLinkedRoleDeletionStatus"
	request := map[string]interface{}{
		"DeletionTaskId": id,
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
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

	return response, nil
}

func (s *ResourceManagerService) ResourceManagerServiceLinkedRoleStateRefreshFunc(id string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.GetServiceLinkedRoleDeletionStatus(id)
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

// DescribeRamServiceLinkedRole replicates the root RamService method within the
// resourcemanager package to avoid a circular import on the root package. It calls
// the RAM SDK GetRole via connectivity.AliyunClient.WithRamClient, preserving the
// original retry / error semantics.
func (s *ResourceManagerService) DescribeRamServiceLinkedRole(id string) (*ram.GetRoleResponse, error) {
	parts, err := helper.ParseResourceId(id, 2)
	if err != nil {
		return nil, errs.WrapErrorf(err, "invalid service linked role id: %s", id)
	}
	id = parts[1]

	response := &ram.GetRoleResponse{}
	request := ram.CreateGetRoleRequest()
	request.RegionId = s.client.RegionId
	request.RoleName = id
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		raw, err := s.client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
			return ramClient.GetRole(request)
		})
		if err != nil {
			if errs.IsExpectedErrors(err, []string{errs.ThrottlingUser}) {
				time.Sleep(2 * time.Second)
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)
		response, _ = raw.(*ram.GetRoleResponse)
		return nil
	})
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"EntityNotExist.Role"}) {
			return response, errs.WrapErrorf(err, errs.NotFoundMsg, errs.AlibabaCloudSdkGoERROR)
		}
		return response, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
	}
	return response, nil
}

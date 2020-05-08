package alicloud

import (
	"log"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/resourcemanager"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

type ResourcemanagerService struct {
	client *connectivity.AliyunClient
}

func (s *ResourcemanagerService) DescribeResourceManagerRole(id string) (object resourcemanager.Role, err error) {
	request := resourcemanager.CreateGetRoleRequest()
	request.RegionId = s.client.RegionId

	request.RoleName = id

	raw, err := s.client.WithResourcemanagerClient(func(resourcemanagerClient *resourcemanager.Client) (interface{}, error) {
		return resourcemanagerClient.GetRole(request)
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"EntityNotExist.Role"}) {
			err = WrapErrorf(Error(GetNotFoundMessage("ResourceManagerRole", id)), NotFoundMsg, ProviderERROR)
			return
		}
		err = WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
		return
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*resourcemanager.GetRoleResponse)
	return response.Role, nil
}

func (s *ResourcemanagerService) DescribeResourceManagerResourceGroup(id string) (object resourcemanager.ResourceGroup, err error) {
	request := resourcemanager.CreateGetResourceGroupRequest()
	request.RegionId = s.client.RegionId

	request.ResourceGroupId = id

	raw, err := s.client.WithResourcemanagerClient(func(resourcemanagerClient *resourcemanager.Client) (interface{}, error) {
		return resourcemanagerClient.GetResourceGroup(request)
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"EntityNotExists.ResourceGroup"}) {
			err = WrapErrorf(Error(GetNotFoundMessage("ResourceManagerResourceGroup", id)), NotFoundMsg, ProviderERROR)
			return
		}
		err = WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
		return
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*resourcemanager.GetResourceGroupResponse)
	if response.ResourceGroup.Status == "PendingDelete" {
		log.Printf("[WARN] Removing ResourceManagerResourceGroup  %s because it's already gone", id)
		return response.ResourceGroup, WrapErrorf(Error(GetNotFoundMessage("ResourceManagerResourceGroup", id)), NotFoundMsg, ProviderERROR)
	}
	return response.ResourceGroup, nil
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
			if object.Status == failState {
				return object, object.Status, WrapError(Error(FailedToReachTargetStatus, object.Status))
			}
		}
		return object, object.Status, nil
	}
}

func (s *ResourcemanagerService) DescribeResourceManagerFolder(id string) (object resourcemanager.Folder, err error) {
	request := resourcemanager.CreateGetFolderRequest()
	request.RegionId = s.client.RegionId

	request.FolderId = id

	raw, err := s.client.WithResourcemanagerClient(func(resourcemanagerClient *resourcemanager.Client) (interface{}, error) {
		return resourcemanagerClient.GetFolder(request)
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"EntityNotExists.Folder", "EntityNotExists.ResourceDirectory"}) {
			err = WrapErrorf(Error(GetNotFoundMessage("ResourceManagerFolder", id)), NotFoundMsg, ProviderERROR)
			return
		}
		err = WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
		return
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*resourcemanager.GetFolderResponse)
	return response.Folder, nil
}

func (s *ResourcemanagerService) DescribeResourceManagerHandshake(id string) (object resourcemanager.Handshake, err error) {
	request := resourcemanager.CreateGetHandshakeRequest()
	request.RegionId = s.client.RegionId

	request.HandshakeId = id

	raw, err := s.client.WithResourcemanagerClient(func(resourcemanagerClient *resourcemanager.Client) (interface{}, error) {
		return resourcemanagerClient.GetHandshake(request)
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"EntityNotExists.Handshake"}) {
			err = WrapErrorf(Error(GetNotFoundMessage("ResourceManagerHandshake", id)), NotFoundMsg, ProviderERROR)
			return
		}
		err = WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
		return
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*resourcemanager.GetHandshakeResponse)
	return response.Handshake, nil
}

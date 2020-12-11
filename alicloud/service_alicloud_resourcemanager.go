package alicloud

import (
	"log"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/resourcemanager"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
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

func (s *ResourcemanagerService) GetPolicyVersion(id string, d *schema.ResourceData) (object resourcemanager.PolicyVersion, err error) {
	request := resourcemanager.CreateGetPolicyVersionRequest()
	request.RegionId = s.client.RegionId

	request.PolicyName = id
	if v, ok := d.GetOk("default_version"); ok {
		request.VersionId = v.(string)
	}
	request.PolicyType = "Custom"

	raw, err := s.client.WithResourcemanagerClient(func(resourcemanagerClient *resourcemanager.Client) (interface{}, error) {
		return resourcemanagerClient.GetPolicyVersion(request)
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"EntityNotExist.Policy", "EntityNotExist.Policy.Version"}) {
			err = WrapErrorf(Error(GetNotFoundMessage("ResourceManagerPolicy", id)), NotFoundMsg, ProviderERROR)
			return
		}
		err = WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
		return
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*resourcemanager.GetPolicyVersionResponse)
	return response.PolicyVersion, nil
}

func (s *ResourcemanagerService) DescribeResourceManagerPolicy(id string) (object resourcemanager.Policy, err error) {
	request := resourcemanager.CreateGetPolicyRequest()
	request.RegionId = s.client.RegionId

	request.PolicyName = id
	request.PolicyType = "Custom"

	raw, err := s.client.WithResourcemanagerClient(func(resourcemanagerClient *resourcemanager.Client) (interface{}, error) {
		return resourcemanagerClient.GetPolicy(request)
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"EntityNotExist.Policy"}) {
			err = WrapErrorf(Error(GetNotFoundMessage("ResourceManagerPolicy", id)), NotFoundMsg, ProviderERROR)
			return
		}
		err = WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
		return
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*resourcemanager.GetPolicyResponse)
	return response.Policy, nil
}

func (s *ResourcemanagerService) DescribeResourceManagerAccount(id string) (object resourcemanager.Account, err error) {
	request := resourcemanager.CreateGetAccountRequest()
	request.RegionId = s.client.RegionId

	request.AccountId = id

	raw, err := s.client.WithResourcemanagerClient(func(resourcemanagerClient *resourcemanager.Client) (interface{}, error) {
		return resourcemanagerClient.GetAccount(request)
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"EntityNotExists.Account", "EntityNotExists.ResourceDirectory"}) {
			err = WrapErrorf(Error(GetNotFoundMessage("ResourceManagerAccount", id)), NotFoundMsg, ProviderERROR)
			return
		}
		err = WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
		return
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*resourcemanager.GetAccountResponse)
	return response.Account, nil
}

func (s *ResourcemanagerService) DescribeResourceManagerResourceDirectory(id string) (object resourcemanager.ResourceDirectory, err error) {
	request := resourcemanager.CreateGetResourceDirectoryRequest()
	request.RegionId = s.client.RegionId

	raw, err := s.client.WithResourcemanagerClient(func(resourcemanagerClient *resourcemanager.Client) (interface{}, error) {
		return resourcemanagerClient.GetResourceDirectory(request)
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"ResourceDirectoryNotInUse"}) {
			err = WrapErrorf(Error(GetNotFoundMessage("ResourceManagerResourceDirectory", id)), NotFoundMsg, ProviderERROR)
			return
		}
		err = WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
		return
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*resourcemanager.GetResourceDirectoryResponse)
	return response.ResourceDirectory, nil
}

func (s *ResourcemanagerService) DescribeResourceManagerPolicyVersion(id string) (object resourcemanager.PolicyVersion, err error) {
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		err = WrapError(err)
		return
	}
	request := resourcemanager.CreateGetPolicyVersionRequest()
	request.RegionId = s.client.RegionId
	request.PolicyName = parts[0]
	request.VersionId = parts[1]
	request.PolicyType = "Custom"

	raw, err := s.client.WithResourcemanagerClient(func(resourcemanagerClient *resourcemanager.Client) (interface{}, error) {
		return resourcemanagerClient.GetPolicyVersion(request)
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"EntityNotExist.Policy", "EntityNotExist.Policy.Version"}) {
			err = WrapErrorf(Error(GetNotFoundMessage("ResourceManagerPolicyVersion", id)), NotFoundMsg, ProviderERROR)
			return
		}
		err = WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
		return
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*resourcemanager.GetPolicyVersionResponse)
	return response.PolicyVersion, nil
}

func (s *ResourcemanagerService) DescribeResourceManagerPolicyAttachment(id string) (object resourcemanager.PolicyAttachment, err error) {
	parts, err := ParseResourceId(id, 5)
	if err != nil {
		err = WrapError(err)
		return
	}
	request := resourcemanager.CreateListPolicyAttachmentsRequest()
	request.RegionId = s.client.RegionId
	request.PolicyName = parts[0]
	request.PolicyType = parts[1]
	request.PrincipalName = parts[2]
	request.PrincipalType = parts[3]
	request.ResourceGroupId = parts[4]

	raw, err := s.client.WithResourcemanagerClient(func(resourcemanagerClient *resourcemanager.Client) (interface{}, error) {
		return resourcemanagerClient.ListPolicyAttachments(request)
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"EntityNotExist.Policy", "EntityNotExists.ResourceGroup"}) {
			err = WrapErrorf(Error(GetNotFoundMessage("ResourceManagerPolicyAttachment", id)), NotFoundMsg, ProviderERROR)
			return
		}
		err = WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
		return
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*resourcemanager.ListPolicyAttachmentsResponse)

	if len(response.PolicyAttachments.PolicyAttachment) < 1 {
		err = WrapErrorf(Error(GetNotFoundMessage("ResourceManagerPolicyAttachment", id)), NotFoundMsg, ProviderERROR, response.RequestId)
		return
	}
	return response.PolicyAttachments.PolicyAttachment[0], nil
}

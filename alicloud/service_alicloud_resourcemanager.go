package alicloud

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/services/resourcemanager"
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

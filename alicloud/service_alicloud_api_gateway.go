package alicloud

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/cloudapi"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

type CloudApiService struct {
	client *connectivity.AliyunClient
}

func (s *CloudApiService) DescribeApiGroup(groupId string) (apiGroup *cloudapi.DescribeApiGroupResponse, err error) {
	req := cloudapi.CreateDescribeApiGroupRequest()
	req.GroupId = groupId

	raw, err := s.client.WithCloudApiClient(func(cloudApiClient *cloudapi.Client) (interface{}, error) {
		return cloudApiClient.DescribeApiGroup(req)
	})
	if err != nil {
		if IsExceptedError(err, ApiGroupNotFound) {
			err = GetNotFoundErrorFromString(GetNotFoundMessage("ApiGroup", groupId))
		}
		return
	}
	apiGroup, _ = raw.(*cloudapi.DescribeApiGroupResponse)
	if apiGroup == nil || apiGroup.GroupId == "" {
		err = GetNotFoundErrorFromString(GetNotFoundMessage("ApiGroup", groupId))
	}
	return
}

func (s *CloudApiService) DescribeApp(appId string) (app *cloudapi.DescribeAppResponse, err error) {
	req := cloudapi.CreateDescribeAppRequest()
	req.AppId = requests.Integer(appId)

	raw, err := s.client.WithCloudApiClient(func(cloudApiClient *cloudapi.Client) (interface{}, error) {
		return cloudApiClient.DescribeApp(req)
	})
	if err != nil {
		if IsExceptedError(err, NotFoundApp) {
			err = GetNotFoundErrorFromString(GetNotFoundMessage("App", appId))
		}
		return
	}
	app, _ = raw.(*cloudapi.DescribeAppResponse)
	if app == nil {
		err = GetNotFoundErrorFromString(GetNotFoundMessage("App", appId))
	}
	return
}

func (s *CloudApiService) DescribeApi(apiId string, groupId string) (api *cloudapi.DescribeApiResponse, err error) {
	req := cloudapi.CreateDescribeApiRequest()
	req.ApiId = apiId
	req.GroupId = groupId

	raw, err := s.client.WithCloudApiClient(func(cloudApiClient *cloudapi.Client) (interface{}, error) {
		return cloudApiClient.DescribeApi(req)
	})
	if err != nil {
		if IsExceptedError(err, ApiNotFound) || IsExceptedError(err, ApiGroupNotFound) {
			err = GetNotFoundErrorFromString(GetNotFoundMessage("Api", apiId))
		}
		return
	}
	api, _ = raw.(*cloudapi.DescribeApiResponse)
	if api == nil || api.ApiId == "" {
		err = GetNotFoundErrorFromString(GetNotFoundMessage("Api", apiId))
	}
	return
}

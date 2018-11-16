package alicloud

import (
	"strconv"
	"strings"
	"time"

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
		if IsExceptedErrors(err, []string{ApiGroupNotFound, ApiNotFound}) {
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

func (s *CloudApiService) DescribeAuthorization(id string) (*cloudapi.AuthorizedApp, error) {
	args := cloudapi.CreateDescribeAuthorizedAppsRequest()
	split := strings.Split(id, COLON_SEPARATED)

	args.GroupId = split[0]
	args.ApiId = split[1]
	args.StageName = split[3]
	appId, _ := strconv.Atoi(split[2])

	var allApps []cloudapi.AuthorizedApp

	for {
		raw, err := s.client.WithCloudApiClient(func(cloudApiClient *cloudapi.Client) (interface{}, error) {
			return cloudApiClient.DescribeAuthorizedApps(args)
		})
		if err != nil {
			if IsExceptedErrors(err, []string{ApiGroupNotFound, ApiNotFound}) {
				err = GetNotFoundErrorFromString(GetNotFoundMessage("Authorization", id))
			}
			return nil, err
		}
		resp, _ := raw.(*cloudapi.DescribeAuthorizedAppsResponse)

		if resp == nil {
			break
		}

		allApps = append(allApps, resp.AuthorizedApps.AuthorizedApp...)

		if len(allApps) < PageSizeLarge {
			break
		}

		if page, err := getNextpageNumber(args.PageNumber); err != nil {
			return nil, err
		} else {
			args.PageNumber = page
		}
	}

	var filteredAppsTemp []cloudapi.AuthorizedApp
	for _, app := range allApps {
		if app.AppId == appId {
			filteredAppsTemp = append(filteredAppsTemp, app)
		}

	}

	if len(filteredAppsTemp) < 1 {
		e := GetNotFoundErrorFromString(GetNotFoundMessage("Authorization", id))
		return nil, e
	}
	return &filteredAppsTemp[0], nil
}

func (s *CloudApiService) DescribeVpcAccess(id string) (vpc *cloudapi.VpcAccessAttribute, e error) {
	args := cloudapi.CreateDescribeVpcAccessesRequest()
	split := strings.Split(id, COLON_SEPARATED)

	var allVpcs []cloudapi.VpcAccessAttribute

	for {
		raw, err := s.client.WithCloudApiClient(func(cloudApiClient *cloudapi.Client) (interface{}, error) {
			return cloudApiClient.DescribeVpcAccesses(args)
		})
		if err != nil {
			return nil, err
		}
		resp, _ := raw.(*cloudapi.DescribeVpcAccessesResponse)

		if resp == nil {
			break
		}

		allVpcs = append(allVpcs, resp.VpcAccessAttributes.VpcAccessAttribute...)

		if len(allVpcs) < PageSizeLarge {
			break
		}

		if page, err := getNextpageNumber(args.PageNumber); err != nil {
			return nil, err
		} else {
			args.PageNumber = page
		}
	}

	var filteredVpcsTemp []cloudapi.VpcAccessAttribute
	for _, vpc := range allVpcs {
		iPort, _ := strconv.Atoi(split[3])
		if vpc.Port == iPort && vpc.InstanceId == split[2] && vpc.VpcId == split[1] && vpc.Name == split[0] {
			filteredVpcsTemp = append(filteredVpcsTemp, vpc)
		}
	}

	if len(filteredVpcsTemp) < 1 {
		e = GetNotFoundErrorFromString(GetNotFoundMessage("VPC", id))
		return nil, e
	}

	return &filteredVpcsTemp[0], nil
}

func (s *CloudApiService) WaitForAppAttachmentAuthorization(id string, timeout int) (err error) {
	if timeout <= 0 {
		timeout = DefaultTimeout
	}

	for {
		_, err = s.DescribeAuthorization(id)
		if err == nil {
			break
		}

		timeout = timeout - DefaultIntervalShort
		if timeout <= 0 {
			return GetTimeErrorFromString(GetTimeoutMessage("Authorization", AuthorizationDone))
		}
		time.Sleep(DefaultIntervalShort * time.Second)
	}

	return err
}

func (s *CloudApiService) DescribeDeployedApi(groupId string, apiId string, stageName string) (api *cloudapi.DescribeDeployedApiResponse, err error) {
	req := cloudapi.CreateDescribeDeployedApiRequest()
	req.ApiId = apiId
	req.GroupId = groupId
	req.StageName = stageName

	raw, err := s.client.WithCloudApiClient(func(cloudApiClient *cloudapi.Client) (interface{}, error) {
		return cloudApiClient.DescribeDeployedApi(req)
	})
	if err != nil {
		if IsExceptedErrors(err, []string{ApiGroupNotFound, ApiNotFound, NotFoundStage}) {
			err = GetNotFoundErrorFromString(GetNotFoundMessage("DeployedApi", apiId))
		}
		return
	}
	api, _ = raw.(*cloudapi.DescribeDeployedApiResponse)
	if api == nil {
		err = GetNotFoundErrorFromString(GetNotFoundMessage("DeployedApi", apiId))
	}
	return
}

func (s *CloudApiService) DeployedApi(groupId string, apiId string, stageName string) (err error) {
	req := cloudapi.CreateDeployApiRequest()
	req.ApiId = apiId
	req.GroupId = groupId
	req.StageName = stageName
	req.Description = DeployCommonDescription

	_, err = s.client.WithCloudApiClient(func(cloudApiClient *cloudapi.Client) (interface{}, error) {
		return cloudApiClient.DeployApi(req)
	})

	return
}

func (s *CloudApiService) AbolishApi(groupId string, apiId string, stageName string) (err error) {
	req := cloudapi.CreateAbolishApiRequest()
	req.ApiId = apiId
	req.GroupId = groupId
	req.StageName = stageName

	_, err = s.client.WithCloudApiClient(func(cloudApiClient *cloudapi.Client) (interface{}, error) {
		return cloudApiClient.AbolishApi(req)
	})

	if err != nil {
		if IsExceptedErrors(err, []string{ApiGroupNotFound, ApiNotFound, NotFoundStage}) {
			err = GetNotFoundErrorFromString(GetNotFoundMessage("DeployedApi", apiId))
		}
	}

	return
}

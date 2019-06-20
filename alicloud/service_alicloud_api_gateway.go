package alicloud

import (
	"fmt"
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

func (s *CloudApiService) DescribeApiGatewayGroup(id string) (apiGroup *cloudapi.DescribeApiGroupResponse, err error) {
	request := cloudapi.CreateDescribeApiGroupRequest()
	request.GroupId = id

	raw, err := s.client.WithCloudApiClient(func(cloudApiClient *cloudapi.Client) (interface{}, error) {
		return cloudApiClient.DescribeApiGroup(request)
	})
	if err != nil {
		if IsExceptedError(err, ApiGroupNotFound) {
			return nil, WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
		}
		return nil, WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	addDebug(request.GetActionName(), raw)
	apiGroup, _ = raw.(*cloudapi.DescribeApiGroupResponse)
	if apiGroup.GroupId == "" {
		return nil, WrapErrorf(Error(GetNotFoundMessage("ApiGatewayGroup", id)), NotFoundMsg, ProviderERROR)
	}
	return
}
func (s *CloudApiService) WaitForApiGatewayGroup(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribeApiGatewayGroup(id)
		if err != nil {
			if NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return WrapError(err)
			}
		}
		if string(object.GroupId) == id && status != Deleted {
			return nil
		}
		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, string(object.GroupId), id, ProviderERROR)
		}
	}
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

func (s *CloudApiService) DescribeApiGatewayApi(id string) (api *cloudapi.DescribeApiResponse, err error) {
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		return nil, WrapError(err)
	}
	request := cloudapi.CreateDescribeApiRequest()
	request.ApiId = parts[1]
	request.GroupId = parts[0]

	raw, err := s.client.WithCloudApiClient(func(cloudApiClient *cloudapi.Client) (interface{}, error) {
		return cloudApiClient.DescribeApi(request)
	})
	if err != nil {
		if IsExceptedErrors(err, []string{ApiGroupNotFound, ApiNotFound}) {
			return nil, WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
		}
		return nil, WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw)
	api, _ = raw.(*cloudapi.DescribeApiResponse)
	if api.ApiId == "" {
		return nil, WrapErrorf(Error(GetNotFoundMessage("ApiGatewayApi", id)), NotFoundMsg, ProviderERROR)
	}
	return
}

func (s *CloudApiService) WaitForApiGatewayApi(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribeApiGatewayApi(id)
		if err != nil {
			if NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return WrapError(err)
			}
		}
		respId := fmt.Sprintf("%s%s%s", object.GroupId, COLON_SEPARATED, object.ApiId)
		if respId == id && status != Deleted {
			return nil
		}
		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, respId, id, ProviderERROR)
		}
	}
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

func (s *CloudApiService) DescribeDeployedApi(id string, stageName string) (api *cloudapi.DescribeDeployedApiResponse, err error) {
	request := cloudapi.CreateDescribeDeployedApiRequest()
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		return nil, WrapError(err)
	}
	request.ApiId = parts[1]
	request.GroupId = parts[0]
	request.StageName = stageName

	raw, err := s.client.WithCloudApiClient(func(cloudApiClient *cloudapi.Client) (interface{}, error) {
		return cloudApiClient.DescribeDeployedApi(request)
	})
	if err != nil {
		if IsExceptedErrors(err, []string{ApiGroupNotFound, ApiNotFound, NotFoundStage}) {
			return nil, WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
		}
		return nil, WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw)
	api, _ = raw.(*cloudapi.DescribeDeployedApiResponse)
	return
}

func (s *CloudApiService) DeployedApi(id string, stageName string) (err error) {
	request := cloudapi.CreateDeployApiRequest()
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		return WrapError(err)
	}
	request.ApiId = parts[1]
	request.GroupId = parts[0]
	request.StageName = stageName
	request.Description = DeployCommonDescription

	raw, err := s.client.WithCloudApiClient(func(cloudApiClient *cloudapi.Client) (interface{}, error) {
		return cloudApiClient.DeployApi(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw)
	return
}

func (s *CloudApiService) AbolishApi(id string, stageName string) (err error) {
	request := cloudapi.CreateAbolishApiRequest()
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		return WrapError(err)
	}
	request.ApiId = parts[1]
	request.GroupId = parts[0]
	request.StageName = stageName

	raw, err := s.client.WithCloudApiClient(func(cloudApiClient *cloudapi.Client) (interface{}, error) {
		return cloudApiClient.AbolishApi(request)
	})

	if err != nil {
		if IsExceptedErrors(err, []string{ApiGroupNotFound, ApiNotFound, NotFoundStage}) {
			return WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
		}
		return WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw)
	return
}

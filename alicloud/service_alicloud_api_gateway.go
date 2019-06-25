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

func (s *CloudApiService) DescribeApiGatewayApp(id string) (app *cloudapi.DescribeAppResponse, err error) {
	request := cloudapi.CreateDescribeAppRequest()
	request.AppId = requests.Integer(id)

	raw, err := s.client.WithCloudApiClient(func(cloudApiClient *cloudapi.Client) (interface{}, error) {
		return cloudApiClient.DescribeApp(request)
	})
	if err != nil {
		if IsExceptedError(err, NotFoundApp) {
			return nil, WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
		}
		return nil, WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw)
	app, _ = raw.(*cloudapi.DescribeAppResponse)
	return
}

func (s *CloudApiService) WaitForApiGatewayApp(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribeApiGatewayApp(id)
		if err != nil {
			if NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return WrapError(err)
			}
		}
		if string(object.AppId) == id && status != Deleted {
			return nil
		}
		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, string(object.AppId), id, ProviderERROR)
		}
	}
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

func (s *CloudApiService) DescribeApiGatewayAppAttachment(id string) (*cloudapi.AuthorizedApp, error) {
	request := cloudapi.CreateDescribeAuthorizedAppsRequest()
	parts, err := ParseResourceId(id, 4)
	if err != nil {
		return nil, WrapError(err)
	}
	request.GroupId = parts[0]
	request.ApiId = parts[1]
	request.StageName = parts[3]
	appId, _ := strconv.Atoi(parts[2])

	var allApps []cloudapi.AuthorizedApp

	for {
		raw, err := s.client.WithCloudApiClient(func(cloudApiClient *cloudapi.Client) (interface{}, error) {
			return cloudApiClient.DescribeAuthorizedApps(request)
		})
		if err != nil {
			if IsExceptedErrors(err, []string{ApiGroupNotFound, ApiNotFound}) {
				return nil, WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
			}
			return nil, WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		response, _ := raw.(*cloudapi.DescribeAuthorizedAppsResponse)

		allApps = append(allApps, response.AuthorizedApps.AuthorizedApp...)

		if len(allApps) < PageSizeLarge {
			break
		}

		if page, err := getNextpageNumber(request.PageNumber); err != nil {
			return nil, WrapError(err)
		} else {
			request.PageNumber = page
		}
	}

	var filteredAppsTemp []cloudapi.AuthorizedApp
	for _, app := range allApps {
		if app.AppId == appId {
			filteredAppsTemp = append(filteredAppsTemp, app)
		}

	}

	if len(filteredAppsTemp) < 1 {
		return nil, WrapErrorf(Error(GetNotFoundMessage("ApigatewayAppAttachment", id)), NotFoundMsg, ProviderERROR)
	}
	return &filteredAppsTemp[0], nil
}

func (s *CloudApiService) DescribeApiGatewayVpcAccess(id string) (vpc *cloudapi.VpcAccessAttribute, e error) {
	request := cloudapi.CreateDescribeVpcAccessesRequest()
	parts, err := ParseResourceId(id, 4)
	if err != nil {
		return nil, WrapError(err)
	}
	var allVpcs []cloudapi.VpcAccessAttribute

	for {
		raw, err := s.client.WithCloudApiClient(func(cloudApiClient *cloudapi.Client) (interface{}, error) {
			return cloudApiClient.DescribeVpcAccesses(request)
		})
		if err != nil {
			return nil, WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw)
		response, _ := raw.(*cloudapi.DescribeVpcAccessesResponse)

		allVpcs = append(allVpcs, response.VpcAccessAttributes.VpcAccessAttribute...)

		if len(allVpcs) < PageSizeLarge {
			break
		}

		if page, err := getNextpageNumber(request.PageNumber); err != nil {
			return nil, WrapError(err)
		} else {
			request.PageNumber = page
		}
	}

	var filteredVpcsTemp []cloudapi.VpcAccessAttribute
	for _, vpc := range allVpcs {
		iPort, _ := strconv.Atoi(parts[3])
		if vpc.Port == iPort && vpc.InstanceId == parts[2] && vpc.VpcId == parts[1] && vpc.Name == parts[0] {
			filteredVpcsTemp = append(filteredVpcsTemp, vpc)
		}
	}

	if len(filteredVpcsTemp) < 1 {
		return nil, WrapErrorf(Error(GetNotFoundMessage("ApiGatewayVpcAccess", id)), NotFoundMsg, ProviderERROR)
	}

	return &filteredVpcsTemp[0], nil
}

func (s *CloudApiService) WaitForApiGatewayAppAttachment(id string, status Status, timeout int) (err error) {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	parts, err := ParseResourceId(id, 4)
	if err != nil {
		return WrapError(err)
	}
	appIds := parts[2]
	for {
		object, err := s.DescribeApiGatewayAppAttachment(id)
		if err != nil {
			if NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return WrapError(err)
			}
		}
		if strconv.Itoa(object.AppId) == appIds && status != Deleted {
			return nil
		}
		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, strconv.Itoa(object.AppId), appIds, ProviderERROR)
		}
		time.Sleep(DefaultIntervalShort * time.Second)
	}
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

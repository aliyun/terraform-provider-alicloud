package apigateway

import (
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/tags"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/cloudapi"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type CloudApiService struct {
	client *connectivity.AliyunClient
}

func (s *CloudApiService) DescribeApiGatewayGroup(id string) (*cloudapi.DescribeApiGroupResponse, error) {
	apiGroup := &cloudapi.DescribeApiGroupResponse{}
	request := cloudapi.CreateDescribeApiGroupRequest()
	request.RegionId = s.client.RegionId
	request.GroupId = id

	raw, err := s.client.WithCloudApiClient(func(cloudApiClient *cloudapi.Client) (interface{}, error) {
		return cloudApiClient.DescribeApiGroup(request)
	})
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"NotFoundApiGroup"}) {
			return apiGroup, errs.WrapErrorf(err, errs.NotFoundMsg, errs.AlibabaCloudSdkGoERROR)
		}
		return apiGroup, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
	}

	logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)
	apiGroup, _ = raw.(*cloudapi.DescribeApiGroupResponse)
	if apiGroup.GroupId == "" {
		return apiGroup, errs.WrapErrorf(errs.NotFoundErr("ApiGatewayGroup", id), errs.NotFoundMsg, errs.ProviderERROR)
	}
	return apiGroup, nil
}

func (s *CloudApiService) WaitForApiGatewayGroup(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribeApiGatewayGroup(id)
		if err != nil {
			if errs.NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return errs.WrapError(err)
			}
		}
		if string(object.GroupId) == id && status != Deleted {
			return nil
		}
		if time.Now().After(deadline) {
			return errs.WrapErrorf(err, errs.WaitTimeoutMsg, id, helper.GetFunc(1), timeout, string(object.GroupId), id, errs.ProviderERROR)
		}
	}
}

func (s *CloudApiService) DescribeApiGatewayApp(id string) (*cloudapi.DescribeAppResponse, error) {
	app := &cloudapi.DescribeAppResponse{}
	request := cloudapi.CreateDescribeAppRequest()
	request.RegionId = s.client.RegionId
	request.AppId = requests.Integer(id)

	raw, err := s.client.WithCloudApiClient(func(cloudApiClient *cloudapi.Client) (interface{}, error) {
		return cloudApiClient.DescribeApp(request)
	})
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"NotFoundApp"}) {
			return app, errs.WrapErrorf(err, errs.NotFoundMsg, errs.AlibabaCloudSdkGoERROR)
		}
		return app, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
	}
	logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)
	app, _ = raw.(*cloudapi.DescribeAppResponse)
	return app, nil
}

func (s *CloudApiService) WaitForApiGatewayApp(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribeApiGatewayApp(id)
		if err != nil {
			if errs.NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return errs.WrapError(err)
			}
		}
		if fmt.Sprint(object.AppId) == id && status != Deleted {
			return nil
		}
		if time.Now().After(deadline) {
			return errs.WrapErrorf(err, errs.WaitTimeoutMsg, id, helper.GetFunc(1), timeout, fmt.Sprint(object.AppId), id, errs.ProviderERROR)
		}
	}
}

func (s *CloudApiService) DescribeApiGatewayApi(id string) (*cloudapi.DescribeApiResponse, error) {
	api := &cloudapi.DescribeApiResponse{}
	parts, err := helper.ParseResourceId(id, 2)
	if err != nil {
		return api, errs.WrapError(err)
	}
	request := cloudapi.CreateDescribeApiRequest()
	request.RegionId = s.client.RegionId
	request.ApiId = parts[1]
	request.GroupId = parts[0]

	raw, err := s.client.WithCloudApiClient(func(cloudApiClient *cloudapi.Client) (interface{}, error) {
		return cloudApiClient.DescribeApi(request)
	})
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"NotFoundApiGroup", "NotFoundApi"}) {
			return api, errs.WrapErrorf(err, errs.NotFoundMsg, errs.AlibabaCloudSdkGoERROR)
		}
		return api, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
	}
	logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)
	api, _ = raw.(*cloudapi.DescribeApiResponse)
	if api.ApiId == "" {
		return api, errs.WrapErrorf(errs.NotFoundErr("ApiGatewayApi", id), errs.NotFoundMsg, errs.ProviderERROR)
	}
	return api, nil
}

func (s *CloudApiService) WaitForApiGatewayApi(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribeApiGatewayApi(id)
		if err != nil {
			if errs.NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			}
			return errs.WrapError(err)
		}
		respId := fmt.Sprintf("%s%s%s", object.GroupId, COLON_SEPARATED, object.ApiId)
		if respId == id && status != Deleted {
			return nil
		}
		if time.Now().After(deadline) {
			return errs.WrapErrorf(err, errs.WaitTimeoutMsg, id, helper.GetFunc(1), timeout, respId, id, errs.ProviderERROR)
		}
	}
}

func (s *CloudApiService) DescribeApiGatewayAppAttachment(id string) (*cloudapi.AuthorizedApp, error) {
	app := &cloudapi.AuthorizedApp{}
	request := cloudapi.CreateDescribeAuthorizedAppsRequest()
	request.RegionId = s.client.RegionId
	parts, err := helper.ParseResourceId(id, 4)
	if err != nil {
		return app, errs.WrapError(err)
	}
	request.GroupId = parts[0]
	request.ApiId = parts[1]
	request.StageName = parts[3]
	appId, _ := strconv.ParseInt(parts[2], 10, 64)

	var allApps []cloudapi.AuthorizedApp

	for {
		raw, err := s.client.WithCloudApiClient(func(cloudApiClient *cloudapi.Client) (interface{}, error) {
			return cloudApiClient.DescribeAuthorizedApps(request)
		})
		if err != nil {
			if errs.IsExpectedErrors(err, []string{"NotFoundApiGroup", "NotFoundApi"}) {
				return app, errs.WrapErrorf(err, errs.NotFoundMsg, errs.AlibabaCloudSdkGoERROR)
			}
			return app, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
		}
		logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)
		response, _ := raw.(*cloudapi.DescribeAuthorizedAppsResponse)

		allApps = append(allApps, response.AuthorizedApps.AuthorizedApp...)

		if len(allApps) < helper.PageSizeLarge {
			break
		}

		page, err := helper.GetNextpageNumber(request.PageNumber)
		if err != nil {
			return app, errs.WrapError(err)
		}
		request.PageNumber = page
	}

	var filteredAppsTemp []cloudapi.AuthorizedApp
	for _, app := range allApps {
		if app.AppId == appId {
			filteredAppsTemp = append(filteredAppsTemp, app)
		}
	}

	if len(filteredAppsTemp) < 1 {
		return app, errs.WrapErrorf(errs.NotFoundErr("ApigatewayAppAttachment", id), errs.NotFoundMsg, errs.ProviderERROR)
	}
	app = &filteredAppsTemp[0]
	return app, nil
}

func (s *CloudApiService) DescribeApiGatewayPluginAttachment(id string) (*cloudapi.PluginAttribute, error) {
	plugin := &cloudapi.PluginAttribute{}
	request := cloudapi.CreateDescribePluginsByApiRequest()
	request.RegionId = s.client.RegionId
	parts, err := helper.ParseResourceId(id, 4)
	if err != nil {
		return plugin, errs.WrapError(err)
	}
	request.GroupId = parts[0]
	request.ApiId = parts[1]
	request.StageName = parts[3]
	pluginId := parts[2]

	var allPlugins []cloudapi.PluginAttribute

	for {
		raw, err := s.client.WithCloudApiClient(func(cloudApiClient *cloudapi.Client) (interface{}, error) {
			return cloudApiClient.DescribePluginsByApi(request)
		})
		if err != nil {
			if errs.IsExpectedErrors(err, []string{"NotFoundApiGroup", "NotFoundApi"}) {
				return plugin, errs.WrapErrorf(err, errs.NotFoundMsg, errs.AlibabaCloudSdkGoERROR)
			}
			return plugin, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
		}
		logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)
		response, _ := raw.(*cloudapi.DescribePluginsByApiResponse)

		allPlugins = append(allPlugins, response.Plugins.PluginAttribute...)

		if len(allPlugins) < helper.PageSizeLarge {
			break
		}

		page, err := helper.GetNextpageNumber(request.PageNumber)
		if err != nil {
			return plugin, errs.WrapError(err)
		}
		request.PageNumber = page
	}

	var filteredPluginsTemp []cloudapi.PluginAttribute
	for _, plugin := range allPlugins {
		if plugin.PluginId == pluginId {
			filteredPluginsTemp = append(filteredPluginsTemp, plugin)
		}
	}

	if len(filteredPluginsTemp) < 1 {
		return plugin, errs.WrapErrorf(errs.NotFoundErr("ApiGatewayPluginAttachment", id), errs.NotFoundMsg, errs.ProviderERROR)
	}

	plugin = &filteredPluginsTemp[0]
	return plugin, nil
}

func (s *CloudApiService) DescribeApiGatewayVpcAccess(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "DescribeVpcAccesses"
	parts, err := helper.ParseResourceId(id, 4)
	if err != nil {
		return nil, errs.WrapError(err)
	}

	request := map[string]interface{}{
		"Name":       parts[0],
		"VpcId":      parts[1],
		"InstanceId": parts[2],
		"Port":       parts[3],
		"PageNumber": 1,
		"PageSize":   helper.PageSizeLarge,
	}

	idExist := false
	for {
		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(5*time.Minute, func() *retry.RetryError {
			response, err = client.RpcPost("CloudAPI", "2016-07-14", action, nil, request, true)
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

		resp, err := jsonpath.Get("$.VpcAccessAttributes.VpcAccessAttribute", response)
		if err != nil {
			return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.VpcAccessAttributes.VpcAccessAttribute", response)
		}

		if v, ok := resp.([]interface{}); !ok || len(v) < 1 {
			return object, errs.WrapErrorf(errs.NotFoundErr("ApiGateway:VpcAccess", id), errs.NotFoundWithResponse, response)
		}

		for _, v := range resp.([]interface{}) {
			if fmt.Sprint(v.(map[string]interface{})["Name"]) == parts[0] && fmt.Sprint(v.(map[string]interface{})["VpcId"]) == parts[1] && fmt.Sprint(v.(map[string]interface{})["InstanceId"]) == parts[2] && fmt.Sprint(v.(map[string]interface{})["Port"]) == parts[3] {
				idExist = true
				return v.(map[string]interface{}), nil
			}
		}

		if len(resp.([]interface{})) < request["PageSize"].(int) {
			break
		}

		request["PageNumber"] = request["PageNumber"].(int) + 1
	}

	if !idExist {
		return object, errs.WrapErrorf(errs.NotFoundErr("ApiGateway:VpcAccess", id), errs.NotFoundWithResponse, response)
	}

	return object, nil
}

func (s *CloudApiService) WaitForApiGatewayAppAttachment(id string, status Status, timeout int) (err error) {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	parts, err := helper.ParseResourceId(id, 4)
	if err != nil {
		return errs.WrapError(err)
	}
	appIds := parts[2]
	for {
		object, err := s.DescribeApiGatewayAppAttachment(id)
		if err != nil {
			if errs.NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return errs.WrapError(err)
			}
		}
		if strconv.FormatInt(object.AppId, 10) == appIds && status != Deleted {
			return nil
		}
		if time.Now().After(deadline) {
			return errs.WrapErrorf(err, errs.WaitTimeoutMsg, id, helper.GetFunc(1), timeout, strconv.FormatInt(object.AppId, 10), appIds, errs.ProviderERROR)
		}
		time.Sleep(helper.DefaultIntervalShort * time.Second)
	}
}

func (s *CloudApiService) DescribeDeployedApi(id string, stageName string) (*cloudapi.DescribeDeployedApiResponse, error) {
	api := &cloudapi.DescribeDeployedApiResponse{}
	request := cloudapi.CreateDescribeDeployedApiRequest()
	request.RegionId = s.client.RegionId
	parts, err := helper.ParseResourceId(id, 2)
	if err != nil {
		return api, errs.WrapError(err)
	}
	request.ApiId = parts[1]
	request.GroupId = parts[0]
	request.StageName = stageName

	raw, err := s.client.WithCloudApiClient(func(cloudApiClient *cloudapi.Client) (interface{}, error) {
		return cloudApiClient.DescribeDeployedApi(request)
	})
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"NotFoundApiGroup", "NotFoundApi", "NotFoundStage"}) {
			return api, errs.WrapErrorf(err, errs.NotFoundMsg, errs.AlibabaCloudSdkGoERROR)
		}
		return api, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
	}
	logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)
	api, _ = raw.(*cloudapi.DescribeDeployedApiResponse)
	return api, nil
}

func (s *CloudApiService) DeployedApi(id string, stageName string) (err error) {
	request := cloudapi.CreateDeployApiRequest()
	request.RegionId = s.client.RegionId
	parts, err := helper.ParseResourceId(id, 2)
	if err != nil {
		return errs.WrapError(err)
	}
	request.ApiId = parts[1]
	request.GroupId = parts[0]
	request.StageName = stageName
	request.Description = DeployCommonDescription

	raw, err := s.client.WithCloudApiClient(func(cloudApiClient *cloudapi.Client) (interface{}, error) {
		return cloudApiClient.DeployApi(request)
	})
	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, id, request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
	}
	logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)
	return
}

func (s *CloudApiService) AbolishApi(id string, stageName string) (err error) {
	request := cloudapi.CreateAbolishApiRequest()
	request.RegionId = s.client.RegionId
	parts, err := helper.ParseResourceId(id, 2)
	if err != nil {
		return errs.WrapError(err)
	}
	request.ApiId = parts[1]
	request.GroupId = parts[0]
	request.StageName = stageName

	raw, err := s.client.WithCloudApiClient(func(cloudApiClient *cloudapi.Client) (interface{}, error) {
		return cloudApiClient.AbolishApi(request)
	})

	if err != nil {
		if errs.IsExpectedErrors(err, []string{"NotFoundApiGroup", "NotFoundApi", "NotFoundStage"}) {
			return errs.WrapErrorf(err, errs.NotFoundMsg, errs.AlibabaCloudSdkGoERROR)
		}
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, id, request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
	}
	logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)
	return
}

func (s *CloudApiService) DescribeTags(resourceId string, resourceTags map[string]interface{}, resourceType tags.TagResourceType) (tags []cloudapi.TagResource, err error) {
	request := cloudapi.CreateListTagResourcesRequest()
	request.RegionId = s.client.RegionId
	request.ResourceType = string(resourceType)
	request.ResourceId = &[]string{resourceId}
	if resourceTags != nil && len(resourceTags) > 0 {
		var reqTags []cloudapi.ListTagResourcesTag
		for key, value := range resourceTags {
			reqTags = append(reqTags, cloudapi.ListTagResourcesTag{
				Key:   key,
				Value: value.(string),
			})
		}
		request.Tag = &reqTags
	}
	raw, err := s.client.WithCloudApiClient(func(cloudApiClient *cloudapi.Client) (interface{}, error) {
		return cloudApiClient.ListTagResources(request)
	})
	if err != nil {
		err = errs.WrapErrorf(err, errs.DefaultErrorMsg, resourceId, request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
		return
	}
	logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*cloudapi.ListTagResourcesResponse)

	return response.TagResources.TagResource, nil
}

func (s *CloudApiService) setInstanceTags(d *schema.ResourceData, resourceType tags.TagResourceType) error {
	oraw, nraw := d.GetChange("tags")
	o := oraw.(map[string]interface{})
	n := nraw.(map[string]interface{})
	create, remove := s.diffTags(s.tagsFromMap(o), s.tagsFromMap(n))

	if len(remove) > 0 {
		var tagKey []string
		for _, v := range remove {
			tagKey = append(tagKey, v.Key)
		}
		request := cloudapi.CreateUntagResourcesRequest()
		request.ResourceId = &[]string{d.Id()}
		request.ResourceType = string(resourceType)
		request.TagKey = &tagKey
		request.RegionId = s.client.RegionId
		raw, err := s.client.WithCloudApiClient(func(client *cloudapi.Client) (interface{}, error) {
			return client.UntagResources(request)
		})
		if err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
		}
		logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)
	}

	if len(create) > 0 {
		request := cloudapi.CreateTagResourcesRequest()
		request.ResourceId = &[]string{d.Id()}
		request.Tag = &create
		request.ResourceType = string(resourceType)
		request.RegionId = s.client.RegionId
		raw, err := s.client.WithCloudApiClient(func(client *cloudapi.Client) (interface{}, error) {
			return client.TagResources(request)
		})
		if err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
		}
		logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)
	}

	return nil
}

func (s *CloudApiService) tagsToMap(tags []cloudapi.TagResource) map[string]string {
	result := make(map[string]string)
	for _, t := range tags {
		if !s.ignoreTag(t) {
			result[t.TagKey] = t.TagValue
		}
	}
	return result
}

func (s *CloudApiService) ignoreTag(t cloudapi.TagResource) bool {
	filter := []string{"^aliyun", "^acs:", "^http://", "^https://"}
	for _, v := range filter {
		log.Printf("[DEBUG] Matching prefix %v with %v\n", v, t.TagKey)
		ok, _ := regexp.MatchString(v, t.TagKey)
		if ok {
			log.Printf("[DEBUG] Found Alibaba Cloud specific t %s (val: %s), ignoring.\n", t.TagKey, t.TagValue)
			return true
		}
	}
	return false
}

func (s *CloudApiService) diffTags(oldTags, newTags []cloudapi.TagResourcesTag) ([]cloudapi.TagResourcesTag, []cloudapi.TagResourcesTag) {
	// First, we're creating everything we have
	create := make(map[string]interface{})
	for _, t := range newTags {
		create[t.Key] = t.Value
	}

	// Build the list of what to remove
	var remove []cloudapi.TagResourcesTag
	for _, t := range oldTags {
		old, ok := create[t.Key]
		if !ok || old != t.Value {
			// Delete it!
			remove = append(remove, t)
		}
	}

	return s.tagsFromMap(create), remove
}

func (s *CloudApiService) tagsFromMap(m map[string]interface{}) []cloudapi.TagResourcesTag {
	result := make([]cloudapi.TagResourcesTag, 0, len(m))
	for k, v := range m {
		result = append(result, cloudapi.TagResourcesTag{
			Key:   k,
			Value: v.(string),
		})
	}

	return result
}

func (s *CloudApiService) DescribeApiGatewayBackend(id string) (object map[string]interface{}, err error) {
	client := s.client

	request := map[string]interface{}{
		"BackendId": id,
	}

	var response map[string]interface{}
	action := "DescribeBackendInfo"
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		resp, err := client.RpcPost("CloudAPI", "2016-07-14", action, nil, request, true)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		response = resp
		logs.AddDebug(action, response, request)
		return nil
	})
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"NotFoundBackend"}) {
			return object, errs.WrapErrorf(err, errs.NotFoundMsg, errs.AlibabaCloudSdkGoERROR)
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.BackendInfo", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.BackendInfo", response)
	}
	return v.(map[string]interface{}), nil
}

func (s *CloudApiService) DescribeApiGatewayLogConfig(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "DescribeLogConfig"
	request := map[string]interface{}{
		"LogType": id,
	}
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("CloudAPI", "2016-07-14", action, nil, request, true)
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
	v, err := jsonpath.Get("$.LogInfos.LogInfo", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.LogInfos.LogInfo", response)
	}
	if len(v.([]interface{})) < 1 {
		return object, errs.WrapErrorf(errs.NotFoundErr("ApiGateway", id), errs.NotFoundWithResponse, response)
	} else {
		if fmt.Sprint(v.([]interface{})[0].(map[string]interface{})["LogType"]) != id {
			return object, errs.WrapErrorf(errs.NotFoundErr("ApiGateway", id), errs.NotFoundWithResponse, response)
		}
	}
	object = v.([]interface{})[0].(map[string]interface{})
	return object, nil
}

func (s *CloudApiService) DescribeApiGatewayModel(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	action := "DescribeModels"

	client := s.client

	parts, err := helper.ParseResourceId(id, 2)
	if err != nil {
		return nil, errs.WrapError(err)
	}
	request := map[string]interface{}{
		"GroupId":   parts[0],
		"ModelName": parts[1],
	}

	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("CloudAPI", "2016-07-14", action, nil, request, true)
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

	v, err := jsonpath.Get("$.ModelDetails.ModelDetail", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.ModelDetails.ModelDetail", response)
	}

	if len(v.([]interface{})) < 1 {
		return object, errs.WrapErrorf(errs.NotFoundErr("ApiGateway", id), errs.NotFoundWithResponse, response)
	} else {
		if fmt.Sprint(v.([]interface{})[0].(map[string]interface{})["ModelName"]) != parts[1] {
			return object, errs.WrapErrorf(errs.NotFoundErr("ApiGateway", id), errs.NotFoundWithResponse, response)
		}
	}

	object = v.([]interface{})[0].(map[string]interface{})

	return object, nil
}

func (s *CloudApiService) DescribeApiGatewayStageModel(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	action := "DescribeStageModels"
	client := s.client

	request := map[string]interface{}{
		"PageNumber": 1,
		"PageSize":   helper.PageSizeXLarge,
	}

	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("CloudAPI", "2016-07-14", action, nil, request, true)
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

	v, err := jsonpath.Get("$.StageModelInfoList", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.StageModelInfoList", response)
	}

	if v == nil {
		return object, errs.WrapErrorf(errs.NotFoundErr("ApiGateway", id), errs.NotFoundWithResponse, response)
	}

	for _, item := range v.([]interface{}) {
		itemMap := item.(map[string]interface{})
		if fmt.Sprint(itemMap["StageModelId"]) == id {
			return itemMap, nil
		}
	}

	return object, errs.WrapErrorf(errs.NotFoundErr("ApiGateway", id), errs.NotFoundWithResponse, response)
}

func (s *CloudApiService) DescribeApiGatewayBackendModel(id string) (object map[string]interface{}, err error) {
	client := s.client

	parts, err := helper.ParseResourceId(id, 2)
	if err != nil {
		return object, errs.WrapError(err)
	}
	backendId := parts[0]
	stageName := parts[1]

	request := map[string]interface{}{
		"BackendId": backendId,
	}

	var response map[string]interface{}
	action := "DescribeBackendInfo"
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		resp, err := client.RpcPost("CloudAPI", "2016-07-14", action, nil, request, true)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		response = resp
		logs.AddDebug(action, response, request)
		return nil
	})
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"NotFoundBackend"}) {
			return object, errs.WrapErrorf(err, errs.NotFoundMsg, errs.AlibabaCloudSdkGoERROR)
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$.BackendInfo", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.BackendInfo", response)
	}

	backendInfo, ok := v.(map[string]interface{})
	if !ok {
		return object, errs.WrapError(errs.Error("DescribeApiGatewayBackendModel: BackendInfo is not a map"))
	}

	var backendModels []interface{}

	if models, ok := backendInfo["BackendModels"]; ok {
		if modelList, ok := models.([]interface{}); ok {
			backendModels = modelList
		}
	} else if models, ok := backendInfo["BackendModelList"]; ok {
		if modelList, ok := models.([]interface{}); ok {
			backendModels = modelList
		}
	} else if models, ok := backendInfo["BackendModel"]; ok {
		if modelList, ok := models.([]interface{}); ok {
			backendModels = modelList
		} else if modelMap, ok := models.(map[string]interface{}); ok {
			backendModels = []interface{}{modelMap}
		}
	}

	backendType := ""
	if v, ok := backendInfo["BackendType"]; ok {
		backendType = fmt.Sprint(v)
	}

	for _, model := range backendModels {
		modelMap, ok := model.(map[string]interface{})
		if !ok {
			continue
		}
		if modelMap["StageName"] == stageName {
			modelMap["BackendId"] = backendId
			modelMap["BackendType"] = backendType

			if backendConfig, ok := modelMap["BackendConfig"].(map[string]interface{}); ok {
				modelData := make(map[string]interface{})
				if v, ok := backendConfig["ServiceTimeout"].(float64); ok && v > 0 {
					modelData["ServiceTimeout"] = int(v)
				}
				switch backendType {
				case "HTTP":
					if v, ok := backendConfig["ServiceAddress"]; ok && v != "" {
						modelData["ServiceAddress"] = v
					}
					if v, ok := backendConfig["HttpTargetHostName"]; ok && v != "" {
						modelData["HttpTargetHostName"] = v
					}
				case "VPC":
					if vpcConfig, ok := backendConfig["VpcConfig"].(map[string]interface{}); ok && len(vpcConfig) > 0 {
						newVpcConfig := make(map[string]interface{})
						if v, ok := vpcConfig["VpcAccessId"]; ok && v != "" {
							newVpcConfig["VpcAccessId"] = v
						}
						if v, ok := vpcConfig["VpcScheme"]; ok && v != "" {
							newVpcConfig["VpcScheme"] = v
						}
						if v, ok := vpcConfig["VpcTargetHostName"]; ok && v != "" {
							newVpcConfig["VpcTargetHostName"] = v
						}
						if v, ok := vpcConfig["ServiceHttpMethod"]; ok && v != "" {
							newVpcConfig["ServiceHttpMethod"] = v
						}
						if v, ok := vpcConfig["ServicePath"]; ok && v != "" {
							newVpcConfig["ServicePath"] = v
						}
						if len(newVpcConfig) > 0 {
							modelData["VpcConfig"] = newVpcConfig
						}
					}
				case "FC_EVENT", "FC_EVENT_V3", "FC_HTTP", "FC_HTTP_V3":
					if fcConfig, ok := backendConfig["FunctionComputeConfig"].(map[string]interface{}); ok && len(fcConfig) > 0 {
						newFcConfig := make(map[string]interface{})
						for _, key := range []string{
							"FcRegionId", "FcType", "ServiceName", "FunctionName", "RoleArn",
							"Qualifier", "FcVersion", "Path", "FcBaseUrl", "Method",
							"ContentTypeCatagory", "ContentTypeValue", "OnlyBusinessPath", "TriggerName",
						} {
							if v, ok := fcConfig[key]; ok && v != "" {
								newFcConfig[key] = v
							}
						}
						if len(newFcConfig) > 0 {
							modelData["FunctionComputeConfig"] = newFcConfig
						}
					}
				case "OSS":
					if ossConfig, ok := backendConfig["OssConfig"].(map[string]interface{}); ok && len(ossConfig) > 0 {
						newOssConfig := make(map[string]interface{})
						for _, key := range []string{"OssRegionId", "BucketName", "Key", "Action"} {
							if v, ok := ossConfig[key]; ok && v != "" {
								newOssConfig[key] = v
							}
						}
						if len(newOssConfig) > 0 {
							modelData["OssConfig"] = newOssConfig
						}
					}
				case "MOCK":
					if mockConfig, ok := backendConfig["MockConfig"].(map[string]interface{}); ok && len(mockConfig) > 0 {
						newMockConfig := make(map[string]interface{})
						if v, ok := mockConfig["MockResult"]; ok && v != "" {
							newMockConfig["MockResult"] = v
						}
						if v, ok := mockConfig["MockStatusCode"]; ok {
							if code, ok := v.(float64); ok && code > 0 {
								newMockConfig["MockStatusCode"] = int(code)
							} else if code, ok := v.(int); ok && code > 0 {
								newMockConfig["MockStatusCode"] = code
							}
						}
						if headers, ok := mockConfig["MockHeaders"].([]interface{}); ok && len(headers) > 0 {
							newMockConfig["MockHeaders"] = headers
						}
						if len(newMockConfig) > 0 {
							modelData["MockConfig"] = newMockConfig
						}
					}
				case "EVENTBRIDGE":
					if ebConfig, ok := backendConfig["EventBridgeConfig"].(map[string]interface{}); ok && len(ebConfig) > 0 {
						newEbConfig := make(map[string]interface{})
						for k, v := range ebConfig {
							if v != "" && v != nil {
								newEbConfig[k] = v
							}
						}
						if len(newEbConfig) > 0 {
							modelData["EventBridgeConfig"] = newEbConfig
						}
					}
				}
				if len(modelData) > 0 {
					bytes, err := json.Marshal(modelData)
					if err == nil {
						modelMap["BackendModelData"] = string(bytes)
					}
				}
			}

			return modelMap, nil
		}
	}

	return object, errs.WrapErrorf(err, errs.NotFoundMsg, errs.AlibabaCloudSdkGoERROR)
}

func (s *CloudApiService) SetResourceTags(d *schema.ResourceData, resourceType string) error {

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
			action := "UntagResources"
			request := map[string]interface{}{
				"RegionId":     s.client.RegionId,
				"ResourceType": resourceType,
				"ResourceId.1": d.Id(),
			}
			for i, key := range removedTagKeys {
				request[fmt.Sprintf("TagKey.%d", i+1)] = key
			}
			wait := helper.IncrementalWait(2*time.Second, 1*time.Second)
			err := retry.Retry(10*time.Minute, func() *retry.RetryError {
				response, err := client.RpcPost("CloudAPI", "2016-07-14", action, nil, request, false)
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
				"RegionId":     s.client.RegionId,
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
				response, err := client.RpcPost("CloudAPI", "2016-07-14", action, nil, request, false)
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

func (s *CloudApiService) DescribeApiGatewayPlugin(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "DescribePlugins"
	request := map[string]interface{}{
		"PluginId": id,
	}
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("CloudAPI", "2016-07-14", action, nil, request, true)
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
	v, err := jsonpath.Get("$.Plugins.PluginAttribute", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.Plugins.PluginAttribute", response)
	}
	if len(v.([]interface{})) < 1 {
		return object, errs.WrapErrorf(errs.NotFoundErr("ApiGateway", id), errs.NotFoundWithResponse, response)
	} else {
		if fmt.Sprint(v.([]interface{})[0].(map[string]interface{})["PluginId"]) != id {
			return object, errs.WrapErrorf(errs.NotFoundErr("ApiGateway", id), errs.NotFoundWithResponse, response)
		}
	}
	object = v.([]interface{})[0].(map[string]interface{})
	return object, nil
}

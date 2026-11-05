package apigateway

import (
	"fmt"
	"strconv"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	logs "github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
)

// appKeyExpectedNotFoundCodes lists the API error codes that signal an AppKey
// is already gone, so Delete can treat them as success.
var appKeyExpectedNotFoundCodes = []string{"NotFoundApp", "NotFoundAppKey"}

// DescribeAppSecurities lists every AppKey/AppSecret pair of the given App.
// CloudAPI/2016-07-14 DescribeAppSecurities returns them under
// $.AppSecuritys.AppSecurity. AppId is an int64 query parameter, so callers
// pass the parsed value; the App resource stores AppId as a decimal string.
func (s *CloudApiService) DescribeAppSecurities(appId int64) ([]map[string]interface{}, error) {
	action := "DescribeAppSecurities"
	request := map[string]interface{}{
		"RegionId": s.client.RegionId,
		"AppId":    appId,
	}
	response, err := s.client.RpcPost("CloudAPI", "2016-07-14", action, nil, request, true)
	logs.AddDebug(action, response, request)
	if err != nil {
		return nil, errs.WrapErrorf(err, errs.DefaultErrorMsg, strconv.FormatInt(appId, 10), action, errs.AlibabaCloudSdkGoERROR)
	}
	sec, err := jsonpath.Get("$.AppSecuritys.AppSecurity", response)
	if err != nil {
		return nil, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, strconv.FormatInt(appId, 10), "$.AppSecuritys.AppSecurity", response)
	}
	items, ok := sec.([]interface{})
	if !ok {
		return nil, errs.WrapErrorf(errs.NotFoundErr("ApiGateway:AppKey", strconv.FormatInt(appId, 10)), errs.NotFoundWithResponse, response)
	}
	result := make([]map[string]interface{}, 0, len(items))
	for _, raw := range items {
		if m, ok := raw.(map[string]interface{}); ok {
			result = append(result, m)
		}
	}
	return result, nil
}

// FindAppSecurity returns the AppSecurity entry whose AppKey matches, or a
// NotFound error when the AppKey no longer exists under the App.
func (s *CloudApiService) FindAppSecurity(appId int64, appKey string) (map[string]interface{}, error) {
	items, err := s.DescribeAppSecurities(appId)
	if err != nil {
		return nil, err
	}
	for _, m := range items {
		if fmt.Sprint(m["AppKey"]) == appKey {
			return m, nil
		}
	}
	return nil, errs.WrapErrorf(errs.NotFoundErr("ApiGateway:AppKey", fmt.Sprintf("%d:%s", appId, appKey)), errs.NotFoundWithResponse, nil)
}

// CreateAppKey creates an AppKey under the App. AppKey and AppSecret are both
// optional on CreateAppKey; when the caller leaves them empty the server
// generates them and they are read back through DescribeAppSecurities, since
// CreateAppKey's response only carries a RequestId.
func (s *CloudApiService) CreateAppKey(appId, appKey, appSecret string) error {
	action := "CreateAppKey"
	request := map[string]interface{}{
		"RegionId": s.client.RegionId,
		"AppId":    appId,
	}
	if appKey != "" {
		request["AppKey"] = appKey
	}
	if appSecret != "" {
		request["AppSecret"] = appSecret
	}
	response, err := s.client.RpcPost("CloudAPI", "2016-07-14", action, nil, request, true)
	logs.AddDebug(action, response, request)
	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, appId, action, errs.AlibabaCloudSdkGoERROR)
	}
	return nil
}

// ResetAppSecret rotates the AppSecret of an AppKey by calling
// CloudAPI/2016-07-14 ResetAppSecret. NewAppKey is intentionally not supplied
// because changing the AppKey identity would change the resource's primary
// identifier (the id is `<app_id>:<app_key>`); rotating the AppKey identity is
// therefore a ForceNew, not an in-place update. NewAppSecret is left empty so
// the server auto-generates a fresh AppSecret. The response only carries a
// RequestId, so the caller reads the new AppSecret back through
// DescribeAppSecurities (FindAppSecurity).
func (s *CloudApiService) ResetAppSecret(appKey string) error {
	action := "ResetAppSecret"
	request := map[string]interface{}{
		"RegionId": s.client.RegionId,
		"AppKey":   appKey,
	}
	response, err := s.client.RpcPost("CloudAPI", "2016-07-14", action, nil, request, true)
	logs.AddDebug(action, response, request)
	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, appKey, action, errs.AlibabaCloudSdkGoERROR)
	}
	return nil
}

// DeleteAppKey removes an AppKey from its App. A missing App or AppKey is
// treated as success so a drifted destroy converges.
func (s *CloudApiService) DeleteAppKey(appId, appKey string) error {
	action := "DeleteAppKey"
	request := map[string]interface{}{
		"RegionId": s.client.RegionId,
		"AppId":    appId,
		"AppKey":   appKey,
	}
	response, err := s.client.RpcPost("CloudAPI", "2016-07-14", action, nil, request, true)
	logs.AddDebug(action, response, request)
	if err != nil {
		if errs.IsExpectedErrors(err, appKeyExpectedNotFoundCodes) {
			return nil
		}
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, fmt.Sprintf("%s:%s", appId, appKey), action, errs.AlibabaCloudSdkGoERROR)
	}
	return nil
}

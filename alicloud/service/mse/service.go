package mse

import (
	"fmt"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/tags"

	"github.com/PaesslerAG/jsonpath"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Service struct {
	client *connectivity.AliyunClient
}

// NewService constructs the mse Service from a connectivity client.
func NewService(client *connectivity.AliyunClient) *Service {
	return &Service{client: client}
}

func (s *Service) DescribeMseCluster(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "QueryClusterDetail"
	request := map[string]interface{}{
		"RegionId":   s.client.RegionId,
		"InstanceId": id,
	}
	request["ClientToken"] = helper.BuildClientToken("QueryClusterDetail")
	response, err = client.RpcPost("mse", "2019-05-31", action, nil, request, true)
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"mse-200-021"}) {
			err = errs.WrapErrorf(errs.NotFoundErr("MseCluster", id), errs.NotFoundMsg, errs.ProviderERROR)
			return object, err
		}
		err = errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
		return object, err
	}
	if fmt.Sprint(response["Success"]) == "false" {
		return object, errs.WrapError(fmt.Errorf("%s failed, response: %v", action, response))
	}
	logs.AddDebug(action, response, request)
	v, err := jsonpath.Get("$.Data", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.Data", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}

func (s *Service) GetInstanceIdBYClusterId(clusterId string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "QueryClusterDetail"
	request := map[string]interface{}{
		"RegionId":  s.client.RegionId,
		"ClusterId": clusterId,
	}
	request["ClientToken"] = helper.BuildClientToken("QueryClusterDetail")
	response, err = client.RpcPost("mse", "2019-05-31", action, nil, request, true)
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"mse-200-021"}) {
			err = errs.WrapErrorf(errs.NotFoundErr("MseCluster", clusterId), errs.NotFoundMsg, errs.ProviderERROR)
			return object, err
		}
		err = errs.WrapErrorf(err, errs.DefaultErrorMsg, clusterId, action, errs.AlibabaCloudSdkGoERROR)
		return object, err
	}
	if fmt.Sprint(response["Success"]) == "false" {
		return object, errs.WrapError(fmt.Errorf("%s failed, response: %v", action, response))
	}
	logs.AddDebug(action, response, request)
	v, err := jsonpath.Get("$.Data", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, clusterId, "$.Data", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}

func (s *Service) SetResourceTags(d *schema.ResourceData, resourceType string) error {
	if d.HasChange("tags") {
		var err error
		var action string
		client := s.client
		var request map[string]interface{}
		var response map[string]interface{}
		query := make(map[string]interface{})

		added, removed := tags.ParsingTags(d)
		removedTagKeys := make([]string, 0)
		for _, v := range removed {
			if !tags.IgnoredTags(v, "") {
				removedTagKeys = append(removedTagKeys, v)
			}
		}
		if len(removedTagKeys) > 0 {
			action = "UnTagResources"
			request = make(map[string]interface{})
			query = make(map[string]interface{})
			request["ResourceId.1"] = d.Id()
			request["RegionId"] = client.RegionId
			for i, key := range removedTagKeys {
				request[fmt.Sprintf("TagKey.%d", i+1)] = key
			}

			request["ResourceType"] = resourceType
			wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
			err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
				response, err = client.RpcPost("mse", "2019-05-31", action, query, request, false)

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
			action = "TagResources"
			request = make(map[string]interface{})
			query = make(map[string]interface{})
			request["ResourceId.1"] = d.Id()
			request["RegionId"] = client.RegionId
			count := 1
			for key, value := range added {
				request[fmt.Sprintf("Tag.%d.Key", count)] = key
				request[fmt.Sprintf("Tag.%d.Value", count)] = value
				count++
			}

			request["ResourceType"] = resourceType
			wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
			err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
				response, err = client.RpcPost("mse", "2019-05-31", action, query, request, false)

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

func (s *Service) MseClusterStateRefreshFunc(id string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeMseCluster(id)
		if err != nil {
			if errs.NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", errs.WrapError(err)
		}

		for _, failState := range failStates {
			if object["InitStatus"].(string) == failState {
				return object, object["InitStatus"].(string), errs.WrapError(errs.Error(errs.FailedToReachTargetStatus, object["InitStatus"].(string)))
			}
		}
		return object, object["InitStatus"].(string), nil
	}
}

func (s *Service) DescribeMseGateway(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "GetGateway"
	request := map[string]interface{}{
		"GatewayUniqueId": id,
	}
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcGet("mse", "2019-05-31", action, request, nil)
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
	if errs.IsExpectedErrorCodes(fmt.Sprint(response["Code"]), []string{"404"}) {
		return object, errs.WrapErrorf(errs.NotFoundErr("MSE:Gateway", id), errs.NotFoundMsg, errs.ProviderERROR)
	}
	if fmt.Sprint(response["Success"]) == "false" {
		return object, errs.WrapError(fmt.Errorf("%s failed, response: %v", action, response))
	}
	v, err := jsonpath.Get("$.Data", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.Data", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}

func (s *Service) MseGatewayStateRefreshFunc(id string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeMseGateway(id)
		if err != nil {
			if errs.NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", errs.WrapError(err)
		}

		for _, failState := range failStates {
			if fmt.Sprint(object["Status"]) == failState {
				return object, fmt.Sprint(object["Status"]), errs.WrapError(errs.Error(errs.FailedToReachTargetStatus, fmt.Sprint(object["Status"])))
			}
		}
		return object, fmt.Sprint(object["Status"]), nil
	}
}

func (s *Service) ListGatewaySlb(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "ListGatewaySlb"
	request := map[string]interface{}{
		"GatewayUniqueId": id,
	}
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcGet("mse", "2019-05-31", action, request, nil)
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
	if errs.IsExpectedErrorCodes(fmt.Sprint(response["Code"]), []string{"404"}) {
		return object, errs.WrapErrorf(errs.NotFoundErr("MSE:Gateway", id), errs.NotFoundMsg, errs.ProviderERROR)
	}
	if fmt.Sprint(response["Success"]) == "false" {
		return object, errs.WrapError(fmt.Errorf("%s failed, response: %v", action, response))
	}
	v, err := jsonpath.Get("$.Data", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.Data", response)
	}
	object = make(map[string]interface{})
	object["SlbList"] = v
	return object, nil
}

func (s *Service) DescribeMseZnode(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "QueryZnodeDetail"
	parts, err := helper.ParseResourceId(id, 2)
	if err != nil {
		err = errs.WrapError(err)
		return
	}
	request := map[string]interface{}{
		"ClusterId": parts[0],
		"Path":      parts[1],
	}
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcGet("mse", "2019-05-31", action, request, nil)
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
	if fmt.Sprint(response["Success"]) == "false" {
		return object, errs.WrapErrorf(errs.NotFoundErr("MSE:Znode", id), errs.NotFoundMsg, errs.ProviderERROR)
	}
	v, err := jsonpath.Get("$.Data", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.Data", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}
func (s *Service) DescribeMseNacosConfig(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	parts, err := helper.ParseResourceIdWithEscaped(id, 4)
	if err != nil {
		return nil, errs.WrapError(err)
	}
	action := "GetNacosConfig"
	request := map[string]interface{}{
		"InstanceId":  parts[0],
		"NamespaceId": parts[1],
		"DataId":      parts[2],
		"Group":       parts[3],
	}
	response, err = client.RpcPost("mse", "2019-05-31", action, nil, request, true)
	// For delete check
	if errs.IsExpectedErrors(err, []string{"InternalError"}) {
		return object, errs.WrapErrorf(errs.NotFoundErr("MSE:MseNacosConfig", id), errs.NotFoundMsg, errs.ProviderERROR)
	}

	if err != nil {
		if errs.IsExpectedErrors(err, []string{"mse-200-021"}) {
			err = errs.WrapErrorf(errs.NotFoundErr("MseNacosConfig", id), errs.NotFoundMsg, errs.ProviderERROR)
			return object, err
		}
		err = errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
		return object, err
	}
	if fmt.Sprint(response["Success"]) == "false" {
		return object, errs.WrapError(fmt.Errorf("%s failed, response: %v", action, response))
	}
	logs.AddDebug(action, response, request)
	v, err := jsonpath.Get("$.Configuration", response)
	if err != nil {
		if err.Error() == "unknown key Configuration" {
			return object, errs.WrapErrorf(errs.NotFoundErr("MSE", id), errs.NotFoundWithResponse, response)
		}
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.Configuration", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}

func (s *Service) DescribeMseEngineNamespace(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	parts, err := helper.ParseResourceId(id, 2)
	if err != nil {
		err = errs.WrapError(err)
		return
	}
	action := "ListEngineNamespaces"
	request := map[string]interface{}{
		"InstanceId": parts[0],
	}
	idExist := false
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcGet("mse", "2019-05-31", action, request, nil)
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
	// There is an API bug while query an Instance deleted.
	if errs.IsExpectedErrors(err, []string{"InternalError"}) {
		return object, errs.WrapErrorf(errs.NotFoundErr("MSE:EngineNamespace", id), errs.NotFoundMsg, errs.ProviderERROR)
	}
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"InvalidParameter"}) {
			return object, errs.WrapErrorf(errs.NotFoundErr("MSE cluster", parts[0]), errs.NotFoundMsg, "")
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}
	if fmt.Sprint(response["Success"]) == "false" {
		return object, errs.WrapErrorf(errs.NotFoundErr("MSE:EngineNamespace", id), errs.NotFoundMsg, errs.ProviderERROR)
	}
	v, err := jsonpath.Get("$.Data", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.Data", response)
	}
	if len(v.([]interface{})) < 1 {
		return object, errs.WrapErrorf(errs.NotFoundErr("MSE", id), errs.NotFoundWithResponse, response)
	}
	for _, v := range v.([]interface{}) {
		if fmt.Sprint(v.(map[string]interface{})["Namespace"]) == parts[1] {
			idExist = true
			return v.(map[string]interface{}), nil
		}
	}
	if !idExist {
		return object, errs.WrapErrorf(errs.NotFoundErr("MSE", id), errs.NotFoundWithResponse, response)
	}
	return object, nil
}

// DescribeVswitch is a local copy of root VpcService.DescribeVswitch, inlined
// here so the mse package does not import the not-yet-migrated root package or
// another product package (migration stock-method-disposition #4). Behavior is
// preserved verbatim.
func (s *Service) DescribeVswitch(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "DescribeVSwitchAttributes"
	request := map[string]interface{}{
		"RegionId":  client.RegionId,
		"VSwitchId": id,
	}
	response, err = client.RpcPost("Vpc", "2016-04-28", action, nil, request, true)
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"InvalidVSwitchId.NotFound", "InvalidVswitchID.NotFound"}) {
			err = errs.WrapErrorf(errs.NotFoundErr("Vswitch", id), errs.NotFoundWithError, err)
			return object, err
		}
		err = errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
		return object, err
	}
	logs.AddDebug(action, response, request)
	v, err := jsonpath.Get("$", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$", response)
	}
	object = v.(map[string]interface{})
	if fmt.Sprint(object["VSwitchId"]) != id {
		return object, errs.WrapErrorf(errs.NotFoundErr("vswitch", id), errs.NotFoundWithResponse, response)
	}
	return object, nil
}

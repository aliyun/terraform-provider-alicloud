package alicloud

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

type MseService struct {
	client *connectivity.AliyunClient
}

func (s *MseService) DescribeMseCluster(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "QueryClusterDetail"
	request := map[string]interface{}{
		"RegionId":   s.client.RegionId,
		"InstanceId": id,
	}
	request["ClientToken"] = buildClientToken("QueryClusterDetail")
	response, err = client.RpcPost("mse", "2019-05-31", action, nil, request, true)
	if err != nil {
		if IsExpectedErrors(err, []string{"mse-200-021"}) {
			err = WrapErrorf(Error(GetNotFoundMessage("MseCluster", id)), NotFoundMsg, ProviderERROR)
			return object, err
		}
		err = WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
		return object, err
	}
	if fmt.Sprint(response["Success"]) == "false" {
		return object, WrapError(fmt.Errorf("%s failed, response: %v", action, response))
	}
	addDebug(action, response, request)
	v, err := jsonpath.Get("$.Data", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Data", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}

func (s *MseService) GetInstanceIdBYClusterId(clusterId string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "QueryClusterDetail"
	request := map[string]interface{}{
		"RegionId":  s.client.RegionId,
		"ClusterId": clusterId,
	}
	request["ClientToken"] = buildClientToken("QueryClusterDetail")
	response, err = client.RpcPost("mse", "2019-05-31", action, nil, request, true)
	if err != nil {
		if IsExpectedErrors(err, []string{"mse-200-021"}) {
			err = WrapErrorf(Error(GetNotFoundMessage("MseCluster", clusterId)), NotFoundMsg, ProviderERROR)
			return object, err
		}
		err = WrapErrorf(err, DefaultErrorMsg, clusterId, action, AlibabaCloudSdkGoERROR)
		return object, err
	}
	if fmt.Sprint(response["Success"]) == "false" {
		return object, WrapError(fmt.Errorf("%s failed, response: %v", action, response))
	}
	addDebug(action, response, request)
	v, err := jsonpath.Get("$.Data", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, clusterId, "$.Data", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}

func (s *MseService) SetResourceTags(d *schema.ResourceData, resourceType string) error {
	if d.HasChange("tags") {
		var err error
		var action string
		client := s.client
		var request map[string]interface{}
		var response map[string]interface{}
		query := make(map[string]interface{})

		added, removed := parsingTags(d)
		removedTagKeys := make([]string, 0)
		for _, v := range removed {
			if !ignoredTags(v, "") {
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
			wait := incrementalWait(3*time.Second, 5*time.Second)
			err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
				response, err = client.RpcPost("mse", "2019-05-31", action, query, request, false)

				if err != nil {
					if NeedRetry(err) {
						wait()
						return resource.RetryableError(err)
					}
					return resource.NonRetryableError(err)
				}
				addDebug(action, response, request)
				return nil
			})
			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
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
			wait := incrementalWait(3*time.Second, 5*time.Second)
			err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
				response, err = client.RpcPost("mse", "2019-05-31", action, query, request, false)

				if err != nil {
					if NeedRetry(err) {
						wait()
						return resource.RetryableError(err)
					}
					return resource.NonRetryableError(err)
				}
				addDebug(action, response, request)
				return nil
			})
			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
			}

		}
		d.SetPartial("tags")
	}

	return nil
}

func (s *MseService) MseClusterStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeMseCluster(id)
		if err != nil {
			if NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}

		for _, failState := range failStates {
			if object["InitStatus"].(string) == failState {
				return object, object["InitStatus"].(string), WrapError(Error(FailedToReachTargetStatus, object["InitStatus"].(string)))
			}
		}
		return object, object["InitStatus"].(string), nil
	}
}

func (s *MseService) DescribeMseGateway(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "GetGateway"
	request := map[string]interface{}{
		"GatewayUniqueId": id,
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = client.RpcGet("mse", "2019-05-31", action, request, nil)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	if IsExpectedErrorCodes(fmt.Sprint(response["Code"]), []string{"404"}) {
		return object, WrapErrorf(Error(GetNotFoundMessage("MSE:Gateway", id)), NotFoundMsg, ProviderERROR)
	}
	if fmt.Sprint(response["Success"]) == "false" {
		return object, WrapError(fmt.Errorf("%s failed, response: %v", action, response))
	}
	v, err := jsonpath.Get("$.Data", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Data", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}

func (s *MseService) MseGatewayStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeMseGateway(id)
		if err != nil {
			if NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}

		for _, failState := range failStates {
			if fmt.Sprint(object["Status"]) == failState {
				return object, fmt.Sprint(object["Status"]), WrapError(Error(FailedToReachTargetStatus, fmt.Sprint(object["Status"])))
			}
		}
		return object, fmt.Sprint(object["Status"]), nil
	}
}

func (s *MseService) ListGatewaySlb(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "ListGatewaySlb"
	request := map[string]interface{}{
		"GatewayUniqueId": id,
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = client.RpcGet("mse", "2019-05-31", action, request, nil)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	if IsExpectedErrorCodes(fmt.Sprint(response["Code"]), []string{"404"}) {
		return object, WrapErrorf(Error(GetNotFoundMessage("MSE:Gateway", id)), NotFoundMsg, ProviderERROR)
	}
	if fmt.Sprint(response["Success"]) == "false" {
		return object, WrapError(fmt.Errorf("%s failed, response: %v", action, response))
	}
	v, err := jsonpath.Get("$.Data", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Data", response)
	}
	object = make(map[string]interface{})
	object["SlbList"] = v
	return object, nil
}

func (s *MseService) DescribeMseZnode(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "QueryZnodeDetail"
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		err = WrapError(err)
		return
	}
	request := map[string]interface{}{
		"ClusterId": parts[0],
		"Path":      parts[1],
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = client.RpcGet("mse", "2019-05-31", action, request, nil)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	if fmt.Sprint(response["Success"]) == "false" {
		return object, WrapErrorf(Error(GetNotFoundMessage("MSE:Znode", id)), NotFoundMsg, ProviderERROR)
	}
	v, err := jsonpath.Get("$.Data", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Data", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}
func (s *MseService) DescribeMseNacosConfig(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	parts, err := ParseResourceIdWithEscaped(id, 4)
	if err != nil {
		return nil, WrapError(err)
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
	if IsExpectedErrors(err, []string{"InternalError"}) {
		return object, WrapErrorf(Error(GetNotFoundMessage("MSE:MseNacosConfig", id)), NotFoundMsg, ProviderERROR)
	}

	if err != nil {
		if IsExpectedErrors(err, []string{"mse-200-021"}) {
			err = WrapErrorf(Error(GetNotFoundMessage("MseNacosConfig", id)), NotFoundMsg, ProviderERROR)
			return object, err
		}
		err = WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
		return object, err
	}
	if fmt.Sprint(response["Success"]) == "false" {
		return object, WrapError(fmt.Errorf("%s failed, response: %v", action, response))
	}
	addDebug(action, response, request)
	v, err := jsonpath.Get("$.Configuration", response)
	if err != nil {
		if err.Error() == "unknown key Configuration" {
			return object, WrapErrorf(Error(GetNotFoundMessage("MSE", id)), NotFoundWithResponse, response)
		}
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Configuration", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}

func (s *MseService) DescribeMseEngineNamespace(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		err = WrapError(err)
		return
	}
	action := "ListEngineNamespaces"
	request := map[string]interface{}{
		"InstanceId": parts[0],
	}
	idExist := false
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = client.RpcGet("mse", "2019-05-31", action, request, nil)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	// There is an API bug while query an Instance deleted.
	if IsExpectedErrors(err, []string{"InternalError"}) {
		return object, WrapErrorf(Error(GetNotFoundMessage("MSE:EngineNamespace", id)), NotFoundMsg, ProviderERROR)
	}
	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidParameter"}) {
			return object, WrapErrorf(Error(GetNotFoundMessage("MSE cluster", parts[0])), NotFoundMsg, "")
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	if fmt.Sprint(response["Success"]) == "false" {
		return object, WrapErrorf(Error(GetNotFoundMessage("MSE:EngineNamespace", id)), NotFoundMsg, ProviderERROR)
	}
	v, err := jsonpath.Get("$.Data", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Data", response)
	}
	if len(v.([]interface{})) < 1 {
		return object, WrapErrorf(Error(GetNotFoundMessage("MSE", id)), NotFoundWithResponse, response)
	}
	for _, v := range v.([]interface{}) {
		if fmt.Sprint(v.(map[string]interface{})["Namespace"]) == parts[1] {
			idExist = true
			return v.(map[string]interface{}), nil
		}
	}
	if !idExist {
		return object, WrapErrorf(Error(GetNotFoundMessage("MSE", id)), NotFoundWithResponse, response)
	}
	return object, nil
}

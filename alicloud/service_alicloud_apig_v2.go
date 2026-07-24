package alicloud

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/tidwall/sjson"
)

type ApigServiceV2 struct {
	client *connectivity.AliyunClient
}

// DescribeApigHttpApi <<< Encapsulated get interface for Apig HttpApi.

func (s *ApigServiceV2) DescribeApigHttpApi(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	httpApiId := id
	request = make(map[string]interface{})
	query = make(map[string]*string)

	action := fmt.Sprintf("/v1/http-apis/%s", httpApiId)

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(1*time.Minute, func() *resource.RetryError {
		response, err = client.RoaGet("APIG", "2024-03-27", action, query, nil, nil)

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
		if IsExpectedErrors(err, []string{"DatabaseError.RecordNotFound"}) {
			return object, WrapErrorf(NotFoundErr("HttpApi", id), NotFoundMsg, err)
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	code, _ := jsonpath.Get("$.code", response)
	if InArray(fmt.Sprint(code), []string{"DatabaseError.RecordNotFound"}) {
		return object, WrapErrorf(NotFoundErr("HttpApi", id), NotFoundMsg, response)
	}

	v, err := jsonpath.Get("$.data", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.data", response)
	}

	return v.(map[string]interface{}), nil
}

func (s *ApigServiceV2) ApigHttpApiStateRefreshFunc(id string, field string, failStates []string) resource.StateRefreshFunc {
	return s.ApigHttpApiStateRefreshFuncWithApi(id, field, failStates, s.DescribeApigHttpApi)
}

func (s *ApigServiceV2) ApigHttpApiStateRefreshFuncWithApi(id string, field string, failStates []string, call func(id string) (map[string]interface{}, error)) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := call(id)
		if err != nil {
			if NotFoundError(err) {
				return object, "", nil
			}
			return nil, "", WrapError(err)
		}
		v, err := jsonpath.Get(field, object)
		currentStatus := fmt.Sprint(v)

		if strings.HasPrefix(field, "#") {
			v, _ := jsonpath.Get(strings.TrimPrefix(field, "#"), object)
			if v != nil {
				currentStatus = "#CHECKSET"
			}
		}

		for _, failState := range failStates {
			if currentStatus == failState {
				return object, currentStatus, WrapError(Error(FailedToReachTargetStatus, currentStatus))
			}
		}
		return object, currentStatus, nil
	}
}

// DescribeApigHttpApi >>> Encapsulated.

// DescribeApigDomain <<< Encapsulated get interface for Apig Domain.

func (s *ApigServiceV2) DescribeApigDomain(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	domainId := id
	request = make(map[string]interface{})
	query = make(map[string]*string)

	action := fmt.Sprintf("/v1/domains/%s", domainId)

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(1*time.Minute, func() *resource.RetryError {
		response, err = client.RoaGet("APIG", "2024-03-27", action, query, nil, nil)

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
	if response == nil {
		return object, WrapErrorf(NotFoundErr("Domain", id), NotFoundMsg, response)
	}
	code, _ := jsonpath.Get("$.code", response)
	if InArray(fmt.Sprint(code), []string{"DatabaseError.RecordNotFound"}) {
		return object, WrapErrorf(NotFoundErr("Domain", id), NotFoundMsg, response)
	}

	v, err := jsonpath.Get("$.data", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.data", response)
	}

	return v.(map[string]interface{}), nil
}

func (s *ApigServiceV2) ApigDomainStateRefreshFunc(id string, field string, failStates []string) resource.StateRefreshFunc {
	return s.ApigDomainStateRefreshFuncWithApi(id, field, failStates, s.DescribeApigDomain)
}

func (s *ApigServiceV2) ApigDomainStateRefreshFuncWithApi(id string, field string, failStates []string, call func(id string) (map[string]interface{}, error)) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := call(id)
		if err != nil {
			if NotFoundError(err) {
				return object, "", nil
			}
			return nil, "", WrapError(err)
		}
		v, err := jsonpath.Get(field, object)
		currentStatus := fmt.Sprint(v)

		if strings.HasPrefix(field, "#") {
			v, _ := jsonpath.Get(strings.TrimPrefix(field, "#"), object)
			if v != nil {
				currentStatus = "#CHECKSET"
			}
		}

		for _, failState := range failStates {
			if currentStatus == failState {
				return object, currentStatus, WrapError(Error(FailedToReachTargetStatus, currentStatus))
			}
		}
		return object, currentStatus, nil
	}
}

// DescribeApigDomain >>> Encapsulated.

// DescribeApigGateway <<< Encapsulated get interface for Apig Gateway.

func (s *ApigServiceV2) DescribeApigGateway(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	gatewayId := id
	request = make(map[string]interface{})
	query = make(map[string]*string)

	action := fmt.Sprintf("/v1/gateways/%s", gatewayId)

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(1*time.Minute, func() *resource.RetryError {
		response, err = client.RoaGet("APIG", "2024-03-27", action, query, nil, nil)

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
	if response == nil {
		return object, WrapErrorf(NotFoundErr("Gateway", id), NotFoundMsg, response)
	}
	code, _ := jsonpath.Get("$.code", response)
	if InArray(fmt.Sprint(code), []string{"NotFound.GatewayNotFound", "Conflict.GatewayIsDeleted"}) {
		return object, WrapErrorf(NotFoundErr("Gateway", id), NotFoundMsg, response)
	}

	v, err := jsonpath.Get("$.data", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.data", response)
	}

	return v.(map[string]interface{}), nil
}

func (s *ApigServiceV2) ApigGatewayStateRefreshFunc(id string, field string, failStates []string) resource.StateRefreshFunc {
	return s.ApigGatewayStateRefreshFuncWithApi(id, field, failStates, s.DescribeApigGateway)
}

func (s *ApigServiceV2) ApigGatewayStateRefreshFuncWithApi(id string, field string, failStates []string, call func(id string) (map[string]interface{}, error)) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := call(id)
		if err != nil {
			if NotFoundError(err) {
				return object, "", nil
			}
			return nil, "", WrapError(err)
		}
		object["chargeType"] = convertApigGatewaydatachargeTypeResponse(object["chargeType"])
		v, err := jsonpath.Get(field, object)
		currentStatus := fmt.Sprint(v)

		if strings.HasPrefix(field, "#") {
			v, _ := jsonpath.Get(strings.TrimPrefix(field, "#"), object)
			if v != nil {
				currentStatus = "#CHECKSET"
			}
		}

		for _, failState := range failStates {
			if currentStatus == failState {
				return object, currentStatus, WrapError(Error(FailedToReachTargetStatus, currentStatus))
			}
		}
		return object, currentStatus, nil
	}
}

// DescribeApigGateway >>> Encapsulated.

// SetResourceTags <<< Encapsulated tag function for Apig.
func (s *ApigServiceV2) SetResourceTags(d *schema.ResourceData, resourceType string) error {
	if d.HasChange("tags") {
		var action string
		var err error
		client := s.client
		var request map[string]interface{}
		var response map[string]interface{}
		query := make(map[string]*string)
		body := make(map[string]interface{})

		added, removed := parsingTags(d)
		removedTagKeys := make([]string, 0)
		for _, v := range removed {
			if !ignoredTags(v, "") {
				removedTagKeys = append(removedTagKeys, v)
			}
		}
		if len(removedTagKeys) > 0 {
			action = fmt.Sprintf("/v1/tags")
			request = make(map[string]interface{})
			query = make(map[string]*string)
			body = make(map[string]interface{})
			query["ResourceId"] = StringPointer(convertListToJsonString(expandSingletonToList(d.Id())))
			query["RegionId"] = StringPointer(client.RegionId)
			query["TagKey"] = StringPointer(convertListToJsonString(convertListStringToListInterface(removedTagKeys)))
			query["ResourceType"] = StringPointer(resourceType)
			body = request
			wait := incrementalWait(3*time.Second, 5*time.Second)
			err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
				response, err = client.RoaDelete("APIG", "2024-03-27", action, query, nil, nil, true)
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
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
			}

		}

		if len(added) > 0 {
			action = fmt.Sprintf("/v1/tags")
			request = make(map[string]interface{})
			query = make(map[string]*string)
			body = make(map[string]interface{})

			count := 1
			tagsMaps := make([]map[string]interface{}, 0)
			for key, value := range added {
				tagsMap := make(map[string]interface{})
				tagsMap["value"] = value
				tagsMap["key"] = key
				tagsMaps = append(tagsMaps, tagsMap)
				count++
			}
			request["tag"] = tagsMaps

			request["resourceType"] = "gateway"
			jsonString := convertObjectToJsonString(request)
			jsonString, _ = sjson.Set(jsonString, "resourceId.0", d.Id())
			_ = json.Unmarshal([]byte(jsonString), &request)

			body = request
			wait := incrementalWait(3*time.Second, 5*time.Second)
			err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
				response, err = client.RoaPost("APIG", "2024-03-27", action, query, nil, body, true)
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
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
			}

		}
	}

	return nil
}

// SetResourceTags >>> tag function encapsulated.

// DescribeApigEnvironment <<< Encapsulated get interface for Apig Environment.

func (s *ApigServiceV2) DescribeApigEnvironment(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	environmentId := id
	action := fmt.Sprintf("/v1/environments/%s", environmentId)
	request = make(map[string]interface{})
	query = make(map[string]*string)
	request["environmentId"] = id

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(1*time.Minute, func() *resource.RetryError {
		response, err = client.RoaGet("APIG", "2024-03-27", action, query, nil, nil)

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
		if IsExpectedErrors(err, []string{"NotFound.EnvironmentNotFound"}) {
			return object, WrapErrorf(NotFoundErr("Environment", id), NotFoundMsg, err)
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$.data", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.data", response)
	}

	return v.(map[string]interface{}), nil
}

func (s *ApigServiceV2) ApigEnvironmentStateRefreshFunc(id string, field string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeApigEnvironment(id)
		if err != nil {
			if NotFoundError(err) {
				return object, "", nil
			}
			return nil, "", WrapError(err)
		}

		v, err := jsonpath.Get(field, object)
		currentStatus := fmt.Sprint(v)

		if strings.HasPrefix(field, "#") {
			v, _ := jsonpath.Get(strings.TrimPrefix(field, "#"), object)
			if v != nil {
				currentStatus = "#CHECKSET"
			}
		}

		for _, failState := range failStates {
			if currentStatus == failState {
				return object, currentStatus, WrapError(Error(FailedToReachTargetStatus, currentStatus))
			}
		}
		return object, currentStatus, nil
	}
}

// DescribeApigEnvironment >>> Encapsulated.

// DescribeApigService <<< Encapsulated get interface for Apig Service.

func (s *ApigServiceV2) DescribeApigService(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	serviceId := id
	request = make(map[string]interface{})
	query = make(map[string]*string)

	action := fmt.Sprintf("/v1/services/%s", serviceId)

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(1*time.Minute, func() *resource.RetryError {
		response, err = client.RoaGet("APIG", "2024-03-27", action, query, nil, nil)

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
	if response == nil {
		return object, WrapErrorf(NotFoundErr("Service", id), NotFoundMsg, response)
	}
	if err != nil {
		if IsExpectedErrors(err, []string{"NotFound.ServiceNotFound"}) {
			return object, WrapErrorf(NotFoundErr("Service", id), NotFoundMsg, response)
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	code, _ := jsonpath.Get("$.code", response)
	if InArray(fmt.Sprint(code), []string{"NotFound.ServiceNotFound"}) {
		return object, WrapErrorf(NotFoundErr("Service", id), NotFoundMsg, response)
	}

	v, err := jsonpath.Get("$.data", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.data", response)
	}

	return v.(map[string]interface{}), nil
}

func (s *ApigServiceV2) ApigServiceStateRefreshFunc(id string, field string, failStates []string) resource.StateRefreshFunc {
	return s.ApigServiceStateRefreshFuncWithApi(id, field, failStates, s.DescribeApigService)
}

func (s *ApigServiceV2) ApigServiceStateRefreshFuncWithApi(id string, field string, failStates []string, call func(id string) (map[string]interface{}, error)) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := call(id)
		if err != nil {
			if NotFoundError(err) {
				return object, "", nil
			}
			return nil, "", WrapError(err)
		}
		v, err := jsonpath.Get(field, object)
		currentStatus := fmt.Sprint(v)

		if strings.HasPrefix(field, "#") {
			v, _ := jsonpath.Get(strings.TrimPrefix(field, "#"), object)
			if v != nil {
				currentStatus = "#CHECKSET"
			}
		}

		for _, failState := range failStates {
			if currentStatus == failState {
				return object, currentStatus, WrapError(Error(FailedToReachTargetStatus, currentStatus))
			}
		}
		return object, currentStatus, nil
	}
}

// DescribeApigService >>> Encapsulated.

// DescribeApigPlugin <<< Encapsulated get interface for Apig Plugin.

func (s *ApigServiceV2) DescribeApigPlugin(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	request = make(map[string]interface{})
	query = make(map[string]*string)

	action := fmt.Sprintf("/v1/plugins")

	pageNumber := 1
	query["pageSize"] = StringPointer(fmt.Sprintf("%d", PageSizeLarge))
	query["pageNumber"] = StringPointer(fmt.Sprintf("%d", pageNumber))

	for {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(1*time.Minute, func() *resource.RetryError {
			response, err = client.RoaGet("APIG", "2024-03-27", action, query, nil, nil)

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
		if response == nil {
			return object, WrapErrorf(NotFoundErr("Plugin", id), NotFoundMsg, response)
		}
		code, _ := jsonpath.Get("$.code", response)
		if InArray(fmt.Sprint(code), []string{"DatabaseError.RecordNotFound"}) {
			return object, WrapErrorf(NotFoundErr("Plugin", id), NotFoundMsg, response)
		}

		v, err := jsonpath.Get("$.data.items[*]", response)
		if err != nil {
			return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.data.items[*]", response)
		}

		result, _ := v.([]interface{})
		for _, vv := range result {
			item := vv.(map[string]interface{})
			if fmt.Sprint(item["pluginId"]) == id {
				return item, nil
			}
		}

		if len(result) < PageSizeLarge {
			break
		}
		pageNumber += 1
		query["pageNumber"] = StringPointer(fmt.Sprintf("%d", pageNumber))
	}

	return object, WrapErrorf(NotFoundErr("Plugin", id), NotFoundMsg, response)
}

func (s *ApigServiceV2) ApigPluginStateRefreshFunc(id string, field string, failStates []string) resource.StateRefreshFunc {
	return s.ApigPluginStateRefreshFuncWithApi(id, field, failStates, s.DescribeApigPlugin)
}

func (s *ApigServiceV2) ApigPluginStateRefreshFuncWithApi(id string, field string, failStates []string, call func(id string) (map[string]interface{}, error)) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := call(id)
		if err != nil {
			if NotFoundError(err) {
				return object, "", nil
			}
			return nil, "", WrapError(err)
		}
		v, err := jsonpath.Get(field, object)
		currentStatus := fmt.Sprint(v)

		if strings.HasPrefix(field, "#") {
			v, _ := jsonpath.Get(strings.TrimPrefix(field, "#"), object)
			if v != nil {
				currentStatus = "#CHECKSET"
			}
		}

		for _, failState := range failStates {
			if currentStatus == failState {
				return object, currentStatus, WrapError(Error(FailedToReachTargetStatus, currentStatus))
			}
		}
		return object, currentStatus, nil
	}
}

// DescribeApigPlugin >>> Encapsulated.
// DescribeApigPluginClass <<< Encapsulated get interface for Apig PluginClass.

func (s *ApigServiceV2) DescribeApigPluginClass(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	pluginClassId := id
	action := fmt.Sprintf("/v1/plugin-classes/%s", pluginClassId)
	request = make(map[string]interface{})
	query = make(map[string]*string)
	request["pluginClassId"] = id

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(1*time.Minute, func() *resource.RetryError {
		response, err = client.RoaGet("APIG", "2024-03-27", action, query, nil, nil)

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
		if IsExpectedErrors(err, []string{"DatabaseError.RecordNotFound"}) {
			return object, WrapErrorf(NotFoundErr("PluginClass", id), NotFoundMsg, err)
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$.data", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.data", response)
	}

	return v.(map[string]interface{}), nil
}

func (s *ApigServiceV2) ApigPluginClassStateRefreshFunc(id string, field string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeApigPluginClass(id)
		if err != nil {
			if NotFoundError(err) {
				return object, "", nil
			}
			return nil, "", WrapError(err)
		}

		v, err := jsonpath.Get(field, object)
		currentStatus := fmt.Sprint(v)

		if strings.HasPrefix(field, "#") {
			v, _ := jsonpath.Get(strings.TrimPrefix(field, "#"), object)
			if v != nil {
				currentStatus = "#CHECKSET"
			}
		}

		for _, failState := range failStates {
			if currentStatus == failState {
				return object, currentStatus, WrapError(Error(FailedToReachTargetStatus, currentStatus))
			}
		}
		return object, currentStatus, nil
	}
}

// DescribeApigPluginClass >>> Encapsulated.
// DescribeApigOperation <<< Encapsulated get interface for Apig Operation.

func (s *ApigServiceV2) DescribeApigOperation(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	parts := strings.Split(id, ":")
	if len(parts) != 2 {
		err = WrapError(fmt.Errorf("invalid Resource Id %s. Expected parts' length %d, got %d", id, 2, len(parts)))
	}
	httpApiId := parts[0]
	operationId := parts[1]
	action := fmt.Sprintf("/v1/http-apis/%s/operations/%s", httpApiId, operationId)
	request = make(map[string]interface{})
	query = make(map[string]*string)

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(1*time.Minute, func() *resource.RetryError {
		response, err = client.RoaGet("APIG", "2024-03-27", action, query, nil, nil)

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
		if IsExpectedErrors(err, []string{"DatabaseError.RecordNotFound"}) {
			return object, WrapErrorf(NotFoundErr("Operation", id), NotFoundMsg, err)
		}
		return object, WrapErrorf(NotFoundErr("Operation", id), NotFoundMsg, response)
	}

	v, err := jsonpath.Get("$.data", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.data", response)
	}

	return v.(map[string]interface{}), nil
}

func (s *ApigServiceV2) ApigOperationStateRefreshFunc(id string, field string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeApigOperation(id)
		if err != nil {
			if NotFoundError(err) {
				return object, "", nil
			}
			return nil, "", WrapError(err)
		}

		v, err := jsonpath.Get(field, object)
		currentStatus := fmt.Sprint(v)

		if strings.HasPrefix(field, "#") {
			v, _ := jsonpath.Get(strings.TrimPrefix(field, "#"), object)
			if v != nil {
				currentStatus = "#CHECKSET"
			}
		}

		for _, failState := range failStates {
			if currentStatus == failState {
				return object, currentStatus, WrapError(Error(FailedToReachTargetStatus, currentStatus))
			}
		}
		return object, currentStatus, nil
	}
}

// DescribeApigOperation >>> Encapsulated.
// DescribeApigApiAttachment <<< Encapsulated get interface for Apig ApiAttachment.

func (s *ApigServiceV2) DescribeApigApiAttachment(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	parts := strings.Split(id, ":")
	if len(parts) != 3 {
		err = WrapError(fmt.Errorf("invalid Resource Id %s. Expected parts' length %d, got %d", id, 3, len(parts)))
	}
	httpApiId := parts[0]
	action := fmt.Sprintf("/v1/http-apis/%s/attachment", httpApiId)
	request = make(map[string]interface{})
	query = make(map[string]*string)
	query["environmentId"] = StringPointer(parts[2])
	query["routeId"] = StringPointer(parts[1])

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(1*time.Minute, func() *resource.RetryError {
		response, err = client.RoaGet("APIG", "2024-03-27", action, query, nil, nil)

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
	if response == nil {
		return object, WrapErrorf(NotFoundErr("Gateway", id), NotFoundMsg, err)
	}

	v, err := jsonpath.Get("$.data", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.data", response)
	}

	currentStatus := v.(map[string]interface{})["deployStatus"]
	if currentStatus == "NotDeployed" {
		return object, WrapErrorf(NotFoundErr("ApiAttachment", id), NotFoundMsg, response)
	}

	return v.(map[string]interface{}), nil
}

func (s *ApigServiceV2) ApigApiAttachmentStateRefreshFunc(id string, field string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeApigApiAttachment(id)
		if err != nil {
			if NotFoundError(err) {
				return object, "", nil
			}
			return nil, "", WrapError(err)
		}

		v, err := jsonpath.Get(field, object)
		currentStatus := fmt.Sprint(v)

		if strings.HasPrefix(field, "#") {
			v, _ := jsonpath.Get(strings.TrimPrefix(field, "#"), object)
			if v != nil {
				currentStatus = "#CHECKSET"
			}
		}

		for _, failState := range failStates {
			if currentStatus == failState {
				return object, currentStatus, WrapError(Error(FailedToReachTargetStatus, currentStatus))
			}
		}
		return object, currentStatus, nil
	}
}

// DescribeApigApiAttachment >>> Encapsulated.
// DescribeApigRoute <<< Encapsulated get interface for Apig Route.

func (s *ApigServiceV2) DescribeApigRoute(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	parts := strings.Split(id, ":")
	if len(parts) != 2 {
		err = WrapError(fmt.Errorf("invalid Resource Id %s. Expected parts' length %d, got %d", id, 2, len(parts)))
	}
	httpApiId := parts[0]
	routeId := parts[1]
	action := fmt.Sprintf("/v1/http-apis/%s/routes/%s", httpApiId, routeId)
	request = make(map[string]interface{})
	query = make(map[string]*string)

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(1*time.Minute, func() *resource.RetryError {
		response, err = client.RoaGet("APIG", "2024-03-27", action, query, nil, nil)

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
		if IsExpectedErrors(err, []string{"NotFound.RouteNotFound"}) {
			return object, WrapErrorf(NotFoundErr("Route", id), NotFoundMsg, err)
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.data", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.data", response)
	}

	return v.(map[string]interface{}), nil
}

func (s *ApigServiceV2) ApigRouteStateRefreshFunc(id string, field string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeApigRoute(id)
		if err != nil {
			if NotFoundError(err) {
				return object, "", nil
			}
			return nil, "", WrapError(err)
		}

		v, err := jsonpath.Get(field, object)
		currentStatus := fmt.Sprint(v)

		if strings.HasPrefix(field, "#") {
			v, _ := jsonpath.Get(strings.TrimPrefix(field, "#"), object)
			if v != nil {
				currentStatus = "#CHECKSET"
			}
		}

		for _, failState := range failStates {
			if currentStatus == failState {
				return object, currentStatus, WrapError(Error(FailedToReachTargetStatus, currentStatus))
			}
		}
		return object, currentStatus, nil
	}
}

// DescribeApigRoute >>> Encapsulated.

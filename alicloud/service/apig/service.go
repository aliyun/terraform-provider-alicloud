package apig

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/convert"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/tags"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tidwall/sjson"
)

type Service struct {
	client *connectivity.AliyunClient
}

func NewService(client *connectivity.AliyunClient) *Service {
	return &Service{client: client}
}

// DescribeApigHttpApi <<< Encapsulated get interface for Apig HttpApi.

func (s *Service) DescribeApigHttpApi(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	httpApiId := id
	request = make(map[string]interface{})
	query = make(map[string]*string)

	action := fmt.Sprintf("/v1/http-apis/%s", httpApiId)

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.RoaGet("APIG", "2024-03-27", action, query, nil, nil)

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
		if errs.IsExpectedErrors(err, []string{"DatabaseError.RecordNotFound"}) {
			return object, errs.WrapErrorf(errs.NotFoundErr("HttpApi", id), errs.NotFoundMsg, err)
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}
	code, _ := jsonpath.Get("$.code", response)
	if helper.InArray(fmt.Sprint(code), []string{"DatabaseError.RecordNotFound"}) {
		return object, errs.WrapErrorf(errs.NotFoundErr("HttpApi", id), errs.NotFoundMsg, response)
	}

	v, err := jsonpath.Get("$.data", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.data", response)
	}

	return v.(map[string]interface{}), nil
}

func (s *Service) ApigHttpApiStateRefreshFunc(id string, field string, failStates []string) retry.StateRefreshFunc {
	return s.ApigHttpApiStateRefreshFuncWithApi(id, field, failStates, s.DescribeApigHttpApi)
}

func (s *Service) ApigHttpApiStateRefreshFuncWithApi(id string, field string, failStates []string, call func(id string) (map[string]interface{}, error)) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := call(id)
		if err != nil {
			if errs.NotFoundError(err) {
				return object, "", nil
			}
			return nil, "", errs.WrapError(err)
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
				return object, currentStatus, errs.WrapError(errs.Error(errs.FailedToReachTargetStatus, currentStatus))
			}
		}
		return object, currentStatus, nil
	}
}

// DescribeApigHttpApi >>> Encapsulated.

// DescribeApigDomain <<< Encapsulated get interface for Apig Domain.

func (s *Service) DescribeApigDomain(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	domainId := id
	request = make(map[string]interface{})
	query = make(map[string]*string)

	action := fmt.Sprintf("/v1/domains/%s", domainId)

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.RoaGet("APIG", "2024-03-27", action, query, nil, nil)

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
	if response == nil {
		return object, errs.WrapErrorf(errs.NotFoundErr("Domain", id), errs.NotFoundMsg, response)
	}
	code, _ := jsonpath.Get("$.code", response)
	if helper.InArray(fmt.Sprint(code), []string{"DatabaseError.RecordNotFound"}) {
		return object, errs.WrapErrorf(errs.NotFoundErr("Domain", id), errs.NotFoundMsg, response)
	}

	v, err := jsonpath.Get("$.data", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.data", response)
	}

	return v.(map[string]interface{}), nil
}

func (s *Service) ApigDomainStateRefreshFunc(id string, field string, failStates []string) retry.StateRefreshFunc {
	return s.ApigDomainStateRefreshFuncWithApi(id, field, failStates, s.DescribeApigDomain)
}

func (s *Service) ApigDomainStateRefreshFuncWithApi(id string, field string, failStates []string, call func(id string) (map[string]interface{}, error)) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := call(id)
		if err != nil {
			if errs.NotFoundError(err) {
				return object, "", nil
			}
			return nil, "", errs.WrapError(err)
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
				return object, currentStatus, errs.WrapError(errs.Error(errs.FailedToReachTargetStatus, currentStatus))
			}
		}
		return object, currentStatus, nil
	}
}

// DescribeApigDomain >>> Encapsulated.

// DescribeApigGateway <<< Encapsulated get interface for Apig Gateway.

func (s *Service) DescribeApigGateway(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	gatewayId := id
	request = make(map[string]interface{})
	query = make(map[string]*string)

	action := fmt.Sprintf("/v1/gateways/%s", gatewayId)

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.RoaGet("APIG", "2024-03-27", action, query, nil, nil)

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
	if response == nil {
		return object, errs.WrapErrorf(errs.NotFoundErr("Gateway", id), errs.NotFoundMsg, response)
	}
	code, _ := jsonpath.Get("$.code", response)
	if helper.InArray(fmt.Sprint(code), []string{"NotFound.GatewayNotFound", "Conflict.GatewayIsDeleted"}) {
		return object, errs.WrapErrorf(errs.NotFoundErr("Gateway", id), errs.NotFoundMsg, response)
	}

	v, err := jsonpath.Get("$.data", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.data", response)
	}

	return v.(map[string]interface{}), nil
}

func (s *Service) ApigGatewayStateRefreshFunc(id string, field string, failStates []string) retry.StateRefreshFunc {
	return s.ApigGatewayStateRefreshFuncWithApi(id, field, failStates, s.DescribeApigGateway)
}

func (s *Service) ApigGatewayStateRefreshFuncWithApi(id string, field string, failStates []string, call func(id string) (map[string]interface{}, error)) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := call(id)
		if err != nil {
			if errs.NotFoundError(err) {
				return object, "", nil
			}
			return nil, "", errs.WrapError(err)
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
				return object, currentStatus, errs.WrapError(errs.Error(errs.FailedToReachTargetStatus, currentStatus))
			}
		}
		return object, currentStatus, nil
	}
}

// DescribeApigGateway >>> Encapsulated.

// SetResourceTags <<< Encapsulated tag function for Apig.
func (s *Service) SetResourceTags(d *schema.ResourceData, resourceType string) error {
	if d.HasChange("tags") {
		var action string
		var err error
		client := s.client
		var request map[string]interface{}
		var response map[string]interface{}
		query := make(map[string]*string)
		body := make(map[string]interface{})

		added, removed := tags.ParsingTags(d)
		removedTagKeys := make([]string, 0)
		for _, v := range removed {
			if !tags.IgnoredTags(v, "") {
				removedTagKeys = append(removedTagKeys, v)
			}
		}
		if len(removedTagKeys) > 0 {
			action = fmt.Sprintf("/v1/tags")
			request = make(map[string]interface{})
			query = make(map[string]*string)
			body = make(map[string]interface{})
			query["ResourceId"] = helper.StringPointer(convert.ConvertListToJsonString(convert.ExpandSingletonToList(d.Id())))
			query["RegionId"] = helper.StringPointer(client.RegionId)
			query["TagKey"] = helper.StringPointer(convert.ConvertListToJsonString(convert.ConvertListStringToListInterface(removedTagKeys)))
			query["ResourceType"] = helper.StringPointer(resourceType)
			body = request
			wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
			err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
				response, err = client.RoaDelete("APIG", "2024-03-27", action, query, nil, nil, true)
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
				return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
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
			jsonString := convert.ConvertObjectToJsonString(request)
			jsonString, _ = sjson.Set(jsonString, "resourceId.0", d.Id())
			_ = json.Unmarshal([]byte(jsonString), &request)

			body = request
			wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
			err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
				response, err = client.RoaPost("APIG", "2024-03-27", action, query, nil, body, true)
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
				return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
			}

		}
	}

	return nil
}

// SetResourceTags >>> tag function encapsulated.

// DescribeApigEnvironment <<< Encapsulated get interface for Apig Environment.

func (s *Service) DescribeApigEnvironment(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	environmentId := id
	action := fmt.Sprintf("/v1/environments/%s", environmentId)
	request = make(map[string]interface{})
	query = make(map[string]*string)
	request["environmentId"] = id

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.RoaGet("APIG", "2024-03-27", action, query, nil, nil)

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
		if errs.IsExpectedErrors(err, []string{"NotFound.EnvironmentNotFound"}) {
			return object, errs.WrapErrorf(errs.NotFoundErr("Environment", id), errs.NotFoundMsg, err)
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$.data", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.data", response)
	}

	return v.(map[string]interface{}), nil
}

func (s *Service) ApigEnvironmentStateRefreshFunc(id string, field string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeApigEnvironment(id)
		if err != nil {
			if errs.NotFoundError(err) {
				return object, "", nil
			}
			return nil, "", errs.WrapError(err)
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
				return object, currentStatus, errs.WrapError(errs.Error(errs.FailedToReachTargetStatus, currentStatus))
			}
		}
		return object, currentStatus, nil
	}
}

// DescribeApigEnvironment >>> Encapsulated.

// DescribeApigService <<< Encapsulated get interface for Apig Service.

func (s *Service) DescribeApigService(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	serviceId := id
	request = make(map[string]interface{})
	query = make(map[string]*string)

	action := fmt.Sprintf("/v1/services/%s", serviceId)

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.RoaGet("APIG", "2024-03-27", action, query, nil, nil)

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
	if response == nil {
		return object, errs.WrapErrorf(errs.NotFoundErr("Service", id), errs.NotFoundMsg, response)
	}
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"NotFound.ServiceNotFound"}) {
			return object, errs.WrapErrorf(errs.NotFoundErr("Service", id), errs.NotFoundMsg, response)
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}
	code, _ := jsonpath.Get("$.code", response)
	if helper.InArray(fmt.Sprint(code), []string{"NotFound.ServiceNotFound"}) {
		return object, errs.WrapErrorf(errs.NotFoundErr("Service", id), errs.NotFoundMsg, response)
	}

	v, err := jsonpath.Get("$.data", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.data", response)
	}

	return v.(map[string]interface{}), nil
}

func (s *Service) ApigServiceStateRefreshFunc(id string, field string, failStates []string) retry.StateRefreshFunc {
	return s.ApigServiceStateRefreshFuncWithApi(id, field, failStates, s.DescribeApigService)
}

func (s *Service) ApigServiceStateRefreshFuncWithApi(id string, field string, failStates []string, call func(id string) (map[string]interface{}, error)) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := call(id)
		if err != nil {
			if errs.NotFoundError(err) {
				return object, "", nil
			}
			return nil, "", errs.WrapError(err)
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
				return object, currentStatus, errs.WrapError(errs.Error(errs.FailedToReachTargetStatus, currentStatus))
			}
		}
		return object, currentStatus, nil
	}
}

// DescribeApigService >>> Encapsulated.

// DescribeApigPlugin <<< Encapsulated get interface for Apig Plugin.

func (s *Service) DescribeApigPlugin(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	request = make(map[string]interface{})
	query = make(map[string]*string)

	action := fmt.Sprintf("/v1/plugins")

	gatewayTypes := []string{"API", "AI"}

	for _, gatewayType := range gatewayTypes {
		query["pageSize"] = helper.StringPointer(strconv.Itoa(helper.PageSizeLarge))
		query["pageNumber"] = helper.StringPointer("1")
		query["gatewayType"] = helper.StringPointer(gatewayType)
		for {
			wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
			err = retry.Retry(1*time.Minute, func() *retry.RetryError {
				response, err = client.RoaGet("APIG", "2024-03-27", action, query, nil, nil)
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
				if errs.IsExpectedErrors(err, []string{"DatabaseError.RecordNotFound"}) {
					return object, errs.WrapErrorf(errs.NotFoundErr("Plugin", id), errs.NotFoundMsg, err)
				}
				return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
			}

			resp, err := jsonpath.Get("$.data.items", response)
			if err != nil {
				return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.data.items", response)
			}

			for _, v := range resp.([]interface{}) {
				if fmt.Sprint(v.(map[string]interface{})["pluginId"]) == id {
					return v.(map[string]interface{}), nil
				}
			}

			if len(resp.([]interface{})) < helper.PageSizeLarge {
				break
			}

			pageNumber, err := strconv.Atoi(*query["pageNumber"])
			if err != nil {
				return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
			}

			query["pageNumber"] = helper.StringPointer(strconv.Itoa(pageNumber + 1))
		}
	}

	return nil, errs.WrapErrorf(errs.NotFoundErr("Plugin", id), errs.NotFoundMsg, response)
}

func (s *Service) ApigPluginStateRefreshFunc(id string, field string, failStates []string) retry.StateRefreshFunc {
	return s.ApigPluginStateRefreshFuncWithApi(id, field, failStates, s.DescribeApigPlugin)
}

func (s *Service) ApigPluginStateRefreshFuncWithApi(id string, field string, failStates []string, call func(id string) (map[string]interface{}, error)) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := call(id)
		if err != nil {
			if errs.NotFoundError(err) {
				return object, "", nil
			}
			return nil, "", errs.WrapError(err)
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
				return object, currentStatus, errs.WrapError(errs.Error(errs.FailedToReachTargetStatus, currentStatus))
			}
		}
		return object, currentStatus, nil
	}
}

// DescribeApigPlugin >>> Encapsulated.
// DescribeApigPluginClass <<< Encapsulated get interface for Apig PluginClass.

func (s *Service) DescribeApigPluginClass(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	pluginClassId := id
	action := fmt.Sprintf("/v1/plugin-classes/%s", pluginClassId)
	request = make(map[string]interface{})
	query = make(map[string]*string)
	request["pluginClassId"] = id

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.RoaGet("APIG", "2024-03-27", action, query, nil, nil)

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
		if errs.IsExpectedErrors(err, []string{"DatabaseError.RecordNotFound"}) {
			return object, errs.WrapErrorf(errs.NotFoundErr("PluginClass", id), errs.NotFoundMsg, err)
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$.data", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.data", response)
	}

	return v.(map[string]interface{}), nil
}

func (s *Service) ApigPluginClassStateRefreshFunc(id string, field string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeApigPluginClass(id)
		if err != nil {
			if errs.NotFoundError(err) {
				return object, "", nil
			}
			return nil, "", errs.WrapError(err)
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
				return object, currentStatus, errs.WrapError(errs.Error(errs.FailedToReachTargetStatus, currentStatus))
			}
		}
		return object, currentStatus, nil
	}
}

// DescribeApigPluginClass >>> Encapsulated.
// DescribeApigOperation <<< Encapsulated get interface for Apig Operation.

func (s *Service) DescribeApigOperation(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	parts := strings.Split(id, ":")
	if len(parts) != 2 {
		err = errs.WrapError(fmt.Errorf("invalid Resource Id %s. Expected parts' length %d, got %d", id, 2, len(parts)))
	}
	httpApiId := parts[0]
	operationId := parts[1]
	action := fmt.Sprintf("/v1/http-apis/%s/operations/%s", httpApiId, operationId)
	request = make(map[string]interface{})
	query = make(map[string]*string)

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.RoaGet("APIG", "2024-03-27", action, query, nil, nil)

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
		if errs.IsExpectedErrors(err, []string{"DatabaseError.RecordNotFound"}) {
			return object, errs.WrapErrorf(errs.NotFoundErr("Operation", id), errs.NotFoundMsg, err)
		}
		return object, errs.WrapErrorf(errs.NotFoundErr("Operation", id), errs.NotFoundMsg, response)
	}

	v, err := jsonpath.Get("$.data", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.data", response)
	}

	return v.(map[string]interface{}), nil
}

func (s *Service) ApigOperationStateRefreshFunc(id string, field string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeApigOperation(id)
		if err != nil {
			if errs.NotFoundError(err) {
				return object, "", nil
			}
			return nil, "", errs.WrapError(err)
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
				return object, currentStatus, errs.WrapError(errs.Error(errs.FailedToReachTargetStatus, currentStatus))
			}
		}
		return object, currentStatus, nil
	}
}

// DescribeApigOperation >>> Encapsulated.
// DescribeApigApiAttachment <<< Encapsulated get interface for Apig ApiAttachment.

func (s *Service) DescribeApigApiAttachment(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	parts := strings.Split(id, ":")
	if len(parts) != 3 {
		err = errs.WrapError(fmt.Errorf("invalid Resource Id %s. Expected parts' length %d, got %d", id, 3, len(parts)))
	}
	httpApiId := parts[0]
	action := fmt.Sprintf("/v1/http-apis/%s/attachment", httpApiId)
	request = make(map[string]interface{})
	query = make(map[string]*string)
	query["environmentId"] = helper.StringPointer(parts[2])
	query["routeId"] = helper.StringPointer(parts[1])

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.RoaGet("APIG", "2024-03-27", action, query, nil, nil)

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
	if response == nil {
		return object, errs.WrapErrorf(errs.NotFoundErr("Gateway", id), errs.NotFoundMsg, err)
	}

	v, err := jsonpath.Get("$.data", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.data", response)
	}

	currentStatus := v.(map[string]interface{})["deployStatus"]
	if currentStatus == "NotDeployed" {
		return object, errs.WrapErrorf(errs.NotFoundErr("ApiAttachment", id), errs.NotFoundMsg, response)
	}

	return v.(map[string]interface{}), nil
}

func (s *Service) ApigApiAttachmentStateRefreshFunc(id string, field string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeApigApiAttachment(id)
		if err != nil {
			if errs.NotFoundError(err) {
				return object, "", nil
			}
			return nil, "", errs.WrapError(err)
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
				return object, currentStatus, errs.WrapError(errs.Error(errs.FailedToReachTargetStatus, currentStatus))
			}
		}
		return object, currentStatus, nil
	}
}

// DescribeApigApiAttachment >>> Encapsulated.

// DescribeApigRoute <<< Encapsulated get interface for Apig Route.

func (s *Service) DescribeApigRoute(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	parts := strings.Split(id, ":")
	if len(parts) != 2 {
		err = errs.WrapError(fmt.Errorf("invalid Resource Id %s. Expected parts' length %d, got %d", id, 2, len(parts)))
		return nil, err
	}
	httpApiId := parts[0]
	routeId := parts[1]
	request = make(map[string]interface{})
	query = make(map[string]*string)

	action := fmt.Sprintf("/v1/http-apis/%s/routes/%s", httpApiId, routeId)

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.RoaGet("APIG", "2024-03-27", action, query, nil, nil)

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
	if response == nil {
		return object, errs.WrapErrorf(errs.NotFoundErr("Route", id), errs.NotFoundMsg, response)
	}
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"NotFound.RouteNotFound"}) {
			return object, errs.WrapErrorf(errs.NotFoundErr("Route", id), errs.NotFoundMsg, response)
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$.data", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.data", response)
	}

	return v.(map[string]interface{}), nil
}

func (s *Service) ApigRouteStateRefreshFunc(id string, field string, failStates []string) retry.StateRefreshFunc {
	return s.ApigRouteStateRefreshFuncWithApi(id, field, failStates, s.DescribeApigRoute)
}

func (s *Service) ApigRouteStateRefreshFuncWithApi(id string, field string, failStates []string, call func(id string) (map[string]interface{}, error)) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := call(id)
		if err != nil {
			if errs.NotFoundError(err) {
				return object, "", nil
			}
			return nil, "", errs.WrapError(err)
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
				return object, currentStatus, errs.WrapError(errs.Error(errs.FailedToReachTargetStatus, currentStatus))
			}
		}
		return object, currentStatus, nil
	}
}

// DescribeApigRoute >>> Encapsulated.

// DescribeApigAiModelProvider <<< Encapsulated get interface for Apig AiModelProvider.

func (s *Service) DescribeApigAiModelProvider(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	modelProviderId := id
	request = make(map[string]interface{})
	query = make(map[string]*string)

	action := fmt.Sprintf("/v1/ai-model-providers/%s", modelProviderId)

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.RoaGet("APIG", "2024-03-27", action, query, nil, nil)

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
	if response == nil {
		return object, errs.WrapErrorf(errs.NotFoundErr("AiModelProvider", id), errs.NotFoundMsg, response)
	}
	if err != nil {
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$.data", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.data", response)
	}

	return v.(map[string]interface{}), nil
}

func (s *Service) ApigAiModelProviderStateRefreshFunc(id string, field string, failStates []string) retry.StateRefreshFunc {
	return s.ApigAiModelProviderStateRefreshFuncWithApi(id, field, failStates, s.DescribeApigAiModelProvider)
}

func (s *Service) ApigAiModelProviderStateRefreshFuncWithApi(id string, field string, failStates []string, call func(id string) (map[string]interface{}, error)) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := call(id)
		if err != nil {
			if errs.NotFoundError(err) {
				return object, "", nil
			}
			return nil, "", errs.WrapError(err)
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
				return object, currentStatus, errs.WrapError(errs.Error(errs.FailedToReachTargetStatus, currentStatus))
			}
		}
		return object, currentStatus, nil
	}
}

// DescribeApigAiModelProvider >>> Encapsulated.

// DescribeApigPolicy <<< Encapsulated get interface for Apig Policy.

func (s *Service) DescribeApigPolicy(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	policyId := id
	request = make(map[string]interface{})
	query = make(map[string]*string)

	action := fmt.Sprintf("/v2/policies/%s", policyId)

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.RoaGet("APIG", "2024-03-27", action, query, nil, nil)

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
	if response == nil {
		return object, errs.WrapErrorf(errs.NotFoundErr("Policy", id), errs.NotFoundMsg, response)
	}
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"DatabaseError.RecordNotFound"}) {
			return object, errs.WrapErrorf(errs.NotFoundErr("Policy", id), errs.NotFoundMsg, response)
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}
	code, _ := jsonpath.Get("$.code", response)
	if helper.InArray(fmt.Sprint(code), []string{"DatabaseError.RecordNotFound"}) {
		return object, errs.WrapErrorf(errs.NotFoundErr("Policy", id), errs.NotFoundMsg, response)
	}

	v, err := jsonpath.Get("$.data", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.data", response)
	}

	return v.(map[string]interface{}), nil
}

func (s *Service) ApigPolicyStateRefreshFunc(id string, field string, failStates []string) retry.StateRefreshFunc {
	return s.ApigPolicyStateRefreshFuncWithApi(id, field, failStates, s.DescribeApigPolicy)
}

func (s *Service) ApigPolicyStateRefreshFuncWithApi(id string, field string, failStates []string, call func(id string) (map[string]interface{}, error)) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := call(id)
		if err != nil {
			if errs.NotFoundError(err) {
				return object, "", nil
			}
			return nil, "", errs.WrapError(err)
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
				return object, currentStatus, errs.WrapError(errs.Error(errs.FailedToReachTargetStatus, currentStatus))
			}
		}
		return object, currentStatus, nil
	}
}

// DescribeApigPolicy >>> Encapsulated.

// DescribeApigSource <<< Encapsulated get interface for Apig Source.

func (s *Service) DescribeApigSource(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	sourceId := id
	request = make(map[string]interface{})
	query = make(map[string]*string)

	action := fmt.Sprintf("/v1/sources/%s", sourceId)

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.RoaGet("APIG", "2024-03-27", action, query, nil, nil)

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
	if response == nil {
		return object, errs.WrapErrorf(errs.NotFoundErr("Source", id), errs.NotFoundMsg, response)
	}
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"NotFound.ServiceSourceNotFound"}) {
			return object, errs.WrapErrorf(errs.NotFoundErr("Source", id), errs.NotFoundMsg, response)
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}
	code, _ := jsonpath.Get("$.code", response)
	if helper.InArray(fmt.Sprint(code), []string{"NotFound.ServiceSourceNotFound"}) {
		return object, errs.WrapErrorf(errs.NotFoundErr("Source", id), errs.NotFoundMsg, response)
	}

	v, err := jsonpath.Get("$.data", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.data", response)
	}

	return v.(map[string]interface{}), nil
}

func (s *Service) ApigSourceStateRefreshFunc(id string, field string, failStates []string) retry.StateRefreshFunc {
	return s.ApigSourceStateRefreshFuncWithApi(id, field, failStates, s.DescribeApigSource)
}

func (s *Service) ApigSourceStateRefreshFuncWithApi(id string, field string, failStates []string, call func(id string) (map[string]interface{}, error)) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := call(id)
		if err != nil {
			if errs.NotFoundError(err) {
				return object, "", nil
			}
			return nil, "", errs.WrapError(err)
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
				return object, currentStatus, errs.WrapError(errs.Error(errs.FailedToReachTargetStatus, currentStatus))
			}
		}
		return object, currentStatus, nil
	}
}

// DescribeApigSource >>> Encapsulated.

// DescribeApigSecret <<< Encapsulated get interface for Apig Secret.

func (s *Service) DescribeApigSecret(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	secretId := id
	request = make(map[string]interface{})
	query = make(map[string]*string)

	action := fmt.Sprintf("/v1/secrets/%s", secretId)

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.RoaGet("APIG", "2024-03-27", action, query, nil, nil)

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
	if response == nil {
		return object, errs.WrapErrorf(errs.NotFoundErr("Secret", id), errs.NotFoundMsg, response)
	}
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"NotFound.Secret"}) {
			return object, errs.WrapErrorf(errs.NotFoundErr("Secret", id), errs.NotFoundMsg, response)
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$.data", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.data", response)
	}

	return v.(map[string]interface{}), nil
}

func (s *Service) ApigSecretStateRefreshFunc(id string, field string, failStates []string) retry.StateRefreshFunc {
	return s.ApigSecretStateRefreshFuncWithApi(id, field, failStates, s.DescribeApigSecret)
}

func (s *Service) ApigSecretStateRefreshFuncWithApi(id string, field string, failStates []string, call func(id string) (map[string]interface{}, error)) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := call(id)
		if err != nil {
			if errs.NotFoundError(err) {
				return object, "", nil
			}
			return nil, "", errs.WrapError(err)
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
				return object, currentStatus, errs.WrapError(errs.Error(errs.FailedToReachTargetStatus, currentStatus))
			}
		}
		return object, currentStatus, nil
	}
}

// DescribeApigSecret >>> Encapsulated.

package alicloud

import (
	"fmt"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	rpc "github.com/alibabacloud-go/tea-rpc/client"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

type ExpressConnectRouterServiceV2 struct {
	client *connectivity.AliyunClient
}

// DescribeExpressConnectRouterExpressConnectRouter <<< Encapsulated get interface for ExpressConnectRouter ExpressConnectRouter.

func (s *ExpressConnectRouterServiceV2) DescribeExpressConnectRouterExpressConnectRouter(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	action := "DescribeExpressConnectRouter"
	conn, err := client.NewExpressconnectrouterClient()
	if err != nil {
		return object, WrapError(err)
	}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["EcrId"] = id

	request["ClientToken"] = buildClientToken(action)

	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(1*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2023-09-01"), StringPointer("AK"), query, request, &runtime)
		request["ClientToken"] = buildClientToken(action)

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
		addDebug(action, response, request)
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$.EcrList[*]", response)
	if err != nil {
		return object, WrapErrorf(Error(GetNotFoundMessage("ExpressConnectRouter", id)), NotFoundMsg, response)
	}

	if len(v.([]interface{})) == 0 {
		return object, WrapErrorf(Error(GetNotFoundMessage("ExpressConnectRouter", id)), NotFoundMsg, response)
	}

	currentStatus := v.([]interface{})[0].(map[string]interface{})["AlibabaSideAsn"]
	if currentStatus == nil {
		return object, WrapErrorf(Error(GetNotFoundMessage("ExpressConnectRouter", id)), NotFoundMsg, response)
	}

	return v.([]interface{})[0].(map[string]interface{}), nil
}
func (s *ExpressConnectRouterServiceV2) DescribeDescribeInstanceGrantedToExpressConnectRouter(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	action := "DescribeInstanceGrantedToExpressConnectRouter"
	conn, err := client.NewExpressconnectrouterClient()
	if err != nil {
		return object, WrapError(err)
	}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["EcrId"] = id

	request["ClientToken"] = buildClientToken(action)

	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(1*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2023-09-01"), StringPointer("AK"), query, request, &runtime)
		request["ClientToken"] = buildClientToken(action)

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
	code, _ := jsonpath.Get("$.Code", response)
	if InArray(fmt.Sprint(code), []string{"ResourceNotFound.EcrId"}) {
		return object, WrapErrorf(Error(GetNotFoundMessage("ExpressConnectRouter", id)), NotFoundMsg, response)
	}

	return response, nil
}
func (s *ExpressConnectRouterServiceV2) DescribeDescribeExpressConnectRouterRouteEntries(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	action := "DescribeExpressConnectRouterRouteEntries"
	conn, err := client.NewExpressconnectrouterClient()
	if err != nil {
		return object, WrapError(err)
	}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["EcrId"] = id

	request["ClientToken"] = buildClientToken(action)

	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(1*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2023-09-01"), StringPointer("AK"), query, request, &runtime)
		request["ClientToken"] = buildClientToken(action)

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
	code, _ := jsonpath.Get("$.Code", response)
	if InArray(fmt.Sprint(code), []string{"ResourceNotFound.EcrId"}) {
		return object, WrapErrorf(Error(GetNotFoundMessage("ExpressConnectRouter", id)), NotFoundMsg, response)
	}

	return response, nil
}
func (s *ExpressConnectRouterServiceV2) DescribeDescribeDisabledExpressConnectRouterRouteEntries(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	action := "DescribeDisabledExpressConnectRouterRouteEntries"
	conn, err := client.NewExpressconnectrouterClient()
	if err != nil {
		return object, WrapError(err)
	}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["EcrId"] = id

	request["ClientToken"] = buildClientToken(action)

	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(1*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2023-09-01"), StringPointer("AK"), query, request, &runtime)
		request["ClientToken"] = buildClientToken(action)

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
	code, _ := jsonpath.Get("$.Code", response)
	if InArray(fmt.Sprint(code), []string{"ResourceNotFound.EcrId"}) {
		return object, WrapErrorf(Error(GetNotFoundMessage("ExpressConnectRouter", id)), NotFoundMsg, response)
	}

	return response, nil
}
func (s *ExpressConnectRouterServiceV2) DescribeDescribeExpressConnectRouterInterRegionTransitMode(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	action := "DescribeExpressConnectRouterInterRegionTransitMode"
	conn, err := client.NewExpressconnectrouterClient()
	if err != nil {
		return object, WrapError(err)
	}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["EcrId"] = id

	request["ClientToken"] = buildClientToken(action)

	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(1*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2023-09-01"), StringPointer("AK"), query, request, &runtime)
		request["ClientToken"] = buildClientToken(action)

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
	code, _ := jsonpath.Get("$.Code", response)
	if InArray(fmt.Sprint(code), []string{"ResourceNotFound.EcrId"}) {
		return object, WrapErrorf(Error(GetNotFoundMessage("ExpressConnectRouter", id)), NotFoundMsg, response)
	}

	return response, nil
}
func (s *ExpressConnectRouterServiceV2) DescribeDescribeExpressConnectRouterRegion(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	action := "DescribeExpressConnectRouterRegion"
	conn, err := client.NewExpressconnectrouterClient()
	if err != nil {
		return object, WrapError(err)
	}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["EcrId"] = id

	request["ClientToken"] = buildClientToken(action)

	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(1*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2023-09-01"), StringPointer("AK"), query, request, &runtime)
		request["ClientToken"] = buildClientToken(action)

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
	code, _ := jsonpath.Get("$.Code", response)
	if InArray(fmt.Sprint(code), []string{"ResourceNotFound.EcrId"}) {
		return object, WrapErrorf(Error(GetNotFoundMessage("ExpressConnectRouter", id)), NotFoundMsg, response)
	}

	return response, nil
}

func (s *ExpressConnectRouterServiceV2) ExpressConnectRouterExpressConnectRouterStateRefreshFunc(id string, field string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeExpressConnectRouterExpressConnectRouter(id)
		if err != nil {
			if NotFoundError(err) {
				return object, "", nil
			}
			return nil, "", WrapError(err)
		}

		v, err := jsonpath.Get(field, object)
		currentStatus := fmt.Sprint(v)

		for _, failState := range failStates {
			if currentStatus == failState {
				return object, currentStatus, WrapError(Error(FailedToReachTargetStatus, currentStatus))
			}
		}
		return object, currentStatus, nil
	}
}

func (s *ExpressConnectRouterServiceV2) DescribeAsyncDescribeInstanceGrantedToExpressConnectRouter(d *schema.ResourceData, res map[string]interface{}) (object map[string]interface{}, err error) {
	client := s.client
	id := d.Id()
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	action := "DescribeInstanceGrantedToExpressConnectRouter"
	conn, err := client.NewExpressconnectrouterClient()
	if err != nil {
		return object, WrapError(err)
	}
	request = make(map[string]interface{})
	query = make(map[string]interface{})

	request["EcrId"] = id
	request["ClientToken"] = buildClientToken(action)

	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(1*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2023-09-01"), StringPointer("AK"), query, request, &runtime)
		request["ClientToken"] = buildClientToken(action)

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
		if IsExpectedErrors(err, []string{"ResourceNotFound.EcrId"}) {
			return object, WrapErrorf(Error(GetNotFoundMessage("ExpressConnectRouter", id)), NotFoundMsg, response)
		}
		addDebug(action, response, request)
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}

	return response, nil
}

func (s *ExpressConnectRouterServiceV2) DescribeAsyncExpressConnectRouterExpressConnectRouterStateRefreshFunc(d *schema.ResourceData, res map[string]interface{}, field string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeAsyncDescribeInstanceGrantedToExpressConnectRouter(d, res)
		if err != nil {
			if NotFoundError(err) {
				return object, "", nil
			}
			return nil, "", WrapError(err)
		}

		v, err := jsonpath.Get(field, object)
		currentStatus := fmt.Sprint(v)

		for _, failState := range failStates {
			if currentStatus == failState {
				if _err, ok := object["error"]; ok {
					return _err, currentStatus, WrapError(Error(FailedToReachTargetStatus, currentStatus))
				}
				return object, currentStatus, WrapError(Error(FailedToReachTargetStatus, currentStatus))
			}
		}
		return object, currentStatus, nil
	}
}
func (s *ExpressConnectRouterServiceV2) DescribeAsyncGetExpressConnectRouter(d *schema.ResourceData, res map[string]interface{}) (object map[string]interface{}, err error) {
	client := s.client
	id := d.Id()
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	action := "GetExpressConnectRouter"
	conn, err := client.NewExpressconnectrouterClient()
	if err != nil {
		return object, WrapError(err)
	}
	request = make(map[string]interface{})
	query = make(map[string]interface{})

	request["EcrId"] = id
	request["ClientToken"] = buildClientToken(action)

	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(1*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2023-09-01"), StringPointer("AK"), query, request, &runtime)
		request["ClientToken"] = buildClientToken(action)

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
		addDebug(action, response, request)
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}

	return response, nil
}

// DescribeExpressConnectRouterExpressConnectRouter >>> Encapsulated.

// SetResourceTags <<< Encapsulated tag function for ExpressConnectRouter.
func (s *ExpressConnectRouterServiceV2) SetResourceTags(d *schema.ResourceData, resourceType string) error {
	if d.HasChange("tags") {
		var err error
		var action string
		var conn *rpc.Client
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
			action = "UntagResources"
			conn, err = client.NewExpressconnectrouterClient()
			if err != nil {
				return WrapError(err)
			}
			request = make(map[string]interface{})
			query = make(map[string]interface{})
			request["ResourceId.1"] = d.Id()
			request["ClientToken"] = buildClientToken(action)
			request["ResourceType"] = resourceType
			for i, key := range removedTagKeys {
				request[fmt.Sprintf("TagKey.%d", i+1)] = key
			}

			runtime := util.RuntimeOptions{}
			runtime.SetAutoretry(true)
			wait := incrementalWait(3*time.Second, 5*time.Second)
			err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
				response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2023-09-01"), StringPointer("AK"), query, request, &runtime)
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
			conn, err = client.NewExpressconnectrouterClient()
			if err != nil {
				return WrapError(err)
			}
			request = make(map[string]interface{})
			query = make(map[string]interface{})
			request["ResourceId.1"] = d.Id()
			request["ClientToken"] = buildClientToken(action)
			request["ResourceType"] = resourceType
			count := 1
			for key, value := range added {
				request[fmt.Sprintf("Tag.%d.Key", count)] = key
				request[fmt.Sprintf("Tag.%d.Value", count)] = value
				count++
			}

			runtime := util.RuntimeOptions{}
			runtime.SetAutoretry(true)
			wait := incrementalWait(3*time.Second, 5*time.Second)
			err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
				response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2023-09-01"), StringPointer("AK"), query, request, &runtime)
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
	}

	return nil
}

// SetResourceTags >>> tag function encapsulated.

// DescribeExpressConnectRouterExpressConnectRouterVbrChildInstance <<< Encapsulated get interface for ExpressConnectRouter ExpressConnectRouterVbrChildInstance.

func (s *ExpressConnectRouterServiceV2) DescribeExpressConnectRouterExpressConnectRouterVbrChildInstance(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	parts := strings.Split(id, ":")
	if len(parts) != 3 {
		err = WrapError(fmt.Errorf("invalid Resource Id %s. Expected parts' length %d, got %d", id, 3, len(parts)))
	}
	action := "DescribeExpressConnectRouterChildInstance"
	conn, err := client.NewExpressconnectrouterClient()
	if err != nil {
		return object, WrapError(err)
	}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["ChildInstanceId"] = parts[1]
	request["ChildInstanceType"] = parts[2]
	request["EcrId"] = parts[0]

	request["ClientToken"] = buildClientToken(action)

	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(1*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2023-09-01"), StringPointer("AK"), query, request, &runtime)
		request["ClientToken"] = buildClientToken(action)

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
		if IsExpectedErrors(err, []string{"ResourceNotFound.EcrId"}) {
			return object, WrapErrorf(Error(GetNotFoundMessage("ExpressConnectRouterVbrChildInstance", id)), NotFoundMsg, response)
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$.ChildInstanceList[*]", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.ChildInstanceList[*]", response)
	}

	if len(v.([]interface{})) == 0 {
		return object, WrapErrorf(Error(GetNotFoundMessage("ExpressConnectRouterVbrChildInstance", id)), NotFoundMsg, response)
	}

	result, _ := v.([]interface{})
	for _, v := range result {
		item := v.(map[string]interface{})
		if fmt.Sprint(item["ChildInstanceId"]) != parts[1] {
			continue
		}
		if fmt.Sprint(item["ChildInstanceType"]) != parts[2] {
			continue
		}
		if fmt.Sprint(item["ChildInstanceType"]) != "VBR" {
			continue
		}
		if fmt.Sprint(item["EcrId"]) != parts[0] {
			continue
		}
		return item, nil
	}
	return object, WrapErrorf(Error(GetNotFoundMessage("ExpressConnectRouterVbrChildInstance", id)), NotFoundMsg, response)
}

func (s *ExpressConnectRouterServiceV2) ExpressConnectRouterExpressConnectRouterVbrChildInstanceStateRefreshFunc(id string, field string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeExpressConnectRouterExpressConnectRouterVbrChildInstance(id)
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

func (s *ExpressConnectRouterServiceV2) DescribeAsyncExpressConnectRouterVbrChildInstanceDescribeExpressConnectRouterChildInstance(d *schema.ResourceData, res map[string]interface{}) (object map[string]interface{}, err error) {
	client := s.client
	id := d.Id()
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	parts := strings.Split(id, ":")
	if len(parts) != 3 {
		err = WrapError(fmt.Errorf("invalid Resource Id %s. Expected parts' length %d, got %d", id, 3, len(parts)))
	}
	action := "DescribeExpressConnectRouterChildInstance"
	conn, err := client.NewExpressconnectrouterClient()
	if err != nil {
		return object, WrapError(err)
	}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["EcrId"] = parts[0]
	request["ChildInstanceId"] = parts[1]

	request["ClientToken"] = buildClientToken(action)

	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(1*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2023-09-01"), StringPointer("AK"), query, request, &runtime)
		request["ClientToken"] = buildClientToken(action)

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
		return response, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}

	return response, nil
}

func (s *ExpressConnectRouterServiceV2) DescribeAsyncExpressConnectRouterExpressConnectRouterVbrChildInstanceStateRefreshFunc(d *schema.ResourceData, res map[string]interface{}, field string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeAsyncExpressConnectRouterVbrChildInstanceDescribeExpressConnectRouterChildInstance(d, res)
		if err != nil {
			if NotFoundError(err) {
				return object, "", nil
			}
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
				if _err, ok := object["error"]; ok {
					return _err, currentStatus, WrapError(Error(FailedToReachTargetStatus, currentStatus))
				}
				return object, currentStatus, WrapError(Error(FailedToReachTargetStatus, currentStatus))
			}
		}
		return object, currentStatus, nil
	}
}

// DescribeExpressConnectRouterExpressConnectRouterVbrChildInstance >>> Encapsulated.
// DescribeExpressConnectRouterExpressConnectRouterTrAssociation <<< Encapsulated get interface for ExpressConnectRouter ExpressConnectRouterTrAssociation.

func (s *ExpressConnectRouterServiceV2) DescribeExpressConnectRouterExpressConnectRouterTrAssociation(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	parts := strings.Split(id, ":")
	if len(parts) != 3 {
		err = WrapError(fmt.Errorf("invalid Resource Id %s. Expected parts' length %d, got %d", id, 3, len(parts)))
	}
	action := "DescribeExpressConnectRouterAssociation"
	conn, err := client.NewExpressconnectrouterClient()
	if err != nil {
		return object, WrapError(err)
	}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["AssociationId"] = parts[1]
	request["EcrId"] = parts[0]
	request["TransitRouterId"] = parts[2]

	request["ClientToken"] = buildClientToken(action)

	request["AssociationNodeType"] = "TR"
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(1*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2023-09-01"), StringPointer("AK"), query, request, &runtime)
		request["ClientToken"] = buildClientToken(action)

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
	code, _ := jsonpath.Get("$.Code", response)
	if InArray(fmt.Sprint(code), []string{"ResourceNotFound.EcrId"}) {
		return object, WrapErrorf(Error(GetNotFoundMessage("ExpressConnectRouterTrAssociation", id)), NotFoundMsg, response)
	}

	v, err := jsonpath.Get("$.AssociationList[*]", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.AssociationList[*]", response)
	}

	if len(v.([]interface{})) == 0 {
		return object, WrapErrorf(Error(GetNotFoundMessage("ExpressConnectRouterTrAssociation", id)), NotFoundMsg, response)
	}

	result, _ := v.([]interface{})
	for _, v := range result {
		item := v.(map[string]interface{})
		if item["AssociationId"] != parts[1] {
			continue
		}
		if item["AssociationNodeType"] != "TR" {
			continue
		}
		if item["EcrId"] != parts[0] {
			continue
		}
		if item["TransitRouterId"] != parts[2] {
			continue
		}
		return item, nil
	}
	return object, WrapErrorf(Error(GetNotFoundMessage("ExpressConnectRouterTrAssociation", id)), NotFoundMsg, response)
}

func (s *ExpressConnectRouterServiceV2) ExpressConnectRouterExpressConnectRouterTrAssociationStateRefreshFunc(id string, field string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeExpressConnectRouterExpressConnectRouterTrAssociation(id)
		if err != nil {
			if NotFoundError(err) {
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}

		v, err := jsonpath.Get(field, object)
		currentStatus := fmt.Sprint(v)

		for _, failState := range failStates {
			if currentStatus == failState {
				return object, currentStatus, WrapError(Error(FailedToReachTargetStatus, currentStatus))
			}
		}
		return object, currentStatus, nil
	}
}

// DescribeExpressConnectRouterExpressConnectRouterTrAssociation >>> Encapsulated.
// DescribeExpressConnectRouterExpressConnectRouterVpcAssociation <<< Encapsulated get interface for ExpressConnectRouter ExpressConnectRouterVpcAssociation.

func (s *ExpressConnectRouterServiceV2) DescribeExpressConnectRouterExpressConnectRouterVpcAssociation(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	parts := strings.Split(id, ":")
	if len(parts) != 3 {
		err = WrapError(fmt.Errorf("invalid Resource Id %s. Expected parts' length %d, got %d", id, 3, len(parts)))
	}
	action := "DescribeExpressConnectRouterAssociation"
	conn, err := client.NewExpressconnectrouterClient()
	if err != nil {
		return object, WrapError(err)
	}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["AssociationId"] = parts[1]
	request["EcrId"] = parts[0]
	request["VpcId"] = parts[2]

	request["ClientToken"] = buildClientToken(action)

	request["AssociationNodeType"] = "VPC"
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(1*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2023-09-01"), StringPointer("AK"), query, request, &runtime)
		request["ClientToken"] = buildClientToken(action)

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
	code, _ := jsonpath.Get("$.Code", response)
	if InArray(fmt.Sprint(code), []string{"ResourceNotFound.EcrId"}) {
		return object, WrapErrorf(Error(GetNotFoundMessage("ExpressConnectRouterVpcAssociation", id)), NotFoundMsg, response)
	}

	v, err := jsonpath.Get("$.AssociationList[*]", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.AssociationList[*]", response)
	}

	if len(v.([]interface{})) == 0 {
		return object, WrapErrorf(Error(GetNotFoundMessage("ExpressConnectRouterVpcAssociation", id)), NotFoundMsg, response)
	}

	result, _ := v.([]interface{})
	for _, v := range result {
		item := v.(map[string]interface{})
		if item["AssociationId"] != parts[1] {
			continue
		}
		if item["AssociationNodeType"] != "VPC" {
			continue
		}
		if item["EcrId"] != parts[0] {
			continue
		}
		if item["VpcId"] != parts[2] {
			continue
		}
		return item, nil
	}
	return object, WrapErrorf(Error(GetNotFoundMessage("ExpressConnectRouterVpcAssociation", id)), NotFoundMsg, response)
}
func (s *ExpressConnectRouterServiceV2) DescribeDescribeExpressConnectRouterAllowedPrefixHistory(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	parts := strings.Split(id, ":")
	if len(parts) != 3 {
		err = WrapError(fmt.Errorf("invalid Resource Id %s. Expected parts' length %d, got %d", id, 3, len(parts)))
	}
	action := "DescribeExpressConnectRouterAllowedPrefixHistory"
	conn, err := client.NewExpressconnectrouterClient()
	if err != nil {
		return object, WrapError(err)
	}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["EcrId"] = parts[0]
	request["InstanceId"] = parts[2]

	request["ClientToken"] = buildClientToken(action)

	request["InstanceType"] = "VPC"
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(1*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2023-09-01"), StringPointer("AK"), query, request, &runtime)
		request["ClientToken"] = buildClientToken(action)

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
		if IsExpectedErrors(err, []string{"ResourceNotFound.EcrId", "ResourceNotFound.AssociationId", "ResourceNotFound.InstanceId"}) {
			return object, WrapErrorf(Error(GetNotFoundMessage("ExpressConnectRouterVpcAssociation", id)), NotFoundMsg, response)
		}
		addDebug(action, response, request)
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$.AllowedPrefixHistoryList", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.AllowedPrefixHistoryList", response)
	}

	return v.(map[string]interface{}), nil
}

func (s *ExpressConnectRouterServiceV2) ExpressConnectRouterExpressConnectRouterVpcAssociationStateRefreshFunc(id string, field string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeExpressConnectRouterExpressConnectRouterVpcAssociation(id)
		if err != nil {
			if NotFoundError(err) {
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}

		v, err := jsonpath.Get(field, object)
		currentStatus := fmt.Sprint(v)

		for _, failState := range failStates {
			if currentStatus == failState {
				return object, currentStatus, WrapError(Error(FailedToReachTargetStatus, currentStatus))
			}
		}
		return object, currentStatus, nil
	}
}

// DescribeExpressConnectRouterExpressConnectRouterVpcAssociation >>> Encapsulated.

// DescribeExpressConnectRouterGrantAssociation <<< Encapsulated get interface for ExpressConnectRouter GrantAssociation.

func (s *ExpressConnectRouterServiceV2) DescribeExpressConnectRouterGrantAssociation(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	parts := strings.Split(id, ":")
	if len(parts) != 3 {
		err = WrapError(fmt.Errorf("invalid Resource Id %s. Expected parts' length %d, got %d", id, 3, len(parts)))
	}
	action := "DescribeInstanceGrantedToExpressConnectRouter"
	conn, err := client.NewExpressconnectrouterClient()
	if err != nil {
		return object, WrapError(err)
	}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["EcrId"] = parts[0]
	request["InstanceId"] = parts[1]
	request["InstanceRegionId"] = parts[2]

	request["ClientToken"] = buildClientToken(action)

	request["CallerType"] = "OTHER"
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(1*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2023-09-01"), StringPointer("AK"), query, request, &runtime)
		request["ClientToken"] = buildClientToken(action)

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

	v, err := jsonpath.Get("$.EcrGrantedInstanceList[*]", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.EcrGrantedInstanceList[*]", response)
	}

	if len(v.([]interface{})) == 0 {
		return object, WrapErrorf(Error(GetNotFoundMessage("GrantAssociation", id)), NotFoundMsg, response)
	}

	return v.([]interface{})[0].(map[string]interface{}), nil
}

func (s *ExpressConnectRouterServiceV2) ExpressConnectRouterGrantAssociationStateRefreshFunc(id string, field string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeExpressConnectRouterGrantAssociation(id)
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

// DescribeExpressConnectRouterGrantAssociation >>> Encapsulated.

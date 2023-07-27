package alicloud

import (
	"fmt"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

type EventBridgeServiceV2 struct {
	client *connectivity.AliyunClient
}

func (s *EventBridgeServiceV2) DescribeEventBridgeConnection(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	action := "GetConnection"

	conn, err := s.client.NewEventbridgeClient()
	if err != nil {
		return object, WrapError(err)
	}

	request := map[string]interface{}{
		"ConnectionName": id,
	}

	idExist := false
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-04-01"), StringPointer("AK"), nil, request, &runtime)
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
		if IsExpectedErrors(err, []string{"ConnectionNotExist"}) {
			return object, WrapErrorf(Error(GetNotFoundMessage("EventBridge:Connection", id)), NotFoundWithResponse, response)
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}

	resp, err := jsonpath.Get("$.Data.Connections", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Data.Connections", response)
	}

	if v, ok := resp.([]interface{}); !ok || len(v) < 1 {
		return object, WrapErrorf(Error(GetNotFoundMessage("EventBridge:Connection", id)), NotFoundWithResponse, response)
	}

	for _, v := range resp.([]interface{}) {
		if fmt.Sprint(v.(map[string]interface{})["ConnectionName"]) == id {
			idExist = true
			return v.(map[string]interface{}), nil
		}
	}

	if !idExist {
		return object, WrapErrorf(Error(GetNotFoundMessage("EventBridge:Connection", id)), NotFoundWithResponse, response)
	}

	return object, nil
}

func (s *EventBridgeServiceV2) DescribeEventBridgeApiDestination(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	action := "GetApiDestination"

	conn, err := s.client.NewEventbridgeClient()
	if err != nil {
		return object, WrapError(err)
	}

	request := map[string]interface{}{
		"ApiDestinationName": id,
	}

	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(1*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-04-01"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
		if IsExpectedErrors(err, []string{"ApiDestinationNotExist"}) {
			return object, WrapErrorf(Error(GetNotFoundMessage("EventBridge:ApiDestination", id)), NotFoundWithResponse, response)
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$.Data", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Data", response)
	}

	object = v.(map[string]interface{})

	return object, nil
}

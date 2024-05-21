package alicloud

import (
	"fmt"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

type MessageServiceServiceV2 struct {
	client *connectivity.AliyunClient
}

// DescribeMessageServiceQueue <<< Encapsulated get interface for MessageService Queue.

func (s *MessageServiceServiceV2) DescribeMessageServiceQueue(id string) (object map[string]interface{}, err error) {
	client := s.client
	var response map[string]interface{}
	var query map[string]interface{}
	action := "GetQueueAttributes"
	conn, err := client.NewMnsClient()
	if err != nil {
		return object, WrapError(err)
	}
	query = make(map[string]interface{})
	query["QueueName"] = id

	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(1*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2022-01-19"), StringPointer("AK"), query, nil, &runtime)

		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, query)
		return nil
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"QueueNotExist"}) {
			return object, WrapErrorf(Error(GetNotFoundMessage("Queue", id)), NotFoundMsg, response)
		}
		addDebug(action, response, query)
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$.Data", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Data", response)
	}

	return v.(map[string]interface{}), nil
}

func (s *MessageServiceServiceV2) MessageServiceQueueStateRefreshFunc(id string, field string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeMessageServiceQueue(id)
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

// DescribeMessageServiceQueue >>> Encapsulated.

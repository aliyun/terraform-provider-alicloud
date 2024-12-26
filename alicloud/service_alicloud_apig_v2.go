package alicloud

import (
	"fmt"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
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
	action := fmt.Sprintf("/v1/http-apis/%s", httpApiId)
	conn, err := client.NewApigClient()
	if err != nil {
		return object, WrapError(err)
	}
	request = make(map[string]interface{})
	query = make(map[string]*string)
	request["httpApiId"] = id

	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(1*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer("2024-03-27"), nil, StringPointer("GET"), StringPointer("AK"), StringPointer(action), query, nil, nil, &runtime)

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
		return object, WrapErrorf(Error(GetNotFoundMessage("HttpApi", id)), NotFoundMsg, response)
	}
	response = response["body"].(map[string]interface{})
	code, _ := jsonpath.Get("$.code", response)
	if InArray(fmt.Sprint(code), []string{"DatabaseError.RecordNotFound"}) {
		return object, WrapErrorf(Error(GetNotFoundMessage("HttpApi", id)), NotFoundMsg, response)
	}

	v, err := jsonpath.Get("$.data", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.data", response)
	}

	return v.(map[string]interface{}), nil
}

func (s *ApigServiceV2) ApigHttpApiStateRefreshFunc(id string, field string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeApigHttpApi(id)
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

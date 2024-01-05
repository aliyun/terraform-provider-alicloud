package alicloud

import (
	"fmt"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

type ExpressConnectServiceV2 struct {
	client *connectivity.AliyunClient
}

// DescribeExpressConnectEcFailoverTestJob <<< Encapsulated get interface for ExpressConnect EcFailoverTestJob.

func (s *ExpressConnectServiceV2) DescribeExpressConnectEcFailoverTestJob(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	action := "DescribeFailoverTestJob"
	conn, err := client.NewExpressconnectClient()
	if err != nil {
		return object, WrapError(err)
	}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	query["JobId"] = id
	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken(action)

	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(1*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), query, request, &runtime)
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
		if IsExpectedErrors(err, []string{"IllegalParam.JobId"}) {
			return object, WrapErrorf(Error(GetNotFoundMessage("EcFailoverTestJob", id)), NotFoundMsg, response)
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$.FailoverTestJobModel", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.FailoverTestJobModel", response)
	}

	return v.(map[string]interface{}), nil
}

func (s *ExpressConnectServiceV2) ExpressConnectEcFailoverTestJobStateRefreshFunc(id string, field string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeExpressConnectEcFailoverTestJob(id)
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

// DescribeExpressConnectEcFailoverTestJob >>> Encapsulated.

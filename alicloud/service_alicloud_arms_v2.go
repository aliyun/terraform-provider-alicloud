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

type ArmsServiceV2 struct {
	client *connectivity.AliyunClient
}

// DescribeArmsPrometheusMonitoring <<< Encapsulated get interface for Arms PrometheusMonitoring.

func (s *ArmsServiceV2) DescribeArmsPrometheusMonitoring(id string) (object map[string]interface{}, err error) {

	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	parts := strings.Split(id, ":")
	if len(parts) != 3 {
		err = WrapError(fmt.Errorf("invalid Resource Id %s. Expected parts' length %d, got %d", id, 3, len(parts)))
	}
	action := "GetPrometheusMonitoring"
	conn, err := client.NewArmsClient()
	if err != nil {
		return object, WrapError(err)
	}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	query["ClusterId"] = parts[0]
	query["MonitoringName"] = parts[1]
	query["Type"] = parts[2]
	request["RegionId"] = client.RegionId

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(1*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-08-08"), StringPointer("AK"), query, request, &util.RuntimeOptions{})

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
	if InArray(fmt.Sprint(code), []string{"404"}) {
		return object, WrapErrorf(Error(GetNotFoundMessage("PrometheusMonitoring", id)), NotFoundMsg, ProviderERROR, fmt.Sprint(response["RequestId"]))
	}

	v, err := jsonpath.Get("$.Data", response)
	if err != nil {
		return object, WrapErrorf(Error(GetNotFoundMessage("PrometheusMonitoring", id)), NotFoundMsg, ProviderERROR, fmt.Sprint(response["RequestId"]))
	}

	return v.(map[string]interface{}), nil
}

func (s *ArmsServiceV2) ArmsPrometheusMonitoringStateRefreshFunc(id string, field string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeArmsPrometheusMonitoring(id)
		if err != nil {
			if NotFoundError(err) {
				return object, "", nil
			}
			return nil, "", WrapError(err)
		}

		currentStatus := fmt.Sprint(object[field])
		for _, failState := range failStates {
			if currentStatus == failState {
				return object, currentStatus, WrapError(Error(FailedToReachTargetStatus, currentStatus))
			}
		}
		return object, currentStatus, nil
	}
}

// DescribeArmsPrometheusMonitoring >>> Encapsulated.
// DescribeArmsRemoteWrite <<< Encapsulated get interface for Arms RemoteWrite.

func (s *ArmsServiceV2) DescribeArmsRemoteWrite(id string) (object map[string]interface{}, err error) {

	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	parts := strings.Split(id, ":")
	if len(parts) != 2 {
		err = WrapError(fmt.Errorf("invalid Resource Id %s. Expected parts' length %d, got %d", id, 2, len(parts)))
	}
	action := "GetPrometheusRemoteWrite"
	conn, err := client.NewArmsClient()
	if err != nil {
		return object, WrapError(err)
	}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	query["ClusterId"] = parts[0]
	query["RemoteWriteName"] = parts[1]
	request["RegionId"] = client.RegionId

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(1*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-08-08"), StringPointer("AK"), query, request, &util.RuntimeOptions{})

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
	if InArray(fmt.Sprint(code), []string{"404"}) {
		return object, WrapErrorf(Error(GetNotFoundMessage("RemoteWrite", id)), NotFoundMsg, ProviderERROR, fmt.Sprint(response["RequestId"]))
	}

	v, err := jsonpath.Get("$.Data", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Data", response)
	}

	return v.(map[string]interface{}), nil
}

func (s *ArmsServiceV2) ArmsRemoteWriteStateRefreshFunc(id string, field string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeArmsRemoteWrite(id)
		if err != nil {
			if NotFoundError(err) {
				return object, "", nil
			}
			return nil, "", WrapError(err)
		}

		currentStatus := fmt.Sprint(object[field])
		for _, failState := range failStates {
			if currentStatus == failState {
				return object, currentStatus, WrapError(Error(FailedToReachTargetStatus, currentStatus))
			}
		}
		return object, currentStatus, nil
	}
}

// DescribeArmsRemoteWrite >>> Encapsulated.

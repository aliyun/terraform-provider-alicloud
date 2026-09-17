package alicloud

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// DescribeRealtimeComputeVariable <<< Encapsulated get interface for RealtimeCompute Variable.

func (s *RealtimeComputeServiceV2) DescribeRealtimeComputeVariable(id string) (object map[string]interface{}, err error) {
	client := s.client
	var query map[string]*string
	var header map[string]*string
	parts := strings.Split(id, ":")
	if len(parts) != 3 {
		err = WrapError(fmt.Errorf("invalid Resource Id %s. Expected parts' length %d, got %d", id, 3, len(parts)))
		return nil, err
	}
	workspace := parts[0]
	namespace := parts[1]
	name := parts[2]
	query = make(map[string]*string)
	header = make(map[string]*string)
	header["workspace"] = StringPointer(workspace)

	action := fmt.Sprintf("/api/v2/namespaces/%s/variables", namespace)

	pageSize := PageSizeLarge
	pageIndex := 1

	for {
		query["pageSize"] = StringPointer(strconv.Itoa(pageSize))
		query["pageIndex"] = StringPointer(strconv.Itoa(pageIndex))

		var response map[string]interface{}
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(1*time.Minute, func() *resource.RetryError {
			response, err = client.RoaGet("ververica", "2022-07-18", action, query, header, nil)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, query)
		if err != nil {
			if IsExpectedErrors(err, []string{"990301"}) {
				return object, WrapErrorf(NotFoundErr("variable", id), NotFoundMsg, response)
			}
			return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
		}

		v, err := jsonpath.Get("$.data", response)
		if err != nil {
			return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.data", response)
		}

		items, _ := v.([]interface{})
		for _, item := range items {
			m, _ := item.(map[string]interface{})
			if fmt.Sprint(m["name"]) == name {
				return m, nil
			}
		}

		if len(items) < pageSize {
			break
		}
		pageIndex++
	}

	return object, WrapErrorf(NotFoundErr("variable", id), NotFoundMsg, nil)
}

func (s *RealtimeComputeServiceV2) RealtimeComputeVariableStateRefreshFunc(id string, field string, failStates []string) resource.StateRefreshFunc {
	return s.RealtimeComputeVariableStateRefreshFuncWithApi(id, field, failStates, s.DescribeRealtimeComputeVariable)
}

func (s *RealtimeComputeServiceV2) RealtimeComputeVariableStateRefreshFuncWithApi(id string, field string, failStates []string, call func(id string) (map[string]interface{}, error)) resource.StateRefreshFunc {
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

// DescribeRealtimeComputeVariable >>> Encapsulated.

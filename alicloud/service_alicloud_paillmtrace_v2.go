package alicloud

import (
	"fmt"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

type PaillmtraceServiceV2 struct {
	client *connectivity.AliyunClient
}

// DescribePaillmtraceEval <<< Encapsulated get interface for Paillmtrace Eval.

func (s *PaillmtraceServiceV2) DescribePaillmtraceEval(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	request = make(map[string]interface{})
	query = make(map[string]*string)
	query["EvaluationId"] = StringPointer(id)

	action := fmt.Sprintf("/api/v1/PAILLMTrace/eval")

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(1*time.Minute, func() *resource.RetryError {
		response, err = client.RoaGet("PaiLLMTrace", "2024-03-11", action, query, nil, nil)
		if err != nil {
			if IsExpectedErrors(err, []string{"NotFound"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		if IsExpectedErrors(err, []string{"NotFound", "EntityNotExist"}) {
			return object, WrapErrorf(NotFoundErr("Eval", id), NotFoundMsg, response)
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	code, _ := jsonpath.Get("$.Code", response)
	if InArray(fmt.Sprint(code), []string{"NotFound", "EntityNotExist"}) {
		return object, WrapErrorf(NotFoundErr("Eval", id), NotFoundMsg, response)
	}

	return response, nil
}

func (s *PaillmtraceServiceV2) PaillmtraceEvalStateRefreshFunc(id string, field string, failStates []string) resource.StateRefreshFunc {
	return s.PaillmtraceEvalStateRefreshFuncWithApi(id, field, failStates, s.DescribePaillmtraceEval)
}

func (s *PaillmtraceServiceV2) PaillmtraceEvalStateRefreshFuncWithApi(id string, field string, failStates []string, call func(id string) (map[string]interface{}, error)) resource.StateRefreshFunc {
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

// DescribePaillmtraceEval >>> Encapsulated.

// ListPaillmtraceEvals <<< Encapsulated list interface for Paillmtrace Eval.

func (s *PaillmtraceServiceV2) ListPaillmtraceEvals(pageNumber, pageSize int) (objects []interface{}, totalCount int, err error) {
	client := s.client
	var response map[string]interface{}
	var query map[string]*string
	query = make(map[string]*string)
	query["PageNumber"] = StringPointer(fmt.Sprint(pageNumber))
	query["PageSize"] = StringPointer(fmt.Sprint(pageSize))

	action := fmt.Sprintf("/api/v1/PAILLMTrace/eval")

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(1*time.Minute, func() *resource.RetryError {
		response, err = client.RoaGet("PaiLLMTrace", "2024-03-11", action, query, nil, nil)
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
		return nil, 0, WrapErrorf(err, DefaultErrorMsg, "ListEvals", action, AlibabaCloudSdkGoERROR)
	}

	v, _ := jsonpath.Get("$.TotalCount", response)
	if v != nil {
		totalCount = int(v.(float64))
	}

	items, _ := jsonpath.Get("$.Evaluations[*]", response)
	if items != nil {
		if arr, ok := items.([]interface{}); ok {
			objects = arr
		}
	}

	return objects, totalCount, nil
}

// ListPaillmtraceEvals >>> Encapsulated.

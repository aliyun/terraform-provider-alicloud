// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"github.com/PaesslerAG/jsonpath"
	"strings"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

type AgentLoopServiceV2 struct {
	client *connectivity.AliyunClient
}

// DescribeAgentLoopAgentSpace <<< Encapsulated get interface for AgentLoop AgentSpace.

func (s *AgentLoopServiceV2) DescribeAgentLoopAgentSpace(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	agentSpace := id
	request = make(map[string]interface{})
	query = make(map[string]*string)

	action := fmt.Sprintf("/agentspace/%s", agentSpace)

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(1*time.Minute, func() *resource.RetryError {
		response, err = client.RoaGet("AgentLoop", "2026-05-20", action, query, nil, nil)

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
		if IsExpectedErrors(err, []string{"AgentSpaceNotExist"}) {
			return object, WrapErrorf(NotFoundErr("AgentSpace", id), NotFoundMsg, response)
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}

	return response, nil
}

func (s *AgentLoopServiceV2) AgentLoopAgentSpaceStateRefreshFunc(id string, field string, failStates []string) resource.StateRefreshFunc {
	return s.AgentLoopAgentSpaceStateRefreshFuncWithApi(id, field, failStates, s.DescribeAgentLoopAgentSpace)
}

func (s *AgentLoopServiceV2) AgentLoopAgentSpaceStateRefreshFuncWithApi(id string, field string, failStates []string, call func(id string) (map[string]interface{}, error)) resource.StateRefreshFunc {
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

// DescribeAgentLoopAgentSpace >>> Encapsulated.

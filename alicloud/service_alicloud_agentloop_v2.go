// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

type AgentloopServiceV2 struct {
	client *connectivity.AliyunClient
}

// DescribeAgentloopAgentSpace <<< Encapsulated get interface for Agentloop AgentSpace.

func (s *AgentloopServiceV2) DescribeAgentloopAgentSpace(id string) (object map[string]interface{}, err error) {
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
			return object, WrapErrorf(NotFoundErr("AgentSpace", id), NotFoundMsg, err)
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	if response == nil {
		return object, WrapErrorf(NotFoundErr("AgentSpace", id), NotFoundMsg, response)
	}

	return response, nil
}

func (s *AgentloopServiceV2) AgentloopAgentSpaceStateRefreshFunc(id string, field string, failStates []string) resource.StateRefreshFunc {
	return s.AgentloopAgentSpaceStateRefreshFuncWithApi(id, field, failStates, s.DescribeAgentloopAgentSpace)
}

func (s *AgentloopServiceV2) AgentloopEvaluationTaskStateRefreshFunc(id string, field string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeAgentloopEvaluationTask(id)
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

func (s *AgentloopServiceV2) AgentloopAgentSpaceStateRefreshFuncWithApi(id string, field string, failStates []string, call func(id string) (map[string]interface{}, error)) resource.StateRefreshFunc {
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

// DescribeAgentloopAgentSpace >>> Encapsulated.

// DescribeAgentloopContextStore <<< Encapsulated get interface for Agentloop ContextStore.

func (s *AgentloopServiceV2) DescribeAgentloopContextStore(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	parts := strings.Split(id, ":")
	if len(parts) != 2 {
		return object, WrapError(fmt.Errorf("invalid Resource Id %s. Expected parts' length %d, got %d", id, 2, len(parts)))
	}
	agentSpace := parts[0]
	contextStoreName := parts[1]
	request = make(map[string]interface{})
	query = make(map[string]*string)

	action := fmt.Sprintf("/agentspace/%s/contextstore/%s", agentSpace, contextStoreName)

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
		if IsExpectedErrors(err, []string{"ContextStoreNotExist"}) {
			return object, WrapErrorf(NotFoundErr("ContextStore", id), NotFoundMsg, err)
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	if response == nil {
		return object, WrapErrorf(NotFoundErr("ContextStore", id), NotFoundMsg, response)
	}

	return response, nil
}

// DescribeAgentloopContextStore >>> Encapsulated.

// DescribeAgentloopContextStoreApiKey <<< Encapsulated get interface for Agentloop ContextStoreApiKey.

func (s *AgentloopServiceV2) DescribeAgentloopContextStoreApiKey(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	parts := strings.Split(id, ":")
	if len(parts) != 3 {
		return object, WrapError(fmt.Errorf("invalid Resource Id %s. Expected parts' length %d, got %d", id, 3, len(parts)))
	}
	agentSpace := parts[0]
	contextStoreName := parts[1]
	name := parts[2]
	request = make(map[string]interface{})
	query = make(map[string]*string)

	action := fmt.Sprintf("/agentspace/%s/contextstore/%s/apikey/%s", agentSpace, contextStoreName, name)

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
		// ContextStoreNotExist / AgentSpaceNotExist mean a parent resource was
		// deleted out-of-band; the API key is gone with it, so treat them as
		// NotFound and let the state be cleared on refresh.
		if IsExpectedErrors(err, []string{"NotFound", "ContextStoreNotExist", "AgentSpaceNotExist"}) {
			return object, WrapErrorf(NotFoundErr("ContextStoreApiKey", id), NotFoundMsg, err)
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	if response == nil {
		return object, WrapErrorf(NotFoundErr("ContextStoreApiKey", id), NotFoundMsg, response)
	}

	return response, nil
}

// DescribeAgentloopContextStoreApiKey >>> Encapsulated.

// DescribeAgentloopDataset <<< Encapsulated get interface for Agentloop Dataset.

func (s *AgentloopServiceV2) DescribeAgentloopDataset(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	parts := strings.Split(id, ":")
	if len(parts) != 2 {
		return object, WrapError(fmt.Errorf("invalid Resource Id %s. Expected parts' length %d, got %d", id, 2, len(parts)))
	}
	agentSpace := parts[0]
	datasetName := parts[1]
	request = make(map[string]interface{})
	query = make(map[string]*string)

	action := fmt.Sprintf("/agentspace/%s/dataset/%s", agentSpace, datasetName)

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
		// AgentSpaceNotExist means the hosting AgentSpace was deleted
		// out-of-band; the dataset is gone with it, so treat it as NotFound
		// and let the state be cleared on refresh.
		if IsExpectedErrors(err, []string{"DatasetNotExist", "WorkspaceNotExist", "AgentSpaceNotExist"}) {
			return object, WrapErrorf(NotFoundErr("Dataset", id), NotFoundMsg, err)
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	if response == nil {
		return object, WrapErrorf(NotFoundErr("Dataset", id), NotFoundMsg, response)
	}

	return response, nil
}

// DescribeAgentloopDataset >>> Encapsulated.

// DescribeAgentloopEndpointConnector <<< Encapsulated get interface for Agentloop EndpointConnector.

func (s *AgentloopServiceV2) DescribeAgentloopEndpointConnector(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	parts := strings.Split(id, ":")
	if len(parts) != 2 {
		return object, WrapError(fmt.Errorf("invalid Resource Id %s. Expected parts' length %d, got %d", id, 2, len(parts)))
	}
	agentSpace := parts[0]
	connectorId := parts[1]
	request = make(map[string]interface{})
	query = make(map[string]*string)

	action := fmt.Sprintf("/api/v1/endpoint-connectors/%s/%s", agentSpace, connectorId)

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
		if IsExpectedErrors(err, []string{"NotFound.EndpointConnector"}) {
			return object, WrapErrorf(NotFoundErr("EndpointConnector", id), NotFoundMsg, err)
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	if response == nil {
		return object, WrapErrorf(NotFoundErr("EndpointConnector", id), NotFoundMsg, response)
	}

	return response, nil
}

// DescribeAgentloopEndpointConnector >>> Encapsulated.

// DescribeAgentloopEvaluationTask <<< Encapsulated get interface for Agentloop EvaluationTask.

func (s *AgentloopServiceV2) DescribeAgentloopEvaluationTask(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	parts := strings.Split(id, ":")
	if len(parts) != 2 {
		return object, WrapError(fmt.Errorf("invalid Resource Id %s. Expected parts' length %d, got %d", id, 2, len(parts)))
	}
	agentSpace := parts[0]
	taskId := parts[1]
	request = make(map[string]interface{})
	query = make(map[string]*string)

	action := fmt.Sprintf("/api/v1/evaluation-task/%s/%s", agentSpace, taskId)

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
		if IsExpectedErrors(err, []string{"EvaluationTaskNotFound"}) {
			return object, WrapErrorf(NotFoundErr("EvaluationTask", id), NotFoundMsg, err)
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	if response == nil {
		return object, WrapErrorf(NotFoundErr("EvaluationTask", id), NotFoundMsg, response)
	}

	return response, nil
}

// DescribeAgentloopEvaluationTask >>> Encapsulated.

// DescribeAgentloopEvaluator <<< Encapsulated get interface for Agentloop Evaluator.

func (s *AgentloopServiceV2) DescribeAgentloopEvaluator(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	parts := strings.Split(id, ":")
	if len(parts) != 2 {
		return object, WrapError(fmt.Errorf("invalid Resource Id %s. Expected parts' length %d, got %d", id, 2, len(parts)))
	}
	agentSpace := parts[0]
	name := parts[1]
	request = make(map[string]interface{})
	query = make(map[string]*string)

	action := fmt.Sprintf("/api/v1/evaluators/%s/%s", agentSpace, name)

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
		if IsExpectedErrors(err, []string{"EvaluatorNotFound"}) {
			return object, WrapErrorf(NotFoundErr("Evaluator", id), NotFoundMsg, err)
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	if response == nil {
		return object, WrapErrorf(NotFoundErr("Evaluator", id), NotFoundMsg, response)
	}

	return response, nil
}

// DescribeAgentloopEvaluator >>> Encapsulated.

// DescribeAgentloopEvaluatorSkill <<< Encapsulated get interface for Agentloop EvaluatorSkill.

func (s *AgentloopServiceV2) DescribeAgentloopEvaluatorSkill(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	parts := strings.Split(id, ":")
	if len(parts) != 3 {
		return object, WrapError(fmt.Errorf("invalid Resource Id %s. Expected parts' length %d, got %d", id, 3, len(parts)))
	}
	agentSpace := parts[0]
	name := parts[1]
	skillName := parts[2]
	request = make(map[string]interface{})
	query = make(map[string]*string)
	query["agentSpace"] = StringPointer(agentSpace)

	action := fmt.Sprintf("/api/v1/evaluator/%s/skill/%s", name, skillName)

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
		if IsExpectedErrors(err, []string{"EvaluatorSkillNotFound"}) {
			return object, WrapErrorf(NotFoundErr("EvaluatorSkill", id), NotFoundMsg, err)
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	if response == nil {
		return object, WrapErrorf(NotFoundErr("EvaluatorSkill", id), NotFoundMsg, response)
	}

	return response, nil
}

// DescribeAgentloopEvaluatorSkill >>> Encapsulated.

// DescribeAgentloopOptimizeTask <<< Encapsulated get interface for Agentloop OptimizeTask.

func (s *AgentloopServiceV2) DescribeAgentloopOptimizeTask(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	parts := strings.Split(id, ":")
	if len(parts) != 2 {
		return object, WrapError(fmt.Errorf("invalid Resource Id %s. Expected parts' length %d, got %d", id, 2, len(parts)))
	}
	agentSpace := parts[0]
	optimizeTaskName := parts[1]
	request = make(map[string]interface{})
	query = make(map[string]*string)

	action := fmt.Sprintf("/agentspace/%s/optimizetasks/%s", agentSpace, optimizeTaskName)

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
		if IsExpectedErrors(err, []string{"OptimizeTaskNotExist"}) {
			return object, WrapErrorf(NotFoundErr("OptimizeTask", id), NotFoundMsg, err)
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	if response == nil {
		return object, WrapErrorf(NotFoundErr("OptimizeTask", id), NotFoundMsg, response)
	}

	return response, nil
}

// DescribeAgentloopOptimizeTask >>> Encapsulated.

// DescribeAgentloopPipeline <<< Encapsulated get interface for Agentloop Pipeline.

func (s *AgentloopServiceV2) DescribeAgentloopPipeline(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	parts := strings.Split(id, ":")
	if len(parts) != 2 {
		return object, WrapError(fmt.Errorf("invalid Resource Id %s. Expected parts' length %d, got %d", id, 2, len(parts)))
	}
	agentSpace := parts[0]
	pipelineName := parts[1]
	request = make(map[string]interface{})
	query = make(map[string]*string)

	action := fmt.Sprintf("/agentspace/%s/pipeline/%s", agentSpace, pipelineName)

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
		if IsExpectedErrors(err, []string{"PipelineNotExist"}) {
			return object, WrapErrorf(NotFoundErr("Pipeline", id), NotFoundMsg, err)
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	if response == nil {
		return object, WrapErrorf(NotFoundErr("Pipeline", id), NotFoundMsg, response)
	}

	return response, nil
}

// DescribeAgentloopPipeline >>> Encapsulated.

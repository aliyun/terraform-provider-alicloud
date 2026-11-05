package agentloop

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

type Service struct {
	client *connectivity.AliyunClient
}

func NewService(client *connectivity.AliyunClient) *Service {
	return &Service{client: client}
}

// DescribeAgentSpace <<< Encapsulated get interface for Agentloop AgentSpace.

func (s *Service) DescribeAgentSpace(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	agentSpace := id
	request = make(map[string]interface{})
	query = make(map[string]*string)

	action := fmt.Sprintf("/agentspace/%s", agentSpace)

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.RetryContext(context.Background(), 1*time.Minute, func() *retry.RetryError {
		response, err = client.RoaGet("AgentLoop", "2026-05-20", action, query, nil, nil)

		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"AgentSpaceNotExist"}) {
			return object, errs.WrapErrorf(errs.NotFoundErr("AgentSpace", id), errs.NotFoundMsg, err)
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}
	if response == nil {
		return object, errs.WrapErrorf(errs.NotFoundErr("AgentSpace", id), errs.NotFoundMsg, response)
	}

	return response, nil
}

func (s *Service) EvaluationTaskStateRefreshFunc(id string, field string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeEvaluationTask(id)
		if err != nil {
			if errs.NotFoundError(err) {
				return object, "", nil
			}
			return nil, "", errs.WrapError(err)
		}
		v, _ := jsonpath.Get(field, object)
		currentStatus := fmt.Sprint(v)

		if strings.HasPrefix(field, "#") {
			v, _ := jsonpath.Get(strings.TrimPrefix(field, "#"), object)
			if v != nil {
				currentStatus = "#CHECKSET"
			}
		}

		for _, failState := range failStates {
			if currentStatus == failState {
				return object, currentStatus, errs.WrapError(errs.Error(errs.FailedToReachTargetStatus, currentStatus))
			}
		}
		return object, currentStatus, nil
	}
}

// DescribeAgentSpace >>> Encapsulated.

// DescribeContextStore <<< Encapsulated get interface for Agentloop ContextStore.

func (s *Service) DescribeContextStore(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	parts := strings.Split(id, ":")
	if len(parts) != 2 {
		return object, errs.WrapError(fmt.Errorf("invalid Resource Id %s. Expected parts' length %d, got %d", id, 2, len(parts)))
	}
	agentSpace := parts[0]
	contextStoreName := parts[1]
	request = make(map[string]interface{})
	query = make(map[string]*string)

	action := fmt.Sprintf("/agentspace/%s/contextstore/%s", agentSpace, contextStoreName)

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.RetryContext(context.Background(), 1*time.Minute, func() *retry.RetryError {
		response, err = client.RoaGet("AgentLoop", "2026-05-20", action, query, nil, nil)

		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"ContextStoreNotExist"}) {
			return object, errs.WrapErrorf(errs.NotFoundErr("ContextStore", id), errs.NotFoundMsg, err)
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}
	if response == nil {
		return object, errs.WrapErrorf(errs.NotFoundErr("ContextStore", id), errs.NotFoundMsg, response)
	}

	return response, nil
}

// DescribeContextStore >>> Encapsulated.

// DescribeContextStoreApiKey <<< Encapsulated get interface for Agentloop ContextStoreApiKey.

func (s *Service) DescribeContextStoreApiKey(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	parts := strings.Split(id, ":")
	if len(parts) != 3 {
		return object, errs.WrapError(fmt.Errorf("invalid Resource Id %s. Expected parts' length %d, got %d", id, 3, len(parts)))
	}
	agentSpace := parts[0]
	contextStoreName := parts[1]
	name := parts[2]
	request = make(map[string]interface{})
	query = make(map[string]*string)

	action := fmt.Sprintf("/agentspace/%s/contextstore/%s/apikey/%s", agentSpace, contextStoreName, name)

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.RetryContext(context.Background(), 1*time.Minute, func() *retry.RetryError {
		response, err = client.RoaGet("AgentLoop", "2026-05-20", action, query, nil, nil)

		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)
	if err != nil {
		// ContextStoreNotExist / AgentSpaceNotExist mean a parent resource was
		// deleted out-of-band; the API key is gone with it, so treat them as
		// NotFound and let the state be cleared on refresh.
		if errs.IsExpectedErrors(err, []string{"NotFound", "ContextStoreNotExist", "AgentSpaceNotExist"}) {
			return object, errs.WrapErrorf(errs.NotFoundErr("ContextStoreApiKey", id), errs.NotFoundMsg, err)
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}
	if response == nil {
		return object, errs.WrapErrorf(errs.NotFoundErr("ContextStoreApiKey", id), errs.NotFoundMsg, response)
	}

	return response, nil
}

// DescribeContextStoreApiKey >>> Encapsulated.

// DescribeDataset <<< Encapsulated get interface for Agentloop Dataset.

func (s *Service) DescribeDataset(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	parts := strings.Split(id, ":")
	if len(parts) != 2 {
		return object, errs.WrapError(fmt.Errorf("invalid Resource Id %s. Expected parts' length %d, got %d", id, 2, len(parts)))
	}
	agentSpace := parts[0]
	datasetName := parts[1]
	request = make(map[string]interface{})
	query = make(map[string]*string)

	action := fmt.Sprintf("/agentspace/%s/dataset/%s", agentSpace, datasetName)

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.RetryContext(context.Background(), 1*time.Minute, func() *retry.RetryError {
		response, err = client.RoaGet("AgentLoop", "2026-05-20", action, query, nil, nil)

		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)
	if err != nil {
		// AgentSpaceNotExist means the hosting AgentSpace was deleted
		// out-of-band; the dataset is gone with it, so treat it as NotFound
		// and let the state be cleared on refresh.
		if errs.IsExpectedErrors(err, []string{"DatasetNotExist", "WorkspaceNotExist", "AgentSpaceNotExist"}) {
			return object, errs.WrapErrorf(errs.NotFoundErr("Dataset", id), errs.NotFoundMsg, err)
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}
	if response == nil {
		return object, errs.WrapErrorf(errs.NotFoundErr("Dataset", id), errs.NotFoundMsg, response)
	}

	return response, nil
}

// DescribeDataset >>> Encapsulated.

// DescribeEndpointConnector <<< Encapsulated get interface for Agentloop EndpointConnector.

func (s *Service) DescribeEndpointConnector(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	parts := strings.Split(id, ":")
	if len(parts) != 2 {
		return object, errs.WrapError(fmt.Errorf("invalid Resource Id %s. Expected parts' length %d, got %d", id, 2, len(parts)))
	}
	agentSpace := parts[0]
	connectorId := parts[1]
	request = make(map[string]interface{})
	query = make(map[string]*string)

	action := fmt.Sprintf("/api/v1/endpoint-connectors/%s/%s", agentSpace, connectorId)

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.RetryContext(context.Background(), 1*time.Minute, func() *retry.RetryError {
		response, err = client.RoaGet("AgentLoop", "2026-05-20", action, query, nil, nil)

		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"NotFound.EndpointConnector"}) {
			return object, errs.WrapErrorf(errs.NotFoundErr("EndpointConnector", id), errs.NotFoundMsg, err)
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}
	if response == nil {
		return object, errs.WrapErrorf(errs.NotFoundErr("EndpointConnector", id), errs.NotFoundMsg, response)
	}

	return response, nil
}

// DescribeEndpointConnector >>> Encapsulated.

// DescribeEvaluationTask <<< Encapsulated get interface for Agentloop EvaluationTask.

func (s *Service) DescribeEvaluationTask(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	parts := strings.Split(id, ":")
	if len(parts) != 2 {
		return object, errs.WrapError(fmt.Errorf("invalid Resource Id %s. Expected parts' length %d, got %d", id, 2, len(parts)))
	}
	agentSpace := parts[0]
	taskId := parts[1]
	request = make(map[string]interface{})
	query = make(map[string]*string)

	action := fmt.Sprintf("/api/v1/evaluation-task/%s/%s", agentSpace, taskId)

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.RetryContext(context.Background(), 1*time.Minute, func() *retry.RetryError {
		response, err = client.RoaGet("AgentLoop", "2026-05-20", action, query, nil, nil)

		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"EvaluationTaskNotFound"}) {
			return object, errs.WrapErrorf(errs.NotFoundErr("EvaluationTask", id), errs.NotFoundMsg, err)
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}
	if response == nil {
		return object, errs.WrapErrorf(errs.NotFoundErr("EvaluationTask", id), errs.NotFoundMsg, response)
	}

	return response, nil
}

// DescribeEvaluationTask >>> Encapsulated.

// DescribeEvaluator <<< Encapsulated get interface for Agentloop Evaluator.

func (s *Service) DescribeEvaluator(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	parts := strings.Split(id, ":")
	if len(parts) != 2 {
		return object, errs.WrapError(fmt.Errorf("invalid Resource Id %s. Expected parts' length %d, got %d", id, 2, len(parts)))
	}
	agentSpace := parts[0]
	name := parts[1]
	request = make(map[string]interface{})
	query = make(map[string]*string)

	action := fmt.Sprintf("/api/v1/evaluators/%s/%s", agentSpace, name)

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.RetryContext(context.Background(), 1*time.Minute, func() *retry.RetryError {
		response, err = client.RoaGet("AgentLoop", "2026-05-20", action, query, nil, nil)

		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"EvaluatorNotFound"}) {
			return object, errs.WrapErrorf(errs.NotFoundErr("Evaluator", id), errs.NotFoundMsg, err)
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}
	if response == nil {
		return object, errs.WrapErrorf(errs.NotFoundErr("Evaluator", id), errs.NotFoundMsg, response)
	}

	return response, nil
}

// DescribeEvaluator >>> Encapsulated.

// DescribeEvaluatorSkill <<< Encapsulated get interface for Agentloop EvaluatorSkill.

func (s *Service) DescribeEvaluatorSkill(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	parts := strings.Split(id, ":")
	if len(parts) != 3 {
		return object, errs.WrapError(fmt.Errorf("invalid Resource Id %s. Expected parts' length %d, got %d", id, 3, len(parts)))
	}
	agentSpace := parts[0]
	name := parts[1]
	skillName := parts[2]
	request = make(map[string]interface{})
	query = make(map[string]*string)
	query["agentSpace"] = helper.StringPointer(agentSpace)

	action := fmt.Sprintf("/api/v1/evaluator/%s/skill/%s", name, skillName)

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.RetryContext(context.Background(), 1*time.Minute, func() *retry.RetryError {
		response, err = client.RoaGet("AgentLoop", "2026-05-20", action, query, nil, nil)

		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"EvaluatorSkillNotFound"}) {
			return object, errs.WrapErrorf(errs.NotFoundErr("EvaluatorSkill", id), errs.NotFoundMsg, err)
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}
	if response == nil {
		return object, errs.WrapErrorf(errs.NotFoundErr("EvaluatorSkill", id), errs.NotFoundMsg, response)
	}

	return response, nil
}

// DescribeEvaluatorSkill >>> Encapsulated.

// DescribeOptimizeTask <<< Encapsulated get interface for Agentloop OptimizeTask.

func (s *Service) DescribeOptimizeTask(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	parts := strings.Split(id, ":")
	if len(parts) != 2 {
		return object, errs.WrapError(fmt.Errorf("invalid Resource Id %s. Expected parts' length %d, got %d", id, 2, len(parts)))
	}
	agentSpace := parts[0]
	optimizeTaskName := parts[1]
	request = make(map[string]interface{})
	query = make(map[string]*string)

	action := fmt.Sprintf("/agentspace/%s/optimizetasks/%s", agentSpace, optimizeTaskName)

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.RetryContext(context.Background(), 1*time.Minute, func() *retry.RetryError {
		response, err = client.RoaGet("AgentLoop", "2026-05-20", action, query, nil, nil)

		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"OptimizeTaskNotExist"}) {
			return object, errs.WrapErrorf(errs.NotFoundErr("OptimizeTask", id), errs.NotFoundMsg, err)
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}
	if response == nil {
		return object, errs.WrapErrorf(errs.NotFoundErr("OptimizeTask", id), errs.NotFoundMsg, response)
	}

	return response, nil
}

// DescribeOptimizeTask >>> Encapsulated.

// DescribePipeline <<< Encapsulated get interface for Agentloop Pipeline.

func (s *Service) DescribePipeline(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	parts := strings.Split(id, ":")
	if len(parts) != 2 {
		return object, errs.WrapError(fmt.Errorf("invalid Resource Id %s. Expected parts' length %d, got %d", id, 2, len(parts)))
	}
	agentSpace := parts[0]
	pipelineName := parts[1]
	request = make(map[string]interface{})
	query = make(map[string]*string)

	action := fmt.Sprintf("/agentspace/%s/pipeline/%s", agentSpace, pipelineName)

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.RetryContext(context.Background(), 1*time.Minute, func() *retry.RetryError {
		response, err = client.RoaGet("AgentLoop", "2026-05-20", action, query, nil, nil)

		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"PipelineNotExist"}) {
			return object, errs.WrapErrorf(errs.NotFoundErr("Pipeline", id), errs.NotFoundMsg, err)
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}
	if response == nil {
		return object, errs.WrapErrorf(errs.NotFoundErr("Pipeline", id), errs.NotFoundMsg, response)
	}

	return response, nil
}

// DescribePipeline >>> Encapsulated.

package alicloud

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/services/oos"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

type OosService struct {
	client *connectivity.AliyunClient
}

func (s *OosService) DescribeOosTemplate(id string) (object oos.Template, err error) {
	request := oos.CreateGetTemplateRequest()
	request.RegionId = s.client.RegionId

	request.TemplateName = id

	raw, err := s.client.WithOosClient(func(oosClient *oos.Client) (interface{}, error) {
		return oosClient.GetTemplate(request)
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"EntityNotExists.Template"}) {
			err = WrapErrorf(Error(GetNotFoundMessage("OosTemplate", id)), NotFoundMsg, ProviderERROR)
			return
		}
		err = WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
		return
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*oos.GetTemplateResponse)
	return response.Template, nil
}

func (s *OosService) DescribeOosExecution(id string) (object oos.Execution, err error) {
	request := oos.CreateListExecutionsRequest()
	request.RegionId = s.client.RegionId

	request.ExecutionId = id

	raw, err := s.client.WithOosClient(func(oosClient *oos.Client) (interface{}, error) {
		return oosClient.ListExecutions(request)
	})
	if err != nil {
		err = WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
		return
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*oos.ListExecutionsResponse)

	if len(response.Executions) < 1 {
		err = WrapErrorf(Error(GetNotFoundMessage("OosExecution", id)), NotFoundMsg, ProviderERROR, response.RequestId)
		return
	}
	return response.Executions[0], nil
}

func (s *OosService) OosExecutionStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeOosExecution(id)
		if err != nil {
			if NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}

		for _, failState := range failStates {
			if object.Status == failState {
				return object, object.Status, WrapError(Error(FailedToReachTargetStatus, object.Status))
			}
		}
		return object, object.Status, nil
	}
}

package alicloud

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/services/mse"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

type MseService struct {
	client *connectivity.AliyunClient
}

func (s *MseService) DescribeMseCluster(id string) (object mse.Data, err error) {
	request := mse.CreateQueryClusterDetailRequest()
	request.RegionId = s.client.RegionId

	request.InstanceId = id

	raw, err := s.client.WithMseClient(func(mseClient *mse.Client) (interface{}, error) {
		return mseClient.QueryClusterDetail(request)
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"mse-200-021"}) {
			err = WrapErrorf(Error(GetNotFoundMessage("MseCluster", id)), NotFoundMsg, ProviderERROR)
			return
		}
		err = WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
		return
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*mse.QueryClusterDetailResponse)
	return response.Data, nil
}

func (s *MseService) MseClusterStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeMseCluster(id)
		if err != nil {
			if NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}

		for _, failState := range failStates {
			if object.InitStatus == failState {
				return object, object.InitStatus, WrapError(Error(FailedToReachTargetStatus, object.InitStatus))
			}
		}
		return object, object.InitStatus, nil
	}
}

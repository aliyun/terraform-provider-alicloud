package alicloud

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/services/actiontrail"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

type ActiontrailService struct {
	client *connectivity.AliyunClient
}

func (s *ActiontrailService) DescribeActiontrailTrail(id string) (object actiontrail.TrailListItem, err error) {
	request := actiontrail.CreateDescribeTrailsRequest()
	request.RegionId = s.client.RegionId

	request.NameList = id

	raw, err := s.client.WithActiontrailClient(func(actiontrailClient *actiontrail.Client) (interface{}, error) {
		return actiontrailClient.DescribeTrails(request)
	})
	if err != nil {
		err = WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
		return
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*actiontrail.DescribeTrailsResponse)

	if len(response.TrailList) < 1 {
		err = WrapErrorf(Error(GetNotFoundMessage("ActiontrailTrail", id)), NotFoundMsg, ProviderERROR, response.RequestId)
		return
	}
	return response.TrailList[0], nil
}

func (s *ActiontrailService) ActiontrailTrailStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeActiontrailTrail(id)
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

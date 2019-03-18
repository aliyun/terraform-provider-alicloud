package alicloud

import (
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/actiontrail"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

type ActionTrailService struct {
	client *connectivity.AliyunClient
}

func (s *ActionTrailService) DescribeActionTrail(id string) (trail actiontrail.TrailListItem, err error) {
	request := actiontrail.CreateDescribeTrailsRequest()

	raw, err := s.client.WithActionTrailClient(func(actiontrailClient *actiontrail.Client) (interface{}, error) {
		return actiontrailClient.DescribeTrails(request)
	})
	if err != nil {
		return trail, WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw)
	response, _ := raw.(*actiontrail.DescribeTrailsResponse)
	if response == nil || len(response.TrailList) < 1 {
		return trail, WrapErrorf(Error(GetNotFoundMessage("ActionTrail", id)), NotFoundMsg, ProviderERROR)
	}
	return response.TrailList[0], nil
}

func (s *ActionTrailService) WaitForActionTrail(id string, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribeActionTrail(id)
		if err != nil {
			return WrapError(err)
		}
		if object.Name == id {
			return nil
		}
		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, object.Name, id, ProviderERROR)
		}
	}
}

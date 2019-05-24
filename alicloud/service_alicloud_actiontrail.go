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
	for _, item := range response.TrailList {
		if item.Name == id {
			return item, nil
		}
	}

	return trail, WrapErrorf(Error(GetNotFoundMessage("ActionTrail", id)), NotFoundMsg, ProviderERROR)
}

func (s *ActionTrailService) startActionTrail(id string) (err error) {
	request := actiontrail.CreateStartLoggingRequest()
	request.Name = id
	request.Method = "GET"
	raw, err := s.client.WithActionTrailClient(func(actiontrailClient *actiontrail.Client) (interface{}, error) {
		return actiontrailClient.StartLogging(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw)
	return nil
}

func (s *ActionTrailService) WaitForActionTrail(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribeActionTrail(id)
		if err != nil {
			if NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return WrapError(err)
			}
		}
		if object.Status == string(status) {
			return nil
		}
		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, object.Status, status, ProviderERROR)
		}
	}
}

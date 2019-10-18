package alicloud

import (
	"time"

	"github.com/hashicorp/terraform/helper/resource"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/emr"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

type EmrService struct {
	client *connectivity.AliyunClient
}

func (s *EmrService) DescribeEmrCluster(id string) (response *emr.DescribeClusterV2Response, err error) {
	request := emr.CreateDescribeClusterV2Request()
	request.Id = id

	raw, err := s.client.WithEmrClient(func(emrClient *emr.Client) (interface{}, error) {
		return emrClient.DescribeClusterV2(request)
	})

	if err != nil {
		err = WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
		return
	}

	addDebug(request.GetActionName(), raw, request.RpcRequest)
	response, _ = raw.(*emr.DescribeClusterV2Response)
	if response.ClusterInfo.Status == "RELEASED" {
		err = WrapErrorf(Error(GetNotFoundMessage("EmrCluster", id)), NotFoundMsg, AlibabaCloudSdkGoERROR)
	}

	return
}

func (s *EmrService) WaitForEmrCluster(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribeEmrCluster(id)
		if err != nil {
			if NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return WrapError(err)
			}
		}

		if object.ClusterInfo.Id == id && status != Deleted {
			break
		}

		time.Sleep(DefaultIntervalShort * time.Second)
		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, object.ClusterInfo.Id, id, ProviderERROR)
		}
	}
	return nil
}

func (s *EmrService) EmrClusterStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeEmrCluster(id)
		if err != nil {
			if NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}

		for _, failState := range failStates {
			if object.ClusterInfo.Status == failState {
				return object, object.ClusterInfo.Status, WrapError(Error(FailedToReachTargetStatus, object.ClusterInfo.Status))
			}
		}

		return object, object.ClusterInfo.Status, nil
	}
}

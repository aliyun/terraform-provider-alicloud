package alicloud

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/services/hbase"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

const (
	Hb_LAUNCHING     = "LAUNCHING"
	Hb_CREATING      = "CREATING"
	Hb_ACTIVATION    = "ACTIVATION"
	Hb_DELETING      = "DELETING"
	Hb_CREATE_FAILED = "CREATE_FAILED"
)

type HBaseService struct {
	client *connectivity.AliyunClient
}

func (s *HBaseService) NotFoundHBaseInstance(err error) bool {
	if NotFoundError(err) || IsExceptedErrors(err, []string{InvalidHBaseInstanceIdNotFound, InvalidHBaseNameNotFound}) {
		return true
	}
	return false
}

func (s *HBaseService) DescribeHBaseInstance(id string) (instance hbase.DescribeInstanceResponse, err error) {
	request := hbase.CreateDescribeInstanceRequest()
	request.RegionId = s.client.RegionId
	request.ClusterId = id
	raw, err := s.client.WithHbaseClient(func(client *hbase.Client) (interface{}, error) {
		return client.DescribeInstance(request)
	})
	response, _ := raw.(*hbase.DescribeInstanceResponse)
	if err != nil {
		if IsExceptedErrors(err, []string{InvalidHBaseInstanceIdNotFound}) {
			return instance, WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
		}
		return instance, WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	if response == nil || response.InstanceId == "" {
		return instance, WrapErrorf(Error(GetNotFoundMessage("HBase Instance", id)), NotFoundMsg, AlibabaCloudSdkGoERROR)
	}
	instance = *response
	return instance, nil
}

/**
pop has limit, support next.
*/
func (s *HBaseService) DescribeIpWhitelist(id string) (instance hbase.DescribeIpWhitelistResponse, err error) {
	request := hbase.CreateDescribeIpWhitelistRequest()
	request.RegionId = s.client.RegionId
	request.ClusterId = id
	raw, err := s.client.WithHbaseClient(func(client *hbase.Client) (interface{}, error) {
		return client.DescribeIpWhitelist(request)
	})
	response, _ := raw.(*hbase.DescribeIpWhitelistResponse)
	if err != nil {
		if IsExceptedErrors(err, []string{InvalidHBaseInstanceIdNotFound}) {
			return instance, WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
		}
		return instance, WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	return *response, nil
}

func (s *HBaseService) HBaseClusterStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeHBaseInstance(id)
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

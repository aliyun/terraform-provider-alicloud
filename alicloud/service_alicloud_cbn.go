package alicloud

import (
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/cbn"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

type CbnService struct {
	client *connectivity.AliyunClient
}

func (s *CbnService) DescribeCenFlowlog(id string) (cenFlowlog cbn.FlowLog, err error) {
	request := cbn.CreateDescribeFlowlogsRequest()
	request.RegionId = s.client.RegionId

	request.FlowLogId = id

	raw, err := s.client.WithCbnClient(func(cbnClient *cbn.Client) (interface{}, error) {
		return cbnClient.DescribeFlowlogs(request)
	})
	if err != nil {
		err = WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
		return
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*cbn.DescribeFlowlogsResponse)

	if len(response.FlowLogs.FlowLog) < 1 {
		err = WrapErrorf(Error(GetNotFoundMessage("CenFlowlog", id)), NotFoundMsg, ProviderERROR)
		return
	}
	return response.FlowLogs.FlowLog[0], nil
}

func (s *CbnService) WaitForCenFlowlog(id string, expected map[string]interface{}, isDelete bool, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribeCenFlowlog(id)
		if err != nil {
			if NotFoundError(err) {
				if isDelete {
					return nil
				}
			} else {
				return WrapError(err)
			}
		}
		ready, current, err := checkWaitForReady(object, expected)
		if err != nil {
			return WrapError(err)
		}
		if ready {
			return nil
		}
		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, current, expected, ProviderERROR)
		}
		time.Sleep(DefaultIntervalShort * time.Second)
	}
}

func (s *CbnService) DescribeCenInstance(id string) (cenInstance cbn.Cen, err error) {
	request := cbn.CreateDescribeCensRequest()
	request.RegionId = s.client.RegionId
	filters := make([]cbn.DescribeCensFilter, 0)
	filters = append(filters, cbn.DescribeCensFilter{
		Key:   "CenId",
		Value: &[]string{id},
	})
	request.Filter = &filters

	raw, err := s.client.WithCbnClient(func(cbnClient *cbn.Client) (interface{}, error) {
		return cbnClient.DescribeCens(request)
	})
	if err != nil {
		err = WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
		return
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*cbn.DescribeCensResponse)

	if len(response.Cens.Cen) < 1 {
		err = WrapErrorf(Error(GetNotFoundMessage("CenInstance", id)), NotFoundMsg, ProviderERROR)
		return
	}
	return response.Cens.Cen[0], nil
}

func (s *CbnService) CenInstanceStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeCenInstance(id)
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

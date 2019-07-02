package alicloud

import (
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ons"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

type OnsService struct {
	client *connectivity.AliyunClient
}

func (s *OnsService) InstanceNotExistFunc(err error) bool {
	return strings.Contains(err.Error(), OnsInstanceNotExist)
}

func (s *OnsService) GetPreventCache() requests.Integer {
	return requests.NewInteger(int(time.Now().UnixNano() / 1e6))
}

func (s *OnsService) DescribeOnsInstance(id string) (*ons.OnsInstanceBaseInfoResponse, error) {
	request := ons.CreateOnsInstanceBaseInfoRequest()
	request.PreventCache = s.GetPreventCache()
	request.InstanceId = id

	raw, err := s.client.WithOnsClient(func(onsClient *ons.Client) (interface{}, error) {
		return onsClient.OnsInstanceBaseInfo(request)
	})

	if err != nil {
		if IsExceptedError(err, InvalidDomainNameNoExist) {
			return nil, WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
		}
		return nil, WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	response, _ := raw.(*ons.OnsInstanceBaseInfoResponse)
	addDebug(request.GetActionName(), raw)
	return response, nil

}

func (s *OnsService) WaitForOnsInstance(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		response, err := s.DescribeOnsInstance(id)

		if err != nil {
			if NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return WrapError(err)
			}
		}

		if response.InstanceBaseInfo.InstanceId == id && status != Deleted {
			return nil
		}

		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, response.InstanceBaseInfo.InstanceId, id, ProviderERROR)
		}
		time.Sleep(DefaultIntervalShort * time.Second)
	}

}

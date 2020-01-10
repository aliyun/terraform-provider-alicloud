package alicloud

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/services/market"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

type MarketService struct {
	client *connectivity.AliyunClient
}

func (s *MarketService) DescribeMarketOrder(id string) (order *market.DescribeOrderResponse, err error) {
	request := market.CreateDescribeOrderRequest()
	request.OrderId = id
	raw, err := s.client.WithMarketClient(func(client *market.Client) (interface{}, error) {
		return client.DescribeOrder(request)
	})
	response, _ := raw.(*market.DescribeOrderResponse)
	if err != nil {
		if IsExceptedErrors(err, []string{"null"}) {
			return order, WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
		}
		return order, WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	if response.OrderId == 0 {
		return order, WrapErrorf(Error(GetNotFoundMessage("Market Order", id)), NotFoundMsg, AlibabaCloudSdkGoERROR)
	}
	return response, nil
}

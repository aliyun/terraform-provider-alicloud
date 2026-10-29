// Package market contains the market service resources and data sources
// migrated from the root alicloud package. The root package cannot be
// imported from a service package, so the encapsulated market calls live
// here.
package market

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/services/market"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
)

// MarketService wraps the market API calls used by the resources of this package.
type MarketService struct {
	client *connectivity.AliyunClient
}

// NewMarketService is the exported constructor: the client field stays
// unexported, so callers outside the package (including the acceptance
// tests) build values with it.
func NewMarketService(client *connectivity.AliyunClient) *MarketService {
	return &MarketService{client: client}
}

func (s *MarketService) DescribeMarketOrder(id string) (order *market.DescribeOrderResponse, err error) {
	request := market.CreateDescribeOrderRequest()
	request.OrderId = id
	raw, err := s.client.WithMarketClient(func(client *market.Client) (interface{}, error) {
		return client.DescribeOrder(request)
	})
	response, _ := raw.(*market.DescribeOrderResponse)
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"null"}) {
			return order, errs.WrapErrorf(err, errs.NotFoundMsg, errs.AlibabaCloudSdkGoERROR)
		}
		return order, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
	}
	logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)
	if response.OrderId == 0 {
		return order, errs.WrapErrorf(errs.NotFoundErr("Market Order", id), errs.NotFoundMsg, errs.AlibabaCloudSdkGoERROR)
	}
	return response, nil
}

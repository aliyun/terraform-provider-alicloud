package natgateway

import (
	"context"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	logs "github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
)

// DescribeNatGatewayCapacityReservation queries the capacity reservation of one
// NAT gateway. The API is a single-object query keyed by NatGatewayId: a gateway
// without a reservation answers with an empty status rather than NotFound.
func DescribeNatGatewayCapacityReservation(ctx context.Context, client *connectivity.AliyunClient, natGatewayId string) (object map[string]interface{}, err error) {
	request := map[string]interface{}{
		"RegionId":     client.RegionId,
		"NatGatewayId": natGatewayId,
	}
	action := "DescribeNatGatewayCapacityReservation"

	response, err := helper.RetryRpcPost(ctx, action, func() (map[string]interface{}, error) {
		return client.RpcPost("Vpc", "2016-04-28", action, nil, request, true)
	}, 1*time.Minute, helper.NeedRetry)
	logs.AddDebug(action, response, request)
	if err != nil {
		return object, err
	}
	return response, nil
}

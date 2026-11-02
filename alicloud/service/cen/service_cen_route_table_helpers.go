package cen

import (
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

// queryRouteTableById and describeRouteTable duplicate the two VPC
// route-table lookups that CreateCenRouteEntryParas needs. They are kept
// local (private) because the VPC product is not yet migrated to a subpackage
// the cen package could import, and the migrated cen package must not import
// the root alicloud package. The root VpcService keeps its own copies for the
// unmigrated consumers. Behavior is identical.

func (s *CenService) queryRouteTableById(routeTableId string) (rt vpc.RouteTable, err error) {
	request := vpc.CreateDescribeRouteTablesRequest()
	request.RegionId = s.client.RegionId
	request.RouteTableId = routeTableId
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		raw, err := s.client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.DescribeRouteTables(request)
		})
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)
		response, _ := raw.(*vpc.DescribeRouteTablesResponse)
		if len(response.RouteTables.RouteTable) == 0 ||
			response.RouteTables.RouteTable[0].RouteTableId != routeTableId {
			return retry.NonRetryableError(errs.WrapErrorf(errs.NotFoundErr("RouteTable", routeTableId), errs.NotFoundMsg, errs.ProviderERROR))
		}
		rt = response.RouteTables.RouteTable[0]
		return nil
	})
	return
}

func (s *CenService) describeRouteTable(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "DescribeRouteTableList"
	request := map[string]interface{}{
		"RegionId":     s.client.RegionId,
		"RouteTableId": id,
	}
	response, err = client.RpcPost("Vpc", "2016-04-28", action, nil, request, true)
	if err != nil {
		err = errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
		return
	}
	logs.AddDebug(action, response, request)
	v, err := jsonpath.Get("$.RouterTableList.RouterTableListType", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.RouterTableList.RouterTableListType", response)
	}
	if len(v.([]interface{})) < 1 {
		return object, errs.WrapErrorf(errs.NotFoundErr("VPC", id), errs.NotFoundWithResponse, response)
	} else {
		if v.([]interface{})[0].(map[string]interface{})["RouteTableId"].(string) != id {
			return object, errs.WrapErrorf(errs.NotFoundErr("VPC", id), errs.NotFoundWithResponse, response)
		}
	}
	object = v.([]interface{})[0].(map[string]interface{})
	return object, nil
}

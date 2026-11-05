package expressconnect

import (
	"fmt"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func (s *ExpressConnectService) DescribeRouterInterface(id, regionId string) (ri vpc.RouterInterfaceType, err error) {
	request := vpc.CreateDescribeRouterInterfacesRequest()
	request.RegionId = regionId
	values := []string{id}
	filter := []vpc.DescribeRouterInterfacesFilter{
		{
			Key:   "RouterInterfaceId",
			Value: &values,
		},
	}
	request.Filter = &filter
	invoker := helper.NewInvoker()
	err = invoker.Run(func() error {
		raw, err := s.client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.DescribeRouterInterfaces(request)
		})
		if err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, id, request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
		}
		logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)
		response, _ := raw.(*vpc.DescribeRouterInterfacesResponse)
		if len(response.RouterInterfaceSet.RouterInterfaceType) <= 0 ||
			response.RouterInterfaceSet.RouterInterfaceType[0].RouterInterfaceId != id {
			return errs.WrapErrorf(errs.NotFoundErr("RouterInterface", id), NotFoundMsg, errs.ProviderERROR)
		}
		ri = response.RouterInterfaceSet.RouterInterfaceType[0]
		return nil
	})
	return
}

func (s *ExpressConnectService) DescribeRouterInterfaceConnection(id, regionId string) (ri vpc.RouterInterfaceType, err error) {
	ri, err = s.DescribeRouterInterface(id, regionId)
	if err != nil {
		return ri, errs.WrapError(err)
	}

	if ri.OppositeInterfaceId == "" || ri.OppositeRouterType == "" ||
		ri.OppositeRouterId == "" || ri.OppositeInterfaceOwnerId == "" {
		return ri, errs.WrapErrorf(errs.NotFoundErr("RouterInterface", id), NotFoundMsg, errs.ProviderERROR)
	}
	return ri, nil
}

func (s *ExpressConnectService) WaitForRouterInterface(id, regionId string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribeRouterInterface(id, regionId)
		if err != nil {
			if errs.NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return errs.WrapError(err)
			}
		}
		if object.Status == string(status) {
			return nil
		}
		if time.Now().After(deadline) {
			return errs.WrapErrorf(err, WaitTimeoutMsg, id, helper.GetFunc(1), timeout, object.Status, string(status), errs.ProviderERROR)
		}
		time.Sleep(helper.DefaultIntervalShort * time.Second)
	}
}

func (s *ExpressConnectService) WaitForRouterInterfaceConnection(id, regionId string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribeRouterInterfaceConnection(id, regionId)
		if err != nil {
			if errs.NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return errs.WrapError(err)
			}
		}
		if object.Status == string(status) {
			return nil
		}
		if time.Now().After(deadline) {
			return errs.WrapErrorf(err, WaitTimeoutMsg, id, helper.GetFunc(1), timeout, object.Status, string(status), errs.ProviderERROR)
		}
		time.Sleep(helper.DefaultIntervalShort * time.Second)
	}
}

func (s *ExpressConnectService) DeactivateRouterInterface(interfaceId string) error {
	request := vpc.CreateDeactivateRouterInterfaceRequest()
	request.RegionId = s.client.RegionId
	request.RouterInterfaceId = interfaceId
	raw, err := s.client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
		return vpcClient.DeactivateRouterInterface(request)
	})
	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "RouterInterface", request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
	}
	logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)
	return nil
}

func (s *ExpressConnectService) ActivateRouterInterface(interfaceId string) error {
	request := vpc.CreateActivateRouterInterfaceRequest()
	request.RegionId = s.client.RegionId
	request.RouterInterfaceId = interfaceId
	raw, err := s.client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
		return vpcClient.ActivateRouterInterface(request)
	})
	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "RouterInterface", request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
	}
	logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)
	return nil
}

func (s *ExpressConnectService) DescribeExpressConnectPhysicalConnection(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	action := "DescribePhysicalConnections"

	client := s.client

	request := map[string]interface{}{
		"RegionId":   s.client.RegionId,
		"PageSize":   helper.PageSizeLarge,
		"PageNumber": 1,
	}

	filterMapList := make([]map[string]interface{}, 0)
	filterMapList = append(filterMapList, map[string]interface{}{
		"Key":   "PhysicalConnectionId",
		"Value": []string{id},
	})
	request["Filter"] = filterMapList

	idExist := false
	for {
		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(5*time.Minute, func() *retry.RetryError {
			response, err = client.RpcPost("Vpc", "2016-04-28", action, nil, request, true)
			if err != nil {
				if helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		logs.AddDebug(action, response, request)

		if err != nil {
			return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
		}

		resp, err := jsonpath.Get("$.PhysicalConnectionSet.PhysicalConnectionType", response)
		if err != nil {
			return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.PhysicalConnectionSet.PhysicalConnectionType", response)
		}

		if v, ok := resp.([]interface{}); !ok || len(v) < 1 {
			return object, errs.WrapErrorf(errs.NotFoundErr("ExpressConnect:PhysicalConnection", id), errs.NotFoundWithResponse, response)
		}

		for _, v := range resp.([]interface{}) {
			if fmt.Sprint(v.(map[string]interface{})["PhysicalConnectionId"]) == id {
				idExist = true
				return v.(map[string]interface{}), nil
			}
		}

		if len(resp.([]interface{})) < request["PageSize"].(int) {
			break
		}

		request["PageNumber"] = request["PageNumber"].(int) + 1
	}

	if !idExist {
		return object, errs.WrapErrorf(errs.NotFoundErr("ExpressConnect:PhysicalConnection", id), errs.NotFoundWithResponse, response)
	}

	return object, nil
}

func (s *ExpressConnectService) ExpressConnectPhysicalConnectionStateRefreshFunc(id string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeExpressConnectPhysicalConnection(id)
		if err != nil {
			if errs.NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", errs.WrapError(err)
		}

		for _, failState := range failStates {
			if fmt.Sprint(object["Status"]) == failState {
				return object, fmt.Sprint(object["Status"]), errs.WrapError(errs.Error(FailedToReachTargetStatus, fmt.Sprint(object["Status"])))
			}
		}

		return object, fmt.Sprint(object["Status"]), nil
	}
}

func (s *ExpressConnectService) ModifyExpressConnectPhysicalConnectionStatus(d *schema.ResourceData, status string) (err error) {
	var response map[string]interface{}
	client := s.client

	switch status {
	case "Confirmed":
		action := "ConfirmPhysicalConnection"

		confirmedReq := map[string]interface{}{
			"RegionId":             s.client.RegionId,
			"ClientToken":          helper.BuildClientToken("ConfirmPhysicalConnection"),
			"PhysicalConnectionId": d.Id(),
		}

		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(20*time.Minute, func() *retry.RetryError {
			response, err = client.RpcPost("Vpc", "2016-04-28", action, nil, confirmedReq, true)
			if err != nil {
				if helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		logs.AddDebug(action, response, confirmedReq)

		if err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
		}

		stateConf := helper.BuildStateConf([]string{}, []string{"Confirmed"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, s.ExpressConnectPhysicalConnectionStateRefreshFunc(d.Id(), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return errs.WrapErrorf(err, IdMsg, d.Id())
		}
	case "Enabled":
		action := "CreatePhysicalConnectionOccupancyOrder"

		createPhysicalConnectionOccupancyOrderReq := map[string]interface{}{
			"RegionId":             s.client.RegionId,
			"ClientToken":          helper.BuildClientToken("CreatePhysicalConnectionOccupancyOrder"),
			"PhysicalConnectionId": d.Id(),
			"InstanceChargeType":   "PrePaid",
			"AutoPay":              true,
		}

		if v, ok := d.GetOkExists("period"); ok {
			createPhysicalConnectionOccupancyOrderReq["Period"] = v
		}

		if v, ok := d.GetOk("pricing_cycle"); ok {
			createPhysicalConnectionOccupancyOrderReq["PricingCycle"] = v
		}

		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(20*time.Minute, func() *retry.RetryError {
			response, err = client.RpcPost("Vpc", "2016-04-28", action, nil, createPhysicalConnectionOccupancyOrderReq, true)
			if err != nil {
				if helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		logs.AddDebug(action, response, createPhysicalConnectionOccupancyOrderReq)

		if err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
		}

		if resp, err := jsonpath.Get("$.Data", response); err != nil || resp == nil {
			return errs.WrapErrorf(err, errs.FailedGetAttributeMsg, d.Id(), "$.Data", response)
		} else {
			orderId := resp.(map[string]interface{})["OrderId"]
			d.Set("order_id", orderId)
		}

		stateConf := helper.BuildStateConf([]string{}, []string{"Enabled"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, s.ExpressConnectPhysicalConnectionStateRefreshFunc(d.Id(), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return errs.WrapErrorf(err, IdMsg, d.Id())
		}
	case "Canceled":
		action := "CancelPhysicalConnection"

		canceledReq := map[string]interface{}{
			"RegionId":             s.client.RegionId,
			"ClientToken":          helper.BuildClientToken("CancelPhysicalConnection"),
			"PhysicalConnectionId": d.Id(),
		}

		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(20*time.Minute, func() *retry.RetryError {
			response, err = client.RpcPost("Vpc", "2016-04-28", action, nil, canceledReq, true)
			if err != nil {
				if helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		logs.AddDebug(action, response, canceledReq)

		if err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
		}

		stateConf := helper.BuildStateConf([]string{}, []string{"Canceled"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, s.ExpressConnectPhysicalConnectionStateRefreshFunc(d.Id(), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return errs.WrapErrorf(err, IdMsg, d.Id())
		}
	case "Terminated":
		action := "TerminatePhysicalConnection"

		terminatedReq := map[string]interface{}{
			"RegionId":             s.client.RegionId,
			"ClientToken":          helper.BuildClientToken("TerminatePhysicalConnection"),
			"PhysicalConnectionId": d.Id(),
		}

		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(20*time.Minute, func() *retry.RetryError {
			response, err = client.RpcPost("Vpc", "2016-04-28", action, nil, terminatedReq, true)
			if err != nil {
				if helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		logs.AddDebug(action, response, terminatedReq)

		if err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
		}

		stateConf := helper.BuildStateConf([]string{}, []string{"Terminated"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, s.ExpressConnectPhysicalConnectionStateRefreshFunc(d.Id(), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return errs.WrapErrorf(err, IdMsg, d.Id())
		}
	}

	return nil
}

func (s *ExpressConnectService) JudgeRegionValidation(key, region string) error {
	request := ecs.CreateDescribeRegionsRequest()
	request.RegionId = s.client.RegionId
	raw, err := s.client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.DescribeRegions(request)
	})
	if err != nil {
		return fmt.Errorf("DescribeRegions got an error: %#v", err)
	}
	logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)
	resp, _ := raw.(*ecs.DescribeRegionsResponse)
	if resp == nil || len(resp.Regions.Region) < 1 {
		return errs.GetNotFoundErrorFromString("There is no any available region.")
	}

	var rs []string
	for _, v := range resp.Regions.Region {
		if v.RegionId == region {
			return nil
		}
		rs = append(rs, v.RegionId)
	}
	return fmt.Errorf("'%s' is invalid. Expected on %v.", key, strings.Join(rs, ", "))
}

// DescribeZone validate zoneId is valid in region

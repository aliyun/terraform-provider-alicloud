package expressconnect

import (
	"fmt"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type ExpressConnectService struct {
	client *connectivity.AliyunClient
}

func NewExpressConnectService(client *connectivity.AliyunClient) *ExpressConnectService {
	return &ExpressConnectService{client: client}
}

func (s *ExpressConnectService) DescribeExpressConnectVbrPconnAssociation(id string) (object map[string]interface{}, err error) {
	client := s.client
	parts, err := helper.ParseResourceId(id, 2)
	if err != nil {
		return object, errs.WrapError(err)
	}

	request := map[string]interface{}{
		"Filter.1.Key":     "VbrId",
		"Filter.1.Value.1": parts[0],
		"RegionId":         s.client.RegionId,
	}

	var response map[string]interface{}
	action := "DescribeVirtualBorderRouters"
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		resp, err := client.RpcPost("Vpc", "2016-04-28", action, nil, request, true)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		response = resp
		logs.AddDebug(action, response, request)
		return nil
	})
	if err != nil {
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.VirtualBorderRouterSet.VirtualBorderRouterType[*].AssociatedPhysicalConnections.AssociatedPhysicalConnection[*]", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.VirtualBorderRouterSet.VirtualBorderRouterType[*].AssociatedPhysicalConnections.AssociatedPhysicalConnection[*]", response)
	}
	if len(v.([]interface{})) < 1 {
		return object, errs.WrapErrorf(errs.NotFoundErr("VbrPconnAssociation", id), errs.NotFoundWithResponse, response)
	} else {
		for _, vv := range v.([]interface{}) {
			item := vv.(map[string]interface{})

			if fmt.Sprint(item["PhysicalConnectionId"]) == parts[1] {
				return item, nil
			}
		}
	}
	return object, errs.WrapErrorf(errs.NotFoundErr("VbrPconnAssociation", id), errs.NotFoundWithResponse, response)
}

func (s *ExpressConnectService) ExpressConnectVbrPconnAssociationStateRefreshFunc(d *schema.ResourceData, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeExpressConnectVbrPconnAssociation(d.Id())
		if err != nil {
			if errs.NotFoundError(err) {
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

func (s *ExpressConnectService) DescribeExpressConnectVirtualPhysicalConnection(id string) (object map[string]interface{}, err error) {
	client := s.client

	request := map[string]interface{}{
		"VirtualPhysicalConnectionIds.1": id,
		"RegionId":                       s.client.RegionId,
	}

	var response map[string]interface{}
	action := "ListVirtualPhysicalConnections"
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		resp, err := client.RpcPost("Vpc", "2016-04-28", action, nil, request, true)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		response = resp
		logs.AddDebug(action, response, request)
		return nil
	})
	if err != nil {
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.VirtualPhysicalConnections", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.VirtualPhysicalConnections", response)
	}
	if len(v.([]interface{})) < 1 {
		return object, errs.WrapErrorf(errs.NotFoundErr("VirtualPhysicalConnections", id), errs.NotFoundWithResponse, response)
	}
	return v.([]interface{})[0].(map[string]interface{}), nil
}

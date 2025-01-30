package alicloud

import (
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	"fmt"
)

type ExpressConnectService struct {
	client *connectivity.AliyunClient
}

func (s *ExpressConnectService) DescribeExpressConnectVbrPconnAssociation(id string) (object map[string]interface{}, err error) {
	client := s.client
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		return object, WrapError(err)
	}

	request := map[string]interface{}{
		"Filter.1.Key":     "VbrId",
		"Filter.1.Value.1": parts[0],
		"RegionId":         s.client.RegionId,
	}

	var response map[string]interface{}
	action := "DescribeVirtualBorderRouters"
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		resp, err := client.RpcPost("Vpc", "2016-04-28", action, nil, request, true)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		response = resp
		addDebug(action, response, request)
		return nil
	})
	if err != nil {
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.VirtualBorderRouterSet.VirtualBorderRouterType[*].AssociatedPhysicalConnections.AssociatedPhysicalConnection[*]", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.VirtualBorderRouterSet.VirtualBorderRouterType[*].AssociatedPhysicalConnections.AssociatedPhysicalConnection[*]", response)
	}
	if len(v.([]interface{})) < 1 {
		return object, WrapErrorf(Error(GetNotFoundMessage("VbrPconnAssociation", id)), NotFoundWithResponse, response)
	} else {
		for _, vv := range v.([]interface{}) {
			item := vv.(map[string]interface{})

			if fmt.Sprint(item["PhysicalConnectionId"]) == parts[1] {
				return item, nil
			}
		}
	}
	return object, WrapErrorf(Error(GetNotFoundMessage("VbrPconnAssociation", id)), NotFoundWithResponse, response)
}

func (s *ExpressConnectService) ExpressConnectVbrPconnAssociationStateRefreshFunc(d *schema.ResourceData, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeExpressConnectVbrPconnAssociation(d.Id())
		if err != nil {
			if NotFoundError(err) {
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}
		for _, failState := range failStates {
			if fmt.Sprint(object["Status"]) == failState {
				return object, fmt.Sprint(object["Status"]), WrapError(Error(FailedToReachTargetStatus, fmt.Sprint(object["Status"])))
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
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		resp, err := client.RpcPost("Vpc", "2016-04-28", action, nil, request, true)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		response = resp
		addDebug(action, response, request)
		return nil
	})
	if err != nil {
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.VirtualPhysicalConnections", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.VirtualPhysicalConnections", response)
	}
	if len(v.([]interface{})) < 1 {
		return object, WrapErrorf(Error(GetNotFoundMessage("VirtualPhysicalConnections", id)), NotFoundWithResponse, response)
	}
	return v.([]interface{})[0].(map[string]interface{}), nil
}

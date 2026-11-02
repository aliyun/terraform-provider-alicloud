package cen

import (
	"fmt"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/cbn"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/convert"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/tags"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type CbnService struct {
	client *connectivity.AliyunClient
}

// NewCbnService is the exported constructor: the client field stays
// unexported, so callers outside the package (including the acceptance
// tests in package cen_test and the cross-product tests that patch
// Cen-grant methods) build values with it.
func NewCbnService(client *connectivity.AliyunClient) *CbnService {
	return &CbnService{client: client}
}

func (s *CbnService) DescribeCenFlowlog(id string) (object cbn.FlowLog, err error) {
	request := cbn.CreateDescribeFlowlogsRequest()
	request.RegionId = s.client.RegionId

	request.FlowLogId = id
	var raw interface{}
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		raw, err = s.client.WithCbnClient(func(cbnClient *cbn.Client) (interface{}, error) {
			return cbnClient.DescribeFlowlogs(request)
		})
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})

	if err != nil {
		err = errs.WrapErrorf(err, errs.DefaultErrorMsg, id, request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
		return
	}
	response, _ := raw.(*cbn.DescribeFlowlogsResponse)

	if len(response.FlowLogs.FlowLog) < 1 {
		err = errs.WrapErrorf(errs.NotFoundErr("CenFlowlog", id), errs.NotFoundMsg, errs.ProviderERROR)
		return
	}
	return response.FlowLogs.FlowLog[0], nil
}

func (s *CbnService) WaitForCenFlowlog(id string, expected map[string]interface{}, isDelete bool, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribeCenFlowlog(id)
		if err != nil {
			if errs.NotFoundError(err) {
				if isDelete {
					return nil
				}
			} else {
				return errs.WrapError(err)
			}
		}
		ready, current, err := helper.CheckWaitForReady(object, expected)
		if err != nil {
			return errs.WrapError(err)
		}
		if ready {
			return nil
		}
		if time.Now().After(deadline) {
			return errs.WrapErrorf(err, errs.WaitTimeoutMsg, id, helper.GetFunc(1), timeout, current, expected, errs.ProviderERROR)
		}
		time.Sleep(helper.DefaultIntervalShort * time.Second)
	}
}

func (s *CbnService) DescribeCenInstance(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	action := "DescribeCens"
	client := s.client

	filterMaps := make([]map[string]interface{}, 0)
	filterMaps = append(filterMaps, map[string]interface{}{
		"Key":   "CenId",
		"Value": []string{id},
	})

	request := map[string]interface{}{
		"Filter":     filterMaps,
		"PageSize":   helper.PageSizeLarge,
		"PageNumber": 1,
	}

	idExist := false
	for {
		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(5*time.Minute, func() *retry.RetryError {
			response, err = client.RpcPost("Cbn", "2017-09-12", action, nil, request, true)
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

		resp, err := jsonpath.Get("$.Cens.Cen", response)
		if err != nil {
			return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.Cens.Cen", response)
		}

		if v, ok := resp.([]interface{}); !ok || len(v) < 1 {
			return object, errs.WrapErrorf(errs.NotFoundErr("Cen:Instance", id), errs.NotFoundWithResponse, response)
		}

		for _, v := range resp.([]interface{}) {
			if fmt.Sprint(v.(map[string]interface{})["CenId"]) == id {
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
		return object, errs.WrapErrorf(errs.NotFoundErr("Cen:Instance", id), errs.NotFoundWithResponse, response)
	}

	return object, nil
}

func (s *CbnService) CenInstanceStateRefreshFunc(d *schema.ResourceData, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeCenInstance(d.Id())
		if err != nil {
			if errs.NotFoundError(err) {
				return nil, "", nil
			}
			return nil, "", errs.WrapError(err)
		}

		for _, failState := range failStates {
			if fmt.Sprint(object["Status"]) == failState {
				return object, fmt.Sprint(object["Status"]), errs.WrapError(errs.Error(errs.FailedToReachTargetStatus, fmt.Sprint(object["Status"])))
			}
		}

		return object, fmt.Sprint(object["Status"]), nil
	}
}

func (s *CbnService) setResourceTags(d *schema.ResourceData, resourceType string) error {
	oldItems, newItems := d.GetChange("tags")
	added := make([]cbn.TagResourcesTag, 0)
	for key, value := range newItems.(map[string]interface{}) {
		added = append(added, cbn.TagResourcesTag{
			Key:   key,
			Value: value.(string),
		})
	}
	removed := make([]string, 0)
	for key := range oldItems.(map[string]interface{}) {
		removed = append(removed, key)
	}
	if len(removed) > 0 {
		request := cbn.CreateUntagResourcesRequest()
		request.RegionId = s.client.RegionId
		request.ResourceId = &[]string{d.Id()}
		request.ResourceType = resourceType
		request.TagKey = &removed
		raw, err := s.client.WithCbnClient(func(cbnClient *cbn.Client) (interface{}, error) {
			return cbnClient.UntagResources(request)
		})
		logs.AddDebug(request.GetActionName(), raw)
		if err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
		}
	}
	if len(added) > 0 {
		request := cbn.CreateTagResourcesRequest()
		request.RegionId = s.client.RegionId
		request.ResourceId = &[]string{d.Id()}
		request.ResourceType = resourceType
		request.Tag = &added
		raw, err := s.client.WithCbnClient(func(cbnClient *cbn.Client) (interface{}, error) {
			return cbnClient.TagResources(request)
		})
		logs.AddDebug(request.GetActionName(), raw)
		if err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
		}
	}
	return nil
}

func (s *CbnService) DescribeCenRouteMap(id string) (object cbn.RouteMap, err error) {
	parts, err := helper.ParseResourceId(id, 2)
	if err != nil {
		err = errs.WrapError(err)
		return
	}
	request := cbn.CreateDescribeCenRouteMapsRequest()
	request.RegionId = s.client.RegionId
	request.CenId = parts[0]
	request.RouteMapId = parts[1]

	var raw interface{}
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		raw, err = s.client.WithCbnClient(func(cbnClient *cbn.Client) (interface{}, error) {
			return cbnClient.DescribeCenRouteMaps(request)
		})
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})

	if err != nil {
		err = errs.WrapErrorf(err, errs.DefaultErrorMsg, id, request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
		return
	}
	response, _ := raw.(*cbn.DescribeCenRouteMapsResponse)

	if len(response.RouteMaps.RouteMap) < 1 {
		err = errs.WrapErrorf(errs.NotFoundErr("CenRouteMap", id), errs.NotFoundMsg, errs.ProviderERROR)
		return
	}
	return response.RouteMaps.RouteMap[0], nil
}

func (s *CbnService) CenRouteMapStateRefreshFunc(id string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeCenRouteMap(id)
		if err != nil {
			if errs.NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", errs.WrapError(err)
		}

		for _, failState := range failStates {
			if object.Status == failState {
				return object, object.Status, errs.WrapError(errs.Error(errs.FailedToReachTargetStatus, object.Status))
			}
		}
		return object, object.Status, nil
	}
}

func (s *CbnService) DescribeCenPrivateZone(id string) (object cbn.PrivateZoneInfo, err error) {
	parts, err := helper.ParseResourceId(id, 2)
	if err != nil {
		err = errs.WrapError(err)
		return
	}
	request := cbn.CreateDescribeCenPrivateZoneRoutesRequest()
	request.CenId = parts[0]
	request.AccessRegionId = parts[1]

	var raw interface{}
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		raw, err = s.client.WithCbnClient(func(cbnClient *cbn.Client) (interface{}, error) {
			return cbnClient.DescribeCenPrivateZoneRoutes(request)
		})
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)

	if err != nil {
		if errs.IsExpectedErrors(err, []string{"ParameterIllegal.CenInstanceId"}) {
			return object, errs.WrapErrorf(errs.NotFoundErr("Cen:PrivateZone", id), errs.NotFoundMsg, errs.ProviderERROR)
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
	}

	response, _ := raw.(*cbn.DescribeCenPrivateZoneRoutesResponse)

	if len(response.PrivateZoneInfos.PrivateZoneInfo) < 1 {
		err = errs.WrapErrorf(errs.NotFoundErr("Cen:PrivateZone", id), errs.NotFoundMsg, errs.ProviderERROR)
		return
	}

	return response.PrivateZoneInfos.PrivateZoneInfo[0], nil
}

func (s *CbnService) CenPrivateZoneStateRefreshFunc(id string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeCenPrivateZone(id)
		if err != nil {
			if errs.NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", errs.WrapError(err)
		}

		for _, failState := range failStates {
			if object.Status == failState {
				return object, object.Status, errs.WrapError(errs.Error(errs.FailedToReachTargetStatus, object.Status))
			}
		}
		return object, object.Status, nil
	}
}

func (s *CbnService) DescribeCenVbrHealthCheck(id string) (object cbn.VbrHealthCheck, err error) {
	parts, err := helper.ParseResourceId(id, 2)
	if err != nil {
		err = errs.WrapError(err)
		return
	}
	request := cbn.CreateDescribeCenVbrHealthCheckRequest()
	request.RegionId = s.client.RegionId
	request.VbrInstanceId = parts[0]
	request.VbrInstanceRegionId = parts[1]

	var raw interface{}
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		raw, err = s.client.WithCbnClient(func(cbnClient *cbn.Client) (interface{}, error) {
			return cbnClient.DescribeCenVbrHealthCheck(request)
		})
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})

	if err != nil {
		err = errs.WrapErrorf(err, errs.DefaultErrorMsg, id, request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
		return
	}
	response, _ := raw.(*cbn.DescribeCenVbrHealthCheckResponse)

	if len(response.VbrHealthChecks.VbrHealthCheck) < 1 {
		err = errs.WrapErrorf(errs.NotFoundErr("CenVbrHealthCheck", id), errs.NotFoundMsg, errs.ProviderERROR, response.RequestId)
		return
	}
	return response.VbrHealthChecks.VbrHealthCheck[0], nil
}

func (s *CbnService) DescribeCenInstanceAttachment(id string) (object cbn.DescribeCenAttachedChildInstanceAttributeResponse, err error) {
	parts, err := helper.ParseResourceId(id, 4)
	if err != nil {
		err = errs.WrapError(err)
		return
	}
	request := cbn.CreateDescribeCenAttachedChildInstanceAttributeRequest()
	request.RegionId = s.client.RegionId
	request.ChildInstanceId = parts[1]
	request.ChildInstanceRegionId = parts[3]
	request.ChildInstanceType = parts[2]
	request.CenId = parts[0]

	var raw interface{}
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		raw, err = s.client.WithCbnClient(func(cbnClient *cbn.Client) (interface{}, error) {
			return cbnClient.DescribeCenAttachedChildInstanceAttribute(request)
		})
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})

	if err != nil {
		if errs.IsExpectedErrors(err, []string{"ParameterCenInstanceId", "ParameterError", "ParameterInstanceId"}) {
			err = errs.WrapErrorf(errs.NotFoundErr("CenInstanceAttachment", id), errs.NotFoundMsg, errs.ProviderERROR)
			return
		}
		err = errs.WrapErrorf(err, errs.DefaultErrorMsg, id, request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
		return
	}
	response, _ := raw.(*cbn.DescribeCenAttachedChildInstanceAttributeResponse)
	return *response, nil
}

func (s *CbnService) CenInstanceAttachmentStateRefreshFunc(id string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeCenInstanceAttachment(id)
		if err != nil {
			if errs.NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", errs.WrapError(err)
		}

		for _, failState := range failStates {
			if object.Status == failState {
				return object, object.Status, errs.WrapError(errs.Error(errs.FailedToReachTargetStatus, object.Status))
			}
		}
		return object, object.Status, nil
	}
}

func (s *CbnService) DescribeCenBandwidthPackage(id string) (object cbn.CenBandwidthPackage, err error) {
	request := cbn.CreateDescribeCenBandwidthPackagesRequest()
	request.RegionId = s.client.RegionId
	filters := make([]cbn.DescribeCenBandwidthPackagesFilter, 0)
	filters = append(filters, cbn.DescribeCenBandwidthPackagesFilter{
		Key:   "CenBandwidthPackageId",
		Value: &[]string{id},
	})
	request.Filter = &filters

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(11*time.Minute, func() *retry.RetryError {
		raw, err := s.client.WithCbnClient(func(cbnClient *cbn.Client) (interface{}, error) {
			return cbnClient.DescribeCenBandwidthPackages(request)
		})
		if err != nil {
			if errs.IsExpectedErrors(err, []string{"AliyunGoClientFailure", "ServiceUnavailable", "Throttling", "Throttling.User"}) {
				wait()
				return retry.RetryableError(err)
			}
			if errs.IsExpectedErrors(err, []string{"ParameterCenInstanceId"}) {
				err = errs.WrapErrorf(errs.NotFoundErr("CenBandwidthPackage", id), errs.NotFoundMsg, errs.ProviderERROR)
				return retry.NonRetryableError(err)
			}
			err = errs.WrapErrorf(err, errs.DefaultErrorMsg, id, request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
			return retry.NonRetryableError(err)
		}
		logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)
		response, _ := raw.(*cbn.DescribeCenBandwidthPackagesResponse)

		if len(response.CenBandwidthPackages.CenBandwidthPackage) < 1 {
			err = errs.WrapErrorf(errs.NotFoundErr("CenBandwidthPackage", id), errs.NotFoundMsg, errs.ProviderERROR, response.RequestId)
			return retry.NonRetryableError(err)
		}
		object = response.CenBandwidthPackages.CenBandwidthPackage[0]
		return nil
	})
	return object, errs.WrapError(err)
}

func (s *CbnService) CenBandwidthPackageStateRefreshFunc(id string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeCenBandwidthPackage(id)
		if err != nil {
			if errs.NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", errs.WrapError(err)
		}

		for _, failState := range failStates {
			if object.Status == failState {
				return object, object.Status, errs.WrapError(errs.Error(errs.FailedToReachTargetStatus, object.Status))
			}
		}
		return object, object.Status, nil
	}
}

func (s *CbnService) ListTagResources(id string, resourceType string) (object interface{}, err error) {
	client := s.client
	action := "ListTagResources"

	request := map[string]interface{}{
		"RegionId":     s.client.RegionId,
		"ResourceType": resourceType,
	}

	resourceIdNum := strings.Count(id, ":")

	switch resourceIdNum {
	case 0:
		request["ResourceId.1"] = id
	case 1:
		parts, err := helper.ParseResourceId(id, 2)
		if err != nil {
			return object, errs.WrapError(err)
		}
		request["ResourceId.1"] = parts[resourceIdNum]
	}

	tags := make([]interface{}, 0)
	var response map[string]interface{}

	for {
		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(5*time.Minute, func() *retry.RetryError {
			response, err := client.RpcPost("Cbn", "2017-09-12", action, nil, request, false)
			if err != nil {
				if helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			logs.AddDebug(action, response, request)
			v, err := jsonpath.Get("$.TagResources.TagResource", response)
			if err != nil {
				return retry.NonRetryableError(errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.TagResources.TagResource", response))
			}

			if v != nil {
				tags = append(tags, v.([]interface{})...)
			}

			return nil
		})
		if err != nil {
			err = errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
			return
		}
		if response["NextToken"] == nil {
			break
		}
		request["NextToken"] = response["NextToken"]
	}

	return tags, nil
}

func (s *CbnService) SetResourceTags(d *schema.ResourceData, resourceType string) error {

	resourceIdNum := strings.Count(d.Id(), ":")

	if d.HasChange("tags") {
		added, removed := tags.ParsingTags(d)
		client := s.client

		removedTagKeys := make([]string, 0)
		for _, v := range removed {
			if !tags.IgnoredTags(v, "") {
				removedTagKeys = append(removedTagKeys, v)
			}
		}

		if len(removedTagKeys) > 0 {
			action := "UntagResources"
			request := map[string]interface{}{
				"RegionId":     s.client.RegionId,
				"ResourceType": resourceType,
			}

			switch resourceIdNum {
			case 0:
				request["ResourceId.1"] = d.Id()
			case 1:
				parts, err := helper.ParseResourceId(d.Id(), 2)
				if err != nil {
					return errs.WrapError(err)
				}
				request["ResourceId.1"] = parts[resourceIdNum]
			}

			for i, key := range removedTagKeys {
				request[fmt.Sprintf("TagKey.%d", i+1)] = key
			}
			wait := helper.IncrementalWait(2*time.Second, 1*time.Second)
			err := retry.Retry(10*time.Minute, func() *retry.RetryError {
				response, err := client.RpcPost("Cbn", "2017-09-12", action, nil, request, false)
				if err != nil {
					if helper.NeedRetry(err) {
						wait()
						return retry.RetryableError(err)

					}
					return retry.NonRetryableError(err)
				}
				logs.AddDebug(action, response, request)
				return nil
			})
			if err != nil {
				return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
			}
		}
		if len(added) > 0 {
			action := "TagResources"
			request := map[string]interface{}{
				"RegionId":     s.client.RegionId,
				"ResourceType": resourceType,
			}

			switch resourceIdNum {
			case 0:
				request["ResourceId.1"] = d.Id()
			case 1:
				parts, err := helper.ParseResourceId(d.Id(), 2)
				if err != nil {
					return errs.WrapError(err)
				}
				request["ResourceId.1"] = parts[resourceIdNum]
			}

			count := 1
			for key, value := range added {
				request[fmt.Sprintf("Tag.%d.Key", count)] = key
				request[fmt.Sprintf("Tag.%d.Value", count)] = value
				count++
			}

			wait := helper.IncrementalWait(2*time.Second, 1*time.Second)
			err := retry.Retry(10*time.Minute, func() *retry.RetryError {
				response, err := client.RpcPost("Cbn", "2017-09-12", action, nil, request, false)
				if err != nil {
					if helper.NeedRetry(err) {
						wait()
						return retry.RetryableError(err)

					}
					return retry.NonRetryableError(err)
				}
				logs.AddDebug(action, response, request)
				return nil
			})
			if err != nil {
				return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
			}
		}
	}
	return nil
}

func (s *CbnService) DescribeCenRouteService(id string) (object cbn.RouteServiceEntry, err error) {
	parts, err := helper.ParseResourceId(id, 4)
	if err != nil {
		err = errs.WrapError(err)
		return
	}
	request := cbn.CreateDescribeRouteServicesInCenRequest()
	request.RegionId = s.client.RegionId
	request.AccessRegionId = parts[3]
	request.CenId = parts[0]
	request.Host = parts[2]
	request.HostRegionId = parts[1]

	var raw interface{}
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		raw, err = s.client.WithCbnClient(func(cbnClient *cbn.Client) (interface{}, error) {
			return cbnClient.DescribeRouteServicesInCen(request)
		})
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})

	if err != nil {
		err = errs.WrapErrorf(err, errs.DefaultErrorMsg, id, request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
		return
	}
	logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*cbn.DescribeRouteServicesInCenResponse)

	if len(response.RouteServiceEntries.RouteServiceEntry) < 1 {
		err = errs.WrapErrorf(errs.NotFoundErr("CenRouteService", id), errs.NotFoundMsg, errs.ProviderERROR, response.RequestId)
		return
	}
	return response.RouteServiceEntries.RouteServiceEntry[0], nil
}

func (s *CbnService) CenRouteServiceStateRefreshFunc(id string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeCenRouteService(id)
		if err != nil {
			if errs.NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", errs.WrapError(err)
		}

		for _, failState := range failStates {
			if object.Status == failState {
				return object, object.Status, errs.WrapError(errs.Error(errs.FailedToReachTargetStatus, object.Status))
			}
		}
		return object, object.Status, nil
	}
}

func (s *CbnService) DescribeCenTransitRouter(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	parts, err := helper.ParseResourceId(id, 2)
	if err != nil {
		return nil, errs.WrapError(err)
	}
	action := "ListTransitRouters"
	request := map[string]interface{}{
		"Region_id":       s.client.RegionId,
		"CenId":           parts[0],
		"TransitRouterId": parts[1],
	}
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("Cbn", "2017-09-12", action, nil, request, true)
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
		if errs.IsExpectedErrors(err, []string{"ParameterCenInstanceId"}) {
			return nil, errs.WrapErrorf(errs.NotFoundErr("CEN Instance ID", id), errs.NotFoundMsg, errs.ProviderERROR)
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.TransitRouters", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.TransitRouters", response)
	}
	if len(v.([]interface{})) < 1 {
		return object, errs.WrapErrorf(errs.NotFoundErr("CEN", id), errs.NotFoundWithResponse, response)
	} else {
		if v.([]interface{})[0].(map[string]interface{})["TransitRouterId"].(string) != parts[1] {
			return object, errs.WrapErrorf(errs.NotFoundErr("CEN", id), errs.NotFoundWithResponse, response)
		}
	}
	object = v.([]interface{})[0].(map[string]interface{})
	return object, nil
}

func (s *CbnService) CenTransitRouterStateRefreshFunc(id string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {

		object, err := s.DescribeCenTransitRouter(id)
		if err != nil {
			if errs.NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", errs.WrapError(err)
		}

		for _, failState := range failStates {
			if fmt.Sprint(object["Status"]) == failState {
				return object, fmt.Sprint(object["Status"]), errs.WrapError(errs.Error(errs.FailedToReachTargetStatus, fmt.Sprint(object["Status"])))
			}
		}
		return object, fmt.Sprint(object["Status"]), nil
	}
}

func (s *CbnService) DescribeCenTransitRouterPeerAttachment(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	parts, opErr := helper.ParseResourceId(id, 2)
	if opErr != nil {
		return nil, errs.WrapError(opErr)
	}
	action := "ListTransitRouterPeerAttachments"
	request := map[string]interface{}{
		"RegionId":                  s.client.RegionId,
		"CenId":                     parts[0],
		"TransitRouterAttachmentId": parts[1],
	}
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("Cbn", "2017-09-12", action, nil, request, true)
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
		if errs.IsExpectedErrors(err, []string{"ParameterCenInstanceId", "IllegalParam.Region"}) {
			return nil, errs.WrapErrorf(errs.NotFoundErr("CEN Instance ID", id), errs.NotFoundMsg, errs.ProviderERROR)
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.TransitRouterAttachments", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.TransitRouterAttachments", response)
	}
	if len(v.([]interface{})) < 1 {
		return object, errs.WrapErrorf(errs.NotFoundErr("CEN", id), errs.NotFoundWithResponse, response)
	} else {
		if v.([]interface{})[0].(map[string]interface{})["TransitRouterAttachmentId"].(string) != parts[1] {
			return object, errs.WrapErrorf(errs.NotFoundErr("CEN", id), errs.NotFoundWithResponse, response)
		}
	}
	object = v.([]interface{})[0].(map[string]interface{})
	return object, nil
}

func (s *CbnService) CenTransitRouterPeerAttachmentStateRefreshFunc(id string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeCenTransitRouterPeerAttachment(id)
		if err != nil {
			if errs.NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", errs.WrapError(err)
		}

		for _, failState := range failStates {
			if fmt.Sprint(object["Status"]) == failState {
				return object, fmt.Sprint(object["Status"]), errs.WrapError(errs.Error(errs.FailedToReachTargetStatus, fmt.Sprint(object["Status"])))
			}
		}
		return object, fmt.Sprint(object["Status"]), nil
	}
}

func (s *CbnService) DescribeCenTransitRouterVbrAttachment(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	action := "ListTransitRouterVbrAttachments"
	client := s.client

	parts, err := helper.ParseResourceId(id, 2)
	if err != nil {
		return nil, errs.WrapError(err)
	}

	request := map[string]interface{}{
		"RegionId":                  s.client.RegionId,
		"CenId":                     parts[0],
		"TransitRouterAttachmentId": parts[1],
		"MaxResults":                helper.PageSizeLarge,
	}

	idExist := false
	for {
		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(5*time.Minute, func() *retry.RetryError {
			response, err = client.RpcPost("Cbn", "2017-09-12", action, nil, request, true)
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
			if errs.IsExpectedErrors(err, []string{"IllegalParam.Region"}) {
				return object, errs.WrapErrorf(errs.NotFoundErr("Cen:TransitRouterVbrAttachment", id), errs.NotFoundWithResponse, response)
			}
			return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
		}

		resp, err := jsonpath.Get("$.TransitRouterAttachments", response)
		if err != nil {
			return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.TransitRouterAttachments", response)
		}

		if v, ok := resp.([]interface{}); !ok || len(v) < 1 {
			return object, errs.WrapErrorf(errs.NotFoundErr("Cen:TransitRouterVbrAttachment", id), errs.NotFoundWithResponse, response)
		}

		for _, v := range resp.([]interface{}) {
			if fmt.Sprint(v.(map[string]interface{})["CenId"]) == parts[0] && fmt.Sprint(v.(map[string]interface{})["TransitRouterAttachmentId"]) == parts[1] {
				idExist = true
				return v.(map[string]interface{}), nil
			}
		}

		if nextToken, ok := response["NextToken"].(string); ok && nextToken != "" {
			request["NextToken"] = nextToken
		} else {
			break
		}
	}

	if !idExist {
		return object, errs.WrapErrorf(errs.NotFoundErr("Cen:TransitRouterVbrAttachment", id), errs.NotFoundWithResponse, response)
	}

	return object, nil
}

func (s *CbnService) CenTransitRouterVbrAttachmentStateRefreshFunc(id string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeCenTransitRouterVbrAttachment(id)
		if err != nil {
			if errs.NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", errs.WrapError(err)
		}

		for _, failState := range failStates {
			if fmt.Sprint(object["Status"]) == failState {
				return object, fmt.Sprint(object["Status"]), errs.WrapError(errs.Error(errs.FailedToReachTargetStatus, fmt.Sprint(object["Status"])))
			}
		}

		return object, fmt.Sprint(object["Status"]), nil
	}
}

func (s *CbnService) DescribeCenTransitRouterVpcAttachment(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	action := "ListTransitRouterVpcAttachments"
	client := s.client

	parts, err := helper.ParseResourceId(id, 2)
	if err != nil {
		return nil, errs.WrapError(err)
	}

	request := map[string]interface{}{
		"RegionId":                  s.client.RegionId,
		"CenId":                     parts[0],
		"TransitRouterAttachmentId": parts[1],
		"MaxResults":                helper.PageSizeLarge,
	}

	idExist := false
	for {
		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(5*time.Minute, func() *retry.RetryError {
			response, err = client.RpcPost("Cbn", "2017-09-12", action, nil, request, true)
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
			if errs.IsExpectedErrors(err, []string{"IllegalParam.Region"}) {
				return object, errs.WrapErrorf(errs.NotFoundErr("Cen:TransitRouterVpcAttachment", id), errs.NotFoundWithResponse, response)
			}
			return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
		}

		resp, err := jsonpath.Get("$.TransitRouterAttachments", response)
		if err != nil {
			return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.TransitRouterAttachments", response)
		}

		if v, ok := resp.([]interface{}); !ok || len(v) < 1 {
			return object, errs.WrapErrorf(errs.NotFoundErr("Cen:TransitRouterVpcAttachment", id), errs.NotFoundWithResponse, response)
		}

		for _, v := range resp.([]interface{}) {
			if fmt.Sprint(v.(map[string]interface{})["CenId"]) == parts[0] && fmt.Sprint(v.(map[string]interface{})["TransitRouterAttachmentId"]) == parts[1] {
				idExist = true
				return v.(map[string]interface{}), nil
			}
		}

		if nextToken, ok := response["NextToken"].(string); ok && nextToken != "" {
			request["NextToken"] = nextToken
		} else {
			break
		}
	}

	if !idExist {
		return object, errs.WrapErrorf(errs.NotFoundErr("Cen:TransitRouterVpcAttachment", id), errs.NotFoundWithResponse, response)
	}

	return object, nil
}

func (s *CbnService) CenTransitRouterVpcAttachmentStateRefreshFunc(id string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeCenTransitRouterVpcAttachment(id)
		if err != nil {
			if errs.NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", errs.WrapError(err)
		}

		for _, failState := range failStates {
			if fmt.Sprint(object["Status"]) == failState {
				return object, fmt.Sprint(object["Status"]), errs.WrapError(errs.Error(errs.FailedToReachTargetStatus, fmt.Sprint(object["Status"])))
			}
		}

		return object, fmt.Sprint(object["Status"]), nil
	}
}

func (s *CbnService) DescribeCenTransitRouterRouteEntry(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	parts, opErr := helper.ParseResourceId(id, 2)
	if opErr != nil {
		return nil, errs.WrapError(opErr)
	}
	action := "ListTransitRouterRouteEntries"
	request := map[string]interface{}{
		"TransitRouterRouteTableId":  parts[0],
		"TransitRouterRouteEntryIds": []string{parts[1]},
	}
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("Cbn", "2017-09-12", action, nil, request, true)
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
		if errs.IsExpectedErrors(err, []string{"InvalidTransitRouterRouteTableId.NotFound"}) {
			return nil, errs.WrapErrorf(errs.NotFoundErr("CEN TransitRouter RouteTable ID", id), errs.NotFoundMsg, errs.ProviderERROR)
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.TransitRouterRouteEntries", response)
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"ParameterCenInstanceId"}) {
			return nil, errs.WrapErrorf(errs.NotFoundErr("CEN Instance ID", id), errs.NotFoundMsg, errs.ProviderERROR)
		}
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.TransitRouterRouteEntries", response)
	}
	if len(v.([]interface{})) < 1 {
		return object, errs.WrapErrorf(errs.NotFoundErr("CEN", id), errs.NotFoundWithResponse, response)
	} else {
		if v.([]interface{})[0].(map[string]interface{})["TransitRouterRouteEntryId"].(string) != parts[1] {
			return object, errs.WrapErrorf(errs.NotFoundErr("CEN", id), errs.NotFoundWithResponse, response)
		}
	}
	object = v.([]interface{})[0].(map[string]interface{})
	return object, nil
}

func (s *CbnService) CenTransitRouterRouteEntryStateRefreshFunc(id string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeCenTransitRouterRouteEntry(id)
		if err != nil {
			if errs.NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", errs.WrapError(err)
		}

		for _, failState := range failStates {
			if fmt.Sprint(object["TransitRouterRouteEntryStatus"]) == failState {
				return object, fmt.Sprint(object["TransitRouterRouteEntryStatus"]), errs.WrapError(errs.Error(errs.FailedToReachTargetStatus, fmt.Sprint(object["TransitRouterRouteEntryStatus"])))
			}
		}
		return object, fmt.Sprint(object["TransitRouterRouteEntryStatus"]), nil
	}
}

func (s *CbnService) DescribeCenTransitRouterRouteTable(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	parts, opErr := helper.ParseResourceId(id, 2)
	if opErr != nil {
		return nil, errs.WrapError(opErr)
	}
	action := "ListTransitRouterRouteTables"
	request := map[string]interface{}{
		"RegionId":                   s.client.RegionId,
		"TransitRouterId":            parts[0],
		"TransitRouterRouteTableIds": []string{parts[1]},
	}
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("Cbn", "2017-09-12", action, nil, request, true)
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
		if errs.IsExpectedErrors(err, []string{"InvalidTransitRouterId.NotFound"}) {
			return nil, errs.WrapErrorf(errs.NotFoundErr("CEN TransitRouter ID", id), errs.NotFoundMsg, errs.ProviderERROR)
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.TransitRouterRouteTables", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.TransitRouterRouteTables", response)
	}
	if len(v.([]interface{})) < 1 {
		return object, errs.WrapErrorf(errs.NotFoundErr("CEN", id), errs.NotFoundWithResponse, response)
	} else {
		if v.([]interface{})[0].(map[string]interface{})["TransitRouterRouteTableId"].(string) != parts[1] {
			return object, errs.WrapErrorf(errs.NotFoundErr("CEN", id), errs.NotFoundWithResponse, response)
		}
	}
	object = v.([]interface{})[0].(map[string]interface{})
	return object, nil
}

func (s *CbnService) CenTransitRouterRouteTableStateRefreshFunc(id string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeCenTransitRouterRouteTable(id)
		if err != nil {
			if errs.NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", errs.WrapError(err)
		}

		for _, failState := range failStates {
			if fmt.Sprint(object["TransitRouterRouteTableStatus"]) == failState {
				return object, fmt.Sprint(object["TransitRouterRouteTableStatus"]), errs.WrapError(errs.Error(errs.FailedToReachTargetStatus, fmt.Sprint(object["TransitRouterRouteTableStatus"])))
			}
		}
		return object, fmt.Sprint(object["TransitRouterRouteTableStatus"]), nil
	}
}

func (s *CbnService) DescribeCenTransitRouterRouteTableAssociation(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "ListTransitRouterRouteTableAssociations"
	parts, err := helper.ParseResourceId(id, 2)
	if err != nil {
		err = errs.WrapError(err)
		return
	}
	request := map[string]interface{}{
		"TransitRouterAttachmentId": parts[0],
		"TransitRouterRouteTableId": parts[1],
	}
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("Cbn", "2017-09-12", action, nil, request, true)
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
		if errs.IsExpectedErrors(err, []string{"ParameterCenInstanceId"}) {
			return nil, errs.WrapErrorf(errs.NotFoundErr("CEN Instance ID", id), errs.NotFoundMsg, errs.ProviderERROR)
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.TransitRouterAssociations", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.TransitRouterAssociations", response)
	}
	if len(v.([]interface{})) < 1 {
		return object, errs.WrapErrorf(errs.NotFoundErr("CEN", id), errs.NotFoundWithResponse, response)
	} else {
		if v.([]interface{})[0].(map[string]interface{})["TransitRouterAttachmentId"].(string) != parts[0] {
			return object, errs.WrapErrorf(errs.NotFoundErr("CEN", id), errs.NotFoundWithResponse, response)
		}
	}
	object = v.([]interface{})[0].(map[string]interface{})
	return object, nil
}

func (s *CbnService) CenTransitRouterRouteTableAssociationStateRefreshFunc(id string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeCenTransitRouterRouteTableAssociation(id)
		if err != nil {
			if errs.NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", errs.WrapError(err)
		}

		for _, failState := range failStates {
			if fmt.Sprint(object["Status"]) == failState {
				return object, fmt.Sprint(object["Status"]), errs.WrapError(errs.Error(errs.FailedToReachTargetStatus, fmt.Sprint(object["Status"])))
			}
		}
		return object, fmt.Sprint(object["Status"]), nil
	}
}

func (s *CbnService) DescribeCenTransitRouterRouteTablePropagation(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "ListTransitRouterRouteTablePropagations"
	parts, err := helper.ParseResourceId(id, 2)
	if err != nil {
		err = errs.WrapError(err)
		return
	}
	request := map[string]interface{}{
		"TransitRouterAttachmentId": parts[0],
		"TransitRouterRouteTableId": parts[1],
	}
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("Cbn", "2017-09-12", action, nil, request, true)
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
		if errs.IsExpectedErrors(err, []string{"ParameterCenInstanceId"}) {
			return nil, errs.WrapErrorf(errs.NotFoundErr("CEN Instance ID", id), errs.NotFoundMsg, errs.ProviderERROR)
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.TransitRouterPropagations", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.TransitRouterPropagations", response)
	}
	if len(v.([]interface{})) < 1 {
		return object, errs.WrapErrorf(errs.NotFoundErr("CEN", id), errs.NotFoundWithResponse, response)
	} else {
		if v.([]interface{})[0].(map[string]interface{})["TransitRouterAttachmentId"].(string) != parts[0] {
			return object, errs.WrapErrorf(errs.NotFoundErr("CEN", id), errs.NotFoundWithResponse, response)
		}
	}
	object = v.([]interface{})[0].(map[string]interface{})
	return object, nil
}

func (s *CbnService) CenTransitRouterRouteTablePropagationStateRefreshFunc(id string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeCenTransitRouterRouteTablePropagation(id)
		if err != nil {
			if errs.NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", errs.WrapError(err)
		}

		for _, failState := range failStates {
			if fmt.Sprint(object["Status"]) == failState {
				return object, fmt.Sprint(object["Status"]), errs.WrapError(errs.Error(errs.FailedToReachTargetStatus, fmt.Sprint(object["Status"])))
			}
		}
		return object, fmt.Sprint(object["Status"]), nil
	}
}

func (s *CbnService) DescribeCenTrafficMarkingPolicy(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	parts, err := helper.ParseResourceId(id, 2)
	if err != nil {
		err = errs.WrapError(err)
		return
	}
	action := "ListTrafficMarkingPolicies"
	request := map[string]interface{}{
		"TransitRouterId":        parts[0],
		"TrafficMarkingPolicyId": parts[1],
	}
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("Cbn", "2017-09-12", action, nil, request, true)
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
	v, err := jsonpath.Get("$.TrafficMarkingPolicies", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.TrafficMarkingPolicies", response)
	}
	if len(v.([]interface{})) < 1 {
		return object, errs.WrapErrorf(errs.NotFoundErr("CEN", id), errs.NotFoundWithResponse, response)
	} else {
		if v.([]interface{})[0].(map[string]interface{})["TrafficMarkingPolicyId"].(string) != parts[1] {
			return object, errs.WrapErrorf(errs.NotFoundErr("CEN", id), errs.NotFoundWithResponse, response)
		}
	}
	object = v.([]interface{})[0].(map[string]interface{})
	return object, nil
}

func (s *CbnService) CenTrafficMarkingPolicyStateRefreshFunc(id string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeCenTrafficMarkingPolicy(id)
		if err != nil {
			if errs.NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", errs.WrapError(err)
		}

		for _, failState := range failStates {
			if fmt.Sprint(object["TrafficMarkingPolicyStatus"]) == failState {
				return object, fmt.Sprint(object["TrafficMarkingPolicyStatus"]), errs.WrapError(errs.Error(errs.FailedToReachTargetStatus, fmt.Sprint(object["TrafficMarkingPolicyStatus"])))
			}
		}
		return object, fmt.Sprint(object["TrafficMarkingPolicyStatus"]), nil
	}
}

func (s *CbnService) DescribeCenTransitRouterVpnAttachment(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "ListTransitRouterVpnAttachments"
	request := map[string]interface{}{
		"RegionId":                  s.client.RegionId,
		"TransitRouterAttachmentId": id,
	}
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("Cbn", "2017-09-12", action, nil, request, true)
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
		if errs.IsExpectedErrors(err, []string{"IllegalParam.Region"}) {
			return nil, errs.WrapErrorf(errs.NotFoundErr("CEN Instance ID", id), errs.NotFoundMsg, errs.ProviderERROR)
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.TransitRouterAttachments", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.TransitRouterAttachments", response)
	}
	if len(v.([]interface{})) < 1 {
		return object, errs.WrapErrorf(errs.NotFoundErr("CEN", id), errs.NotFoundWithResponse, response)
	}
	object = v.([]interface{})[0].(map[string]interface{})
	return object, nil
}

func (s *CbnService) CenTransitRouterVpnAttachmentStateRefreshFunc(id string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeCenTransitRouterVpnAttachment(id)
		if err != nil {
			if errs.NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", errs.WrapError(err)
		}

		for _, failState := range failStates {
			if fmt.Sprint(object["Status"]) == failState {
				return object, fmt.Sprint(object["Status"]), errs.WrapError(errs.Error(errs.FailedToReachTargetStatus, fmt.Sprint(object["Status"])))
			}
		}
		return object, fmt.Sprint(object["Status"]), nil
	}
}

func (s *CbnService) DescribeCenTransitRouterGrantAttachment(id string) (object map[string]interface{}, err error) {
	client := s.client
	if err != nil {
		return object, errs.WrapError(err)
	}
	parts, err := helper.ParseResourceId(id, 4)
	if err != nil {
		return object, errs.WrapError(err)
	}

	request := map[string]interface{}{
		"ResourceId":  parts[1],
		"ProductType": parts[0],
		"RegionId":    s.client.RegionId,
	}

	var response map[string]interface{}
	action := "DescribeGrantRulesToResource"
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		resp, err := client.RpcPost("Cbn", "2017-09-12", action, nil, request, true)
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
	v, err := jsonpath.Get("$.GrantRules", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.GrantRules", response)
	}
	if len(v.([]interface{})) < 1 {
		return object, errs.WrapErrorf(errs.NotFoundErr("TransitRouterGrantAttachment", id), errs.NotFoundWithResponse, response)
	} else {
		for _, vv := range v.([]interface{}) {
			item := vv.(map[string]interface{})
			if fmt.Sprint(item["CenId"]) == parts[3] && fmt.Sprint(item["CenOwnerId"]) == parts[2] {
				return item, nil
			}
		}
	}
	return object, errs.WrapErrorf(errs.NotFoundErr("TransitRouterGrantAttachment", id), errs.NotFoundWithResponse, response)
}

func (s *CbnService) DescribeCenTransitRouterPrefixListAssociation(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	action := "ListTransitRouterPrefixListAssociation"
	client := s.client

	parts, err := helper.ParseResourceId(id, 4)
	if err != nil {
		return nil, errs.WrapError(err)
	}
	request := map[string]interface{}{
		"RegionId":             s.client.RegionId,
		"PrefixListId":         parts[0],
		"TransitRouterId":      parts[1],
		"TransitRouterTableId": parts[2],
	}

	idExist := false
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("Cbn", "2017-09-12", action, nil, request, true)
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

	resp, err := jsonpath.Get("$.PrefixLists", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.PrefixLists", response)
	}

	if v, ok := resp.([]interface{}); !ok || len(v) < 1 {
		return object, errs.WrapErrorf(errs.NotFoundErr("CEN:TransitRouterPrefixListAssociation", id), errs.NotFoundWithResponse, response)
	}

	for _, v := range resp.([]interface{}) {
		if fmt.Sprint(v.(map[string]interface{})["NextHop"]) == parts[3] {
			idExist = true
			return v.(map[string]interface{}), nil
		}
	}

	if !idExist {
		return object, errs.WrapErrorf(errs.NotFoundErr("CEN:TransitRouterPrefixListAssociation", id), errs.NotFoundWithResponse, response)
	}

	return object, nil
}

func (s *CbnService) CenTransitRouterPrefixListAssociationStateRefreshFunc(d *schema.ResourceData, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeCenTransitRouterPrefixListAssociation(d.Id())
		if err != nil {
			if errs.NotFoundError(err) {
				return nil, "", nil
			}
			return nil, "", errs.WrapError(err)
		}

		for _, failState := range failStates {
			if fmt.Sprint(object["Status"]) == failState {
				return object, fmt.Sprint(object["Status"]), errs.WrapError(errs.Error(errs.FailedToReachTargetStatus, fmt.Sprint(object["Status"])))
			}
		}
		return object, fmt.Sprint(object["Status"]), nil
	}
}

func (s *CbnService) DescribeCenTransitRouterCidr(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	action := "ListTransitRouterCidr"
	client := s.client

	parts, err := helper.ParseResourceId(id, 2)
	if err != nil {
		return nil, errs.WrapError(err)
	}

	request := map[string]interface{}{
		"RegionId":            s.client.RegionId,
		"ClientToken":         helper.BuildClientToken("ListTransitRouterCidr"),
		"TransitRouterId":     parts[0],
		"TransitRouterCidrId": parts[1],
	}

	idExist := false
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("Cbn", "2017-09-12", action, nil, request, true)
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

	resp, err := jsonpath.Get("$.CidrLists", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.CidrLists", response)
	}

	if v, ok := resp.([]interface{}); !ok || len(v) < 1 {
		return object, errs.WrapErrorf(errs.NotFoundErr("CEN:TransitRouterCidr", id), errs.NotFoundWithResponse, response)
	}

	for _, v := range resp.([]interface{}) {
		if fmt.Sprint(v.(map[string]interface{})["TransitRouterCidrId"]) == parts[1] {
			idExist = true
			return v.(map[string]interface{}), nil
		}
	}

	if !idExist {
		return object, errs.WrapErrorf(errs.NotFoundErr("CEN:TransitRouterCidr", id), errs.NotFoundWithResponse, response)
	}

	return object, nil
}

func (s *CbnService) DescribeCenTransitRouterMulticastDomainSource(id string) (object map[string]interface{}, err error) {
	client := s.client
	if err != nil {
		return object, errs.WrapError(err)
	}
	parts, err := helper.ParseResourceId(id, 3)
	if err != nil {
		return object, errs.WrapError(err)
	}

	request := map[string]interface{}{
		"TransitRouterMulticastDomainId": parts[0],
		"GroupIpAddress":                 parts[1],
		"NetworkInterfaceIds.1":          parts[2],
	}

	var response map[string]interface{}
	action := "ListTransitRouterMulticastGroups"
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		resp, err := client.RpcPost("Cbn", "2017-09-12", action, nil, request, true)
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
		if errs.IsExpectedErrors(err, []string{"IllegalParam.TransitRouterMulticastDomainId"}) {
			return object, errs.WrapErrorf(errs.NotFoundErr("CEN:TransitRouterMulticastGroups", id), errs.NotFoundWithResponse, response)
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.TransitRouterMulticastGroups", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.TransitRouterMulticastGroups", response)
	}
	resp := v.([]interface{})
	if len(resp) < 1 {
		return object, errs.WrapErrorf(errs.NotFoundErr("CEN:TransitRouterMulticastGroups", id), errs.NotFoundWithResponse, response)
	}
	return resp[0].(map[string]interface{}), nil
}

func (s *CbnService) CenTransitRouterMulticastDomainSourceStateRefreshFunc(d *schema.ResourceData, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeCenTransitRouterMulticastDomainSource(d.Id())
		if err != nil {
			if errs.NotFoundError(err) {
				return nil, "", nil
			}
			return nil, "", errs.WrapError(err)
		}
		for _, failState := range failStates {
			if fmt.Sprint(object["Status"]) == failState {
				return object, fmt.Sprint(object["Status"]), errs.WrapError(errs.Error(errs.FailedToReachTargetStatus, fmt.Sprint(object["Status"])))
			}
		}
		return object, fmt.Sprint(object["Status"]), nil
	}
}

func (s *CbnService) DescribeCenTransitRouterMulticastDomain(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	action := "ListTransitRouterMulticastDomains"
	client := s.client

	request := map[string]interface{}{
		"RegionId":                       s.client.RegionId,
		"ClientToken":                    helper.BuildClientToken("ListTransitRouterMulticastDomains"),
		"TransitRouterMulticastDomainId": id,
	}

	idExist := false
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("Cbn", "2017-09-12", action, nil, request, true)
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

	resp, err := jsonpath.Get("$.TransitRouterMulticastDomains", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.TransitRouterMulticastDomains", response)
	}

	if v, ok := resp.([]interface{}); !ok || len(v) < 1 {
		return object, errs.WrapErrorf(errs.NotFoundErr("CEN:TransitRouterMulticastDomain", id), errs.NotFoundWithResponse, response)
	}

	for _, v := range resp.([]interface{}) {
		if fmt.Sprint(v.(map[string]interface{})["TransitRouterMulticastDomainId"]) == id {
			idExist = true
			return v.(map[string]interface{}), nil
		}
	}

	if !idExist {
		return object, errs.WrapErrorf(errs.NotFoundErr("CEN:TransitRouterMulticastDomain", id), errs.NotFoundWithResponse, response)
	}

	return object, nil
}

func (s *CbnService) CenTransitRouterMulticastDomainStateRefreshFunc(id string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeCenTransitRouterMulticastDomain(id)
		if err != nil {
			if errs.NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", errs.WrapError(err)
		}

		for _, failState := range failStates {
			if fmt.Sprint(object["Status"]) == failState {
				return object, fmt.Sprint(object["Status"]), errs.WrapError(errs.Error(errs.FailedToReachTargetStatus, fmt.Sprint(object["Status"])))
			}
		}
		return object, fmt.Sprint(object["Status"]), nil
	}
}

func (s *CbnService) DescribeCenInterRegionTrafficQosQueue(id string) (object map[string]interface{}, err error) {
	client := s.client
	if err != nil {
		return object, errs.WrapError(err)
	}

	request := map[string]interface{}{
		"TrafficQosQueueId": id,
	}

	var response map[string]interface{}
	action := "ListCenInterRegionTrafficQosQueues"
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		resp, err := client.RpcPost("Cbn", "2017-09-12", action, nil, request, true)
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
	v, err := jsonpath.Get("$.TrafficQosQueues", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.TrafficQosQueues", response)
	}
	resp := v.([]interface{})
	if len(resp) < 1 {
		return object, errs.WrapErrorf(errs.NotFoundErr("CEN:InterRegionTrafficQosQueue", id), errs.NotFoundWithResponse, response)
	}
	return resp[0].(map[string]interface{}), nil
}

func (s *CbnService) CenInterRegionTrafficQosQueueStateRefreshFunc(d *schema.ResourceData, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeCenInterRegionTrafficQosQueue(d.Id())
		if err != nil {
			if errs.NotFoundError(err) {
				return nil, "", nil
			}
			return nil, "", errs.WrapError(err)
		}
		for _, failState := range failStates {
			if fmt.Sprint(object["Status"]) == failState {
				return object, fmt.Sprint(object["Status"]), errs.WrapError(errs.Error(errs.FailedToReachTargetStatus, fmt.Sprint(object["Status"])))
			}
		}
		return object, fmt.Sprint(object["Status"]), nil
	}
}

func (s *CbnService) DescribeCenInterRegionTrafficQosPolicy(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	action := "ListCenInterRegionTrafficQosPolicies"
	client := s.client

	request := map[string]interface{}{
		"TrafficQosPolicyId": id,
	}

	idExist := false
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("Cbn", "2017-09-12", action, nil, request, true)
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

	resp, err := jsonpath.Get("$.TrafficQosPolicies", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.TrafficQosPolicies", response)
	}

	if v, ok := resp.([]interface{}); !ok || len(v) < 1 {
		return object, errs.WrapErrorf(errs.NotFoundErr("CEN:InterRegionTrafficQosPolicy", id), errs.NotFoundWithResponse, response)
	}

	for _, v := range resp.([]interface{}) {
		if fmt.Sprint(v.(map[string]interface{})["TrafficQosPolicyId"]) == id {
			idExist = true
			return v.(map[string]interface{}), nil
		}
	}

	if !idExist {
		return object, errs.WrapErrorf(errs.NotFoundErr("CEN:InterRegionTrafficQosPolicy", id), errs.NotFoundWithResponse, response)
	}

	return object, nil
}

func (s *CbnService) CenInterRegionTrafficQosPolicyStateRefreshFunc(id string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeCenInterRegionTrafficQosPolicy(id)
		if err != nil {
			if errs.NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", errs.WrapError(err)
		}

		for _, failState := range failStates {
			if fmt.Sprint(object["TrafficQosPolicyStatus"]) == failState {
				return object, fmt.Sprint(object["TrafficQosPolicyStatus"]), errs.WrapError(errs.Error(errs.FailedToReachTargetStatus, fmt.Sprint(object["TrafficQosPolicyStatus"])))
			}
		}

		return object, fmt.Sprint(object["TrafficQosPolicyStatus"]), nil
	}
}

func (s *CbnService) DescribeCenTransitRouterMulticastDomainPeerMember(id string) (object map[string]interface{}, err error) {
	client := s.client
	if err != nil {
		return object, errs.WrapError(err)
	}
	parts, err := helper.ParseResourceId(id, 3)
	if err != nil {
		return object, errs.WrapError(err)
	}

	request := map[string]interface{}{
		"TransitRouterMulticastDomainId":      parts[0],
		"GroupIpAddress":                      parts[1],
		"PeerTransitRouterMulticastDomains.1": parts[2],
	}

	var response map[string]interface{}
	action := "ListTransitRouterMulticastGroups"
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		resp, err := client.RpcPost("Cbn", "2017-09-12", action, nil, request, true)
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
	count, err := jsonpath.Get("$.TotalCount", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.TotalCount", response)
	}
	totalCount, err := convert.FormatInt(count)
	if err != nil {
		return nil, err
	}
	if totalCount == 0 {
		return object, errs.WrapErrorf(errs.NotFoundErr("CEN.TransitRouterMulticastDomainPeerMember", id), errs.NotFoundWithResponse, response)
	}
	v, err := jsonpath.Get("$.TransitRouterMulticastGroups[0]", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.TransitRouterMulticastGroups[0]", response)
	}
	return v.(map[string]interface{}), nil
}

func (s *CbnService) CenTransitRouterMulticastDomainPeerMemberStateRefreshFunc(d *schema.ResourceData, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeCenTransitRouterMulticastDomainPeerMember(d.Id())
		if err != nil {
			if errs.NotFoundError(err) {
				return nil, "", nil
			}
			return nil, "", errs.WrapError(err)
		}
		for _, failState := range failStates {
			if fmt.Sprint(object["Status"]) == failState {
				return object, fmt.Sprint(object["Status"]), errs.WrapError(errs.Error(errs.FailedToReachTargetStatus, fmt.Sprint(object["Status"])))
			}
		}
		return object, fmt.Sprint(object["Status"]), nil
	}
}

func (s *CbnService) DescribeCenTransitRouterMulticastDomainMember(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	action := "ListTransitRouterMulticastGroups"

	client := s.client
	if err != nil {
		return object, errs.WrapError(err)
	}

	parts, err := helper.ParseResourceId(id, 3)
	if err != nil {
		return object, errs.WrapError(err)
	}

	request := map[string]interface{}{
		"ClientToken":                    helper.BuildClientToken("ListTransitRouterMulticastGroups"),
		"TransitRouterMulticastDomainId": parts[0],
		"GroupIpAddress":                 parts[1],
		"NetworkInterfaceIds":            []string{parts[2]},
		"MaxResults":                     helper.PageSizeLarge,
	}

	idExist := false
	for {
		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(5*time.Minute, func() *retry.RetryError {
			response, err = client.RpcPost("Cbn", "2017-09-12", action, nil, request, true)
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
			if errs.IsExpectedErrors(err, []string{"IllegalParam.TransitRouterMulticastDomainId"}) {
				return object, errs.WrapErrorf(errs.NotFoundErr("Cen:TransitRouterMulticastDomainMember", id), errs.NotFoundWithResponse, response)
			}
			return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
		}

		totalCount, err := jsonpath.Get("$.TotalCount", response)
		if err != nil {
			return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.TotalCount", response)
		}

		if fmt.Sprint(totalCount) == "0" {
			return object, errs.WrapErrorf(errs.NotFoundErr("Cen:TransitRouterMulticastDomainMember", id), errs.NotFoundWithResponse, response)
		}

		resp, err := jsonpath.Get("$.TransitRouterMulticastGroups", response)
		if err != nil {
			return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.TransitRouterMulticastGroups", response)
		}

		if v, ok := resp.([]interface{}); !ok || len(v) < 1 {
			return object, errs.WrapErrorf(errs.NotFoundErr("Cen:TransitRouterMulticastDomainMember", id), errs.NotFoundWithResponse, response)
		}

		for _, v := range resp.([]interface{}) {
			if fmt.Sprint(v.(map[string]interface{})["TransitRouterMulticastDomainId"]) == parts[0] && fmt.Sprint(v.(map[string]interface{})["GroupIpAddress"]) == parts[1] && fmt.Sprint(v.(map[string]interface{})["NetworkInterfaceId"]) == parts[2] {
				idExist = true
				return v.(map[string]interface{}), nil
			}
		}

		if nextToken, ok := response["NextToken"].(string); ok && nextToken != "" {
			request["NextToken"] = nextToken
		} else {
			break
		}
	}

	if !idExist {
		return object, errs.WrapErrorf(errs.NotFoundErr("Cen:TransitRouterMulticastDomainMember", id), errs.NotFoundWithResponse, response)
	}

	return object, nil
}

func (s *CbnService) CenTransitRouterMulticastDomainMemberStateRefreshFunc(d *schema.ResourceData, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeCenTransitRouterMulticastDomainMember(d.Id())
		if err != nil {
			if errs.NotFoundError(err) {
				return nil, "", nil
			}
			return nil, "", errs.WrapError(err)
		}

		for _, failState := range failStates {
			if fmt.Sprint(object["Status"]) == failState {
				return object, fmt.Sprint(object["Status"]), errs.WrapError(errs.Error(errs.FailedToReachTargetStatus, fmt.Sprint(object["Status"])))
			}
		}

		return object, fmt.Sprint(object["Status"]), nil
	}
}

func (s *CbnService) DescribeCenChildInstanceRouteEntryToAttachment(id string) (object map[string]interface{}, err error) {
	client := s.client
	if err != nil {
		return object, errs.WrapError(err)
	}
	parts, err := helper.ParseResourceId(id, 4)
	if err != nil {
		return object, errs.WrapError(err)
	}

	request := map[string]interface{}{
		"CenId":                     parts[0],
		"ChildInstanceRouteTableId": parts[1],
		"TransitRouterAttachmentId": parts[2],
		"RouteFilter.1.Key":         "prefix-exact-match",
		"RouteFilter.1.Value.1":     parts[3],
	}

	var response map[string]interface{}
	action := "ListCenChildInstanceRouteEntriesToAttachment"
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		resp, err := client.RpcPost("Cbn", "2017-09-12", action, nil, request, true)
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
		if errs.IsExpectedErrors(err, []string{"ResourceNotFound.TransitRouterAttachment"}) {
			err = errs.WrapErrorf(errs.NotFoundErr("Cen", id), errs.NotFoundMsg, errs.ProviderERROR)
			return object, err
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}
	all, err := jsonpath.Get("$.RouteEntry", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.RouteEntry", response)
	}
	if len(all.([]interface{})) == 0 {
		return object, errs.WrapErrorf(errs.NotFoundErr("Cen", id), errs.NotFoundWithResponse, response)
	}
	v, err := jsonpath.Get("$.RouteEntry[0]", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.RouteEntry[0]", response)
	}
	return v.(map[string]interface{}), nil
}

func (s *CbnService) CenChildInstanceRouteEntryToAttachmentStateRefreshFunc(d *schema.ResourceData, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeCenChildInstanceRouteEntryToAttachment(d.Id())
		if err != nil {
			if errs.NotFoundError(err) {
				return nil, "", nil
			}
			return nil, "", errs.WrapError(err)
		}
		for _, failState := range failStates {
			if fmt.Sprint(object["Status"]) == failState {
				return object, fmt.Sprint(object["Status"]), errs.WrapError(errs.Error(errs.FailedToReachTargetStatus, fmt.Sprint(object["Status"])))
			}
		}
		return object, fmt.Sprint(object["Status"]), nil
	}
}

func (s *CbnService) DescribeCenTransitRouterMulticastDomainAssociation(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	action := "ListTransitRouterMulticastDomainAssociations"
	client := s.client

	parts, err := helper.ParseResourceId(id, 3)
	if err != nil {
		return nil, errs.WrapError(err)
	}

	request := map[string]interface{}{
		"ClientToken":               helper.BuildClientToken("ListTransitRouterMulticastDomainAssociations"),
		"TransitRouterAttachmentId": parts[1],
		"VSwitchIds":                []string{parts[2]},
	}

	idExist := false
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("Cbn", "2017-09-12", action, nil, request, true)
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

	resp, err := jsonpath.Get("$.TransitRouterMulticastAssociations", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.TransitRouterMulticastAssociations", response)
	}

	if v, ok := resp.([]interface{}); !ok || len(v) < 1 {
		return object, errs.WrapErrorf(errs.NotFoundErr("CEN:TransitRouterMulticastDomainAssociation", id), errs.NotFoundWithResponse, response)
	}

	for _, v := range resp.([]interface{}) {
		if fmt.Sprint(v.(map[string]interface{})["VSwitchId"]) == parts[2] {
			idExist = true
			return v.(map[string]interface{}), nil
		}
	}

	if !idExist {
		return object, errs.WrapErrorf(errs.NotFoundErr("CEN:TransitRouterMulticastDomainAssociation", id), errs.NotFoundWithResponse, response)
	}

	return object, nil
}

func (s *CbnService) CenTransitRouterMulticastDomainAssociationStateRefreshFunc(id string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeCenTransitRouterMulticastDomainAssociation(id)
		if err != nil {
			if errs.NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", errs.WrapError(err)
		}

		for _, failState := range failStates {
			if fmt.Sprint(object["Status"]) == failState {
				return object, fmt.Sprint(object["Status"]), errs.WrapError(errs.Error(errs.FailedToReachTargetStatus, fmt.Sprint(object["Status"])))
			}
		}

		return object, fmt.Sprint(object["Status"]), nil
	}
}

func (s *CbnService) DescribeCenTransitRouteTableAggregation(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	action := "DescribeTransitRouteTableAggregation"

	client := s.client
	if err != nil {
		return object, errs.WrapError(err)
	}

	parts, _ := helper.ParseResourceId(id, 2)
	if strings.Contains(id, "#") {
		parts = strings.Split(id, "#")
	}

	request := map[string]interface{}{
		"ClientToken":                      helper.BuildClientToken("DescribeTransitRouteTableAggregation"),
		"TransitRouteTableId":              parts[0],
		"TransitRouteTableAggregationCidr": parts[1],
	}

	idExist := false
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("Cbn", "2017-09-12", action, nil, request, true)
		if err != nil {
			if errs.IsExpectedErrors(err, []string{"Operation.Blocking"}) || helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)

	if err != nil {
		if errs.IsExpectedErrors(err, []string{"InstanceNotExist.AggregationRoute", "InstanceNotExist.TransitRouterRouteTable"}) {
			return nil, errs.WrapErrorf(errs.NotFoundErr("CEN:TransitRouteTableAggregation", id), errs.NotFoundWithResponse, response)
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}

	resp, err := jsonpath.Get("$.Data", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.Data", response)
	}

	if v, ok := resp.([]interface{}); !ok || len(v) < 1 {
		return object, errs.WrapErrorf(errs.NotFoundErr("CEN:TransitRouteTableAggregation", id), errs.NotFoundWithResponse, response)
	}

	for _, v := range resp.([]interface{}) {
		if fmt.Sprint(v.(map[string]interface{})["TransitRouteTableAggregationCidr"]) == parts[1] {
			idExist = true
			return v.(map[string]interface{}), nil
		}
	}

	if !idExist {
		return object, errs.WrapErrorf(errs.NotFoundErr("CEN:TransitRouteTableAggregation", id), errs.NotFoundWithResponse, response)
	}

	return object, nil
}

func (s *CbnService) CenTransitRouteTableAggregationStateRefreshFunc(d *schema.ResourceData, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeCenTransitRouteTableAggregation(d.Id())
		if err != nil {
			if errs.NotFoundError(err) {
				return nil, "", nil
			}
			return nil, "", errs.WrapError(err)
		}

		for _, failState := range failStates {
			if fmt.Sprint(object["Status"]) == failState {
				return object, fmt.Sprint(object["Status"]), errs.WrapError(errs.Error(errs.FailedToReachTargetStatus, fmt.Sprint(object["Status"])))
			}
		}
		return object, fmt.Sprint(object["Status"]), nil
	}
}

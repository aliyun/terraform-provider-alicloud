package alicloud

import (
	"fmt"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/cbn"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

type CbnService struct {
	client *connectivity.AliyunClient
}

func (s *CbnService) DescribeCenFlowlog(id string) (object cbn.FlowLog, err error) {
	request := cbn.CreateDescribeFlowlogsRequest()
	request.RegionId = s.client.RegionId

	request.FlowLogId = id
	var raw interface{}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err = s.client.WithCbnClient(func(cbnClient *cbn.Client) (interface{}, error) {
			return cbnClient.DescribeFlowlogs(request)
		})
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})

	if err != nil {
		err = WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
		return
	}
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

func (s *CbnService) DescribeCenInstance(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	action := "DescribeCens"

	conn, err := s.client.NewCbnClient()
	if err != nil {
		return nil, WrapError(err)
	}

	filterMaps := make([]map[string]interface{}, 0)
	filterMaps = append(filterMaps, map[string]interface{}{
		"Key":   "CenId",
		"Value": []string{id},
	})

	request := map[string]interface{}{
		"Filter":     filterMaps,
		"PageSize":   PageSizeLarge,
		"PageNumber": 1,
	}

	idExist := false
	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-09-12"), StringPointer("AK"), nil, request, &runtime)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, request)

		if err != nil {
			return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
		}

		resp, err := jsonpath.Get("$.Cens.Cen", response)
		if err != nil {
			return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Cens.Cen", response)
		}

		if v, ok := resp.([]interface{}); !ok || len(v) < 1 {
			return object, WrapErrorf(Error(GetNotFoundMessage("Cen:Instance", id)), NotFoundWithResponse, response)
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
		return object, WrapErrorf(Error(GetNotFoundMessage("Cen:Instance", id)), NotFoundWithResponse, response)
	}

	return object, nil
}

func (s *CbnService) CenInstanceStateRefreshFunc(d *schema.ResourceData, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeCenInstance(d.Id())
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
		addDebug(request.GetActionName(), raw)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
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
		addDebug(request.GetActionName(), raw)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
	}
	return nil
}

func (s *CbnService) DescribeCenRouteMap(id string) (object cbn.RouteMap, err error) {
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		err = WrapError(err)
		return
	}
	request := cbn.CreateDescribeCenRouteMapsRequest()
	request.RegionId = s.client.RegionId
	request.CenId = parts[0]
	request.RouteMapId = parts[1]

	var raw interface{}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err = s.client.WithCbnClient(func(cbnClient *cbn.Client) (interface{}, error) {
			return cbnClient.DescribeCenRouteMaps(request)
		})
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})

	if err != nil {
		err = WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
		return
	}
	response, _ := raw.(*cbn.DescribeCenRouteMapsResponse)

	if len(response.RouteMaps.RouteMap) < 1 {
		err = WrapErrorf(Error(GetNotFoundMessage("CenRouteMap", id)), NotFoundMsg, ProviderERROR)
		return
	}
	return response.RouteMaps.RouteMap[0], nil
}

func (s *CbnService) CenRouteMapStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeCenRouteMap(id)
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

func (s *CbnService) DescribeCenPrivateZone(id string) (object cbn.PrivateZoneInfo, err error) {
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		err = WrapError(err)
		return
	}
	request := cbn.CreateDescribeCenPrivateZoneRoutesRequest()
	request.CenId = parts[0]
	request.AccessRegionId = parts[1]

	var raw interface{}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err = s.client.WithCbnClient(func(cbnClient *cbn.Client) (interface{}, error) {
			return cbnClient.DescribeCenPrivateZoneRoutes(request)
		})
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)

	if err != nil {
		if IsExpectedErrors(err, []string{"ParameterIllegal.CenInstanceId"}) {
			return object, WrapErrorf(Error(GetNotFoundMessage("Cen:PrivateZone", id)), NotFoundMsg, ProviderERROR)
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	response, _ := raw.(*cbn.DescribeCenPrivateZoneRoutesResponse)

	if len(response.PrivateZoneInfos.PrivateZoneInfo) < 1 {
		err = WrapErrorf(Error(GetNotFoundMessage("Cen:PrivateZone", id)), NotFoundMsg, ProviderERROR)
		return
	}

	return response.PrivateZoneInfos.PrivateZoneInfo[0], nil
}

func (s *CbnService) CenPrivateZoneStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeCenPrivateZone(id)
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

func (s *CbnService) DescribeCenVbrHealthCheck(id string) (object cbn.VbrHealthCheck, err error) {
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		err = WrapError(err)
		return
	}
	request := cbn.CreateDescribeCenVbrHealthCheckRequest()
	request.RegionId = s.client.RegionId
	request.VbrInstanceId = parts[0]
	request.VbrInstanceRegionId = parts[1]

	var raw interface{}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err = s.client.WithCbnClient(func(cbnClient *cbn.Client) (interface{}, error) {
			return cbnClient.DescribeCenVbrHealthCheck(request)
		})
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})

	if err != nil {
		err = WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
		return
	}
	response, _ := raw.(*cbn.DescribeCenVbrHealthCheckResponse)

	if len(response.VbrHealthChecks.VbrHealthCheck) < 1 {
		err = WrapErrorf(Error(GetNotFoundMessage("CenVbrHealthCheck", id)), NotFoundMsg, ProviderERROR, response.RequestId)
		return
	}
	return response.VbrHealthChecks.VbrHealthCheck[0], nil
}

func (s *CbnService) DescribeCenInstanceAttachment(id string) (object cbn.DescribeCenAttachedChildInstanceAttributeResponse, err error) {
	parts, err := ParseResourceId(id, 4)
	if err != nil {
		err = WrapError(err)
		return
	}
	request := cbn.CreateDescribeCenAttachedChildInstanceAttributeRequest()
	request.RegionId = s.client.RegionId
	request.ChildInstanceId = parts[1]
	request.ChildInstanceRegionId = parts[3]
	request.ChildInstanceType = parts[2]
	request.CenId = parts[0]

	var raw interface{}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err = s.client.WithCbnClient(func(cbnClient *cbn.Client) (interface{}, error) {
			return cbnClient.DescribeCenAttachedChildInstanceAttribute(request)
		})
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})

	if err != nil {
		if IsExpectedErrors(err, []string{"ParameterCenInstanceId", "ParameterError", "ParameterInstanceId"}) {
			err = WrapErrorf(Error(GetNotFoundMessage("CenInstanceAttachment", id)), NotFoundMsg, ProviderERROR)
			return
		}
		err = WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
		return
	}
	response, _ := raw.(*cbn.DescribeCenAttachedChildInstanceAttributeResponse)
	return *response, nil
}

func (s *CbnService) CenInstanceAttachmentStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeCenInstanceAttachment(id)
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

func (s *CbnService) DescribeCenBandwidthPackage(id string) (object cbn.CenBandwidthPackage, err error) {
	request := cbn.CreateDescribeCenBandwidthPackagesRequest()
	request.RegionId = s.client.RegionId
	filters := make([]cbn.DescribeCenBandwidthPackagesFilter, 0)
	filters = append(filters, cbn.DescribeCenBandwidthPackagesFilter{
		Key:   "CenBandwidthPackageId",
		Value: &[]string{id},
	})
	request.Filter = &filters

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(11*time.Minute, func() *resource.RetryError {
		raw, err := s.client.WithCbnClient(func(cbnClient *cbn.Client) (interface{}, error) {
			return cbnClient.DescribeCenBandwidthPackages(request)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{"AliyunGoClientFailure", "ServiceUnavailable", "Throttling", "Throttling.User"}) {
				wait()
				return resource.RetryableError(err)
			}
			if IsExpectedErrors(err, []string{"ParameterCenInstanceId"}) {
				err = WrapErrorf(Error(GetNotFoundMessage("CenBandwidthPackage", id)), NotFoundMsg, ProviderERROR)
				return resource.NonRetryableError(err)
			}
			err = WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		response, _ := raw.(*cbn.DescribeCenBandwidthPackagesResponse)

		if len(response.CenBandwidthPackages.CenBandwidthPackage) < 1 {
			err = WrapErrorf(Error(GetNotFoundMessage("CenBandwidthPackage", id)), NotFoundMsg, ProviderERROR, response.RequestId)
			return resource.NonRetryableError(err)
		}
		object = response.CenBandwidthPackages.CenBandwidthPackage[0]
		return nil
	})
	return object, WrapError(err)
}

func (s *CbnService) CenBandwidthPackageStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeCenBandwidthPackage(id)
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

func (s *CbnService) ListTagResources(id string, resourceType string) (object interface{}, err error) {
	conn, err := s.client.NewCbnClient()
	if err != nil {
		return nil, WrapError(err)
	}
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
		parts, err := ParseResourceId(id, 2)
		if err != nil {
			return object, WrapError(err)
		}
		request["ResourceId.1"] = parts[resourceIdNum]
	}

	tags := make([]interface{}, 0)
	var response map[string]interface{}

	for {
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err := conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-09-12"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(action, response, request)
			v, err := jsonpath.Get("$.TagResources.TagResource", response)
			if err != nil {
				return resource.NonRetryableError(WrapErrorf(err, FailedGetAttributeMsg, id, "$.TagResources.TagResource", response))
			}

			if v != nil {
				tags = append(tags, v.([]interface{})...)
			}

			return nil
		})
		if err != nil {
			err = WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
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
		added, removed := parsingTags(d)
		conn, err := s.client.NewCbnClient()
		if err != nil {
			return WrapError(err)
		}

		removedTagKeys := make([]string, 0)
		for _, v := range removed {
			if !ignoredTags(v, "") {
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
				parts, err := ParseResourceId(d.Id(), 2)
				if err != nil {
					return WrapError(err)
				}
				request["ResourceId.1"] = parts[resourceIdNum]
			}

			for i, key := range removedTagKeys {
				request[fmt.Sprintf("TagKey.%d", i+1)] = key
			}
			wait := incrementalWait(2*time.Second, 1*time.Second)
			err := resource.Retry(10*time.Minute, func() *resource.RetryError {
				response, err := conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-09-12"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
				if err != nil {
					if NeedRetry(err) {
						wait()
						return resource.RetryableError(err)

					}
					return resource.NonRetryableError(err)
				}
				addDebug(action, response, request)
				return nil
			})
			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
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
				parts, err := ParseResourceId(d.Id(), 2)
				if err != nil {
					return WrapError(err)
				}
				request["ResourceId.1"] = parts[resourceIdNum]
			}

			count := 1
			for key, value := range added {
				request[fmt.Sprintf("Tag.%d.Key", count)] = key
				request[fmt.Sprintf("Tag.%d.Value", count)] = value
				count++
			}

			wait := incrementalWait(2*time.Second, 1*time.Second)
			err := resource.Retry(10*time.Minute, func() *resource.RetryError {
				response, err := conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-09-12"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
				if err != nil {
					if NeedRetry(err) {
						wait()
						return resource.RetryableError(err)

					}
					return resource.NonRetryableError(err)
				}
				addDebug(action, response, request)
				return nil
			})
			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
			}
		}
		d.SetPartial("tags")
	}
	return nil
}

func (s *CbnService) DescribeCenRouteService(id string) (object cbn.RouteServiceEntry, err error) {
	parts, err := ParseResourceId(id, 4)
	if err != nil {
		err = WrapError(err)
		return
	}
	request := cbn.CreateDescribeRouteServicesInCenRequest()
	request.RegionId = s.client.RegionId
	request.AccessRegionId = parts[3]
	request.CenId = parts[0]
	request.Host = parts[2]
	request.HostRegionId = parts[1]

	var raw interface{}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err = s.client.WithCbnClient(func(cbnClient *cbn.Client) (interface{}, error) {
			return cbnClient.DescribeRouteServicesInCen(request)
		})
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})

	if err != nil {
		err = WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
		return
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*cbn.DescribeRouteServicesInCenResponse)

	if len(response.RouteServiceEntries.RouteServiceEntry) < 1 {
		err = WrapErrorf(Error(GetNotFoundMessage("CenRouteService", id)), NotFoundMsg, ProviderERROR, response.RequestId)
		return
	}
	return response.RouteServiceEntries.RouteServiceEntry[0], nil
}

func (s *CbnService) CenRouteServiceStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeCenRouteService(id)
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

func (s *CbnService) DescribeCenTransitRouter(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewCbnClient()
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		return nil, WrapError(err)
	}
	action := "ListTransitRouters"
	request := map[string]interface{}{
		"Region_id":       s.client.RegionId,
		"CenId":           parts[0],
		"TransitRouterId": parts[1],
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-09-12"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		if IsExpectedErrors(err, []string{"ParameterCenInstanceId"}) {
			return nil, WrapErrorf(Error(GetNotFoundMessage("CEN Instance ID", id)), NotFoundMsg, ProviderERROR)
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.TransitRouters", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.TransitRouters", response)
	}
	if len(v.([]interface{})) < 1 {
		return object, WrapErrorf(Error(GetNotFoundMessage("CEN", id)), NotFoundWithResponse, response)
	} else {
		if v.([]interface{})[0].(map[string]interface{})["TransitRouterId"].(string) != parts[1] {
			return object, WrapErrorf(Error(GetNotFoundMessage("CEN", id)), NotFoundWithResponse, response)
		}
	}
	object = v.([]interface{})[0].(map[string]interface{})
	return object, nil
}

func (s *CbnService) CenTransitRouterStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {

		object, err := s.DescribeCenTransitRouter(id)
		if err != nil {
			if NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
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

func (s *CbnService) DescribeCenTransitRouterPeerAttachment(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewCbnClient()
	if err != nil {
		return nil, WrapError(err)
	}
	parts, err1 := ParseResourceId(id, 2)
	if err1 != nil {
		return nil, WrapError(err1)
	}
	action := "ListTransitRouterPeerAttachments"
	request := map[string]interface{}{
		"RegionId":                  s.client.RegionId,
		"CenId":                     parts[0],
		"TransitRouterAttachmentId": parts[1],
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-09-12"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		if IsExpectedErrors(err, []string{"ParameterCenInstanceId", "IllegalParam.Region"}) {
			return nil, WrapErrorf(Error(GetNotFoundMessage("CEN Instance ID", id)), NotFoundMsg, ProviderERROR)
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.TransitRouterAttachments", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.TransitRouterAttachments", response)
	}
	if len(v.([]interface{})) < 1 {
		return object, WrapErrorf(Error(GetNotFoundMessage("CEN", id)), NotFoundWithResponse, response)
	} else {
		if v.([]interface{})[0].(map[string]interface{})["TransitRouterAttachmentId"].(string) != parts[1] {
			return object, WrapErrorf(Error(GetNotFoundMessage("CEN", id)), NotFoundWithResponse, response)
		}
	}
	object = v.([]interface{})[0].(map[string]interface{})
	return object, nil
}

func (s *CbnService) CenTransitRouterPeerAttachmentStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeCenTransitRouterPeerAttachment(id)
		if err != nil {
			if NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
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

func (s *CbnService) DescribeCenTransitRouterVbrAttachment(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	action := "ListTransitRouterVbrAttachments"

	conn, err := s.client.NewCbnClient()
	if err != nil {
		return nil, WrapError(err)
	}

	parts, err := ParseResourceId(id, 2)
	if err != nil {
		return nil, WrapError(err)
	}

	request := map[string]interface{}{
		"RegionId":                  s.client.RegionId,
		"CenId":                     parts[0],
		"TransitRouterAttachmentId": parts[1],
		"MaxResults":                PageSizeLarge,
	}

	idExist := false
	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-09-12"), StringPointer("AK"), nil, request, &runtime)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, request)

		if err != nil {
			if IsExpectedErrors(err, []string{"IllegalParam.Region"}) {
				return object, WrapErrorf(Error(GetNotFoundMessage("Cen:TransitRouterVbrAttachment", id)), NotFoundWithResponse, response)
			}
			return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
		}

		resp, err := jsonpath.Get("$.TransitRouterAttachments", response)
		if err != nil {
			return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.TransitRouterAttachments", response)
		}

		if v, ok := resp.([]interface{}); !ok || len(v) < 1 {
			return object, WrapErrorf(Error(GetNotFoundMessage("Cen:TransitRouterVbrAttachment", id)), NotFoundWithResponse, response)
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
		return object, WrapErrorf(Error(GetNotFoundMessage("Cen:TransitRouterVbrAttachment", id)), NotFoundWithResponse, response)
	}

	return object, nil
}

func (s *CbnService) CenTransitRouterVbrAttachmentStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeCenTransitRouterVbrAttachment(id)
		if err != nil {
			if NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
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

func (s *CbnService) DescribeCenTransitRouterVpcAttachment(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	action := "ListTransitRouterVpcAttachments"

	conn, err := s.client.NewCbnClient()
	if err != nil {
		return nil, WrapError(err)
	}

	parts, err := ParseResourceId(id, 2)
	if err != nil {
		return nil, WrapError(err)
	}

	request := map[string]interface{}{
		"RegionId":                  s.client.RegionId,
		"CenId":                     parts[0],
		"TransitRouterAttachmentId": parts[1],
		"MaxResults":                PageSizeLarge,
	}

	idExist := false
	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-09-12"), StringPointer("AK"), nil, request, &runtime)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, request)

		if err != nil {
			if IsExpectedErrors(err, []string{"IllegalParam.Region"}) {
				return object, WrapErrorf(Error(GetNotFoundMessage("Cen:TransitRouterVpcAttachment", id)), NotFoundWithResponse, response)
			}
			return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
		}

		resp, err := jsonpath.Get("$.TransitRouterAttachments", response)
		if err != nil {
			return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.TransitRouterAttachments", response)
		}

		if v, ok := resp.([]interface{}); !ok || len(v) < 1 {
			return object, WrapErrorf(Error(GetNotFoundMessage("Cen:TransitRouterVpcAttachment", id)), NotFoundWithResponse, response)
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
		return object, WrapErrorf(Error(GetNotFoundMessage("Cen:TransitRouterVpcAttachment", id)), NotFoundWithResponse, response)
	}

	return object, nil
}

func (s *CbnService) CenTransitRouterVpcAttachmentStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeCenTransitRouterVpcAttachment(id)
		if err != nil {
			if NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
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

func (s *CbnService) DescribeCenTransitRouterRouteEntry(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewCbnClient()
	if err != nil {
		return nil, WrapError(err)
	}
	parts, err1 := ParseResourceId(id, 2)
	if err1 != nil {
		return nil, WrapError(err1)
	}
	action := "ListTransitRouterRouteEntries"
	request := map[string]interface{}{
		"TransitRouterRouteTableId":  parts[0],
		"TransitRouterRouteEntryIds": []string{parts[1]},
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-09-12"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidTransitRouterRouteTableId.NotFound"}) {
			return nil, WrapErrorf(Error(GetNotFoundMessage("CEN TransitRouter RouteTable ID", id)), NotFoundMsg, ProviderERROR)
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.TransitRouterRouteEntries", response)
	if err != nil {
		if IsExpectedErrors(err, []string{"ParameterCenInstanceId"}) {
			return nil, WrapErrorf(Error(GetNotFoundMessage("CEN Instance ID", id)), NotFoundMsg, ProviderERROR)
		}
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.TransitRouterRouteEntries", response)
	}
	if len(v.([]interface{})) < 1 {
		return object, WrapErrorf(Error(GetNotFoundMessage("CEN", id)), NotFoundWithResponse, response)
	} else {
		if v.([]interface{})[0].(map[string]interface{})["TransitRouterRouteEntryId"].(string) != parts[1] {
			return object, WrapErrorf(Error(GetNotFoundMessage("CEN", id)), NotFoundWithResponse, response)
		}
	}
	object = v.([]interface{})[0].(map[string]interface{})
	return object, nil
}

func (s *CbnService) CenTransitRouterRouteEntryStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeCenTransitRouterRouteEntry(id)
		if err != nil {
			if NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}

		for _, failState := range failStates {
			if fmt.Sprint(object["TransitRouterRouteEntryStatus"]) == failState {
				return object, fmt.Sprint(object["TransitRouterRouteEntryStatus"]), WrapError(Error(FailedToReachTargetStatus, fmt.Sprint(object["TransitRouterRouteEntryStatus"])))
			}
		}
		return object, fmt.Sprint(object["TransitRouterRouteEntryStatus"]), nil
	}
}

func (s *CbnService) DescribeCenTransitRouterRouteTable(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewCbnClient()
	if err != nil {
		return nil, WrapError(err)
	}
	parts, err1 := ParseResourceId(id, 2)
	if err1 != nil {
		return nil, WrapError(err1)
	}
	action := "ListTransitRouterRouteTables"
	request := map[string]interface{}{
		"RegionId":                   s.client.RegionId,
		"TransitRouterId":            parts[0],
		"TransitRouterRouteTableIds": []string{parts[1]},
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-09-12"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidTransitRouterId.NotFound"}) {
			return nil, WrapErrorf(Error(GetNotFoundMessage("CEN TransitRouter ID", id)), NotFoundMsg, ProviderERROR)
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.TransitRouterRouteTables", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.TransitRouterRouteTables", response)
	}
	if len(v.([]interface{})) < 1 {
		return object, WrapErrorf(Error(GetNotFoundMessage("CEN", id)), NotFoundWithResponse, response)
	} else {
		if v.([]interface{})[0].(map[string]interface{})["TransitRouterRouteTableId"].(string) != parts[1] {
			return object, WrapErrorf(Error(GetNotFoundMessage("CEN", id)), NotFoundWithResponse, response)
		}
	}
	object = v.([]interface{})[0].(map[string]interface{})
	return object, nil
}

func (s *CbnService) CenTransitRouterRouteTableStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeCenTransitRouterRouteTable(id)
		if err != nil {
			if NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}

		for _, failState := range failStates {
			if fmt.Sprint(object["TransitRouterRouteTableStatus"]) == failState {
				return object, fmt.Sprint(object["TransitRouterRouteTableStatus"]), WrapError(Error(FailedToReachTargetStatus, fmt.Sprint(object["TransitRouterRouteTableStatus"])))
			}
		}
		return object, fmt.Sprint(object["TransitRouterRouteTableStatus"]), nil
	}
}

func (s *CbnService) DescribeCenTransitRouterRouteTableAssociation(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewCbnClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "ListTransitRouterRouteTableAssociations"
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		err = WrapError(err)
		return
	}
	request := map[string]interface{}{
		"TransitRouterAttachmentId": parts[0],
		"TransitRouterRouteTableId": parts[1],
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-09-12"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		if IsExpectedErrors(err, []string{"ParameterCenInstanceId"}) {
			return nil, WrapErrorf(Error(GetNotFoundMessage("CEN Instance ID", id)), NotFoundMsg, ProviderERROR)
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.TransitRouterAssociations", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.TransitRouterAssociations", response)
	}
	if len(v.([]interface{})) < 1 {
		return object, WrapErrorf(Error(GetNotFoundMessage("CEN", id)), NotFoundWithResponse, response)
	} else {
		if v.([]interface{})[0].(map[string]interface{})["TransitRouterAttachmentId"].(string) != parts[0] {
			return object, WrapErrorf(Error(GetNotFoundMessage("CEN", id)), NotFoundWithResponse, response)
		}
	}
	object = v.([]interface{})[0].(map[string]interface{})
	return object, nil
}

func (s *CbnService) CenTransitRouterRouteTableAssociationStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeCenTransitRouterRouteTableAssociation(id)
		if err != nil {
			if NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
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

func (s *CbnService) DescribeCenTransitRouterRouteTablePropagation(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewCbnClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "ListTransitRouterRouteTablePropagations"
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		err = WrapError(err)
		return
	}
	request := map[string]interface{}{
		"TransitRouterAttachmentId": parts[0],
		"TransitRouterRouteTableId": parts[1],
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-09-12"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		if IsExpectedErrors(err, []string{"ParameterCenInstanceId"}) {
			return nil, WrapErrorf(Error(GetNotFoundMessage("CEN Instance ID", id)), NotFoundMsg, ProviderERROR)
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.TransitRouterPropagations", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.TransitRouterPropagations", response)
	}
	if len(v.([]interface{})) < 1 {
		return object, WrapErrorf(Error(GetNotFoundMessage("CEN", id)), NotFoundWithResponse, response)
	} else {
		if v.([]interface{})[0].(map[string]interface{})["TransitRouterAttachmentId"].(string) != parts[0] {
			return object, WrapErrorf(Error(GetNotFoundMessage("CEN", id)), NotFoundWithResponse, response)
		}
	}
	object = v.([]interface{})[0].(map[string]interface{})
	return object, nil
}

func (s *CbnService) CenTransitRouterRouteTablePropagationStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeCenTransitRouterRouteTablePropagation(id)
		if err != nil {
			if NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
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

func (s *CbnService) DescribeCenTrafficMarkingPolicy(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewCbnClient()
	if err != nil {
		return nil, WrapError(err)
	}
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		err = WrapError(err)
		return
	}
	action := "ListTrafficMarkingPolicies"
	request := map[string]interface{}{
		"TransitRouterId":        parts[0],
		"TrafficMarkingPolicyId": parts[1],
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-09-12"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.TrafficMarkingPolicies", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.TrafficMarkingPolicies", response)
	}
	if len(v.([]interface{})) < 1 {
		return object, WrapErrorf(Error(GetNotFoundMessage("CEN", id)), NotFoundWithResponse, response)
	} else {
		if v.([]interface{})[0].(map[string]interface{})["TrafficMarkingPolicyId"].(string) != parts[1] {
			return object, WrapErrorf(Error(GetNotFoundMessage("CEN", id)), NotFoundWithResponse, response)
		}
	}
	object = v.([]interface{})[0].(map[string]interface{})
	return object, nil
}

func (s *CbnService) CenTrafficMarkingPolicyStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeCenTrafficMarkingPolicy(id)
		if err != nil {
			if NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}

		for _, failState := range failStates {
			if fmt.Sprint(object["TrafficMarkingPolicyStatus"]) == failState {
				return object, fmt.Sprint(object["TrafficMarkingPolicyStatus"]), WrapError(Error(FailedToReachTargetStatus, fmt.Sprint(object["TrafficMarkingPolicyStatus"])))
			}
		}
		return object, fmt.Sprint(object["TrafficMarkingPolicyStatus"]), nil
	}
}

func (s *CbnService) DescribeCenTransitRouterVpnAttachment(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewCbnClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "ListTransitRouterVpnAttachments"
	request := map[string]interface{}{
		"RegionId":                  s.client.RegionId,
		"TransitRouterAttachmentId": id,
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-09-12"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		if IsExpectedErrors(err, []string{"IllegalParam.Region"}) {
			return nil, WrapErrorf(Error(GetNotFoundMessage("CEN Instance ID", id)), NotFoundMsg, ProviderERROR)
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.TransitRouterAttachments", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.TransitRouterAttachments", response)
	}
	if len(v.([]interface{})) < 1 {
		return object, WrapErrorf(Error(GetNotFoundMessage("CEN", id)), NotFoundWithResponse, response)
	}
	object = v.([]interface{})[0].(map[string]interface{})
	return object, nil
}

func (s *CbnService) CenTransitRouterVpnAttachmentStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeCenTransitRouterVpnAttachment(id)
		if err != nil {
			if NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
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

func (s *CbnService) DescribeCenTransitRouterGrantAttachment(id string) (object map[string]interface{}, err error) {
	conn, err := s.client.NewCbnClient()
	if err != nil {
		return object, WrapError(err)
	}
	parts, err := ParseResourceId(id, 4)
	if err != nil {
		return object, WrapError(err)
	}

	request := map[string]interface{}{
		"ResourceId":  parts[1],
		"ProductType": parts[0],
		"RegionId":    s.client.RegionId,
	}

	var response map[string]interface{}
	action := "DescribeGrantRulesToResource"
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		resp, err := conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-09-12"), StringPointer("AK"), nil, request, &runtime)
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
	v, err := jsonpath.Get("$.GrantRules", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.GrantRules", response)
	}
	if len(v.([]interface{})) < 1 {
		return object, WrapErrorf(Error(GetNotFoundMessage("TransitRouterGrantAttachment", id)), NotFoundWithResponse, response)
	} else {
		for _, vv := range v.([]interface{}) {
			item := vv.(map[string]interface{})
			if fmt.Sprint(item["CenId"]) == parts[3] && fmt.Sprint(item["CenOwnerId"]) == parts[2] {
				return item, nil
			}
		}
	}
	return object, WrapErrorf(Error(GetNotFoundMessage("TransitRouterGrantAttachment", id)), NotFoundWithResponse, response)
}

func (s *CbnService) DescribeCenTransitRouterPrefixListAssociation(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	action := "ListTransitRouterPrefixListAssociation"

	conn, err := s.client.NewCbnClient()
	if err != nil {
		return nil, WrapError(err)
	}

	parts, err := ParseResourceId(id, 4)
	if err != nil {
		return nil, WrapError(err)
	}
	request := map[string]interface{}{
		"RegionId":             s.client.RegionId,
		"PrefixListId":         parts[0],
		"TransitRouterId":      parts[1],
		"TransitRouterTableId": parts[2],
	}

	idExist := false
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-09-12"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}

	resp, err := jsonpath.Get("$.PrefixLists", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.PrefixLists", response)
	}

	if v, ok := resp.([]interface{}); !ok || len(v) < 1 {
		return object, WrapErrorf(Error(GetNotFoundMessage("CEN:TransitRouterPrefixListAssociation", id)), NotFoundWithResponse, response)
	}

	for _, v := range resp.([]interface{}) {
		if fmt.Sprint(v.(map[string]interface{})["NextHop"]) == parts[3] {
			idExist = true
			return v.(map[string]interface{}), nil
		}
	}

	if !idExist {
		return object, WrapErrorf(Error(GetNotFoundMessage("CEN:TransitRouterPrefixListAssociation", id)), NotFoundWithResponse, response)
	}

	return object, nil
}

func (s *CbnService) CenTransitRouterPrefixListAssociationStateRefreshFunc(d *schema.ResourceData, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeCenTransitRouterPrefixListAssociation(d.Id())
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

func (s *CbnService) DescribeCenTransitRouterCidr(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	action := "ListTransitRouterCidr"

	conn, err := s.client.NewCbnClient()
	if err != nil {
		return nil, WrapError(err)
	}

	parts, err := ParseResourceId(id, 2)
	if err != nil {
		return nil, WrapError(err)
	}

	request := map[string]interface{}{
		"RegionId":            s.client.RegionId,
		"ClientToken":         buildClientToken("ListTransitRouterCidr"),
		"TransitRouterId":     parts[0],
		"TransitRouterCidrId": parts[1],
	}

	idExist := false
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-09-12"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}

	resp, err := jsonpath.Get("$.CidrLists", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.CidrLists", response)
	}

	if v, ok := resp.([]interface{}); !ok || len(v) < 1 {
		return object, WrapErrorf(Error(GetNotFoundMessage("CEN:TransitRouterCidr", id)), NotFoundWithResponse, response)
	}

	for _, v := range resp.([]interface{}) {
		if fmt.Sprint(v.(map[string]interface{})["TransitRouterCidrId"]) == parts[1] {
			idExist = true
			return v.(map[string]interface{}), nil
		}
	}

	if !idExist {
		return object, WrapErrorf(Error(GetNotFoundMessage("CEN:TransitRouterCidr", id)), NotFoundWithResponse, response)
	}

	return object, nil
}

func (s *CbnService) DescribeCenTransitRouterMulticastDomainSource(id string) (object map[string]interface{}, err error) {
	conn, err := s.client.NewCbnClient()
	if err != nil {
		return object, WrapError(err)
	}
	parts, err := ParseResourceId(id, 3)
	if err != nil {
		return object, WrapError(err)
	}

	request := map[string]interface{}{
		"TransitRouterMulticastDomainId": parts[0],
		"GroupIpAddress":                 parts[1],
		"NetworkInterfaceIds.1":          parts[2],
	}

	var response map[string]interface{}
	action := "ListTransitRouterMulticastGroups"
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		resp, err := conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-09-12"), StringPointer("AK"), nil, request, &runtime)
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
		if IsExpectedErrors(err, []string{"IllegalParam.TransitRouterMulticastDomainId"}) {
			return object, WrapErrorf(Error(GetNotFoundMessage("CEN:TransitRouterMulticastGroups", id)), NotFoundWithResponse, response)
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.TransitRouterMulticastGroups", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.TransitRouterMulticastGroups", response)
	}
	resp := v.([]interface{})
	if len(resp) < 1 {
		return object, WrapErrorf(Error(GetNotFoundMessage("CEN:TransitRouterMulticastGroups", id)), NotFoundWithResponse, response)
	}
	return resp[0].(map[string]interface{}), nil
}

func (s *CbnService) CenTransitRouterMulticastDomainSourceStateRefreshFunc(d *schema.ResourceData, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeCenTransitRouterMulticastDomainSource(d.Id())
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

func (s *CbnService) DescribeCenTransitRouterMulticastDomain(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	action := "ListTransitRouterMulticastDomains"

	conn, err := s.client.NewCbnClient()
	if err != nil {
		return nil, WrapError(err)
	}

	request := map[string]interface{}{
		"RegionId":                       s.client.RegionId,
		"ClientToken":                    buildClientToken("ListTransitRouterMulticastDomains"),
		"TransitRouterMulticastDomainId": id,
	}

	idExist := false
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-09-12"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}

	resp, err := jsonpath.Get("$.TransitRouterMulticastDomains", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.TransitRouterMulticastDomains", response)
	}

	if v, ok := resp.([]interface{}); !ok || len(v) < 1 {
		return object, WrapErrorf(Error(GetNotFoundMessage("CEN:TransitRouterMulticastDomain", id)), NotFoundWithResponse, response)
	}

	for _, v := range resp.([]interface{}) {
		if fmt.Sprint(v.(map[string]interface{})["TransitRouterMulticastDomainId"]) == id {
			idExist = true
			return v.(map[string]interface{}), nil
		}
	}

	if !idExist {
		return object, WrapErrorf(Error(GetNotFoundMessage("CEN:TransitRouterMulticastDomain", id)), NotFoundWithResponse, response)
	}

	return object, nil
}

func (s *CbnService) CenTransitRouterMulticastDomainStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeCenTransitRouterMulticastDomain(id)
		if err != nil {
			if NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
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

func (s *CbnService) DescribeCenInterRegionTrafficQosQueue(id string) (object map[string]interface{}, err error) {
	conn, err := s.client.NewCbnClient()
	if err != nil {
		return object, WrapError(err)
	}

	request := map[string]interface{}{
		"TrafficQosQueueId": id,
	}

	var response map[string]interface{}
	action := "ListCenInterRegionTrafficQosQueues"
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		resp, err := conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-09-12"), StringPointer("AK"), nil, request, &runtime)
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
	v, err := jsonpath.Get("$.TrafficQosQueues", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.TrafficQosQueues", response)
	}
	resp := v.([]interface{})
	if len(resp) < 1 {
		return object, WrapErrorf(Error(GetNotFoundMessage("CEN:InterRegionTrafficQosQueue", id)), NotFoundWithResponse, response)
	}
	return resp[0].(map[string]interface{}), nil
}

func (s *CbnService) CenInterRegionTrafficQosQueueStateRefreshFunc(d *schema.ResourceData, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeCenInterRegionTrafficQosQueue(d.Id())
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

func (s *CbnService) DescribeCenInterRegionTrafficQosPolicy(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	action := "ListCenInterRegionTrafficQosPolicies"

	conn, err := s.client.NewCbnClient()
	if err != nil {
		return nil, WrapError(err)
	}

	request := map[string]interface{}{
		"TrafficQosPolicyId": id,
	}

	idExist := false
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-09-12"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}

	resp, err := jsonpath.Get("$.TrafficQosPolicies", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.TrafficQosPolicies", response)
	}

	if v, ok := resp.([]interface{}); !ok || len(v) < 1 {
		return object, WrapErrorf(Error(GetNotFoundMessage("CEN:InterRegionTrafficQosPolicy", id)), NotFoundWithResponse, response)
	}

	for _, v := range resp.([]interface{}) {
		if fmt.Sprint(v.(map[string]interface{})["TrafficQosPolicyId"]) == id {
			idExist = true
			return v.(map[string]interface{}), nil
		}
	}

	if !idExist {
		return object, WrapErrorf(Error(GetNotFoundMessage("CEN:InterRegionTrafficQosPolicy", id)), NotFoundWithResponse, response)
	}

	return object, nil
}

func (s *CbnService) CenInterRegionTrafficQosPolicyStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeCenInterRegionTrafficQosPolicy(id)
		if err != nil {
			if NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}

		for _, failState := range failStates {
			if fmt.Sprint(object["TrafficQosPolicyStatus"]) == failState {
				return object, fmt.Sprint(object["TrafficQosPolicyStatus"]), WrapError(Error(FailedToReachTargetStatus, fmt.Sprint(object["TrafficQosPolicyStatus"])))
			}
		}

		return object, fmt.Sprint(object["TrafficQosPolicyStatus"]), nil
	}
}

func (s *CbnService) DescribeCenTransitRouterMulticastDomainPeerMember(id string) (object map[string]interface{}, err error) {
	conn, err := s.client.NewCbnClient()
	if err != nil {
		return object, WrapError(err)
	}
	parts, err := ParseResourceId(id, 3)
	if err != nil {
		return object, WrapError(err)
	}

	request := map[string]interface{}{
		"TransitRouterMulticastDomainId":      parts[0],
		"GroupIpAddress":                      parts[1],
		"PeerTransitRouterMulticastDomains.1": parts[2],
	}

	var response map[string]interface{}
	action := "ListTransitRouterMulticastGroups"
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		resp, err := conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-09-12"), StringPointer("AK"), nil, request, &runtime)
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
	count, err := jsonpath.Get("$.TotalCount", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.TotalCount", response)
	}
	if formatInt(count) == 0 {
		return object, WrapErrorf(Error(GetNotFoundMessage("CEN.TransitRouterMulticastDomainPeerMember", id)), NotFoundWithResponse, response)
	}
	v, err := jsonpath.Get("$.TransitRouterMulticastGroups[0]", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.TransitRouterMulticastGroups[0]", response)
	}
	return v.(map[string]interface{}), nil
}

func (s *CbnService) CenTransitRouterMulticastDomainPeerMemberStateRefreshFunc(d *schema.ResourceData, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeCenTransitRouterMulticastDomainPeerMember(d.Id())
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

func (s *CbnService) DescribeCenTransitRouterMulticastDomainMember(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	action := "ListTransitRouterMulticastGroups"

	conn, err := s.client.NewCbnClient()
	if err != nil {
		return object, WrapError(err)
	}

	parts, err := ParseResourceId(id, 3)
	if err != nil {
		return object, WrapError(err)
	}

	request := map[string]interface{}{
		"ClientToken":                    buildClientToken("ListTransitRouterMulticastGroups"),
		"TransitRouterMulticastDomainId": parts[0],
		"GroupIpAddress":                 parts[1],
		"NetworkInterfaceIds":            []string{parts[2]},
		"MaxResults":                     PageSizeLarge,
	}

	idExist := false
	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-09-12"), StringPointer("AK"), nil, request, &runtime)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, request)

		if err != nil {
			if IsExpectedErrors(err, []string{"IllegalParam.TransitRouterMulticastDomainId"}) {
				return object, WrapErrorf(Error(GetNotFoundMessage("Cen:TransitRouterMulticastDomainMember", id)), NotFoundWithResponse, response)
			}
			return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
		}

		resp, err := jsonpath.Get("$.TransitRouterMulticastGroups", response)
		if err != nil {
			return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.TransitRouterMulticastGroups", response)
		}

		if v, ok := resp.([]interface{}); !ok || len(v) < 1 {
			return object, WrapErrorf(Error(GetNotFoundMessage("Cen:TransitRouterMulticastDomainMember", id)), NotFoundWithResponse, response)
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
		return object, WrapErrorf(Error(GetNotFoundMessage("Cen:TransitRouterMulticastDomainMember", id)), NotFoundWithResponse, response)
	}

	return object, nil
}

func (s *CbnService) CenTransitRouterMulticastDomainMemberStateRefreshFunc(d *schema.ResourceData, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeCenTransitRouterMulticastDomainMember(d.Id())
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

func (s *CbnService) DescribeCenChildInstanceRouteEntryToAttachment(id string) (object map[string]interface{}, err error) {
	conn, err := s.client.NewCbnClient()
	if err != nil {
		return object, WrapError(err)
	}
	parts, err := ParseResourceId(id, 4)
	if err != nil {
		return object, WrapError(err)
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
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		resp, err := conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-09-12"), StringPointer("AK"), nil, request, &runtime)
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
		if IsExpectedErrors(err, []string{"ResourceNotFound.TransitRouterAttachment"}) {
			err = WrapErrorf(Error(GetNotFoundMessage("Cen", id)), NotFoundMsg, ProviderERROR)
			return object, err
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	all, err := jsonpath.Get("$.RouteEntry", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.RouteEntry", response)
	}
	if len(all.([]interface{})) == 0 {
		return object, WrapErrorf(Error(GetNotFoundMessage("Cen", id)), NotFoundWithResponse, response)
	}
	v, err := jsonpath.Get("$.RouteEntry[0]", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.RouteEntry[0]", response)
	}
	return v.(map[string]interface{}), nil
}

func (s *CbnService) CenChildInstanceRouteEntryToAttachmentStateRefreshFunc(d *schema.ResourceData, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeCenChildInstanceRouteEntryToAttachment(d.Id())
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

func (s *CbnService) DescribeCenTransitRouterMulticastDomainAssociation(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	action := "ListTransitRouterMulticastDomainAssociations"

	conn, err := s.client.NewCbnClient()
	if err != nil {
		return nil, WrapError(err)
	}

	parts, err := ParseResourceId(id, 3)
	if err != nil {
		return nil, WrapError(err)
	}

	request := map[string]interface{}{
		"ClientToken":               buildClientToken("ListTransitRouterMulticastDomainAssociations"),
		"TransitRouterAttachmentId": parts[1],
		"VSwitchIds":                []string{parts[2]},
	}

	idExist := false
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-09-12"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}

	resp, err := jsonpath.Get("$.TransitRouterMulticastAssociations", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.TransitRouterMulticastAssociations", response)
	}

	if v, ok := resp.([]interface{}); !ok || len(v) < 1 {
		return object, WrapErrorf(Error(GetNotFoundMessage("CEN:TransitRouterMulticastDomainAssociation", id)), NotFoundWithResponse, response)
	}

	for _, v := range resp.([]interface{}) {
		if fmt.Sprint(v.(map[string]interface{})["VSwitchId"]) == parts[2] {
			idExist = true
			return v.(map[string]interface{}), nil
		}
	}

	if !idExist {
		return object, WrapErrorf(Error(GetNotFoundMessage("CEN:TransitRouterMulticastDomainAssociation", id)), NotFoundWithResponse, response)
	}

	return object, nil
}

func (s *CbnService) CenTransitRouterMulticastDomainAssociationStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeCenTransitRouterMulticastDomainAssociation(id)
		if err != nil {
			if NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
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

func (s *CbnService) DescribeCenTransitRouteTableAggregation(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	action := "DescribeTransitRouteTableAggregation"

	conn, err := s.client.NewCbnClient()
	if err != nil {
		return object, WrapError(err)
	}

	parts, err := ParseResourceId(id, 2)
	if err != nil {
		return nil, WrapError(err)
	}

	request := map[string]interface{}{
		"ClientToken":                      buildClientToken("DescribeTransitRouteTableAggregation"),
		"TransitRouteTableId":              parts[0],
		"TransitRouteTableAggregationCidr": parts[1],
	}

	idExist := false
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-09-12"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			if IsExpectedErrors(err, []string{"Operation.Blocking"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		if IsExpectedErrors(err, []string{"InstanceNotExist.AggregationRoute", "InstanceNotExist.TransitRouterRouteTable"}) {
			return nil, WrapErrorf(Error(GetNotFoundMessage("CEN:TransitRouteTableAggregation", id)), NotFoundWithResponse, response)
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}

	resp, err := jsonpath.Get("$.Data", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Data", response)
	}

	if v, ok := resp.([]interface{}); !ok || len(v) < 1 {
		return object, WrapErrorf(Error(GetNotFoundMessage("CEN:TransitRouteTableAggregation", id)), NotFoundWithResponse, response)
	}

	for _, v := range resp.([]interface{}) {
		if fmt.Sprint(v.(map[string]interface{})["TransitRouteTableAggregationCidr"]) == parts[1] {
			idExist = true
			return v.(map[string]interface{}), nil
		}
	}

	if !idExist {
		return object, WrapErrorf(Error(GetNotFoundMessage("CEN:TransitRouteTableAggregation", id)), NotFoundWithResponse, response)
	}

	return object, nil
}

func (s *CbnService) CenTransitRouteTableAggregationStateRefreshFunc(d *schema.ResourceData, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeCenTransitRouteTableAggregation(d.Id())
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

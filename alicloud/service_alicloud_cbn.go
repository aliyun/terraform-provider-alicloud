package alicloud

import (
	"time"

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

	raw, err := s.client.WithCbnClient(func(cbnClient *cbn.Client) (interface{}, error) {
		return cbnClient.DescribeFlowlogs(request)
	})
	if err != nil {
		err = WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
		return
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
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

func (s *CbnService) DescribeCenInstance(id string) (object cbn.Cen, err error) {
	request := cbn.CreateDescribeCensRequest()
	request.RegionId = s.client.RegionId
	filters := make([]cbn.DescribeCensFilter, 0)
	filters = append(filters, cbn.DescribeCensFilter{
		Key:   "CenId",
		Value: &[]string{id},
	})
	request.Filter = &filters

	raw, err := s.client.WithCbnClient(func(cbnClient *cbn.Client) (interface{}, error) {
		return cbnClient.DescribeCens(request)
	})
	if err != nil {
		err = WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
		return
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*cbn.DescribeCensResponse)

	if len(response.Cens.Cen) < 1 {
		err = WrapErrorf(Error(GetNotFoundMessage("CenInstance", id)), NotFoundMsg, ProviderERROR, response.RequestId)
		return
	}
	return response.Cens.Cen[0], nil
}

func (s *CbnService) CenInstanceStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeCenInstance(id)
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
	for key, _ := range oldItems.(map[string]interface{}) {
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

	raw, err := s.client.WithCbnClient(func(cbnClient *cbn.Client) (interface{}, error) {
		return cbnClient.DescribeCenRouteMaps(request)
	})
	if err != nil {
		err = WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
		return
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
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
	request.RegionId = s.client.RegionId
	request.AccessRegionId = parts[1]
	request.CenId = parts[0]

	raw, err := s.client.WithCbnClient(func(cbnClient *cbn.Client) (interface{}, error) {
		return cbnClient.DescribeCenPrivateZoneRoutes(request)
	})
	if err != nil {
		err = WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
		return
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*cbn.DescribeCenPrivateZoneRoutesResponse)

	if len(response.PrivateZoneInfos.PrivateZoneInfo) < 1 {
		err = WrapErrorf(Error(GetNotFoundMessage("CenPrivateZone", id)), NotFoundMsg, ProviderERROR)
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

	raw, err := s.client.WithCbnClient(func(cbnClient *cbn.Client) (interface{}, error) {
		return cbnClient.DescribeCenVbrHealthCheck(request)
	})
	if err != nil {
		err = WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
		return
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
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

	raw, err := s.client.WithCbnClient(func(cbnClient *cbn.Client) (interface{}, error) {
		return cbnClient.DescribeCenAttachedChildInstanceAttribute(request)
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"ParameterCenInstanceId", "ParameterError", "ParameterInstanceId"}) {
			err = WrapErrorf(Error(GetNotFoundMessage("CenInstanceAttachment", id)), NotFoundMsg, ProviderERROR)
			return
		}
		err = WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
		return
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
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

func (s *CbnService) SetResourceTags(d *schema.ResourceData, resourceType string) error {
	oldItems, newItems := d.GetChange("tags")
	added := make([]cbn.TagResourcesTag, 0)
	for key, value := range newItems.(map[string]interface{}) {
		added = append(added, cbn.TagResourcesTag{
			Key:   key,
			Value: value.(string),
		})
	}
	removed := make([]string, 0)
	for key, _ := range oldItems.(map[string]interface{}) {
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

	raw, err := s.client.WithCbnClient(func(cbnClient *cbn.Client) (interface{}, error) {
		return cbnClient.DescribeRouteServicesInCen(request)
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

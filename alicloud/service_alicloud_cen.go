package alicloud

import (
	"fmt"
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/cbn"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

type CenService struct {
	client *connectivity.AliyunClient
}

const DefaultCenTimeout = 60
const DefaultCenTimeoutLong = 180

const ChildInstanceTypeVpc = "VPC"
const ChildInstanceTypeVbr = "VBR"

func (s *CenService) DescribeCenInstance(cenId string) (c cbn.Cen, err error) {
	request := cbn.CreateDescribeCensRequest()

	values := []string{cenId}
	filters := []cbn.DescribeCensFilter{{
		Key:   "CenId",
		Value: &values,
	}}

	request.Filter = &filters

	invoker := NewInvoker()
	err = invoker.Run(func() error {
		raw, err := s.client.WithCenClient(func(cbnClient *cbn.Client) (interface{}, error) {
			return cbnClient.DescribeCens(request)
		})
		if err != nil {
			if IsExceptedError(err, ParameterCenInstanceIdNotExist) {
				return WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
			}
			return WrapErrorf(err, DefaultErrorMsg, cenId, request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw)

		response, _ := raw.(*cbn.DescribeCensResponse)
		if len(response.Cens.Cen) <= 0 || response.Cens.Cen[0].CenId != cenId {
			return WrapErrorf(Error(GetNotFoundMessage("Cen Instance", cenId)), NotFoundMsg, ProviderERROR)
		}
		c = response.Cens.Cen[0]
		return nil
	})

	return
}

func (s *CenService) WaitForCenInstance(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)

	for {
		object, err := s.DescribeCenInstance(id)
		if err != nil {
			if NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return WrapError(err)
			}
		}
		if object.Status == string(status) {
			return nil
		}
		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, object.Status, string(status), ProviderERROR)
		}
	}
}

func (s *CenService) DescribeCenAttachedChildInstanceById(instanceId, cenId string) (c cbn.ChildInstance, err error) {
	request := cbn.CreateDescribeCenAttachedChildInstancesRequest()
	request.CenId = cenId

	for pageNum := 1; ; pageNum++ {
		request.PageNumber = requests.NewInteger(pageNum)
		raw, err := s.client.WithCenClient(func(cbnClient *cbn.Client) (interface{}, error) {
			return cbnClient.DescribeCenAttachedChildInstances(request)
		})
		response, _ := raw.(*cbn.DescribeCenAttachedChildInstancesResponse)
		if err != nil {
			return c, err
		}

		instanceList := response.ChildInstances.ChildInstance
		for instanceNum := 0; instanceNum <= len(instanceList)-1; instanceNum++ {
			if instanceList[instanceNum].ChildInstanceId == instanceId {
				return instanceList[instanceNum], nil
			}
		}

		if pageNum*response.PageSize >= response.TotalCount {
			return c, GetNotFoundErrorFromString(GetNotFoundMessage("CEN Child Instance", instanceId))
		}
	}
}

func (s *CenService) WaitForCenChildInstanceAttached(instanceId string, cenId string, status Status, timeout int) error {
	if timeout <= 0 {
		timeout = DefaultTimeout
	}

	for {
		instance, err := s.DescribeCenAttachedChildInstanceById(instanceId, cenId)
		if err != nil {
			return err
		}
		if instance.Status == string(status) {
			break
		}
		timeout = timeout - DefaultIntervalShort
		if timeout <= 0 {
			return GetTimeErrorFromString(GetTimeoutMessage("CEN Child Instance Attachment", string(status)))
		}
		time.Sleep(DefaultIntervalShort * time.Second)
	}

	return nil
}

func (s *CenService) WaitForCenChildInstanceDetached(instanceId string, cenId string, timeout int) error {
	if timeout <= 0 {
		timeout = DefaultTimeout
	}

	for {
		_, err := s.DescribeCenAttachedChildInstanceById(instanceId, cenId)
		if err != nil {
			if NotFoundError(err) {
				break
			} else {
				return err
			}
		}

		timeout = timeout - DefaultIntervalShort
		if timeout <= 0 {
			return GetTimeErrorFromString(fmt.Sprintf("Waitting for %s detach timeout.", instanceId))
		}
		time.Sleep(DefaultIntervalShort * time.Second)
	}

	return nil
}

func (s *CenService) DescribeCenBandwidthPackage(cenBwpId string) (c cbn.CenBandwidthPackage, err error) {
	request := cbn.CreateDescribeCenBandwidthPackagesRequest()

	values := []string{cenBwpId}
	filters := []cbn.DescribeCenBandwidthPackagesFilter{{
		Key:   "CenBandwidthPackageId",
		Value: &values,
	}}
	request.Filter = &filters

	invoker := NewInvoker()
	err = invoker.Run(func() error {
		raw, err := s.client.WithCenClient(func(cbnClient *cbn.Client) (interface{}, error) {
			return cbnClient.DescribeCenBandwidthPackages(request)
		})
		if err != nil {
			if IsExceptedError(err, ParameterCenInstanceIdNotExist) {
				return GetNotFoundErrorFromString(GetNotFoundMessage("CEN Bandwidth Package", cenBwpId))
			}
			return err
		}
		resp, _ := raw.(*cbn.DescribeCenBandwidthPackagesResponse)
		if resp == nil || len(resp.CenBandwidthPackages.CenBandwidthPackage) <= 0 || resp.CenBandwidthPackages.CenBandwidthPackage[0].CenBandwidthPackageId != cenBwpId {
			return GetNotFoundErrorFromString(GetNotFoundMessage("CEN Bandwidth Package", cenBwpId))
		}
		c = resp.CenBandwidthPackages.CenBandwidthPackage[0]
		return nil
	})

	return
}

func (s *CenService) WaitForCenBandwidthPackage(cenBwpId string, status Status, timeout int) error {
	if timeout <= 0 {
		timeout = DefaultTimeout
	}

	for {
		cenBwp, err := s.DescribeCenBandwidthPackage(cenBwpId)
		if err != nil && !NotFoundError(err) {
			return err
		}
		if cenBwp.Status == string(status) {
			break
		}
		timeout = timeout - DefaultIntervalShort
		if timeout <= 0 {
			return GetTimeErrorFromString(GetTimeoutMessage("CEN Bandwidth Package", string(status)))
		}
		time.Sleep(DefaultIntervalShort * time.Second)
	}

	return nil
}

func (s *CenService) WaitForCenBandwidthPackageUpdate(cenBwpId string, bandwidth int, timeout int) error {
	if timeout <= 0 {
		timeout = DefaultTimeout
	}

	for {
		cenBwp, err := s.DescribeCenBandwidthPackage(cenBwpId)
		if err != nil {
			return err
		}
		if cenBwp.Bandwidth == bandwidth {
			break
		}

		timeout = timeout - DefaultIntervalShort
		if timeout <= 0 {
			return GetTimeErrorFromString(fmt.Sprintf("Waitting for CEN bandwidth package update is timeout"))
		}
		time.Sleep(DefaultIntervalShort * time.Second)
	}

	return nil
}

func (s *CenService) WaitForCenBandwidthPackageAttachment(cenBwpId string, status Status, timeout int) error {
	if timeout <= 0 {
		timeout = DefaultTimeout
	}

	for {
		cenBwp, err := s.DescribeCenBandwidthPackage(cenBwpId)
		if err != nil {
			return err
		}
		if cenBwp.Status == string(status) {
			break
		}
		timeout = timeout - DefaultIntervalShort
		if timeout <= 0 {
			return GetTimeErrorFromString(GetTimeoutMessage("CEN Bandwidth Package Attachment", string(status)))
		}
		time.Sleep(DefaultIntervalShort * time.Second)
	}

	return nil
}

func (s *CenService) DescribeCenBandwidthPackageById(cenBwpId string) (c cbn.CenBandwidthPackage, err error) {
	resp, err := s.DescribeCenBandwidthPackage(cenBwpId)
	if err != nil {
		return c, err
	}

	if len(resp.CenIds.CenId) != 1 || resp.Status != string(InUse) {
		return c, GetNotFoundErrorFromString(GetNotFoundMessage("CEN bandwidth package attachment", cenBwpId))
	}

	return resp, nil
}

func (s *CenService) SetCenInterRegionBandwidthLimit(cenId, localRegionId, oppositeRegionId string, bandwidthLimit int) (err error) {
	request := cbn.CreateSetCenInterRegionBandwidthLimitRequest()
	request.CenId = cenId
	request.LocalRegionId = localRegionId
	request.OppositeRegionId = oppositeRegionId
	request.BandwidthLimit = requests.NewInteger(bandwidthLimit)

	_, err = s.client.WithCenClient(func(cbnClient *cbn.Client) (interface{}, error) {
		return cbnClient.SetCenInterRegionBandwidthLimit(request)
	})

	return err
}

func (s *CenService) DescribeCenBandwidthLimit(cenId, localRegionId, oppositeRegionId string) (c cbn.CenInterRegionBandwidthLimit, err error) {
	request := cbn.CreateDescribeCenInterRegionBandwidthLimitsRequest()
	request.CenId = cenId

	for pageNum := 1; ; pageNum++ {
		request.PageNumber = requests.NewInteger(pageNum)
		request.PageSize = requests.NewInteger(PageSizeLarge)
		raw, err := s.client.WithCenClient(func(cbnClient *cbn.Client) (interface{}, error) {
			return cbnClient.DescribeCenInterRegionBandwidthLimits(request)
		})
		if err != nil {
			return c, err
		}
		resp, _ := raw.(*cbn.DescribeCenInterRegionBandwidthLimitsResponse)

		cenBandwidthLimitList := resp.CenInterRegionBandwidthLimits.CenInterRegionBandwidthLimit
		for limitNum := 0; limitNum <= len(cenBandwidthLimitList)-1; limitNum++ {
			ifMatch := cenBandwidthLimitList[limitNum].LocalRegionId == localRegionId && cenBandwidthLimitList[limitNum].OppositeRegionId == oppositeRegionId
			if !ifMatch {
				ifMatch = cenBandwidthLimitList[limitNum].LocalRegionId == oppositeRegionId && cenBandwidthLimitList[limitNum].OppositeRegionId == localRegionId
			}
			if ifMatch {
				return cenBandwidthLimitList[limitNum], nil
			}
		}

		if pageNum*resp.PageSize >= resp.TotalCount {
			return c, GetNotFoundErrorFromString(fmt.Sprintf("The specified CEN bandwith limit CEN Id %s localRegionId %s oppositeRegionId %s is not found", cenId, localRegionId, oppositeRegionId))
		}
	}
}

func (s *CenService) WaitForCenInterRegionBandwidthLimitActive(cenId string, localRegionId string, oppositeRegionId string, timeout int) error {
	if timeout <= 0 {
		timeout = DefaultTimeout
	}

	for {
		cenBandwidthLimit, err := s.DescribeCenBandwidthLimit(cenId, localRegionId, oppositeRegionId)
		if err != nil {
			return err
		}

		if cenBandwidthLimit.Status == string(Active) {
			break
		}

		timeout = timeout - DefaultIntervalShort
		if timeout <= 0 {
			return GetTimeErrorFromString(fmt.Sprintf("Waitting for bandwidth limit CenId %s localRegionId %s oppositeRegionId %s timeout.", cenId, localRegionId, oppositeRegionId))
		}
		time.Sleep(DefaultIntervalShort * time.Second)
	}

	return nil
}

func (s *CenService) WaitForCenInterRegionBandwidthLimitDestroy(cenId string, localRegionId string, oppositeRegionId string, timeout int) error {
	if timeout <= 0 {
		timeout = DefaultTimeout
	}

	for {
		_, err := s.DescribeCenBandwidthLimit(cenId, localRegionId, oppositeRegionId)
		if err != nil {
			if NotFoundError(err) {
				break
			}
			return err
		}

		timeout = timeout - DefaultIntervalShort
		if timeout <= 0 {
			return GetTimeErrorFromString(fmt.Sprintf("Waitting for bandwidth limit CenId %s localRegionId %s oppositeRegionId %s timeout.", cenId, localRegionId, oppositeRegionId))
		}
		time.Sleep(DefaultIntervalShort * time.Second)
	}

	return nil
}

func (s *CenService) CreateCenRouteEntryParas(vtbId string) (childInstanceId string, instanceType string, err error) {
	vpcService := VpcService{s.client}
	routeTableService := RouteTableService{s.client}
	//Query VRouterId and judge whether it is a vbr
	vtb1, err := vpcService.QueryRouteTableById(vtbId)
	if err != nil {
		return childInstanceId, instanceType, err
	}

	if strings.HasPrefix(vtb1.VRouterId, "vbr") {
		return vtb1.VRouterId, ChildInstanceTypeVbr, nil
	}
	//if the VRouterId belonged to a VPC, get the VPC ID
	vtb2, err := routeTableService.DescribeRouteTable(vtbId)
	if err != nil {
		return childInstanceId, instanceType, err
	}
	return vtb2.VpcId, ChildInstanceTypeVpc, nil
}

func (s *CenService) DescribePublishedRouteEntriesById(id string) (c cbn.PublishedRouteEntry, err error) {
	parts := strings.Split(id, COLON_SEPARATED)
	if len(parts) != 3 {
		return c, fmt.Errorf("invalid resource id")
	}
	cenId := parts[0]
	vtbId := parts[1]
	cidr := parts[2]

	childInstanceId, childInstanceType, err := s.CreateCenRouteEntryParas(vtbId)
	if err != nil {
		return c, err
	}

	request := cbn.CreateDescribePublishedRouteEntriesRequest()
	request.CenId = cenId
	request.ChildInstanceId = childInstanceId
	request.ChildInstanceType = childInstanceType
	request.ChildInstanceRegionId = s.client.RegionId
	request.ChildInstanceRouteTableId = vtbId
	request.DestinationCidrBlock = cidr

	invoker := NewInvoker()
	err = invoker.Run(func() error {
		raw, err := s.client.WithCenClient(func(cbnClient *cbn.Client) (interface{}, error) {
			return cbnClient.DescribePublishedRouteEntries(request)
		})
		if err != nil {
			if IsExceptedErrors(err, []string{ParameterIllegal, ParameterIllegalCenInstanceId, InstanceNotExist}) {
				return GetNotFoundErrorFromString(GetNotFoundMessage("CEN RouteEntries", id))
			}
			return err
		}
		resp, _ := raw.(*cbn.DescribePublishedRouteEntriesResponse)
		if resp == nil || len(resp.PublishedRouteEntries.PublishedRouteEntry) <= 0 {
			return GetNotFoundErrorFromString(GetNotFoundMessage("CEN RouteEntries", id))
		}
		c = resp.PublishedRouteEntries.PublishedRouteEntry[0]

		return nil
	})

	return
}

func (s *CenService) WaitForRouterEntryPublished(id string, status Status, timeout int) error {
	if timeout <= 0 {
		timeout = DefaultTimeout
	}

	for {
		routeEntry, err := s.DescribePublishedRouteEntriesById(id)
		if err != nil {
			return nil
		}

		if string(status) == routeEntry.PublishStatus {
			break
		}

		timeout = timeout - DefaultIntervalShort
		if timeout <= 0 {
			return GetTimeErrorFromString(GetTimeoutMessage("CEN RouteEntries", string(status)))
		}
		time.Sleep(DefaultIntervalShort * time.Second)
	}

	return nil
}

func (s *CenService) GetCenIdAndAnotherId(id string) (string, string, error) {
	parts := strings.Split(id, COLON_SEPARATED)

	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid resource id")
	}

	return parts[0], parts[1], nil
}

func (s *CenService) GetCenAndRegionIds(id string) (retString []string, err error) {
	parts := strings.Split(id, COLON_SEPARATED)

	if len(parts) != 3 {
		return retString, fmt.Errorf("invalid resource id")
	}

	return parts, nil
}

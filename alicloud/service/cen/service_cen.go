package cen

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/cbn"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

type CenService struct {
	client *connectivity.AliyunClient
}

// NewCenService is the exported constructor: the client field stays
// unexported, so callers outside the package (including the acceptance
// tests in package cen_test) build values with it.
func NewCenService(client *connectivity.AliyunClient) *CenService {
	return &CenService{client: client}
}

// Status mirrors the root-package string kind used by the legacy Cen
// lifecycle waiters; it is local to this package so the migrated code does
// not import the root alicloud package.
type Status string

const DefaultCenTimeout = 60
const DefaultCenTimeoutLong = 180

const ChildInstanceTypeVpc = "VPC"
const ChildInstanceTypeVbr = "VBR"
const ChildInstanceTypeCcn = "CCN"

// Status lifecycle constants used by the legacy Cen waiters. They mirror the
// root-package values and are local to this package so the migrated code does
// not import the root alicloud package.
const COLON_SEPARATED = ":"
const Deleted Status = "Deleted"
const InUse Status = "InUse"
const NOPUBLISHED Status = "NonPublished"
const Idle Status = "Idle"
const PUBLISHED Status = "Published"

// PayType mirrors the root-package string kind used by the legacy CEN
// bandwidth-package charge-type logic; local to this package so the migrated
// code does not import the root alicloud package.
type PayType string

const PostPaid = PayType("PostPaid")

// GetCenChildInstanceType resolves a CEN child-instance kind from its ID
// prefix. It mirrors the root-package helper (which is still used by the
// unmigrated root common unit test) so this package stays self-contained.
func GetCenChildInstanceType(id string) (c string, e error) {
	if strings.HasPrefix(id, "vpc") {
		return ChildInstanceTypeVpc, nil
	} else if strings.HasPrefix(id, "vbr") {
		return ChildInstanceTypeVbr, nil
	} else if strings.HasPrefix(id, "ccn") {
		return ChildInstanceTypeCcn, nil
	} else {
		return c, fmt.Errorf("CEN child instance ID invalid. Now, it only supports VPC or VBR or CCN instance.")
	}
}

func (s *CenService) DescribeCenInstanceAttachment(id string) (*cbn.ChildInstance, error) {
	c := &cbn.ChildInstance{}
	request := cbn.CreateDescribeCenAttachedChildInstancesRequest()
	request.RegionId = s.client.RegionId
	cenId, instanceId, err := s.GetCenIdAndAnotherId(id)
	if err != nil {
		return c, errs.WrapError(err)
	}
	request.CenId = cenId

	for pageNum := 1; ; pageNum++ {
		request.PageNumber = requests.NewInteger(pageNum)
		var raw interface{}
		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(5*time.Minute, func() *retry.RetryError {
			raw, err = s.client.WithCbnClient(func(cbnClient *cbn.Client) (interface{}, error) {
				return cbnClient.DescribeCenAttachedChildInstances(request)
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
			if errs.IsExpectedErrors(err, []string{"ParameterInstanceId"}) {
				return nil, errs.WrapErrorf(errs.NotFoundErr("CEN Instance Attachment", instanceId), errs.NotFoundMsg, errs.ProviderERROR)
			}
			return c, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, request.GetActionName(), errs.ProviderERROR)
		}
		response, _ := raw.(*cbn.DescribeCenAttachedChildInstancesResponse)

		instanceList := response.ChildInstances.ChildInstance
		for instanceNum := 0; instanceNum <= len(instanceList)-1; instanceNum++ {
			if instanceList[instanceNum].ChildInstanceId == instanceId {
				return &instanceList[instanceNum], nil
			}
		}

		if pageNum*response.PageSize >= response.TotalCount {
			return c, errs.WrapErrorf(errs.NotFoundErr("CEN Instance Attachment", instanceId), errs.NotFoundMsg, errs.ProviderERROR)
		}
	}
}

func (s *CenService) WaitForCenInstanceAttachment(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribeCenInstanceAttachment(id)
		if err != nil {
			if errs.NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return errs.WrapError(err)
			}
		}
		if object.Status == "Attached" {
			break
		}
		if time.Now().After(deadline) {
			return errs.WrapErrorf(err, errs.WaitTimeoutMsg, id, helper.GetFunc(1), timeout, object.Status, status, errs.ProviderERROR)
		}
		time.Sleep(helper.DefaultIntervalShort * time.Second)
	}

	return nil
}

func (s *CenService) DescribeCenBandwidthPackage(id string) (c cbn.CenBandwidthPackage, err error) {
	request := cbn.CreateDescribeCenBandwidthPackagesRequest()
	request.RegionId = s.client.RegionId
	values := []string{id}
	filters := []cbn.DescribeCenBandwidthPackagesFilter{{
		Key:   "CenBandwidthPackageId",
		Value: &values,
	}}
	request.Filter = &filters

	var raw interface{}
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(10*time.Minute, func() *retry.RetryError {
		raw, err = s.client.WithCbnClient(func(cbnClient *cbn.Client) (interface{}, error) {
			return cbnClient.DescribeCenBandwidthPackages(request)
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
		if errs.IsExpectedErrors(err, []string{"ParameterCenInstanceId"}) {
			return c, errs.WrapErrorf(err, errs.NotFoundMsg, errs.AlibabaCloudSdkGoERROR)
		}
		return c, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
	}

	response, _ := raw.(*cbn.DescribeCenBandwidthPackagesResponse)
	if len(response.CenBandwidthPackages.CenBandwidthPackage) <= 0 || response.CenBandwidthPackages.CenBandwidthPackage[0].CenBandwidthPackageId != id {
		return c, errs.WrapErrorf(errs.NotFoundErr("CEN Bandwidth Package", id), errs.NotFoundMsg, errs.ProviderERROR)
	}
	c = response.CenBandwidthPackages.CenBandwidthPackage[0]
	return c, nil
}

func (s *CenService) WaitForCenBandwidthPackage(id string, status Status, bandwidth, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)

	for {
		object, err := s.DescribeCenBandwidthPackage(id)
		if err != nil {
			if errs.NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return errs.WrapError(err)
			}
		}
		if object.Status == string(status) && strconv.FormatInt(object.Bandwidth, 10) == strconv.Itoa(bandwidth) {
			return nil
		}
		if time.Now().After(deadline) {
			return errs.WrapErrorf(err, errs.WaitTimeoutMsg, id, helper.GetFunc(1), timeout, object.Status, status, errs.ProviderERROR)
		}
		time.Sleep(helper.DefaultIntervalShort * time.Second)
	}
}

func (s *CenService) DescribeCenBandwidthPackageAttachment(id string) (c cbn.CenBandwidthPackage, err error) {
	object, err := s.DescribeCenBandwidthPackage(id)
	if err != nil {
		return c, errs.WrapError(err)
	}

	if len(object.CenIds.CenId) != 1 || object.Status != string(InUse) {
		return c, errs.WrapErrorf(errs.NotFoundErr("CenBandwidthPackageAttachment", id), errs.NotFoundMsg, errs.ProviderERROR)
	}

	return object, nil
}

func (s *CenService) WaitForCenBandwidthPackageAttachment(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)

	for {
		object, err := s.DescribeCenBandwidthPackageAttachment(id)
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
			break
		}
		if time.Now().After(deadline) {
			return errs.WrapErrorf(err, errs.WaitTimeoutMsg, id, helper.GetFunc(1), timeout, object.Status, status, errs.ProviderERROR)
		}
		time.Sleep(helper.DefaultIntervalShort * time.Second)
	}

	return nil
}

func (s *CenService) SetCenInterRegionBandwidthLimit(cenId, localRegionId, oppositeRegionId string, bandwidthLimit int) (err error) {
	request := cbn.CreateSetCenInterRegionBandwidthLimitRequest()
	request.RegionId = s.client.RegionId
	request.CenId = cenId
	request.LocalRegionId = localRegionId
	request.OppositeRegionId = oppositeRegionId
	request.BandwidthLimit = requests.NewInteger(bandwidthLimit)

	raw, err := s.client.WithCbnClient(func(cbnClient *cbn.Client) (interface{}, error) {
		return cbnClient.SetCenInterRegionBandwidthLimit(request)
	})
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"InvalidOperation.CenInstanceStatus"}) {
			return errs.WrapError(err)
		}
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_cen_bandwidth_limit", request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
	}
	logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)
	return nil
}

func (s *CenService) DescribeCenBandwidthLimit(id string) (c cbn.CenInterRegionBandwidthLimit, err error) {
	request := cbn.CreateDescribeCenInterRegionBandwidthLimitsRequest()
	request.RegionId = s.client.RegionId
	paras, err := s.GetCenAndRegionIds(id)
	if err != nil {
		return c, errs.WrapError(err)
	}

	cenId := paras[0]
	localRegionId := paras[1]
	oppositeRegionId := paras[2]
	if strings.Compare(localRegionId, oppositeRegionId) > 0 {
		localRegionId, oppositeRegionId = oppositeRegionId, localRegionId
	}
	request.CenId = cenId

	for pageNum := 1; ; pageNum++ {
		request.PageNumber = requests.NewInteger(pageNum)
		request.PageSize = requests.NewInteger(helper.PageSizeLarge)
		var raw interface{}
		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(10*time.Minute, func() *retry.RetryError {
			raw, err = s.client.WithCbnClient(func(cbnClient *cbn.Client) (interface{}, error) {
				return cbnClient.DescribeCenInterRegionBandwidthLimits(request)
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

		if err != nil {
			return c, errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_cen_bandwidth_limit", request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
		}
		logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)
		response, _ := raw.(*cbn.DescribeCenInterRegionBandwidthLimitsResponse)

		cenBandwidthLimitList := response.CenInterRegionBandwidthLimits.CenInterRegionBandwidthLimit
		for limitNum := 0; limitNum <= len(cenBandwidthLimitList)-1; limitNum++ {
			ifMatch := cenBandwidthLimitList[limitNum].LocalRegionId == localRegionId && cenBandwidthLimitList[limitNum].OppositeRegionId == oppositeRegionId
			if !ifMatch {
				ifMatch = cenBandwidthLimitList[limitNum].LocalRegionId == oppositeRegionId && cenBandwidthLimitList[limitNum].OppositeRegionId == localRegionId
			}
			if ifMatch {
				return cenBandwidthLimitList[limitNum], nil
			}
		}

		if pageNum*response.PageSize >= response.TotalCount {
			return c, errs.WrapErrorf(errs.NotFoundErr("CenBandwidthLimit", id), errs.NotFoundMsg, errs.ProviderERROR)
		}
	}
}

func (s *CenService) CenBandwidthLimitStateRefreshFunc(id string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeCenBandwidthLimit(id)
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

func (s *CenService) CreateCenRouteEntryParas(vtbId string) (childInstanceId string, instanceType string, err error) {
	//Query VRouterId and judge whether it is a vbr
	vtb1, err := s.queryRouteTableById(vtbId)
	if err != nil {
		return childInstanceId, instanceType, errs.WrapError(err)
	}

	if strings.HasPrefix(vtb1.VRouterId, "vbr") {
		return vtb1.VRouterId, ChildInstanceTypeVbr, nil
	}
	//if the VRouterId belonged to a VPC, get the VPC ID
	vtb2, err := s.describeRouteTable(vtbId)
	if err != nil {
		return childInstanceId, instanceType, errs.WrapError(err)
	}
	return fmt.Sprint(vtb2["VpcId"]), ChildInstanceTypeVpc, nil
}

func (s *CenService) DescribeCenRouteEntry(id string) (c cbn.PublishedRouteEntry, err error) {
	parts, err := helper.ParseResourceId(id, 3)
	if err != nil {
		return c, errs.WrapError(err)
	}
	cenId := parts[0]
	vtbId := parts[1]
	cidr := parts[2]

	childInstanceId, childInstanceType, err := s.CreateCenRouteEntryParas(vtbId)
	if err != nil {
		return c, errs.WrapError(err)
	}

	request := cbn.CreateDescribePublishedRouteEntriesRequest()
	request.RegionId = s.client.RegionId
	request.CenId = cenId
	request.ChildInstanceId = childInstanceId
	request.ChildInstanceType = childInstanceType
	request.ChildInstanceRegionId = s.client.RegionId
	request.ChildInstanceRouteTableId = vtbId
	request.DestinationCidrBlock = cidr

	var raw interface{}
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		raw, err = s.client.WithCbnClient(func(cbnClient *cbn.Client) (interface{}, error) {
			return cbnClient.DescribePublishedRouteEntries(request)
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
		if errs.IsExpectedErrors(err, []string{"ParameterIllegal", "ParameterIllegal.CenInstanceId", "Instance.NotExist"}) {
			return c, errs.WrapErrorf(err, errs.NotFoundMsg, errs.AlibabaCloudSdkGoERROR)
		}
		return c, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, request.GetActionName(), errs.AlibabaCloudSdkGoERROR)

	}
	response, _ := raw.(*cbn.DescribePublishedRouteEntriesResponse)
	if len(response.PublishedRouteEntries.PublishedRouteEntry) <= 0 || response.PublishedRouteEntries.PublishedRouteEntry[0].PublishStatus == string(NOPUBLISHED) {
		return c, errs.WrapErrorf(errs.NotFoundErr("CenRouteEntries", id), errs.NotFoundMsg, errs.ProviderERROR)
	}
	c = response.PublishedRouteEntries.PublishedRouteEntry[0]

	return c, nil
}

func (s *CenService) WaitForCenRouterEntry(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)

	for {
		object, err := s.DescribeCenRouteEntry(id)
		if err != nil {
			if errs.NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return errs.WrapError(err)
			}
		}

		if object.PublishStatus == string(status) {
			break
		}

		if time.Now().After(deadline) {
			return errs.WrapErrorf(err, errs.WaitTimeoutMsg, id, helper.GetFunc(1), timeout, object.PublishStatus, string(status), errs.ProviderERROR)
		}
		time.Sleep(helper.DefaultIntervalShort * time.Second)
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

func (s *CenService) CenInstanceAttachmentStateRefreshFunc(id string, failStates []string) retry.StateRefreshFunc {
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

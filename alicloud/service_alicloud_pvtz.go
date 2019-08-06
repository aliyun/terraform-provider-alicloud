package alicloud

import (
	"strconv"

	"github.com/hashicorp/terraform/helper/resource"

	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/pvtz"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

type PvtzService struct {
	client *connectivity.AliyunClient
}

func (s *PvtzService) DescribePvtzZone(id string) (zone pvtz.DescribeZoneInfoResponse, err error) {
	request := pvtz.CreateDescribeZoneInfoRequest()
	request.ZoneId = id

	var response *pvtz.DescribeZoneInfoResponse
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err := s.client.WithPvtzClient(func(pvtzClient *pvtz.Client) (interface{}, error) {
			return pvtzClient.DescribeZoneInfo(request)
		})
		if err != nil {
			if IsExceptedErrors(err, []string{ServiceUnavailable, PvtzThrottlingUser, PvtzSystemBusy}) {
				time.Sleep(5 * time.Second)
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw)
		response, _ = raw.(*pvtz.DescribeZoneInfoResponse)
		return nil
	})
	if err != nil {
		if IsExceptedErrors(err, []string{ZoneNotExists, ZoneVpcNotExists}) {
			return zone, WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
		}
		return zone, WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	if response.ZoneId != id {
		return zone, WrapErrorf(Error(GetNotFoundMessage("PvtzZone", id)), NotFoundMsg, ProviderERROR)
	}
	zone = *response

	return
}

func (s *PvtzService) DescribePvtzZoneAttachment(id string) (object pvtz.DescribeZoneInfoResponse, err error) {
	object, err = s.DescribePvtzZone(id)
	if err != nil {
		err = WrapError(err)
		return
	}

	if len(object.BindVpcs.Vpc) < 1 {
		err = WrapErrorf(Error(GetNotFoundMessage("PvtzZoneAttachment", id)), NotFoundMsg, ProviderERROR)
	}

	return
}

func (s *PvtzService) WaitForZoneAttachment(id string, vpcIdMap map[string]string, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	var vpcId string
	for {
		object, err := s.DescribePvtzZoneAttachment(id)
		if err != nil && !NotFoundError(err) {
			return WrapError(err)
		}

		equal := true
		if len(object.BindVpcs.Vpc) == len(vpcIdMap) {
			for _, vpc := range object.BindVpcs.Vpc {
				if _, ok := vpcIdMap[vpc.VpcId]; !ok {
					equal = false
					vpcId = vpc.VpcId
					break
				}
			}
		}
		if equal {
			break
		}
		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, "", vpcId, ProviderERROR)
		}
	}
	return nil
}

func (s *PvtzService) DescribePvtzZoneRecord(id string) (record pvtz.Record, err error) {
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		return record, WrapError(err)
	}
	request := pvtz.CreateDescribeZoneRecordsRequest()
	request.ZoneId = parts[1]
	request.PageNumber = requests.NewInteger(1)
	request.PageSize = requests.NewInteger(PageSizeLarge)

	recordIdStr := parts[0]
	var response *pvtz.DescribeZoneRecordsResponse

	for {
		err := resource.Retry(5*time.Minute, func() *resource.RetryError {
			raw, err := s.client.WithPvtzClient(func(pvtzClient *pvtz.Client) (interface{}, error) {
				return pvtzClient.DescribeZoneRecords(request)
			})
			if err != nil {
				if IsExceptedErrors(err, []string{ServiceUnavailable, PvtzThrottlingUser, PvtzSystemBusy}) {
					time.Sleep(5 * time.Second)
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(request.GetActionName(), raw)
			response, _ = raw.(*pvtz.DescribeZoneRecordsResponse)
			return nil
		})
		if err != nil {
			if IsExceptedErrors(err, []string{ZoneNotExists, ZoneVpcNotExists}) {
				return record, WrapErrorf(Error(GetNotFoundMessage("ZoneRecord", id)), NotFoundMsg, AlibabaCloudSdkGoERROR)
			}
			return record, WrapErrorf(err, DefaultErrorMsg, recordIdStr, request.GetActionName(), AlibabaCloudSdkGoERROR)
		}

		if len(response.Records.Record) < 1 {
			return record, WrapErrorf(Error(GetNotFoundMessage("ZoneRecord", id)), NotFoundMsg, ProviderERROR)
		}

		for _, rec := range response.Records.Record {
			if strconv.Itoa(rec.RecordId) == parts[0] {
				record = rec
				return record, nil
			}
		}
		if len(response.Records.Record) < PageSizeLarge {
			break
		}

		if page, err := getNextpageNumber(request.PageNumber); err != nil {
			return record, WrapError(err)
		} else {
			request.PageNumber = page
		}
	}

	return record, WrapErrorf(Error(GetNotFoundMessage("ZoneRecord", recordIdStr)), NotFoundMsg, ProviderERROR)

}

func (s *PvtzService) WaitForPvtzZone(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)

	for {
		object, err := s.DescribePvtzZone(id)
		if err != nil {
			if NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return WrapError(err)
			}
		}
		if object.ZoneId == id {
			return nil
		}
		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, object.ZoneId, id, ProviderERROR)
		}

	}
}

func (s *PvtzService) WaitForPvtzZoneAttachment(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)

	for {
		object, err := s.DescribePvtzZoneAttachment(id)
		if err != nil {
			if NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return WrapError(err)
			}
		}
		if object.ZoneId == id {
			return nil
		}
		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, object.ZoneId, id, ProviderERROR)
		}

	}
}

func (s *PvtzService) WaitForPvtzZoneRecord(id string, status Status, timeout int) error {
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		return WrapError(err)
	}
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)

	for {
		object, err := s.DescribePvtzZoneRecord(id)
		if err != nil {
			if NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return WrapError(err)
			}
		}
		if strconv.Itoa(object.RecordId) == parts[0] {
			return nil
		}
		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, strconv.Itoa(object.RecordId), id, ProviderERROR)
		}

	}
}

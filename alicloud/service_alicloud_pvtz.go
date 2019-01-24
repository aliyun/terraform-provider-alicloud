package alicloud

import (
	"strconv"

	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/pvtz"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

type PvtzService struct {
	client *connectivity.AliyunClient
}

func (s *PvtzService) DescribePvtzZoneInfo(zoneId string) (zone pvtz.DescribeZoneInfoResponse, err error) {
	request := pvtz.CreateDescribeZoneInfoRequest()
	request.ZoneId = zoneId

	invoker := PvtzInvoker()
	err = invoker.Run(func() error {
		raw, err := s.client.WithPvtzClient(func(pvtzClient *pvtz.Client) (interface{}, error) {
			return pvtzClient.DescribeZoneInfo(request)
		})
		if err != nil {
			if IsExceptedErrors(err, []string{ZoneNotExists, ZoneVpcNotExists}) {
				return WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
			}
			return WrapErrorf(err, DefaultErrorMsg, zoneId, request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		resp, _ := raw.(*pvtz.DescribeZoneInfoResponse)
		if resp == nil || resp.ZoneId != zoneId {
			return WrapErrorf(Error(GetNotFoundMessage("PrivateZone", zoneId)), NotFoundMsg, ProviderERROR)
		}
		zone = *resp
		return nil
	})

	return
}

func (s *PvtzService) DescribePvtzZoneAttachment(zoneId string) (zone pvtz.DescribeZoneInfoResponse, err error) {
	zone, err = s.DescribePvtzZoneInfo(zoneId)
	if err != nil {
		err = WrapError(err)
		return
	}

	if len(zone.BindVpcs.Vpc) < 1 {
		err = WrapErrorf(Error(GetNotFoundMessage("PrivateZone Attachment", zoneId)), NotFoundMsg, ProviderERROR)
	}

	return
}

func (s *PvtzService) WaitZoneAttachment(zoneId string, vpcIdMap map[string]string, timeout int) error {
	if timeout <= 0 {
		timeout = DefaultTimeout
	}
	for {
		zone, err := s.DescribePvtzZoneAttachment(zoneId)
		if err != nil && !NotFoundError(err) {
			return WrapError(err)
		}

		equal := true
		if len(zone.BindVpcs.Vpc) == len(vpcIdMap) {
			for _, vpc := range zone.BindVpcs.Vpc {
				if _, ok := vpcIdMap[vpc.VpcId]; !ok {
					equal = false
					break
				}
			}
		}
		if equal {
			break
		}
		timeout = timeout - DefaultIntervalShort
		if timeout <= 0 {
			return WrapError(Error(GetTimeoutMessage("PrivateZone Attachment", "Ready")))
		}
		time.Sleep(DefaultIntervalShort * time.Second)
	}
	return nil
}

func (s *PvtzService) DescribeZoneRecord(recordId int, zoneId string) (record pvtz.Record, err error) {
	request := pvtz.CreateDescribeZoneRecordsRequest()
	request.ZoneId = zoneId
	request.PageNumber = requests.NewInteger(1)
	request.PageSize = requests.NewInteger(PageSizeLarge)

	recordIdStr := strconv.Itoa(recordId)

	invoker := PvtzInvoker()
	err = invoker.Run(func() error {
		for {
			raw, err := s.client.WithPvtzClient(func(pvtzClient *pvtz.Client) (interface{}, error) {
				return pvtzClient.DescribeZoneRecords(request)
			})

			if err != nil {
				if IsExceptedErrors(err, []string{ZoneNotExists, ZoneVpcNotExists}) {
					return GetNotFoundErrorFromString(GetNotFoundMessage("PrivateZoneRecord", recordIdStr))
				}
				return err
			}
			resp, _ := raw.(*pvtz.DescribeZoneRecordsResponse)
			if resp == nil || len(resp.Records.Record) < 1 {
				return GetNotFoundErrorFromString(GetNotFoundMessage("PrivateZoneRecord", recordIdStr))
			}

			for _, rec := range resp.Records.Record {
				if rec.RecordId == recordId {
					record = rec
					return nil
				}
			}
			if len(resp.Records.Record) < PageSizeLarge {
				break
			}

			if page, err := getNextpageNumber(request.PageNumber); err != nil {
				return err
			} else {
				request.PageNumber = page
			}
		}

		return GetNotFoundErrorFromString(GetNotFoundMessage("PrivateZoneRecord", recordIdStr))
	})

	return
}

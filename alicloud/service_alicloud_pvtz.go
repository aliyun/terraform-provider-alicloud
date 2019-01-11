package alicloud

import (
	"strconv"

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
				return GetNotFoundErrorFromString(GetNotFoundMessage("PrivateZone", zoneId))
			}
			return err
		}
		resp, _ := raw.(*pvtz.DescribeZoneInfoResponse)
		if resp == nil || resp.ZoneId != zoneId {
			return GetNotFoundErrorFromString(GetNotFoundMessage("PrivateZone", zoneId))
		}
		zone = *resp
		return nil
	})

	return

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
				if IsExceptedErrors(err, []string{ZoneNotExists}) {
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

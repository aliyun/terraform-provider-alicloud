package alicloud

import (
	"strconv"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/pvtz"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

type PvtzService struct {
	client *connectivity.AliyunClient
}

func (s *PvtzService) DescribePvtzZoneInfo(zoneId string) (zone pvtz.DescribeZoneInfoResponse, err error) {
	request := pvtz.CreateDescribeZoneInfoRequest()
	request.ZoneId = zoneId

	invoker := NewInvoker()
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

	invoker := NewInvoker()
	err = invoker.Run(func() error {
		raw, err := s.client.WithPvtzClient(func(pvtzClient *pvtz.Client) (interface{}, error) {
			return pvtzClient.DescribeZoneRecords(request)
		})

		recordIdStr := strconv.Itoa(recordId)

		if err != nil {
			if IsExceptedErrors(err, []string{ZoneNotExists}) {
				return GetNotFoundErrorFromString(GetNotFoundMessage("PrivateZoneRecord", recordIdStr))
			}
			return err
		}
		resp, _ := raw.(*pvtz.DescribeZoneRecordsResponse)
		if resp == nil {
			return GetNotFoundErrorFromString(GetNotFoundMessage("PrivateZoneRecord", recordIdStr))
		}

		var found bool
		for _, rec := range resp.Records.Record {
			if rec.RecordId == recordId {
				record = rec
				found = true
			}
		}

		if found == false {
			return GetNotFoundErrorFromString(GetNotFoundMessage("PrivateZoneRecord", recordIdStr))
		}

		return nil
	})

	return
}

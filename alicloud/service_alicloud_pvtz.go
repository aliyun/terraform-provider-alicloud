package alicloud

import (
	"strconv"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/pvtz"
)

func (client *AliyunClient) DescribePvtzZoneInfo(zoneId string) (zone pvtz.DescribeZoneInfoResponse, err error) {
	conn := client.pvtzconn
	request := pvtz.CreateDescribeZoneInfoRequest()
	request.ZoneId = zoneId

	invoker := NewInvoker()
	err = invoker.Run(func() error {
		resp, err := conn.DescribeZoneInfo(request)
		if err != nil {
			if IsExceptedErrors(err, []string{ZoneNotExists, ZoneVpcNotExists}) {
				return GetNotFoundErrorFromString(GetNotFoundMessage("PrivateZone", zoneId))
			}
			return err
		}
		if resp == nil || resp.ZoneId != zoneId {
			return GetNotFoundErrorFromString(GetNotFoundMessage("PrivateZone", zoneId))
		}
		zone = *resp
		return nil
	})

	return

}

func (client *AliyunClient) DescribeZoneRecord(recordId int, zoneId string) (record pvtz.Record, err error) {
	conn := client.pvtzconn
	request := pvtz.CreateDescribeZoneRecordsRequest()
	request.ZoneId = zoneId

	invoker := NewInvoker()
	err = invoker.Run(func() error {
		resp, err := conn.DescribeZoneRecords(request)

		recordIdStr := strconv.Itoa(recordId)

		if err != nil {
			if IsExceptedErrors(err, []string{ZoneNotExists}) {
				return GetNotFoundErrorFromString(GetNotFoundMessage("PrivateZoneRecord", recordIdStr))
			}
			return err
		}
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

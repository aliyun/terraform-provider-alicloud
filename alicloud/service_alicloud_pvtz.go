package alicloud

import (
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

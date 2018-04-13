package alicloud

import (
	"time"

	"strconv"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/cms"
)

func BuildCmsCommonRequest(region string) *requests.CommonRequest {

	request := requests.NewCommonRequest()

	return request
}

func (client *AliyunClient) BuildCmsAlarmRequest(id string) *requests.CommonRequest {

	request := BuildCmsCommonRequest(client.RegionId)
	request.QueryParams["Id"] = id

	return request
}

func (client *AliyunClient) DescribeAlarm(id string) (alarm cms.AlarmInListAlarm, err error) {

	request := cms.CreateListAlarmRequest()

	request.Id = id
	response, err := client.cmsconn.ListAlarm(request)
	if err != nil {
		return alarm, err
	}

	if len(response.AlarmList.Alarm) < 1 {
		return alarm, GetNotFoundErrorFromString(GetNotFoundMessage("Alarm Rule", id))
	}

	return response.AlarmList.Alarm[0], nil
}

func (client *AliyunClient) WaitForCmsAlarm(id string, enabled bool, timeout int) error {
	if timeout <= 0 {
		timeout = DefaultTimeout
	}

	for {
		alarm, err := client.DescribeAlarm(id)
		if err != nil {
			return err
		}

		if alarm.Enable == enabled {
			break
		}
		timeout = timeout - DefaultIntervalShort
		if timeout <= 0 {
			return GetTimeErrorFromString(GetTimeoutMessage("Alarm", strconv.FormatBool(enabled)))
		}
		time.Sleep(DefaultIntervalShort * time.Second)
	}
	return nil
}

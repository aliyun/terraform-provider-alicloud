package alicloud

import (
	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
)

type QuotasService struct {
	client *connectivity.AliyunClient
}

func (s *QuotasService) DescribeQuotasQuotaAlarm(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "GetQuotaAlarm"
	request := map[string]interface{}{
		"SourceIp": s.client.SourceIp,
		"RegionId": s.client.RegionId,
		"AlarmId":  id,
	}
	response, err = client.Do("quotas", rpcParam("POST", "2020-05-10", action), nil, request, nil, nil, false)
	if err != nil {
		err = WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
		return
	}
	addDebug(action, response, request)
	v, err := jsonpath.Get("$.QuotaAlarm", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.QuotaAlarm", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}

func (s *QuotasService) DescribeQuotasQuotaApplication(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "GetQuotaApplication"
	request := map[string]interface{}{
		"SourceIp":      s.client.SourceIp,
		"RegionId":      s.client.RegionId,
		"ApplicationId": id,
	}
	response, err = client.Do("quotas", rpcParam("POST", "2020-05-10", action), nil, request, nil, nil, false)
	if err != nil {
		err = WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
		return
	}
	addDebug(action, response, request)
	v, err := jsonpath.Get("$.QuotaApplication", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.QuotaApplication", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}

package alicloud

import (
	"github.com/PaesslerAG/jsonpath"
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	util "github.com/alibabacloud-go/tea-utils/service"
	utilv2 "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
)

type QuotasService struct {
	client *connectivity.AliyunClient
}

func (s *QuotasService) DescribeQuotasQuotaAlarm(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewQuotasClientV2()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "GetQuotaAlarm"
	request := map[string]interface{}{
		"SourceIp": s.client.SourceIp,
		"RegionId": s.client.RegionId,
		"AlarmId":  id,
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	response, err = conn.CallApi(rpcParam(action, "POST", "2020-05-10"), &openapi.OpenApiRequest{Query: nil, Body: request, HostMap: nil}, &utilv2.RuntimeOptions{})
	if err != nil {
		err = WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
		return
	}
	response = response["body"].(map[string]interface{})
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
	conn, err := s.client.NewQuotasClientV2()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "GetQuotaApplication"
	request := map[string]interface{}{
		"SourceIp":      s.client.SourceIp,
		"RegionId":      s.client.RegionId,
		"ApplicationId": id,
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	response, err = conn.CallApi(rpcParam(action, "POST", "2020-05-10"), &openapi.OpenApiRequest{Query: nil, Body: request, HostMap: nil}, &utilv2.RuntimeOptions{})
	if err != nil {
		err = WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
		return
	}
	response = response["body"].(map[string]interface{})
	addDebug(action, response, request)
	v, err := jsonpath.Get("$.QuotaApplication", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.QuotaApplication", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}

package alicloud

import (
	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
)

type QuotasService struct {
	client *connectivity.AliyunClient
}

func (s *QuotasService) DescribeQuotasApplicationInfo(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewQuotasClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "GetQuotaApplication"
	request := map[string]interface{}{
		"RegionId":      s.client.RegionId,
		"ApplicationId": id,
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-05-10"), StringPointer("AK"), nil, request, &runtime)
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

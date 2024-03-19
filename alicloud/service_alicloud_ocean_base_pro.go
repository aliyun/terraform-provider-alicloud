package alicloud

import (
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

type OceanBaseProService struct {
	client *connectivity.AliyunClient
}

func (s *OceanBaseProService) DescribeOceanBaseInstance(id string) (object map[string]interface{}, err error) {
	conn, err := s.client.NewOceanbaseClient()
	if err != nil {
		return object, WrapError(err)
	}

	request := map[string]interface{}{
		"InstanceId": id,
		"RegionId":   s.client.RegionId,
	}

	var response map[string]interface{}
	action := "DescribeInstance"
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		resp, err := conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-09-01"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		response = resp
		addDebug(action, response, request)
		return nil
	})
	if err != nil {
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	if _, ok := response["Instance"]; !ok {
		return object, WrapErrorf(Error(GetNotFoundMessage("Instance", id)), NotFoundWithResponse, response)
	}
	v, err := jsonpath.Get("$.Instance", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Instance", response)
	}
	return v.(map[string]interface{}), nil
}

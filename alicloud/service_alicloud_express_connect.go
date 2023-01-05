package alicloud

import (
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

type ExpressConnectService struct {
	client *connectivity.AliyunClient
}

func (s *ExpressConnectService) DescribeExpressConnectVirtualPhysicalConnection(id string) (object map[string]interface{}, err error) {
	conn, err := s.client.NewVpcClient()
	if err != nil {
		return object, WrapError(err)
	}

	request := map[string]interface{}{
		"VirtualPhysicalConnectionIds.1": id,
		"RegionId":                       s.client.RegionId,
	}

	var response map[string]interface{}
	action := "ListVirtualPhysicalConnections"
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		resp, err := conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, request, &runtime)
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
	v, err := jsonpath.Get("$.VirtualPhysicalConnections", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.VirtualPhysicalConnections", response)
	}
	if len(v.([]interface{})) < 1 {
		return object, WrapErrorf(Error(GetNotFoundMessage("VirtualPhysicalConnections", id)), NotFoundWithResponse, response)
	}
	return v.([]interface{})[0].(map[string]interface{}), nil
}

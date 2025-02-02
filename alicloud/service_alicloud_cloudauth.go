package alicloud

import (
	"fmt"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

type CloudauthService struct {
	client *connectivity.AliyunClient
}

func (s *CloudauthService) DescribeCloudauthFaceConfig(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "DescribeFaceConfig"
	request := map[string]interface{}{
		"Lang": id,
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("Cloudauth", "2019-03-07", action, nil, request, true)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.Items", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Items", response)
	}
	if len(v.([]interface{})) < 1 {
		return object, WrapErrorf(Error(GetNotFoundMessage("Cloudauth", id)), NotFoundWithResponse, response)
	} else {
		for _, obj := range v.([]interface{}) {
			if fmt.Sprint(obj.(map[string]interface{})["BizType"]) == id {
				object = v.([]interface{})[0].(map[string]interface{})
				return object, nil
			}
		}
		return object, WrapErrorf(Error(GetNotFoundMessage("Cloudauth", id)), NotFoundWithResponse, response)
	}
}

package alicloud

import (
	"fmt"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

type FcOpenService struct {
	client *connectivity.AliyunClient
}

func (s *FcOpenService) DescribeFcLayerVersion(id string) (object map[string]interface{}, err error) {
	conn, err := s.client.NewFcClient()
	if err != nil {
		return object, WrapError(err)
	}
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		return object, WrapError(err)
	}

	var response map[string]interface{}
	action := fmt.Sprintf("/2021-04-06/layers/%s/versions/%s", parts[0], parts[1])
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer("2021-04-06"), nil, StringPointer("GET"), StringPointer("AK"), StringPointer(action), nil, nil, nil, &util.RuntimeOptions{})
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response)
		return nil
	})
	if err != nil {
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.body", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.body", response)
	}

	return v.(map[string]interface{}), nil
}

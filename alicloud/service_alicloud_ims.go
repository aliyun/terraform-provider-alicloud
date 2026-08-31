package alicloud

import (
	"fmt"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

type ImsService struct {
	client *connectivity.AliyunClient
}

func (s *ImsService) DescribeRamSamlProvider(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	action := "GetSAMLProvider"
	client := s.client
	request := map[string]interface{}{
		"SAMLProviderName": id,
	}
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("Ims", "2019-08-15", action, nil, request, true)
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
		if IsExpectedErrors(err, []string{"EntityNotExist.SAMLProviderError"}) {
			return object, WrapErrorf(NotFoundErr("Ram:SamlProvider", id), NotFoundMsg, ProviderERROR, fmt.Sprint(response["RequestId"]))
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$.SAMLProvider", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.SAMLProvider", response)
	}

	object = v.(map[string]interface{})

	return object, nil
}

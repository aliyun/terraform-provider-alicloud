package alicloud

import (
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

type MhubService struct {
	client *connectivity.AliyunClient
}

func (s *MhubService) DescribeMhubProduct(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewMhubClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "QueryProductInfo"
	request := map[string]interface{}{
		"ProductId": id,
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-08-25"), StringPointer("AK"), nil, request, &runtime)
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
		if IsExpectedErrors(err, []string{"ProductNotExist"}) {
			err = WrapErrorf(Error(GetNotFoundMessage("MHUBProduct", id)), NotFoundMsg, ProviderERROR)
			return object, err
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$.ProductInfo", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.ProductInfo", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}

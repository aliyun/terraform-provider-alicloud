package alicloud

import (
	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"time"
)

type IaCServiceService struct {
	client *connectivity.AliyunClient
}

func (s *IaCServiceService) GetKeyPairResource(id string) (object map[string]interface{}, err error) {

	var response map[string]interface{}
	conn, err := s.client.NewIaCServiceClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "/api/v1/providers/Aliyun/products/ECP/resourceTypes/KeyPair/resources/" + id
	request := map[string]*string{
		"regionId": StringPointer(s.client.RegionId),
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer("2021-07-22"), nil, StringPointer("GET"), StringPointer("AK"), StringPointer(action), request, nil, nil, &util.RuntimeOptions{})

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
		if IsExpectedErrors(err, []string{"InvalidResource.NotFound"}) {
			return response, WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$.body.resource.resourceAttributes", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.body.resource.resourceAttributes", response)
	}

	object, _ = convertJsonStringToMap(v.(string))
	return object, nil
}

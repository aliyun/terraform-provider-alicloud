package alicloud

import (
	"fmt"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

type TagService struct {
	client *connectivity.AliyunClient
}

func (s *TagService) ListTagValues(key string) (object []interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewTagClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "ListTagValues"
	request := map[string]interface{}{
		"RegionId":  s.client.RegionId,
		"Key":       key,
		"QueryType": "MetaTag",
	}
	values := make([]interface{}, 0)
	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2018-08-28"), StringPointer("AK"), nil, request, &runtime)
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
			if IsExpectedErrors(err, []string{"MetaTagKeyNotFound"}) {
				return object, WrapErrorf(Error(GetNotFoundMessage("TAG:MetaTag", key)), NotFoundMsg, ProviderERROR, fmt.Sprint(response["RequestId"]))
			}
			return object, WrapErrorf(err, DefaultErrorMsg, key, action, AlibabaCloudSdkGoERROR)
		}
		v, err := jsonpath.Get("$.Values.Value", response)
		if err != nil {
			return object, WrapErrorf(err, FailedGetAttributeMsg, key, "$.Values.Value", response)
		}
		values = append(values, v.([]interface{})...)
		if nextToken, ok := response["NextToken"].(string); ok && nextToken != "" {
			request["NextToken"] = nextToken
		} else {
			break
		}
		request["NextToken"] = response["NextToken"]
	}

	return values, nil
}

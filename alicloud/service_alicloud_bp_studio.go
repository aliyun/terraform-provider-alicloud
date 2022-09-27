package alicloud

import (
	"fmt"
	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"time"
)

type BpStudioService struct {
	client *connectivity.AliyunClient
}

func (s *BpStudioService) DescribeBpStudioApplication(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewBpstudioClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "GetApplication"
	request := map[string]interface{}{
		"ApplicationId": id,
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2021-09-31"), StringPointer("AK"), nil, request, &runtime)
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

	if fmt.Sprint(response["Code"]) == "8004" {
		return object, WrapErrorf(Error(GetNotFoundMessage("BPStudio:Application", id)), NotFoundMsg, ProviderERROR)
	}
	v, err := jsonpath.Get("$.Data", response)

	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Data", response)
	}

	object = v.(map[string]interface{})
	if len(object) < 1 {
		return object, WrapErrorf(Error(GetNotFoundMessage("BPStudio:Application", id)), NotFoundWithResponse, response)
	}

	return object, nil
}

func (s *BpStudioService) BpStudioApplicationStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeBpStudioApplication(id)
		if err != nil {
			if NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}
		for _, failState := range failStates {
			if fmt.Sprint(object["Status"]) == failState {
				return object, fmt.Sprint(object["Status"]), WrapError(Error(FailedToReachTargetStatus, fmt.Sprint(object["Status"])))
			}
		}
		return object, fmt.Sprint(object["Status"]), nil
	}
}

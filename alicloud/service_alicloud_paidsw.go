package alicloud

import (
	"fmt"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

type PaiDswService struct {
	client *connectivity.AliyunClient
}

// DescribePaiDswTempFileTask Encapsulated get interface for PaiDsw TempFileTask.
func (s *PaiDswService) DescribePaiDswTempFileTask(id string) (object map[string]interface{}, err error) {
	client := s.client
	var response map[string]interface{}
	var query map[string]*string
	tempFileTaskId := id
	query = make(map[string]*string)
	query["RegionId"] = StringPointer(client.RegionId)

	action := fmt.Sprintf("/api/v2/tempfiletasks/%s", tempFileTaskId)

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(1*time.Minute, func() *resource.RetryError {
		response, err = client.RoaGet("pai-dsw", "2022-01-01", action, query, nil, nil)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, query)
	if err != nil {
		if IsExpectedErrors(err, []string{"ValidationError"}) {
			return object, WrapErrorf(NotFoundErr("PaiDswTempFileTask", id), NotFoundMsg, err)
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	code, _ := jsonpath.Get("$.HttpStatusCode", response)
	if InArray(fmt.Sprint(code), []string{"400"}) {
		return object, WrapErrorf(NotFoundErr("PaiDswTempFileTask", id), NotFoundMsg, response)
	}

	v, err := jsonpath.Get("$", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$", response)
	}
	if v == nil {
		return object, WrapErrorf(NotFoundErr("PaiDswTempFileTask", id), NotFoundMsg, response)
	}
	object = v.(map[string]interface{})
	return object, nil
}

// PaiDswTempFileTaskStateRefreshFunc ...
func (s *PaiDswService) PaiDswTempFileTaskStateRefreshFunc(id string, field string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribePaiDswTempFileTask(id)
		if err != nil {
			if NotFoundError(err) {
				return 0, "", nil
			}
			return nil, "", err
		}
		if value, ok := object[field]; ok {
			return value, "", nil
		}
		return object, "", nil
	}
}

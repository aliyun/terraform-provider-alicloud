package alicloud

import (
	"fmt"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

type Schedulerx2Service struct {
	client *connectivity.AliyunClient
}

func (s *Schedulerx2Service) DescribeSchedulerxNamespace(id string) (object map[string]interface{}, err error) {
	conn, err := s.client.NewEdasschedulerxClient()
	if err != nil {
		return object, WrapError(err)
	}

	request := map[string]interface{}{
		"RegionId": s.client.RegionId,
	}

	var response map[string]interface{}
	action := "ListNamespaces"
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		resp, err := conn.DoRequest(StringPointer(action), nil, StringPointer("GET"), StringPointer("2019-04-30"), StringPointer("AK"), request, nil, &runtime)
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
	v, err := jsonpath.Get("$.Data.Namespaces", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Data.Namespaces", response)
	}
	if len(v.([]interface{})) < 1 {
		return object, WrapErrorf(Error(GetNotFoundMessage("Namespace", id)), NotFoundWithResponse, response)
	} else {
		for _, vv := range v.([]interface{}) {
			item := vv.(map[string]interface{})
			if fmt.Sprint(item["UId"]) == id {
				return item, nil
			}
		}
	}
	return object, WrapErrorf(Error(GetNotFoundMessage("Namespace", id)), NotFoundWithResponse, response)
}

package alicloud

import (
	"fmt"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

type ResourceManagerService struct {
	client *connectivity.AliyunClient
}

func (s *ResourceManagerService) DescribeResourceManagerDelegatedAdministrator(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewResourcemanagerClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "ListDelegatedAdministrators"
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		err = WrapError(err)
		return
	}
	request := map[string]interface{}{
		"ServicePrincipal": parts[1],
		"PageNumber":       1,
		"PageSize":         PageSizeLarge,
	}
	idExist := false
	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-03-31"), StringPointer("AK"), nil, request, &runtime)
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
		v, err := jsonpath.Get("$.Accounts.Account", response)
		if err != nil {
			return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Accounts.Account", response)
		}
		if len(v.([]interface{})) < 1 {
			return object, WrapErrorf(Error(GetNotFoundMessage("ResourceManager", id)), NotFoundWithResponse, response)
		}
		for _, v := range v.([]interface{}) {
			if fmt.Sprint(v.(map[string]interface{})["AccountId"]) == parts[0] {
				idExist = true
				return v.(map[string]interface{}), nil
			}
		}
		if len(v.([]interface{})) < request["PageSize"].(int) {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}
	if !idExist {
		return object, WrapErrorf(Error(GetNotFoundMessage("ResourceManager", id)), NotFoundWithResponse, response)
	}
	return
}

package alicloud

import (
	"fmt"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

type MaxComputeService struct {
	client *connectivity.AliyunClient
}

func (s *MaxComputeService) DescribeMaxcomputeProject(id string) (object map[string]interface{}, err error) {
	client := s.client

	var response map[string]interface{}
	action := "/api/v1/projects/" + id
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		resp, err := client.RoaGet("MaxCompute", "2022-01-04", action, nil, nil, nil)
		if err != nil {
			return resource.NonRetryableError(err)
		}
		response = resp
		addDebug(action, response)
		return nil
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"INTERNAL_SERVER_ERROR", "OBJECT_NOT_EXIST"}) {
			return object, WrapErrorf(Error(GetNotFoundMessage("MaxCompute", id)), NotFoundWithResponse, response)
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.data", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.data", response)
	}
	status, err := jsonpath.Get("$.status", v)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.status", response)
	}
	if status == "DELETING" {
		return object, WrapErrorf(Error(GetNotFoundMessage("MaxCompute:Project", id)), NotFoundWithResponse, response)
	}
	return v.(map[string]interface{}), nil
}

func (s *MaxComputeService) MaxcomputeProjectStateRefreshFunc(d *schema.ResourceData, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeMaxcomputeProject(d.Id())
		if err != nil {
			if NotFoundError(err) {
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}
		for _, failState := range failStates {
			if fmt.Sprint(object["status"]) == failState {
				return object, fmt.Sprint(object["status"]), WrapError(Error(FailedToReachTargetStatus, fmt.Sprint(object["status"])))
			}
		}
		return object, fmt.Sprint(object["status"]), nil
	}
}

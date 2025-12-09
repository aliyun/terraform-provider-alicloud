// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

type OpenApiExplorerServiceV2 struct {
	client *connectivity.AliyunClient
}

// DescribeOpenApiExplorerApiMcpServer <<< Encapsulated get interface for OpenApiExplorer ApiMcpServer.

func (s *OpenApiExplorerServiceV2) DescribeOpenApiExplorerApiMcpServer(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	var header map[string]*string
	request = make(map[string]interface{})
	query = make(map[string]*string)
	header = make(map[string]*string)
	query["id"] = StringPointer(id)

	action := fmt.Sprintf("/apimcpserver")

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(1*time.Minute, func() *resource.RetryError {
		response, err = client.RoaGet("OpenAPIExplorer", "2024-11-30", action, query, header, nil)

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
		if IsExpectedErrors(err, []string{"InvalidApiMcpServer.NotFound"}) {
			return object, WrapErrorf(NotFoundErr("ApiMcpServer", id), NotFoundMsg, response)
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}

	return response, nil
}

func (s *OpenApiExplorerServiceV2) OpenApiExplorerApiMcpServerStateRefreshFunc(id string, field string, failStates []string) resource.StateRefreshFunc {
	return s.OpenApiExplorerApiMcpServerStateRefreshFuncWithApi(id, field, failStates, s.DescribeOpenApiExplorerApiMcpServer)
}

func (s *OpenApiExplorerServiceV2) OpenApiExplorerApiMcpServerStateRefreshFuncWithApi(id string, field string, failStates []string, call func(id string) (map[string]interface{}, error)) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := call(id)
		if err != nil {
			if NotFoundError(err) {
				return object, "", nil
			}
			return nil, "", WrapError(err)
		}
		v, err := jsonpath.Get(field, object)
		currentStatus := fmt.Sprint(v)

		if strings.HasPrefix(field, "#") {
			v, _ := jsonpath.Get(strings.TrimPrefix(field, "#"), object)
			if v != nil {
				currentStatus = "#CHECKSET"
			}
		}

		for _, failState := range failStates {
			if currentStatus == failState {
				return object, currentStatus, WrapError(Error(FailedToReachTargetStatus, currentStatus))
			}
		}
		return object, currentStatus, nil
	}
}

// DescribeOpenApiExplorerApiMcpServer >>> Encapsulated.

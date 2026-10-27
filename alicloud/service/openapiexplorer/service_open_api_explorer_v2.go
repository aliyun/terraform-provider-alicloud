// Package openapiexplorer. This file is generated automatically. Please do not modify it manually, thank you!
package openapiexplorer

import (
	"fmt"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

type OpenApiExplorerServiceV2 struct {
	client *connectivity.AliyunClient
}

// NewOpenApiExplorerServiceV2 is the exported constructor used by the
// package's external tests (package openapiexplorer_test); production wiring
// reaches the service through conns.ServicePackage.
func NewOpenApiExplorerServiceV2(client *connectivity.AliyunClient) *OpenApiExplorerServiceV2 {
	return &OpenApiExplorerServiceV2{client: client}
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
	query["id"] = helper.StringPointer(id)

	action := fmt.Sprintf("/apimcpserver")

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.RoaGet("OpenAPIExplorer", "2024-11-30", action, query, header, nil)

		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"InvalidApiMcpServer.NotFound"}) {
			return object, errs.WrapErrorf(errs.NotFoundErr("ApiMcpServer", id), errs.NotFoundMsg, response)
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}

	return response, nil
}

func (s *OpenApiExplorerServiceV2) OpenApiExplorerApiMcpServerStateRefreshFunc(id string, field string, failStates []string) retry.StateRefreshFunc {
	return s.OpenApiExplorerApiMcpServerStateRefreshFuncWithApi(id, field, failStates, s.DescribeOpenApiExplorerApiMcpServer)
}

func (s *OpenApiExplorerServiceV2) OpenApiExplorerApiMcpServerStateRefreshFuncWithApi(id string, field string, failStates []string, call func(id string) (map[string]interface{}, error)) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := call(id)
		if err != nil {
			if errs.NotFoundError(err) {
				return object, "", nil
			}
			return nil, "", errs.WrapError(err)
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
				return object, currentStatus, errs.WrapError(errs.Error(errs.FailedToReachTargetStatus, currentStatus))
			}
		}
		return object, currentStatus, nil
	}
}

// DescribeOpenApiExplorerApiMcpServer >>> Encapsulated.

package servicemesh

import (
	"fmt"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

type ServicemeshService struct {
	client *connectivity.AliyunClient
}

// NewServicemeshService constructs a ServicemeshService for the SDKv2
// servicemesh sub-package. Exported so the external servicemesh_test package
// can build the service without referencing the unexported client field; the
// two service structs are kept separate because they expose colliding method
// names (DescribeServiceMeshServiceMesh / ServiceMeshServiceMeshStateRefreshFunc).
func NewServicemeshService(client *connectivity.AliyunClient) *ServicemeshService {
	return &ServicemeshService{client}
}

func (s *ServicemeshService) DescribeServiceMeshServiceMesh(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "DescribeServiceMeshDetail"
	request := map[string]interface{}{
		"ServiceMeshId": id,
	}
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("servicemesh", "2020-01-11", action, nil, request, true)
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
		if errs.IsExpectedErrors(err, []string{"ServiceMesh.NotFound", "StatusForbidden"}) {
			return object, errs.WrapErrorf(errs.NotFoundErr("ServiceMesh:ServiceMesh", id), errs.NotFoundMsg, errs.ProviderERROR, fmt.Sprint(response["RequestId"]))
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.ServiceMesh", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.ServiceMesh", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}

func (s *ServicemeshService) DescribeServiceMeshDetail(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "DescribeServiceMeshDetail"
	request := map[string]interface{}{
		"ServiceMeshId": id,
	}
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("servicemesh", "2020-01-11", action, nil, request, true)
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
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.ServiceMesh", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.ServiceMesh", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}

func (s *ServicemeshService) ServiceMeshServiceMeshStateRefreshFunc(id string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeServiceMeshServiceMesh(id)
		if err != nil {
			if errs.NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", errs.WrapError(err)
		}

		for _, failState := range failStates {
			if fmt.Sprint(object["ServiceMeshInfo"].(map[string]interface{})["State"]) == failState {
				return object, fmt.Sprint(object["ServiceMeshInfo"].(map[string]interface{})["State"]), errs.WrapError(errs.Error(errs.FailedToReachTargetStatus, fmt.Sprint(object["ServiceMeshInfo.State"])))
			}
		}
		return object, fmt.Sprint(object["ServiceMeshInfo"].(map[string]interface{})["State"]), nil
	}
}

func (s *ServicemeshService) DescribeUserPermissions(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "DescribeUserPermissions"
	request := map[string]interface{}{
		"SubAccountUserId": id,
	}
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("servicemesh", "2020-01-11", action, nil, request, true)
		if err != nil {
			if helper.NeedRetry(err) || errs.IsExpectedErrors(err, []string{"StatusForbidden"}) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$", response)
	}
	object = v.(map[string]interface{})
	return object, nil

}

func (s *ServicemeshService) DescribeServiceMeshExtensionProvider(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	parts, err := helper.ParseResourceId(id, 3)
	if err != nil {
		return object, errs.WrapError(err)
	}
	action := "DescribeExtensionProvider"
	request := map[string]interface{}{
		"ServiceMeshId": parts[0],
		"Type":          parts[1],
	}

	idExist := false
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("servicemesh", "2020-01-11", action, nil, request, true)
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
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}

	resp, err := jsonpath.Get("$.ExtensionProviderArray", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.ExtensionProviderArray", response)
	}

	if v, ok := resp.([]interface{}); !ok || len(v) < 1 {
		return object, errs.WrapErrorf(errs.NotFoundErr("ServiceMesh:ExtensionProvider", id), errs.NotFoundWithResponse, response)
	}

	for _, v := range resp.([]interface{}) {
		if fmt.Sprint(v.(map[string]interface{})["Name"]) == parts[2] {
			idExist = true
			return v.(map[string]interface{}), nil
		}
	}

	if !idExist {
		return object, errs.WrapErrorf(errs.NotFoundErr("ServiceMesh:ExtensionProvider", id), errs.NotFoundWithResponse, response)
	}

	return object, nil
}

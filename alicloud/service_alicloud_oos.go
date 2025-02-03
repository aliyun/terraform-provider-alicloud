package alicloud

import (
	"fmt"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

type OosService struct {
	client *connectivity.AliyunClient
}

func (s *OosService) DescribeOosTemplate(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "GetTemplate"
	request := map[string]interface{}{
		"RegionId":     s.client.RegionId,
		"TemplateName": id,
	}
	response, err = client.RpcPost("oos", "2019-06-01", action, nil, request, true)
	if err != nil {
		if IsExpectedErrors(err, []string{"EntityNotExists.Template"}) {
			err = WrapErrorf(Error(GetNotFoundMessage("OosTemplate", id)), NotFoundMsg, ProviderERROR)
			return object, err
		}
		err = WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
		return object, err
	}
	addDebug(action, response, request)
	v, err := jsonpath.Get("$.Template", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Template", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}

func (s *OosService) DescribeOosExecution(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "ListExecutions"
	request := map[string]interface{}{
		"RegionId":    s.client.RegionId,
		"ExecutionId": id,
	}
	response, err = client.RpcPost("oos", "2019-06-01", action, nil, request, true)
	if err != nil {
		err = WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
		return
	}
	addDebug(action, response, request)
	v, err := jsonpath.Get("$.Executions", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Executions", response)
	}
	if len(v.([]interface{})) < 1 {
		return object, WrapErrorf(Error(GetNotFoundMessage("OOS", id)), NotFoundWithResponse, response)
	} else {
		if v.([]interface{})[0].(map[string]interface{})["ExecutionId"].(string) != id {
			return object, WrapErrorf(Error(GetNotFoundMessage("OOS", id)), NotFoundWithResponse, response)
		}
	}
	object = v.([]interface{})[0].(map[string]interface{})
	return object, nil
}

func (s *OosService) OosExecutionStateRefreshFunc(id string, field string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeOosExecution(id)
		if err != nil {
			if NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}

		currentStatus := fmt.Sprint(object[field])
		for _, failState := range failStates {
			if currentStatus == failState {
				return object, currentStatus, WrapError(Error(FailedToReachTargetStatus, currentStatus))
			}
		}
		return object, currentStatus, nil
	}
}

func (s *OosService) DescribeOosApplication(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "GetApplication"
	request := map[string]interface{}{
		"RegionId": s.client.RegionId,
		"Name":     id,
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("oos", "2019-06-01", action, nil, request, true)
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
		if IsExpectedErrors(err, []string{"EntityNotExists.Application"}) {
			return object, WrapErrorf(Error(GetNotFoundMessage("Oos:Application", id)), NotFoundMsg, ProviderERROR)
		}

		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.Application", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Application", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}

func (s *OosService) DescribeOosApplicationGroup(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "GetApplicationGroup"
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		err = WrapError(err)
		return
	}
	request := map[string]interface{}{
		"RegionId":        s.client.RegionId,
		"Name":            parts[1],
		"ApplicationName": parts[0],
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("oos", "2019-06-01", action, nil, request, true)
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
		if IsExpectedErrors(err, []string{"EntityNotExists.ApplicationGroup"}) {
			return object, WrapErrorf(Error(GetNotFoundMessage("Oos:ApplicationGroup", id)), NotFoundMsg, ProviderERROR)
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.ApplicationGroup", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.ApplicationGroup", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}

func (s *OosService) DescribeOosPatchBaseline(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "GetPatchBaseline"
	request := map[string]interface{}{
		"Name": id,
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("oos", "2019-06-01", action, nil, request, true)
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
		if IsExpectedErrors(err, []string{"EntityNotExists.PatchBaseline"}) {
			return object, WrapErrorf(Error(GetNotFoundMessage("Oos:Application", id)), NotFoundMsg, ProviderERROR)
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.PatchBaseline", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.PatchBaseline", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}
func (s *OosService) DescribeOosStateConfiguration(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "ListStateConfigurations"
	request := map[string]interface{}{
		"RegionId": s.client.RegionId,
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("oos", "2019-06-01", action, nil, request, true)
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
	v, err := jsonpath.Get("$.StateConfigurations", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.StateConfigurations", response)
	}
	if len(v.([]interface{})) < 1 {
		return object, WrapErrorf(Error(GetNotFoundMessage("OOS", id)), NotFoundWithResponse, response)
	} else {
		if fmt.Sprint(v.([]interface{})[0].(map[string]interface{})["StateConfigurationId"]) != id {
			return object, WrapErrorf(Error(GetNotFoundMessage("OOS", id)), NotFoundWithResponse, response)
		}
	}
	object = v.([]interface{})[0].(map[string]interface{})
	return object, nil
}

func (s *OosService) DescribeOosServiceSetting(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "GetServiceSettings"
	request := map[string]interface{}{
		"RegionId": s.client.RegionId,
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("oos", "2019-06-01", action, nil, request, true)
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
	v, err := jsonpath.Get("$.ServiceSettings", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.ServiceSettings", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}

func (s *OosService) DescribeOosParameter(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "GetParameter"
	request := map[string]interface{}{
		"Name": id,
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("oos", "2019-06-01", action, nil, request, true)
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
		if IsExpectedErrors(err, []string{"EntityNotExists.Parameter"}) {
			return object, WrapErrorf(Error(GetNotFoundMessage("Oos:Application", id)), NotFoundMsg, ProviderERROR)
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.Parameter", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Parameter", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}

func (s *OosService) DescribeOosSecretParameter(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	action := "GetSecretParameter"

	client := s.client

	request := map[string]interface{}{
		"RegionId": s.client.RegionId,
		"Name":     id,
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("oos", "2019-06-01", action, nil, request, true)
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
		if IsExpectedErrors(err, []string{"EntityNotExists.Parameter"}) {
			return object, WrapErrorf(Error(GetNotFoundMessage("Oos:SecretParameter", id)), NotFoundWithResponse, response)
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$.Parameter", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Parameter", response)
	}

	object = v.(map[string]interface{})

	return object, nil
}

func (s *OosService) DescribeOosDefaultPatchBaseline(id string) (object map[string]interface{}, err error) {
	client := s.client

	request := map[string]interface{}{
		"Name":     id,
		"RegionId": s.client.RegionId,
	}

	var response map[string]interface{}
	action := "GetPatchBaseline"
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		resp, err := client.RpcPost("oos", "2019-06-01", action, nil, request, true)
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
	v, err := jsonpath.Get("$.PatchBaseline", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.PatchBaseline", response)
	}
	if fmt.Sprint(v.(map[string]interface{})["IsDefault"]) != "true" {
		return response, WrapErrorf(Error(GetNotFoundMessage("OosDefaultPatchBaseline", id)), NotFoundMsg, AlibabaCloudSdkGoERROR)
	}
	return v.(map[string]interface{}), nil
}

func (s *OosService) DescribeOosSecretParameterForDataSource(id string, d *schema.ResourceData) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	action := "GetSecretParameter"
	client := s.client
	request := map[string]interface{}{
		"RegionId": s.client.RegionId,
		"Name":     id,
	}

	if v, ok := d.GetOkExists("with_decryption"); ok {
		request["WithDecryption"] = v
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("oos", "2019-06-01", action, nil, request, true)
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
		if IsExpectedErrors(err, []string{"EntityNotExists.Parameter"}) {
			return object, WrapErrorf(Error(GetNotFoundMessage("Oos:SecretParameter", id)), NotFoundWithResponse, response)
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$.Parameter", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Parameter", response)
	}

	object = v.(map[string]interface{})

	return object, nil
}

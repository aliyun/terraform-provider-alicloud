package scdn

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Service wraps the SCDN API calls used by the resources and data sources of
// this package. The root alicloud package cannot be imported from a service
// package, so the encapsulated calls live here with the resources.
type Service struct {
	client *connectivity.AliyunClient
}

// NewService is the exported constructor: the client field stays unexported,
// so callers outside the package (including the acceptance tests) build
// values with it.
func NewService(client *connectivity.AliyunClient) *Service {
	return &Service{client: client}
}

func (s *Service) DescribeScdnDomainCertificateInfo(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "DescribeScdnDomainCertificateInfo"
	request := map[string]interface{}{
		"DomainName": id,
	}
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("scdn", "2017-11-15", action, nil, request, true)
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
	v, err := jsonpath.Get("$.CertInfos.CertInfo", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.CertInfos.CertInfo", response)
	}
	if len(v.([]interface{})) < 1 {
		return object, errs.WrapErrorf(errs.NotFoundErr("SCDN", id), errs.NotFoundWithResponse, response)
	} else {
		if fmt.Sprint(v.([]interface{})[0].(map[string]interface{})["DomainName"]) != id {
			return object, errs.WrapErrorf(errs.NotFoundErr("SCDN", id), errs.NotFoundWithResponse, response)
		}
	}
	object = v.([]interface{})[0].(map[string]interface{})
	return object, nil
}

func (s *Service) convertFunctionsToString(v []interface{}) (string, error) {
	arrayMaps := make([]interface{}, len(v))
	for i, vv := range v {
		item := vv.(map[string]interface{})
		temp := map[string]interface{}{
			"ConfigId":     item["config_id"],
			"FunctionName": item["function_name"],
			"Status":       item["status"],
		}

		functionArgs := item["function_args"].(*schema.Set).List()
		functionArgsMaps := make([]interface{}, len(functionArgs))
		for j, functionArgsVal := range functionArgs {
			functionArgsItem := functionArgsVal.(map[string]interface{})
			functionArgsTemp := map[string]interface{}{
				"ArgName":  functionArgsItem["arg_name"],
				"ArgValue": functionArgsItem["arg_value"],
			}

			functionArgsMaps[j] = functionArgsTemp
		}
		temp["function_args"] = functionArgsMaps
		arrayMaps[i] = temp
	}
	maps, err := json.Marshal(arrayMaps)
	if err != nil {
		return "", errs.WrapError(err)
	}
	return string(maps), nil
}

func (s *Service) DescribeScdnDomain(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "DescribeScdnDomainDetail"
	request := map[string]interface{}{
		"DomainName": id,
	}
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("scdn", "2017-11-15", action, nil, request, true)
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
		if errs.IsExpectedErrors(err, []string{"InvalidDomain.NotFound"}) {
			return object, errs.WrapErrorf(errs.NotFoundErr("SCDN:Domain", id), errs.NotFoundMsg, errs.ProviderERROR, fmt.Sprint(response["RequestId"]))
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.DomainDetail", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.DomainDetail", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}

func (s *Service) ScdnDomainStateRefreshFunc(id string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeScdnDomain(id)
		if err != nil {
			if errs.NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", errs.WrapError(err)
		}

		for _, failState := range failStates {
			if fmt.Sprint(object["DomainStatus"]) == failState {
				return object, fmt.Sprint(object["DomainStatus"]), errs.WrapError(errs.Error(errs.FailedToReachTargetStatus, fmt.Sprint(object["DomainStatus"])))
			}
		}
		return object, fmt.Sprint(object["DomainStatus"]), nil
	}
}

func (s *Service) DescribeScdnCertificateList(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "DescribeScdnCertificateList"
	request := map[string]interface{}{
		"DomainName": id,
	}
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("scdn", "2017-11-15", action, nil, request, true)
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
	v, err := jsonpath.Get("$.CertificateListModel.CertList.Cert", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.CertificateListModel.CertList.Cert", response)
	}
	if len(v.([]interface{})) < 1 {
		return object, errs.WrapErrorf(errs.NotFoundErr("SCDN", id), errs.NotFoundWithResponse, response)
	} else {
		if fmt.Sprint(v.([]interface{})[0].(map[string]interface{})["DomainName"]) != id {
			return object, errs.WrapErrorf(errs.NotFoundErr("SCDN", id), errs.NotFoundWithResponse, response)
		}
	}
	object = v.([]interface{})[0].(map[string]interface{})
	return object, nil
}

func (s *Service) DescribeScdnDomainCname(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "DescribeScdnDomainCname"
	request := map[string]interface{}{
		"DomainName": id,
	}
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcGet("scdn", "2017-11-15", action, request, nil)
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
	v, err := jsonpath.Get("$.CnameDatas.Data", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.CnameDatas.Data", response)
	}
	if len(v.([]interface{})) < 1 {
		return object, errs.WrapErrorf(errs.NotFoundErr("SCDN", id), errs.NotFoundWithResponse, response)
	} else {
		if fmt.Sprint(v.([]interface{})[0].(map[string]interface{})["DomainName"]) != id {
			return object, errs.WrapErrorf(errs.NotFoundErr("SCDN", id), errs.NotFoundWithResponse, response)
		}
	}
	object = v.([]interface{})[0].(map[string]interface{})
	return object, nil
}

func (s *Service) DescribeScdnDomainConfigs(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "DescribeScdnDomainConfigs"
	request := map[string]interface{}{
		"DomainName": id,
	}
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("scdn", "2017-11-15", action, nil, request, true)
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
	v, err := jsonpath.Get("$.DomainConfigs.DomainConfig", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.DomainConfigs.DomainConfig", response)
	}
	if len(v.([]interface{})) < 1 {
		return object, errs.WrapErrorf(errs.NotFoundErr("SCDN", id), errs.NotFoundWithResponse, response)
	} else {
		if fmt.Sprint(v.([]interface{})[0].(map[string]interface{})["DomainName"]) != id {
			return object, errs.WrapErrorf(errs.NotFoundErr("SCDN", id), errs.NotFoundWithResponse, response)
		}
	}
	object = v.([]interface{})[0].(map[string]interface{})
	return object, nil
}

func (s *Service) DescribeScdnDomainDetail(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "DescribeScdnDomainDetail"
	request := map[string]interface{}{
		"DomainName": id,
	}
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("scdn", "2017-11-15", action, nil, request, true)
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
		if errs.IsExpectedErrors(err, []string{"InvalidDomain.NotFound"}) {
			return object, errs.WrapErrorf(errs.NotFoundErr("SCDN:Domain", id), errs.NotFoundMsg, errs.ProviderERROR, fmt.Sprint(response["RequestId"]))
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.DomainDetail", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.DomainDetail", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}

func (s *Service) DescribeScdnDomainConfig(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	parts, err := helper.ParseResourceId(id, 3)
	if err != nil {
		err = errs.WrapError(err)
		return
	}
	action := "DescribeScdnDomainConfigs"
	request := map[string]interface{}{
		"DomainName":    parts[0],
		"FunctionNames": parts[1],
	}
	if parts[2] != "" {
		request["ConfigId"] = parts[2]
	}
	response, err = client.RpcPost("scdn", "2017-11-15", action, nil, request, true)
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"InvalidDomain.NotFound"}) {
			err = errs.WrapErrorf(errs.NotFoundErr("ScdnDomainConfig", id), errs.NotFoundMsg, errs.ProviderERROR)
			return object, err
		}
		err = errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
		return object, err
	}
	logs.AddDebug(action, response, request)
	v, err := jsonpath.Get("$.DomainConfigs.DomainConfig", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.DomainConfigs.DomainConfig", response)
	}
	if len(v.([]interface{})) < 1 {
		return object, errs.WrapErrorf(errs.NotFoundErr("DCDN:DomainConfig", id), errs.NotFoundWithResponse, response)
	} else if len(v.([]interface{})) > 0 {
		object = v.([]interface{})[0].(map[string]interface{})
	}
	return object, nil
}

func (s *Service) ScdnDomainConfigStateRefreshFunc(id string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeScdnDomainConfig(id)
		if err != nil {
			if errs.NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", errs.WrapError(err)
		}

		for _, failState := range failStates {
			if object["Status"].(string) == failState {
				return object, object["Status"].(string), errs.WrapError(errs.Error(errs.FailedToReachTargetStatus, object["Status"].(string)))
			}
		}
		return object, object["Status"].(string), nil
	}
}

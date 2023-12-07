package alicloud

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

type DcdnService struct {
	client *connectivity.AliyunClient
}

func (s *DcdnService) convertSourcesToString(v []interface{}) (string, error) {
	arrayMaps := make([]interface{}, len(v))
	for i, vv := range v {
		item := vv.(map[string]interface{})
		arrayMaps[i] = map[string]interface{}{
			"Content":  item["content"],
			"Port":     item["port"],
			"Priority": item["priority"],
			"Type":     item["type"],
			"Weight":   item["weight"],
		}
	}
	maps, err := json.Marshal(arrayMaps)
	if err != nil {
		return "", WrapError(err)
	}
	return string(maps), nil
}

func (s *DcdnService) DescribeDcdnDomainCertificateInfo(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewDcdnClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "DescribeDcdnDomainCertificateInfo"
	request := map[string]interface{}{
		"RegionId":   s.client.RegionId,
		"DomainName": id,
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2018-01-15"), StringPointer("AK"), nil, request, &runtime)
	if err != nil {
		err = WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
		return
	}
	addDebug(action, response, request)
	v, err := jsonpath.Get("$.CertInfos.CertInfo", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.CertInfos.CertInfo", response)
	}
	if len(v.([]interface{})) < 1 {
		return object, WrapErrorf(Error(GetNotFoundMessage("DCDN", id)), NotFoundWithResponse, response)
	} else {
		if v.([]interface{})[0].(map[string]interface{})["DomainName"].(string) != id {
			return object, WrapErrorf(Error(GetNotFoundMessage("DCDN", id)), NotFoundWithResponse, response)
		}
	}
	object = v.([]interface{})[0].(map[string]interface{})
	return object, nil
}

func (s *DcdnService) DescribeDcdnDomain(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewDcdnClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "DescribeDcdnDomainDetail"
	request := map[string]interface{}{
		"RegionId":   s.client.RegionId,
		"DomainName": id,
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2018-01-15"), StringPointer("AK"), nil, request, &runtime)
	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidDomain.NotFound"}) {
			err = WrapErrorf(Error(GetNotFoundMessage("DcdnDomain", id)), NotFoundMsg, ProviderERROR)
			return object, err
		}
		err = WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
		return object, err
	}
	addDebug(action, response, request)
	v, err := jsonpath.Get("$.DomainDetail", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.DomainDetail", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}

func (s *DcdnService) DcdnDomainStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeDcdnDomain(id)
		if err != nil {
			if NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}

		for _, failState := range failStates {
			if object["DomainStatus"].(string) == failState {
				return object, object["DomainStatus"].(string), WrapError(Error(FailedToReachTargetStatus, object["DomainStatus"].(string)))
			}
		}
		return object, object["DomainStatus"].(string), nil
	}
}

func (s *DcdnService) DescribeDcdnDomainConfig(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	action := "DescribeDcdnDomainConfigs"

	conn, err := s.client.NewDcdnClient()
	if err != nil {
		return nil, WrapError(err)
	}

	parts, err := ParseResourceId(id, 3)
	if err != nil {
		return nil, WrapError(err)
	}

	request := map[string]interface{}{
		"DomainName":    parts[0],
		"FunctionNames": parts[1],
		"ConfigId":      parts[2],
	}

	idExist := false
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2018-01-15"), StringPointer("AK"), nil, request, &runtime)
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
		if IsExpectedErrors(err, []string{"InvalidDomain.NotFound"}) {
			return object, WrapErrorf(Error(GetNotFoundMessage("Dcdn:DomainConfig", id)), NotFoundMsg, ProviderERROR, fmt.Sprint(response["RequestId"]))
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}

	resp, err := jsonpath.Get("$.DomainConfigs.DomainConfig", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.DomainConfigs.DomainConfig", response)
	}

	if v, ok := resp.([]interface{}); !ok || len(v) < 1 {
		return object, WrapErrorf(Error(GetNotFoundMessage("Dcdn:DomainConfig", id)), NotFoundWithResponse, response)
	}

	for _, v := range resp.([]interface{}) {
		if fmt.Sprint(v.(map[string]interface{})["FunctionName"]) == parts[1] && fmt.Sprint(v.(map[string]interface{})["ConfigId"]) == parts[2] {
			idExist = true
			return v.(map[string]interface{}), nil
		}
	}

	if !idExist {
		return object, WrapErrorf(Error(GetNotFoundMessage("Dcdn:DomainConfig", id)), NotFoundWithResponse, response)
	}

	return object, nil
}

func (s *DcdnService) DcdnDomainConfigStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeDcdnDomainConfig(id)
		if err != nil {
			if NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}

		for _, failState := range failStates {
			if object["Status"].(string) == failState {
				return object, object["Status"].(string), WrapError(Error(FailedToReachTargetStatus, object["Status"].(string)))
			}
		}
		return object, object["Status"].(string), nil
	}
}

func (s *DcdnService) DescribeDcdnIpaDomain(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewDcdnClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "DescribeDcdnIpaDomainDetail"
	request := map[string]interface{}{
		"DomainName": id,
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2018-01-15"), StringPointer("AK"), nil, request, &runtime)
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
		if IsExpectedErrors(err, []string{"InvalidDomain.NotFound"}) {
			err = WrapErrorf(Error(GetNotFoundMessage("DcdnIpaDomain", id)), NotFoundMsg, ProviderERROR)
			return object, err
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.DomainDetail", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.DomainDetail", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}

func (s *DcdnService) DcdnIpaDomainStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeDcdnIpaDomain(id)
		if err != nil {
			if NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}

		for _, failState := range failStates {
			if object["DomainStatus"].(string) == failState {
				return object, object["DomainStatus"].(string), WrapError(Error(FailedToReachTargetStatus, object["DomainStatus"].(string)))
			}
		}
		return object, object["DomainStatus"].(string), nil
	}
}

func (s *DcdnService) DescribeDcdnWafPolicy(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewDcdnClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "DescribeDcdnWafPolicies"
	request := map[string]interface{}{
		"QueryArgs":  fmt.Sprintf("{\"PolicyIds\":\"%s\"}", id),
		"PageNumber": 1,
		"PageSize":   PageSizeLarge,
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2018-01-15"), StringPointer("AK"), nil, request, &runtime)
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
	v, err := jsonpath.Get("$.Policies", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Policies", response)
	}
	result, _ := v.([]interface{})
	if len(result) < 1 {
		return object, WrapErrorf(Error(GetNotFoundMessage("DcdnWafPolicy", id)), NotFoundMsg, ProviderERROR)
	}
	object = result[0].(map[string]interface{})
	return object, nil
}

func (s *DcdnService) DescribeDcdnWafDomain(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewDcdnClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "DescribeDcdnWafDomains"
	request := map[string]interface{}{
		"PageNumber": 1,
		"PageSize":   PageSizeLarge,
	}

	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2018-01-15"), StringPointer("AK"), nil, request, &runtime)
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
		v, err := jsonpath.Get("$.Domains", response)
		if err != nil {
			return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Domains", response)
		}
		if len(v.([]interface{})) < 1 {
			return object, WrapErrorf(Error(GetNotFoundMessage("DcdnWafDomains", id)), NotFoundWithResponse, response)
		}
		for _, v := range v.([]interface{}) {
			if fmt.Sprint(v.(map[string]interface{})["DomainName"]) == id {
				return v.(map[string]interface{}), nil
			}
		}
		if len(v.([]interface{})) < request["PageSize"].(int) {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}
	return object, nil
}

func (s *DcdnService) DescribeDcdnWafDomainDefenseScenes(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewDcdnClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "DescribeDcdnWafDomainDetail"
	request := map[string]interface{}{
		"DomainName": id,
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2018-01-15"), StringPointer("AK"), nil, request, &runtime)
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
	v, err := jsonpath.Get("$.Domain", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Domain", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}

func (s *DcdnService) DescribeDcdnWafPolicyDomainAttachment(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewDcdnClient()
	if err != nil {
		return nil, WrapError(err)
	}
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		return object, WrapError(err)
	}
	action := "DescribeDcdnWafPolicyDomains"
	request := map[string]interface{}{
		"PolicyId": parts[0],
	}
	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2018-01-15"), StringPointer("AK"), nil, request, &runtime)
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
		v, err := jsonpath.Get("$.Domains", response)
		if err != nil {
			return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Domains", response)
		}
		if len(v.([]interface{})) < 1 {
			return object, WrapErrorf(Error(GetNotFoundMessage("DCDN", id)), NotFoundWithResponse, response)
		}
		for _, v := range v.([]interface{}) {
			if fmt.Sprint(v.(map[string]interface{})["DomainName"]) == parts[1] {
				return v.(map[string]interface{}), nil
			}
		}
		if len(v.([]interface{})) < request["PageSize"].(int) {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}
	return object, nil
}

func (s *DcdnService) DescribeDcdnKv(id string) (object map[string]interface{}, err error) {
	conn, err := s.client.NewDcdnClient()
	if err != nil {
		return object, WrapError(err)
	}
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		return object, WrapError(err)
	}

	request := map[string]interface{}{
		"Namespace": parts[0],
		"Key":       parts[1],
	}

	var response map[string]interface{}
	action := "GetDcdnKv"
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		resp, err := conn.DoRequest(StringPointer(action), nil, StringPointer("GET"), StringPointer("2018-01-15"), StringPointer("AK"), request, nil, &runtime)
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
		if IsExpectedErrors(err, []string{"InvalidKey.Malformed", "InvalidKey.NotFound"}) {
			return object, WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$", response)
	}
	return v.(map[string]interface{}), nil
}

func (s *DcdnService) DescribeDcdnKvNamespace(id string) (object map[string]interface{}, err error) {
	conn, err := s.client.NewDcdnClient()
	if err != nil {
		return object, WrapError(err)
	}

	request := map[string]interface{}{
		"Namespace": id,
	}

	var response map[string]interface{}
	action := "DescribeDcdnKvNamespace"
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		resp, err := conn.DoRequest(StringPointer(action), nil, StringPointer("GET"), StringPointer("2018-01-15"), StringPointer("AK"), request, nil, &runtime)
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
		if IsExpectedErrors(err, []string{"InvalidNameSpace.Malformed", "InvalidNameSpace.NotFound"}) {
			return object, WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$", response)
	}
	return v.(map[string]interface{}), nil
}

func (s *DcdnService) DcdnKvNamespaceStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeDcdnKvNamespace(id)
		if err != nil {
			if NotFoundError(err) {
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

func (s *DcdnService) DescribeDcdnKvAccountStatus() (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewDcdnClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "DescribeDcdnKvAccountStatus"
	request := map[string]interface{}{}

	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("GET"), StringPointer("2018-01-15"), StringPointer("AK"), request, nil, &runtime)
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
		return object, WrapErrorf(err, DefaultErrorMsg, "", action, AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, "", "$", response)
	}
	object = v.(map[string]interface{})
	return object, nil

}

func (s *DcdnService) DescribeDcdnWafRule(id string) (object map[string]interface{}, err error) {
	conn, err := s.client.NewDcdnClient()
	if err != nil {
		return object, WrapError(err)
	}

	request := map[string]interface{}{
		"RuleId": id,
	}

	var response map[string]interface{}
	action := "DescribeDcdnWafRule"
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		resp, err := conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2018-01-15"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			if NeedRetry(err) || IsExpectedErrors(err, []string{"InternalError"}) {
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

	if _, ok := response["Rule"]; !ok {
		return object, WrapErrorf(Error(GetNotFoundMessage("WafRule", id)), NotFoundWithResponse, response)
	}

	v, err := jsonpath.Get("$.Rule", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Rule", response)
	}
	return v.(map[string]interface{}), nil
}

func (s *DcdnService) DescribeDcdnEr(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	action := "DescribeRoutine"

	conn, err := s.client.NewDcdnClient()
	if err != nil {
		return object, WrapError(err)
	}

	request := map[string]interface{}{
		"Name": id,
	}

	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2018-01-15"), StringPointer("AK"), nil, request, &runtime)
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
		if _, ok := response["Content"]; !ok {
			return object, WrapErrorf(Error(GetNotFoundMessage("Dcdn:Er", id)), NotFoundMsg, ProviderERROR)
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$.Content", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Content", response)
	}

	object = v.(map[string]interface{})

	return object, nil
}

func (s *DcdnService) ListTagResources(id string, resourceType string) (object interface{}, err error) {
	conn, err := s.client.NewDcdnClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "DescribeDcdnTagResources"

	request := map[string]interface{}{
		"ResourceType": resourceType,
	}

	resourceIdNum := strings.Count(id, ":")

	switch resourceIdNum {
	case 0:
		request["ResourceId.1"] = id
	case 1:
		parts, err := ParseResourceId(id, 2)
		if err != nil {
			return object, WrapError(err)
		}
		request["ResourceId.1"] = parts[resourceIdNum]
	}

	tags := make([]interface{}, 0)
	var response map[string]interface{}

	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2018-01-15"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)

		resp, err := jsonpath.Get("$.TagResources", response)
		if err != nil {
			return resource.NonRetryableError(WrapErrorf(err, FailedGetAttributeMsg, id, "$.TagResources", response))
		}

		if v, ok := resp.([]interface{}); ok && len(v) > 0 {
			tags = append(tags, v[0].(map[string]interface{})["Tag"].([]interface{})...)
		}

		return nil
	})

	if err != nil {
		err = WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
		return
	}

	return tags, nil
}

func (s *DcdnService) SetResourceTags(d *schema.ResourceData, resourceType string) error {

	resourceIdNum := strings.Count(d.Id(), ":")

	if d.HasChange("tags") {
		added, removed := parsingTags(d)
		conn, err := s.client.NewDcdnClient()
		if err != nil {
			return WrapError(err)
		}

		removedTagKeys := make([]string, 0)
		for _, v := range removed {
			if !ignoredTags(v, "") {
				removedTagKeys = append(removedTagKeys, v)
			}
		}

		if len(removedTagKeys) > 0 {
			action := "UntagDcdnResources"
			request := map[string]interface{}{
				"ResourceType": resourceType,
			}

			switch resourceIdNum {
			case 0:
				request["ResourceId.1"] = d.Id()
			case 1:
				parts, err := ParseResourceId(d.Id(), 2)
				if err != nil {
					return WrapError(err)
				}
				request["ResourceId.1"] = parts[resourceIdNum]
			}

			for i, key := range removedTagKeys {
				request[fmt.Sprintf("TagKey.%d", i+1)] = key
			}

			runtime := util.RuntimeOptions{}
			runtime.SetAutoretry(true)
			wait := incrementalWait(2*time.Second, 1*time.Second)
			err := resource.Retry(10*time.Minute, func() *resource.RetryError {
				response, err := conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2018-01-15"), StringPointer("AK"), nil, request, &runtime)
				if err != nil {
					if NeedRetry(err) {
						wait()
						return resource.RetryableError(err)

					}
					return resource.NonRetryableError(err)
				}
				addDebug(action, response, request)
				return nil
			})
			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
			}
		}
		if len(added) > 0 {
			action := "TagDcdnResources"
			request := map[string]interface{}{
				"ResourceType": resourceType,
			}

			switch resourceIdNum {
			case 0:
				request["ResourceId.1"] = d.Id()
			case 1:
				parts, err := ParseResourceId(d.Id(), 2)
				if err != nil {
					return WrapError(err)
				}
				request["ResourceId.1"] = parts[resourceIdNum]
			}

			count := 1
			for key, value := range added {
				request[fmt.Sprintf("Tag.%d.Key", count)] = key
				request[fmt.Sprintf("Tag.%d.Value", count)] = value
				count++
			}

			runtime := util.RuntimeOptions{}
			runtime.SetAutoretry(true)
			wait := incrementalWait(2*time.Second, 1*time.Second)
			err := resource.Retry(10*time.Minute, func() *resource.RetryError {
				response, err := conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2018-01-15"), StringPointer("AK"), nil, request, &runtime)
				if err != nil {
					if NeedRetry(err) {
						wait()
						return resource.RetryableError(err)

					}
					return resource.NonRetryableError(err)
				}
				addDebug(action, response, request)
				return nil
			})
			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
			}
		}
		d.SetPartial("tags")
	}
	return nil
}

package alicloud

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

type SaeService struct {
	client *connectivity.AliyunClient
}

func (s *SaeService) DescribeSaeNamespace(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewServerlessClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "/pop/v1/paas/namespace"
	request := map[string]*string{
		"NamespaceId": StringPointer(id),
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer("2019-05-06"), nil, StringPointer("GET"), StringPointer("AK"), StringPointer(action), request, nil, nil, &util.RuntimeOptions{})
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidNamespaceId.NotFound"}) {
			return object, WrapErrorf(Error(GetNotFoundMessage("SAE:Namespace", id)), NotFoundMsg, ProviderERROR)
		}
		return object, WrapErrorf(err, DefaultErrorMsg, "alicloud_sae_namespace", "GET "+action, AlibabaCloudSdkGoERROR)
	}
	if respBody, isExist := response["body"]; isExist {
		response = respBody.(map[string]interface{})
	} else {
		return object, WrapError(fmt.Errorf("%s failed, response: %v", action, response))
	}
	addDebug(action, response, request)
	if IsExpectedErrorCodes(fmt.Sprint(response["Code"]), []string{"InvalidNamespaceId.NotFound"}) {
		return object, WrapErrorf(Error(GetNotFoundMessage("SAE:Namespace", id)), NotFoundMsg, ProviderERROR)
	}
	if fmt.Sprint(response["Success"]) == "false" {
		return object, WrapError(fmt.Errorf("%s failed, response: %v", "GET "+action, response))
	}
	v, err := jsonpath.Get("$.Data", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Data", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}

func (s *SaeService) DescribeSaeConfigMap(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewServerlessClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "/pop/v1/sam/configmap/configMap"
	request := map[string]*string{
		"ConfigMapId": StringPointer(id),
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer("2019-05-06"), nil, StringPointer("GET"), StringPointer("AK"), StringPointer(action), request, nil, nil, &util.RuntimeOptions{})
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
		if IsExpectedErrors(err, []string{"NotFound.ConfigMap"}) {
			return object, WrapErrorf(Error(GetNotFoundMessage("SAE:ConfigMap", id)), NotFoundMsg, ProviderERROR)
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	if respBody, isExist := response["body"]; isExist {
		response = respBody.(map[string]interface{})
	} else {
		return object, WrapError(fmt.Errorf("%s failed, response: %v", action, response))
	}
	addDebug(action, response, request)

	if fmt.Sprint(response["Success"]) == "false" {
		return object, WrapError(fmt.Errorf("%s failed, response: %v", action, response))
	}
	v, err := jsonpath.Get("$.Data", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Data", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}

func (s *SaeService) DescribeApplicationStatus(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewServerlessClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "/pop/v1/sam/app/describeApplicationStatus"
	request := map[string]*string{
		"AppId": StringPointer(id),
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer("2019-05-06"), nil, StringPointer("GET"), StringPointer("AK"), StringPointer(action), request, nil, nil, &util.RuntimeOptions{})
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
	if respBody, isExist := response["body"]; isExist {
		response = respBody.(map[string]interface{})
	} else {
		return object, WrapError(fmt.Errorf("%s failed, response: %v", "Put "+action, response))
	}
	if fmt.Sprint(response["Success"]) == "false" {
		return object, WrapError(fmt.Errorf("%s failed, response: %v", action, response))
	}
	v, err := jsonpath.Get("$.Data", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Data", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}

func (s *SaeService) DescribeSaeApplication(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewServerlessClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "/pop/v1/sam/app/describeApplicationConfig"
	request := map[string]*string{
		"AppId": StringPointer(id),
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer("2019-05-06"), nil, StringPointer("GET"), StringPointer("AK"), StringPointer(action), request, nil, nil, &util.RuntimeOptions{})
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
		if IsExpectedErrors(err, []string{"InvalidAppId.NotFound"}) {
			return object, WrapErrorf(Error(GetNotFoundMessage("SAE:Application", id)), NotFoundMsg, ProviderERROR)
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	if respBody, isExist := response["body"]; isExist {
		response = respBody.(map[string]interface{})
	} else {
		return object, WrapError(fmt.Errorf("%s failed, response: %v", "Put "+action, response))
	}
	if fmt.Sprint(response["Success"]) == "false" {
		return object, WrapError(fmt.Errorf("%s failed, response: %v", action, response))
	}
	v, err := jsonpath.Get("$.Data", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Data", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}

func (s *SaeService) DescribeApplicationSlb(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewServerlessClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "/pop/v1/sam/app/slb"
	request := map[string]*string{
		"AppId": StringPointer(id),
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer("2019-05-06"), nil, StringPointer("GET"), StringPointer("AK"), StringPointer(action), request, nil, nil, &util.RuntimeOptions{})
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
	if respBody, isExist := response["body"]; isExist {
		response = respBody.(map[string]interface{})
	} else {
		return object, WrapError(fmt.Errorf("%s failed, response: %v", "Put "+action, response))
	}
	if fmt.Sprint(response["Success"]) == "false" {
		return object, WrapError(fmt.Errorf("%s failed, response: %v", action, response))
	}
	v, err := jsonpath.Get("$.Data", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Data", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}

func (s *SaeService) DescribeApplicationImage(id, imageUrl string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewServerlessClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "/pop/v1/sam/container/describeApplicationImage"
	request := map[string]*string{
		"AppId":    StringPointer(id),
		"ImageUrl": StringPointer(imageUrl),
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer("2019-05-06"), nil, StringPointer("GET"), StringPointer("AK"), StringPointer(action), request, nil, nil, &util.RuntimeOptions{})
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
	if respBody, isExist := response["body"]; isExist {
		response = respBody.(map[string]interface{})
	} else {
		return object, WrapError(fmt.Errorf("%s failed, response: %v", "Put "+action, response))
	}
	if fmt.Sprint(response["Success"]) == "false" {
		return object, WrapError(fmt.Errorf("%s failed, response: %v", action, response))
	}
	v, err := jsonpath.Get("$.Data", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Data", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}

func (s *SaeService) UpdateSlb(d *schema.ResourceData) error {
	if d.HasChange("intranet") || d.HasChange("internet") || d.HasChange("internet_ip") || d.HasChange("intranet_ip") {
		update := false
		request := map[string]*string{
			"AppId": StringPointer(d.Id()),
		}

		if v, ok := d.GetOk("intranet"); ok {
			nrawIntranet := v.(*schema.Set).List()
			if len(nrawIntranet) == 0 {
				// need to unbind intranet
				update = true
				request["Intranet"] = StringPointer(strconv.FormatBool(true))
			}
		}
		if v, ok := d.GetOk("internet"); ok {
			nrawInteranet := v.(*schema.Set).List()
			if len(nrawInteranet) == 0 {
				// need to unbind intranet
				update = true
				request["Internet"] = StringPointer(strconv.FormatBool(true))
			}
		}
		if update {
			action := "/pop/v1/sam/app/slb"
			conn, err := s.client.NewServerlessClient()
			if err != nil {
				return WrapError(err)
			}
			wait := incrementalWait(3*time.Second, 3*time.Second)
			var response map[string]interface{}
			err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
				response, err = conn.DoRequest(StringPointer("2019-05-06"), nil, StringPointer("DELETE"), StringPointer("AK"), StringPointer(action), request, nil, nil, &util.RuntimeOptions{})
				if err != nil {
					if IsExpectedErrors(err, []string{"Application.InvalidStatus", "Application.ChangerOrderRunning"}) || NeedRetry(err) {
						wait()
						return resource.RetryableError(err)
					}
					return resource.NonRetryableError(err)
				}
				return nil
			})
			addDebug(action, response, request)
			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), "POST "+action, AlibabaCloudSdkGoERROR)
			}
			if respBody, isExist := response["body"]; isExist {
				response = respBody.(map[string]interface{})
			} else {
				return WrapError(fmt.Errorf("%s failed, response: %v", "POST "+action, response))
			}
			if fmt.Sprint(response["Success"]) == "false" {
				return WrapError(fmt.Errorf("%s failed, response: %v", "POST "+action, response))
			}
			return nil
		}

		update = false
		request = map[string]*string{
			"AppId": StringPointer(d.Id()),
		}
		if v, ok := d.GetOk("internet_ip"); ok {
			update = true
			request["InternetSlbId"] = StringPointer(v.(string))
		}
		if v, ok := d.GetOk("intranet_ip"); ok {
			update = true
			request["IntranetSlbId"] = StringPointer(v.(string))
		}
		if v, ok := d.GetOk("intranet"); ok {
			update = true
			for _, intranet := range v.(*schema.Set).List() {
				intranetMap := intranet.(map[string]interface{})
				intranetReq := []interface{}{
					map[string]interface{}{
						"httpsCertId": intranetMap["https_cert_id"],
						"protocol":    intranetMap["protocol"],
						"targetPort":  intranetMap["target_port"],
						"port":        intranetMap["port"],
					},
				}
				obj, err := json.Marshal(intranetReq)
				if err != nil {
					return WrapError(err)
				}
				request["Intranet"] = StringPointer(string(obj))
			}
		}

		if v, ok := d.GetOk("internet"); ok {
			update = true
			for _, internet := range v.(*schema.Set).List() {
				internetMap := internet.(map[string]interface{})
				internetReq := []interface{}{
					map[string]interface{}{
						"httpsCertId": internetMap["https_cert_id"],
						"protocol":    internetMap["protocol"],
						"targetPort":  internetMap["target_port"],
						"port":        internetMap["port"],
					},
				}
				obj, err := json.Marshal(internetReq)
				if err != nil {
					return WrapError(err)
				}
				request["Internet"] = StringPointer(string(obj))
			}
		}

		if update {
			action := "/pop/v1/sam/app/slb"
			conn, err := s.client.NewServerlessClient()
			if err != nil {
				return WrapError(err)
			}
			wait := incrementalWait(3*time.Second, 3*time.Second)
			var response map[string]interface{}
			err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
				response, err = conn.DoRequest(StringPointer("2019-05-06"), nil, StringPointer("POST"), StringPointer("AK"), StringPointer(action), request, nil, nil, &util.RuntimeOptions{})
				if err != nil {
					if IsExpectedErrors(err, []string{"Application.InvalidStatus", "Application.ChangerOrderRunning"}) || NeedRetry(err) {
						wait()
						return resource.RetryableError(err)
					}
					return resource.NonRetryableError(err)
				}
				return nil
			})
			addDebug(action, response, request)
			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), "POST "+action, AlibabaCloudSdkGoERROR)
			}
			if respBody, isExist := response["body"]; isExist {
				response = respBody.(map[string]interface{})
			} else {
				return WrapError(fmt.Errorf("%s failed, response: %v", "POST "+action, response))
			}
			if fmt.Sprint(response["Success"]) == "false" {
				return WrapError(fmt.Errorf("%s failed, response: %v", "POST "+action, response))
			}
		}
	}

	return nil
}

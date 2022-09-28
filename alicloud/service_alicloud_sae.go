package alicloud

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
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
		return object, WrapError(fmt.Errorf("%s failed, response: %v", "GET "+action, response))
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

func (s *SaeService) DescribeSaeIngress(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewServerlessClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "/pop/v1/sam/ingress/Ingress"
	request := map[string]*string{
		"IngressId": StringPointer(id),
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
	if IsExpectedErrors(err, []string{"InvalidParameter.WithMessage"}) {
		return object, WrapErrorf(Error(GetNotFoundMessage("SAE:Ingress", id)), NotFoundMsg, ProviderERROR)
	}
	addDebug(action, response, request)

	if respBody, isExist := response["body"]; isExist {
		response = respBody.(map[string]interface{})
	} else {
		return object, WrapError(fmt.Errorf("%s failed, response: %v", "GET "+action, response))
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

func (s *SaeService) DescribeIngress(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewServerlessClient()
	if err != nil {
		return nil, WrapError(err)
	}

	action := "DescribeIngress"
	request := map[string]*string{
		"IngressId": StringPointer(id),
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
		return object, WrapError(fmt.Errorf("%s failed, response: %v", "GET "+action, response))
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

func (s *SaeService) UpdateSlb(d *schema.ResourceData) error {
	if d.HasChange("intranet") || d.HasChange("internet") || d.HasChange("internet_slb_id") || d.HasChange("intranet_slb_id") {
		update := false
		request := map[string]*string{
			"AppId": StringPointer(d.Id()),
		}
		//unbind intranet
		if d.HasChange("intranet") {
			oraw, nraw := d.GetChange("intranet")
			remove := oraw.(*schema.Set).Difference(nraw.(*schema.Set)).List()
			if len(remove) != 0 {
				update = true
				request["Intranet"] = StringPointer(strconv.FormatBool(true))
			}
		}
		//unbind internet
		if d.HasChange("internet") {
			oraw, nraw := d.GetChange("internet")
			remove := oraw.(*schema.Set).Difference(nraw.(*schema.Set)).List()
			if len(remove) != 0 {
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
		if d.HasChange("internet_slb_id") {
			update = true
			if v, exist := d.GetOk("internet_slb_id"); exist {
				request["InternetSlbId"] = StringPointer(v.(string))
			}

		}
		if d.HasChange("intranet_slb_id") {
			update = true
			if v, exist := d.GetOk("intranet_slb_id"); exist {
				request["IntranetSlbId"] = StringPointer(v.(string))
			}

		}
		if d.HasChange("intranet") {
			update = true
			for _, intranet := range d.Get("intranet").(*schema.Set).List() {
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

		if d.HasChange("internet") {
			update = true
			for _, internet := range d.Get("internet").(*schema.Set).List() {
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

func (s *SaeService) SaeApplicationStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeApplicationStatus(id)
		if err != nil {
			if NotFoundError(err) {
				// Set this to nil if nothing matched
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}

		for _, failState := range failStates {
			if fmt.Sprint(object["LastChangeOrderStatus"]) == failState {
				return object, fmt.Sprint(object["LastChangeOrderStatus"]), WrapError(Error(FailedToReachTargetStatus, fmt.Sprint(object["LastChangeOrderStatus"])))
			}
		}
		if fmt.Sprint(object["LastChangeOrderStatus"]) == "RUNNING" && fmt.Sprint(object["SubStatus"]) == "runningButHasError" {
			return object, fmt.Sprint(object["LastChangeOrderStatus"]), WrapError(Error(FailedToReachTargetStatus, fmt.Sprint(object["LastChangeOrderStatus"])))
		}
		return object, fmt.Sprint(object["LastChangeOrderStatus"]), nil
	}
}

func (s *SaeService) DescribeSaeApplicationScalingRule(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewServerlessClient()
	if err != nil {
		return nil, WrapError(err)
	}
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		err = WrapError(err)
		return
	}
	action := "/pop/v1/sam/scale/applicationScalingRules"
	idExist := false
	request := map[string]*string{
		"AppId":       StringPointer(parts[0]),
		"PageSize":    StringPointer(strconv.Itoa(PageSizeLarge)),
		"CurrentPage": StringPointer("1"),
	}
	for {
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
				return object, WrapErrorf(Error(GetNotFoundMessage("SAE:ApplicationScalingRule", id)), NotFoundMsg, ProviderERROR)
			}
			return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
		}
		if respBody, isExist := response["body"]; isExist {
			response = respBody.(map[string]interface{})
		} else {
			return object, WrapError(fmt.Errorf("%s failed, response: %v", "Put "+action, response))
		}
		v, err := jsonpath.Get("$.Data.ApplicationScalingRules", response)
		if err != nil {
			return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Data.ApplicationScalingRules", response)
		}
		if len(v.([]interface{})) < 1 {
			return object, WrapErrorf(Error(GetNotFoundMessage("SAE:ApplicationScalingRule", id)), NotFoundMsg, ProviderERROR)
		}
		for _, v := range v.([]interface{}) {
			if fmt.Sprint(v.(map[string]interface{})["ScaleRuleName"]) == parts[1] {
				idExist = true
				return v.(map[string]interface{}), nil
			}
		}
		if len(v.([]interface{})) < PageSizeLarge {
			break
		}
		currentPage, err := strconv.Atoi(*request["CurrentPage"])
		if err != nil {
			return object, WrapError(err)
		}
		request["CurrentPage"] = StringPointer(strconv.Itoa(currentPage + 1))
	}

	if !idExist {
		return object, WrapErrorf(Error(GetNotFoundMessage("SAE:ApplicationScalingRule", id)), NotFoundMsg, ProviderERROR)
	}
	return object, nil
}

func (s *SaeService) DescribeSaeGreyTagRoute(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewServerlessClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "/pop/v1/sam/tagroute/greyTagRoute"
	request := map[string]*string{
		"GreyTagRouteId": StringPointer(id),
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
		if IsExpectedErrors(err, []string{"InvalidParameter.WithMessage"}) {
			return object, WrapErrorf(Error(GetNotFoundMessage("SAE", id)), NotFoundWithResponse, response)
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

func (s *SaeService) ListTagResources(id string, resourceType string) (object interface{}, err error) {
	conn, err := s.client.NewServerlessClient()
	if err != nil {
		return nil, WrapError(err)
	}
	ids, err := json.Marshal([]string{id})
	if err != nil {
		return object, err
	}
	action := "/tags"
	request := map[string]*string{
		"RegionId":     StringPointer(s.client.RegionId),
		"ResourceType": StringPointer(resourceType),
		"ResourceIds":  StringPointer(string(ids)),
	}
	tags := make([]interface{}, 0)
	var response map[string]interface{}

	for {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer("2019-05-06"), nil, StringPointer("GET"), StringPointer("AK"), StringPointer(action), request, nil, nil, &util.RuntimeOptions{})
			if err != nil {
				if IsExpectedErrors(err, []string{Throttling}) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, request)
		if err != nil {
			err = WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
			return
		}
		if respBody, isExist := response["body"]; isExist {
			response = respBody.(map[string]interface{})
		} else {
			return object, WrapError(fmt.Errorf("%s failed, response: %v", "Put "+action, response))
		}
		if fmt.Sprint(response["Success"]) == "false" {
			return object, WrapError(fmt.Errorf("%s failed, response: %v", action, response))
		}
		v, err := jsonpath.Get("$.Data.TagResources", response)
		if err != nil {
			return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Data.TagResources", response)
		}
		if v != nil {
			tags = append(tags, v.([]interface{})...)
		}
		if response["NextToken"] == nil {
			break
		}
		request["NextToken"] = StringPointer(response["NextToken"].(string))
	}

	return tags, nil
}

func (s *SaeService) SetResourceTags(d *schema.ResourceData, resourceType string) error {

	if d.HasChange("tags") {
		added, removed := parsingTags(d)
		conn, err := s.client.NewServerlessClient()
		if err != nil {
			return WrapError(err)
		}
		ids, err := json.Marshal([]string{d.Id()})
		if err != nil {
			return err
		}

		removedTagKeys := make([]interface{}, 0)
		for _, v := range removed {
			if !ignoredTags(v, "") {
				removedTagKeys = append(removedTagKeys, v)
			}
		}
		if len(removedTagKeys) > 0 {
			action := "/tags"
			request := map[string]*string{
				"RegionId":     StringPointer(s.client.RegionId),
				"ResourceType": StringPointer(resourceType),
				"ResourceIds":  StringPointer(string(ids)),
				"TagKeys":      StringPointer(convertListToJsonString(removedTagKeys)),
			}

			wait := incrementalWait(2*time.Second, 1*time.Second)
			err := resource.Retry(10*time.Minute, func() *resource.RetryError {
				response, err := conn.DoRequest(StringPointer("2019-05-06"), nil, StringPointer("DELETE"), StringPointer("AK"), StringPointer(action), request, nil, nil, &util.RuntimeOptions{})
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
			action := "/tags"
			request := map[string]*string{
				"RegionId":     StringPointer(s.client.RegionId),
				"ResourceType": StringPointer(resourceType),
				"ResourceIds":  StringPointer(string(ids)),
			}
			tags := make([]map[string]interface{}, len(added))
			for key, value := range added {
				tags = append(tags, map[string]interface{}{
					"key":   key,
					"value": value,
				})
			}
			jsonString, err := convertListMapToJsonString(tags)
			if err != nil {
				return WrapError(err)
			}
			request["Tags"] = StringPointer(jsonString)

			wait := incrementalWait(2*time.Second, 1*time.Second)
			err = resource.Retry(10*time.Minute, func() *resource.RetryError {
				response, err := conn.DoRequest(StringPointer("2019-05-06"), nil, StringPointer("POST"), StringPointer("AK"), StringPointer(action), request, nil, nil, &util.RuntimeOptions{})
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

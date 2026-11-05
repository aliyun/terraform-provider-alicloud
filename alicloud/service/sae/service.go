package sae

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/convert"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/tags"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Service struct {
	client *connectivity.AliyunClient
}

// NewService constructs the SAE (2019-05-06) Service from a connectivity client.
func NewService(client *connectivity.AliyunClient) *Service {
	return &Service{client: client}
}

func (s *Service) DescribeSaeNamespace(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "/pop/v1/paas/namespace"
	request := map[string]*string{
		"NamespaceId": helper.StringPointer(id),
	}
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RoaGet("sae", "2019-05-06", action, request, nil, nil)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"InvalidNamespaceId.NotFound"}) {
			return object, errs.WrapErrorf(errs.NotFoundErr("SAE:Namespace", id), errs.NotFoundMsg, errs.ProviderERROR)
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_sae_namespace", "GET "+action, errs.AlibabaCloudSdkGoERROR)
	}
	logs.AddDebug(action, response, request)
	v, err := jsonpath.Get("$.Data", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.Data", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}

func (s *Service) DescribeSaeConfigMap(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "/pop/v1/sam/configmap/configMap"
	request := map[string]*string{
		"ConfigMapId": helper.StringPointer(id),
	}
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RoaGet("sae", "2019-05-06", action, request, nil, nil)
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
		if errs.IsExpectedErrors(err, []string{"NotFound.ConfigMap"}) {
			return object, errs.WrapErrorf(errs.NotFoundErr("SAE:ConfigMap", id), errs.NotFoundMsg, errs.ProviderERROR)
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}
	logs.AddDebug(action, response, request)
	v, err := jsonpath.Get("$.Data", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.Data", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}

func (s *Service) DescribeApplicationStatus(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "/pop/v1/sam/app/describeApplicationStatus"
	request := map[string]*string{
		"AppId": helper.StringPointer(id),
	}
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RoaGet("sae", "2019-05-06", action, request, nil, nil)
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
	v, err := jsonpath.Get("$.Data", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.Data", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}

func (s *Service) DescribeSaeApplication(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "/pop/v1/sam/app/describeApplicationConfig"
	request := map[string]*string{
		"AppId": helper.StringPointer(id),
	}
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RoaGet("sae", "2019-05-06", action, request, nil, nil)
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
		if errs.IsExpectedErrors(err, []string{"InvalidAppId.NotFound"}) {
			return object, errs.WrapErrorf(errs.NotFoundErr("SAE:Application", id), errs.NotFoundMsg, errs.ProviderERROR)
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.Data", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.Data", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}

func (s *Service) DescribeSaeApplicationChangeOrder(orderId string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	action := "/pop/v1/sam/changeorder/DescribeChangeOrder"

	client := s.client

	request := map[string]*string{
		"ChangeOrderId": helper.StringPointer(orderId),
	}

	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RoaGet("sae", "2019-05-06", action, request, nil, nil)
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
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, orderId, action, errs.AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.Data", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, orderId, "$.Data", response)
	}

	object = v.(map[string]interface{})

	return object, nil
}

func (s *Service) SaeApplicationChangeOrderStateRefreshFunc(orderId string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeSaeApplicationChangeOrder(orderId)
		if err != nil {
			if errs.NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", errs.WrapError(err)
		}

		for _, failState := range failStates {
			if fmt.Sprint(object["Status"]) == failState {
				return object, fmt.Sprint(object["Status"]), errs.WrapError(errs.Error(errs.FailedToReachTargetStatus, fmt.Sprint(object["Status"])))
			}
		}
		return object, fmt.Sprint(object["Status"]), nil
	}
}

func (s *Service) DescribeSaeIngress(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "/pop/v1/sam/ingress/Ingress"
	request := map[string]*string{
		"IngressId": helper.StringPointer(id),
	}
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RoaGet("sae", "2019-05-06", action, request, nil, nil)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"InvalidParameter.WithMessage"}) {
			return object, errs.WrapErrorf(errs.NotFoundErr("SAE:Ingress", id), errs.NotFoundMsg, errs.ProviderERROR)
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}
	logs.AddDebug(action+"-Describe", response, fmt.Sprint(request))
	v, err := jsonpath.Get("$.Data", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.Data", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}

func (s *Service) DescribeApplicationImage(id, imageUrl string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "/pop/v1/sam/container/describeApplicationImage"
	request := map[string]*string{
		"AppId":    helper.StringPointer(id),
		"ImageUrl": helper.StringPointer(imageUrl),
	}
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RoaGet("sae", "2019-05-06", action, request, nil, nil)
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
	v, err := jsonpath.Get("$.Data", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.Data", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}

func (s *Service) DescribeIngress(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client

	action := "DescribeIngress"
	request := map[string]*string{
		"IngressId": helper.StringPointer(id),
	}
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RoaGet("sae", "2019-05-06", action, request, nil, nil)
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
	v, err := jsonpath.Get("$.Data", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.Data", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}

func (s *Service) DescribeApplicationSlb(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "/pop/v1/sam/app/slb"
	request := map[string]*string{
		"AppId": helper.StringPointer(id),
	}
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RoaGet("sae", "2019-05-06", action, request, nil, nil)
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
	v, err := jsonpath.Get("$.Data", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.Data", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}

func (s *Service) UpdateSlb(d *schema.ResourceData) error {
	if d.HasChange("intranet") || d.HasChange("internet") || d.HasChange("internet_slb_id") || d.HasChange("intranet_slb_id") {
		update := false
		client := s.client
		var err error
		request := map[string]*string{
			"AppId": helper.StringPointer(d.Id()),
		}
		//unbind intranet
		if d.HasChange("intranet") {
			oraw, nraw := d.GetChange("intranet")
			remove := oraw.(*schema.Set).Difference(nraw.(*schema.Set)).List()
			if len(remove) != 0 {
				update = true
				request["Intranet"] = helper.StringPointer(strconv.FormatBool(true))
			}
		}
		//unbind internet
		if d.HasChange("internet") {
			oraw, nraw := d.GetChange("internet")
			remove := oraw.(*schema.Set).Difference(nraw.(*schema.Set)).List()
			if len(remove) != 0 {
				update = true
				request["Internet"] = helper.StringPointer(strconv.FormatBool(true))
			}
		}
		if update {
			action := "/pop/v1/sam/app/slb"
			wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
			var response map[string]interface{}
			err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
				response, err = client.RoaDelete("sae", "2019-05-06", action, request, nil, nil, false)
				if err != nil {
					if errs.IsExpectedErrors(err, []string{"Application.InvalidStatus", "Application.ChangerOrderRunning"}) || helper.NeedRetry(err) {
						wait()
						return retry.RetryableError(err)
					}
					return retry.NonRetryableError(err)
				}
				return nil
			})
			logs.AddDebug(action, response, request)
			if err != nil {
				return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), "POST "+action, errs.AlibabaCloudSdkGoERROR)
			}
			return nil
		}

		update = false
		request = map[string]*string{
			"AppId": helper.StringPointer(d.Id()),
		}
		if d.HasChange("internet_slb_id") {
			update = true
			if v, exist := d.GetOk("internet_slb_id"); exist {
				request["InternetSlbId"] = helper.StringPointer(v.(string))
			}

		}
		if d.HasChange("intranet_slb_id") {
			update = true
			if v, exist := d.GetOk("intranet_slb_id"); exist {
				request["IntranetSlbId"] = helper.StringPointer(v.(string))
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
					return errs.WrapError(err)
				}
				request["Intranet"] = helper.StringPointer(string(obj))
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
					return errs.WrapError(err)
				}
				request["Internet"] = helper.StringPointer(string(obj))
			}
		}

		if update {
			action := "/pop/v1/sam/app/slb"
			wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
			var response map[string]interface{}
			err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
				response, err = client.RoaPost("sae", "2019-05-06", action, request, nil, nil, false)
				if err != nil {
					if errs.IsExpectedErrors(err, []string{"Application.InvalidStatus", "Application.ChangerOrderRunning"}) || helper.NeedRetry(err) {
						wait()
						return retry.RetryableError(err)
					}
					return retry.NonRetryableError(err)
				}
				return nil
			})
			logs.AddDebug(action, response, request)
			if err != nil {
				return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), "POST "+action, errs.AlibabaCloudSdkGoERROR)
			}
		}
	}
	return nil
}

func (s *Service) SaeApplicationStateRefreshFunc(id string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeApplicationStatus(id)
		if err != nil {
			if errs.NotFoundError(err) {
				// Set this to nil if nothing matched
				return nil, "", nil
			}
			return nil, "", errs.WrapError(err)
		}

		for _, failState := range failStates {
			if fmt.Sprint(object["LastChangeOrderStatus"]) == failState {
				return object, fmt.Sprint(object["LastChangeOrderStatus"]), errs.WrapError(errs.Error(errs.FailedToReachTargetStatus, fmt.Sprint(object["LastChangeOrderStatus"])))
			}
		}
		if fmt.Sprint(object["LastChangeOrderStatus"]) == "RUNNING" && fmt.Sprint(object["SubStatus"]) == "runningButHasError" {
			return object, fmt.Sprint(object["LastChangeOrderStatus"]), errs.WrapError(errs.Error(errs.FailedToReachTargetStatus, fmt.Sprint(object["LastChangeOrderStatus"])))
		}
		return object, fmt.Sprint(object["LastChangeOrderStatus"]), nil
	}
}

func (s *Service) DescribeSaeApplicationScalingRule(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	parts, err := helper.ParseResourceId(id, 2)
	if err != nil {
		err = errs.WrapError(err)
		return
	}
	action := "/pop/v1/sam/scale/applicationScalingRules"
	idExist := false
	request := map[string]*string{
		"AppId":       helper.StringPointer(parts[0]),
		"PageSize":    helper.StringPointer(strconv.Itoa(helper.PageSizeLarge)),
		"CurrentPage": helper.StringPointer("1"),
	}
	for {
		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(5*time.Minute, func() *retry.RetryError {
			response, err = client.RoaGet("sae", "2019-05-06", action, request, nil, nil)
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
			if errs.IsExpectedErrors(err, []string{"InvalidAppId.NotFound"}) {
				return object, errs.WrapErrorf(errs.NotFoundErr("SAE:ApplicationScalingRule", id), errs.NotFoundMsg, errs.ProviderERROR)
			}
			return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
		}
		v, err := jsonpath.Get("$.Data.ApplicationScalingRules", response)
		if err != nil {
			return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.Data.ApplicationScalingRules", response)
		}
		if len(v.([]interface{})) < 1 {
			return object, errs.WrapErrorf(errs.NotFoundErr("SAE:ApplicationScalingRule", id), errs.NotFoundMsg, errs.ProviderERROR)
		}
		for _, v := range v.([]interface{}) {
			if fmt.Sprint(v.(map[string]interface{})["ScaleRuleName"]) == parts[1] {
				idExist = true
				return v.(map[string]interface{}), nil
			}
		}
		if len(v.([]interface{})) < helper.PageSizeLarge {
			break
		}
		currentPage, err := strconv.Atoi(*request["CurrentPage"])
		if err != nil {
			return object, errs.WrapError(err)
		}
		request["CurrentPage"] = helper.StringPointer(strconv.Itoa(currentPage + 1))
	}

	if !idExist {
		return object, errs.WrapErrorf(errs.NotFoundErr("SAE:ApplicationScalingRule", id), errs.NotFoundMsg, errs.ProviderERROR)
	}
	return object, nil
}

func (s *Service) DescribeSaeGreyTagRoute(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "/pop/v1/sam/tagroute/greyTagRoute"
	request := map[string]*string{
		"GreyTagRouteId": helper.StringPointer(id),
	}
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RoaGet("sae", "2019-05-06", action, request, nil, nil)
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
		if errs.IsExpectedErrors(err, []string{"InvalidParameter.WithMessage"}) {
			return object, errs.WrapErrorf(errs.NotFoundErr("SAE", id), errs.NotFoundWithResponse, response)
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.Data", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.Data", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}

func (s *Service) ListTagResources(id string, resourceType string) (object interface{}, err error) {
	client := s.client
	ids, err := json.Marshal([]string{id})
	if err != nil {
		return object, err
	}
	action := "/tags"
	request := map[string]*string{
		"RegionId":     helper.StringPointer(s.client.RegionId),
		"ResourceType": helper.StringPointer(resourceType),
		"ResourceIds":  helper.StringPointer(string(ids)),
	}
	tags := make([]interface{}, 0)
	var response map[string]interface{}

	for {
		wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(5*time.Minute, func() *retry.RetryError {
			response, err = client.RoaGet("sae", "2019-05-06", action, request, nil, nil)
			if err != nil {
				if errs.IsExpectedErrors(err, []string{errs.Throttling}) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		logs.AddDebug(action, response, request)
		if err != nil {
			err = errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
			return
		}
		v, err := jsonpath.Get("$.Data.TagResources", response)
		if err != nil {
			return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.Data.TagResources", response)
		}
		if v != nil {
			tags = append(tags, v.([]interface{})...)
		}
		if response["NextToken"] == nil {
			break
		}
		request["NextToken"] = helper.StringPointer(response["NextToken"].(string))
	}

	return tags, nil
}

func (s *Service) SetResourceTags(d *schema.ResourceData, resourceType string) error {

	if d.HasChange("tags") {
		added, removed := tags.ParsingTags(d)
		client := s.client
		ids, err := json.Marshal([]string{d.Id()})
		if err != nil {
			return err
		}

		removedTagKeys := make([]interface{}, 0)
		for _, v := range removed {
			if !tags.IgnoredTags(v, "") {
				removedTagKeys = append(removedTagKeys, v)
			}
		}
		if len(removedTagKeys) > 0 {
			action := "/tags"
			request := map[string]*string{
				"RegionId":     helper.StringPointer(s.client.RegionId),
				"ResourceType": helper.StringPointer(resourceType),
				"ResourceIds":  helper.StringPointer(string(ids)),
				"TagKeys":      helper.StringPointer(convert.ConvertListToJsonString(removedTagKeys)),
			}

			wait := helper.IncrementalWait(2*time.Second, 1*time.Second)
			err := retry.Retry(10*time.Minute, func() *retry.RetryError {
				response, err := client.RoaDelete("sae", "2019-05-06", action, request, nil, nil, false)
				if err != nil {
					if helper.NeedRetry(err) {
						wait()
						return retry.RetryableError(err)

					}
					return retry.NonRetryableError(err)
				}
				logs.AddDebug(action, response, request)
				return nil
			})
			if err != nil {
				return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
			}
		}
		if len(added) > 0 {
			action := "/tags"
			request := map[string]*string{
				"RegionId":     helper.StringPointer(s.client.RegionId),
				"ResourceType": helper.StringPointer(resourceType),
				"ResourceIds":  helper.StringPointer(string(ids)),
			}
			tags := make([]map[string]interface{}, len(added))
			for key, value := range added {
				tags = append(tags, map[string]interface{}{
					"key":   key,
					"value": value,
				})
			}
			jsonString, err := convert.ConvertListMapToJsonString(tags)
			if err != nil {
				return errs.WrapError(err)
			}
			request["Tags"] = helper.StringPointer(jsonString)

			wait := helper.IncrementalWait(2*time.Second, 1*time.Second)
			err = retry.Retry(10*time.Minute, func() *retry.RetryError {
				response, err := client.RoaPost("sae", "2019-05-06", action, request, nil, nil, false)
				if err != nil {
					if helper.NeedRetry(err) {
						wait()
						return retry.RetryableError(err)

					}
					return retry.NonRetryableError(err)
				}
				logs.AddDebug(action, response, request)
				return nil
			})
			if err != nil {
				return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
			}
		}
	}
	return nil
}

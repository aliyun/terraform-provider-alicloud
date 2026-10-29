package ddoscoo

import (
	"fmt"
	"strconv"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/convert"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Service struct {
	client *connectivity.AliyunClient
}

func NewService(client *connectivity.AliyunClient) *Service {
	return &Service{client: client}
}

func (s *Service) DescribeDdoscooInstance(id string) (object map[string]interface{}, err error) {
	client := s.client
	var response map[string]interface{}
	action := "DescribeInstances"
	request := map[string]interface{}{
		"InstanceIds": []string{id},
		"PageSize":    helper.PageSizeLarge,
		"PageNumber":  1,
	}

	idExist := false
	for {
		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(5*time.Minute, func() *retry.RetryError {
			response, err = client.RpcPost("ddoscoo", "2020-01-01", action, nil, request, true)
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
			if errs.IsExpectedErrors(err, []string{"InstanceNotFound"}) {
				return object, errs.WrapErrorf(errs.NotFoundErr("DdosCoo:Instance", id), errs.NotFoundWithResponse, response)
			}
			return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
		}

		resp, err := jsonpath.Get("$.Instances", response)
		if err != nil {
			return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.Instances", response)
		}

		if v, ok := resp.([]interface{}); !ok || len(v) < 1 {
			return object, errs.WrapErrorf(errs.NotFoundErr("DdosCoo:Instance", id), errs.NotFoundWithResponse, response)
		}

		for _, v := range resp.([]interface{}) {
			if fmt.Sprint(v.(map[string]interface{})["InstanceId"]) == id {
				idExist = true
				return v.(map[string]interface{}), nil
			}
		}

		if len(resp.([]interface{})) < request["PageSize"].(int) {
			break
		}

		request["PageNumber"] = request["PageNumber"].(int) + 1
	}

	if !idExist {
		return object, errs.WrapErrorf(errs.NotFoundErr("DdosCoo:Instance", id), errs.NotFoundWithResponse, response)
	}

	return object, nil
}

func (s *Service) DdosStateRefreshFunc(id string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeDdoscooInstance(id)
		if err != nil {
			if errs.NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", errs.WrapError(err)
		}
		status := ""
		actualStatus, err := convert.FormatInt(object["Status"])
		if err != nil {
			return nil, "", err
		}
		if actualStatus == 1 {
			status = "Available"
		} else {
			status = "Unavailable"
		}
		return object, status, nil
	}
}

func (s *Service) DescribeDdoscooInstanceSpec(id string) (object map[string]interface{}, err error) {
	client := s.client
	var response map[string]interface{}
	action := "DescribeInstanceSpecs"
	request := map[string]interface{}{
		"InstanceIds": []string{id},
	}

	idExist := false
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("ddoscoo", "2020-01-01", action, nil, request, true)
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
		if errs.IsExpectedErrors(err, []string{"InstanceNotFound", "ddos_coop3301"}) || errs.NotFoundError(err) {
			return object, errs.WrapErrorf(errs.NotFoundErr("DdosCoo:Instance", id), errs.NotFoundWithResponse, response)
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}

	resp, err := jsonpath.Get("$.InstanceSpecs", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.InstanceSpecs", response)
	}

	if v, ok := resp.([]interface{}); !ok || len(v) < 1 {
		return object, errs.WrapErrorf(errs.NotFoundErr("DdosCoo:Instance", id), errs.NotFoundWithResponse, response)
	}

	for _, v := range resp.([]interface{}) {
		if fmt.Sprint(v.(map[string]interface{})["InstanceId"]) == id {
			idExist = true
			return v.(map[string]interface{}), nil
		}
	}

	if !idExist {
		return object, errs.WrapErrorf(errs.NotFoundErr("DdosCoo:Instance", id), errs.NotFoundWithResponse, response)
	}

	return object, nil
}

func (s *Service) DescribeDdoscooInstanceExt(id string) (object map[string]interface{}, err error) {
	client := s.client
	var response map[string]interface{}
	action := "DescribeInstanceExt"
	request := map[string]interface{}{
		"InstanceId": id,
		"PageSize":   helper.PageSizeLarge,
		"PageNumber": 1,
	}

	idExist := false
	for {
		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(5*time.Minute, func() *retry.RetryError {
			response, err = client.RpcPost("ddoscoo", "2020-01-01", action, nil, request, true)
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

		resp, err := jsonpath.Get("$.InstanceExtSpecs", response)
		if err != nil {
			return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.InstanceExtSpecs", response)
		}

		if v, ok := resp.([]interface{}); !ok || len(v) < 1 {
			return object, errs.WrapErrorf(errs.NotFoundErr("DdosCoo:Instance", id), errs.NotFoundWithResponse, response)
		}

		for _, v := range resp.([]interface{}) {
			if fmt.Sprint(v.(map[string]interface{})["InstanceId"]) == id {
				idExist = true
				return v.(map[string]interface{}), nil
			}
		}

		if len(resp.([]interface{})) < request["PageSize"].(int) {
			break
		}

		request["PageNumber"] = request["PageNumber"].(int) + 1
	}

	if !idExist {
		return object, errs.WrapErrorf(errs.NotFoundErr("DdosCoo:Instance", id), errs.NotFoundWithResponse, response)
	}

	return object, nil
}

func (s *Service) UpdateInstanceSpec(schemaName string, specName string, d *schema.ResourceData, meta interface{}) (err error) {
	client := s.client
	var response map[string]interface{}
	action := "ModifyInstance"
	request := map[string]interface{}{
		"InstanceId":       d.Id(),
		"RegionId":         s.client.RegionId,
		"ProductCode":      "ddos",
		"ProductType":      "ddoscoo",
		"SubscriptionType": "Subscription",
		"ModifyType":       "Upgrade",
	}
	if d.Get("product_type").(string) == "ddoscoo_intl" {
		request["RegionId"] = "ap-southeast-1"
	} else {
		request["RegionId"] = "cn-hangzhou"
	}

	if v, ok := d.GetOk("product_type"); ok {
		request["ProductType"] = v.(string)
	}

	o, n := d.GetChange(schemaName)
	oi, _ := strconv.Atoi(o.(string))
	ni, _ := strconv.Atoi(n.(string))
	if ni < oi {
		request["ModifyType"] = "Downgrade"
	}

	request["Parameter"] = []map[string]string{
		{
			"Code":  specName,
			"Value": d.Get(schemaName).(string),
		},
	}

	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
		response, err = client.RpcPost("BssOpenApi", "2017-12-14", action, nil, request, true)
		if err != nil {
			if helper.NeedRetry(err) || errs.IsExpectedErrors(err, []string{"SYSTEM.CONCURRENT_OPERATE"}) {
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
	return nil
}

func (s *Service) DescribeDdoscooSchedulerRule(id string) (object map[string]interface{}, err error) {
	client := s.client
	var response map[string]interface{}
	action := "DescribeSchedulerRules"
	request := map[string]interface{}{
		"RegionId":   s.client.RegionId,
		"Domain":     id,
		"PageNumber": 1,
		"PageSize":   10,
	}
	idExist := false
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("ddoscoo", "2020-01-01", action, nil, request, true)
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
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.SchedulerRules", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.SchedulerRules", response)
	}
	if len(v.([]interface{})) < 1 {
		return object, errs.WrapErrorf(errs.NotFoundErr("DdosCoo", id), errs.NotFoundWithResponse, response)
	}

	for _, v := range v.([]interface{}) {
		if fmt.Sprint(v.(map[string]interface{})["RuleName"]) == id {
			idExist = true
			return v.(map[string]interface{}), nil
		}
	}

	if !idExist {
		return object, errs.WrapErrorf(errs.NotFoundErr("DdosCoo", id), errs.NotFoundWithResponse, response)
	}
	return object, nil
}

func (s *Service) DescribeDdoscooDomainResource(id string) (object map[string]interface{}, err error) {
	client := s.client
	var response map[string]interface{}
	action := "DescribeDomainResource"
	request := map[string]interface{}{
		"RegionId":   s.client.RegionId,
		"Domain":     id,
		"PageNumber": 1,
		"PageSize":   10,
	}
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("ddoscoo", "2020-01-01", action, nil, request, true)
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
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.WebRules", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.WebRules", response)
	}
	if len(v.([]interface{})) < 1 {
		return object, errs.WrapErrorf(errs.NotFoundErr("DdosCoo", id), errs.NotFoundWithResponse, response)
	} else {
		if v.([]interface{})[0].(map[string]interface{})["Domain"].(string) != id {
			return object, errs.WrapErrorf(errs.NotFoundErr("DdosCoo", id), errs.NotFoundWithResponse, response)
		}
	}
	object = v.([]interface{})[0].(map[string]interface{})
	return object, nil
}

func (s *Service) DescribeDdoscooPort(id string) (object map[string]interface{}, err error) {
	client := s.client
	var response map[string]interface{}
	action := "DescribePort"
	parts, err := helper.ParseResourceId(id, 3)
	if err != nil {
		err = errs.WrapError(err)
		return
	}
	request := map[string]interface{}{
		"RegionId":         s.client.RegionId,
		"FrontendPort":     parts[1],
		"FrontendProtocol": parts[2],
		"InstanceId":       parts[0],
		"PageNumber":       1,
		"PageSize":         10,
	}
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("ddoscoo", "2020-01-01", action, nil, request, true)
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
		if errs.IsExpectedErrors(err, []string{"anycast_controller3006"}) || errs.NotFoundError(err) {
			return object, errs.WrapErrorf(errs.NotFoundErr("DdosCoo", id), errs.NotFoundWithResponse, response)
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.NetworkRules", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.NetworkRules", response)
	}
	if len(v.([]interface{})) < 1 {
		return object, errs.WrapErrorf(errs.NotFoundErr("DdosCoo", id), errs.NotFoundWithResponse, response)
	} else {
		if v.([]interface{})[0].(map[string]interface{})["FrontendProtocol"].(string) != parts[2] {
			return object, errs.WrapErrorf(errs.NotFoundErr("DdosCoo", id), errs.NotFoundWithResponse, response)
		}
	}
	object = v.([]interface{})[0].(map[string]interface{})
	return object, nil
}

package dmsenterprise

import (
	"fmt"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tferrs "github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	tfhelper "github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	tflogs "github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type DmsEnterpriseService struct {
	client *connectivity.AliyunClient
}

func (s *DmsEnterpriseService) DescribeDmsEnterpriseInstance(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "GetInstance"
	parts, err := tfhelper.ParseResourceId(id, 2)
	if err != nil {
		err = tferrs.WrapError(err)
		return
	}
	request := map[string]interface{}{
		"RegionId": s.client.RegionId,
		"Host":     parts[0],
		"Port":     parts[1],
	}
	response, err = client.RpcPost("dms-enterprise", "2018-11-01", action, nil, request, true)
	if err != nil {
		if tferrs.IsExpectedErrors(err, []string{"InstanceNoEnoughNumber"}) {
			err = tferrs.WrapErrorf(tferrs.NotFoundErr("DmsEnterpriseInstance", id), tferrs.NotFoundMsg, tferrs.ProviderERROR)
			return object, err
		}
		err = tferrs.WrapErrorf(err, tferrs.DefaultErrorMsg, id, action, tferrs.AlibabaCloudSdkGoERROR)
		return object, err
	}
	tflogs.AddDebug(action, response, request)
	v, err := jsonpath.Get("$.Instance", response)
	if err != nil {
		return object, tferrs.WrapErrorf(err, tferrs.FailedGetAttributeMsg, id, "$.Instance", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}

func (s *DmsEnterpriseService) DescribeDmsEnterpriseUser(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "GetUser"
	request := map[string]interface{}{
		"RegionId": s.client.RegionId,
		"Uid":      id,
	}
	response, err = client.RpcPost("dms-enterprise", "2018-11-01", action, nil, request, true)
	if err != nil {
		err = tferrs.WrapErrorf(err, tferrs.DefaultErrorMsg, id, action, tferrs.AlibabaCloudSdkGoERROR)
		return
	}
	tflogs.AddDebug(action, response, request)
	v, err := jsonpath.Get("$.User", response)
	if err != nil {
		return object, tferrs.WrapErrorf(err, tferrs.FailedGetAttributeMsg, id, "$.User", response)
	}
	object = v.(map[string]interface{})
	if fmt.Sprint(object["State"]) == "DELETE" {
		return object, tferrs.GetNotFoundErrorFromString(tferrs.GetNotFoundMessage("DMS Enterprise User", id))
	}
	return object, nil
}

func (s *DmsEnterpriseService) DescribeDmsEnterpriseProxy(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "GetProxy"
	request := map[string]interface{}{
		"RegionId": s.client.RegionId,
		"ProxyId":  id,
	}
	wait := tfhelper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("dms-enterprise", "2018-11-01", action, nil, request, true)
		if err != nil {
			if tfhelper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	tflogs.AddDebug(action, response, request)
	if err != nil {
		if tferrs.IsExpectedErrors(err, []string{"InvalidParameterValid"}) {
			return object, tferrs.WrapErrorf(tferrs.NotFoundErr("DMSEnterprise:Proxy", id), tferrs.NotFoundMsg, tferrs.ProviderERROR, fmt.Sprint(response["RequestId"]))
		}
		return object, tferrs.WrapErrorf(err, tferrs.DefaultErrorMsg, id, action, tferrs.AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$", response)
	if err != nil {
		return object, tferrs.WrapErrorf(err, tferrs.FailedGetAttributeMsg, id, "$", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}

func (s *DmsEnterpriseService) DescribeDmsEnterpriseProxyAccess(id string) (object map[string]interface{}, err error) {
	client := s.client

	request := map[string]interface{}{
		"ProxyAccessId": id,
		"RegionId":      s.client.RegionId,
	}

	var response map[string]interface{}
	action := "GetProxyAccess"
	wait := tfhelper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		resp, err := client.RpcPost("dms-enterprise", "2018-11-01", action, nil, request, true)
		if err != nil {
			if tfhelper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		response = resp
		tflogs.AddDebug(action, response, request)
		return nil
	})
	if err != nil {
		if tferrs.IsExpectedErrors(err, []string{"InvalidParameterValid"}) {
			return object, tferrs.WrapErrorf(err, tferrs.NotFoundMsg, tferrs.AlibabaCloudSdkGoERROR)
		}
		return object, tferrs.WrapErrorf(err, tferrs.DefaultErrorMsg, id, action, tferrs.AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.ProxyAccess", response)
	if err != nil {
		return object, tferrs.WrapErrorf(err, tferrs.FailedGetAttributeMsg, id, "$.ProxyAccess", response)
	}
	return v.(map[string]interface{}), nil
}

func (s *DmsEnterpriseService) DmsEnterpriseProxyAccessStateRefreshFunc(d *schema.ResourceData, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeDmsEnterpriseProxyAccess(d.Id())
		if err != nil {
			if tferrs.NotFoundError(err) {
				return nil, "", nil
			}
			return nil, "", tferrs.WrapError(err)
		}
		for _, failState := range failStates {
			if fmt.Sprint(object[""]) == failState {
				return object, fmt.Sprint(object[""]), tferrs.WrapError(tferrs.Error(tferrs.FailedToReachTargetStatus, fmt.Sprint(object[""])))
			}
		}
		return object, fmt.Sprint(object[""]), nil
	}
}

func (s *DmsEnterpriseService) InspectProxyAccessSecret(id string) (object map[string]interface{}, err error) {
	client := s.client

	request := map[string]interface{}{
		"ProxyAccessId": id,
		"RegionId":      s.client.RegionId,
	}

	var response map[string]interface{}
	action := "InspectProxyAccessSecret"
	wait := tfhelper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		resp, err := client.RpcPost("dms-enterprise", "2018-11-01", action, nil, request, true)
		if err != nil {
			if tfhelper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		response = resp
		tflogs.AddDebug(action, response, request)
		return nil
	})
	if err != nil {
		if tferrs.IsExpectedErrors(err, []string{"InvalidParameterValid"}) {
			return object, tferrs.WrapErrorf(err, tferrs.NotFoundMsg, tferrs.AlibabaCloudSdkGoERROR)
		}
		return object, tferrs.WrapErrorf(err, tferrs.DefaultErrorMsg, id, action, tferrs.AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$", response)
	if err != nil {
		return object, tferrs.WrapErrorf(err, tferrs.FailedGetAttributeMsg, id, "$", response)
	}
	return v.(map[string]interface{}), nil
}

func (s *DmsEnterpriseService) DescribeDmsEnterpriseLogicDatabase(id string) (object map[string]interface{}, err error) {
	client := s.client

	request := map[string]interface{}{
		"DbId":     id,
		"RegionId": s.client.RegionId,
	}

	var response map[string]interface{}
	action := "GetLogicDatabase"
	wait := tfhelper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		resp, err := client.RpcPost("dms-enterprise", "2018-11-01", action, nil, request, true)
		if err != nil {
			if tfhelper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		response = resp
		tflogs.AddDebug(action, response, request)
		return nil
	})
	if err != nil {
		return object, tferrs.WrapErrorf(err, tferrs.DefaultErrorMsg, id, action, tferrs.AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.LogicDatabase", response)
	success, _ := jsonpath.Get("$.Success", response)
	if err != nil && success.(bool) {
		return object, tferrs.WrapErrorf(tferrs.NotFoundErr("DmsEnterprise", id), tferrs.NotFoundWithResponse, response)
	}
	return v.(map[string]interface{}), nil
}

func (s *DmsEnterpriseService) DmsEnterpriseLogicDatabaseStateRefreshFunc(d *schema.ResourceData, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeDmsEnterpriseLogicDatabase(d.Id())
		if err != nil {
			if tferrs.NotFoundError(err) {
				return nil, "", nil
			}
			return nil, "", tferrs.WrapError(err)
		}
		for _, failState := range failStates {
			if fmt.Sprint(object[""]) == failState {
				return object, fmt.Sprint(object[""]), tferrs.WrapError(tferrs.Error(tferrs.FailedToReachTargetStatus, fmt.Sprint(object[""])))
			}
		}
		return object, fmt.Sprint(object[""]), nil
	}
}

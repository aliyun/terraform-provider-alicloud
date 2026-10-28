package ecp

import (
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

// Service wraps the Cloudphone (ECP) API calls used by the resources and
// data sources of this package.
type Service struct {
	client *connectivity.AliyunClient
}

// NewService is the exported constructor: the client field stays unexported,
// so callers outside the package (including the flat layout's acceptance
// tests) build values with it.
func NewService(client *connectivity.AliyunClient) *Service {
	return &Service{client: client}
}

func (s *Service) DescribeEcpKeyPair(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "ListKeyPairs"
	request := map[string]interface{}{
		"RegionId":    s.client.RegionId,
		"KeyPairName": id,
	}
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("cloudphone", "2020-12-30", action, nil, request, true)
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
	v, err := jsonpath.Get("$.KeyPairs.KeyPair", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.KeyPairs.KeyPair", response)
	}
	if len(v.([]interface{})) < 1 {
		return object, errs.WrapErrorf(errs.NotFoundErr("ECP", id), errs.NotFoundWithResponse, response)
	} else {
		if fmt.Sprint(v.([]interface{})[0].(map[string]interface{})["KeyPairName"]) != id {
			return object, errs.WrapErrorf(errs.NotFoundErr("ECP", id), errs.NotFoundWithResponse, response)
		}
	}
	object = v.([]interface{})[0].(map[string]interface{})
	return object, nil
}

func (s *Service) DescribeEcpInstance(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	action := "ListInstances"

	client := s.client

	request := map[string]interface{}{
		"RegionId":   s.client.RegionId,
		"InstanceId": []string{id},
		"MaxResults": helper.PageSizeLarge,
	}

	idExist := false
	for {
		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(5*time.Minute, func() *retry.RetryError {
			response, err = client.RpcPost("cloudphone", "2020-12-30", action, nil, request, true)
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

		resp, err := jsonpath.Get("$.Instances.Instance", response)
		if err != nil {
			return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.Instances.Instance", response)
		}

		if v, ok := resp.([]interface{}); !ok || len(v) < 1 {
			return object, errs.WrapErrorf(errs.NotFoundErr("Ecp:Instance", id), errs.NotFoundWithResponse, response)
		}

		for _, v := range resp.([]interface{}) {
			if fmt.Sprint(v.(map[string]interface{})["InstanceId"]) == id {
				idExist = true
				return v.(map[string]interface{}), nil
			}
		}

		if nextToken, ok := response["NextToken"].(string); ok && nextToken != "" {
			request["NextToken"] = nextToken
		} else {
			break
		}
	}

	if !idExist {
		return object, errs.WrapErrorf(errs.NotFoundErr("Ecp:Instance", id), errs.NotFoundWithResponse, response)
	}

	return object, nil
}

func (s *Service) EcpInstanceStateRefreshFunc(id string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeEcpInstance(id)
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

func (s *Service) ModifyEcpInstanceStatus(d *schema.ResourceData, status string) (err error) {
	var response map[string]interface{}
	client := s.client

	switch status {
	case "Running":
		action := "StartInstances"

		startInstanceReq := map[string]interface{}{
			"RegionId":   s.client.RegionId,
			"InstanceId": []string{d.Id()},
		}
		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(20*time.Minute, func() *retry.RetryError {
			response, err = client.RpcPost("cloudphone", "2020-12-30", action, nil, startInstanceReq, true)
			if err != nil {
				if helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		logs.AddDebug(action, response, startInstanceReq)

		if err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
		}

		stateConf := helper.BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, s.EcpInstanceStateRefreshFunc(d.Id(), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return errs.WrapErrorf(err, errs.IdMsg, d.Id())
		}
	case "Stopped":
		action := "StopInstances"

		stopInstanceReq := map[string]interface{}{
			"RegionId":   s.client.RegionId,
			"InstanceId": []string{d.Id()},
		}

		if v, ok := d.GetOkExists("force"); ok {
			stopInstanceReq["Force"] = v
		}

		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(20*time.Minute, func() *retry.RetryError {
			response, err = client.RpcPost("cloudphone", "2020-12-30", action, nil, stopInstanceReq, true)
			if err != nil {
				if helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		logs.AddDebug(action, response, stopInstanceReq)

		if err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
		}

		stateConf := helper.BuildStateConf([]string{}, []string{"Stopped"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, s.EcpInstanceStateRefreshFunc(d.Id(), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return errs.WrapErrorf(err, errs.IdMsg, d.Id())
		}
	}

	return nil
}

func (s *Service) ListInstanceVncUrl(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "ListInstanceVncUrl"
	request := map[string]interface{}{
		"RegionId":   s.client.RegionId,
		"InstanceId": id,
	}
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("cloudphone", "2020-12-30", action, nil, request, true)
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
	v, err := jsonpath.Get("$", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}

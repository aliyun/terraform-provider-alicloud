package alicloud

import (
	"fmt"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

type CloudphoneService struct {
	client *connectivity.AliyunClient
}

func (s *CloudphoneService) DescribeEcpKeyPair(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "ListKeyPairs"
	request := map[string]interface{}{
		"RegionId":    s.client.RegionId,
		"KeyPairName": id,
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("cloudphone", "2020-12-30", action, nil, request, true)
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
	v, err := jsonpath.Get("$.KeyPairs.KeyPair", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.KeyPairs.KeyPair", response)
	}
	if len(v.([]interface{})) < 1 {
		return object, WrapErrorf(Error(GetNotFoundMessage("ECP", id)), NotFoundWithResponse, response)
	} else {
		if fmt.Sprint(v.([]interface{})[0].(map[string]interface{})["KeyPairName"]) != id {
			return object, WrapErrorf(Error(GetNotFoundMessage("ECP", id)), NotFoundWithResponse, response)
		}
	}
	object = v.([]interface{})[0].(map[string]interface{})
	return object, nil
}

func (s *CloudphoneService) DescribeEcpInstance(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	action := "ListInstances"

	client := s.client

	request := map[string]interface{}{
		"RegionId":   s.client.RegionId,
		"InstanceId": []string{id},
		"MaxResults": PageSizeLarge,
	}

	idExist := false
	for {
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = client.RpcPost("cloudphone", "2020-12-30", action, nil, request, true)
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

		resp, err := jsonpath.Get("$.Instances.Instance", response)
		if err != nil {
			return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Instances.Instance", response)
		}

		if v, ok := resp.([]interface{}); !ok || len(v) < 1 {
			return object, WrapErrorf(Error(GetNotFoundMessage("Ecp:Instance", id)), NotFoundWithResponse, response)
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
		return object, WrapErrorf(Error(GetNotFoundMessage("Ecp:Instance", id)), NotFoundWithResponse, response)
	}

	return object, nil
}

func (s *CloudphoneService) EcpInstanceStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeEcpInstance(id)
		if err != nil {
			if NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
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

func (s *CloudphoneService) ModifyEcpInstanceStatus(d *schema.ResourceData, status string) (err error) {
	var response map[string]interface{}
	client := s.client

	switch status {
	case "Running":
		action := "StartInstances"

		startInstanceReq := map[string]interface{}{
			"RegionId":   s.client.RegionId,
			"InstanceId": []string{d.Id()},
		}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(20*time.Minute, func() *resource.RetryError {
			response, err = client.RpcPost("cloudphone", "2020-12-30", action, nil, startInstanceReq, true)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, startInstanceReq)

		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}

		stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, s.EcpInstanceStateRefreshFunc(d.Id(), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
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

		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(20*time.Minute, func() *resource.RetryError {
			response, err = client.RpcPost("cloudphone", "2020-12-30", action, nil, stopInstanceReq, true)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, stopInstanceReq)

		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}

		stateConf := BuildStateConf([]string{}, []string{"Stopped"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, s.EcpInstanceStateRefreshFunc(d.Id(), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}

	return nil
}

func (s *CloudphoneService) ListInstanceVncUrl(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "ListInstanceVncUrl"
	request := map[string]interface{}{
		"RegionId":   s.client.RegionId,
		"InstanceId": id,
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("cloudphone", "2020-12-30", action, nil, request, true)
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
	v, err := jsonpath.Get("$", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}

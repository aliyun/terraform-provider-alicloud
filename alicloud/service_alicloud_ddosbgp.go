package alicloud

import (
	"fmt"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

type DdosbgpService struct {
	client *connectivity.AliyunClient
}

func (s *DdosbgpService) DescribeDdosbgpInstance(id string) (object map[string]interface{}, err error) {
	client := s.client
	var response map[string]interface{}
	action := "DescribeInstanceList"
	request := map[string]interface{}{
		"InstanceIdList": "[\"" + id + "\"]",
		"RegionId":       s.client.RegionId,
		"PageNo":         "1",
		"PageSize":       "10",
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(6*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("ddosbgp", "2018-07-20", action, nil, request, true)
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
	addDebug(action, response, request)
	if err != nil {
		if IsExpectedErrors(err, []string{"InstanceNotFound", "InvalidInstance"}) {
			return object, WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	v, _ := jsonpath.Get("$.InstanceList", response)
	if v == nil {
		return object, nil
	}
	if len(v.([]interface{})) < 1 || fmt.Sprint(v.([]interface{})[0].(map[string]interface{})["InstanceId"]) != id {
		return object, WrapErrorf(NotFoundErr("DdosBgp", id), NotFoundWithResponse, response)
	}
	object = v.([]interface{})[0].(map[string]interface{})
	return object, nil
}

func (s *DdosbgpService) DescribeDdosbgpInstanceSpec(id string) (object map[string]interface{}, err error) {
	client := s.client
	var response map[string]interface{}
	action := "DescribeInstanceSpecs"
	request := map[string]interface{}{
		"InstanceIdList": "[\"" + id + "\"]",
		"RegionId":       s.client.RegionId,
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(6*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("ddosbgp", "2018-07-20", action, nil, request, true)
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
	addDebug(action, response, request)
	if err != nil {
		if IsExpectedErrors(err, []string{"InstanceNotFound", "InvalidInstance"}) {
			return object, WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	v, _ := jsonpath.Get("$.InstanceSpecs", response)
	if v == nil {
		return object, nil
	}
	if len(v.([]interface{})) < 1 {
		return object, WrapErrorf(NotFoundErr("Ddos", id), NotFoundWithResponse, response)
	}
	object = v.([]interface{})[0].(map[string]interface{})
	return object, nil
}

func (s *DdosbgpService) DescribeDdosbgpIp(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "DescribePackIpList"
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		err = WrapError(err)
		return
	}
	request := map[string]interface{}{
		"RegionId":   s.client.RegionId,
		"InstanceId": parts[0],
		"Ip":         parts[1],
		"PageNo":     "1",
		"PageSize":   PageSizeSmall,
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("ddosbgp", "2018-07-20", action, nil, request, true)
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
	if fmt.Sprint(response["Success"]) == "false" {
		return object, WrapError(fmt.Errorf("%s failed, response: %v", action, response))
	}
	v, err := jsonpath.Get("$.IpList", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.IpList", response)
	}
	if len(v.([]interface{})) == 0 {
		return object, WrapErrorf(NotFoundErr("Ddosbgp Instance", id), NotFoundMsg, ProviderERROR)
	}
	object = v.([]interface{})[0].(map[string]interface{})
	return object, nil
}

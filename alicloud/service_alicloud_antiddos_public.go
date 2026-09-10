package alicloud

import (
	"fmt"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

type AntiddosPublicService struct {
	client *connectivity.AliyunClient
}

func (s *AntiddosPublicService) DescribeDdosBasicAntiddos(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "DescribeDdosThreshold"
	parts, err := ParseResourceId(id, 3)
	if err != nil {
		err = WrapError(err)
		return
	}
	request := map[string]interface{}{
		"DdosRegionId": s.client.RegionId,
		"InstanceIds":  []string{parts[0]},
		"InstanceType": parts[1],
	}
	request["DdosType"] = parts[2]
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("antiddos-public", "2017-05-18", action, nil, request, true)
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
	v, err := jsonpath.Get("$.Thresholds.Threshold", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Thresholds.Threshold", response)
	}
	if len(v.([]interface{})) < 1 {
		return object, WrapErrorf(NotFoundErr("DdosBasic", id), NotFoundWithResponse, response)
	} else {
		if fmt.Sprint(v.([]interface{})[0].(map[string]interface{})["InstanceId"]) != parts[0] {
			return object, WrapErrorf(NotFoundErr("DdosBasic", id), NotFoundWithResponse, response)
		}
	}
	object = v.([]interface{})[0].(map[string]interface{})
	return object, nil
}

func (s *AntiddosPublicService) DdosBasicAntiDdosStateRefreshFunc(id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeDdosBasicAntiddos(id)
		if err != nil {
			if NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}
		return object, fmt.Sprint(object["DdosType"]), nil
	}
}

func (s *AntiddosPublicService) DescribeDdosBasicThreshold(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	action := "DescribeIpDdosThreshold"
	client := s.client
	parts, err := ParseResourceId(id, 3)
	if err != nil {
		return nil, WrapError(err)
	}

	request := map[string]interface{}{
		"DdosRegionId": s.client.RegionId,
		"DdosType":     "defense",
		"InstanceType": parts[0],
		"InstanceId":   parts[1],
		"InternetIp":   parts[2],
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("antiddos-public", "2017-05-18", action, nil, request, true)
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

	v, err := jsonpath.Get("$.Threshold", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Threshold", response)
	}

	object = v.(map[string]interface{})

	return object, nil
}

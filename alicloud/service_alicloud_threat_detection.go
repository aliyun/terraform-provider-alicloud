package alicloud

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

type ThreatDetectionService struct {
	client *connectivity.AliyunClient
}

func (s *ThreatDetectionService) DescribeThreatDetectionBackupPolicy(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	action := "DescribeBackupPolicy"

	client := s.client

	request := map[string]interface{}{
		"Id": id,
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("Sas", "2018-12-03", action, nil, request, true)
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
		if IsExpectedErrors(err, []string{"DataNotExist", "InvalidId"}) {
			return nil, WrapErrorf(NotFoundErr("TDS:BackupPolicy", id), NotFoundMsg, ProviderERROR)
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$.BackupPolicyDetail", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.BackupPolicyDetail", response)
	}

	object = v.(map[string]interface{})

	return object, nil
}

func (s *ThreatDetectionService) ThreatDetectionBackupPolicyStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeThreatDetectionBackupPolicy(id)
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

func (s *ThreatDetectionService) DescribeThreatDetectionVulWhitelist(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	action := "GetVulWhitelist"
	client := s.client
	request := map[string]interface{}{
		"VulWhitelistId": id,
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("Sas", "2018-12-03", action, nil, request, true)
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
		if IsExpectedErrors(err, []string{"VulWhitelistNotExist", "InvalidId"}) {
			return nil, WrapErrorf(NotFoundErr("TDS:VulWhitelist", id), NotFoundMsg, ProviderERROR)
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$.VulWhitelist", response)
	if err != nil {
		if fmt.Sprint(response["Code"]) == "VulWhitelistNotExist" {
			return nil, WrapErrorf(NotFoundErr("TDS:VulWhitelist", id), NotFoundMsg, ProviderERROR)
		}
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.VulWhitelist", response)
	}

	object = v.(map[string]interface{})

	return object, nil
}

func (s *ThreatDetectionService) DescribeThreatDetectionHoneypotNode(id string) (object map[string]interface{}, err error) {
	client := s.client

	request := map[string]interface{}{
		"NodeId": id,
	}

	var response map[string]interface{}
	action := "GetHoneypotNode"
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		resp, err := client.RpcPost("Sas", "2018-12-03", action, nil, request, true)
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
		if IsExpectedErrors(err, []string{"NodeNotExist"}) {
			return object, WrapErrorf(NotFoundErr("HoneypotNode", id), NotFoundWithResponse, response)
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.HoneypotNode", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.HoneypotNode", response)
	}

	if fmt.Sprint(v.(map[string]interface{})["TotalStatus"]) == "-1" {
		return object, WrapErrorf(NotFoundErr("HoneypotNode", id), NotFoundWithResponse, response)
	}

	return v.(map[string]interface{}), nil
}

func (s *ThreatDetectionService) ThreatDetectionHoneypotNodeStateRefreshFunc(d *schema.ResourceData, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeThreatDetectionHoneypotNode(d.Id())
		if err != nil {
			if NotFoundError(err) {
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}
		for _, failState := range failStates {
			if fmt.Sprint(object["TotalStatus"]) == failState {
				return object, fmt.Sprint(object["TotalStatus"]), WrapError(Error(FailedToReachTargetStatus, fmt.Sprint(object["TotalStatus"])))
			}
		}
		return object, fmt.Sprint(object["TotalStatus"]), nil
	}
}

func (s *ThreatDetectionService) DescribeThreatDetectionLogShipper(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	action := "DescribeLogShipperStatus"

	client := s.client

	request := map[string]interface{}{
		"From": "sas",
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("Sas", "2018-12-03", action, nil, request, true)
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

	v, err := jsonpath.Get("$.LogShipperStatus", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.LogShipperStatus", response)
	}

	object = v.(map[string]interface{})

	return object, nil
}

func (s *ThreatDetectionService) DescribeThreatDetectionHoneypotPreset(id string) (object map[string]interface{}, err error) {
	client := s.client

	request := map[string]interface{}{
		"HoneypotPresetId": id,
	}

	var response map[string]interface{}
	action := "GetHoneypotPreset"
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		resp, err := client.RpcPost("Sas", "2018-12-03", action, nil, request, true)
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
		if IsExpectedErrors(err, []string{"HoneypotPresetNotExist"}) {
			return object, WrapErrorf(NotFoundErr("HoneypotPreset", id), NotFoundWithResponse, response)
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.Data", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Data", response)
	}
	return v.(map[string]interface{}), nil
}

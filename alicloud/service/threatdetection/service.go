package threatdetection

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

type Service struct {
	client *connectivity.AliyunClient
}

func NewService(client *connectivity.AliyunClient) *Service {
	return &Service{client: client}
}

func (s *Service) DescribeThreatDetectionBackupPolicy(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	action := "DescribeBackupPolicy"

	client := s.client

	request := map[string]interface{}{
		"Id": id,
	}

	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("Sas", "2018-12-03", action, nil, request, true)
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
		if errs.IsExpectedErrors(err, []string{"DataNotExist", "InvalidId"}) {
			return nil, errs.WrapErrorf(errs.NotFoundErr("TDS:BackupPolicy", id), errs.NotFoundMsg, errs.ProviderERROR)
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$.BackupPolicyDetail", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.BackupPolicyDetail", response)
	}

	object = v.(map[string]interface{})

	return object, nil
}

func (s *Service) ThreatDetectionBackupPolicyStateRefreshFunc(id string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeThreatDetectionBackupPolicy(id)
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

func (s *Service) DescribeThreatDetectionVulWhitelist(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	action := "GetVulWhitelist"
	client := s.client
	request := map[string]interface{}{
		"VulWhitelistId": id,
	}

	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("Sas", "2018-12-03", action, nil, request, true)
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
		if errs.IsExpectedErrors(err, []string{"VulWhitelistNotExist", "InvalidId"}) {
			return nil, errs.WrapErrorf(errs.NotFoundErr("TDS:VulWhitelist", id), errs.NotFoundMsg, errs.ProviderERROR)
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$.VulWhitelist", response)
	if err != nil {
		if fmt.Sprint(response["Code"]) == "VulWhitelistNotExist" {
			return nil, errs.WrapErrorf(errs.NotFoundErr("TDS:VulWhitelist", id), errs.NotFoundMsg, errs.ProviderERROR)
		}
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.VulWhitelist", response)
	}

	object = v.(map[string]interface{})

	return object, nil
}

func (s *Service) DescribeThreatDetectionHoneypotNode(id string) (object map[string]interface{}, err error) {
	client := s.client

	request := map[string]interface{}{
		"NodeId": id,
	}

	var response map[string]interface{}
	action := "GetHoneypotNode"
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		resp, err := client.RpcPost("Sas", "2018-12-03", action, nil, request, true)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		response = resp
		logs.AddDebug(action, response, request)
		return nil
	})
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"NodeNotExist"}) {
			return object, errs.WrapErrorf(errs.NotFoundErr("HoneypotNode", id), errs.NotFoundWithResponse, response)
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.HoneypotNode", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.HoneypotNode", response)
	}

	if fmt.Sprint(v.(map[string]interface{})["TotalStatus"]) == "-1" {
		return object, errs.WrapErrorf(errs.NotFoundErr("HoneypotNode", id), errs.NotFoundWithResponse, response)
	}

	return v.(map[string]interface{}), nil
}

func (s *Service) ThreatDetectionHoneypotNodeStateRefreshFunc(d *schema.ResourceData, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeThreatDetectionHoneypotNode(d.Id())
		if err != nil {
			if errs.NotFoundError(err) {
				return nil, "", nil
			}
			return nil, "", errs.WrapError(err)
		}
		for _, failState := range failStates {
			if fmt.Sprint(object["TotalStatus"]) == failState {
				return object, fmt.Sprint(object["TotalStatus"]), errs.WrapError(errs.Error(errs.FailedToReachTargetStatus, fmt.Sprint(object["TotalStatus"])))
			}
		}
		return object, fmt.Sprint(object["TotalStatus"]), nil
	}
}

func (s *Service) DescribeThreatDetectionLogShipper(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	action := "DescribeLogShipperStatus"

	client := s.client

	request := map[string]interface{}{
		"From": "sas",
	}

	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("Sas", "2018-12-03", action, nil, request, true)
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

	v, err := jsonpath.Get("$.LogShipperStatus", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.LogShipperStatus", response)
	}

	object = v.(map[string]interface{})

	return object, nil
}

func (s *Service) DescribeThreatDetectionHoneypotPreset(id string) (object map[string]interface{}, err error) {
	client := s.client

	request := map[string]interface{}{
		"HoneypotPresetId": id,
	}

	var response map[string]interface{}
	action := "GetHoneypotPreset"
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		resp, err := client.RpcPost("Sas", "2018-12-03", action, nil, request, true)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		response = resp
		logs.AddDebug(action, response, request)
		return nil
	})
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"HoneypotPresetNotExist"}) {
			return object, errs.WrapErrorf(errs.NotFoundErr("HoneypotPreset", id), errs.NotFoundWithResponse, response)
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.Data", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.Data", response)
	}
	return v.(map[string]interface{}), nil
}

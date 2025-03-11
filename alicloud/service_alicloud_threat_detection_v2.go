package alicloud

import (
	"fmt"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

type ThreatDetectionServiceV2 struct {
	client *connectivity.AliyunClient
}

// DescribeThreatDetectionInstance <<< Encapsulated get interface for ThreatDetection Instance.

func (s *ThreatDetectionServiceV2) DescribeThreatDetectionInstance(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	var endpoint string
	action := "GetInstance"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	query["InstanceId"] = id

	request["CommodityCode"] = "sas"
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(1*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPostWithEndpoint("Sas", "2018-12-03", action, query, request, true, endpoint)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			if !client.IsInternationalAccount() && IsExpectedErrors(err, []string{"NODATA"}) {
				endpoint = connectivity.SaSOpenAPIEndpointInternational
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		if IsExpectedErrors(err, []string{}) {
			return object, WrapErrorf(Error(GetNotFoundMessage("Instance", id)), NotFoundMsg, ProviderERROR, fmt.Sprint(response["RequestId"]))
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}

	currentStatus := response["InstanceId"]
	if currentStatus == nil {
		return object, WrapErrorf(Error(GetNotFoundMessage("Instance", id)), NotFoundMsg, response)
	}

	return response, nil
}
func (s *ThreatDetectionServiceV2) DescribeQueryAvailableInstances(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var endpoint string
	var query map[string]interface{}
	action := "QueryAvailableInstances"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	query["InstanceIDs"] = id

	request["ProductCode"] = "sas"
	request["ProductType"] = "sas"
	if client.IsInternationalAccount() {
		request["ProductType"] = ""
	}
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(1*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPostWithEndpoint("BssOpenApi", "2017-12-14", action, query, request, true, endpoint)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			if !client.IsInternationalAccount() && IsExpectedErrors(err, []string{"NotApplicable"}) {
				request["ProductType"] = ""
				endpoint = connectivity.BssOpenAPIEndpointInternational
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})
	if err != nil {
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$.Data.InstanceList[*]", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Data.InstanceList[*]", response)
	}

	if len(v.([]interface{})) == 0 {
		return object, WrapErrorf(Error(GetNotFoundMessage("Instance", id)), NotFoundMsg, ProviderERROR, fmt.Sprint(response["RequestId"]))
	}

	currentStatus := v.([]interface{})[0].(map[string]interface{})["InstanceID"]
	if currentStatus == nil {
		return object, WrapErrorf(Error(GetNotFoundMessage("Instance", id)), NotFoundMsg, ProviderERROR, fmt.Sprint(response["RequestId"]))
	}

	return v.([]interface{})[0].(map[string]interface{}), nil
}

func (s *ThreatDetectionServiceV2) ThreatDetectionInstanceStateRefreshFunc(id string, field string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeThreatDetectionInstance(id)
		if err != nil {
			if NotFoundError(err) {
				return object, "", nil
			}
			return nil, "", WrapError(err)
		}

		v, err := jsonpath.Get(field, object)
		currentStatus := fmt.Sprint(v)

		for _, failState := range failStates {
			if currentStatus == failState {
				return object, currentStatus, WrapError(Error(FailedToReachTargetStatus, currentStatus))
			}
		}
		return object, currentStatus, nil
	}
}

// DescribeThreatDetectionInstance >>> Encapsulated.

// DescribeThreatDetectionClientUserDefineRule <<< Encapsulated get interface for ThreatDetection ClientUserDefineRule.

func (s *ThreatDetectionServiceV2) DescribeThreatDetectionClientUserDefineRule(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	action := "GetClientUserDefineRule"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	query["Id"] = id

	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(1*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("Sas", "2018-12-03", action, query, request, true)

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
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$.UserDefineRuleDetail", response)
	if err != nil {
		return object, WrapErrorf(Error(GetNotFoundMessage("ClientUserDefineRule", id)), NotFoundMsg, response)
	}

	currentStatus := v.(map[string]interface{})["Id"]
	if currentStatus == nil {
		return object, WrapErrorf(Error(GetNotFoundMessage("ClientUserDefineRule", id)), NotFoundMsg, response)
	}

	return v.(map[string]interface{}), nil
}

func (s *ThreatDetectionServiceV2) ThreatDetectionClientUserDefineRuleStateRefreshFunc(id string, field string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeThreatDetectionClientUserDefineRule(id)
		if err != nil {
			if NotFoundError(err) {
				return object, "", nil
			}
			return nil, "", WrapError(err)
		}

		v, err := jsonpath.Get(field, object)
		currentStatus := fmt.Sprint(v)
		for _, failState := range failStates {
			if currentStatus == failState {
				return object, currentStatus, WrapError(Error(FailedToReachTargetStatus, currentStatus))
			}
		}
		return object, currentStatus, nil
	}
}

// DescribeThreatDetectionClientUserDefineRule >>> Encapsulated.

// DescribeThreatDetectionLogMeta <<< Encapsulated get interface for ThreatDetection LogMeta.

func (s *ThreatDetectionServiceV2) DescribeThreatDetectionLogMeta(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["LogStore"] = id

	action := "GetLogMeta"

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(1*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("Sas", "2018-12-03", action, query, request, true)

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

	v, err := jsonpath.Get("$.LogMeta", response)
	if err != nil {
		return object, WrapErrorf(Error(GetNotFoundMessage("LogMeta", id)), NotFoundMsg, response)
	}

	currentStatus := v.(map[string]interface{})["LogStore"]
	if currentStatus == nil {
		return object, WrapErrorf(Error(GetNotFoundMessage("LogMeta", id)), NotFoundMsg, response)
	}

	return v.(map[string]interface{}), nil
}

func (s *ThreatDetectionServiceV2) ThreatDetectionLogMetaStateRefreshFunc(id string, field string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeThreatDetectionLogMeta(id)
		if err != nil {
			if NotFoundError(err) {
				return object, "", nil
			}
			return nil, "", WrapError(err)
		}

		v, err := jsonpath.Get(field, object)
		currentStatus := fmt.Sprint(v)

		if strings.HasPrefix(field, "#") {
			v, _ := jsonpath.Get(strings.TrimPrefix(field, "#"), object)
			if v != nil {
				currentStatus = "#CHECKSET"
			}
		}

		for _, failState := range failStates {
			if currentStatus == failState {
				return object, currentStatus, WrapError(Error(FailedToReachTargetStatus, currentStatus))
			}
		}
		return object, currentStatus, nil
	}
}

// DescribeThreatDetectionLogMeta >>> Encapsulated.

// DescribeThreatDetectionClientFileProtect <<< Encapsulated get interface for ThreatDetection ClientFileProtect.

func (s *ThreatDetectionServiceV2) DescribeThreatDetectionClientFileProtect(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	action := "GetFileProtectRule"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	query["Id"] = id

	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(1*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("Sas", "2018-12-03", action, query, request, true)

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
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$.Data", response)
	if err != nil {
		return object, WrapErrorf(Error(GetNotFoundMessage("ClientFileProtect", id)), NotFoundMsg, response)
	}

	return v.(map[string]interface{}), nil
}

func (s *ThreatDetectionServiceV2) ThreatDetectionClientFileProtectStateRefreshFunc(id string, field string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeThreatDetectionClientFileProtect(id)
		if err != nil {
			if NotFoundError(err) {
				return object, "", nil
			}
			return nil, "", WrapError(err)
		}

		v, err := jsonpath.Get(field, object)
		currentStatus := fmt.Sprint(v)
		for _, failState := range failStates {
			if currentStatus == failState {
				return object, currentStatus, WrapError(Error(FailedToReachTargetStatus, currentStatus))
			}
		}
		return object, currentStatus, nil
	}
}

// DescribeThreatDetectionClientFileProtect >>> Encapsulated.

// DescribeThreatDetectionFileUploadLimit <<< Encapsulated get interface for ThreatDetection FileUploadLimit.

func (s *ThreatDetectionServiceV2) DescribeThreatDetectionFileUploadLimit(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	parts := strings.Split(id, ":")
	if len(parts) != 0 {
		err = WrapError(fmt.Errorf("invalid Resource Id %s. Expected parts' length %d, got %d", id, 0, len(parts)))
	}
	action := "GetFileUploadLimit"
	request = make(map[string]interface{})
	query = make(map[string]interface{})

	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(1*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("Sas", "2018-12-03", action, query, request, true)

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
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$.Data", response)
	if err != nil {
		return object, WrapErrorf(Error(GetNotFoundMessage("FileUploadLimit", id)), NotFoundMsg, response)
	}

	currentStatus := v.(map[string]interface{})["Limit"]
	if currentStatus == nil {
		return object, WrapErrorf(Error(GetNotFoundMessage("FileUploadLimit", id)), NotFoundMsg, response)
	}

	return v.(map[string]interface{}), nil
}

func (s *ThreatDetectionServiceV2) ThreatDetectionFileUploadLimitStateRefreshFunc(id string, field string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeThreatDetectionFileUploadLimit(id)
		if err != nil {
			if NotFoundError(err) {
				return object, "", nil
			}
			return nil, "", WrapError(err)
		}

		v, err := jsonpath.Get(field, object)
		currentStatus := fmt.Sprint(v)
		for _, failState := range failStates {
			if currentStatus == failState {
				return object, currentStatus, WrapError(Error(FailedToReachTargetStatus, currentStatus))
			}
		}
		return object, currentStatus, nil
	}
}

// DescribeThreatDetectionFileUploadLimit >>> Encapsulated.

// DescribeThreatDetectionMaliciousFileWhitelistConfig <<< Encapsulated get interface for ThreatDetection MaliciousFileWhitelistConfig.

func (s *ThreatDetectionServiceV2) DescribeThreatDetectionMaliciousFileWhitelistConfig(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	action := "GetMaliciousFileWhitelistConfig"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	query["ConfigId"] = id

	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(1*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("Sas", "2018-12-03", action, query, request, true)

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
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$.Data", response)
	if err != nil {
		return object, WrapErrorf(Error(GetNotFoundMessage("MaliciousFileWhitelistConfig", id)), NotFoundMsg, response)
	}

	return v.(map[string]interface{}), nil
}

func (s *ThreatDetectionServiceV2) ThreatDetectionMaliciousFileWhitelistConfigStateRefreshFunc(id string, field string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeThreatDetectionMaliciousFileWhitelistConfig(id)
		if err != nil {
			if NotFoundError(err) {
				return object, "", nil
			}
			return nil, "", WrapError(err)
		}

		v, err := jsonpath.Get(field, object)
		currentStatus := fmt.Sprint(v)

		for _, failState := range failStates {
			if currentStatus == failState {
				return object, currentStatus, WrapError(Error(FailedToReachTargetStatus, currentStatus))
			}
		}
		return object, currentStatus, nil
	}
}

// DescribeThreatDetectionMaliciousFileWhitelistConfig >>> Encapsulated.
// DescribeThreatDetectionImageEventOperation <<< Encapsulated get interface for ThreatDetection ImageEventOperation.

func (s *ThreatDetectionServiceV2) DescribeThreatDetectionImageEventOperation(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	action := "GetImageEventOperation"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	query["Id"] = id

	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(1*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("Sas", "2018-12-03", action, query, request, true)

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
		if IsExpectedErrors(err, []string{"DataNotExists"}) {
			return object, WrapErrorf(Error(GetNotFoundMessage("ImageEventOperation", id)), NotFoundMsg, response)
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$.Data", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Data", response)
	}

	return v.(map[string]interface{}), nil
}

func (s *ThreatDetectionServiceV2) ThreatDetectionImageEventOperationStateRefreshFunc(id string, field string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeThreatDetectionImageEventOperation(id)
		if err != nil {
			if NotFoundError(err) {
				return object, "", nil
			}
			return nil, "", WrapError(err)
		}

		v, err := jsonpath.Get(field, object)
		currentStatus := fmt.Sprint(v)
		for _, failState := range failStates {
			if currentStatus == failState {
				return object, currentStatus, WrapError(Error(FailedToReachTargetStatus, currentStatus))
			}
		}
		return object, currentStatus, nil
	}
}

// DescribeThreatDetectionImageEventOperation >>> Encapsulated.

// DescribeThreatDetectionSasTrail <<< Encapsulated get interface for ThreatDetection SasTrail.

func (s *ThreatDetectionServiceV2) DescribeThreatDetectionSasTrail(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	parts := strings.Split(id, ":")
	if len(parts) != 0 {
		err = WrapError(fmt.Errorf("invalid Resource Id %s. Expected parts' length %d, got %d", id, 0, len(parts)))
	}
	action := "GetServiceTrail"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["RegionId"] = client.RegionId

	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(1*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("Sas", "2018-12-03", action, query, request, true)

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
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$.ServiceTrail", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.ServiceTrail", response)
	}

	currentStatus := v.(map[string]interface{})["Config"]
	if currentStatus == "off" {
		return object, WrapErrorf(Error(GetNotFoundMessage("SasTrail", id)), NotFoundMsg, response)
	}

	return v.(map[string]interface{}), nil
}

func (s *ThreatDetectionServiceV2) ThreatDetectionSasTrailStateRefreshFunc(id string, field string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeThreatDetectionSasTrail(id)
		if err != nil {
			if NotFoundError(err) {
				return object, "", nil
			}
			return nil, "", WrapError(err)
		}

		v, err := jsonpath.Get(field, object)
		currentStatus := fmt.Sprint(v)
		for _, failState := range failStates {
			if currentStatus == failState {
				return object, currentStatus, WrapError(Error(FailedToReachTargetStatus, currentStatus))
			}
		}
		return object, currentStatus, nil
	}
}

// DescribeThreatDetectionSasTrail >>> Encapsulated.

// DescribeThreatDetectionOssScanConfig <<< Encapsulated get interface for ThreatDetection OssScanConfig.

func (s *ThreatDetectionServiceV2) DescribeThreatDetectionOssScanConfig(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	action := "GetOssScanConfig"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	query["Id"] = id

	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(1*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("Sas", "2018-12-03", action, query, request, true)

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
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$.Data", response)
	if err != nil {
		return object, WrapErrorf(Error(GetNotFoundMessage("OssScanConfig", id)), NotFoundMsg, response)
	}

	currentStatus := v.(map[string]interface{})["Enable"]
	if currentStatus == nil {
		return object, WrapErrorf(Error(GetNotFoundMessage("OssScanConfig", id)), NotFoundMsg, response)
	}

	return v.(map[string]interface{}), nil
}

func (s *ThreatDetectionServiceV2) ThreatDetectionOssScanConfigStateRefreshFunc(id string, field string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeThreatDetectionOssScanConfig(id)
		if err != nil {
			if NotFoundError(err) {
				return object, "", nil
			}
			return nil, "", WrapError(err)
		}

		v, err := jsonpath.Get(field, object)
		currentStatus := fmt.Sprint(v)

		for _, failState := range failStates {
			if currentStatus == failState {
				return object, currentStatus, WrapError(Error(FailedToReachTargetStatus, currentStatus))
			}
		}
		return object, currentStatus, nil
	}
}

// DescribeThreatDetectionOssScanConfig >>> Encapsulated.

// DescribeThreatDetectionAntiBruteForceRule <<< Encapsulated get interface for ThreatDetection AntiBruteForceRule.

func (s *ThreatDetectionServiceV2) DescribeThreatDetectionAntiBruteForceRule(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["Id"] = id

	action := "DescribeAntiBruteForceRules"

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(1*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("Sas", "2018-12-03", action, query, request, true)

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

	v, err := jsonpath.Get("$.Rules[*]", response)
	if err != nil {
		return object, WrapErrorf(Error(GetNotFoundMessage("AntiBruteForceRule", id)), NotFoundMsg, response)
	}

	if len(v.([]interface{})) == 0 {
		return object, WrapErrorf(Error(GetNotFoundMessage("AntiBruteForceRule", id)), NotFoundMsg, response)
	}

	currentStatus := v.([]interface{})[0].(map[string]interface{})["Id"]
	if currentStatus == nil {
		return object, WrapErrorf(Error(GetNotFoundMessage("AntiBruteForceRule", id)), NotFoundMsg, response)
	}

	return v.([]interface{})[0].(map[string]interface{}), nil
}

func (s *ThreatDetectionServiceV2) ThreatDetectionAntiBruteForceRuleStateRefreshFunc(id string, field string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeThreatDetectionAntiBruteForceRule(id)
		if err != nil {
			if NotFoundError(err) {
				return object, "", nil
			}
			return nil, "", WrapError(err)
		}

		v, err := jsonpath.Get(field, object)
		currentStatus := fmt.Sprint(v)

		if strings.HasPrefix(field, "#") {
			v, _ := jsonpath.Get(strings.TrimPrefix(field, "#"), object)
			if v != nil {
				currentStatus = "#CHECKSET"
			}
		}

		for _, failState := range failStates {
			if currentStatus == failState {
				return object, currentStatus, WrapError(Error(FailedToReachTargetStatus, currentStatus))
			}
		}
		return object, currentStatus, nil
	}
}

// DescribeThreatDetectionAntiBruteForceRule >>> Encapsulated.

// DescribeThreatDetectionAssetSelectionConfig <<< Encapsulated get interface for ThreatDetection AssetSelectionConfig.

func (s *ThreatDetectionServiceV2) DescribeThreatDetectionAssetSelectionConfig(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["BusinessType"] = id

	action := "GetAssetSelectionConfig"

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(1*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("Sas", "2018-12-03", action, query, request, true)

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

	v, err := jsonpath.Get("$.Data", response)
	if err != nil {
		return object, WrapErrorf(Error(GetNotFoundMessage("AssetSelectionConfig", id)), NotFoundMsg, response)
	}

	currentStatus := v.(map[string]interface{})["TargetType"]
	if currentStatus == nil {
		return object, WrapErrorf(Error(GetNotFoundMessage("AssetSelectionConfig", id)), NotFoundMsg, response)
	}

	return v.(map[string]interface{}), nil
}

func (s *ThreatDetectionServiceV2) ThreatDetectionAssetSelectionConfigStateRefreshFunc(id string, field string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeThreatDetectionAssetSelectionConfig(id)
		if err != nil {
			if NotFoundError(err) {
				return object, "", nil
			}
			return nil, "", WrapError(err)
		}

		v, err := jsonpath.Get(field, object)
		currentStatus := fmt.Sprint(v)

		if strings.HasPrefix(field, "#") {
			v, _ := jsonpath.Get(strings.TrimPrefix(field, "#"), object)
			if v != nil {
				currentStatus = "#CHECKSET"
			}
		}

		for _, failState := range failStates {
			if currentStatus == failState {
				return object, currentStatus, WrapError(Error(FailedToReachTargetStatus, currentStatus))
			}
		}
		return object, currentStatus, nil
	}
}

// DescribeThreatDetectionAssetSelectionConfig >>> Encapsulated.

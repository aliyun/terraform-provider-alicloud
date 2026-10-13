// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type SslCertificatesServiceServiceV2 struct {
	client *connectivity.AliyunClient
}

// DescribeSslCertificatesServicePcaCertificate <<< Encapsulated get interface for SslCertificatesService PcaCertificate.

func (s *SslCertificatesServiceServiceV2) DescribeSslCertificatesServicePcaCertificate(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["Identifier"] = id

	action := "DescribeCACertificate"

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("cas", "2020-06-30", action, query, request, true)

		if err != nil {
			if NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		if IsExpectedErrors(err, []string{"NotFound"}) {
			return object, WrapErrorf(NotFoundErr("PcaCertificate", id), NotFoundMsg, response)
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$.Certificate", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Certificate", response)
	}

	return v.(map[string]interface{}), nil
}

func (s *SslCertificatesServiceServiceV2) SslCertificatesServicePcaCertificateStateRefreshFunc(id string, field string, failStates []string) retry.StateRefreshFunc {
	return s.SslCertificatesServicePcaCertificateStateRefreshFuncWithApi(id, field, failStates, s.DescribeSslCertificatesServicePcaCertificate)
}

func (s *SslCertificatesServiceServiceV2) SslCertificatesServicePcaCertificateStateRefreshFuncWithApi(id string, field string, failStates []string, call func(id string) (map[string]interface{}, error)) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := call(id)
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

// DescribeSslCertificatesServicePcaCertificate >>> Encapsulated.

// DescribeSslCertificatesServiceCertificate <<< Encapsulated get interface for SslCertificatesService Certificate.

func (s *SslCertificatesServiceServiceV2) DescribeSslCertificatesServiceCertificate(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["CertId"] = id

	action := "GetUserCertificateDetail"

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("cas", "2020-04-07", action, query, request, true)

		if err != nil {
			if NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		if IsExpectedErrors(err, []string{"NotFound"}) {
			return object, WrapErrorf(NotFoundErr("Certificate", id), NotFoundMsg, response)
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}

	return response, nil
}

func (s *SslCertificatesServiceServiceV2) SslCertificatesServiceCertificateStateRefreshFunc(id string, field string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeSslCertificatesServiceCertificate(id)
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

// DescribeSslCertificatesServiceCertificate >>> Encapsulated.

// SetResourceTags <<< Encapsulated tag function for SslCertificates.
func (s *SslCertificatesServiceServiceV2) SetResourceTags(d *schema.ResourceData, resourceType string) error {
	if d.HasChange("tags") {
		var action string
		var err error
		client := s.client
		var request map[string]interface{}
		var response map[string]interface{}
		query := make(map[string]interface{})

		added, removed := parsingTags(d)
		removedTagKeys := make([]string, 0)
		for _, v := range removed {
			if !ignoredTags(v, "") {
				removedTagKeys = append(removedTagKeys, v)
			}
		}
		if len(removedTagKeys) > 0 {
			action = "UntagResources"
			request = make(map[string]interface{})
			query = make(map[string]interface{})
			request["ResourceId.1"] = d.Id()
			request["RegionId"] = client.RegionId
			for i, key := range removedTagKeys {
				request[fmt.Sprintf("TagKey.%d", i+1)] = key
			}

			request["ResourceType"] = resourceType
			wait := incrementalWait(3*time.Second, 5*time.Second)
			err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
				response, err = client.RpcPost("cas", "2020-06-30", action, query, request, true)
				if err != nil {
					if NeedRetry(err) {
						wait()
						return retry.RetryableError(err)
					}
					return retry.NonRetryableError(err)
				}
				return nil
			})
			addDebug(action, response, request)
			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
			}

		}

		if len(added) > 0 {
			action = "TagResources"
			request = make(map[string]interface{})
			query = make(map[string]interface{})
			request["ResourceId.1"] = d.Id()
			request["RegionId"] = client.RegionId
			count := 1
			for key, value := range added {
				request[fmt.Sprintf("Tag.%d.Key", count)] = key
				request[fmt.Sprintf("Tag.%d.Value", count)] = value
				count++
			}

			request["ResourceType"] = resourceType
			wait := incrementalWait(3*time.Second, 5*time.Second)
			err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
				response, err = client.RpcPost("cas", "2020-06-30", action, query, request, true)
				if err != nil {
					if NeedRetry(err) {
						wait()
						return retry.RetryableError(err)
					}
					return retry.NonRetryableError(err)
				}
				return nil
			})
			addDebug(action, response, request)
			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
			}

		}
	}

	return nil
}

// SetResourceTags >>> tag function encapsulated.

func (s *SslCertificatesServiceServiceV2) DescribeSslCertificatesServicePcaCert(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["Identifier"] = id

	action := "DescribeClientCertificate"

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("cas", "2020-06-30", action, query, request, true)

		if err != nil {
			if NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		if IsExpectedErrors(err, []string{"NotFound"}) {
			return object, WrapErrorf(NotFoundErr("PcaCert", id), NotFoundMsg, response)
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$.Certificate", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Certificate", response)
	}

	return v.(map[string]interface{}), nil
}

func (s *SslCertificatesServiceServiceV2) SslCertificatesServicePcaCertStateRefreshFunc(id string, field string, failStates []string) retry.StateRefreshFunc {
	return s.SslCertificatesServicePcaCertStateRefreshFuncWithApi(id, field, failStates, s.DescribeSslCertificatesServicePcaCert)
}

func (s *SslCertificatesServiceServiceV2) SslCertificatesServicePcaCertStateRefreshFuncWithApi(id string, field string, failStates []string, call func(id string) (map[string]interface{}, error)) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := call(id)
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

// DescribeSslCertificatesServicePcaCert >>> Encapsulated.

// DescribeSslCertificatesServiceCompany <<< Encapsulated get interface for SslCertificatesService Company.

func (s *SslCertificatesServiceServiceV2) DescribeSslCertificatesServiceCompany(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["CompanyId"] = id

	action := "GetCompany"

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("cas", "2020-04-07", action, query, request, true)

		if err != nil {
			if NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		if IsExpectedErrors(err, []string{"NotFound"}) {
			return object, WrapErrorf(NotFoundErr("Company", id), NotFoundMsg, response)
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}

	return response, nil
}

func (s *SslCertificatesServiceServiceV2) SslCertificatesServiceCompanyStateRefreshFunc(id string, field string, failStates []string) retry.StateRefreshFunc {
	return s.SslCertificatesServiceCompanyStateRefreshFuncWithApi(id, field, failStates, s.DescribeSslCertificatesServiceCompany)
}

func (s *SslCertificatesServiceServiceV2) SslCertificatesServiceCompanyStateRefreshFuncWithApi(id string, field string, failStates []string, call func(id string) (map[string]interface{}, error)) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := call(id)
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

// DescribeSslCertificatesServiceCompany >>> Encapsulated.

// DescribeSslCertificatesServiceContact <<< Encapsulated get interface for SslCertificatesService Contact.

func (s *SslCertificatesServiceServiceV2) DescribeSslCertificatesServiceContact(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["ContactId"] = id

	action := "GetContact"

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("cas", "2020-04-07", action, query, request, true)

		if err != nil {
			if NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		if IsExpectedErrors(err, []string{"NotFound"}) {
			return object, WrapErrorf(NotFoundErr("Contact", id), NotFoundMsg, response)
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}

	return response, nil
}

func (s *SslCertificatesServiceServiceV2) SslCertificatesServiceContactStateRefreshFunc(id string, field string, failStates []string) retry.StateRefreshFunc {
	return s.SslCertificatesServiceContactStateRefreshFuncWithApi(id, field, failStates, s.DescribeSslCertificatesServiceContact)
}

func (s *SslCertificatesServiceServiceV2) SslCertificatesServiceContactStateRefreshFuncWithApi(id string, field string, failStates []string, call func(id string) (map[string]interface{}, error)) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := call(id)
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

// DescribeSslCertificatesServiceContact >>> Encapsulated.

// DescribeSslCertificatesServiceInstance <<< Encapsulated get interface for SslCertificatesService Instance.
func (s *SslCertificatesServiceServiceV2) DescribeSslCertificatesServiceInstance(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["InstanceId"] = id

	action := "GetInstanceDetail"

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("cas", "2020-04-07", action, query, request, true)

		if err != nil {
			if NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		if IsExpectedErrors(err, []string{"NotFound"}) {
			return object, WrapErrorf(NotFoundErr("Instance", id), NotFoundMsg, response)
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}

	return response, nil
}

func (s *SslCertificatesServiceServiceV2) SslCertificatesServiceInstanceStateRefreshFunc(id string, field string, failStates []string) retry.StateRefreshFunc {
	return s.SslCertificatesServiceInstanceStateRefreshFuncWithApi(id, field, failStates, s.DescribeSslCertificatesServiceInstance)
}

func (s *SslCertificatesServiceServiceV2) SslCertificatesServiceInstanceStateRefreshFuncWithApi(id string, field string, failStates []string, call func(id string) (map[string]interface{}, error)) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := call(id)
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

// DescribeSslCertificatesServiceInstance >>> Encapsulated.

// DescribeSslCertificatesServiceCertificateApply <<< Encapsulated get interface for SslCertificatesService CertificateApply.
func (s *SslCertificatesServiceServiceV2) DescribeSslCertificatesServiceCertificateApply(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["InstanceId"] = id

	action := "GetInstanceDetail"

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("cas", "2020-04-07", action, query, request, true)

		if err != nil {
			if NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		if IsExpectedErrors(err, []string{"NotFound"}) {
			return object, WrapErrorf(NotFoundErr("CertificateApply", id), NotFoundMsg, response)
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}

	return response, nil
}

func (s *SslCertificatesServiceServiceV2) SslCertificatesServiceCertificateApplyStateRefreshFunc(id string, field string, failStates []string) retry.StateRefreshFunc {
	return s.SslCertificatesServiceCertificateApplyStateRefreshFuncWithApi(id, field, failStates, s.DescribeSslCertificatesServiceCertificateApply)
}

func (s *SslCertificatesServiceServiceV2) SslCertificatesServiceCertificateApplyStateRefreshFuncWithApi(id string, field string, failStates []string, call func(id string) (map[string]interface{}, error)) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := call(id)
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

// DescribeSslCertificatesServiceCertificateApply >>> Encapsulated.

// certificateApplicationIsOutstanding reports whether the instance still has a certificate
// application in flight, which is what makes a withdrawal both possible and necessary.
//
// An instance that cannot be read counts as still having one. Not knowing is not the same as
// knowing there is nothing left to withdraw, and reporting a withdrawal that never happened leaves
// the instance sitting in pending for whatever tries to delete it next — a state in which it can be
// neither refunded nor deleted. An instance that is genuinely gone is the one exception.
func (s *SslCertificatesServiceServiceV2) certificateApplicationIsOutstanding(id string) bool {
	object, err := s.DescribeSslCertificatesServiceInstance(id)
	if err != nil {
		return !NotFoundError(err)
	}
	return fmt.Sprint(object["Status"]) == "pending"
}

// DescribeSslCertificatesServiceInstanceCertificate <<< Encapsulated get interface for SslCertificatesService InstanceCertificate.
func (s *SslCertificatesServiceServiceV2) DescribeSslCertificatesServiceInstanceCertificate(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["CertificateId"] = id

	action := "GetCertificateDetail"

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("cas", "2020-04-07", action, query, request, true)

		if err != nil {
			if NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		if IsExpectedErrors(err, []string{"NotFound"}) {
			return object, WrapErrorf(NotFoundErr("InstanceCertificate", id), NotFoundMsg, response)
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}

	return response, nil
}

func (s *SslCertificatesServiceServiceV2) SslCertificatesServiceInstanceCertificateStateRefreshFunc(id string, field string, failStates []string) retry.StateRefreshFunc {
	return s.SslCertificatesServiceInstanceCertificateStateRefreshFuncWithApi(id, field, failStates, s.DescribeSslCertificatesServiceInstanceCertificate)
}

func (s *SslCertificatesServiceServiceV2) SslCertificatesServiceInstanceCertificateStateRefreshFuncWithApi(id string, field string, failStates []string, call func(id string) (map[string]interface{}, error)) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := call(id)
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

// DescribeSslCertificatesServiceInstanceCertificate >>> Encapsulated.

// DescribeSslCertificatesServiceCertificateValidation <<< Encapsulated get interface for SslCertificatesService CertificateValidation.
func (s *SslCertificatesServiceServiceV2) DescribeSslCertificatesServiceCertificateValidation(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["InstanceId"] = id

	action := "GetInstanceDetail"

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("cas", "2020-04-07", action, query, request, true)

		if err != nil {
			if NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		if IsExpectedErrors(err, []string{"NotFound"}) {
			return object, WrapErrorf(NotFoundErr("CertificateValidation", id), NotFoundMsg, response)
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}

	return response, nil
}

func (s *SslCertificatesServiceServiceV2) SslCertificatesServiceCertificateValidationStateRefreshFunc(id string, field string, failStates []string) retry.StateRefreshFunc {
	return s.SslCertificatesServiceCertificateValidationStateRefreshFuncWithApi(id, field, failStates, s.DescribeSslCertificatesServiceCertificateValidation)
}

func (s *SslCertificatesServiceServiceV2) SslCertificatesServiceCertificateValidationStateRefreshFuncWithApi(id string, field string, failStates []string, call func(id string) (map[string]interface{}, error)) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := call(id)
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

// DescribeSslCertificatesServiceCertificateValidation >>> Encapsulated.

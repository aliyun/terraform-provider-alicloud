// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

type CmsServiceV2 struct {
	client *connectivity.AliyunClient
}

// DescribeCmsWorkspace <<< Encapsulated get interface for Cms Workspace.

func (s *CmsServiceV2) DescribeCmsWorkspace(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	workspaceName := id
	request = make(map[string]interface{})
	query = make(map[string]*string)

	action := fmt.Sprintf("/workspace/%s", workspaceName)

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(1*time.Minute, func() *resource.RetryError {
		response, err = client.RoaGet("Cms", "2024-03-30", action, query, nil, nil)

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
		if IsExpectedErrors(err, []string{"WorkspaceNotExist"}) {
			return object, WrapErrorf(NotFoundErr("Workspace", id), NotFoundMsg, response)
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}

	return response, nil
}

func (s *CmsServiceV2) CmsWorkspaceStateRefreshFunc(id string, field string, failStates []string) resource.StateRefreshFunc {
	return s.CmsWorkspaceStateRefreshFuncWithApi(id, field, failStates, s.DescribeCmsWorkspace)
}

func (s *CmsServiceV2) CmsWorkspaceStateRefreshFuncWithApi(id string, field string, failStates []string, call func(id string) (map[string]interface{}, error)) resource.StateRefreshFunc {
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

// DescribeCmsWorkspace >>> Encapsulated.

// DescribeCmsIntegrationPolicy <<< Encapsulated get interface for Cms IntegrationPolicy.

func (s *CmsServiceV2) DescribeCmsIntegrationPolicy(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	policyId := id
	request = make(map[string]interface{})
	query = make(map[string]*string)

	action := fmt.Sprintf("/integration-policies/%s", policyId)

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(1*time.Minute, func() *resource.RetryError {
		response, err = client.RoaGet("Cms", "2024-03-30", action, query, nil, nil)

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
		if IsExpectedErrors(err, []string{"404", "15007"}) {
			return object, WrapErrorf(NotFoundErr("IntegrationPolicy", id), NotFoundMsg, response)
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}

	return response, nil
}

func (s *CmsServiceV2) CmsIntegrationPolicyStateRefreshFunc(id string, field string, failStates []string) resource.StateRefreshFunc {
	return s.CmsIntegrationPolicyStateRefreshFuncWithApi(id, field, failStates, s.DescribeCmsIntegrationPolicy)
}

func (s *CmsServiceV2) CmsIntegrationPolicyStateRefreshFuncWithApi(id string, field string, failStates []string, call func(id string) (map[string]interface{}, error)) resource.StateRefreshFunc {
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

// DescribeCmsIntegrationPolicy >>> Encapsulated.

// DescribeCmsPrometheusInstance <<< Encapsulated get interface for Cms PrometheusInstance.

func (s *CmsServiceV2) DescribeCmsPrometheusInstance(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	prometheusInstanceId := id
	request = make(map[string]interface{})
	query = make(map[string]*string)

	action := fmt.Sprintf("/prometheus-instances/%s", prometheusInstanceId)

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(1*time.Minute, func() *resource.RetryError {
		response, err = client.RoaGet("Cms", "2024-03-30", action, query, nil, nil)

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
		if IsExpectedErrors(err, []string{"404"}) {
			return object, WrapErrorf(NotFoundErr("PrometheusInstance", id), NotFoundMsg, response)
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$.prometheusInstance", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.prometheusInstance", response)
	}

	return v.(map[string]interface{}), nil
}

func (s *CmsServiceV2) CmsPrometheusInstanceStateRefreshFunc(id string, field string, failStates []string) resource.StateRefreshFunc {
	return s.CmsPrometheusInstanceStateRefreshFuncWithApi(id, field, failStates, s.DescribeCmsPrometheusInstance)
}

func (s *CmsServiceV2) CmsPrometheusInstanceStateRefreshFuncWithApi(id string, field string, failStates []string, call func(id string) (map[string]interface{}, error)) resource.StateRefreshFunc {
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

// DescribeCmsPrometheusInstance >>> Encapsulated.

// DescribeCmsPrometheusView <<< Encapsulated get interface for Cms PrometheusView.

func (s *CmsServiceV2) DescribeCmsPrometheusView(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	prometheusViewId := id
	request = make(map[string]interface{})
	query = make(map[string]*string)

	action := fmt.Sprintf("/prometheus-views/%s", prometheusViewId)

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(1*time.Minute, func() *resource.RetryError {
		response, err = client.RoaGet("Cms", "2024-03-30", action, query, nil, nil)

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
		if IsExpectedErrors(err, []string{"404"}) {
			return object, WrapErrorf(NotFoundErr("PrometheusView", id), NotFoundMsg, response)
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$.prometheusView", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.prometheusView", response)
	}

	return v.(map[string]interface{}), nil
}

func (s *CmsServiceV2) CmsPrometheusViewStateRefreshFunc(id string, field string, failStates []string) resource.StateRefreshFunc {
	return s.CmsPrometheusViewStateRefreshFuncWithApi(id, field, failStates, s.DescribeCmsPrometheusView)
}

func (s *CmsServiceV2) CmsPrometheusViewStateRefreshFuncWithApi(id string, field string, failStates []string, call func(id string) (map[string]interface{}, error)) resource.StateRefreshFunc {
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

// DescribeCmsPrometheusView >>> Encapsulated.

// DescribeCmsAddonRelease <<< Encapsulated get interface for Cms AddonRelease.

func (s *CmsServiceV2) DescribeCmsAddonRelease(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	parts := strings.Split(id, ":")
	if len(parts) != 2 {
		err = WrapError(fmt.Errorf("invalid Resource Id %s. Expected parts' length %d, got %d", id, 2, len(parts)))
		return nil, err
	}
	releaseName := parts[1]
	policyId := parts[0]
	request = make(map[string]interface{})
	query = make(map[string]*string)

	action := fmt.Sprintf("/integration-policies/%s/addon-releases/%s", policyId, releaseName)

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(1*time.Minute, func() *resource.RetryError {
		response, err = client.RoaGet("Cms", "2024-03-30", action, query, nil, nil)

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
		if IsExpectedErrors(err, []string{"404", "12002"}) {
			return object, WrapErrorf(NotFoundErr("AddonRelease", id), NotFoundMsg, response)
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$.release", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.release", response)
	}

	return v.(map[string]interface{}), nil
}

func (s *CmsServiceV2) CmsAddonReleaseStateRefreshFunc(id string, field string, failStates []string) resource.StateRefreshFunc {
	return s.CmsAddonReleaseStateRefreshFuncWithApi(id, field, failStates, s.DescribeCmsAddonRelease)
}

func (s *CmsServiceV2) CmsAddonReleaseStateRefreshFuncWithApi(id string, field string, failStates []string, call func(id string) (map[string]interface{}, error)) resource.StateRefreshFunc {
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

// DescribeCmsAddonRelease >>> Encapsulated.

// DescribeCmsAggTaskGroup <<< Encapsulated get interface for Cms AggTaskGroup.

func (s *CmsServiceV2) DescribeCmsAggTaskGroup(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	parts := strings.Split(id, ":")
	if len(parts) != 2 {
		err = WrapError(fmt.Errorf("invalid Resource Id %s. Expected parts' length %d, got %d", id, 2, len(parts)))
		return nil, err
	}
	groupId := parts[1]
	instanceId := parts[0]
	request = make(map[string]interface{})
	query = make(map[string]*string)

	action := fmt.Sprintf("/prometheus-instances/%s/agg-task-groups/%s", instanceId, groupId)

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(1*time.Minute, func() *resource.RetryError {
		response, err = client.RoaGet("Cms", "2024-03-30", action, query, nil, nil)

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
		if IsExpectedErrors(err, []string{"404"}) {
			return object, WrapErrorf(NotFoundErr("AggTaskGroup", id), NotFoundMsg, response)
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$.aggTaskGroup", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.aggTaskGroup", response)
	}

	return v.(map[string]interface{}), nil
}

func (s *CmsServiceV2) CmsAggTaskGroupStateRefreshFunc(id string, field string, failStates []string) resource.StateRefreshFunc {
	return s.CmsAggTaskGroupStateRefreshFuncWithApi(id, field, failStates, s.DescribeCmsAggTaskGroup)
}

func (s *CmsServiceV2) CmsAggTaskGroupStateRefreshFuncWithApi(id string, field string, failStates []string, call func(id string) (map[string]interface{}, error)) resource.StateRefreshFunc {
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

// DescribeCmsAggTaskGroup >>> Encapsulated.

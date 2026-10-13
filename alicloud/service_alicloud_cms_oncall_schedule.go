// Package alicloud. This file is hand-written from the Cms 2024-03-30 OpenAPI definition.
package alicloud

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// DescribeCmsOncallSchedule <<< Encapsulated get interface for Cms OncallSchedule.
func (s *CmsServiceV2) DescribeCmsOncallSchedule(id string) (object map[string]interface{}, err error) {
	client := s.client
	var response map[string]interface{}
	var query map[string]*string
	query = make(map[string]*string)

	action := fmt.Sprintf("/oncallSchedule/%s", id)

	wait := incrementalWait(3*time.Second, 3*time.Second)
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
	addDebug(action, response, query)
	if err != nil {
		if IsExpectedErrors(err, []string{"ResourceNotFound"}) || NotFoundError(err) {
			return object, WrapErrorf(NotFoundErr("Cms:OncallSchedule", id), NotFoundMsg, response)
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}

	return response, nil
}

// CmsOncallScheduleStateRefreshFunc ...
func (s *CmsServiceV2) CmsOncallScheduleStateRefreshFunc(id string, field string, failStates []string) resource.StateRefreshFunc {
	return s.CmsOncallScheduleStateRefreshFuncWithApi(id, field, failStates, s.DescribeCmsOncallSchedule)
}

// CmsOncallScheduleStateRefreshFuncWithApi ...
func (s *CmsServiceV2) CmsOncallScheduleStateRefreshFuncWithApi(id string, field string, failStates []string, call func(id string) (map[string]interface{}, error)) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := call(id)
		if err != nil {
			if NotFoundError(err) {
				return object, "", nil
			}
			return nil, "", WrapError(err)
		}
		return object, "", nil
	}
}

// DescribeCmsOncallSchedule >>> Encapsulated.

// flattenOncallScheduleRotations converts the API response rotations list (camelCase)
// into the schema-compatible list (snake_case) for state storage.
func flattenOncallScheduleRotations(raw interface{}) []map[string]interface{} {
	result := make([]map[string]interface{}, 0)
	if raw == nil {
		return result
	}
	list, ok := raw.([]interface{})
	if !ok {
		return result
	}
	for _, item := range list {
		m, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		rotation := make(map[string]interface{})
		rotation["active_days"] = flattenStringList(m["activeDays"])
		rotation["contacts"] = flattenStringList(m["contacts"])
		rotation["rotation_end_time"] = m["rotationEndTime"]
		rotation["rotation_name"] = m["rotationName"]
		rotation["rotation_start_time"] = m["rotationStartTime"]
		rotation["shift_length"] = m["shiftLength"]
		rotation["shift_recurrence_frequency"] = m["shiftRecurrenceFrequency"]
		rotation["start_date"] = m["startDate"]
		rotation["time_zone"] = m["timeZone"]
		result = append(result, rotation)
	}
	return result
}

// expandOncallScheduleRotations converts the schema rotations list (snake_case)
// into the API request body list (camelCase).
func expandOncallScheduleRotations(raw interface{}) []map[string]interface{} {
	result := make([]map[string]interface{}, 0)
	if raw == nil {
		return result
	}
	list, ok := raw.([]interface{})
	if !ok {
		return result
	}
	for _, item := range list {
		m, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		rotation := make(map[string]interface{})
		rotation["activeDays"] = flattenStringList(m["active_days"])
		rotation["contacts"] = flattenStringList(m["contacts"])
		if v, ok := m["rotation_end_time"]; ok && v != nil {
			rotation["rotationEndTime"] = v
		}
		if v, ok := m["rotation_name"]; ok && v != nil {
			rotation["rotationName"] = v
		}
		if v, ok := m["rotation_start_time"]; ok && v != nil {
			rotation["rotationStartTime"] = v
		}
		if v, ok := m["shift_length"]; ok && v != nil {
			rotation["shiftLength"] = v
		}
		if v, ok := m["shift_recurrence_frequency"]; ok && v != nil {
			rotation["shiftRecurrenceFrequency"] = v
		}
		if v, ok := m["start_date"]; ok && v != nil {
			rotation["startDate"] = v
		}
		if v, ok := m["time_zone"]; ok && v != nil {
			rotation["timeZone"] = v
		}
		result = append(result, rotation)
	}
	return result
}

// expandOncallScheduleSubstitudes passes the schema substitudes list (strings)
// straight through to the API body. The API field is named "substitudes".
func expandOncallScheduleSubstitudes(raw interface{}) []interface{} {
	return flattenStringList(raw)
}

// flattenOncallScheduleSubstitudes returns the API substitudes list as-is.
func flattenOncallScheduleSubstitudes(raw interface{}) []interface{} {
	return flattenStringList(raw)
}

func flattenStringList(raw interface{}) []interface{} {
	result := make([]interface{}, 0)
	if raw == nil {
		return result
	}
	list, ok := raw.([]interface{})
	if !ok {
		return result
	}
	for _, v := range list {
		result = append(result, v)
	}
	return result
}

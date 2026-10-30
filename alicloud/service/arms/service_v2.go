package arms

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/tags"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type ServiceV2 struct {
	client *connectivity.AliyunClient
}

func NewServiceV2(client *connectivity.AliyunClient) *ServiceV2 {
	return &ServiceV2{client: client}
}

// DescribeArmsPrometheusMonitoring <<< Encapsulated get interface for Arms PrometheusMonitoring.

func (s *ServiceV2) DescribeArmsPrometheusMonitoring(id string) (object map[string]interface{}, err error) {

	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	parts := strings.Split(id, ":")
	if len(parts) != 3 {
		err = errs.WrapError(fmt.Errorf("invalid Resource Id %s. Expected parts' length %d, got %d", id, 3, len(parts)))
	}
	action := "GetPrometheusMonitoring"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	query["ClusterId"] = parts[0]
	query["MonitoringName"] = parts[1]
	query["Type"] = parts[2]
	request["RegionId"] = client.RegionId

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("ARMS", "2019-08-08", action, query, request, false)

		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		logs.AddDebug(action, response, request)
		return nil
	})

	code, _ := jsonpath.Get("$.Code", response)
	if helper.InArray(fmt.Sprint(code), []string{"404"}) {
		return object, errs.WrapErrorf(errs.NotFoundErr("PrometheusMonitoring", id), errs.NotFoundMsg, errs.ProviderERROR, fmt.Sprint(response["RequestId"]))
	}

	v, err := jsonpath.Get("$.Data", response)
	if err != nil {
		return object, errs.WrapErrorf(errs.NotFoundErr("PrometheusMonitoring", id), errs.NotFoundMsg, errs.ProviderERROR, fmt.Sprint(response["RequestId"]))
	}

	return v.(map[string]interface{}), nil
}

func (s *ServiceV2) ArmsPrometheusMonitoringStateRefreshFunc(id string, field string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeArmsPrometheusMonitoring(id)
		if err != nil {
			if errs.NotFoundError(err) {
				return object, "", nil
			}
			return nil, "", errs.WrapError(err)
		}

		currentStatus := fmt.Sprint(object[field])
		for _, failState := range failStates {
			if currentStatus == failState {
				return object, currentStatus, errs.WrapError(errs.Error(errs.FailedToReachTargetStatus, currentStatus))
			}
		}
		return object, currentStatus, nil
	}
}

// DescribeArmsPrometheusMonitoring >>> Encapsulated.
// DescribeArmsRemoteWrite <<< Encapsulated get interface for Arms RemoteWrite.

func (s *ServiceV2) DescribeArmsRemoteWrite(id string) (object map[string]interface{}, err error) {

	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	parts := strings.Split(id, ":")
	if len(parts) != 2 {
		err = errs.WrapError(fmt.Errorf("invalid Resource Id %s. Expected parts' length %d, got %d", id, 2, len(parts)))
	}
	action := "GetPrometheusRemoteWrite"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	query["ClusterId"] = parts[0]
	query["RemoteWriteName"] = parts[1]
	request["RegionId"] = client.RegionId

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("ARMS", "2019-08-08", action, query, request, false)

		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		logs.AddDebug(action, response, request)
		return nil
	})

	code, _ := jsonpath.Get("$.Code", response)
	if helper.InArray(fmt.Sprint(code), []string{"404"}) {
		return object, errs.WrapErrorf(errs.NotFoundErr("RemoteWrite", id), errs.NotFoundMsg, errs.ProviderERROR, fmt.Sprint(response["RequestId"]))
	}

	v, err := jsonpath.Get("$.Data", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.Data", response)
	}

	return v.(map[string]interface{}), nil
}

func (s *ServiceV2) ArmsRemoteWriteStateRefreshFunc(id string, field string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeArmsRemoteWrite(id)
		if err != nil {
			if errs.NotFoundError(err) {
				return object, "", nil
			}
			return nil, "", errs.WrapError(err)
		}

		currentStatus := fmt.Sprint(object[field])
		for _, failState := range failStates {
			if currentStatus == failState {
				return object, currentStatus, errs.WrapError(errs.Error(errs.FailedToReachTargetStatus, currentStatus))
			}
		}
		return object, currentStatus, nil
	}
}

// DescribeArmsRemoteWrite >>> Encapsulated.

// DescribeArmsEnvironment <<< Encapsulated get interface for Arms Environment.

func (s *ServiceV2) DescribeArmsEnvironment(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	action := "DescribeEnvironment"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	query["EnvironmentId"] = id
	query["RegionId"] = client.RegionId

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("ARMS", "2019-08-08", action, query, request, true)

		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		logs.AddDebug(action, response, request)
		return nil
	})

	code, _ := jsonpath.Get("$.Code", response)
	if helper.InArray(fmt.Sprint(code), []string{"404"}) {
		return object, errs.WrapErrorf(errs.NotFoundErr("Environment", id), errs.NotFoundMsg, response)
	}

	v, err := jsonpath.Get("$.Data", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.Data", response)
	}

	return v.(map[string]interface{}), nil
}

func (s *ServiceV2) ArmsEnvironmentStateRefreshFunc(id string, field string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeArmsEnvironment(id)
		if err != nil {
			if errs.NotFoundError(err) {
				return object, "", nil
			}
			return nil, "", errs.WrapError(err)
		}

		v, err := jsonpath.Get(field, object)
		currentStatus := fmt.Sprint(v)

		for _, failState := range failStates {
			if currentStatus == failState {
				return object, currentStatus, errs.WrapError(errs.Error(errs.FailedToReachTargetStatus, currentStatus))
			}
		}
		return object, currentStatus, nil
	}
}

// DescribeArmsEnvironment >>> Encapsulated.

// SetResourceTags <<< Encapsulated tag function for Arms.
func (s *ServiceV2) SetResourceTags(d *schema.ResourceData, resourceType string) error {
	if d.HasChange("tags") {
		var action string
		var err error
		client := s.client
		var request map[string]interface{}
		var response map[string]interface{}
		query := make(map[string]interface{})

		added, removed := tags.ParsingTags(d)
		removedTagKeys := make([]string, 0)
		for _, v := range removed {
			if !tags.IgnoredTags(v, "") {
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
			for i, key := range removedTagKeys {
				request[fmt.Sprintf("TagKey.%d", i+1)] = key
			}

			wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
			err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
				response, err = client.RpcPost("ARMS", "2019-08-08", action, query, request, true)
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
				return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
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
			wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
			err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
				response, err = client.RpcPost("ARMS", "2019-08-08", action, query, request, true)
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
				return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
			}

		}
	}

	return nil
}

// SetResourceTags >>> tag function encapsulated.
// DescribeArmsPrometheus <<< Encapsulated get interface for Arms Prometheus.

func (s *ServiceV2) DescribeArmsPrometheus(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	action := "GetPrometheusInstance"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	query["ClusterId"] = id
	request["RegionId"] = client.RegionId

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("ARMS", "2019-08-08", action, query, request, true)

		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		logs.AddDebug(action, response, request)
		return nil
	})

	code, _ := jsonpath.Get("$.Code", response)
	if helper.InArray(fmt.Sprint(code), []string{"404"}) {
		return object, errs.WrapErrorf(errs.NotFoundErr("Prometheus", id), errs.NotFoundMsg, response)
	}

	v, err := jsonpath.Get("$.Data", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.Data", response)
	}

	return v.(map[string]interface{}), nil
}

func (s *ServiceV2) ArmsPrometheusStateRefreshFunc(id string, field string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeArmsPrometheus(id)
		if err != nil {
			if errs.NotFoundError(err) {
				return object, "", nil
			}
			return nil, "", errs.WrapError(err)
		}

		v, err := jsonpath.Get(field, object)
		currentStatus := fmt.Sprint(v)
		for _, failState := range failStates {
			if currentStatus == failState {
				return object, currentStatus, errs.WrapError(errs.Error(errs.FailedToReachTargetStatus, currentStatus))
			}
		}
		return object, currentStatus, nil
	}
}

// DescribeArmsPrometheus >>> Encapsulated.

// DescribeArmsEnvFeature <<< Encapsulated get interface for Arms EnvFeature.

func (s *ServiceV2) DescribeArmsEnvFeature(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	parts := strings.Split(id, ":")
	if len(parts) != 2 {
		err = errs.WrapError(fmt.Errorf("invalid Resource Id %s. Expected parts' length %d, got %d", id, 2, len(parts)))
	}
	action := "DescribeEnvironmentFeature"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	query["FeatureName"] = parts[1]
	query["EnvironmentId"] = parts[0]
	request["RegionId"] = client.RegionId

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("ARMS", "2019-08-08", action, query, request, true)

		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		logs.AddDebug(action, response, request)
		return nil
	})

	code, _ := jsonpath.Get("$.Code", response)
	if helper.InArray(fmt.Sprint(code), []string{"404", "14000"}) {
		return object, errs.WrapErrorf(errs.NotFoundErr("EnvFeature", id), errs.NotFoundMsg, response)
	}

	v, err := jsonpath.Get("$.Data", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.Data", response)
	}

	return v.(map[string]interface{}), nil
}

func (s *ServiceV2) ArmsEnvFeatureStateRefreshFunc(id string, field string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeArmsEnvFeature(id)
		if err != nil {
			if errs.NotFoundError(err) {
				return object, "", nil
			}
			return nil, "", errs.WrapError(err)
		}

		v, err := jsonpath.Get(field, object)
		currentStatus := fmt.Sprint(v)
		for _, failState := range failStates {
			if currentStatus == failState {
				return object, currentStatus, errs.WrapError(errs.Error(errs.FailedToReachTargetStatus, currentStatus))
			}
		}
		return object, currentStatus, nil
	}
}

// DescribeArmsEnvFeature >>> Encapsulated.

// DescribeArmsAddonRelease <<< Encapsulated get interface for Arms AddonRelease.

func (s *ServiceV2) DescribeArmsAddonRelease(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	parts := strings.Split(id, ":")
	if len(parts) != 2 {
		err = errs.WrapError(fmt.Errorf("invalid Resource Id %s. Expected parts' length %d, got %d", id, 2, len(parts)))
	}
	action := "DescribeAddonRelease"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	query["ReleaseName"] = parts[1]
	query["EnvironmentId"] = parts[0]
	request["RegionId"] = client.RegionId

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("ARMS", "2019-08-08", action, query, request, true)

		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		logs.AddDebug(action, response, request)
		return nil
	})

	code, _ := jsonpath.Get("$.Code", response)
	if helper.InArray(fmt.Sprint(code), []string{"404"}) {
		return object, errs.WrapErrorf(errs.NotFoundErr("AddonRelease", id), errs.NotFoundMsg, response)
	}

	v, err := jsonpath.Get("$.Data", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.Data", response)
	}

	return v.(map[string]interface{}), nil
}

func (s *ServiceV2) ArmsAddonReleaseStateRefreshFunc(id string, field string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeArmsAddonRelease(id)
		if err != nil {
			if errs.NotFoundError(err) {
				return nil, "", nil
			}
			return nil, "", errs.WrapError(err)
		}

		v, err := jsonpath.Get(field, object)
		currentStatus := fmt.Sprint(v)
		for _, failState := range failStates {
			if currentStatus == failState {
				return object, currentStatus, errs.WrapError(errs.Error(errs.FailedToReachTargetStatus, currentStatus))
			}
		}
		return object, currentStatus, nil
	}
}

// DescribeArmsAddonRelease >>> Encapsulated.

// DescribeArmsEnvPodMonitor <<< Encapsulated get interface for Arms EnvPodMonitor.

func (s *ServiceV2) DescribeArmsEnvPodMonitor(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	parts := strings.Split(id, ":")
	if len(parts) != 3 {
		err = errs.WrapError(fmt.Errorf("invalid Resource Id %s. Expected parts' length %d, got %d", id, 3, len(parts)))
	}
	action := "DescribeEnvPodMonitor"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	query["PodMonitorName"] = parts[2]
	query["EnvironmentId"] = parts[0]
	query["Namespace"] = parts[1]
	request["RegionId"] = client.RegionId

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("ARMS", "2019-08-08", action, query, request, true)

		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		logs.AddDebug(action, response, request)
		return nil
	})

	code, _ := jsonpath.Get("$.Code", response)
	if helper.InArray(fmt.Sprint(code), []string{"404"}) {
		return object, errs.WrapErrorf(errs.NotFoundErr("EnvPodMonitor", id), errs.NotFoundMsg, response)
	}

	v, err := jsonpath.Get("$.Data", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.Data", response)
	}

	return v.(map[string]interface{}), nil
}

func (s *ServiceV2) ArmsEnvPodMonitorStateRefreshFunc(id string, field string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeArmsEnvPodMonitor(id)
		if err != nil {
			if errs.NotFoundError(err) {
				return object, "", nil
			}
			return nil, "", errs.WrapError(err)
		}

		v, err := jsonpath.Get(field, object)
		currentStatus := fmt.Sprint(v)
		for _, failState := range failStates {
			if currentStatus == failState {
				return object, currentStatus, errs.WrapError(errs.Error(errs.FailedToReachTargetStatus, currentStatus))
			}
		}
		return object, currentStatus, nil
	}
}

// DescribeArmsEnvPodMonitor >>> Encapsulated.

// DescribeArmsEnvServiceMonitor <<< Encapsulated get interface for Arms EnvServiceMonitor.

func (s *ServiceV2) DescribeArmsEnvServiceMonitor(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	parts := strings.Split(id, ":")
	if len(parts) != 3 {
		err = errs.WrapError(fmt.Errorf("invalid Resource Id %s. Expected parts' length %d, got %d", id, 3, len(parts)))
	}
	action := "DescribeEnvServiceMonitor"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	query["ServiceMonitorName"] = parts[2]
	query["EnvironmentId"] = parts[0]
	query["Namespace"] = parts[1]
	request["RegionId"] = client.RegionId

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("ARMS", "2019-08-08", action, query, request, true)

		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		logs.AddDebug(action, response, request)
		return nil
	})

	code, _ := jsonpath.Get("$.Code", response)
	if helper.InArray(fmt.Sprint(code), []string{"404"}) {
		return object, errs.WrapErrorf(errs.NotFoundErr("EnvServiceMonitor", id), errs.NotFoundMsg, response)
	}

	v, err := jsonpath.Get("$.Data", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.Data", response)
	}

	return v.(map[string]interface{}), nil
}

func (s *ServiceV2) ArmsEnvServiceMonitorStateRefreshFunc(id string, field string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeArmsEnvServiceMonitor(id)
		if err != nil {
			if errs.NotFoundError(err) {
				return object, "", nil
			}
			return nil, "", errs.WrapError(err)
		}

		v, err := jsonpath.Get(field, object)
		currentStatus := fmt.Sprint(v)
		for _, failState := range failStates {
			if currentStatus == failState {
				return object, currentStatus, errs.WrapError(errs.Error(errs.FailedToReachTargetStatus, currentStatus))
			}
		}
		return object, currentStatus, nil
	}
}

// DescribeArmsEnvServiceMonitor >>> Encapsulated.

// DescribeArmsEnvCustomJob <<< Encapsulated get interface for Arms EnvCustomJob.

func (s *ServiceV2) DescribeArmsEnvCustomJob(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	parts := strings.Split(id, ":")
	if len(parts) != 2 {
		err = errs.WrapError(fmt.Errorf("invalid Resource Id %s. Expected parts' length %d, got %d", id, 2, len(parts)))
	}
	action := "DescribeEnvCustomJob"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	query["CustomJobName"] = parts[1]
	query["EnvironmentId"] = parts[0]
	request["RegionId"] = client.RegionId

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("ARMS", "2019-08-08", action, query, request, true)

		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		logs.AddDebug(action, response, request)
		return nil
	})

	code, _ := jsonpath.Get("$.Code", response)
	if helper.InArray(fmt.Sprint(code), []string{"404"}) {
		return object, errs.WrapErrorf(errs.NotFoundErr("EnvCustomJob", id), errs.NotFoundMsg, response)
	}

	v, err := jsonpath.Get("$.Data", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.Data", response)
	}

	return v.(map[string]interface{}), nil
}

func (s *ServiceV2) ArmsEnvCustomJobStateRefreshFunc(id string, field string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeArmsEnvCustomJob(id)
		if err != nil {
			if errs.NotFoundError(err) {
				return object, "", nil
			}
			return nil, "", errs.WrapError(err)
		}

		v, err := jsonpath.Get(field, object)
		currentStatus := fmt.Sprint(v)
		for _, failState := range failStates {
			if currentStatus == failState {
				return object, currentStatus, errs.WrapError(errs.Error(errs.FailedToReachTargetStatus, currentStatus))
			}
		}
		return object, currentStatus, nil
	}
}

// DescribeArmsEnvCustomJob >>> Encapsulated.

// DescribeArmsSyntheticTask <<< Encapsulated get interface for Arms SyntheticTask.

func (s *ServiceV2) DescribeArmsSyntheticTask(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	action := "GetTimingSyntheticTask"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	query["TaskId"] = id
	query["RegionId"] = client.RegionId

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.RpcGet("ARMS", "2019-08-08", action, query, request)

		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		logs.AddDebug(action, response, request)
		return nil
	})

	if err != nil {
		if errs.IsExpectedErrors(err, []string{"ParameterIllegal"}) {
			return object, errs.WrapErrorf(errs.NotFoundErr("SyntheticTask", id), errs.NotFoundMsg, response)
		}
		logs.AddDebug(action, response, request)
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$.Data", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.Data", response)
	}

	return v.(map[string]interface{}), nil
}

func (s *ServiceV2) ArmsSyntheticTaskStateRefreshFunc(id string, field string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeArmsSyntheticTask(id)
		if err != nil {
			if errs.NotFoundError(err) {
				return object, "", nil
			}
			return nil, "", errs.WrapError(err)
		}

		v, err := jsonpath.Get(field, object)
		currentStatus := fmt.Sprint(v)

		for _, failState := range failStates {
			if currentStatus == failState {
				return object, currentStatus, errs.WrapError(errs.Error(errs.FailedToReachTargetStatus, currentStatus))
			}
		}
		return object, currentStatus, nil
	}
}

// DescribeArmsSyntheticTask >>> Encapsulated.

// DescribeArmsGrafanaWorkspace <<< Encapsulated get interface for Arms GrafanaWorkspace.

func (s *ServiceV2) DescribeArmsGrafanaWorkspace(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["GrafanaWorkspaceId"] = id
	request["RegionId"] = client.RegionId
	action := "GetGrafanaWorkspace"

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("ARMS", "2019-08-08", action, query, request, true)

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

	v, err := jsonpath.Get("$.Data", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.Data", response)
	}

	currentStatus := v.(map[string]interface{})["status"]
	if fmt.Sprint(currentStatus) == "DeleteSucceed" {
		return object, errs.WrapErrorf(errs.NotFoundErr("GrafanaWorkspace", id), errs.NotFoundMsg, response)
	}

	return v.(map[string]interface{}), nil
}

func (s *ServiceV2) ArmsGrafanaWorkspaceStateRefreshFunc(id string, field string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeArmsGrafanaWorkspace(id)
		if err != nil {
			if errs.NotFoundError(err) {
				return object, "", nil
			}
			return nil, "", errs.WrapError(err)
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
				return object, currentStatus, errs.WrapError(errs.Error(errs.FailedToReachTargetStatus, currentStatus))
			}
		}
		return object, currentStatus, nil
	}
}

// DescribeArmsGrafanaWorkspace >>> Encapsulated.

// DescribeArmsNotificationPolicy <<< Encapsulated get interface for Arms NotificationPolicy.
func (s *ServiceV2) DescribeArmsNotificationPolicy(id string) (object map[string]interface{}, err error) {
	return s.DescribeArmsNotificationPolicyContext(context.Background(), id)
}

func (s *ServiceV2) DescribeArmsNotificationPolicyContext(ctx context.Context, id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	query["Ids"] = id
	query["Page"] = 1
	query["Size"] = helper.PageSizeXLarge
	query["IsDetail"] = true
	query["RegionId"] = client.RegionId

	action := "ListNotificationPolicies"

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.RetryContext(ctx, 1*time.Minute, func() *retry.RetryError {
		response, err = client.RpcGet("ARMS", "2019-08-08", action, query, request)

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
		// A transport-level 404 is definitive. A successful RPC that carries
		// Code 404 only in its body is classified after the collection
		// container is validated below, so a malformed body cannot be turned
		// into NotFound.
		if errs.IsTypedNotFound(err) && !errs.IsResponseCodeNotFound(err) {
			return object, errs.WrapErrorf(errs.NotFoundErr("NotificationPolicy", id), errs.NotFoundMsg, response)
		}
		if !errs.IsResponseCodeNotFound(err) {
			return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
		}
	}
	code, _ := jsonpath.Get("$.Code", response)

	// The documented contract returns NotificationPolicies as an array. Validate
	// the container before any NotFound classification, including the response
	// level Code: a null, scalar or object shape is a malformed response, not
	// evidence that the policy no longer exists. Only a confirmed array that is
	// empty or lacks the requested Id means the policy is gone.
	raw, err := jsonpath.Get("$.PageBean.NotificationPolicies", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.PageBean.NotificationPolicies", response)
	}
	list, ok := raw.([]interface{})
	if !ok {
		return object, errs.WrapErrorf(fmt.Errorf("unexpected NotificationPolicies response type %T", raw), errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}

	if helper.InArray(fmt.Sprint(code), []string{"404"}) || len(list) == 0 {
		return object, errs.WrapErrorf(errs.NotFoundErr("NotificationPolicy", id), errs.NotFoundMsg, response)
	}

	for _, item := range list {
		policy, ok := item.(map[string]interface{})
		if !ok {
			return object, errs.WrapErrorf(fmt.Errorf("unexpected NotificationPolicies element type %T", item), errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
		}
		if fmt.Sprint(policy["Id"]) == id {
			return policy, nil
		}
	}

	return object, errs.WrapErrorf(errs.NotFoundErr("NotificationPolicy", id), errs.NotFoundMsg, response)
}

// DescribeArmsNotificationPolicy >>> Encapsulated.

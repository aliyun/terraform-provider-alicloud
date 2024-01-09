package alicloud

import (
	"fmt"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	rpc "github.com/alibabacloud-go/tea-rpc/client"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

type ArmsServiceV2 struct {
	client *connectivity.AliyunClient
}

// DescribeArmsPrometheusMonitoring <<< Encapsulated get interface for Arms PrometheusMonitoring.

func (s *ArmsServiceV2) DescribeArmsPrometheusMonitoring(id string) (object map[string]interface{}, err error) {

	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	parts := strings.Split(id, ":")
	if len(parts) != 3 {
		err = WrapError(fmt.Errorf("invalid Resource Id %s. Expected parts' length %d, got %d", id, 3, len(parts)))
	}
	action := "GetPrometheusMonitoring"
	conn, err := client.NewArmsClient()
	if err != nil {
		return object, WrapError(err)
	}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	query["ClusterId"] = parts[0]
	query["MonitoringName"] = parts[1]
	query["Type"] = parts[2]
	request["RegionId"] = client.RegionId

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(1*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-08-08"), StringPointer("AK"), query, request, &util.RuntimeOptions{})

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

	code, _ := jsonpath.Get("$.Code", response)
	if InArray(fmt.Sprint(code), []string{"404"}) {
		return object, WrapErrorf(Error(GetNotFoundMessage("PrometheusMonitoring", id)), NotFoundMsg, ProviderERROR, fmt.Sprint(response["RequestId"]))
	}

	v, err := jsonpath.Get("$.Data", response)
	if err != nil {
		return object, WrapErrorf(Error(GetNotFoundMessage("PrometheusMonitoring", id)), NotFoundMsg, ProviderERROR, fmt.Sprint(response["RequestId"]))
	}

	return v.(map[string]interface{}), nil
}

func (s *ArmsServiceV2) ArmsPrometheusMonitoringStateRefreshFunc(id string, field string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeArmsPrometheusMonitoring(id)
		if err != nil {
			if NotFoundError(err) {
				return object, "", nil
			}
			return nil, "", WrapError(err)
		}

		currentStatus := fmt.Sprint(object[field])
		for _, failState := range failStates {
			if currentStatus == failState {
				return object, currentStatus, WrapError(Error(FailedToReachTargetStatus, currentStatus))
			}
		}
		return object, currentStatus, nil
	}
}

// DescribeArmsPrometheusMonitoring >>> Encapsulated.
// DescribeArmsRemoteWrite <<< Encapsulated get interface for Arms RemoteWrite.

func (s *ArmsServiceV2) DescribeArmsRemoteWrite(id string) (object map[string]interface{}, err error) {

	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	parts := strings.Split(id, ":")
	if len(parts) != 2 {
		err = WrapError(fmt.Errorf("invalid Resource Id %s. Expected parts' length %d, got %d", id, 2, len(parts)))
	}
	action := "GetPrometheusRemoteWrite"
	conn, err := client.NewArmsClient()
	if err != nil {
		return object, WrapError(err)
	}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	query["ClusterId"] = parts[0]
	query["RemoteWriteName"] = parts[1]
	request["RegionId"] = client.RegionId

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(1*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-08-08"), StringPointer("AK"), query, request, &util.RuntimeOptions{})

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

	code, _ := jsonpath.Get("$.Code", response)
	if InArray(fmt.Sprint(code), []string{"404"}) {
		return object, WrapErrorf(Error(GetNotFoundMessage("RemoteWrite", id)), NotFoundMsg, ProviderERROR, fmt.Sprint(response["RequestId"]))
	}

	v, err := jsonpath.Get("$.Data", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Data", response)
	}

	return v.(map[string]interface{}), nil
}

func (s *ArmsServiceV2) ArmsRemoteWriteStateRefreshFunc(id string, field string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeArmsRemoteWrite(id)
		if err != nil {
			if NotFoundError(err) {
				return object, "", nil
			}
			return nil, "", WrapError(err)
		}

		currentStatus := fmt.Sprint(object[field])
		for _, failState := range failStates {
			if currentStatus == failState {
				return object, currentStatus, WrapError(Error(FailedToReachTargetStatus, currentStatus))
			}
		}
		return object, currentStatus, nil
	}
}

// DescribeArmsRemoteWrite >>> Encapsulated.

// DescribeArmsEnvironment <<< Encapsulated get interface for Arms Environment.

func (s *ArmsServiceV2) DescribeArmsEnvironment(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	action := "DescribeEnvironment"
	conn, err := client.NewArmsClient()
	if err != nil {
		return object, WrapError(err)
	}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	query["EnvironmentId"] = id
	request["RegionId"] = client.RegionId

	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(1*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-08-08"), StringPointer("AK"), query, request, &runtime)

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

	code, _ := jsonpath.Get("$.Code", response)
	if InArray(fmt.Sprint(code), []string{"404"}) {
		return object, WrapErrorf(Error(GetNotFoundMessage("Environment", id)), NotFoundMsg, response)
	}

	v, err := jsonpath.Get("$.Data", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Data", response)
	}

	return v.(map[string]interface{}), nil
}

func (s *ArmsServiceV2) ArmsEnvironmentStateRefreshFunc(id string, field string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeArmsEnvironment(id)
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

// DescribeArmsEnvironment >>> Encapsulated.

// SetResourceTags <<< Encapsulated tag function for Arms.
func (s *ArmsServiceV2) SetResourceTags(d *schema.ResourceData, resourceType string) error {
	if d.HasChange("tags") {
		var err error
		var action string
		var conn *rpc.Client
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
			conn, err = client.NewArmsClient()
			if err != nil {
				return WrapError(err)
			}
			request = make(map[string]interface{})
			query = make(map[string]interface{})
			request["ResourceId.1"] = d.Id()
			for i, key := range removedTagKeys {
				request[fmt.Sprintf("TagKey.%d", i+1)] = key
			}

			request["ResourceType"] = resourceType
			runtime := util.RuntimeOptions{}
			runtime.SetAutoretry(true)
			wait := incrementalWait(3*time.Second, 5*time.Second)
			err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
				response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-08-08"), StringPointer("AK"), query, request, &runtime)

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
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
			}

		}

		if len(added) > 0 {
			action = "TagResources"
			conn, err = client.NewArmsClient()
			if err != nil {
				return WrapError(err)
			}
			request = make(map[string]interface{})
			query = make(map[string]interface{})
			request["ResourceId.1"] = d.Id()
			count := 1
			for key, value := range added {
				request[fmt.Sprintf("Tag.%d.Key", count)] = key
				request[fmt.Sprintf("Tag.%d.Value", count)] = value
				count++
			}

			request["ResourceType"] = resourceType
			runtime := util.RuntimeOptions{}
			runtime.SetAutoretry(true)
			wait := incrementalWait(3*time.Second, 5*time.Second)
			err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
				response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-08-08"), StringPointer("AK"), query, request, &runtime)

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
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
			}

		}
		d.SetPartial("tags")
	}

	return nil
}

// SetResourceTags >>> tag function encapsulated.
// DescribeArmsPrometheus <<< Encapsulated get interface for Arms Prometheus.

func (s *ArmsServiceV2) DescribeArmsPrometheus(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	action := "GetPrometheusInstance"
	conn, err := client.NewArmsClient()
	if err != nil {
		return object, WrapError(err)
	}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	query["ClusterId"] = id
	request["RegionId"] = client.RegionId

	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(1*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-08-08"), StringPointer("AK"), query, request, &runtime)

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

	code, _ := jsonpath.Get("$.Code", response)
	if InArray(fmt.Sprint(code), []string{"404"}) {
		return object, WrapErrorf(Error(GetNotFoundMessage("Prometheus", id)), NotFoundMsg, response)
	}

	v, err := jsonpath.Get("$.Data", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Data", response)
	}

	return v.(map[string]interface{}), nil
}

func (s *ArmsServiceV2) ArmsPrometheusStateRefreshFunc(id string, field string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeArmsPrometheus(id)
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

// DescribeArmsPrometheus >>> Encapsulated.

// DescribeArmsEnvFeature <<< Encapsulated get interface for Arms EnvFeature.

func (s *ArmsServiceV2) DescribeArmsEnvFeature(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	parts := strings.Split(id, ":")
	if len(parts) != 2 {
		err = WrapError(fmt.Errorf("invalid Resource Id %s. Expected parts' length %d, got %d", id, 2, len(parts)))
	}
	action := "DescribeEnvironmentFeature"
	conn, err := client.NewArmsClient()
	if err != nil {
		return object, WrapError(err)
	}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	query["FeatureName"] = parts[1]
	query["EnvironmentId"] = parts[0]
	request["RegionId"] = client.RegionId

	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(1*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-08-08"), StringPointer("AK"), query, request, &runtime)

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

	code, _ := jsonpath.Get("$.Code", response)
	if InArray(fmt.Sprint(code), []string{"404", "14000"}) {
		return object, WrapErrorf(Error(GetNotFoundMessage("EnvFeature", id)), NotFoundMsg, response)
	}

	v, err := jsonpath.Get("$.Data", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Data", response)
	}

	return v.(map[string]interface{}), nil
}

func (s *ArmsServiceV2) ArmsEnvFeatureStateRefreshFunc(id string, field string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeArmsEnvFeature(id)
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

// DescribeArmsEnvFeature >>> Encapsulated.

// DescribeArmsAddonRelease <<< Encapsulated get interface for Arms AddonRelease.

func (s *ArmsServiceV2) DescribeArmsAddonRelease(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	parts := strings.Split(id, ":")
	if len(parts) != 2 {
		err = WrapError(fmt.Errorf("invalid Resource Id %s. Expected parts' length %d, got %d", id, 2, len(parts)))
	}
	action := "DescribeAddonRelease"
	conn, err := client.NewArmsClient()
	if err != nil {
		return object, WrapError(err)
	}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	query["ReleaseName"] = parts[1]
	query["EnvironmentId"] = parts[0]
	request["RegionId"] = client.RegionId

	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(1*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-08-08"), StringPointer("AK"), query, request, &runtime)

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

	code, _ := jsonpath.Get("$.Code", response)
	if InArray(fmt.Sprint(code), []string{"404"}) {
		return object, WrapErrorf(Error(GetNotFoundMessage("AddonRelease", id)), NotFoundMsg, response)
	}

	v, err := jsonpath.Get("$.Data", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Data", response)
	}

	return v.(map[string]interface{}), nil
}

func (s *ArmsServiceV2) ArmsAddonReleaseStateRefreshFunc(id string, field string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeArmsAddonRelease(id)
		if err != nil {
			if NotFoundError(err) {
				return nil, "", nil
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

// DescribeArmsAddonRelease >>> Encapsulated.

// DescribeArmsEnvPodMonitor <<< Encapsulated get interface for Arms EnvPodMonitor.

func (s *ArmsServiceV2) DescribeArmsEnvPodMonitor(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	parts := strings.Split(id, ":")
	if len(parts) != 3 {
		err = WrapError(fmt.Errorf("invalid Resource Id %s. Expected parts' length %d, got %d", id, 3, len(parts)))
	}
	action := "DescribeEnvPodMonitor"
	conn, err := client.NewArmsClient()
	if err != nil {
		return object, WrapError(err)
	}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	query["PodMonitorName"] = parts[2]
	query["EnvironmentId"] = parts[0]
	query["Namespace"] = parts[1]
	request["RegionId"] = client.RegionId

	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(1*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-08-08"), StringPointer("AK"), query, request, &runtime)

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

	code, _ := jsonpath.Get("$.Code", response)
	if InArray(fmt.Sprint(code), []string{"404"}) {
		return object, WrapErrorf(Error(GetNotFoundMessage("EnvPodMonitor", id)), NotFoundMsg, response)
	}

	v, err := jsonpath.Get("$.Data", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Data", response)
	}

	return v.(map[string]interface{}), nil
}

func (s *ArmsServiceV2) ArmsEnvPodMonitorStateRefreshFunc(id string, field string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeArmsEnvPodMonitor(id)
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

// DescribeArmsEnvPodMonitor >>> Encapsulated.

// DescribeArmsEnvServiceMonitor <<< Encapsulated get interface for Arms EnvServiceMonitor.

func (s *ArmsServiceV2) DescribeArmsEnvServiceMonitor(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	parts := strings.Split(id, ":")
	if len(parts) != 3 {
		err = WrapError(fmt.Errorf("invalid Resource Id %s. Expected parts' length %d, got %d", id, 3, len(parts)))
	}
	action := "DescribeEnvServiceMonitor"
	conn, err := client.NewArmsClient()
	if err != nil {
		return object, WrapError(err)
	}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	query["ServiceMonitorName"] = parts[2]
	query["EnvironmentId"] = parts[0]
	query["Namespace"] = parts[1]
	request["RegionId"] = client.RegionId

	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(1*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-08-08"), StringPointer("AK"), query, request, &runtime)

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

	code, _ := jsonpath.Get("$.Code", response)
	if InArray(fmt.Sprint(code), []string{"404"}) {
		return object, WrapErrorf(Error(GetNotFoundMessage("EnvServiceMonitor", id)), NotFoundMsg, response)
	}

	v, err := jsonpath.Get("$.Data", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Data", response)
	}

	return v.(map[string]interface{}), nil
}

func (s *ArmsServiceV2) ArmsEnvServiceMonitorStateRefreshFunc(id string, field string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeArmsEnvServiceMonitor(id)
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

// DescribeArmsEnvServiceMonitor >>> Encapsulated.

// DescribeArmsEnvCustomJob <<< Encapsulated get interface for Arms EnvCustomJob.

func (s *ArmsServiceV2) DescribeArmsEnvCustomJob(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	parts := strings.Split(id, ":")
	if len(parts) != 2 {
		err = WrapError(fmt.Errorf("invalid Resource Id %s. Expected parts' length %d, got %d", id, 2, len(parts)))
	}
	action := "DescribeEnvCustomJob"
	conn, err := client.NewArmsClient()
	if err != nil {
		return object, WrapError(err)
	}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	query["CustomJobName"] = parts[1]
	query["EnvironmentId"] = parts[0]
	request["RegionId"] = client.RegionId

	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(1*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-08-08"), StringPointer("AK"), query, request, &runtime)

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

	code, _ := jsonpath.Get("$.Code", response)
	if InArray(fmt.Sprint(code), []string{"404"}) {
		return object, WrapErrorf(Error(GetNotFoundMessage("EnvCustomJob", id)), NotFoundMsg, response)
	}

	v, err := jsonpath.Get("$.Data", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Data", response)
	}

	return v.(map[string]interface{}), nil
}

func (s *ArmsServiceV2) ArmsEnvCustomJobStateRefreshFunc(id string, field string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeArmsEnvCustomJob(id)
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

// DescribeArmsEnvCustomJob >>> Encapsulated.


// DescribeArmsXTraceApp <<< Encapsulated get interface for Arms XTraceApp.

	func (s *ArmsServiceV2) DescribeArmsXTraceApp(id string) (object map[string]interface{}, err error) {
			client := s.client
			var request map[string]interface{}
	var response map[string]interface{}
			var query map[string]interface{}
				                    action := "DescribeXTraceApp"
            conn, err := client.NewArmsClient()
        if err != nil {
        return object, WrapError(err)
    }
    request = make(map[string]interface{})
            query = make(map[string]interface{})
                    query["Pid"] = id
            query["RegionId"] = client.RegionId
    
            runtime := util.RuntimeOptions{}
    runtime.SetAutoretry(true)
    wait := incrementalWait(3*time.Second, 5*time.Second)
        err = resource.Retry(1*time.Minute,  func() *resource.RetryError {
                        response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("GET"), StringPointer("2019-08-08"), StringPointer("AK"), query, request, &runtime)
                
        if err != nil {
                        if  NeedRetry(err) {
                wait()
                return resource.RetryableError(err)
            }
            return resource.NonRetryableError(err)
        }
        addDebug(action, response, request)
        return nil
    })

            if err != nil {
                        if IsExpectedErrors(err, []string{"AppNotExist"}) {
            return object, WrapErrorf(Error(GetNotFoundMessage("XTraceApp", id)), NotFoundMsg, response)
        }
                return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
    }
    
		
		
										
		
		
										
		return response, nil
	}

	func (s *ArmsServiceV2) ArmsXTraceAppStateRefreshFunc(id string, field string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
							object, err := s.DescribeArmsXTraceApp(id)
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


// DescribeArmsXTraceApp >>> Encapsulated.

package dts

import (
	"encoding/json"
	"fmt"
	"sort"
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

type Service struct {
	client *connectivity.AliyunClient
}

func NewService(client *connectivity.AliyunClient) *Service {
	return &Service{client: client}
}

func (s *Service) DescribeSubscriptionJob(id string) (object map[string]interface{}, err error) {
	client := s.client
	var response map[string]interface{}
	action := "DescribeDtsJobDetail"
	request := map[string]interface{}{
		"RegionId": s.client.RegionId,
		"DtsJobId": id,
	}
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("Dts", "2020-01-01", action, nil, request, true)
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
		if errs.IsExpectedErrors(err, []string{"Forbidden.InstanceNotFound"}) {
			return object, errs.WrapErrorf(errs.NotFoundErr("DTS:SubscriptionJob", id), errs.NotFoundMsg, errs.ProviderERROR, fmt.Sprint(response["RequestId"]))
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}
	if fmt.Sprint(response["Success"]) == "false" {
		return object, errs.WrapError(fmt.Errorf("%s failed, response: %v", action, response))
	}
	v, err := jsonpath.Get("$", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$", response)
	}
	object = v.(map[string]interface{})
	if object["Status"] == "Starting" {
		object["Status"] = "Normal"
	}
	return object, nil
}

func (s *Service) SubscriptionJobStateRefreshFunc(id string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeSubscriptionJob(id)
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

func (s *Service) ListTagResources(id string, resourceType string) (object interface{}, err error) {
	client := s.client
	action := "ListTagResources"
	request := map[string]interface{}{
		"RegionId":     s.client.RegionId,
		"ResourceType": resourceType,
		"ResourceId.1": id,
	}
	tags := make([]interface{}, 0)
	var response map[string]interface{}

	for {
		wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(5*time.Minute, func() *retry.RetryError {
			response, err := client.RpcPost("Dts", "2020-01-01", action, nil, request, false)
			if err != nil {
				if errs.IsExpectedErrors(err, []string{errs.Throttling}) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			logs.AddDebug(action, response, request)
			v, err := jsonpath.Get("$.TagResources.TagResource", response)
			if err != nil {
				return retry.NonRetryableError(errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.TagResources.TagResource", response))
			}
			if v != nil {
				tags = append(tags, v.([]interface{})...)
			}
			return nil
		})
		if err != nil {
			err = errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
			return
		}
		if response["NextToken"] == nil {
			break
		}
		request["NextToken"] = response["NextToken"]
	}

	return tags, nil
}

func (s *Service) SetResourceTags(d *schema.ResourceData, resourceType string) error {

	if d.HasChange("tags") {
		added, removed := tags.ParsingTags(d)
		client := s.client

		removedTagKeys := make([]string, 0)
		for _, v := range removed {
			if !tags.IgnoredTags(v, "") {
				removedTagKeys = append(removedTagKeys, v)
			}
		}
		if len(removedTagKeys) > 0 {
			action := "UntagResources"
			request := map[string]interface{}{
				"RegionId":     s.client.RegionId,
				"ResourceType": resourceType,
				"ResourceId.1": d.Id(),
			}

			if resourceType == "ALIYUN::DTS::INSTANCE:JOB" {
				request["ResourceType"] = "ALIYUN::DTS::INSTANCE"
				request["ResourceId.1"] = d.Get("dts_instance_id")
			}

			for i, key := range removedTagKeys {
				request[fmt.Sprintf("TagKey.%d", i+1)] = key
			}
			wait := helper.IncrementalWait(2*time.Second, 1*time.Second)
			err := retry.Retry(10*time.Minute, func() *retry.RetryError {
				response, err := client.RpcPost("Dts", "2020-01-01", action, nil, request, false)
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
				return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
			}
		}
		if len(added) > 0 {
			action := "TagResources"
			request := map[string]interface{}{
				"RegionId":     s.client.RegionId,
				"ResourceType": resourceType,
				"ResourceId.1": d.Id(),
			}

			if resourceType == "ALIYUN::DTS::INSTANCE:JOB" {
				request["ResourceType"] = "ALIYUN::DTS::INSTANCE"
				request["ResourceId.1"] = d.Get("dts_instance_id")
			}

			count := 1
			for key, value := range added {
				request[fmt.Sprintf("Tag.%d.Key", count)] = key
				request[fmt.Sprintf("Tag.%d.Value", count)] = value
				count++
			}

			wait := helper.IncrementalWait(2*time.Second, 1*time.Second)
			err := retry.Retry(10*time.Minute, func() *retry.RetryError {
				response, err := client.RpcPost("Dts", "2020-01-01", action, nil, request, false)
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
				return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
			}
		}
	}
	return nil
}

func (s *Service) DescribeSynchronizationJob(id string) (object map[string]interface{}, err error) {
	client := s.client
	var response map[string]interface{}
	action := "DescribeDtsJobDetail"
	request := map[string]interface{}{
		"RegionId": s.client.RegionId,
		"DtsJobId": id,
	}
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("Dts", "2020-01-01", action, nil, request, true)
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
		if errs.IsExpectedErrors(err, []string{"Forbidden.InstanceNotFound"}) {
			return object, errs.WrapErrorf(errs.NotFoundErr("DTS:SynchronizationJob", id), errs.NotFoundMsg, errs.ProviderERROR, fmt.Sprint(response["RequestId"]))
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}
	if fmt.Sprint(response["Success"]) == "false" {
		return object, errs.WrapError(fmt.Errorf("%s failed, response: %v", action, response))
	}
	v, err := jsonpath.Get("$", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$", response)
	}
	object = v.(map[string]interface{})
	// From startup to synchronizing, the resource goes through the pre-check and initialization phases.
	// Because the initialization phase takes too long, and once it passes the pre-check, it can be considered that the task can be executed normally.
	// Therefore, the initialization state is regarded as an equivalent synchronization state.
	if object["Status"] == "synchronizing" || object["Status"] == "Initializing" {
		object["Status"] = "Synchronizing"
	}
	// After calling the delete interface, soft delete is actually performed,
	// and instance data will still be queried from the query interface.
	// Therefore, it is judged that the instance does not exist when the end state is "Finished".
	if object["Status"] == "Finished" {
		return object, errs.WrapErrorf(errs.NotFoundErr("DTS:SynchronizationJob", id), errs.NotFoundMsg, errs.ProviderERROR, fmt.Sprint(response["RequestId"]))
	}
	return object, nil
}

func (s *Service) SynchronizationJobStateRefreshFunc(id string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeSynchronizationJob(id)
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

func (s *Service) QueryChangedJobParameters(dtsJobId string) (string, error) {
	client := s.client
	action := "DescribeDtsJobConfig"
	request := map[string]interface{}{
		"RegionId": client.RegionId,
		"DtsJobId": dtsJobId,
	}
	resp, err := client.RpcPost("Dts", "2020-01-01", action, nil, request, true)
	if err != nil {
		return "", err
	}

	logs.AddDebug(action, resp, request)

	// Parameters 必须为 []interface{}
	parameters, ok := resp["Parameters"].([]interface{})
	if !ok {
		return "", fmt.Errorf("no Parameters found in response")
	}

	changed := make([]map[string]interface{}, 0)

	// filter 被修改的项
	for _, item := range parameters {
		param, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		rv, rvOk := param["RunningValue"]
		dv, dvOk := param["DefaultValue"]
		if rvOk && dvOk && rv != dv {
			changed = append(changed, map[string]interface{}{
				"module": param["Module"],
				"name":   param["Name"],
				"value":  rv,
			})
		}
	}

	if len(changed) == 0 {
		return "", nil
	}

	// 按 name 排序
	sort.Slice(changed, func(i, j int) bool {
		ni, _ := changed[i]["name"].(string)
		nj, _ := changed[j]["name"].(string)
		return ni < nj
	})

	// 转一行字符串
	b, err := json.Marshal(changed)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// DescribeJobsByInstanceId lists the DTS jobs currently attached to the given DTS instance.
// An instance that carries historical deletion records can lose the server-side mapping between
// the job ID stored in state and the instance, so listing the jobs by instance ID is the
// reliable way to recover the currently valid job ID.
func (s *Service) DescribeJobsByInstanceId(instanceId, jobType string) (jobs []map[string]interface{}, err error) {
	client := s.client
	var response map[string]interface{}
	action := "DescribeDtsJobs"
	request := map[string]interface{}{
		"RegionId":      client.RegionId,
		"DtsInstanceId": instanceId,
		"JobType":       jobType,
		"PageSize":      helper.PageSizeLarge,
		"PageNumber":    1,
	}
	for {
		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(5*time.Minute, func() *retry.RetryError {
			response, err = client.RpcPost("Dts", "2020-01-01", action, nil, request, true)
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
			if errs.IsExpectedErrors(err, []string{"Forbidden.InstanceNotFound", "InvalidJobId"}) {
				return jobs, errs.WrapErrorf(errs.NotFoundErr("DTS", instanceId), errs.NotFoundWithResponse, response)
			}
			return jobs, errs.WrapErrorf(err, errs.DefaultErrorMsg, instanceId, action, errs.AlibabaCloudSdkGoERROR)
		}
		v, err := jsonpath.Get("$.DtsJobList", response)
		if err != nil {
			return jobs, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, instanceId, "$.DtsJobList", response)
		}
		list, ok := v.([]interface{})
		if !ok {
			return jobs, errs.WrapErrorf(errs.NotFoundErr("DTS", instanceId), errs.NotFoundWithResponse, response)
		}
		for _, item := range list {
			if job, ok := item.(map[string]interface{}); ok {
				jobs = append(jobs, job)
			}
		}
		if len(list) < request["PageSize"].(int) {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}
	return jobs, nil
}

// isJobNotFoundSignal reports whether a failed DTS job modification means that the job ID
// stored in state no longer maps to a valid DTS instance, for example when the instance carries
// historical deletion records and the job ID and the instance ID no longer match.
func isJobNotFoundSignal(err error, response map[string]interface{}) bool {
	notFoundCodes := []string{"Forbidden.InstanceNotFound", "InvalidJobId"}
	if err != nil {
		return errs.IsExpectedErrors(err, notFoundCodes)
	}
	if fmt.Sprint(response["Success"]) == "false" {
		errCode := fmt.Sprint(response["ErrCode"])
		for _, code := range notFoundCodes {
			if errCode == code {
				return true
			}
		}
	}
	return false
}

// setEndpointSSL folds source_endpoint_ssl and destination_endpoint_ssl into the DTS Reserve
// parameter as its srcSSL and destSSL keys. ConfigureDtsJob has no dedicated SSL parameter, so
// Reserve is the only place the connection mode can be sent. Values already present under those
// keys in a user-supplied reserve are overwritten by the dedicated attributes.
func setEndpointSSL(d *schema.ResourceData, request map[string]interface{}) error {
	srcSSL, srcOk := d.GetOk("source_endpoint_ssl")
	destSSL, destOk := d.GetOk("destination_endpoint_ssl")
	if !srcOk && !destOk {
		return nil
	}
	reserve := make(map[string]interface{})
	if v, ok := request["Reserve"]; ok {
		if err := json.Unmarshal([]byte(fmt.Sprint(v)), &reserve); err != nil {
			return errs.WrapError(err)
		}
	}
	if srcOk {
		reserve["srcSSL"] = srcSSL
	}
	if destOk {
		reserve["destSSL"] = destSSL
	}
	reserveJson, err := json.Marshal(reserve)
	if err != nil {
		return errs.WrapError(err)
	}
	request["Reserve"] = string(reserveJson)
	return nil
}

// convertEndpointSslResponse maps the SslSolutionEnum returned by DescribeDtsJobDetail and
// DescribeDtsJobs back to the srcSSL/destSSL code accepted by ConfigureDtsJob, so that
// source_endpoint_ssl and destination_endpoint_ssl round-trip on read and import.
//
// The MongoDB Atlas value is matched under both spellings because the API metadata disagrees with
// itself: SslSolutionEnum carries no machine-readable value list, and its documented description
// says ENABLE_ONLY_4_MONGODB_ALTAS while its title says ENABLE_ONLY_4_MONGODB_ATLAS.
//
// Anything unrecognised returns nil rather than the raw value, so that callers leave the attribute
// alone. Storing an unmapped enum would put a value outside the ValidateFunc domain into state,
// which can never equal a legal config value and so leaves a diff that no apply can settle.
func convertEndpointSslResponse(source interface{}) interface{} {
	switch source {
	case "DISABLE":
		return "0"
	case "ENABLE_WITH_CERTIFICATE", "ENABLE_ONLY_4_MONGODB_ALTAS", "ENABLE_ONLY_4_MONGODB_ATLAS":
		return "1"
	case "ENABLE_ONLY_4_KAFKA_SCRAM_SHA_256":
		return "3"
	}
	return nil
}

// endpointSSLReserved builds the Reserved payload for a ModifyDtsJob call issued with
// ModifyTypeEnum=UPDATE_RESERVED. ModifyDtsJob merges Reserved into the job's existing reserve
// instead of replacing it, so only the SSL keys need to be sent and any other reserve keys the
// job already carries are left alone.
//
// Each key is gated on its own change. Both attributes are Computed, so an endpoint the user never
// configured still reports a value here - one derived from the last read, not from the config.
// Sending it back would re-assert a connection mode the user never asked for, and because the read
// mapping is lossy it need not be the mode the endpoint actually had.
func endpointSSLReserved(d *schema.ResourceData) (string, error) {
	reserved := make(map[string]interface{})
	if d.HasChange("source_endpoint_ssl") {
		reserved["srcSSL"] = d.Get("source_endpoint_ssl")
	}
	if d.HasChange("destination_endpoint_ssl") {
		reserved["destSSL"] = d.Get("destination_endpoint_ssl")
	}
	reservedJson, err := json.Marshal(reserved)
	if err != nil {
		return "", errs.WrapError(err)
	}
	return string(reservedJson), nil
}

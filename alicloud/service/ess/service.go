package ess

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ess"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/tags"
)

type EssService struct {
	client *connectivity.AliyunClient
}

// NewEssService constructs an EssService from an AliyunClient. The client field is
// unexported, so external packages (including the ess_test ACC package) must build
// instances through this constructor rather than via a struct literal.
func NewEssService(client *connectivity.AliyunClient) *EssService {
	return &EssService{client: client}
}

func (s *EssService) DescribeInstances(scalingGroupId string, lifecycleState string) (instances []ess.ScalingInstance, err error) {
	request := ess.CreateDescribeScalingInstancesRequest()
	request.RegionId = s.client.RegionId
	request.ScalingGroupId = scalingGroupId
	request.LifecycleState = lifecycleState
	v, err := s.client.WithEssClient(func(essClient *ess.Client) (interface{}, error) {
		return essClient.DescribeScalingInstances(request)
	})
	if err != nil {
		return instances, errs.WrapErrorf(err, errs.DefaultErrorMsg, scalingGroupId, request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
	}

	resp, _ := v.(*ess.DescribeScalingInstancesResponse)
	return resp.ScalingInstances.ScalingInstance, nil
}

func (s *EssService) DescribeEssAlarm(id string) (alarm ess.Alarm, err error) {
	request := ess.CreateDescribeAlarmsRequest()
	request.RegionId = s.client.RegionId
	request.AlarmTaskId = id
	request.MetricType = "system"

	var Alarms interface{}
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		Alarms, err = s.client.WithEssClient(func(essClient *ess.Client) (interface{}, error) {
			return essClient.DescribeAlarms(request)
		})
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	if err != nil {
		return alarm, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
	}
	logs.AddDebug(request.GetActionName(), Alarms, request.RpcRequest, request)
	AlarmsResponse, _ := Alarms.(*ess.DescribeAlarmsResponse)
	systemAlarms := AlarmsResponse.AlarmList.Alarm

	if len(systemAlarms) > 0 {
		return systemAlarms[0], nil
	}

	AlarmsRequest := ess.CreateDescribeAlarmsRequest()
	AlarmsRequest.RegionId = s.client.RegionId
	AlarmsRequest.AlarmTaskId = id
	AlarmsRequest.MetricType = "custom"

	var raw interface{}
	wait = helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		raw, err = s.client.WithEssClient(func(essClient *ess.Client) (interface{}, error) {
			return essClient.DescribeAlarms(AlarmsRequest)
		})
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	if err != nil {
		return alarm, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, AlarmsRequest.GetActionName(), errs.AlibabaCloudSdkGoERROR)
	}
	logs.AddDebug(AlarmsRequest.GetActionName(), raw, AlarmsRequest.RpcRequest, AlarmsRequest)
	response, _ := raw.(*ess.DescribeAlarmsResponse)
	customAlarms := response.AlarmList.Alarm

	if len(customAlarms) > 0 {
		return customAlarms[0], nil
	}
	return alarm, errs.WrapErrorf(errs.NotFoundErr("EssAlarm", id), errs.NotFoundMsg, errs.ProviderERROR)
}

func (s *EssService) DescribeEssLifecycleHook(id string) (hook ess.LifecycleHook, err error) {
	request := ess.CreateDescribeLifecycleHooksRequest()
	request.LifecycleHookId = &[]string{id}
	request.RegionId = s.client.RegionId
	raw, err := s.client.WithEssClient(func(essClient *ess.Client) (interface{}, error) {
		return essClient.DescribeLifecycleHooks(request)
	})
	if err != nil {
		err = errs.WrapErrorf(err, errs.DefaultErrorMsg, id, request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
		return
	}
	logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*ess.DescribeLifecycleHooksResponse)
	for _, v := range response.LifecycleHooks.LifecycleHook {
		if v.LifecycleHookId == id {
			return v, nil
		}
	}
	err = errs.WrapErrorf(errs.NotFoundErr("EssLifecycleHook", id), errs.NotFoundMsg, errs.ProviderERROR)
	return
}

func (s *EssService) WaitForEssLifecycleHook(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)

	for {
		object, err := s.DescribeEssLifecycleHook(id)
		if err != nil {
			if errs.NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return errs.WrapError(err)
			}
		}
		if object.LifecycleHookId == id && status != Deleted {
			return nil
		}
		if time.Now().After(deadline) {
			return errs.WrapErrorf(err, errs.WaitTimeoutMsg, id, helper.GetFunc(1), timeout, object.LifecycleHookId, id, errs.ProviderERROR)
		}
	}
}

func (s *EssService) DescribeEssNotification(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	parts := strings.SplitN(id, ":", 2)
	scalingGroupId, notificationArn := parts[0], parts[1]
	request := map[string]interface{}{
		"ScalingGroupId": scalingGroupId,
		"RegionId":       s.client.RegionId,
	}
	response, err = client.RpcPost("Ess", "2014-08-28", "DescribeNotificationConfigurations", nil, request, true)
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"NotificationConfigurationNotExist", "InvalidScalingGroupId.NotFound"}) {
			err = errs.WrapErrorf(errs.NotFoundErr("EssNotification", id), errs.NotFoundMsg, errs.ProviderERROR)
		}
		err = errs.WrapErrorf(err, errs.DefaultErrorMsg, id, "DescribeNotificationConfigurations", errs.AlibabaCloudSdkGoERROR)
		return
	}

	logs.AddDebug("DescribeNotificationConfigurations", response, request)

	v, err := jsonpath.Get("$.NotificationConfigurationModels", response)

	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.NotificationConfigurationModels", response)
	}

	vv, err := jsonpath.Get("$.NotificationConfigurationModel", v)

	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.NotificationConfigurationModel", response)
	}

	for _, w := range vv.([]interface{}) {
		m := w.(map[string]interface{})
		if m["NotificationArn"] == notificationArn {
			return m, nil
		}
	}

	err = errs.WrapErrorf(errs.NotFoundErr("EssNotificationConfiguration", id), errs.NotFoundMsg, errs.ProviderERROR)
	return

}

func (s *EssService) WaitForEssNotification(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribeEssNotification(id)
		if err != nil {
			if errs.NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return errs.WrapError(err)
			}
		}
		resourceId := fmt.Sprintf("%s:%s", object["ScalingGroupId"], object["NotificationArn"])
		if resourceId == id && status != Deleted {
			return nil
		}
		if time.Now().After(deadline) {
			return errs.WrapErrorf(err, errs.WaitTimeoutMsg, id, helper.GetFunc(1), timeout, resourceId, id, errs.ProviderERROR)
		}
	}
}

func (s *EssService) ActivityStateRefreshFunc(activityId string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {

		request := ess.CreateDescribeScalingActivitiesRequest()
		request.ScalingActivityId = &[]string{activityId}
		request.RegionId = s.client.RegionId
		raw, e := s.client.WithEssClient(func(essClient *ess.Client) (interface{}, error) {
			return essClient.DescribeScalingActivities(request)
		})
		if e != nil {
			return nil, "", errs.WrapErrorf(e, activityId, request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
		}

		response, _ := raw.(*ess.DescribeScalingActivitiesResponse)
		for _, v := range response.ScalingActivities.ScalingActivity {
			if v.ScalingActivityId == activityId {
				for _, failState := range failStates {
					if fmt.Sprint(v.StatusCode) == failState {
						return v, fmt.Sprint(v.StatusCode), errs.WrapError(errs.Error(errs.FailedToReachTargetStatus, fmt.Sprint(v.StatusCode)))
					}
				}
				return v, fmt.Sprint(v.StatusCode), nil
			}
		}
		return nil, "", errs.Error("activity not found")
	}
}

func (s *EssService) DescribeEssScalingGroupById(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	request := map[string]interface{}{
		"ScalingGroupId.1": id,
		"RegionId":         s.client.RegionId,
	}

	response, err = client.RpcPost("Ess", "2014-08-28", "DescribeScalingGroups", nil, request, true)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, "DescribeScalingGroups", errs.AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.ScalingGroups.ScalingGroup", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.ScalingGroups.ScalingGroup", response)
	}
	if len(v.([]interface{})) == 0 {
		return object, errs.WrapErrorf(errs.NotFoundErr("ScalingGroups", id), errs.NotFoundMsg, errs.ProviderERROR)
	}

	object = v.([]interface{})[0].(map[string]interface{})
	return object, nil

}

// Deprecated: Use method DescribeEssScalingGroupById instead.
func (s *EssService) DescribeEssScalingGroup(id string) (group ess.ScalingGroup, err error) {
	request := ess.CreateDescribeScalingGroupsRequest()
	request.ScalingGroupId = &[]string{id}
	request.RegionId = s.client.RegionId
	raw, e := s.client.WithEssClient(func(essClient *ess.Client) (interface{}, error) {
		return essClient.DescribeScalingGroups(request)
	})
	if e != nil {
		err = errs.WrapErrorf(e, id, request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
		return
	}
	logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*ess.DescribeScalingGroupsResponse)
	for _, v := range response.ScalingGroups.ScalingGroup {
		if v.ScalingGroupId == id {
			return v, nil
		}
	}
	err = errs.WrapErrorf(errs.NotFoundErr("EssScalingGroup", id), errs.NotFoundMsg, errs.ProviderERROR)
	return
}

func (s *EssService) DescribeEssScalingGroupSuspendProcess(id string) (object map[string]interface{}, err error) {
	strs, err := helper.ParseResourceId(id, 2)
	if err != nil {
		return object, errs.WrapError(err)
	}
	idExist := false
	scalingGroupId, processTemp := strs[0], strs[1]
	request := ess.CreateDescribeScalingGroupsRequest()
	request.ScalingGroupId = &[]string{scalingGroupId}
	request.RegionId = s.client.RegionId
	raw, e := s.client.WithEssClient(func(essClient *ess.Client) (interface{}, error) {
		return essClient.DescribeScalingGroups(request)
	})
	if e != nil {
		err = errs.WrapErrorf(e, id, request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
		return object, err
	}
	logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*ess.DescribeScalingGroupsResponse)
	for _, v := range response.ScalingGroups.ScalingGroup {
		if v.ScalingGroupId == scalingGroupId {
			process := v.SuspendedProcesses.SuspendedProcess
			for i := 0; i < len(process); i++ {
				if strings.EqualFold(processTemp, process[i]) {
					idExist = true
					return object, nil
				}
			}
		}
	}
	if !idExist {
		return object, errs.WrapErrorf(errs.NotFoundErr("EssScalingGroup", id), errs.NotFoundWithResponse, response)
	}
	return
}

func (s *EssService) DescribeEssScalingConfiguration(id string) (config ess.ScalingConfigurationInDescribeScalingConfigurations, err error) {
	request := ess.CreateDescribeScalingConfigurationsRequest()
	request.ScalingConfigurationId = &[]string{id}
	request.RegionId = s.client.RegionId
	raw, err := s.client.WithEssClient(func(essClient *ess.Client) (interface{}, error) {
		return essClient.DescribeScalingConfigurations(request)
	})
	if err != nil {
		err = errs.WrapErrorf(err, errs.DefaultErrorMsg, id, request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
		return
	}
	logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*ess.DescribeScalingConfigurationsResponse)
	for _, v := range response.ScalingConfigurations.ScalingConfiguration {
		if v.ScalingConfigurationId == id {
			return v, nil
		}
	}

	err = errs.GetNotFoundErrorFromString(errs.GetNotFoundMessage("Scaling Configuration", id))
	return
}

func (s *EssService) DescribeEssEciScalingConfiguration(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	request := map[string]interface{}{
		"ScalingConfigurationId.1": id,
		"RegionId":                 s.client.RegionId,
	}

	response, err = client.RpcPost("Ess", "2014-08-28", "DescribeEciScalingConfigurations", nil, request, true)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, "DescribeEciScalingConfigurations", errs.AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.ScalingConfigurations", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.ScalingConfigurations", response)
	}
	if len(v.([]interface{})) == 0 {
		return object, errs.WrapErrorf(errs.NotFoundErr("EciScalingConfiguration", id), errs.NotFoundMsg, errs.ProviderERROR)
	}

	object = v.([]interface{})[0].(map[string]interface{})
	return object, nil
}

//********

func (s *EssService) DescribeEssScalingConfigurationByCommonApi(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	request := map[string]interface{}{
		"ScalingConfigurationId.1": id,
		"RegionId":                 s.client.RegionId,
	}

	response, err = client.RpcPost("Ess", "2014-08-28", "DescribeScalingConfigurations", nil, request, true)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, "DescribeScalingConfigurations", errs.AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$.ScalingConfigurations", response)

	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.ScalingConfigurations", response)
	}

	vv, err := jsonpath.Get("$.ScalingConfiguration", v)

	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.ScalingConfigurations", response)
	}

	if len(vv.([]interface{})) == 0 {
		return object, errs.WrapErrorf(errs.NotFoundErr("ScalingConfiguration", id), errs.NotFoundMsg, errs.ProviderERROR)
	}

	object = vv.([]interface{})[0].(map[string]interface{})
	//.(map[string]interface{}) .(map[string]interface{})
	w := object["WeightedCapacities"]
	ww := object["InstanceTypes"]
	if w != nil && ww != nil {
		weightedCapacity := w.(map[string]interface{})["WeightedCapacity"].([]interface{})
		instanceType := ww.(map[string]interface{})["InstanceType"].([]interface{})
		instanceTypeOverride := make([]ess.ModifyScalingConfigurationInstanceTypeOverride, 0)
		if len(weightedCapacity) != 0 && len(instanceType) != 0 {
			for i := 0; i < len(weightedCapacity); i++ {
				l := ess.ModifyScalingConfigurationInstanceTypeOverride{
					InstanceType:     instanceType[i].(string),
					WeightedCapacity: weightedCapacity[i].(json.Number).String(),
				}
				instanceTypeOverride = append(instanceTypeOverride, l)
			}
			m := make(map[string][]ess.ModifyScalingConfigurationInstanceTypeOverride)
			m["InstanceTypeOverride"] = instanceTypeOverride
			object["InstanceTypeOverrides"] = m
		}
	}
	return object, nil
}

func (s *EssService) ActiveEssScalingConfiguration(sgId, id string) error {
	request := ess.CreateModifyScalingGroupRequest()
	request.ScalingGroupId = sgId
	request.ActiveScalingConfigurationId = id
	request.RegionId = s.client.RegionId
	raw, err := s.client.WithEssClient(func(essClient *ess.Client) (interface{}, error) {
		return essClient.ModifyScalingGroup(request)
	})
	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, id, request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
	}
	logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)
	return err
}

func (s *EssService) WaitForScalingConfiguration(id string, status Status, timeout int) (err error) {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)

	for {
		object, err := s.DescribeEssScalingConfiguration(id)
		if err != nil {
			if errs.NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return errs.WrapError(err)
			}
		}

		if object.ScalingConfigurationId == id && status != Deleted {
			return nil
		}

		if time.Now().After(deadline) {
			return errs.WrapErrorf(err, errs.WaitTimeoutMsg, id, helper.GetFunc(1), timeout, object.ScalingConfigurationId, id, errs.ProviderERROR)
		}
	}
}

// Flattens an array of datadisk into a []map[string]interface{}
func (s *EssService) flattenDataDiskMappings(list []ess.DataDisk) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(list))
	for _, i := range list {
		l := map[string]interface{}{
			"size":                    i.Size,
			"category":                i.Category,
			"snapshot_id":             i.SnapshotId,
			"device":                  i.Device,
			"delete_with_instance":    i.DeleteWithInstance,
			"encrypted":               i.Encrypted,
			"kms_key_id":              i.KMSKeyId,
			"disk_name":               i.DiskName,
			"description":             i.Description,
			"auto_snapshot_policy_id": i.AutoSnapshotPolicyId,
			"performance_level":       i.PerformanceLevel,
		}
		result = append(result, l)
	}
	return result
}

func (s *EssService) flattenSpotPriceLimitMappings(list []ess.SpotPriceModel) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(list))
	for _, i := range list {
		p, _ := strconv.ParseFloat(strconv.FormatFloat(i.PriceLimit, 'f', 2, 64), 64)
		l := map[string]interface{}{
			"instance_type": i.InstanceType,
			"price_limit":   p,
		}
		result = append(result, l)
	}
	return result
}

func (s *EssService) flattenInstancePatternInfoMappings(list []ess.InstancePatternInfo) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(list))
	for _, i := range list {
		memeory, _ := strconv.ParseFloat(strconv.FormatFloat(i.Memory, 'f', 2, 64), 64)
		maxPrice, _ := strconv.ParseFloat(strconv.FormatFloat(i.MaxPrice, 'f', 2, 64), 64)
		l := map[string]interface{}{
			"instance_family_level":   i.InstanceFamilyLevel,
			"memory":                  memeory,
			"cores":                   i.Cores,
			"max_price":               maxPrice,
			"burstable_performance":   i.BurstablePerformance,
			"architectures":           i.Architectures.Architecture,
			"excluded_instance_types": i.ExcludedInstanceTypes.ExcludedInstanceType,
		}
		result = append(result, l)
	}
	return result
}

func (s *EssService) flattenVserverGroupList(vServerGroups []ess.VServerGroup) []map[string]interface{} {
	groups := make([]map[string]interface{}, 0, len(vServerGroups))
	for _, v := range vServerGroups {
		vserverGroupAttributes := v.VServerGroupAttributes.VServerGroupAttribute
		attrs := make([]map[string]interface{}, 0, len(vserverGroupAttributes))
		for _, a := range vserverGroupAttributes {
			attr := map[string]interface{}{
				"vserver_group_id": a.VServerGroupId,
				"port":             a.Port,
				"weight":           a.Weight,
			}
			attrs = append(attrs, attr)
		}
		group := map[string]interface{}{
			"loadbalancer_id":    v.LoadBalancerId,
			"vserver_attributes": attrs,
		}
		groups = append(groups, group)
	}
	return groups
}

func (s *EssService) DescribeEssScalingRule(id string) (object map[string]interface{}, err error) {

	var response map[string]interface{}
	client := s.client
	request := map[string]interface{}{
		"ScalingRuleId.1": id,
		"RegionId":        s.client.RegionId,
		"ShowAlarmRules":  "false",
	}

	response, err = client.RpcPost("Ess", "2014-08-28", "DescribeScalingRules", nil, request, true)
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"InvalidScalingRuleId.NotFound"}) {
			return object, errs.WrapErrorf(err, errs.NotFoundMsg, errs.AlibabaCloudSdkGoERROR)
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, "DescribeScalingRules", errs.AlibabaCloudSdkGoERROR)
	}
	logs.AddDebug("DescribeScalingRules", response, request)
	get, err := jsonpath.Get("$.ScalingRules", response)
	if err != nil || get == nil {
		return object, errs.WrapErrorf(errs.NotFoundErr("ScalingRules", id), errs.NotFoundMsg, errs.ProviderERROR)
	}
	v, err := jsonpath.Get("$.ScalingRules.ScalingRule", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.ScalingRules.ScalingRule", response)
	}
	if len(v.([]interface{})) == 0 {
		return object, errs.WrapErrorf(errs.NotFoundErr("ScalingRules", id), errs.NotFoundMsg, errs.ProviderERROR)
	}

	object = v.([]interface{})[0].(map[string]interface{})
	return object, nil
}

func (s *EssService) DescribeEssInstanceRefresh(id string) (object map[string]interface{}, err error) {

	var response map[string]interface{}
	client := s.client
	strs, _ := helper.ParseResourceId(id, 2)
	scalingGroupId, instanceRefreshTaskId := strs[0], strs[1]

	request := map[string]interface{}{
		"InstanceRefreshTaskId.1": instanceRefreshTaskId,
		"RegionId":                s.client.RegionId,
		"ScalingGroupId":          scalingGroupId,
	}

	response, err = client.RpcPost("Ess", "2014-08-28", "DescribeInstanceRefreshes", nil, request, true)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, "DescribeInstanceRefreshes", errs.AlibabaCloudSdkGoERROR)
	}
	logs.AddDebug("DescribeInstanceRefreshes", response, request)

	get, err := jsonpath.Get("$.InstanceRefreshTasks", response)
	if err != nil || get == nil {
		return object, errs.WrapErrorf(errs.NotFoundErr("InstanceRefreshTasks", id), errs.NotFoundMsg, errs.ProviderERROR)
	}
	v, err := jsonpath.Get("$.InstanceRefreshTasks.InstanceRefreshTask", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.InstanceRefreshTasks.InstanceRefreshTask", response)
	}
	if len(v.([]interface{})) == 0 {
		return object, errs.WrapErrorf(errs.NotFoundErr("InstanceRefreshTasks", id), errs.NotFoundMsg, errs.ProviderERROR)
	}

	object = v.([]interface{})[0].(map[string]interface{})
	return object, nil
}

func (s *EssService) DescribeEssScalingRuleWithAlarm(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	request := map[string]interface{}{
		"ScalingRuleId.1": id,
		"RegionId":        s.client.RegionId,
		"ShowAlarmRules":  "true",
	}

	response, err = client.RpcPost("Ess", "2014-08-28", "DescribeScalingRules", nil, request, true)
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"InvalidScalingRuleId.NotFound"}) {
			return object, errs.WrapErrorf(err, errs.NotFoundMsg, errs.AlibabaCloudSdkGoERROR)
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, "DescribeScalingRules", errs.AlibabaCloudSdkGoERROR)
	}
	logs.AddDebug("DescribeScalingRules", response, request)
	w, err := jsonpath.Get("$.ScalingRules", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.NotFoundMsg, id, "$.ScalingRules", response)
	}
	v, err := jsonpath.Get("$.ScalingRule", w)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.ScalingRules.ScalingRule", response)
	}
	if len(v.([]interface{})) == 0 {
		return object, errs.WrapErrorf(errs.NotFoundErr("ScalingRules", id), errs.NotFoundMsg, errs.ProviderERROR)
	}

	object = v.([]interface{})[0].(map[string]interface{})
	return object, nil
}

func (s *EssService) WaitForEssScalingRule(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)

	for {
		object, err := s.DescribeEssScalingRule(id)
		if err != nil {
			if errs.NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			}
			return errs.WrapError(err)
		}

		if object["ScalingRuleId"] == id && status != Deleted {
			return nil
		}
		if time.Now().After(deadline) {
			return errs.WrapErrorf(err, errs.WaitTimeoutMsg, id, helper.GetFunc(1), timeout, object["ScalingRuleId"], id, errs.ProviderERROR)
		}
	}
}

func (s *EssService) DescribeEssScheduledTask(id string) (object map[string]interface{}, err error) {

	var response map[string]interface{}
	client := s.client
	request := map[string]interface{}{
		"ScheduledTaskId.1": id,
		"RegionId":          s.client.RegionId,
	}
	response, err = client.RpcPost("Ess", "2014-08-28", "DescribeScheduledTasks", nil, request, true)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, "DescribeScheduledTasks", errs.AlibabaCloudSdkGoERROR)
	}

	logs.AddDebug("DescribeScheduledTasks", response, request, request)

	v, err := jsonpath.Get("$.ScheduledTasks.ScheduledTask", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.ScheduledTasks.ScheduledTask", response)
	}
	if len(v.([]interface{})) == 0 {
		return object, errs.WrapErrorf(errs.NotFoundErr("ScheduledTasks", id), errs.NotFoundMsg, errs.ProviderERROR)
	}

	for _, w := range v.([]interface{}) {
		m := w.(map[string]interface{})
		if m["ScheduledTaskId"] == id {
			return m, nil
		}
	}
	err = errs.WrapErrorf(errs.NotFoundErr("EssScheduledTask", id), errs.NotFoundMsg, errs.ProviderERROR)
	return

}

func (s *EssService) WaitForEssScheduledTask(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribeEssScheduledTask(id)
		if err != nil {
			if errs.NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return errs.WrapError(err)
			}
		}

		if object["TaskEnabled"] == "true" {
			return nil
		}

		time.Sleep(helper.DefaultIntervalShort * time.Second)
		if time.Now().After(deadline) {
			return errs.WrapErrorf(err, errs.WaitTimeoutMsg, id, helper.GetFunc(1), timeout, object["ScheduledTaskId"], id, errs.ProviderERROR)
		}
	}
}

func (srv *EssService) DescribeEssAttachment(id string, instanceIds []string) (instances []ess.ScalingInstance, err error) {
	request := ess.CreateDescribeScalingInstancesRequest()
	request.RegionId = srv.client.RegionId
	request.ScalingGroupId = id
	request.CreationType = "Attached"
	if len(instanceIds) > 0 {
		request.InstanceId = &instanceIds
	}

	raw, err := srv.client.WithEssClient(func(essClient *ess.Client) (interface{}, error) {
		return essClient.DescribeScalingInstances(request)
	})
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"InvalidScalingGroupId.NotFound"}) {
			err = errs.WrapErrorf(err, errs.NotFoundMsg, errs.AlibabaCloudSdkGoERROR)
		} else {
			err = errs.WrapErrorf(err, errs.DefaultErrorMsg, id, request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
		}
		return
	}
	logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*ess.DescribeScalingInstancesResponse)
	if len(response.ScalingInstances.ScalingInstance) < 1 {
		err = errs.WrapErrorf(errs.NotFoundErr("EssAttachment", id), errs.NotFoundMsg, errs.ProviderERROR)
		return
	}
	return response.ScalingInstances.ScalingInstance, nil
}

func (s *EssService) DescribeEssScalingConfifurations(id string) (configs []ess.ScalingConfigurationInDescribeScalingConfigurations, err error) {
	request := ess.CreateDescribeScalingConfigurationsRequest()
	request.ScalingGroupId = id
	request.PageNumber = requests.NewInteger(1)
	request.PageSize = requests.NewInteger(helper.PageSizeLarge)
	request.RegionId = s.client.RegionId
	for {
		raw, err := s.client.WithEssClient(func(essClient *ess.Client) (interface{}, error) {
			return essClient.DescribeScalingConfigurations(request)
		})
		if err != nil {
			return configs, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
		}
		logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)
		response, _ := raw.(*ess.DescribeScalingConfigurationsResponse)
		if len(response.ScalingConfigurations.ScalingConfiguration) < 1 {
			break
		}
		configs = append(configs, response.ScalingConfigurations.ScalingConfiguration...)
		if len(response.ScalingConfigurations.ScalingConfiguration) < helper.PageSizeLarge {
			break
		}

		if page, err := helper.GetNextpageNumber(request.PageNumber); err != nil {
			return configs, errs.WrapError(err)
		} else {
			request.PageNumber = page
		}
	}

	if len(configs) < 1 {
		return configs, errs.WrapErrorf(errs.NotFoundErr("EssScalingConfifuration", id), errs.NotFoundMsg, errs.ProviderERROR)
	}

	return
}

func (srv *EssService) EssRemoveInstances(client *connectivity.AliyunClient, d *schema.ResourceData, id string, instanceIds []string) error {

	if len(instanceIds) < 1 {
		return nil
	}
	group, err := srv.DescribeEssScalingGroup(id)

	if err != nil {
		return errs.WrapError(err)
	}

	if group.LifecycleState == string(Inactive) {
		return errs.WrapError(errs.Error("Scaling group current status is %s, please active it before attaching or removing ECS instances.", group.LifecycleState))
	} else {
		if err := srv.WaitForEssScalingGroup(group.ScalingGroupId, Active, DefaultTimeout); err != nil {
			if errs.NotFoundError(err) {
				return nil
			}
			return errs.WrapError(err)
		}
	}

	removed := instanceIds
	if err := retry.Retry(5*time.Minute, func() *retry.RetryError {
		request := ess.CreateRemoveInstancesRequest()
		request.ScalingGroupId = id
		request.RegionId = srv.client.RegionId
		if len(removed) > 0 {
			request.InstanceId = &removed
		} else {
			return nil
		}
		raw, err := srv.client.WithEssClient(func(essClient *ess.Client) (interface{}, error) {
			return essClient.RemoveInstances(request)
		})
		if err != nil {
			if errs.IsExpectedErrors(err, []string{"InvalidScalingGroupId.NotFound"}) {
				return nil
			}
			if errs.IsExpectedErrors(err, []string{"IncorrectCapacity.MinSize"}) {
				instances, err := srv.DescribeEssAttachment(id, instanceIds)
				if len(instances) > 0 {
					if group.MinSize == 0 {
						return retry.RetryableError(errs.WrapError(err))
					}
					return retry.NonRetryableError(errs.WrapError(err))
				}
			}
			if errs.IsExpectedErrors(err, []string{"ScalingActivityInProgress", "IncorrectScalingGroupStatus"}) {
				time.Sleep(5)
				return retry.RetryableError(errs.WrapErrorf(err, errs.DefaultErrorMsg, id, request.GetActionName(), errs.AlibabaCloudSdkGoERROR))
			}
			return retry.NonRetryableError(errs.WrapErrorf(err, errs.DefaultErrorMsg, id, request.GetActionName(), errs.AlibabaCloudSdkGoERROR))
		}

		response, _ := raw.(*ess.RemoveInstancesResponse)
		essService := EssService{client}
		stateConf := helper.BuildStateConf([]string{}, []string{"Successful"}, d.Timeout(schema.TimeoutCreate), 1*time.Minute, essService.ActivityStateRefreshFunc(response.ScalingActivityId, []string{"Failed", "Rejected"}))
		if _, err := stateConf.WaitForState(); err != nil {
			return retry.NonRetryableError(errs.WrapError(err))
		}

		logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)
		time.Sleep(3 * time.Second)
		instances, err := srv.DescribeEssAttachment(id, instanceIds)
		if err != nil {
			if errs.NotFoundError(err) {
				return nil
			}
			return retry.NonRetryableError(errs.WrapError(err))
		}
		if len(instances) > 0 {
			removed = make([]string, 0)
			for _, inst := range instances {
				removed = append(removed, inst.InstanceId)
			}
			return retry.RetryableError(errs.WrapError(errs.Error("There are still ECS instances in the scaling group.")))
		}

		return nil
	}); err != nil {
		return errs.WrapError(err)
	}
	return nil
}

// WaitForScalingGroup waits for group to given status
func (s *EssService) WaitForEssScalingGroup(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)

	for {
		object, err := s.DescribeEssScalingGroup(id)
		if err != nil {
			if errs.NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return errs.WrapError(err)
			}
		}
		if object.LifecycleState == string(status) {
			return nil
		}
		time.Sleep(helper.DefaultIntervalShort * time.Second)
		if time.Now().After(deadline) {
			return errs.WrapErrorf(err, errs.WaitTimeoutMsg, id, helper.GetFunc(1), timeout, object.LifecycleState, string(status), errs.ProviderERROR)
		}
	}
}

// ess dimensions to map
func (s *EssService) flattenDimensionsToMap(dimensions []ess.Dimension) map[string]string {
	result := make(map[string]string)
	for _, dimension := range dimensions {
		if dimension.DimensionKey == UserId || dimension.DimensionKey == ScalingGroup {
			continue
		}
		result[dimension.DimensionKey] = dimension.DimensionValue
	}
	return result
}

func (s *EssService) WaitForEssAttachment(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)

	for {
		object, err := s.DescribeEssAttachment(id, make([]string, 0))
		if err != nil {
			if errs.NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return errs.WrapError(err)
			}
		}
		if len(object) > 0 && status != Deleted {
			return nil
		}
		time.Sleep(helper.DefaultIntervalShort * time.Second)
		if time.Now().After(deadline) {
			return errs.WrapErrorf(err, errs.WaitTimeoutMsg, id, helper.GetFunc(1), timeout, Null, string(status), errs.ProviderERROR)
		}
	}
}

func (s *EssService) WaitForEssAlarm(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribeEssAlarm(id)
		if err != nil {
			if errs.NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return errs.WrapError(err)
			}
		}
		if object.AlarmTaskId == id && status != Deleted {
			return nil
		}
		time.Sleep(helper.DefaultIntervalShort * time.Second)
		if time.Now().After(deadline) {
			return errs.WrapErrorf(err, errs.WaitTimeoutMsg, id, helper.GetFunc(1), timeout, object.AlarmTaskId, id, errs.ProviderERROR)
		}
	}
}

func (s *EssService) SetResourceTags(d *schema.ResourceData, scalingGroupId string, client *connectivity.AliyunClient) error {

	if d.HasChange("tags") {
		added, removed := tags.ParsingTags(d)

		// untag resources
		if len(removed) > 0 {
			removedTagKeys := make([]string, 0)
			for _, v := range removed {
				if !tags.IgnoredTags(v, "") {
					removedTagKeys = append(removedTagKeys, v)
				}
			}
			untagRequest := ess.CreateUntagResourcesRequest()
			resourceId := []string{scalingGroupId}
			untagRequest.ResourceId = &resourceId
			untagRequest.ResourceType = "scalinggroup"
			untagRequest.TagKey = &removedTagKeys
			raw, err := client.WithEssClient(func(essClient *ess.Client) (interface{}, error) {
				return essClient.UntagResources(untagRequest)
			})
			if err != nil {
				return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), untagRequest.GetActionName(), errs.AlibabaCloudSdkGoERROR)
			}
			logs.AddDebug(untagRequest.GetActionName(), raw, untagRequest.RpcRequest, untagRequest)
		}

		// tag resources
		if len(added) > 0 {
			tagRequest := ess.CreateTagResourcesRequest()
			resourceId := []string{scalingGroupId}
			tagRequest.ResourceId = &resourceId
			tagRequest.ResourceType = "scalinggroup"
			tagList := make([]ess.TagResourcesTag, 0)
			for k, v := range added {
				if !tags.IgnoredTags(v.(string), "") {
					tagList = append(tagList, ess.TagResourcesTag{
						Key:   k,
						Value: v.(string),
					})
				}
			}
			tagRequest.Tag = &tagList
			raw, err := client.WithEssClient(func(essClient *ess.Client) (interface{}, error) {
				return essClient.TagResources(tagRequest)
			})
			if err != nil {
				return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), tagRequest.GetActionName(), errs.AlibabaCloudSdkGoERROR)
			}
			logs.AddDebug(tagRequest.GetActionName(), raw, tagRequest.RpcRequest, tagRequest)
		}
	}
	return nil
}

func (s *EssService) ChangeResourceGroup(d *schema.ResourceData, scalingGroupId string, client *connectivity.AliyunClient) error {

	if d.HasChange("resource_group_id") {
		request := ess.CreateChangeResourceGroupRequest()
		request.ResourceId = scalingGroupId
		request.ResourceType = "scalinggroup"
		request.RegionId = client.RegionId
		if v, ok := d.GetOk("resource_group_id"); ok {
			request.NewResourceGroupId = v.(string)
		}
		raw, err := client.WithEssClient(func(essClient *ess.Client) (interface{}, error) {
			return essClient.ChangeResourceGroup(request)
		})
		if err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
		}
		logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)
	}
	return nil
}

func (s *EssService) ListTagResources(scalingGroupId string, client *connectivity.AliyunClient) (object interface{}, err error) {
	listTagsRequest := ess.CreateListTagResourcesRequest()
	listTagsRequest.ResourceType = "scalinggroup"
	resourceIds := []string{scalingGroupId}
	listTagsRequest.ResourceId = &resourceIds
	raw, err := client.WithEssClient(func(essClient *ess.Client) (interface{}, error) {
		return essClient.ListTagResources(listTagsRequest)
	})
	if err != nil {
		return nil, errs.WrapErrorf(err, errs.DefaultErrorMsg, scalingGroupId, listTagsRequest.GetActionName(), errs.AlibabaCloudSdkGoERROR)
	}
	logs.AddDebug(listTagsRequest.GetActionName(), raw, listTagsRequest.RpcRequest, listTagsRequest)

	response, _ := raw.(*ess.ListTagResourcesResponse)

	tags := make([]interface{}, 0)
	for _, v := range response.TagResources.TagResource {
		tags = append(tags, map[string]interface{}{
			"TagKey":   v.TagKey,
			"TagValue": v.TagValue,
		})
	}
	return tags, nil
}

package alicloud

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ess"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

type EssService struct {
	client *connectivity.AliyunClient
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
		return instances, WrapErrorf(err, DefaultErrorMsg, scalingGroupId, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	resp, _ := v.(*ess.DescribeScalingInstancesResponse)
	return resp.ScalingInstances.ScalingInstance, nil
}

func (s *EssService) DescribeEssAlarm(id string) (alarm ess.Alarm, err error) {
	request := ess.CreateDescribeAlarmsRequest()
	request.RegionId = s.client.RegionId
	request.AlarmTaskId = id
	request.MetricType = "system"
	Alarms, err := s.client.WithEssClient(func(essClient *ess.Client) (interface{}, error) {
		return essClient.DescribeAlarms(request)
	})
	if err != nil {
		return alarm, WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), Alarms, request.RpcRequest, request)
	AlarmsResponse, _ := Alarms.(*ess.DescribeAlarmsResponse)
	systemAlarms := AlarmsResponse.AlarmList.Alarm

	if len(systemAlarms) > 0 {
		return systemAlarms[0], nil
	}

	AlarmsRequest := ess.CreateDescribeAlarmsRequest()
	AlarmsRequest.RegionId = s.client.RegionId
	AlarmsRequest.AlarmTaskId = id
	AlarmsRequest.MetricType = "custom"
	raw, err := s.client.WithEssClient(func(essClient *ess.Client) (interface{}, error) {
		return essClient.DescribeAlarms(AlarmsRequest)
	})
	if err != nil {
		return alarm, WrapErrorf(err, DefaultErrorMsg, id, AlarmsRequest.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(AlarmsRequest.GetActionName(), raw, AlarmsRequest.RpcRequest, AlarmsRequest)
	response, _ := raw.(*ess.DescribeAlarmsResponse)
	customAlarms := response.AlarmList.Alarm

	if len(customAlarms) > 0 {
		return customAlarms[0], nil
	}
	return alarm, WrapErrorf(Error(GetNotFoundMessage("EssAlarm", id)), NotFoundMsg, ProviderERROR)
}

func (s *EssService) DescribeEssLifecycleHook(id string) (hook ess.LifecycleHook, err error) {
	request := ess.CreateDescribeLifecycleHooksRequest()
	request.LifecycleHookId = &[]string{id}
	request.RegionId = s.client.RegionId
	raw, err := s.client.WithEssClient(func(essClient *ess.Client) (interface{}, error) {
		return essClient.DescribeLifecycleHooks(request)
	})
	if err != nil {
		err = WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
		return
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*ess.DescribeLifecycleHooksResponse)
	for _, v := range response.LifecycleHooks.LifecycleHook {
		if v.LifecycleHookId == id {
			return v, nil
		}
	}
	err = WrapErrorf(Error(GetNotFoundMessage("EssLifecycleHook", id)), NotFoundMsg, ProviderERROR)
	return
}

func (s *EssService) WaitForEssLifecycleHook(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)

	for {
		object, err := s.DescribeEssLifecycleHook(id)
		if err != nil {
			if NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return WrapError(err)
			}
		}
		if object.LifecycleHookId == id && status != Deleted {
			return nil
		}
		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, object.LifecycleHookId, id, ProviderERROR)
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
		if IsExpectedErrors(err, []string{"NotificationConfigurationNotExist", "InvalidScalingGroupId.NotFound"}) {
			err = WrapErrorf(Error(GetNotFoundMessage("EssNotification", id)), NotFoundMsg, ProviderERROR)
		}
		err = WrapErrorf(err, DefaultErrorMsg, id, "DescribeNotificationConfigurations", AlibabaCloudSdkGoERROR)
		return
	}

	addDebug("DescribeNotificationConfigurations", response, request)

	v, err := jsonpath.Get("$.NotificationConfigurationModels", response)

	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.NotificationConfigurationModels", response)
	}

	vv, err := jsonpath.Get("$.NotificationConfigurationModel", v)

	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.NotificationConfigurationModel", response)
	}

	for _, w := range vv.([]interface{}) {
		m := w.(map[string]interface{})
		if m["NotificationArn"] == notificationArn {
			return m, nil
		}
	}

	err = WrapErrorf(Error(GetNotFoundMessage("EssNotificationConfiguration", id)), NotFoundMsg, ProviderERROR)
	return

}

func (s *EssService) WaitForEssNotification(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribeEssNotification(id)
		if err != nil {
			if NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return WrapError(err)
			}
		}
		resourceId := fmt.Sprintf("%s:%s", object["ScalingGroupId"], object["NotificationArn"])
		if resourceId == id && status != Deleted {
			return nil
		}
		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, resourceId, id, ProviderERROR)
		}
	}
}

func (s *EssService) ActivityStateRefreshFunc(activityId string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {

		request := ess.CreateDescribeScalingActivitiesRequest()
		request.ScalingActivityId = &[]string{activityId}
		request.RegionId = s.client.RegionId
		raw, e := s.client.WithEssClient(func(essClient *ess.Client) (interface{}, error) {
			return essClient.DescribeScalingActivities(request)
		})
		if e != nil {
			return nil, "", WrapErrorf(e, activityId, request.GetActionName(), AlibabaCloudSdkGoERROR)
		}

		response, _ := raw.(*ess.DescribeScalingActivitiesResponse)
		for _, v := range response.ScalingActivities.ScalingActivity {
			if v.ScalingActivityId == activityId {
				for _, failState := range failStates {
					if fmt.Sprint(v.StatusCode) == failState {
						return v, fmt.Sprint(v.StatusCode), WrapError(Error(FailedToReachTargetStatus, fmt.Sprint(v.StatusCode)))
					}
				}
				return v, fmt.Sprint(v.StatusCode), nil
			}
		}
		return nil, "", Error("activity not found")
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
		return object, WrapErrorf(err, DefaultErrorMsg, id, "DescribeScalingGroups", AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.ScalingGroups.ScalingGroup", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.ScalingGroups.ScalingGroup", response)
	}
	if len(v.([]interface{})) == 0 {
		return object, WrapErrorf(Error(GetNotFoundMessage("ScalingGroups", id)), NotFoundMsg, ProviderERROR)
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
		err = WrapErrorf(e, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
		return
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*ess.DescribeScalingGroupsResponse)
	for _, v := range response.ScalingGroups.ScalingGroup {
		if v.ScalingGroupId == id {
			return v, nil
		}
	}
	err = WrapErrorf(Error(GetNotFoundMessage("EssScalingGroup", id)), NotFoundMsg, ProviderERROR)
	return
}

func (s *EssService) DescribeEssScalingGroupSuspendProcess(id string) (object map[string]interface{}, err error) {
	strs, err := ParseResourceId(id, 2)
	if err != nil {
		return object, WrapError(err)
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
		err = WrapErrorf(e, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
		return object, err
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
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
		return object, WrapErrorf(Error(GetNotFoundMessage("EssScalingGroup", id)), NotFoundWithResponse, response)
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
		err = WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
		return
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*ess.DescribeScalingConfigurationsResponse)
	for _, v := range response.ScalingConfigurations.ScalingConfiguration {
		if v.ScalingConfigurationId == id {
			return v, nil
		}
	}

	err = GetNotFoundErrorFromString(GetNotFoundMessage("Scaling Configuration", id))
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
		return object, WrapErrorf(err, DefaultErrorMsg, id, "DescribeEciScalingConfigurations", AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.ScalingConfigurations", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.ScalingConfigurations", response)
	}
	if len(v.([]interface{})) == 0 {
		return object, WrapErrorf(Error(GetNotFoundMessage("EciScalingConfiguration", id)), NotFoundMsg, ProviderERROR)
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
		return object, WrapErrorf(err, DefaultErrorMsg, id, "DescribeScalingConfigurations", AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$.ScalingConfigurations", response)

	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.ScalingConfigurations", response)
	}

	vv, err := jsonpath.Get("$.ScalingConfiguration", v)

	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.ScalingConfigurations", response)
	}

	if len(vv.([]interface{})) == 0 {
		return object, WrapErrorf(Error(GetNotFoundMessage("ScalingConfiguration", id)), NotFoundMsg, ProviderERROR)
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
		return WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	return err
}

func (s *EssService) WaitForScalingConfiguration(id string, status Status, timeout int) (err error) {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)

	for {
		object, err := s.DescribeEssScalingConfiguration(id)
		if err != nil {
			if NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return WrapError(err)
			}
		}

		if object.ScalingConfigurationId == id && status != Deleted {
			return nil
		}

		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, object.ScalingConfigurationId, id, ProviderERROR)
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

func (s *EssService) DescribeEssScalingRule(id string) (rule ess.ScalingRule, err error) {
	request := ess.CreateDescribeScalingRulesRequest()
	request.ScalingRuleId = &[]string{id}
	request.ShowAlarmRules = "false"
	request.RegionId = s.client.RegionId
	raw, err := s.client.WithEssClient(func(essClient *ess.Client) (interface{}, error) {
		return essClient.DescribeScalingRules(request)
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidScalingRuleId.NotFound"}) {
			return rule, WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
		}
		return rule, WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*ess.DescribeScalingRulesResponse)
	for _, v := range response.ScalingRules.ScalingRule {
		if v.ScalingRuleId == id {
			return v, nil
		}
	}

	return rule, WrapErrorf(Error(GetNotFoundMessage("EssScalingRule", id)), NotFoundMsg, ProviderERROR)
}

func (s *EssService) DescribeEssScalingRuleWithAlarm(id string) (rule ess.ScalingRule, err error) {
	request := ess.CreateDescribeScalingRulesRequest()
	request.ScalingRuleId = &[]string{id}
	request.ShowAlarmRules = "true"
	request.RegionId = s.client.RegionId
	raw, err := s.client.WithEssClient(func(essClient *ess.Client) (interface{}, error) {
		return essClient.DescribeScalingRules(request)
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidScalingRuleId.NotFound"}) {
			return rule, WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
		}
		return rule, WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*ess.DescribeScalingRulesResponse)
	for _, v := range response.ScalingRules.ScalingRule {
		if v.ScalingRuleId == id {
			return v, nil
		}
	}

	return rule, WrapErrorf(Error(GetNotFoundMessage("EssScalingRule", id)), NotFoundMsg, ProviderERROR)
}

func (s *EssService) WaitForEssScalingRule(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)

	for {
		object, err := s.DescribeEssScalingRule(id)
		if err != nil {
			if NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			}
			return WrapError(err)
		}

		if object.ScalingRuleId == id && status != Deleted {
			return nil
		}
		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, object.ScalingRuleId, id, ProviderERROR)
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
		return object, WrapErrorf(err, DefaultErrorMsg, id, "DescribeScheduledTasks", AlibabaCloudSdkGoERROR)
	}

	addDebug("DescribeScheduledTasks", response, request, request)

	v, err := jsonpath.Get("$.ScheduledTasks.ScheduledTask", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.ScheduledTasks.ScheduledTask", response)
	}
	if len(v.([]interface{})) == 0 {
		return object, WrapErrorf(Error(GetNotFoundMessage("ScheduledTasks", id)), NotFoundMsg, ProviderERROR)
	}

	for _, w := range v.([]interface{}) {
		m := w.(map[string]interface{})
		if m["ScheduledTaskId"] == id {
			return m, nil
		}
	}
	err = WrapErrorf(Error(GetNotFoundMessage("EssScheduledTask", id)), NotFoundMsg, ProviderERROR)
	return

}

func (s *EssService) WaitForEssScheduledTask(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribeEssScheduledTask(id)
		if err != nil {
			if NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return WrapError(err)
			}
		}

		if object["TaskEnabled"] == "true" {
			return nil
		}

		time.Sleep(DefaultIntervalShort * time.Second)
		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, object["ScheduledTaskId"], id, ProviderERROR)
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
		if IsExpectedErrors(err, []string{"InvalidScalingGroupId.NotFound"}) {
			err = WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
		} else {
			err = WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		return
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*ess.DescribeScalingInstancesResponse)
	if len(response.ScalingInstances.ScalingInstance) < 1 {
		err = WrapErrorf(Error(GetNotFoundMessage("EssAttachment", id)), NotFoundMsg, ProviderERROR)
		return
	}
	return response.ScalingInstances.ScalingInstance, nil
}

func (s *EssService) DescribeEssScalingConfifurations(id string) (configs []ess.ScalingConfigurationInDescribeScalingConfigurations, err error) {
	request := ess.CreateDescribeScalingConfigurationsRequest()
	request.ScalingGroupId = id
	request.PageNumber = requests.NewInteger(1)
	request.PageSize = requests.NewInteger(PageSizeLarge)
	request.RegionId = s.client.RegionId
	for {
		raw, err := s.client.WithEssClient(func(essClient *ess.Client) (interface{}, error) {
			return essClient.DescribeScalingConfigurations(request)
		})
		if err != nil {
			return configs, WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		response, _ := raw.(*ess.DescribeScalingConfigurationsResponse)
		if len(response.ScalingConfigurations.ScalingConfiguration) < 1 {
			break
		}
		configs = append(configs, response.ScalingConfigurations.ScalingConfiguration...)
		if len(response.ScalingConfigurations.ScalingConfiguration) < PageSizeLarge {
			break
		}

		if page, err := getNextpageNumber(request.PageNumber); err != nil {
			return configs, WrapError(err)
		} else {
			request.PageNumber = page
		}
	}

	if len(configs) < 1 {
		return configs, WrapErrorf(Error(GetNotFoundMessage("EssScalingConfifuration", id)), NotFoundMsg, ProviderERROR)
	}

	return
}

func (srv *EssService) EssRemoveInstances(client *connectivity.AliyunClient, d *schema.ResourceData, id string, instanceIds []string) error {

	if len(instanceIds) < 1 {
		return nil
	}
	group, err := srv.DescribeEssScalingGroup(id)

	if err != nil {
		return WrapError(err)
	}

	if group.LifecycleState == string(Inactive) {
		return WrapError(Error("Scaling group current status is %s, please active it before attaching or removing ECS instances.", group.LifecycleState))
	} else {
		if err := srv.WaitForEssScalingGroup(group.ScalingGroupId, Active, DefaultTimeout); err != nil {
			if NotFoundError(err) {
				return nil
			}
			return WrapError(err)
		}
	}

	removed := instanceIds
	if err := resource.Retry(5*time.Minute, func() *resource.RetryError {
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
			if IsExpectedErrors(err, []string{"InvalidScalingGroupId.NotFound"}) {
				return nil
			}
			if IsExpectedErrors(err, []string{"IncorrectCapacity.MinSize"}) {
				instances, err := srv.DescribeEssAttachment(id, instanceIds)
				if len(instances) > 0 {
					if group.MinSize == 0 {
						return resource.RetryableError(WrapError(err))
					}
					return resource.NonRetryableError(WrapError(err))
				}
			}
			if IsExpectedErrors(err, []string{"ScalingActivityInProgress", "IncorrectScalingGroupStatus"}) {
				time.Sleep(5)
				return resource.RetryableError(WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR))
			}
			return resource.NonRetryableError(WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR))
		}

		response, _ := raw.(*ess.RemoveInstancesResponse)
		essService := EssService{client}
		stateConf := BuildStateConf([]string{}, []string{"Successful"}, d.Timeout(schema.TimeoutCreate), 1*time.Minute, essService.ActivityStateRefreshFunc(response.ScalingActivityId, []string{"Failed", "Rejected"}))
		if _, err := stateConf.WaitForState(); err != nil {
			return resource.NonRetryableError(WrapError(err))
		}

		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		time.Sleep(3 * time.Second)
		instances, err := srv.DescribeEssAttachment(id, instanceIds)
		if err != nil {
			if NotFoundError(err) {
				return nil
			}
			return resource.NonRetryableError(WrapError(err))
		}
		if len(instances) > 0 {
			removed = make([]string, 0)
			for _, inst := range instances {
				removed = append(removed, inst.InstanceId)
			}
			return resource.RetryableError(WrapError(Error("There are still ECS instances in the scaling group.")))
		}

		return nil
	}); err != nil {
		return WrapError(err)
	}
	return nil
}

// WaitForScalingGroup waits for group to given status
func (s *EssService) WaitForEssScalingGroup(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)

	for {
		object, err := s.DescribeEssScalingGroup(id)
		if err != nil {
			if NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return WrapError(err)
			}
		}
		if object.LifecycleState == string(status) {
			return nil
		}
		time.Sleep(DefaultIntervalShort * time.Second)
		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, object.LifecycleState, string(status), ProviderERROR)
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
			if NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return WrapError(err)
			}
		}
		if len(object) > 0 && status != Deleted {
			return nil
		}
		time.Sleep(DefaultIntervalShort * time.Second)
		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, Null, string(status), ProviderERROR)
		}
	}
}

func (s *EssService) WaitForEssAlarm(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribeEssAlarm(id)
		if err != nil {
			if NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return WrapError(err)
			}
		}
		if object.AlarmTaskId == id && status != Deleted {
			return nil
		}
		time.Sleep(DefaultIntervalShort * time.Second)
		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, object.AlarmTaskId, id, ProviderERROR)
		}
	}
}

func (s *EssService) SetResourceTags(d *schema.ResourceData, scalingGroupId string, client *connectivity.AliyunClient) error {

	if d.HasChange("tags") {
		added, removed := parsingTags(d)

		// untag resources
		if len(removed) > 0 {
			removedTagKeys := make([]string, 0)
			for _, v := range removed {
				if !ignoredTags(v, "") {
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
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), untagRequest.GetActionName(), AlibabaCloudSdkGoERROR)
			}
			addDebug(untagRequest.GetActionName(), raw, untagRequest.RpcRequest, untagRequest)
		}

		// tag resources
		if len(added) > 0 {
			tagRequest := ess.CreateTagResourcesRequest()
			resourceId := []string{scalingGroupId}
			tagRequest.ResourceId = &resourceId
			tagRequest.ResourceType = "scalinggroup"
			tags := make([]ess.TagResourcesTag, 0)
			for k, v := range added {
				if !ignoredTags(v.(string), "") {
					tags = append(tags, ess.TagResourcesTag{
						Key:   k,
						Value: v.(string),
					})
				}
			}
			tagRequest.Tag = &tags
			raw, err := client.WithEssClient(func(essClient *ess.Client) (interface{}, error) {
				return essClient.TagResources(tagRequest)
			})
			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), tagRequest.GetActionName(), AlibabaCloudSdkGoERROR)
			}
			addDebug(tagRequest.GetActionName(), raw, tagRequest.RpcRequest, tagRequest)
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
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
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
		return nil, WrapErrorf(err, DefaultErrorMsg, scalingGroupId, listTagsRequest.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(listTagsRequest.GetActionName(), raw, listTagsRequest.RpcRequest, listTagsRequest)

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

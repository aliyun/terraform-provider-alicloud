package alicloud

import (
	"fmt"
	"reflect"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ess"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

type EssService struct {
	client *connectivity.AliyunClient
}

func (s *EssService) DescribeEssAlarmById(alarmTaskId string) (alarm ess.Alarm, err error) {
	systemAlarms, err := s.DescribeSystemEssAlarmById(alarmTaskId)
	if err != nil {
		return
	}
	customAlarms, err := s.DescribeCustomEssAlarmById(alarmTaskId)
	systemAlarmsEmpty := systemAlarms == nil || len(systemAlarms) == 0
	customAlarmsEmpty := customAlarms == nil || len(customAlarms) == 0
	if systemAlarmsEmpty && customAlarmsEmpty {
		err = GetNotFoundErrorFromString(GetNotFoundMessage("Ess alarm", alarmTaskId))
		return
	}
	if systemAlarmsEmpty {
		return customAlarms[0], nil
	} else {
		return systemAlarms[0], nil
	}
}

func (s *EssService) DescribeSystemEssAlarmById(alarmTaskId string) (alarm []ess.Alarm, err error) {
	args := ess.CreateDescribeAlarmsRequest()
	args.AlarmTaskId = alarmTaskId
	args.MetricType = string(System)
	raw, err := s.client.WithEssClient(func(essClient *ess.Client) (interface{}, error) {
		return essClient.DescribeAlarms(args)
	})
	if err != nil {
		return
	}
	resp, _ := raw.(*ess.DescribeAlarmsResponse)
	return resp.AlarmList.Alarm, nil
}

func (s *EssService) DescribeCustomEssAlarmById(alarmTaskId string) (alarm []ess.Alarm, err error) {
	args := ess.CreateDescribeAlarmsRequest()
	args.AlarmTaskId = alarmTaskId
	args.MetricType = string(Custom)
	raw, err := s.client.WithEssClient(func(essClient *ess.Client) (interface{}, error) {
		return essClient.DescribeAlarms(args)
	})
	if err != nil {
		return
	}
	resp, _ := raw.(*ess.DescribeAlarmsResponse)
	return resp.AlarmList.Alarm, nil
}

func (s *EssService) DescribeEssLifecycleHook(id string) (hook ess.LifecycleHook, err error) {
	request := ess.CreateDescribeLifecycleHooksRequest()
	request.LifecycleHookId = &[]string{id}

	raw, err := s.client.WithEssClient(func(essClient *ess.Client) (interface{}, error) {
		return essClient.DescribeLifecycleHooks(request)
	})
	if err != nil {
		err = WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
		return
	}
	addDebug(request.GetActionName(), raw)
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

func (s *EssService) DescribeEssScalingGroup(id string) (group ess.ScalingGroup, err error) {
	request := ess.CreateDescribeScalingGroupsRequest()
	request.ScalingGroupId1 = id

	raw, e := s.client.WithEssClient(func(essClient *ess.Client) (interface{}, error) {
		return essClient.DescribeScalingGroups(request)
	})
	if e != nil {
		err = WrapErrorf(e, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
		return
	}
	addDebug(request.GetActionName(), raw)
	response, _ := raw.(*ess.DescribeScalingGroupsResponse)
	for _, v := range response.ScalingGroups.ScalingGroup {
		if v.ScalingGroupId == id {
			return v, nil
		}
	}
	err = WrapErrorf(Error(GetNotFoundMessage("EssScalingGroup", id)), NotFoundMsg, ProviderERROR)
	return
}

func (s *EssService) DescribeScalingConfigurationById(configId string) (config ess.ScalingConfiguration, err error) {
	args := ess.CreateDescribeScalingConfigurationsRequest()
	args.ScalingConfigurationId1 = configId

	raw, err := s.client.WithEssClient(func(essClient *ess.Client) (interface{}, error) {
		return essClient.DescribeScalingConfigurations(args)
	})
	if err != nil {
		return
	}
	resp, _ := raw.(*ess.DescribeScalingConfigurationsResponse)
	if resp == nil || len(resp.ScalingConfigurations.ScalingConfiguration) < 1 {
		err = GetNotFoundErrorFromString(GetNotFoundMessage("Scaling Configuration", configId))
		return
	}

	return resp.ScalingConfigurations.ScalingConfiguration[0], nil
}

func (s *EssService) ActiveScalingConfigurationById(sgId, configId string) error {
	args := ess.CreateModifyScalingGroupRequest()
	args.ScalingGroupId = sgId
	args.ActiveScalingConfigurationId = configId

	_, err := s.client.WithEssClient(func(essClient *ess.Client) (interface{}, error) {
		return essClient.ModifyScalingGroup(args)
	})
	return err
}

// Flattens an array of datadisk into a []map[string]interface{}
func (s *EssService) flattenDataDiskMappings(list []ess.DataDisk) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(list))
	for _, i := range list {
		l := map[string]interface{}{
			"size":                 i.Size,
			"category":             i.Category,
			"snapshot_id":          i.SnapshotId,
			"device":               i.Device,
			"delete_with_instance": i.DeleteWithInstance,
		}
		result = append(result, l)
	}
	return result
}

func (s *EssService) DescribeEssScalingRule(id string) (rule ess.ScalingRule, err error) {
	request := ess.CreateDescribeScalingRulesRequest()
	request.ScalingRuleId1 = id

	raw, err := s.client.WithEssClient(func(essClient *ess.Client) (interface{}, error) {
		return essClient.DescribeScalingRules(request)
	})
	if err != nil {
		if IsExceptedErrors(err, []string{InvalidScalingRuleIdNotFound}) {
			return rule, WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
		}
		return rule, WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
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
			} else {
				return WrapError(err)
			}
		}

		if object.ScalingRuleId == id && status != Deleted {
			return nil
		}
		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, object.ScalingRuleId, id, ProviderERROR)
		}
	}
}

func (s *EssService) DescribeEssScheduledTask(id string) (task ess.ScheduledTask, err error) {
	request := ess.CreateDescribeScheduledTasksRequest()
	request.ScheduledTaskId1 = id

	raw, err := s.client.WithEssClient(func(essClient *ess.Client) (interface{}, error) {
		return essClient.DescribeScheduledTasks(request)
	})
	if err != nil {
		return task, WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw)
	response, _ := raw.(*ess.DescribeScheduledTasksResponse)

	for _, v := range response.ScheduledTasks.ScheduledTask {
		if v.ScheduledTaskId == id {
			task = v
			return
		}
	}
	err = WrapErrorf(Error(GetNotFoundMessage("EssSchedule", id)), NotFoundMsg, ProviderERROR)
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

		if object.TaskEnabled {
			return nil
		}

		time.Sleep(DefaultIntervalShort * time.Second)
		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, object.ScheduledTaskId, id, ProviderERROR)
		}
	}
}

func (s *EssService) DeleteScalingGroupById(sgId string) error {
	request := ess.CreateDeleteScalingGroupRequest()
	request.ScalingGroupId = sgId
	request.ForceDelete = requests.NewBoolean(true)
	return resource.Retry(10*time.Minute, func() *resource.RetryError {

		response, err := s.client.WithEssClient(func(essClient *ess.Client) (interface{}, error) {
			return essClient.DeleteScalingGroup(request)
		})

		if err != nil {
			if IsExceptedErrors(err, []string{InvalidScalingGroupIdNotFound}) {
				return nil
			}
			return resource.NonRetryableError(WrapError(err))
		}
		addDebug(request.GetActionName(), response)
		_, err = s.DescribeEssScalingGroup(sgId)
		if err != nil {
			if NotFoundError(err) {
				return nil
			}
			return resource.NonRetryableError(WrapError(err))
		}

		return resource.RetryableError(WrapErrorf(err, DeleteTimeoutMsg, sgId, request.GetActionName(), ProviderERROR))
	})
}

func (srv *EssService) DescribeScalingInstances(groupId, configurationId string, instanceIds []string, creationType string) (instances []ess.ScalingInstance, err error) {
	req := ess.CreateDescribeScalingInstancesRequest()

	req.ScalingGroupId = groupId
	req.ScalingConfigurationId = configurationId
	s := reflect.ValueOf(req).Elem()

	if len(instanceIds) > 0 {
		for i, id := range instanceIds {
			s.FieldByName(fmt.Sprintf("InstanceId%d", i+1)).Set(reflect.ValueOf(id))
		}
	}

	raw, err := srv.client.WithEssClient(func(essClient *ess.Client) (interface{}, error) {
		return essClient.DescribeScalingInstances(req)
	})
	if err != nil {
		return
	}
	resp, _ := raw.(*ess.DescribeScalingInstancesResponse)
	if resp == nil || len(resp.ScalingInstances.ScalingInstance) < 1 {
		return instances, GetNotFoundErrorFromString(fmt.Sprintf("There is no any instances in the scaling group %s.", groupId))
	}

	return resp.ScalingInstances.ScalingInstance, nil
}

func (s *EssService) DescribeScalingConfifurations(groupId string) (configs []ess.ScalingConfiguration, err error) {
	req := ess.CreateDescribeScalingConfigurationsRequest()
	req.ScalingGroupId = groupId
	req.PageNumber = requests.NewInteger(1)
	req.PageSize = requests.NewInteger(PageSizeLarge)

	for {
		raw, err := s.client.WithEssClient(func(essClient *ess.Client) (interface{}, error) {
			return essClient.DescribeScalingConfigurations(req)
		})
		if err != nil {
			return configs, err
		}
		resp, _ := raw.(*ess.DescribeScalingConfigurationsResponse)
		if resp == nil || len(resp.ScalingConfigurations.ScalingConfiguration) < 1 {
			break
		}
		configs = append(configs, resp.ScalingConfigurations.ScalingConfiguration...)
		if len(resp.ScalingConfigurations.ScalingConfiguration) < PageSizeLarge {
			break
		}

		if page, err := getNextpageNumber(req.PageNumber); err != nil {
			return configs, err
		} else {
			req.PageNumber = page
		}
	}

	if len(configs) < 1 {
		return configs, GetNotFoundErrorFromString(fmt.Sprintf("There is no any scaling confifurations in the scaling group %s.", groupId))
	}

	return
}

func (srv *EssService) EssRemoveInstances(groupId string, instanceIds []string) error {

	if len(instanceIds) < 1 {
		return nil
	}
	group, err := srv.DescribeEssScalingGroup(groupId)

	if err != nil {
		return WrapError(err)
	}

	if group.LifecycleState == string(Inactive) {
		return fmt.Errorf("Scaling group current status is %s, please active it before attaching or removing ECS instances.", group.LifecycleState)
	} else {
		if err := srv.WaitForEssScalingGroup(group.ScalingGroupId, Active, DefaultTimeout); err != nil {
			if NotFoundError(err) {
				return nil
			}
			return fmt.Errorf("WaitForScalingGroup is %#v got an error: %#v.", Active, err)
		}
	}

	removed := instanceIds
	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		req := ess.CreateRemoveInstancesRequest()
		req.ScalingGroupId = groupId

		if len(removed) > 0 {
			req.InstanceId = &removed
		} else {
			return nil
		}
		_, err := srv.client.WithEssClient(func(essClient *ess.Client) (interface{}, error) {
			return essClient.RemoveInstances(req)
		})
		if err != nil {
			if IsExceptedError(err, IncorrectCapacityMinSize) {
				instances, err := srv.DescribeScalingInstances(groupId, "", instanceIds, "")
				if len(instances) > 0 {
					if group.MinSize == 0 {
						return resource.RetryableError(fmt.Errorf("Removing instances got an error: %#v", err))
					}
					return resource.NonRetryableError(fmt.Errorf("To remove %d instances, the total capacity will be lesser than the scaling group min size %d. "+
						"Please shorten scaling group min size and try again.", len(instanceIds), group.MinSize))
				}
			}
			if IsExceptedError(err, ScalingActivityInProgress) || IsExceptedError(err, IncorrectScalingGroupStatus) {
				time.Sleep(5)
				return resource.RetryableError(fmt.Errorf("Removing instances got an error: %#v", err))
			}
			if IsExceptedError(err, InvalidScalingGroupIdNotFound) {
				return nil
			}
			return resource.NonRetryableError(fmt.Errorf("Removing instances got an error: %#v", err))
		}
		time.Sleep(3 * time.Second)
		instances, err := srv.DescribeScalingInstances(groupId, "", instanceIds, "")
		if err != nil {
			if NotFoundError(err) || IsExceptedErrors(err, []string{InvalidScalingGroupIdNotFound}) {
				return nil
			}
			return resource.NonRetryableError(err)
		}
		if len(instances) > 0 {
			removed = make([]string, 0)
			for _, inst := range instances {
				removed = append(removed, inst.InstanceId)
			}
			return resource.RetryableError(fmt.Errorf("There are still ECS instances in the scaling group."))
		}

		return nil
	})
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

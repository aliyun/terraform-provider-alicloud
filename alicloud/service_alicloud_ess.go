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
	args := ess.CreateDescribeAlarmsRequest()
	args.AlarmTaskId = alarmTaskId

	raw, err := s.client.WithEssClient(func(essClient *ess.Client) (interface{}, error) {
		return essClient.DescribeAlarms(args)
	})
	if err != nil {
		return
	}

	resp, _ := raw.(*ess.DescribeAlarmsResponse)
	if resp == nil || len(resp.AlarmList.Alarm) == 0 {
		err = GetNotFoundErrorFromString(GetNotFoundMessage("Ess alarm", alarmTaskId))
		return
	}
	return resp.AlarmList.Alarm[0], nil
}

func (s *EssService) DescribeLifecycleHookById(hookId string) (hook ess.LifecycleHook, err error) {
	args := ess.CreateDescribeLifecycleHooksRequest()
	hookIds := []string{hookId}
	var hookIdsPtr *[]string
	hookIdsPtr = &hookIds
	args.LifecycleHookId = hookIdsPtr

	raw, err := s.client.WithEssClient(func(essClient *ess.Client) (interface{}, error) {
		return essClient.DescribeLifecycleHooks(args)
	})
	if err != nil {
		return
	}
	resp, _ := raw.(*ess.DescribeLifecycleHooksResponse)
	if resp == nil || len(resp.LifecycleHooks.LifecycleHook) == 0 {
		err = GetNotFoundErrorFromString(GetNotFoundMessage("Lifecycle Hook", hookId))
		return
	}

	return resp.LifecycleHooks.LifecycleHook[0], nil
}

func (s *EssService) DescribeScalingGroup(sgId string) (group ess.ScalingGroup, err error) {
	request := ess.CreateDescribeScalingGroupsRequest()
	request.ScalingGroupId1 = sgId

	raw, e := s.client.WithEssClient(func(essClient *ess.Client) (interface{}, error) {
		return essClient.DescribeScalingGroups(request)
	})
	if e != nil {
		err = WrapErrorf(e, sgId, request.GetActionName(), AlibabaCloudSdkGoERROR)
		return
	}
	response, _ := raw.(*ess.DescribeScalingGroupsResponse)
	addDebug(request.GetActionName(), response)
	if response == nil || len(response.ScalingGroups.ScalingGroup) == 0 {
		err = WrapErrorf(Error(GetNotFoundMessage("Scaling Group", sgId)), NotFoundMsg, ProviderERROR)
		return
	}

	return response.ScalingGroups.ScalingGroup[0], nil
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

func (s *EssService) DescribeScalingRuleById(sgId, ruleId string) (rule ess.ScalingRule, err error) {
	args := ess.CreateDescribeScalingRulesRequest()
	args.ScalingGroupId = sgId
	args.ScalingRuleId1 = ruleId

	raw, err := s.client.WithEssClient(func(essClient *ess.Client) (interface{}, error) {
		return essClient.DescribeScalingRules(args)
	})
	if err != nil {
		if IsExceptedErrors(err, []string{InvalidScalingRuleIdNotFound}) {
			err = GetNotFoundErrorFromString(GetNotFoundMessage("Scaling rule", ruleId))
		}
		return
	}
	resp, _ := raw.(*ess.DescribeScalingRulesResponse)
	if resp == nil || len(resp.ScalingRules.ScalingRule) < 1 {
		err = GetNotFoundErrorFromString(GetNotFoundMessage("Scaling rule", ruleId))
		return
	}

	return resp.ScalingRules.ScalingRule[0], nil
}

func (s *EssService) DeleteScalingRuleById(ruleId string) error {
	args := ess.CreateDeleteScalingRuleRequest()
	args.ScalingRuleId = ruleId

	_, err := s.client.WithEssClient(func(essClient *ess.Client) (interface{}, error) {
		return essClient.DeleteScalingRule(args)
	})
	return err
}

func (s *EssService) DescribeScheduleById(scheduleId string) (task ess.ScheduledTask, err error) {
	args := ess.CreateDescribeScheduledTasksRequest()
	args.ScheduledTaskId1 = scheduleId

	raw, err := s.client.WithEssClient(func(essClient *ess.Client) (interface{}, error) {
		return essClient.DescribeScheduledTasks(args)
	})
	if err != nil {
		return
	}
	resp, _ := raw.(*ess.DescribeScheduledTasksResponse)
	if resp == nil || len(resp.ScheduledTasks.ScheduledTask) < 1 {
		err = GetNotFoundErrorFromString(GetNotFoundMessage("Schedule task", scheduleId))
		return
	}

	return resp.ScheduledTasks.ScheduledTask[0], nil
}

func (s *EssService) DeleteScheduleById(scheduleId string) error {
	args := ess.CreateDeleteScheduledTaskRequest()
	args.ScheduledTaskId = scheduleId

	_, err := s.client.WithEssClient(func(essClient *ess.Client) (interface{}, error) {
		return essClient.DeleteScheduledTask(args)
	})
	return err
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
		_, err = s.DescribeScalingGroup(sgId)
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
	group, err := srv.DescribeScalingGroup(groupId)

	if err != nil {
		return WrapError(err)
	}

	if group.LifecycleState == string(Inactive) {
		return fmt.Errorf("Scaling group current status is %s, please active it before attaching or removing ECS instances.", group.LifecycleState)
	} else {
		if err := srv.WaitForScalingGroup(group.ScalingGroupId, Active, DefaultTimeout); err != nil {
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
		s := reflect.ValueOf(req).Elem()

		if len(removed) > 0 {
			for i, id := range removed {
				s.FieldByName(fmt.Sprintf("InstanceId%d", i+1)).Set(reflect.ValueOf(id))
			}
		} else {
			return nil
		}
		_, err := srv.client.WithEssClient(func(essClient *ess.Client) (interface{}, error) {
			return essClient.RemoveInstances(req)
		})
		if err != nil {
			if IsExceptedError(err, IncorrectCapacityMinSize) {
				if group.MinSize == 0 {
					return resource.RetryableError(fmt.Errorf("Removing instances got an error: %#v", err))
				}
				return resource.NonRetryableError(fmt.Errorf("To remove %d instances, the total capacity will be lesser than the scaling group min size %d. "+
					"Please shorten scaling group min size and try again.", len(instanceIds), group.MinSize))
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
func (s *EssService) WaitForScalingGroup(groupId string, status Status, timeout int) error {
	if timeout <= 0 {
		timeout = DefaultTimeout
	}
	for {
		sg, err := s.DescribeScalingGroup(groupId)
		if err != nil {
			return WrapError(err)
		}

		if sg.LifecycleState == string(status) {
			break
		}

		timeout = timeout - DefaultIntervalShort
		if timeout <= 0 {
			return WrapError(Error(GetTimeoutMessage("Scaling Group", string(status))))
		}

		time.Sleep(DefaultIntervalShort * time.Second)

	}
	return nil
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

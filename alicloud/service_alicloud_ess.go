package alicloud

import (
	"fmt"
	"time"

	"reflect"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ess"
	"github.com/hashicorp/terraform/helper/resource"
)

func (client *AliyunClient) DescribeLifecycleHookById(hookId string) (hook ess.LifecycleHook, err error) {
	args := ess.CreateDescribeLifecycleHooksRequest()
	hookIds := []string{hookId}
	var hookIdsPtr *[]string
	hookIdsPtr = &hookIds
	args.LifecycleHookId = hookIdsPtr

	resp, err := client.essconn.DescribeLifecycleHooks(args)
	if err != nil {
		return
	}

	if resp == nil || len(resp.LifecycleHooks.LifecycleHook) == 0 {
		err = GetNotFoundErrorFromString(GetNotFoundMessage("Lifecycle Hook", hookId))
		return
	}

	return resp.LifecycleHooks.LifecycleHook[0], nil
}

func (client *AliyunClient) DescribeScalingGroupById(sgId string) (group ess.ScalingGroup, err error) {
	args := ess.CreateDescribeScalingGroupsRequest()
	args.ScalingGroupId1 = sgId

	resp, err := client.essconn.DescribeScalingGroups(args)
	if err != nil {
		return
	}

	if resp == nil || len(resp.ScalingGroups.ScalingGroup) == 0 {
		err = GetNotFoundErrorFromString(GetNotFoundMessage("Scaling Group", sgId))
		return
	}

	return resp.ScalingGroups.ScalingGroup[0], nil
}

func (client *AliyunClient) DescribeScalingConfigurationById(configId string) (config ess.ScalingConfiguration, err error) {
	args := ess.CreateDescribeScalingConfigurationsRequest()
	args.ScalingConfigurationId1 = configId

	resp, err := client.essconn.DescribeScalingConfigurations(args)
	if err != nil {
		return
	}

	if resp == nil || len(resp.ScalingConfigurations.ScalingConfiguration) < 1 {
		err = GetNotFoundErrorFromString(GetNotFoundMessage("Scaling Configuration", configId))
		return
	}

	return resp.ScalingConfigurations.ScalingConfiguration[0], nil
}

func (client *AliyunClient) ActiveScalingConfigurationById(sgId, configId string) error {
	args := ess.CreateModifyScalingGroupRequest()
	args.ScalingGroupId = sgId
	args.ActiveScalingConfigurationId = configId

	_, err := client.essconn.ModifyScalingGroup(args)
	return err
}

// Flattens an array of datadisk into a []map[string]interface{}
func flattenDataDiskMappings(list []ess.DataDisk) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(list))
	for _, i := range list {
		l := map[string]interface{}{
			"size":        i.Size,
			"category":    i.Category,
			"snapshot_id": i.SnapshotId,
			"device":      i.Device,
		}
		result = append(result, l)
	}
	return result
}

func (client *AliyunClient) DescribeScalingRuleById(sgId, ruleId string) (rule ess.ScalingRule, err error) {
	args := ess.CreateDescribeScalingRulesRequest()
	args.ScalingGroupId = sgId
	args.ScalingRuleId1 = ruleId

	resp, err := client.essconn.DescribeScalingRules(args)
	if err != nil {
		if IsExceptedErrors(err, []string{InvalidScalingRuleIdNotFound}) {
			err = GetNotFoundErrorFromString(GetNotFoundMessage("Scaling rule", ruleId))
		}
		return
	}

	if resp == nil || len(resp.ScalingRules.ScalingRule) < 1 {
		err = GetNotFoundErrorFromString(GetNotFoundMessage("Scaling rule", ruleId))
		return
	}

	return resp.ScalingRules.ScalingRule[0], nil
}

func (client *AliyunClient) DeleteScalingRuleById(ruleId string) error {
	args := ess.CreateDeleteScalingRuleRequest()
	args.ScalingRuleId = ruleId

	_, err := client.essconn.DeleteScalingRule(args)
	return err
}

func (client *AliyunClient) DescribeScheduleById(scheduleId string) (task ess.ScheduledTask, err error) {
	args := ess.CreateDescribeScheduledTasksRequest()
	args.ScheduledTaskId1 = scheduleId

	resp, err := client.essconn.DescribeScheduledTasks(args)
	if err != nil {
		return
	}

	if resp == nil || len(resp.ScheduledTasks.ScheduledTask) < 1 {
		err = GetNotFoundErrorFromString(GetNotFoundMessage("Schedule task", scheduleId))
		return
	}

	return resp.ScheduledTasks.ScheduledTask[0], nil
}

func (client *AliyunClient) DeleteScheduleById(scheduleId string) error {
	args := ess.CreateDeleteScheduledTaskRequest()
	args.ScheduledTaskId = scheduleId

	_, err := client.essconn.DeleteScheduledTask(args)
	return err
}

func (client *AliyunClient) DeleteScalingGroupById(sgId string) error {
	req := ess.CreateDeleteScalingGroupRequest()
	req.ScalingGroupId = sgId
	req.ForceDelete = requests.NewBoolean(true)
	return resource.Retry(5*time.Minute, func() *resource.RetryError {

		_, err := client.essconn.DeleteScalingGroup(req)

		if err != nil {
			if !IsExceptedErrors(err, []string{InvalidScalingGroupIdNotFound}) {
				return resource.RetryableError(fmt.Errorf("Delete scaling group timeout and got an error:%#v.", err))
			}
		}

		_, err = client.DescribeScalingGroupById(sgId)
		if err != nil {
			if NotFoundError(err) {
				return nil
			}
			return resource.NonRetryableError(err)
		}

		return resource.RetryableError(fmt.Errorf("Delete scaling group timeout and got an error:%#v.", err))
	})
}

func (client *AliyunClient) DescribeScalingInstances(groupId, configurationId string, instanceIds []string, creationType string) (instances []ess.ScalingInstance, err error) {
	req := ess.CreateDescribeScalingInstancesRequest()

	req.ScalingGroupId = groupId
	req.ScalingConfigurationId = configurationId
	s := reflect.ValueOf(req).Elem()

	if len(instanceIds) > 0 {
		for i, id := range instanceIds {
			s.FieldByName(fmt.Sprintf("InstanceId%d", i+1)).Set(reflect.ValueOf(id))
		}
	}

	resp, err := client.essconn.DescribeScalingInstances(req)
	if err != nil {
		return
	}
	if resp == nil || len(resp.ScalingInstances.ScalingInstance) < 1 {
		return instances, GetNotFoundErrorFromString(fmt.Sprintf("There is no any instances in the scaling group %s.", groupId))
	}

	return resp.ScalingInstances.ScalingInstance, nil
}

func (client *AliyunClient) DescribeScalingConfifurations(groupId string) (configs []ess.ScalingConfiguration, err error) {
	req := ess.CreateDescribeScalingConfigurationsRequest()
	req.ScalingGroupId = groupId
	req.PageNumber = requests.NewInteger(1)
	req.PageSize = requests.NewInteger(PageSizeLarge)

	for {
		resp, err := client.essconn.DescribeScalingConfigurations(req)
		if err != nil {
			return configs, err
		}
		if resp == nil || len(resp.ScalingConfigurations.ScalingConfiguration) < 1 {
			break
		}
		configs = append(configs, resp.ScalingConfigurations.ScalingConfiguration...)
		if len(resp.ScalingConfigurations.ScalingConfiguration) < PageSizeLarge {
			break
		}
		req.PageNumber = req.PageNumber + requests.NewInteger(1)
	}

	if len(configs) < 1 {
		return configs, GetNotFoundErrorFromString(fmt.Sprintf("There is no any scaling confifurations in the scaling group %s.", groupId))
	}

	return
}

func (client *AliyunClient) EssRemoveInstances(groupId string, instanceIds []string) error {

	if len(instanceIds) < 1 {
		return nil
	}
	group, err := client.DescribeScalingGroupById(groupId)

	if err != nil {
		return fmt.Errorf("DescribeScalingGroupById %s error: %#v", groupId, err)
	}

	if group.LifecycleState == string(Inactive) {
		return fmt.Errorf("Scaling group current status is %s, please active it before attaching or removing ECS instances.", group.LifecycleState)
	} else {
		if err := client.WaitForScalingGroup(group.ScalingGroupId, Active, DefaultTimeout); err != nil {
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
		if _, err := client.essconn.RemoveInstances(req); err != nil {
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

		instances, err := client.DescribeScalingInstances(groupId, "", instanceIds, "")
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
func (client *AliyunClient) WaitForScalingGroup(groupId string, status Status, timeout int) error {
	if timeout <= 0 {
		timeout = DefaultTimeout
	}
	for {
		sg, err := client.DescribeScalingGroupById(groupId)
		if err != nil {
			return err
		}

		if sg.LifecycleState == string(status) {
			break
		}

		timeout = timeout - DefaultIntervalShort
		if timeout <= 0 {
			return GetTimeErrorFromString(GetTimeoutMessage("Scaling Group", string(status)))
		}

		time.Sleep(DefaultIntervalShort * time.Second)

	}
	return nil
}

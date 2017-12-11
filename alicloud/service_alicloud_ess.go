package alicloud

import (
	"fmt"
	"github.com/denverdino/aliyungo/common"
	"github.com/denverdino/aliyungo/ess"
	"github.com/hashicorp/terraform/helper/resource"
	"time"
)

func (client *AliyunClient) DescribeScalingGroupById(sgId string) (*ess.ScalingGroupItemType, error) {
	args := ess.DescribeScalingGroupsArgs{
		RegionId:       client.Region,
		ScalingGroupId: []string{sgId},
	}

	sgs, _, err := client.essconn.DescribeScalingGroups(&args)
	if err != nil {
		return nil, err
	}

	if len(sgs) == 0 {
		return nil, GetNotFoundErrorFromString("Scaling group not found")
	}

	return &sgs[0], nil
}

func (client *AliyunClient) DescribeScalingConfigurationById(configId string) (*ess.ScalingConfigurationItemType, error) {
	args := ess.DescribeScalingConfigurationsArgs{
		RegionId: client.Region,
		//ScalingGroupId:         sgId,
		ScalingConfigurationId: []string{configId},
	}

	cs, _, err := client.essconn.DescribeScalingConfigurations(&args)
	if err != nil {
		return nil, err
	}

	if len(cs) == 0 {
		return nil, GetNotFoundErrorFromString("Scaling configuration not found")
	}

	return &cs[0], nil
}

func (client *AliyunClient) ActiveScalingConfigurationById(sgId, configId string) error {
	args := ess.ModifyScalingGroupArgs{
		ScalingGroupId:               sgId,
		ActiveScalingConfigurationId: configId,
	}

	_, err := client.essconn.ModifyScalingGroup(&args)
	return err
}

// Flattens an array of datadisk into a []map[string]interface{}
func flattenDataDiskMappings(list []ess.DataDiskItemType) []map[string]interface{} {
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

func (client *AliyunClient) DescribeScalingRuleById(sgId, ruleId string) (*ess.ScalingRuleItemType, error) {
	args := ess.DescribeScalingRulesArgs{
		RegionId:       client.Region,
		ScalingGroupId: sgId,
		ScalingRuleId:  []string{ruleId},
	}

	cs, _, err := client.essconn.DescribeScalingRules(&args)
	if err != nil {
		return nil, err
	}

	if len(cs) == 0 {
		return nil, GetNotFoundErrorFromString("Scaling rule not found")
	}

	return &cs[0], nil
}

func (client *AliyunClient) DeleteScalingRuleById(ruleId string) error {
	args := ess.DeleteScalingRuleArgs{
		RegionId:      client.Region,
		ScalingRuleId: ruleId,
	}

	_, err := client.essconn.DeleteScalingRule(&args)
	return err
}

func (client *AliyunClient) DescribeScheduleById(scheduleId string) (*ess.ScheduledTaskItemType, error) {
	args := ess.DescribeScheduledTasksArgs{
		RegionId:        client.Region,
		ScheduledTaskId: []string{scheduleId},
	}

	cs, _, err := client.essconn.DescribeScheduledTasks(&args)
	if err != nil {
		return nil, err
	}

	if len(cs) == 0 {
		return nil, GetNotFoundErrorFromString("Schedule not found")
	}

	return &cs[0], nil
}

func (client *AliyunClient) DeleteScheduleById(scheduleId string) error {
	args := ess.DeleteScheduledTaskArgs{
		RegionId:        client.Region,
		ScheduledTaskId: scheduleId,
	}

	_, err := client.essconn.DeleteScheduledTask(&args)
	return err
}

func (client *AliyunClient) DeleteScalingGroupById(sgId string) error {
	return resource.Retry(5*time.Minute, func() *resource.RetryError {

		_, err := client.essconn.DeleteScalingGroup(&ess.DeleteScalingGroupArgs{
			ScalingGroupId: sgId,
			ForceDelete:    true,
		})

		if err != nil {
			e, _ := err.(*common.Error)
			if e.ErrorResponse.Code != InvalidScalingGroupIdNotFound {
				return resource.RetryableError(fmt.Errorf("Scaling group in use - trying again while it is deleted."))
			}
		}

		_, err = client.DescribeScalingGroupById(sgId)
		if err != nil {
			if NotFoundError(err) {
				return nil
			}
			return resource.NonRetryableError(err)
		}

		return resource.RetryableError(fmt.Errorf("Scaling group in use - trying again while it is deleted."))
	})
}

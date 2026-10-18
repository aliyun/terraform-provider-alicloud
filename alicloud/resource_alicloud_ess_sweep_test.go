package alicloud

import (
	"fmt"
	"log"
	"strings"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ess"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// This file 补齐's missing ESS sweepers for the persistent child resources that
// hang off a scaling group: scaling configurations, scaling rules, lifecycle
// hooks, notifications and alarms. The scaling-group sweeper already
// force-deletes the group (which cascades to children); these sweepers clean up
// any orphaned children that survive a failed/incomplete group delete.
//
// All sweepers reuse protectedSweepVpcName for NODELETING/do-not-delete name
// protection (same generic matcher used by VPC sweepers) so SWEEP_ALL=1 is safe.
// Child resources are swept before their parent scaling group via Dependencies.

func init() {
	// Children are registered without Dependencies; the parent scaling-group
	// sweeper (in resource_alicloud_ess_scaling_group_test.go) lists these as
	// Dependencies so children are swept BEFORE the parent group, matching the
	// "sweep child before parent" guidance.
	resource.AddTestSweepers("alicloud_ess_scaling_configuration", &resource.Sweeper{
		Name: "alicloud_ess_scaling_configuration",
		F:    testSweepEssScalingConfigurations,
	})
	resource.AddTestSweepers("alicloud_ess_scaling_rule", &resource.Sweeper{
		Name: "alicloud_ess_scaling_rule",
		F:    testSweepEssScalingRules,
	})
	resource.AddTestSweepers("alicloud_ess_lifecycle_hook", &resource.Sweeper{
		Name: "alicloud_ess_lifecycle_hook",
		F:    testSweepEssLifecycleHooks,
	})
	resource.AddTestSweepers("alicloud_ess_notification", &resource.Sweeper{
		Name: "alicloud_ess_notification",
		F:    testSweepEssNotifications,
	})
	resource.AddTestSweepers("alicloud_ess_alarm", &resource.Sweeper{
		Name: "alicloud_ess_alarm",
		F:    testSweepEssAlarms,
	})
}

var essSweepPrefixes = []string{
	"tf-testAcc",
	"tf_testAcc",
	"tf-test",
}

// essSkipByName applies the standard sweeper name gate: NODELETING/do-not-delete
// markers always skip; without SWEEP_ALL only the test prefixes are deleted.
func essSkipByName(name string) bool {
	if protectedSweepVpcName(name) {
		return true
	}
	if !sweepAll() {
		matches := false
		for _, prefix := range essSweepPrefixes {
			if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
				matches = true
				break
			}
		}
		return !matches
	}
	return false
}

// essScalingGroupNames lists every scaling group in the region and returns a
// ScalingGroupId -> ScalingGroupName map. Used by child sweepers that have no
// name of their own (e.g. notifications) to inherit the parent's protection.
func essScalingGroupNames(client *connectivity.AliyunClient) (map[string]string, error) {
	names := make(map[string]string)
	req := ess.CreateDescribeScalingGroupsRequest()
	req.RegionId = client.RegionId
	req.PageSize = requests.NewInteger(PageSizeLarge)
	req.PageNumber = requests.NewInteger(1)
	for {
		raw, err := client.WithEssClient(func(essClient *ess.Client) (interface{}, error) {
			return essClient.DescribeScalingGroups(req)
		})
		if err != nil {
			return names, fmt.Errorf("error retrieving scaling groups: %s", err)
		}
		resp, _ := raw.(*ess.DescribeScalingGroupsResponse)
		if resp == nil || len(resp.ScalingGroups.ScalingGroup) < 1 {
			break
		}
		for _, g := range resp.ScalingGroups.ScalingGroup {
			names[g.ScalingGroupId] = g.ScalingGroupName
		}
		if len(resp.ScalingGroups.ScalingGroup) < PageSizeLarge {
			break
		}
		page, err := getNextpageNumber(req.PageNumber)
		if err != nil {
			return names, err
		}
		req.PageNumber = page
	}
	return names, nil
}

func testSweepEssScalingConfigurations(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	var configs []ess.ScalingConfigurationInDescribeScalingConfigurations
	req := ess.CreateDescribeScalingConfigurationsRequest()
	req.RegionId = client.RegionId
	req.PageSize = requests.NewInteger(PageSizeLarge)
	req.PageNumber = requests.NewInteger(1)
	for {
		raw, err := client.WithEssClient(func(essClient *ess.Client) (interface{}, error) {
			return essClient.DescribeScalingConfigurations(req)
		})
		if err != nil {
			return fmt.Errorf("error retrieving scaling configurations: %s", err)
		}
		resp, _ := raw.(*ess.DescribeScalingConfigurationsResponse)
		if resp == nil || len(resp.ScalingConfigurations.ScalingConfiguration) < 1 {
			break
		}
		configs = append(configs, resp.ScalingConfigurations.ScalingConfiguration...)
		if len(resp.ScalingConfigurations.ScalingConfiguration) < PageSizeLarge {
			break
		}
		page, err := getNextpageNumber(req.PageNumber)
		if err != nil {
			return err
		}
		req.PageNumber = page
	}

	for _, v := range configs {
		name := v.ScalingConfigurationName
		id := v.ScalingConfigurationId
		if essSkipByName(name) {
			log.Printf("[INFO] Skipping Scaling Configuration: %s (%s)", name, id)
			continue
		}
		log.Printf("[INFO] Deleting Scaling Configuration: %s (%s)", name, id)
		req := ess.CreateDeleteScalingConfigurationRequest()
		req.ScalingConfigurationId = id
		req.RegionId = client.RegionId
		_, err := client.WithEssClient(func(essClient *ess.Client) (interface{}, error) {
			return essClient.DeleteScalingConfiguration(req)
		})
		if err != nil {
			log.Printf("[ERROR] Failed to delete Scaling Configuration (%s (%s)): %s", name, id, err)
		}
	}
	return nil
}

func testSweepEssScalingRules(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	var rules []ess.ScalingRule
	req := ess.CreateDescribeScalingRulesRequest()
	req.RegionId = client.RegionId
	req.PageSize = requests.NewInteger(PageSizeLarge)
	req.PageNumber = requests.NewInteger(1)
	for {
		raw, err := client.WithEssClient(func(essClient *ess.Client) (interface{}, error) {
			return essClient.DescribeScalingRules(req)
		})
		if err != nil {
			return fmt.Errorf("error retrieving scaling rules: %s", err)
		}
		resp, _ := raw.(*ess.DescribeScalingRulesResponse)
		if resp == nil || len(resp.ScalingRules.ScalingRule) < 1 {
			break
		}
		rules = append(rules, resp.ScalingRules.ScalingRule...)
		if len(resp.ScalingRules.ScalingRule) < PageSizeLarge {
			break
		}
		page, err := getNextpageNumber(req.PageNumber)
		if err != nil {
			return err
		}
		req.PageNumber = page
	}

	for _, v := range rules {
		name := v.ScalingRuleName
		id := v.ScalingRuleId
		if essSkipByName(name) {
			log.Printf("[INFO] Skipping Scaling Rule: %s (%s)", name, id)
			continue
		}
		log.Printf("[INFO] Deleting Scaling Rule: %s (%s)", name, id)
		req := ess.CreateDeleteScalingRuleRequest()
		req.ScalingRuleId = id
		req.RegionId = client.RegionId
		_, err := client.WithEssClient(func(essClient *ess.Client) (interface{}, error) {
			return essClient.DeleteScalingRule(req)
		})
		if err != nil {
			log.Printf("[ERROR] Failed to delete Scaling Rule (%s (%s)): %s", name, id, err)
		}
	}
	return nil
}

func testSweepEssLifecycleHooks(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	// DescribeLifecycleHooks rejects an empty ScalingGroupId, so iterate per
	// scaling group and inherit the parent group's NODELETING/prefix gate.
	groupNames, err := essScalingGroupNames(client)
	if err != nil {
		return err
	}

	for scalingGroupId, groupName := range groupNames {
		if essSkipByName(groupName) {
			log.Printf("[INFO] Skipping Lifecycle Hooks of Scaling Group (protected): %s (%s)", groupName, scalingGroupId)
			continue
		}
		req := ess.CreateDescribeLifecycleHooksRequest()
		req.RegionId = client.RegionId
		req.ScalingGroupId = scalingGroupId
		req.PageSize = requests.NewInteger(PageSizeLarge)
		req.PageNumber = requests.NewInteger(1)
		for {
			raw, err := client.WithEssClient(func(essClient *ess.Client) (interface{}, error) {
				return essClient.DescribeLifecycleHooks(req)
			})
			if err != nil {
				log.Printf("[ERROR] Failed to list lifecycle hooks of Scaling Group %s (%s): %s", groupName, scalingGroupId, err)
				break
			}
			resp, _ := raw.(*ess.DescribeLifecycleHooksResponse)
			if resp == nil || len(resp.LifecycleHooks.LifecycleHook) < 1 {
				break
			}
			for _, v := range resp.LifecycleHooks.LifecycleHook {
				name := v.LifecycleHookName
				id := v.LifecycleHookId
				// Double-check the hook name itself for NODELETING markers.
				if protectedSweepVpcName(name) {
					log.Printf("[INFO] Skipping Lifecycle Hook (NODELETING): %s (%s)", name, id)
					continue
				}
				log.Printf("[INFO] Deleting Lifecycle Hook: %s (%s)", name, id)
				delReq := ess.CreateDeleteLifecycleHookRequest()
				delReq.LifecycleHookId = id
				delReq.RegionId = client.RegionId
				_, err := client.WithEssClient(func(essClient *ess.Client) (interface{}, error) {
					return essClient.DeleteLifecycleHook(delReq)
				})
				if err != nil {
					log.Printf("[ERROR] Failed to delete Lifecycle Hook (%s (%s)): %s", name, id, err)
				}
			}
			if len(resp.LifecycleHooks.LifecycleHook) < PageSizeLarge {
				break
			}
			page, err := getNextpageNumber(req.PageNumber)
			if err != nil {
				break
			}
			req.PageNumber = page
		}
	}
	return nil
}

// testSweepEssNotifications lists notification configs per scaling group
// (DescribeNotificationConfigurations requires a ScalingGroupId) and inherits
// the parent scaling group's NODELETING/prefix protection so real存 notification
// configs are never swept, even under SWEEP_ALL=1.
func testSweepEssNotifications(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	groupNames, err := essScalingGroupNames(client)
	if err != nil {
		return err
	}

	action := "DescribeNotificationConfigurations"
	for scalingGroupId, groupName := range groupNames {
		// Inherit the parent scaling group's protection; a notification has no
		// name of its own, so a NODELETING/prefix-mismatched group shields its
		// notification configs.
		if essSkipByName(groupName) {
			log.Printf("[INFO] Skipping Notifications of Scaling Group (protected): %s (%s)", groupName, scalingGroupId)
			continue
		}
		request := map[string]interface{}{
			"RegionId":       client.RegionId,
			"ScalingGroupId": scalingGroupId,
			"PageSize":       PageSizeLarge,
			"PageNumber":     1,
		}
		for {
			response, err := client.RpcPost("Ess", "2014-08-28", action, nil, request, true)
			if err != nil {
				log.Printf("[ERROR] Failed to list notifications of Scaling Group %s (%s): %s", groupName, scalingGroupId, err)
				break
			}
			addDebug(action, response, request)
			raw, err := jsonpath.Get("$.NotificationConfigurationModels.NotificationConfigurationModel", response)
			if err != nil {
				break
			}
			result, _ := raw.([]interface{})
			if len(result) < 1 {
				break
			}
			for _, v := range result {
				item, ok := v.(map[string]interface{})
				if !ok {
					continue
				}
				arn := fmt.Sprint(item["NotificationArn"])
				log.Printf("[INFO] Deleting Notification %s:%s (group=%s)", scalingGroupId, arn, groupName)
				req := ess.CreateDeleteNotificationConfigurationRequest()
				req.RegionId = client.RegionId
				req.ScalingGroupId = scalingGroupId
				req.NotificationArn = arn
				_, err := client.WithEssClient(func(essClient *ess.Client) (interface{}, error) {
					return essClient.DeleteNotificationConfiguration(req)
				})
				if err != nil {
					log.Printf("[ERROR] Failed to delete Notification %s:%s: %s", scalingGroupId, arn, err)
				}
			}
			if len(result) < PageSizeLarge {
				break
			}
			request["PageNumber"] = request["PageNumber"].(int) + 1
		}
	}
	return nil
}

func testSweepEssAlarms(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	var alarms []ess.Alarm
	req := ess.CreateDescribeAlarmsRequest()
	req.RegionId = client.RegionId
	req.PageSize = requests.NewInteger(PageSizeLarge)
	req.PageNumber = requests.NewInteger(1)
	for {
		raw, err := client.WithEssClient(func(essClient *ess.Client) (interface{}, error) {
			return essClient.DescribeAlarms(req)
		})
		if err != nil {
			return fmt.Errorf("error retrieving ess alarms: %s", err)
		}
		resp, _ := raw.(*ess.DescribeAlarmsResponse)
		if resp == nil || len(resp.AlarmList.Alarm) < 1 {
			break
		}
		alarms = append(alarms, resp.AlarmList.Alarm...)
		if len(resp.AlarmList.Alarm) < PageSizeLarge {
			break
		}
		page, err := getNextpageNumber(req.PageNumber)
		if err != nil {
			return err
		}
		req.PageNumber = page
	}

	for _, v := range alarms {
		// The resource schema maps "name" -> Alarm.Name; prefer it, falling
		// back to AlarmTaskName for older tasks that set only that field.
		name := v.Name
		if name == "" {
			name = v.AlarmTaskName
		}
		id := v.AlarmTaskId
		if essSkipByName(name) {
			log.Printf("[INFO] Skipping Ess Alarm: %s (%s)", name, id)
			continue
		}
		log.Printf("[INFO] Deleting Ess Alarm: %s (%s)", name, id)
		req := ess.CreateDeleteAlarmRequest()
		req.AlarmTaskId = id
		req.RegionId = client.RegionId
		_, err := client.WithEssClient(func(essClient *ess.Client) (interface{}, error) {
			return essClient.DeleteAlarm(req)
		})
		if err != nil {
			log.Printf("[ERROR] Failed to delete Ess Alarm (%s (%s)): %s", name, id, err)
		}
	}
	return nil
}

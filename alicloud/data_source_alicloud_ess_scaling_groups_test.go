package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAliCloudEssScalingGroupsDataSource(t *testing.T) {
	rand := acctest.RandInt()
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudEssScalinggroupsDataSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_ess_scaling_group.default.scaling_group_name}"`,
		}),
		fakeConfig: testAccCheckAliCloudEssScalinggroupsDataSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_ess_scaling_group.default.scaling_group_name}_fake"`,
		}),
	}

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudEssScalinggroupsDataSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_ess_scaling_group.default.id}"]`,
		}),
		fakeConfig: testAccCheckAliCloudEssScalinggroupsDataSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_ess_scaling_group.default.id}_fake"]`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudEssScalinggroupsDataSourceConfig(rand, map[string]string{
			"ids":        `["${alicloud_ess_scaling_group.default.id}"]`,
			"name_regex": `"${alicloud_ess_scaling_group.default.scaling_group_name}"`,
		}),
		fakeConfig: testAccCheckAliCloudEssScalinggroupsDataSourceConfig(rand, map[string]string{
			"ids":        `["${alicloud_ess_scaling_group.default.id}_fake"]`,
			"name_regex": `"${alicloud_ess_scaling_group.default.scaling_group_name}"`,
		}),
	}

	var existEssScalingGroupsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"groups.#":                           "1",
			"ids.#":                              "1",
			"names.#":                            "1",
			"groups.0.id":                        CHECKSET,
			"groups.0.name":                      fmt.Sprintf("tf-testAccDataSourceEssScalingGroups-%d", rand),
			"groups.0.region_id":                 CHECKSET,
			"groups.0.min_size":                  "0",
			"groups.0.desired_capacity":          "1",
			"groups.0.max_size":                  "2",
			"groups.0.cooldown_time":             "20",
			"groups.0.scaling_policy":            "release",
			"groups.0.stop_instance_timeout":     "30",
			"groups.0.max_instance_lifetime":     "86400",
			"groups.0.removal_policies.#":        "2",
			"groups.0.removal_policies.0":        "OldestInstance",
			"groups.0.removal_policies.1":        "NewestInstance",
			"groups.0.load_balancer_ids.#":       "0",
			"groups.0.db_instance_ids.#":         "0",
			"groups.0.vswitch_ids.#":             "1",
			"groups.0.enable_desired_capacity":   CHECKSET,
			"groups.0.monitor_group_id":          CHECKSET,
			"groups.0.system_suspended":          CHECKSET,
			"groups.0.resource_group_id":         CHECKSET,
			"groups.0.group_type":                CHECKSET,
			"groups.0.total_capacity":            CHECKSET,
			"groups.0.init_capacity":             CHECKSET,
			"groups.0.pending_wait_capacity":     CHECKSET,
			"groups.0.removing_wait_capacity":    CHECKSET,
			"groups.0.protected_capacity":        CHECKSET,
			"groups.0.standby_capacity":          CHECKSET,
			"groups.0.spot_capacity":             CHECKSET,
			"groups.0.stopped_capacity":          CHECKSET,
			"groups.0.multi_az_policy":           CHECKSET,
			"groups.0.active_capacity":           CHECKSET,
			"groups.0.pending_capacity":          CHECKSET,
			"groups.0.removing_capacity":         CHECKSET,
			"groups.0.creation_time":             CHECKSET,
			"groups.0.vpc_id":                    CHECKSET,
			"groups.0.vswitch_id":                CHECKSET,
			"groups.0.health_check_type":         CHECKSET,
			"groups.0.suspended_processes.#":     "0",
			"groups.0.group_deletion_protection": CHECKSET,
			"groups.0.spot_instance_remedy":      CHECKSET,
			"groups.0.modification_time":         CHECKSET,
			"groups.0.total_instance_count":      CHECKSET,
			"groups.0.lifecycle_state":           CHECKSET,
			"groups.0.tags.key":                  "value",
		}
	}

	var fakeEssScalingGroupsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"groups.#": "0",
			"ids.#":    "0",
			"names.#":  "0",
		}
	}

	var essScalingGroupsCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_ess_scaling_groups.default",
		existMapFunc: existEssScalingGroupsMapFunc,
		fakeMapFunc:  fakeEssScalingGroupsMapFunc,
	}

	essScalingGroupsCheckInfo.dataSourceTestCheck(t, rand, nameRegexConf, idsConf, allConf)
}

func TestAccAliCloudEssScalingGroupsDataSourceSupply(t *testing.T) {
	rand := acctest.RandInt()
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudEssScalinggroupsDataSourceConfigSupply(rand, map[string]string{
			"name_regex": `"${alicloud_ess_scaling_group.default.scaling_group_name}"`,
		}),
		fakeConfig: testAccCheckAliCloudEssScalinggroupsDataSourceConfigSupply(rand, map[string]string{
			"name_regex": `"${alicloud_ess_scaling_group.default.scaling_group_name}_fake"`,
		}),
	}

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudEssScalinggroupsDataSourceConfigSupply(rand, map[string]string{
			"ids": `["${alicloud_ess_scaling_group.default.id}"]`,
		}),
		fakeConfig: testAccCheckAliCloudEssScalinggroupsDataSourceConfigSupply(rand, map[string]string{
			"ids": `["${alicloud_ess_scaling_group.default.id}_fake"]`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudEssScalinggroupsDataSourceConfigSupply(rand, map[string]string{
			"ids":        `["${alicloud_ess_scaling_group.default.id}"]`,
			"name_regex": `"${alicloud_ess_scaling_group.default.scaling_group_name}"`,
		}),
		fakeConfig: testAccCheckAliCloudEssScalinggroupsDataSourceConfig(rand, map[string]string{
			"ids":        `["${alicloud_ess_scaling_group.default.id}_fake"]`,
			"name_regex": `"${alicloud_ess_scaling_group.default.scaling_group_name}"`,
		}),
	}

	var existEssScalingGroupsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"groups.#":                                          "1",
			"ids.#":                                             "1",
			"names.#":                                           "1",
			"groups.0.id":                                       CHECKSET,
			"groups.0.name":                                     fmt.Sprintf("tf-testAccDataSourceEssScalingGroups-%d", rand),
			"groups.0.region_id":                                CHECKSET,
			"groups.0.min_size":                                 "0",
			"groups.0.desired_capacity":                         "1",
			"groups.0.max_size":                                 "2",
			"groups.0.cooldown_time":                            "20",
			"groups.0.stop_instance_timeout":                    "30",
			"groups.0.max_instance_lifetime":                    "86400",
			"groups.0.az_balance":                               "true",
			"groups.0.removal_policies.#":                       "2",
			"groups.0.removal_policies.0":                       "OldestInstance",
			"groups.0.removal_policies.1":                       "NewestInstance",
			"groups.0.load_balancer_ids.#":                      "0",
			"groups.0.db_instance_ids.#":                        "0",
			"groups.0.vswitch_ids.#":                            "1",
			"groups.0.monitor_group_id":                         CHECKSET,
			"groups.0.enable_desired_capacity":                  CHECKSET,
			"groups.0.system_suspended":                         CHECKSET,
			"groups.0.resource_group_id":                        CHECKSET,
			"groups.0.scaling_policy":                           "release",
			"groups.0.group_type":                               CHECKSET,
			"groups.0.total_capacity":                           CHECKSET,
			"groups.0.init_capacity":                            CHECKSET,
			"groups.0.pending_wait_capacity":                    CHECKSET,
			"groups.0.removing_wait_capacity":                   CHECKSET,
			"groups.0.protected_capacity":                       CHECKSET,
			"groups.0.standby_capacity":                         CHECKSET,
			"groups.0.spot_capacity":                            CHECKSET,
			"groups.0.stopped_capacity":                         CHECKSET,
			"groups.0.multi_az_policy":                          CHECKSET,
			"groups.0.spot_instance_pools":                      "2",
			"groups.0.on_demand_percentage_above_base_capacity": "100",
			"groups.0.spot_allocation_strategy":                 "lowestPrice",
			"groups.0.allocation_strategy":                      "lowestPrice",
			"groups.0.on_demand_base_capacity":                  "1",
			"groups.0.active_capacity":                          CHECKSET,
			"groups.0.pending_capacity":                         CHECKSET,
			"groups.0.removing_capacity":                        CHECKSET,
			"groups.0.creation_time":                            CHECKSET,
			"groups.0.vpc_id":                                   CHECKSET,
			"groups.0.vswitch_id":                               CHECKSET,
			"groups.0.health_check_type":                        CHECKSET,
			"groups.0.suspended_processes.#":                    "0",
			"groups.0.group_deletion_protection":                CHECKSET,
			"groups.0.spot_instance_remedy":                     CHECKSET,
			"groups.0.modification_time":                        CHECKSET,
			"groups.0.total_instance_count":                     CHECKSET,
			"groups.0.lifecycle_state":                          CHECKSET,
			"groups.0.tags.key":                                 "value",
		}
	}

	var fakeEssScalingGroupsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"groups.#": "0",
			"ids.#":    "0",
			"names.#":  "0",
		}
	}

	var essScalingGroupsCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_ess_scaling_groups.default",
		existMapFunc: existEssScalingGroupsMapFunc,
		fakeMapFunc:  fakeEssScalingGroupsMapFunc,
	}

	essScalingGroupsCheckInfo.dataSourceTestCheck(t, rand, nameRegexConf, idsConf, allConf)
}

func TestAccAliCloudEssScalingGroupsDataSourceCapacityOptions(t *testing.T) {
	rand := acctest.RandInt()
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudEssScalinggroupsDataSourceConfigCapacityOptions(rand, map[string]string{
			"name_regex": `"${alicloud_ess_scaling_group.default.scaling_group_name}"`,
		}),
		fakeConfig: testAccCheckAliCloudEssScalinggroupsDataSourceConfigCapacityOptions(rand, map[string]string{
			"name_regex": `"${alicloud_ess_scaling_group.default.scaling_group_name}_fake"`,
		}),
	}

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudEssScalinggroupsDataSourceConfigCapacityOptions(rand, map[string]string{
			"ids": `["${alicloud_ess_scaling_group.default.id}"]`,
		}),
		fakeConfig: testAccCheckAliCloudEssScalinggroupsDataSourceConfigCapacityOptions(rand, map[string]string{
			"ids": `["${alicloud_ess_scaling_group.default.id}_fake"]`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudEssScalinggroupsDataSourceConfigCapacityOptions(rand, map[string]string{
			"ids":        `["${alicloud_ess_scaling_group.default.id}"]`,
			"name_regex": `"${alicloud_ess_scaling_group.default.scaling_group_name}"`,
		}),
		fakeConfig: testAccCheckAliCloudEssScalinggroupsDataSourceConfigCapacityOptions(rand, map[string]string{
			"ids":        `["${alicloud_ess_scaling_group.default.id}_fake"]`,
			"name_regex": `"${alicloud_ess_scaling_group.default.scaling_group_name}"`,
		}),
	}

	var existEssScalingGroupsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"groups.#":                         "1",
			"ids.#":                            "1",
			"names.#":                          "1",
			"groups.0.id":                      CHECKSET,
			"groups.0.name":                    fmt.Sprintf("tf-testAccDataSourceEssScalingGroupsCapacityOptions-%d", rand),
			"groups.0.region_id":               CHECKSET,
			"groups.0.min_size":                "0",
			"groups.0.desired_capacity":        "1",
			"groups.0.max_size":                "2",
			"groups.0.cooldown_time":           "20",
			"groups.0.on_demand_base_capacity": "10",
			"groups.0.on_demand_percentage_above_base_capacity":                  "10",
			"groups.0.removal_policies.#":                                        "2",
			"groups.0.removal_policies.0":                                        "OldestInstance",
			"groups.0.removal_policies.1":                                        "NewestInstance",
			"groups.0.load_balancer_ids.#":                                       "0",
			"groups.0.db_instance_ids.#":                                         "0",
			"groups.0.vswitch_ids.#":                                             "1",
			"groups.0.monitor_group_id":                                          CHECKSET,
			"groups.0.enable_desired_capacity":                                   CHECKSET,
			"groups.0.system_suspended":                                          CHECKSET,
			"groups.0.resource_group_id":                                         CHECKSET,
			"groups.0.scaling_policy":                                            "release",
			"groups.0.group_type":                                                CHECKSET,
			"groups.0.total_capacity":                                            CHECKSET,
			"groups.0.init_capacity":                                             CHECKSET,
			"groups.0.pending_wait_capacity":                                     CHECKSET,
			"groups.0.removing_wait_capacity":                                    CHECKSET,
			"groups.0.protected_capacity":                                        CHECKSET,
			"groups.0.standby_capacity":                                          CHECKSET,
			"groups.0.spot_capacity":                                             CHECKSET,
			"groups.0.stopped_capacity":                                          CHECKSET,
			"groups.0.multi_az_policy":                                           CHECKSET,
			"groups.0.spot_instance_pools":                                       "2",
			"groups.0.active_capacity":                                           CHECKSET,
			"groups.0.pending_capacity":                                          CHECKSET,
			"groups.0.removing_capacity":                                         CHECKSET,
			"groups.0.creation_time":                                             CHECKSET,
			"groups.0.vpc_id":                                                    CHECKSET,
			"groups.0.vswitch_id":                                                CHECKSET,
			"groups.0.health_check_type":                                         CHECKSET,
			"groups.0.suspended_processes.#":                                     "0",
			"groups.0.group_deletion_protection":                                 CHECKSET,
			"groups.0.spot_instance_remedy":                                      CHECKSET,
			"groups.0.modification_time":                                         CHECKSET,
			"groups.0.total_instance_count":                                      CHECKSET,
			"groups.0.lifecycle_state":                                           CHECKSET,
			"groups.0.tags.key":                                                  "value",
			"groups.0.capacity_options_on_demand_base_capacity":                  "10",
			"groups.0.capacity_options_on_demand_percentage_above_base_capacity": "10",
			"groups.0.capacity_options_compensate_with_on_demand":                "true",
			"groups.0.capacity_options_spot_auto_replace_on_demand":              "true",
		}
	}

	var fakeEssScalingGroupsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"groups.#": "0",
			"ids.#":    "0",
			"names.#":  "0",
		}
	}

	var essScalingGroupsCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_ess_scaling_groups.default",
		existMapFunc: existEssScalingGroupsMapFunc,
		fakeMapFunc:  fakeEssScalingGroupsMapFunc,
	}

	essScalingGroupsCheckInfo.dataSourceTestCheck(t, rand, nameRegexConf, idsConf, allConf)
}

func TestAccAliCloudEssScalingGroupsDataSourceLaunchTemplateOverrides(t *testing.T) {
	rand := acctest.RandInt()
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudEssScalinggroupsDataSourceConfigLaunchTemplateOverrides(rand, map[string]string{
			"name_regex": `"${alicloud_ess_scaling_group.default.scaling_group_name}"`,
		}),
		fakeConfig: testAccCheckAliCloudEssScalinggroupsDataSourceConfigLaunchTemplateOverrides(rand, map[string]string{
			"name_regex": `"${alicloud_ess_scaling_group.default.scaling_group_name}_fake"`,
		}),
	}

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudEssScalinggroupsDataSourceConfigLaunchTemplateOverrides(rand, map[string]string{
			"ids": `["${alicloud_ess_scaling_group.default.id}"]`,
		}),
		fakeConfig: testAccCheckAliCloudEssScalinggroupsDataSourceConfigLaunchTemplateOverrides(rand, map[string]string{
			"ids": `["${alicloud_ess_scaling_group.default.id}_fake"]`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudEssScalinggroupsDataSourceConfigLaunchTemplateOverrides(rand, map[string]string{
			"ids":        `["${alicloud_ess_scaling_group.default.id}"]`,
			"name_regex": `"${alicloud_ess_scaling_group.default.scaling_group_name}"`,
		}),
		fakeConfig: testAccCheckAliCloudEssScalinggroupsDataSourceConfigLaunchTemplateOverrides(rand, map[string]string{
			"ids":        `["${alicloud_ess_scaling_group.default.id}_fake"]`,
			"name_regex": `"${alicloud_ess_scaling_group.default.scaling_group_name}"`,
		}),
	}

	var existEssScalingGroupsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"groups.#":                         "1",
			"ids.#":                            "1",
			"names.#":                          "1",
			"groups.0.id":                      CHECKSET,
			"groups.0.name":                    fmt.Sprintf("tf-testAccDataSourceEssScalingGroups-%d", rand),
			"groups.0.region_id":               CHECKSET,
			"groups.0.min_size":                "0",
			"groups.0.desired_capacity":        "1",
			"groups.0.max_size":                "2",
			"groups.0.cooldown_time":           "20",
			"groups.0.on_demand_base_capacity": "10",
			"groups.0.on_demand_percentage_above_base_capacity":     "10",
			"groups.0.removal_policies.#":                           "2",
			"groups.0.removal_policies.0":                           "OldestInstance",
			"groups.0.removal_policies.1":                           "NewestInstance",
			"groups.0.load_balancer_ids.#":                          "0",
			"groups.0.db_instance_ids.#":                            "0",
			"groups.0.vswitch_ids.#":                                "1",
			"groups.0.monitor_group_id":                             CHECKSET,
			"groups.0.enable_desired_capacity":                      CHECKSET,
			"groups.0.system_suspended":                             CHECKSET,
			"groups.0.resource_group_id":                            CHECKSET,
			"groups.0.scaling_policy":                               "release",
			"groups.0.group_type":                                   CHECKSET,
			"groups.0.total_capacity":                               CHECKSET,
			"groups.0.init_capacity":                                CHECKSET,
			"groups.0.pending_wait_capacity":                        CHECKSET,
			"groups.0.removing_wait_capacity":                       CHECKSET,
			"groups.0.protected_capacity":                           CHECKSET,
			"groups.0.standby_capacity":                             CHECKSET,
			"groups.0.spot_capacity":                                CHECKSET,
			"groups.0.stopped_capacity":                             CHECKSET,
			"groups.0.multi_az_policy":                              CHECKSET,
			"groups.0.spot_instance_pools":                          "2",
			"groups.0.active_capacity":                              CHECKSET,
			"groups.0.pending_capacity":                             CHECKSET,
			"groups.0.removing_capacity":                            CHECKSET,
			"groups.0.creation_time":                                CHECKSET,
			"groups.0.vpc_id":                                       CHECKSET,
			"groups.0.vswitch_id":                                   CHECKSET,
			"groups.0.health_check_type":                            CHECKSET,
			"groups.0.suspended_processes.#":                        "0",
			"groups.0.group_deletion_protection":                    CHECKSET,
			"groups.0.spot_instance_remedy":                         CHECKSET,
			"groups.0.modification_time":                            CHECKSET,
			"groups.0.total_instance_count":                         CHECKSET,
			"groups.0.lifecycle_state":                              CHECKSET,
			"groups.0.tags.key":                                     "value",
			"groups.0.launch_template_override.#":                   "1",
			"groups.0.launch_template_override.0.instance_type":     CHECKSET,
			"groups.0.launch_template_override.0.weighted_capacity": "4",
			"groups.0.launch_template_override.0.spot_price_limit":  "2.1",
		}
	}

	var fakeEssScalingGroupsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"groups.#": "0",
			"ids.#":    "0",
			"names.#":  "0",
		}
	}

	var essScalingGroupsCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_ess_scaling_groups.default",
		existMapFunc: existEssScalingGroupsMapFunc,
		fakeMapFunc:  fakeEssScalingGroupsMapFunc,
	}

	essScalingGroupsCheckInfo.dataSourceTestCheck(t, rand, nameRegexConf, idsConf, allConf)
}

func TestAccAliCloudEssScalingGroupsDataSourceCompensateWithOnDemand(t *testing.T) {
	rand := acctest.RandInt()
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudEssScalinggroupsDataSourceConfigCompensateWithOnDemand(rand, map[string]string{
			"name_regex": `"${alicloud_ess_scaling_group.default.scaling_group_name}"`,
		}),
		fakeConfig: testAccCheckAliCloudEssScalinggroupsDataSourceConfigCompensateWithOnDemand(rand, map[string]string{
			"name_regex": `"${alicloud_ess_scaling_group.default.scaling_group_name}_fake"`,
		}),
	}

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudEssScalinggroupsDataSourceConfigCompensateWithOnDemand(rand, map[string]string{
			"ids": `["${alicloud_ess_scaling_group.default.id}"]`,
		}),
		fakeConfig: testAccCheckAliCloudEssScalinggroupsDataSourceConfigCompensateWithOnDemand(rand, map[string]string{
			"ids": `["${alicloud_ess_scaling_group.default.id}_fake"]`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudEssScalinggroupsDataSourceConfigCompensateWithOnDemand(rand, map[string]string{
			"ids":        `["${alicloud_ess_scaling_group.default.id}"]`,
			"name_regex": `"${alicloud_ess_scaling_group.default.scaling_group_name}"`,
		}),
		fakeConfig: testAccCheckAliCloudEssScalinggroupsDataSourceConfigCompensateWithOnDemand(rand, map[string]string{
			"ids":        `["${alicloud_ess_scaling_group.default.id}_fake"]`,
			"name_regex": `"${alicloud_ess_scaling_group.default.scaling_group_name}"`,
		}),
	}

	var existEssScalingGroupsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"groups.#":                         "1",
			"ids.#":                            "1",
			"names.#":                          "1",
			"groups.0.id":                      CHECKSET,
			"groups.0.name":                    fmt.Sprintf("tf-testAccDataSourceEssScalingGroups-%d", rand),
			"groups.0.region_id":               CHECKSET,
			"groups.0.min_size":                "0",
			"groups.0.desired_capacity":        "1",
			"groups.0.max_size":                "2",
			"groups.0.cooldown_time":           "20",
			"groups.0.on_demand_base_capacity": "10",
			"groups.0.on_demand_percentage_above_base_capacity": "10",
			"groups.0.removal_policies.#":                       "2",
			"groups.0.removal_policies.0":                       "OldestInstance",
			"groups.0.removal_policies.1":                       "NewestInstance",
			"groups.0.load_balancer_ids.#":                      "0",
			"groups.0.db_instance_ids.#":                        "0",
			"groups.0.vswitch_ids.#":                            "1",
			"groups.0.monitor_group_id":                         CHECKSET,
			"groups.0.enable_desired_capacity":                  CHECKSET,
			"groups.0.system_suspended":                         CHECKSET,
			"groups.0.resource_group_id":                        CHECKSET,
			"groups.0.scaling_policy":                           "release",
			"groups.0.group_type":                               CHECKSET,
			"groups.0.total_capacity":                           CHECKSET,
			"groups.0.init_capacity":                            CHECKSET,
			"groups.0.pending_wait_capacity":                    CHECKSET,
			"groups.0.removing_wait_capacity":                   CHECKSET,
			"groups.0.protected_capacity":                       CHECKSET,
			"groups.0.standby_capacity":                         CHECKSET,
			"groups.0.spot_capacity":                            CHECKSET,
			"groups.0.stopped_capacity":                         CHECKSET,
			"groups.0.multi_az_policy":                          CHECKSET,
			"groups.0.spot_instance_pools":                      "2",
			"groups.0.active_capacity":                          CHECKSET,
			"groups.0.pending_capacity":                         CHECKSET,
			"groups.0.removing_capacity":                        CHECKSET,
			"groups.0.creation_time":                            CHECKSET,
			"groups.0.vpc_id":                                   CHECKSET,
			"groups.0.vswitch_id":                               CHECKSET,
			"groups.0.suspended_processes.#":                    "0",
			"groups.0.group_deletion_protection":                CHECKSET,
			"groups.0.spot_instance_remedy":                     CHECKSET,
			"groups.0.modification_time":                        CHECKSET,
			"groups.0.total_instance_count":                     CHECKSET,
			"groups.0.lifecycle_state":                          CHECKSET,
			"groups.0.tags.key":                                 "value",
			"groups.0.compensate_with_on_demand":                "true",
		}
	}

	var fakeEssScalingGroupsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"groups.#": "0",
			"ids.#":    "0",
			"names.#":  "0",
		}
	}

	var essScalingGroupsCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_ess_scaling_groups.default",
		existMapFunc: existEssScalingGroupsMapFunc,
		fakeMapFunc:  fakeEssScalingGroupsMapFunc,
	}

	essScalingGroupsCheckInfo.dataSourceTestCheck(t, rand, nameRegexConf, idsConf, allConf)
}

func testAccCheckAliCloudEssScalinggroupsDataSourceConfigCompensateWithOnDemand(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`
%s

variable "name" {
	default = "tf-testAccDataSourceEssScalingGroups-%d"
}

resource "alicloud_ess_scaling_group" "default" {
	min_size = 0
	max_size = 2
	desired_capacity = 1
	scaling_group_name = "${var.name}"
	default_cooldown = 20
    multi_az_policy = "COST_OPTIMIZED"
    spot_instance_remedy = "true"
    compensate_with_on_demand = "true"
	on_demand_base_capacity = 10
    spot_instance_pools = 2
    on_demand_percentage_above_base_capacity = 10
	removal_policies = ["OldestInstance", "NewestInstance"]
	vswitch_ids = ["${alicloud_vswitch.default.id}"]
    stop_instance_timeout = 30
	tags = {"key": "value"}
	
}


data "alicloud_ess_scaling_groups" "default" {
  %s
}
`, EcsInstanceCommonTestCase, rand, strings.Join(pairs, "\n  "))
	return config
}

func testAccCheckAliCloudEssScalinggroupsDataSourceConfigLaunchTemplateOverrides(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`
%s

variable "name" {
	default = "tf-testAccDataSourceEssScalingGroups-%d"
}

data "alicloud_zones" "default1" {
  available_disk_category     = "cloud_efficiency"
  available_resource_creation = "VSwitch"
}
data "alicloud_instance_types" "default1" {
  availability_zone = "${data.alicloud_zones.default1.zones.0.id}"
}
resource "alicloud_security_group" "default1" {
  		name   = var.name
  		vpc_id = "${alicloud_vpc.default.id}"
	}
data "alicloud_images" "default3" {
		name_regex  = "^centos.*_64"
  		most_recent = true
  		owners      = "system"
	}
resource "alicloud_ecs_launch_template" "default3" {
  		launch_template_name = "tf-test3"
  		image_id             =  data.alicloud_images.default3.images.0.id
  		instance_charge_type = "PrePaid"
  		instance_type        =  data.alicloud_instance_types.default1.instance_types.0.id
  		internet_charge_type          = "PayByBandwidth"
  		internet_max_bandwidth_in     = "5"
  		internet_max_bandwidth_out    = "0"
  		io_optimized                  = "optimized"
  		network_type                  = "vpc"
  		security_enhancement_strategy = "Active"
  		spot_price_limit              = "5"
  		spot_strategy                 = "SpotWithPriceLimit"
  		security_group_id             = alicloud_security_group.default1.id
	}

resource "alicloud_ess_scaling_group" "default" {
	min_size = 0
	max_size = 2
	desired_capacity = 1
	scaling_group_name = "${var.name}"
	default_cooldown = 20
    multi_az_policy = "COST_OPTIMIZED"
    spot_instance_remedy = "true"
	on_demand_base_capacity = 10
    spot_instance_pools = 2
    on_demand_percentage_above_base_capacity = 10
	removal_policies = ["OldestInstance", "NewestInstance"]
	vswitch_ids = ["${alicloud_vswitch.default.id}"]
    stop_instance_timeout = 30
	tags = {"key": "value"}
	launch_template_id = "${alicloud_ecs_launch_template.default3.id}"
    launch_template_version = "Latest"
	scaling_policy = "release"
	launch_template_override {
        instance_type = "${data.alicloud_instance_types.default.instance_types.1.id}"
        weighted_capacity = "4"
        spot_price_limit = "2.1"
	}
}



data "alicloud_ess_scaling_groups" "default" {
  %s
}
`, EcsInstanceCommonTestCase, rand, strings.Join(pairs, "\n  "))
	return config
}

func testAccCheckAliCloudEssScalinggroupsDataSourceConfigCapacityOptions(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`
%s

variable "name" {
	default = "tf-testAccDataSourceEssScalingGroupsCapacityOptions-%d"
}

resource "alicloud_ess_scaling_group" "default" {
	min_size = 0
	max_size = 2
	desired_capacity = 1
	scaling_group_name = "${var.name}"
	default_cooldown = 20
    multi_az_policy = "COST_OPTIMIZED"
    spot_instance_remedy = "true"
	on_demand_base_capacity = 10
    spot_instance_pools = 2
    on_demand_percentage_above_base_capacity = 10
	removal_policies = ["OldestInstance", "NewestInstance"]
	vswitch_ids = ["${alicloud_vswitch.default.id}"]
    stop_instance_timeout = 30
	tags = {"key": "value"}
	scaling_policy = "release"
	capacity_options_on_demand_base_capacity = 10
    capacity_options_on_demand_percentage_above_base_capacity = 10
    capacity_options_compensate_with_on_demand = "true"
    capacity_options_spot_auto_replace_on_demand = "true"
}


data "alicloud_ess_scaling_groups" "default" {
  %s
}
`, EcsInstanceCommonTestCase, rand, strings.Join(pairs, "\n  "))
	return config
}

func testAccCheckAliCloudEssScalinggroupsDataSourceConfigSupply(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`
%s

variable "name" {
	default = "tf-testAccDataSourceEssScalingGroups-%d"
}

resource "alicloud_ess_scaling_group" "default" {
	min_size = 0
	max_size = 2
	desired_capacity = 1
	scaling_group_name = "${var.name}"
	default_cooldown = 20
    multi_az_policy = "COMPOSABLE"
    on_demand_percentage_above_base_capacity = "100"
	on_demand_base_capacity = 1
	removal_policies = ["OldestInstance", "NewestInstance"]
	vswitch_ids = ["${alicloud_vswitch.default.id}"]
    stop_instance_timeout = 30
	tags = {"key": "value"}
	max_instance_lifetime = 86400
	spot_allocation_strategy = "lowestPrice"
	allocation_strategy = "lowestPrice"
	az_balance = "true"
	scaling_policy = "release"
}


data "alicloud_ess_scaling_groups" "default" {
  %s
}
`, EcsInstanceCommonTestCase, rand, strings.Join(pairs, "\n  "))
	return config
}

func testAccCheckAliCloudEssScalinggroupsDataSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`
%s

variable "name" {
	default = "tf-testAccDataSourceEssScalingGroups-%d"
}

resource "alicloud_ess_scaling_group" "default" {
	min_size = 0
	max_size = 2
	desired_capacity = 1
	scaling_group_name = "${var.name}"
	default_cooldown = 20
	removal_policies = ["OldestInstance", "NewestInstance"]
	vswitch_ids = ["${alicloud_vswitch.default.id}"]
    stop_instance_timeout = 30
	tags = {"key": "value"}
	max_instance_lifetime = 86400
	scaling_policy = "release"
}

data "alicloud_ess_scaling_groups" "default" {
  %s
}
`, EcsInstanceCommonTestCase, rand, strings.Join(pairs, "\n  "))
	return config
}

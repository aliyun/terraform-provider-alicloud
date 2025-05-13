package alicloud

import (
	"fmt"
	"strings"
	"testing"
)

func TestAccAliCloudEssScalingrulesDataSource(t *testing.T) {

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudEssScalingrulesDataSourceConfig(map[string]string{
			"ids": `["${alicloud_ess_scaling_rule.default.id}"]`,
		}),
		fakeConfig: testAccCheckAliCloudEssScalingrulesDataSourceConfig(map[string]string{
			"ids": `["${alicloud_ess_scaling_rule.default.id}_fake"]`,
		}),
	}

	scalingGroupIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudEssScalingrulesDataSourceConfig(map[string]string{
			"scaling_group_id": `"${alicloud_ess_scaling_rule.default.scaling_group_id}"`,
		}),
		fakeConfig: testAccCheckAliCloudEssScalingrulesDataSourceConfig(map[string]string{
			"scaling_group_id": `"${alicloud_ess_scaling_rule.default.scaling_group_id}_fake"`,
		}),
	}

	typeConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudEssScalingrulesDataSourceConfig(map[string]string{
			"scaling_group_id": `"${alicloud_ess_scaling_rule.default.scaling_group_id}"`,
			"type":             `"SimpleScalingRule"`,
		}),
		fakeConfig: testAccCheckAliCloudEssScalingrulesDataSourceConfig(map[string]string{
			"scaling_group_id": `"${alicloud_ess_scaling_rule.default.scaling_group_id}_fake"`,
			"type":             `"TargetTrackingScalingRule"`,
		}),
	}

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudEssScalingrulesDataSourceConfig(map[string]string{
			"name_regex": `"${alicloud_ess_scaling_rule.default.scaling_rule_name}"`,
		}),
		fakeConfig: testAccCheckAliCloudEssScalingrulesDataSourceConfig(map[string]string{
			"name_regex": `"${alicloud_ess_scaling_rule.default.scaling_rule_name}_fake"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudEssScalingrulesDataSourceConfig(map[string]string{
			"ids":              `["${alicloud_ess_scaling_rule.default.id}"]`,
			"scaling_group_id": `"${alicloud_ess_scaling_rule.default.scaling_group_id}"`,
			"type":             `"SimpleScalingRule"`,
			"name_regex":       `"${alicloud_ess_scaling_rule.default.scaling_rule_name}"`,
		}),
		fakeConfig: testAccCheckAliCloudEssScalingrulesDataSourceConfig(map[string]string{
			"ids":              `["${alicloud_ess_scaling_rule.default.id}"]`,
			"scaling_group_id": `"${alicloud_ess_scaling_rule.default.scaling_group_id}_fake"`,
			"type":             `"SimpleScalingRule"`,
			"name_regex":       `"${alicloud_ess_scaling_rule.default.scaling_rule_name}"`,
		}),
	}

	var existEssRecordsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"rules.#":                          "1",
			"ids.#":                            "1",
			"names.#":                          "1",
			"rules.0.id":                       CHECKSET,
			"rules.0.scaling_group_id":         CHECKSET,
			"rules.0.name":                     CHECKSET,
			"rules.0.type":                     CHECKSET,
			"rules.0.cooldown":                 "30",
			"rules.0.adjustment_type":          "PercentChangeInCapacity",
			"rules.0.adjustment_value":         "1",
			"rules.0.min_adjustment_magnitude": "1",
			"rules.0.scaling_rule_ari":         CHECKSET,
		}
	}

	var fakeEssRecordsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"rules.#": "0",
			"ids.#":   "0",
			"names.#": "0",
		}
	}

	var EssScalingrulesRecordsCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_ess_scaling_rules.default",
		existMapFunc: existEssRecordsMapFunc,
		fakeMapFunc:  fakeEssRecordsMapFunc,
	}

	EssScalingrulesRecordsCheckInfo.dataSourceTestCheck(t, -1, idsConf, scalingGroupIdConf, typeConf, nameRegexConf, allConf)
}

func TestAccAliCloudEssScalingrulesDataSourcePredictiveRule(t *testing.T) {

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudEssScalingrulesDataSourcePredictiveRuleConfig(map[string]string{
			"ids": `["${alicloud_ess_scaling_rule.default.id}"]`,
		}),
		fakeConfig: testAccCheckAliCloudEssScalingrulesDataSourcePredictiveRuleConfig(map[string]string{
			"ids": `["${alicloud_ess_scaling_rule.default.id}_fake"]`,
		}),
	}

	scalingGroupIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudEssScalingrulesDataSourcePredictiveRuleConfig(map[string]string{
			"scaling_group_id": `"${alicloud_ess_scaling_rule.default.scaling_group_id}"`,
		}),
		fakeConfig: testAccCheckAliCloudEssScalingrulesDataSourcePredictiveRuleConfig(map[string]string{
			"scaling_group_id": `"${alicloud_ess_scaling_rule.default.scaling_group_id}_fake"`,
		}),
	}
	typeConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudEssScalingrulesDataSourcePredictiveRuleConfig(map[string]string{
			"scaling_group_id": `"${alicloud_ess_scaling_rule.default.scaling_group_id}"`,
			"type":             `"PredictiveScalingRule"`,
		}),
		fakeConfig: testAccCheckAliCloudEssScalingrulesDataSourcePredictiveRuleConfig(map[string]string{
			"scaling_group_id": `"${alicloud_ess_scaling_rule.default.scaling_group_id}_fake"`,
			"type":             `"TargetTrackingScalingRule"`,
		}),
	}
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudEssScalingrulesDataSourcePredictiveRuleConfig(map[string]string{
			"name_regex": `"${alicloud_ess_scaling_rule.default.scaling_rule_name}"`,
		}),
		fakeConfig: testAccCheckAliCloudEssScalingrulesDataSourcePredictiveRuleConfig(map[string]string{
			"name_regex": `"${alicloud_ess_scaling_rule.default.scaling_rule_name}_fake"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudEssScalingrulesDataSourcePredictiveRuleConfig(map[string]string{
			"ids":              `["${alicloud_ess_scaling_rule.default.id}"]`,
			"scaling_group_id": `"${alicloud_ess_scaling_rule.default.scaling_group_id}"`,
			"type":             `"PredictiveScalingRule"`,
			"name_regex":       `"${alicloud_ess_scaling_rule.default.scaling_rule_name}"`,
		}),
		fakeConfig: testAccCheckAliCloudEssScalingrulesDataSourcePredictiveRuleConfig(map[string]string{
			"ids":              `["${alicloud_ess_scaling_rule.default.id}"]`,
			"scaling_group_id": `"${alicloud_ess_scaling_rule.default.scaling_group_id}_fake"`,
			"type":             `"SimpleScalingRule"`,
			"name_regex":       `"${alicloud_ess_scaling_rule.default.scaling_rule_name}"`,
		}),
	}

	var existEssRecordsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"rules.#":                             "1",
			"ids.#":                               "1",
			"names.#":                             "1",
			"rules.0.id":                          CHECKSET,
			"rules.0.scaling_group_id":            CHECKSET,
			"rules.0.name":                        CHECKSET,
			"rules.0.type":                        CHECKSET,
			"rules.0.predictive_task_buffer_time": "0",
			"rules.0.target_value":                "20.1",
			"rules.0.predictive_value_behavior":   "MaxOverridePredictiveValue",
			"rules.0.predictive_scaling_mode":     "PredictAndScale",
			"rules.0.metric_name":                 "CpuUtilization",
			"rules.0.predictive_value_buffer":     "0",
			"rules.0.initial_max_size":            "1",
			"rules.0.scaling_rule_ari":            CHECKSET,
		}
	}

	var fakeEssRecordsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"rules.#": "0",
			"ids.#":   "0",
			"names.#": "0",
		}
	}

	var EssScalingrulesRecordsCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_ess_scaling_rules.default",
		existMapFunc: existEssRecordsMapFunc,
		fakeMapFunc:  fakeEssRecordsMapFunc,
	}

	EssScalingrulesRecordsCheckInfo.dataSourceTestCheck(t, -1, idsConf, scalingGroupIdConf, typeConf, nameRegexConf, allConf)
}

func TestAccAliCloudEssScalingrulesDataSourceStepRule(t *testing.T) {

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudEssScalingrulesDataSourceStepRuleConfig(map[string]string{
			"ids": `["${alicloud_ess_scaling_rule.default.id}"]`,
		}),
		fakeConfig: testAccCheckAliCloudEssScalingrulesDataSourceStepRuleConfig(map[string]string{
			"ids": `["${alicloud_ess_scaling_rule.default.id}_fake"]`,
		}),
	}

	scalingGroupIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudEssScalingrulesDataSourceStepRuleConfig(map[string]string{
			"scaling_group_id": `"${alicloud_ess_scaling_rule.default.scaling_group_id}"`,
		}),
		fakeConfig: testAccCheckAliCloudEssScalingrulesDataSourceStepRuleConfig(map[string]string{
			"scaling_group_id": `"${alicloud_ess_scaling_rule.default.scaling_group_id}_fake"`,
		}),
	}
	typeConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudEssScalingrulesDataSourceStepRuleConfig(map[string]string{
			"scaling_group_id": `"${alicloud_ess_scaling_rule.default.scaling_group_id}"`,
			"type":             `"StepScalingRule"`,
		}),
		fakeConfig: testAccCheckAliCloudEssScalingrulesDataSourceStepRuleConfig(map[string]string{
			"scaling_group_id": `"${alicloud_ess_scaling_rule.default.scaling_group_id}_fake"`,
			"type":             `"TargetTrackingScalingRule"`,
		}),
	}
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudEssScalingrulesDataSourceStepRuleConfig(map[string]string{
			"name_regex": `"${alicloud_ess_scaling_rule.default.scaling_rule_name}"`,
		}),
		fakeConfig: testAccCheckAliCloudEssScalingrulesDataSourceStepRuleConfig(map[string]string{
			"name_regex": `"${alicloud_ess_scaling_rule.default.scaling_rule_name}_fake"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudEssScalingrulesDataSourceStepRuleConfig(map[string]string{
			"ids":              `["${alicloud_ess_scaling_rule.default.id}"]`,
			"scaling_group_id": `"${alicloud_ess_scaling_rule.default.scaling_group_id}"`,
			"type":             `"StepScalingRule"`,
			"name_regex":       `"${alicloud_ess_scaling_rule.default.scaling_rule_name}"`,
		}),
		fakeConfig: testAccCheckAliCloudEssScalingrulesDataSourceStepRuleConfig(map[string]string{
			"ids":              `["${alicloud_ess_scaling_rule.default.id}"]`,
			"scaling_group_id": `"${alicloud_ess_scaling_rule.default.scaling_group_id}_fake"`,
			"type":             `"SimpleScalingRule"`,
			"name_regex":       `"${alicloud_ess_scaling_rule.default.scaling_rule_name}"`,
		}),
	}

	var existEssRecordsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"rules.#":                           "1",
			"ids.#":                             "1",
			"names.#":                           "1",
			"rules.0.id":                        CHECKSET,
			"rules.0.scaling_group_id":          CHECKSET,
			"rules.0.name":                      CHECKSET,
			"rules.0.type":                      CHECKSET,
			"rules.0.estimated_instance_warmup": "200",
			"rules.0.adjustment_type":           "TotalCapacity",
			"rules.0.step_adjustment.#":         "2",
			"rules.0.step_adjustment.0.metric_interval_lower_bound": "10.3",
			"rules.0.step_adjustment.0.metric_interval_upper_bound": "20.1",
			"rules.0.step_adjustment.0.scaling_adjustment":          "1",
			"rules.0.step_adjustment.1.metric_interval_lower_bound": "20.1",
			"rules.0.step_adjustment.1.scaling_adjustment":          "2",
		}
	}

	var fakeEssRecordsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"rules.#": "0",
			"ids.#":   "0",
			"names.#": "0",
		}
	}

	var EssScalingrulesRecordsCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_ess_scaling_rules.default",
		existMapFunc: existEssRecordsMapFunc,
		fakeMapFunc:  fakeEssRecordsMapFunc,
	}

	EssScalingrulesRecordsCheckInfo.dataSourceTestCheck(t, -1, idsConf, scalingGroupIdConf, typeConf, nameRegexConf, allConf)
}

func TestAccAliCloudEssScalingrulesDataSourceTargetTrackingRule(t *testing.T) {

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudEssScalingrulesDataSourceTargetTrackingRuleConfig(map[string]string{
			"ids": `["${alicloud_ess_scaling_rule.default.id}"]`,
		}),
		fakeConfig: testAccCheckAliCloudEssScalingrulesDataSourceTargetTrackingRuleConfig(map[string]string{
			"ids": `["${alicloud_ess_scaling_rule.default.id}_fake"]`,
		}),
	}

	scalingGroupIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudEssScalingrulesDataSourceTargetTrackingRuleConfig(map[string]string{
			"scaling_group_id": `"${alicloud_ess_scaling_rule.default.scaling_group_id}"`,
		}),
		fakeConfig: testAccCheckAliCloudEssScalingrulesDataSourceTargetTrackingRuleConfig(map[string]string{
			"scaling_group_id": `"${alicloud_ess_scaling_rule.default.scaling_group_id}_fake"`,
		}),
	}
	typeConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudEssScalingrulesDataSourceTargetTrackingRuleConfig(map[string]string{
			"scaling_group_id": `"${alicloud_ess_scaling_rule.default.scaling_group_id}"`,
			"type":             `"TargetTrackingScalingRule"`,
		}),
		fakeConfig: testAccCheckAliCloudEssScalingrulesDataSourceTargetTrackingRuleConfig(map[string]string{
			"scaling_group_id": `"${alicloud_ess_scaling_rule.default.scaling_group_id}_fake"`,
			"type":             `"StepScalingRule"`,
		}),
	}
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudEssScalingrulesDataSourceTargetTrackingRuleConfig(map[string]string{
			"name_regex": `"${alicloud_ess_scaling_rule.default.scaling_rule_name}"`,
		}),
		fakeConfig: testAccCheckAliCloudEssScalingrulesDataSourceTargetTrackingRuleConfig(map[string]string{
			"name_regex": `"${alicloud_ess_scaling_rule.default.scaling_rule_name}_fake"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudEssScalingrulesDataSourceTargetTrackingRuleConfig(map[string]string{
			"ids":              `["${alicloud_ess_scaling_rule.default.id}"]`,
			"scaling_group_id": `"${alicloud_ess_scaling_rule.default.scaling_group_id}"`,
			"type":             `"TargetTrackingScalingRule"`,
			"name_regex":       `"${alicloud_ess_scaling_rule.default.scaling_rule_name}"`,
		}),
		fakeConfig: testAccCheckAliCloudEssScalingrulesDataSourceTargetTrackingRuleConfig(map[string]string{
			"ids":              `["${alicloud_ess_scaling_rule.default.id}"]`,
			"scaling_group_id": `"${alicloud_ess_scaling_rule.default.scaling_group_id}_fake"`,
			"type":             `"SimpleScalingRule"`,
			"name_regex":       `"${alicloud_ess_scaling_rule.default.scaling_rule_name}"`,
		}),
	}

	var existEssRecordsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"rules.#":                            "1",
			"ids.#":                              "1",
			"names.#":                            "1",
			"rules.0.id":                         CHECKSET,
			"rules.0.scaling_group_id":           CHECKSET,
			"rules.0.name":                       CHECKSET,
			"rules.0.type":                       CHECKSET,
			"rules.0.estimated_instance_warmup":  "200",
			"rules.0.metric_name":                "CpuUtilization",
			"rules.0.metric_type":                "system",
			"rules.0.scale_in_evaluation_count":  "5",
			"rules.0.scale_out_evaluation_count": "1",
			"rules.0.target_value":               "20.2",
		}
	}

	var fakeEssRecordsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"rules.#": "0",
			"ids.#":   "0",
			"names.#": "0",
		}
	}

	var EssScalingrulesRecordsCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_ess_scaling_rules.default",
		existMapFunc: existEssRecordsMapFunc,
		fakeMapFunc:  fakeEssRecordsMapFunc,
	}

	EssScalingrulesRecordsCheckInfo.dataSourceTestCheck(t, -1, idsConf, scalingGroupIdConf, typeConf, nameRegexConf, allConf)
}

func TestAccAliCloudEssScalingrulesDataSourceTargetTrackingRuleDisableScaleIn(t *testing.T) {

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudEssScalingrulesDataSourceTargetTrackingRuleDisableScaleInConfig(map[string]string{
			"ids": `["${alicloud_ess_scaling_rule.default.id}"]`,
		}),
		fakeConfig: testAccCheckAliCloudEssScalingrulesDataSourceTargetTrackingRuleDisableScaleInConfig(map[string]string{
			"ids": `["${alicloud_ess_scaling_rule.default.id}_fake"]`,
		}),
	}

	scalingGroupIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudEssScalingrulesDataSourceTargetTrackingRuleDisableScaleInConfig(map[string]string{
			"scaling_group_id": `"${alicloud_ess_scaling_rule.default.scaling_group_id}"`,
		}),
		fakeConfig: testAccCheckAliCloudEssScalingrulesDataSourceTargetTrackingRuleDisableScaleInConfig(map[string]string{
			"scaling_group_id": `"${alicloud_ess_scaling_rule.default.scaling_group_id}_fake"`,
		}),
	}
	typeConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudEssScalingrulesDataSourceTargetTrackingRuleDisableScaleInConfig(map[string]string{
			"scaling_group_id": `"${alicloud_ess_scaling_rule.default.scaling_group_id}"`,
			"type":             `"TargetTrackingScalingRule"`,
		}),
		fakeConfig: testAccCheckAliCloudEssScalingrulesDataSourceTargetTrackingRuleDisableScaleInConfig(map[string]string{
			"scaling_group_id": `"${alicloud_ess_scaling_rule.default.scaling_group_id}_fake"`,
			"type":             `"StepScalingRule"`,
		}),
	}
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudEssScalingrulesDataSourceTargetTrackingRuleDisableScaleInConfig(map[string]string{
			"name_regex": `"${alicloud_ess_scaling_rule.default.scaling_rule_name}"`,
		}),
		fakeConfig: testAccCheckAliCloudEssScalingrulesDataSourceTargetTrackingRuleDisableScaleInConfig(map[string]string{
			"name_regex": `"${alicloud_ess_scaling_rule.default.scaling_rule_name}_fake"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudEssScalingrulesDataSourceTargetTrackingRuleDisableScaleInConfig(map[string]string{
			"ids":              `["${alicloud_ess_scaling_rule.default.id}"]`,
			"scaling_group_id": `"${alicloud_ess_scaling_rule.default.scaling_group_id}"`,
			"type":             `"TargetTrackingScalingRule"`,
			"name_regex":       `"${alicloud_ess_scaling_rule.default.scaling_rule_name}"`,
		}),
		fakeConfig: testAccCheckAliCloudEssScalingrulesDataSourceTargetTrackingRuleDisableScaleInConfig(map[string]string{
			"ids":              `["${alicloud_ess_scaling_rule.default.id}"]`,
			"scaling_group_id": `"${alicloud_ess_scaling_rule.default.scaling_group_id}_fake"`,
			"type":             `"SimpleScalingRule"`,
			"name_regex":       `"${alicloud_ess_scaling_rule.default.scaling_rule_name}"`,
		}),
	}

	var existEssRecordsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"rules.#":                           "1",
			"ids.#":                             "1",
			"names.#":                           "1",
			"rules.0.id":                        CHECKSET,
			"rules.0.scaling_group_id":          CHECKSET,
			"rules.0.name":                      CHECKSET,
			"rules.0.type":                      CHECKSET,
			"rules.0.estimated_instance_warmup": "200",
			"rules.0.metric_name":               "CpuUtilization",
			"rules.0.metric_type":               "system",
			"rules.0.target_value":              "20.2",
			"rules.0.disable_scale_in":          "true",
		}
	}

	var fakeEssRecordsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"rules.#": "0",
			"ids.#":   "0",
			"names.#": "0",
		}
	}

	var EssScalingrulesRecordsCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_ess_scaling_rules.default",
		existMapFunc: existEssRecordsMapFunc,
		fakeMapFunc:  fakeEssRecordsMapFunc,
	}

	EssScalingrulesRecordsCheckInfo.dataSourceTestCheck(t, -1, idsConf, scalingGroupIdConf, typeConf, nameRegexConf, allConf)
}

func TestAccAliCloudEssScalingrulesDataSourceTargetTrackingRuleWithHybrid(t *testing.T) {

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudEssScalingrulesDataSourceTargetTrackingRuleWithHybridConfig(map[string]string{
			"ids": `["${alicloud_ess_scaling_rule.default.id}"]`,
		}),
		fakeConfig: testAccCheckAliCloudEssScalingrulesDataSourceTargetTrackingRuleWithHybridConfig(map[string]string{
			"ids": `["${alicloud_ess_scaling_rule.default.id}_fake"]`,
		}),
	}

	scalingGroupIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudEssScalingrulesDataSourceTargetTrackingRuleWithHybridConfig(map[string]string{
			"scaling_group_id": `"${alicloud_ess_scaling_rule.default.scaling_group_id}"`,
		}),
		fakeConfig: testAccCheckAliCloudEssScalingrulesDataSourceTargetTrackingRuleWithHybridConfig(map[string]string{
			"scaling_group_id": `"${alicloud_ess_scaling_rule.default.scaling_group_id}_fake"`,
		}),
	}
	typeConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudEssScalingrulesDataSourceTargetTrackingRuleWithHybridConfig(map[string]string{
			"scaling_group_id": `"${alicloud_ess_scaling_rule.default.scaling_group_id}"`,
			"type":             `"TargetTrackingScalingRule"`,
		}),
		fakeConfig: testAccCheckAliCloudEssScalingrulesDataSourceTargetTrackingRuleWithHybridConfig(map[string]string{
			"scaling_group_id": `"${alicloud_ess_scaling_rule.default.scaling_group_id}_fake"`,
			"type":             `"StepScalingRule"`,
		}),
	}
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudEssScalingrulesDataSourceTargetTrackingRuleWithHybridConfig(map[string]string{
			"name_regex": `"${alicloud_ess_scaling_rule.default.scaling_rule_name}"`,
		}),
		fakeConfig: testAccCheckAliCloudEssScalingrulesDataSourceTargetTrackingRuleWithHybridConfig(map[string]string{
			"name_regex": `"${alicloud_ess_scaling_rule.default.scaling_rule_name}_fake"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudEssScalingrulesDataSourceTargetTrackingRuleWithHybridConfig(map[string]string{
			"ids":              `["${alicloud_ess_scaling_rule.default.id}"]`,
			"scaling_group_id": `"${alicloud_ess_scaling_rule.default.scaling_group_id}"`,
			"type":             `"TargetTrackingScalingRule"`,
			"name_regex":       `"${alicloud_ess_scaling_rule.default.scaling_rule_name}"`,
		}),
		fakeConfig: testAccCheckAliCloudEssScalingrulesDataSourceTargetTrackingRuleWithHybridConfig(map[string]string{
			"ids":              `["${alicloud_ess_scaling_rule.default.id}"]`,
			"scaling_group_id": `"${alicloud_ess_scaling_rule.default.scaling_group_id}_fake"`,
			"type":             `"SimpleScalingRule"`,
			"name_regex":       `"${alicloud_ess_scaling_rule.default.scaling_rule_name}"`,
		}),
	}

	var existEssRecordsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"rules.#":                               "1",
			"ids.#":                                 "1",
			"names.#":                               "1",
			"rules.0.id":                            CHECKSET,
			"rules.0.scaling_group_id":              CHECKSET,
			"rules.0.name":                          CHECKSET,
			"rules.0.type":                          CHECKSET,
			"rules.0.estimated_instance_warmup":     "200",
			"rules.0.metric_type":                   "hybrid",
			"rules.0.target_value":                  "20.2",
			"rules.0.hybrid_monitor_namespace":      CHECKSET,
			"rules.0.hybrid_metrics.#":              "2",
			"rules.0.hybrid_metrics.0.metric_name":  "AliyunEcs_CPUUtilization",
			"rules.0.hybrid_metrics.0.id":           "a",
			"rules.0.hybrid_metrics.0.statistic":    "Average",
			"rules.0.hybrid_metrics.0.dimensions.#": "1",
			"rules.0.hybrid_metrics.1.id":           "Expression",
			"rules.0.hybrid_metrics.1.expression":   "2*a",
		}
	}

	var fakeEssRecordsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"rules.#": "0",
			"ids.#":   "0",
			"names.#": "0",
		}
	}

	var EssScalingrulesRecordsCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_ess_scaling_rules.default",
		existMapFunc: existEssRecordsMapFunc,
		fakeMapFunc:  fakeEssRecordsMapFunc,
	}

	EssScalingrulesRecordsCheckInfo.dataSourceTestCheck(t, -1, idsConf, scalingGroupIdConf, typeConf, nameRegexConf, allConf)
}

func testAccCheckAliCloudEssScalingrulesDataSourceTargetTrackingRuleWithHybridConfig(attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`

variable "name" {
	default = "tf-testAccDataSourceEssScalingTargetTrackingConfig"
}

data "alicloud_zones" "default" {
  available_disk_category     = "cloud_efficiency"
  available_resource_creation = "VSwitch"
}
data "alicloud_account" "default" {} 
data "alicloud_instance_types" "default" {
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
  cpu_core_count    = 2
  memory_size       = 4
}

data "alicloud_images" "default" {
  name_regex  = "^ubuntu_[0-9]+_[0-9]+_x64*"
  most_recent = true
  owners      = "system"
}

data "alicloud_vpcs" "default" {
	
}
data "alicloud_vswitches" "default" {
	vpc_id = data.alicloud_vpcs.default.ids.0
	zone_id      = data.alicloud_zones.default.zones.0.id
}

resource "alicloud_vswitch" "vswitch" {
  count             = length(data.alicloud_vswitches.default.ids) > 0 ? 0 : 1
  vpc_id            = data.alicloud_vpcs.default.ids.0
  cidr_block        = cidrsubnet(data.alicloud_vpcs.default.vpcs[0].cidr_block, 8, 8)
  zone_id           = data.alicloud_zones.default.zones.0.id
  vswitch_name      = var.name
}

locals {
  vswitch_id = length(data.alicloud_vswitches.default.ids) > 0 ? data.alicloud_vswitches.default.ids[0] : concat(alicloud_vswitch.vswitch.*.id, [""])[0]
}

resource "alicloud_security_group" "default" {
  name   = "${var.name}"
  vpc_id = data.alicloud_vpcs.default.ids.0
}

resource "alicloud_security_group_rule" "default" {
  	type = "ingress"
  	ip_protocol = "tcp"
  	nic_type = "intranet"
  	policy = "accept"
  	port_range = "22/22"
  	priority = 1
  	security_group_id = "${alicloud_security_group.default.id}"
  	cidr_ip = "172.16.0.0/24"
}

resource "alicloud_ess_scaling_group" "default" {
	min_size = 0
	max_size = 2
	default_cooldown = 20
	removal_policies = ["OldestInstance", "NewestInstance"]
	scaling_group_name = "${var.name}"
	vswitch_ids = ["${local.vswitch_id}"]
}

resource "alicloud_cms_namespace" "default" {
	  description   = "ttt"
	  namespace     = "ttt"
	  specification = "cms.s1.large"
	}
	
resource "alicloud_cms_hybrid_monitor_fc_task" "default" {
  namespace      = alicloud_cms_namespace.default.id
  yarm_config    = <<EOF
---
products:
- namespace: "acs_ecs_dashboard"
  metric_info:
  - metric_list:
    - "CPUUtilization"
    - "DiskReadBPS"
    - "InternetOut"
    - "IntranetOut"
    - "cpu_idle"
    - "cpu_system"
    - "cpu_total"
    - "diskusage_utilization"
- namespace: "acs_rds_dashboard"
  metric_info:
  - metric_list:
    - "MySQL_QPS"
    - "MySQL_TPS"
EOF
  target_user_id = data.alicloud_account.default.id
}

resource "alicloud_ess_scaling_rule" "default"{
	scaling_group_id = "${alicloud_ess_scaling_group.default.id}"
	scaling_rule_name = "${var.name}"
	scaling_rule_type = "TargetTrackingScalingRule"
    target_value = "20.2"
    metric_type = "hybrid"
    estimated_instance_warmup = 200
	hybrid_monitor_namespace = "${alicloud_cms_namespace.default.id}"
    hybrid_metrics {
        id = "a"
        metric_name = "AliyunEcs_CPUUtilization"
        statistic = "Average"
        dimensions {
			dimension_key = "tag_ESS"	
            dimension_value = "ESS"
		}
	}
    hybrid_metrics {
        id = "Expression"
        expression = "2*a"
	}
}

data "alicloud_ess_scaling_rules" "default" {
  %s
}
`, strings.Join(pairs, "\n  "))
	return config
}

func testAccCheckAliCloudEssScalingrulesDataSourceTargetTrackingRuleDisableScaleInConfig(attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`

variable "name" {
	default = "tf-testAccDataSourceEssScalingTargetTrackingConfig"
}

data "alicloud_zones" "default" {
  available_disk_category     = "cloud_efficiency"
  available_resource_creation = "VSwitch"
}

data "alicloud_instance_types" "default" {
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
  cpu_core_count    = 2
  memory_size       = 4
}

data "alicloud_images" "default" {
  name_regex  = "^ubuntu_[0-9]+_[0-9]+_x64*"
  most_recent = true
  owners      = "system"
}

data "alicloud_vpcs" "default" {
	
}
data "alicloud_vswitches" "default" {
	vpc_id = data.alicloud_vpcs.default.ids.0
	zone_id      = data.alicloud_zones.default.zones.0.id
}

resource "alicloud_vswitch" "vswitch" {
  count             = length(data.alicloud_vswitches.default.ids) > 0 ? 0 : 1
  vpc_id            = data.alicloud_vpcs.default.ids.0
  cidr_block        = cidrsubnet(data.alicloud_vpcs.default.vpcs[0].cidr_block, 8, 8)
  zone_id           = data.alicloud_zones.default.zones.0.id
  vswitch_name      = var.name
}

locals {
  vswitch_id = length(data.alicloud_vswitches.default.ids) > 0 ? data.alicloud_vswitches.default.ids[0] : concat(alicloud_vswitch.vswitch.*.id, [""])[0]
}

resource "alicloud_security_group" "default" {
  name   = "${var.name}"
  vpc_id = data.alicloud_vpcs.default.ids.0
}

resource "alicloud_security_group_rule" "default" {
  	type = "ingress"
  	ip_protocol = "tcp"
  	nic_type = "intranet"
  	policy = "accept"
  	port_range = "22/22"
  	priority = 1
  	security_group_id = "${alicloud_security_group.default.id}"
  	cidr_ip = "172.16.0.0/24"
}

resource "alicloud_ess_scaling_group" "default" {
	min_size = 0
	max_size = 2
	default_cooldown = 20
	removal_policies = ["OldestInstance", "NewestInstance"]
	scaling_group_name = "${var.name}"
	vswitch_ids = ["${local.vswitch_id}"]
}

resource "alicloud_ess_scaling_rule" "default"{
	scaling_group_id = "${alicloud_ess_scaling_group.default.id}"
	scaling_rule_name = "${var.name}"
	scaling_rule_type = "TargetTrackingScalingRule"
    target_value = "20.2"
	metric_name = "CpuUtilization"
    metric_type = "system"
	estimated_instance_warmup = 200
	disable_scale_in = true
}

data "alicloud_ess_scaling_rules" "default" {
  %s
}
`, strings.Join(pairs, "\n  "))
	return config
}

func testAccCheckAliCloudEssScalingrulesDataSourceTargetTrackingRuleConfig(attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`

variable "name" {
	default = "tf-testAccDataSourceEssScalingTargetTrackingConfig"
}

data "alicloud_zones" "default" {
  available_disk_category     = "cloud_efficiency"
  available_resource_creation = "VSwitch"
}

data "alicloud_instance_types" "default" {
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
  cpu_core_count    = 2
  memory_size       = 4
}

data "alicloud_images" "default" {
  name_regex  = "^ubuntu_[0-9]+_[0-9]+_x64*"
  most_recent = true
  owners      = "system"
}

data "alicloud_vpcs" "default" {
	
}
data "alicloud_vswitches" "default" {
	vpc_id = data.alicloud_vpcs.default.ids.0
	zone_id      = data.alicloud_zones.default.zones.0.id
}

resource "alicloud_vswitch" "vswitch" {
  count             = length(data.alicloud_vswitches.default.ids) > 0 ? 0 : 1
  vpc_id            = data.alicloud_vpcs.default.ids.0
  cidr_block        = cidrsubnet(data.alicloud_vpcs.default.vpcs[0].cidr_block, 8, 8)
  zone_id           = data.alicloud_zones.default.zones.0.id
  vswitch_name      = var.name
}

locals {
  vswitch_id = length(data.alicloud_vswitches.default.ids) > 0 ? data.alicloud_vswitches.default.ids[0] : concat(alicloud_vswitch.vswitch.*.id, [""])[0]
}

resource "alicloud_security_group" "default" {
  name   = "${var.name}"
  vpc_id = data.alicloud_vpcs.default.ids.0
}

resource "alicloud_security_group_rule" "default" {
  	type = "ingress"
  	ip_protocol = "tcp"
  	nic_type = "intranet"
  	policy = "accept"
  	port_range = "22/22"
  	priority = 1
  	security_group_id = "${alicloud_security_group.default.id}"
  	cidr_ip = "172.16.0.0/24"
}

resource "alicloud_ess_scaling_group" "default" {
	min_size = 0
	max_size = 2
	default_cooldown = 20
	removal_policies = ["OldestInstance", "NewestInstance"]
	scaling_group_name = "${var.name}"
	vswitch_ids = ["${local.vswitch_id}"]
}

resource "alicloud_ess_scaling_rule" "default"{
	scaling_group_id = "${alicloud_ess_scaling_group.default.id}"
	scaling_rule_name = "${var.name}"
	scaling_rule_type = "TargetTrackingScalingRule"
    target_value = "20.2"
	metric_name = "CpuUtilization"
    metric_type = "system"
	estimated_instance_warmup = 200
	scale_in_evaluation_count = 5
    scale_out_evaluation_count = 1
}

data "alicloud_ess_scaling_rules" "default" {
  %s
}
`, strings.Join(pairs, "\n  "))
	return config
}

func testAccCheckAliCloudEssScalingrulesDataSourceStepRuleConfig(attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`

variable "name" {
	default = "tf-testAccDataSourceEssScalingStepRulesCn"
}

data "alicloud_zones" "default" {
  available_disk_category     = "cloud_efficiency"
  available_resource_creation = "VSwitch"
}

data "alicloud_instance_types" "default" {
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
  cpu_core_count    = 2
  memory_size       = 4
}

data "alicloud_images" "default" {
  name_regex  = "^ubuntu_[0-9]+_[0-9]+_x64*"
  most_recent = true
  owners      = "system"
}

data "alicloud_vpcs" "default" {
	
}
data "alicloud_vswitches" "default" {
	vpc_id = data.alicloud_vpcs.default.ids.0
	zone_id      = data.alicloud_zones.default.zones.0.id
}

resource "alicloud_vswitch" "vswitch" {
  count             = length(data.alicloud_vswitches.default.ids) > 0 ? 0 : 1
  vpc_id            = data.alicloud_vpcs.default.ids.0
  cidr_block        = cidrsubnet(data.alicloud_vpcs.default.vpcs[0].cidr_block, 8, 8)
  zone_id           = data.alicloud_zones.default.zones.0.id
  vswitch_name      = var.name
}

locals {
  vswitch_id = length(data.alicloud_vswitches.default.ids) > 0 ? data.alicloud_vswitches.default.ids[0] : concat(alicloud_vswitch.vswitch.*.id, [""])[0]
}

resource "alicloud_security_group" "default" {
  name   = "${var.name}"
  vpc_id = data.alicloud_vpcs.default.ids.0
}

resource "alicloud_security_group_rule" "default" {
  	type = "ingress"
  	ip_protocol = "tcp"
  	nic_type = "intranet"
  	policy = "accept"
  	port_range = "22/22"
  	priority = 1
  	security_group_id = "${alicloud_security_group.default.id}"
  	cidr_ip = "172.16.0.0/24"
}

resource "alicloud_ess_scaling_group" "default" {
	min_size = 0
	max_size = 2
	default_cooldown = 20
	removal_policies = ["OldestInstance", "NewestInstance"]
	scaling_group_name = "${var.name}"
	vswitch_ids = ["${local.vswitch_id}"]
}

resource "alicloud_ess_scaling_rule" "default"{
	scaling_group_id = "${alicloud_ess_scaling_group.default.id}"
	scaling_rule_name = "${var.name}"
	scaling_rule_type = "StepScalingRule"
    adjustment_type = "TotalCapacity"
	estimated_instance_warmup = 200
    step_adjustment {
        metric_interval_lower_bound = "10.3"
        metric_interval_upper_bound = "20.1"
        scaling_adjustment = "1"
	}
    step_adjustment {
        metric_interval_lower_bound = "20.1"
        scaling_adjustment = "2"
	}
}

data "alicloud_ess_scaling_rules" "default" {
  %s
}
`, strings.Join(pairs, "\n  "))
	return config
}

func testAccCheckAliCloudEssScalingrulesDataSourcePredictiveRuleConfig(attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`

variable "name" {
	default = "tf-testAccDataSourceEssScalingPredictiveRule"
}

data "alicloud_zones" "default" {
  available_disk_category     = "cloud_efficiency"
  available_resource_creation = "VSwitch"
}

data "alicloud_instance_types" "default" {
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
  cpu_core_count    = 2
  memory_size       = 4
}

data "alicloud_images" "default" {
  name_regex  = "^ubuntu_[0-9]+_[0-9]+_x64*"
  most_recent = true
  owners      = "system"
}

data "alicloud_vpcs" "default" {
	
}
data "alicloud_vswitches" "default" {
	vpc_id = data.alicloud_vpcs.default.ids.0
	zone_id      = data.alicloud_zones.default.zones.0.id
}

resource "alicloud_vswitch" "vswitch" {
  count             = length(data.alicloud_vswitches.default.ids) > 0 ? 0 : 1
  vpc_id            = data.alicloud_vpcs.default.ids.0
  cidr_block        = cidrsubnet(data.alicloud_vpcs.default.vpcs[0].cidr_block, 8, 8)
  zone_id           = data.alicloud_zones.default.zones.0.id
  vswitch_name      = var.name
}

locals {
  vswitch_id = length(data.alicloud_vswitches.default.ids) > 0 ? data.alicloud_vswitches.default.ids[0] : concat(alicloud_vswitch.vswitch.*.id, [""])[0]
}

resource "alicloud_security_group" "default" {
  name   = "${var.name}"
  vpc_id = data.alicloud_vpcs.default.ids.0
}

resource "alicloud_security_group_rule" "default" {
  	type = "ingress"
  	ip_protocol = "tcp"
  	nic_type = "intranet"
  	policy = "accept"
  	port_range = "22/22"
  	priority = 1
  	security_group_id = "${alicloud_security_group.default.id}"
  	cidr_ip = "172.16.0.0/24"
}

resource "alicloud_ess_scaling_group" "default" {
	min_size = 0
	max_size = 2
	default_cooldown = 20
	removal_policies = ["OldestInstance", "NewestInstance"]
	scaling_group_name = "${var.name}"
	vswitch_ids = ["${local.vswitch_id}"]
}

resource "alicloud_ess_scaling_rule" "default"{
	scaling_group_id = "${alicloud_ess_scaling_group.default.id}"
	scaling_rule_name = "${var.name}"
	initial_max_size = 1
	predictive_scaling_mode = "PredictAndScale"
	predictive_value_behavior = "MaxOverridePredictiveValue"
	predictive_value_buffer = 0
    predictive_task_buffer_time = 0
	scaling_rule_type = "PredictiveScalingRule"
    metric_name = "CpuUtilization"
	target_value = 20.1
}

data "alicloud_ess_scaling_rules" "default" {
  %s
}
`, strings.Join(pairs, "\n  "))
	return config
}

func testAccCheckAliCloudEssScalingrulesDataSourceConfig(attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`

variable "name" {
	default = "tf-testAccDataSourceEssScalingRules"
}

data "alicloud_zones" "default" {
  available_disk_category     = "cloud_efficiency"
  available_resource_creation = "VSwitch"
}

data "alicloud_instance_types" "default" {
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
  cpu_core_count    = 2
  memory_size       = 4
}

data "alicloud_images" "default" {
  name_regex  = "^ubuntu_[0-9]+_[0-9]+_x64*"
  most_recent = true
  owners      = "system"
}

data "alicloud_vpcs" "default" {
	
}
data "alicloud_vswitches" "default" {
	vpc_id = data.alicloud_vpcs.default.ids.0
	zone_id      = data.alicloud_zones.default.zones.0.id
}

resource "alicloud_vswitch" "vswitch" {
  count             = length(data.alicloud_vswitches.default.ids) > 0 ? 0 : 1
  vpc_id            = data.alicloud_vpcs.default.ids.0
  cidr_block        = cidrsubnet(data.alicloud_vpcs.default.vpcs[0].cidr_block, 8, 8)
  zone_id           = data.alicloud_zones.default.zones.0.id
  vswitch_name      = var.name
}

locals {
  vswitch_id = length(data.alicloud_vswitches.default.ids) > 0 ? data.alicloud_vswitches.default.ids[0] : concat(alicloud_vswitch.vswitch.*.id, [""])[0]
}

resource "alicloud_security_group" "default" {
  name   = "${var.name}"
  vpc_id = data.alicloud_vpcs.default.ids.0
}

resource "alicloud_security_group_rule" "default" {
  	type = "ingress"
  	ip_protocol = "tcp"
  	nic_type = "intranet"
  	policy = "accept"
  	port_range = "22/22"
  	priority = 1
  	security_group_id = "${alicloud_security_group.default.id}"
  	cidr_ip = "172.16.0.0/24"
}

resource "alicloud_ess_scaling_group" "default" {
	min_size = 0
	max_size = 2
	default_cooldown = 20
	removal_policies = ["OldestInstance", "NewestInstance"]
	scaling_group_name = "${var.name}"
	vswitch_ids = ["${local.vswitch_id}"]
}

resource "alicloud_ess_scaling_rule" "default"{
	scaling_group_id = "${alicloud_ess_scaling_group.default.id}"
	scaling_rule_name = "${var.name}"
	adjustment_type = "PercentChangeInCapacity"
	adjustment_value = 1
	cooldown = 30
    min_adjustment_magnitude = 1
}

data "alicloud_ess_scaling_rules" "default" {
  %s
}
`, strings.Join(pairs, "\n  "))
	return config
}

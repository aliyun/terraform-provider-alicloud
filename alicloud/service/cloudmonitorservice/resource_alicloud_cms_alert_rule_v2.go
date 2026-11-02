// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package cloudmonitorservice

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/convert"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// cmsAlertRuleV2WorkspaceNotReady reports whether the API rejected the request
// because a freshly created workspace is not visible yet; the call is safe to retry.
func cmsAlertRuleV2WorkspaceNotReady(err error) bool {
	if err == nil {
		return false
	}
	message := err.Error()
	return strings.Contains(message, "工作空间校验失败") && strings.Contains(message, "404")
}

// cmsAlertRuleV2UserLabels drops the server-managed labels (prefixed "_cms_")
// that the API attaches on read, so state only mirrors the configured map.
func cmsAlertRuleV2UserLabels(raw interface{}) map[string]interface{} {
	source, ok := raw.(map[string]interface{})
	if !ok {
		return nil
	}
	userLabels := make(map[string]interface{}, len(source))
	for key, value := range source {
		if strings.HasPrefix(key, "_cms_") {
			continue
		}
		userLabels[key] = value
	}
	return userLabels
}

func cmsAlertRuleV2NonZeroNumber(v interface{}) bool {
	switch value := v.(type) {
	case int:
		return value != 0
	case int64:
		return value != 0
	case float32:
		return value != 0
	case float64:
		return value != 0
	}
	return true
}

// cmsAlertRuleV2NestedList converts a nested Terraform list attribute into its API
// request form, renaming the Terraform keys to the API keys and dropping the unset
// or zero values the API rejects.
func cmsAlertRuleV2NestedList(v interface{}, field string, keyMapping map[string]string, zeroAllowed ...string) []map[string]interface{} {
	allowZero := make(map[string]bool, len(zeroAllowed))
	for _, key := range zeroAllowed {
		allowZero[key] = true
	}
	raw, _ := jsonpath.Get("$[0]."+field, v)
	if raw == nil {
		return nil
	}
	result := make([]map[string]interface{}, 0)
	for _, childRaw := range convert.ConvertToInterfaceArray(raw) {
		child, ok := childRaw.(map[string]interface{})
		if !ok {
			continue
		}
		item := make(map[string]interface{})
		for tfKey, apiKey := range keyMapping {
			value, exists := child[tfKey]
			if !exists || value == nil || value == "" {
				continue
			}
			switch value.(type) {
			case int, int64, float32, float64:
				if !cmsAlertRuleV2NonZeroNumber(value) && !allowZero[tfKey] {
					continue
				}
			case []interface{}:
				if len(value.([]interface{})) == 0 {
					continue
				}
			}
			item[apiKey] = value
		}
		if len(item) > 0 {
			result = append(result, item)
		}
	}
	if len(result) == 0 {
		return nil
	}
	return result
}

// cmsAlertRuleV2JoiningList converts the joinings attribute into its API form,
// including the nested join conditions.
func cmsAlertRuleV2JoiningList(v interface{}) []map[string]interface{} {
	raw, _ := jsonpath.Get("$[0].joinings", v)
	if raw == nil {
		return nil
	}
	result := make([]map[string]interface{}, 0)
	for _, childRaw := range convert.ConvertToInterfaceArray(raw) {
		child, ok := childRaw.(map[string]interface{})
		if !ok {
			continue
		}
		item := make(map[string]interface{})
		if value, ok := child["type"].(string); ok && value != "" {
			item["type"] = value
		}
		conditions := make([]map[string]interface{}, 0)
		for _, conditionRaw := range convert.ConvertToInterfaceArray(child["conditions"]) {
			condition, ok := conditionRaw.(map[string]interface{})
			if !ok {
				continue
			}
			conditionItem := make(map[string]interface{})
			if value, ok := condition["lhs_field"].(string); ok && value != "" {
				conditionItem["lhsField"] = value
			}
			if value, ok := condition["rhs_field"].(string); ok && value != "" {
				conditionItem["rhsField"] = value
			}
			if value, ok := condition["operator"].(string); ok && value != "" {
				conditionItem["operator"] = value
			}
			if len(conditionItem) > 0 {
				conditions = append(conditions, conditionItem)
			}
		}
		if len(conditions) > 0 {
			item["conditions"] = conditions
		}
		if len(item) > 0 {
			result = append(result, item)
		}
	}
	if len(result) == 0 {
		return nil
	}
	return result
}

func cmsAlertRuleV2CustomizeDiff(_ context.Context, d *schema.ResourceDiff, _ interface{}) error {
	for _, block := range []string{"datasource_config", "query_config", "condition_config", "schedule_config", "notify_config"} {
		if v, ok := d.GetOk(block); !ok || len(v.([]interface{})) == 0 {
			return fmt.Errorf("%s is required", block)
		}
	}
	if v, ok := d.GetOk("query_config"); ok {
		queries := v.([]interface{})
		if len(queries) > 0 && queries[0] != nil {
			query := queries[0].(map[string]interface{})
			if query["type"] == "APM_MULTI_QUERY" {
				if identifiers, _ := query["service_id_list"].([]interface{}); len(identifiers) == 0 {
					return fmt.Errorf("query_config.service_id_list is required when query_config.type is APM_MULTI_QUERY")
				}
				if measures, _ := query["measure_list"].([]interface{}); len(measures) == 0 {
					return fmt.Errorf("query_config.measure_list is required when query_config.type is APM_MULTI_QUERY")
				}
			}
			if query["type"] == "UMODEL_LOGSET_QUERY" {
				if logSet, _ := query["log_set"].(string); logSet == "" {
					return fmt.Errorf("query_config.log_set is required when query_config.type is UMODEL_LOGSET_QUERY")
				}
				if expr, _ := query["expr"].(string); expr == "" {
					return fmt.Errorf("query_config.expr is required when query_config.type is UMODEL_LOGSET_QUERY")
				}
				if windowSecs, _ := query["window_secs"].(int); windowSecs == 0 {
					return fmt.Errorf("query_config.window_secs is required when query_config.type is UMODEL_LOGSET_QUERY")
				}
			}
			if query["type"] == "UMODEL_METRICSET_QUERY" {
				if windowSecs, _ := query["window_secs"].(int); windowSecs == 0 {
					return fmt.Errorf("query_config.window_secs is required when query_config.type is UMODEL_METRICSET_QUERY")
				}
				if aggregate, _ := query["aggregate"].(string); aggregate == "" {
					return fmt.Errorf("query_config.aggregate is required when query_config.type is UMODEL_METRICSET_QUERY")
				}
			}
		}
	}
	if v, ok := d.GetOk("datasource_config"); ok {
		configs := v.([]interface{})
		if len(configs) > 0 && configs[0] != nil {
			config := configs[0].(map[string]interface{})
			if config["type"] == "SLS" {
				if project, _ := config["project"].(string); project == "" {
					return fmt.Errorf("datasource_config.project is required when datasource_config.type is SLS")
				}
				if stores, _ := config["stores"].([]interface{}); len(stores) == 0 {
					return fmt.Errorf("datasource_config.stores requires at least one store when datasource_config.type is SLS")
				}
			}
		}
	}
	if v, ok := d.GetOk("condition_config"); ok {
		configs := v.([]interface{})
		if len(configs) > 0 && configs[0] != nil {
			config := configs[0].(map[string]interface{})
			if config["type"] == "SLS_MULTI_CONDITION" {
				if triggers, _ := config["triggers"].([]interface{}); len(triggers) == 0 {
					return fmt.Errorf("condition_config.triggers requires at least one trigger when condition_config.type is SLS_MULTI_CONDITION")
				}
			}
			if config["type"] == "UMODEL_LOGSET_CONDITION" {
				matchOperator, _ := config["match_operator"].(string)
				if matchOperator == "" {
					return fmt.Errorf("condition_config.match_operator is required when condition_config.type is UMODEL_LOGSET_CONDITION")
				}
				if durationSecs, _ := config["duration_secs"].(int); durationSecs == 0 {
					return fmt.Errorf("condition_config.duration_secs is required when condition_config.type is UMODEL_LOGSET_CONDITION")
				}
				if severity, _ := config["severity"].(string); severity == "" {
					return fmt.Errorf("condition_config.severity is required when condition_config.type is UMODEL_LOGSET_CONDITION")
				}
				if matchOperator == "CONTAINS" || matchOperator == "EQUALS" || matchOperator == "REGEX" {
					if matchField, _ := config["match_field"].(string); matchField == "" {
						return fmt.Errorf("condition_config.match_field is required when condition_config.match_operator is CONTAINS, EQUALS or REGEX")
					}
					if matchValue, _ := config["match_value"].(string); matchValue == "" {
						return fmt.Errorf("condition_config.match_value is required when condition_config.match_operator is CONTAINS, EQUALS or REGEX")
					}
				}
			}
			if config["type"] == "UMODEL_METRICSET_CONDITION" || config["type"] == "UMODEL_METRICSET_MULTI_CONDITION" {
				noDataPolicy, _ := config["no_data_policy"].(string)
				noDataAlertSeverity, _ := config["no_data_alert_severity"].(string)
				if noDataPolicy == "NO_DATA_TO_ALERT" && noDataAlertSeverity == "" {
					return fmt.Errorf("condition_config.no_data_alert_severity is required when condition_config.no_data_policy is NO_DATA_TO_ALERT for %s", config["type"])
				}
				if noDataPolicy != "" && noDataPolicy != "NO_DATA_TO_ALERT" && noDataAlertSeverity != "" {
					return fmt.Errorf("condition_config.no_data_alert_severity must not be set when condition_config.no_data_policy is not NO_DATA_TO_ALERT for %s", config["type"])
				}
			}
		}
	}
	return nil
}

func resourceAliCloudCmsAlertRuleV2() *schema.Resource {
	return &schema.Resource{
		Create:        resourceAliCloudCmsAlertRuleV2Create,
		Read:          resourceAliCloudCmsAlertRuleV2Read,
		Update:        resourceAliCloudCmsAlertRuleV2Update,
		Delete:        resourceAliCloudCmsAlertRuleV2Delete,
		CustomizeDiff: cmsAlertRuleV2CustomizeDiff,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"action_integration_config": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"actions": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							ForceNew: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"enabled": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
							ForceNew: true,
						},
					},
				},
			},
			"alert_rule_v2_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"annotations": {
				Type:     schema.TypeMap,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"arms_integration_config": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
							ForceNew: true,
						},
					},
				},
			},
			"condition_config": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"operator": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							ForceNew: true,
						},
						"legacy_raw": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"prometheus": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"prom_ql": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"times": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"severity": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"no_data_duration_secs": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
							ForceNew: true,
						},
						"match_field": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"count_threshold": {
							Type:     schema.TypeInt,
							Optional: true,
							ForceNew: true,
						},
						"compare_list": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							ForceNew: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"operator": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
										ForceNew: true,
									},
									"aggregate": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
										ForceNew: true,
									},
									"yoy_time_value": {
										Type:     schema.TypeInt,
										Optional: true,
										Computed: true,
										ForceNew: true,
									},
									"threshold": {
										Type:     schema.TypeFloat,
										Optional: true,
										Computed: true,
										ForceNew: true,
									},
									"yoy_time_unit": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
										ForceNew: true,
									},
								},
							},
						},
						"relation": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"express_escalation": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"raw_expression": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"times": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"severity": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"yoy_time_value": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
							ForceNew: true,
						},
						"count_operator": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"threshold": {
							Type:     schema.TypeFloat,
							Optional: true,
							Computed: true,
							ForceNew: true,
						},
						"legacy_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"severity": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							ForceNew: true,
						},
						"no_data_policy": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"no_data_alert_severity": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"simple_escalation": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"metric_name": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
									"escalations": {
										Type:     schema.TypeList,
										Optional: true,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"comparison_operator": {
													Type:     schema.TypeString,
													Optional: true,
													Computed: true,
												},
												"times": {
													Type:     schema.TypeInt,
													Optional: true,
													Computed: true,
												},
												"pre_condition": {
													Type:     schema.TypeString,
													Optional: true,
													Computed: true,
												},
												"severity": {
													Type:     schema.TypeString,
													Optional: true,
													Computed: true,
												},
												"statistics": {
													Type:     schema.TypeString,
													Optional: true,
													Computed: true,
												},
												"threshold": {
													Type:     schema.TypeString,
													Optional: true,
													Computed: true,
												},
											},
										},
									},
									"period": {
										Type:     schema.TypeInt,
										Optional: true,
										Computed: true,
									},
								},
							},
						},
						"yoy_time_unit": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							ForceNew: true,
						},
						"escalation_type": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"composite_escalation": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"relation": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
									"times": {
										Type:     schema.TypeInt,
										Optional: true,
										Computed: true,
									},
									"escalations": {
										Type:     schema.TypeList,
										Optional: true,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"metric_name": {
													Type:     schema.TypeString,
													Optional: true,
													Computed: true,
												},
												"comparison_operator": {
													Type:     schema.TypeString,
													Optional: true,
													Computed: true,
												},
												"pre_condition": {
													Type:     schema.TypeString,
													Optional: true,
													Computed: true,
												},
												"period": {
													Type:     schema.TypeInt,
													Optional: true,
													Computed: true,
												},
												"statistics": {
													Type:     schema.TypeString,
													Optional: true,
													Computed: true,
												},
												"threshold": {
													Type:     schema.TypeString,
													Optional: true,
													Computed: true,
												},
											},
										},
									},
									"severity": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
								},
							},
						},
						"duration_secs": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
							ForceNew: true,
						},
						"match_operator": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"match_value": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"aggregate": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							ForceNew: true,
						},
						"threshold_list": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							ForceNew: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"severity": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
										ForceNew: true,
									},
									"threshold": {
										Type:     schema.TypeFloat,
										Optional: true,
										Computed: true,
										ForceNew: true,
									},
								},
							},
						},
						"triggers": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							ForceNew: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"operator": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
										ForceNew: true,
									},
									"condition": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
										ForceNew: true,
									},
									"count_condition": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
										ForceNew: true,
									},
									"severity": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
										ForceNew: true,
									},
								},
							},
						},
						"alert_count": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
							ForceNew: true,
						},
						"no_data_alert_level": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							ForceNew: true,
						},
					},
				},
			},
			"content_template": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"datasource_config": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"legacy_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"project": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"type": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"instance_id": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							ForceNew: true,
						},
						"legacy_raw": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"product_category": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"region_id": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							ForceNew: true,
						},
						"stores": {
							Type:     schema.TypeList,
							Optional: true,
							ForceNew: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"store_type": {
										Type:     schema.TypeString,
										Optional: true,
										ForceNew: true,
									},
									"store": {
										Type:     schema.TypeString,
										Optional: true,
										ForceNew: true,
									},
								},
							},
						},
					},
				},
			},
			"datasource_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"display_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"labels": {
				Type:     schema.TypeMap,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"notify_config": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"utc_offset": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							ForceNew: true,
						},
						"notify_strategies": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							ForceNew: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"type": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"active_days": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							ForceNew: true,
							Elem:     &schema.Schema{Type: schema.TypeInt},
						},
						"active_end_time": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							ForceNew: true,
						},
						"active_start_time": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							ForceNew: true,
						},
						"silence_time_secs": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
							ForceNew: true,
						},
						"channels": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							ForceNew: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
										ForceNew: true,
									},
									"identifiers": {
										Type:     schema.TypeList,
										Optional: true,
										Computed: true,
										ForceNew: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
					},
				},
			},
			"notify_strategy_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"observe_resource_config": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"entity_type": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"entity_domain": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"relation_type": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"product_category": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"resources": {
							Type:     schema.TypeList,
							Optional: true,
							ForceNew: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"namespace": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"observe_resource_global_scope": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"observe_resource_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"partition_key": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"query_config": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"legacy_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"log_set": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"offset_secs": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
							ForceNew: true,
						},
						"legacy_raw": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"prom_ql": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"entity_filters": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							ForceNew: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"operator": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
										ForceNew: true,
									},
									"field": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
										ForceNew: true,
									},
									"value": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
										ForceNew: true,
									},
								},
							},
						},
						"window_secs": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
							ForceNew: true,
						},
						"dimensions": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeMap, Elem: &schema.Schema{Type: schema.TypeString}},
						},
						"filter_list": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							ForceNew: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
										ForceNew: true,
									},
									"value": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
										ForceNew: true,
									},
									"key": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
										ForceNew: true,
									},
								},
							},
						},
						"label_filters": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							ForceNew: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"operator": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
										ForceNew: true,
									},
									"value": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
										ForceNew: true,
									},
									"name": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
										ForceNew: true,
									},
								},
							},
						},
						"namespace": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"group_id": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"metric_set": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							ForceNew: true,
						},
						"entity_type": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							ForceNew: true,
						},
						"measure_list": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							ForceNew: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"measure_code": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
										ForceNew: true,
									},
									"group_by": {
										Type:     schema.TypeList,
										Optional: true,
										Computed: true,
										ForceNew: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"window_secs": {
										Type:     schema.TypeInt,
										Optional: true,
										Computed: true,
										ForceNew: true,
									},
								},
							},
						},
						"type": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"expr": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							ForceNew: true,
						},
						"entity_domain": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							ForceNew: true,
						},
						"relation_type": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"metric": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							ForceNew: true,
						},
						"service_id_list": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							ForceNew: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"aggregate": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"enable_data_complete_check": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
							ForceNew: true,
						},
						"entity_fields": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							ForceNew: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"field": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
										ForceNew: true,
									},
									"value": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
										ForceNew: true,
									},
								},
							},
						},
						"group_type": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							ForceNew: true,
						},
						"queries": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							ForceNew: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"expr": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
										ForceNew: true,
									},
									"time_unit": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
										ForceNew: true,
									},
									"start": {
										Type:     schema.TypeInt,
										Optional: true,
										Computed: true,
										ForceNew: true,
									},
									"end": {
										Type:     schema.TypeInt,
										Optional: true,
										Computed: true,
										ForceNew: true,
									},
									"window": {
										Type:     schema.TypeInt,
										Optional: true,
										Computed: true,
										ForceNew: true,
									},
								},
							},
						},
						"joinings": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							ForceNew: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
										ForceNew: true,
									},
									"conditions": {
										Type:     schema.TypeList,
										Optional: true,
										Computed: true,
										ForceNew: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"lhs_field": {
													Type:     schema.TypeString,
													Optional: true,
													Computed: true,
													ForceNew: true,
												},
												"rhs_field": {
													Type:     schema.TypeString,
													Optional: true,
													Computed: true,
													ForceNew: true,
												},
												"operator": {
													Type:     schema.TypeString,
													Optional: true,
													Computed: true,
													ForceNew: true,
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			"rca_config": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"digital_employee_name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"enable_rca": {
							Type:     schema.TypeBool,
							Optional: true,
							ForceNew: true,
						},
					},
				},
			},
			"schedule_config": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"interval_secs": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
							ForceNew: true,
						},
					},
				},
			},
			"severity_levels": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"workspace": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAliCloudCmsAlertRuleV2Create(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := fmt.Sprintf("/manageAlertRules")
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	body := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	query["RegionId"] = helper.StringPointer(client.RegionId)

	actionIntegrationConfig := make(map[string]interface{})

	if v := d.Get("action_integration_config"); !helper.IsNil(v) {
		actions1, _ := jsonpath.Get("$[0].actions", v)
		if actions1 != nil && actions1 != "" {
			actionIntegrationConfig["actions"] = actions1
		}
		enabled1, _ := jsonpath.Get("$[0].enabled", v)
		if enabled1 != nil && enabled1 != "" {
			actionIntegrationConfig["enabled"] = enabled1
		}

		request["actionIntegrationConfig"] = actionIntegrationConfig
	}

	if v, ok := d.GetOk("annotations"); ok {
		request["annotations"] = v
	}
	armsIntegrationConfig := make(map[string]interface{})

	if v := d.Get("arms_integration_config"); !helper.IsNil(v) {
		enabled3, _ := jsonpath.Get("$[0].enabled", v)
		if enabled3 != nil && enabled3 != "" {
			armsIntegrationConfig["enabled"] = enabled3
		}

		request["armsIntegrationConfig"] = armsIntegrationConfig
	}

	conditionConfig := make(map[string]interface{})

	if v := d.Get("condition_config"); !helper.IsNil(v) {
		aggregate1, _ := jsonpath.Get("$[0].aggregate", v)
		if aggregate1 != nil && aggregate1 != "" {
			conditionConfig["aggregate"] = aggregate1
		}
		compositeEscalation := make(map[string]interface{})
		relation1, _ := jsonpath.Get("$[0].composite_escalation[0].relation", d.Get("condition_config"))
		if relation1 != nil && relation1 != "" {
			compositeEscalation["relation"] = relation1
		}
		severity1, _ := jsonpath.Get("$[0].composite_escalation[0].severity", d.Get("condition_config"))
		if severity1 != nil && severity1 != "" {
			compositeEscalation["severity"] = severity1
		}
		times1, _ := jsonpath.Get("$[0].composite_escalation[0].times", d.Get("condition_config"))
		if times1 != nil && times1 != "" {
			compositeEscalation["times"] = times1
		}

		if escalations := cmsAlertRuleV2NestedList(d.Get("condition_config"), "composite_escalation[0].escalations", map[string]string{
			"comparison_operator": "comparisonOperator",
			"metric_name":         "metricName",
			"statistics":          "statistics",
			"threshold":           "threshold",
			"pre_condition":       "preCondition",
			"period":              "period",
		}, "threshold"); len(escalations) > 0 {
			compositeEscalation["escalations"] = escalations
		}

		if len(compositeEscalation) > 0 {
			conditionConfig["compositeEscalation"] = compositeEscalation
		}
		countOperator1, _ := jsonpath.Get("$[0].count_operator", v)
		if countOperator1 != nil && countOperator1 != "" {
			conditionConfig["countOperator"] = countOperator1
		}
		countThreshold1, _ := jsonpath.Get("$[0].count_threshold", v)
		if countThreshold1 != nil && countThreshold1 != "" && cmsAlertRuleV2NonZeroNumber(countThreshold1) {
			conditionConfig["countThreshold"] = countThreshold1
		}
		durationSecs1, _ := jsonpath.Get("$[0].duration_secs", v)
		if durationSecs1 != nil && durationSecs1 != "" && cmsAlertRuleV2NonZeroNumber(durationSecs1) {
			conditionConfig["durationSecs"] = durationSecs1
		}
		escalationType1, _ := jsonpath.Get("$[0].escalation_type", v)
		if escalationType1 != nil && escalationType1 != "" {
			conditionConfig["escalationType"] = escalationType1
		}
		matchField1, _ := jsonpath.Get("$[0].match_field", v)
		if matchField1 != nil && matchField1 != "" {
			conditionConfig["matchField"] = matchField1
		}
		matchOperator1, _ := jsonpath.Get("$[0].match_operator", v)
		if matchOperator1 != nil && matchOperator1 != "" {
			conditionConfig["matchOperator"] = matchOperator1
		}
		matchValue1, _ := jsonpath.Get("$[0].match_value", v)
		if matchValue1 != nil && matchValue1 != "" {
			conditionConfig["matchValue"] = matchValue1
		}
		noDataAlertSeverity1, _ := jsonpath.Get("$[0].no_data_alert_severity", v)
		if noDataAlertSeverity1 != nil && noDataAlertSeverity1 != "" {
			conditionConfig["noDataAlertSeverity"] = noDataAlertSeverity1
		}
		noDataDurationSecs1, _ := jsonpath.Get("$[0].no_data_duration_secs", v)
		if noDataDurationSecs1 != nil && noDataDurationSecs1 != "" && cmsAlertRuleV2NonZeroNumber(noDataDurationSecs1) {
			conditionConfig["noDataDurationSecs"] = noDataDurationSecs1
		}
		noDataPolicy1, _ := jsonpath.Get("$[0].no_data_policy", v)
		if noDataPolicy1 != nil && noDataPolicy1 != "" {
			conditionConfig["noDataPolicy"] = noDataPolicy1
		}
		operator1, _ := jsonpath.Get("$[0].operator", v)
		if operator1 != nil && operator1 != "" {
			conditionConfig["operator"] = operator1
		}
		relation3, _ := jsonpath.Get("$[0].relation", v)
		if relation3 != nil && relation3 != "" {
			conditionConfig["relation"] = relation3
		}
		severity7, _ := jsonpath.Get("$[0].severity", v)
		if severity7 != nil && severity7 != "" {
			conditionConfig["severity"] = severity7
		}
		simpleEscalation := make(map[string]interface{})
		metricName1, _ := jsonpath.Get("$[0].simple_escalation[0].metric_name", d.Get("condition_config"))
		if metricName1 != nil && metricName1 != "" {
			simpleEscalation["metricName"] = metricName1
		}
		period1, _ := jsonpath.Get("$[0].simple_escalation[0].period", d.Get("condition_config"))
		if period1 != nil && period1 != "" {
			simpleEscalation["period"] = period1
		}

		if escalations := cmsAlertRuleV2NestedList(d.Get("condition_config"), "simple_escalation[0].escalations", map[string]string{
			"comparison_operator": "comparisonOperator",
			"statistics":          "statistics",
			"threshold":           "threshold",
			"severity":            "severity",
			"pre_condition":       "preCondition",
			"times":               "times",
		}, "threshold"); len(escalations) > 0 {
			simpleEscalation["escalations"] = escalations
		}

		if len(simpleEscalation) > 0 {
			conditionConfig["simpleEscalation"] = simpleEscalation
		}
		threshold1, _ := jsonpath.Get("$[0].threshold", v)
		if threshold1 != nil && threshold1 != "" && cmsAlertRuleV2NonZeroNumber(threshold1) {
			conditionConfig["threshold"] = threshold1
		}
		type1, _ := jsonpath.Get("$[0].type", v)
		if type1 != nil && type1 != "" {
			conditionConfig["type"] = type1
		}
		yoyTimeUnit1, _ := jsonpath.Get("$[0].yoy_time_unit", v)
		if yoyTimeUnit1 != nil && yoyTimeUnit1 != "" {
			conditionConfig["yoyTimeUnit"] = yoyTimeUnit1
		}
		yoyTimeValue1, _ := jsonpath.Get("$[0].yoy_time_value", v)
		if yoyTimeValue1 != nil && yoyTimeValue1 != "" && cmsAlertRuleV2NonZeroNumber(yoyTimeValue1) {
			conditionConfig["yoyTimeValue"] = yoyTimeValue1
		}
		if compareList := cmsAlertRuleV2NestedList(v, "compare_list", map[string]string{
			"aggregate":      "aggregate",
			"operator":       "operator",
			"threshold":      "threshold",
			"yoy_time_unit":  "yoyTimeUnit",
			"yoy_time_value": "yoyTimeValue",
		}, "threshold"); len(compareList) > 0 {
			conditionConfig["compareList"] = compareList
		}
		if thresholdList := cmsAlertRuleV2NestedList(v, "threshold_list", map[string]string{
			"severity":  "severity",
			"threshold": "threshold",
		}, "threshold"); len(thresholdList) > 0 {
			conditionConfig["thresholdList"] = thresholdList
		}
		if triggers := cmsAlertRuleV2NestedList(v, "triggers", map[string]string{
			"operator":        "operator",
			"condition":       "condition",
			"count_condition": "countCondition",
			"severity":        "severity",
		}); len(triggers) > 0 {
			conditionConfig["triggers"] = triggers
		}
		if alertCount1, _ := jsonpath.Get("$[0].alert_count", v); alertCount1 != nil && cmsAlertRuleV2NonZeroNumber(alertCount1) {
			conditionConfig["alertCount"] = alertCount1
		}
		if noDataAlertLevel1, _ := jsonpath.Get("$[0].no_data_alert_level", v); noDataAlertLevel1 != nil && noDataAlertLevel1 != "" {
			conditionConfig["noDataAlertLevel"] = noDataAlertLevel1
		}

		request["conditionConfig"] = conditionConfig
	}

	if v, ok := d.GetOk("content_template"); ok {
		request["contentTemplate"] = v
	}
	datasourceConfig := make(map[string]interface{})

	if v := d.Get("datasource_config"); !helper.IsNil(v) {
		instanceId1, _ := jsonpath.Get("$[0].instance_id", v)
		if instanceId1 != nil && instanceId1 != "" {
			datasourceConfig["instanceId"] = instanceId1
		}
		productCategory1, _ := jsonpath.Get("$[0].product_category", v)
		if productCategory1 != nil && productCategory1 != "" {
			datasourceConfig["productCategory"] = productCategory1
		}
		project1, _ := jsonpath.Get("$[0].project", v)
		if project1 != nil && project1 != "" {
			datasourceConfig["project"] = project1
		}
		regionId1, _ := jsonpath.Get("$[0].region_id", v)
		if regionId1 != nil && regionId1 != "" {
			datasourceConfig["regionId"] = regionId1
		}
		type3, _ := jsonpath.Get("$[0].type", v)
		if type3 != nil && type3 != "" {
			datasourceConfig["type"] = type3
		}
		stores3, _ := jsonpath.Get("$[0].stores", d.Get("datasource_config"))
		if stores3 != nil {
			stores := make([]map[string]interface{}, 0)
			for _, rawStore := range convert.ConvertToInterfaceArray(stores3) {
				storeRaw, ok := rawStore.(map[string]interface{})
				if !ok {
					continue
				}
				store := make(map[string]interface{})
				if v, ok := storeRaw["store"].(string); ok && v != "" {
					store["store"] = v
				}
				if v, ok := storeRaw["store_type"].(string); ok && v != "" {
					store["storeType"] = v
				}
				if len(store) > 0 {
					stores = append(stores, store)
				}
			}
			if len(stores) > 0 {
				datasourceConfig["stores"] = stores
			}
		}

		request["datasourceConfig"] = datasourceConfig
	}

	if v, ok := d.GetOk("display_name"); ok {
		request["displayName"] = v
	}
	if v, ok := d.GetOkExists("enabled"); ok {
		request["enabled"] = v
	}
	if v, ok := d.GetOk("labels"); ok {
		request["labels"] = v
	}
	notifyConfig := make(map[string]interface{})

	if v := d.Get("notify_config"); !helper.IsNil(v) {
		activeDays1, _ := jsonpath.Get("$[0].active_days", v)
		if activeDays1 != nil && activeDays1 != "" {
			notifyConfig["activeDays"] = activeDays1
		}
		activeEndTime1, _ := jsonpath.Get("$[0].active_end_time", v)
		if activeEndTime1 != nil && activeEndTime1 != "" {
			notifyConfig["activeEndTime"] = activeEndTime1
		}
		activeStartTime1, _ := jsonpath.Get("$[0].active_start_time", v)
		if activeStartTime1 != nil && activeStartTime1 != "" {
			notifyConfig["activeStartTime"] = activeStartTime1
		}
		notifyStrategies1, _ := jsonpath.Get("$[0].notify_strategies", v)
		if notifyStrategies1 != nil && notifyStrategies1 != "" {
			notifyConfig["notifyStrategies"] = notifyStrategies1
		}
		silenceTimeSecs1, _ := jsonpath.Get("$[0].silence_time_secs", v)
		if silenceTimeSecs1 != nil && silenceTimeSecs1 != "" && cmsAlertRuleV2NonZeroNumber(silenceTimeSecs1) {
			notifyConfig["silenceTimeSecs"] = silenceTimeSecs1
		}
		type5, _ := jsonpath.Get("$[0].type", v)
		if type5 != nil && type5 != "" {
			notifyConfig["type"] = type5
		}
		utcOffset1, _ := jsonpath.Get("$[0].utc_offset", v)
		if utcOffset1 != nil && utcOffset1 != "" {
			notifyConfig["utcOffset"] = utcOffset1
		}
		if channels := cmsAlertRuleV2NestedList(v, "channels", map[string]string{
			"type":        "type",
			"identifiers": "identifiers",
		}); len(channels) > 0 {
			notifyConfig["channels"] = channels
		}

		request["notifyConfig"] = notifyConfig
	}

	observeResourceConfig := make(map[string]interface{})

	if v := d.Get("observe_resource_config"); !helper.IsNil(v) {
		entityDomain1, _ := jsonpath.Get("$[0].entity_domain", v)
		if entityDomain1 != nil && entityDomain1 != "" {
			observeResourceConfig["entityDomain"] = entityDomain1
		}
		entityType1, _ := jsonpath.Get("$[0].entity_type", v)
		if entityType1 != nil && entityType1 != "" {
			observeResourceConfig["entityType"] = entityType1
		}
		namespace1, _ := jsonpath.Get("$[0].namespace", v)
		if namespace1 != nil && namespace1 != "" {
			observeResourceConfig["namespace"] = namespace1
		}
		productCategory3, _ := jsonpath.Get("$[0].product_category", v)
		if productCategory3 != nil && productCategory3 != "" {
			observeResourceConfig["productCategory"] = productCategory3
		}
		relationType1, _ := jsonpath.Get("$[0].relation_type", v)
		if relationType1 != nil && relationType1 != "" {
			observeResourceConfig["relationType"] = relationType1
		}
		resources1, _ := jsonpath.Get("$[0].resources", v)
		if resources1 != nil && resources1 != "" {
			observeResourceConfig["resources"] = resources1
		}

		request["observeResourceConfig"] = observeResourceConfig
	}

	queryConfig := make(map[string]interface{})

	if v := d.Get("query_config"); !helper.IsNil(v) {
		aggregate3, _ := jsonpath.Get("$[0].aggregate", v)
		if aggregate3 != nil && aggregate3 != "" {
			queryConfig["aggregate"] = aggregate3
		}
		dimensions1, _ := jsonpath.Get("$[0].dimensions", v)
		if dimensions1 != nil && dimensions1 != "" {
			queryConfig["dimensions"] = dimensions1
		}
		enableDataCompleteCheck1, _ := jsonpath.Get("$[0].enable_data_complete_check", v)
		if enableDataCompleteCheck1 != nil && enableDataCompleteCheck1 != "" && enableDataCompleteCheck1 != false {
			queryConfig["enableDataCompleteCheck"] = enableDataCompleteCheck1
		}
		entityDomain3, _ := jsonpath.Get("$[0].entity_domain", v)
		if entityDomain3 != nil && entityDomain3 != "" {
			queryConfig["entityDomain"] = entityDomain3
		}
		entityType3, _ := jsonpath.Get("$[0].entity_type", v)
		if entityType3 != nil && entityType3 != "" {
			queryConfig["entityType"] = entityType3
		}
		expr1, _ := jsonpath.Get("$[0].expr", v)
		if expr1 != nil && expr1 != "" {
			queryConfig["expr"] = expr1
		}
		groupId1, _ := jsonpath.Get("$[0].group_id", v)
		if groupId1 != nil && groupId1 != "" {
			queryConfig["groupId"] = groupId1
		}
		logSet1, _ := jsonpath.Get("$[0].log_set", v)
		if logSet1 != nil && logSet1 != "" {
			queryConfig["logSet"] = logSet1
		}
		metric1, _ := jsonpath.Get("$[0].metric", v)
		if metric1 != nil && metric1 != "" {
			queryConfig["metric"] = metric1
		}
		metricSet1, _ := jsonpath.Get("$[0].metric_set", v)
		if metricSet1 != nil && metricSet1 != "" {
			queryConfig["metricSet"] = metricSet1
		}
		namespace3, _ := jsonpath.Get("$[0].namespace", v)
		if namespace3 != nil && namespace3 != "" {
			queryConfig["namespace"] = namespace3
		}
		offsetSecs1, _ := jsonpath.Get("$[0].offset_secs", v)
		if offsetSecs1 != nil && offsetSecs1 != "" && cmsAlertRuleV2NonZeroNumber(offsetSecs1) {
			queryConfig["offsetSecs"] = offsetSecs1
		}
		relationType3, _ := jsonpath.Get("$[0].relation_type", v)
		if relationType3 != nil && relationType3 != "" {
			queryConfig["relationType"] = relationType3
		}
		serviceIdList1, _ := jsonpath.Get("$[0].service_id_list", v)
		if serviceIdList1 != nil && serviceIdList1 != "" {
			queryConfig["serviceIdList"] = serviceIdList1
		}
		type7, _ := jsonpath.Get("$[0].type", v)
		if type7 != nil && type7 != "" {
			queryConfig["type"] = type7
		}
		windowSecs1, _ := jsonpath.Get("$[0].window_secs", v)
		if windowSecs1 != nil && windowSecs1 != "" && cmsAlertRuleV2NonZeroNumber(windowSecs1) {
			queryConfig["windowSecs"] = windowSecs1
		}
		if entityFields := cmsAlertRuleV2NestedList(v, "entity_fields", map[string]string{
			"field": "field",
			"value": "value",
		}); len(entityFields) > 0 {
			queryConfig["entityFields"] = entityFields
		}
		if entityFilters := cmsAlertRuleV2NestedList(v, "entity_filters", map[string]string{
			"field":    "field",
			"operator": "operator",
			"value":    "value",
		}); len(entityFilters) > 0 {
			queryConfig["entityFilters"] = entityFilters
		}
		if filterList := cmsAlertRuleV2NestedList(v, "filter_list", map[string]string{
			"key":   "key",
			"type":  "type",
			"value": "value",
		}); len(filterList) > 0 {
			queryConfig["filterList"] = filterList
		}
		if labelFilters := cmsAlertRuleV2NestedList(v, "label_filters", map[string]string{
			"name":     "name",
			"operator": "operator",
			"value":    "value",
		}); len(labelFilters) > 0 {
			queryConfig["labelFilters"] = labelFilters
		}
		if measureList := cmsAlertRuleV2NestedList(v, "measure_list", map[string]string{
			"measure_code": "measureCode",
			"window_secs":  "windowSecs",
			"group_by":     "groupBy",
		}); len(measureList) > 0 {
			queryConfig["measureList"] = measureList
		}
		if groupType, _ := jsonpath.Get("$[0].group_type", v); groupType != nil && groupType != "" {
			queryConfig["groupType"] = groupType
		}
		if queries := cmsAlertRuleV2NestedList(v, "queries", map[string]string{
			"expr":      "expr",
			"time_unit": "timeUnit",
			"start":     "start",
			"end":       "end",
			"window":    "window",
		}); len(queries) > 0 {
			queryConfig["queries"] = queries
		}
		if joinings := cmsAlertRuleV2JoiningList(v); len(joinings) > 0 {
			queryConfig["joinings"] = joinings
		}

		request["queryConfig"] = queryConfig
	}

	rcaConfig := make(map[string]interface{})

	if v := d.Get("rca_config"); !helper.IsNil(v) {
		digitalEmployeeName1, _ := jsonpath.Get("$[0].digital_employee_name", v)
		if digitalEmployeeName1 != nil && digitalEmployeeName1 != "" {
			rcaConfig["digitalEmployeeName"] = digitalEmployeeName1
		}
		enableRca1, _ := jsonpath.Get("$[0].enable_rca", v)
		if enableRca1 != nil && enableRca1 != "" && enableRca1 != false {
			rcaConfig["enableRca"] = enableRca1
		}

		request["rcaConfig"] = rcaConfig
	}

	scheduleConfig := make(map[string]interface{})

	if v := d.Get("schedule_config"); !helper.IsNil(v) {
		intervalSecs1, _ := jsonpath.Get("$[0].interval_secs", v)
		if intervalSecs1 != nil && intervalSecs1 != "" && cmsAlertRuleV2NonZeroNumber(intervalSecs1) {
			scheduleConfig["intervalSecs"] = intervalSecs1
		}
		type9, _ := jsonpath.Get("$[0].type", v)
		if type9 != nil && type9 != "" {
			scheduleConfig["type"] = type9
		}

		request["scheduleConfig"] = scheduleConfig
	}

	if v, ok := d.GetOk("workspace"); ok {
		request["workspace"] = v
	}
	request["action"] = "CREATE"
	body = request
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		response, err = client.RoaPost("Cms", "2024-03-30", action, query, nil, body, true)
		if err != nil {
			if helper.NeedRetry(err) || cmsAlertRuleV2WorkspaceNotReady(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)

	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_cms_alert_rule_v2", action, errs.AlibabaCloudSdkGoERROR)
	}

	id, _ := jsonpath.Get("$.data.alertRule.uuid", response)
	d.SetId(fmt.Sprint(id))

	return resourceAliCloudCmsAlertRuleV2Read(d, meta)
}

func CmsAlertRuleV2ParseQueryRelations(resourceId string, queryConfigRaw map[string]interface{}) ([]map[string]interface{}, []map[string]interface{}, error) {
	queriesMaps := make([]map[string]interface{}, 0)
	if queriesRaw := queryConfigRaw["queries"]; queriesRaw != nil {
		for index, queriesChildRaw := range convert.ConvertToInterfaceArray(queriesRaw) {
			queriesChild, ok := queriesChildRaw.(map[string]interface{})
			if !ok {
				return nil, nil, fmt.Errorf("reading alicloud_cms_alert_rule_v2 (%s): query_config.0.queries.%d is not an object: %#v", resourceId, index, queriesChildRaw)
			}
			queriesMap := make(map[string]interface{})
			queriesMap["expr"] = queriesChild["expr"]
			queriesMap["time_unit"] = queriesChild["timeUnit"]
			queriesMap["start"] = queriesChild["start"]
			queriesMap["end"] = queriesChild["end"]
			queriesMap["window"] = queriesChild["window"]
			queriesMaps = append(queriesMaps, queriesMap)
		}
	}
	joiningsMaps := make([]map[string]interface{}, 0)
	if joiningsRaw := queryConfigRaw["joinings"]; joiningsRaw != nil {
		for index, joiningsChildRaw := range convert.ConvertToInterfaceArray(joiningsRaw) {
			joiningsMap := make(map[string]interface{})
			joiningsChild, ok := joiningsChildRaw.(map[string]interface{})
			if !ok {
				return nil, nil, fmt.Errorf("reading alicloud_cms_alert_rule_v2 (%s): query_config.0.joinings.%d is not an object: %#v", resourceId, index, joiningsChildRaw)
			}
			joiningsMap["type"] = joiningsChild["type"]
			conditionsMaps := make([]map[string]interface{}, 0)
			if conditionsRaw := joiningsChild["conditions"]; conditionsRaw != nil {
				for conditionsIndex, conditionsChildRaw := range convert.ConvertToInterfaceArray(conditionsRaw) {
					conditionsMap := make(map[string]interface{})
					conditionsChild, ok := conditionsChildRaw.(map[string]interface{})
					if !ok {
						return nil, nil, fmt.Errorf("reading alicloud_cms_alert_rule_v2 (%s): query_config.0.joinings.%d.conditions.%d is not an object: %#v", resourceId, index, conditionsIndex, conditionsChildRaw)
					}
					conditionsMap["lhs_field"] = conditionsChild["lhsField"]
					conditionsMap["rhs_field"] = conditionsChild["rhsField"]
					conditionsMap["operator"] = conditionsChild["operator"]
					conditionsMaps = append(conditionsMaps, conditionsMap)
				}
			}
			joiningsMap["conditions"] = conditionsMaps
			joiningsMaps = append(joiningsMaps, joiningsMap)
		}
	}
	return queriesMaps, joiningsMaps, nil
}

func resourceAliCloudCmsAlertRuleV2Read(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cmsServiceV2 := CmsServiceV2{client}

	objectRaw, err := cmsServiceV2.DescribeCmsAlertRuleV2(d.Id())
	if err != nil {
		if !d.IsNewResource() && errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_cms_alert_rule_v2 DescribeCmsAlertRuleV2 Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}

	d.Set("annotations", objectRaw["annotations"])
	d.Set("content_template", objectRaw["contentTemplate"])
	d.Set("created_at", objectRaw["createdAt"])
	d.Set("datasource_type", objectRaw["datasourceType"])
	d.Set("display_name", objectRaw["displayName"])
	d.Set("enabled", objectRaw["enabled"])
	d.Set("labels", cmsAlertRuleV2UserLabels(objectRaw["labels"]))
	d.Set("notify_strategy_id", objectRaw["notifyStrategyId"])
	d.Set("observe_resource_global_scope", objectRaw["observeResourceGlobalScope"])
	d.Set("observe_resource_type", objectRaw["observeResourceType"])
	d.Set("partition_key", objectRaw["partitionKey"])
	d.Set("severity_levels", objectRaw["severityLevels"])
	d.Set("status", objectRaw["status"])
	d.Set("updated_at", objectRaw["updatedAt"])
	d.Set("workspace", objectRaw["workspace"])
	d.Set("alert_rule_v2_id", objectRaw["uuid"])

	actionIntegrationConfigMaps := make([]map[string]interface{}, 0)
	actionIntegrationConfigMap := make(map[string]interface{})
	actionIntegrationConfigRaw := make(map[string]interface{})
	if objectRaw["actionIntegrationConfig"] != nil {
		actionIntegrationConfigRaw = objectRaw["actionIntegrationConfig"].(map[string]interface{})
	}
	if len(actionIntegrationConfigRaw) > 0 {
		actionIntegrationConfigMap["enabled"] = actionIntegrationConfigRaw["enabled"]

		actionsRaw := make([]interface{}, 0)
		if actionIntegrationConfigRaw["actions"] != nil {
			actionsRaw = convert.ConvertToInterfaceArray(actionIntegrationConfigRaw["actions"])
		}

		actionIntegrationConfigMap["actions"] = actionsRaw
		actionIntegrationConfigMaps = append(actionIntegrationConfigMaps, actionIntegrationConfigMap)
	}
	if err := d.Set("action_integration_config", actionIntegrationConfigMaps); err != nil {
		return err
	}
	armsIntegrationConfigMaps := make([]map[string]interface{}, 0)
	armsIntegrationConfigMap := make(map[string]interface{})
	armsIntegrationConfigRaw := make(map[string]interface{})
	if objectRaw["armsIntegrationConfig"] != nil {
		armsIntegrationConfigRaw = objectRaw["armsIntegrationConfig"].(map[string]interface{})
	}
	if len(armsIntegrationConfigRaw) > 0 {
		armsIntegrationConfigMap["enabled"] = armsIntegrationConfigRaw["enabled"]

		armsIntegrationConfigMaps = append(armsIntegrationConfigMaps, armsIntegrationConfigMap)
	}
	if err := d.Set("arms_integration_config", armsIntegrationConfigMaps); err != nil {
		return err
	}
	conditionConfigMaps := make([]map[string]interface{}, 0)
	conditionConfigMap := make(map[string]interface{})
	conditionConfigRaw := make(map[string]interface{})
	if objectRaw["conditionConfig"] != nil {
		conditionConfigRaw = objectRaw["conditionConfig"].(map[string]interface{})
	}
	if len(conditionConfigRaw) > 0 {
		conditionConfigMap["aggregate"] = conditionConfigRaw["aggregate"]
		conditionConfigMap["alert_count"] = conditionConfigRaw["alertCount"]
		conditionConfigMap["count_operator"] = conditionConfigRaw["countOperator"]
		conditionConfigMap["count_threshold"] = conditionConfigRaw["countThreshold"]
		conditionConfigMap["duration_secs"] = conditionConfigRaw["durationSecs"]
		conditionConfigMap["escalation_type"] = conditionConfigRaw["escalationType"]
		conditionConfigMap["legacy_raw"] = conditionConfigRaw["legacyRaw"]
		conditionConfigMap["legacy_type"] = conditionConfigRaw["legacyType"]
		conditionConfigMap["match_field"] = conditionConfigRaw["matchField"]
		conditionConfigMap["match_operator"] = conditionConfigRaw["matchOperator"]
		conditionConfigMap["match_value"] = conditionConfigRaw["matchValue"]
		conditionConfigMap["no_data_alert_level"] = conditionConfigRaw["noDataAlertLevel"]
		conditionConfigMap["no_data_alert_severity"] = conditionConfigRaw["noDataAlertSeverity"]
		conditionConfigMap["no_data_duration_secs"] = conditionConfigRaw["noDataDurationSecs"]
		conditionConfigMap["no_data_policy"] = conditionConfigRaw["noDataPolicy"]
		conditionConfigMap["operator"] = conditionConfigRaw["operator"]
		conditionConfigMap["relation"] = conditionConfigRaw["relation"]
		conditionConfigMap["severity"] = conditionConfigRaw["severity"]
		conditionConfigMap["threshold"] = conditionConfigRaw["threshold"]
		conditionConfigMap["type"] = conditionConfigRaw["type"]
		conditionConfigMap["yoy_time_unit"] = conditionConfigRaw["yoyTimeUnit"]
		conditionConfigMap["yoy_time_value"] = conditionConfigRaw["yoyTimeValue"]

		compareListRaw := conditionConfigRaw["compareList"]
		compareListMaps := make([]map[string]interface{}, 0)
		if compareListRaw != nil {
			for _, compareListChildRaw := range convert.ConvertToInterfaceArray(compareListRaw) {
				compareListMap := make(map[string]interface{})
				compareListChildRaw := compareListChildRaw.(map[string]interface{})
				compareListMap["aggregate"] = compareListChildRaw["aggregate"]
				compareListMap["operator"] = compareListChildRaw["operator"]
				compareListMap["threshold"] = compareListChildRaw["threshold"]
				compareListMap["yoy_time_unit"] = compareListChildRaw["yoyTimeUnit"]
				compareListMap["yoy_time_value"] = compareListChildRaw["yoyTimeValue"]

				compareListMaps = append(compareListMaps, compareListMap)
			}
		}
		conditionConfigMap["compare_list"] = compareListMaps
		compositeEscalationMaps := make([]map[string]interface{}, 0)
		compositeEscalationMap := make(map[string]interface{})
		compositeEscalationRaw := make(map[string]interface{})
		if conditionConfigRaw["compositeEscalation"] != nil {
			compositeEscalationRaw = conditionConfigRaw["compositeEscalation"].(map[string]interface{})
		}
		if len(compositeEscalationRaw) > 0 {
			compositeEscalationMap["relation"] = compositeEscalationRaw["relation"]
			compositeEscalationMap["severity"] = compositeEscalationRaw["severity"]
			compositeEscalationMap["times"] = compositeEscalationRaw["times"]

			escalationsRaw := compositeEscalationRaw["escalations"]
			escalationsMaps := make([]map[string]interface{}, 0)
			if escalationsRaw != nil {
				for _, escalationsChildRaw := range convert.ConvertToInterfaceArray(escalationsRaw) {
					escalationsMap := make(map[string]interface{})
					escalationsChildRaw := escalationsChildRaw.(map[string]interface{})
					escalationsMap["comparison_operator"] = escalationsChildRaw["comparisonOperator"]
					escalationsMap["metric_name"] = escalationsChildRaw["metricName"]
					escalationsMap["period"] = escalationsChildRaw["period"]
					escalationsMap["pre_condition"] = escalationsChildRaw["preCondition"]
					escalationsMap["statistics"] = escalationsChildRaw["statistics"]
					escalationsMap["threshold"] = escalationsChildRaw["threshold"]

					escalationsMaps = append(escalationsMaps, escalationsMap)
				}
			}
			compositeEscalationMap["escalations"] = escalationsMaps
			compositeEscalationMaps = append(compositeEscalationMaps, compositeEscalationMap)
		}
		conditionConfigMap["composite_escalation"] = compositeEscalationMaps
		expressEscalationMaps := make([]map[string]interface{}, 0)
		expressEscalationMap := make(map[string]interface{})
		expressEscalationRaw := make(map[string]interface{})
		if conditionConfigRaw["expressEscalation"] != nil {
			expressEscalationRaw = conditionConfigRaw["expressEscalation"].(map[string]interface{})
		}
		if len(expressEscalationRaw) > 0 {
			expressEscalationMap["raw_expression"] = expressEscalationRaw["rawExpression"]
			expressEscalationMap["severity"] = expressEscalationRaw["severity"]
			expressEscalationMap["times"] = expressEscalationRaw["times"]

			expressEscalationMaps = append(expressEscalationMaps, expressEscalationMap)
		}
		conditionConfigMap["express_escalation"] = expressEscalationMaps
		prometheusMaps := make([]map[string]interface{}, 0)
		prometheusMap := make(map[string]interface{})
		prometheusRaw := make(map[string]interface{})
		if conditionConfigRaw["prometheus"] != nil {
			prometheusRaw = conditionConfigRaw["prometheus"].(map[string]interface{})
		}
		if len(prometheusRaw) > 0 {
			prometheusMap["prom_ql"] = prometheusRaw["promQl"]
			prometheusMap["severity"] = prometheusRaw["severity"]
			prometheusMap["times"] = prometheusRaw["times"]

			prometheusMaps = append(prometheusMaps, prometheusMap)
		}
		conditionConfigMap["prometheus"] = prometheusMaps
		simpleEscalationMaps := make([]map[string]interface{}, 0)
		simpleEscalationMap := make(map[string]interface{})
		simpleEscalationRaw := make(map[string]interface{})
		if conditionConfigRaw["simpleEscalation"] != nil {
			simpleEscalationRaw = conditionConfigRaw["simpleEscalation"].(map[string]interface{})
		}
		if len(simpleEscalationRaw) > 0 {
			simpleEscalationMap["metric_name"] = simpleEscalationRaw["metricName"]
			simpleEscalationMap["period"] = simpleEscalationRaw["period"]

			escalationsRaw := simpleEscalationRaw["escalations"]
			escalationsMaps := make([]map[string]interface{}, 0)
			if escalationsRaw != nil {
				for _, escalationsChildRaw := range convert.ConvertToInterfaceArray(escalationsRaw) {
					escalationsMap := make(map[string]interface{})
					escalationsChildRaw := escalationsChildRaw.(map[string]interface{})
					escalationsMap["comparison_operator"] = escalationsChildRaw["comparisonOperator"]
					escalationsMap["pre_condition"] = escalationsChildRaw["preCondition"]
					escalationsMap["severity"] = escalationsChildRaw["severity"]
					escalationsMap["statistics"] = escalationsChildRaw["statistics"]
					escalationsMap["threshold"] = escalationsChildRaw["threshold"]
					escalationsMap["times"] = escalationsChildRaw["times"]

					escalationsMaps = append(escalationsMaps, escalationsMap)
				}
			}
			simpleEscalationMap["escalations"] = escalationsMaps
			simpleEscalationMaps = append(simpleEscalationMaps, simpleEscalationMap)
		}
		conditionConfigMap["simple_escalation"] = simpleEscalationMaps
		thresholdListRaw := conditionConfigRaw["thresholdList"]
		thresholdListMaps := make([]map[string]interface{}, 0)
		if thresholdListRaw != nil {
			for _, thresholdListChildRaw := range convert.ConvertToInterfaceArray(thresholdListRaw) {
				thresholdListMap := make(map[string]interface{})
				thresholdListChildRaw := thresholdListChildRaw.(map[string]interface{})
				thresholdListMap["severity"] = thresholdListChildRaw["severity"]
				thresholdListMap["threshold"] = thresholdListChildRaw["threshold"]

				thresholdListMaps = append(thresholdListMaps, thresholdListMap)
			}
		}
		conditionConfigMap["threshold_list"] = thresholdListMaps
		triggersRaw := conditionConfigRaw["triggers"]
		triggersMaps := make([]map[string]interface{}, 0)
		if triggersRaw != nil {
			for i, triggersChildRaw := range convert.ConvertToInterfaceArray(triggersRaw) {
				triggersChild, ok := triggersChildRaw.(map[string]interface{})
				if !ok {
					return fmt.Errorf("alicloud_cms_alert_rule_v2 %s: unexpected response type %T for condition_config.triggers[%d]", d.Id(), triggersChildRaw, i)
				}
				triggersMap := make(map[string]interface{})
				triggersMap["operator"] = triggersChild["operator"]
				triggersMap["condition"] = triggersChild["condition"]
				triggersMap["count_condition"] = triggersChild["countCondition"]
				triggersMap["severity"] = triggersChild["severity"]

				triggersMaps = append(triggersMaps, triggersMap)
			}
		}
		conditionConfigMap["triggers"] = triggersMaps
		conditionConfigMaps = append(conditionConfigMaps, conditionConfigMap)
	}
	if err := d.Set("condition_config", conditionConfigMaps); err != nil {
		return err
	}
	datasourceConfigMaps := make([]map[string]interface{}, 0)
	datasourceConfigMap := make(map[string]interface{})
	datasourceConfigRaw := make(map[string]interface{})
	if objectRaw["datasourceConfig"] != nil {
		datasourceConfigRaw = objectRaw["datasourceConfig"].(map[string]interface{})
	}
	if len(datasourceConfigRaw) > 0 {
		datasourceConfigMap["instance_id"] = datasourceConfigRaw["instanceId"]
		datasourceConfigMap["legacy_raw"] = datasourceConfigRaw["legacyRaw"]
		datasourceConfigMap["legacy_type"] = datasourceConfigRaw["legacyType"]
		datasourceConfigMap["product_category"] = datasourceConfigRaw["productCategory"]
		datasourceConfigMap["project"] = datasourceConfigRaw["project"]
		datasourceConfigMap["region_id"] = datasourceConfigRaw["regionId"]
		datasourceConfigMap["type"] = datasourceConfigRaw["type"]

		storesRaw := datasourceConfigRaw["stores"]
		storesMaps := make([]map[string]interface{}, 0)
		if storesRaw != nil {
			for i, storesChildRaw := range convert.ConvertToInterfaceArray(storesRaw) {
				storesChild, ok := storesChildRaw.(map[string]interface{})
				if !ok {
					return fmt.Errorf("alicloud_cms_alert_rule_v2 %s: unexpected response type %T for datasource_config.stores[%d]", d.Id(), storesChildRaw, i)
				}
				storesMap := make(map[string]interface{})
				storesMap["store"] = storesChild["store"]
				storesMap["store_type"] = storesChild["storeType"]

				storesMaps = append(storesMaps, storesMap)
			}
		}
		datasourceConfigMap["stores"] = storesMaps
		datasourceConfigMaps = append(datasourceConfigMaps, datasourceConfigMap)
	}
	if err := d.Set("datasource_config", datasourceConfigMaps); err != nil {
		return err
	}
	notifyConfigMaps := make([]map[string]interface{}, 0)
	notifyConfigMap := make(map[string]interface{})
	notifyConfigRaw := make(map[string]interface{})
	if objectRaw["notifyConfig"] != nil {
		notifyConfigRaw = objectRaw["notifyConfig"].(map[string]interface{})
	}
	if len(notifyConfigRaw) > 0 {
		notifyConfigMap["active_end_time"] = notifyConfigRaw["activeEndTime"]
		notifyConfigMap["active_start_time"] = notifyConfigRaw["activeStartTime"]
		notifyConfigMap["silence_time_secs"] = notifyConfigRaw["silenceTimeSecs"]
		notifyConfigMap["type"] = notifyConfigRaw["type"]
		notifyConfigMap["utc_offset"] = notifyConfigRaw["utcOffset"]

		activeDaysRaw := make([]interface{}, 0)
		if notifyConfigRaw["activeDays"] != nil {
			activeDaysRaw = convert.ConvertToInterfaceArray(notifyConfigRaw["activeDays"])
		}

		notifyConfigMap["active_days"] = activeDaysRaw
		channelsRaw := notifyConfigRaw["channels"]
		channelsMaps := make([]map[string]interface{}, 0)
		if channelsRaw != nil {
			for _, channelsChildRaw := range convert.ConvertToInterfaceArray(channelsRaw) {
				channelsMap := make(map[string]interface{})
				channelsChildRaw := channelsChildRaw.(map[string]interface{})
				channelsMap["type"] = channelsChildRaw["type"]

				identifiersRaw := make([]interface{}, 0)
				if channelsChildRaw["identifiers"] != nil {
					identifiersRaw = convert.ConvertToInterfaceArray(channelsChildRaw["identifiers"])
				}

				channelsMap["identifiers"] = identifiersRaw
				channelsMaps = append(channelsMaps, channelsMap)
			}
		}
		notifyConfigMap["channels"] = channelsMaps
		notifyStrategiesRaw := make([]interface{}, 0)
		if notifyConfigRaw["notifyStrategies"] != nil {
			notifyStrategiesRaw = convert.ConvertToInterfaceArray(notifyConfigRaw["notifyStrategies"])
		}

		notifyConfigMap["notify_strategies"] = notifyStrategiesRaw
		notifyConfigMaps = append(notifyConfigMaps, notifyConfigMap)
	}
	if err := d.Set("notify_config", notifyConfigMaps); err != nil {
		return err
	}
	observeResourceConfigMaps := make([]map[string]interface{}, 0)
	observeResourceConfigMap := make(map[string]interface{})
	observeResourceConfigRaw := make(map[string]interface{})
	if objectRaw["observeResourceConfig"] != nil {
		raw, ok := objectRaw["observeResourceConfig"].(map[string]interface{})
		if !ok {
			return fmt.Errorf("alicloud_cms_alert_rule_v2 %s: unexpected response type %T for observe_resource_config", d.Id(), objectRaw["observeResourceConfig"])
		}
		observeResourceConfigRaw = raw
	}
	if len(observeResourceConfigRaw) > 0 {
		observeResourceConfigMap["entity_domain"] = observeResourceConfigRaw["entityDomain"]
		observeResourceConfigMap["entity_type"] = observeResourceConfigRaw["entityType"]
		observeResourceConfigMap["namespace"] = observeResourceConfigRaw["namespace"]
		observeResourceConfigMap["product_category"] = observeResourceConfigRaw["productCategory"]
		observeResourceConfigMap["relation_type"] = observeResourceConfigRaw["relationType"]

		resourcesRaw := make([]interface{}, 0)
		if observeResourceConfigRaw["resources"] != nil {
			resourcesRaw = convert.ConvertToInterfaceArray(observeResourceConfigRaw["resources"])
		}

		observeResourceConfigMap["resources"] = resourcesRaw
		observeResourceConfigMaps = append(observeResourceConfigMaps, observeResourceConfigMap)
	}
	if err := d.Set("observe_resource_config", observeResourceConfigMaps); err != nil {
		return err
	}
	queryConfigMaps := make([]map[string]interface{}, 0)
	queryConfigMap := make(map[string]interface{})
	queryConfigRaw := make(map[string]interface{})
	if objectRaw["queryConfig"] != nil {
		queryConfigRaw = objectRaw["queryConfig"].(map[string]interface{})
	}
	if len(queryConfigRaw) > 0 {
		queryConfigMap["aggregate"] = queryConfigRaw["aggregate"]
		queryConfigMap["enable_data_complete_check"] = queryConfigRaw["enableDataCompleteCheck"]
		queryConfigMap["entity_domain"] = queryConfigRaw["entityDomain"]
		queryConfigMap["entity_type"] = queryConfigRaw["entityType"]
		queryConfigMap["expr"] = queryConfigRaw["expr"]
		queryConfigMap["group_id"] = queryConfigRaw["groupId"]
		queryConfigMap["legacy_raw"] = queryConfigRaw["legacyRaw"]
		queryConfigMap["legacy_type"] = queryConfigRaw["legacyType"]
		queryConfigMap["log_set"] = queryConfigRaw["logSet"]
		queryConfigMap["metric"] = queryConfigRaw["metric"]
		queryConfigMap["metric_set"] = queryConfigRaw["metricSet"]
		queryConfigMap["namespace"] = queryConfigRaw["namespace"]
		queryConfigMap["offset_secs"] = queryConfigRaw["offsetSecs"]
		queryConfigMap["prom_ql"] = queryConfigRaw["promQl"]
		queryConfigMap["relation_type"] = queryConfigRaw["relationType"]
		queryConfigMap["type"] = queryConfigRaw["type"]
		queryConfigMap["window_secs"] = queryConfigRaw["windowSecs"]

		queryConfigMap["dimensions"] = queryConfigRaw["dimensions"]
		entityFieldsRaw := queryConfigRaw["entityFields"]
		entityFieldsMaps := make([]map[string]interface{}, 0)
		if entityFieldsRaw != nil {
			for _, entityFieldsChildRaw := range convert.ConvertToInterfaceArray(entityFieldsRaw) {
				entityFieldsMap := make(map[string]interface{})
				entityFieldsChildRaw := entityFieldsChildRaw.(map[string]interface{})
				entityFieldsMap["field"] = entityFieldsChildRaw["field"]
				entityFieldsMap["value"] = entityFieldsChildRaw["value"]

				entityFieldsMaps = append(entityFieldsMaps, entityFieldsMap)
			}
		}
		queryConfigMap["entity_fields"] = entityFieldsMaps
		entityFiltersRaw := queryConfigRaw["entityFilters"]
		entityFiltersMaps := make([]map[string]interface{}, 0)
		if entityFiltersRaw != nil {
			for _, entityFiltersChildRaw := range convert.ConvertToInterfaceArray(entityFiltersRaw) {
				entityFiltersMap := make(map[string]interface{})
				entityFiltersChildRaw := entityFiltersChildRaw.(map[string]interface{})
				entityFiltersMap["field"] = entityFiltersChildRaw["field"]
				entityFiltersMap["operator"] = entityFiltersChildRaw["operator"]
				entityFiltersMap["value"] = entityFiltersChildRaw["value"]

				entityFiltersMaps = append(entityFiltersMaps, entityFiltersMap)
			}
		}
		queryConfigMap["entity_filters"] = entityFiltersMaps
		filterListRaw := queryConfigRaw["filterList"]
		filterListMaps := make([]map[string]interface{}, 0)
		if filterListRaw != nil {
			for _, filterListChildRaw := range convert.ConvertToInterfaceArray(filterListRaw) {
				filterListMap := make(map[string]interface{})
				filterListChildRaw := filterListChildRaw.(map[string]interface{})
				filterListMap["key"] = filterListChildRaw["key"]
				filterListMap["type"] = filterListChildRaw["type"]
				filterListMap["value"] = filterListChildRaw["value"]

				filterListMaps = append(filterListMaps, filterListMap)
			}
		}
		queryConfigMap["filter_list"] = filterListMaps
		labelFiltersRaw := queryConfigRaw["labelFilters"]
		labelFiltersMaps := make([]map[string]interface{}, 0)
		if labelFiltersRaw != nil {
			for _, labelFiltersChildRaw := range convert.ConvertToInterfaceArray(labelFiltersRaw) {
				labelFiltersMap := make(map[string]interface{})
				labelFiltersChildRaw := labelFiltersChildRaw.(map[string]interface{})
				labelFiltersMap["name"] = labelFiltersChildRaw["name"]
				labelFiltersMap["operator"] = labelFiltersChildRaw["operator"]
				labelFiltersMap["value"] = labelFiltersChildRaw["value"]

				labelFiltersMaps = append(labelFiltersMaps, labelFiltersMap)
			}
		}
		queryConfigMap["label_filters"] = labelFiltersMaps
		measureListRaw := queryConfigRaw["measureList"]
		measureListMaps := make([]map[string]interface{}, 0)
		if measureListRaw != nil {
			for _, measureListChildRaw := range convert.ConvertToInterfaceArray(measureListRaw) {
				measureListMap := make(map[string]interface{})
				measureListChildRaw := measureListChildRaw.(map[string]interface{})
				measureListMap["measure_code"] = measureListChildRaw["measureCode"]
				measureListMap["window_secs"] = measureListChildRaw["windowSecs"]

				groupByRaw := make([]interface{}, 0)
				if measureListChildRaw["groupBy"] != nil {
					groupByRaw = convert.ConvertToInterfaceArray(measureListChildRaw["groupBy"])
				}

				measureListMap["group_by"] = groupByRaw
				measureListMaps = append(measureListMaps, measureListMap)
			}
		}
		queryConfigMap["measure_list"] = measureListMaps
		queryConfigMap["group_type"] = queryConfigRaw["groupType"]
		queriesMaps, joiningsMaps, err := CmsAlertRuleV2ParseQueryRelations(d.Id(), queryConfigRaw)
		if err != nil {
			return err
		}
		queryConfigMap["queries"] = queriesMaps
		queryConfigMap["joinings"] = joiningsMaps
		serviceIdListRaw := make([]interface{}, 0)
		if queryConfigRaw["serviceIdList"] != nil {
			serviceIdListRaw = convert.ConvertToInterfaceArray(queryConfigRaw["serviceIdList"])
		}

		queryConfigMap["service_id_list"] = serviceIdListRaw
		queryConfigMaps = append(queryConfigMaps, queryConfigMap)
	}
	if err := d.Set("query_config", queryConfigMaps); err != nil {
		return err
	}
	rcaConfigMaps := make([]map[string]interface{}, 0)
	rcaConfigMap := make(map[string]interface{})
	rcaConfigRaw := make(map[string]interface{})
	if objectRaw["rcaConfig"] != nil {
		raw, ok := objectRaw["rcaConfig"].(map[string]interface{})
		if !ok {
			return fmt.Errorf("alicloud_cms_alert_rule_v2 %s: unexpected response type %T for rca_config", d.Id(), objectRaw["rcaConfig"])
		}
		rcaConfigRaw = raw
	}
	if len(rcaConfigRaw) > 0 {
		rcaConfigMap["digital_employee_name"] = rcaConfigRaw["digitalEmployeeName"]
		rcaConfigMap["enable_rca"] = rcaConfigRaw["enableRca"]

		rcaConfigMaps = append(rcaConfigMaps, rcaConfigMap)
	}
	if err := d.Set("rca_config", rcaConfigMaps); err != nil {
		return err
	}
	scheduleConfigMaps := make([]map[string]interface{}, 0)
	scheduleConfigMap := make(map[string]interface{})
	scheduleConfigRaw := make(map[string]interface{})
	if objectRaw["scheduleConfig"] != nil {
		scheduleConfigRaw = objectRaw["scheduleConfig"].(map[string]interface{})
	}
	if len(scheduleConfigRaw) > 0 {
		scheduleConfigMap["interval_secs"] = scheduleConfigRaw["intervalSecs"]
		scheduleConfigMap["type"] = scheduleConfigRaw["type"]

		scheduleConfigMaps = append(scheduleConfigMaps, scheduleConfigMap)
	}
	if err := d.Set("schedule_config", scheduleConfigMaps); err != nil {
		return err
	}

	d.Set("alert_rule_v2_id", d.Id())

	return nil
}

func resourceAliCloudCmsAlertRuleV2Update(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	var body map[string]interface{}
	update := false

	var err error
	action := fmt.Sprintf("/manageAlertRules")
	request = make(map[string]interface{})
	query = make(map[string]*string)
	body = make(map[string]interface{})
	request["uuid"] = d.Id()
	query["RegionId"] = helper.StringPointer(client.RegionId)
	if d.HasChange("action_integration_config") {
		update = true
	}
	actionIntegrationConfig := make(map[string]interface{})

	if v := d.Get("action_integration_config"); !helper.IsNil(v) || d.HasChange("action_integration_config") {
		actions1, _ := jsonpath.Get("$[0].actions", v)
		if actions1 != nil && actions1 != "" {
			actionIntegrationConfig["actions"] = actions1
		}
		enabled1, _ := jsonpath.Get("$[0].enabled", v)
		if enabled1 != nil && enabled1 != "" {
			actionIntegrationConfig["enabled"] = enabled1
		}

		request["actionIntegrationConfig"] = actionIntegrationConfig
	}

	if d.HasChange("annotations") {
		update = true
	}
	if v, ok := d.GetOk("annotations"); ok || d.HasChange("annotations") {
		request["annotations"] = v
	}
	if d.HasChange("arms_integration_config") {
		update = true
	}
	armsIntegrationConfig := make(map[string]interface{})

	if v := d.Get("arms_integration_config"); !helper.IsNil(v) || d.HasChange("arms_integration_config") {
		enabled3, _ := jsonpath.Get("$[0].enabled", v)
		if enabled3 != nil && enabled3 != "" {
			armsIntegrationConfig["enabled"] = enabled3
		}

		request["armsIntegrationConfig"] = armsIntegrationConfig
	}

	if d.HasChange("condition_config") {
		update = true
	}
	conditionConfig := make(map[string]interface{})

	if v := d.Get("condition_config"); !helper.IsNil(v) || d.HasChange("condition_config") {
		aggregate1, _ := jsonpath.Get("$[0].aggregate", v)
		if aggregate1 != nil && aggregate1 != "" {
			conditionConfig["aggregate"] = aggregate1
		}
		compositeEscalation := make(map[string]interface{})
		relation1, _ := jsonpath.Get("$[0].composite_escalation[0].relation", d.Get("condition_config"))
		if relation1 != nil && relation1 != "" {
			compositeEscalation["relation"] = relation1
		}
		severity1, _ := jsonpath.Get("$[0].composite_escalation[0].severity", d.Get("condition_config"))
		if severity1 != nil && severity1 != "" {
			compositeEscalation["severity"] = severity1
		}
		times1, _ := jsonpath.Get("$[0].composite_escalation[0].times", d.Get("condition_config"))
		if times1 != nil && times1 != "" {
			compositeEscalation["times"] = times1
		}

		if escalations := cmsAlertRuleV2NestedList(d.Get("condition_config"), "composite_escalation[0].escalations", map[string]string{
			"comparison_operator": "comparisonOperator",
			"metric_name":         "metricName",
			"statistics":          "statistics",
			"threshold":           "threshold",
			"pre_condition":       "preCondition",
			"period":              "period",
		}, "threshold"); len(escalations) > 0 {
			compositeEscalation["escalations"] = escalations
		}

		if len(compositeEscalation) > 0 {
			conditionConfig["compositeEscalation"] = compositeEscalation
		}
		countOperator1, _ := jsonpath.Get("$[0].count_operator", v)
		if countOperator1 != nil && countOperator1 != "" {
			conditionConfig["countOperator"] = countOperator1
		}
		countThreshold1, _ := jsonpath.Get("$[0].count_threshold", v)
		if countThreshold1 != nil && countThreshold1 != "" && cmsAlertRuleV2NonZeroNumber(countThreshold1) {
			conditionConfig["countThreshold"] = countThreshold1
		}
		durationSecs1, _ := jsonpath.Get("$[0].duration_secs", v)
		if durationSecs1 != nil && durationSecs1 != "" && cmsAlertRuleV2NonZeroNumber(durationSecs1) {
			conditionConfig["durationSecs"] = durationSecs1
		}
		escalationType1, _ := jsonpath.Get("$[0].escalation_type", v)
		if escalationType1 != nil && escalationType1 != "" {
			conditionConfig["escalationType"] = escalationType1
		}
		matchField1, _ := jsonpath.Get("$[0].match_field", v)
		if matchField1 != nil && matchField1 != "" {
			conditionConfig["matchField"] = matchField1
		}
		matchOperator1, _ := jsonpath.Get("$[0].match_operator", v)
		if matchOperator1 != nil && matchOperator1 != "" {
			conditionConfig["matchOperator"] = matchOperator1
		}
		matchValue1, _ := jsonpath.Get("$[0].match_value", v)
		if matchValue1 != nil && matchValue1 != "" {
			conditionConfig["matchValue"] = matchValue1
		}
		noDataAlertSeverity1, _ := jsonpath.Get("$[0].no_data_alert_severity", v)
		if noDataAlertSeverity1 != nil && noDataAlertSeverity1 != "" {
			conditionConfig["noDataAlertSeverity"] = noDataAlertSeverity1
		}
		noDataDurationSecs1, _ := jsonpath.Get("$[0].no_data_duration_secs", v)
		if noDataDurationSecs1 != nil && noDataDurationSecs1 != "" && cmsAlertRuleV2NonZeroNumber(noDataDurationSecs1) {
			conditionConfig["noDataDurationSecs"] = noDataDurationSecs1
		}
		noDataPolicy1, _ := jsonpath.Get("$[0].no_data_policy", v)
		if noDataPolicy1 != nil && noDataPolicy1 != "" {
			conditionConfig["noDataPolicy"] = noDataPolicy1
		}
		operator1, _ := jsonpath.Get("$[0].operator", v)
		if operator1 != nil && operator1 != "" {
			conditionConfig["operator"] = operator1
		}
		relation3, _ := jsonpath.Get("$[0].relation", v)
		if relation3 != nil && relation3 != "" {
			conditionConfig["relation"] = relation3
		}
		severity7, _ := jsonpath.Get("$[0].severity", v)
		if severity7 != nil && severity7 != "" {
			conditionConfig["severity"] = severity7
		}
		simpleEscalation := make(map[string]interface{})
		metricName1, _ := jsonpath.Get("$[0].simple_escalation[0].metric_name", d.Get("condition_config"))
		if metricName1 != nil && metricName1 != "" {
			simpleEscalation["metricName"] = metricName1
		}
		period1, _ := jsonpath.Get("$[0].simple_escalation[0].period", d.Get("condition_config"))
		if period1 != nil && period1 != "" {
			simpleEscalation["period"] = period1
		}

		if escalations := cmsAlertRuleV2NestedList(d.Get("condition_config"), "simple_escalation[0].escalations", map[string]string{
			"comparison_operator": "comparisonOperator",
			"statistics":          "statistics",
			"threshold":           "threshold",
			"severity":            "severity",
			"pre_condition":       "preCondition",
			"times":               "times",
		}, "threshold"); len(escalations) > 0 {
			simpleEscalation["escalations"] = escalations
		}

		if len(simpleEscalation) > 0 {
			conditionConfig["simpleEscalation"] = simpleEscalation
		}
		threshold1, _ := jsonpath.Get("$[0].threshold", v)
		if threshold1 != nil && threshold1 != "" && cmsAlertRuleV2NonZeroNumber(threshold1) {
			conditionConfig["threshold"] = threshold1
		}
		type1, _ := jsonpath.Get("$[0].type", v)
		if type1 != nil && type1 != "" {
			conditionConfig["type"] = type1
		}
		yoyTimeUnit1, _ := jsonpath.Get("$[0].yoy_time_unit", v)
		if yoyTimeUnit1 != nil && yoyTimeUnit1 != "" {
			conditionConfig["yoyTimeUnit"] = yoyTimeUnit1
		}
		yoyTimeValue1, _ := jsonpath.Get("$[0].yoy_time_value", v)
		if yoyTimeValue1 != nil && yoyTimeValue1 != "" && cmsAlertRuleV2NonZeroNumber(yoyTimeValue1) {
			conditionConfig["yoyTimeValue"] = yoyTimeValue1
		}
		if compareList := cmsAlertRuleV2NestedList(v, "compare_list", map[string]string{
			"aggregate":      "aggregate",
			"operator":       "operator",
			"threshold":      "threshold",
			"yoy_time_unit":  "yoyTimeUnit",
			"yoy_time_value": "yoyTimeValue",
		}, "threshold"); len(compareList) > 0 {
			conditionConfig["compareList"] = compareList
		}
		if thresholdList := cmsAlertRuleV2NestedList(v, "threshold_list", map[string]string{
			"severity":  "severity",
			"threshold": "threshold",
		}, "threshold"); len(thresholdList) > 0 {
			conditionConfig["thresholdList"] = thresholdList
		}
		if triggers := cmsAlertRuleV2NestedList(v, "triggers", map[string]string{
			"operator":        "operator",
			"condition":       "condition",
			"count_condition": "countCondition",
			"severity":        "severity",
		}); len(triggers) > 0 {
			conditionConfig["triggers"] = triggers
		}
		if alertCount1, _ := jsonpath.Get("$[0].alert_count", v); alertCount1 != nil && cmsAlertRuleV2NonZeroNumber(alertCount1) {
			conditionConfig["alertCount"] = alertCount1
		}
		if noDataAlertLevel1, _ := jsonpath.Get("$[0].no_data_alert_level", v); noDataAlertLevel1 != nil && noDataAlertLevel1 != "" {
			conditionConfig["noDataAlertLevel"] = noDataAlertLevel1
		}

		request["conditionConfig"] = conditionConfig
	}

	if d.HasChange("content_template") {
		update = true
	}
	if v, ok := d.GetOk("content_template"); ok || d.HasChange("content_template") {
		request["contentTemplate"] = v
	}
	if d.HasChange("datasource_config") {
		update = true
	}
	datasourceConfig := make(map[string]interface{})

	if v := d.Get("datasource_config"); !helper.IsNil(v) || d.HasChange("datasource_config") {
		instanceId1, _ := jsonpath.Get("$[0].instance_id", v)
		if instanceId1 != nil && instanceId1 != "" {
			datasourceConfig["instanceId"] = instanceId1
		}
		productCategory1, _ := jsonpath.Get("$[0].product_category", v)
		if productCategory1 != nil && productCategory1 != "" {
			datasourceConfig["productCategory"] = productCategory1
		}
		project1, _ := jsonpath.Get("$[0].project", v)
		if project1 != nil && project1 != "" {
			datasourceConfig["project"] = project1
		}
		regionId1, _ := jsonpath.Get("$[0].region_id", v)
		if regionId1 != nil && regionId1 != "" {
			datasourceConfig["regionId"] = regionId1
		}
		type3, _ := jsonpath.Get("$[0].type", v)
		if type3 != nil && type3 != "" {
			datasourceConfig["type"] = type3
		}
		stores3, _ := jsonpath.Get("$[0].stores", d.Get("datasource_config"))
		if stores3 != nil {
			stores := make([]map[string]interface{}, 0)
			for _, rawStore := range convert.ConvertToInterfaceArray(stores3) {
				storeRaw, ok := rawStore.(map[string]interface{})
				if !ok {
					continue
				}
				store := make(map[string]interface{})
				if v, ok := storeRaw["store"].(string); ok && v != "" {
					store["store"] = v
				}
				if v, ok := storeRaw["store_type"].(string); ok && v != "" {
					store["storeType"] = v
				}
				if len(store) > 0 {
					stores = append(stores, store)
				}
			}
			if len(stores) > 0 {
				datasourceConfig["stores"] = stores
			}
		}

		request["datasourceConfig"] = datasourceConfig
	}

	if d.HasChange("display_name") {
		update = true
	}
	if v, ok := d.GetOk("display_name"); ok || d.HasChange("display_name") {
		request["displayName"] = v
	}
	if d.HasChange("enabled") {
		update = true
	}
	if v, ok := d.GetOkExists("enabled"); ok || d.HasChange("enabled") {
		request["enabled"] = v
	}
	if d.HasChange("labels") {
		update = true
	}
	if v, ok := d.GetOk("labels"); ok || d.HasChange("labels") {
		request["labels"] = v
	}
	if d.HasChange("notify_config") {
		update = true
	}
	notifyConfig := make(map[string]interface{})

	if v := d.Get("notify_config"); !helper.IsNil(v) || d.HasChange("notify_config") {
		activeDays1, _ := jsonpath.Get("$[0].active_days", v)
		if activeDays1 != nil && activeDays1 != "" {
			notifyConfig["activeDays"] = activeDays1
		}
		activeEndTime1, _ := jsonpath.Get("$[0].active_end_time", v)
		if activeEndTime1 != nil && activeEndTime1 != "" {
			notifyConfig["activeEndTime"] = activeEndTime1
		}
		activeStartTime1, _ := jsonpath.Get("$[0].active_start_time", v)
		if activeStartTime1 != nil && activeStartTime1 != "" {
			notifyConfig["activeStartTime"] = activeStartTime1
		}
		notifyStrategies1, _ := jsonpath.Get("$[0].notify_strategies", v)
		if notifyStrategies1 != nil && notifyStrategies1 != "" {
			notifyConfig["notifyStrategies"] = notifyStrategies1
		}
		silenceTimeSecs1, _ := jsonpath.Get("$[0].silence_time_secs", v)
		if silenceTimeSecs1 != nil && silenceTimeSecs1 != "" && cmsAlertRuleV2NonZeroNumber(silenceTimeSecs1) {
			notifyConfig["silenceTimeSecs"] = silenceTimeSecs1
		}
		type5, _ := jsonpath.Get("$[0].type", v)
		if type5 != nil && type5 != "" {
			notifyConfig["type"] = type5
		}
		utcOffset1, _ := jsonpath.Get("$[0].utc_offset", v)
		if utcOffset1 != nil && utcOffset1 != "" {
			notifyConfig["utcOffset"] = utcOffset1
		}
		if channels := cmsAlertRuleV2NestedList(v, "channels", map[string]string{
			"type":        "type",
			"identifiers": "identifiers",
		}); len(channels) > 0 {
			notifyConfig["channels"] = channels
		}

		request["notifyConfig"] = notifyConfig
	}

	if d.HasChange("observe_resource_config") {
		update = true
	}
	observeResourceConfig := make(map[string]interface{})

	if v := d.Get("observe_resource_config"); !helper.IsNil(v) || d.HasChange("observe_resource_config") {
		entityDomain1, _ := jsonpath.Get("$[0].entity_domain", v)
		if entityDomain1 != nil && entityDomain1 != "" {
			observeResourceConfig["entityDomain"] = entityDomain1
		}
		entityType1, _ := jsonpath.Get("$[0].entity_type", v)
		if entityType1 != nil && entityType1 != "" {
			observeResourceConfig["entityType"] = entityType1
		}
		namespace1, _ := jsonpath.Get("$[0].namespace", v)
		if namespace1 != nil && namespace1 != "" {
			observeResourceConfig["namespace"] = namespace1
		}
		productCategory3, _ := jsonpath.Get("$[0].product_category", v)
		if productCategory3 != nil && productCategory3 != "" {
			observeResourceConfig["productCategory"] = productCategory3
		}
		relationType1, _ := jsonpath.Get("$[0].relation_type", v)
		if relationType1 != nil && relationType1 != "" {
			observeResourceConfig["relationType"] = relationType1
		}
		resources1, _ := jsonpath.Get("$[0].resources", v)
		if resources1 != nil && resources1 != "" {
			observeResourceConfig["resources"] = resources1
		}

		request["observeResourceConfig"] = observeResourceConfig
	}

	if d.HasChange("query_config") {
		update = true
	}
	queryConfig := make(map[string]interface{})

	if v := d.Get("query_config"); !helper.IsNil(v) || d.HasChange("query_config") {
		aggregate3, _ := jsonpath.Get("$[0].aggregate", v)
		if aggregate3 != nil && aggregate3 != "" {
			queryConfig["aggregate"] = aggregate3
		}
		dimensions1, _ := jsonpath.Get("$[0].dimensions", v)
		if dimensions1 != nil && dimensions1 != "" {
			queryConfig["dimensions"] = dimensions1
		}
		enableDataCompleteCheck1, _ := jsonpath.Get("$[0].enable_data_complete_check", v)
		if enableDataCompleteCheck1 != nil && enableDataCompleteCheck1 != "" && enableDataCompleteCheck1 != false {
			queryConfig["enableDataCompleteCheck"] = enableDataCompleteCheck1
		}
		entityDomain3, _ := jsonpath.Get("$[0].entity_domain", v)
		if entityDomain3 != nil && entityDomain3 != "" {
			queryConfig["entityDomain"] = entityDomain3
		}
		entityType3, _ := jsonpath.Get("$[0].entity_type", v)
		if entityType3 != nil && entityType3 != "" {
			queryConfig["entityType"] = entityType3
		}
		expr1, _ := jsonpath.Get("$[0].expr", v)
		if expr1 != nil && expr1 != "" {
			queryConfig["expr"] = expr1
		}
		groupId1, _ := jsonpath.Get("$[0].group_id", v)
		if groupId1 != nil && groupId1 != "" {
			queryConfig["groupId"] = groupId1
		}
		logSet1, _ := jsonpath.Get("$[0].log_set", v)
		if logSet1 != nil && logSet1 != "" {
			queryConfig["logSet"] = logSet1
		}
		metric1, _ := jsonpath.Get("$[0].metric", v)
		if metric1 != nil && metric1 != "" {
			queryConfig["metric"] = metric1
		}
		metricSet1, _ := jsonpath.Get("$[0].metric_set", v)
		if metricSet1 != nil && metricSet1 != "" {
			queryConfig["metricSet"] = metricSet1
		}
		namespace3, _ := jsonpath.Get("$[0].namespace", v)
		if namespace3 != nil && namespace3 != "" {
			queryConfig["namespace"] = namespace3
		}
		offsetSecs1, _ := jsonpath.Get("$[0].offset_secs", v)
		if offsetSecs1 != nil && offsetSecs1 != "" && cmsAlertRuleV2NonZeroNumber(offsetSecs1) {
			queryConfig["offsetSecs"] = offsetSecs1
		}
		relationType3, _ := jsonpath.Get("$[0].relation_type", v)
		if relationType3 != nil && relationType3 != "" {
			queryConfig["relationType"] = relationType3
		}
		serviceIdList1, _ := jsonpath.Get("$[0].service_id_list", v)
		if serviceIdList1 != nil && serviceIdList1 != "" {
			queryConfig["serviceIdList"] = serviceIdList1
		}
		type7, _ := jsonpath.Get("$[0].type", v)
		if type7 != nil && type7 != "" {
			queryConfig["type"] = type7
		}
		windowSecs1, _ := jsonpath.Get("$[0].window_secs", v)
		if windowSecs1 != nil && windowSecs1 != "" && cmsAlertRuleV2NonZeroNumber(windowSecs1) {
			queryConfig["windowSecs"] = windowSecs1
		}
		if entityFields := cmsAlertRuleV2NestedList(v, "entity_fields", map[string]string{
			"field": "field",
			"value": "value",
		}); len(entityFields) > 0 {
			queryConfig["entityFields"] = entityFields
		}
		if entityFilters := cmsAlertRuleV2NestedList(v, "entity_filters", map[string]string{
			"field":    "field",
			"operator": "operator",
			"value":    "value",
		}); len(entityFilters) > 0 {
			queryConfig["entityFilters"] = entityFilters
		}
		if filterList := cmsAlertRuleV2NestedList(v, "filter_list", map[string]string{
			"key":   "key",
			"type":  "type",
			"value": "value",
		}); len(filterList) > 0 {
			queryConfig["filterList"] = filterList
		}
		if labelFilters := cmsAlertRuleV2NestedList(v, "label_filters", map[string]string{
			"name":     "name",
			"operator": "operator",
			"value":    "value",
		}); len(labelFilters) > 0 {
			queryConfig["labelFilters"] = labelFilters
		}
		if measureList := cmsAlertRuleV2NestedList(v, "measure_list", map[string]string{
			"measure_code": "measureCode",
			"window_secs":  "windowSecs",
			"group_by":     "groupBy",
		}); len(measureList) > 0 {
			queryConfig["measureList"] = measureList
		}
		if groupType, _ := jsonpath.Get("$[0].group_type", v); groupType != nil && groupType != "" {
			queryConfig["groupType"] = groupType
		}
		if queries := cmsAlertRuleV2NestedList(v, "queries", map[string]string{
			"expr":      "expr",
			"time_unit": "timeUnit",
			"start":     "start",
			"end":       "end",
			"window":    "window",
		}); len(queries) > 0 {
			queryConfig["queries"] = queries
		}
		if joinings := cmsAlertRuleV2JoiningList(v); len(joinings) > 0 {
			queryConfig["joinings"] = joinings
		}

		request["queryConfig"] = queryConfig
	}

	rcaConfig := make(map[string]interface{})

	if v := d.Get("rca_config"); !helper.IsNil(v) {
		digitalEmployeeName1, _ := jsonpath.Get("$[0].digital_employee_name", v)
		if digitalEmployeeName1 != nil && digitalEmployeeName1 != "" {
			rcaConfig["digitalEmployeeName"] = digitalEmployeeName1
		}
		enableRca1, _ := jsonpath.Get("$[0].enable_rca", v)
		if enableRca1 != nil && enableRca1 != "" && enableRca1 != false {
			rcaConfig["enableRca"] = enableRca1
		}

		request["rcaConfig"] = rcaConfig
	}

	if d.HasChange("rca_config") {
		update = true
	}
	if d.HasChange("schedule_config") {
		update = true
	}
	scheduleConfig := make(map[string]interface{})

	if v := d.Get("schedule_config"); !helper.IsNil(v) || d.HasChange("schedule_config") {
		intervalSecs1, _ := jsonpath.Get("$[0].interval_secs", v)
		if intervalSecs1 != nil && intervalSecs1 != "" && cmsAlertRuleV2NonZeroNumber(intervalSecs1) {
			scheduleConfig["intervalSecs"] = intervalSecs1
		}
		type9, _ := jsonpath.Get("$[0].type", v)
		if type9 != nil && type9 != "" {
			scheduleConfig["type"] = type9
		}

		request["scheduleConfig"] = scheduleConfig
	}

	if d.HasChange("workspace") {
		update = true
	}
	if v, ok := d.GetOk("workspace"); ok || d.HasChange("workspace") {
		request["workspace"] = v
	}
	request["action"] = "UPDATE"
	body = request
	if update {
		wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RoaPost("Cms", "2024-03-30", action, query, nil, body, true)
			if err != nil {
				if helper.NeedRetry(err) || cmsAlertRuleV2WorkspaceNotReady(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		logs.AddDebug(action, response, request)
		if err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
		}
	}

	return resourceAliCloudCmsAlertRuleV2Read(d, meta)
}

func resourceAliCloudCmsAlertRuleV2Delete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	action := fmt.Sprintf("/manageAlertRules")
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	body := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["uuid"] = d.Id()
	query["RegionId"] = helper.StringPointer(client.RegionId)

	if v, ok := d.GetOk("workspace"); ok {
		request["workspace"] = v
	}
	request["action"] = "DELETE"
	body = request
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		response, err = client.RoaPost("Cms", "2024-03-30", action, query, nil, body, true)
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
		if errs.NotFoundError(err) {
			return nil
		}
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
	}

	return nil
}

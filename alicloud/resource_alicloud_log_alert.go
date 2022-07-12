package alicloud

import (
	"fmt"
	"time"

	sls "github.com/aliyun/aliyun-log-go-sdk"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAlicloudLogAlert() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudLogAlertCreate,
		Read:   resourceAlicloudLogAlertRead,
		Update: resourceAlicloudLogAlertUpdate,
		Delete: resourceAlicloudLogAlertDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"version": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"project_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"alert_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"alert_displayname": {
				Type:     schema.TypeString,
				Required: true,
			},
			"alert_description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"condition": {
				Type:       schema.TypeString,
				Optional:   true,
				Deprecated: "Deprecated from 1.161.0+, use eval_condition in severity_configurations",
			},
			"dashboard": {
				Type:       schema.TypeString,
				Optional:   true,
				Deprecated: "Deprecated from 1.161.0+, use dashboardId in query_list",
			},
			"mute_until": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntAtLeast(0),
				Default:      time.Now().Unix(),
			},
			"throttling": {
				Type:       schema.TypeString,
				Optional:   true,
				Deprecated: "Deprecated from 1.161.0+, use repeat_interval in policy_configuration",
			},
			"notify_threshold": {
				Type:       schema.TypeInt,
				Optional:   true,
				Deprecated: "Deprecated from 1.161.0+, use threshold",
			},
			"threshold": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"no_data_fire": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"no_data_severity": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntInSlice([]int{2, 4, 6, 8, 10}),
			},
			"send_resolved": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"auto_annotation": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"query_list": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"chart_title": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"project": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"store": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"store_type": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"query": {
							Type:     schema.TypeString,
							Required: true,
						},
						"logstore": {
							Type:       schema.TypeString,
							Optional:   true,
							Deprecated: "Deprecated from 1.161.0+, use store",
						},
						"start": {
							Type:     schema.TypeString,
							Required: true,
						},
						"end": {
							Type:     schema.TypeString,
							Required: true,
						},
						"time_span_type": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "Custom",
						},
						"region": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"role_arn": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"dashboard_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"power_sql_mode": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringInSlice([]string{"auto", "enable", "disable"}, false),
						},
					},
				},
			},

			"notification_list": {
				Type:       schema.TypeList,
				Optional:   true,
				Deprecated: "Deprecated from 1.161.0+, use policy_configuration for notification",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								sls.NotificationTypeSMS,
								sls.NotificationTypeDingTalk,
								sls.NotificationTypeEmail,
								sls.NotificationTypeMessageCenter},
								false),
						},
						"content": {
							Type:     schema.TypeString,
							Required: true,
						},
						"service_uri": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"mobile_list": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"email_list": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"labels": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:     schema.TypeString,
							Required: true,
						},
						"value": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"annotations": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:     schema.TypeString,
							Required: true,
						},
						"value": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"severity_configurations": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"severity": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntInSlice([]int{2, 4, 6, 8, 10}),
						},
						"eval_condition": {
							Type:     schema.TypeMap,
							Required: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
			"join_configurations": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"cross_join", "inner_join", "left_join", "right_join", "full_join", "left_exclude", "right_exclude", "concat", "no_join"}, false),
						},
						"condition": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"policy_configuration": {
				Type:     schema.TypeSet,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"alert_policy_id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"action_policy_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"repeat_interval": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"group_configuration": {
				Type:     schema.TypeSet,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Required: true,
						},
						"fields": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"schedule_interval": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				Deprecated:    "Field 'schedule_interval' has been deprecated from provider version 1.176.0. New field 'schedule' instead.",
				ConflictsWith: []string{"schedule"},
			},
			"schedule_type": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				Deprecated:    "Field 'schedule_type' has been deprecated from provider version 1.176.0. New field 'schedule' instead.",
				ConflictsWith: []string{"schedule"},
			},
			"schedule": {
				Type:     schema.TypeSet,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"FixedRate", "Cron", "Hourly", "Daily", "Weekly"}, false),
						},
						"interval": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"cron_expression": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"day_of_week": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"hour": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"delay": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"run_immediately": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"time_zone": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
				ConflictsWith: []string{"schedule_type", "schedule_interval"},
			},
		},
	}
}

func resourceAlicloudLogAlertCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	project_name := d.Get("project_name").(string)
	alert_name := d.Get("alert_name").(string)
	alert_displayname := d.Get("alert_displayname").(string)

	schedule := &sls.Schedule{
		Type:     d.Get("schedule_type").(string),
		Interval: d.Get("schedule_interval").(string),
	}
	if v, ok := d.GetOk("schedule"); ok {
		for _, e := range v.(*schema.Set).List() {
			scheduleMap := e.(map[string]interface{})
			schedule.Type, _ = scheduleMap["type"].(string)
			schedule.Interval, _ = scheduleMap["interval"].(string)
			schedule.CronExpression, _ = scheduleMap["cron_expression"].(string)
			schedule.Hour = int32(scheduleMap["hour"].(int))
			schedule.DayOfWeek = int32(scheduleMap["day_of_week"].(int))
			schedule.Delay = int32(scheduleMap["delay"].(int))
			schedule.RunImmediately, _ = scheduleMap["run_immediately"].(bool)
			schedule.TimeZone, _ = scheduleMap["time_zone"].(string)
		}
	}

	alert := &sls.Alert{
		Name:        alert_name,
		DisplayName: alert_displayname,
		Description: d.Get("alert_description").(string),
		State:       "Enabled",
		Schedule:    schedule,
	}
	if err := resource.Retry(2*time.Minute, func() *resource.RetryError {
		_, err := client.WithLogClient(func(slsClient *sls.Client) (interface{}, error) {
			if _, ok := d.GetOk("version"); ok {
				alert.Configuration = createAlert2Config(d)
			} else {
				dashboard := d.Get("dashboard").(string)
				err := CreateDashboard(project_name, dashboard, slsClient)
				if err != nil {
					return nil, err
				}
				alert.Configuration = createAlertConfig(d, project_name, dashboard, slsClient)
			}
			return nil, slsClient.CreateAlert(project_name, alert)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{LogClientTimeout}) {
				time.Sleep(5 * time.Second)
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	}); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_log_alert", "CreateLogstoreAlert", AliyunLogGoSdkERROR)
	}

	d.SetId(fmt.Sprintf("%s%s%s", project_name, COLON_SEPARATED, alert_name))
	return resourceAlicloudLogAlertRead(d, meta)
}

func resourceAlicloudLogAlertRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	logService := LogService{client}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	object, err := logService.DescribeLogAlert(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("project_name", parts[0])
	d.Set("alert_name", parts[1])
	d.Set("alert_displayname", object.DisplayName)
	d.Set("alert_description", object.Description)
	d.Set("mute_until", object.Configuration.MuteUntil)
	d.Set("schedule_interval", object.Schedule.Interval)
	d.Set("schedule_type", object.Schedule.Type)
	d.Set("throttling", object.Configuration.Throttling)
	d.Set("notify_threshold", object.Configuration.NotifyThreshold)
	d.Set("condition", object.Configuration.Condition)
	d.Set("dashboard", object.Configuration.Dashboard)

	var schedules []map[string]interface{}
	scheduleConf := map[string]interface{}{
		"type":            object.Schedule.Type,
		"interval":        object.Schedule.Interval,
		"cron_expression": object.Schedule.CronExpression,
		"delay":           object.Schedule.Delay,
		"day_of_week":     object.Schedule.DayOfWeek,
		"hour":            object.Schedule.Hour,
		"run_immediately": object.Schedule.RunImmediately,
		"time_zone":       object.Schedule.TimeZone,
	}
	schedules = append(schedules, scheduleConf)
	d.Set("schedule", schedules)

	var notiList []map[string]interface{}
	for _, v := range object.Configuration.NotificationList {
		mapping := getNotiMap(v)
		notiList = append(notiList, mapping)
	}

	var queryList []map[string]interface{}
	d.Set("notification_list", notiList)
	isV1 := object.Configuration.Version == ""
	if isV1 {
		for _, v := range object.Configuration.QueryList {
			mapping := map[string]interface{}{
				"chart_title":    v.ChartTitle,
				"logstore":       v.LogStore,
				"query":          v.Query,
				"start":          v.Start,
				"end":            v.End,
				"time_span_type": v.TimeSpanType,
			}
			queryList = append(queryList, mapping)
		}
		d.Set("query_list", queryList)
		return nil
	}
	d.Set("version", object.Configuration.Version)
	d.Set("type", object.Configuration.Type)
	d.Set("threshold", object.Configuration.Threshold)

	var labels []map[string]string
	for _, v := range object.Configuration.Labels {
		labels = append(labels, map[string]string{
			"key":   v.Key,
			"value": v.Value,
		})
	}
	d.Set("labels", labels)

	var annotations []map[string]string
	for _, v := range object.Configuration.Annotations {
		annotations = append(annotations, map[string]string{
			"key":   v.Key,
			"value": v.Value,
		})
	}
	d.Set("annotations", annotations)

	var severityConfigurations []map[string]interface{}
	for _, v := range object.Configuration.SeverityConfigurations {
		severityConf := map[string]interface{}{
			"severity": v.Severity,
			"eval_condition": map[string]interface{}{
				"condition":       v.EvalCondition.Condition,
				"count_condition": v.EvalCondition.CountCondition,
			},
		}
		severityConfigurations = append(severityConfigurations, severityConf)
	}
	d.Set("severity_configurations", severityConfigurations)
	d.Set("no_data_fire", object.Configuration.NoDataFire)
	d.Set("no_data_severity", object.Configuration.NoDataSeverity)
	d.Set("send_resolved", object.Configuration.SendResolved)
	d.Set("auto_annotation", object.Configuration.AutoAnnotation)
	var joinConfigurations []map[string]string
	for _, v := range object.Configuration.JoinConfigurations {
		joinConf := map[string]string{
			"type":      v.Type,
			"condition": v.Condition,
		}
		joinConfigurations = append(joinConfigurations, joinConf)
	}
	d.Set("join_configurations", joinConfigurations)

	var groupConfigurations []map[string]interface{}
	groupConf := map[string]interface{}{
		"type":   object.Configuration.GroupConfiguration.Type,
		"fields": object.Configuration.GroupConfiguration.Fields,
	}
	groupConfigurations = append(groupConfigurations, groupConf)
	d.Set("group_configuration", groupConfigurations)

	var policyConfigurations []map[string]interface{}
	policyConf := map[string]interface{}{
		"alert_policy_id":  object.Configuration.PolicyConfiguration.AlertPolicyId,
		"action_policy_id": object.Configuration.PolicyConfiguration.ActionPolicyId,
		"repeat_interval":  object.Configuration.PolicyConfiguration.RepeatInterval,
	}
	policyConfigurations = append(policyConfigurations, policyConf)
	d.Set("policy_configuration", policyConfigurations)

	for _, v := range object.Configuration.QueryList {
		mapping := map[string]interface{}{
			"chart_title":    v.ChartTitle,
			"region":         v.Region,
			"project":        v.Project,
			"store":          v.Store,
			"store_type":     v.StoreType,
			"query":          v.Query,
			"start":          v.Start,
			"end":            v.End,
			"time_span_type": v.TimeSpanType,
			"role_arn":       v.RoleArn,
			"dashboard_id":   v.DashboardId,
			"power_sql_mode": string(v.PowerSqlMode),
		}
		queryList = append(queryList, mapping)
	}
	d.Set("query_list", queryList)
	return nil
}

func resourceAlicloudLogAlertUpdate(d *schema.ResourceData, meta interface{}) error {

	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}

	client := meta.(*connectivity.AliyunClient)
	schedule := &sls.Schedule{
		Type:     d.Get("schedule_type").(string),
		Interval: d.Get("schedule_interval").(string),
	}
	if v, ok := d.GetOk("schedule"); ok {
		for _, e := range v.(*schema.Set).List() {
			scheduleMap := e.(map[string]interface{})
			schedule.Type, _ = scheduleMap["type"].(string)
			schedule.Interval, _ = scheduleMap["interval"].(string)
			schedule.CronExpression, _ = scheduleMap["cron_expression"].(string)
			schedule.Hour = int32(scheduleMap["hour"].(int))
			schedule.DayOfWeek = int32(scheduleMap["day_of_week"].(int))
			schedule.Delay = int32(scheduleMap["delay"].(int))
			schedule.RunImmediately, _ = scheduleMap["run_immediately"].(bool)
			schedule.TimeZone, _ = scheduleMap["time_zone"].(string)
		}
	}
	params := &sls.Alert{
		Name:        parts[1],
		DisplayName: d.Get("alert_displayname").(string),
		Description: d.Get("alert_description").(string),
		State:       "Enabled",
		Schedule:    schedule,
	}

	if err := resource.Retry(2*time.Minute, func() *resource.RetryError {
		_, err := client.WithLogClient(func(slsClient *sls.Client) (interface{}, error) {
			if _, ok := d.GetOk("version"); ok {
				params.Configuration = createAlert2Config(d)
			} else {
				project_name := d.Get("project_name").(string)
				dashboard := d.Get("dashboard").(string)
				err := CreateDashboard(project_name, dashboard, slsClient)
				if err != nil {
					return nil, err
				}
				params.Configuration = createAlertConfig(d, project_name, dashboard, slsClient)
			}
			return nil, slsClient.UpdateAlert(parts[0], params)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{LogClientTimeout}) {
				time.Sleep(5 * time.Second)
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	}); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), "UpdateAlert", AliyunLogGoSdkERROR)
	}

	return resourceAlicloudLogAlertRead(d, meta)
}

func resourceAlicloudLogAlertDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	logService := LogService{client}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	var requestInfo *sls.Client
	err = resource.Retry(3*time.Minute, func() *resource.RetryError {
		raw, err := client.WithLogClient(func(slsClient *sls.Client) (interface{}, error) {
			requestInfo = slsClient
			return nil, slsClient.DeleteAlert(parts[0], parts[1])
		})
		if err != nil {
			if IsExpectedErrors(err, []string{LogClientTimeout}) {
				time.Sleep(5 * time.Second)
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		if debugOn() {
			addDebug("DeleteAlert", raw, requestInfo, map[string]interface{}{
				"project_name": parts[0],
				"alert":        parts[1],
			})
		}
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_log_alert", "DeleteAlert", AliyunLogGoSdkERROR)
	}
	return WrapError(logService.WaitForLogstoreAlert(d.Id(), Deleted, DefaultTimeout))
}

func createAlertConfig(d *schema.ResourceData, project, dashboard string, client *sls.Client) *sls.AlertConfiguration {

	noti := []*sls.Notification{}
	if v, ok := d.GetOk("notification_list"); ok {
		for _, e := range v.([]interface{}) {
			noti_map := e.(map[string]interface{})
			content := noti_map["content"].(string)

			email_list := []string{}
			email_list_temp := noti_map["email_list"].(*schema.Set).List()
			for _, v := range email_list_temp {
				new_v := v.(string)
				email_list = append(email_list, new_v)
			}
			mobile_list_temp := noti_map["mobile_list"].(*schema.Set).List()
			mobile_list := []string{}
			if len(mobile_list_temp) > 0 {
				for _, v := range mobile_list_temp {
					new_v := v.(string)
					mobile_list = append(mobile_list, new_v)
				}
			}

			if noti_map["type"].(string) == sls.NotificationTypeEmail {
				email := &sls.Notification{
					Type:      sls.NotificationTypeEmail,
					EmailList: email_list,
					Content:   content,
				}
				noti = append(noti, email)
			}

			if noti_map["type"].(string) == sls.NotificationTypeSMS {
				sms := &sls.Notification{
					Type:       sls.NotificationTypeSMS,
					MobileList: mobile_list,
					Content:    content,
				}
				noti = append(noti, sms)
			}
			if noti_map["type"].(string) == sls.NotificationTypeDingTalk {
				ding := &sls.Notification{
					Type:       sls.NotificationTypeDingTalk,
					ServiceUri: noti_map["service_uri"].(string),
					Content:    content,
				}
				noti = append(noti, ding)
			}
			if noti_map["type"].(string) == sls.NotificationTypeMessageCenter {
				messageCenter := &sls.Notification{
					Type:    sls.NotificationTypeMessageCenter,
					Content: content,
				}
				noti = append(noti, messageCenter)
			}
		}
	}

	queryList := []*sls.AlertQuery{}

	if v, ok := d.GetOk("query_list"); ok {
		for _, e := range v.([]interface{}) {
			query_map := e.(map[string]interface{})
			query := &sls.AlertQuery{
				ChartTitle:   GetCharTitile(project, dashboard, query_map["chart_title"].(string), client),
				LogStore:     query_map["logstore"].(string),
				Query:        query_map["query"].(string),
				Start:        query_map["start"].(string),
				End:          query_map["end"].(string),
				TimeSpanType: query_map["time_span_type"].(string),
			}
			queryList = append(queryList, query)

		}
	}

	config := &sls.AlertConfiguration{
		Condition:        d.Get("condition").(string),
		Dashboard:        d.Get("dashboard").(string),
		QueryList:        queryList,
		MuteUntil:        int64(d.Get("mute_until").(int)),
		NotificationList: noti,
		Throttling:       d.Get("throttling").(string),
		NotifyThreshold:  int32(d.Get("notify_threshold").(int)),
	}
	return config
}

func getNotiMap(v *sls.Notification) map[string]interface{} {
	mapping := make(map[string]interface{})

	mapping["content"] = v.Content
	if v.Type == sls.NotificationTypeSMS {
		mapping["type"] = sls.NotificationTypeSMS
		mapping["mobile_list"] = v.MobileList
	}

	if v.Type == sls.NotificationTypeEmail {
		mapping["type"] = sls.NotificationTypeEmail
		mapping["email_list"] = v.EmailList
	}

	if v.Type == sls.NotificationTypeDingTalk {
		mapping["type"] = sls.NotificationTypeDingTalk
		mapping["service_uri"] = v.ServiceUri
	}

	if v.Type == sls.NotificationTypeMessageCenter {
		mapping["type"] = sls.NotificationTypeMessageCenter
	}
	return mapping

}

func createAlert2Config(d *schema.ResourceData) *sls.AlertConfiguration {
	labels := []*sls.Tag{}
	if v, ok := d.GetOk("labels"); ok {
		for _, e := range v.([]interface{}) {
			labelMap := e.(map[string]interface{})
			label := new(sls.Tag)
			label.Key, _ = labelMap["key"].(string)
			label.Value, _ = labelMap["value"].(string)
			labels = append(labels, label)
		}
	}

	annotations := []*sls.Tag{}
	if v, ok := d.GetOk("annotations"); ok {
		for _, e := range v.([]interface{}) {
			annotationMap := e.(map[string]interface{})
			annotation := new(sls.Tag)
			annotation.Key, _ = annotationMap["key"].(string)
			annotation.Value, _ = annotationMap["value"].(string)
			annotations = append(annotations, annotation)
		}
	}

	severityConfigurations := []*sls.SeverityConfiguration{}
	if v, ok := d.GetOk("severity_configurations"); ok {
		for _, e := range v.([]interface{}) {
			severityConfiguration := new(sls.SeverityConfiguration)
			severityConfigurationMap := e.(map[string]interface{})
			severityConfiguration.Severity = sls.Severity(severityConfigurationMap["severity"].(int))
			evalConditionMap := severityConfigurationMap["eval_condition"].(map[string]interface{})
			condition, _ := evalConditionMap["condition"].(string)
			countCondition, _ := evalConditionMap["count_condition"].(string)
			severityConfiguration.EvalCondition = sls.ConditionConfiguration{
				Condition:      condition,
				CountCondition: countCondition,
			}
			severityConfigurations = append(severityConfigurations, severityConfiguration)
		}
	}

	joinConfigurations := []*sls.JoinConfiguration{}
	if v, ok := d.GetOk("join_configurations"); ok {
		for _, e := range v.([]interface{}) {
			joinConfigurationMap := e.(map[string]interface{})
			joinConfiguration := new(sls.JoinConfiguration)
			joinConfiguration.Type, _ = joinConfigurationMap["type"].(string)
			joinConfiguration.Condition, _ = joinConfigurationMap["condition"].(string)
			joinConfigurations = append(joinConfigurations, joinConfiguration)
		}
	}

	groupConfiguration := sls.GroupConfiguration{}
	if v, ok := d.GetOk("group_configuration"); ok {
		for _, e := range v.(*schema.Set).List() {
			groupConfigurationMap := e.(map[string]interface{})
			groupConfiguration.Type, _ = groupConfigurationMap["type"].(string)

			fields := []string{}
			for _, v := range groupConfigurationMap["fields"].(*schema.Set).List() {
				fields = append(fields, v.(string))
			}
			groupConfiguration.Fields = fields
		}
	}

	policyConfiguration := sls.PolicyConfiguration{}
	if v, ok := d.GetOk("policy_configuration"); ok {
		for _, e := range v.(*schema.Set).List() {
			policyConfigurationMap := e.(map[string]interface{})
			policyConfiguration.AlertPolicyId, _ = policyConfigurationMap["alert_policy_id"].(string)
			policyConfiguration.ActionPolicyId, _ = policyConfigurationMap["action_policy_id"].(string)
			policyConfiguration.RepeatInterval, _ = policyConfigurationMap["repeat_interval"].(string)
		}
	}

	queryList := []*sls.AlertQuery{}

	if v, ok := d.GetOk("query_list"); ok {
		for _, e := range v.([]interface{}) {
			query_map := e.(map[string]interface{})
			query := &sls.AlertQuery{
				DashboardId:  query_map["dashboard_id"].(string),
				ChartTitle:   query_map["chart_title"].(string),
				RoleArn:      query_map["role_arn"].(string),
				Region:       query_map["region"].(string),
				Project:      query_map["project"].(string),
				StoreType:    query_map["store_type"].(string),
				Store:        query_map["store"].(string),
				Query:        query_map["query"].(string),
				Start:        query_map["start"].(string),
				End:          query_map["end"].(string),
				TimeSpanType: query_map["time_span_type"].(string),
				PowerSqlMode: sls.PowerSqlMode(query_map["power_sql_mode"].(string)),
			}
			queryList = append(queryList, query)
		}
	}

	version, ok := d.GetOk("version")
	if !ok {
		version = "2.0"
	}
	configType, ok := d.GetOk("type")
	if !ok {
		configType = "default"
	}
	autoAnnotation := false
	if v, ok := d.GetOk("auto_annotation"); ok {
		autoAnnotation = v.(bool)
	}

	config := &sls.AlertConfiguration{
		Version:                version.(string),
		Type:                   configType.(string),
		Labels:                 labels,
		Annotations:            annotations,
		SeverityConfigurations: severityConfigurations,
		JoinConfigurations:     joinConfigurations,
		GroupConfiguration:     groupConfiguration,
		PolicyConfiguration:    policyConfiguration,
		QueryList:              queryList,
		Threshold:              d.Get("threshold").(int),
		Dashboard:              d.Get("dashboard").(string),
		NoDataFire:             d.Get("no_data_fire").(bool),
		NoDataSeverity:         sls.Severity(d.Get("no_data_severity").(int)),
		SendResolved:           d.Get("send_resolved").(bool),
		MuteUntil:              int64(d.Get("mute_until").(int)),
		AutoAnnotation:         autoAnnotation,
	}
	return config
}

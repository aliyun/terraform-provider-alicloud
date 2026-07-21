// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudSlsScheduledSql() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudSlsScheduledSqlCreate,
		Read:   resourceAliCloudSlsScheduledSqlRead,
		Update: resourceAliCloudSlsScheduledSqlUpdate,
		Delete: resourceAliCloudSlsScheduledSqlDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"display_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"project": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"schedule": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"time_zone": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"run_immediately": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"cron_expression": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"delay": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"interval": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"scheduled_sql_configuration": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"max_retries": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"script": {
							Type:     schema.TypeString,
							Optional: true,
						},
						// lintignore: S006
						"parameters": {
							Type:     schema.TypeMap,
							Optional: true,
						},
						"resource_pool": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"from_time_expr": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"dest_role_arn": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"role_arn": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"to_time": {
							Type:     schema.TypeInt,
							Optional: true,
							ForceNew: true,
						},
						"max_run_time_in_seconds": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"data_format": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"sql_type": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"source_logstore": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"dest_logstore": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"dest_endpoint": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"to_time_expr": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"from_time": {
							Type:     schema.TypeInt,
							Optional: true,
							ForceNew: true,
						},
						"dest_project": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"scheduled_sql_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func resourceAliCloudSlsScheduledSqlCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := fmt.Sprintf("/scheduledsqls")
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	body := make(map[string]interface{})
	hostMap := make(map[string]*string)
	var err error
	request = make(map[string]interface{})
	if v, ok := d.GetOk("scheduled_sql_name"); ok {
		request["name"] = v
	}
	hostMap["project"] = StringPointer(d.Get("project").(string))

	configuration := make(map[string]interface{})

	if v := d.Get("scheduled_sql_configuration"); v != nil {
		destEndpoint1, _ := jsonpath.Get("$[0].dest_endpoint", v)
		if destEndpoint1 != nil && destEndpoint1 != "" {
			configuration["destEndpoint"] = destEndpoint1
		}
		sourceLogstore1, _ := jsonpath.Get("$[0].source_logstore", v)
		if sourceLogstore1 != nil && sourceLogstore1 != "" {
			configuration["sourceLogstore"] = sourceLogstore1
		}
		destProject1, _ := jsonpath.Get("$[0].dest_project", v)
		if destProject1 != nil && destProject1 != "" {
			configuration["destProject"] = destProject1
		}
		maxRunTimeInSeconds1, _ := jsonpath.Get("$[0].max_run_time_in_seconds", v)
		if maxRunTimeInSeconds1 != nil && maxRunTimeInSeconds1 != "" {
			configuration["maxRunTimeInSeconds"] = maxRunTimeInSeconds1
		}
		parameters1, _ := jsonpath.Get("$[0].parameters", v)
		if parameters1 != nil && parameters1 != "" {
			configuration["parameters"] = parameters1
		}
		destRoleArn1, _ := jsonpath.Get("$[0].dest_role_arn", v)
		if destRoleArn1 != nil && destRoleArn1 != "" {
			configuration["destRoleArn"] = destRoleArn1
		}
		destLogstore1, _ := jsonpath.Get("$[0].dest_logstore", v)
		if destLogstore1 != nil && destLogstore1 != "" {
			configuration["destLogstore"] = destLogstore1
		}
		fromTimeExpr1, _ := jsonpath.Get("$[0].from_time_expr", v)
		if fromTimeExpr1 != nil && fromTimeExpr1 != "" {
			configuration["fromTimeExpr"] = fromTimeExpr1
		}
		sqlType1, _ := jsonpath.Get("$[0].sql_type", v)
		if sqlType1 != nil && sqlType1 != "" {
			configuration["sqlType"] = sqlType1
		}
		toTime1, _ := jsonpath.Get("$[0].to_time", v)
		if toTime1 != nil && toTime1 != "" {
			configuration["toTime"] = toTime1
		}
		script1, _ := jsonpath.Get("$[0].script", v)
		if script1 != nil && script1 != "" {
			configuration["script"] = script1
		}
		fromTime1, _ := jsonpath.Get("$[0].from_time", v)
		if fromTime1 != nil && fromTime1 != "" {
			configuration["fromTime"] = fromTime1
		}
		dataFormat1, _ := jsonpath.Get("$[0].data_format", v)
		if dataFormat1 != nil && dataFormat1 != "" {
			configuration["dataFormat"] = dataFormat1
		}
		roleArn1, _ := jsonpath.Get("$[0].role_arn", v)
		if roleArn1 != nil && roleArn1 != "" {
			configuration["roleArn"] = roleArn1
		}
		maxRetries1, _ := jsonpath.Get("$[0].max_retries", v)
		if maxRetries1 != nil && maxRetries1 != "" {
			configuration["maxRetries"] = maxRetries1
		}
		toTimeExpr1, _ := jsonpath.Get("$[0].to_time_expr", v)
		if toTimeExpr1 != nil && toTimeExpr1 != "" {
			configuration["toTimeExpr"] = toTimeExpr1
		}
		resourcePool1, _ := jsonpath.Get("$[0].resource_pool", v)
		if resourcePool1 != nil && resourcePool1 != "" {
			configuration["resourcePool"] = resourcePool1
		}

		request["configuration"] = configuration
	}

	schedule := make(map[string]interface{})

	if v := d.Get("schedule"); v != nil {
		cronExpression1, _ := jsonpath.Get("$[0].cron_expression", v)
		if cronExpression1 != nil && cronExpression1 != "" {
			schedule["cronExpression"] = cronExpression1
		}
		delay1, _ := jsonpath.Get("$[0].delay", v)
		if delay1 != nil && delay1 != "" {
			schedule["delay"] = delay1
		}
		timeZone1, _ := jsonpath.Get("$[0].time_zone", v)
		if timeZone1 != nil && timeZone1 != "" {
			schedule["timeZone"] = timeZone1
		}
		type1, _ := jsonpath.Get("$[0].type", v)
		if type1 != nil && type1 != "" {
			schedule["type"] = type1
		}
		interval1, _ := jsonpath.Get("$[0].interval", v)
		if interval1 != nil && interval1 != "" {
			schedule["interval"] = interval1
		}
		runImmediately1, _ := jsonpath.Get("$[0].run_immediately", v)
		if runImmediately1 != nil && runImmediately1 != "" {
			schedule["runImmediately"] = runImmediately1
		}

		request["schedule"] = schedule
	}

	if v, ok := d.GetOk("description"); ok {
		request["description"] = v
	}
	request["displayName"] = d.Get("display_name")
	body = request
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.Do("Sls", roaParam("POST", "2020-12-30", "CreateScheduledSQL", action), query, body, nil, hostMap, false)
		if err != nil {
			if IsExpectedErrors(err, []string{"403"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_sls_scheduled_sql", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%v:%v", *hostMap["project"], request["name"]))

	return resourceAliCloudSlsScheduledSqlUpdate(d, meta)
}

func resourceAliCloudSlsScheduledSqlRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	slsServiceV2 := SlsServiceV2{client}

	objectRaw, err := slsServiceV2.DescribeSlsScheduledSql(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_sls_scheduled_sql DescribeSlsScheduledSql Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("description", objectRaw["description"])
	d.Set("display_name", objectRaw["displayName"])
	d.Set("status", objectRaw["status"])

	scheduleMaps := make([]map[string]interface{}, 0)
	scheduleMap := make(map[string]interface{})
	scheduleRaw := make(map[string]interface{})
	if objectRaw["schedule"] != nil {
		scheduleRaw = objectRaw["schedule"].(map[string]interface{})
	}
	if len(scheduleRaw) > 0 {
		scheduleMap["cron_expression"] = scheduleRaw["cronExpression"]
		scheduleMap["delay"] = scheduleRaw["delay"]
		scheduleMap["interval"] = scheduleRaw["interval"]
		scheduleMap["run_immediately"] = scheduleRaw["runImmediately"]
		scheduleMap["time_zone"] = scheduleRaw["timeZone"]
		scheduleMap["type"] = scheduleRaw["type"]

		scheduleMaps = append(scheduleMaps, scheduleMap)
	}
	if err := d.Set("schedule", scheduleMaps); err != nil {
		return err
	}
	scheduledSQLConfigurationMaps := make([]map[string]interface{}, 0)
	scheduledSQLConfigurationMap := make(map[string]interface{})
	configurationRaw := make(map[string]interface{})
	if objectRaw["configuration"] != nil {
		configurationRaw = objectRaw["configuration"].(map[string]interface{})
	}
	if len(configurationRaw) > 0 {
		scheduledSQLConfigurationMap["data_format"] = configurationRaw["dataFormat"]
		scheduledSQLConfigurationMap["dest_endpoint"] = configurationRaw["destEndpoint"]
		scheduledSQLConfigurationMap["dest_logstore"] = configurationRaw["destLogstore"]
		scheduledSQLConfigurationMap["dest_project"] = configurationRaw["destProject"]
		scheduledSQLConfigurationMap["dest_role_arn"] = configurationRaw["destRoleArn"]
		scheduledSQLConfigurationMap["from_time"] = configurationRaw["fromTime"]
		scheduledSQLConfigurationMap["from_time_expr"] = configurationRaw["fromTimeExpr"]
		scheduledSQLConfigurationMap["max_retries"] = configurationRaw["maxRetries"]
		scheduledSQLConfigurationMap["max_run_time_in_seconds"] = configurationRaw["maxRunTimeInSeconds"]
		scheduledSQLConfigurationMap["parameters"] = configurationRaw["parameters"]
		scheduledSQLConfigurationMap["resource_pool"] = configurationRaw["resourcePool"]
		scheduledSQLConfigurationMap["role_arn"] = configurationRaw["roleArn"]
		scheduledSQLConfigurationMap["script"] = configurationRaw["script"]
		scheduledSQLConfigurationMap["source_logstore"] = configurationRaw["sourceLogstore"]
		scheduledSQLConfigurationMap["sql_type"] = configurationRaw["sqlType"]
		scheduledSQLConfigurationMap["to_time"] = configurationRaw["toTime"]
		scheduledSQLConfigurationMap["to_time_expr"] = configurationRaw["toTimeExpr"]

		scheduledSQLConfigurationMaps = append(scheduledSQLConfigurationMaps, scheduledSQLConfigurationMap)
	}
	if err := d.Set("scheduled_sql_configuration", scheduledSQLConfigurationMaps); err != nil {
		return err
	}

	parts := strings.Split(d.Id(), ":")
	d.Set("project", parts[0])
	d.Set("scheduled_sql_name", parts[1])

	return nil
}

func resourceAliCloudSlsScheduledSqlUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	var body map[string]interface{}
	update := false

	slsServiceV2 := SlsServiceV2{client}
	objectRaw, _ := slsServiceV2.DescribeSlsScheduledSql(d.Id())

	if d.HasChange("status") {
		var err error
		target := d.Get("status").(string)

		currentStatus, err := jsonpath.Get("status", objectRaw)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, d.Id(), "status", objectRaw)
		}
		if fmt.Sprint(currentStatus) != target {
			if target == "ENABLED" {
				parts := strings.Split(d.Id(), ":")
				scheduledSQLName := parts[1]
				action := fmt.Sprintf("/scheduledsqls/%s?action=enable", scheduledSQLName)
				request = make(map[string]interface{})
				query = make(map[string]*string)
				body = make(map[string]interface{})
				hostMap := make(map[string]*string)
				hostMap["project"] = StringPointer(parts[0])

				body = request
				wait := incrementalWait(3*time.Second, 5*time.Second)
				err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
					response, err = client.Do("Sls", roaParam("PUT", "2020-12-30", "EnableScheduledSQL", action), query, body, nil, hostMap, false)
					if err != nil {
						if NeedRetry(err) {
							wait()
							return resource.RetryableError(err)
						}
						return resource.NonRetryableError(err)
					}
					return nil
				})
				addDebug(action, response, request)
				if err != nil {
					return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
				}

			}
			if target == "DISABLED" {
				parts := strings.Split(d.Id(), ":")
				scheduledSQLName := parts[1]
				action := fmt.Sprintf("/scheduledsqls/%s?action=disable", scheduledSQLName)
				request = make(map[string]interface{})
				query = make(map[string]*string)
				body = make(map[string]interface{})
				hostMap := make(map[string]*string)
				hostMap["project"] = StringPointer(parts[0])

				body = request
				wait := incrementalWait(3*time.Second, 5*time.Second)
				err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
					response, err = client.Do("Sls", roaParam("PUT", "2020-12-30", "DisableScheduledSQL", action), query, body, nil, hostMap, false)
					if err != nil {
						if NeedRetry(err) {
							wait()
							return resource.RetryableError(err)
						}
						return resource.NonRetryableError(err)
					}
					return nil
				})
				addDebug(action, response, request)
				if err != nil {
					return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
				}

			}
		}
	}

	var err error
	parts := strings.Split(d.Id(), ":")
	scheduledSQLName := parts[1]
	action := fmt.Sprintf("/scheduledsqls/%s", scheduledSQLName)
	request = make(map[string]interface{})
	query = make(map[string]*string)
	body = make(map[string]interface{})
	hostMap := make(map[string]*string)
	hostMap["project"] = StringPointer(parts[0])

	if !d.IsNewResource() && d.HasChange("scheduled_sql_configuration") {
		update = true
	}
	configuration := make(map[string]interface{})

	if v := d.Get("scheduled_sql_configuration"); v != nil {
		destEndpoint1, _ := jsonpath.Get("$[0].dest_endpoint", v)
		if destEndpoint1 != nil && destEndpoint1 != "" {
			configuration["destEndpoint"] = destEndpoint1
		}
		sourceLogstore1, _ := jsonpath.Get("$[0].source_logstore", v)
		if sourceLogstore1 != nil && sourceLogstore1 != "" {
			configuration["sourceLogstore"] = sourceLogstore1
		}
		destProject1, _ := jsonpath.Get("$[0].dest_project", v)
		if destProject1 != nil && destProject1 != "" {
			configuration["destProject"] = destProject1
		}
		maxRunTimeInSeconds1, _ := jsonpath.Get("$[0].max_run_time_in_seconds", v)
		if maxRunTimeInSeconds1 != nil && maxRunTimeInSeconds1 != "" {
			configuration["maxRunTimeInSeconds"] = maxRunTimeInSeconds1
		}
		parameters1, _ := jsonpath.Get("$[0].parameters", v)
		if parameters1 != nil && parameters1 != "" {
			configuration["parameters"] = parameters1
		}
		destRoleArn1, _ := jsonpath.Get("$[0].dest_role_arn", v)
		if destRoleArn1 != nil && destRoleArn1 != "" {
			configuration["destRoleArn"] = destRoleArn1
		}
		destLogstore1, _ := jsonpath.Get("$[0].dest_logstore", v)
		if destLogstore1 != nil && destLogstore1 != "" {
			configuration["destLogstore"] = destLogstore1
		}
		fromTimeExpr1, _ := jsonpath.Get("$[0].from_time_expr", v)
		if fromTimeExpr1 != nil && fromTimeExpr1 != "" {
			configuration["fromTimeExpr"] = fromTimeExpr1
		}
		sqlType1, _ := jsonpath.Get("$[0].sql_type", v)
		if sqlType1 != nil && sqlType1 != "" {
			configuration["sqlType"] = sqlType1
		}
		toTime1, _ := jsonpath.Get("$[0].to_time", v)
		if toTime1 != nil && toTime1 != "" {
			configuration["toTime"] = toTime1
		}
		script1, _ := jsonpath.Get("$[0].script", v)
		if script1 != nil && script1 != "" {
			configuration["script"] = script1
		}
		fromTime1, _ := jsonpath.Get("$[0].from_time", v)
		if fromTime1 != nil && fromTime1 != "" {
			configuration["fromTime"] = fromTime1
		}
		dataFormat1, _ := jsonpath.Get("$[0].data_format", v)
		if dataFormat1 != nil && dataFormat1 != "" {
			configuration["dataFormat"] = dataFormat1
		}
		roleArn1, _ := jsonpath.Get("$[0].role_arn", v)
		if roleArn1 != nil && roleArn1 != "" {
			configuration["roleArn"] = roleArn1
		}
		maxRetries1, _ := jsonpath.Get("$[0].max_retries", v)
		if maxRetries1 != nil && maxRetries1 != "" {
			configuration["maxRetries"] = maxRetries1
		}
		toTimeExpr1, _ := jsonpath.Get("$[0].to_time_expr", v)
		if toTimeExpr1 != nil && toTimeExpr1 != "" {
			configuration["toTimeExpr"] = toTimeExpr1
		}
		resourcePool1, _ := jsonpath.Get("$[0].resource_pool", v)
		if resourcePool1 != nil && resourcePool1 != "" {
			configuration["resourcePool"] = resourcePool1
		}

		request["configuration"] = configuration
	}

	if !d.IsNewResource() && d.HasChange("schedule") {
		update = true
	}
	schedule := make(map[string]interface{})

	if v := d.Get("schedule"); v != nil {
		cronExpression1, _ := jsonpath.Get("$[0].cron_expression", v)
		if cronExpression1 != nil && cronExpression1 != "" {
			schedule["cronExpression"] = cronExpression1
		}
		delay1, _ := jsonpath.Get("$[0].delay", v)
		if delay1 != nil && delay1 != "" {
			schedule["delay"] = delay1
		}
		timeZone1, _ := jsonpath.Get("$[0].time_zone", v)
		if timeZone1 != nil && timeZone1 != "" {
			schedule["timeZone"] = timeZone1
		}
		type1, _ := jsonpath.Get("$[0].type", v)
		if type1 != nil && type1 != "" {
			schedule["type"] = type1
		}
		interval1, _ := jsonpath.Get("$[0].interval", v)
		if interval1 != nil && interval1 != "" {
			schedule["interval"] = interval1
		}
		runImmediately1, _ := jsonpath.Get("$[0].run_immediately", v)
		if runImmediately1 != nil && runImmediately1 != "" {
			schedule["runImmediately"] = runImmediately1
		}

		request["schedule"] = schedule
	}

	if !d.IsNewResource() && d.HasChange("description") {
		update = true
	}
	if v, ok := d.GetOk("description"); ok || d.HasChange("description") {
		request["description"] = v
	}
	if !d.IsNewResource() && d.HasChange("display_name") {
		update = true
	}
	request["displayName"] = d.Get("display_name")
	body = request
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.Do("Sls", roaParam("PUT", "2020-12-30", "UpdateScheduledSQL", action), query, body, nil, hostMap, false)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, request)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
	}

	return resourceAliCloudSlsScheduledSqlRead(d, meta)
}

func resourceAliCloudSlsScheduledSqlDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	parts := strings.Split(d.Id(), ":")
	scheduledSQLName := parts[1]
	action := fmt.Sprintf("/scheduledsqls/%s", scheduledSQLName)
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	hostMap := make(map[string]*string)
	var err error
	request = make(map[string]interface{})
	hostMap["project"] = StringPointer(parts[0])

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.Do("Sls", roaParam("DELETE", "2020-12-30", "DeleteScheduledSQL", action), query, nil, nil, hostMap, false)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		if IsExpectedErrors(err, []string{"403"}) || NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return nil
}

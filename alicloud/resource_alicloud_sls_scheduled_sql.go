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

func resourceAliCloudSlsScheduledSQL() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudSlsScheduledSQLCreate,
		Read:   resourceAliCloudSlsScheduledSQLRead,
		Update: resourceAliCloudSlsScheduledSQLUpdate,
		Delete: resourceAliCloudSlsScheduledSQLDelete,
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
		},
	}
}

func resourceAliCloudSlsScheduledSQLCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := fmt.Sprintf("/scheduledsqls")
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	body := make(map[string]interface{})
	hostMap := make(map[string]*string)
	var err error
	request = make(map[string]interface{})
	hostMap["project"] = StringPointer(d.Get("project").(string))
	request["name"] = d.Get("scheduled_sql_name")

	request["displayName"] = d.Get("display_name")
	if v, ok := d.GetOk("description"); ok {
		request["description"] = v
	}
	objectDataLocalMap := make(map[string]interface{})
	if v := d.Get("schedule"); v != nil {
		nodeNative, _ := jsonpath.Get("$[0].type", d.Get("schedule"))
		if nodeNative != nil && nodeNative != "" {
			objectDataLocalMap["type"] = nodeNative
		}
		nodeNative1, _ := jsonpath.Get("$[0].cron_expression", d.Get("schedule"))
		if nodeNative1 != nil && nodeNative1 != "" {
			objectDataLocalMap["cronExpression"] = nodeNative1
		}
		nodeNative2, _ := jsonpath.Get("$[0].run_immediately", d.Get("schedule"))
		if nodeNative2 != nil && nodeNative2 != "" {
			objectDataLocalMap["runImmediately"] = nodeNative2
		}
		nodeNative3, _ := jsonpath.Get("$[0].time_zone", d.Get("schedule"))
		if nodeNative3 != nil && nodeNative3 != "" {
			objectDataLocalMap["timeZone"] = nodeNative3
		}
		nodeNative4, _ := jsonpath.Get("$[0].delay", d.Get("schedule"))
		if nodeNative4 != nil && nodeNative4 != "" {
			objectDataLocalMap["delay"] = nodeNative4
		}
		nodeNative5, _ := jsonpath.Get("$[0].interval", d.Get("schedule"))
		if nodeNative5 != nil && nodeNative5 != "" {
			objectDataLocalMap["interval"] = nodeNative5
		}

		request["schedule"] = objectDataLocalMap
	}

	objectDataLocalMap1 := make(map[string]interface{})
	if v := d.Get("scheduled_sql_configuration"); v != nil {
		nodeNative6, _ := jsonpath.Get("$[0].script", d.Get("scheduled_sql_configuration"))
		if nodeNative6 != nil && nodeNative6 != "" {
			objectDataLocalMap1["script"] = nodeNative6
		}
		nodeNative7, _ := jsonpath.Get("$[0].sql_type", d.Get("scheduled_sql_configuration"))
		if nodeNative7 != nil && nodeNative7 != "" {
			objectDataLocalMap1["sqlType"] = nodeNative7
		}
		nodeNative8, _ := jsonpath.Get("$[0].dest_endpoint", d.Get("scheduled_sql_configuration"))
		if nodeNative8 != nil && nodeNative8 != "" {
			objectDataLocalMap1["destEndpoint"] = nodeNative8
		}
		nodeNative9, _ := jsonpath.Get("$[0].dest_project", d.Get("scheduled_sql_configuration"))
		if nodeNative9 != nil && nodeNative9 != "" {
			objectDataLocalMap1["destProject"] = nodeNative9
		}
		nodeNative10, _ := jsonpath.Get("$[0].source_logstore", d.Get("scheduled_sql_configuration"))
		if nodeNative10 != nil && nodeNative10 != "" {
			objectDataLocalMap1["sourceLogstore"] = nodeNative10
		}
		nodeNative11, _ := jsonpath.Get("$[0].dest_logstore", d.Get("scheduled_sql_configuration"))
		if nodeNative11 != nil && nodeNative11 != "" {
			objectDataLocalMap1["destLogstore"] = nodeNative11
		}
		nodeNative12, _ := jsonpath.Get("$[0].role_arn", d.Get("scheduled_sql_configuration"))
		if nodeNative12 != nil && nodeNative12 != "" {
			objectDataLocalMap1["roleArn"] = nodeNative12
		}
		nodeNative13, _ := jsonpath.Get("$[0].dest_role_arn", d.Get("scheduled_sql_configuration"))
		if nodeNative13 != nil && nodeNative13 != "" {
			objectDataLocalMap1["destRoleArn"] = nodeNative13
		}
		nodeNative14, _ := jsonpath.Get("$[0].from_time_expr", d.Get("scheduled_sql_configuration"))
		if nodeNative14 != nil && nodeNative14 != "" {
			objectDataLocalMap1["fromTimeExpr"] = nodeNative14
		}
		nodeNative15, _ := jsonpath.Get("$[0].to_time_expr", d.Get("scheduled_sql_configuration"))
		if nodeNative15 != nil && nodeNative15 != "" {
			objectDataLocalMap1["toTimeExpr"] = nodeNative15
		}
		nodeNative16, _ := jsonpath.Get("$[0].max_run_time_in_seconds", d.Get("scheduled_sql_configuration"))
		if nodeNative16 != nil && nodeNative16 != "" {
			objectDataLocalMap1["maxRunTimeInSeconds"] = nodeNative16
		}
		nodeNative17, _ := jsonpath.Get("$[0].resource_pool", d.Get("scheduled_sql_configuration"))
		if nodeNative17 != nil && nodeNative17 != "" {
			objectDataLocalMap1["resourcePool"] = nodeNative17
		}
		nodeNative18, _ := jsonpath.Get("$[0].max_retries", d.Get("scheduled_sql_configuration"))
		if nodeNative18 != nil && nodeNative18 != "" {
			objectDataLocalMap1["maxRetries"] = nodeNative18
		}
		nodeNative19, _ := jsonpath.Get("$[0].from_time", d.Get("scheduled_sql_configuration"))
		if nodeNative19 != nil && nodeNative19 != "" {
			objectDataLocalMap1["fromTime"] = nodeNative19
		}
		nodeNative20, _ := jsonpath.Get("$[0].to_time", d.Get("scheduled_sql_configuration"))
		if nodeNative20 != nil && nodeNative20 != "" {
			objectDataLocalMap1["toTime"] = nodeNative20
		}
		nodeNative21, _ := jsonpath.Get("$[0].data_format", d.Get("scheduled_sql_configuration"))
		if nodeNative21 != nil && nodeNative21 != "" {
			objectDataLocalMap1["dataFormat"] = nodeNative21
		}
		nodeNative22, _ := jsonpath.Get("$[0].parameters", d.Get("scheduled_sql_configuration"))
		if nodeNative22 != nil && nodeNative22 != "" {
			objectDataLocalMap1["parameters"] = nodeNative22
		}

		request["configuration"] = objectDataLocalMap1
	}

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
		addDebug(action, response, request)
		return nil
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_sls_scheduled_sql", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%v:%v", *hostMap["project"], request["name"]))

	return resourceAliCloudSlsScheduledSQLUpdate(d, meta)
}

func resourceAliCloudSlsScheduledSQLRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	slsServiceV2 := SlsServiceV2{client}

	objectRaw, err := slsServiceV2.DescribeSlsScheduledSQL(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_sls_scheduled_sql DescribeSlsScheduledSQL Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("description", objectRaw["description"])
	d.Set("display_name", objectRaw["displayName"])

	scheduleMaps := make([]map[string]interface{}, 0)
	scheduleMap := make(map[string]interface{})
	schedule1Raw := make(map[string]interface{})
	if objectRaw["schedule"] != nil {
		schedule1Raw = objectRaw["schedule"].(map[string]interface{})
	}
	if len(schedule1Raw) > 0 {
		scheduleMap["cron_expression"] = schedule1Raw["cronExpression"]
		scheduleMap["delay"] = schedule1Raw["delay"]
		scheduleMap["interval"] = schedule1Raw["interval"]
		scheduleMap["run_immediately"] = schedule1Raw["runImmediately"]
		scheduleMap["time_zone"] = schedule1Raw["timeZone"]
		scheduleMap["type"] = schedule1Raw["type"]

		scheduleMaps = append(scheduleMaps, scheduleMap)
	}
	d.Set("schedule", scheduleMaps)
	scheduledSQLConfigurationMaps := make([]map[string]interface{}, 0)
	scheduledSQLConfigurationMap := make(map[string]interface{})
	configuration1Raw := make(map[string]interface{})
	if objectRaw["configuration"] != nil {
		configuration1Raw = objectRaw["configuration"].(map[string]interface{})
	}
	if len(configuration1Raw) > 0 {
		scheduledSQLConfigurationMap["data_format"] = configuration1Raw["dataFormat"]
		scheduledSQLConfigurationMap["dest_endpoint"] = configuration1Raw["destEndpoint"]
		scheduledSQLConfigurationMap["dest_logstore"] = configuration1Raw["destLogstore"]
		scheduledSQLConfigurationMap["dest_project"] = configuration1Raw["destProject"]
		scheduledSQLConfigurationMap["dest_role_arn"] = configuration1Raw["destRoleArn"]
		scheduledSQLConfigurationMap["from_time"] = configuration1Raw["fromTime"]
		scheduledSQLConfigurationMap["from_time_expr"] = configuration1Raw["fromTimeExpr"]
		scheduledSQLConfigurationMap["max_retries"] = configuration1Raw["maxRetries"]
		scheduledSQLConfigurationMap["max_run_time_in_seconds"] = configuration1Raw["maxRunTimeInSeconds"]
		scheduledSQLConfigurationMap["parameters"] = configuration1Raw["parameters"]
		scheduledSQLConfigurationMap["resource_pool"] = configuration1Raw["resourcePool"]
		scheduledSQLConfigurationMap["role_arn"] = configuration1Raw["roleArn"]
		scheduledSQLConfigurationMap["script"] = configuration1Raw["script"]
		scheduledSQLConfigurationMap["source_logstore"] = configuration1Raw["sourceLogstore"]
		scheduledSQLConfigurationMap["sql_type"] = configuration1Raw["sqlType"]
		scheduledSQLConfigurationMap["to_time"] = configuration1Raw["toTime"]
		scheduledSQLConfigurationMap["to_time_expr"] = configuration1Raw["toTimeExpr"]

		scheduledSQLConfigurationMaps = append(scheduledSQLConfigurationMaps, scheduledSQLConfigurationMap)
	}
	d.Set("scheduled_sql_configuration", scheduledSQLConfigurationMaps)

	parts := strings.Split(d.Id(), ":")
	d.Set("project", parts[0])
	d.Set("scheduled_sql_name", parts[1])

	return nil
}

func resourceAliCloudSlsScheduledSQLUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	var body map[string]interface{}
	update := false
	parts := strings.Split(d.Id(), ":")
	scheduledSQLName := parts[1]
	action := fmt.Sprintf("/scheduledsqls/%s", scheduledSQLName)
	var err error
	request = make(map[string]interface{})
	query = make(map[string]*string)
	body = make(map[string]interface{})
	hostMap := make(map[string]*string)
	hostMap["project"] = StringPointer(parts[0])
	if !d.IsNewResource() && d.HasChange("display_name") {
		update = true
	}
	request["displayName"] = d.Get("display_name")
	if !d.IsNewResource() && d.HasChange("description") {
		update = true
	}
	request["description"] = d.Get("description")
	if d.HasChange("schedule") {
		update = true
	}
	objectDataLocalMap := make(map[string]interface{})
	if v := d.Get("schedule"); v != nil {
		nodeNative, _ := jsonpath.Get("$[0].type", v)
		if nodeNative != nil && nodeNative != "" {
			objectDataLocalMap["type"] = nodeNative
		}
		nodeNative1, _ := jsonpath.Get("$[0].cron_expression", v)
		if nodeNative1 != nil && nodeNative1 != "" {
			objectDataLocalMap["cronExpression"] = nodeNative1
		}
		nodeNative2, _ := jsonpath.Get("$[0].run_immediately", v)
		if nodeNative2 != nil && nodeNative2 != "" {
			objectDataLocalMap["runImmediately"] = nodeNative2
		}
		nodeNative3, _ := jsonpath.Get("$[0].time_zone", v)
		if nodeNative3 != nil && nodeNative3 != "" {
			objectDataLocalMap["timeZone"] = nodeNative3
		}
		nodeNative4, _ := jsonpath.Get("$[0].delay", v)
		if nodeNative4 != nil && nodeNative4 != "" {
			objectDataLocalMap["delay"] = nodeNative4
		}
		nodeNative5, _ := jsonpath.Get("$[0].interval", v)
		if nodeNative5 != nil && nodeNative5 != "" {
			objectDataLocalMap["interval"] = nodeNative5
		}

		request["schedule"] = objectDataLocalMap
	}

	if d.HasChange("scheduled_sql_configuration") {
		update = true
	}
	objectDataLocalMap1 := make(map[string]interface{})
	if v := d.Get("scheduled_sql_configuration"); v != nil {
		nodeNative6, _ := jsonpath.Get("$[0].script", v)
		if nodeNative6 != nil && nodeNative6 != "" {
			objectDataLocalMap1["script"] = nodeNative6
		}
		nodeNative7, _ := jsonpath.Get("$[0].sql_type", v)
		if nodeNative7 != nil && nodeNative7 != "" {
			objectDataLocalMap1["sqlType"] = nodeNative7
		}
		nodeNative8, _ := jsonpath.Get("$[0].dest_endpoint", v)
		if nodeNative8 != nil && nodeNative8 != "" {
			objectDataLocalMap1["destEndpoint"] = nodeNative8
		}
		nodeNative9, _ := jsonpath.Get("$[0].dest_project", v)
		if nodeNative9 != nil && nodeNative9 != "" {
			objectDataLocalMap1["destProject"] = nodeNative9
		}
		nodeNative10, _ := jsonpath.Get("$[0].source_logstore", v)
		if nodeNative10 != nil && nodeNative10 != "" {
			objectDataLocalMap1["sourceLogstore"] = nodeNative10
		}
		nodeNative11, _ := jsonpath.Get("$[0].dest_logstore", v)
		if nodeNative11 != nil && nodeNative11 != "" {
			objectDataLocalMap1["destLogstore"] = nodeNative11
		}
		nodeNative12, _ := jsonpath.Get("$[0].role_arn", v)
		if nodeNative12 != nil && nodeNative12 != "" {
			objectDataLocalMap1["roleArn"] = nodeNative12
		}
		nodeNative13, _ := jsonpath.Get("$[0].dest_role_arn", v)
		if nodeNative13 != nil && nodeNative13 != "" {
			objectDataLocalMap1["destRoleArn"] = nodeNative13
		}
		nodeNative14, _ := jsonpath.Get("$[0].from_time_expr", v)
		if nodeNative14 != nil && nodeNative14 != "" {
			objectDataLocalMap1["fromTimeExpr"] = nodeNative14
		}
		nodeNative15, _ := jsonpath.Get("$[0].to_time_expr", v)
		if nodeNative15 != nil && nodeNative15 != "" {
			objectDataLocalMap1["toTimeExpr"] = nodeNative15
		}
		nodeNative16, _ := jsonpath.Get("$[0].max_run_time_in_seconds", v)
		if nodeNative16 != nil && nodeNative16 != "" {
			objectDataLocalMap1["maxRunTimeInSeconds"] = nodeNative16
		}
		nodeNative17, _ := jsonpath.Get("$[0].resource_pool", v)
		if nodeNative17 != nil && nodeNative17 != "" {
			objectDataLocalMap1["resourcePool"] = nodeNative17
		}
		nodeNative18, _ := jsonpath.Get("$[0].max_retries", v)
		if nodeNative18 != nil && nodeNative18 != "" {
			objectDataLocalMap1["maxRetries"] = nodeNative18
		}
		nodeNative19, _ := jsonpath.Get("$[0].from_time", v)
		if nodeNative19 != nil && nodeNative19 != "" {
			objectDataLocalMap1["fromTime"] = nodeNative19
		}
		nodeNative20, _ := jsonpath.Get("$[0].to_time", v)
		if nodeNative20 != nil && nodeNative20 != "" {
			objectDataLocalMap1["toTime"] = nodeNative20
		}
		nodeNative21, _ := jsonpath.Get("$[0].data_format", v)
		if nodeNative21 != nil && nodeNative21 != "" {
			objectDataLocalMap1["dataFormat"] = nodeNative21
		}
		nodeNative22, _ := jsonpath.Get("$[0].parameters", v)
		if nodeNative22 != nil && nodeNative22 != "" {
			objectDataLocalMap1["parameters"] = nodeNative22
		}

		request["configuration"] = objectDataLocalMap1
	}

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
			addDebug(action, response, request)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
	}

	if d.HasChange("status") {
		client := meta.(*connectivity.AliyunClient)
		slsServiceV2 := SlsServiceV2{client}
		object, err := slsServiceV2.DescribeSlsScheduledSQL(d.Id())
		if err != nil {
			return WrapError(err)
		}

		target := d.Get("status").(string)
		if object["status"].(string) != target {
			if target == "ENABLED" {
				parts = strings.Split(d.Id(), ":")
				scheduledSQLName = parts[1]
				action = fmt.Sprintf("/scheduledsqls/{scheduledSQLName}?action=enable")
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
					addDebug(action, response, request)
					return nil
				})
				if err != nil {
					return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
				}

			}
			if target == "DISABLED" {
				parts = strings.Split(d.Id(), ":")
				scheduledSQLName = parts[1]
				action = fmt.Sprintf("/scheduledsqls/{scheduledSQLName}?action=disable")
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
					addDebug(action, response, request)
					return nil
				})
				if err != nil {
					return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
				}

			}
		}
	}

	return resourceAliCloudSlsScheduledSQLRead(d, meta)
}

func resourceAliCloudSlsScheduledSQLDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	parts := strings.Split(d.Id(), ":")
	scheduledSQLName := parts[1]
	action := fmt.Sprintf("/scheduledsqls/%s", scheduledSQLName)
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	body := make(map[string]interface{})
	hostMap := make(map[string]*string)
	var err error
	request = make(map[string]interface{})
	hostMap["project"] = StringPointer(parts[0])

	body = request
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.Do("Sls", roaParam("DELETE", "2020-12-30", "DeleteScheduledSQL", action), query, body, nil, hostMap, false)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})

	if err != nil {
		if IsExpectedErrors(err, []string{"403"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return nil
}

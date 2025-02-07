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

func resourceAliCloudSchedulerxJob() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudSchedulerxJobCreate,
		Read:   resourceAliCloudSchedulerxJobRead,
		Update: resourceAliCloudSchedulerxJobUpdate,
		Delete: resourceAliCloudSchedulerxJobDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"attempt_interval": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"class_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"content": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"execute_mode": {
				Type:     schema.TypeString,
				Required: true,
			},
			"fail_times": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"group_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"job_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"job_monitor_info": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"contact_info": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"user_name": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"user_phone": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"user_mail": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"ding": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
						"monitor_config": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"timeout_enable": {
										Type:     schema.TypeBool,
										Optional: true,
										Computed: true,
									},
									"timeout": {
										Type:     schema.TypeInt,
										Optional: true,
										Computed: true,
									},
									"miss_worker_enable": {
										Type:     schema.TypeBool,
										Optional: true,
										Computed: true,
									},
									"fail_enable": {
										Type:     schema.TypeBool,
										Optional: true,
										Computed: true,
									},
									"send_channel": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
									"timeout_kill_enable": {
										Type:     schema.TypeBool,
										Optional: true,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
			"job_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"job_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"map_task_xattrs": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"task_attempt_interval": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"queue_size": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"dispatcher_size": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"page_size": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"task_max_attempt": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"consumer_size": {
							Type:     schema.TypeInt,
							Optional: true,
						},
					},
				},
			},
			"max_attempt": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"max_concurrency": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"namespace": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"namespace_source": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"parameters": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: StringInSlice([]string{"Enable", "Disable"}, false),
			},
			"success_notice_enable": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"task_dispatch_mode": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"template": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"time_config": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"calendar": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"time_expression": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"data_offset": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"time_type": {
							Type:     schema.TypeInt,
							Required: true,
						},
					},
				},
			},
			"timezone": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"x_attrs": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func resourceAliCloudSchedulerxJobCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateJob"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["Namespace"] = d.Get("namespace")
	request["GroupId"] = d.Get("group_id")
	request["RegionId"] = client.RegionId

	if v, ok := d.GetOk("job_monitor_info"); ok {
		jsonPathResult, err := jsonpath.Get("$[0].monitor_config[0].miss_worker_enable", v)
		if err == nil && jsonPathResult != "" {
			request["MissWorkerEnable"] = jsonPathResult
		}
	}
	jsonPathResult1, err := jsonpath.Get("$[0].time_expression", d.Get("time_config"))
	if err == nil {
		request["TimeExpression"] = jsonPathResult1
	}

	if v, ok := d.GetOk("x_attrs"); ok {
		request["XAttrs"] = v
	}
	if v, ok := d.GetOk("content"); ok {
		request["Content"] = v
	}
	if v, ok := d.GetOk("max_concurrency"); ok {
		request["MaxConcurrency"] = v
	}
	if v, ok := d.GetOk("class_name"); ok {
		request["ClassName"] = v
	}
	if v, ok := d.GetOk("job_monitor_info"); ok {
		localData, err := jsonpath.Get("$[0].contact_info", v)
		if err != nil {
			return WrapError(err)
		}
		contactInfoMapsArray := make([]interface{}, 0)
		for _, dataLoop := range localData.([]interface{}) {
			dataLoopTmp := dataLoop.(map[string]interface{})
			dataLoopMap := make(map[string]interface{})
			dataLoopMap["Ding"] = dataLoopTmp["ding"]
			dataLoopMap["UserName"] = dataLoopTmp["user_name"]
			dataLoopMap["UserPhone"] = dataLoopTmp["user_phone"]
			dataLoopMap["UserMail"] = dataLoopTmp["user_mail"]
			contactInfoMapsArray = append(contactInfoMapsArray, dataLoopMap)
		}
		request["ContactInfo"] = contactInfoMapsArray
	}

	request["ExecuteMode"] = d.Get("execute_mode")
	jsonPathResult7, err := jsonpath.Get("$[0].data_offset", d.Get("time_config"))
	if err == nil {
		request["DataOffset"] = jsonPathResult7
	}

	if v, ok := d.GetOk("timezone"); ok {
		request["Timezone"] = v
	}
	if v, ok := d.GetOk("job_monitor_info"); ok {
		jsonPathResult9, err := jsonpath.Get("$[0].monitor_config[0].timeout_enable", v)
		if err == nil && jsonPathResult9 != "" {
			request["TimeoutEnable"] = jsonPathResult9
		}
	}
	if v, ok := d.GetOkExists("success_notice_enable"); ok {
		request["SuccessNoticeEnable"] = v
	}
	if v, ok := d.GetOk("map_task_xattrs"); ok {
		jsonPathResult11, err := jsonpath.Get("$[0].dispatcher_size", v)
		if err == nil && jsonPathResult11 != "" {
			request["DispatcherSize"] = jsonPathResult11
		}
	}
	if v, ok := d.GetOk("job_monitor_info"); ok {
		jsonPathResult12, err := jsonpath.Get("$[0].monitor_config[0].timeout_kill_enable", v)
		if err == nil && jsonPathResult12 != "" {
			request["TimeoutKillEnable"] = jsonPathResult12
		}
	}
	if v, ok := d.GetOk("map_task_xattrs"); ok {
		jsonPathResult13, err := jsonpath.Get("$[0].page_size", v)
		if err == nil && jsonPathResult13 != "" {
			request["PageSize"] = jsonPathResult13
		}
	}
	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}
	if v, ok := d.GetOk("status"); ok {
		request["Status"] = convertSchedulerxJobStatusRequest(v.(string))
	}
	if v, ok := d.GetOk("parameters"); ok {
		request["Parameters"] = v
	}
	request["Name"] = d.Get("job_name")
	if v, ok := d.GetOkExists("attempt_interval"); ok {
		request["AttemptInterval"] = v
	}
	if v, ok := d.GetOk("namespace_source"); ok {
		request["NamespaceSource"] = v
	}
	jsonPathResult20, err := jsonpath.Get("$[0].calendar", d.Get("time_config"))
	if err == nil {
		request["Calendar"] = jsonPathResult20
	}

	if v, ok := d.GetOkExists("max_attempt"); ok {
		request["MaxAttempt"] = v
	}
	if v, ok := d.GetOk("map_task_xattrs"); ok {
		jsonPathResult22, err := jsonpath.Get("$[0].task_max_attempt", v)
		if err == nil && jsonPathResult22 != "" {
			request["TaskMaxAttempt"] = jsonPathResult22
		}
	}
	if v, ok := d.GetOk("map_task_xattrs"); ok {
		jsonPathResult23, err := jsonpath.Get("$[0].task_attempt_interval", v)
		if err == nil && jsonPathResult23 != "" {
			request["TaskAttemptInterval"] = jsonPathResult23
		}
	}
	if v, ok := d.GetOk("map_task_xattrs"); ok {
		jsonPathResult24, err := jsonpath.Get("$[0].consumer_size", v)
		if err == nil && jsonPathResult24 != "" {
			request["ConsumerSize"] = jsonPathResult24
		}
	}
	if v, ok := d.GetOkExists("fail_times"); ok {
		request["FailTimes"] = v
	}
	jsonPathResult26, err := jsonpath.Get("$[0].time_type", d.Get("time_config"))
	if err == nil {
		request["TimeType"] = jsonPathResult26
	}

	if v, ok := d.GetOk("job_monitor_info"); ok {
		jsonPathResult27, err := jsonpath.Get("$[0].monitor_config[0].timeout", v)
		if err == nil && jsonPathResult27 != "" {
			request["Timeout"] = jsonPathResult27
		}
	}
	request["JobType"] = d.Get("job_type")
	if v, ok := d.GetOk("job_monitor_info"); ok {
		jsonPathResult29, err := jsonpath.Get("$[0].monitor_config[0].fail_enable", v)
		if err == nil && jsonPathResult29 != "" {
			request["FailEnable"] = jsonPathResult29
		}
	}
	if v, ok := d.GetOk("job_monitor_info"); ok {
		jsonPathResult30, err := jsonpath.Get("$[0].monitor_config[0].send_channel", v)
		if err == nil && jsonPathResult30 != "" {
			request["SendChannel"] = jsonPathResult30
		}
	}
	if v, ok := d.GetOk("map_task_xattrs"); ok {
		jsonPathResult31, err := jsonpath.Get("$[0].queue_size", v)
		if err == nil && jsonPathResult31 != "" {
			request["QueueSize"] = jsonPathResult31
		}
	}
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("schedulerx2", "2019-04-30", action, query, request, true)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_schedulerx_job", action, AlibabaCloudSdkGoERROR)
	}

	DataJobIdVar, _ := jsonpath.Get("$.Data.JobId", response)
	d.SetId(fmt.Sprintf("%v:%v:%v", request["Namespace"], request["GroupId"], DataJobIdVar))

	return resourceAliCloudSchedulerxJobRead(d, meta)
}

func resourceAliCloudSchedulerxJobRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	schedulerxServiceV2 := SchedulerxServiceV2{client}

	objectRaw, err := schedulerxServiceV2.DescribeSchedulerxJob(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_schedulerx_job DescribeSchedulerxJob Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	if objectRaw["AttemptInterval"] != nil {
		d.Set("attempt_interval", objectRaw["AttemptInterval"])
	}
	if objectRaw["ClassName"] != nil {
		d.Set("class_name", objectRaw["ClassName"])
	}
	if objectRaw["Content"] != nil {
		d.Set("content", objectRaw["Content"])
	}
	if objectRaw["Description"] != nil {
		d.Set("description", objectRaw["Description"])
	}
	if objectRaw["ExecuteMode"] != nil {
		d.Set("execute_mode", objectRaw["ExecuteMode"])
	}
	if objectRaw["Name"] != nil {
		d.Set("job_name", objectRaw["Name"])
	}
	if objectRaw["JobType"] != nil {
		d.Set("job_type", objectRaw["JobType"])
	}
	if objectRaw["MaxAttempt"] != nil {
		d.Set("max_attempt", objectRaw["MaxAttempt"])
	}
	if objectRaw["MaxConcurrency"] != nil {
		d.Set("max_concurrency", objectRaw["MaxConcurrency"])
	}
	if objectRaw["Parameters"] != nil {
		d.Set("parameters", objectRaw["Parameters"])
	}
	if objectRaw["Status"] != nil {
		d.Set("status", convertSchedulerxJobDataJobConfigInfoStatusResponse(objectRaw["Status"]))
	}
	if objectRaw["XAttrs"] != nil {
		d.Set("x_attrs", objectRaw["XAttrs"])
	}
	if objectRaw["JobId"] != nil {
		d.Set("job_id", objectRaw["JobId"])
	}

	jobMonitorInfoMaps := make([]map[string]interface{}, 0)
	jobMonitorInfoMap := make(map[string]interface{})
	jobMonitorInfo1Raw := make(map[string]interface{})
	if objectRaw["JobMonitorInfo"] != nil {
		jobMonitorInfo1Raw = objectRaw["JobMonitorInfo"].(map[string]interface{})
	}
	if len(jobMonitorInfo1Raw) > 0 {

		contactInfo1Raw := jobMonitorInfo1Raw["ContactInfo"]
		contactInfoMaps := make([]map[string]interface{}, 0)
		if contactInfo1Raw != nil {
			for _, contactInfoChild1Raw := range contactInfo1Raw.([]interface{}) {
				contactInfoMap := make(map[string]interface{})
				contactInfoChild1Raw := contactInfoChild1Raw.(map[string]interface{})
				contactInfoMap["ding"] = contactInfoChild1Raw["Ding"]
				contactInfoMap["user_mail"] = contactInfoChild1Raw["UserMail"]
				contactInfoMap["user_name"] = contactInfoChild1Raw["UserName"]
				contactInfoMap["user_phone"] = contactInfoChild1Raw["UserPhone"]

				contactInfoMaps = append(contactInfoMaps, contactInfoMap)
			}
		}
		jobMonitorInfoMap["contact_info"] = contactInfoMaps
		monitorConfigMaps := make([]map[string]interface{}, 0)
		monitorConfigMap := make(map[string]interface{})
		monitorConfig1Raw := make(map[string]interface{})
		if jobMonitorInfo1Raw["MonitorConfig"] != nil {
			monitorConfig1Raw = jobMonitorInfo1Raw["MonitorConfig"].(map[string]interface{})
		}
		if len(monitorConfig1Raw) > 0 {
			monitorConfigMap["fail_enable"] = monitorConfig1Raw["FailEnable"]
			monitorConfigMap["miss_worker_enable"] = monitorConfig1Raw["MissWorkerEnable"]
			monitorConfigMap["send_channel"] = monitorConfig1Raw["SendChannel"]
			monitorConfigMap["timeout"] = monitorConfig1Raw["Timeout"]
			monitorConfigMap["timeout_enable"] = monitorConfig1Raw["TimeoutEnable"]
			monitorConfigMap["timeout_kill_enable"] = monitorConfig1Raw["TimeoutKillEnable"]

			monitorConfigMaps = append(monitorConfigMaps, monitorConfigMap)
		}
		jobMonitorInfoMap["monitor_config"] = monitorConfigMaps
		jobMonitorInfoMaps = append(jobMonitorInfoMaps, jobMonitorInfoMap)
	}
	if objectRaw["JobMonitorInfo"] != nil {
		if err := d.Set("job_monitor_info", jobMonitorInfoMaps); err != nil {
			return err
		}
	}
	mapTaskXattrsMaps := make([]map[string]interface{}, 0)
	mapTaskXattrsMap := make(map[string]interface{})
	mapTaskXAttrs1Raw := make(map[string]interface{})
	if objectRaw["MapTaskXAttrs"] != nil {
		mapTaskXAttrs1Raw = objectRaw["MapTaskXAttrs"].(map[string]interface{})
	}
	if len(mapTaskXAttrs1Raw) > 0 {
		mapTaskXattrsMap["consumer_size"] = mapTaskXAttrs1Raw["ConsumerSize"]
		mapTaskXattrsMap["dispatcher_size"] = mapTaskXAttrs1Raw["DispatcherSize"]
		mapTaskXattrsMap["page_size"] = mapTaskXAttrs1Raw["PageSize"]
		mapTaskXattrsMap["queue_size"] = mapTaskXAttrs1Raw["QueueSize"]
		mapTaskXattrsMap["task_attempt_interval"] = mapTaskXAttrs1Raw["TaskAttemptInterval"]
		mapTaskXattrsMap["task_max_attempt"] = mapTaskXAttrs1Raw["TaskMaxAttempt"]

		mapTaskXattrsMaps = append(mapTaskXattrsMaps, mapTaskXattrsMap)
	}
	if objectRaw["MapTaskXAttrs"] != nil {
		if err := d.Set("map_task_xattrs", mapTaskXattrsMaps); err != nil {
			return err
		}
	}
	timeConfigMaps := make([]map[string]interface{}, 0)
	timeConfigMap := make(map[string]interface{})
	timeConfig1Raw := make(map[string]interface{})
	if objectRaw["TimeConfig"] != nil {
		timeConfig1Raw = objectRaw["TimeConfig"].(map[string]interface{})
	}
	if len(timeConfig1Raw) > 0 {
		timeConfigMap["calendar"] = timeConfig1Raw["Calendar"]
		timeConfigMap["data_offset"] = timeConfig1Raw["DataOffset"]
		timeConfigMap["time_expression"] = timeConfig1Raw["TimeExpression"]
		timeConfigMap["time_type"] = timeConfig1Raw["TimeType"]

		timeConfigMaps = append(timeConfigMaps, timeConfigMap)
	}
	if objectRaw["TimeConfig"] != nil {
		if err := d.Set("time_config", timeConfigMaps); err != nil {
			return err
		}
	}

	parts := strings.Split(d.Id(), ":")
	d.Set("namespace", parts[0])
	d.Set("group_id", parts[1])

	return nil
}

func resourceAliCloudSchedulerxJobUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false

	if d.HasChange("status") {
		schedulerxServiceV2 := SchedulerxServiceV2{client}
		object, err := schedulerxServiceV2.DescribeSchedulerxJob(d.Id())
		if err != nil {
			return WrapError(err)
		}

		target := d.Get("status").(string)
		if convertSchedulerxJobDataJobConfigInfoStatusResponse(object["Status"]).(string) != target {
			if target == "Enable" {
				parts := strings.Split(d.Id(), ":")
				action := "EnableJob"
				request = make(map[string]interface{})
				query = make(map[string]interface{})
				query["Namespace"] = parts[0]
				query["JobId"] = parts[2]
				query["GroupId"] = parts[1]
				query["RegionId"] = client.RegionId
				if v, ok := d.GetOk("namespace_source"); ok {
					query["NamespaceSource"] = StringPointer(v.(string))
				}

				wait := incrementalWait(3*time.Second, 5*time.Second)
				err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
					response, err = client.RpcGet("schedulerx2", "2019-04-30", action, query, request)
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
			if target == "Disable" {
				parts := strings.Split(d.Id(), ":")
				action := "DisableJob"
				request = make(map[string]interface{})
				query = make(map[string]interface{})
				query["Namespace"] = parts[0]
				query["GroupId"] = parts[1]
				query["JobId"] = parts[2]
				query["RegionId"] = client.RegionId
				if v, ok := d.GetOk("namespace_source"); ok {
					query["NamespaceSource"] = StringPointer(v.(string))
				}

				wait := incrementalWait(3*time.Second, 5*time.Second)
				err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
					response, err = client.RpcGet("schedulerx2", "2019-04-30", action, query, request)
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

	parts := strings.Split(d.Id(), ":")
	action := "UpdateJob"
	var err error
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["Namespace"] = parts[0]
	request["GroupId"] = parts[1]
	request["JobId"] = parts[2]
	request["RegionId"] = client.RegionId
	if v, ok := d.GetOk("task_dispatch_mode"); ok {
		request["TaskDispatchMode"] = v
	}
	if d.HasChange("job_monitor_info.0.monitor_config.0.miss_worker_enable") {
		update = true
		jsonPathResult1, err := jsonpath.Get("$[0].monitor_config[0].miss_worker_enable", d.Get("job_monitor_info"))
		if err == nil {
			request["MissWorkerEnable"] = jsonPathResult1
		}
	}

	if d.HasChange("time_config.0.time_expression") {
		update = true
	}
	jsonPathResult2, err := jsonpath.Get("$[0].time_expression", d.Get("time_config"))
	if err == nil {
		request["TimeExpression"] = jsonPathResult2
	}

	if d.HasChange("x_attrs") {
		update = true
		request["XAttrs"] = d.Get("x_attrs")
	}

	if d.HasChange("content") {
		update = true
		request["Content"] = d.Get("content")
	}

	if d.HasChange("max_concurrency") {
		update = true
		request["MaxConcurrency"] = d.Get("max_concurrency")
	}

	if d.HasChange("class_name") {
		update = true
		request["ClassName"] = d.Get("class_name")
	}

	if d.HasChange("job_monitor_info") {
		update = true
		if v, ok := d.GetOk("job_monitor_info"); ok || d.HasChange("job_monitor_info") {
			localData, err := jsonpath.Get("$[0].contact_info", v)
			if err != nil {
				return WrapError(err)
			}
			contactInfoMapsArray := make([]interface{}, 0)
			for _, dataLoop := range localData.([]interface{}) {
				dataLoopTmp := dataLoop.(map[string]interface{})
				dataLoopMap := make(map[string]interface{})
				dataLoopMap["Ding"] = dataLoopTmp["ding"]
				dataLoopMap["UserName"] = dataLoopTmp["user_name"]
				dataLoopMap["UserPhone"] = dataLoopTmp["user_phone"]
				dataLoopMap["UserMail"] = dataLoopTmp["user_mail"]
				contactInfoMapsArray = append(contactInfoMapsArray, dataLoopMap)
			}
			request["ContactInfo"] = contactInfoMapsArray
		}
	}

	if v, ok := d.GetOk("template"); ok {
		request["Template"] = v
	}
	if d.HasChange("execute_mode") {
		update = true
	}
	request["ExecuteMode"] = d.Get("execute_mode")
	if d.HasChange("time_config.0.data_offset") {
		update = true
	}
	jsonPathResult9, err := jsonpath.Get("$[0].data_offset", d.Get("time_config"))
	if err == nil {
		request["DataOffset"] = jsonPathResult9
	}

	if d.HasChange("job_monitor_info.0.monitor_config.0.timeout_enable") {
		update = true
		jsonPathResult10, err := jsonpath.Get("$[0].monitor_config[0].timeout_enable", d.Get("job_monitor_info"))
		if err == nil {
			request["TimeoutEnable"] = jsonPathResult10
		}
	}

	if v, ok := d.GetOk("timezone"); ok {
		request["Timezone"] = v
	}
	if v, ok := d.GetOkExists("success_notice_enable"); ok {
		request["SuccessNoticeEnable"] = v
	}
	if d.HasChange("map_task_xattrs.0.dispatcher_size") {
		update = true
		jsonPathResult13, err := jsonpath.Get("$[0].dispatcher_size", d.Get("map_task_xattrs"))
		if err == nil {
			request["DispatcherSize"] = jsonPathResult13
		}
	}

	if d.HasChange("job_monitor_info.0.monitor_config.0.timeout_kill_enable") {
		update = true
		jsonPathResult14, err := jsonpath.Get("$[0].monitor_config[0].timeout_kill_enable", d.Get("job_monitor_info"))
		if err == nil {
			request["TimeoutKillEnable"] = jsonPathResult14
		}
	}

	if d.HasChange("map_task_xattrs.0.page_size") {
		update = true
		jsonPathResult15, err := jsonpath.Get("$[0].page_size", d.Get("map_task_xattrs"))
		if err == nil {
			request["PageSize"] = jsonPathResult15
		}
	}

	if d.HasChange("description") {
		update = true
		request["Description"] = d.Get("description")
	}

	if d.HasChange("parameters") {
		update = true
		request["Parameters"] = d.Get("parameters")
	}

	if d.HasChange("job_name") {
		update = true
	}
	request["Name"] = d.Get("job_name")
	if d.HasChange("attempt_interval") {
		update = true
		request["AttemptInterval"] = d.Get("attempt_interval")
	}

	if v, ok := d.GetOk("namespace_source"); ok {
		request["NamespaceSource"] = v
	}
	if d.HasChange("time_config.0.calendar") {
		update = true
	}
	jsonPathResult21, err := jsonpath.Get("$[0].calendar", d.Get("time_config"))
	if err == nil {
		request["Calendar"] = jsonPathResult21
	}

	if d.HasChange("max_attempt") {
		update = true
		request["MaxAttempt"] = d.Get("max_attempt")
	}

	if d.HasChange("map_task_xattrs.0.task_max_attempt") {
		update = true
		jsonPathResult23, err := jsonpath.Get("$[0].task_max_attempt", d.Get("map_task_xattrs"))
		if err == nil {
			request["TaskMaxAttempt"] = jsonPathResult23
		}
	}

	if d.HasChange("map_task_xattrs.0.task_attempt_interval") {
		update = true
		jsonPathResult24, err := jsonpath.Get("$[0].task_attempt_interval", d.Get("map_task_xattrs"))
		if err == nil {
			request["TaskAttemptInterval"] = jsonPathResult24
		}
	}

	if v, ok := d.GetOkExists("fail_times"); ok {
		request["FailTimes"] = v
	}
	if d.HasChange("map_task_xattrs.0.consumer_size") {
		update = true
		jsonPathResult26, err := jsonpath.Get("$[0].consumer_size", d.Get("map_task_xattrs"))
		if err == nil {
			request["ConsumerSize"] = jsonPathResult26
		}
	}

	if d.HasChange("time_config.0.time_type") {
		update = true
	}
	jsonPathResult27, err := jsonpath.Get("$[0].time_type", d.Get("time_config"))
	if err == nil {
		request["TimeType"] = jsonPathResult27
	}

	if d.HasChange("job_monitor_info.0.monitor_config.0.timeout") {
		update = true
		jsonPathResult28, err := jsonpath.Get("$[0].monitor_config[0].timeout", d.Get("job_monitor_info"))
		if err == nil {
			request["Timeout"] = jsonPathResult28
		}
	}

	if d.HasChange("job_monitor_info.0.monitor_config.0.fail_enable") {
		update = true
		jsonPathResult29, err := jsonpath.Get("$[0].monitor_config[0].fail_enable", d.Get("job_monitor_info"))
		if err == nil {
			request["FailEnable"] = jsonPathResult29
		}
	}

	if d.HasChange("job_monitor_info.0.monitor_config.0.send_channel") {
		update = true
		jsonPathResult30, err := jsonpath.Get("$[0].monitor_config[0].send_channel", d.Get("job_monitor_info"))
		if err == nil {
			request["SendChannel"] = jsonPathResult30
		}
	}

	if d.HasChange("map_task_xattrs.0.queue_size") {
		update = true
		jsonPathResult31, err := jsonpath.Get("$[0].queue_size", d.Get("map_task_xattrs"))
		if err == nil {
			request["QueueSize"] = jsonPathResult31
		}
	}

	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("schedulerx2", "2019-04-30", action, query, request, true)
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

	return resourceAliCloudSchedulerxJobRead(d, meta)
}

func resourceAliCloudSchedulerxJobDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	parts := strings.Split(d.Id(), ":")
	action := "DeleteJob"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	query["Namespace"] = parts[0]
	query["GroupId"] = parts[1]
	query["JobId"] = parts[2]
	query["RegionId"] = client.RegionId

	if v, ok := d.GetOk("namespace_source"); ok {
		query["NamespaceSource"] = StringPointer(v.(string))
	}

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcGet("schedulerx2", "2019-04-30", action, query, request)

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
		if NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return nil
}

func convertSchedulerxJobDataJobConfigInfoStatusResponse(source interface{}) interface{} {
	source = fmt.Sprint(source)
	switch source {
	case "1":
		return "Enable"
	case "0":
		return "Disable"
	}
	return source
}
func convertSchedulerxJobStatusRequest(source interface{}) interface{} {
	source = fmt.Sprint(source)
	switch source {
	case "Enable":
		return "1"
	case "Disable":
		return "0"
	}
	return source
}

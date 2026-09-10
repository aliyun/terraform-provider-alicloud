// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudSchedulerxAppGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudSchedulerxAppGroupCreate,
		Read:   resourceAliCloudSchedulerxAppGroupRead,
		Update: resourceAliCloudSchedulerxAppGroupUpdate,
		Delete: resourceAliCloudSchedulerxAppGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"app_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"app_type": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"app_version": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"delete_jobs": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enable_log": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"group_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"max_concurrency": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"max_jobs": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"monitor_config_json": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"monitor_contacts_json": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"namespace": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"namespace_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"namespace_source": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"schedule_busy_workers": {
				Type:     schema.TypeBool,
				Optional: true,
			},
		},
	}
}

func resourceAliCloudSchedulerxAppGroupCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateAppGroup"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	if v, ok := d.GetOk("namespace"); ok {
		query["Namespace"] = v
	}
	if v, ok := d.GetOk("group_id"); ok {
		query["GroupId"] = v
	}
	query["RegionId"] = client.RegionId

	if v, ok := d.GetOk("description"); ok {
		query["Description"] = StringPointer(v.(string))
	}

	if v, ok := d.GetOk("monitor_config_json"); ok {
		query["MonitorConfigJson"] = StringPointer(v.(string))
	}

	if v, ok := d.GetOk("max_jobs"); ok {
		query["MaxJobs"] = StringPointer(strconv.Itoa(v.(int)))
	}

	if v, ok := d.GetOk("app_type"); ok {
		query["AppType"] = StringPointer(strconv.Itoa(v.(int)))
	}

	if v, ok := d.GetOk("schedule_busy_workers"); ok {
		query["ScheduleBusyWorkers"] = StringPointer(strconv.FormatBool(v.(bool)))
	}

	if v, ok := d.GetOk("namespace_source"); ok {
		query["NamespaceSource"] = StringPointer(v.(string))
	}

	if v, ok := d.GetOk("namespace_name"); ok {
		query["NamespaceName"] = StringPointer(v.(string))
	}

	if v, ok := d.GetOkExists("enable_log"); ok {
		query["EnableLog"] = StringPointer(strconv.FormatBool(v.(bool)))
	}

	if v, ok := d.GetOk("app_name"); ok {
		query["AppName"] = StringPointer(v.(string))
	}

	if v, ok := d.GetOk("app_version"); ok {
		query["AppVersion"] = StringPointer(v.(string))
	}

	if v, ok := d.GetOk("monitor_contacts_json"); ok {
		query["MonitorContactsJson"] = StringPointer(v.(string))
	}

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_schedulerx_app_group", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%v:%v", query["Namespace"], query["GroupId"]))

	return resourceAliCloudSchedulerxAppGroupRead(d, meta)
}

func resourceAliCloudSchedulerxAppGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	schedulerxServiceV2 := SchedulerxServiceV2{client}

	objectRaw, err := schedulerxServiceV2.DescribeSchedulerxAppGroup(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_schedulerx_app_group DescribeSchedulerxAppGroup Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("app_name", objectRaw["AppName"])
	d.Set("app_version", objectRaw["AppVersion"])
	d.Set("description", objectRaw["Description"])
	d.Set("enable_log", objectRaw["EnableLog"])
	d.Set("max_jobs", objectRaw["MaxJobs"])
	d.Set("monitor_config_json", objectRaw["MonitorConfigJson"])
	d.Set("monitor_contacts_json", objectRaw["MonitorContactsJson"])
	d.Set("group_id", objectRaw["GroupId"])
	d.Set("namespace", objectRaw["Namespace"])

	return nil
}

func resourceAliCloudSchedulerxAppGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false

	var err error
	parts := strings.Split(d.Id(), ":")
	action := "UpdateAppGroup"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["Namespace"] = parts[0]
	request["GroupId"] = parts[1]
	request["RegionId"] = client.RegionId
	if d.HasChange("description") {
		update = true
		request["Description"] = d.Get("description")
	}

	if d.HasChange("monitor_config_json") {
		update = true
	}
	if v, ok := d.GetOk("monitor_config_json"); ok {
		request["MonitorConfigJson"] = v
	}

	if d.HasChange("enable_log") {
		update = true
	}
	if v, ok := d.GetOkExists("enable_log"); ok {
		request["EnableLog"] = v
	}

	if d.HasChange("app_version") {
		update = true
	}
	if v, ok := d.GetOk("app_version"); ok {
		request["AppVersion"] = v
	}

	if d.HasChange("monitor_contacts_json") {
		update = true
	}
	if v, ok := d.GetOk("monitor_contacts_json"); ok {
		request["MonitorContactsJson"] = v
	}

	if v, ok := d.GetOkExists("max_concurrency"); ok {
		request["MaxConcurrency"] = v
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

	return resourceAliCloudSchedulerxAppGroupRead(d, meta)
}

func resourceAliCloudSchedulerxAppGroupDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	parts := strings.Split(d.Id(), ":")
	action := "DeleteAppGroup"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["Namespace"] = parts[0]
	request["GroupId"] = parts[1]
	request["RegionId"] = client.RegionId

	if v, ok := d.GetOkExists("delete_jobs"); ok {
		request["DeleteJobs"] = v
	}
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
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
		if NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return nil
}

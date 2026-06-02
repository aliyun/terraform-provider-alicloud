package alicloud

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAliCloudCmsAggTaskGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudCmsAggTaskGroupCreate,
		Read:   resourceAliCloudCmsAggTaskGroupRead,
		Update: resourceAliCloudCmsAggTaskGroupUpdate,
		Delete: resourceAliCloudCmsAggTaskGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"agg_task_group_config": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateYamlString,
				StateFunc: func(v interface{}) string {
					yaml, _ := normalizeYamlString(v)
					return yaml
				},
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if d.Get("cron_expr").(string) != "" {
						return d.Id() != ""
					}
					equal, _ := compareYamlTemplateAreEquivalent(old, new)
					return equal
				},
			},
			"agg_task_group_config_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"agg_task_group_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"agg_task_group_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"cron_expr": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"delay": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"max_retries": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"max_run_time_in_seconds": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"override_if_exists": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"precheck_string": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringIsJSON,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					equal, _ := compareJsonTemplateAreEquivalent(old, new)
					return equal
				},
			},
			"region_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"schedule_mode": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"schedule_time_expr": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"source_prometheus_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: StringInSlice([]string{"Running", "Stopped"}, false),
			},
			"target_prometheus_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"to_time": {
				Type:     schema.TypeInt,
				Optional: true,
			},
		},
	}
}

func resourceAliCloudCmsAggTaskGroupCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	instanceId := d.Get("source_prometheus_id")
	action := fmt.Sprintf("/prometheus-instances/%s/agg-task-groups", instanceId)
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	body := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	status := "Running"

	if v, ok := d.GetOk("schedule_time_expr"); ok {
		request["scheduleTimeExpr"] = v
	}
	if v, ok := d.GetOk("cron_expr"); ok {
		request["cronExpr"] = v
	}
	if v, ok := d.GetOk("description"); ok {
		request["description"] = v
	}
	if v, ok := d.GetOk("status"); ok {
		request["status"] = v
		status = v.(string)
	}
	request["aggTaskGroupConfig"] = d.Get("agg_task_group_config")
	if v, ok := d.GetOk("schedule_mode"); ok {
		request["scheduleMode"] = v
	}
	if v, ok := d.GetOkExists("override_if_exists"); ok {
		query["overrideIfExists"] = StringPointer(strconv.FormatBool(v.(bool)))
	}

	request["aggTaskGroupName"] = d.Get("agg_task_group_name")
	if v, ok := d.GetOkExists("delay"); ok {
		request["delay"] = v
	}
	if v, ok := d.GetOkExists("max_run_time_in_seconds"); ok {
		request["maxRunTimeInSeconds"] = v
	}
	request["targetPrometheusId"] = d.Get("target_prometheus_id")
	if v, ok := d.GetOkExists("to_time"); ok {
		request["toTime"] = v
	}
	if v, ok := d.GetOkExists("max_retries"); ok {
		request["maxRetries"] = v
	}
	if v, ok := d.GetOk("agg_task_group_config_type"); ok {
		request["aggTaskGroupConfigType"] = v
	}
	if v, ok := d.GetOk("precheck_string"); ok {
		request["precheckString"] = v
	}
	body = request
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RoaPost("Cms", "2024-03-30", action, query, nil, body, true)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_cms_agg_task_group", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%v:%v", response["sourcePrometheusId"], response["aggTaskGroupId"]))

	cmsServiceV2 := CmsServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{status}, d.Timeout(schema.TimeoutCreate), 5*time.Second, cmsServiceV2.CmsAggTaskGroupStateRefreshFunc(d.Id(), "status", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAliCloudCmsAggTaskGroupRead(d, meta)
}

func resourceAliCloudCmsAggTaskGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cmsServiceV2 := CmsServiceV2{client}

	objectRaw, err := cmsServiceV2.DescribeCmsAggTaskGroup(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_cms_agg_task_group DescribeCmsAggTaskGroup Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("agg_task_group_config", objectRaw["aggTaskGroupConfig"])
	d.Set("agg_task_group_name", objectRaw["aggTaskGroupName"])
	d.Set("cron_expr", objectRaw["cronExpr"])
	d.Set("delay", objectRaw["delay"])
	d.Set("description", objectRaw["description"])
	d.Set("max_retries", objectRaw["maxRetries"])
	d.Set("max_run_time_in_seconds", objectRaw["maxRunTimeInSeconds"])
	d.Set("precheck_string", objectRaw["precheckString"])
	d.Set("region_id", objectRaw["regionId"])
	d.Set("schedule_mode", objectRaw["scheduleMode"])
	d.Set("schedule_time_expr", objectRaw["scheduleTimeExpr"])
	d.Set("status", objectRaw["status"])
	d.Set("target_prometheus_id", objectRaw["targetPrometheusId"])
	d.Set("to_time", objectRaw["toTime"])
	d.Set("agg_task_group_id", objectRaw["aggTaskGroupId"])
	d.Set("source_prometheus_id", objectRaw["sourcePrometheusId"])

	return nil
}

func resourceAliCloudCmsAggTaskGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	var body map[string]interface{}
	update := false

	var err error
	parts := strings.Split(d.Id(), ":")
	instanceId := parts[0]
	groupId := parts[1]
	action := fmt.Sprintf("/prometheus-instances/%s/agg-task-groups/%s", instanceId, groupId)
	request = make(map[string]interface{})
	query = make(map[string]*string)
	body = make(map[string]interface{})
	status := "Running"

	if d.HasChange("schedule_time_expr") {
		update = true
	}
	if v, ok := d.GetOk("schedule_time_expr"); ok || d.HasChange("schedule_time_expr") {
		request["scheduleTimeExpr"] = v
	}
	if d.HasChange("cron_expr") {
		update = true
	}
	if v, ok := d.GetOk("cron_expr"); ok || d.HasChange("cron_expr") {
		request["cronExpr"] = v
	}
	if d.HasChange("description") {
		update = true
	}
	if v, ok := d.GetOk("description"); ok || d.HasChange("description") {
		request["description"] = v
	}
	if d.HasChange("status") {
		update = true
	}
	if v, ok := d.GetOk("status"); ok || d.HasChange("status") {
		request["status"] = v
		status = v.(string)
	}
	if d.HasChange("agg_task_group_config") {
		update = true
	}
	request["aggTaskGroupConfig"] = d.Get("agg_task_group_config")
	if d.HasChange("schedule_mode") {
		update = true
	}
	if v, ok := d.GetOk("schedule_mode"); ok || d.HasChange("schedule_mode") {
		request["scheduleMode"] = v
	}
	if d.HasChange("agg_task_group_name") {
		update = true
	}
	request["aggTaskGroupName"] = d.Get("agg_task_group_name")
	if d.HasChange("delay") {
		update = true
	}
	if v, ok := d.GetOkExists("delay"); ok || d.HasChange("delay") {
		request["delay"] = v
	}
	if d.HasChange("max_run_time_in_seconds") {
		update = true
	}
	if v, ok := d.GetOkExists("max_run_time_in_seconds"); ok || d.HasChange("max_run_time_in_seconds") {
		request["maxRunTimeInSeconds"] = v
	}
	if d.HasChange("target_prometheus_id") {
		update = true
	}
	request["targetPrometheusId"] = d.Get("target_prometheus_id")
	if d.HasChange("to_time") {
		update = true
	}
	if v, ok := d.GetOkExists("to_time"); ok || d.HasChange("to_time") {
		request["toTime"] = v
	}
	if d.HasChange("max_retries") {
		update = true
	}
	if v, ok := d.GetOkExists("max_retries"); ok || d.HasChange("max_retries") {
		request["maxRetries"] = v
	}
	if v, ok := d.GetOk("agg_task_group_config_type"); ok {
		request["aggTaskGroupConfigType"] = v
	}
	if d.HasChange("precheck_string") {
		update = true
	}
	if v, ok := d.GetOk("precheck_string"); ok || d.HasChange("precheck_string") {
		request["precheckString"] = v
	}
	body = request
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RoaPut("Cms", "2024-03-30", action, query, nil, body, true)
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
		cmsServiceV2 := CmsServiceV2{client}
		stateConf := BuildStateConf([]string{}, []string{status}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, cmsServiceV2.CmsAggTaskGroupStateRefreshFunc(d.Id(), "status", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}

	return resourceAliCloudCmsAggTaskGroupRead(d, meta)
}

func resourceAliCloudCmsAggTaskGroupDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	parts := strings.Split(d.Id(), ":")
	instanceId := parts[0]
	groupId := parts[1]
	action := fmt.Sprintf("/prometheus-instances/%s/agg-task-groups/%s", instanceId, groupId)
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	var err error
	request = make(map[string]interface{})

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RoaDelete("Cms", "2024-03-30", action, query, nil, nil, true)
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

	cmsServiceV2 := CmsServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{""}, d.Timeout(schema.TimeoutDelete), 30*time.Second, cmsServiceV2.CmsAggTaskGroupStateRefreshFunc(d.Id(), "aggTaskGroupId", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return nil
}

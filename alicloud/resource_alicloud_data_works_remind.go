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
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAlicloudDataWorksRemind() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudDataWorksRemindCreate,
		Read:   resourceAlicloudDataWorksRemindRead,
		Update: resourceAlicloudDataWorksRemindUpdate,
		Delete: resourceAlicloudDataWorksRemindDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"remind_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"remind_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"remind_unit": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"NODE", "BASELINE", "PROJECT", "BIZPROCESS"}, false),
			},
			"remind_type": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"FINISHED", "UNFINISHED", "ERROR", "CYCLE_UNFINISHED", "TIMEOUT",
				}, false),
			},
			"alert_unit": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"OWNER", "OTHER"}, false),
			},
			"dnd_end": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"alert_methods": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"nodes": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"node_id": {
							Type:     schema.TypeInt,
							Required: true,
						},
					},
				},
			},
			"baselines": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"baseline_id": {
							Type:     schema.TypeInt,
							Required: true,
						},
					},
				},
			},
			"alert_targets": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"useflag": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"biz_processes": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"biz_process_id": {
							Type:     schema.TypeInt,
							Required: true,
						},
					},
				},
			},
			"max_alert_times": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  3,
			},
			"alert_interval": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  3,
			},
			"detail": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"robots": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"web_url": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"webhooks": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"projects": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"project_id": {
							Type:     schema.TypeInt,
							Required: true,
						},
					},
				},
			},
			"founder": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"dnd_start": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"region_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAlicloudDataWorksRemindCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateRemind"
	request := make(map[string]interface{})
	var err error

	request["RemindName"] = d.Get("remind_name")
	request["RemindUnit"] = d.Get("remind_unit")
	request["RemindType"] = d.Get("remind_type")
	if v, ok := d.GetOk("alert_unit"); ok {
		request["AlertUnit"] = v
	}
	if v, ok := d.GetOk("dnd_end"); ok {
		request["DndEnd"] = v
	}
	if v, ok := d.GetOk("max_alert_times"); ok {
		request["MaxAlertTimes"] = v
	}
	if v, ok := d.GetOk("alert_interval"); ok {
		request["AlertInterval"] = v
	}
	if v, ok := d.GetOk("detail"); ok {
		request["Detail"] = v
	}
	request["UseFlag"] = d.Get("useflag")

	if v, ok := d.GetOk("alert_methods"); ok {
		request["AlertMethods"] = strings.Join(expandStringList(v.([]interface{})), ",")
	}
	if v, ok := d.GetOk("alert_targets"); ok {
		request["AlertTargets"] = strings.Join(expandStringList(v.([]interface{})), ",")
	}
	if v, ok := d.GetOk("webhooks"); ok {
		request["Webhooks"] = strings.Join(expandStringList(v.([]interface{})), ",")
	}
	if v, ok := d.GetOk("nodes"); ok {
		request["NodeIds"] = joinIntField(v.([]interface{}), "node_id")
	}
	if v, ok := d.GetOk("baselines"); ok {
		request["BaselineIds"] = joinIntField(v.([]interface{}), "baseline_id")
	}
	if v, ok := d.GetOk("biz_processes"); ok {
		request["BizProcessIds"] = joinIntField(v.([]interface{}), "biz_process_id")
	}
	if v, ok := d.GetOk("robots"); ok {
		request["RobotUrls"] = joinStringField(v.([]interface{}), "web_url")
	}
	if v, ok := d.GetOk("projects"); ok {
		ids := joinIntField(v.([]interface{}), "project_id")
		if ids != "" {
			request["ProjectId"] = ids
		}
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("dataworks-public", "2020-05-18", action, nil, request, false)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_data_works_remind", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["Data"]))

	return resourceAlicloudDataWorksRemindRead(d, meta)
}

func resourceAlicloudDataWorksRemindRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	dataworksPublicService := DataworksPublicService{client}
	object, err := dataworksPublicService.DescribeDataWorksRemind(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_data_works_remind DescribeDataWorksRemind Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("remind_id", fmt.Sprint(object["remindId"]))
	d.Set("remind_name", object["remindName"])
	d.Set("remind_unit", object["remindUnit"])
	d.Set("remind_type", object["remindType"])
	d.Set("alert_unit", object["alertUnit"])
	d.Set("dnd_end", object["dndEnd"])
	d.Set("dnd_start", object["dndStart"])
	d.Set("founder", object["founder"])
	d.Set("detail", object["detail"])
	d.Set("max_alert_times", object["maxAlertTimes"])
	d.Set("alert_interval", object["alertInterval"])
	d.Set("useflag", object["useflag"])
	d.Set("region_id", object["regionId"])

	if v, err := jsonpath.Get("$.alertMethods", object); err == nil {
		d.Set("alert_methods", expandStringList(v.([]interface{})))
	}
	if v, err := jsonpath.Get("$.alertTargets", object); err == nil {
		d.Set("alert_targets", expandStringList(v.([]interface{})))
	}
	if v, err := jsonpath.Get("$.webhooks", object); err == nil {
		d.Set("webhooks", expandStringList(v.([]interface{})))
	}
	if v, err := jsonpath.Get("$.nodes", object); err == nil {
		d.Set("nodes", flattenRemindNodes(v.([]interface{})))
	}
	if v, err := jsonpath.Get("$.baselines", object); err == nil {
		d.Set("baselines", flattenRemindBaselines(v.([]interface{})))
	}
	if v, err := jsonpath.Get("$.bizProcesses", object); err == nil {
		d.Set("biz_processes", flattenRemindBizProcesses(v.([]interface{})))
	}
	if v, err := jsonpath.Get("$.robots", object); err == nil {
		d.Set("robots", flattenRemindRobots(v.([]interface{})))
	}
	if v, err := jsonpath.Get("$.projects", object); err == nil {
		d.Set("projects", flattenRemindProjects(v.([]interface{})))
	}

	return nil
}

func resourceAlicloudDataWorksRemindUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "UpdateRemind"
	request := make(map[string]interface{})
	var err error

	request["RemindId"] = d.Id()
	if d.HasChange("remind_name") {
		request["RemindName"] = d.Get("remind_name")
	}
	if d.HasChange("remind_unit") {
		request["RemindUnit"] = d.Get("remind_unit")
	}
	if d.HasChange("remind_type") {
		request["RemindType"] = d.Get("remind_type")
	}
	if d.HasChange("alert_unit") {
		request["AlertUnit"] = d.Get("alert_unit")
	}
	if d.HasChange("dnd_end") {
		request["DndEnd"] = d.Get("dnd_end")
	}
	if d.HasChange("max_alert_times") {
		request["MaxAlertTimes"] = d.Get("max_alert_times")
	}
	if d.HasChange("alert_interval") {
		request["AlertInterval"] = d.Get("alert_interval")
	}
	if d.HasChange("detail") {
		request["Detail"] = d.Get("detail")
	}
	if d.HasChange("useflag") {
		request["UseFlag"] = d.Get("useflag")
	}
	if d.HasChange("alert_methods") {
		request["AlertMethods"] = strings.Join(expandStringList(d.Get("alert_methods").([]interface{})), ",")
	}
	if d.HasChange("alert_targets") {
		request["AlertTargets"] = strings.Join(expandStringList(d.Get("alert_targets").([]interface{})), ",")
	}
	if d.HasChange("webhooks") {
		request["Webhooks"] = strings.Join(expandStringList(d.Get("webhooks").([]interface{})), ",")
	}
	if d.HasChange("nodes") {
		request["NodeIds"] = joinIntField(d.Get("nodes").([]interface{}), "node_id")
	}
	if d.HasChange("baselines") {
		request["BaselineIds"] = joinIntField(d.Get("baselines").([]interface{}), "baseline_id")
	}
	if d.HasChange("biz_processes") {
		request["BizProcessIds"] = joinIntField(d.Get("biz_processes").([]interface{}), "biz_process_id")
	}
	if d.HasChange("robots") {
		request["RobotUrls"] = joinStringField(d.Get("robots").([]interface{}), "web_url")
	}
	if d.HasChange("projects") {
		ids := joinIntField(d.Get("projects").([]interface{}), "project_id")
		if ids != "" {
			request["ProjectId"] = ids
		}
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
		response, err = client.RpcPost("dataworks-public", "2020-05-18", action, nil, request, false)
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

	return resourceAlicloudDataWorksRemindRead(d, meta)
}

func resourceAlicloudDataWorksRemindDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "DeleteRemind"
	request := map[string]interface{}{
		"RemindId": d.Id(),
	}
	var err error

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("dataworks-public", "2020-05-18", action, nil, request, false)
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
	return nil
}

func joinIntField(items []interface{}, field string) string {
	ids := make([]string, 0, len(items))
	for _, item := range items {
		m := item.(map[string]interface{})
		ids = append(ids, fmt.Sprint(m[field]))
	}
	return strings.Join(ids, ",")
}

func joinStringField(items []interface{}, field string) string {
	urls := make([]string, 0, len(items))
	for _, item := range items {
		m := item.(map[string]interface{})
		urls = append(urls, m[field].(string))
	}
	return strings.Join(urls, ",")
}

func flattenRemindNodes(items []interface{}) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(items))
	for _, item := range items {
		m := item.(map[string]interface{})
		result = append(result, map[string]interface{}{
			"node_id": m["nodeId"],
		})
	}
	return result
}

func flattenRemindBaselines(items []interface{}) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(items))
	for _, item := range items {
		m := item.(map[string]interface{})
		result = append(result, map[string]interface{}{
			"baseline_id": m["baselineId"],
		})
	}
	return result
}

func flattenRemindBizProcesses(items []interface{}) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(items))
	for _, item := range items {
		m := item.(map[string]interface{})
		result = append(result, map[string]interface{}{
			"biz_process_id": m["bizId"],
		})
	}
	return result
}

func flattenRemindRobots(items []interface{}) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(items))
	for _, item := range items {
		m := item.(map[string]interface{})
		result = append(result, map[string]interface{}{
			"web_url": m["webUrl"],
		})
	}
	return result
}

func flattenRemindProjects(items []interface{}) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(items))
	for _, item := range items {
		m := item.(map[string]interface{})
		result = append(result, map[string]interface{}{
			"project_id": m["projectId"],
		})
	}
	return result
}

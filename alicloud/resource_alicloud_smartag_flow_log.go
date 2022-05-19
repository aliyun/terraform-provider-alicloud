package alicloud

import (
	"fmt"
	"log"
	"time"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAlicloudSmartagFlowLog() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudSmartagFlowLogCreate,
		Read:   resourceAlicloudSmartagFlowLogRead,
		Update: resourceAlicloudSmartagFlowLogUpdate,
		Delete: resourceAlicloudSmartagFlowLogDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(1 * time.Minute),
			Delete: schema.DefaultTimeout(1 * time.Minute),
			Update: schema.DefaultTimeout(1 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"active_aging": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntBetween(60, 6000),
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"flow_log_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"inactive_aging": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntBetween(10, 600),
			},
			"logstore_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"netflow_server_ip": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"netflow_server_port": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"netflow_version": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"V10", "V5", "V9"}, false),
			},
			"output_type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"all", "netflow", "sls"}, false),
			},
			"project_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"sls_region_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"Active", "Inactive"}, false),
			},
		},
	}
}

func resourceAlicloudSmartagFlowLogCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateFlowLog"
	request := make(map[string]interface{})
	conn, err := client.NewSmartagClient()
	if err != nil {
		return WrapError(err)
	}
	if v, ok := d.GetOk("active_aging"); ok {
		request["ActiveAging"] = v
	}
	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}
	if v, ok := d.GetOk("flow_log_name"); ok {
		request["Name"] = v
	}
	if v, ok := d.GetOk("inactive_aging"); ok {
		request["InactiveAging"] = v
	}
	request["OutputType"] = d.Get("output_type")
	if v, ok := d.GetOk("logstore_name"); ok {
		request["LogstoreName"] = v
	}
	if v, ok := d.GetOk("netflow_server_ip"); ok {
		request["NetflowServerIp"] = v
	}
	if v, ok := d.GetOk("netflow_server_port"); ok {
		request["NetflowServerPort"] = v
	}
	if v, ok := d.GetOk("project_name"); ok {
		request["ProjectName"] = v
	}
	if v, ok := d.GetOk("netflow_version"); ok {
		request["NetflowVersion"] = v
	}

	request["RegionId"] = client.RegionId

	if v, ok := d.GetOk("sls_region_id"); ok {
		request["SlsRegionId"] = v
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2018-03-13"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_smartag_flow_log", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["FlowLogId"]))

	return resourceAlicloudSmartagFlowLogUpdate(d, meta)
}
func resourceAlicloudSmartagFlowLogRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	SagService := SagService{client}
	object, err := SagService.DescribeSmartagFlowLog(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_smartag_flow_log SagService.DescribeSmartagFlowLog Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	if v, ok := object["ActiveAging"]; ok && fmt.Sprint(v) != "0" {
		d.Set("active_aging", formatInt(v))
	}
	d.Set("description", object["Description"])
	d.Set("flow_log_name", object["Name"])
	if v, ok := object["InactiveAging"]; ok && fmt.Sprint(v) != "0" {
		d.Set("inactive_aging", formatInt(v))
	}
	d.Set("logstore_name", object["LogstoreName"])
	d.Set("netflow_server_ip", object["NetflowServerIp"])
	if v, ok := object["NetflowServerPort"]; ok {
		d.Set("netflow_server_port", formatInt(v))
	}
	d.Set("netflow_version", object["NetflowVersion"])
	d.Set("output_type", object["OutputType"])
	d.Set("project_name", object["ProjectName"])
	d.Set("sls_region_id", object["SlsRegionId"])
	d.Set("status", object["Status"])
	return nil
}
func resourceAlicloudSmartagFlowLogUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	SagService := SagService{client}
	conn, err := client.NewSmartagClient()
	if err != nil {
		return WrapError(err)
	}
	var response map[string]interface{}
	d.Partial(true)

	update := false
	request := map[string]interface{}{
		"FlowLogId": d.Id(),
	}
	if !d.IsNewResource() && d.HasChange("description") {
		update = true
	}
	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}
	if !d.IsNewResource() && d.HasChange("flow_log_name") {
		update = true
	}
	if v, ok := d.GetOk("flow_log_name"); ok {
		request["Name"] = v
	}
	if !d.IsNewResource() && d.HasChange("inactive_aging") {
		update = true
	}
	if v, ok := d.GetOk("inactive_aging"); ok {
		request["InactiveAging"] = v
	}
	if !d.IsNewResource() && d.HasChange("logstore_name") {
		update = true
	}
	request["LogstoreName"] = d.Get("logstore_name")
	if !d.IsNewResource() && d.HasChange("netflow_server_ip") {
		update = true
	}
	if v, ok := d.GetOk("netflow_server_ip"); ok {
		request["NetflowServerIp"] = v
	}
	if !d.IsNewResource() && d.HasChange("netflow_server_port") {
		update = true
	}
	if v, ok := d.GetOk("netflow_server_port"); ok {
		request["NetflowServerPort"] = v
	}
	if !d.IsNewResource() && d.HasChange("netflow_version") {
		update = true
	}
	if v, ok := d.GetOk("netflow_version"); ok {
		request["NetflowVersion"] = v
	}
	if !d.IsNewResource() && d.HasChange("output_type") {
		update = true
	}
	if v, ok := d.GetOk("output_type"); ok {
		request["OutputType"] = v
	}
	if !d.IsNewResource() && d.HasChange("project_name") {
		update = true
	}
	request["ProjectName"] = d.Get("project_name")
	request["RegionId"] = client.RegionId
	if !d.IsNewResource() && d.HasChange("sls_region_id") {
		update = true
	}
	request["SlsRegionId"] = d.Get("sls_region_id")
	if !d.IsNewResource() && d.HasChange("active_aging") {
		update = true
		if v, ok := d.GetOk("active_aging"); ok {
			request["ActiveAging"] = v
		}
	}
	if update {
		action := "ModifyFlowLogAttribute"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2018-03-13"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
		d.SetPartial("description")
		d.SetPartial("flow_log_name")
		d.SetPartial("inactive_aging")
		d.SetPartial("logstore_name")
		d.SetPartial("netflow_server_ip")
		d.SetPartial("netflow_server_port")
		d.SetPartial("netflow_version")
		d.SetPartial("output_type")
		d.SetPartial("project_name")
		d.SetPartial("sls_region_id")
		d.SetPartial("active_aging")
	}
	if d.HasChange("status") {

		if d.Get("status") == "Active" {
			request := map[string]interface{}{
				"FlowLogId": d.Id(),
			}
			request["RegionId"] = client.RegionId
			action := "ActiveFlowLog"
			wait := incrementalWait(3*time.Second, 3*time.Second)
			err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
				response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2018-03-13"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
			stateConf := BuildStateConf([]string{}, []string{"Active"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, SagService.SmartagFlowLogStateRefreshFunc(d.Id(), []string{}))
			if _, err := stateConf.WaitForState(); err != nil {
				return WrapErrorf(err, IdMsg, d.Id())
			}
		}
		if d.Get("status") == "Inactive" {
			request := map[string]interface{}{
				"FlowLogId": d.Id(),
			}
			request["RegionId"] = client.RegionId
			action := "DeactiveFlowLog"
			wait := incrementalWait(3*time.Second, 3*time.Second)
			err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
				response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2018-03-13"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
			stateConf := BuildStateConf([]string{}, []string{"Inactive"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, SagService.SmartagFlowLogStateRefreshFunc(d.Id(), []string{}))
			if _, err := stateConf.WaitForState(); err != nil {
				return WrapErrorf(err, IdMsg, d.Id())
			}
		}
		d.SetPartial("status")
	}
	d.Partial(false)
	return resourceAlicloudSmartagFlowLogRead(d, meta)
}
func resourceAlicloudSmartagFlowLogDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DeleteFlowLog"
	var response map[string]interface{}
	conn, err := client.NewSmartagClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"FlowLogId": d.Id(),
	}

	request["RegionId"] = client.RegionId
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2018-03-13"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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

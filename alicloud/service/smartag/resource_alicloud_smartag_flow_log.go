package smartag

import (
	"fmt"
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/convert"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
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
	var err error
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

	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		response, err = client.RpcPost("Smartag", "2018-03-13", action, nil, request, false)
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
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_smartag_flow_log", action, errs.AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["FlowLogId"]))

	return resourceAlicloudSmartagFlowLogUpdate(d, meta)
}
func resourceAlicloudSmartagFlowLogRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	Service := Service{client}
	object, err := Service.DescribeSmartagFlowLog(d.Id())
	if err != nil {
		if errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_smartag_flow_log Service.DescribeSmartagFlowLog Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}
	if v, ok := object["ActiveAging"]; ok && fmt.Sprint(v) != "0" {
		active_agingVal, err := convert.FormatInt(v)
		if err != nil {
			return err
		}
		d.Set("active_aging", active_agingVal)
	}
	d.Set("description", object["Description"])
	d.Set("flow_log_name", object["Name"])
	if v, ok := object["InactiveAging"]; ok && fmt.Sprint(v) != "0" {
		inactive_agingVal, err := convert.FormatInt(v)
		if err != nil {
			return err
		}
		d.Set("inactive_aging", inactive_agingVal)
	}
	d.Set("logstore_name", object["LogstoreName"])
	d.Set("netflow_server_ip", object["NetflowServerIp"])
	if v, ok := object["NetflowServerPort"]; ok {
		netflow_server_portVal, err := convert.FormatInt(v)
		if err != nil {
			return err
		}
		d.Set("netflow_server_port", netflow_server_portVal)
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
	Service := Service{client}
	var err error
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
		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RpcPost("Smartag", "2018-03-13", action, nil, request, false)
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
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
		}
	}
	if d.HasChange("status") {

		if d.Get("status") == "Active" {
			request := map[string]interface{}{
				"FlowLogId": d.Id(),
			}
			request["RegionId"] = client.RegionId
			action := "ActiveFlowLog"
			wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
			err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
				response, err = client.RpcPost("Smartag", "2018-03-13", action, nil, request, false)
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
				return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
			}
			stateConf := helper.BuildStateConf([]string{}, []string{"Active"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, Service.SmartagFlowLogStateRefreshFunc(d.Id(), []string{}))
			if _, err := stateConf.WaitForState(); err != nil {
				return errs.WrapErrorf(err, errs.IdMsg, d.Id())
			}
		}
		if d.Get("status") == "Inactive" {
			request := map[string]interface{}{
				"FlowLogId": d.Id(),
			}
			request["RegionId"] = client.RegionId
			action := "DeactiveFlowLog"
			wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
			err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
				response, err = client.RpcPost("Smartag", "2018-03-13", action, nil, request, false)
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
				return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
			}
			stateConf := helper.BuildStateConf([]string{}, []string{"Inactive"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, Service.SmartagFlowLogStateRefreshFunc(d.Id(), []string{}))
			if _, err := stateConf.WaitForState(); err != nil {
				return errs.WrapErrorf(err, errs.IdMsg, d.Id())
			}
		}
	}
	d.Partial(false)
	return resourceAlicloudSmartagFlowLogRead(d, meta)
}
func resourceAlicloudSmartagFlowLogDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DeleteFlowLog"
	var response map[string]interface{}
	var err error
	request := map[string]interface{}{
		"FlowLogId": d.Id(),
	}

	request["RegionId"] = client.RegionId
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		response, err = client.RpcPost("Smartag", "2018-03-13", action, nil, request, false)
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
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
	}

	return nil
}

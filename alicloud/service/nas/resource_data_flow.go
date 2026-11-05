package nas

import (
	"fmt"
	"log"
	"regexp"
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

func resourceAlicloudNasDataFlow() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudNasDataFlowCreate,
		Read:   resourceAlicloudNasDataFlowRead,
		Update: resourceAlicloudNasDataFlowUpdate,
		Delete: resourceAlicloudNasDataFlowDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"data_flow_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.All(validation.StringLenBetween(2, 256), validation.StringDoesNotMatch(regexp.MustCompile(`(^http://.*)|(^https://.*)`), "It cannot begin with \"http://\", \"https://\".")),
			},
			"dry_run": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"file_system_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"fset_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"source_security_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"SSL", "NONE"}, false),
			},
			"source_storage": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"Running", "Stopped"}, false),
			},
			"throughput": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntInSlice([]int{600, 1200, 1500}),
			},
		},
	}
}

func resourceAlicloudNasDataFlowCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateDataFlow"
	request := make(map[string]interface{})
	var err error
	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}
	if v, ok := d.GetOkExists("dry_run"); ok {
		request["DryRun"] = v
	}
	request["FileSystemId"] = d.Get("file_system_id")
	request["FsetId"] = d.Get("fset_id")
	if v, ok := d.GetOk("source_security_type"); ok {
		request["SourceSecurityType"] = v
	}
	request["SourceStorage"] = d.Get("source_storage")
	request["Throughput"] = d.Get("throughput")
	request["ClientToken"] = helper.BuildClientToken("CreateDataFlow")
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		response, err = client.RpcPost("NAS", "2017-06-26", action, nil, request, true)
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
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_nas_data_flow", action, errs.AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(request["FileSystemId"], ":", response["DataFlowId"]))
	nasService := NasService{client}
	stateConf := helper.BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, nasService.NasDataFlowStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errs.WrapErrorf(err, errs.IdMsg, d.Id())
	}

	return resourceAlicloudNasDataFlowRead(d, meta)
}
func resourceAlicloudNasDataFlowRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	nasService := NasService{client}
	object, err := nasService.DescribeNasDataFlow(d.Id())
	if err != nil {
		if errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_nas_data_flow nasService.DescribeNasDataFlow Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}
	d.Set("data_flow_id", object["DataFlowId"])
	d.Set("file_system_id", object["FileSystemId"])
	d.Set("description", object["Description"])
	d.Set("fset_id", object["FsetId"])
	d.Set("source_security_type", object["SourceSecurityType"])
	d.Set("source_storage", object["SourceStorage"])
	d.Set("status", object["Status"])
	throughputVal, err := convert.FormatInt(object["Throughput"])
	if err != nil {
		return err
	}
	d.Set("throughput", throughputVal)
	return nil
}
func resourceAlicloudNasDataFlowUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	nasService := NasService{client}
	var err error
	var response map[string]interface{}
	parts, err := helper.ParseResourceId(d.Id(), 2)
	if err != nil {
		return errs.WrapError(err)
	}
	d.Partial(true)

	update := false
	modifyDataFlowReq := map[string]interface{}{
		"DataFlowId":   parts[1],
		"FileSystemId": parts[0],
	}
	if d.HasChange("description") {
		update = true
		if v, ok := d.GetOk("description"); ok {
			modifyDataFlowReq["Description"] = v
		}
	}
	if d.HasChange("throughput") {
		update = true
		modifyDataFlowReq["Throughput"] = d.Get("throughput")
	}
	if update {
		if v, ok := d.GetOkExists("dry_run"); ok {
			modifyDataFlowReq["DryRun"] = v
		}
		action := "ModifyDataFlow"
		modifyDataFlowReq["ClientToken"] = helper.BuildClientToken("ModifyDataFlow")
		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RpcPost("NAS", "2017-06-26", action, nil, modifyDataFlowReq, true)
			if err != nil {
				if helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		logs.AddDebug(action, response, modifyDataFlowReq)
		if err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
		}
		stateConf := helper.BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, nasService.NasDataFlowStateRefreshFunc(d.Id(), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return errs.WrapErrorf(err, errs.IdMsg, d.Id())
		}
	}

	if d.HasChange("status") {
		object, err := nasService.DescribeNasDataFlow(d.Id())
		if err != nil {
			return errs.WrapError(err)
		}
		target := d.Get("status").(string)
		if object["Status"].(string) != target {
			if target == "Running" {
				request := map[string]interface{}{
					"DataFlowId":   parts[1],
					"FileSystemId": parts[0],
				}
				if v, ok := d.GetOkExists("dry_run"); ok {
					request["DryRun"] = v
				}
				action := "StartDataFlow"
				request["ClientToken"] = helper.BuildClientToken("StartDataFlow")
				wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
				err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
					response, err = client.RpcPost("NAS", "2017-06-26", action, nil, request, true)
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
				stateConf := helper.BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, nasService.NasDataFlowStateRefreshFunc(d.Id(), []string{}))
				if _, err := stateConf.WaitForState(); err != nil {
					return errs.WrapErrorf(err, errs.IdMsg, d.Id())
				}
			}
			if target == "Stopped" {
				request := map[string]interface{}{
					"DataFlowId":   parts[1],
					"FileSystemId": parts[0],
				}
				if v, ok := d.GetOkExists("dry_run"); ok {
					request["DryRun"] = v
				}
				action := "StopDataFlow"
				request["ClientToken"] = helper.BuildClientToken("StopDataFlow")
				wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
				err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
					response, err = client.RpcPost("NAS", "2017-06-26", action, nil, request, true)
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
				stateConf := helper.BuildStateConf([]string{}, []string{"Stopped"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, nasService.NasDataFlowStateRefreshFunc(d.Id(), []string{}))
				if _, err := stateConf.WaitForState(); err != nil {
					return errs.WrapErrorf(err, errs.IdMsg, d.Id())
				}
			}
		}
	}
	d.Partial(false)
	return resourceAlicloudNasDataFlowRead(d, meta)
}
func resourceAlicloudNasDataFlowDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	parts, err := helper.ParseResourceId(d.Id(), 2)
	if err != nil {
		return errs.WrapError(err)
	}
	action := "DeleteDataFlow"
	var response map[string]interface{}
	request := map[string]interface{}{
		"DataFlowId":   parts[1],
		"FileSystemId": parts[0],
	}

	if v, ok := d.GetOkExists("dry_run"); ok {
		request["DryRun"] = v
	}
	request["ClientToken"] = helper.BuildClientToken("DeleteDataFlow")
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		response, err = client.RpcPost("NAS", "2017-06-26", action, nil, request, true)
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
	nasService := NasService{client}
	stateConf := helper.BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 5*time.Second, nasService.NasDataFlowStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errs.WrapErrorf(err, errs.IdMsg, d.Id())
	}
	return nil
}

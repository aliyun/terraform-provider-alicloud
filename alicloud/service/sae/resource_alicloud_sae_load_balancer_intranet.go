package sae

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAlicloudSaeLoadBalancerIntranet() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudSaeSaeLoadBalancerIntranetCreate,
		Read:   resourceAlicloudSaeSaeLoadBalancerIntranetRead,
		Update: resourceAlicloudSaeSaeLoadBalancerIntranetUpdate,
		Delete: resourceAlicloudSaeSaeLoadBalancerIntranetDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"app_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"intranet_slb_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"intranet": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"https_cert_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"protocol": {
							Type:         schema.TypeString,
							ValidateFunc: validation.StringInSlice([]string{"TCP", "HTTP", "HTTPS"}, false),
							Optional:     true,
						},
						"target_port": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"port": {
							Type:     schema.TypeInt,
							Optional: true,
						},
					},
				},
			},
			"intranet_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAlicloudSaeSaeLoadBalancerIntranetCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	saeService := NewService(client)
	var response map[string]interface{}
	action := "/pop/v1/sam/app/slb"
	request := make(map[string]*string)
	var err error
	request["AppId"] = helper.StringPointer(d.Get("app_id").(string))
	if v, ok := d.GetOk("intranet_slb_id"); ok {
		request["IntranetSlbId"] = helper.StringPointer(v.(string))
	}
	intranetReq := make([]interface{}, 0)
	for _, intranet := range d.Get("intranet").(*schema.Set).List() {
		intranetMap := intranet.(map[string]interface{})
		intranetReq = append(intranetReq, map[string]interface{}{
			"httpsCertId": intranetMap["https_cert_id"],
			"protocol":    intranetMap["protocol"],
			"targetPort":  intranetMap["target_port"],
			"port":        intranetMap["port"],
		})
	}
	obj, err := json.Marshal(intranetReq)
	if err != nil {
		return errs.WrapError(err)
	}
	request["Intranet"] = helper.StringPointer(string(obj))
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
		response, err = client.RoaPost("sae", "2019-05-06", action, request, nil, nil, false)
		if err != nil {
			if errs.IsExpectedErrors(err, []string{"Application.InvalidStatus", "Application.ChangerOrderRunning"}) || helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)
	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), "POST "+action, errs.AlibabaCloudSdkGoERROR)
	}
	d.SetId(fmt.Sprint(d.Get("app_id")))

	stateConf := helper.BuildStateConf([]string{}, []string{"SUCCESS"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, saeService.SaeApplicationStateRefreshFunc(d.Get("app_id").(string), []string{"FAIL", "AUTO_BATCH_WAIT", "APPROVED", "WAIT_APPROVAL", "WAIT_BATCH_CONFIRM", "ABORT", "SYSTEM_FAIL"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errs.WrapErrorf(err, errs.IdMsg, d.Id())
	}
	return resourceAlicloudSaeSaeLoadBalancerIntranetRead(d, meta)
}
func resourceAlicloudSaeSaeLoadBalancerIntranetRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	saeService := NewService(client)

	describeApplicationSlbObject, err := saeService.DescribeApplicationSlb(d.Id())
	if err != nil {
		return errs.WrapError(err)
	}

	d.Set("intranet_ip", describeApplicationSlbObject["IntranetIp"])
	d.Set("intranet_slb_id", describeApplicationSlbObject["IntranetSlbId"])
	d.Set("app_id", d.Id())
	intranetArray := make([]interface{}, 0)
	if v, ok := describeApplicationSlbObject["Intranet"]; ok {
		for _, intranet := range v.([]interface{}) {
			intranetObject := intranet.(map[string]interface{})
			intranetObj := map[string]interface{}{
				"https_cert_id": intranetObject["HttpsCertId"],
				"protocol":      intranetObject["Protocol"],
				"target_port":   intranetObject["TargetPort"],
				"port":          intranetObject["Port"],
			}
			intranetArray = append(intranetArray, intranetObj)
		}
	}
	d.Set("intranet", intranetArray)
	return nil
}
func resourceAlicloudSaeSaeLoadBalancerIntranetUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	saeService := NewService(client)
	var err error
	var response map[string]interface{}
	update := false
	request := map[string]*string{
		"AppId": helper.StringPointer(d.Id()),
	}
	if d.HasChange("intranet_slb_id") {
		update = true
	}
	if v, ok := d.GetOk("intranet_slb_id"); ok {
		request["IntranetSlbId"] = helper.StringPointer(v.(string))
	}

	if d.HasChange("intranet") {
		update = true
	}
	intranetReq := make([]interface{}, 0)
	for _, intranet := range d.Get("intranet").(*schema.Set).List() {
		intranetMap := intranet.(map[string]interface{})
		intranetReq = append(intranetReq, map[string]interface{}{
			"httpsCertId": intranetMap["https_cert_id"],
			"protocol":    intranetMap["protocol"],
			"targetPort":  intranetMap["target_port"],
			"port":        intranetMap["port"],
		})
	}
	obj, err := json.Marshal(intranetReq)
	if err != nil {
		return errs.WrapError(err)
	}
	request["Intranet"] = helper.StringPointer(string(obj))

	if update {
		action := "/pop/v1/sam/app/slb"
		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RoaPost("sae", "2019-05-06", action, request, nil, nil, false)
			if err != nil {
				if errs.IsExpectedErrors(err, []string{"Application.InvalidStatus", "Application.ChangerOrderRunning"}) || helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		logs.AddDebug(action, response, request)
		if err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), "POST "+action, errs.AlibabaCloudSdkGoERROR)
		}
	}
	stateConf := helper.BuildStateConf([]string{}, []string{"SUCCESS"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, saeService.SaeApplicationStateRefreshFunc(d.Get("app_id").(string), []string{"FAIL", "AUTO_BATCH_WAIT", "APPROVED", "WAIT_APPROVAL", "WAIT_BATCH_CONFIRM", "ABORT", "SYSTEM_FAIL"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errs.WrapErrorf(err, errs.IdMsg, d.Id())
	}
	return resourceAlicloudSaeSaeLoadBalancerIntranetRead(d, meta)
}
func resourceAlicloudSaeSaeLoadBalancerIntranetDelete(d *schema.ResourceData, meta interface{}) error {
	request := map[string]*string{
		"AppId":    helper.StringPointer(d.Id()),
		"Intranet": helper.StringPointer(strconv.FormatBool(true)),
	}
	client := meta.(*connectivity.AliyunClient)
	var err error

	action := "/pop/v1/sam/app/slb"
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	var response map[string]interface{}
	err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
		response, err = client.RoaDelete("sae", "2019-05-06", action, request, nil, nil, false)
		if err != nil {
			if errs.IsExpectedErrors(err, []string{"Application.InvalidStatus", "Application.ChangerOrderRunning"}) || helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)
	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), "POST "+action, errs.AlibabaCloudSdkGoERROR)
	}
	return nil
}

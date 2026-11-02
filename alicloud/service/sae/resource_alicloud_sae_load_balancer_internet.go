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

func resourceAlicloudSaeLoadBalancerInternet() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudSaeSaeLoadBalancerInternetCreate,
		Read:   resourceAlicloudSaeSaeLoadBalancerInternetRead,
		Update: resourceAlicloudSaeSaeLoadBalancerInternetUpdate,
		Delete: resourceAlicloudSaeSaeLoadBalancerInternetDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"app_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"internet_slb_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"internet": {
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
			"internet_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAlicloudSaeSaeLoadBalancerInternetCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	saeService := NewService(client)
	var response map[string]interface{}
	action := "/pop/v1/sam/app/slb"
	request := make(map[string]*string)
	var err error
	request["AppId"] = helper.StringPointer(d.Get("app_id").(string))
	if v, ok := d.GetOk("internet_slb_id"); ok {
		request["InternetSlbId"] = helper.StringPointer(v.(string))
	}
	internetReq := make([]interface{}, 0)
	for _, internet := range d.Get("internet").(*schema.Set).List() {
		internetMap := internet.(map[string]interface{})
		internetReq = append(internetReq, map[string]interface{}{
			"httpsCertId": internetMap["https_cert_id"],
			"protocol":    internetMap["protocol"],
			"targetPort":  internetMap["target_port"],
			"port":        internetMap["port"],
		})
	}
	obj, err := json.Marshal(internetReq)
	if err != nil {
		return errs.WrapError(err)
	}
	request["Internet"] = helper.StringPointer(string(obj))
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
	return resourceAlicloudSaeSaeLoadBalancerInternetRead(d, meta)
}
func resourceAlicloudSaeSaeLoadBalancerInternetRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	saeService := NewService(client)

	describeApplicationSlbObject, err := saeService.DescribeApplicationSlb(d.Id())
	if err != nil {
		return errs.WrapError(err)
	}

	d.Set("internet_ip", describeApplicationSlbObject["InternetIp"])
	d.Set("internet_slb_id", describeApplicationSlbObject["InternetSlbId"])
	d.Set("app_id", d.Id())
	internetArray := make([]interface{}, 0)
	if v, ok := describeApplicationSlbObject["Internet"]; ok {
		for _, internet := range v.([]interface{}) {
			internetObject := internet.(map[string]interface{})
			internetObj := map[string]interface{}{
				"https_cert_id": internetObject["HttpsCertId"],
				"protocol":      internetObject["Protocol"],
				"target_port":   internetObject["TargetPort"],
				"port":          internetObject["Port"],
			}
			internetArray = append(internetArray, internetObj)
		}
	}
	d.Set("internet", internetArray)
	return nil
}
func resourceAlicloudSaeSaeLoadBalancerInternetUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	saeService := NewService(client)
	var err error
	var response map[string]interface{}
	update := false
	request := map[string]*string{
		"AppId": helper.StringPointer(d.Id()),
	}
	if d.HasChange("internet_slb_id") {
		update = true
	}
	if v, ok := d.GetOk("internet_slb_id"); ok {
		request["InternetSlbId"] = helper.StringPointer(v.(string))
	}

	if d.HasChange("internet") {
		update = true
	}
	internetReq := make([]interface{}, 0)
	for _, internet := range d.Get("internet").(*schema.Set).List() {
		internetMap := internet.(map[string]interface{})
		internetReq = append(internetReq, map[string]interface{}{
			"httpsCertId": internetMap["https_cert_id"],
			"protocol":    internetMap["protocol"],
			"targetPort":  internetMap["target_port"],
			"port":        internetMap["port"],
		})
	}
	obj, err := json.Marshal(internetReq)
	if err != nil {
		return errs.WrapError(err)
	}
	request["Internet"] = helper.StringPointer(string(obj))

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
		if err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), "POST "+action, errs.AlibabaCloudSdkGoERROR)
		}
		logs.AddDebug(action, response, request)
	}
	stateConf := helper.BuildStateConf([]string{}, []string{"SUCCESS"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, saeService.SaeApplicationStateRefreshFunc(d.Get("app_id").(string), []string{"FAIL", "AUTO_BATCH_WAIT", "APPROVED", "WAIT_APPROVAL", "WAIT_BATCH_CONFIRM", "ABORT", "SYSTEM_FAIL"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errs.WrapErrorf(err, errs.IdMsg, d.Id())
	}
	return resourceAlicloudSaeSaeLoadBalancerInternetRead(d, meta)
}
func resourceAlicloudSaeSaeLoadBalancerInternetDelete(d *schema.ResourceData, meta interface{}) error {
	request := map[string]*string{
		"AppId":    helper.StringPointer(d.Id()),
		"Internet": helper.StringPointer(strconv.FormatBool(true)),
	}
	client := meta.(*connectivity.AliyunClient)
	var err error

	action := "/pop/v1/sam/app/slb"
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	var response map[string]interface{}
	err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
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

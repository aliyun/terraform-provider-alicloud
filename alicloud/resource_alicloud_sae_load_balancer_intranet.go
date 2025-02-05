package alicloud

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
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
	saeService := SaeService{client}
	var response map[string]interface{}
	action := "/pop/v1/sam/app/slb"
	request := make(map[string]*string)
	var err error
	request["AppId"] = StringPointer(d.Get("app_id").(string))
	if v, ok := d.GetOk("intranet_slb_id"); ok {
		request["IntranetSlbId"] = StringPointer(v.(string))
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
		return WrapError(err)
	}
	request["Intranet"] = StringPointer(string(obj))
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
		response, err = client.RoaPost("sae", "2019-05-06", action, request, nil, nil, false)
		if err != nil {
			if IsExpectedErrors(err, []string{"Application.InvalidStatus", "Application.ChangerOrderRunning"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), "POST "+action, AlibabaCloudSdkGoERROR)
	}
	d.SetId(fmt.Sprint(d.Get("app_id")))

	stateConf := BuildStateConf([]string{}, []string{"SUCCESS"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, saeService.SaeApplicationStateRefreshFunc(d.Get("app_id").(string), []string{"FAIL", "AUTO_BATCH_WAIT", "APPROVED", "WAIT_APPROVAL", "WAIT_BATCH_CONFIRM", "ABORT", "SYSTEM_FAIL"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return resourceAlicloudSaeSaeLoadBalancerIntranetRead(d, meta)
}
func resourceAlicloudSaeSaeLoadBalancerIntranetRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	saeService := SaeService{client}

	describeApplicationSlbObject, err := saeService.DescribeApplicationSlb(d.Id())
	if err != nil {
		return WrapError(err)
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
	saeService := SaeService{client}
	var err error
	var response map[string]interface{}
	update := false
	request := map[string]*string{
		"AppId": StringPointer(d.Id()),
	}
	if d.HasChange("intranet_slb_id") {
		update = true
	}
	if v, ok := d.GetOk("intranet_slb_id"); ok {
		request["IntranetSlbId"] = StringPointer(v.(string))
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
		return WrapError(err)
	}
	request["Intranet"] = StringPointer(string(obj))

	if update {
		action := "/pop/v1/sam/app/slb"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RoaPost("sae", "2019-05-06", action, request, nil, nil, false)
			if err != nil {
				if IsExpectedErrors(err, []string{"Application.InvalidStatus", "Application.ChangerOrderRunning"}) || NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, request)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), "POST "+action, AlibabaCloudSdkGoERROR)
		}
	}
	stateConf := BuildStateConf([]string{}, []string{"SUCCESS"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, saeService.SaeApplicationStateRefreshFunc(d.Get("app_id").(string), []string{"FAIL", "AUTO_BATCH_WAIT", "APPROVED", "WAIT_APPROVAL", "WAIT_BATCH_CONFIRM", "ABORT", "SYSTEM_FAIL"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return resourceAlicloudSaeSaeLoadBalancerIntranetRead(d, meta)
}
func resourceAlicloudSaeSaeLoadBalancerIntranetDelete(d *schema.ResourceData, meta interface{}) error {
	request := map[string]*string{
		"AppId":    StringPointer(d.Id()),
		"Intranet": StringPointer(strconv.FormatBool(true)),
	}
	client := meta.(*connectivity.AliyunClient)
	var err error

	action := "/pop/v1/sam/app/slb"
	wait := incrementalWait(3*time.Second, 3*time.Second)
	var response map[string]interface{}
	err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
		response, err = client.RoaDelete("sae", "2019-05-06", action, request, nil, nil, false)
		if err != nil {
			if IsExpectedErrors(err, []string{"Application.InvalidStatus", "Application.ChangerOrderRunning"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), "POST "+action, AlibabaCloudSdkGoERROR)
	}
	return nil
}

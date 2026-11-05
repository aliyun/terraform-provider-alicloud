package alidns

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	alidnssdk "github.com/aliyun/alibaba-cloud-sdk-go/services/alidns"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/validators"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAliCloudAlidnsRecord() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudAlidnsRecordCreate,
		Read:   resourceAliCloudAlidnsRecordRead,
		Update: resourceAliCloudAlidnsRecordUpdate,
		Delete: resourceAliCloudAlidnsRecordDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(4 * time.Minute),
			Update: schema.DefaultTimeout(3 * time.Minute),
			Delete: schema.DefaultTimeout(6 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"domain_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"lang": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"line": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "default",
			},
			"priority": {
				Type:             schema.TypeInt,
				Optional:         true,
				DiffSuppressFunc: dnsPriorityDiffSuppressFunc,
			},
			"rr": {
				Type:     schema.TypeString,
				Required: true,
			},
			"remark": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validators.StringInSlice([]string{"ENABLE", "DISABLE"}, false),
				Default:      "ENABLE",
			},
			"ttl": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  600,
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"user_client_ip": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"value": {
				Type:             schema.TypeString,
				Required:         true,
				DiffSuppressFunc: dnsValueDiffSuppressFunc,
			},
		},
	}
}

func resourceAliCloudAlidnsRecordCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := alidnssdk.CreateAddDomainRecordRequest()
	request.DomainName = d.Get("domain_name").(string)
	if v, ok := d.GetOk("lang"); ok {
		request.Lang = v.(string)
	}
	if v, ok := d.GetOk("line"); ok {
		line := v.(string)
		if line != "default" && d.Get("type").(string) == "FORWORD_URL" {
			return errs.WrapError(errs.Error("The ForwordURLRecord only support default line."))
		}
		request.Line = line
	}
	if v, ok := d.GetOk("priority"); !ok && d.Get("type").(string) == "MX" {
		return errs.WrapError(errs.Error("'priority': required field when 'type' is MX."))
	} else if ok {
		request.Priority = requests.Integer(strconv.Itoa(v.(int)))
	}
	request.RR = d.Get("rr").(string)
	if v, ok := d.GetOk("ttl"); ok {
		request.TTL = requests.NewInteger(v.(int))
	}
	request.Type = d.Get("type").(string)
	if v, ok := d.GetOk("user_client_ip"); ok {
		request.UserClientIp = v.(string)
	}
	request.Value = d.Get("value").(string)
	wait := helper.IncrementalWait(3*time.Second, 10*time.Second)
	err := retry.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *retry.RetryError {
		raw, err := client.WithAlidnsClient(func(alidnsClient *alidnssdk.Client) (interface{}, error) {
			return alidnsClient.AddDomainRecord(request)
		})
		if err != nil {
			if errs.IsExpectedErrors(err, []string{"InternalError", "LastOperationNotFinished"}) || helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		logs.AddDebug(request.GetActionName(), raw)
		response, _ := raw.(*alidnssdk.AddDomainRecordResponse)
		d.SetId(fmt.Sprintf("%v", response.RecordId))
		return nil
	})
	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_alidns_record", request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
	}

	return resourceAliCloudAlidnsRecordUpdate(d, meta)
}

func resourceAliCloudAlidnsRecordRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	alidnsService := NewAlidnsService(client)
	object, err := alidnsService.DescribeAlidnsRecord(d.Id())
	if err != nil {
		if errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_alidns_record alidnsService.DescribeAlidnsRecord Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}

	d.Set("domain_name", object.DomainName)
	d.Set("line", object.Line)
	d.Set("priority", object.Priority)
	d.Set("rr", object.RR)
	d.Set("status", object.Status)
	d.Set("ttl", object.TTL)
	d.Set("type", object.Type)
	d.Set("value", object.Value)

	return nil
}

func resourceAliCloudAlidnsRecordUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	d.Partial(true)

	update := false
	request := alidnssdk.CreateSetDomainRecordStatusRequest()
	request.RecordId = d.Id()
	if d.HasChange("status") {
		update = true
	}
	request.Status = d.Get("status").(string)
	if !d.IsNewResource() && d.HasChange("lang") {
		update = true
	}
	request.Lang = d.Get("lang").(string)
	if !d.IsNewResource() && d.HasChange("user_client_ip") {
		update = true
	}
	request.UserClientIp = d.Get("user_client_ip").(string)
	if update {
		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err := retry.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *retry.RetryError {
			raw, err := client.WithAlidnsClient(func(alidnsClient *alidnssdk.Client) (interface{}, error) {
				return alidnsClient.SetDomainRecordStatus(request)
			})
			if err != nil {
				if errs.IsExpectedErrors(err, []string{"LastOperationNotFinished"}) || helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			logs.AddDebug(request.GetActionName(), raw)
			return nil
		})
		if err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
		}
	}
	update = false
	updateDomainRecordRemarkReq := alidnssdk.CreateUpdateDomainRecordRemarkRequest()
	updateDomainRecordRemarkReq.RecordId = d.Id()
	updateDomainRecordRemarkReq.Lang = d.Get("lang").(string)
	if d.HasChange("remark") {
		update = true
	}
	updateDomainRecordRemarkReq.Remark = d.Get("remark").(string)
	updateDomainRecordRemarkReq.UserClientIp = d.Get("user_client_ip").(string)
	if update {
		raw, err := client.WithAlidnsClient(func(alidnsClient *alidnssdk.Client) (interface{}, error) {
			return alidnsClient.UpdateDomainRecordRemark(updateDomainRecordRemarkReq)
		})
		logs.AddDebug(updateDomainRecordRemarkReq.GetActionName(), raw)
		if err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), updateDomainRecordRemarkReq.GetActionName(), errs.AlibabaCloudSdkGoERROR)
		}
	}
	update = false
	updateDomainRecordReq := alidnssdk.CreateUpdateDomainRecordRequest()
	updateDomainRecordReq.RecordId = d.Id()
	if !d.IsNewResource() && d.HasChange("rr") {
		update = true
	}
	updateDomainRecordReq.RR = d.Get("rr").(string)
	if !d.IsNewResource() && d.HasChange("type") {
		update = true
	}
	updateDomainRecordReq.Type = d.Get("type").(string)
	if !d.IsNewResource() && d.HasChange("value") {
		update = true
	}
	updateDomainRecordReq.Value = d.Get("value").(string)
	updateDomainRecordReq.Lang = d.Get("lang").(string)
	if !d.IsNewResource() && d.HasChange("line") {
		update = true
	}
	updateDomainRecordReq.Line = d.Get("line").(string)
	if updateDomainRecordReq.Type == "MX" {
		if !d.IsNewResource() && d.HasChange("priority") {
			update = true
		}
		updateDomainRecordReq.Priority = requests.NewInteger(d.Get("priority").(int))
	}

	if !d.IsNewResource() && d.HasChange("ttl") {
		update = true
	}
	updateDomainRecordReq.TTL = requests.NewInteger(d.Get("ttl").(int))
	updateDomainRecordReq.UserClientIp = d.Get("user_client_ip").(string)
	if update {
		raw, err := client.WithAlidnsClient(func(alidnsClient *alidnssdk.Client) (interface{}, error) {
			return alidnsClient.UpdateDomainRecord(updateDomainRecordReq)
		})
		logs.AddDebug(updateDomainRecordReq.GetActionName(), raw)
		if err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), updateDomainRecordReq.GetActionName(), errs.AlibabaCloudSdkGoERROR)
		}
	}
	d.Partial(false)

	return resourceAliCloudAlidnsRecordRead(d, meta)
}

func resourceAliCloudAlidnsRecordDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	request := alidnssdk.CreateDeleteDomainRecordRequest()
	request.RecordId = d.Id()
	if v, ok := d.GetOk("lang"); ok {
		request.Lang = v.(string)
	}
	if v, ok := d.GetOk("user_client_ip"); ok {
		request.UserClientIp = v.(string)
	}
	wait := helper.IncrementalWait(3*time.Second, 10*time.Second)
	err := retry.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *retry.RetryError {
		raw, err := client.WithAlidnsClient(func(alidnsClient *alidnssdk.Client) (interface{}, error) {
			return alidnsClient.DeleteDomainRecord(request)
		})
		if err != nil {
			if errs.IsExpectedErrors(err, []string{"InternalError", "RecordForbidden.DNSChange", "LastOperationNotFinished"}) || helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		logs.AddDebug(request.GetActionName(), raw)
		return nil
	})
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"DomainRecordNotBelongToUser"}) {
			return nil
		}
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
	}

	return nil
}

package alidns

import (
	"strconv"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	alidnssdk "github.com/aliyun/alibaba-cloud-sdk-go/services/alidns"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAlicloudDnsRecord() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudDnsRecordCreate,
		Read:   resourceAlicloudDnsRecordRead,
		Update: resourceAlicloudDnsRecordUpdate,
		Delete: resourceAlicloudDnsRecordDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		DeprecationMessage: "This resource has been deprecated since v1.85.0 and will be removed in the future. Please use 'alicloud_alidns_record' instead.",

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"host_record": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateRR,
			},
			"type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"A", "NS", "MX", "TXT", "CNAME", "SRV", "AAAA", "CAA", "REDIRECT_URL", "FORWORD_URL"}, false),
			},
			"value": {
				Type:             schema.TypeString,
				Required:         true,
				DiffSuppressFunc: dnsValueDiffSuppressFunc,
			},
			"ttl": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  600,
			},
			"priority": {
				Type:             schema.TypeInt,
				Optional:         true,
				DiffSuppressFunc: dnsPriorityDiffSuppressFunc,
			},
			"routing": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "default",
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"locked": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func resourceAlicloudDnsRecordCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	request := alidnssdk.CreateAddDomainRecordRequest()
	request.RegionId = client.RegionId
	request.DomainName = d.Get("name").(string)
	request.RR = d.Get("host_record").(string)
	request.Type = d.Get("type").(string)
	request.Value = d.Get("value").(string)
	request.TTL = requests.NewInteger(d.Get("ttl").(int))

	if v, ok := d.GetOk("priority"); !ok && request.Type == "MX" {
		return errs.WrapError(errs.Error("'priority': required field when 'type' is MX."))
	} else if ok {
		request.Priority = requests.Integer(strconv.Itoa(v.(int)))
	}

	if v, ok := d.GetOk("routing"); ok {
		routing := v.(string)
		if routing != "default" && request.Type == "FORWORD_URL" {
			return errs.WrapError(errs.Error("The ForwordURLRecord only support default line."))
		}
		request.Line = routing
	}

	if err := retry.Retry(3*time.Minute, func() *retry.RetryError {
		raw, err := client.WithDnsClient(func(dnsClient *alidnssdk.Client) (interface{}, error) {
			return dnsClient.AddDomainRecord(request)
		})
		if err != nil {
			if errs.IsExpectedErrors(err, []string{"InternalError", "LastOperationNotFinished"}) {
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)
		response, _ := raw.(*alidnssdk.AddDomainRecordResponse)
		d.SetId(response.RecordId)
		return nil
	}); err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_dns_record", request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
	}

	return resourceAlicloudDnsRecordRead(d, meta)
}

func resourceAlicloudDnsRecordUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := alidnssdk.CreateUpdateDomainRecordRequest()
	request.RegionId = client.RegionId
	request.RecordId = d.Id()
	request.RR = d.Get("host_record").(string)
	request.Type = d.Get("type").(string)
	if request.Type == "MX" {
		request.Priority = requests.NewInteger(d.Get("priority").(int))
	}
	request.TTL = requests.NewInteger(d.Get("ttl").(int))
	request.Line = d.Get("routing").(string)

	request.Value = d.Get("value").(string)

	if err := retry.Retry(3*time.Minute, func() *retry.RetryError {
		raw, err := client.WithDnsClient(func(dnsClient *alidnssdk.Client) (interface{}, error) {
			return dnsClient.UpdateDomainRecord(request)
		})
		if err != nil {
			if errs.IsExpectedErrors(err, []string{"LastOperationNotFinished"}) {
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	}); err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
	}

	return resourceAlicloudDnsRecordRead(d, meta)
}

func resourceAlicloudDnsRecordRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	dnsService := NewDnsService(client)
	object, err := dnsService.DescribeDnsRecord(d.Id())
	if err != nil {
		if errs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}
	d.Set("ttl", object.TTL)
	d.Set("priority", object.Priority)
	d.Set("name", object.DomainName)
	d.Set("host_record", object.RR)
	d.Set("type", object.Type)
	d.Set("value", object.Value)
	d.Set("routing", object.Line)
	d.Set("status", object.Status)
	d.Set("locked", object.Locked)

	return nil
}

func resourceAlicloudDnsRecordDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	dnsService := NewDnsService(client)
	request := alidnssdk.CreateDeleteDomainRecordRequest()
	request.RegionId = client.RegionId
	request.RecordId = d.Id()
	return retry.Retry(5*time.Minute, func() *retry.RetryError {
		raw, err := client.WithDnsClient(func(dnsClient *alidnssdk.Client) (interface{}, error) {
			return dnsClient.DeleteDomainRecord(request)
		})
		if err != nil {
			if errs.IsExpectedErrors(err, []string{"DomainRecordNotBelongToUser"}) {
				return nil
			}
			if errs.IsExpectedErrors(err, []string{"RecordForbidden.DNSChange", "InternalError", "LastOperationNotFinished"}) {
				return retry.RetryableError(errs.WrapErrorf(err, errs.DefaultTimeoutMsg, d.Id(), request.GetActionName(), errs.AlibabaCloudSdkGoERROR))
			}
			return retry.NonRetryableError(errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), request.GetActionName(), errs.AlibabaCloudSdkGoERROR))
		}
		logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)
		_, err = dnsService.DescribeDnsRecord(d.Id())
		if err != nil {
			if errs.NotFoundError(err) {
				return nil
			}
			return retry.NonRetryableError(errs.WrapError(err))
		}
		return retry.RetryableError(errs.WrapErrorf(err, errs.DefaultTimeoutMsg, d.Id(), request.GetActionName(), errs.AlibabaCloudSdkGoERROR))
	})
}

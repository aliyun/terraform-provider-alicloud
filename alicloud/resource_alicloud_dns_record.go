package alicloud

import (
	"strconv"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/alidns"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
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
				ValidateFunc: validateDomainRecordType,
			},
			"value": {
				Type:     schema.TypeString,
				Required: true,
			},
			"ttl": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  600,
			},
			"priority": {
				Type:             schema.TypeInt,
				Optional:         true,
				ValidateFunc:     validateDomainRecordPriority,
				DiffSuppressFunc: dnsPriorityDiffSuppressFunc,
			},
			"routing": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateDomainRecordLine,
				Default:      "default",
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

	request := alidns.CreateAddDomainRecordRequest()
	request.DomainName = d.Get("name").(string)
	request.RR = d.Get("host_record").(string)
	request.Type = d.Get("type").(string)
	request.Value = d.Get("value").(string)

	if v, ok := d.GetOk("priority"); !ok && request.Type == "MX" {
		return WrapError(Error("'priority': required field when 'type' is MX."))
	} else if ok {
		request.Priority = requests.Integer(strconv.Itoa(v.(int)))
	}

	if v, ok := d.GetOk("routing"); ok {
		routing := v.(string)
		if routing != "default" && request.Type == "FORWORD_URL" {
			return WrapError(Error("The ForwordURLRecord only support default line."))
		}
		request.Line = routing
	}

	if err := resource.Retry(3*time.Minute, func() *resource.RetryError {
		raw, err := client.WithDnsClient(func(dnsClient *alidns.Client) (interface{}, error) {
			return dnsClient.AddDomainRecord(request)
		})
		if err != nil {
			if IsExceptedError(err, DnsInternalError) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw)
		response, _ := raw.(*alidns.AddDomainRecordResponse)
		d.SetId(response.RecordId)
		return nil
	}); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "dns_record", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	return resourceAlicloudDnsRecordUpdate(d, meta)
}

func resourceAlicloudDnsRecordUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	d.Partial(true)
	attributeUpdate := false
	request := alidns.CreateUpdateDomainRecordRequest()
	request.RegionId = d.Id()
	request.RR = d.Get("host_record").(string)
	request.Type = d.Get("type").(string)
	request.Value = d.Get("value").(string)

	if !d.IsNewResource() {
		requiredParams := []string{"host_record", "type", "value"}
		for _, v := range requiredParams {
			if d.HasChange(v) {
				d.SetPartial(v)
				attributeUpdate = true
			}
		}
	}
	if d.HasChange("priority") && !d.IsNewResource() {
		d.SetPartial("priority")
		request.Priority = requests.Integer(strconv.Itoa(d.Get("priority").(int)))
		attributeUpdate = true
	}

	if d.HasChange("ttl") && !d.IsNewResource() {
		d.SetPartial("ttl")
		request.TTL = requests.Integer(strconv.Itoa(d.Get("ttl").(int)))
		attributeUpdate = true
	}

	if d.HasChange("routing") && !d.IsNewResource() {
		d.SetPartial("routing")
		request.Line = d.Get("routing").(string)
		attributeUpdate = true
	}

	if attributeUpdate {
		raw, err := client.WithDnsClient(func(dnsClient *alidns.Client) (interface{}, error) {
			return dnsClient.UpdateDomainRecord(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw)
	}

	d.Partial(false)

	return resourceAlicloudDnsRecordRead(d, meta)
}

func resourceAlicloudDnsRecordRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	dnsService := &DnsService{client: client}
	recordInfo, err := dnsService.DescribeDnsRecord(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("ttl", recordInfo.TTL)
	d.Set("priority", recordInfo.Priority)
	d.Set("name", recordInfo.DomainName)
	d.Set("host_record", recordInfo.RR)
	d.Set("type", recordInfo.Type)
	d.Set("value", recordInfo.Value)
	d.Set("routing", recordInfo.Line)
	d.Set("status", recordInfo.Status)
	d.Set("locked", recordInfo.Locked)

	return nil
}

func resourceAlicloudDnsRecordDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	request := alidns.CreateDeleteDomainRecordRequest()
	request.RecordId = d.Id()
	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err := client.WithDnsClient(func(dnsClient *alidns.Client) (interface{}, error) {
			return dnsClient.DeleteDomainRecord(request)
		})
		if err != nil {
			if IsExceptedErrors(err, []string{DomainRecordNotBelongToUser}) {
				return nil
			}
			if IsExceptedErrors(err, []string{RecordForbiddenDNSChange}) {
				return resource.RetryableError(WrapErrorf(err, DeleteTimeoutMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR))
			}
			return resource.NonRetryableError(WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR))
		}

		addDebug(request.GetActionName(), raw)
		dnsService := &DnsService{client: client}
		_, err = dnsService.DescribeDnsRecord(d.Id())
		if err != nil {
			if NotFoundError(err) || IsExceptedError(err, DomainRecordNotBelongToUser) {
				return nil
			}
			return resource.NonRetryableError(WrapError(err))
		}

		return nil
	})
}

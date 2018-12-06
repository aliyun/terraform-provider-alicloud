package alicloud

import (
	"fmt"
	"time"

	"github.com/denverdino/aliyungo/dns"
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
			"name": &schema.Schema{
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

	args := &dns.AddDomainRecordArgs{
		DomainName: d.Get("name").(string),
		RR:         d.Get("host_record").(string),
		Type:       d.Get("type").(string),
		Value:      d.Get("value").(string),
		Priority:   int32(d.Get("priority").(int)),
	}

	if _, ok := d.GetOk("priority"); !ok && args.Type == dns.MXRecord {
		return fmt.Errorf("'priority': required field when 'type' is MX.")
	}

	if v, ok := d.GetOk("routing"); ok {
		routing := v.(string)
		if routing != "default" && args.Type == dns.ForwordURLRecord {
			return fmt.Errorf("The ForwordURLRecord only support default line.")
		}
		args.Line = routing
	}

	if err := resource.Retry(3*time.Minute, func() *resource.RetryError {
		raw, err := client.WithDnsClient(func(dnsClient *dns.Client) (interface{}, error) {
			return dnsClient.AddDomainRecord(args)
		})
		if err != nil {
			if IsExceptedError(err, DnsInternalError) {
				return resource.RetryableError(fmt.Errorf("create resource failure for lock conflict:%v", err))
			}
			return resource.NonRetryableError(err)
		}
		response, _ := raw.(*dns.AddDomainRecordResponse)
		d.SetId(response.RecordId)
		return nil
	}); err != nil {
		return fmt.Errorf("AddDomainRecord got a error: %#v", err)
	}

	return resourceAlicloudDnsRecordUpdate(d, meta)
}

func resourceAlicloudDnsRecordUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	d.Partial(true)
	attributeUpdate := false
	args := &dns.UpdateDomainRecordArgs{
		RecordId: d.Id(),
		RR:       d.Get("host_record").(string),
		Type:     d.Get("type").(string),
		Value:    d.Get("value").(string),
	}

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
		args.Priority = int32(d.Get("priority").(int))
		attributeUpdate = true
	}

	if d.HasChange("ttl") && !d.IsNewResource() {
		d.SetPartial("ttl")
		args.TTL = int32(d.Get("ttl").(int))
		attributeUpdate = true
	}

	if d.HasChange("routing") && !d.IsNewResource() {
		d.SetPartial("routing")
		args.Line = d.Get("routing").(string)
		attributeUpdate = true
	}

	if attributeUpdate {
		_, err := client.WithDnsClient(func(dnsClient *dns.Client) (interface{}, error) {
			return dnsClient.UpdateDomainRecord(args)
		})
		if err != nil {
			return fmt.Errorf("UpdateDomainRecord got an error: %#v", err)
		}
	}

	d.Partial(false)

	return resourceAlicloudDnsRecordRead(d, meta)
}

func resourceAlicloudDnsRecordRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	args := &dns.DescribeDomainRecordInfoNewArgs{
		RecordId: d.Id(),
	}
	raw, err := client.WithDnsClient(func(dnsClient *dns.Client) (interface{}, error) {
		return dnsClient.DescribeDomainRecordInfoNew(args)
	})
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return err
	}
	response, _ := raw.(*dns.DescribeDomainRecordInfoNewResponse)
	record := response.RecordTypeNew
	d.Set("ttl", record.TTL)
	d.Set("priority", record.Priority)
	d.Set("name", record.DomainName)
	d.Set("host_record", record.RR)
	d.Set("type", record.Type)
	d.Set("value", record.Value)
	d.Set("routing", record.Line)
	d.Set("status", record.Status)
	d.Set("locked", record.Locked)

	return nil
}

func resourceAlicloudDnsRecordDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	args := &dns.DeleteDomainRecordArgs{
		RecordId: d.Id(),
	}
	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		_, err := client.WithDnsClient(func(dnsClient *dns.Client) (interface{}, error) {
			return dnsClient.DeleteDomainRecord(args)
		})
		if err != nil {
			if IsExceptedErrors(err, []string{DomainRecordNotBelongToUser}) {
				return nil
			}
			if IsExceptedErrors(err, []string{RecordForbiddenDNSChange}) {
				return resource.RetryableError(fmt.Errorf("Operation forbidden because DNS is changing - trying again after change complete."))
			}
			return resource.NonRetryableError(fmt.Errorf("Error deleting domain record %s: %#v", d.Id(), err))
		}

		raw, err := client.WithDnsClient(func(dnsClient *dns.Client) (interface{}, error) {
			return dnsClient.DescribeDomainRecordInfoNew(&dns.DescribeDomainRecordInfoNewArgs{
				RecordId: d.Id(),
			})
		})
		if err != nil {
			if NotFoundError(err) || IsExceptedError(err, DomainRecordNotBelongToUser) {
				return nil
			}
			return resource.NonRetryableError(fmt.Errorf("Describe domain record got an error: %#v.", err))
		}
		response, _ := raw.(*dns.DescribeDomainRecordInfoNewResponse)
		if response == nil {
			return nil
		}

		return nil
	})
}

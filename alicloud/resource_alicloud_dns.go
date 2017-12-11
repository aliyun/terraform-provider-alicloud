package alicloud

import (
	"fmt"
	"time"

	"github.com/denverdino/aliyungo/common"
	"github.com/denverdino/aliyungo/dns"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceAlicloudDns() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudDnsCreate,
		Read:   resourceAlicloudDnsRead,
		Update: resourceAlicloudDnsUpdate,
		Delete: resourceAlicloudDnsDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateDomainName,
			},
			"group_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"dns_server": &schema.Schema{
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceAlicloudDnsCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AliyunClient).dnsconn

	args := &dns.AddDomainArgs{
		DomainName: d.Get("name").(string),
	}

	response, err := conn.AddDomain(args)
	if err != nil {
		return fmt.Errorf("AddDomain got an error: %#v", err)
	}

	d.SetId(response.DomainName)
	return resourceAlicloudDnsUpdate(d, meta)
}

func resourceAlicloudDnsUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AliyunClient).dnsconn

	d.Partial(true)

	args := &dns.ChangeDomainGroupArgs{
		DomainName: d.Get("name").(string),
	}

	if d.HasChange("group_id") {
		d.SetPartial("group_id")
		args.GroupId = d.Get("group_id").(string)

		_, err := conn.ChangeDomainGroup(args)
		if err != nil {
			return fmt.Errorf("ChangeDomainGroup got an error: %#v", err)
		}
	}

	d.Partial(false)
	return resourceAlicloudDnsRead(d, meta)
}

func resourceAlicloudDnsRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AliyunClient).dnsconn

	args := &dns.DescribeDomainInfoArgs{
		DomainName: d.Id(),
	}

	domain, err := conn.DescribeDomainInfo(args)
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("DescribeDomainInfo got an error: %#v", err)
	}

	d.Set("group_id", domain.GroupId)
	d.Set("name", domain.DomainName)
	d.Set("dns_server", domain.DnsServers.DnsServer)
	return nil
}

func resourceAlicloudDnsDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AliyunClient).dnsconn

	args := &dns.DeleteDomainArgs{
		DomainName: d.Id(),
	}

	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		_, err := conn.DeleteDomain(args)
		if err != nil {
			e, _ := err.(*common.Error)
			if e.ErrorResponse.Code == RecordForbiddenDNSChange {
				return resource.RetryableError(fmt.Errorf("Operation forbidden because DNS is changing - trying again after change complete."))
			}
			return resource.NonRetryableError(fmt.Errorf("Error deleting domain %s: %#v", d.Id(), err))
		}
		return nil
	})
}

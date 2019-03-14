package alicloud

import (
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/alidns"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
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
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateDomainName,
				ForceNew:     true,
			},
			"group_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"dns_server": {
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
	client := meta.(*connectivity.AliyunClient)

	request := alidns.CreateAddDomainRequest()
	request.DomainName = d.Get("name").(string)

	raw, err := client.WithDnsClient(func(dnsClient *alidns.Client) (interface{}, error) {
		return dnsClient.AddDomain(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "dns_domains", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw)
	response, _ := raw.(*alidns.AddDomainResponse)
	d.SetId(response.DomainName)
	return resourceAlicloudDnsUpdate(d, meta)
}

func resourceAlicloudDnsUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := alidns.CreateChangeDomainGroupRequest()
	request.DomainName = d.Get("name").(string)
	request.GroupId = d.Get("group_id").(string)

	if d.HasChange("group_id") {
		raw, err := client.WithDnsClient(func(dnsClient *alidns.Client) (interface{}, error) {
			return dnsClient.ChangeDomainGroup(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw)
	}
	return resourceAlicloudDnsRead(d, meta)
}

func resourceAlicloudDnsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	dnsService := &DnsService{client: client}
	domain, err := dnsService.DescribeDns(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("group_id", domain.GroupId)
	d.Set("name", domain.DomainName)
	d.Set("dns_server", domain.DnsServers.DnsServer)
	return nil
}

func resourceAlicloudDnsDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := alidns.CreateDeleteDomainRequest()
	request.DomainName = d.Id()

	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err := client.WithDnsClient(func(dnsClient *alidns.Client) (interface{}, error) {
			return dnsClient.DeleteDomain(request)
		})
		if err != nil {
			if IsExceptedError(err, RecordForbiddenDNSChange) {
				return resource.RetryableError(WrapErrorf(err, DeleteTimeoutMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR))
			}
			return resource.NonRetryableError(WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR))
		}
		addDebug(request.GetActionName(), raw)
		return nil
	})
}

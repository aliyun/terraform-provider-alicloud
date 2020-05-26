package alicloud

import (
	"fmt"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/alidns"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func resourceAlicloudAlidnsDomain() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudAlidnsDomainCreate,
		Read:   resourceAlicloudAlidnsDomainRead,
		Update: resourceAlicloudAlidnsDomainUpdate,
		Delete: resourceAlicloudAlidnsDomainDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"dns_servers": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"domain_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"domain_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"group_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"lang": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"remark": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"tags": tagsSchema(),
		},
	}
}

func resourceAlicloudAlidnsDomainCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := alidns.CreateAddDomainRequest()
	request.DomainName = d.Get("domain_name").(string)
	if v, ok := d.GetOk("group_id"); ok {
		request.GroupId = v.(string)
	}
	if v, ok := d.GetOk("lang"); ok {
		request.Lang = v.(string)
	}
	if v, ok := d.GetOk("resource_group_id"); ok {
		request.ResourceGroupId = v.(string)
	}
	raw, err := client.WithAlidnsClient(func(alidnsClient *alidns.Client) (interface{}, error) {
		return alidnsClient.AddDomain(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_alidns_domain", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw)
	response, _ := raw.(*alidns.AddDomainResponse)
	d.SetId(fmt.Sprintf("%v", response.DomainName))

	return resourceAlicloudAlidnsDomainUpdate(d, meta)
}
func resourceAlicloudAlidnsDomainRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	alidnsService := AlidnsService{client}
	object, err := alidnsService.DescribeAlidnsDomain(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("domain_name", d.Id())
	d.Set("dns_servers", object.DnsServers.DnsServer)
	d.Set("domain_id", object.DomainId)
	d.Set("group_id", object.GroupId)
	d.Set("remark", object.Remark)

	listTagResourcesObject, err := alidnsService.ListTagResources(d.Id())
	if err != nil {
		return WrapError(err)
	}

	tags := make(map[string]string)
	for _, t := range listTagResourcesObject.TagResources {
		tags[t.TagKey] = t.TagValue
	}
	d.Set("tags", tags)
	return nil
}
func resourceAlicloudAlidnsDomainUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	alidnsService := AlidnsService{client}
	d.Partial(true)

	if d.HasChange("tags") {
		if err := alidnsService.SetResourceTags(d, "DOMAIN"); err != nil {
			return WrapError(err)
		}
		d.SetPartial("tags")
	}
	update := false
	request := alidns.CreateChangeDomainGroupRequest()
	request.DomainName = d.Id()
	if !d.IsNewResource() && d.HasChange("group_id") {
		update = true
		request.GroupId = d.Get("group_id").(string)
	}
	if !d.IsNewResource() && d.HasChange("lang") {
		update = true
		request.Lang = d.Get("lang").(string)
	}
	if update {
		raw, err := client.WithAlidnsClient(func(alidnsClient *alidns.Client) (interface{}, error) {
			return alidnsClient.ChangeDomainGroup(request)
		})
		addDebug(request.GetActionName(), raw)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		d.SetPartial("group_id")
		d.SetPartial("lang")
	}
	update = false
	updateDomainRemarkReq := alidns.CreateUpdateDomainRemarkRequest()
	updateDomainRemarkReq.DomainName = d.Id()
	updateDomainRemarkReq.Lang = d.Get("lang").(string)
	if d.HasChange("remark") {
		update = true
		updateDomainRemarkReq.Remark = d.Get("remark").(string)
	}
	if update {
		raw, err := client.WithAlidnsClient(func(alidnsClient *alidns.Client) (interface{}, error) {
			return alidnsClient.UpdateDomainRemark(updateDomainRemarkReq)
		})
		addDebug(updateDomainRemarkReq.GetActionName(), raw)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), updateDomainRemarkReq.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		d.SetPartial("remark")
	}
	d.Partial(false)
	return resourceAlicloudAlidnsDomainRead(d, meta)
}
func resourceAlicloudAlidnsDomainDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	request := alidns.CreateDeleteDomainRequest()
	request.DomainName = d.Id()
	if v, ok := d.GetOk("lang"); ok {
		request.Lang = v.(string)
	}
	raw, err := client.WithAlidnsClient(func(alidnsClient *alidns.Client) (interface{}, error) {
		return alidnsClient.DeleteDomain(request)
	})
	addDebug(request.GetActionName(), raw)
	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidDomainName.NoExist"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	return nil
}

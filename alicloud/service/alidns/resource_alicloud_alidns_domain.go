package alidns

import (
	"fmt"
	"log"
	"time"

	alidnssdk "github.com/aliyun/alibaba-cloud-sdk-go/services/alidns"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/tags"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
		Timeouts: &schema.ResourceTimeout{
			Delete: schema.DefaultTimeout(6 * time.Minute),
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
			"group_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"lang": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"puny_code": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"remark": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"tags": tags.TagsSchema(),
		},
	}
}

func resourceAlicloudAlidnsDomainCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := alidnssdk.CreateAddDomainRequest()
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

	raw, err := client.WithAlidnsClient(func(alidnsClient *alidnssdk.Client) (interface{}, error) {
		return alidnsClient.AddDomain(request)
	})
	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_alidns_domain", request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
	}
	logs.AddDebug(request.GetActionName(), raw)
	response, _ := raw.(*alidnssdk.AddDomainResponse)
	d.SetId(fmt.Sprintf("%v", response.DomainName))

	return resourceAlicloudAlidnsDomainUpdate(d, meta)
}
func resourceAlicloudAlidnsDomainRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	alidnsService := NewAlidnsService(client)
	object, err := alidnsService.DescribeAlidnsDomain(d.Id())
	if err != nil {
		if errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_alidns_domain alidnsService.DescribeAlidnsDomain Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}

	d.Set("domain_name", d.Id())
	d.Set("dns_servers", object.DnsServers.DnsServer)
	d.Set("domain_id", object.DomainId)
	d.Set("group_id", object.GroupId)
	d.Set("group_name", object.GroupName)
	d.Set("puny_code", object.PunyCode)
	d.Set("remark", object.Remark)
	d.Set("resource_group_id", object.ResourceGroupId)

	listTagResourcesObject, err := alidnsService.ListTagResources(d.Id())
	if err != nil {
		return errs.WrapError(err)
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
	alidnsService := NewAlidnsService(client)
	d.Partial(true)

	if d.HasChange("tags") {
		if err := alidnsService.SetResourceTags(d, "DOMAIN"); err != nil {
			return errs.WrapError(err)
		}
	}
	update := false
	request := alidnssdk.CreateUpdateDomainRemarkRequest()
	request.DomainName = d.Id()
	if !d.IsNewResource() && d.HasChange("lang") {
		update = true
		request.Lang = d.Get("lang").(string)
	}
	if d.HasChange("remark") {
		update = true
		request.Remark = d.Get("remark").(string)
	}
	if update {
		raw, err := client.WithAlidnsClient(func(alidnsClient *alidnssdk.Client) (interface{}, error) {
			return alidnsClient.UpdateDomainRemark(request)
		})
		logs.AddDebug(request.GetActionName(), raw)
		if err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
		}
	}
	update = false
	changeDomainGroupReq := alidnssdk.CreateChangeDomainGroupRequest()
	changeDomainGroupReq.DomainName = d.Id()
	if !d.IsNewResource() && d.HasChange("group_id") {
		update = true
		changeDomainGroupReq.GroupId = d.Get("group_id").(string)
	}
	if !d.IsNewResource() && d.HasChange("lang") {
		update = true
		changeDomainGroupReq.Lang = d.Get("lang").(string)
	}
	if update {
		raw, err := client.WithAlidnsClient(func(alidnsClient *alidnssdk.Client) (interface{}, error) {
			return alidnsClient.ChangeDomainGroup(changeDomainGroupReq)
		})
		logs.AddDebug(changeDomainGroupReq.GetActionName(), raw)
		if err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), changeDomainGroupReq.GetActionName(), errs.AlibabaCloudSdkGoERROR)
		}
	}
	d.Partial(false)
	return resourceAlicloudAlidnsDomainRead(d, meta)
}
func resourceAlicloudAlidnsDomainDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	request := alidnssdk.CreateDeleteDomainRequest()
	request.DomainName = d.Id()
	if v, ok := d.GetOk("lang"); ok {
		request.Lang = v.(string)
	}
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err := retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		raw, err := client.WithAlidnsClient(func(alidnsClient *alidnssdk.Client) (interface{}, error) {
			return alidnsClient.DeleteDomain(request)
		})
		if err != nil {
			if errs.IsExpectedErrors(err, []string{"RecordForbidden.DNSChange", "DnsSystemBusyness", "InternalError"}) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		logs.AddDebug(request.GetActionName(), raw)
		return nil
	})
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"InvalidDomainName.NoExist"}) {
			return nil
		}
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
	}
	return nil
}

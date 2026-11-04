package alidns

import (
	"fmt"
	"log"

	alidnssdk "github.com/aliyun/alibaba-cloud-sdk-go/services/alidns"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlicloudAlidnsDomainGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudAlidnsDomainGroupCreate,
		Read:   resourceAlicloudAlidnsDomainGroupRead,
		Update: resourceAlicloudAlidnsDomainGroupUpdate,
		Delete: resourceAlicloudAlidnsDomainGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"domain_group_name": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"group_name"},
			},
			"group_name": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				Deprecated:    "Field 'group_name' has been deprecated from version 1.97.0. Use 'domain_group_name' instead.",
				ConflictsWith: []string{"domain_group_name"},
			},
			"lang": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceAlicloudAlidnsDomainGroupCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := alidnssdk.CreateAddDomainGroupRequest()
	if v, ok := d.GetOk("domain_group_name"); ok {
		request.GroupName = v.(string)
	} else if v, ok := d.GetOk("group_name"); ok {
		request.GroupName = v.(string)
	} else {
		return errs.WrapError(errs.Error(`[ERROR] Argument "group_name" or "domain_group_name" must be set one!`))
	}

	if v, ok := d.GetOk("lang"); ok {
		request.Lang = v.(string)
	}

	raw, err := client.WithAlidnsClient(func(alidnsClient *alidnssdk.Client) (interface{}, error) {
		return alidnsClient.AddDomainGroup(request)
	})
	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_alidns_domain_group", request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
	}
	logs.AddDebug(request.GetActionName(), raw)
	response, _ := raw.(*alidnssdk.AddDomainGroupResponse)
	d.SetId(fmt.Sprintf("%v", response.GroupId))

	return resourceAlicloudAlidnsDomainGroupRead(d, meta)
}
func resourceAlicloudAlidnsDomainGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	alidnsService := NewAlidnsService(client)
	object, err := alidnsService.DescribeAlidnsDomainGroup(d.Id())
	if err != nil {
		if errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_alidns_domain_group alidnsService.DescribeAlidnsDomainGroup Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}

	d.Set("domain_group_name", object.GroupName)
	d.Set("group_name", object.GroupName)
	return nil
}
func resourceAlicloudAlidnsDomainGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	update := false
	request := alidnssdk.CreateUpdateDomainGroupRequest()
	request.GroupId = d.Id()
	if d.HasChange("domain_group_name") {
		update = true
		request.GroupName = d.Get("domain_group_name").(string)
	}
	if d.HasChange("group_name") {
		update = true
		request.GroupName = d.Get("group_name").(string)
	}
	if request.GroupName == "" {
		request.GroupName = d.Get("domain_group_name").(string)
	}
	if d.HasChange("lang") {
		update = true
		request.Lang = d.Get("lang").(string)
	}
	if update {
		raw, err := client.WithAlidnsClient(func(alidnsClient *alidnssdk.Client) (interface{}, error) {
			return alidnsClient.UpdateDomainGroup(request)
		})
		logs.AddDebug(request.GetActionName(), raw)
		if err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
		}
	}
	return resourceAlicloudAlidnsDomainGroupRead(d, meta)
}
func resourceAlicloudAlidnsDomainGroupDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	request := alidnssdk.CreateDeleteDomainGroupRequest()
	request.GroupId = d.Id()
	if v, ok := d.GetOk("lang"); ok {
		request.Lang = v.(string)
	}
	raw, err := client.WithAlidnsClient(func(alidnsClient *alidnssdk.Client) (interface{}, error) {
		return alidnsClient.DeleteDomainGroup(request)
	})
	logs.AddDebug(request.GetActionName(), raw)
	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
	}
	return nil
}

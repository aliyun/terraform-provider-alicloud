package alicloud

import (
	"fmt"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/alidns"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
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
			"group_name": {
				Type:     schema.TypeString,
				Required: true,
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

	request := alidns.CreateAddDomainGroupRequest()
	request.GroupName = d.Get("group_name").(string)
	if v, ok := d.GetOk("lang"); ok {
		request.Lang = v.(string)
	}
	raw, err := client.WithAlidnsClient(func(alidnsClient *alidns.Client) (interface{}, error) {
		return alidnsClient.AddDomainGroup(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_alidns_domain_group", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw)
	response, _ := raw.(*alidns.AddDomainGroupResponse)
	d.SetId(fmt.Sprintf("%v", response.GroupId))

	return resourceAlicloudAlidnsDomainGroupRead(d, meta)
}
func resourceAlicloudAlidnsDomainGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	alidnsService := AlidnsService{client}
	object, err := alidnsService.DescribeAlidnsDomainGroup(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("group_name", object.GroupName)
	return nil
}
func resourceAlicloudAlidnsDomainGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	update := false
	request := alidns.CreateUpdateDomainGroupRequest()
	request.GroupId = d.Id()
	if d.HasChange("group_name") {
		update = true
	}
	request.GroupName = d.Get("group_name").(string)
	if d.HasChange("lang") {
		update = true
		request.Lang = d.Get("lang").(string)
	}
	if update {
		raw, err := client.WithAlidnsClient(func(alidnsClient *alidns.Client) (interface{}, error) {
			return alidnsClient.UpdateDomainGroup(request)
		})
		addDebug(request.GetActionName(), raw)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
	}
	return resourceAlicloudAlidnsDomainGroupRead(d, meta)
}
func resourceAlicloudAlidnsDomainGroupDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	request := alidns.CreateDeleteDomainGroupRequest()
	request.GroupId = d.Id()
	if v, ok := d.GetOk("lang"); ok {
		request.Lang = v.(string)
	}
	raw, err := client.WithAlidnsClient(func(alidnsClient *alidns.Client) (interface{}, error) {
		return alidnsClient.DeleteDomainGroup(request)
	})
	addDebug(request.GetActionName(), raw)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	return nil
}

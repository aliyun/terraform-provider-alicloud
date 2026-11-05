package alidns

import (
	"time"

	alidnssdk "github.com/aliyun/alibaba-cloud-sdk-go/services/alidns"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlicloudDnsGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudDnsGroupCreate,
		Read:   resourceAlicloudDnsGroupRead,
		Update: resourceAlicloudDnsGroupUpdate,
		Delete: resourceAlicloudDnsGroupDelete,

		DeprecationMessage: "This resource has been deprecated since v1.84.0 and will be removed in the future. Please use 'alicloud_alidns_domain_group' instead.",

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceAlicloudDnsGroupCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	request := alidnssdk.CreateAddDomainGroupRequest()
	request.RegionId = client.RegionId
	request.GroupName = d.Get("name").(string)

	raw, err := client.WithDnsClient(func(dnsClient *alidnssdk.Client) (interface{}, error) {
		return dnsClient.AddDomainGroup(request)
	})
	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_dns_group", request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
	}
	logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*alidnssdk.AddDomainGroupResponse)
	d.SetId(response.GroupId)
	return resourceAlicloudDnsGroupRead(d, meta)
}

func resourceAlicloudDnsGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := alidnssdk.CreateUpdateDomainGroupRequest()
	request.RegionId = client.RegionId
	request.GroupId = d.Id()

	if d.HasChange("name") {
		request.GroupName = d.Get("name").(string)
		raw, err := client.WithDnsClient(func(dnsClient *alidnssdk.Client) (interface{}, error) {
			return dnsClient.UpdateDomainGroup(request)
		})
		if err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
		}
		logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)
	}

	return resourceAlicloudDnsGroupRead(d, meta)
}

func resourceAlicloudDnsGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	dnsService := NewDnsService(client)
	object, err := dnsService.DescribeDnsGroup(d.Id())
	if err != nil {
		if errs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}
	d.Set("name", object.GroupName)
	return nil
}

func resourceAlicloudDnsGroupDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := alidnssdk.CreateDeleteDomainGroupRequest()
	request.RegionId = client.RegionId
	request.GroupId = d.Id()

	return retry.Retry(5*time.Minute, func() *retry.RetryError {
		raw, err := client.WithDnsClient(func(dnsClient *alidnssdk.Client) (interface{}, error) {
			return dnsClient.DeleteDomainGroup(request)
		})
		if err != nil {
			if errs.IsExpectedErrors(err, []string{"Fobidden.NotEmptyGroup"}) {
				return retry.RetryableError(errs.WrapErrorf(err, errs.DefaultTimeoutMsg, d.Id(), request.GetActionName(), errs.AlibabaCloudSdkGoERROR))
			}
			return retry.NonRetryableError(errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), request.GetActionName(), errs.AlibabaCloudSdkGoERROR))
		}
		logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})
}

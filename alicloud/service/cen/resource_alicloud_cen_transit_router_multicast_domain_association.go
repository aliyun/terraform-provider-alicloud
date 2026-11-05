package cen

import (
	"fmt"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAliCloudCenTransitRouterMulticastDomainAssociation() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudCenTransitRouterMulticastDomainAssociationCreate,
		Read:   resourceAliCloudCenTransitRouterMulticastDomainAssociationRead,
		Delete: resourceAliCloudCenTransitRouterMulticastDomainAssociationDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(3 * time.Minute),
			Delete: schema.DefaultTimeout(3 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"transit_router_multicast_domain_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"transit_router_attachment_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"vswitch_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAliCloudCenTransitRouterMulticastDomainAssociationCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cbnService := NewCbnService(client)
	var response map[string]interface{}
	action := "AssociateTransitRouterMulticastDomain"
	request := make(map[string]interface{})
	var err error

	request["ClientToken"] = helper.BuildClientToken("AssociateTransitRouterMulticastDomain")
	request["TransitRouterMulticastDomainId"] = d.Get("transit_router_multicast_domain_id")
	request["TransitRouterAttachmentId"] = d.Get("transit_router_attachment_id")

	vSwitchId := d.Get("vswitch_id").(string)
	request["VSwitchIds"] = []string{vSwitchId}

	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *retry.RetryError {
		response, err = client.RpcPost("Cbn", "2017-09-12", action, nil, request, true)
		if err != nil {
			if errs.IsExpectedErrors(err, []string{"Operation.Blocking"}) || helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)

	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_cen_transit_router_multicast_domain_association", action, errs.AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%v:%v:%v", request["TransitRouterMulticastDomainId"], request["TransitRouterAttachmentId"], vSwitchId))

	stateConf := helper.BuildStateConf([]string{}, []string{"Associated"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, cbnService.CenTransitRouterMulticastDomainAssociationStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errs.WrapErrorf(err, errs.IdMsg, d.Id())
	}

	return resourceAliCloudCenTransitRouterMulticastDomainAssociationRead(d, meta)
}

func resourceAliCloudCenTransitRouterMulticastDomainAssociationRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cbnService := NewCbnService(client)

	object, err := cbnService.DescribeCenTransitRouterMulticastDomainAssociation(d.Id())
	if err != nil {
		if !d.IsNewResource() && errs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}

	d.Set("transit_router_multicast_domain_id", object["TransitRouterMulticastDomainId"])
	d.Set("transit_router_attachment_id", object["TransitRouterAttachmentId"])
	d.Set("vswitch_id", object["VSwitchId"])
	d.Set("status", object["Status"])

	return nil
}

func resourceAliCloudCenTransitRouterMulticastDomainAssociationDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cbnService := NewCbnService(client)
	action := "DisassociateTransitRouterMulticastDomain"
	var response map[string]interface{}

	var err error

	parts, err := helper.ParseResourceId(d.Id(), 3)
	if err != nil {
		return errs.WrapError(err)
	}

	request := map[string]interface{}{
		"ClientToken":                    helper.BuildClientToken("DisassociateTransitRouterMulticastDomain"),
		"TransitRouterMulticastDomainId": parts[0],
		"TransitRouterAttachmentId":      parts[1],
		"VSwitchIds":                     []string{parts[2]},
	}

	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *retry.RetryError {
		response, err = client.RpcPost("Cbn", "2017-09-12", action, nil, request, true)
		if err != nil {
			if errs.IsExpectedErrors(err, []string{"Operation.Blocking", "InvalidOperation.MulticastGroupExist", "IncorrectStatus.Attachment"}) || helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)

	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
	}

	stateConf := helper.BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 5*time.Second, cbnService.CenTransitRouterMulticastDomainAssociationStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errs.WrapErrorf(err, errs.IdMsg, d.Id())
	}

	return nil
}

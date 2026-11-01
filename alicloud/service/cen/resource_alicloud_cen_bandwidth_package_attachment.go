package cen

import (
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/cbn"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlicloudCenBandwidthPackageAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudCenBandwidthPackageAttachmentCreate,
		Read:   resourceAlicloudCenBandwidthPackageAttachmentRead,
		Delete: resourceAlicloudCenBandwidthPackageAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"bandwidth_package_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAlicloudCenBandwidthPackageAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cenService := NewCenService(client)

	cenId := d.Get("instance_id").(string)
	cenBwpId := d.Get("bandwidth_package_id").(string)

	request := cbn.CreateAssociateCenBandwidthPackageRequest()
	request.CenId = cenId
	request.CenBandwidthPackageId = cenBwpId

	err := retry.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *retry.RetryError {
		raw, err := client.WithCbnClient(func(cbnClient *cbn.Client) (interface{}, error) {
			return cbnClient.AssociateCenBandwidthPackage(request)
		})
		if err != nil {
			if errs.IsExpectedErrors(err, []string{"InvalidOperation.BwpInstanceStatus", "InvalidOperation.BwpBusinessStatus", "InvalidOperation.CenInstanceStatus", "Operation.Blocking"}) || helper.NeedRetry(err) {
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})

	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_cen_bandwidth_package_attachment", request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
	}

	d.SetId(cenBwpId)

	if err := cenService.WaitForCenBandwidthPackageAttachment(d.Id(), InUse, DefaultCenTimeout); err != nil {
		return errs.WrapError(err)
	}

	return resourceAlicloudCenBandwidthPackageAttachmentRead(d, meta)
}

func resourceAlicloudCenBandwidthPackageAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cenService := NewCenService(client)

	object, err := cenService.DescribeCenBandwidthPackageAttachment(d.Id())
	if err != nil {
		if !d.IsNewResource() && errs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}

	d.Set("instance_id", object.CenIds.CenId[0])
	d.Set("bandwidth_package_id", object.CenBandwidthPackageId)

	return nil
}

func resourceAlicloudCenBandwidthPackageAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cenService := NewCenService(client)
	cenId := d.Get("instance_id").(string)
	cenBwpId := d.Get("bandwidth_package_id").(string)

	request := cbn.CreateUnassociateCenBandwidthPackageRequest()
	request.CenId = cenId
	request.CenBandwidthPackageId = cenBwpId

	err := retry.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *retry.RetryError {
		raw, err := client.WithCbnClient(func(cbnClient *cbn.Client) (interface{}, error) {
			return cbnClient.UnassociateCenBandwidthPackage(request)
		})

		if err != nil {
			if errs.IsExpectedErrors(err, []string{"InvalidOperation.BwpInstanceStatus", "InvalidOperation.BwpBusinessStatus", "InvalidOperation.CenInstanceStatus", "InvalidOperation.CenBandwidthLimitsNotZero", "Operation.Blocking"}) || helper.NeedRetry(err) {
				return retry.RetryableError(err)
			}

			return retry.NonRetryableError(err)
		}

		logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)

		return nil
	})

	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
	}

	return errs.WrapError(cenService.WaitForCenBandwidthPackageAttachment(cenBwpId, Deleted, DefaultCenTimeout))
}

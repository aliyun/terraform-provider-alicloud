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

func resourceAlicloudCenRouteEntry() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudCenRouteEntryCreate,
		Read:   resourceAlicloudCenRouteEntryRead,
		Delete: resourceAlicloudCenRouteEntryDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"route_table_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"cidr_block": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAlicloudCenRouteEntryCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cenService := NewCenService(client)
	cenId := d.Get("instance_id").(string)
	vtbId := d.Get("route_table_id").(string)
	cidr := d.Get("cidr_block").(string)
	childInstanceId, childInstanceType, err := cenService.CreateCenRouteEntryParas(vtbId)
	if err != nil {
		return errs.WrapError(err)
	}

	request := cbn.CreatePublishRouteEntriesRequest()
	request.RegionId = client.RegionId
	request.CenId = cenId
	request.ChildInstanceId = childInstanceId
	request.ChildInstanceType = childInstanceType
	request.ChildInstanceRegionId = client.RegionId
	request.ChildInstanceRouteTableId = vtbId
	request.DestinationCidrBlock = cidr

	err = retry.Retry(3*time.Minute, func() *retry.RetryError {
		raw, err := client.WithCbnClient(func(cbnClient *cbn.Client) (interface{}, error) {
			return cbnClient.PublishRouteEntries(request)
		})
		if err != nil {
			if errs.IsExpectedErrors(err, []string{"Operation.Blocking", "not in a valid state for the operation"}) {
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})
	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_cen_route_entry", request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
	}

	d.SetId(cenId + COLON_SEPARATED + vtbId + COLON_SEPARATED + cidr)

	err = cenService.WaitForCenRouterEntry(d.Id(), PUBLISHED, DefaultCenTimeout)
	if err != nil {
		return errs.WrapError(err)
	}

	return resourceAlicloudCenRouteEntryRead(d, meta)
}

func resourceAlicloudCenRouteEntryRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cenService := NewCenService(client)

	parts, err := helper.ParseResourceId(d.Id(), 3)
	if err != nil {
		return errs.WrapError(err)
	}
	cenId := parts[0]

	object, err := cenService.DescribeCenRouteEntry(d.Id())
	if err != nil {
		if errs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}

	if object.PublishStatus == string(NOPUBLISHED) {
		d.SetId("")
		return nil
	}

	d.Set("instance_id", cenId)
	d.Set("route_table_id", object.ChildInstanceRouteTableId)
	d.Set("cidr_block", object.DestinationCidrBlock)

	return nil
}

func resourceAlicloudCenRouteEntryDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cenService := NewCenService(client)

	cenId := d.Get("instance_id").(string)
	vtbId := d.Get("route_table_id").(string)
	cidr := d.Get("cidr_block").(string)
	childInstanceId, childInstanceType, err := cenService.CreateCenRouteEntryParas(vtbId)
	if err != nil {
		if errs.NotFoundError(err) {
			return nil
		}
		return errs.WrapError(err)
	}

	request := cbn.CreateWithdrawPublishedRouteEntriesRequest()
	request.RegionId = client.RegionId
	request.CenId = cenId
	request.ChildInstanceId = childInstanceId
	request.ChildInstanceType = childInstanceType
	request.ChildInstanceRegionId = client.RegionId
	request.ChildInstanceRouteTableId = vtbId
	request.DestinationCidrBlock = cidr

	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		raw, err := client.WithCbnClient(func(cbnClient *cbn.Client) (interface{}, error) {
			return cbnClient.WithdrawPublishedRouteEntries(request)
		})
		if err != nil {
			if errs.IsExpectedErrors(err, []string{"InvalidOperation.CenInstanceStatus", "InternalError"}) {
				return retry.RetryableError(err)
			}

			return retry.NonRetryableError(err)
		}
		logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"InvalidOperation.NotFoundRoute", "The instance is not exist"}) {
			return nil
		}
		return errs.WrapErrorf(err, errs.DataDefaultErrorMsg, d.Id(), request.GetActionName(), errs.AlibabaCloudSdkGoERROR)

	}
	return errs.WrapError(cenService.WaitForCenRouterEntry(d.Id(), Deleted, DefaultCenTimeoutLong))
}

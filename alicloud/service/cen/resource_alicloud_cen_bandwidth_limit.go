package cen

import (
	"fmt"
	"strings"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAlicloudCenBandwidthLimit() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudCenBandwidthLimitCreate,
		Read:   resourceAlicloudCenBandwidthLimitRead,
		Update: resourceAlicloudCenBandwidthLimitUpdate,
		Delete: resourceAlicloudCenBandwidthLimitDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"region_ids": {
				Type:     schema.TypeSet,
				Required: true,
				ForceNew: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				MaxItems: 2,
				MinItems: 2,
			},
			"bandwidth_limit": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntAtLeast(1),
			},
		},
	}
}

func resourceAlicloudCenBandwidthLimitCreate(d *schema.ResourceData, meta interface{}) error {
	cenId := d.Get("instance_id").(string)
	regionIds := d.Get("region_ids").(*schema.Set).List()
	if len(regionIds) != 2 {
		return errs.WrapError(errs.Error("Two different region ids should be set for bandwidth limit. "))
	}

	localRegionId := regionIds[0].(string)
	oppositeRegionId := regionIds[1].(string)

	if strings.Compare(localRegionId, oppositeRegionId) <= 0 {
		d.SetId(cenId + COLON_SEPARATED + localRegionId + COLON_SEPARATED + oppositeRegionId)
	} else {
		d.SetId(cenId + COLON_SEPARATED + oppositeRegionId + COLON_SEPARATED + localRegionId)
	}

	return resourceAlicloudCenBandwidthLimitUpdate(d, meta)
}

func resourceAlicloudCenBandwidthLimitRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cenService := NewCenService(client)

	object, err := cenService.DescribeCenBandwidthLimit(d.Id())
	if err != nil {
		if errs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}

	respRegionIds := make([]string, 0)
	respRegionIds = append(respRegionIds, object.LocalRegionId, object.OppositeRegionId)

	d.Set("region_ids", respRegionIds)
	d.Set("instance_id", object.CenId)
	d.Set("bandwidth_limit", object.BandwidthLimit)

	return nil
}

func resourceAlicloudCenBandwidthLimitUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cenService := NewCenService(client)
	cenId := d.Get("instance_id").(string)

	regionIds := d.Get("region_ids").(*schema.Set).List()
	if len(regionIds) != 2 {
		return errs.WrapError(errs.Error("Two different region ids should be set for bandwidth limit. "))
	}

	localRegionId := regionIds[0].(string)
	oppositeRegionId := regionIds[1].(string)
	var bandwidthLimit int

	if d.HasChange("bandwidth_limit") {
		bandwidthLimit = d.Get("bandwidth_limit").(int)
		if bandwidthLimit == 0 {
			return errs.WrapError(errs.Error("the bandwidth limit should be at least than 1 Mbps"))
		}

		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err := retry.Retry(5*time.Minute, func() *retry.RetryError {
			err := cenService.SetCenInterRegionBandwidthLimit(cenId, localRegionId, oppositeRegionId, bandwidthLimit)
			if err != nil {
				if errs.IsExpectedErrors(err, []string{"InvalidOperation.CenInstanceStatus", "Operation.Blocking"}) || helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		if err != nil {
			return errs.WrapError(err)
		}

		err = retry.Retry(5*time.Minute, func() *retry.RetryError {

			stateConf := helper.BuildStateConf([]string{"Modifying"}, []string{"Active"}, d.Timeout(schema.TimeoutUpdate), 3*time.Second, cenService.CenBandwidthLimitStateRefreshFunc(d.Id(), []string{}))

			if _, err = stateConf.WaitForState(); err != nil {
				if errs.IsExpectedErrors(err, []string{errs.ThrottlingUser}) {
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		if err != nil {
			return errs.WrapErrorf(err, errs.IdMsg, d.Id())
		}
	}
	return resourceAlicloudCenBandwidthLimitRead(d, meta)
}

func resourceAlicloudCenBandwidthLimitDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cenService := NewCenService(client)
	cenId := d.Get("instance_id").(string)

	regionIds := d.Get("region_ids").(*schema.Set).List()
	if len(regionIds) != 2 {
		return fmt.Errorf("Two different region ids should be set for bandwidth limit")
	}

	localRegionId := regionIds[0].(string)
	oppositeRegionId := regionIds[1].(string)

	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err := retry.Retry(5*time.Minute, func() *retry.RetryError {
		err := cenService.SetCenInterRegionBandwidthLimit(cenId, localRegionId, oppositeRegionId, 0)
		if err != nil {
			if errs.IsExpectedErrors(err, []string{"InvalidOperation.CenInstanceStatus", "Operation.Blocking"}) || helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	if err != nil {
		return errs.WrapError(err)
	}

	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		stateConf := helper.BuildStateConf([]string{"Active", "Modifying"}, []string{}, d.Timeout(schema.TimeoutDelete), 3*time.Second, cenService.CenBandwidthLimitStateRefreshFunc(d.Id(), []string{}))

		_, err = stateConf.WaitForState()
		if errs.IsExpectedErrors(err, []string{errs.ThrottlingUser}) {
			return retry.RetryableError(err)
		}
		return retry.NonRetryableError(err)
	})
	if err != nil {
		return errs.WrapErrorf(err, errs.IdMsg, d.Id())
	}
	return nil
}

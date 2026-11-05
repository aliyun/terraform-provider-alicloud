package alb

import (
	"fmt"
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceLoadBalancerCommonBandwidthPackageAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceLoadBalancerCommonBandwidthPackageAttachmentCreate,
		Read:   resourceLoadBalancerCommonBandwidthPackageAttachmentRead,
		Update: resourceLoadBalancerCommonBandwidthPackageAttachmentUpdate,
		Delete: resourceLoadBalancerCommonBandwidthPackageAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"bandwidth_package_id": {
				Required: true,
				ForceNew: true,
				Type:     schema.TypeString,
			},
			"dry_run": {
				Optional: true,
				Type:     schema.TypeBool,
			},
			"load_balancer_id": {
				Required: true,
				ForceNew: true,
				Type:     schema.TypeString,
			},
			"status": {
				Computed: true,
				Type:     schema.TypeString,
			},
		},
	}
}

func resourceLoadBalancerCommonBandwidthPackageAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	albService := AlbService{client}
	request := map[string]interface{}{
		"RegionId": client.RegionId,
	}
	var err error

	if v, ok := d.GetOk("bandwidth_package_id"); ok {
		request["BandwidthPackageId"] = v
	}
	if v, ok := d.GetOk("dry_run"); ok {
		request["DryRun"] = v
	}
	if v, ok := d.GetOk("load_balancer_id"); ok {
		request["LoadBalancerId"] = v
	}

	request["ClientToken"] = helper.BuildClientToken("AttachCommonBandwidthPackageToLoadBalancer")
	var response map[string]interface{}
	action := "AttachCommonBandwidthPackageToLoadBalancer"
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *retry.RetryError {
		resp, err := client.RpcPost("Alb", "2020-06-16", action, nil, request, true)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		response = resp
		logs.AddDebug(action, response, request)
		return nil
	})
	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_alb_load_balancer_common_bandwidth_package_attachment", action, errs.AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(request["LoadBalancerId"], ":", request["BandwidthPackageId"]))
	stateConf := helper.BuildStateConf([]string{}, []string{"Active"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, albService.AlbLoadBalancerCommonBandwidthPackageAttachmentStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errs.WrapErrorf(err, errs.IdMsg, d.Id())
	}
	return resourceLoadBalancerCommonBandwidthPackageAttachmentRead(d, meta)
}

func resourceLoadBalancerCommonBandwidthPackageAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	albService := AlbService{client}

	object, err := albService.DescribeAlbLoadBalancerCommonBandwidthPackageAttachment(d.Id())
	if err != nil {
		if errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_alb_load_balancer_common_bandwidth_package_attachment albService.DescribeAlbLoadBalancerCommonBandwidthPackageAttachment Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}
	parts, err := helper.ParseResourceId(d.Id(), 2)
	if err != nil {
		return errs.WrapError(err)
	}
	d.Set("load_balancer_id", parts[0])
	d.Set("bandwidth_package_id", parts[1])

	loadBalancerStatus := object["LoadBalancerStatus"]
	d.Set("status", loadBalancerStatus)

	return nil
}
func resourceLoadBalancerCommonBandwidthPackageAttachmentUpdate(d *schema.ResourceData, meta interface{}) error {
	return resourceLoadBalancerCommonBandwidthPackageAttachmentRead(d, meta)
}

func resourceLoadBalancerCommonBandwidthPackageAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	albService := AlbService{client}
	var err error
	parts, err := helper.ParseResourceId(d.Id(), 2)
	if err != nil {
		return errs.WrapError(err)
	}

	request := map[string]interface{}{
		"LoadBalancerId": parts[0], "BandwidthPackageId": parts[1],
		"RegionId": client.RegionId,
	}

	if v, ok := d.GetOk("dry_run"); ok {
		request["DryRun"] = v
	}

	request["ClientToken"] = helper.BuildClientToken("DetachCommonBandwidthPackageFromLoadBalancer")
	action := "DetachCommonBandwidthPackageFromLoadBalancer"
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *retry.RetryError {
		resp, err := client.RpcPost("Alb", "2020-06-16", action, nil, request, true)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		logs.AddDebug(action, resp, request)
		return nil
	})
	if err != nil {
		if errs.NotFoundError(err) {
			return nil
		}
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
	}
	stateConf := helper.BuildStateConf([]string{}, []string{"Active"}, d.Timeout(schema.TimeoutDelete), 5*time.Second, albService.AlbLoadBalancerCommonBandwidthPackageAttachmentStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errs.WrapErrorf(err, errs.IdMsg, d.Id())
	}
	return nil
}

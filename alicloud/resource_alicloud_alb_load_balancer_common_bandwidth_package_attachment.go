package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudAlbLoadBalancerCommonBandwidthPackageAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudAlbLoadBalancerCommonBandwidthPackageAttachmentCreate,
		Read:   resourceAlicloudAlbLoadBalancerCommonBandwidthPackageAttachmentRead,
		Update: resourceAlicloudAlbLoadBalancerCommonBandwidthPackageAttachmentUpdate,
		Delete: resourceAlicloudAlbLoadBalancerCommonBandwidthPackageAttachmentDelete,
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

func resourceAlicloudAlbLoadBalancerCommonBandwidthPackageAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
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

	request["ClientToken"] = buildClientToken("AttachCommonBandwidthPackageToLoadBalancer")
	var response map[string]interface{}
	action := "AttachCommonBandwidthPackageToLoadBalancer"
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *resource.RetryError {
		resp, err := client.RpcPost("Alb", "2020-06-16", action, nil, request, true)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		response = resp
		addDebug(action, response, request)
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_alb_load_balancer_common_bandwidth_package_attachment", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(request["LoadBalancerId"], ":", request["BandwidthPackageId"]))
	stateConf := BuildStateConf([]string{}, []string{"Active"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, albService.AlbLoadBalancerCommonBandwidthPackageAttachmentStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return resourceAlicloudAlbLoadBalancerCommonBandwidthPackageAttachmentRead(d, meta)
}

func resourceAlicloudAlbLoadBalancerCommonBandwidthPackageAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	albService := AlbService{client}

	object, err := albService.DescribeAlbLoadBalancerCommonBandwidthPackageAttachment(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_alb_load_balancer_common_bandwidth_package_attachment albService.DescribeAlbLoadBalancerCommonBandwidthPackageAttachment Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	d.Set("load_balancer_id", parts[0])
	d.Set("bandwidth_package_id", parts[1])

	loadBalancerStatus := object["LoadBalancerStatus"]
	d.Set("status", loadBalancerStatus)

	return nil
}
func resourceAlicloudAlbLoadBalancerCommonBandwidthPackageAttachmentUpdate(d *schema.ResourceData, meta interface{}) error {
	return resourceAlicloudAlbLoadBalancerCommonBandwidthPackageAttachmentRead(d, meta)
}

func resourceAlicloudAlbLoadBalancerCommonBandwidthPackageAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	albService := AlbService{client}
	var err error
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}

	request := map[string]interface{}{
		"LoadBalancerId": parts[0], "BandwidthPackageId": parts[1],
		"RegionId": client.RegionId,
	}

	if v, ok := d.GetOk("dry_run"); ok {
		request["DryRun"] = v
	}

	request["ClientToken"] = buildClientToken("DetachCommonBandwidthPackageFromLoadBalancer")
	action := "DetachCommonBandwidthPackageFromLoadBalancer"
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *resource.RetryError {
		resp, err := client.RpcPost("Alb", "2020-06-16", action, nil, request, true)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, resp, request)
		return nil
	})
	if err != nil {
		if NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	stateConf := BuildStateConf([]string{}, []string{"Active"}, d.Timeout(schema.TimeoutDelete), 5*time.Second, albService.AlbLoadBalancerCommonBandwidthPackageAttachmentStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return nil
}

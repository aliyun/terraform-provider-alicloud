package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAlicloudCenTransitRouterGrantAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudCenTransitRouterGrantAttachmentCreate,
		Read:   resourceAlicloudCenTransitRouterGrantAttachmentRead,
		Delete: resourceAlicloudCenTransitRouterGrantAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(1 * time.Minute),
			Delete: schema.DefaultTimeout(1 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"cen_id": {
				Required: true,
				ForceNew: true,
				Type:     schema.TypeString,
			},
			"cen_owner_id": {
				Required: true,
				ForceNew: true,
				Type:     schema.TypeString,
			},
			"instance_id": {
				Required: true,
				ForceNew: true,
				Type:     schema.TypeString,
			},
			"instance_type": {
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"VPC", "ExpressConnect", "VPN"}, false),
				Type:         schema.TypeString,
			},
			"order_type": {
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"PayByCenOwner", "PayByResourceOwner"}, false),
				Type:         schema.TypeString,
			},
		},
	}
}

func resourceAlicloudCenTransitRouterGrantAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	request := map[string]interface{}{
		"RegionId": client.RegionId,
	}
	var err error

	if v, ok := d.GetOk("cen_id"); ok {
		request["CenId"] = v
	}
	if v, ok := d.GetOk("cen_owner_id"); ok {
		request["CenOwnerId"] = v
	}
	if v, ok := d.GetOk("instance_id"); ok {
		request["InstanceId"] = v
	}
	if v, ok := d.GetOk("instance_type"); ok {
		request["InstanceType"] = v
	}
	if v, ok := d.GetOk("order_type"); ok {
		request["OrderType"] = v
	}

	var response map[string]interface{}
	action := "GrantInstanceToTransitRouter"
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *resource.RetryError {
		resp, err := client.RpcPost("Cbn", "2017-09-12", action, nil, request, false)
		if err != nil {
			if IsExpectedErrors(err, []string{"OperationFailed.TaskConflict"}) || NeedRetry(err) {
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_cen_transit_router_grant_attachment", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(request["InstanceType"], ":", request["InstanceId"], ":", request["CenOwnerId"], ":", request["CenId"]))

	return resourceAlicloudCenTransitRouterGrantAttachmentRead(d, meta)
}

func resourceAlicloudCenTransitRouterGrantAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cbnService := CbnService{client}

	object, err := cbnService.DescribeCenTransitRouterGrantAttachment(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_cen_transit_router_grant_attachment cbnService.DescribeCenTransitRouterGrantAttachment Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	parts, err := ParseResourceId(d.Id(), 4)
	if err != nil {
		return WrapError(err)
	}
	d.Set("instance_type", parts[0])
	d.Set("instance_id", parts[1])
	d.Set("cen_owner_id", parts[2])
	d.Set("cen_id", parts[3])
	d.Set("order_type", object["OrderType"])

	return nil
}

func resourceAlicloudCenTransitRouterGrantAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var err error
	parts, err := ParseResourceId(d.Id(), 4)
	if err != nil {
		return WrapError(err)
	}

	request := map[string]interface{}{
		"InstanceType": parts[0],
		"InstanceId":   parts[1],
		"CenOwnerId":   parts[2],
		"CenId":        parts[3],
		"RegionId":     client.RegionId,
	}

	action := "RevokeInstanceFromTransitRouter"
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *resource.RetryError {
		resp, err := client.RpcPost("Cbn", "2017-09-12", action, nil, request, false)
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
	return nil
}

package alicloud

import (
	"fmt"
	"log"
	"time"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudCenChildInstanceRouteEntryToAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudCenChildInstanceRouteEntryToAttachmentCreate,
		Update: resourceAlicloudCenChildInstanceRouteEntryToAttachmentUpdate,
		Read:   resourceAlicloudCenChildInstanceRouteEntryToAttachmentRead,
		Delete: resourceAlicloudCenChildInstanceRouteEntryToAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(6 * time.Minute),
			Delete: schema.DefaultTimeout(6 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"cen_id": {
				Required: true,
				ForceNew: true,
				Type:     schema.TypeString,
			},
			"child_instance_route_table_id": {
				Required: true,
				ForceNew: true,
				Type:     schema.TypeString,
			},
			"destination_cidr_block": {
				Required: true,
				ForceNew: true,
				Type:     schema.TypeString,
			},
			"dry_run": {
				Optional: true,
				Type:     schema.TypeBool,
			},
			"service_type": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"status": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"transit_router_attachment_id": {
				Required: true,
				ForceNew: true,
				Type:     schema.TypeString,
			},
		},
	}
}

func resourceAlicloudCenChildInstanceRouteEntryToAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cbnService := CbnService{client}
	request := make(map[string]interface{})
	var err error

	if v, ok := d.GetOk("cen_id"); ok {
		request["CenId"] = v
	}
	if v, ok := d.GetOk("child_instance_route_table_id"); ok {
		request["RouteTableId"] = v
	}
	if v, ok := d.GetOk("destination_cidr_block"); ok {
		request["DestinationCidrBlock"] = v
	}
	if v, ok := d.GetOk("dry_run"); ok {
		request["DryRun"] = v
	}
	if v, ok := d.GetOk("transit_router_attachment_id"); ok {
		request["TransitRouterAttachmentId"] = v
	}

	request["ClientToken"] = buildClientToken("CreateCenChildInstanceRouteEntryToAttachment")
	var response map[string]interface{}
	action := "CreateCenChildInstanceRouteEntryToAttachment"
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *resource.RetryError {
		resp, err := client.RpcPost("Cbn", "2017-09-12", action, nil, request, true)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_cen_child_instance_route_entry_to_attachment", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(request["CenId"], ":", request["RouteTableId"], ":", request["TransitRouterAttachmentId"], ":", request["DestinationCidrBlock"]))
	stateConf := BuildStateConf([]string{}, []string{"Available"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, cbnService.CenChildInstanceRouteEntryToAttachmentStateRefreshFunc(d, []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return resourceAlicloudCenChildInstanceRouteEntryToAttachmentRead(d, meta)
}

func resourceAlicloudCenChildInstanceRouteEntryToAttachmentUpdate(d *schema.ResourceData, meta interface{}) error {
	return resourceAlicloudCenChildInstanceRouteEntryToAttachmentRead(d, meta)
}

func resourceAlicloudCenChildInstanceRouteEntryToAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cbnService := CbnService{client}

	object, err := cbnService.DescribeCenChildInstanceRouteEntryToAttachment(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_cen_child_instance_route_entry_to_attachment cbnService.DescribeCenChildInstanceRouteEntryToAttachment Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	parts, err := ParseResourceId(d.Id(), 4)
	if err != nil {
		return WrapError(err)
	}
	d.Set("cen_id", parts[0])
	d.Set("child_instance_route_table_id", parts[1])
	d.Set("transit_router_attachment_id", parts[2])
	d.Set("destination_cidr_block", parts[3])
	d.Set("service_type", object["ServiceType"])
	d.Set("status", object["Status"])

	return nil
}

func resourceAlicloudCenChildInstanceRouteEntryToAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cbnService := CbnService{client}
	var err error
	parts, err := ParseResourceId(d.Id(), 4)
	if err != nil {
		return WrapError(err)
	}

	request := map[string]interface{}{
		"CenId": parts[0], "RouteTableId": parts[1], "TransitRouterAttachmentId": parts[2], "DestinationCidrBlock": parts[3],
	}

	if v, ok := d.GetOk("dry_run"); ok {
		request["DryRun"] = v
	}

	request["ClientToken"] = buildClientToken("DeleteCenChildInstanceRouteEntryToAttachment")
	action := "DeleteCenChildInstanceRouteEntryToAttachment"
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *resource.RetryError {
		resp, err := client.RpcPost("Cbn", "2017-09-12", action, nil, request, true)
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
	stateConf := BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 5*time.Second, cbnService.CenChildInstanceRouteEntryToAttachmentStateRefreshFunc(d, []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return nil
}

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

func resourceAliCloudCenTransitRouterVbrAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudCenTransitRouterVbrAttachmentCreate,
		Read:   resourceAliCloudCenTransitRouterVbrAttachmentRead,
		Update: resourceAliCloudCenTransitRouterVbrAttachmentUpdate,
		Delete: resourceAliCloudCenTransitRouterVbrAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"cen_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"vbr_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"transit_router_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"resource_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      "VBR",
				ValidateFunc: StringInSlice([]string{"VBR"}, false),
			},
			"vbr_owner_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"auto_publish_route_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"transit_router_attachment_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"transit_router_attachment_description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"tags": tagsSchema(),
			"dry_run": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"transit_router_attachment_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"route_table_association_enabled": {
				Type:       schema.TypeBool,
				Optional:   true,
				Deprecated: "Field `route_table_association_enabled` has been deprecated from provider version 1.233.1. Please use the resource `alicloud_cen_transit_router_route_table_association` instead.",
			},
			"route_table_propagation_enabled": {
				Type:       schema.TypeBool,
				Optional:   true,
				Deprecated: "Field `route_table_propagation_enabled` has been deprecated from provider version 1.233.1. Please use the resource `alicloud_cen_transit_router_route_table_propagation` instead.",
			},
		},
	}
}

func resourceAliCloudCenTransitRouterVbrAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cbnService := CbnService{client}
	var response map[string]interface{}
	action := "CreateTransitRouterVbrAttachment"
	request := make(map[string]interface{})
	conn, err := client.NewCbnClient()
	if err != nil {
		return WrapError(err)
	}

	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken("CreateTransitRouterVbrAttachment")
	request["CenId"] = d.Get("cen_id")
	request["VbrId"] = d.Get("vbr_id")

	if v, ok := d.GetOk("transit_router_id"); ok {
		request["TransitRouterId"] = v
	}

	if v, ok := d.GetOk("resource_type"); ok {
		request["ResourceType"] = v
	}

	if v, ok := d.GetOk("vbr_owner_id"); ok {
		request["VbrOwnerId"] = v
	}

	if v, ok := d.GetOkExists("auto_publish_route_enabled"); ok {
		request["AutoPublishRouteEnabled"] = v
	}

	if v, ok := d.GetOk("transit_router_attachment_name"); ok {
		request["TransitRouterAttachmentName"] = v
	}

	if v, ok := d.GetOk("transit_router_attachment_description"); ok {
		request["TransitRouterAttachmentDescription"] = v
	}

	if v, ok := d.GetOk("tags"); ok {
		tagsMap := ConvertTags(v.(map[string]interface{}))
		request["Tag"] = tagsMap
	}

	if v, ok := d.GetOkExists("route_table_association_enabled"); ok {
		request["RouteTableAssociationEnabled"] = v
	}

	if v, ok := d.GetOkExists("route_table_propagation_enabled"); ok {
		request["RouteTablePropagationEnabled"] = v
	}

	if v, ok := d.GetOkExists("dry_run"); ok {
		request["DryRun"] = v
	}

	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-09-12"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			if IsExpectedErrors(err, []string{"Operation.Blocking", "IncorrectStatus.Status", "InstanceStatus.NotSupport"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_cen_transit_router_vbr_attachment", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%v:%v", request["CenId"], response["TransitRouterAttachmentId"]))

	stateConf := BuildStateConf([]string{}, []string{"Attached"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, cbnService.CenTransitRouterVbrAttachmentStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAliCloudCenTransitRouterVbrAttachmentRead(d, meta)
}

func resourceAliCloudCenTransitRouterVbrAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cbnService := CbnService{client}

	object, err := cbnService.DescribeCenTransitRouterVbrAttachment(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_cen_transit_router_vbr_attachment cbnService.DescribeCenTransitRouterVbrAttachment Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("cen_id", object["CenId"])
	d.Set("vbr_id", object["VbrId"])
	d.Set("transit_router_id", object["TransitRouterId"])
	d.Set("resource_type", object["ResourceType"])
	d.Set("vbr_owner_id", object["VbrOwnerId"])
	d.Set("auto_publish_route_enabled", object["AutoPublishRouteEnabled"])
	d.Set("transit_router_attachment_name", object["TransitRouterAttachmentName"])
	d.Set("transit_router_attachment_description", object["TransitRouterAttachmentDescription"])
	d.Set("transit_router_attachment_id", object["TransitRouterAttachmentId"])
	d.Set("status", object["Status"])

	listTagResourcesObject, err := cbnService.ListTagResources(d.Id(), "TransitRouterVbrAttachment")
	if err != nil {
		return WrapError(err)
	}

	d.Set("tags", tagsToMap(listTagResourcesObject))

	return nil
}

func resourceAliCloudCenTransitRouterVbrAttachmentUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cbnService := CbnService{client}
	var response map[string]interface{}
	d.Partial(true)

	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}

	update := false
	request := map[string]interface{}{
		"ClientToken":               buildClientToken("UpdateTransitRouterVbrAttachmentAttribute"),
		"TransitRouterAttachmentId": parts[1],
	}

	if d.HasChange("auto_publish_route_enabled") {
		update = true

		if v, ok := d.GetOkExists("auto_publish_route_enabled"); ok {
			request["AutoPublishRouteEnabled"] = v
		}
	}

	if d.HasChange("transit_router_attachment_name") {
		update = true
	}
	if v, ok := d.GetOk("transit_router_attachment_name"); ok {
		request["TransitRouterAttachmentName"] = v
	}

	if d.HasChange("transit_router_attachment_description") {
		update = true
	}
	if v, ok := d.GetOk("transit_router_attachment_description"); ok {
		request["TransitRouterAttachmentDescription"] = v
	}

	if v, ok := d.GetOkExists("dry_run"); ok {
		request["DryRun"] = v
	}

	if update {
		action := "UpdateTransitRouterVbrAttachmentAttribute"
		conn, err := client.NewCbnClient()
		if err != nil {
			return WrapError(err)
		}

		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-09-12"), StringPointer("AK"), nil, request, &runtime)
			if err != nil {
				if IsExpectedErrors(err, []string{"Operation.Blocking", "IncorrectStatus.Status"}) || NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, request)

		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}

		stateConf := BuildStateConf([]string{}, []string{"Attached"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, cbnService.CenTransitRouterVbrAttachmentStateRefreshFunc(d.Id(), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}

		d.SetPartial("auto_publish_route_enabled")
		d.SetPartial("transit_router_attachment_name")
		d.SetPartial("transit_router_attachment_description")
	}

	if d.HasChange("tags") {
		if err := cbnService.SetResourceTags(d, "TransitRouterVbrAttachment"); err != nil {
			return WrapError(err)
		}

		d.SetPartial("tags")
	}

	d.Partial(false)

	return resourceAliCloudCenTransitRouterVbrAttachmentRead(d, meta)
}

func resourceAliCloudCenTransitRouterVbrAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cbnService := CbnService{client}
	action := "DeleteTransitRouterVbrAttachment"
	var response map[string]interface{}

	conn, err := client.NewCbnClient()
	if err != nil {
		return WrapError(err)
	}

	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}

	request := map[string]interface{}{
		"ClientToken":               buildClientToken("DeleteTransitRouterVbrAttachment"),
		"TransitRouterAttachmentId": parts[1],
	}

	if v, ok := d.GetOkExists("dry_run"); ok {
		request["DryRun"] = v
	}

	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-09-12"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			if IsExpectedErrors(err, []string{"Operation.Blocking", "IncorrectStatus.Status", "InstanceStatus.NotSupport"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	stateConf := BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 5*time.Second, cbnService.CenTransitRouterVbrAttachmentStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return nil
}

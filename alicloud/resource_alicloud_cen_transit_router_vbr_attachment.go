package alicloud

import (
	"fmt"
	"log"
	"strings"
	"time"

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
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(15 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"auto_publish_route_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"cen_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"order_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: StringInSlice([]string{"PayByCenOwner", "PayByResourceOwner"}, false),
			},
			"region_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": tagsSchema(),
			"transit_router_attachment_description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"transit_router_attachment_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"transit_router_attachment_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"transit_router_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"vbr_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"vbr_owner_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"resource_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      "VBR",
				ValidateFunc: StringInSlice([]string{"VBR"}, false),
			},
			"dry_run": {
				Type:     schema.TypeBool,
				Optional: true,
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

	action := "CreateTransitRouterVbrAttachment"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken(action)

	if v, ok := d.GetOkExists("dry_run"); ok {
		request["DryRun"] = v
	}
	if v, ok := d.GetOk("resource_type"); ok {
		request["ResourceType"] = v
	}
	if v, ok := d.GetOkExists("route_table_association_enabled"); ok {
		request["RouteTableAssociationEnabled"] = v
	}

	if v, ok := d.GetOkExists("route_table_propagation_enabled"); ok {
		request["RouteTablePropagationEnabled"] = v
	}

	if v, ok := d.GetOkExists("vbr_owner_id"); ok {
		request["VbrOwnerId"] = v
	}
	if v, ok := d.GetOkExists("auto_publish_route_enabled"); ok {
		request["AutoPublishRouteEnabled"] = v
	}
	if v, ok := d.GetOk("cen_id"); ok {
		request["CenId"] = v
	}
	if v, ok := d.GetOk("tags"); ok {
		tagsMap := ConvertTags(v.(map[string]interface{}))
		request = expandTagsToMap(request, tagsMap)
	}

	if v, ok := d.GetOk("transit_router_id"); ok {
		request["TransitRouterId"] = v
	}
	request["VbrId"] = d.Get("vbr_id")
	if v, ok := d.GetOk("transit_router_attachment_description"); ok {
		request["TransitRouterAttachmentDescription"] = v
	}
	if v, ok := d.GetOk("transit_router_attachment_name"); ok {
		request["TransitRouterAttachmentName"] = v
	}
	wait := incrementalWait(5*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("Cbn", "2017-09-12", action, query, request, true)
		if err != nil {
			if IsExpectedErrors(err, []string{"Operation.Blocking", "IncorrectStatus.Status", "InstanceStatus.NotSupport", "Throttling.User"}) || NeedRetry(err) {
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

	cenServiceV2 := CenServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{"Attached"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, cenServiceV2.CenTransitRouterVbrAttachmentStateRefreshFunc(d.Id(), "Status", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAliCloudCenTransitRouterVbrAttachmentUpdate(d, meta)
}

func resourceAliCloudCenTransitRouterVbrAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cenServiceV2 := CenServiceV2{client}

	objectRaw, err := cenServiceV2.DescribeCenTransitRouterVbrAttachment(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_cen_transit_router_vbr_attachment DescribeCenTransitRouterVbrAttachment Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("auto_publish_route_enabled", objectRaw["AutoPublishRouteEnabled"])
	d.Set("cen_id", objectRaw["CenId"])
	d.Set("create_time", objectRaw["CreationTime"])
	d.Set("order_type", objectRaw["OrderType"])
	d.Set("region_id", objectRaw["VbrRegionId"])
	d.Set("status", objectRaw["Status"])
	d.Set("transit_router_attachment_description", objectRaw["TransitRouterAttachmentDescription"])
	d.Set("transit_router_attachment_name", objectRaw["TransitRouterAttachmentName"])
	d.Set("transit_router_id", objectRaw["TransitRouterId"])
	d.Set("vbr_id", objectRaw["VbrId"])
	d.Set("vbr_owner_id", objectRaw["VbrOwnerId"])
	d.Set("transit_router_attachment_id", objectRaw["TransitRouterAttachmentId"])
	d.Set("resource_type", objectRaw["ResourceType"])

	cbnService := CbnService{client}
	listTagResourcesObject, err := cbnService.ListTagResources(d.Id(), "TransitRouterVbrAttachment")
	if err != nil {
		return WrapError(err)
	}

	d.Set("tags", tagsToMap(listTagResourcesObject))

	return nil
}

func resourceAliCloudCenTransitRouterVbrAttachmentUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false

	var err error
	parts := strings.Split(d.Id(), ":")
	action := "UpdateTransitRouterVbrAttachmentAttribute"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["TransitRouterAttachmentId"] = parts[1]
	if v, ok := d.GetOkExists("dry_run"); ok {
		request["DryRun"] = v
	}

	request["ClientToken"] = buildClientToken(action)
	if !d.IsNewResource() && d.HasChange("auto_publish_route_enabled") {
		update = true
		request["AutoPublishRouteEnabled"] = d.Get("auto_publish_route_enabled")
	}

	if d.HasChange("order_type") {
		update = true
		request["OrderType"] = d.Get("order_type")
	}

	if !d.IsNewResource() && d.HasChange("transit_router_attachment_description") {
		update = true
		request["TransitRouterAttachmentDescription"] = d.Get("transit_router_attachment_description")
	}

	if !d.IsNewResource() && d.HasChange("transit_router_attachment_name") {
		update = true
		request["TransitRouterAttachmentName"] = d.Get("transit_router_attachment_name")
	}

	if update {
		wait := incrementalWait(5*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("Cbn", "2017-09-12", action, query, request, true)
			if err != nil {
				if IsExpectedErrors(err, []string{"Operation.Blocking", "IncorrectStatus.Status", "Throttling.User"}) || NeedRetry(err) {
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
		cenServiceV2 := CenServiceV2{client}
		stateConf := BuildStateConf([]string{}, []string{"Attached"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, cenServiceV2.CenTransitRouterVbrAttachmentStateRefreshFunc(d.Id(), "Status", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}

	if d.HasChange("tags") {
		cbnService := CbnService{client}
		if err := cbnService.SetResourceTags(d, "TransitRouterVbrAttachment"); err != nil {
			return WrapError(err)
		}
	}
	return resourceAliCloudCenTransitRouterVbrAttachmentRead(d, meta)
}

func resourceAliCloudCenTransitRouterVbrAttachmentDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	parts := strings.Split(d.Id(), ":")
	action := "DeleteTransitRouterVbrAttachment"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["TransitRouterAttachmentId"] = parts[1]
	if v, ok := d.GetOkExists("dry_run"); ok {
		request["DryRun"] = v
	}

	request["ClientToken"] = buildClientToken(action)

	wait := incrementalWait(5*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("Cbn", "2017-09-12", action, query, request, true)
		if err != nil {
			if IsExpectedErrors(err, []string{"Operation.Blocking", "IncorrectStatus.Status", "InstanceStatus.NotSupport", "Throttling.User"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		if NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	cenServiceV2 := CenServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{""}, d.Timeout(schema.TimeoutDelete), 5*time.Second, cenServiceV2.CenTransitRouterVbrAttachmentStateRefreshFunc(d.Id(), "Status", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return nil
}

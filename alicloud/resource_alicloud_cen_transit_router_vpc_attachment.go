// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
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

func resourceAliCloudCenTransitRouterVpcAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudCenTransitRouterVpcAttachmentCreate,
		Read:   resourceAliCloudCenTransitRouterVpcAttachmentRead,
		Update: resourceAliCloudCenTransitRouterVpcAttachmentUpdate,
		Delete: resourceAliCloudCenTransitRouterVpcAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"auto_publish_route_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"cen_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"dry_run": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"force_delete": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"payment_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"PayAsYouGo"}, false),
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
			"transit_router_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"transit_router_vpc_attachment_name": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"transit_router_attachment_name"},
			},
			"transit_router_vpc_attachment_options": {
				Type:     schema.TypeMap,
				Computed: true,
				Optional: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"vpc_owner_id": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"zone_mappings": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"zone_id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"vswitch_id": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"transit_router_attachment_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"resource_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Computed:     true,
				ValidateFunc: StringInSlice([]string{"VPC"}, false),
			},
			"transit_router_attachment_name": {
				Type:       schema.TypeString,
				Optional:   true,
				Computed:   true,
				Deprecated: "Field 'transit_router_attachment_name' has been deprecated since provider version 1.230.1. New field 'transit_router_vpc_attachment_name' instead.",
			},
			"route_table_association_enabled": {
				Type:       schema.TypeBool,
				Optional:   true,
				Deprecated: "Field `route_table_association_enabled` has been deprecated from provider version 1.192.0. Please use the resource `alicloud_cen_transit_router_route_table_association` instead.",
			},
			"route_table_propagation_enabled": {
				Type:       schema.TypeBool,
				Optional:   true,
				Deprecated: "Field `route_table_propagation_enabled` has been deprecated from provider version 1.192.0. Please use the resource `alicloud_cen_transit_router_route_table_propagation` instead.",
			},
		},
	}
}

func resourceAliCloudCenTransitRouterVpcAttachmentCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateTransitRouterVpcAttachment"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	query["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken(action)

	if v, ok := d.GetOk("cen_id"); ok {
		request["CenId"] = v
	}
	if v, ok := d.GetOk("zone_mappings"); ok {
		zoneMappingsMaps := make([]interface{}, 0)
		for _, dataLoop := range v.(*schema.Set).List() {
			dataLoopTmp := dataLoop.(map[string]interface{})
			dataLoopMap := make(map[string]interface{})
			dataLoopMap["ZoneId"] = dataLoopTmp["zone_id"]
			dataLoopMap["VSwitchId"] = dataLoopTmp["vswitch_id"]
			zoneMappingsMaps = append(zoneMappingsMaps, dataLoopMap)
		}
		request["ZoneMappings"] = zoneMappingsMaps
	}

	if v, ok := d.GetOk("resource_type"); ok {
		request["ResourceType"] = v
	}
	if v, ok := d.GetOk("transit_router_id"); ok {
		request["TransitRouterId"] = v
	}
	request["VpcId"] = d.Get("vpc_id")
	if v, ok := d.GetOk("transit_router_attachment_description"); ok {
		request["TransitRouterAttachmentDescription"] = v
	}
	if v, ok := d.GetOk("tags"); ok {
		tagsMap := ConvertTags(v.(map[string]interface{}))
		request = expandTagsToMap(request, tagsMap)
	}

	if v, ok := d.GetOk("transit_router_attachment_name"); ok || d.HasChange("transit_router_attachment_name") {
		request["TransitRouterAttachmentName"] = v
	}

	if v, ok := d.GetOk("transit_router_vpc_attachment_name"); ok {
		request["TransitRouterAttachmentName"] = v
	}
	if v, ok := d.GetOkExists("auto_publish_route_enabled"); ok {
		request["AutoPublishRouteEnabled"] = v
	}
	if v, ok := d.GetOk("transit_router_vpc_attachment_options"); ok {
		options, _ := convertMaptoJsonString(v.(map[string]interface{}))
		request["TransitRouterVPCAttachmentOptions"] = options
	}
	if v, ok := d.GetOk("payment_type"); ok {
		request["ChargeType"] = convertCenTransitRouterVpcAttachmentChargeTypeRequest(v.(string))
	}
	if v, ok := d.GetOk("vpc_owner_id"); ok {
		request["VpcOwnerId"] = v
	}
	if v, ok := d.GetOkExists("dry_run"); ok {
		request["DryRun"] = v
	}
	if v, ok := d.GetOkExists("route_table_association_enabled"); ok {
		request["RouteTableAssociationEnabled"] = v
	}

	if v, ok := d.GetOkExists("route_table_propagation_enabled"); ok {
		request["RouteTablePropagationEnabled"] = v
	}

	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("Cbn", "2017-09-12", action, query, request, true)
		if err != nil {
			if IsExpectedErrors(err, []string{"Operation.Blocking", "Throttling.User", "IncorrectStatus.Attachment", "IncorrectStatus.VpcOrVswitch", "InstanceStatus.NotSupport", "IncorrectStatus.Status", "IncorrectStatus.VpcResource", "IncorrectStatus.VpcRouteTable", "IncorrectStatus.VpcSwitch"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_cen_transit_router_vpc_attachment", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%v:%v", request["CenId"], response["TransitRouterAttachmentId"]))

	cenServiceV2 := CenServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{"Attached"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, cenServiceV2.CenTransitRouterVpcAttachmentStateRefreshFunc(d.Id(), "Status", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAliCloudCenTransitRouterVpcAttachmentUpdate(d, meta)
}

func resourceAliCloudCenTransitRouterVpcAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cenServiceV2 := CenServiceV2{client}

	objectRaw, err := cenServiceV2.DescribeCenTransitRouterVpcAttachment(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_cen_transit_router_vpc_attachment DescribeCenTransitRouterVpcAttachment Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	if objectRaw["AutoPublishRouteEnabled"] != nil {
		d.Set("auto_publish_route_enabled", objectRaw["AutoPublishRouteEnabled"])
	}
	if objectRaw["CreationTime"] != nil {
		d.Set("create_time", objectRaw["CreationTime"])
	}
	if convertCenTransitRouterVpcAttachmentTransitRouterAttachmentsChargeTypeResponse(objectRaw["ChargeType"]) != nil {
		d.Set("payment_type", convertCenTransitRouterVpcAttachmentTransitRouterAttachmentsChargeTypeResponse(objectRaw["ChargeType"]))
	}
	if objectRaw["Status"] != nil {
		d.Set("status", objectRaw["Status"])
	}
	if objectRaw["TransitRouterAttachmentDescription"] != nil {
		d.Set("transit_router_attachment_description", objectRaw["TransitRouterAttachmentDescription"])
	}
	if objectRaw["TransitRouterId"] != nil {
		d.Set("transit_router_id", objectRaw["TransitRouterId"])
	}
	if objectRaw["TransitRouterAttachmentName"] != nil {
		d.Set("transit_router_vpc_attachment_name", objectRaw["TransitRouterAttachmentName"])
	}
	if objectRaw["TransitRouterVPCAttachmentOptions"] != nil {
		d.Set("transit_router_vpc_attachment_options", objectRaw["TransitRouterVPCAttachmentOptions"])
	}
	if objectRaw["VpcId"] != nil {
		d.Set("vpc_id", objectRaw["VpcId"])
	}
	if objectRaw["VpcOwnerId"] != nil {
		d.Set("vpc_owner_id", objectRaw["VpcOwnerId"])
	}
	d.Set("transit_router_attachment_id", objectRaw["TransitRouterAttachmentId"])
	d.Set("resource_type", objectRaw["ResourceType"])
	d.Set("cen_id", objectRaw["CenId"])

	tagsMaps := objectRaw["Tags"]
	d.Set("tags", tagsToMap(tagsMaps))
	zoneMappings1Raw := objectRaw["ZoneMappings"]
	zoneMappingsMaps := make([]map[string]interface{}, 0)
	if zoneMappings1Raw != nil {
		for _, zoneMappingsChild1Raw := range zoneMappings1Raw.([]interface{}) {
			zoneMappingsMap := make(map[string]interface{})
			zoneMappingsChild1Raw := zoneMappingsChild1Raw.(map[string]interface{})
			zoneMappingsMap["vswitch_id"] = zoneMappingsChild1Raw["VSwitchId"]
			zoneMappingsMap["zone_id"] = zoneMappingsChild1Raw["ZoneId"]

			zoneMappingsMaps = append(zoneMappingsMaps, zoneMappingsMap)
		}
	}
	if objectRaw["ZoneMappings"] != nil {
		if err := d.Set("zone_mappings", zoneMappingsMaps); err != nil {
			return err
		}
	}

	d.Set("transit_router_attachment_name", d.Get("transit_router_vpc_attachment_name"))
	return nil
}

func resourceAliCloudCenTransitRouterVpcAttachmentUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false
	action := "UpdateTransitRouterVpcAttachmentAttribute"
	var err error
	parts, _ := ParseResourceId(d.Id(), 2)
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	query["TransitRouterAttachmentId"] = parts[1]

	request["ClientToken"] = buildClientToken(action)
	if !d.IsNewResource() && d.HasChange("transit_router_attachment_description") {
		update = true
		request["TransitRouterAttachmentDescription"] = d.Get("transit_router_attachment_description")
	}

	if !d.IsNewResource() && d.HasChange("transit_router_attachment_name") {
		update = true
		request["TransitRouterAttachmentName"] = d.Get("transit_router_attachment_name")
	}

	if !d.IsNewResource() && d.HasChange("transit_router_vpc_attachment_name") {
		update = true
		request["TransitRouterAttachmentName"] = d.Get("transit_router_vpc_attachment_name")
	}

	if !d.IsNewResource() && d.HasChange("auto_publish_route_enabled") {
		update = true
		request["AutoPublishRouteEnabled"] = d.Get("auto_publish_route_enabled")
	}

	if !d.IsNewResource() && d.HasChange("transit_router_vpc_attachment_options") {
		update = true
		options, _ := convertMaptoJsonString(d.Get("transit_router_vpc_attachment_options").(map[string]interface{}))
		request["TransitRouterVPCAttachmentOptions"] = options
	}

	if v, ok := d.GetOkExists("dry_run"); ok {
		request["DryRun"] = v
	}
	if update {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("Cbn", "2017-09-12", action, query, request, true)
			if err != nil {
				if IsExpectedErrors(err, []string{"Operation.Blocking", "Throttling.User", "IncorrectStatus.Status", "IncorrectStatus.Vpc", "IncorrectStatus.VpcOrVswitch", "IncorrectStatus.Attachment"}) || NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(action, response, request)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		cenServiceV2 := CenServiceV2{client}
		stateConf := BuildStateConf([]string{}, []string{"Attached"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, cenServiceV2.CenTransitRouterVpcAttachmentStateRefreshFunc(d.Id(), "Status", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}

	if !d.IsNewResource() && d.HasChange("zone_mappings") {
		oldEntry, newEntry := d.GetChange("zone_mappings")
		oldEntrySet := oldEntry.(*schema.Set)
		newEntrySet := newEntry.(*schema.Set)
		removed := oldEntrySet.Difference(newEntrySet)
		added := newEntrySet.Difference(oldEntrySet)
		request = make(map[string]interface{})
		query = make(map[string]interface{})
		query["TransitRouterAttachmentId"] = parts[1]
		request["ClientToken"] = buildClientToken(action)
		action := "UpdateTransitRouterVpcAttachmentZones"
		if v, ok := d.GetOkExists("dry_run"); ok {
			request["DryRun"] = v
		}

		if removed.Len() > 0 {
			localData := removed.List()
			removeZoneMappingsMaps := make([]interface{}, 0)
			for _, dataLoop := range localData {
				dataLoopTmp := dataLoop.(map[string]interface{})
				dataLoopMap := make(map[string]interface{})
				dataLoopMap["VSwitchId"] = dataLoopTmp["vswitch_id"]
				dataLoopMap["ZoneId"] = dataLoopTmp["zone_id"]
				removeZoneMappingsMaps = append(removeZoneMappingsMaps, dataLoopMap)
			}
			request["RemoveZoneMappings"] = removeZoneMappingsMaps
		}

		if added.Len() > 0 {
			localData := added.List()
			addZoneMappingsMaps := make([]interface{}, 0)
			for _, dataLoop := range localData {
				dataLoopTmp := dataLoop.(map[string]interface{})
				dataLoopMap := make(map[string]interface{})
				dataLoopMap["VSwitchId"] = dataLoopTmp["vswitch_id"]
				dataLoopMap["ZoneId"] = dataLoopTmp["zone_id"]
				addZoneMappingsMaps = append(addZoneMappingsMaps, dataLoopMap)
			}
			request["AddZoneMappings"] = addZoneMappingsMaps
		}

		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("Cbn", "2017-09-12", action, query, request, true)
			if err != nil {
				if IsExpectedErrors(err, []string{"Operation.Blocking", "IncorrectStatus.Status", "IncorrectStatus.Vpc", "IncorrectStatus.VpcOrVswitch", "IncorrectStatus.Attachment"}) || NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(action, response, request)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		cenServiceV2 := CenServiceV2{client}
		stateConf := BuildStateConf([]string{}, []string{"Attached"}, d.Timeout(schema.TimeoutUpdate), 15*time.Second, cenServiceV2.CenTransitRouterVpcAttachmentStateRefreshFunc(d.Id(), "Status", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}

	}
	if d.HasChange("tags") {
		cbnService := CbnService{client}
		if err := cbnService.SetResourceTags(d, "TransitRouterVpcAttachment"); err != nil {
			return WrapError(err)
		}
	}
	return resourceAliCloudCenTransitRouterVpcAttachmentRead(d, meta)
}

func resourceAliCloudCenTransitRouterVpcAttachmentDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	action := "DeleteTransitRouterVpcAttachment"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	parts, _ := ParseResourceId(d.Id(), 2)
	request = make(map[string]interface{})
	query["TransitRouterAttachmentId"] = parts[1]

	request["ClientToken"] = buildClientToken(action)

	if v, ok := d.GetOkExists("dry_run"); ok {
		request["DryRun"] = v
	}
	if v, ok := d.GetOkExists("force_delete"); ok {
		request["Force"] = v
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("Cbn", "2017-09-12", action, query, request, true)
		request["ClientToken"] = buildClientToken(action)

		if err != nil {
			if IsExpectedErrors(err, []string{"Operation.Blocking", "Throttling.User", "InstanceStatus.NotSupport", "IncorrectStatus.VpcOrVswitch", "IncorrectStatus.Attachment", "IncorrectStatus.Status", "IncorrectStatus.VpcRouteEntry", "IncorrectStatus.Vpc", "IncorrectStatus.VpcRouteTable", "TokenProcessing", "IncorrectStatus.VpcSwitch"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	cenServiceV2 := CenServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{""}, d.Timeout(schema.TimeoutDelete), 5*time.Second, cenServiceV2.CenTransitRouterVpcAttachmentStateRefreshFunc(d.Id(), "Status", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return nil
}

func convertCenTransitRouterVpcAttachmentPaymentTypeResponse(source string) string {
	switch source {
	case "POSTPAY":
		return "PayAsYouGo"
	}
	return source
}

func convertCenTransitRouterVpcAttachmentPaymentTypeRequest(source string) string {
	switch source {
	case "PayAsYouGo":
		return "POSTPAY"
	}
	return source
}

func convertCenTransitRouterVpcAttachmentTransitRouterAttachmentsChargeTypeResponse(source interface{}) interface{} {
	switch source {
	case "POSTPAY":
		return "PayAsYouGo"
	}
	return source
}
func convertCenTransitRouterVpcAttachmentChargeTypeRequest(source interface{}) interface{} {
	switch source {
	case "PayAsYouGo":
		return "POSTPAY"
	}
	return source
}

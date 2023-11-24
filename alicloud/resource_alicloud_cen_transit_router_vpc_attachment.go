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
			Create: schema.DefaultTimeout(3 * time.Minute),
			Update: schema.DefaultTimeout(3 * time.Minute),
			Delete: schema.DefaultTimeout(3 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"cen_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"vpc_id": {
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
				Computed:     true,
				ValidateFunc: StringInSlice([]string{"VPC"}, false),
			},
			"payment_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Computed:     true,
				ValidateFunc: StringInSlice([]string{"PayAsYouGo"}, false),
			},
			"vpc_owner_id": {
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
			"dry_run": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"tags": tagsSchema(),
			"zone_mappings": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"vswitch_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"zone_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
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
	cbnService := CbnService{client}
	var response map[string]interface{}
	action := "CreateTransitRouterVpcAttachment"
	request := make(map[string]interface{})
	conn, err := client.NewCbnClient()
	if err != nil {
		return WrapError(err)
	}

	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken("CreateTransitRouterVpcAttachment")
	request["CenId"] = d.Get("cen_id")
	request["VpcId"] = d.Get("vpc_id")

	if v, ok := d.GetOk("transit_router_id"); ok {
		request["TransitRouterId"] = v
	}

	if v, ok := d.GetOk("resource_type"); ok {
		request["ResourceType"] = v
	}

	if v, ok := d.GetOk("payment_type"); ok {
		request["ChargeType"] = convertCenTransitRouterVpcAttachmentPaymentTypeRequest(v.(string))
	}

	if v, ok := d.GetOk("vpc_owner_id"); ok {
		request["VpcOwnerId"] = v
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

	if v, ok := d.GetOkExists("route_table_association_enabled"); ok {
		request["RouteTableAssociationEnabled"] = v
	}

	if v, ok := d.GetOkExists("route_table_propagation_enabled"); ok {
		request["RouteTablePropagationEnabled"] = v
	}

	if v, ok := d.GetOkExists("dry_run"); ok {
		request["DryRun"] = v
	}

	zoneMappings := d.Get("zone_mappings")
	zoneMappingsMaps := make([]map[string]interface{}, 0)
	for _, zoneMappingsList := range zoneMappings.(*schema.Set).List() {
		zoneMappingsMap := make(map[string]interface{})
		zoneMappingsArg := zoneMappingsList.(map[string]interface{})

		if vSwitchId, ok := zoneMappingsArg["vswitch_id"]; ok {
			zoneMappingsMap["VSwitchId"] = vSwitchId
		}

		if zoneId, ok := zoneMappingsArg["zone_id"]; ok {
			zoneMappingsMap["ZoneId"] = zoneId
		}

		zoneMappingsMaps = append(zoneMappingsMaps, zoneMappingsMap)
	}

	request["ZoneMappings"] = zoneMappingsMaps

	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-09-12"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			if IsExpectedErrors(err, []string{"Operation.Blocking", "InstanceStatus.NotSupport", "IncorrectStatus.Status", "IncorrectStatus.VpcOrVswitch", "IncorrectStatus.Attachment", "IncorrectStatus.VpcResource", "IncorrectStatus.VpcRouteTable", "IncorrectStatus.VpcSwitch"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_cen_transit_router_vpc_attachment", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%v:%v", request["CenId"], response["TransitRouterAttachmentId"]))

	stateConf := BuildStateConf([]string{}, []string{"Attached"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, cbnService.CenTransitRouterVpcAttachmentStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAliCloudCenTransitRouterVpcAttachmentUpdate(d, meta)
}

func resourceAliCloudCenTransitRouterVpcAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cbnService := CbnService{client}

	object, err := cbnService.DescribeCenTransitRouterVpcAttachment(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_cen_transit_router_vpc_attachment cbnService.DescribeCenTransitRouterVpcAttachment Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("cen_id", object["CenId"])
	d.Set("vpc_id", object["VpcId"])
	d.Set("transit_router_id", object["TransitRouterId"])
	d.Set("resource_type", object["ResourceType"])
	d.Set("payment_type", convertCenTransitRouterVpcAttachmentPaymentTypeResponse(object["ChargeType"].(string)))
	d.Set("vpc_owner_id", fmt.Sprint(object["VpcOwnerId"]))
	d.Set("auto_publish_route_enabled", object["AutoPublishRouteEnabled"])
	d.Set("transit_router_attachment_name", object["TransitRouterAttachmentName"])
	d.Set("transit_router_attachment_description", object["TransitRouterAttachmentDescription"])
	d.Set("transit_router_attachment_id", object["TransitRouterAttachmentId"])
	d.Set("status", object["Status"])

	if zoneMappings, ok := object["ZoneMappings"]; ok {
		zoneMappingsMaps := make([]map[string]interface{}, 0)
		for _, zoneMappingsList := range zoneMappings.([]interface{}) {
			zoneMappingsArg := zoneMappingsList.(map[string]interface{})
			zoneMappingsMap := map[string]interface{}{}

			if vSwitchId, ok := zoneMappingsArg["VSwitchId"]; ok {
				zoneMappingsMap["vswitch_id"] = vSwitchId
			}

			if zoneId, ok := zoneMappingsArg["ZoneId"]; ok {
				zoneMappingsMap["zone_id"] = zoneId
			}

			zoneMappingsMaps = append(zoneMappingsMaps, zoneMappingsMap)
		}

		d.Set("zone_mappings", zoneMappingsMaps)
	}

	listTagResourcesObject, err := cbnService.ListTagResources(d.Id(), "TransitRouterVpcAttachment")
	if err != nil {
		return WrapError(err)
	}

	d.Set("tags", tagsToMap(listTagResourcesObject))

	return nil
}

func resourceAliCloudCenTransitRouterVpcAttachmentUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cbnService := CbnService{client}
	var response map[string]interface{}
	d.Partial(true)

	if d.HasChange("tags") {
		if err := cbnService.SetResourceTags(d, "TransitRouterVpcAttachment"); err != nil {
			return WrapError(err)
		}
		d.SetPartial("tags")
	}

	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}

	update := false
	request := map[string]interface{}{
		"TransitRouterAttachmentId": parts[1],
	}

	if !d.IsNewResource() && d.HasChange("auto_publish_route_enabled") {
		update = true
	}
	if v, ok := d.GetOkExists("auto_publish_route_enabled"); ok {
		request["AutoPublishRouteEnabled"] = v
	}

	if !d.IsNewResource() && d.HasChange("transit_router_attachment_name") {
		update = true
	}
	if v, ok := d.GetOk("transit_router_attachment_name"); ok {
		request["TransitRouterAttachmentName"] = v
	}

	if !d.IsNewResource() && d.HasChange("transit_router_attachment_description") {
		update = true
	}
	if v, ok := d.GetOk("transit_router_attachment_description"); ok {
		request["TransitRouterAttachmentDescription"] = v
	}

	if update {
		if _, ok := d.GetOkExists("dry_run"); ok {
			request["DryRun"] = d.Get("dry_run")
		}

		action := "UpdateTransitRouterVpcAttachmentAttribute"
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
				if IsExpectedErrors(err, []string{"Operation.Blocking", "IncorrectStatus.Status", "IncorrectStatus.Vpc", "IncorrectStatus.VpcOrVswitch", "IncorrectStatus.Attachment"}) || NeedRetry(err) {
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

		stateConf := BuildStateConf([]string{}, []string{"Attached"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, cbnService.CenTransitRouterVpcAttachmentStateRefreshFunc(d.Id(), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}

		d.SetPartial("auto_publish_route_enabled")
		d.SetPartial("transit_router_attachment_name")
		d.SetPartial("transit_router_attachment_description")
	}

	if !d.IsNewResource() && d.HasChange("zone_mappings") {
		oraw, nraw := d.GetChange("zone_mappings")
		remove := oraw.(*schema.Set).Difference(nraw.(*schema.Set)).List()
		create := nraw.(*schema.Set).Difference(oraw.(*schema.Set)).List()

		updateZonesRequest := map[string]interface{}{
			"TransitRouterAttachmentId": parts[1],
		}

		action := "UpdateTransitRouterVpcAttachmentZones"
		conn, err := client.NewCbnClient()
		if err != nil {
			return WrapError(err)
		}

		if len(remove) > 0 {
			zoneMappingsMaps := make([]map[string]interface{}, 0)
			for _, zoneMappingsList := range remove {
				zoneMappingsMap := make(map[string]interface{})
				zoneMappingsArg := zoneMappingsList.(map[string]interface{})

				if vSwitchId, ok := zoneMappingsArg["vswitch_id"]; ok {
					zoneMappingsMap["VSwitchId"] = vSwitchId
				}

				if zoneId, ok := zoneMappingsArg["zone_id"]; ok {
					zoneMappingsMap["ZoneId"] = zoneId
				}

				zoneMappingsMaps = append(zoneMappingsMaps, zoneMappingsMap)
			}

			updateZonesRequest["RemoveZoneMappings"] = zoneMappingsMaps
		}

		if len(create) > 0 {
			zoneMappingsMaps := make([]map[string]interface{}, 0)
			for _, zoneMappingsList := range create {
				zoneMappingsMap := make(map[string]interface{})
				zoneMappingsArg := zoneMappingsList.(map[string]interface{})

				if vSwitchId, ok := zoneMappingsArg["vswitch_id"]; ok {
					zoneMappingsMap["VSwitchId"] = vSwitchId
				}

				if zoneId, ok := zoneMappingsArg["zone_id"]; ok {
					zoneMappingsMap["ZoneId"] = zoneId
				}

				zoneMappingsMaps = append(zoneMappingsMaps, zoneMappingsMap)
			}

			updateZonesRequest["AddZoneMappings"] = zoneMappingsMaps
		}

		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-09-12"), StringPointer("AK"), nil, updateZonesRequest, &runtime)
			if err != nil {
				if IsExpectedErrors(err, []string{"Operation.Blocking", "IncorrectStatus.Status", "IncorrectStatus.Vpc", "IncorrectStatus.VpcOrVswitch", "IncorrectStatus.Attachment"}) || NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, updateZonesRequest)

		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}

		stateConf := BuildStateConf([]string{}, []string{"Attached"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, cbnService.CenTransitRouterVpcAttachmentStateRefreshFunc(d.Id(), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}

		d.SetPartial("zone_mappings")
	}

	d.Partial(false)

	return resourceAliCloudCenTransitRouterVpcAttachmentRead(d, meta)
}

func resourceAliCloudCenTransitRouterVpcAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cbnService := CbnService{client}
	action := "DeleteTransitRouterVpcAttachment"
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
			if IsExpectedErrors(err, []string{"Operation.Blocking", "InstanceStatus.NotSupport", "IncorrectStatus.Status", "IncorrectStatus.VpcOrVswitch", "IncorrectStatus.Vpc", "IncorrectStatus.VpcRouteEntry", "IncorrectStatus.Attachment", "IncorrectStatus.VpcRouteTable", "IncorrectStatus.VpcSwitch", "TokenProcessing"}) || NeedRetry(err) {
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

	stateConf := BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 5*time.Second, cbnService.CenTransitRouterVpcAttachmentStateRefreshFunc(d.Id(), []string{}))
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

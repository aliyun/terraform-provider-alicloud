package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudCenTransitRouterVpcAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudCenTransitRouterVpcAttachmentCreate,
		Read:   resourceAlicloudCenTransitRouterVpcAttachmentRead,
		Update: resourceAlicloudCenTransitRouterVpcAttachmentUpdate,
		Delete: resourceAlicloudCenTransitRouterVpcAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(3 * time.Minute),
			Delete: schema.DefaultTimeout(3 * time.Minute),
			Update: schema.DefaultTimeout(3 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"cen_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"dry_run": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"resource_type": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "VPC",
			},
			"route_table_association_enabled": {
				Type:       schema.TypeBool,
				Optional:   true,
				Deprecated: "Field 'route_table_association_enabled' has been deprecated from provider version 1.192.0. Please use the resource 'alicloud_cen_transit_router_route_table_association' instead.",
			},
			"route_table_propagation_enabled": {
				Type:       schema.TypeBool,
				Optional:   true,
				Deprecated: "Field 'route_table_propagation_enabled' has been deprecated from provider version 1.192.0. Please use the resource 'alicloud_cen_transit_router_route_table_propagation' instead.",
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"transit_router_attachment_description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"transit_router_attachment_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"transit_router_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
			},
			"transit_router_attachment_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"vpc_owner_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"payment_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"PayAsYouGo"}, false),
				ForceNew:     true,
				Computed:     true,
			},
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
		},
	}
}

func resourceAlicloudCenTransitRouterVpcAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cbnService := CbnService{client}
	var response map[string]interface{}
	action := "CreateTransitRouterVpcAttachment"
	request := make(map[string]interface{})
	conn, err := client.NewCbnClient()
	if err != nil {
		return WrapError(err)
	}
	request["CenId"] = d.Get("cen_id")

	if v, ok := d.GetOkExists("dry_run"); ok {
		request["DryRun"] = v
	}

	request["RegionId"] = client.RegionId
	if v, ok := d.GetOk("resource_type"); ok {
		request["ResourceType"] = v
	}

	if v, ok := d.GetOkExists("route_table_association_enabled"); ok {
		request["RouteTableAssociationEnabled"] = v
	}

	if v, ok := d.GetOkExists("route_table_propagation_enabled"); ok {
		request["RouteTablePropagationEnabled"] = v
	}

	if v, ok := d.GetOk("transit_router_attachment_description"); ok {
		request["TransitRouterAttachmentDescription"] = v
	}

	if v, ok := d.GetOk("transit_router_attachment_name"); ok {
		request["TransitRouterAttachmentName"] = v
	}

	if v, ok := d.GetOk("payment_type"); ok {
		request["ChargeType"] = convertCenTransitRouterVpcAttachmentPaymentTypeRequest(v.(string))
	}

	if v, ok := d.GetOk("transit_router_id"); ok {
		request["TransitRouterId"] = v
	}

	request["VpcId"] = d.Get("vpc_id")
	if v, ok := d.GetOk("vpc_owner_id"); ok {
		request["VpcOwnerId"] = v
	}

	zoneMappingsMaps := make([]map[string]interface{}, 0)
	for _, zoneMappings := range d.Get("zone_mappings").(*schema.Set).List() {
		zoneMappingsMap := make(map[string]interface{})
		zoneMappingsArg := zoneMappings.(map[string]interface{})
		zoneMappingsMap["VSwitchId"] = zoneMappingsArg["vswitch_id"]
		zoneMappingsMap["ZoneId"] = zoneMappingsArg["zone_id"]
		zoneMappingsMaps = append(zoneMappingsMaps, zoneMappingsMap)
	}
	request["ZoneMappings"] = zoneMappingsMaps

	request["ClientToken"] = buildClientToken("CreateTransitRouterVpcAttachment")
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-09-12"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			if IsExpectedErrors(err, []string{"Operation.Blocking", "InstanceStatus.NotSupport", "IncorrectStatus.Status", "IncorrectStatus.VpcOrVswitch", "IncorrectStatus.Attachment"}) || NeedRetry(err) {
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

	return resourceAlicloudCenTransitRouterVpcAttachmentRead(d, meta)
}

func resourceAlicloudCenTransitRouterVpcAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cbnService := CbnService{client}
	object, err := cbnService.DescribeCenTransitRouterVpcAttachment(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_cen_transit_router_vpc_attachment cbnService.DescribeCenTransitRouterVpcAttachment Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	parts, err1 := ParseResourceId(d.Id(), 2)
	if err1 != nil {
		return WrapError(err1)
	}
	d.Set("cen_id", parts[0])
	d.Set("resource_type", object["ResourceType"])
	d.Set("status", object["Status"])
	d.Set("transit_router_attachment_description", object["TransitRouterAttachmentDescription"])
	d.Set("transit_router_attachment_name", object["TransitRouterAttachmentName"])
	d.Set("transit_router_attachment_id", object["TransitRouterAttachmentId"])
	d.Set("vpc_id", object["VpcId"])
	d.Set("payment_type", convertCenTransitRouterVpcAttachmentPaymentTypeResponse(object["ChargeType"].(string)))
	d.Set("vpc_owner_id", fmt.Sprint(object["VpcOwnerId"]))

	zoneMappings := make([]map[string]interface{}, 0)
	if zoneMappingsList, ok := object["ZoneMappings"].([]interface{}); ok {
		for _, v := range zoneMappingsList {
			if m1, ok := v.(map[string]interface{}); ok {
				temp1 := map[string]interface{}{
					"vswitch_id": m1["VSwitchId"],
					"zone_id":    m1["ZoneId"],
				}
				zoneMappings = append(zoneMappings, temp1)

			}
		}
	}
	if err := d.Set("zone_mappings", zoneMappings); err != nil {
		return WrapError(err)
	}
	return nil
}

func resourceAlicloudCenTransitRouterVpcAttachmentUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cbnService := CbnService{client}
	conn, err := client.NewCbnClient()
	if err != nil {
		return WrapError(err)
	}
	var response map[string]interface{}
	d.Partial(true)
	update := false
	parts, err1 := ParseResourceId(d.Id(), 2)
	if err1 != nil {
		return WrapError(err1)
	}
	request := map[string]interface{}{
		"TransitRouterAttachmentId": parts[1],
	}
	if d.HasChange("resource_type") {
		update = true
		request["ResourceType"] = d.Get("resource_type")
	}
	if d.HasChange("transit_router_attachment_description") {
		update = true
		request["TransitRouterAttachmentDescription"] = d.Get("transit_router_attachment_description")
	}
	if d.HasChange("transit_router_attachment_name") {
		update = true
		request["TransitRouterAttachmentName"] = d.Get("transit_router_attachment_name")
	}
	if update {
		if _, ok := d.GetOkExists("dry_run"); ok {
			request["DryRun"] = d.Get("dry_run")
		}
		action := "UpdateTransitRouterVpcAttachmentAttribute"
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-09-12"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
		d.SetPartial("resource_type")
		d.SetPartial("transit_router_attachment_description")
		d.SetPartial("transit_router_attachment_name")
	}

	if d.HasChange("zone_mappings") {
		oraw, nraw := d.GetChange("zone_mappings")
		remove := oraw.(*schema.Set).Difference(nraw.(*schema.Set)).List()
		create := nraw.(*schema.Set).Difference(oraw.(*schema.Set)).List()
		updateZonesRequest := map[string]interface{}{
			"TransitRouterAttachmentId": parts[1],
		}
		if len(remove) > 0 {
			zoneMappingsMaps := make([]map[string]interface{}, 0)
			for _, zoneMappings := range remove {
				zoneMappingsMap := make(map[string]interface{})
				zoneMappingsArg := zoneMappings.(map[string]interface{})
				zoneMappingsMap["VSwitchId"] = zoneMappingsArg["vswitch_id"]
				zoneMappingsMap["ZoneId"] = zoneMappingsArg["zone_id"]
				zoneMappingsMaps = append(zoneMappingsMaps, zoneMappingsMap)
			}
			updateZonesRequest["RemoveZoneMappings"] = zoneMappingsMaps
		}

		if len(create) > 0 {
			zoneMappingsMaps := make([]map[string]interface{}, 0)
			for _, zoneMappings := range create {
				zoneMappingsMap := make(map[string]interface{})
				zoneMappingsArg := zoneMappings.(map[string]interface{})
				zoneMappingsMap["VSwitchId"] = zoneMappingsArg["vswitch_id"]
				zoneMappingsMap["ZoneId"] = zoneMappingsArg["zone_id"]
				zoneMappingsMaps = append(zoneMappingsMaps, zoneMappingsMap)
			}
			updateZonesRequest["AddZoneMappings"] = zoneMappingsMaps
		}

		action := "UpdateTransitRouterVpcAttachmentZones"
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-09-12"), StringPointer("AK"), nil, updateZonesRequest, &util.RuntimeOptions{})
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
	return resourceAlicloudCenTransitRouterVpcAttachmentRead(d, meta)
}

func resourceAlicloudCenTransitRouterVpcAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cbnService := CbnService{client}
	action := "DeleteTransitRouterVpcAttachment"
	var response map[string]interface{}
	conn, err := client.NewCbnClient()
	if err != nil {
		return WrapError(err)
	}
	parts, err1 := ParseResourceId(d.Id(), 2)
	if err1 != nil {
		return WrapError(err1)
	}
	request := map[string]interface{}{
		"TransitRouterAttachmentId": parts[1],
	}

	if v, ok := d.GetOkExists("dry_run"); ok {
		request["DryRun"] = v
	}
	if v, ok := d.GetOk("resource_type"); ok {
		request["ResourceType"] = v
	}
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-09-12"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
		if err != nil {
			if IsExpectedErrors(err, []string{"Operation.Blocking", "InstanceStatus.NotSupport", "IncorrectStatus.Status", "IncorrectStatus.VpcOrVswitch", "IncorrectStatus.Attachment"}) || NeedRetry(err) {
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

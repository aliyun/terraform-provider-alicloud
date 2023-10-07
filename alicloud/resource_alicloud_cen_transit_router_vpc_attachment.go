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
			"charge_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"resource_type": {
				Type:     schema.TypeString,
				Optional: true,
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
				ForceNew: true,
			},
			"transit_router_vpc_attachment_name": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"transit_router_attachment_name"},
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"zone_mappings": {
				Type:     schema.TypeList,
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
			"transit_router_attachment_name": {
				Type:       schema.TypeString,
				Optional:   true,
				Computed:   true,
				Deprecated: "Field 'transit_router_attachment_name' has been deprecated since provider version 1.212.0. New field 'transit_router_vpc_attachment_name' instead.",
			},
		},
	}
}

func resourceAliCloudCenTransitRouterVpcAttachmentCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateTransitRouterVpcAttachment"
	var request map[string]interface{}
	var response map[string]interface{}
	conn, err := client.NewCenClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken(action)

	if v, ok := d.GetOk("cen_id"); ok {
		request["CenId"] = v
	}
	if v, ok := d.GetOk("zone_mappings"); ok {
		zoneMappingsMaps := make([]map[string]interface{}, 0)
		for _, dataLoop := range v.([]interface{}) {
			dataLoopTmp := dataLoop.(map[string]interface{})
			dataLoopMap := make(map[string]interface{})
			dataLoopMap["ZoneId"] = dataLoopTmp["zone_id"]
			dataLoopMap["VSwitchId"] = dataLoopTmp["vswitch_id"]
			zoneMappingsMaps = append(zoneMappingsMaps, dataLoopMap)
		}
		request["ZoneMappings"] = zoneMappingsMaps
	}

	if v, ok := d.GetOk("transit_router_id"); ok {
		request["TransitRouterId"] = v
	}
	request["VpcId"] = d.Get("vpc_id")
	if v, ok := d.GetOk("transit_router_attachment_description"); ok {
		request["TransitRouterAttachmentDescription"] = v
	}
	if v, ok := d.GetOk("charge_type"); ok {
		request["ChargeType"] = v
	}
	if v, ok := d.GetOk("transit_router_attachment_name"); ok {
		request["TransitRouterAttachmentName"] = v
	}

	if v, ok := d.GetOk("transit_router_vpc_attachment_name"); ok {
		request["TransitRouterAttachmentName"] = v
	}
	if v, ok := d.GetOkExists("auto_publish_route_enabled"); ok {
		request["AutoPublishRouteEnabled"] = v
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-09-12"), StringPointer("AK"), nil, request, &runtime)
		request["ClientToken"] = buildClientToken(action)

		if err != nil {
			if IsExpectedErrors(err, []string{"Operation.Blocking", "Throttling.User", "IncorrectStatus.Attachment", "IncorrectStatus.VpcOrVswitch", "InstanceStatus.NotSupport", "IncorrectStatus.Status", "IncorrectStatus.VpcResource", "IncorrectStatus.VpcRouteTable"}) || NeedRetry(err) {
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

	d.Set("auto_publish_route_enabled", objectRaw["AutoPublishRouteEnabled"])
	d.Set("charge_type", objectRaw["ChargeType"])
	d.Set("create_time", objectRaw["CreationTime"])
	d.Set("status", objectRaw["Status"])
	d.Set("transit_router_attachment_description", objectRaw["TransitRouterAttachmentDescription"])
	d.Set("transit_router_id", objectRaw["TransitRouterId"])
	d.Set("transit_router_vpc_attachment_name", objectRaw["TransitRouterAttachmentName"])
	d.Set("vpc_id", objectRaw["VpcId"])
	d.Set("resource_type", objectRaw["ResourceType"])
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
	d.Set("zone_mappings", zoneMappingsMaps)

	d.Set("transit_router_attachment_name", d.Get("transit_router_vpc_attachment_name"))
	return nil
}

func resourceAliCloudCenTransitRouterVpcAttachmentUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	update := false
	action := "UpdateTransitRouterVpcAttachmentAttribute"
	conn, err := client.NewCenClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	parts, _ := ParseResourceId(d.Id(), 2)
	request["TransitRouterAttachmentId"] = parts[1]
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

	if update {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-09-12"), StringPointer("AK"), nil, request, &runtime)
			request["ClientToken"] = buildClientToken(action)

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
		removed := oldEntry
		added := newEntry

		if len(removed.([]interface{})) > 0 {
			action := "UpdateTransitRouterVpcAttachmentZones"
			conn, err := client.NewCenClient()
			if err != nil {
				return WrapError(err)
			}
			request = make(map[string]interface{})
			parts, _ := ParseResourceId(d.Id(), 2)
			request["TransitRouterAttachmentId"] = parts[1]
			request["ClientToken"] = buildClientToken(action)
			localData := removed.([]interface{})
			removeZoneMappingsMaps := make([]map[string]interface{}, 0)
			for _, dataLoop := range localData {
				dataLoopTmp := dataLoop.(map[string]interface{})
				dataLoopMap := make(map[string]interface{})
				dataLoopMap["VSwitchId"] = dataLoopTmp["vswitch_id"]
				dataLoopMap["ZoneId"] = dataLoopTmp["zone_id"]
				removeZoneMappingsMaps = append(removeZoneMappingsMaps, dataLoopMap)
			}
			request["RemoveZoneMappings"] = removeZoneMappingsMaps

			runtime := util.RuntimeOptions{}
			runtime.SetAutoretry(true)
			wait := incrementalWait(3*time.Second, 5*time.Second)
			err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
				response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-09-12"), StringPointer("AK"), nil, request, &runtime)
				request["ClientToken"] = buildClientToken(action)

				if err != nil {
					if NeedRetry(err) {
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

		if len(added.([]interface{})) > 0 {
			action := "UpdateTransitRouterVpcAttachmentZones"
			conn, err := client.NewCenClient()
			if err != nil {
				return WrapError(err)
			}
			request = make(map[string]interface{})
			parts, _ := ParseResourceId(d.Id(), 2)
			request["TransitRouterAttachmentId"] = parts[1]
			request["ClientToken"] = buildClientToken(action)
			localData := added.([]interface{})
			addZoneMappingsMaps := make([]map[string]interface{}, 0)
			for _, dataLoop := range localData {
				dataLoopTmp := dataLoop.(map[string]interface{})
				dataLoopMap := make(map[string]interface{})
				dataLoopMap["VSwitchId"] = dataLoopTmp["vswitch_id"]
				dataLoopMap["ZoneId"] = dataLoopTmp["zone_id"]
				addZoneMappingsMaps = append(addZoneMappingsMaps, dataLoopMap)
			}
			request["AddZoneMappings"] = addZoneMappingsMaps

			runtime := util.RuntimeOptions{}
			runtime.SetAutoretry(true)
			wait := incrementalWait(3*time.Second, 5*time.Second)
			err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
				response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-09-12"), StringPointer("AK"), nil, request, &runtime)
				request["ClientToken"] = buildClientToken(action)

				if err != nil {
					if NeedRetry(err) {
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
	}
	if d.HasChange("tags") {
		cenServiceV2 := CenServiceV2{client}
		if err := cenServiceV2.SetResourceTags(d, ""); err != nil {
			return WrapError(err)
		}
		d.SetPartial("tags")
	}
	return resourceAliCloudCenTransitRouterVpcAttachmentRead(d, meta)
}

func resourceAliCloudCenTransitRouterVpcAttachmentDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	action := "DeleteTransitRouterVpcAttachment"
	var request map[string]interface{}
	var response map[string]interface{}
	conn, err := client.NewCenClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	parts, _ := ParseResourceId(d.Id(), 2)
	request["TransitRouterAttachmentId"] = parts[1]

	request["ClientToken"] = buildClientToken(action)

	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-09-12"), StringPointer("AK"), nil, request, &runtime)
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
		if NotFoundError(err) {
			return nil
		}
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

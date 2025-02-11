package alicloud

import (
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudCenTransitRouterVpnAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudCenTransitRouterVpnAttachmentCreate,
		Read:   resourceAlicloudCenTransitRouterVpnAttachmentRead,
		Update: resourceAlicloudCenTransitRouterVpnAttachmentUpdate,
		Delete: resourceAlicloudCenTransitRouterVpnAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(40 * time.Minute),
			Update: schema.DefaultTimeout(1 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"auto_publish_route_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"cen_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"transit_router_attachment_description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.All(validation.StringLenBetween(2, 256), validation.StringDoesNotMatch(regexp.MustCompile(`(^http://.*)|(^https://.*)`), "It must be `2` to `256` characters in length and cannot start with `https://` or `https://`.")),
			},
			"transit_router_attachment_name": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringMatch(regexp.MustCompile(`^[a-zA-Z][-_a-zA-Z0-9]{1,127}$`), "The name can be up to 128 characters in length and can contain digits, letters, hyphens (-), and underscores (_). It must start with a digit or letter."),
			},
			"transit_router_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"vpn_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"vpn_owner_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"zone": {
				Type:     schema.TypeSet,
				Required: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"zone_id": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"tags": tagsSchema(),
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAlicloudCenTransitRouterVpnAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateTransitRouterVpnAttachment"
	request := make(map[string]interface{})
	var err error
	if v, ok := d.GetOkExists("auto_publish_route_enabled"); ok {
		request["AutoPublishRouteEnabled"] = v
	}
	request["RegionId"] = client.RegionId
	if v, ok := d.GetOk("transit_router_attachment_description"); ok {
		request["TransitRouterAttachmentDescription"] = v
	}
	if v, ok := d.GetOk("transit_router_attachment_name"); ok {
		request["TransitRouterAttachmentName"] = v
	}
	if v, ok := d.GetOk("cen_id"); ok {
		request["CenId"] = v
	}
	request["TransitRouterId"] = d.Get("transit_router_id")
	request["VpnId"] = d.Get("vpn_id")
	if v, ok := d.GetOk("vpn_owner_id"); ok {
		request["VpnOwnerId"] = v
	}
	for zonePtr, zone := range d.Get("zone").(*schema.Set).List() {
		zoneArg := zone.(map[string]interface{})
		request["Zone."+fmt.Sprint(zonePtr+1)+".ZoneId"] = zoneArg["zone_id"]
	}
	request["ClientToken"] = buildClientToken("CreateTransitRouterVpnAttachment")
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *resource.RetryError {
		response, err = client.RpcPost("Cbn", "2017-09-12", action, nil, request, true)
		if err != nil {
			if IsExpectedErrors(err, []string{"Operation.Blocking", "OperationFailed.AllocateCidrFailed", "IncorrectStatus.Status"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_cen_transit_router_vpn_attachment", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["TransitRouterAttachmentId"]))
	cbnService := CbnService{client}
	stateConf := BuildStateConf([]string{}, []string{"Attached"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, cbnService.CenTransitRouterVpnAttachmentStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAlicloudCenTransitRouterVpnAttachmentUpdate(d, meta)
}

func resourceAlicloudCenTransitRouterVpnAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cbnService := CbnService{client}
	object, err := cbnService.DescribeCenTransitRouterVpnAttachment(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_cen_transit_router_vpn_attachment cbnService.DescribeCenTransitRouterVpnAttachment Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("auto_publish_route_enabled", object["AutoPublishRouteEnabled"])
	d.Set("status", object["Status"])
	d.Set("transit_router_id", object["TransitRouterId"])
	d.Set("transit_router_attachment_description", object["TransitRouterAttachmentDescription"])
	d.Set("transit_router_attachment_name", object["TransitRouterAttachmentName"])
	d.Set("vpn_id", object["VpnId"])
	d.Set("vpn_owner_id", fmt.Sprint(object["VpnOwnerId"]))

	if zonesList, ok := object["Zones"]; ok && zonesList != nil {
		zoneMaps := make([]map[string]interface{}, 0)
		for _, zonesListItem := range zonesList.([]interface{}) {
			if zonesListItemMap, ok := zonesListItem.(map[string]interface{}); ok {
				zonesListItemMap["zone_id"] = zonesListItemMap["ZoneId"]
				zoneMaps = append(zoneMaps, map[string]interface{}{
					"zone_id": zonesListItemMap["ZoneId"],
				})
			}
			d.Set("zone", zoneMaps)
		}
	}

	listTagResourcesObject, err := cbnService.ListTagResources(d.Id(), "TransitRouterVpnAttachment")
	if err != nil {
		return WrapError(err)
	}

	d.Set("tags", tagsToMap(listTagResourcesObject))

	return nil
}

func resourceAlicloudCenTransitRouterVpnAttachmentUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cbnService := CbnService{client}
	var err error
	var response map[string]interface{}
	update := false
	d.Partial(true)
	request := map[string]interface{}{
		"TransitRouterAttachmentId": d.Id(),
	}

	if d.HasChange("tags") {
		if err := cbnService.SetResourceTags(d, "TransitRouterVpnAttachment"); err != nil {
			return WrapError(err)
		}
		d.SetPartial("tags")
	}

	if !d.IsNewResource() && d.HasChange("auto_publish_route_enabled") {
		update = true
		if v, ok := d.GetOkExists("auto_publish_route_enabled"); ok {
			request["AutoPublishRouteEnabled"] = v
		}
	}
	if !d.IsNewResource() && d.HasChange("transit_router_attachment_description") {
		update = true
		if v, ok := d.GetOk("transit_router_attachment_description"); ok {
			request["TransitRouterAttachmentDescription"] = v
		}
	}
	if !d.IsNewResource() && d.HasChange("transit_router_attachment_name") {
		update = true
		if v, ok := d.GetOk("transit_router_attachment_name"); ok {
			request["TransitRouterAttachmentName"] = v
		}
	}
	if update {
		action := "UpdateTransitRouterVpnAttachmentAttribute"
		request["ClientToken"] = buildClientToken("UpdateTransitRouterVpnAttachmentAttribute")
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			response, err = client.RpcPost("Cbn", "2017-09-12", action, nil, request, true)
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
		stateConf := BuildStateConf([]string{}, []string{"Attached"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, cbnService.CenTransitRouterVpnAttachmentStateRefreshFunc(d.Id(), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}

		d.SetPartial("auto_publish_route_enabled")
		d.SetPartial("transit_router_attachment_description")
		d.SetPartial("transit_router_attachment_name")
	}

	d.Partial(false)

	return resourceAlicloudCenTransitRouterVpnAttachmentRead(d, meta)
}

func resourceAlicloudCenTransitRouterVpnAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cbnService := CbnService{client}
	action := "DeleteTransitRouterVpnAttachment"
	var response map[string]interface{}
	var err error
	request := map[string]interface{}{
		"TransitRouterAttachmentId": d.Id(),
	}

	request["ClientToken"] = buildClientToken("DeleteTransitRouterVpnAttachment")
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *resource.RetryError {
		response, err = client.RpcPost("Cbn", "2017-09-12", action, nil, request, true)
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
	stateConf := BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 5*time.Second, cbnService.CenTransitRouterVpnAttachmentStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return nil
}

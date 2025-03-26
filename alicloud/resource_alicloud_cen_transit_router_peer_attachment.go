package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudCenTransitRouterPeerAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudCenTransitRouterPeerAttachmentCreate,
		Read:   resourceAliCloudCenTransitRouterPeerAttachmentRead,
		Update: resourceAliCloudCenTransitRouterPeerAttachmentUpdate,
		Delete: resourceAliCloudCenTransitRouterPeerAttachmentDelete,
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
			"bandwidth": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"bandwidth_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"cen_bandwidth_package_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"cen_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"default_link_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: StringInSlice([]string{"Platinum", "Gold"}, false),
			},
			"dry_run": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"peer_transit_router_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"peer_transit_router_region_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
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
			"transit_router_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"transit_router_peer_attachment_name": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"transit_router_attachment_name"},
				Computed:      true,
			},
			"transit_router_attachment_name": {
				Type:       schema.TypeString,
				Optional:   true,
				Computed:   true,
				Deprecated: "Field 'transit_router_attachment_name' has been deprecated since provider version 1.247.0. New field 'transit_router_peer_attachment_name' instead.",
			},
			"route_table_association_enabled": {
				Type:       schema.TypeBool,
				Optional:   true,
				Deprecated: "Field `route_table_association_enabled` has been deprecated from provider version 1.230.0.",
			},
			"route_table_propagation_enabled": {
				Type:       schema.TypeBool,
				Optional:   true,
				Deprecated: "Field `route_table_propagation_enabled` has been deprecated from provider version 1.230.0.",
			},
			"resource_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  "TR",
			},
		},
	}
}

func resourceAliCloudCenTransitRouterPeerAttachmentCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateTransitRouterPeerAttachment"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken(action)

	if v, ok := d.GetOk("transit_router_id"); ok {
		request["TransitRouterId"] = v
	}
	if v, ok := d.GetOk("transit_router_attachment_description"); ok {
		request["TransitRouterAttachmentDescription"] = v
	}
	request["PeerTransitRouterId"] = d.Get("peer_transit_router_id")
	if v, ok := d.GetOk("peer_transit_router_region_id"); ok {
		request["PeerTransitRouterRegionId"] = v
	}
	if v, ok := d.GetOkExists("auto_publish_route_enabled"); ok {
		request["AutoPublishRouteEnabled"] = v
	}
	if v, ok := d.GetOkExists("bandwidth"); ok {
		request["Bandwidth"] = v
	}
	if v, ok := d.GetOk("cen_id"); ok {
		request["CenId"] = v
	}
	if v, ok := d.GetOk("cen_bandwidth_package_id"); ok {
		request["CenBandwidthPackageId"] = v
	}
	if v, ok := d.GetOk("tags"); ok {
		tagsMap := ConvertTags(v.(map[string]interface{}))
		request = expandTagsToMap(request, tagsMap)
	}

	if v, ok := d.GetOk("bandwidth_type"); ok {
		request["BandwidthType"] = v
	}
	if v, ok := d.GetOk("transit_router_attachment_name"); ok || d.HasChange("transit_router_attachment_name") {
		request["TransitRouterAttachmentName"] = v
	}

	if v, ok := d.GetOk("transit_router_peer_attachment_name"); ok {
		request["TransitRouterAttachmentName"] = v
	}
	if v, ok := d.GetOk("default_link_type"); ok {
		request["DefaultLinkType"] = v
	}
	if v, ok := d.GetOkExists("dry_run"); ok {
		request["DryRun"] = v
	}
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("Cbn", "2017-09-12", action, query, request, true)
		if err != nil {
			if IsExpectedErrors(err, []string{"Operation.Blocking", "Throttling.User", "IncorrectStatus.Status"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_cen_transit_router_peer_attachment", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%v:%v", request["CenId"], response["TransitRouterAttachmentId"]))

	cenServiceV2 := CenServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{"Attached"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, cenServiceV2.CenTransitRouterPeerAttachmentStateRefreshFunc(d.Id(), "Status", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAliCloudCenTransitRouterPeerAttachmentUpdate(d, meta)
}

func resourceAliCloudCenTransitRouterPeerAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cenServiceV2 := CenServiceV2{client}

	objectRaw, err := cenServiceV2.DescribeCenTransitRouterPeerAttachment(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_cen_transit_router_peer_attachment DescribeCenTransitRouterPeerAttachment Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("auto_publish_route_enabled", objectRaw["AutoPublishRouteEnabled"])
	d.Set("bandwidth", objectRaw["Bandwidth"])
	d.Set("bandwidth_type", objectRaw["BandwidthType"])
	d.Set("cen_bandwidth_package_id", objectRaw["CenBandwidthPackageId"])
	d.Set("cen_id", objectRaw["CenId"])
	d.Set("create_time", objectRaw["CreationTime"])
	d.Set("default_link_type", objectRaw["DefaultLinkType"])
	d.Set("peer_transit_router_id", objectRaw["PeerTransitRouterId"])
	d.Set("peer_transit_router_region_id", objectRaw["PeerTransitRouterRegionId"])
	d.Set("region_id", objectRaw["RegionId"])
	d.Set("status", objectRaw["Status"])
	d.Set("transit_router_attachment_description", objectRaw["TransitRouterAttachmentDescription"])
	d.Set("transit_router_id", objectRaw["TransitRouterId"])
	d.Set("transit_router_peer_attachment_name", objectRaw["TransitRouterAttachmentName"])
	d.Set("transit_router_attachment_id", objectRaw["TransitRouterAttachmentId"])

	tagsMaps := objectRaw["Tags"]
	d.Set("tags", tagsToMap(tagsMaps))

	d.Set("transit_router_attachment_name", d.Get("transit_router_peer_attachment_name"))
	return nil
}

func resourceAliCloudCenTransitRouterPeerAttachmentUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false

	var err error
	action := "UpdateTransitRouterPeerAttachmentAttribute"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	parts, _ := ParseResourceId(d.Id(), 2)
	request["TransitRouterAttachmentId"] = parts[1]
	if v, ok := d.GetOkExists("dry_run"); ok {
		request["DryRun"] = v
	}
	request["ClientToken"] = buildClientToken(action)
	if !d.IsNewResource() && d.HasChange("transit_router_attachment_description") {
		update = true
		request["TransitRouterAttachmentDescription"] = d.Get("transit_router_attachment_description")
	}

	if !d.IsNewResource() && d.HasChange("auto_publish_route_enabled") {
		update = true
		request["AutoPublishRouteEnabled"] = d.Get("auto_publish_route_enabled")
	}

	if !d.IsNewResource() && d.HasChange("bandwidth") {
		update = true
		request["Bandwidth"] = d.Get("bandwidth")
	}

	if !d.IsNewResource() && d.HasChange("cen_bandwidth_package_id") {
		update = true
		request["CenBandwidthPackageId"] = d.Get("cen_bandwidth_package_id")
	}

	if !d.IsNewResource() && d.HasChange("bandwidth_type") {
		update = true
		request["BandwidthType"] = d.Get("bandwidth_type")
	}

	if !d.IsNewResource() && d.HasChange("transit_router_attachment_name") {
		update = true
		request["TransitRouterAttachmentName"] = d.Get("transit_router_attachment_name")
	}

	if !d.IsNewResource() && d.HasChange("transit_router_peer_attachment_name") {
		update = true
		request["TransitRouterAttachmentName"] = d.Get("transit_router_peer_attachment_name")
	}

	if !d.IsNewResource() && d.HasChange("default_link_type") {
		update = true
		request["DefaultLinkType"] = d.Get("default_link_type")
	}

	if v, ok := d.GetOk("dry_run"); ok {
		request["DryRun"] = v
	}
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("Cbn", "2017-09-12", action, query, request, true)
			if err != nil {
				if IsExpectedErrors(err, []string{"Operation.Blocking", "Throttling.User", "IncorrectStatus.Status"}) || NeedRetry(err) {
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
		stateConf := BuildStateConf([]string{}, []string{"Attached"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, cenServiceV2.CenTransitRouterPeerAttachmentStateRefreshFunc(d.Id(), "Status", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}

	if d.HasChange("tags") {
		cbnService := CbnService{client}
		if err := cbnService.SetResourceTags(d, "TRANSITROUTERPEERATTACHMENT"); err != nil {
			return WrapError(err)
		}
	}
	return resourceAliCloudCenTransitRouterPeerAttachmentRead(d, meta)
}

func resourceAliCloudCenTransitRouterPeerAttachmentDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	action := "DeleteTransitRouterPeerAttachment"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	parts, _ := ParseResourceId(d.Id(), 2)
	request["TransitRouterAttachmentId"] = parts[1]
	if v, ok := d.GetOkExists("dry_run"); ok {
		request["DryRun"] = v
	}
	request["ClientToken"] = buildClientToken(action)

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("Cbn", "2017-09-12", action, query, request, true)
		request["ClientToken"] = buildClientToken(action)

		if err != nil {
			if IsExpectedErrors(err, []string{"Operation.Blocking", "Throttling.User", "IncorrectStatus.Status"}) || NeedRetry(err) {
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
	stateConf := BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 5*time.Second, cenServiceV2.CenTransitRouterPeerAttachmentStateRefreshFunc(d.Id(), "Status", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return nil
}

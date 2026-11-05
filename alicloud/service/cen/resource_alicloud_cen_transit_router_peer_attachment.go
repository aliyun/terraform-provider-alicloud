package cen

import (
	"fmt"
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/tags"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/validators"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAliCloudCenTransitRouterPeerAttachment() *schema.Resource {
	return &schema.Resource{
		Create: ResourceAliCloudCenTransitRouterPeerAttachmentCreate,
		Read:   ResourceAliCloudCenTransitRouterPeerAttachmentRead,
		Update: ResourceAliCloudCenTransitRouterPeerAttachmentUpdate,
		Delete: ResourceAliCloudCenTransitRouterPeerAttachmentDelete,
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
				ValidateFunc: validators.StringInSlice([]string{"Platinum", "Gold"}, false),
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
			"tags": tags.TagsSchema(),
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

func ResourceAliCloudCenTransitRouterPeerAttachmentCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateTransitRouterPeerAttachment"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["RegionId"] = client.RegionId
	request["ClientToken"] = helper.BuildClientToken(action)

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
		tagsMap := tags.ConvertTags(v.(map[string]interface{}))
		request = tags.ExpandTagsToMap(request, tagsMap)
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
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		response, err = client.RpcPost("Cbn", "2017-09-12", action, query, request, true)
		if err != nil {
			if errs.IsExpectedErrors(err, []string{"Operation.Blocking", "Throttling.User", "IncorrectStatus.Status"}) || helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)

	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_cen_transit_router_peer_attachment", action, errs.AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%v:%v", request["CenId"], response["TransitRouterAttachmentId"]))

	cenServiceV2 := NewCenServiceV2(client)
	stateConf := helper.BuildStateConf([]string{}, []string{"Attached"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, cenServiceV2.CenTransitRouterPeerAttachmentStateRefreshFunc(d.Id(), "Status", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errs.WrapErrorf(err, errs.IdMsg, d.Id())
	}

	return ResourceAliCloudCenTransitRouterPeerAttachmentUpdate(d, meta)
}

func ResourceAliCloudCenTransitRouterPeerAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cenServiceV2 := NewCenServiceV2(client)

	objectRaw, err := cenServiceV2.DescribeCenTransitRouterPeerAttachment(d.Id())
	if err != nil {
		if !d.IsNewResource() && errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_cen_transit_router_peer_attachment DescribeCenTransitRouterPeerAttachment Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
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
	d.Set("tags", tags.TagsToMap(tagsMaps))

	d.Set("transit_router_attachment_name", d.Get("transit_router_peer_attachment_name"))
	return nil
}

func ResourceAliCloudCenTransitRouterPeerAttachmentUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false

	var err error
	action := "UpdateTransitRouterPeerAttachmentAttribute"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	parts, _ := helper.ParseResourceId(d.Id(), 2)
	request["TransitRouterAttachmentId"] = parts[1]
	if v, ok := d.GetOkExists("dry_run"); ok {
		request["DryRun"] = v
	}
	request["ClientToken"] = helper.BuildClientToken(action)
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
		wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RpcPost("Cbn", "2017-09-12", action, query, request, true)
			if err != nil {
				if errs.IsExpectedErrors(err, []string{"Operation.Blocking", "Throttling.User", "IncorrectStatus.Status"}) || helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		logs.AddDebug(action, response, request)
		if err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
		}
		cenServiceV2 := NewCenServiceV2(client)
		stateConf := helper.BuildStateConf([]string{}, []string{"Attached"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, cenServiceV2.CenTransitRouterPeerAttachmentStateRefreshFunc(d.Id(), "Status", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return errs.WrapErrorf(err, errs.IdMsg, d.Id())
		}
	}

	if d.HasChange("tags") {
		cbnService := NewCbnService(client)
		if err := cbnService.SetResourceTags(d, "TRANSITROUTERPEERATTACHMENT"); err != nil {
			return errs.WrapError(err)
		}
	}
	return ResourceAliCloudCenTransitRouterPeerAttachmentRead(d, meta)
}

func ResourceAliCloudCenTransitRouterPeerAttachmentDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	action := "DeleteTransitRouterPeerAttachment"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	parts, _ := helper.ParseResourceId(d.Id(), 2)
	request["TransitRouterAttachmentId"] = parts[1]
	if v, ok := d.GetOkExists("dry_run"); ok {
		request["DryRun"] = v
	}
	request["ClientToken"] = helper.BuildClientToken(action)

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		response, err = client.RpcPost("Cbn", "2017-09-12", action, query, request, true)
		request["ClientToken"] = helper.BuildClientToken(action)

		if err != nil {
			if errs.IsExpectedErrors(err, []string{"Operation.Blocking", "Throttling.User", "IncorrectStatus.Status"}) || helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)

	if err != nil {
		if errs.NotFoundError(err) {
			return nil
		}
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
	}

	cenServiceV2 := NewCenServiceV2(client)
	stateConf := helper.BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 5*time.Second, cenServiceV2.CenTransitRouterPeerAttachmentStateRefreshFunc(d.Id(), "Status", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errs.WrapErrorf(err, errs.IdMsg, d.Id())
	}

	return nil
}

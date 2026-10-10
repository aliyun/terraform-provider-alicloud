package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudPrivateLinkVpcEndpointService() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudPrivateLinkVpcEndpointServiceCreate,
		Read:   resourceAliCloudPrivateLinkVpcEndpointServiceRead,
		Update: resourceAliCloudPrivateLinkVpcEndpointServiceUpdate,
		Delete: resourceAliCloudPrivateLinkVpcEndpointServiceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"address_ip_version": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: StringInSlice([]string{"IPv4", "IPv6", "DualStack"}, false),
			},
			"auto_accept_connection": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"connect_bandwidth": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: IntBetween(100, 10240),
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"dry_run": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"payer": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"Endpoint", "EndpointService"}, false),
			},
			"region_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"resource": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				MaxItems: 10,
				// Hash an element on resource_id + resource_type + zone_id. zone_id
				// is part of the backend identity: the dedicated subresource
				// alicloud_privatelink_vpc_endpoint_service_resource keys on
				// ServiceId:ResourceId:ZoneId, and Attach/Detach send ZoneId per
				// backend. Hashing only on resource_id + resource_type silently
				// deduped two backends with the same resource_id + resource_type
				// in different zones, dropping the cross-zone backend and
				// producing a perpetual diff. To keep a server-assigned zone_id
				// (Optional + Computed) from churning element identity now that
				// zone_id is in the hash, Read preserves the user-configured
				// zone_id instead of overwriting an empty one with the
				// server-returned zone (see Read below).
				Set: func(v interface{}) int {
					m := v.(map[string]interface{})
					return helper.Hashcode(fmt.Sprintf("%s|%s|%s", m["resource_id"], m["resource_type"], m["zone_id"]))
				},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"resource_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"resource_type": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: StringInSlice([]string{"slb", "alb", "nlb", "gwlb"}, false),
						},
						"zone_id": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"service_business_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"service_description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"service_domain": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"service_resource_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"slb", "alb", "nlb", "gwlb"}, false),
			},
			"service_support_ipv6": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"supported_region_list": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": tagsSchema(),
			"vpc_endpoint_service_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"zone_affinity_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func resourceAliCloudPrivateLinkVpcEndpointServiceCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateVpcEndpointService"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken(action)

	if v, ok := d.GetOk("tags"); ok {
		tagsMap := ConvertTags(v.(map[string]interface{}))
		request = expandTagsToMap(request, tagsMap)
	}

	if v, ok := d.GetOk("resource_group_id"); ok {
		request["ResourceGroupId"] = v
	}
	if v, ok := d.GetOkExists("service_support_ipv6"); ok {
		request["ServiceSupportIPv6"] = v
	}
	if v, ok := d.GetOk("service_resource_type"); ok {
		request["ServiceResourceType"] = v
	}
	if v, ok := d.GetOkExists("zone_affinity_enabled"); ok {
		request["ZoneAffinityEnabled"] = v
	}
	if v, ok := d.GetOk("address_ip_version"); ok {
		request["AddressIpVersion"] = v
	}
	if v, ok := d.GetOk("service_description"); ok {
		request["ServiceDescription"] = v
	}
	if v, ok := d.GetOkExists("auto_accept_connection"); ok {
		request["AutoAcceptEnabled"] = v
	}
	if v, ok := d.GetOk("supported_region_list"); ok {
		request["SupportedRegionList"] = v.(*schema.Set).List()
	}
	if v, ok := d.GetOk("resource"); ok {
		resourceMapsArray := make([]interface{}, 0)
		for _, dataLoop := range v.(*schema.Set).List() {
			dataLoopTmp := dataLoop.(map[string]interface{})
			dataLoopMap := make(map[string]interface{})
			dataLoopMap["ResourceId"] = dataLoopTmp["resource_id"]
			dataLoopMap["ResourceType"] = dataLoopTmp["resource_type"]
			if zoneId := vpcEndpointServiceResourceZoneId(dataLoopTmp["zone_id"]); zoneId != "" {
				dataLoopMap["ZoneId"] = zoneId
			}
			resourceMapsArray = append(resourceMapsArray, dataLoopMap)
		}
		request["Resource"] = resourceMapsArray
	}
	if v, ok := d.GetOkExists("dry_run"); ok {
		request["DryRun"] = v
	}
	if v, ok := d.GetOk("payer"); ok {
		request["Payer"] = v
	}
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("Privatelink", "2020-04-15", action, query, request, true)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_privatelink_vpc_endpoint_service", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["ServiceId"]))

	privateLinkServiceV2 := PrivateLinkServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{"Active"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, privateLinkServiceV2.PrivateLinkVpcEndpointServiceStateRefreshFunc(d.Id(), "ServiceStatus", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAliCloudPrivateLinkVpcEndpointServiceUpdate(d, meta)
}

func resourceAliCloudPrivateLinkVpcEndpointServiceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	privateLinkServiceV2 := PrivateLinkServiceV2{client}

	objectRaw, err := privateLinkServiceV2.DescribePrivateLinkVpcEndpointService(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_privatelink_vpc_endpoint_service DescribePrivateLinkVpcEndpointService Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("address_ip_version", objectRaw["AddressIpVersion"])
	d.Set("auto_accept_connection", objectRaw["AutoAcceptEnabled"])
	setVpcEndpointServiceConnectBandwidth(d, objectRaw)
	d.Set("create_time", objectRaw["CreateTime"])
	d.Set("payer", objectRaw["Payer"])
	d.Set("region_id", objectRaw["RegionId"])
	d.Set("resource_group_id", objectRaw["ResourceGroupId"])
	d.Set("service_business_status", objectRaw["ServiceBusinessStatus"])
	d.Set("service_description", objectRaw["ServiceDescription"])
	d.Set("service_domain", objectRaw["ServiceDomain"])
	d.Set("service_resource_type", objectRaw["ServiceResourceType"])
	d.Set("service_support_ipv6", objectRaw["ServiceSupportIPv6"])
	d.Set("supported_region_list", convertSupportedRegionSetToStringList(objectRaw["SupportedRegionSet"]))
	d.Set("status", objectRaw["ServiceStatus"])
	d.Set("vpc_endpoint_service_name", objectRaw["ServiceName"])
	d.Set("zone_affinity_enabled", objectRaw["ZoneAffinityEnabled"])

	resourcesRaw, err := privateLinkServiceV2.ListPrivateLinkVpcEndpointServiceResources(d.Id())
	if err != nil {
		return WrapError(err)
	}
	// Collect the zone_id values the user already configured for each
	// (resource_id, resource_type) pair. zone_id is now part of the TypeSet
	// element identity (see the Set hash), so a server-assigned zone read back
	// into state would change an element's hash and cause perpetual diff /
	// detach-attach churn. For a backend the user left unpinned (zone_id == ""),
	// keep "" instead of adopting the server-returned zone; for a backend the
	// user pinned to a zone, confirm that exact zone.
	configuredZones := map[string]map[string]bool{}
	if rs, ok := d.Get("resource").(*schema.Set); ok {
		for _, e := range rs.List() {
			m := e.(map[string]interface{})
			pair := fmt.Sprintf("%s|%s", m["resource_id"], m["resource_type"])
			z := vpcEndpointServiceResourceZoneId(m["zone_id"])
			if configuredZones[pair] == nil {
				configuredZones[pair] = map[string]bool{}
			}
			configuredZones[pair][z] = true
		}
	}
	resourceMaps := make([]map[string]interface{}, 0, len(resourcesRaw))
	for _, v := range resourcesRaw {
		item := v.(map[string]interface{})
		resourceId := item["ResourceId"]
		resourceType := item["ResourceType"]
		// Normalize ZoneId so an absent/nil API value is stored as "" rather
		// than nil; this keeps TypeSet element hashing stable (nil vs "" hash
		// differently).
		apiZoneId := vpcEndpointServiceResourceZoneId(item["ZoneId"])
		zoneId := apiZoneId
		if zones, ok := configuredZones[fmt.Sprintf("%s|%s", resourceId, resourceType)]; ok && apiZoneId != "" {
			// Only adopt the server-returned zone when it matches one the user
			// configured; if the user left zone_id empty, keep it empty so the
			// element's hash stays stable and no churn is introduced.
			if !zones[apiZoneId] && zones[""] {
				zoneId = ""
			}
		}
		resourceMaps = append(resourceMaps, map[string]interface{}{
			"resource_id":   resourceId,
			"resource_type": resourceType,
			"zone_id":       zoneId,
		})
	}
	d.Set("resource", resourceMaps)

	objectRaw, err = privateLinkServiceV2.DescribeVpcEndpointServiceListTagResources(d.Id())
	if err != nil && !NotFoundError(err) {
		return WrapError(err)
	}

	tagsMaps := objectRaw["TagResources"]
	d.Set("tags", tagsToMap(tagsMaps))

	return nil
}

func setVpcEndpointServiceConnectBandwidth(d *schema.ResourceData, objectRaw map[string]interface{}) {
	connectBandwidth, ok := objectRaw["ConnectBandwidth"]
	serviceResourceType := fmt.Sprint(objectRaw["ServiceResourceType"])
	if serviceResourceType == "" || serviceResourceType == "<nil>" {
		if v, ok := d.GetOk("service_resource_type"); ok {
			serviceResourceType = v.(string)
		}
	}

	if serviceResourceType == "slb" {
		d.Set("connect_bandwidth", connectBandwidth)
		return
	}
	if ok && !isEmptyVpcEndpointServiceConnectBandwidth(connectBandwidth) {
		d.Set("connect_bandwidth", connectBandwidth)
	}
}

func isEmptyVpcEndpointServiceConnectBandwidth(v interface{}) bool {
	switch vv := v.(type) {
	case nil:
		return true
	case int:
		return vv == 0
	case int64:
		return vv == 0
	case float64:
		return vv == 0
	case string:
		return vv == "" || vv == "0"
	default:
		return fmt.Sprint(vv) == "" || fmt.Sprint(vv) == "0"
	}
}

// vpcEndpointServiceResourceZoneId returns a non-empty string form of a zone_id
// set value, or "" when the value is absent, nil, empty, or the "<nil>" string
// produced by fmt.Sprint(nil). Callers use the empty result to omit ZoneId from
// upstream requests, so a nil zone_id can never be sent as a literal ZoneId.
func vpcEndpointServiceResourceZoneId(zoneId interface{}) string {
	if zoneId == nil {
		return ""
	}
	s := fmt.Sprint(zoneId)
	if s == "" || s == "<nil>" {
		return ""
	}
	return s
}

func resourceAliCloudPrivateLinkVpcEndpointServiceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false
	d.Partial(true)

	var err error
	action := "UpdateVpcEndpointServiceAttribute"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["ServiceId"] = d.Id()
	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken(action)
	if !d.IsNewResource() && d.HasChange("service_description") {
		update = true
		request["ServiceDescription"] = d.Get("service_description")
	}

	if !d.IsNewResource() && d.HasChange("auto_accept_connection") {
		update = true
		request["AutoAcceptEnabled"] = d.Get("auto_accept_connection")
	}

	if !d.IsNewResource() && d.HasChange("service_support_ipv6") {
		update = true
		request["ServiceSupportIPv6"] = d.Get("service_support_ipv6")
	}

	if !d.IsNewResource() && d.HasChange("zone_affinity_enabled") {
		update = true
		request["ZoneAffinityEnabled"] = d.Get("zone_affinity_enabled")
	}

	if d.HasChange("connect_bandwidth") {
		update = true
		request["ConnectBandwidth"] = d.Get("connect_bandwidth")
	}

	if v, ok := d.GetOkExists("dry_run"); ok {
		request["DryRun"] = v
	}
	if !d.IsNewResource() && d.HasChange("address_ip_version") {
		update = true
		request["AddressIpVersion"] = d.Get("address_ip_version")
	}

	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("Privatelink", "2020-04-15", action, query, request, true)
			if err != nil {
				if IsExpectedErrors(err, []string{"EndpointServiceOperationDenied", "ConcurrentCallNotSupported", "EndpointServiceLocked"}) || NeedRetry(err) {
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
	}
	if !d.IsNewResource() && d.HasChange("supported_region_list") {
		oldEntry, newEntry := d.GetChange("supported_region_list")
		oldEntrySet := oldEntry.(*schema.Set)
		newEntrySet := newEntry.(*schema.Set)
		removed := oldEntrySet.Difference(newEntrySet)
		added := newEntrySet.Difference(oldEntrySet)

		if added.Len() > 0 {
			action = "UpdateVpcEndpointServiceAttribute"
			request = make(map[string]interface{})
			query = make(map[string]interface{})
			request["ServiceId"] = d.Id()
			request["RegionId"] = client.RegionId
			request["ClientToken"] = buildClientToken(action)
			request["AddSupportedRegionSet"] = added.List()
			if v, ok := d.GetOkExists("dry_run"); ok {
				request["DryRun"] = v
			}
			wait := incrementalWait(3*time.Second, 5*time.Second)
			err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
				response, err = client.RpcPost("Privatelink", "2020-04-15", action, query, request, true)
				if err != nil {
					if IsExpectedErrors(err, []string{"EndpointServiceOperationDenied", "ConcurrentCallNotSupported", "EndpointServiceLocked"}) || NeedRetry(err) {
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
		}

		if removed.Len() > 0 {
			action = "UpdateVpcEndpointServiceAttribute"
			request = make(map[string]interface{})
			query = make(map[string]interface{})
			request["ServiceId"] = d.Id()
			request["RegionId"] = client.RegionId
			request["ClientToken"] = buildClientToken(action)
			request["DeleteSupportedRegionSet"] = removed.List()
			if v, ok := d.GetOkExists("dry_run"); ok {
				request["DryRun"] = v
			}
			wait := incrementalWait(3*time.Second, 5*time.Second)
			err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
				response, err = client.RpcPost("Privatelink", "2020-04-15", action, query, request, true)
				if err != nil {
					if IsExpectedErrors(err, []string{"EndpointServiceOperationDenied", "ConcurrentCallNotSupported", "EndpointServiceLocked"}) || NeedRetry(err) {
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
		}

		privateLinkServiceV2 := PrivateLinkServiceV2{client}
		stateConf := BuildStateConf([]string{}, []string{"Active"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, privateLinkServiceV2.PrivateLinkVpcEndpointServiceStateRefreshFunc(d.Id(), "ServiceStatus", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}
	if !d.IsNewResource() && d.HasChange("resource") {
		oldEntry, newEntry := d.GetChange("resource")
		removed := oldEntry.(*schema.Set).Difference(newEntry.(*schema.Set))
		added := newEntry.(*schema.Set).Difference(oldEntry.(*schema.Set))

		for _, item := range removed.List() {
			resourceItem := item.(map[string]interface{})
			action = "DetachResourceFromVpcEndpointService"
			request = make(map[string]interface{})
			query = make(map[string]interface{})
			request["ServiceId"] = d.Id()
			request["RegionId"] = client.RegionId
			request["ClientToken"] = buildClientToken(action)
			request["ResourceId"] = resourceItem["resource_id"]
			request["ResourceType"] = resourceItem["resource_type"]
			if zoneId := vpcEndpointServiceResourceZoneId(resourceItem["zone_id"]); zoneId != "" {
				request["ZoneId"] = zoneId
			}
			if v, ok := d.GetOkExists("dry_run"); ok {
				request["DryRun"] = v
			}
			wait := incrementalWait(3*time.Second, 5*time.Second)
			err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
				response, err = client.RpcPost("Privatelink", "2020-04-15", action, query, request, true)
				if err != nil {
					if IsExpectedErrors(err, []string{"EndpointServiceOperationDenied", "ConcurrentCallNotSupported", "EndpointServiceLocked"}) || NeedRetry(err) {
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
		}

		for _, item := range added.List() {
			resourceItem := item.(map[string]interface{})
			action = "AttachResourceToVpcEndpointService"
			request = make(map[string]interface{})
			query = make(map[string]interface{})
			request["ServiceId"] = d.Id()
			request["RegionId"] = client.RegionId
			request["ClientToken"] = buildClientToken(action)
			request["ResourceId"] = resourceItem["resource_id"]
			request["ResourceType"] = resourceItem["resource_type"]
			if zoneId := vpcEndpointServiceResourceZoneId(resourceItem["zone_id"]); zoneId != "" {
				request["ZoneId"] = zoneId
			}
			if v, ok := d.GetOkExists("dry_run"); ok {
				request["DryRun"] = v
			}
			wait := incrementalWait(3*time.Second, 5*time.Second)
			err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
				response, err = client.RpcPost("Privatelink", "2020-04-15", action, query, request, true)
				if err != nil {
					if IsExpectedErrors(err, []string{"EndpointServiceOperationDenied", "ConcurrentCallNotSupported", "EndpointServiceLocked"}) || NeedRetry(err) {
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
		}

		privateLinkServiceV2 := PrivateLinkServiceV2{client}
		stateConf := BuildStateConf([]string{}, []string{"Active"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, privateLinkServiceV2.PrivateLinkVpcEndpointServiceStateRefreshFunc(d.Id(), "ServiceStatus", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}
	update = false
	action = "ChangeResourceGroup"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["ResourceId"] = d.Id()

	if _, ok := d.GetOk("resource_group_id"); ok && !d.IsNewResource() && d.HasChange("resource_group_id") {
		update = true
	}
	request["ResourceGroupId"] = d.Get("resource_group_id")
	if !d.IsNewResource() && d.HasChange("region_id") {
		update = true
		request["ResourceRegionId"] = d.Get("region_id")
	}

	request["ResourceType"] = "VpcEndpointService"
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("Privatelink", "2020-04-15", action, query, request, true)
			if err != nil {
				if NeedRetry(err) {
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
		privateLinkServiceV2 := PrivateLinkServiceV2{client}
		stateConf := BuildStateConf([]string{}, []string{fmt.Sprint(d.Get("resource_group_id"))}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, privateLinkServiceV2.PrivateLinkVpcEndpointServiceStateRefreshFunc(d.Id(), "ResourceGroupId", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}

	if d.HasChange("tags") {
		privateLinkServiceV2 := PrivateLinkServiceV2{client}
		if err := privateLinkServiceV2.SetResourceTags(d, "VpcEndpointService"); err != nil {
			return WrapError(err)
		}
	}
	d.Partial(false)
	return resourceAliCloudPrivateLinkVpcEndpointServiceRead(d, meta)
}

func resourceAliCloudPrivateLinkVpcEndpointServiceDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	action := "DeleteVpcEndpointService"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["ServiceId"] = d.Id()
	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken(action)

	if v, ok := d.GetOkExists("dry_run"); ok {
		request["DryRun"] = v
	}
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("Privatelink", "2020-04-15", action, query, request, true)
		request["ClientToken"] = buildClientToken(action)

		if err != nil {
			if IsExpectedErrors(err, []string{"EndpointServiceConnectionDependence", "ConcurrentCallNotSupported"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		if IsExpectedErrors(err, []string{"EndpointServiceNotFound"}) || NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	privateLinkServiceV2 := PrivateLinkServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 5*time.Second, privateLinkServiceV2.PrivateLinkVpcEndpointServiceStateRefreshFunc(d.Id(), "ServiceStatus", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return nil
}

func convertSupportedRegionSetToStringList(src interface{}) (result []interface{}) {
	if src == nil {
		return
	}
	for _, v := range convertToInterfaceArray(src) {
		vv, ok := v.(map[string]interface{})
		if !ok {
			continue
		}
		if fmt.Sprint(vv["RegionServiceStatus"]) == "Closed" {
			continue
		}
		if regionId, ok := vv["SupportedRegionId"]; ok {
			result = append(result, fmt.Sprint(regionId))
		}
	}
	return
}

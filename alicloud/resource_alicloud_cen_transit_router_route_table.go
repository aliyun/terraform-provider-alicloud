// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudCenTransitRouterRouteTable() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudCenTransitRouterRouteTableCreate,
		Read:   resourceAliCloudCenTransitRouterRouteTableRead,
		Update: resourceAliCloudCenTransitRouterRouteTableUpdate,
		Delete: resourceAliCloudCenTransitRouterRouteTableDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"region_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"route_table_options": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"multi_region_ecmp": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": tagsSchema(),
			"transit_router_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"transit_router_route_table_description": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"transit_router_route_table_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"transit_router_route_table_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"transit_router_route_table_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAliCloudCenTransitRouterRouteTableCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateTransitRouterRouteTable"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})

	request["ClientToken"] = buildClientToken(action)

	if v, ok := d.GetOk("transit_router_route_table_description"); ok {
		request["TransitRouterRouteTableDescription"] = v
	}
	if v, ok := d.GetOk("transit_router_route_table_name"); ok {
		request["TransitRouterRouteTableName"] = v
	}
	request["TransitRouterId"] = d.Get("transit_router_id")
	if v, ok := d.GetOk("tags"); ok {
		tagsMap := ConvertTags(v.(map[string]interface{}))
		request = expandTagsToMap(request, tagsMap)
	}

	if v, ok := d.GetOk("route_table_options"); ok {
		routeTableOptionsMultiRegionEcmpJsonPath, err := jsonpath.Get("$[0].multi_region_ecmp", v)
		if err == nil && routeTableOptionsMultiRegionEcmpJsonPath != "" {
			request["RouteTableOptions.MultiRegionECMP"] = routeTableOptionsMultiRegionEcmpJsonPath
		}
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_cen_transit_router_route_table", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["TransitRouterRouteTableId"]))

	cenServiceV2 := CenServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{"Active"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, cenServiceV2.CenTransitRouterRouteTableStateRefreshFunc(d.Id(), "TransitRouterRouteTableStatus", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAliCloudCenTransitRouterRouteTableRead(d, meta)
}

func resourceAliCloudCenTransitRouterRouteTableRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cenServiceV2 := CenServiceV2{client}

	objectRaw, err := cenServiceV2.DescribeCenTransitRouterRouteTable(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_cen_transit_router_route_table DescribeCenTransitRouterRouteTable Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("create_time", objectRaw["CreateTime"])
	d.Set("region_id", objectRaw["RegionId"])
	d.Set("status", objectRaw["TransitRouterRouteTableStatus"])
	d.Set("transit_router_id", objectRaw["TransitRouterId"])
	d.Set("transit_router_route_table_description", objectRaw["TransitRouterRouteTableDescription"])
	d.Set("transit_router_route_table_name", objectRaw["TransitRouterRouteTableName"])
	d.Set("transit_router_route_table_type", objectRaw["TransitRouterRouteTableType"])
	d.Set("transit_router_route_table_id", objectRaw["TransitRouterRouteTableId"])

	routeTableOptionsMaps := make([]map[string]interface{}, 0)
	routeTableOptionsMap := make(map[string]interface{})
	routeTableOptionsRaw := make(map[string]interface{})
	if objectRaw["RouteTableOptions"] != nil {
		routeTableOptionsRaw = objectRaw["RouteTableOptions"].(map[string]interface{})
	}
	if len(routeTableOptionsRaw) > 0 {
		routeTableOptionsMap["multi_region_ecmp"] = routeTableOptionsRaw["MultiRegionECMP"]

		routeTableOptionsMaps = append(routeTableOptionsMaps, routeTableOptionsMap)
	}
	if err := d.Set("route_table_options", routeTableOptionsMaps); err != nil {
		return err
	}
	tagsMaps := objectRaw["Tags"]
	d.Set("tags", tagsToMap(tagsMaps))

	return nil
}

func resourceAliCloudCenTransitRouterRouteTableUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false

	var err error
	action := "UpdateTransitRouterRouteTable"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["TransitRouterRouteTableId"] = d.Id()

	request["ClientToken"] = buildClientToken(action)
	if d.HasChange("transit_router_route_table_description") {
		update = true
		request["TransitRouterRouteTableDescription"] = d.Get("transit_router_route_table_description")
	}

	if d.HasChange("transit_router_route_table_name") {
		update = true
		request["TransitRouterRouteTableName"] = d.Get("transit_router_route_table_name")
	}

	if d.HasChange("route_table_options.0.multi_region_ecmp") {
		update = true
		routeTableOptionsMultiRegionEcmpJsonPath, err := jsonpath.Get("$[0].multi_region_ecmp", d.Get("route_table_options"))
		if err == nil {
			request["RouteTableOptions.MultiRegionECMP"] = routeTableOptionsMultiRegionEcmpJsonPath
		}
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
		stateConf := BuildStateConf([]string{}, []string{"Active"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, cenServiceV2.CenTransitRouterRouteTableStateRefreshFunc(d.Id(), "TransitRouterRouteTableStatus", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}

	if d.HasChange("tags") {
		cenServiceV2 := CenServiceV2{client}
		if err := cenServiceV2.SetResourceTags(d, "TRANSITROUTERROUTETABLE"); err != nil {
			return WrapError(err)
		}
	}
	return resourceAliCloudCenTransitRouterRouteTableRead(d, meta)
}

func resourceAliCloudCenTransitRouterRouteTableDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	action := "DeleteTransitRouterRouteTable"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["TransitRouterRouteTableId"] = d.Id()

	request["ClientToken"] = buildClientToken(action)

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("Cbn", "2017-09-12", action, query, request, true)
		if err != nil {
			if IsExpectedErrors(err, []string{"Operation.Blocking", "Throttling.User", "IncorrectStatus.Status", "TokenProcessing"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidTransitRouterRouteTableId.NotFound"}) || NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	cenServiceV2 := CenServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{""}, d.Timeout(schema.TimeoutDelete), 5*time.Second, cenServiceV2.CenTransitRouterRouteTableStateRefreshFunc(d.Id(), "TransitRouterRouteTableStatus", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return nil
}

package alicloud

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudCenTransitRouteTableAggregation() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudCenTransitRouteTableAggregationCreate,
		Read:   resourceAliCloudCenTransitRouteTableAggregationRead,
		Update: resourceAliCloudCenTransitRouteTableAggregationUpdate,
		Delete: resourceAliCloudCenTransitRouteTableAggregationDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"transit_route_table_aggregation_cidr": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"transit_route_table_aggregation_description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"transit_route_table_aggregation_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"transit_route_table_aggregation_scope": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"transit_route_table_aggregation_scope_list": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"transit_route_table_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAliCloudCenTransitRouteTableAggregationCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateTransitRouteTableAggregation"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	if v, ok := d.GetOk("transit_route_table_aggregation_cidr"); ok {
		request["TransitRouteTableAggregationCidr"] = v
	}
	if v, ok := d.GetOk("transit_route_table_id"); ok {
		request["TransitRouteTableId"] = v
	}

	request["ClientToken"] = buildClientToken(action)

	if v, ok := d.GetOk("transit_route_table_aggregation_name"); ok {
		request["TransitRouteTableAggregationName"] = v
	}
	if v, ok := d.GetOk("transit_route_table_aggregation_description"); ok {
		request["TransitRouteTableAggregationDescription"] = v
	}
	if v, ok := d.GetOk("transit_route_table_aggregation_scope"); ok {
		request["TransitRouteTableAggregationScope"] = v
	}
	if v, ok := d.GetOk("transit_route_table_aggregation_scope_list"); ok {
		transitRouteTableAggregationScopeListMapsArray := v.(*schema.Set).List()
		transitRouteTableAggregationScopeListMapsJson, err := json.Marshal(transitRouteTableAggregationScopeListMapsArray)
		if err != nil {
			return WrapError(err)
		}
		request["TransitRouteTableAggregationScopeList"] = string(transitRouteTableAggregationScopeListMapsJson)
	}

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("Cbn", "2017-09-12", action, query, request, true)
		if err != nil {
			if IsExpectedErrors(err, []string{"Operation.Blocking", "IncorrectStatus.TransitRouteTable", "IncorrectStatus.TransitRouter"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_cen_transit_route_table_aggregation", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%v#%v", request["TransitRouteTableId"], request["TransitRouteTableAggregationCidr"]))

	cenServiceV2 := CenServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{"AllConfigured"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, cenServiceV2.CenTransitRouteTableAggregationStateRefreshFunc(d.Id(), "Status", []string{"ConfigFailed", "PartialConfigured"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAliCloudCenTransitRouteTableAggregationRead(d, meta)
}

func resourceAliCloudCenTransitRouteTableAggregationRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cenServiceV2 := CenServiceV2{client}

	objectRaw, err := cenServiceV2.DescribeCenTransitRouteTableAggregation(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_cen_transit_route_table_aggregation DescribeCenTransitRouteTableAggregation Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("status", objectRaw["Status"])
	d.Set("transit_route_table_aggregation_description", objectRaw["Description"])
	d.Set("transit_route_table_aggregation_name", objectRaw["Name"])
	d.Set("transit_route_table_aggregation_scope", objectRaw["Scope"])
	d.Set("transit_route_table_aggregation_cidr", objectRaw["TransitRouteTableAggregationCidr"])
	d.Set("transit_route_table_id", objectRaw["TrRouteTableId"])

	scopeListRaw := make([]interface{}, 0)
	if objectRaw["ScopeList"] != nil {
		scopeListRaw = objectRaw["ScopeList"].([]interface{})
	}

	d.Set("transit_route_table_aggregation_scope_list", scopeListRaw)

	return nil
}

func resourceAliCloudCenTransitRouteTableAggregationUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false

	var err error
	parts := strings.Split(d.Id(), "#")
	if len(parts) != 2 {
		parts = strings.Split(d.Id(), ":")
	}
	action := "ModifyTransitRouteTableAggregation"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["TransitRouteTableId"] = parts[0]
	request["TransitRouteTableAggregationCidr"] = parts[1]

	request["ClientToken"] = buildClientToken(action)
	if d.HasChange("transit_route_table_aggregation_name") {
		update = true
		request["TransitRouteTableAggregationName"] = d.Get("transit_route_table_aggregation_name")
	}

	if d.HasChange("transit_route_table_aggregation_scope") {
		update = true
		request["TransitRouteTableAggregationScope"] = d.Get("transit_route_table_aggregation_scope")
	}

	if d.HasChange("transit_route_table_aggregation_scope_list") {
		update = true
		if v, ok := d.GetOk("transit_route_table_aggregation_scope_list"); ok || d.HasChange("transit_route_table_aggregation_scope_list") {
			transitRouteTableAggregationScopeListMapsArray := v.(*schema.Set).List()
			transitRouteTableAggregationScopeListMapsJson, err := json.Marshal(transitRouteTableAggregationScopeListMapsArray)
			if err != nil {
				return WrapError(err)
			}
			request["TransitRouteTableAggregationScopeList"] = string(transitRouteTableAggregationScopeListMapsJson)
		}
	}

	if d.HasChange("transit_route_table_aggregation_description") {
		update = true
		request["TransitRouteTableAggregationDescription"] = d.Get("transit_route_table_aggregation_description")
	}

	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("Cbn", "2017-09-12", action, query, request, true)
			if err != nil {
				if IsExpectedErrors(err, []string{"Operation.Blocking", "IncorrectStatus.TransitRouteTable", "IncorrectStatus.TransitRouter", "IncorrectStatus.AggregationRoute"}) || NeedRetry(err) {
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
		stateConf := BuildStateConf([]string{}, []string{"AllConfigured"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, cenServiceV2.CenTransitRouteTableAggregationStateRefreshFunc(d.Id(), "Status", []string{"PartialConfigured", "ConfigFailed"}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}

	return resourceAliCloudCenTransitRouteTableAggregationRead(d, meta)
}

func resourceAliCloudCenTransitRouteTableAggregationDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	parts := strings.Split(d.Id(), "#")
	if len(parts) != 2 {
		parts = strings.Split(d.Id(), ":")
	}
	action := "DeleteTransitRouteTableAggregation"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["TransitRouteTableAggregationCidr"] = parts[1]
	request["TransitRouteTableId"] = parts[0]

	request["ClientToken"] = buildClientToken(action)

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("Cbn", "2017-09-12", action, query, request, true)
		request["ClientToken"] = buildClientToken(action)

		if err != nil {
			if IsExpectedErrors(err, []string{"Operation.Blocking", "IncorrectStatus.TransitRouter", "IncorrectStatus.TransitRouteTable"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		if IsExpectedErrors(err, []string{"InstanceNotExist.AggregationRoute"}) || NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	cenServiceV2 := CenServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 5*time.Second, cenServiceV2.CenTransitRouteTableAggregationStateRefreshFunc(d.Id(), "TransitRouteTableAggregationCidr", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return nil
}

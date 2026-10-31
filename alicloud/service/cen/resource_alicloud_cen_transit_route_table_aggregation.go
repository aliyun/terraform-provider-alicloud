package cen

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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

	request["ClientToken"] = helper.BuildClientToken(action)

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
			return errs.WrapError(err)
		}
		request["TransitRouteTableAggregationScopeList"] = string(transitRouteTableAggregationScopeListMapsJson)
	}

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		response, err = client.RpcPost("Cbn", "2017-09-12", action, query, request, true)
		if err != nil {
			if errs.IsExpectedErrors(err, []string{"Operation.Blocking", "IncorrectStatus.TransitRouteTable", "IncorrectStatus.TransitRouter"}) || helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)

	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_cen_transit_route_table_aggregation", action, errs.AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%v#%v", request["TransitRouteTableId"], request["TransitRouteTableAggregationCidr"]))

	cenServiceV2 := NewCenServiceV2(client)
	stateConf := helper.BuildStateConf([]string{}, []string{"AllConfigured"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, cenServiceV2.CenTransitRouteTableAggregationStateRefreshFunc(d.Id(), "Status", []string{"ConfigFailed", "PartialConfigured"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errs.WrapErrorf(err, errs.IdMsg, d.Id())
	}

	return resourceAliCloudCenTransitRouteTableAggregationRead(d, meta)
}

func resourceAliCloudCenTransitRouteTableAggregationRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cenServiceV2 := NewCenServiceV2(client)

	objectRaw, err := cenServiceV2.DescribeCenTransitRouteTableAggregation(d.Id())
	if err != nil {
		if !d.IsNewResource() && errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_cen_transit_route_table_aggregation DescribeCenTransitRouteTableAggregation Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
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

	request["ClientToken"] = helper.BuildClientToken(action)
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
				return errs.WrapError(err)
			}
			request["TransitRouteTableAggregationScopeList"] = string(transitRouteTableAggregationScopeListMapsJson)
		}
	}

	if d.HasChange("transit_route_table_aggregation_description") {
		update = true
		request["TransitRouteTableAggregationDescription"] = d.Get("transit_route_table_aggregation_description")
	}

	if update {
		wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RpcPost("Cbn", "2017-09-12", action, query, request, true)
			if err != nil {
				if errs.IsExpectedErrors(err, []string{"Operation.Blocking", "IncorrectStatus.TransitRouteTable", "IncorrectStatus.TransitRouter", "IncorrectStatus.AggregationRoute"}) || helper.NeedRetry(err) {
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
		stateConf := helper.BuildStateConf([]string{}, []string{"AllConfigured"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, cenServiceV2.CenTransitRouteTableAggregationStateRefreshFunc(d.Id(), "Status", []string{"PartialConfigured", "ConfigFailed"}))
		if _, err := stateConf.WaitForState(); err != nil {
			return errs.WrapErrorf(err, errs.IdMsg, d.Id())
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

	request["ClientToken"] = helper.BuildClientToken(action)

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		response, err = client.RpcPost("Cbn", "2017-09-12", action, query, request, true)
		request["ClientToken"] = helper.BuildClientToken(action)

		if err != nil {
			if errs.IsExpectedErrors(err, []string{"Operation.Blocking", "IncorrectStatus.TransitRouter", "IncorrectStatus.TransitRouteTable"}) || helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)

	if err != nil {
		if errs.IsExpectedErrors(err, []string{"InstanceNotExist.AggregationRoute"}) || errs.NotFoundError(err) {
			return nil
		}
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
	}

	cenServiceV2 := NewCenServiceV2(client)
	stateConf := helper.BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 5*time.Second, cenServiceV2.CenTransitRouteTableAggregationStateRefreshFunc(d.Id(), "TransitRouteTableAggregationCidr", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errs.WrapErrorf(err, errs.IdMsg, d.Id())
	}

	return nil
}

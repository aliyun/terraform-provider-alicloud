package alicloud

import (
	"fmt"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAlicloudCenTransitRouteTableAggregation() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudCenTransitRouteTableAggregationCreate,
		Read:   resourceAlicloudCenTransitRouteTableAggregationRead,
		Delete: resourceAlicloudCenTransitRouteTableAggregationDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(3 * time.Minute),
			Delete: schema.DefaultTimeout(3 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"transit_route_table_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"transit_route_table_aggregation_cidr": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"transit_route_table_aggregation_scope": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"VPC"}, false),
			},
			"transit_route_table_aggregation_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"transit_route_table_aggregation_description": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAlicloudCenTransitRouteTableAggregationCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cbnService := CbnService{client}
	var response map[string]interface{}
	action := "CreateTransitRouteTableAggregation"
	request := make(map[string]interface{})
	var err error

	request["ClientToken"] = buildClientToken("CreateTransitRouteTableAggregation")
	request["TransitRouteTableId"] = d.Get("transit_route_table_id")
	request["TransitRouteTableAggregationCidr"] = d.Get("transit_route_table_aggregation_cidr")
	request["TransitRouteTableAggregationScope"] = d.Get("transit_route_table_aggregation_scope")

	if v, ok := d.GetOk("transit_route_table_aggregation_name"); ok {
		request["TransitRouteTableAggregationName"] = v
	}

	if v, ok := d.GetOk("transit_route_table_aggregation_description"); ok {
		request["TransitRouteTableAggregationDescription"] = v
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *resource.RetryError {
		response, err = client.RpcPost("Cbn", "2017-09-12", action, nil, request, true)
		if err != nil {
			if IsExpectedErrors(err, []string{"Operation.Blocking"}) || NeedRetry(err) {
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

	d.SetId(fmt.Sprintf("%v:%v", request["TransitRouteTableId"], request["TransitRouteTableAggregationCidr"]))

	stateConf := BuildStateConf([]string{}, []string{"AllConfigured"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, cbnService.CenTransitRouteTableAggregationStateRefreshFunc(d, []string{"ConfigFailed", "PartialConfigured"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAlicloudCenTransitRouteTableAggregationRead(d, meta)
}

func resourceAlicloudCenTransitRouteTableAggregationRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cbnService := CbnService{client}

	object, err := cbnService.DescribeCenTransitRouteTableAggregation(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("transit_route_table_id", object["TrRouteTableId"])
	d.Set("transit_route_table_aggregation_cidr", object["TransitRouteTableAggregationCidr"])
	d.Set("transit_route_table_aggregation_scope", object["Scope"])
	d.Set("transit_route_table_aggregation_name", object["Name"])
	d.Set("transit_route_table_aggregation_description", object["Description"])
	d.Set("status", object["Status"])

	return nil
}

func resourceAlicloudCenTransitRouteTableAggregationDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cbnService := CbnService{client}
	action := "DeleteTransitRouteTableAggregation"
	var response map[string]interface{}

	var err error

	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}

	request := map[string]interface{}{
		"ClientToken":                      buildClientToken("DeleteTransitRouteTableAggregation"),
		"TransitRouteTableId":              parts[0],
		"TransitRouteTableAggregationCidr": parts[1],
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *resource.RetryError {
		response, err = client.RpcPost("Cbn", "2017-09-12", action, nil, request, true)
		if err != nil {
			if IsExpectedErrors(err, []string{"IncorrectStatus.TransitRouteTable", "IncorrectStatus.TransitRouter", "Operation.Blocking"}) || NeedRetry(err) {
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

	stateConf := BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 5*time.Second, cbnService.CenTransitRouteTableAggregationStateRefreshFunc(d, []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return nil
}

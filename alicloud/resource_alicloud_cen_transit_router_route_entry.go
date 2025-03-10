package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudCenTransitRouterRouteEntry() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudCenTransitRouterRouteEntryCreate,
		Read:   resourceAlicloudCenTransitRouterRouteEntryRead,
		Update: resourceAlicloudCenTransitRouterRouteEntryUpdate,
		Delete: resourceAlicloudCenTransitRouterRouteEntryDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(3 * time.Minute),
			Delete: schema.DefaultTimeout(3 * time.Minute),
			Update: schema.DefaultTimeout(3 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"dry_run": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"transit_router_route_entry_description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"transit_router_route_entry_destination_cidr_block": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"transit_router_route_entry_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"transit_router_route_entry_next_hop_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"transit_router_route_entry_next_hop_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"Attachment", "BlackHole"}, false),
			},
			"transit_router_route_table_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"transit_router_route_entry_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAlicloudCenTransitRouterRouteEntryCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cbnService := CbnService{client}
	var response map[string]interface{}
	action := "CreateTransitRouterRouteEntry"
	request := make(map[string]interface{})
	var err error
	if v, ok := d.GetOkExists("dry_run"); ok {
		request["DryRun"] = v
	}

	if v, ok := d.GetOk("transit_router_route_entry_description"); ok {
		request["TransitRouterRouteEntryDescription"] = v
	}

	request["TransitRouterRouteEntryDestinationCidrBlock"] = d.Get("transit_router_route_entry_destination_cidr_block")
	if v, ok := d.GetOk("transit_router_route_entry_name"); ok {
		request["TransitRouterRouteEntryName"] = v
	}

	if v, ok := d.GetOk("transit_router_route_entry_next_hop_id"); ok {
		request["TransitRouterRouteEntryNextHopId"] = v
	}

	request["TransitRouterRouteEntryNextHopType"] = d.Get("transit_router_route_entry_next_hop_type")
	request["TransitRouterRouteTableId"] = d.Get("transit_router_route_table_id")
	request["ClientToken"] = buildClientToken("CreateTransitRouterRouteEntry")
	wait := incrementalWait(3*time.Second, 5*time.Second)
	shortWait := incrementalWait(1*time.Second, 0*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *resource.RetryError {
		response, err = client.RpcPost("Cbn", "2017-09-12", action, nil, request, true)
		if err != nil {
			if IsExpectedErrors(err, []string{"Operation.Blocking", "IncorrectStatus.Status", "InstanceStatus.NotSupport"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			if IsExpectedErrors(err, []string{"OperationFailed.CreateRouteEntryWithSameDestinationCidrBlock"}) {
				shortWait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_cen_transit_router_route_entry", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%v:%v", request["TransitRouterRouteTableId"], response["TransitRouterRouteEntryId"]))
	stateConf := BuildStateConf([]string{}, []string{"Active"}, d.Timeout(schema.TimeoutCreate), 1*time.Second, cbnService.CenTransitRouterRouteEntryStateRefreshFunc(d.Id(), []string{}))
	stateConf.PollInterval = 1 * time.Second
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAlicloudCenTransitRouterRouteEntryRead(d, meta)
}

func resourceAlicloudCenTransitRouterRouteEntryRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cbnService := CbnService{client}
	object, err := cbnService.DescribeCenTransitRouterRouteEntry(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_cen_transit_router_route_entry cbnService.DescribeCenTransitRouterRouteEntry Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	parts, err1 := ParseResourceId(d.Id(), 2)
	if err1 != nil {
		return WrapError(err1)
	}
	d.Set("status", object["TransitRouterRouteEntryStatus"])
	d.Set("transit_router_route_entry_description", object["TransitRouterRouteEntryDescription"])
	d.Set("transit_router_route_entry_destination_cidr_block", object["TransitRouterRouteEntryDestinationCidrBlock"])
	d.Set("transit_router_route_entry_name", object["TransitRouterRouteEntryName"])
	d.Set("transit_router_route_entry_next_hop_id", object["TransitRouterRouteEntryNextHopId"])
	d.Set("transit_router_route_entry_next_hop_type", object["TransitRouterRouteEntryNextHopType"])
	d.Set("transit_router_route_entry_id", object["TransitRouterRouteEntryId"])
	d.Set("transit_router_route_table_id", parts[0])
	return nil
}

func resourceAlicloudCenTransitRouterRouteEntryUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cbnService := CbnService{client}
	var response map[string]interface{}
	var err error
	parts, err1 := ParseResourceId(d.Id(), 2)
	if err1 != nil {
		return WrapError(err1)
	}
	update := false
	request := map[string]interface{}{
		"TransitRouterRouteEntryId": parts[1],
	}
	if d.HasChange("transit_router_route_entry_description") {
		update = true
		request["TransitRouterRouteEntryDescription"] = d.Get("transit_router_route_entry_description")
	}
	if d.HasChange("transit_router_route_entry_name") {
		update = true
		request["TransitRouterRouteEntryName"] = d.Get("transit_router_route_entry_name")
	}
	if update {
		if _, ok := d.GetOkExists("dry_run"); ok {
			request["DryRun"] = d.Get("dry_run")
		}
		action := "UpdateTransitRouterRouteEntry"
		request["ClientToken"] = buildClientToken("UpdateTransitRouterRouteEntry")
		wait := incrementalWait(3*time.Second, 5*time.Second)
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
		stateConf := BuildStateConf([]string{}, []string{"Active"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, cbnService.CenTransitRouterRouteEntryStateRefreshFunc(d.Id(), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}
	return resourceAlicloudCenTransitRouterRouteEntryRead(d, meta)
}

func resourceAlicloudCenTransitRouterRouteEntryDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cbnService := CbnService{client}
	action := "DeleteTransitRouterRouteEntry"
	var response map[string]interface{}
	var err error
	parts, err1 := ParseResourceId(d.Id(), 2)
	if err1 != nil {
		return WrapError(err1)
	}
	request := map[string]interface{}{
		"TransitRouterRouteEntryId": parts[1],
	}

	if v, ok := d.GetOkExists("dry_run"); ok {
		request["DryRun"] = v
	}
	if v, ok := d.GetOk("transit_router_route_entry_next_hop_id"); ok {
		request["TransitRouterRouteEntryNextHopId"] = v
	}
	request["TransitRouterRouteEntryNextHopType"] = d.Get("transit_router_route_entry_next_hop_type")
	request["ClientToken"] = buildClientToken("DeleteTransitRouterRouteEntry")
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *resource.RetryError {
		response, err = client.RpcPost("Cbn", "2017-09-12", action, nil, request, true)
		if err != nil {
			if IsExpectedErrors(err, []string{"Operation.Blocking", "IncorrectStatus.Status", "InstanceStatus.NotSupport"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		if IsExpectedErrors(err, []string{"IllegalParam.TransitRouterRouteEntryId"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	stateConf := BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 1*time.Second, cbnService.CenTransitRouterRouteEntryStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return nil
}

package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudCenTransitRouterRouteTable() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudCenTransitRouterRouteTableCreate,
		Read:   resourceAlicloudCenTransitRouterRouteTableRead,
		Update: resourceAlicloudCenTransitRouterRouteTableUpdate,
		Delete: resourceAlicloudCenTransitRouterRouteTableDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(3 * time.Minute),
			Delete: schema.DefaultTimeout(3 * time.Minute),
			Update: schema.DefaultTimeout(3 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"transit_router_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"transit_router_route_table_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"transit_router_route_table_description": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"tags": tagsSchema(),
			"dry_run": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"transit_router_route_table_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"transit_router_route_table_type": {
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

func resourceAlicloudCenTransitRouterRouteTableCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cbnService := CbnService{client}
	var response map[string]interface{}
	action := "CreateTransitRouterRouteTable"
	request := make(map[string]interface{})
	var err error
	if v, ok := d.GetOkExists("dry_run"); ok {
		request["DryRun"] = v
	}

	request["TransitRouterId"] = d.Get("transit_router_id")
	if v, ok := d.GetOk("transit_router_route_table_description"); ok {
		request["TransitRouterRouteTableDescription"] = v
	}

	if v, ok := d.GetOk("transit_router_route_table_name"); ok {
		request["TransitRouterRouteTableName"] = v
	}

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *resource.RetryError {
		response, err = client.RpcPost("Cbn", "2017-09-12", action, nil, request, false)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_cen_transit_router_route_table", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%v:%v", request["TransitRouterId"], response["TransitRouterRouteTableId"]))
	stateConf := BuildStateConf([]string{}, []string{"Active"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, cbnService.CenTransitRouterRouteTableStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAlicloudCenTransitRouterRouteTableUpdate(d, meta)
}

func resourceAlicloudCenTransitRouterRouteTableRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cbnService := CbnService{client}
	object, err := cbnService.DescribeCenTransitRouterRouteTable(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_cen_transit_router_route_table cbnService.DescribeCenTransitRouterRouteTable Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	parts, err1 := ParseResourceId(d.Id(), 2)
	if err1 != nil {
		return WrapError(err1)
	}

	d.Set("transit_router_id", parts[0])
	d.Set("status", object["TransitRouterRouteTableStatus"])
	d.Set("transit_router_route_table_description", object["TransitRouterRouteTableDescription"])
	d.Set("transit_router_route_table_name", object["TransitRouterRouteTableName"])
	d.Set("transit_router_route_table_id", object["TransitRouterRouteTableId"])
	d.Set("transit_router_route_table_type", object["TransitRouterRouteTableType"])

	listTagResourcesObject, err := cbnService.ListTagResources(d.Id(), "TransitRouterRouteTable")
	if err != nil {
		return WrapError(err)
	}

	d.Set("tags", tagsToMap(listTagResourcesObject))

	return nil
}

func resourceAlicloudCenTransitRouterRouteTableUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cbnService := CbnService{client}
	var err error
	var response map[string]interface{}

	parts, err1 := ParseResourceId(d.Id(), 2)
	if err1 != nil {
		return WrapError(err1)
	}

	update := false
	request := map[string]interface{}{
		"TransitRouterRouteTableId": parts[1],
	}

	if d.HasChange("tags") {
		if err := cbnService.SetResourceTags(d, "TransitRouterRouteTable"); err != nil {
			return WrapError(err)
		}
		d.SetPartial("tags")
	}

	if !d.IsNewResource() && d.HasChange("transit_router_route_table_description") {
		update = true
		request["TransitRouterRouteTableDescription"] = d.Get("transit_router_route_table_description")
	}

	if !d.IsNewResource() && d.HasChange("transit_router_route_table_name") {
		update = true
		request["TransitRouterRouteTableName"] = d.Get("transit_router_route_table_name")
	}

	if update {
		if _, ok := d.GetOkExists("dry_run"); ok {
			request["DryRun"] = d.Get("dry_run")
		}
		action := "UpdateTransitRouterRouteTable"
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			response, err = client.RpcPost("Cbn", "2017-09-12", action, nil, request, false)
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
		stateConf := BuildStateConf([]string{}, []string{"Active"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, cbnService.CenTransitRouterRouteTableStateRefreshFunc(d.Id(), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}
	return resourceAlicloudCenTransitRouterRouteTableRead(d, meta)
}

func resourceAlicloudCenTransitRouterRouteTableDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cbnService := CbnService{client}
	action := "DeleteTransitRouterRouteTable"
	var response map[string]interface{}
	var err error
	parts, err1 := ParseResourceId(d.Id(), 2)
	if err1 != nil {
		return WrapError(err1)
	}
	request := map[string]interface{}{
		"TransitRouterRouteTableId": parts[1],
	}

	if v, ok := d.GetOkExists("dry_run"); ok {
		request["DryRun"] = v
	}
	request["ClientToken"] = buildClientToken("DeleteTransitRouterRouteTable")
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *resource.RetryError {
		response, err = client.RpcPost("Cbn", "2017-09-12", action, nil, request, true)
		if err != nil {
			if IsExpectedErrors(err, []string{"Operation.Blocking", "IncorrectStatus.TransitRouterRouteTable", "IncorrectStatus.Status"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidTransitRouterRouteTableId.NotFound"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	stateConf := BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutUpdate), 1*time.Second, cbnService.CenTransitRouterRouteTableStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return nil
}

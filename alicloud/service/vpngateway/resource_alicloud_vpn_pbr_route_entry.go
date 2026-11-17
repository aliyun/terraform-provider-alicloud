package vpngateway

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/validators"
)

func resourceAlicloudVpnPbrRouteEntry() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudVpnPbrRouteEntryCreate,
		Read:   resourceAlicloudVpnPbrRouteEntryRead,
		Update: resourceAlicloudVpnPbrRouteEntryUpdate,
		Delete: resourceAlicloudVpnPbrRouteEntryDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"next_hop": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"route_dest": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"route_source": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},

			"publish_vpc": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"weight": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntInSlice([]int{0, 100}),
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validators.StringLenBetweenRune(1, 100),
			},
			// priority is create/read/update; the server default is 10. Keep it
			// Optional+Computed so an omitted value follows the value returned by
			// DescribeVpnPbrRouteEntries instead of forcing a local default that could
			// drift for existing route entries or regions that ignore the parameter.
			"priority": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntBetween(1, 100),
			},
			"vpn_gateway_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAlicloudVpnPbrRouteEntryCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vpcService := VpcService{client}
	var response map[string]interface{}
	action := "CreateVpnPbrRouteEntry"
	request := make(map[string]interface{})
	var err error
	request["RegionId"] = client.RegionId
	request["VpnGatewayId"] = d.Get("vpn_gateway_id")
	request["RouteDest"] = d.Get("route_dest")
	request["NextHop"] = d.Get("next_hop")
	request["Weight"] = d.Get("weight")
	request["PublishVpc"] = d.Get("publish_vpc")
	request["RouteSource"] = d.Get("route_source")
	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}
	if v, ok := d.GetOk("priority"); ok {
		request["Priority"] = v
	}

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		request["ClientToken"] = helper.BuildClientToken(action)
		response, err = client.RpcPost("Vpc", "2016-04-28", action, nil, request, true)
		if err != nil {
			if errs.IsExpectedErrors(err, []string{"VpnGateway.Configuring", "TaskConflict", "Appliance.Configuring", "VpnTask.CONFLICT", "VpnConnection.Configuring"}) || helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		logs.AddDebug(action, response, request)
		return nil
	})
	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_vpn_pbr_route_entry", action, errs.AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(request["VpnGatewayId"], ":", request["NextHop"], ":", request["RouteSource"], ":", request["RouteDest"]))

	stateConf := helper.BuildStateConf([]string{}, []string{"published", "normal"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, vpcService.VpnPbrRouteEntryStateRefreshFunc(d.Id(), []string{}, int64(d.Get("weight").(int)), int64(d.Get("priority").(int))))
	if _, err := stateConf.WaitForState(); err != nil {
		return errs.WrapErrorf(err, errs.IdMsg, d.Id())
	}

	return resourceAlicloudVpnPbrRouteEntryRead(d, meta)
}

func vpnPbrRouteEntryIdentityAttribute(value cty.Value, attribute string) (int64, bool) {
	if value == cty.NilVal || value.IsNull() || !value.Type().IsObjectType() || !value.Type().HasAttribute(attribute) {
		return 0, false
	}
	attributeValue := value.GetAttr(attribute)
	if !attributeValue.IsKnown() || attributeValue.IsNull() {
		return 0, false
	}
	result, _ := attributeValue.AsBigFloat().Int64()
	return result, true
}

func vpnPbrRouteEntryReadIdentity(d *schema.ResourceData) []int64 {
	identity := make([]int64, 0, 2)
	rawConfig := d.GetRawConfig()
	rawState := d.GetRawState()
	weight, ok := vpnPbrRouteEntryIdentityAttribute(rawConfig, "weight")
	if !ok {
		weight, ok = vpnPbrRouteEntryIdentityAttribute(rawState, "weight")
	}
	if !ok {
		return identity
	}
	identity = append(identity, weight)
	priority, ok := vpnPbrRouteEntryIdentityAttribute(rawConfig, "priority")
	if !ok {
		priority, ok = vpnPbrRouteEntryIdentityAttribute(rawState, "priority")
	}
	if ok {
		identity = append(identity, priority)
	}
	return identity
}

func resourceAlicloudVpnPbrRouteEntryRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vpcService := VpcService{client}
	object, err := vpcService.DiscoverVpnPbrRouteEntry(d.Id(), vpnPbrRouteEntryReadIdentity(d)...)
	if err != nil {
		if errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_vpn_pbr_route_entry VpcService.DescribeVpnPbrRouteEntry Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}

	d.Set("weight", object["Weight"])
	d.Set("next_hop", object["NextHop"])
	d.Set("route_dest", object["RouteDest"])
	d.Set("vpn_gateway_id", object["VpnInstanceId"])
	d.Set("route_source", object["RouteSource"])
	if object["State"] == "published" {
		d.Set("publish_vpc", true)
	} else {
		d.Set("publish_vpc", false)
	}
	d.Set("status", object["State"])
	d.Set("priority", object["Priority"])
	// description is not returned by DescribeVpnPbrRouteEntries (create-only);
	// keep the configured state value instead of overwriting it with an empty read.
	return nil
}

func resourceAlicloudVpnPbrRouteEntryUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	updateDeadline := time.Now().Add(d.Timeout(schema.TimeoutUpdate))
	var response map[string]interface{}
	d.Partial(true)
	var err error
	parts, err := helper.ParseResourceId(d.Id(), 4)
	if err != nil {
		return errs.WrapError(err)
	}

	request := map[string]interface{}{
		"RegionId":     client.RegionId,
		"RouteDest":    parts[3],
		"RouteSource":  parts[2],
		"NextHop":      parts[1],
		"VpnGatewayId": parts[0],
	}

	update := false
	if d.HasChange("publish_vpc") {
		update = true
		if v, ok := d.GetOkExists("publish_vpc"); ok {
			request["PublishVpc"] = v
		}
	}

	if update {
		request["RouteType"] = "pbr"
		action := "PublishVpnRouteEntry"
		wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
		request["ClientToken"] = helper.BuildClientToken(action)
		err = retry.RetryContext(context.Background(), time.Until(updateDeadline), func() *retry.RetryError {
			response, err = client.RpcPost("Vpc", "2016-04-28", action, nil, request, true)
			if err != nil {
				if errs.IsExpectedErrors(err, []string{"VpnGateway.Configuring", "TaskConflict", "Appliance.Configuring", "VpnTask.CONFLICT", "VpnConnection.Configuring"}) || helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			logs.AddDebug(action, response, request)
			return nil
		})
		if err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_vpn_pbr_route_entry", action, errs.AlibabaCloudSdkGoERROR)
		}
	}

	weightRequest := map[string]interface{}{
		"RegionId":     client.RegionId,
		"RouteDest":    parts[3],
		"RouteSource":  parts[2],
		"NextHop":      parts[1],
		"VpnGatewayId": parts[0],
	}
	update = false
	if d.HasChange("weight") {
		update = true
		oldWeight, newWeight := d.GetChange("weight")
		weightRequest["Weight"] = oldWeight
		weightRequest["NewWeight"] = newWeight
		// ModifyVpnPbrRouteEntryWeight also locates the route entry by its current
		// priority; the old priority is the value still stored on the server when
		// weight and priority change in the same apply.
		if oldPriority, _ := d.GetChange("priority"); oldPriority.(int) > 0 {
			weightRequest["Priority"] = oldPriority
		}
	}

	if update {
		action := "ModifyVpnPbrRouteEntryWeight"
		wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
		weightRequest["ClientToken"] = helper.BuildClientToken(action)
		err = retry.RetryContext(context.Background(), time.Until(updateDeadline), func() *retry.RetryError {
			response, err = client.RpcPost("Vpc", "2016-04-28", action, nil, weightRequest, true)
			if err != nil {
				if errs.IsExpectedErrors(err, []string{"VpnGateway.Configuring", "TaskConflict", "Appliance.Configuring", "VpnTask.CONFLICT", "VpnConnection.Configuring"}) || helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			logs.AddDebug(action, response, weightRequest)
			return nil
		})
		if err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_vpn_pbr_route_entry", action, errs.AlibabaCloudSdkGoERROR)
		}
		// ModifyVpnPbrRouteEntryWeight is asynchronous and the API rejects
		// concurrent modifications on the same gateway, so wait for the gateway
		// to return to active before the priority update below or the final Read.
		if err := waitForActive(client, parts[0], time.Until(updateDeadline)); err != nil {
			return errs.WrapErrorf(err, errs.IdMsg, d.Id())
		}
	}

	priorityRequest := map[string]interface{}{
		"RegionId":     client.RegionId,
		"RouteDest":    parts[3],
		"RouteSource":  parts[2],
		"NextHop":      parts[1],
		"VpnGatewayId": parts[0],
	}
	update = false
	if d.HasChange("priority") {
		update = true
		oldPriority, newPriority := d.GetChange("priority")
		// Priority is optional on the API; omit the previous value when it is unknown
		// (legacy state or a route entry whose priority was never read back).
		if oldPriority.(int) > 0 {
			priorityRequest["Priority"] = oldPriority
		}
		priorityRequest["NewPriority"] = newPriority
		// Weight is part of the resource unique key and is required to locate the
		// route entry; d.Get returns the current (post-weight-change) value.
		priorityRequest["Weight"] = d.Get("weight")
	}

	if update {
		action := "ModifyVpnPbrRouteEntryPriority"
		wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
		priorityRequest["ClientToken"] = helper.BuildClientToken(action)
		err = retry.RetryContext(context.Background(), time.Until(updateDeadline), func() *retry.RetryError {
			response, err = client.RpcPost("Vpc", "2016-04-28", action, nil, priorityRequest, true)
			if err != nil {
				if errs.IsExpectedErrors(err, []string{"VpnGateway.Configuring", "TaskConflict", "Appliance.Configuring", "VpnTask.CONFLICT", "VpnConnection.Configuring"}) || helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			logs.AddDebug(action, response, priorityRequest)
			return nil
		})
		if err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_vpn_pbr_route_entry", action, errs.AlibabaCloudSdkGoERROR)
		}
		// ModifyVpnPbrRouteEntryPriority is asynchronous; reading the route entry
		// before the gateway returns to active can store the previous priority.
		if err := waitForActive(client, parts[0], time.Until(updateDeadline)); err != nil {
			return errs.WrapErrorf(err, errs.IdMsg, d.Id())
		}
	}

	d.Partial(false)
	return resourceAlicloudVpnPbrRouteEntryRead(d, meta)
}

// waitForActive waits until the VPN gateway finishes an asynchronous
// policy route modification. Both ModifyVpnPbrRouteEntryPriority and
// ModifyVpnPbrRouteEntryWeight move the gateway to updating and complete when it
// returns to active; they cannot run concurrently on the same gateway.
func waitForActive(client *connectivity.AliyunClient, vpnGatewayId string, timeout time.Duration) error {
	vpnGatewayServiceV2 := VPNGatewayServiceV2{client}
	stateConf := helper.BuildStateConf([]string{}, []string{"active"}, timeout, 5*time.Second, vpnGatewayServiceV2.VPNGatewayVPNGatewayStateRefreshFunc(vpnGatewayId, "Status", []string{}))
	_, err := stateConf.WaitForStateContext(context.Background())
	return err
}

func resourceAlicloudVpnPbrRouteEntryDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vpcService := VpcService{client}
	var response map[string]interface{}
	var err error
	parts, err := helper.ParseResourceId(d.Id(), 4)
	if err != nil {
		return errs.WrapError(err)
	}

	request := map[string]interface{}{
		"RegionId":     client.RegionId,
		"RouteDest":    parts[3],
		"RouteSource":  parts[2],
		"NextHop":      parts[1],
		"VpnGatewayId": parts[0],
	}

	if v, ok := d.GetOkExists("weight"); ok {
		request["Weight"] = v
	}
	// DeleteVpnPbrRouteEntry locates the entry by its priority as well; pass the
	// state value when the route entry has one, otherwise omit the parameter.
	if v, ok := d.GetOk("priority"); ok {
		request["Priority"] = v
	}

	action := "DeleteVpnPbrRouteEntry"
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	request["ClientToken"] = helper.BuildClientToken(action)
	err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		response, err = client.RpcPost("Vpc", "2016-04-28", action, nil, request, true)
		if err != nil {
			if errs.NotFoundError(err) || errs.IsExpectedErrors(err, []string{"InvalidRouteEntry.NotFound"}) {
				return nil
			}
			if errs.IsExpectedErrors(err, []string{"VpnGateway.Configuring", "VpnRouteEntry.Configuring", "TaskConflict", "Appliance.Configuring", "VpnTask.CONFLICT", "VpnConnection.Configuring"}) || helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		logs.AddDebug(action, response, request)
		return nil
	})
	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_vpn_pbr_route_entry", action, errs.AlibabaCloudSdkGoERROR)
	}

	stateConf := helper.BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 5*time.Second, vpcService.VpnPbrRouteEntryStateRefreshFunc(d.Id(), []string{}, int64(d.Get("weight").(int)), int64(d.Get("priority").(int))))
	if _, err := stateConf.WaitForState(); err != nil {
		return errs.WrapErrorf(err, errs.IdMsg, d.Id())
	}

	return nil
}

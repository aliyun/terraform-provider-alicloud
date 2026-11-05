package vpngateway

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
)

func resourceAlicloudVpnGatewayVcoRoute() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudVpnGatewayVcoRouteCreate,
		Read:   resourceAlicloudVpnGatewayVcoRouteRead,
		Delete: resourceAlicloudVpnGatewayVcoRouteDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(1 * time.Minute),
			Delete: schema.DefaultTimeout(1 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"weight": {
				Type:         schema.TypeInt,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntInSlice([]int{0, 100}),
			},
			"route_dest": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"next_hop": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"vpn_connection_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"overlay_mode": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAlicloudVpnGatewayVcoRouteCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateVcoRouteEntry"
	request := make(map[string]interface{})
	var err error
	request["RegionId"] = client.RegionId
	request["Weight"] = d.Get("weight")
	request["NextHop"] = d.Get("next_hop")
	request["RouteDest"] = d.Get("route_dest")
	request["VpnConnectionId"] = d.Get("vpn_connection_id")
	if v, ok := d.GetOk("overlay_mode"); ok {
		request["OverlayMode"] = v
	}
	request["ClientToken"] = helper.BuildClientToken("CreateVcoRouteEntry")
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *retry.RetryError {
		response, err = client.RpcPost("Vpc", "2016-04-28", action, nil, request, true)
		if err != nil {
			if errs.IsExpectedErrors(err, []string{"VpnGateway.Configuring", "TaskConflict", "Appliance.Configuring", "VpnTask.CONFLICT", "VpnConnection.Configuring"}) || helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)
	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_vpn_gateway_vco_route", action, errs.AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(request["VpnConnectionId"], ":", request["RouteDest"], ":", request["NextHop"], ":", request["Weight"]))

	return resourceAlicloudVpnGatewayVcoRouteRead(d, meta)
}
func resourceAlicloudVpnGatewayVcoRouteRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vpcService := VpcService{client}
	object, err := vpcService.DescribeVpnGatewayVcoRoute(d.Id())
	if err != nil {
		if errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_vpn_gateway_vco_route vpcService.DescribeVpnGatewayVcoRoute Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}
	d.Set("weight", object["Weight"])
	d.Set("status", object["State"])
	d.Set("next_hop", object["NextHop"])
	d.Set("route_dest", object["RouteDest"])
	d.Set("vpn_connection_id", object["VpnConnectionId"])
	return nil
}
func resourceAlicloudVpnGatewayVcoRouteDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vpcService := VpcService{client}
	var err error
	parts, err := helper.ParseResourceId(d.Id(), 4)
	if err != nil {
		return errs.WrapError(err)
	}
	action := "DeleteVcoRouteEntry"
	var response map[string]interface{}
	request := map[string]interface{}{}

	request["RegionId"] = client.RegionId
	request["VpnConnectionId"] = parts[0]
	request["RouteDest"] = parts[1]
	request["NextHop"] = parts[2]
	request["Weight"] = parts[3]

	request["ClientToken"] = helper.BuildClientToken("DeleteVcoRouteEntry")
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *retry.RetryError {
		response, err = client.RpcPost("Vpc", "2016-04-28", action, nil, request, true)
		if err != nil {
			if errs.IsExpectedErrors(err, []string{"VpnGateway.Configuring", "TaskConflict", "Appliance.Configuring", "VpnTask.CONFLICT", "VpnConnection.Configuring"}) || helper.NeedRetry(err) {
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
	if err := WaitForVpnGatewayVcoRouteDeleted(d.Timeout(schema.TimeoutDelete), 3*time.Second, vpcService.VpnGatewayVcoRouteStateRefreshFunc(d.Id(), []string{})); err != nil {
		return errs.WrapErrorf(err, errs.IdMsg, d.Id())
	}
	return nil
}

func WaitForVpnGatewayVcoRouteDeleted(timeout, pollInterval time.Duration, refresh retry.StateRefreshFunc) error {
	stableRefresh := func() (interface{}, string, error) {
		object, _, err := refresh()
		if err != nil || object == nil {
			return object, "", err
		}
		return object, "present", nil
	}
	stateConf := helper.BuildStateConf([]string{"present"}, []string{}, timeout, 0, stableRefresh)
	stateConf.PollInterval = pollInterval
	stateConf.ContinuousTargetOccurence = 3
	_, err := stateConf.WaitForState()
	return err
}

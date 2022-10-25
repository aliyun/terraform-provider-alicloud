package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
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
		},
	}
}

func resourceAlicloudVpnGatewayVcoRouteCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateVcoRouteEntry"
	request := make(map[string]interface{})
	conn, err := client.NewVpcClient()
	if err != nil {
		return WrapError(err)
	}
	request["RegionId"] = client.RegionId
	request["Weight"] = d.Get("weight")
	request["NextHop"] = d.Get("next_hop")
	request["RouteDest"] = d.Get("route_dest")
	request["VpnConnectionId"] = d.Get("vpn_connection_id")
	request["ClientToken"] = buildClientToken("CreateVcoRouteEntry")
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, request, &runtime)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_vpn_gateway_vco_route", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(request["VpnConnectionId"], ":", request["RouteDest"], ":", request["NextHop"], ":", request["Weight"]))

	return resourceAlicloudVpnGatewayVcoRouteRead(d, meta)
}
func resourceAlicloudVpnGatewayVcoRouteRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vpcService := VpcService{client}
	object, err := vpcService.DescribeVpnGatewayVcoRoute(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_vpn_gateway_vco_route vpcService.DescribeVpnGatewayVcoRoute Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
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
	conn, err := client.NewVpcClient()
	if err != nil {
		return WrapError(err)
	}
	parts, err := ParseResourceId(d.Id(), 4)
	if err != nil {
		return WrapError(err)
	}
	action := "DeleteVcoRouteEntry"
	var response map[string]interface{}
	request := map[string]interface{}{}

	request["RegionId"] = client.RegionId
	request["VpnConnectionId"] = parts[0]
	request["RouteDest"] = parts[1]
	request["NextHop"] = parts[2]
	request["Weight"] = parts[3]

	request["ClientToken"] = buildClientToken("DeleteVcoRouteEntry")
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, request, &runtime)
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
	return nil
}

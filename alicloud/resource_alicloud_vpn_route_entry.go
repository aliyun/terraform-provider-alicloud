package alicloud

import (
	"fmt"
	"log"
	"time"

	util "github.com/alibabacloud-go/tea-utils/service"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliyunVpnRouteEntry() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliyunVpnRouteEntryCreate,
		Read:   resourceAliyunVpnRouteEntryRead,
		Update: resourceAliyunVpnRouteEntryUpdate,
		Delete: resourceAliyunVpnRouteEntryDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"vpn_gateway_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

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

			"weight": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntInSlice([]int{0, 100}),
			},

			"publish_vpc": {
				Type:     schema.TypeBool,
				Required: true,
			},

			"route_entry_type": {
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

func resourceAliyunVpnRouteEntryCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vpcService := VpcService{client}
	var response map[string]interface{}
	action := "CreateVpnRouteEntry"
	request := make(map[string]interface{})
	conn, err := client.NewVpcClient()
	if err != nil {
		return WrapError(err)
	}

	request["RegionId"] = client.RegionId

	request["RegionId"] = client.RegionId
	request["VpnGatewayId"] = d.Get("vpn_gateway_id")
	request["RouteDest"] = d.Get("route_dest")
	request["NextHop"] = d.Get("next_hop")
	request["Weight"] = d.Get("weight")
	request["PublishVpc"] = d.Get("publish_vpc")

	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		request["ClientToken"] = buildClientToken(action)
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			if IsExpectedErrors(err, []string{"VpnGateway.Configuring"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_vpn_route_entry", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["VpnInstanceId"], ":", response["NextHop"], ":", response["RouteDest"]))

	stateConf := BuildStateConf([]string{}, []string{"published", "normal"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, vpcService.VpnRouteEntryStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAliyunVpnRouteEntryRead(d, meta)
}

func resourceAliyunVpnRouteEntryRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vpcService := VpcService{client}
	object, err := vpcService.DescribeVpnRouteEntry(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_vpn_route_entry VpcService.DescribeVpnRouteEntry Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("weight", object["Weight"])
	d.Set("next_hop", object["NextHop"])
	d.Set("route_dest", object["RouteDest"])
	d.Set("vpn_gateway_id", object["VpnInstanceId"])

	if object["State"] == "published" {
		d.Set("publish_vpc", true)
	} else {
		d.Set("publish_vpc", false)
	}

	d.Set("status", object["State"])
	d.Set("route_entry_type", object["RouteEntryType"])
	return nil
}

func resourceAliyunVpnRouteEntryUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vpcService := VpcService{client}
	var response map[string]interface{}
	d.Partial(true)
	conn, err := client.NewVpcClient()
	if err != nil {
		return WrapError(err)
	}
	parts, err := ParseResourceId(d.Id(), 3)
	if err != nil {
		return WrapError(err)
	}

	request := map[string]interface{}{
		"RegionId":     client.RegionId,
		"RouteDest":    parts[2],
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
		request["RouteType"] = "dbr"
		action := "PublishVpnRouteEntry"
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 5*time.Second)
		request["ClientToken"] = buildClientToken(action)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, request, &runtime)
			if err != nil {
				if IsExpectedErrors(err, []string{"VpnGateway.Configuring"}) || NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(action, response, request)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "alicloud_vpn_route_entry", action, AlibabaCloudSdkGoERROR)
		}

		status := "normal"
		if request["PublishVpc"].(bool) {
			status = "published"
		}

		stateConf := BuildStateConf([]string{}, []string{status}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, vpcService.VpnRouteEntryStateRefreshFunc(d.Id(), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
		d.SetPartial("public_vpc")
	}

	weightRequest := map[string]interface{}{
		"RegionId":     client.RegionId,
		"RouteDest":    parts[2],
		"NextHop":      parts[1],
		"VpnGatewayId": parts[0],
	}

	update = false
	if d.HasChange("weight") {
		update = true
		oldWeight, newWeight := d.GetChange("weight")
		weightRequest["Weight"] = oldWeight
		weightRequest["NewWeight"] = newWeight
	}

	if update {
		action := "ModifyVpnRouteEntryWeight"
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 5*time.Second)
		request["ClientToken"] = buildClientToken(action)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, weightRequest, &runtime)
			if err != nil {
				if IsExpectedErrors(err, []string{"VpnGateway.Configuring"}) || NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(action, response, weightRequest)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "alicloud_vpn_route_entry", action, AlibabaCloudSdkGoERROR)
		}
		d.SetPartial("weight")
	}

	d.Partial(false)
	return resourceAliyunVpnRouteEntryRead(d, meta)
}

func resourceAliyunVpnRouteEntryDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vpcService := VpcService{client}
	var response map[string]interface{}
	conn, err := client.NewVpcClient()
	if err != nil {
		return WrapError(err)
	}
	parts, err := ParseResourceId(d.Id(), 3)
	if err != nil {
		return WrapError(err)
	}

	request := map[string]interface{}{
		"RegionId":     client.RegionId,
		"RouteDest":    parts[2],
		"NextHop":      parts[1],
		"VpnGatewayId": parts[0],
	}

	if v, ok := d.GetOkExists("weight"); ok {
		request["Weight"] = v
	}

	action := "DeleteVpnRouteEntry"
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	request["ClientToken"] = buildClientToken(action)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			if IsExpectedErrors(err, []string{"VpnGateway.Configuring"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_vpn_route_entry", action, AlibabaCloudSdkGoERROR)
	}

	stateConf := BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 5*time.Second, vpcService.VpnRouteEntryStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return nil
}

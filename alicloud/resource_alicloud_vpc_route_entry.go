// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudVpcRouteEntry() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudVpcRouteEntryCreate,
		Read:   resourceAliCloudVpcRouteEntryRead,
		Update: resourceAliCloudVpcRouteEntryUpdate,
		Delete: resourceAliCloudVpcRouteEntryDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"destination_cidrblock": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"nexthop_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"nexthop_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringInSlice([]string{"HaVip", "RouterInterface", "NetworkInterface", "VpnGateway", "IPv6Gateway", "NatGateway", "Attachment", "VpcPeer", "Instance"}, false),
			},
			"route_entry_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"route_entry_name": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"name"},
			},
			"route_table_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:       schema.TypeString,
				Optional:   true,
				Computed:   true,
				Deprecated: "Field 'name' has been deprecated since provider version 1.211.0. New field 'route_entry_name' instead.",
			},
		},
	}
}

func resourceAliCloudVpcRouteEntryCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateRouteEntry"
	var request map[string]interface{}
	var response map[string]interface{}
	conn, err := client.NewVpcClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	request["RouteTableId"] = d.Get("route_table_id")
	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken(action)

	if v, ok := d.GetOk("name"); ok {
		request["RouteEntryName"] = v
	}

	if v, ok := d.GetOk("route_entry_name"); ok {
		request["RouteEntryName"] = v
	}
	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}
	if v, ok := d.GetOk("destination_cidrblock"); ok {
		request["DestinationCidrBlock"] = v
	}
	if v, ok := d.GetOk("nexthop_id"); ok {
		request["NextHopId"] = v
	}
	if v, ok := d.GetOk("nexthop_type"); ok {
		request["NextHopType"] = v
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, request, &runtime)
		request["ClientToken"] = buildClientToken(action)

		if err != nil {
			if IsExpectedErrors(err, []string{"IncorrectInstanceStatus", "IncorrectRouteEntryStatus", "SystemBusy", "InvalidVBRStatus", "LastTokenProcessing", "IncorrectStatus.Ipv6Address", "IncorrectStatus", "OperationFailed.DistibuteLock", "ServiceUnavailable", "IncorrectStatus.RouteTableStatus", "IncorrectStatus.MultiScopeRiRouteEntry", "IncorrectVpcStatus", "IncorrectHaVipStatus", "OperationConflict", "TaskConflict", "IncorrectStatus.Ipv4Gateway", "IncorrectStatus.VpcPeer", "IncorrectVSwitchStatus"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_vpc_route_entry", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%v:%v", request["RouteTableId"], response["RouteEntryId"]))

	return resourceAliCloudVpcRouteEntryRead(d, meta)
}

func resourceAliCloudVpcRouteEntryRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vpcServiceV2 := VpcServiceV2{client}

	objectRaw, err := vpcServiceV2.DescribeVpcRouteEntry(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_vpc_route_entry DescribeVpcRouteEntry Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("description", objectRaw["Description"])
	d.Set("destination_cidrblock", objectRaw["DestinationCidrBlock"])
	d.Set("route_entry_name", objectRaw["RouteEntryName"])
	d.Set("status", objectRaw["Status"])
	d.Set("route_entry_id", objectRaw["RouteEntryId"])
	d.Set("route_table_id", objectRaw["RouteTableId"])
	nextHop1RawObj, _ := jsonpath.Get("$.NextHops.NextHop[*]", objectRaw)
	nextHop1Raw := make([]interface{}, 0)
	if nextHop1RawObj != nil {
		nextHop1Raw = nextHop1RawObj.([]interface{})
	}
	nextHopChild1Raw := nextHop1Raw[0].(map[string]interface{})
	d.Set("nexthop_id", nextHopChild1Raw["NextHopId"])
	d.Set("nexthop_type", nextHopChild1Raw["NextHopType"])

	d.Set("name", d.Get("route_entry_name"))
	return nil
}

func resourceAliCloudVpcRouteEntryUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	update := false
	parts := strings.Split(d.Id(), ":")
	action := "ModifyRouteEntry"
	conn, err := client.NewVpcClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	request["RouteEntryId"] = parts[1]
	request["RegionId"] = client.RegionId
	if d.HasChange("name") {
		update = true
		request["RouteEntryName"] = d.Get("name")
	}

	if d.HasChange("route_entry_name") {
		update = true
		request["RouteEntryName"] = d.Get("route_entry_name")
	}

	if d.HasChange("description") {
		update = true
		request["Description"] = d.Get("description")
	}

	if d.HasChange("nexthop_type") {
		update = true
		request["NewNextHopType"] = d.Get("nexthop_type")
	}

	if d.HasChange("nexthop_id") {
		update = true
		request["NewNextHopId"] = d.Get("nexthop_id")
	}

	if update {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, request, &runtime)

			if err != nil {
				if IsExpectedErrors(err, []string{"IncorrectInstanceStatus", "IncorrectRouteEntryStatus", "SystemBusy", "InvalidVBRStatus", "LastTokenProcessing", "IncorrectStatus.Ipv6Address", "IncorrectStatus", "OperationFailed.DistibuteLock", "ServiceUnavailable", "IncorrectStatus.RouteTableStatus", "IncorrectStatus.MultiScopeRiRouteEntry", "IncorrectVpcStatus", "IncorrectHaVipStatus", "OperationConflict", "TaskConflict", "IncorrectStatus.Ipv4Gateway", "IncorrectStatus.VpcPeer"}) || NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(action, response, request)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
	}

	return resourceAliCloudVpcRouteEntryRead(d, meta)
}

func resourceAliCloudVpcRouteEntryDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	parts := strings.Split(d.Id(), ":")
	action := "DeleteRouteEntry"
	var request map[string]interface{}
	var response map[string]interface{}
	conn, err := client.NewVpcClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	request["RouteTableId"] = parts[0]
	request["RegionId"] = client.RegionId

	if v, ok := d.GetOk("destination_cidrblock"); ok {
		request["DestinationCidrBlock"] = v
	}
	if v, ok := d.GetOk("nexthop_id"); ok {
		request["NextHopId"] = v
	}
	if v, ok := d.GetOk("next_hops"); ok {
		nextHopListMaps := make([]map[string]interface{}, 0)
		for _, dataLoop := range v.([]interface{}) {
			dataLoopTmp := dataLoop.(map[string]interface{})
			dataLoopMap := make(map[string]interface{})
			dataLoopMap["NextHopId"] = dataLoopTmp["nexthop_id"]
			dataLoopMap["NextHopType"] = dataLoopTmp["nexthop_type"]
			nextHopListMaps = append(nextHopListMaps, dataLoopMap)
		}
		request["NextHopList"] = nextHopListMaps
	}

	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, request, &runtime)

		if err != nil {
			if IsExpectedErrors(err, []string{"IncorrectInstanceStatus", "IncorrectRouteEntryStatus", "SystemBusy", "InvalidVBRStatus", "LastTokenProcessing", "IncorrectStatus.Ipv6Address", "IncorrectStatus", "OperationFailed.DistibuteLock", "ServiceUnavailable", "IncorrectStatus.RouteTableStatus", "IncorrectStatus.MultiScopeRiRouteEntry", "IncorrectVpcStatus", "IncorrectHaVipStatus", "OperationConflict", "TaskConflict", "IncorrectStatus.Ipv4Gateway", "IncorrectStatus.VpcPeer"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})

	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidRouteTable.NotFound"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return nil
}

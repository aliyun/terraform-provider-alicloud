package alicloud

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/tidwall/sjson"
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
			"destination_cidr_block": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"nexthop_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"nexthop_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: StringInSlice([]string{"HaVip", "RouterInterface", "NetworkInterface", "VpnGateway", "IPv6Gateway", "NatGateway", "Attachment", "VpcPeer", "Ipv4Gateway", "GatewayEndpoint", "Ecr", "GatewayLoadBalancerEndpoint", "Instance"}, false),
			},
			"next_hops": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"nexthop_type": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							ForceNew: true,
						},
						"next_hop_region_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"nexthop_id": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							ForceNew: true,
						},
						"enabled": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"next_hop_related_info": {
							Type:     schema.TypeList,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"instance_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"region_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"instance_type": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"weight": {
							Type:     schema.TypeInt,
							Optional: true,
							ForceNew: true,
						},
					},
				},
			},
			"route_entry_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"route_publish_targets": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"target_type": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: StringInSlice([]string{"ECR"}, false),
						},
						"target_instance_id": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"publish_status": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
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
		},
	}
}

func resourceAliCloudVpcRouteEntryCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateRouteEntry"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	if v, ok := d.GetOk("route_table_id"); ok {
		request["RouteTableId"] = v
	}
	if v, ok := d.GetOk("destination_cidr_block"); ok {
		request["DestinationCidrBlock"] = v
	}
	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken(action)

	if v, ok := d.GetOk("name"); ok || d.HasChange("name") {
		request["RouteEntryName"] = v
	}

	if v, ok := d.GetOk("route_entry_name"); ok {
		request["RouteEntryName"] = v
	}
	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}
	if v, ok := d.GetOk("next_hops"); ok {
		nextHopListMapsArray := make([]interface{}, 0)
		for _, dataLoop := range v.([]interface{}) {
			dataLoopTmp := dataLoop.(map[string]interface{})
			dataLoopMap := make(map[string]interface{})
			dataLoopMap["NextHopId"] = dataLoopTmp["nexthop_id"]
			dataLoopMap["Weight"] = dataLoopTmp["weight"]
			dataLoopMap["NextHopType"] = dataLoopTmp["nexthop_type"]
			nextHopListMapsArray = append(nextHopListMapsArray, dataLoopMap)
		}
		request["NextHopList"] = nextHopListMapsArray
	}

	if v, ok := d.GetOk("nexthop_id"); ok {
		request["NextHopId"] = v
	}
	if v, ok := d.GetOk("nexthop_type"); ok {
		request["NextHopType"] = v
	}
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("Vpc", "2016-04-28", action, query, request, true)
		if err != nil {
			if IsExpectedErrors(err, []string{"IncorrectInstanceStatus", "IncorrectRouteEntryStatus", "SystemBusy", "InvalidVBRStatus", "LastTokenProcessing", "IncorrectStatus.Ipv6Address", "IncorrectStatus", "OperationFailed.DistibuteLock", "ServiceUnavailable", "IncorrectStatus.RouteTableStatus", "IncorrectStatus.MultiScopeRiRouteEntry", "IncorrectVpcStatus", "IncorrectHaVipStatus", "OperationConflict", "TaskConflict", "IncorrectStatus.Ipv4Gateway", "IncorrectStatus.VpcPeer", "IncorrectVSwitchStatus"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_vpc_route_entry", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%v:%v", request["RouteTableId"], request["DestinationCidrBlock"]))

	return resourceAliCloudVpcRouteEntryUpdate(d, meta)
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
	d.Set("route_entry_name", objectRaw["RouteEntryName"])
	d.Set("status", objectRaw["Status"])
	d.Set("destination_cidr_block", objectRaw["DestinationCidrBlock"])
	d.Set("route_table_id", objectRaw["RouteTableId"])

	nextHopsRaw, _ := jsonpath.Get("$.NextHops.NextHop[0]", objectRaw)
	d.Set("nexthop_id", nextHopsRaw.(map[string]interface{})["NextHopId"])
	d.Set("nexthop_type", nextHopsRaw.(map[string]interface{})["NextHopType"])

	nextHopRaw, _ := jsonpath.Get("$.NextHops.NextHop", objectRaw)
	nextHopsMaps := make([]map[string]interface{}, 0)
	if nextHopRaw != nil {
		for _, nextHopChildRaw := range nextHopRaw.([]interface{}) {
			nextHopsMap := make(map[string]interface{})
			nextHopChildRaw := nextHopChildRaw.(map[string]interface{})
			nextHopsMap["enabled"] = nextHopChildRaw["Enabled"]
			nextHopsMap["nexthop_id"] = nextHopChildRaw["NextHopId"]
			nextHopsMap["next_hop_region_id"] = nextHopChildRaw["NextHopRegionId"]
			nextHopsMap["nexthop_type"] = nextHopChildRaw["NextHopType"]
			nextHopsMap["weight"] = nextHopChildRaw["Weight"]

			nextHopRelatedInfoMaps := make([]map[string]interface{}, 0)
			nextHopRelatedInfoMap := make(map[string]interface{})
			nextHopRelatedInfoRawObj, _ := jsonpath.Get("$.NextHopRelatedInfo", nextHopChildRaw)
			nextHopRelatedInfoRaw := make(map[string]interface{})
			if nextHopRelatedInfoRawObj != nil {
				nextHopRelatedInfoRaw = nextHopRelatedInfoRawObj.(map[string]interface{})
			}
			if len(nextHopRelatedInfoRaw) > 0 {
				nextHopRelatedInfoMap["instance_id"] = nextHopRelatedInfoRaw["InstanceId"]
				nextHopRelatedInfoMap["instance_type"] = nextHopRelatedInfoRaw["InstanceType"]
				nextHopRelatedInfoMap["region_id"] = nextHopRelatedInfoRaw["RegionId"]

				nextHopRelatedInfoMaps = append(nextHopRelatedInfoMaps, nextHopRelatedInfoMap)
			}
			nextHopsMap["next_hop_related_info"] = nextHopRelatedInfoMaps
			nextHopsMaps = append(nextHopsMaps, nextHopsMap)
		}
	}
	if err := d.Set("next_hops", nextHopsMaps); err != nil {
		return err
	}

	objectRaw, err = vpcServiceV2.DescribeRouteEntryListVpcPublishedRouteEntries(d.Id())
	if err != nil && !NotFoundError(err) {
		if !IsExpectedErrors(err, []string{"ResourceNotAssociated.TargetInstance"}) {
			return WrapError(err)
		}
		objectRaw = make(map[string]interface{})
	}

	routePublishTargetsRaw, _ := jsonpath.Get("$.RoutePublishTargets", objectRaw)

	routePublishTargetsMaps := make([]map[string]interface{}, 0)
	if routePublishTargetsRaw != nil {
		for _, routePublishTargetsChildRaw := range routePublishTargetsRaw.([]interface{}) {
			routePublishTargetsMap := make(map[string]interface{})
			routePublishTargetsChildRaw := routePublishTargetsChildRaw.(map[string]interface{})
			if routePublishTargetsChildRaw["PublishStatus"] == "NonPublished" {
				continue
			}
			routePublishTargetsMap["publish_status"] = routePublishTargetsChildRaw["PublishStatus"]
			routePublishTargetsMap["target_instance_id"] = routePublishTargetsChildRaw["PublishTargetInstanceId"]
			routePublishTargetsMap["target_type"] = routePublishTargetsChildRaw["PublishTargetType"]

			routePublishTargetsMaps = append(routePublishTargetsMaps, routePublishTargetsMap)
		}
	}
	if err := d.Set("route_publish_targets", routePublishTargetsMaps); err != nil {
		return err
	}

	d.Set("name", d.Get("route_entry_name"))
	return nil
}

func resourceAliCloudVpcRouteEntryUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false

	var err error
	parts := strings.Split(d.Id(), ":")
	action := "ModifyRouteEntry"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["DestinationCidrBlock"] = parts[1]
	request["RouteTableId"] = parts[0]
	request["RegionId"] = client.RegionId
	if !d.IsNewResource() && d.HasChange("name") {
		update = true
		request["RouteEntryName"] = d.Get("name")
	}

	if !d.IsNewResource() && d.HasChange("route_entry_name") {
		update = true
		request["RouteEntryName"] = d.Get("route_entry_name")
	}

	if !d.IsNewResource() && d.HasChange("description") {
		update = true
		request["Description"] = d.Get("description")
	}

	if !d.IsNewResource() && d.HasChange("nexthop_id") {
		update = true
		request["NewNextHopId"] = d.Get("nexthop_id")
	}

	if !d.IsNewResource() && d.HasChange("nexthop_type") {
		update = true
		request["NewNextHopType"] = d.Get("nexthop_type")
	}

	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("Vpc", "2016-04-28", action, query, request, true)
			if err != nil {
				if IsExpectedErrors(err, []string{"IncorrectInstanceStatus", "IncorrectRouteEntryStatus", "SystemBusy", "InvalidVBRStatus", "LastTokenProcessing", "IncorrectStatus.Ipv6Address", "IncorrectStatus", "OperationFailed.DistibuteLock", "ServiceUnavailable", "IncorrectStatus.RouteTableStatus", "IncorrectStatus.MultiScopeRiRouteEntry", "IncorrectVpcStatus", "IncorrectHaVipStatus", "OperationConflict", "TaskConflict", "IncorrectStatus.Ipv4Gateway", "IncorrectStatus.VpcPeer"}) || NeedRetry(err) {
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
	}

	if d.HasChange("route_publish_targets") {
		oldEntry, newEntry := d.GetChange("route_publish_targets")
		removed := oldEntry
		added := newEntry

		if len(removed.([]interface{})) > 0 {
			routePublishTargets := removed.([]interface{})

			for _, item := range routePublishTargets {
				parts := strings.Split(d.Id(), ":")
				action := "WithdrawVpcPublishedRouteEntries"
				request = make(map[string]interface{})
				query = make(map[string]interface{})
				request["RegionId"] = client.RegionId
				jsonPathResult, err := jsonpath.Get("$.target_type", item)
				if err == nil {
					request["TargetType"] = jsonPathResult
				}

				jsonPathResult1, err := jsonpath.Get("$.target_instance_id", item)
				if err == nil {
					request["TargetInstanceId"] = jsonPathResult1
				}

				jsonString := convertObjectToJsonString(request)
				jsonString, _ = sjson.Set(jsonString, "RouteEntries.0.RouteTableId", parts[0])
				jsonString, _ = sjson.Set(jsonString, "RouteEntries.0.DestinationCidrBlock", parts[1])
				_ = json.Unmarshal([]byte(jsonString), &request)

				wait := incrementalWait(3*time.Second, 5*time.Second)
				err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
					response, err = client.RpcPost("Vpc", "2016-04-28", action, query, request, true)
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

			}
		}

		if len(added.([]interface{})) > 0 {
			routePublishTargets := added.([]interface{})

			for _, item := range routePublishTargets {
				parts := strings.Split(d.Id(), ":")
				action := "PublishVpcRouteEntries"
				request = make(map[string]interface{})
				query = make(map[string]interface{})
				request["RegionId"] = client.RegionId
				jsonPathResult, err := jsonpath.Get("$.target_type", item)
				if err == nil {
					request["TargetType"] = jsonPathResult
				}

				jsonPathResult1, err := jsonpath.Get("$.target_instance_id", item)
				if err == nil {
					request["TargetInstanceId"] = jsonPathResult1
				}

				jsonString := convertObjectToJsonString(request)
				jsonString, _ = sjson.Set(jsonString, "RouteEntries.0.RouteTableId", parts[0])
				jsonString, _ = sjson.Set(jsonString, "RouteEntries.0.DestinationCidrBlock", parts[1])
				_ = json.Unmarshal([]byte(jsonString), &request)

				wait := incrementalWait(3*time.Second, 5*time.Second)
				err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
					response, err = client.RpcPost("Vpc", "2016-04-28", action, query, request, true)
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

			}
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
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["RouteTableId"] = parts[0]
	request["DestinationCidrBlock"] = parts[1]
	request["RegionId"] = client.RegionId

	if v, ok := d.GetOk("next_hops"); !ok || len(v.([]interface{})) == 1 {
		if v, ok := d.GetOk("route_entry_id"); ok {
			request["RouteEntryId"] = v
		}
		if v, ok := d.GetOk("nexthop_id"); ok {
			request["NextHopId"] = v
		}
	}
	if v, ok := d.GetOk("next_hops"); ok && len(v.([]interface{})) > 1 {
		nextHopListMapsArray := make([]interface{}, 0)
		for _, dataLoop := range v.([]interface{}) {
			dataLoopTmp := dataLoop.(map[string]interface{})
			dataLoopMap := make(map[string]interface{})
			dataLoopMap["NextHopId"] = dataLoopTmp["nexthop_id"]
			dataLoopMap["NextHopType"] = dataLoopTmp["nexthop_type"]
			nextHopListMapsArray = append(nextHopListMapsArray, dataLoopMap)
		}
		request["NextHopList"] = nextHopListMapsArray
	}

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("Vpc", "2016-04-28", action, query, request, true)

		if err != nil {
			if IsExpectedErrors(err, []string{"IncorrectInstanceStatus", "IncorrectRouteEntryStatus", "SystemBusy", "InvalidVBRStatus", "LastTokenProcessing", "IncorrectStatus.Ipv6Address", "IncorrectStatus", "OperationFailed.DistibuteLock", "ServiceUnavailable", "IncorrectStatus.RouteTableStatus", "IncorrectStatus.MultiScopeRiRouteEntry", "IncorrectVpcStatus", "IncorrectHaVipStatus", "OperationConflict", "TaskConflict", "IncorrectStatus.Ipv4Gateway", "IncorrectStatus.VpcPeer", "IncorrectStatus.PrefixList"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		if NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	vpcServiceV2 := VpcServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 5*time.Second, vpcServiceV2.VpcRouteEntryStateRefreshFunc(d.Id(), "Status", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return nil
}

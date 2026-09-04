// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudVpcRouteTargetGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudVpcRouteTargetGroupCreate,
		Read:   resourceAliCloudVpcRouteTargetGroupRead,
		Update: resourceAliCloudVpcRouteTargetGroupUpdate,
		Delete: resourceAliCloudVpcRouteTargetGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"config_mode": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"Active-Standby"}, false),
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"route_target_group_description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"route_target_group_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"route_target_member_list": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"member_id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"member_type": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: StringInSlice([]string{"GatewayLoadBalancerEndpoint"}, false),
						},
						"enable_status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"health_check_status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"weight": {
							Type:     schema.TypeInt,
							Required: true,
						},
					},
				},
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": tagsSchema(),
			"vpc_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAliCloudVpcRouteTargetGroupCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateRouteTargetGroup"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken(action)

	if v, ok := d.GetOk("route_target_group_name"); ok {
		request["RouteTargetGroupName"] = v
	}
	if v, ok := d.GetOk("route_target_member_list"); ok {
		routeTargetMemberListMapsArray := make([]interface{}, 0)
		for _, dataLoop := range convertToInterfaceArray(v) {
			dataLoopTmp := dataLoop.(map[string]interface{})
			dataLoopMap := make(map[string]interface{})
			dataLoopMap["Weight"] = dataLoopTmp["weight"]
			dataLoopMap["MemberType"] = dataLoopTmp["member_type"]
			dataLoopMap["MemberId"] = dataLoopTmp["member_id"]
			routeTargetMemberListMapsArray = append(routeTargetMemberListMapsArray, dataLoopMap)
		}
		request["RouteTargetMemberList"] = routeTargetMemberListMapsArray
	}

	if v, ok := d.GetOk("tags"); ok {
		tagsMap := ConvertTags(v.(map[string]interface{}))
		request["Tags"] = tagsMap
	}

	if v, ok := d.GetOk("resource_group_id"); ok {
		request["ResourceGroupId"] = v
	}
	if v, ok := d.GetOk("route_target_group_description"); ok {
		request["RouteTargetGroupDescription"] = v
	}
	request["VpcId"] = d.Get("vpc_id")
	request["ConfigMode"] = d.Get("config_mode")
	wait := incrementalWait(5*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("Vpc", "2016-04-28", action, query, request, true)
		if err != nil {
			// ResourceNotFound.GatewayLoadBalancerEndpoint is retryable: the
			// referenced GWLB endpoint may be Active in the Privatelink service
			// but not yet propagated to the VPC RouteTargetGroup backend's
			// endpoint registry (cross-service eventual consistency). Retry to
			// give the backend time to register the freshly-created endpoint.
			if IsExpectedErrors(err, []string{"TaskConflict", "ResourceNotFound.GatewayLoadBalancerEndpoint"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_vpc_route_target_group", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["RouteTargetGroupId"]))

	vpcServiceV2 := VpcServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{"Available", "Switched", "Unavailable", "Abnormal"}, d.Timeout(schema.TimeoutCreate), 10*time.Second, vpcServiceV2.VpcRouteTargetGroupStateRefreshFunc(d.Id(), "Status", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	// GetRouteTargetGroup does not return Tags (confirmed via ACC tf-debug),
	// so the Tags parameter carried in the CreateRouteTargetGroup request is
	// not observable through the Get response the Read maps. Apply tags
	// through the dedicated tagging API (TagResources, via SetResourceTags) —
	// the same path Update uses — so the Read (which fetches tags via
	// ListTagResources, the strongly-consistent tagging service) observes the
	// configured tags right after create. This mirrors the Create-tag pattern
	// used by alb/alidns/arms resources in this provider.
	if d.HasChange("tags") {
		if err := vpcServiceV2.SetResourceTags(d, "ROUTETARGETGROUP"); err != nil {
			return WrapError(err)
		}
	}

	return resourceAliCloudVpcRouteTargetGroupRead(d, meta)
}

func resourceAliCloudVpcRouteTargetGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vpcServiceV2 := VpcServiceV2{client}

	objectRaw, err := vpcServiceV2.DescribeVpcRouteTargetGroup(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_vpc_route_target_group DescribeVpcRouteTargetGroup Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("config_mode", objectRaw["ConfigMode"])
	d.Set("create_time", objectRaw["CreateTime"])
	d.Set("resource_group_id", objectRaw["ResourceGroupId"])
	d.Set("route_target_group_description", objectRaw["RouteTargetGroupDescription"])
	d.Set("route_target_group_name", objectRaw["RouteTargetGroupName"])
	d.Set("status", objectRaw["Status"])
	d.Set("vpc_id", objectRaw["VpcId"])

	routeTargetMemberListRaw := objectRaw["RouteTargetMemberList"]
	routeTargetMemberListMaps := make([]map[string]interface{}, 0)
	if routeTargetMemberListRaw != nil {
		for _, routeTargetMemberListChildRaw := range convertToInterfaceArray(routeTargetMemberListRaw) {
			routeTargetMemberListMap := make(map[string]interface{})
			routeTargetMemberListChildRaw := routeTargetMemberListChildRaw.(map[string]interface{})
			routeTargetMemberListMap["enable_status"] = routeTargetMemberListChildRaw["EnableStatus"]
			routeTargetMemberListMap["health_check_status"] = routeTargetMemberListChildRaw["HealthCheckStatus"]
			routeTargetMemberListMap["member_id"] = routeTargetMemberListChildRaw["MemberId"]
			routeTargetMemberListMap["member_type"] = routeTargetMemberListChildRaw["MemberType"]
			routeTargetMemberListMap["weight"] = routeTargetMemberListChildRaw["Weight"]

			routeTargetMemberListMaps = append(routeTargetMemberListMaps, routeTargetMemberListMap)
		}
	}
	if err := d.Set("route_target_member_list", routeTargetMemberListMaps); err != nil {
		return err
	}
	// GetRouteTargetGroup does not return Tags (confirmed via ACC tf-debug:
	// the Get response omits the Tags field even after tags are applied, and
	// the Tags parameter sent in CreateRouteTargetGroup is likewise not
	// reflected in subsequent Get responses). Reading tags from
	// objectRaw["Tags"] therefore always yields an empty set and produces a
	// spurious plan diff after create/update-with-tags. Fetch tags via the
	// dedicated tagging API (ListTagResources) — the same tagging service
	// SetResourceTags writes to — which is the reliable tag source for this
	// resource.
	tagsMaps, err := vpcServiceV2.ListTagResources(d.Id(), "ROUTETARGETGROUP")
	if err != nil {
		return WrapError(err)
	}
	d.Set("tags", tagsToMap(tagsMaps))

	return nil
}

func resourceAliCloudVpcRouteTargetGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false

	var err error
	action := "UpdateRouteTargetGroup"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["RouteTargetGroupId"] = d.Id()
	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken(action)
	if d.HasChange("route_target_group_name") {
		update = true
		request["RouteTargetGroupName"] = d.Get("route_target_group_name")
	}

	// Only send RouteTargetMemberList when the member configuration
	// actually changed. UpdateRouteTargetGroup forbids updating an
	// already-enabled member (OperationDenied.UpdateEnableRouteTargetMember)
	// and weight is immutable, so re-sending the unchanged list on a
	// name/description update would needlessly risk rejection. Send the
	// list only when it changed; switching active/standby is a separate
	// operation (SwitchActiveRouteTarget) not covered by this update path.
	if d.HasChange("route_target_member_list") {
		update = true
		routeTargetMemberListMapsArray := make([]interface{}, 0)
		for _, dataLoop := range convertToInterfaceArray(d.Get("route_target_member_list")) {
			dataLoopTmp := dataLoop.(map[string]interface{})
			dataLoopMap := make(map[string]interface{})
			dataLoopMap["Weight"] = dataLoopTmp["weight"]
			dataLoopMap["MemberType"] = dataLoopTmp["member_type"]
			dataLoopMap["MemberId"] = dataLoopTmp["member_id"]
			routeTargetMemberListMapsArray = append(routeTargetMemberListMapsArray, dataLoopMap)
		}
		request["RouteTargetMemberList"] = routeTargetMemberListMapsArray
	}

	if d.HasChange("route_target_group_description") {
		update = true
		request["RouteTargetGroupDescription"] = d.Get("route_target_group_description")
	}

	if update {
		wait := incrementalWait(5*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("Vpc", "2016-04-28", action, query, request, true)
			if err != nil {
				// ResourceNotFound.GatewayLoadBalancerEndpoint is retryable on
				// Update for the same cross-service propagation reason as Create
				// (see CreateRouteTargetGroup): a referenced GWLB endpoint may be
				// Active in Privatelink but not yet registered in the RouteTargetGroup
				// backend. Retry to let the backend catch up.
				if IsExpectedErrors(err, []string{"TaskConflict", "ResourceNotFound.GatewayLoadBalancerEndpoint"}) || NeedRetry(err) {
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
		vpcServiceV2 := VpcServiceV2{client}
		stateConf := BuildStateConf([]string{}, []string{"Available", "Switched", "Unavailable", "Abnormal"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, vpcServiceV2.VpcRouteTargetGroupStateRefreshFunc(d.Id(), "Status", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}
	update = false
	action = "MoveResourceGroup"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["ResourceId"] = d.Id()
	request["RegionId"] = client.RegionId
	if _, ok := d.GetOk("resource_group_id"); ok && d.HasChange("resource_group_id") {
		update = true
	}
	request["NewResourceGroupId"] = d.Get("resource_group_id")
	request["ResourceType"] = "ROUTETARGETGROUP"
	if update {
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

	if d.HasChange("tags") {
		vpcServiceV2 := VpcServiceV2{client}
		if err := vpcServiceV2.SetResourceTags(d, "ROUTETARGETGROUP"); err != nil {
			return WrapError(err)
		}
	}
	return resourceAliCloudVpcRouteTargetGroupRead(d, meta)
}

func resourceAliCloudVpcRouteTargetGroupDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	action := "DeleteRouteTargetGroup"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["RouteTargetGroupId"] = d.Id()
	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken(action)

	if v, ok := d.GetOk("tags"); ok {
		tagsMap := ConvertTags(v.(map[string]interface{}))
		request = expandTagsToMap(request, tagsMap)
	}

	wait := incrementalWait(5*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("Vpc", "2016-04-28", action, query, request, true)
		if err != nil {
			if IsExpectedErrors(err, []string{"TaskConflict"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		if IsExpectedErrors(err, []string{"ResourceNotFound.RouteTargetGroup"}) || NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	vpcServiceV2 := VpcServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{""}, d.Timeout(schema.TimeoutDelete), 5*time.Second, vpcServiceV2.VpcRouteTargetGroupStateRefreshFunc(d.Id(), "$.RouteTargetGroupId", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return nil
}

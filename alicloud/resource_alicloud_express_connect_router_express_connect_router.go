// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudExpressConnectRouterExpressConnectRouter() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudExpressConnectRouterExpressConnectRouterCreate,
		Read:   resourceAliCloudExpressConnectRouterExpressConnectRouterRead,
		Update: resourceAliCloudExpressConnectRouterExpressConnectRouterUpdate,
		Delete: resourceAliCloudExpressConnectRouterExpressConnectRouterDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"alibaba_side_asn": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringMatch(regexp.MustCompile("^.{0,256}$"), "Represents the description of the leased line gateway."),
			},
			"ecr_name": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringMatch(regexp.MustCompile("^.{0,128}$"), "Name of the Gateway representing the leased line"),
			},
			"regions": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"transit_mode": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: StringInSlice([]string{"NearBy", "ECMP"}, false),
						},
						"region_id": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: StringMatch(regexp.MustCompile("^[A-Za-z0-9_-]+$"), "Representative region ID"),
						},
					},
				},
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": tagsSchema(),
		},
	}
}

func resourceAliCloudExpressConnectRouterExpressConnectRouterCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateExpressConnectRouter"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})

	request["ClientToken"] = buildClientToken(action)

	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}
	request["AlibabaSideAsn"] = d.Get("alibaba_side_asn")
	if v, ok := d.GetOk("ecr_name"); ok {
		request["Name"] = v
	}
	if v, ok := d.GetOk("resource_group_id"); ok {
		request["ResourceGroupId"] = v
	}
	if v, ok := d.GetOk("tags"); ok {
		tagsMap := ConvertTags(v.(map[string]interface{}))
		request = expandTagsToMap(request, tagsMap)
	}

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("ExpressConnectRouter", "2023-09-01", action, query, request, true)
		if err != nil {
			if IsExpectedErrors(err, []string{"Conflict.Lock"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_express_connect_router_express_connect_router", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["EcrId"]))

	return resourceAliCloudExpressConnectRouterExpressConnectRouterUpdate(d, meta)
}

func resourceAliCloudExpressConnectRouterExpressConnectRouterRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	expressConnectRouterServiceV2 := ExpressConnectRouterServiceV2{client}

	objectRaw, err := expressConnectRouterServiceV2.DescribeExpressConnectRouterExpressConnectRouter(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_express_connect_router_express_connect_router DescribeExpressConnectRouterExpressConnectRouter Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("alibaba_side_asn", objectRaw["AlibabaSideAsn"])
	d.Set("create_time", objectRaw["GmtCreate"])
	d.Set("description", objectRaw["Description"])
	d.Set("ecr_name", objectRaw["Name"])
	d.Set("resource_group_id", objectRaw["ResourceGroupId"])
	d.Set("status", objectRaw["Status"])

	tagsMaps := objectRaw["Tags"]
	d.Set("tags", tagsToMap(tagsMaps))

	objectRaw, err = expressConnectRouterServiceV2.DescribeDescribeExpressConnectRouterInterRegionTransitMode(d.Id())
	if err != nil {
		return WrapError(err)
	}
	interRegionTransitModeList1Raw := objectRaw["InterRegionTransitModeList"]
	regionsMaps := make([]map[string]interface{}, 0)
	if interRegionTransitModeList1Raw != nil {
		for _, interRegionTransitModeListChild1Raw := range interRegionTransitModeList1Raw.([]interface{}) {
			regionsMap := make(map[string]interface{})
			interRegionTransitModeListChild1Raw := interRegionTransitModeListChild1Raw.(map[string]interface{})
			regionsMap["region_id"] = interRegionTransitModeListChild1Raw["RegionId"]
			regionsMap["transit_mode"] = interRegionTransitModeListChild1Raw["Mode"]

			regionsMaps = append(regionsMaps, regionsMap)
		}
	}
	d.Set("regions", regionsMaps)

	objectRaw, err = expressConnectRouterServiceV2.DescribeDescribeExpressConnectRouterRegion(d.Id())
	if err != nil {
		return WrapError(err)
	}

	return nil
}

func resourceAliCloudExpressConnectRouterExpressConnectRouterUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false
	d.Partial(true)
	action := "ModifyExpressConnectRouter"
	var err error
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["EcrId"] = d.Id()
	request["ClientToken"] = buildClientToken(action)
	if !d.IsNewResource() && d.HasChange("ecr_name") {
		update = true
		request["Name"] = d.Get("ecr_name")
	}

	if !d.IsNewResource() && d.HasChange("description") {
		update = true
		request["Description"] = d.Get("description")
	}

	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("ExpressConnectRouter", "2023-09-01", action, query, request, true)
			if err != nil {
				if NeedRetry(err) {
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
	update = false
	action = "EnableExpressConnectRouterRouteEntries"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["EcrId"] = d.Id()
	request["ClientToken"] = buildClientToken(action)
	if d.HasChange("disabled_route_entries") {
		update = true
		jsonPathResult, err := jsonpath.Get("$.destination_cidr_block", d.Get("disabled_route_entries"))
		if err == nil {
			request["DestinationCidrBlock"] = jsonPathResult
		}
	}

	if d.HasChange("disabled_route_entries") {
		update = true
		jsonPathResult1, err := jsonpath.Get("$.nexthop_instance_id", d.Get("disabled_route_entries"))
		if err == nil {
			request["NexthopInstanceId"] = jsonPathResult1
		}
	}

	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("ExpressConnectRouter", "2023-09-01", action, query, request, true)
			if err != nil {
				if NeedRetry(err) {
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
	update = false
	action = "DisableExpressConnectRouterRouteEntries"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["EcrId"] = d.Id()
	request["ClientToken"] = buildClientToken(action)
	if d.HasChange("disabled_route_entries") {
		update = true
		jsonPathResult, err := jsonpath.Get("$.destination_cidr_block", d.Get("disabled_route_entries"))
		if err == nil {
			request["DestinationCidrBlock"] = jsonPathResult
		}
	}

	if d.HasChange("disabled_route_entries") {
		update = true
		jsonPathResult1, err := jsonpath.Get("$.nexthop_instance_id", d.Get("disabled_route_entries"))
		if err == nil {
			request["NexthopInstanceId"] = jsonPathResult1
		}
	}

	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("ExpressConnectRouter", "2023-09-01", action, query, request, true)
			if err != nil {
				if NeedRetry(err) {
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
	update = false
	action = "ModifyExpressConnectRouterInterRegionTransitMode"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["EcrId"] = d.Id()
	request["ClientToken"] = buildClientToken(action)
	if d.HasChange("regions") {
		update = true
		if v, ok := d.GetOk("regions"); ok {
			transitModeListMaps := make([]interface{}, 0)
			for _, dataLoop := range v.([]interface{}) {
				dataLoopTmp := dataLoop.(map[string]interface{})
				dataLoopMap := make(map[string]interface{})
				dataLoopMap["Mode"] = dataLoopTmp["transit_mode"]
				dataLoopMap["RegionId"] = dataLoopTmp["region_id"]
				transitModeListMaps = append(transitModeListMaps, dataLoopMap)
			}
			request["TransitModeList"] = transitModeListMaps
		}
	}

	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("ExpressConnectRouter", "2023-09-01", action, query, request, true)
			if err != nil {
				if NeedRetry(err) {
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
		expressConnectRouterServiceV2 := ExpressConnectRouterServiceV2{client}
		stateConf := BuildStateConf([]string{}, []string{"ACTIVE"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, expressConnectRouterServiceV2.ExpressConnectRouterExpressConnectRouterStateRefreshFunc(d.Id(), "Status", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}
	update = false
	action = "GrantInstanceToExpressConnectRouter"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["EcrId"] = d.Id()
	request["ClientToken"] = buildClientToken(action)
	if d.HasChange("granted_instances") {
		update = true
		jsonPathResult, err := jsonpath.Get("$.node_id", d.Get("granted_instances"))
		if err == nil {
			request["InstanceId"] = jsonPathResult
		}
	}

	if d.HasChange("granted_instances") {
		update = true
		jsonPathResult1, err := jsonpath.Get("$.node_type", d.Get("granted_instances"))
		if err == nil {
			request["InstanceType"] = jsonPathResult1
		}
	}

	if d.HasChange("granted_instances") {
		update = true
		jsonPathResult2, err := jsonpath.Get("$.node_region_id", d.Get("granted_instances"))
		if err == nil {
			request["InstanceRegionId"] = jsonPathResult2
		}
	}

	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("ExpressConnectRouter", "2023-09-01", action, query, request, true)
			if err != nil {
				if IsExpectedErrors(err, []string{"Conflict.Lock"}) || NeedRetry(err) {
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
		expressConnectRouterServiceV2 := ExpressConnectRouterServiceV2{client}
		stateConf := BuildStateConf([]string{}, []string{"ACTIVE"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, expressConnectRouterServiceV2.DescribeAsyncExpressConnectRouterExpressConnectRouterStateRefreshFunc(d, response, "$.EcrGrantedInstanceList[*].Status", []string{}))
		if jobDetail, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id(), jobDetail)
		}
	}
	update = false
	action = "RevokeInstanceFromExpressConnectRouter"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["EcrId"] = d.Id()
	request["ClientToken"] = buildClientToken(action)
	if d.HasChange("granted_instances") {
		update = true
		jsonPathResult, err := jsonpath.Get("$.node_id", d.Get("granted_instances"))
		if err == nil {
			request["InstanceId"] = jsonPathResult
		}
	}

	if d.HasChange("granted_instances") {
		update = true
		jsonPathResult1, err := jsonpath.Get("$.node_type", d.Get("granted_instances"))
		if err == nil {
			request["InstanceType"] = jsonPathResult1
		}
	}

	if d.HasChange("granted_instances") {
		update = true
		jsonPathResult2, err := jsonpath.Get("$.node_region_id", d.Get("granted_instances"))
		if err == nil {
			request["InstanceRegionId"] = jsonPathResult2
		}
	}

	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("ExpressConnectRouter", "2023-09-01", action, query, request, true)
			if err != nil {
				if IsExpectedErrors(err, []string{"Conflict.Lock"}) || NeedRetry(err) {
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
		expressConnectRouterServiceV2 := ExpressConnectRouterServiceV2{client}
		stateConf := BuildStateConf([]string{}, []string{"0"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, expressConnectRouterServiceV2.DescribeAsyncExpressConnectRouterExpressConnectRouterStateRefreshFunc(d, response, "$.TotalCount", []string{}))
		if jobDetail, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id(), jobDetail)
		}
	}
	update = false
	action = "MoveResourceGroup"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["ResourceId"] = d.Id()
	request["ClientToken"] = buildClientToken(action)
	if _, ok := d.GetOk("resource_group_id"); ok && !d.IsNewResource() && d.HasChange("resource_group_id") {
		update = true
		request["NewResourceGroupId"] = d.Get("resource_group_id")
	}

	request["ResourceType"] = "EXPRESSCONNECTROUTER"
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("ExpressConnectRouter", "2023-09-01", action, query, request, true)
			if err != nil {
				if IsExpectedErrors(err, []string{"Conflict.Lock"}) || NeedRetry(err) {
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

	if d.HasChange("tags") {
		expressConnectRouterServiceV2 := ExpressConnectRouterServiceV2{client}
		if err := expressConnectRouterServiceV2.SetResourceTags(d, "EXPRESSCONNECTROUTER"); err != nil {
			return WrapError(err)
		}
	}
	d.Partial(false)
	return resourceAliCloudExpressConnectRouterExpressConnectRouterRead(d, meta)
}

func resourceAliCloudExpressConnectRouterExpressConnectRouterDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	action := "DeleteExpressConnectRouter"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["EcrId"] = d.Id()

	request["ClientToken"] = buildClientToken(action)

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("ExpressConnectRouter", "2023-09-01", action, query, request, true)
		request["ClientToken"] = buildClientToken(action)

		if err != nil {
			if IsExpectedErrors(err, []string{"Conflict.Lock"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})

	if err != nil {
		if IsExpectedErrors(err, []string{"ResourceNotFound.EcrId"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	expressConnectRouterServiceV2 := ExpressConnectRouterServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{""}, d.Timeout(schema.TimeoutCreate), 5*time.Second, expressConnectRouterServiceV2.DescribeAsyncExpressConnectRouterExpressConnectRouterStateRefreshFunc(d, response, "$.Code", []string{}))
	if jobDetail, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id(), jobDetail)
	}
	return nil
}

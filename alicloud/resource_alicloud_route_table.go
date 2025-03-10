// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudVpcRouteTable() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudVpcRouteTableCreate,
		Read:   resourceAliCloudVpcRouteTableRead,
		Update: resourceAliCloudVpcRouteTableUpdate,
		Delete: resourceAliCloudVpcRouteTableDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"associate_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"VSwitch", "Gateway"}, false),
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"route_propagation_enable": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"route_table_name": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"name"},
				Computed:      true,
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
			"name": {
				Type:       schema.TypeString,
				Optional:   true,
				Computed:   true,
				Deprecated: "Field 'name' has been deprecated since provider version 1.119.1. New field 'route_table_name' instead.",
			},
		},
	}
}

func resourceAliCloudVpcRouteTableCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateRouteTable"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken(action)

	if v, ok := d.GetOk("name"); ok || d.HasChange("name") {
		request["RouteTableName"] = v
	}

	if v, ok := d.GetOk("route_table_name"); ok {
		request["RouteTableName"] = v
	}
	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}
	request["VpcId"] = d.Get("vpc_id")
	if v, ok := d.GetOk("associate_type"); ok {
		request["AssociateType"] = v
	}
	if v, ok := d.GetOk("tags"); ok {
		tagsMap := ConvertTags(v.(map[string]interface{}))
		request = expandTagsToMap(request, tagsMap)
	}

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("Vpc", "2016-04-28", action, query, request, true)
		if err != nil {
			if IsExpectedErrors(err, []string{"OperationConflict", "LastTokenProcessing", "IncorrectStatus", "SystemBusy", "ServiceUnavailable"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_route_table", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["RouteTableId"]))

	vpcServiceV2 := VpcServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{"Available"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, vpcServiceV2.VpcRouteTableStateRefreshFunc(d.Id(), "Status", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAliCloudVpcRouteTableUpdate(d, meta)
}

func resourceAliCloudVpcRouteTableRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vpcServiceV2 := VpcServiceV2{client}

	objectRaw, err := vpcServiceV2.DescribeVpcRouteTable(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_route_table DescribeVpcRouteTable Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("associate_type", objectRaw["AssociateType"])
	d.Set("create_time", objectRaw["CreationTime"])
	d.Set("description", objectRaw["Description"])
	d.Set("resource_group_id", objectRaw["ResourceGroupId"])
	d.Set("route_propagation_enable", objectRaw["RoutePropagationEnable"])
	d.Set("route_table_name", objectRaw["RouteTableName"])
	d.Set("status", objectRaw["Status"])
	d.Set("vpc_id", objectRaw["VpcId"])

	tagsMaps, _ := jsonpath.Get("$.Tags.Tag", objectRaw)
	d.Set("tags", tagsToMap(tagsMaps))

	d.Set("name", d.Get("route_table_name"))
	return nil
}

func resourceAliCloudVpcRouteTableUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false

	var err error
	action := "ModifyRouteTableAttributes"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["RouteTableId"] = d.Id()
	request["RegionId"] = client.RegionId
	if !d.IsNewResource() && d.HasChange("name") {
		update = true
		request["RouteTableName"] = d.Get("name")
	}

	if !d.IsNewResource() && d.HasChange("route_table_name") {
		update = true
		request["RouteTableName"] = d.Get("route_table_name")
	}

	if !d.IsNewResource() && d.HasChange("description") {
		update = true
		request["Description"] = d.Get("description")
	}

	vpcServiceV2 := VpcServiceV2{client}
	objectRaw, err := vpcServiceV2.DescribeVpcRouteTable(d.Id())
	if d.HasChange("route_propagation_enable") && (objectRaw["RoutePropagationEnable"] != d.Get("route_propagation_enable")) {
		update = true
		request["RoutePropagationEnable"] = d.Get("route_propagation_enable")
	}

	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("Vpc", "2016-04-28", action, query, request, true)
			if err != nil {
				if IsExpectedErrors(err, []string{"OperationConflict", "LastTokenProcessing", "IncorrectStatus", "SystemBusy", "ServiceUnavailable", "IncorrectStatus.RoutePropagationStatus"}) || NeedRetry(err) {
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
		if err := vpcServiceV2.SetResourceTags(d, "ROUTETABLE"); err != nil {
			return WrapError(err)
		}
	}
	return resourceAliCloudVpcRouteTableRead(d, meta)
}

func resourceAliCloudVpcRouteTableDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	action := "DeleteRouteTable"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["RouteTableId"] = d.Id()
	request["RegionId"] = client.RegionId

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("Vpc", "2016-04-28", action, query, request, true)

		if err != nil {
			if IsExpectedErrors(err, []string{"OperationConflict", "DependencyViolation.RouteEntry", "IncorrectRouteTableStatus", "IncorrectStatus.cbnStatus", "OperationDenied.GatewayAssociated"}) || NeedRetry(err) {
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
	stateConf := BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 10*time.Second, vpcServiceV2.VpcRouteTableStateRefreshFunc(d.Id(), "Status", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return nil
}

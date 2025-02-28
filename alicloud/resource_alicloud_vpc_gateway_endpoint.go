// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAliCloudVpcGatewayEndpoint() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudVpcGatewayEndpointCreate,
		Read:   resourceAliCloudVpcGatewayEndpointRead,
		Update: resourceAliCloudVpcGatewayEndpointUpdate,
		Delete: resourceAliCloudVpcGatewayEndpointDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(8 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"gateway_endpoint_descrption": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"gateway_endpoint_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"policy_document": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.ValidateJsonString,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					equal, _ := compareJsonTemplateAreEquivalent(old, new)
					return equal
				},
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"route_tables": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"service_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
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

func resourceAliCloudVpcGatewayEndpointCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateVpcGatewayEndpoint"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken(action)

	if v, ok := d.GetOk("gateway_endpoint_name"); ok {
		request["EndpointName"] = v
	}
	request["ServiceName"] = d.Get("service_name")
	if v, ok := d.GetOk("gateway_endpoint_descrption"); ok {
		request["EndpointDescription"] = v
	}
	if v, ok := d.GetOk("policy_document"); ok {
		request["PolicyDocument"] = v
	}
	request["VpcId"] = d.Get("vpc_id")
	if v, ok := d.GetOk("resource_group_id"); ok {
		request["ResourceGroupId"] = v
	}
	if v, ok := d.GetOk("tags"); ok {
		tagsMap := ConvertTags(v.(map[string]interface{}))
		request = expandTagsToMap(request, tagsMap)
	}

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("Vpc", "2016-04-28", action, query, request, true)
		if err != nil {
			if IsExpectedErrors(err, []string{"OperationFailed.ConcurrentOperation", "IdempotentParameterMismatch", "IncorrectStatus.Vpc", "OperationConflict", "IncorrectStatus", "ServiceUnavailable", "LastTokenProcessing", "SystemBusy"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_vpc_gateway_endpoint", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["EndpointId"]))

	vpcServiceV2 := VpcServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{"Created"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, vpcServiceV2.VpcGatewayEndpointStateRefreshFunc(d.Id(), "EndpointStatus", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAliCloudVpcGatewayEndpointUpdate(d, meta)
}

func resourceAliCloudVpcGatewayEndpointRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vpcServiceV2 := VpcServiceV2{client}

	objectRaw, err := vpcServiceV2.DescribeVpcGatewayEndpoint(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_vpc_gateway_endpoint DescribeVpcGatewayEndpoint Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("create_time", objectRaw["CreationTime"])
	d.Set("gateway_endpoint_descrption", objectRaw["EndpointDescription"])
	d.Set("gateway_endpoint_name", objectRaw["EndpointName"])
	d.Set("policy_document", objectRaw["PolicyDocument"])
	d.Set("resource_group_id", objectRaw["ResourceGroupId"])
	d.Set("service_name", objectRaw["ServiceName"])
	d.Set("status", objectRaw["EndpointStatus"])
	d.Set("vpc_id", objectRaw["VpcId"])

	routeTablesRaw := make([]interface{}, 0)
	if objectRaw["RouteTables"] != nil {
		routeTablesRaw = objectRaw["RouteTables"].([]interface{})
	}

	d.Set("route_tables", routeTablesRaw)
	tagsMaps := objectRaw["Tags"]
	d.Set("tags", tagsToMap(tagsMaps))

	return nil
}

func resourceAliCloudVpcGatewayEndpointUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false
	d.Partial(true)

	var err error
	action := "UpdateVpcGatewayEndpointAttribute"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["EndpointId"] = d.Id()
	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken(action)
	if !d.IsNewResource() && d.HasChange("gateway_endpoint_name") {
		update = true
		request["EndpointName"] = d.Get("gateway_endpoint_name")
	}

	if !d.IsNewResource() && d.HasChange("gateway_endpoint_descrption") {
		update = true
		request["EndpointDescription"] = d.Get("gateway_endpoint_descrption")
	}

	if !d.IsNewResource() && d.HasChange("policy_document") {
		update = true
		request["PolicyDocument"] = d.Get("policy_document")
	}

	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("Vpc", "2016-04-28", action, query, request, true)
			if err != nil {
				if IsExpectedErrors(err, []string{"OperationConflict", "IncorrectStatus", "ServiceUnavailable", "LastTokenProcessing", "SystemBusy", "OperationFailed.ConcurrentOperation", "IdempotentParameterMismatch", "IncorrectStatus.Vpc"}) || NeedRetry(err) {
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
		stateConf := BuildStateConf([]string{}, []string{"Created"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, vpcServiceV2.VpcGatewayEndpointStateRefreshFunc(d.Id(), "EndpointStatus", []string{}))
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
	if _, ok := d.GetOk("resource_group_id"); ok && !d.IsNewResource() && d.HasChange("resource_group_id") {
		update = true
	}
	request["NewResourceGroupId"] = d.Get("resource_group_id")
	request["ResourceType"] = "GatewayEndpoint"
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

	if d.HasChange("route_tables") {
		oldEntry, newEntry := d.GetChange("route_tables")
		oldEntrySet := oldEntry.(*schema.Set)
		newEntrySet := newEntry.(*schema.Set)
		removed := oldEntrySet.Difference(newEntrySet)
		added := newEntrySet.Difference(oldEntrySet)

		if removed.Len() > 0 {
			action := "DissociateRouteTablesFromVpcGatewayEndpoint"
			request = make(map[string]interface{})
			query = make(map[string]interface{})
			request["EndpointId"] = d.Id()
			request["RegionId"] = client.RegionId
			request["ClientToken"] = buildClientToken(action)
			localData := removed.List()
			routeTableIdsMapsArray := localData
			request["RouteTableIds"] = routeTableIdsMapsArray

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
			vpcServiceV2 := VpcServiceV2{client}
			stateConf := BuildStateConf([]string{}, []string{"Created"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, vpcServiceV2.VpcGatewayEndpointStateRefreshFunc(d.Id(), "EndpointStatus", []string{}))
			if _, err := stateConf.WaitForState(); err != nil {
				return WrapErrorf(err, IdMsg, d.Id())
			}

		}

		if added.Len() > 0 {
			action := "AssociateRouteTablesWithVpcGatewayEndpoint"
			request = make(map[string]interface{})
			query = make(map[string]interface{})
			request["EndpointId"] = d.Id()
			request["RegionId"] = client.RegionId
			request["ClientToken"] = buildClientToken(action)
			localData := added.List()
			routeTableIdsMapsArray := localData
			request["RouteTableIds"] = routeTableIdsMapsArray

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
			vpcServiceV2 := VpcServiceV2{client}
			stateConf := BuildStateConf([]string{}, []string{"Created"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, vpcServiceV2.VpcGatewayEndpointStateRefreshFunc(d.Id(), "EndpointStatus", []string{}))
			if _, err := stateConf.WaitForState(); err != nil {
				return WrapErrorf(err, IdMsg, d.Id())
			}

		}

	}
	if d.HasChange("tags") {
		vpcServiceV2 := VpcServiceV2{client}
		if err := vpcServiceV2.SetResourceTags(d, "GatewayEndpoint"); err != nil {
			return WrapError(err)
		}
	}
	d.Partial(false)
	return resourceAliCloudVpcGatewayEndpointRead(d, meta)
}

func resourceAliCloudVpcGatewayEndpointDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	action := "DeleteVpcGatewayEndpoint"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["EndpointId"] = d.Id()
	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken(action)

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("Vpc", "2016-04-28", action, query, request, true)
		request["ClientToken"] = buildClientToken(action)

		if err != nil {
			if IsExpectedErrors(err, []string{"OperationConflict", "IncorrectStatus", "ServiceUnavailable", "LastTokenProcessing", "SystemBusy", "OperationFailed.ConcurrentOperation", "IdempotentParameterMismatch", "IncorrectStatus.Vpc"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		if IsExpectedErrors(err, []string{"ResourceNotFound.GatewayEndpoint"}) || NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	vpcServiceV2 := VpcServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{""}, d.Timeout(schema.TimeoutDelete), 5*time.Second, vpcServiceV2.VpcGatewayEndpointStateRefreshFunc(d.Id(), "EndpointStatus", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return nil
}

// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudVpcIpv4Gateway() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudVpcIpv4GatewayCreate,
		Read:   resourceAlicloudVpcIpv4GatewayRead,
		Update: resourceAlicloudVpcIpv4GatewayUpdate,
		Delete: resourceAlicloudVpcIpv4GatewayDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"dry_run": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"ipv4_gateway_description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ipv4_gateway_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ipv4_gateway_name": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringMatch(regexp.MustCompile("^[a-zA-Z][a-zA-Z0-9-._]{1,127}$"), "The name of the IPv4 gateway. The name must be 2 to 128 characters in length, and can contain letters, digits, periods (.), underscores (_), and hyphens (-). It must start with a letter."),
			},
			"ipv4_gateway_route_table_id": {
				Type:     schema.TypeString,
				Computed: true,
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
			"vpc_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAlicloudVpcIpv4GatewayCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "CreateIpv4Gateway"
	var request map[string]interface{}
	var response map[string]interface{}
	var err error
	request = make(map[string]interface{})
	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken(action)

	if v, ok := d.GetOk("ipv4_gateway_name"); ok {
		request["Ipv4GatewayName"] = v
	}
	if v, ok := d.GetOk("ipv4_gateway_description"); ok {
		request["Ipv4GatewayDescription"] = v
	}
	request["VpcId"] = d.Get("vpc_id")
	if v, ok := d.GetOk("resource_group_id"); ok {
		request["ResourceGroupId"] = v
	}
	request["DryRun"] = d.Get("dry_run")
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("Vpc", "2016-04-28", action, nil, request, true)
		request["ClientToken"] = buildClientToken(action)

		if err != nil {
			if IsExpectedErrors(err, []string{"IncorrectStatus.Vpc", "OperationConflict", "IncorrectStatus", "ServiceUnavailable", "LastTokenProcessing", "SystemBusy"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_vpc_ipv4_gateway", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["Ipv4GatewayId"]))

	vpcServiceV2 := VpcServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{"Created"}, d.Timeout(schema.TimeoutCreate), 0, vpcServiceV2.VpcIpv4GatewayStateRefreshFunc(d.Id(), "Status", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAlicloudVpcIpv4GatewayUpdate(d, meta)
}

func resourceAlicloudVpcIpv4GatewayRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vpcServiceV2 := VpcServiceV2{client}

	objectRaw, err := vpcServiceV2.DescribeVpcIpv4Gateway(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_vpc_ipv4_gateway DescribeVpcIpv4Gateway Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("create_time", objectRaw["CreateTime"])
	d.Set("enabled", objectRaw["Enabled"])
	d.Set("ipv4_gateway_description", objectRaw["Ipv4GatewayDescription"])
	d.Set("ipv4_gateway_name", objectRaw["Ipv4GatewayName"])
	d.Set("ipv4_gateway_route_table_id", objectRaw["Ipv4GatewayRouteTableId"])
	d.Set("resource_group_id", objectRaw["ResourceGroupId"])
	d.Set("status", objectRaw["Status"])
	d.Set("vpc_id", objectRaw["VpcId"])
	d.Set("ipv4_gateway_id", objectRaw["Ipv4GatewayId"])
	tagsMaps := objectRaw["Tags"]
	d.Set("tags", tagsToMap(tagsMaps))

	return nil
}

func resourceAlicloudVpcIpv4GatewayUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	update := false
	d.Partial(true)
	action := "UpdateIpv4GatewayAttribute"
	var err error
	request = make(map[string]interface{})

	request["Ipv4GatewayId"] = d.Id()
	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken(action)
	if !d.IsNewResource() && d.HasChange("ipv4_gateway_name") {
		update = true
		request["Ipv4GatewayName"] = d.Get("ipv4_gateway_name")
	}

	if !d.IsNewResource() && d.HasChange("ipv4_gateway_description") {
		update = true
		request["Ipv4GatewayDescription"] = d.Get("ipv4_gateway_description")
	}

	request["DryRun"] = d.Get("dry_run")
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("Vpc", "2016-04-28", action, nil, request, true)
			request["ClientToken"] = buildClientToken(action)

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
		d.SetPartial("ipv4_gateway_name")
		d.SetPartial("ipv4_gateway_description")
	}
	update = false
	action = "MoveResourceGroup"
	request = make(map[string]interface{})

	request["ResourceId"] = d.Id()
	request["RegionId"] = client.RegionId
	if !d.IsNewResource() && d.HasChange("resource_group_id") {
		update = true
		request["NewResourceGroupId"] = d.Get("resource_group_id")
	}

	request["ResourceType"] = "IPV4GATEWAY"
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("Vpc", "2016-04-28", action, nil, request, false)

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
		d.SetPartial("resource_group_id")
	}

	if d.HasChange("enabled") {
		client := meta.(*connectivity.AliyunClient)
		vpcServiceV2 := VpcServiceV2{client}
		object, err := vpcServiceV2.DescribeVpcIpv4Gateway(d.Id())
		if err != nil {
			return WrapError(err)
		}

		target := d.Get("enabled").(bool)
		if object["Enabled"].(bool) != target {
			if target == true {
				action = "EnableVpcIpv4Gateway"
				request = make(map[string]interface{})

				request["Ipv4GatewayId"] = d.Id()
				request["RegionId"] = client.RegionId
				request["ClientToken"] = buildClientToken(action)
				request["DryRun"] = d.Get("dry_run")
				wait := incrementalWait(3*time.Second, 5*time.Second)
				err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
					response, err = client.RpcPost("Vpc", "2016-04-28", action, nil, request, true)
					request["ClientToken"] = buildClientToken(action)

					if err != nil {
						if IsExpectedErrors(err, []string{"OperationConflict", "IncorrectStatus", "LastTokenProcessing", "SystemBusy", "OperationFailed.LastTokenProcessing", "ServiceUnavailable"}) || NeedRetry(err) {
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
				vpcServiceV2 := VpcServiceV2{client}
				stateConf := BuildStateConf([]string{}, []string{"Created"}, d.Timeout(schema.TimeoutUpdate), 0, vpcServiceV2.VpcIpv4GatewayStateRefreshFunc(d.Id(), "Status", []string{}))
				if _, err := stateConf.WaitForState(); err != nil {
					return WrapErrorf(err, IdMsg, d.Id())
				}

			}
		}
	}

	update = false
	if d.HasChange("tags") {
		update = true
		vpcServiceV2 := VpcServiceV2{client}
		if err := vpcServiceV2.SetResourceTags(d, "IPV4GATEWAY"); err != nil {
			return WrapError(err)
		}
		d.SetPartial("tags")
	}
	d.Partial(false)
	return resourceAlicloudVpcIpv4GatewayRead(d, meta)
}

func resourceAlicloudVpcIpv4GatewayDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "DeleteIpv4Gateway"
	var request map[string]interface{}
	var response map[string]interface{}
	var err error
	request = make(map[string]interface{})

	request["Ipv4GatewayId"] = d.Id()
	request["RegionId"] = client.RegionId

	request["ClientToken"] = buildClientToken(action)

	request["DryRun"] = d.Get("dry_run")
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("Vpc", "2016-04-28", action, nil, request, true)
		request["ClientToken"] = buildClientToken(action)

		if err != nil {
			if IsExpectedErrors(err, []string{"OperationConflict", "ServiceUnavailable", "SystemBusy", "LastTokenProcessing", "IncorrectStatus.Ipv4Gateway", "IncorrectStatus"}) || NeedRetry(err) {
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

	vpcServiceV2 := VpcServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 0, vpcServiceV2.VpcIpv4GatewayStateRefreshFunc(d.Id(), "Status", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return nil
}

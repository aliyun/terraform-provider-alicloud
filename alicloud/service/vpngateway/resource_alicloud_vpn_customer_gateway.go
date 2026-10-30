// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package vpngateway

import (
	"fmt"
	"log"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/tags"
)

func resourceAliCloudVPNGatewayCustomerGateway() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudVPNGatewayCustomerGatewayCreate,
		Read:   resourceAliCloudVPNGatewayCustomerGatewayRead,
		Update: resourceAliCloudVPNGatewayCustomerGatewayUpdate,
		Delete: resourceAliCloudVPNGatewayCustomerGatewayDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"asn": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"create_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"customer_gateway_name": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"name"},
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ip_address": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"tags": tags.TagsSchema(),
			"name": {
				Type:       schema.TypeString,
				Optional:   true,
				Computed:   true,
				Deprecated: "Field 'name' has been deprecated since provider version 1.210.0. New field 'customer_gateway_name' instead.",
			},
		},
	}
}

func resourceAliCloudVPNGatewayCustomerGatewayCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "CreateCustomerGateway"
	var request map[string]interface{}
	var response map[string]interface{}
	var err error
	request = make(map[string]interface{})
	request["RegionId"] = client.RegionId
	request["ClientToken"] = helper.BuildClientToken(action)

	request["IpAddress"] = d.Get("ip_address")
	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}
	if v, ok := d.GetOk("asn"); ok {
		request["Asn"] = v
	}
	if v, ok := d.GetOk("name"); ok {
		request["Name"] = v
	}

	if v, ok := d.GetOk("customer_gateway_name"); ok {
		request["Name"] = v
	}
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		response, err = client.RpcPost("Vpc", "2016-04-28", action, nil, request, true)
		request["ClientToken"] = helper.BuildClientToken(action)

		if err != nil {
			if errs.IsExpectedErrors(err, []string{"OperationConflict"}) || helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		logs.AddDebug(action, response, request)
		return nil
	})

	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_vpn_customer_gateway", action, errs.AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["CustomerGatewayId"]))

	return resourceAliCloudVPNGatewayCustomerGatewayUpdate(d, meta)
}

func resourceAliCloudVPNGatewayCustomerGatewayRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vPNGatewayServiceV2 := VPNGatewayServiceV2{client}

	objectRaw, err := vPNGatewayServiceV2.DescribeVPNGatewayCustomerGateway(d.Id())
	if err != nil {
		if !d.IsNewResource() && errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_vpn_customer_gateway DescribeVPNGatewayCustomerGateway Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}

	d.Set("asn", objectRaw["Asn"])
	d.Set("create_time", objectRaw["CreateTime"])
	d.Set("customer_gateway_name", objectRaw["Name"])
	d.Set("description", objectRaw["Description"])
	d.Set("ip_address", objectRaw["IpAddress"])
	tagsMaps, _ := jsonpath.Get("$.Tags.Tag", objectRaw)
	d.Set("tags", tags.TagsToMap(tagsMaps))

	d.Set("name", d.Get("customer_gateway_name"))
	return nil
}

func resourceAliCloudVPNGatewayCustomerGatewayUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	update := false
	action := "ModifyCustomerGatewayAttribute"
	var err error
	request = make(map[string]interface{})
	request["CustomerGatewayId"] = d.Id()
	request["RegionId"] = client.RegionId
	request["ClientToken"] = helper.BuildClientToken(action)
	if !d.IsNewResource() && d.HasChange("description") {
		update = true
		request["Description"] = d.Get("description")
	}

	if !d.IsNewResource() && d.HasChange("name") {
		update = true
		request["Name"] = d.Get("name")
	}

	if !d.IsNewResource() && d.HasChange("customer_gateway_name") {
		update = true
		request["Name"] = d.Get("customer_gateway_name")
	}

	if update {
		wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RpcPost("Vpc", "2016-04-28", action, nil, request, true)
			request["ClientToken"] = helper.BuildClientToken(action)

			if err != nil {
				if helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			logs.AddDebug(action, response, request)
			return nil
		})
		if err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
		}
	}

	if d.HasChange("tags") {
		vPNGatewayServiceV2 := VPNGatewayServiceV2{client}
		if err := vPNGatewayServiceV2.SetResourceTags(d, "CUSTOMERGATEWAY"); err != nil {
			return errs.WrapError(err)
		}
	}
	return resourceAliCloudVPNGatewayCustomerGatewayRead(d, meta)
}

func resourceAliCloudVPNGatewayCustomerGatewayDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	action := "DeleteCustomerGateway"
	var request map[string]interface{}
	var response map[string]interface{}
	var err error
	request = make(map[string]interface{})
	request["CustomerGatewayId"] = d.Id()
	request["RegionId"] = client.RegionId

	request["ClientToken"] = helper.BuildClientToken(action)

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		response, err = client.RpcPost("Vpc", "2016-04-28", action, nil, request, true)
		request["ClientToken"] = helper.BuildClientToken(action)

		if err != nil {
			if errs.IsExpectedErrors(err, []string{"VpnGateway.Configuring"}) || helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		logs.AddDebug(action, response, request)
		return nil
	})

	if err != nil {
		if errs.IsExpectedErrors(err, []string{"InvalidCustomerGatewayInstanceId.NotFound"}) {
			return nil
		}
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
	}

	return nil
}

// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudNatGatewayNatIpCidr() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudNatGatewayNatIpCidrCreate,
		Read:   resourceAliCloudNatGatewayNatIpCidrRead,
		Update: resourceAliCloudNatGatewayNatIpCidrUpdate,
		Delete: resourceAliCloudNatGatewayNatIpCidrDelete,
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
			"nat_gateway_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"nat_ip_cidr": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"nat_ip_cidr_description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"nat_ip_cidr_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAliCloudNatGatewayNatIpCidrCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateNatIpCidr"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	if v, ok := d.GetOk("nat_ip_cidr"); ok {
		request["NatIpCidr"] = v
	}
	if v, ok := d.GetOk("nat_gateway_id"); ok {
		request["NatGatewayId"] = v
	}
	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken(action)

	if v, ok := d.GetOk("nat_ip_cidr_description"); ok {
		request["NatIpCidrDescription"] = v
	}
	if v, ok := d.GetOkExists("dry_run"); ok {
		request["DryRun"] = v
	}
	request["NatIpCidrName"] = d.Get("nat_ip_cidr_name")
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("Vpc", "2016-04-28", action, query, request, true)
		if err != nil {
			if IsExpectedErrors(err, []string{"OperationConflict"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_vpc_nat_ip_cidr", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%v:%v", request["NatGatewayId"], request["NatIpCidr"]))

	return resourceAliCloudNatGatewayNatIpCidrRead(d, meta)
}

func resourceAliCloudNatGatewayNatIpCidrRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	nATGatewayServiceV2 := NATGatewayServiceV2{client}

	objectRaw, err := nATGatewayServiceV2.DescribeNatGatewayNatIpCidr(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_vpc_nat_ip_cidr DescribeNatGatewayNatIpCidr Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("create_time", objectRaw["CreationTime"])
	d.Set("nat_ip_cidr_description", objectRaw["NatIpCidrDescription"])
	d.Set("nat_ip_cidr_name", objectRaw["NatIpCidrName"])
	d.Set("status", objectRaw["NatIpCidrStatus"])
	d.Set("nat_gateway_id", objectRaw["NatGatewayId"])
	d.Set("nat_ip_cidr", objectRaw["NatIpCidr"])

	return nil
}

func resourceAliCloudNatGatewayNatIpCidrUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false

	var err error
	parts := strings.Split(d.Id(), ":")
	action := "ModifyNatIpCidrAttribute"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["NatIpCidr"] = parts[1]
	request["NatGatewayId"] = parts[0]
	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken(action)
	if d.HasChange("nat_ip_cidr_description") {
		update = true
		request["NatIpCidrDescription"] = d.Get("nat_ip_cidr_description")
	}

	if v, ok := d.GetOkExists("dry_run"); ok {
		request["DryRun"] = v
	}
	if d.HasChange("nat_ip_cidr_name") {
		update = true
	}
	request["NatIpCidrName"] = d.Get("nat_ip_cidr_name")
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

	return resourceAliCloudNatGatewayNatIpCidrRead(d, meta)
}

func resourceAliCloudNatGatewayNatIpCidrDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	parts := strings.Split(d.Id(), ":")
	action := "DeleteNatIpCidr"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["NatIpCidr"] = parts[1]
	request["NatGatewayId"] = parts[0]
	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken(action)

	if v, ok := d.GetOkExists("dry_run"); ok {
		request["DryRun"] = v
	}
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("Vpc", "2016-04-28", action, query, request, true)
		if err != nil {
			if IsExpectedErrors(err, []string{"OperationConflict"}) || NeedRetry(err) {
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

	return nil
}

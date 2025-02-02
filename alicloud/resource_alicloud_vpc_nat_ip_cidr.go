package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudVpcNatIpCidr() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudVpcNatIpCidrCreate,
		Read:   resourceAlicloudVpcNatIpCidrRead,
		Update: resourceAlicloudVpcNatIpCidrUpdate,
		Delete: resourceAlicloudVpcNatIpCidrDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"dry_run": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"nat_gateway_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"nat_ip_cidr": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"nat_ip_cidr_description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"nat_ip_cidr_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAlicloudVpcNatIpCidrCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateNatIpCidr"
	request := make(map[string]interface{})
	var err error
	if v, ok := d.GetOkExists("dry_run"); ok {
		request["DryRun"] = v
	}
	request["NatGatewayId"] = d.Get("nat_gateway_id")
	if v, ok := d.GetOk("nat_ip_cidr"); ok {
		request["NatIpCidr"] = v
	}
	if v, ok := d.GetOk("nat_ip_cidr_description"); ok {
		request["NatIpCidrDescription"] = v
	}
	if v, ok := d.GetOk("nat_ip_cidr_name"); ok {
		request["NatIpCidrName"] = v
	}
	request["RegionId"] = client.RegionId
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		request["ClientToken"] = buildClientToken("CreateNatIpCidr")
		response, err = client.RpcPost("Vpc", "2016-04-28", action, nil, request, true)
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

	d.SetId(fmt.Sprint(request["NatGatewayId"], ":", request["NatIpCidr"]))

	return resourceAlicloudVpcNatIpCidrRead(d, meta)
}
func resourceAlicloudVpcNatIpCidrRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vpcService := VpcService{client}
	object, err := vpcService.DescribeVpcNatIpCidr(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_vpc_nat_ip_cidr vpcService.DescribeVpcNatIpCidr Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	d.Set("nat_gateway_id", parts[0])
	d.Set("nat_ip_cidr", parts[1])
	d.Set("nat_ip_cidr_description", object["NatIpCidrDescription"])
	d.Set("nat_ip_cidr_name", object["NatIpCidrName"])
	d.Set("status", object["NatIpCidrStatus"])
	return nil
}
func resourceAlicloudVpcNatIpCidrUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	update := false
	request := map[string]interface{}{
		"NatGatewayId": parts[0],
		"NatIpCidr":    parts[1],
	}
	if d.HasChange("nat_ip_cidr_name") {
		update = true
	}
	if v, ok := d.GetOk("nat_ip_cidr_name"); ok {
		request["NatIpCidrName"] = v
	}
	request["RegionId"] = client.RegionId
	if d.HasChange("dry_run") || d.IsNewResource() {
		update = true
		if v, ok := d.GetOkExists("dry_run"); ok {
			request["DryRun"] = v
		}
	}
	if d.HasChange("nat_ip_cidr_description") {
		update = true
		if v, ok := d.GetOk("nat_ip_cidr_description"); ok {
			request["NatIpCidrDescription"] = v
		}
	}
	if update {
		action := "ModifyNatIpCidrAttribute"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			request["ClientToken"] = buildClientToken("ModifyNatIpCidrAttribute")
			response, err = client.RpcPost("Vpc", "2016-04-28", action, nil, request, true)
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
	return resourceAlicloudVpcNatIpCidrRead(d, meta)
}
func resourceAlicloudVpcNatIpCidrDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	action := "DeleteNatIpCidr"
	var response map[string]interface{}
	request := map[string]interface{}{
		"NatGatewayId": parts[0],
		"NatIpCidr":    parts[1],
	}

	if v, ok := d.GetOkExists("dry_run"); ok {
		request["DryRun"] = v
	}
	request["RegionId"] = client.RegionId
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		request["ClientToken"] = buildClientToken("DeleteNatIpCidr")
		response, err = client.RpcPost("Vpc", "2016-04-28", action, nil, request, true)
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
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	return nil
}

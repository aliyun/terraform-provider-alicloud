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

func resourceAliCloudVpcIpv6CidrBlock() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudVpcIpv6CidrBlockCreate,
		Read:   resourceAliCloudVpcIpv6CidrBlockRead,
		Update: resourceAliCloudVpcIpv6CidrBlockUpdate,
		Delete: resourceAliCloudVpcIpv6CidrBlockDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"ipv6_cidr_block": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"ipv6_cidr_mask": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"ipv6_ipam_pool_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAliCloudVpcIpv6CidrBlockCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "AssociateVpcCidrBlock"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	if v, ok := d.GetOk("vpc_id"); ok {
		request["VpcId"] = v
	}
	if v, ok := d.GetOk("ipv6_cidr_block"); ok {
		request["IPv6CidrBlock"] = v
	}
	request["RegionId"] = client.RegionId

	request["IpVersion"] = "IPv6"
	if v, ok := d.GetOkExists("ipv6_cidr_mask"); ok {
		request["Ipv6CidrMask"] = v
	}
	if v, ok := d.GetOk("ipv6_ipam_pool_id"); ok {
		request["IpamPoolId"] = v
	}
	wait := incrementalWait(5*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("Vpc", "2016-04-28", action, query, request, true)
		if err != nil {
			if IsExpectedErrors(err, []string{"LastTokenProcessing", "OperationConflict", "SystemBusy", "ServiceUnavailable", "IncorrectStatus.Vpc", "IncorrectStatus"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_vpc_ipv6_cidr_block", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%v#%v", request["VpcId"], response["CidrBlock"]))

	return resourceAliCloudVpcIpv6CidrBlockRead(d, meta)
}

func resourceAliCloudVpcIpv6CidrBlockRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vpcServiceV2 := VpcServiceV2{client}

	objectRaw, err := vpcServiceV2.DescribeVpcIpv6CidrBlock(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_vpc_ipv6_cidr_block DescribeVpcIpv6CidrBlock Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("ipv6_cidr_block", objectRaw["Ipv6CidrBlock"])

	parts := strings.Split(d.Id(), "#")
	d.Set("vpc_id", parts[0])

	return nil
}

func resourceAliCloudVpcIpv6CidrBlockUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Cannot update resource Alicloud Resource Ipv6 Cidr Block.")
	return nil
}

func resourceAliCloudVpcIpv6CidrBlockDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	parts := strings.Split(d.Id(), "#")
	action := "UnassociateVpcCidrBlock"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["VpcId"] = parts[0]
	request["IPv6CidrBlock"] = parts[1]
	request["RegionId"] = client.RegionId

	wait := incrementalWait(5*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("Vpc", "2016-04-28", action, query, request, true)
		if err != nil {
			if IsExpectedErrors(err, []string{"LastTokenProcessing", "OperationFailed.CidrInUse", "OperationConflict", "SystemBusy", "ServiceUnavailable", "IncorrectStatus.Vpc", "IncorrectStatus"}) || NeedRetry(err) {
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

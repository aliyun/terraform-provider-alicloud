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

func resourceAliCloudVpcPublicIpAddressPoolCidrBlock() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudVpcPublicIpAddressPoolCidrBlockCreate,
		Read:   resourceAliCloudVpcPublicIpAddressPoolCidrBlockRead,
		Update: resourceAliCloudVpcPublicIpAddressPoolCidrBlockUpdate,
		Delete: resourceAliCloudVpcPublicIpAddressPoolCidrBlockDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"cidr_block": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"cidr_mask": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"public_ip_address_pool_id": {
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

func resourceAliCloudVpcPublicIpAddressPoolCidrBlockCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "AddPublicIpAddressPoolCidrBlock"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	query["PublicIpAddressPoolId"] = d.Get("public_ip_address_pool_id")
	query["CidrBlock"] = d.Get("cidr_block")
	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken(action)

	if v, ok := d.GetOk("cidr_mask"); ok {
		request["CidrMask"] = v
	}
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("Vpc", "2016-04-28", action, query, request, true)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_vpc_public_ip_address_pool_cidr_block", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%v:%v", query["PublicIpAddressPoolId"], response["CidrBlock"]))

	vpcServiceV2 := VpcServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{"Created"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, vpcServiceV2.VpcPublicIpAddressPoolCidrBlockStateRefreshFunc(d.Id(), "Status", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAliCloudVpcPublicIpAddressPoolCidrBlockRead(d, meta)
}

func resourceAliCloudVpcPublicIpAddressPoolCidrBlockRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vpcServiceV2 := VpcServiceV2{client}

	objectRaw, err := vpcServiceV2.DescribeVpcPublicIpAddressPoolCidrBlock(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_vpc_public_ip_address_pool_cidr_block DescribeVpcPublicIpAddressPoolCidrBlock Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("create_time", objectRaw["CreationTime"])
	d.Set("status", objectRaw["Status"])
	d.Set("cidr_block", objectRaw["CidrBlock"])
	d.Set("public_ip_address_pool_id", objectRaw["PublicIpAddressPoolId"])

	parts := strings.Split(d.Id(), ":")
	d.Set("public_ip_address_pool_id", parts[0])
	d.Set("cidr_block", parts[1])

	return nil
}

func resourceAliCloudVpcPublicIpAddressPoolCidrBlockUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Cannot update resource Alicloud Resource Public Ip Address Pool Cidr Block.")
	return nil
}

func resourceAliCloudVpcPublicIpAddressPoolCidrBlockDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	parts := strings.Split(d.Id(), ":")
	action := "DeletePublicIpAddressPoolCidrBlock"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	query["PublicIpAddressPoolId"] = parts[0]
	query["CidrBlock"] = parts[1]
	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken(action)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("Vpc", "2016-04-28", action, query, request, true)
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

	vpcServiceV2 := VpcServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 5*time.Second, vpcServiceV2.VpcPublicIpAddressPoolCidrBlockStateRefreshFunc(d.Id(), "Status", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return nil
}

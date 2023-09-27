// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"log"
	"strings"
	"time"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudVpcIpv4CidrBlock() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudVpcIpv4CidrBlockCreate,
		Read:   resourceAliCloudVpcIpv4CidrBlockRead,
		Delete: resourceAliCloudVpcIpv4CidrBlockDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"secondary_cidr_block": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAliCloudVpcIpv4CidrBlockCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "AssociateVpcCidrBlock"
	var request map[string]interface{}
	var response map[string]interface{}
	conn, err := client.NewVpcClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	request["VpcId"] = d.Get("vpc_id")
	request["SecondaryCidrBlock"] = d.Get("secondary_cidr_block")
	request["RegionId"] = client.RegionId

	request["IpVersion"] = "IPV4"
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, request, &runtime)

		if err != nil {
			if IsExpectedErrors(err, []string{"IncorrectStatus.Vpc", "OperationConflict", "IncorrectStatus", "ServiceUnavailable", "SystemBusy", "LastTokenProcessing"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_vpc_ipv4_cidr_block", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%v:%v", request["VpcId"], request["SecondaryCidrBlock"]))

	return resourceAliCloudVpcIpv4CidrBlockRead(d, meta)
}

func resourceAliCloudVpcIpv4CidrBlockRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vpcServiceV2 := VpcServiceV2{client}

	_, err := vpcServiceV2.DescribeVpcIpv4CidrBlock(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_vpc_ipv4_cidr_block DescribeVpcIpv4CidrBlock Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	d.Set("secondary_cidr_block", parts[1])
	d.Set("vpc_id", parts[0])

	return nil
}

func resourceAliCloudVpcIpv4CidrBlockDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	parts := strings.Split(d.Id(), ":")
	action := "UnassociateVpcCidrBlock"
	var request map[string]interface{}
	var response map[string]interface{}
	conn, err := client.NewVpcClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	request["VpcId"] = parts[0]
	request["SecondaryCidrBlock"] = parts[1]
	request["RegionId"] = client.RegionId

	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, request, &runtime)

		if err != nil {
			if IsExpectedErrors(err, []string{"IncorrectStatus.Vpc", "OperationConflict", "IncorrectStatus", "ServiceUnavailable", "SystemBusy", "LastTokenProcessing", "OperationFailed.CidrInUse"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})

	if err != nil {
		if NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return nil
}

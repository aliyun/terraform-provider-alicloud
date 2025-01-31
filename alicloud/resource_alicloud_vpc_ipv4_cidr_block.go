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

func resourceAliCloudVpcIpv4CidrBlock() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudVpcIpv4CidrBlockCreate,
		Read:   resourceAliCloudVpcIpv4CidrBlockRead,
		Update: resourceAliCloudVpcIpv4CidrBlockUpdate,
		Delete: resourceAliCloudVpcIpv4CidrBlockDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"ipv4_ipam_pool_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"region_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"secondary_cidr_block": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"secondary_cidr_mask": {
				Type:     schema.TypeInt,
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

func resourceAliCloudVpcIpv4CidrBlockCreate(d *schema.ResourceData, meta interface{}) error {

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
	if v, ok := d.GetOk("secondary_cidr_block"); ok {
		request["SecondaryCidrBlock"] = v
	}
	request["RegionId"] = client.RegionId

	if v, ok := d.GetOk("ipv6_isp"); ok {
		request["Ipv6Isp"] = v
	}
	if v, ok := d.GetOk("ipv6_cidr_block"); ok {
		request["IPv6CidrBlock"] = v
	}
	if v, ok := d.GetOk("ipv4_ipam_pool_id"); ok {
		request["IpamPoolId"] = v
	}
	if v, ok := d.GetOk("secondary_cidr_mask"); ok {
		request["SecondaryCidrMask"] = v
	}
	request["IpVersion"] = "IPV4"
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("Vpc", "2016-04-28", action, query, request, false)
		if err != nil {
			if IsExpectedErrors(err, []string{"IncorrectStatus.Vpc", "OperationConflict", "IncorrectStatus", "ServiceUnavailable", "SystemBusy", "LastTokenProcessing"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_vpc_ipv4_cidr_block", action, AlibabaCloudSdkGoERROR)
	}

	if response["CidrBlock"] != nil {
		d.SetId(fmt.Sprintf("%v:%v", request["VpcId"], response["CidrBlock"]))
	} else {
		d.SetId(fmt.Sprintf("%v:%v", request["VpcId"], request["SecondaryCidrBlock"]))
	}

	return resourceAliCloudVpcIpv4CidrBlockRead(d, meta)
}

func resourceAliCloudVpcIpv4CidrBlockRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vpcServiceV2 := VpcServiceV2{client}

	objectRaw, err := vpcServiceV2.DescribeVpcIpv4CidrBlock(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_vpc_ipv4_cidr_block DescribeVpcIpv4CidrBlock Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	if objectRaw["RegionId"] != nil {
		d.Set("region_id", objectRaw["RegionId"])
	}
	if objectRaw["VpcId"] != nil {
		d.Set("vpc_id", objectRaw["VpcId"])
	}

	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	d.Set("secondary_cidr_block", parts[1])
	d.Set("vpc_id", parts[0])
	return nil
}

func resourceAliCloudVpcIpv4CidrBlockUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Cannot update resource Alicloud Resource Ipv4 Cidr Block.")
	return nil
}

func resourceAliCloudVpcIpv4CidrBlockDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	parts := strings.Split(d.Id(), ":")
	action := "UnassociateVpcCidrBlock"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["VpcId"] = parts[0]
	request["SecondaryCidrBlock"] = parts[1]
	request["RegionId"] = client.RegionId

	if v, ok := d.GetOk("ipv6_cidr_block"); ok {
		request["IPv6CidrBlock"] = v
	}
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("Vpc", "2016-04-28", action, query, request, false)

		if err != nil {
			if IsExpectedErrors(err, []string{"IncorrectStatus.Vpc", "OperationConflict", "IncorrectStatus", "ServiceUnavailable", "SystemBusy", "LastTokenProcessing", "OperationFailed.CidrInUse"}) || NeedRetry(err) {
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

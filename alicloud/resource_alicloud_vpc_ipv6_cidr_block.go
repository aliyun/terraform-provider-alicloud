package alicloud

import (
	"fmt"
	"log"
	"time"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAlicloudVpcIpv6CidrBlock() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudVpcIpv6CidrBlockCreate,
		Read:   resourceAlicloudVpcIpv6CidrBlockRead,
		Delete: resourceAlicloudVpcIpv6CidrBlockDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"ipv6_isp": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"BGP", "ChinaMobile", "ChinaUnicom"}, false),
			},
			"secondary_cidr_block": {
				Type:     schema.TypeString,
				Computed: true,
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

func resourceAlicloudVpcIpv6CidrBlockCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "AssociateVpcCidrBlock"
	request := make(map[string]interface{})
	conn, err := client.NewVpcClient()
	if err != nil {
		return WrapError(err)
	}
	request["IpVersion"] = "IPV6"
	if v, ok := d.GetOk("ipv6_isp"); ok {
		request["Ipv6Isp"] = v
	}
	request["RegionId"] = client.RegionId
	if v, ok := d.GetOk("secondary_cidr_block"); ok {
		request["IPv6CidrBlock"] = v
	}
	request["VpcId"] = d.Get("vpc_id")
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_vpc_ipv6_cidr_block", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(request["VpcId"], ":", request["IPv6CidrBlock"]))

	return resourceAlicloudVpcIpv6CidrBlockRead(d, meta)
}
func resourceAlicloudVpcIpv6CidrBlockRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vpcService := VpcService{client}
	object, err := vpcService.DescribeVpcIpv6CidrBlock(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_vpc_ipv6_cidr_block vpcService.DescribeVpcIpv6CidrBlock Failed!!! %s", err)
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
	d.Set("ipv6_isp", object["Ipv6Isp"])
	return nil
}
func resourceAlicloudVpcIpv6CidrBlockDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	action := "UnassociateVpcCidrBlock"
	var response map[string]interface{}
	conn, err := client.NewVpcClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"SecondaryCidrBlock": parts[1],
		"VpcId":              parts[0],
	}

	request["RegionId"] = client.RegionId
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
	return nil
}

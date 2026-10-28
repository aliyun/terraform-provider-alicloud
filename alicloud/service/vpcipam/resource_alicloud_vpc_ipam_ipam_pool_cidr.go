package vpcipam

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
)

func ResourceVpcIpamIpamPoolCidr() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudVpcIpamIpamPoolCidrCreate,
		Read:   resourceAliCloudVpcIpamIpamPoolCidrRead,
		Update: resourceAliCloudVpcIpamIpamPoolCidrUpdate,
		Delete: resourceAliCloudVpcIpamIpamPoolCidrDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"cidr": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"ipam_pool_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"netmask_length": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAliCloudVpcIpamIpamPoolCidrCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "AddIpamPoolCidr"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	if v, ok := d.GetOk("ipam_pool_id"); ok {
		request["IpamPoolId"] = v
	}
	if v, ok := d.GetOk("cidr"); ok {
		request["Cidr"] = v
	}
	request["RegionId"] = client.RegionId
	request["ClientToken"] = helper.BuildClientToken(action)

	if v, ok := d.GetOkExists("netmask_length"); ok {
		request["NetmaskLength"] = v
	}
	wait := helper.IncrementalWait(5*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		response, err = client.RpcPost("VpcIpam", "2023-02-28", action, query, request, true)
		if err != nil {
			if errs.IsExpectedErrors(err, []string{"OperationConflict"}) || helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)

	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_vpc_ipam_ipam_pool_cidr", action, errs.AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%v#%v", request["IpamPoolId"], response["Cidr"]))

	return resourceAliCloudVpcIpamIpamPoolCidrRead(d, meta)
}

func resourceAliCloudVpcIpamIpamPoolCidrRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vpcIpamServiceV2 := NewService(client)

	objectRaw, err := vpcIpamServiceV2.DescribeVpcIpamIpamPoolCidr(d.Id())
	if err != nil {
		if !d.IsNewResource() && errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_vpc_ipam_ipam_pool_cidr DescribeVpcIpamIpamPoolCidr Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}

	d.Set("status", objectRaw["Status"])
	d.Set("cidr", objectRaw["Cidr"])
	d.Set("ipam_pool_id", objectRaw["IpamPoolId"])

	return nil
}

func resourceAliCloudVpcIpamIpamPoolCidrUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Cannot update resource Alicloud Resource Ipam Pool Cidr.")
	return nil
}

func resourceAliCloudVpcIpamIpamPoolCidrDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	parts, err := splitIpamPoolCidrId(d.Id())
	if err != nil {
		return errs.WrapError(err)
	}
	action := "DeleteIpamPoolCidr"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	request = make(map[string]interface{})
	request["IpamPoolId"] = parts[0]
	request["Cidr"] = parts[1]
	request["RegionId"] = client.RegionId
	request["ClientToken"] = helper.BuildClientToken(action)

	wait := helper.IncrementalWait(5*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		response, err = client.RpcPost("VpcIpam", "2023-02-28", action, query, request, true)
		if err != nil {
			if errs.IsExpectedErrors(err, []string{"OperationConflict", "OperationDenied.IpamPoolAllocationExist"}) || helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)

	if err != nil {
		if errs.NotFoundError(err) {
			return nil
		}
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
	}

	return nil
}

// splitIpamPoolCidrId splits the resource ID with backward compatibility
// Returns [ipamPoolId, cidr]
// Supports both "#" (new) and ":" (old) separators
func splitIpamPoolCidrId(id string) ([]string, error) {
	var parts []string
	// Try new separator "#" first
	if strings.Contains(id, "#") {
		parts = strings.Split(id, "#")
	} else if strings.Contains(id, ":") {
		// Fall back to old separator ":"
		parts = strings.Split(id, ":")
	} else {
		return nil, fmt.Errorf("invalid Resource Id %s. Expected separator '#' or ':'", id)
	}
	return parts, nil
}

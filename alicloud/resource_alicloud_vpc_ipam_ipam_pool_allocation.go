// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudVpcIpamIpamPoolAllocation() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudVpcIpamIpamPoolAllocationCreate,
		Read:   resourceAliCloudVpcIpamIpamPoolAllocationRead,
		Update: resourceAliCloudVpcIpamIpamPoolAllocationUpdate,
		Delete: resourceAliCloudVpcIpamIpamPoolAllocationDelete,
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
			"cidr_mask": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ipam_pool_allocation_description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ipam_pool_allocation_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ipam_pool_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"region_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAliCloudVpcIpamIpamPoolAllocationCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateIpamPoolAllocation"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken(action)

	request["IpamPoolId"] = d.Get("ipam_pool_id")
	if v, ok := d.GetOk("cidr"); ok {
		request["Cidr"] = v
	}
	if v, ok := d.GetOk("ipam_pool_allocation_description"); ok {
		request["IpamPoolAllocationDescription"] = v
	}
	if v, ok := d.GetOk("ipam_pool_allocation_name"); ok {
		request["IpamPoolAllocationName"] = v
	}
	if v, ok := d.GetOkExists("cidr_mask"); ok {
		request["CidrMask"] = v
	}
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("VpcIpam", "2023-02-28", action, query, request, true)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_vpc_ipam_ipam_pool_allocation", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["IpamPoolAllocationId"]))

	return resourceAliCloudVpcIpamIpamPoolAllocationRead(d, meta)
}

func resourceAliCloudVpcIpamIpamPoolAllocationRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vpcIpamServiceV2 := VpcIpamServiceV2{client}

	objectRaw, err := vpcIpamServiceV2.DescribeVpcIpamIpamPoolAllocation(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_vpc_ipam_ipam_pool_allocation DescribeVpcIpamIpamPoolAllocation Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	if objectRaw["Cidr"] != nil {
		d.Set("cidr", objectRaw["Cidr"])
	}
	if objectRaw["CreationTime"] != nil {
		d.Set("create_time", objectRaw["CreationTime"])
	}
	if objectRaw["IpamPoolAllocationDescription"] != nil {
		d.Set("ipam_pool_allocation_description", objectRaw["IpamPoolAllocationDescription"])
	}
	if objectRaw["IpamPoolAllocationName"] != nil {
		d.Set("ipam_pool_allocation_name", objectRaw["IpamPoolAllocationName"])
	}
	if objectRaw["IpamPoolId"] != nil {
		d.Set("ipam_pool_id", objectRaw["IpamPoolId"])
	}
	if objectRaw["RegionId"] != nil {
		d.Set("region_id", objectRaw["RegionId"])
	}
	if objectRaw["Status"] != nil {
		d.Set("status", objectRaw["Status"])
	}

	return nil
}

func resourceAliCloudVpcIpamIpamPoolAllocationUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false

	action := "UpdateIpamPoolAllocation"
	var err error
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["IpamPoolAllocationId"] = d.Id()
	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken(action)
	if d.HasChange("ipam_pool_allocation_description") {
		update = true
		request["IpamPoolAllocationDescription"] = d.Get("ipam_pool_allocation_description")
	}

	if d.HasChange("ipam_pool_allocation_name") {
		update = true
		request["IpamPoolAllocationName"] = d.Get("ipam_pool_allocation_name")
	}

	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("VpcIpam", "2023-02-28", action, query, request, true)
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

	return resourceAliCloudVpcIpamIpamPoolAllocationRead(d, meta)
}

func resourceAliCloudVpcIpamIpamPoolAllocationDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	action := "DeleteIpamPoolAllocation"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["IpamPoolAllocationId"] = d.Id()
	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken(action)

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("VpcIpam", "2023-02-28", action, query, request, true)
		request["ClientToken"] = buildClientToken(action)

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
		if IsExpectedErrors(err, []string{"ResourceNotFound.IpamPoolAllocation"}) || NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return nil
}

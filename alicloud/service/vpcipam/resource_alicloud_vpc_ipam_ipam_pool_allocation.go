// Package vpcipam. This file is generated automatically. Please do not modify it manually, thank you!
package vpcipam

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
)

func ResourceVpcIpamIpamPoolAllocation() *schema.Resource {
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
	request["ClientToken"] = helper.BuildClientToken(action)

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
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		response, err = client.RpcPost("VpcIpam", "2023-02-28", action, query, request, true)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)

	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_vpc_ipam_ipam_pool_allocation", action, errs.AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["IpamPoolAllocationId"]))

	return resourceAliCloudVpcIpamIpamPoolAllocationRead(d, meta)
}

func resourceAliCloudVpcIpamIpamPoolAllocationRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vpcIpamServiceV2 := NewService(client)

	objectRaw, err := vpcIpamServiceV2.DescribeVpcIpamIpamPoolAllocation(d.Id())
	if err != nil {
		if !d.IsNewResource() && errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_vpc_ipam_ipam_pool_allocation DescribeVpcIpamIpamPoolAllocation Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
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
	request["ClientToken"] = helper.BuildClientToken(action)
	if d.HasChange("ipam_pool_allocation_description") {
		update = true
		request["IpamPoolAllocationDescription"] = d.Get("ipam_pool_allocation_description")
	}

	if d.HasChange("ipam_pool_allocation_name") {
		update = true
		request["IpamPoolAllocationName"] = d.Get("ipam_pool_allocation_name")
	}

	if update {
		wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RpcPost("VpcIpam", "2023-02-28", action, query, request, true)
			if err != nil {
				if helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		logs.AddDebug(action, response, request)
		if err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
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
	request["ClientToken"] = helper.BuildClientToken(action)

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		response, err = client.RpcPost("VpcIpam", "2023-02-28", action, query, request, true)
		request["ClientToken"] = helper.BuildClientToken(action)

		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)

	if err != nil {
		if errs.IsExpectedErrors(err, []string{"ResourceNotFound.IpamPoolAllocation"}) || errs.NotFoundError(err) {
			return nil
		}
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
	}

	return nil
}

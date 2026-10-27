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
	"github.com/aliyun/terraform-provider-alicloud/alicloud/tags"
)

func ResourceVpcIpamIpamPool() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudVpcIpamIpamPoolCreate,
		Read:   resourceAliCloudVpcIpamIpamPoolRead,
		Update: resourceAliCloudVpcIpamIpamPoolUpdate,
		Delete: resourceAliCloudVpcIpamIpamPoolDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"allocation_default_cidr_mask": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"allocation_max_cidr_mask": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"allocation_min_cidr_mask": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"auto_import": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"clear_allocation_default_cidr_mask": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ip_version": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"ipam_pool_description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ipam_pool_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"ipam_scope_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"ipv6_isp": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"pool_region_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"region_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"source_ipam_pool_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": tags.TagsSchema(),
		},
	}
}

func resourceAliCloudVpcIpamIpamPoolCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateIpamPool"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["RegionId"] = client.RegionId
	request["ClientToken"] = helper.BuildClientToken(action)

	if v, ok := d.GetOk("pool_region_id"); ok {
		request["PoolRegionId"] = v
	}
	if v, ok := d.GetOk("resource_group_id"); ok {
		request["ResourceGroupId"] = v
	}
	if v, ok := d.GetOk("tags"); ok {
		tagsMap := tags.ConvertTags(v.(map[string]interface{}))
		request = tags.ExpandTagsToMap(request, tagsMap)
	}

	if v, ok := d.GetOk("source_ipam_pool_id"); ok {
		request["SourceIpamPoolId"] = v
	}
	if v, ok := d.GetOk("ipv6_isp"); ok {
		request["Ipv6Isp"] = v
	}
	if v, ok := d.GetOkExists("allocation_max_cidr_mask"); ok {
		request["AllocationMaxCidrMask"] = v
	}
	if v, ok := d.GetOkExists("allocation_default_cidr_mask"); ok {
		request["AllocationDefaultCidrMask"] = v
	}
	if v, ok := d.GetOkExists("allocation_min_cidr_mask"); ok {
		request["AllocationMinCidrMask"] = v
	}
	request["IpamScopeId"] = d.Get("ipam_scope_id")
	if v, ok := d.GetOk("ipam_pool_description"); ok {
		request["IpamPoolDescription"] = v
	}
	if v, ok := d.GetOk("ip_version"); ok {
		request["IpVersion"] = v
	}
	if v, ok := d.GetOkExists("auto_import"); ok {
		request["AutoImport"] = v
	}
	if v, ok := d.GetOk("ipam_pool_name"); ok {
		request["IpamPoolName"] = v
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
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_vpc_ipam_ipam_pool", action, errs.AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["IpamPoolId"]))

	return resourceAliCloudVpcIpamIpamPoolRead(d, meta)
}

func resourceAliCloudVpcIpamIpamPoolRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vpcIpamServiceV2 := NewService(client)

	objectRaw, err := vpcIpamServiceV2.DescribeVpcIpamIpamPool(d.Id())
	if err != nil {
		if !d.IsNewResource() && errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_vpc_ipam_ipam_pool DescribeVpcIpamIpamPool Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}

	d.Set("allocation_default_cidr_mask", objectRaw["AllocationDefaultCidrMask"])
	d.Set("allocation_max_cidr_mask", objectRaw["AllocationMaxCidrMask"])
	d.Set("allocation_min_cidr_mask", objectRaw["AllocationMinCidrMask"])
	d.Set("auto_import", objectRaw["AutoImport"])
	d.Set("create_time", objectRaw["CreateTime"])
	d.Set("ip_version", objectRaw["IpVersion"])
	d.Set("ipam_pool_description", objectRaw["IpamPoolDescription"])
	d.Set("ipam_pool_name", objectRaw["IpamPoolName"])
	d.Set("ipam_scope_id", objectRaw["IpamScopeId"])
	d.Set("ipv6_isp", objectRaw["Ipv6Isp"])
	d.Set("pool_region_id", objectRaw["PoolRegionId"])
	d.Set("resource_group_id", objectRaw["ResourceGroupId"])
	d.Set("source_ipam_pool_id", objectRaw["SourceIpamPoolId"])
	d.Set("status", objectRaw["Status"])
	d.Set("region_id", objectRaw["IpamRegionId"])

	tagsMaps := objectRaw["Tags"]
	d.Set("tags", tags.TagsToMap(tagsMaps))

	return nil
}

func resourceAliCloudVpcIpamIpamPoolUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false

	var err error
	action := "UpdateIpamPool"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["IpamPoolId"] = d.Id()
	request["RegionId"] = client.RegionId
	request["ClientToken"] = helper.BuildClientToken(action)
	if d.HasChange("allocation_min_cidr_mask") {
		update = true
		request["AllocationMinCidrMask"] = d.Get("allocation_min_cidr_mask")
	}

	if d.HasChange("ipam_pool_description") {
		update = true
		request["IpamPoolDescription"] = d.Get("ipam_pool_description")
	}

	if d.HasChange("allocation_max_cidr_mask") {
		update = true
		request["AllocationMaxCidrMask"] = d.Get("allocation_max_cidr_mask")
	}

	if d.HasChange("allocation_default_cidr_mask") {
		update = true
		request["AllocationDefaultCidrMask"] = d.Get("allocation_default_cidr_mask")
	}

	if v, ok := d.GetOkExists("clear_allocation_default_cidr_mask"); ok {
		request["ClearAllocationDefaultCidrMask"] = v
	}
	if d.HasChange("auto_import") {
		update = true
		request["AutoImport"] = d.Get("auto_import")
	}

	if d.HasChange("ipam_pool_name") {
		update = true
		request["IpamPoolName"] = d.Get("ipam_pool_name")
	}

	if update {
		wait := helper.IncrementalWait(5*time.Second, 5*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
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
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
		}
	}
	update = false
	action = "ChangeResourceGroup"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["ResourceId"] = d.Id()
	request["RegionId"] = client.RegionId
	if _, ok := d.GetOk("resource_group_id"); ok && d.HasChange("resource_group_id") {
		update = true
	}
	request["NewResourceGroupId"] = d.Get("resource_group_id")
	request["ResourceType"] = "IPAMPOOL"
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

	if d.HasChange("tags") {
		vpcIpamServiceV2 := NewService(client)
		if err := vpcIpamServiceV2.SetResourceTags(d, "IPAMPOOL"); err != nil {
			return errs.WrapError(err)
		}
	}
	return resourceAliCloudVpcIpamIpamPoolRead(d, meta)
}

func resourceAliCloudVpcIpamIpamPoolDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	action := "DeleteIpamPool"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["IpamPoolId"] = d.Id()
	request["RegionId"] = client.RegionId
	request["ClientToken"] = helper.BuildClientToken(action)

	wait := helper.IncrementalWait(5*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
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
		if errs.NotFoundError(err) {
			return nil
		}
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
	}

	return nil
}

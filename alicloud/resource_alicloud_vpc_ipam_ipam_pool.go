// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"log"
	"time"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudVpcIpamIpamPool() *schema.Resource {
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
			"pool_region_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"region_id": {
				Type:     schema.TypeString,
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
			"tags": tagsSchema(),
		},
	}
}

func resourceAliCloudVpcIpamIpamPoolCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateIpamPool"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	conn, err := client.NewVpcipamClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken(action)

	request["IpamScopeId"] = d.Get("ipam_scope_id")
	if v, ok := d.GetOk("ipam_pool_name"); ok {
		request["IpamPoolName"] = v
	}
	if v, ok := d.GetOk("ipam_pool_description"); ok {
		request["IpamPoolDescription"] = v
	}
	if v, ok := d.GetOk("pool_region_id"); ok {
		request["PoolRegionId"] = v
	}
	if v, ok := d.GetOkExists("allocation_default_cidr_mask"); ok {
		request["AllocationDefaultCidrMask"] = v
	}
	if v, ok := d.GetOkExists("allocation_max_cidr_mask"); ok {
		request["AllocationMaxCidrMask"] = v
	}
	if v, ok := d.GetOkExists("allocation_min_cidr_mask"); ok {
		request["AllocationMinCidrMask"] = v
	}
	if v, ok := d.GetOk("ip_version"); ok {
		request["IpVersion"] = v
	}
	if v, ok := d.GetOk("source_ipam_pool_id"); ok {
		request["SourceIpamPoolId"] = v
	}
	if v, ok := d.GetOkExists("auto_import"); ok {
		request["AutoImport"] = v
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2023-02-28"), StringPointer("AK"), query, request, &runtime)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_vpc_ipam_ipam_pool", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["IpamPoolId"]))

	return resourceAliCloudVpcIpamIpamPoolUpdate(d, meta)
}

func resourceAliCloudVpcIpamIpamPoolRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vpcIpamServiceV2 := VpcIpamServiceV2{client}

	objectRaw, err := vpcIpamServiceV2.DescribeVpcIpamIpamPool(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_vpc_ipam_ipam_pool DescribeVpcIpamIpamPool Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	if objectRaw["AllocationDefaultCidrMask"] != nil {
		d.Set("allocation_default_cidr_mask", objectRaw["AllocationDefaultCidrMask"])
	}
	if objectRaw["AllocationMaxCidrMask"] != nil {
		d.Set("allocation_max_cidr_mask", objectRaw["AllocationMaxCidrMask"])
	}
	if objectRaw["AllocationMinCidrMask"] != nil {
		d.Set("allocation_min_cidr_mask", objectRaw["AllocationMinCidrMask"])
	}
	if objectRaw["AutoImport"] != nil {
		d.Set("auto_import", objectRaw["AutoImport"])
	}
	if objectRaw["CreateTime"] != nil {
		d.Set("create_time", objectRaw["CreateTime"])
	}
	if objectRaw["IpVersion"] != nil {
		d.Set("ip_version", objectRaw["IpVersion"])
	}
	if objectRaw["IpamPoolDescription"] != nil {
		d.Set("ipam_pool_description", objectRaw["IpamPoolDescription"])
	}
	if objectRaw["IpamPoolName"] != nil {
		d.Set("ipam_pool_name", objectRaw["IpamPoolName"])
	}
	if objectRaw["IpamScopeId"] != nil {
		d.Set("ipam_scope_id", objectRaw["IpamScopeId"])
	}
	if objectRaw["PoolRegionId"] != nil {
		d.Set("pool_region_id", objectRaw["PoolRegionId"])
	}
	if objectRaw["SourceIpamPoolId"] != nil {
		d.Set("source_ipam_pool_id", objectRaw["SourceIpamPoolId"])
	}
	if objectRaw["Status"] != nil {
		d.Set("status", objectRaw["Status"])
	}
	if objectRaw["IpamRegionId"] != nil {
		d.Set("region_id", objectRaw["IpamRegionId"])
	}

	tagsMaps := objectRaw["Tags"]
	d.Set("tags", tagsToMap(tagsMaps))

	return nil
}

func resourceAliCloudVpcIpamIpamPoolUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false
	action := "UpdateIpamPool"
	conn, err := client.NewVpcipamClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["IpamPoolId"] = d.Id()
	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken(action)
	if !d.IsNewResource() && d.HasChange("ipam_pool_name") {
		update = true
		request["IpamPoolName"] = d.Get("ipam_pool_name")
	}

	if !d.IsNewResource() && d.HasChange("ipam_pool_description") {
		update = true
		request["IpamPoolDescription"] = d.Get("ipam_pool_description")
	}

	if !d.IsNewResource() && d.HasChange("allocation_default_cidr_mask") {
		update = true
		request["AllocationDefaultCidrMask"] = d.Get("allocation_default_cidr_mask")
	}

	if !d.IsNewResource() && d.HasChange("allocation_max_cidr_mask") {
		update = true
		request["AllocationMaxCidrMask"] = d.Get("allocation_max_cidr_mask")
	}

	if !d.IsNewResource() && d.HasChange("allocation_min_cidr_mask") {
		update = true
		request["AllocationMinCidrMask"] = d.Get("allocation_min_cidr_mask")
	}

	if v, ok := d.GetOkExists("clear_allocation_default_cidr_mask"); ok {
		request["ClearAllocationDefaultCidrMask"] = v
	}
	if !d.IsNewResource() && d.HasChange("auto_import") {
		update = true
		request["AutoImport"] = d.Get("auto_import")
	}

	if update {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2023-02-28"), StringPointer("AK"), query, request, &runtime)
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

	if d.HasChange("tags") {
		vpcIpamServiceV2 := VpcIpamServiceV2{client}
		if err := vpcIpamServiceV2.SetResourceTags(d, "IpamPool"); err != nil {
			return WrapError(err)
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
	conn, err := client.NewVpcipamClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	request["IpamPoolId"] = d.Id()
	request["RegionId"] = client.RegionId

	request["ClientToken"] = buildClientToken(action)

	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2023-02-28"), StringPointer("AK"), query, request, &runtime)
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
		if NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return nil
}

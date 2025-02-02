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

func resourceAliCloudVpcPublicIpAddressPool() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudVpcPublicIpAddressPoolCreate,
		Read:   resourceAliCloudVpcPublicIpAddressPoolRead,
		Update: resourceAliCloudVpcPublicIpAddressPoolUpdate,
		Delete: resourceAliCloudVpcPublicIpAddressPoolDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"biz_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"Default", "CloudBox"}, false),
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ip_address_remaining": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"isp": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"public_ip_address_pool_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"public_ip_address_pool_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"security_protection_types": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": tagsSchema(),
			"total_ip_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"used_ip_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func resourceAliCloudVpcPublicIpAddressPoolCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreatePublicIpAddressPool"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	query["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken(action)

	if v, ok := d.GetOk("isp"); ok {
		request["Isp"] = v
	}
	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}
	if v, ok := d.GetOk("public_ip_address_pool_name"); ok {
		request["Name"] = v
	}
	if v, ok := d.GetOk("resource_group_id"); ok {
		request["ResourceGroupId"] = v
	}
	if v, ok := d.GetOk("tags"); ok {
		tagsMap := ConvertTags(v.(map[string]interface{}))
		request = expandTagsToMap(request, tagsMap)
	}

	if v, ok := d.GetOk("biz_type"); ok {
		request["BizType"] = v
	}
	if v, ok := d.GetOk("zones"); ok {
		zonesMaps := v.([]interface{})
		request["Zones"] = zonesMaps
	}

	if v, ok := d.GetOk("security_protection_types"); ok {
		securityProtectionTypesMaps := v.([]interface{})
		request["SecurityProtectionTypes"] = securityProtectionTypesMaps
	}

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("Vpc", "2016-04-28", action, query, request, true)
		if err != nil {
			if IsExpectedErrors(err, []string{"OperationConflict", "IncorrectStatus", "SystemBusy", "ServiceUnavailable", "OperationFailed.LastTokenProcessing", "LastTokenProcessing"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_vpc_public_ip_address_pool", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["PulbicIpAddressPoolId"]))

	vpcServiceV2 := VpcServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{"Created"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, vpcServiceV2.VpcPublicIpAddressPoolStateRefreshFunc(d.Id(), "Status", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAliCloudVpcPublicIpAddressPoolRead(d, meta)
}

func resourceAliCloudVpcPublicIpAddressPoolRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vpcServiceV2 := VpcServiceV2{client}

	objectRaw, err := vpcServiceV2.DescribeVpcPublicIpAddressPool(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_vpc_public_ip_address_pool DescribeVpcPublicIpAddressPool Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	if objectRaw["BizType"] != nil {
		d.Set("biz_type", objectRaw["BizType"])
	}
	if objectRaw["CreationTime"] != nil {
		d.Set("create_time", objectRaw["CreationTime"])
	}
	if objectRaw["Description"] != nil {
		d.Set("description", objectRaw["Description"])
	}
	if objectRaw["IpAddressRemaining"] != nil {
		d.Set("ip_address_remaining", objectRaw["IpAddressRemaining"])
	}
	if objectRaw["Isp"] != nil {
		d.Set("isp", objectRaw["Isp"])
	}
	if objectRaw["Name"] != nil {
		d.Set("public_ip_address_pool_name", objectRaw["Name"])
	}
	if objectRaw["ResourceGroupId"] != nil {
		d.Set("resource_group_id", objectRaw["ResourceGroupId"])
	}
	if objectRaw["Status"] != nil {
		d.Set("status", objectRaw["Status"])
	}
	if objectRaw["TotalIpNum"] != nil {
		d.Set("total_ip_num", objectRaw["TotalIpNum"])
	}
	if objectRaw["UsedIpNum"] != nil {
		d.Set("used_ip_num", objectRaw["UsedIpNum"])
	}
	if objectRaw["PublicIpAddressPoolId"] != nil {
		d.Set("public_ip_address_pool_id", objectRaw["PublicIpAddressPoolId"])
	}

	securityProtectionTypes1Raw := make([]interface{}, 0)
	if objectRaw["SecurityProtectionTypes"] != nil {
		securityProtectionTypes1Raw = objectRaw["SecurityProtectionTypes"].([]interface{})
	}

	d.Set("security_protection_types", securityProtectionTypes1Raw)
	tagsMaps := objectRaw["Tags"]
	d.Set("tags", tagsToMap(tagsMaps))
	zones1Raw := make([]interface{}, 0)
	if objectRaw["Zones"] != nil {
		zones1Raw = objectRaw["Zones"].([]interface{})
	}

	d.Set("zones", zones1Raw)

	return nil
}

func resourceAliCloudVpcPublicIpAddressPoolUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false
	d.Partial(true)
	action := "UpdatePublicIpAddressPoolAttribute"
	var err error
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	query["PublicIpAddressPoolId"] = d.Id()
	query["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken(action)
	if d.HasChange("description") {
		update = true
		request["Description"] = d.Get("description")
	}

	if d.HasChange("public_ip_address_pool_name") {
		update = true
		request["Name"] = d.Get("public_ip_address_pool_name")
	}

	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("Vpc", "2016-04-28", action, query, request, true)
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
	}
	update = false
	action = "MoveResourceGroup"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	query["ResourceId"] = d.Id()
	query["RegionId"] = client.RegionId
	if _, ok := d.GetOk("resource_group_id"); ok && d.HasChange("resource_group_id") {
		update = true
		request["NewResourceGroupId"] = d.Get("resource_group_id")
	}

	request["ResourceType"] = "PUBLICIPADDRESSPOOL"
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("Vpc", "2016-04-28", action, query, request, false)
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
	}

	if d.HasChange("tags") {
		vpcServiceV2 := VpcServiceV2{client}
		if err := vpcServiceV2.SetResourceTags(d, "PUBLICIPADDRESSPOOL"); err != nil {
			return WrapError(err)
		}
	}
	d.Partial(false)
	return resourceAliCloudVpcPublicIpAddressPoolRead(d, meta)
}

func resourceAliCloudVpcPublicIpAddressPoolDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	action := "DeletePublicIpAddressPool"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	query["PublicIpAddressPoolId"] = d.Id()
	query["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken(action)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("Vpc", "2016-04-28", action, query, request, true)
		request["ClientToken"] = buildClientToken(action)

		if err != nil {
			if IsExpectedErrors(err, []string{"OperationConflict", "IncorrectStatus", "SystemBusy", "ServiceUnavailable", "OperationFailed.LastTokenProcessing", "LastTokenProcessing", "IncorrectStatus.PublicIpAddressPool"}) || NeedRetry(err) {
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
	stateConf := BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 5*time.Second, vpcServiceV2.VpcPublicIpAddressPoolStateRefreshFunc(d.Id(), "Status", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return nil
}

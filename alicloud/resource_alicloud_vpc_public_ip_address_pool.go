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

func resourceAliCloudVpcPublicIpAddressPool() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudVpcPublicIpAddressPoolCreate,
		Read:   resourceAlicloudVpcPublicIpAddressPoolRead,
		Update: resourceAlicloudVpcPublicIpAddressPoolUpdate,
		Delete: resourceAlicloudVpcPublicIpAddressPoolDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
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

func resourceAlicloudVpcPublicIpAddressPoolCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "CreatePublicIpAddressPool"
	var request map[string]interface{}
	var response map[string]interface{}
	conn, err := client.NewVpcClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	request["RegionId"] = client.RegionId
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

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
		request["ClientToken"] = buildClientToken(action)

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
	stateConf := BuildStateConf([]string{}, []string{"Created"}, d.Timeout(schema.TimeoutCreate), 0, vpcServiceV2.VpcPublicIpAddressPoolStateRefreshFunc(d.Id(), "Status", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAlicloudVpcPublicIpAddressPoolUpdate(d, meta)
}

func resourceAlicloudVpcPublicIpAddressPoolRead(d *schema.ResourceData, meta interface{}) error {
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

	d.Set("create_time", objectRaw["CreationTime"])
	d.Set("description", objectRaw["Description"])
	d.Set("ip_address_remaining", objectRaw["IpAddressRemaining"])
	d.Set("isp", objectRaw["Isp"])
	d.Set("public_ip_address_pool_name", objectRaw["Name"])
	d.Set("resource_group_id", objectRaw["ResourceGroupId"])
	d.Set("status", objectRaw["Status"])
	d.Set("total_ip_num", objectRaw["TotalIpNum"])
	d.Set("used_ip_num", objectRaw["UsedIpNum"])
	d.Set("public_ip_address_pool_id", objectRaw["PublicIpAddressPoolId"])
	d.Set("region_id", objectRaw["RegionId"])
	tagsMaps := objectRaw["Tags"]
	d.Set("tags", tagsToMap(tagsMaps))

	return nil
}

func resourceAlicloudVpcPublicIpAddressPoolUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	update := false
	d.Partial(true)
	update = false
	action := "UpdatePublicIpAddressPoolAttribute"
	conn, err := client.NewVpcClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})

	request["PublicIpAddressPoolId"] = d.Id()
	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken(action)

	if !d.IsNewResource() && d.HasChange("description") {
		update = true
		if v, ok := d.GetOk("description"); ok {
			request["Description"] = v
		}
	}

	if !d.IsNewResource() && d.HasChange("public_ip_address_pool_name") {
		update = true
		if v, ok := d.GetOk("public_ip_address_pool_name"); ok {
			request["Name"] = v
		}
	}

	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
		d.SetPartial("description")
		d.SetPartial("public_ip_address_pool_name")
	}
	update = false
	action = "MoveResourceGroup"
	conn, err = client.NewVpcClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})

	request["ResourceId"] = d.Id()
	request["RegionId"] = client.RegionId

	if !d.IsNewResource() && d.HasChange("resource_group_id") {
		update = true
		if v, ok := d.GetOk("resource_group_id"); ok {
			request["NewResourceGroupId"] = v
		}
	}

	request["ResourceType"] = "PUBLICIPADDRESSPOOL"
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})

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
		d.SetPartial("resource_group_id")
	}

	update = false
	if d.HasChange("tags") {
		update = true
		vpcServiceV2 := VpcServiceV2{client}
		if err := vpcServiceV2.SetResourceTags(d, "PUBLICIPADDRESSPOOL"); err != nil {
			return WrapError(err)
		}
		d.SetPartial("tags")
	}
	d.Partial(false)
	return resourceAlicloudVpcPublicIpAddressPoolRead(d, meta)
}

func resourceAlicloudVpcPublicIpAddressPoolDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "DeletePublicIpAddressPool"
	var request map[string]interface{}
	var response map[string]interface{}
	conn, err := client.NewVpcClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})

	request["PublicIpAddressPoolId"] = d.Id()
	request["RegionId"] = client.RegionId

	request["ClientToken"] = buildClientToken(action)

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
	stateConf := BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 0, vpcServiceV2.VpcPublicIpAddressPoolStateRefreshFunc(d.Id(), "Status", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return nil
}

// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudPolardbGlobalSecurityIpGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudPolardbGlobalSecurityIpGroupCreate,
		Read:   resourceAliCloudPolardbGlobalSecurityIpGroupRead,
		Update: resourceAliCloudPolardbGlobalSecurityIpGroupUpdate,
		Delete: resourceAliCloudPolardbGlobalSecurityIpGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"global_ip_group_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"global_ip_list": {
				Type:     schema.TypeString,
				Required: true,
			},
			"region_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAliCloudPolardbGlobalSecurityIpGroupCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateGlobalSecurityIPGroup"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["RegionId"] = client.RegionId

	request["GlobalIgName"] = d.Get("global_ip_group_name")
	request["GIpList"] = d.Get("global_ip_list")
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("polardb", "2017-08-01", action, query, request, true)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_polardb_global_security_ip_group", action, AlibabaCloudSdkGoERROR)
	}

	id, _ := jsonpath.Get("$.GlobalSecurityIPGroup[0].GlobalSecurityGroupId", response)
	d.SetId(fmt.Sprint(id))

	return resourceAliCloudPolardbGlobalSecurityIpGroupRead(d, meta)
}

func resourceAliCloudPolardbGlobalSecurityIpGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	polardbServiceV2 := PolardbServiceV2{client}

	objectRaw, err := polardbServiceV2.DescribePolardbGlobalSecurityIpGroup(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_polardb_global_security_ip_group DescribePolardbGlobalSecurityIpGroup Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("global_ip_group_name", objectRaw["GlobalIgName"])
	d.Set("global_ip_list", objectRaw["GIpList"])
	d.Set("region_id", objectRaw["RegionId"])

	return nil
}

func resourceAliCloudPolardbGlobalSecurityIpGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false
	d.Partial(true)

	var err error
	action := "ModifyGlobalSecurityIPGroupName"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["GlobalSecurityGroupId"] = d.Id()
	request["RegionId"] = client.RegionId
	if d.HasChange("global_ip_group_name") {
		update = true
	}
	request["GlobalIgName"] = d.Get("global_ip_group_name")
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("polardb", "2017-08-01", action, query, request, true)
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
	update = false
	action = "ModifyGlobalSecurityIPGroup"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["GlobalSecurityGroupId"] = d.Id()
	request["RegionId"] = client.RegionId
	if d.HasChange("global_ip_group_name") {
		update = true
	}
	request["GlobalIgName"] = d.Get("global_ip_group_name")
	if d.HasChange("global_ip_list") {
		update = true
	}
	request["GIpList"] = d.Get("global_ip_list")
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("polardb", "2017-08-01", action, query, request, true)
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

	d.Partial(false)
	return resourceAliCloudPolardbGlobalSecurityIpGroupRead(d, meta)
}

func resourceAliCloudPolardbGlobalSecurityIpGroupDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	action := "DeleteGlobalSecurityIPGroup"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["GlobalSecurityGroupId"] = d.Id()
	request["RegionId"] = client.RegionId

	request["GlobalIgName"] = d.Get("global_ip_group_name")
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("polardb", "2017-08-01", action, query, request, true)
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
		if IsExpectedErrors(err, []string{"SecAPICallingFailed"}) || NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return nil
}

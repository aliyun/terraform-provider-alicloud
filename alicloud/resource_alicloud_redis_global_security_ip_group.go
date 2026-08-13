// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAliCloudRedisGlobalSecurityIpGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudRedisGlobalSecurityIpGroupCreate,
		Read:   resourceAliCloudRedisGlobalSecurityIpGroupRead,
		Update: resourceAliCloudRedisGlobalSecurityIpGroupUpdate,
		Delete: resourceAliCloudRedisGlobalSecurityIpGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"global_ip_list": {
				Type:     schema.TypeString,
				Required: true,
			},
			"global_ip_group_name": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceAliCloudRedisGlobalSecurityIpGroupCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateGlobalSecurityIPGroup"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["RegionId"] = client.RegionId

	request["GIpList"] = d.Get("global_ip_list")
	request["GlobalIgName"] = d.Get("global_ip_group_name")
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("R-kvstore", "2015-01-01", action, query, request, true)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_redis_global_security_ip_group", action, AlibabaCloudSdkGoERROR)
	}

	id, _ := jsonpath.Get("$.GlobalSecurityIPGroup[0].GlobalSecurityGroupId", response)
	d.SetId(fmt.Sprint(id))

	return resourceAliCloudRedisGlobalSecurityIpGroupRead(d, meta)
}

func resourceAliCloudRedisGlobalSecurityIpGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	redisServiceV2 := RedisServiceV2{client}

	objectRaw, err := redisServiceV2.DescribeRedisGlobalSecurityIpGroup(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_redis_global_security_ip_group DescribeRedisGlobalSecurityIpGroup Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("global_ip_list", objectRaw["GIpList"])
	d.Set("global_ip_group_name", objectRaw["GlobalIgName"])

	return nil
}

func resourceAliCloudRedisGlobalSecurityIpGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false

	var err error
	action := "ModifyGlobalSecurityIPGroup"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["GlobalSecurityGroupId"] = d.Id()
	request["RegionId"] = client.RegionId
	if d.HasChange("global_ip_list") {
		update = true
	}
	request["GIpList"] = d.Get("global_ip_list")
	if d.HasChange("global_ip_group_name") {
		update = true
	}
	request["GlobalIgName"] = d.Get("global_ip_group_name")
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("R-kvstore", "2015-01-01", action, query, request, true)
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
	action = "ModifyGlobalSecurityIPGroupName"
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
			response, err = client.RpcPost("R-kvstore", "2015-01-01", action, query, request, true)
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

	return resourceAliCloudRedisGlobalSecurityIpGroupRead(d, meta)
}

func resourceAliCloudRedisGlobalSecurityIpGroupDelete(d *schema.ResourceData, meta interface{}) error {

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
		response, err = client.RpcPost("R-kvstore", "2015-01-01", action, query, request, true)
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

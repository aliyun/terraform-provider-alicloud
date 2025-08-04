// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"log"
	"time"
)

func resourceAliCloudMongodbGlobalSecurityIPGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudMongodbGlobalSecurityIPGroupCreate,
		Read:   resourceAliCloudMongodbGlobalSecurityIPGroupRead,
		Update: resourceAliCloudMongodbGlobalSecurityIPGroupUpdate,
		Delete: resourceAliCloudMongodbGlobalSecurityIPGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"global_ig_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"global_security_ip_list": {
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

func resourceAliCloudMongodbGlobalSecurityIPGroupCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateGlobalSecurityIPGroup"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["RegionId"] = client.RegionId

	request["GIpList"] = d.Get("global_security_ip_list")
	request["GlobalIgName"] = d.Get("global_ig_name")
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("Dds", "2015-12-01", action, query, request, true)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_mongodb_global_security_ip_group", action, AlibabaCloudSdkGoERROR)
	}

	id, _ := jsonpath.Get("$.GlobalSecurityIPGroup[0].GlobalSecurityGroupId", response)
	d.SetId(fmt.Sprint(id))

	return resourceAliCloudMongodbGlobalSecurityIPGroupRead(d, meta)
}

func resourceAliCloudMongodbGlobalSecurityIPGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	mongodbServiceV2 := MongodbServiceV2{client}

	objectRaw, err := mongodbServiceV2.DescribeMongodbGlobalSecurityIPGroup(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_mongodb_global_security_ip_group DescribeMongodbGlobalSecurityIPGroup Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("global_ig_name", objectRaw["GlobalIgName"])
	d.Set("global_security_ip_list", objectRaw["GIpList"])
	d.Set("region_id", objectRaw["RegionId"])

	return nil
}

func resourceAliCloudMongodbGlobalSecurityIPGroupUpdate(d *schema.ResourceData, meta interface{}) error {
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
	if d.HasChange("global_security_ip_list") {
		update = true
	}
	request["GIpList"] = d.Get("global_security_ip_list")
	if d.HasChange("global_ig_name") {
		update = true
	}
	request["GlobalIgName"] = d.Get("global_ig_name")
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("Dds", "2015-12-01", action, query, request, true)
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

	return resourceAliCloudMongodbGlobalSecurityIPGroupRead(d, meta)
}

func resourceAliCloudMongodbGlobalSecurityIPGroupDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	action := "DeleteGlobalSecurityIPGroup"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["GlobalSecurityGroupId"] = d.Id()
	request["RegionId"] = client.RegionId

	request["GlobalIgName"] = d.Get("global_ig_name")
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("Dds", "2015-12-01", action, query, request, true)

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

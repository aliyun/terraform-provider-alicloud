// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package polardb

import (
	"fmt"
	"log"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceGlobalSecurityIpGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceGlobalSecurityIpGroupCreate,
		Read:   resourceGlobalSecurityIpGroupRead,
		Update: resourceGlobalSecurityIpGroupUpdate,
		Delete: resourceGlobalSecurityIpGroupDelete,
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

func resourceGlobalSecurityIpGroupCreate(d *schema.ResourceData, meta interface{}) error {

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
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		response, err = client.RpcPost("polardb", "2017-08-01", action, query, request, true)
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
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_polardb_global_security_ip_group", action, errs.AlibabaCloudSdkGoERROR)
	}

	id, _ := jsonpath.Get("$.GlobalSecurityIPGroup[0].GlobalSecurityGroupId", response)
	d.SetId(fmt.Sprint(id))

	return resourceGlobalSecurityIpGroupRead(d, meta)
}

func resourceGlobalSecurityIpGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	polardbServiceV2 := PolardbServiceV2{client}

	objectRaw, err := polardbServiceV2.DescribePolardbGlobalSecurityIpGroup(d.Id())
	if err != nil {
		if !d.IsNewResource() && errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_polardb_global_security_ip_group DescribePolardbGlobalSecurityIpGroup Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}

	d.Set("global_ip_group_name", objectRaw["GlobalIgName"])
	d.Set("global_ip_list", objectRaw["GIpList"])
	d.Set("region_id", objectRaw["RegionId"])

	return nil
}

func resourceGlobalSecurityIpGroupUpdate(d *schema.ResourceData, meta interface{}) error {
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
		wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RpcPost("polardb", "2017-08-01", action, query, request, true)
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
		wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RpcPost("polardb", "2017-08-01", action, query, request, true)
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

	d.Partial(false)
	return resourceGlobalSecurityIpGroupRead(d, meta)
}

func resourceGlobalSecurityIpGroupDelete(d *schema.ResourceData, meta interface{}) error {

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
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		response, err = client.RpcPost("polardb", "2017-08-01", action, query, request, true)
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
		if errs.IsExpectedErrors(err, []string{"SecAPICallingFailed"}) || errs.NotFoundError(err) {
			return nil
		}
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
	}

	return nil
}

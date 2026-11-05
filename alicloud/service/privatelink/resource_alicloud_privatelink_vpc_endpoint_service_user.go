// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package privatelink

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceVpcEndpointServiceUser() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudPrivateLinkVpcEndpointServiceUserCreate,
		Read:   resourceAliCloudPrivateLinkVpcEndpointServiceUserRead,
		Update: resourceAliCloudPrivateLinkVpcEndpointServiceUserUpdate,
		Delete: resourceAliCloudPrivateLinkVpcEndpointServiceUserDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"dry_run": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"service_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"user_arn": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"user_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAliCloudPrivateLinkVpcEndpointServiceUserCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "AddUserToVpcEndpointService"
	var request map[string]interface{}
	var response map[string]interface{}
	var err error
	query := make(map[string]interface{})
	request = make(map[string]interface{})
	request["UserId"] = d.Get("user_id")
	request["ServiceId"] = d.Get("service_id")
	request["RegionId"] = client.RegionId
	request["ClientToken"] = helper.BuildClientToken(action)

	if v, ok := d.GetOk("user_arn"); ok {
		request["UserARN"] = v
	}
	if v, ok := d.GetOkExists("dry_run"); ok {
		request["DryRun"] = v
	}
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		response, err = client.RpcPost("Privatelink", "2020-04-15", action, query, request, true)
		if err != nil {
			if errs.IsExpectedErrors(err, []string{"ConcurrentCallNotSupported"}) || helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)

	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_privatelink_vpc_endpoint_service_user", action, errs.AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%v:%v", request["ServiceId"], request["UserId"]))

	return resourceAliCloudPrivateLinkVpcEndpointServiceUserRead(d, meta)
}

func resourceAliCloudPrivateLinkVpcEndpointServiceUserRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	privateLinkServiceV2 := PrivateLinkServiceV2{client}

	objectRaw, err := privateLinkServiceV2.DescribePrivateLinkVpcEndpointServiceUser(d.Id())
	if err != nil {
		if !d.IsNewResource() && errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_privatelink_vpc_endpoint_service_user DescribePrivateLinkVpcEndpointServiceUser Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}

	if objectRaw["UserId"] != nil {
		d.Set("user_id", objectRaw["UserId"])
	}

	parts := strings.Split(d.Id(), ":")
	d.Set("service_id", parts[0])

	return nil
}

func resourceAliCloudPrivateLinkVpcEndpointServiceUserUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Cannot update resource Alicloud Resource Vpc Endpoint Service User.")
	return nil
}

func resourceAliCloudPrivateLinkVpcEndpointServiceUserDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	parts := strings.Split(d.Id(), ":")
	action := "RemoveUserFromVpcEndpointService"
	var request map[string]interface{}
	var response map[string]interface{}
	var err error
	query := make(map[string]interface{})
	request = make(map[string]interface{})
	request["UserId"] = parts[1]
	request["ServiceId"] = parts[0]
	request["RegionId"] = client.RegionId
	request["ClientToken"] = helper.BuildClientToken(action)

	if v, ok := d.GetOk("user_arn"); ok {
		request["UserARN"] = v
	}
	if v, ok := d.GetOkExists("dry_run"); ok {
		request["DryRun"] = v
	}
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		response, err = client.RpcPost("Privatelink", "2020-04-15", action, query, request, true)
		request["ClientToken"] = helper.BuildClientToken(action)

		if err != nil {
			if errs.IsExpectedErrors(err, []string{"ConcurrentCallNotSupported"}) || helper.NeedRetry(err) {
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

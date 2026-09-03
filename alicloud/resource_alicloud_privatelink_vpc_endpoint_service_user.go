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

func resourceAliCloudPrivatelinkVpcEndpointServiceUser() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudPrivatelinkVpcEndpointServiceUserCreate,
		Read:   resourceAliCloudPrivatelinkVpcEndpointServiceUserRead,
		Update: resourceAliCloudPrivatelinkVpcEndpointServiceUserUpdate,
		Delete: resourceAliCloudPrivatelinkVpcEndpointServiceUserDelete,
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
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				AtLeastOneOf: []string{"user_id", "user_arn"},
			},
			"user_id": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				AtLeastOneOf: []string{"user_id", "user_arn"},
			},
		},
	}
}

func resourceAliCloudPrivatelinkVpcEndpointServiceUserCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "AddUserToVpcEndpointService"
	var request map[string]interface{}
	var response map[string]interface{}
	var err error
	query := make(map[string]interface{})
	request = make(map[string]interface{})
	request["ServiceId"] = d.Get("service_id")
	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken(action)

	if v, ok := d.GetOk("user_id"); ok {
		request["UserId"] = v
	}
	if v, ok := d.GetOk("user_arn"); ok {
		request["UserARN"] = v
	}
	if v, ok := d.GetOkExists("dry_run"); ok {
		request["DryRun"] = v
	}
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("Privatelink", "2020-04-15", action, query, request, true)
		if err != nil {
			if IsExpectedErrors(err, []string{"ConcurrentCallNotSupported"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_privatelink_vpc_endpoint_service_user", action, AlibabaCloudSdkGoERROR)
	}

	serviceId := fmt.Sprint(request["ServiceId"])
	userId := ""
	if v, ok := d.GetOk("user_id"); ok {
		userId = v.(string)
	}
	userArn := ""
	if v, ok := d.GetOk("user_arn"); ok {
		userArn = v.(string)
	}
	if userArn != "" {
		d.SetId(fmt.Sprintf("%s:%s:%s", serviceId, userId, userArn))
	} else {
		d.SetId(fmt.Sprintf("%s:%s", serviceId, userId))
	}

	return resourceAliCloudPrivatelinkVpcEndpointServiceUserRead(d, meta)
}

func resourceAliCloudPrivatelinkVpcEndpointServiceUserRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	privatelinkServiceV2 := PrivatelinkServiceV2{client}

	objectRaw, err := privatelinkServiceV2.DescribePrivatelinkVpcEndpointServiceUser(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_privatelink_vpc_endpoint_service_user DescribePrivatelinkVpcEndpointServiceUser Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	serviceId, userId, userArn, err := parsePrivatelinkVpcEndpointServiceUserId(d.Id())
	if err != nil {
		return WrapError(err)
	}
	d.Set("service_id", serviceId)

	if userId != "" {
		d.Set("user_id", formatPrivatelinkUserId(objectRaw["UserId"]))
	}
	if userArn != "" {
		// The whitelist list API does not return UserARN for account ID entries,
		// and UserARN is immutable, so the value carried by the resource ID is authoritative.
		d.Set("user_arn", userArn)
	}

	return nil
}

func resourceAliCloudPrivatelinkVpcEndpointServiceUserUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Cannot update resource Alicloud Resource Vpc Endpoint Service User.")
	return nil
}

func resourceAliCloudPrivatelinkVpcEndpointServiceUserDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	serviceId, userId, userArn, err := parsePrivatelinkVpcEndpointServiceUserId(d.Id())
	if err != nil {
		return WrapError(err)
	}
	action := "RemoveUserFromVpcEndpointService"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	request = make(map[string]interface{})
	request["ServiceId"] = serviceId
	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken(action)

	if userId != "" {
		request["UserId"] = userId
	}
	if userArn != "" {
		request["UserARN"] = userArn
	}
	if v, ok := d.GetOkExists("dry_run"); ok {
		request["DryRun"] = v
	}
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("Privatelink", "2020-04-15", action, query, request, true)
		request["ClientToken"] = buildClientToken(action)

		if err != nil {
			if IsExpectedErrors(err, []string{"ConcurrentCallNotSupported"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		if IsExpectedErrors(err, []string{"EndpointServiceNotFound"}) || NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return nil
}

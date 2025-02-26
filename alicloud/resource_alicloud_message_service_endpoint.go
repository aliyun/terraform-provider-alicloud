// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"log"
	"time"
)

func resourceAliCloudMessageServiceEndpoint() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudMessageServiceEndpointCreate,
		Read:   resourceAliCloudMessageServiceEndpointRead,
		Update: resourceAliCloudMessageServiceEndpointUpdate,
		Delete: resourceAliCloudMessageServiceEndpointDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"endpoint_enabled": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"endpoint_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"public"}, false),
			},
		},
	}
}

func resourceAliCloudMessageServiceEndpointCreate(d *schema.ResourceData, meta interface{}) error {
	d.SetId(fmt.Sprint(d.Get("endpoint_type")))
	return resourceAliCloudMessageServiceEndpointUpdate(d, meta)
}

func resourceAliCloudMessageServiceEndpointRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	messageServiceServiceV2 := MessageServiceServiceV2{client}

	objectRaw, err := messageServiceServiceV2.DescribeMessageServiceEndpoint(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_message_service_endpoint DescribeMessageServiceEndpoint Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	if objectRaw["EndpointEnabled"] != nil {
		d.Set("endpoint_enabled", objectRaw["EndpointEnabled"])
	}

	d.Set("endpoint_type", d.Id())

	return nil
}

func resourceAliCloudMessageServiceEndpointUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	messageServiceServiceV2 := MessageServiceServiceV2{client}
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}

	object, err := messageServiceServiceV2.DescribeMessageServiceEndpoint(d.Id())
	if err != nil {
		return WrapError(err)
	}

	target, ok := d.GetOkExists("endpoint_enabled")
	if ok && fmt.Sprint(object["EndpointEnabled"]) != fmt.Sprint(target) {
		if target == true {
			action := "EnableEndpoint"
			request = make(map[string]interface{})
			query = make(map[string]interface{})
			request["EndpointType"] = d.Id()
			request["RegionId"] = client.RegionId
			wait := incrementalWait(3*time.Second, 5*time.Second)
			err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
				response, err = client.RpcPost("Mns-open", "2022-01-19", action, query, request, true)
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
		if target == false {
			action := "DisableEndpoint"
			request = make(map[string]interface{})
			query = make(map[string]interface{})
			request["EndpointType"] = d.Id()
			request["RegionId"] = client.RegionId
			wait := incrementalWait(3*time.Second, 5*time.Second)
			err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
				response, err = client.RpcPost("Mns-open", "2022-01-19", action, query, request, true)
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
	}

	return resourceAliCloudMessageServiceEndpointRead(d, meta)
}

func resourceAliCloudMessageServiceEndpointDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[WARN] Cannot destroy resource AliCloud Resource Endpoint. Terraform will remove this resource from the state file, however resources may remain.")
	return nil
}

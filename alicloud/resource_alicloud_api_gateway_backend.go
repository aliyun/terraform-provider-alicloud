package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudApiGatewayBackend() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudApiGatewayBackendCreate,
		Read:   resourceAlicloudApiGatewayBackendRead,
		Update: resourceAlicloudApiGatewayBackendUpdate,
		Delete: resourceAlicloudApiGatewayBackendDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"backend_name": {
				Required: true,
				Type:     schema.TypeString,
			},
			"backend_type": {
				Required:     true,
				ForceNew:     true,
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice([]string{"HTTP", "VPC", "FC_EVENT", "FC_EVENT_V3", "FC_HTTP", "FC_HTTP_V3", "OSS", "MOCK"}, false),
			},
			"create_event_bridge_service_linked_role": {
				Optional: true,
				Computed: true,
				Type:     schema.TypeBool,
			},
			"description": {
				Optional: true,
				Type:     schema.TypeString,
			},
		},
	}
}

func resourceAlicloudApiGatewayBackendCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	request := make(map[string]interface{})
	var err error

	if v, ok := d.GetOk("backend_name"); ok {
		request["BackendName"] = v
	}
	if v, ok := d.GetOk("backend_type"); ok {
		request["BackendType"] = v
	}
	if v, ok := d.GetOk("create_event_bridge_service_linked_role"); ok {
		request["CreateEventBridgeServiceLinkedRole"] = v
	}
	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}

	var response map[string]interface{}
	action := "CreateBackend"
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		resp, err := client.RpcPost("CloudAPI", "2016-07-14", action, nil, request, false)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		response = resp
		addDebug(action, response, request)
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_api_gateway_backend", action, AlibabaCloudSdkGoERROR)
	}

	if v, err := jsonpath.Get("$.BackendId", response); err != nil || v == nil {
		return WrapErrorf(err, IdMsg, "alicloud_api_gateway_backend")
	} else {
		d.SetId(fmt.Sprint(v))
	}

	return resourceAlicloudApiGatewayBackendRead(d, meta)
}

func resourceAlicloudApiGatewayBackendRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cloudApiService := CloudApiService{client}

	object, err := cloudApiService.DescribeApiGatewayBackend(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_api_gateway_backend cloudApiService.DescribeApiGatewayBackend Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("backend_name", object["BackendName"])
	d.Set("backend_type", object["BackendType"])
	d.Set("description", object["Description"])
	return nil
}

func resourceAlicloudApiGatewayBackendUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var err error
	update := false
	request := map[string]interface{}{
		"BackendId": d.Id(),
	}

	if d.HasChange("description") {
		update = true
		if v, ok := d.GetOk("description"); ok {
			request["Description"] = v
		}
	}

	if d.HasChange("backend_name") {
		update = true
	}
	request["BackendName"] = d.Get("backend_name")
	request["BackendType"] = d.Get("backend_type")
	if update {
		action := "ModifyBackend"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			resp, err := client.RpcPost("CloudAPI", "2016-07-14", action, nil, request, false)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(action, resp, request)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
	}

	return resourceAlicloudApiGatewayBackendRead(d, meta)
}

func resourceAlicloudApiGatewayBackendDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var err error

	request := map[string]interface{}{

		"BackendId": d.Id(),
	}

	action := "DeleteBackend"
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		resp, err := client.RpcPost("CloudAPI", "2016-07-14", action, nil, request, false)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, resp, request)
		return nil
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"NotFoundBackend"}) || NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	return nil
}

package alicloud

import (
	"fmt"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAlicloudServiceMeshExtensionProvider() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudServiceMeshExtensionProviderCreate,
		Read:   resourceAlicloudServiceMeshExtensionProviderRead,
		Update: resourceAlicloudServiceMeshExtensionProviderUpdate,
		Delete: resourceAlicloudServiceMeshExtensionProviderDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(3 * time.Minute),
			Update: schema.DefaultTimeout(3 * time.Minute),
			Delete: schema.DefaultTimeout(3 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"service_mesh_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"extension_provider_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"httpextauth", "grpcextauth"}, false),
			},
			"config": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceAlicloudServiceMeshExtensionProviderCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateExtensionProvider"
	request := make(map[string]interface{})
	var err error

	request["ServiceMeshId"] = d.Get("service_mesh_id")
	request["Type"] = d.Get("type")
	request["Name"] = d.Get("extension_provider_name")
	request["Config"] = d.Get("config")

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *resource.RetryError {
		response, err = client.RpcPost("servicemesh", "2020-01-11", action, nil, request, false)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_service_mesh_extension_provider", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%v:%v:%v", request["ServiceMeshId"], request["Type"], request["Name"]))

	return resourceAlicloudServiceMeshExtensionProviderRead(d, meta)
}

func resourceAlicloudServiceMeshExtensionProviderRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	servicemeshService := ServicemeshService{client}
	object, err := servicemeshService.DescribeServiceMeshExtensionProvider(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	parts, err := ParseResourceId(d.Id(), 3)
	if err != nil {
		return WrapError(err)
	}

	d.Set("service_mesh_id", parts[0])
	d.Set("type", object["Type"])
	d.Set("extension_provider_name", object["Name"])
	d.Set("config", object["Config"])

	return nil
}

func resourceAlicloudServiceMeshExtensionProviderUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	update := false
	parts, err := ParseResourceId(d.Id(), 3)
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"ServiceMeshId": parts[0],
		"Type":          parts[1],
		"Name":          parts[2],
	}

	if d.HasChange("config") {
		update = true
	}
	if v, ok := d.GetOk("config"); ok {
		request["Config"] = v
	}

	if update {
		action := "UpdateExtensionProvider"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			response, err = client.RpcPost("servicemesh", "2020-01-11", action, nil, request, false)
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

	return resourceAlicloudServiceMeshExtensionProviderRead(d, meta)
}

func resourceAlicloudServiceMeshExtensionProviderDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DeleteExtensionProvider"
	var response map[string]interface{}
	var err error

	parts, err := ParseResourceId(d.Id(), 3)
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"ServiceMeshId": parts[0],
		"Type":          parts[1],
		"Name":          parts[2],
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *resource.RetryError {
		response, err = client.RpcPost("servicemesh", "2020-01-11", action, nil, request, false)
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

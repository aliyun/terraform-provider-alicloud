package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudResourceManagerDelegatedAdministrator() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudResourceManagerDelegatedAdministratorCreate,
		Read:   resourceAlicloudResourceManagerDelegatedAdministratorRead,
		Delete: resourceAlicloudResourceManagerDelegatedAdministratorDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(1 * time.Minute),
			Delete: schema.DefaultTimeout(1 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"account_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"service_principal": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAlicloudResourceManagerDelegatedAdministratorCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "RegisterDelegatedAdministrator"
	request := make(map[string]interface{})
	var err error
	request["AccountId"] = d.Get("account_id")
	request["ServicePrincipal"] = d.Get("service_principal")
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("ResourceManager", "2020-03-31", action, nil, request, false)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_resource_manager_delegated_administrator", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(request["AccountId"], ":", request["ServicePrincipal"]))

	return resourceAlicloudResourceManagerDelegatedAdministratorRead(d, meta)
}
func resourceAlicloudResourceManagerDelegatedAdministratorRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	resourceManagerService := ResourceManagerService{client}
	_, err := resourceManagerService.DescribeResourceManagerDelegatedAdministrator(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_resource_manager_delegated_administrator resourceManagerService.DescribeResourceManagerDelegatedAdministrator Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	d.Set("account_id", parts[0])
	d.Set("service_principal", parts[1])
	return nil
}
func resourceAlicloudResourceManagerDelegatedAdministratorDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	action := "DeregisterDelegatedAdministrator"
	var response map[string]interface{}
	request := map[string]interface{}{
		"AccountId":        parts[0],
		"ServicePrincipal": parts[1],
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("ResourceManager", "2020-03-31", action, nil, request, false)
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
	return nil
}

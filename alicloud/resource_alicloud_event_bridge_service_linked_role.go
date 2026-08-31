package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudEventBridgeServiceLinkedRole() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudEventBridgeServiceLinkedRoleCreate,
		Read:   resourceAliCloudEventBridgeServiceLinkedRoleRead,
		Delete: resourceAliCloudEventBridgeServiceLinkedRoleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"product_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}
func resourceAliCloudEventBridgeServiceLinkedRoleCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	eventbridgeService := EventbridgeService{client}
	var response map[string]interface{}
	action := "CreateServiceLinkedRoleForProduct"
	request := make(map[string]interface{})
	var err error

	request["ProductName"] = d.Get("product_name")

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *resource.RetryError {
		response, err = client.RpcPost("eventbridge", "2020-04-01", action, nil, request, false)
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
		return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_event_bridge_service_linked_role", action, AlibabaCloudSdkGoERROR)
	}

	if fmt.Sprint(response["Code"]) != "Success" {
		return WrapError(fmt.Errorf("CreateServiceLinkedRoleForProduct failed, response: %v", response))
	}

	d.SetId(fmt.Sprint(request["ProductName"]))

	stateConf := BuildStateConf([]string{}, []string{"true"}, d.Timeout(schema.TimeoutCreate), 10*time.Second, eventbridgeService.CheckRoleForProductRefreshFunc(d.Id(), []string{"false"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAliCloudEventBridgeServiceLinkedRoleRead(d, meta)
}

func resourceAliCloudEventBridgeServiceLinkedRoleRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	eventbridgeService := EventbridgeService{client}

	object, err := eventbridgeService.DescribeEventBridgeServiceLinkedRole(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_event_bridge_event_bus eventbridgeService.DescribeEventBridgeServiceLinkedRole Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("product_name", object["StsRoleName"])

	return nil
}

func resourceAliCloudEventBridgeServiceLinkedRoleDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DeleteServiceLinkedRole"
	var response map[string]interface{}
	var err error

	request := map[string]interface{}{
		"RoleName": d.Id(),
	}

	wait := incrementalWait(3*time.Second, 0*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("ResourceManager", "2020-03-31", action, request, nil, false)
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
		if IsExpectedErrors(err, []string{"EntityNotExist.Role"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return nil
}

package alicloud

import (
	"fmt"
	"log"
	"time"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudCloudFirewallControlPolicyOrder() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudCloudFirewallControlPolicyOrderCreate,
		Read:   resourceAliCloudCloudFirewallControlPolicyOrderRead,
		Update: resourceAliCloudCloudFirewallControlPolicyOrderUpdate,
		Delete: resourceAliCloudCloudFirewallControlPolicyOrderDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"acl_uuid": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"direction": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"in", "out"}, false),
			},
			"order": {
				Type:     schema.TypeInt,
				Required: true,
			},
		},
	}
}

func resourceAliCloudCloudFirewallControlPolicyOrderCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "ModifyControlPolicyPriority"
	request := make(map[string]interface{})
	conn, err := client.NewCloudfwClient()
	if err != nil {
		return WrapError(err)
	}

	request["AclUuid"] = d.Get("acl_uuid")
	request["Direction"] = d.Get("direction")
	request["Order"] = d.Get("order")

	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-12-07"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}

		if fmt.Sprint(response["Message"]) == "not buy user" {
			conn.Endpoint = String(connectivity.CloudFirewallOpenAPIEndpointControlPolicy)
			return resource.RetryableError(fmt.Errorf("%s", response))
		}

		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_cloud_firewall_control_policy_order", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%v:%v", request["AclUuid"], request["Direction"]))

	return resourceAliCloudCloudFirewallControlPolicyOrderRead(d, meta)
}

func resourceAliCloudCloudFirewallControlPolicyOrderRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cloudfwService := CloudfwService{client}

	object, err := cloudfwService.DescribeCloudFirewallControlPolicy(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_cloud_firewall_control_policy_order cloudfwService.DescribeCloudFirewallControlPolicy Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("acl_uuid", object["AclUuid"])
	d.Set("direction", object["Direction"])
	d.Set("order", formatInt(object["Order"]))

	return nil
}

func resourceAliCloudCloudFirewallControlPolicyOrderUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	update := false

	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}

	request := map[string]interface{}{
		"AclUuid":   parts[0],
		"Direction": parts[1],
	}

	if d.HasChange("order") {
		update = true
	}
	request["Order"] = d.Get("order")

	if update {
		action := "ModifyControlPolicyPriority"
		conn, err := client.NewCloudfwClient()
		if err != nil {
			return WrapError(err)
		}

		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-12-07"), StringPointer("AK"), nil, request, &runtime)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}

			if fmt.Sprint(response["Message"]) == "not buy user" {
				conn.Endpoint = String(connectivity.CloudFirewallOpenAPIEndpointControlPolicy)
				return resource.RetryableError(fmt.Errorf("%s", response))
			}

			return nil
		})
		addDebug(action, response, request)

		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "alicloud_cloud_firewall_control_policy_order", action, AlibabaCloudSdkGoERROR)
		}
	}

	return resourceAliCloudCloudFirewallControlPolicyOrderRead(d, meta)
}

func resourceAliCloudCloudFirewallControlPolicyOrderDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Resource alicloud_cloud_firewall_control_policy_order [%s]  will not be deleted", d.Id())

	return nil
}

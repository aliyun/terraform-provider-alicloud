package cloudfirewall

import (
	"fmt"
	"log"
	"time"

	connectivity "github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	convert "github.com/aliyun/terraform-provider-alicloud/alicloud/convert"
	errs "github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	helper "github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	logs "github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	validators "github.com/aliyun/terraform-provider-alicloud/alicloud/validators"
	retry "github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	schema "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceCloudFirewallControlPolicyOrder() *schema.Resource {
	return &schema.Resource{
		Create: resourceControlPolicyOrderCreate,
		Read:   resourceControlPolicyOrderRead,
		Update: resourceControlPolicyOrderUpdate,
		Delete: resourceControlPolicyOrderDelete,
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
				ValidateFunc: validators.StringInSlice([]string{"in", "out"}, false),
			},
			"order": {
				Type:     schema.TypeInt,
				Required: true,
			},
		},
	}
}

func resourceControlPolicyOrderCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	var err error
	var endpoint string
	action := "ModifyControlPolicyPriority"
	request := make(map[string]interface{})

	request["AclUuid"] = d.Get("acl_uuid")
	request["Direction"] = d.Get("direction")
	request["Order"] = d.Get("order")

	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *retry.RetryError {
		response, err = client.RpcPostWithEndpoint("Cloudfw", "2017-12-07", action, nil, request, false, endpoint)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			} else if errs.IsExpectedErrors(err, []string{"not buy user"}) {
				endpoint = connectivity.CloudFirewallOpenAPIEndpointControlPolicy
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}

		return nil
	})
	logs.AddDebug(action, response, request)

	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_cloud_firewall_control_policy_order", action, errs.AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%v:%v", request["AclUuid"], request["Direction"]))

	return resourceControlPolicyOrderRead(d, meta)
}

func resourceControlPolicyOrderRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cloudfwService := NewService(client)

	object, err := cloudfwService.DescribeCloudFirewallControlPolicy(d.Id())
	if err != nil {
		if !d.IsNewResource() && errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_cloud_firewall_control_policy_order cloudfwService.DescribeCloudFirewallControlPolicy Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}

	d.Set("acl_uuid", object["AclUuid"])
	d.Set("direction", object["Direction"])
	orderVal, err := convert.FormatInt(object["Order"])
	if err != nil {
		return err
	}
	d.Set("order", orderVal)

	return nil
}

func resourceControlPolicyOrderUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	update := false

	parts, err := helper.ParseResourceId(d.Id(), 2)
	if err != nil {
		return errs.WrapError(err)
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
		var endpoint string
		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *retry.RetryError {
			response, err = client.RpcPostWithEndpoint("Cloudfw", "2017-12-07", action, nil, request, false, endpoint)
			if err != nil {
				if helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				} else if errs.IsExpectedErrors(err, []string{"not buy user"}) {
					endpoint = connectivity.CloudFirewallOpenAPIEndpointControlPolicy
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}

			return nil
		})
		logs.AddDebug(action, response, request)

		if err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_cloud_firewall_control_policy_order", action, errs.AlibabaCloudSdkGoERROR)
		}
	}

	return resourceControlPolicyOrderRead(d, meta)
}

func resourceControlPolicyOrderDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Resource alicloud_cloud_firewall_control_policy_order [%s]  will not be deleted", d.Id())

	return nil
}

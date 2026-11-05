package cloudfirewall

import (
	"fmt"
	"log"
	"strings"
	"time"

	connectivity "github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	errs "github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	helper "github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	logs "github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	retry "github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	schema "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceCloudFirewallVpcFirewallControlPolicyOrder() *schema.Resource {
	return &schema.Resource{
		Create: resourceVpcFirewallControlPolicyOrderCreate,
		Read:   resourceVpcFirewallControlPolicyOrderRead,
		Update: resourceVpcFirewallControlPolicyOrderUpdate,
		Delete: resourceVpcFirewallControlPolicyOrderDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"acl_uuid": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"lang": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"order": {
				Type:     schema.TypeString,
				Required: true,
			},
			"vpc_firewall_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceVpcFirewallControlPolicyOrderCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "ModifyVpcFirewallControlPolicyPosition"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	if v, ok := d.GetOk("vpc_firewall_id"); ok {
		request["VpcFirewallId"] = v
	}
	if v, ok := d.GetOk("acl_uuid"); ok {
		request["AclUuid"] = v
	}

	request["NewOrder"] = d.Get("order")
	wait := helper.IncrementalWait(3*time.Second, 0*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		response, err = client.RpcPost("Cloudfw", "2017-12-07", action, query, request, true)
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
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_cloud_firewall_vpc_firewall_control_policy_order", action, errs.AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%v:%v", request["VpcFirewallId"], request["AclUuid"]))

	return resourceVpcFirewallControlPolicyOrderRead(d, meta)
}

func resourceVpcFirewallControlPolicyOrderRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cloudFirewallServiceV2 := NewService(client)

	objectRaw, err := cloudFirewallServiceV2.DescribeCloudFirewallVpcFirewallControlPolicyOrder(d.Id())
	if err != nil {
		if !d.IsNewResource() && errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_cloud_firewall_vpc_firewall_control_policy_order DescribeCloudFirewallVpcFirewallControlPolicyOrder Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}

	d.Set("order", objectRaw["Order"])

	parts := strings.Split(d.Id(), ":")
	d.Set("vpc_firewall_id", parts[0])
	d.Set("acl_uuid", parts[1])

	return nil
}

func resourceVpcFirewallControlPolicyOrderUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false

	var err error
	parts := strings.Split(d.Id(), ":")
	action := "ModifyVpcFirewallControlPolicyPosition"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["VpcFirewallId"] = parts[0]
	request["AclUuid"] = parts[1]

	if d.HasChange("order") {
		update = true
	}
	request["NewOrder"] = d.Get("order")
	if update {
		wait := helper.IncrementalWait(3*time.Second, 0*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RpcPost("Cloudfw", "2017-12-07", action, query, request, true)
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

	return resourceVpcFirewallControlPolicyOrderRead(d, meta)
}

func resourceVpcFirewallControlPolicyOrderDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[WARN] Cannot destroy resource AliCloud Resource Vpc Firewall Control Policy Order. Terraform will remove this resource from the state file, however resources may remain.")
	return nil
}

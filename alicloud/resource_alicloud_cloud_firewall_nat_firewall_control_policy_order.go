package alicloud

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudCloudFirewallNatFirewallControlPolicyOrder() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudCloudFirewallNatFirewallControlPolicyOrderCreate,
		Read:   resourceAliCloudCloudFirewallNatFirewallControlPolicyOrderRead,
		Update: resourceAliCloudCloudFirewallNatFirewallControlPolicyOrderUpdate,
		Delete: resourceAliCloudCloudFirewallNatFirewallControlPolicyOrderDelete,
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
				Required: true,
				ForceNew: true,
			},
			"current_page": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"direction": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"nat_gateway_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"order": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceAliCloudCloudFirewallNatFirewallControlPolicyOrderCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "ModifyNatFirewallControlPolicyPosition"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	if v, ok := d.GetOk("nat_gateway_id"); ok {
		request["NatGatewayId"] = v
	}
	if v, ok := d.GetOk("direction"); ok {
		request["Direction"] = v
	}
	if v, ok := d.GetOk("acl_uuid"); ok {
		request["AclUuid"] = v
	}

	request["NewOrder"] = d.Get("order")
	wait := incrementalWait(3*time.Second, 0*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("Cloudfw", "2017-12-07", action, query, request, true)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_cloud_firewall_nat_firewall_control_policy_order", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%v:%v:%v", request["AclUuid"], request["NatGatewayId"], request["Direction"]))

	return resourceAliCloudCloudFirewallNatFirewallControlPolicyOrderRead(d, meta)
}

func resourceAliCloudCloudFirewallNatFirewallControlPolicyOrderRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cloudFirewallServiceV2 := CloudFirewallServiceV2{client}

	objectRaw, err := cloudFirewallServiceV2.DescribeCloudFirewallNatFirewallControlPolicyOrder(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_cloud_firewall_nat_firewall_control_policy_order DescribeCloudFirewallNatFirewallControlPolicyOrder Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("order", objectRaw["Order"])

	parts := strings.Split(d.Id(), ":")
	d.Set("acl_uuid", parts[0])
	d.Set("nat_gateway_id", parts[1])
	d.Set("direction", parts[2])

	return nil
}

func resourceAliCloudCloudFirewallNatFirewallControlPolicyOrderUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false

	var err error
	parts := strings.Split(d.Id(), ":")
	action := "ModifyNatFirewallControlPolicyPosition"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["NatGatewayId"] = parts[1]
	request["Direction"] = parts[2]
	request["AclUuid"] = parts[0]

	if d.HasChange("order") {
		update = true
	}
	request["NewOrder"] = d.Get("order")
	if update {
		wait := incrementalWait(3*time.Second, 0*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("Cloudfw", "2017-12-07", action, query, request, true)
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

	return resourceAliCloudCloudFirewallNatFirewallControlPolicyOrderRead(d, meta)
}

func resourceAliCloudCloudFirewallNatFirewallControlPolicyOrderDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[WARN] Cannot destroy resource AliCloud Resource Nat Firewall Control Policy Order. Terraform will remove this resource from the state file, however resources may remain.")
	return nil
}

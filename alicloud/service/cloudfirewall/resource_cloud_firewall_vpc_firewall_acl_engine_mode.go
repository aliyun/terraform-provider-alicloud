package cloudfirewall

import (
	"fmt"
	"log"
	"time"

	jsonpath "github.com/PaesslerAG/jsonpath"
	connectivity "github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	errs "github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	helper "github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	logs "github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	retry "github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	schema "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	validation "github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourceCloudFirewallVpcFirewallAclEngineMode() *schema.Resource {
	return &schema.Resource{
		Create: resourceVpcFirewallAclEngineModeCreate,
		Read:   resourceVpcFirewallAclEngineModeRead,
		Update: resourceVpcFirewallAclEngineModeUpdate,
		Delete: resourceVpcFirewallAclEngineModeDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"strict_mode": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntInSlice([]int{0, 1}),
			},
			"vpc_firewall_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"member_uid": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

func resourceVpcFirewallAclEngineModeCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "ModifyVpcFirewallAclEngineMode"
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error

	request := map[string]interface{}{
		"VpcFirewallId": d.Get("vpc_firewall_id"),
		"StrictMode":    d.Get("strict_mode"),
	}
	if v, ok := d.GetOk("member_uid"); ok {
		request["MemberUid"] = v
	}

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
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
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_cloud_firewall_vpc_firewall_acl_engine_mode", action, errs.AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(request["VpcFirewallId"]))

	return resourceVpcFirewallAclEngineModeRead(d, meta)
}

func resourceVpcFirewallAclEngineModeRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cloudFirewallServiceV2 := NewService(client)

	objectRaw, err := cloudFirewallServiceV2.DescribeCloudFirewallVpcFirewallAclEngineMode(d.Id())
	if err != nil {
		if !d.IsNewResource() && errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_cloud_firewall_vpc_firewall_acl_engine_mode DescribeCloudFirewallVpcFirewallAclEngineMode Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}

	d.Set("member_uid", objectRaw["MemberUid"])
	d.Set("vpc_firewall_id", objectRaw["AclGroupId"])

	aclConfigRawObj, _ := jsonpath.Get("$.AclConfig", objectRaw)
	aclConfigRaw := make(map[string]interface{})
	if aclConfigRawObj != nil {
		aclConfigRaw = aclConfigRawObj.(map[string]interface{})
	}
	d.Set("strict_mode", aclConfigRaw["StrictMode"])
	return nil
}

func resourceVpcFirewallAclEngineModeUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false

	var err error
	action := "ModifyVpcFirewallAclEngineMode"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["VpcFirewallId"] = d.Id()

	if d.HasChange("strict_mode") {
		update = true
		request["StrictMode"] = d.Get("strict_mode")
	}

	if update {
		wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
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

	return resourceVpcFirewallAclEngineModeRead(d, meta)
}

func resourceVpcFirewallAclEngineModeDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[WARN] Cannot destroy resource AliCloud Resource Vpc Firewall Acl Engine Mode. Terraform will remove this resource from the state file, however resources may remain.")
	return nil
}

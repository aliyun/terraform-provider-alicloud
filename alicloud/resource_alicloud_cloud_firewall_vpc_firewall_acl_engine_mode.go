package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAliCloudCloudFirewallVpcFirewallAclEngineMode() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudCloudFirewallVpcFirewallAclEngineModeCreate,
		Read:   resourceAliCloudCloudFirewallVpcFirewallAclEngineModeRead,
		Update: resourceAliCloudCloudFirewallVpcFirewallAclEngineModeUpdate,
		Delete: resourceAliCloudCloudFirewallVpcFirewallAclEngineModeDelete,
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

func resourceAliCloudCloudFirewallVpcFirewallAclEngineModeCreate(d *schema.ResourceData, meta interface{}) error {
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

	wait := incrementalWait(3*time.Second, 5*time.Second)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_cloud_firewall_vpc_firewall_acl_engine_mode", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(request["VpcFirewallId"]))

	return resourceAliCloudCloudFirewallVpcFirewallAclEngineModeRead(d, meta)
}

func resourceAliCloudCloudFirewallVpcFirewallAclEngineModeRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cloudFirewallServiceV2 := CloudFirewallServiceV2{client}

	objectRaw, err := cloudFirewallServiceV2.DescribeCloudFirewallVpcFirewallAclEngineMode(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_cloud_firewall_vpc_firewall_acl_engine_mode DescribeCloudFirewallVpcFirewallAclEngineMode Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
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

func resourceAliCloudCloudFirewallVpcFirewallAclEngineModeUpdate(d *schema.ResourceData, meta interface{}) error {
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
		wait := incrementalWait(3*time.Second, 5*time.Second)
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

	return resourceAliCloudCloudFirewallVpcFirewallAclEngineModeRead(d, meta)
}

func resourceAliCloudCloudFirewallVpcFirewallAclEngineModeDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[WARN] Cannot destroy resource AliCloud Resource Vpc Firewall Acl Engine Mode. Terraform will remove this resource from the state file, however resources may remain.")
	return nil
}

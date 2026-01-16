// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudCloudFirewallVpcFirewallIpsConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudCloudFirewallVpcFirewallIpsConfigCreate,
		Read:   resourceAliCloudCloudFirewallVpcFirewallIpsConfigRead,
		Update: resourceAliCloudCloudFirewallVpcFirewallIpsConfigUpdate,
		Delete: resourceAliCloudCloudFirewallVpcFirewallIpsConfigDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"basic_rules": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: IntInSlice([]int{1, 0}),
			},
			"enable_all_patch": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: IntInSlice([]int{1, 0}),
			},
			"lang": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"member_uid": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"rule_class": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"run_mode": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: IntInSlice([]int{1, 0}),
			},
			"vpc_firewall_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAliCloudCloudFirewallVpcFirewallIpsConfigCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "ModifyVpcFirewallDefaultIPSConfig"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	if v, ok := d.GetOk("vpc_firewall_id"); ok {
		request["VpcFirewallId"] = v
	}

	request["BasicRules"] = d.Get("basic_rules")
	request["EnableAllPatch"] = d.Get("enable_all_patch")
	request["RunMode"] = d.Get("run_mode")
	if v, ok := d.GetOk("member_uid"); ok {
		request["MemberUid"] = v
	}
	if v, ok := d.GetOk("rule_class"); ok {
		request["RuleClass"] = v
	}
	if v, ok := d.GetOk("lang"); ok {
		request["Lang"] = v
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_cloud_firewall_vpc_firewall_ips_config", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(request["VpcFirewallId"]))
	return resourceAliCloudCloudFirewallVpcFirewallIpsConfigRead(d, meta)
}

func resourceAliCloudCloudFirewallVpcFirewallIpsConfigRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cloudFirewallServiceV2 := CloudFirewallServiceV2{client}

	objectRaw, err := cloudFirewallServiceV2.DescribeCloudFirewallVpcFirewallIpsConfig(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_cloud_firewall_vpc_firewall_ips_config DescribeCloudFirewallVpcFirewallIpsConfig Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("basic_rules", objectRaw["BasicRules"])
	d.Set("enable_all_patch", objectRaw["EnableAllPatch"])
	d.Set("rule_class", objectRaw["RuleClass"])
	d.Set("run_mode", objectRaw["RunMode"])
	d.Set("vpc_firewall_id", d.Id())
	return nil
}

func resourceAliCloudCloudFirewallVpcFirewallIpsConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false

	var err error
	action := "ModifyVpcFirewallDefaultIPSConfig"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["VpcFirewallId"] = d.Id()

	if v, ok := d.GetOk("lang"); ok {
		request["Lang"] = v
	}
	if d.HasChange("basic_rules") {
		update = true
	}
	request["BasicRules"] = d.Get("basic_rules")
	if d.HasChange("enable_all_patch") {
		update = true
	}
	request["EnableAllPatch"] = d.Get("enable_all_patch")
	if d.HasChange("run_mode") {
		update = true
	}
	request["RunMode"] = d.Get("run_mode")
	if v, ok := d.GetOk("member_uid"); ok {
		request["MemberUid"] = v
	}
	if d.HasChange("rule_class") {
		update = true
		request["RuleClass"] = d.Get("rule_class")
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

	return resourceAliCloudCloudFirewallVpcFirewallIpsConfigRead(d, meta)
}

func resourceAliCloudCloudFirewallVpcFirewallIpsConfigDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[WARN] Cannot destroy resource AliCloud Resource Vpc Firewall Ips Config. Terraform will remove this resource from the state file, however resources may remain.")
	return nil
}

package alicloud

import (
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudCloudFirewallIPSConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudCloudFirewallIPSConfigCreate,
		Read:   resourceAliCloudCloudFirewallIPSConfigRead,
		Update: resourceAliCloudCloudFirewallIPSConfigUpdate,
		Delete: resourceAliCloudCloudFirewallIPSConfigDelete,
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
				Type:     schema.TypeInt,
				Optional: true,
			},
			"cti_rules": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"lang": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"max_sdl": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"patch_rules": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"rule_class": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"run_mode": {
				Type:     schema.TypeInt,
				Optional: true,
			},
		},
	}
}

func resourceAliCloudCloudFirewallIPSConfigCreate(d *schema.ResourceData, meta interface{}) error {
	accountId, err := meta.(*connectivity.AliyunClient).AccountId()
	if err != nil {
		return err
	}
	d.SetId(accountId)
	return resourceAliCloudCloudFirewallIPSConfigUpdate(d, meta)
}

func resourceAliCloudCloudFirewallIPSConfigRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cloudFirewallServiceV2 := CloudFirewallServiceV2{client}

	objectRaw, err := cloudFirewallServiceV2.DescribeCloudFirewallIPSConfig(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_cloud_firewall_ips_config DescribeCloudFirewallIPSConfig Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("basic_rules", objectRaw["BasicRules"])
	d.Set("cti_rules", objectRaw["CtiRules"])
	d.Set("max_sdl", objectRaw["MaxSdl"])
	d.Set("patch_rules", objectRaw["PatchRules"])
	d.Set("rule_class", objectRaw["RuleClass"])
	d.Set("run_mode", objectRaw["RunMode"])

	return nil
}

func resourceAliCloudCloudFirewallIPSConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false

	var err error
	action := "ModifyDefaultIPSConfig"
	request = make(map[string]interface{})
	query = make(map[string]interface{})

	if d.HasChange("cti_rules") {
		update = true
	}
	query["CtiRules"] = d.Get("cti_rules")
	if d.HasChange("patch_rules") {
		update = true
	}
	query["PatchRules"] = d.Get("patch_rules")
	if v, ok := d.GetOk("lang"); ok {
		query["Lang"] = v
	}
	if d.HasChange("basic_rules") {
		update = true
	}
	if v, ok := d.GetOk("basic_rules"); ok || (d.IsNewResource() || d.HasChange("basic_rules")) {
		query["BasicRules"] = v
	}
	if d.HasChange("run_mode") {
		update = true
	}
	query["RunMode"] = d.Get("run_mode")
	if d.HasChange("max_sdl") {
		update = true
	}
	if v, ok := d.GetOk("max_sdl"); ok || (d.IsNewResource() || d.HasChange("max_sdl")) {
		query["MaxSdl"] = v
	}
	if d.HasChange("rule_class") {
		update = true
	}
	if v, ok := d.GetOk("rule_class"); ok || (d.IsNewResource() || d.HasChange("rule_class")) {
		query["RuleClass"] = v
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

	return resourceAliCloudCloudFirewallIPSConfigRead(d, meta)
}

func resourceAliCloudCloudFirewallIPSConfigDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[WARN] Cannot destroy resource AliCloud Resource I P S Config. Terraform will remove this resource from the state file, however resources may remain.")
	return nil
}

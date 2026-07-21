// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudRamPasswordPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudRamPasswordPolicyCreate,
		Read:   resourceAliCloudRamPasswordPolicyRead,
		Update: resourceAliCloudRamPasswordPolicyUpdate,
		Delete: resourceAliCloudRamPasswordPolicyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"hard_expiry": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"max_login_attemps": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"max_password_age": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"minimum_password_different_character": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"minimum_password_length": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"password_not_contain_user_name": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"password_reuse_prevention": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"require_lowercase_characters": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"require_numbers": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"require_symbols": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"require_uppercase_characters": {
				Type:     schema.TypeBool,
				Optional: true,
			},
		},
	}
}

func resourceAliCloudRamPasswordPolicyCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "SetPasswordPolicy"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})

	if v, ok := d.GetOkExists("minimum_password_length"); ok {
		request["MinimumPasswordLength"] = v
	}
	if v, ok := d.GetOkExists("require_lowercase_characters"); ok {
		request["RequireLowercaseCharacters"] = v
	}
	if v, ok := d.GetOkExists("require_uppercase_characters"); ok {
		request["RequireUppercaseCharacters"] = v
	}
	if v, ok := d.GetOkExists("require_numbers"); ok {
		request["RequireNumbers"] = v
	}
	if v, ok := d.GetOkExists("require_symbols"); ok {
		request["RequireSymbols"] = v
	}
	if v, ok := d.GetOkExists("hard_expiry"); ok {
		request["HardExpire"] = v
	}
	if v, ok := d.GetOkExists("max_login_attemps"); ok {
		request["MaxLoginAttemps"] = v
	}
	if v, ok := d.GetOkExists("password_reuse_prevention"); ok {
		request["PasswordReusePrevention"] = v
	}
	if v, ok := d.GetOkExists("max_password_age"); ok {
		request["MaxPasswordAge"] = v
	}
	if v, ok := d.GetOkExists("minimum_password_different_character"); ok {
		request["MinimumPasswordDifferentCharacter"] = v
	}
	if v, ok := d.GetOkExists("password_not_contain_user_name"); ok {
		request["PasswordNotContainUserName"] = v
	}
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("Ims", "2019-08-15", action, query, request, true)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_ram_password_policy", action, AlibabaCloudSdkGoERROR)
	}

	accountId, err := client.AccountId()
	d.SetId(accountId)

	return resourceAliCloudRamPasswordPolicyRead(d, meta)
}

func resourceAliCloudRamPasswordPolicyRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ramServiceV2 := RamServiceV2{client}

	objectRaw, err := ramServiceV2.DescribeRamPasswordPolicy(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_ram_password_policy DescribeRamPasswordPolicy Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("hard_expiry", objectRaw["HardExpire"])
	d.Set("max_login_attemps", objectRaw["MaxLoginAttemps"])
	d.Set("max_password_age", objectRaw["MaxPasswordAge"])
	d.Set("minimum_password_different_character", objectRaw["MinimumPasswordDifferentCharacter"])
	d.Set("minimum_password_length", objectRaw["MinimumPasswordLength"])
	d.Set("password_not_contain_user_name", objectRaw["PasswordNotContainUserName"])
	d.Set("password_reuse_prevention", objectRaw["PasswordReusePrevention"])
	d.Set("require_lowercase_characters", objectRaw["RequireLowercaseCharacters"])
	d.Set("require_numbers", objectRaw["RequireNumbers"])
	d.Set("require_symbols", objectRaw["RequireSymbols"])
	d.Set("require_uppercase_characters", objectRaw["RequireUppercaseCharacters"])

	return nil
}

func resourceAliCloudRamPasswordPolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false

	var err error
	action := "SetPasswordPolicy"
	request = make(map[string]interface{})
	query = make(map[string]interface{})

	if d.HasChange("minimum_password_length") {
		update = true
		request["MinimumPasswordLength"] = d.Get("minimum_password_length")
	}

	if d.HasChange("require_lowercase_characters") {
		update = true
		request["RequireLowercaseCharacters"] = d.Get("require_lowercase_characters")
	}

	if d.HasChange("require_uppercase_characters") {
		update = true
		request["RequireUppercaseCharacters"] = d.Get("require_uppercase_characters")
	}

	if d.HasChange("require_numbers") {
		update = true
		request["RequireNumbers"] = d.Get("require_numbers")
	}

	if d.HasChange("require_symbols") {
		update = true
		request["RequireSymbols"] = d.Get("require_symbols")
	}

	if d.HasChange("hard_expiry") {
		update = true
		request["HardExpire"] = d.Get("hard_expiry")
	}

	if d.HasChange("max_login_attemps") {
		update = true
		request["MaxLoginAttemps"] = d.Get("max_login_attemps")
	}

	if d.HasChange("password_reuse_prevention") {
		update = true
		request["PasswordReusePrevention"] = d.Get("password_reuse_prevention")
	}

	if d.HasChange("max_password_age") {
		update = true
		request["MaxPasswordAge"] = d.Get("max_password_age")
	}

	if d.HasChange("minimum_password_different_character") {
		update = true
		request["MinimumPasswordDifferentCharacter"] = d.Get("minimum_password_different_character")
	}

	if d.HasChange("password_not_contain_user_name") {
		update = true
		request["PasswordNotContainUserName"] = d.Get("password_not_contain_user_name")
	}

	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("Ims", "2019-08-15", action, query, request, true)
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

	return resourceAliCloudRamPasswordPolicyRead(d, meta)
}

func resourceAliCloudRamPasswordPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[WARN] Cannot destroy resource AliCloud Resource Password Policy. Terraform will remove this resource from the state file, however resources may remain.")
	return nil
}

package alicloud

import (
	"fmt"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ram"
	"github.com/hashicorp/terraform/helper/schema"

	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

// defaults from API docs https://help.aliyun.com/document_detail/28739.html
// when resource is created it sets not specifies value in the resource block to defaults
// also during deletion it rollbacks changes to defaults (API has only a Set method)
var (
	default_minimum_password_length      = 12
	default_require_lowercase_characters = true
	default_require_uppercase_characters = true
	default_require_numbers              = true
	default_require_symbols              = true
	default_hard_expiry                  = false
	default_max_password_age             = 0 // means disable
	default_password_reuse_prevention    = 0 // means disable
	default_max_login_attempts           = 5
)

func resourceAlicloudRamAccountPasswordPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudRamAccountPasswordPolicyUpdate,
		Read:   resourceAlicloudRamAccountPasswordPolicyRead,
		Update: resourceAlicloudRamAccountPasswordPolicyUpdate,
		Delete: resourceAlicloudRamAccountPasswordPolicyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"minimum_password_length": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      default_minimum_password_length,
				ValidateFunc: intBetween(8, 32),
			},
			"require_lowercase_characters": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  default_require_lowercase_characters,
			},
			"require_uppercase_characters": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  default_require_uppercase_characters,
			},
			"require_numbers": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  default_require_numbers,
			},
			"require_symbols": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  default_require_symbols,
			},
			"hard_expiry": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  default_hard_expiry,
			},
			"max_password_age": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      default_max_password_age,
				ValidateFunc: intBetween(0, 1095),
			},
			"password_reuse_prevention": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      default_password_reuse_prevention,
				ValidateFunc: intBetween(0, 24),
			},
			"max_login_attempts": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      default_max_login_attempts,
				ValidateFunc: intBetween(0, 32),
			},
		},
	}
}

func resourceAlicloudRamAccountPasswordPolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := ram.CreateSetPasswordPolicyRequest()
	if v, ok := d.GetOk("minimum_password_length"); ok {
		request.MinimumPasswordLength = requests.NewInteger(v.(int))
	}
	// due to lack of GetOkExists() in old Terraform included in vendor dir
	// GetOk() doesn't distinguish beetwen false/0 and empty, so in else clause set 0 or false
	if v, ok := d.GetOk("require_lowercase_characters"); ok {
		request.RequireLowercaseCharacters = requests.NewBoolean(v.(bool))
	} else {
		request.RequireLowercaseCharacters = requests.NewBoolean(false)
	}
	if v, ok := d.GetOk("require_uppercase_characters"); ok {
		request.RequireUppercaseCharacters = requests.NewBoolean(v.(bool))
	} else {
		request.RequireUppercaseCharacters = requests.NewBoolean(false)
	}
	if v, ok := d.GetOk("require_numbers"); ok {
		request.RequireNumbers = requests.NewBoolean(v.(bool))
	} else {
		request.RequireNumbers = requests.NewBoolean(false)
	}
	if v, ok := d.GetOk("require_symbols"); ok {
		request.RequireSymbols = requests.NewBoolean(v.(bool))
	} else {
		request.RequireSymbols = requests.NewBoolean(false)
	}
	if v, ok := d.GetOk("max_login_attempts"); ok {
		request.MaxLoginAttemps = requests.NewInteger(v.(int))
	} else {
		request.MaxLoginAttemps = requests.NewInteger(0)
	}
	if v, ok := d.GetOk("hard_expiry"); ok {
		request.HardExpiry = requests.NewBoolean(v.(bool))
	} else {
		request.HardExpiry = requests.NewBoolean(false)
	}
	if v, ok := d.GetOk("max_password_age"); ok {
		request.MaxPasswordAge = requests.NewInteger(v.(int))
	} else {
		request.MaxPasswordAge = requests.NewInteger(0)
	}
	if v, ok := d.GetOk("password_reuse_prevention"); ok {
		request.PasswordReusePrevention = requests.NewInteger(v.(int))
	} else {
		request.PasswordReusePrevention = requests.NewInteger(0)
	}
	raw, err := client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
		return ramClient.SetPasswordPolicy(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_ram_account_password_policy", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw)

	d.SetId("ram-account-password-policy")

	return resourceAlicloudRamAccountPasswordPolicyRead(d, meta)
}

func resourceAlicloudRamAccountPasswordPolicyRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := ram.CreateGetPasswordPolicyRequest()
	raw, err := client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
		return ramClient.GetPasswordPolicy(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_ram_account_password_policy", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw)
	response, _ := raw.(*ram.GetPasswordPolicyResponse)
	passwordPolicy := response.PasswordPolicy
	d.Set("minimum_password_length", passwordPolicy.MinimumPasswordLength)
	d.Set("require_lowercase_characters", passwordPolicy.RequireLowercaseCharacters)
	d.Set("require_uppercase_characters", passwordPolicy.RequireUppercaseCharacters)
	d.Set("require_numbers", passwordPolicy.RequireNumbers)
	d.Set("require_symbols", passwordPolicy.RequireSymbols)
	d.Set("hard_expiry", passwordPolicy.HardExpiry)
	d.Set("max_password_age", passwordPolicy.MaxPasswordAge)
	d.Set("password_reuse_prevention", passwordPolicy.PasswordReusePrevention)
	d.Set("max_login_attempts", passwordPolicy.MaxLoginAttemps)
	return nil
}

func resourceAlicloudRamAccountPasswordPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := ram.CreateSetPasswordPolicyRequest()
	request.MinimumPasswordLength = requests.NewInteger(default_minimum_password_length)
	request.RequireLowercaseCharacters = requests.NewBoolean(default_require_lowercase_characters)
	request.RequireUppercaseCharacters = requests.NewBoolean(default_require_uppercase_characters)
	request.RequireNumbers = requests.NewBoolean(default_require_numbers)
	request.RequireSymbols = requests.NewBoolean(default_require_symbols)
	request.HardExpiry = requests.NewBoolean(default_hard_expiry)
	request.MaxPasswordAge = requests.NewInteger(default_max_password_age)
	request.PasswordReusePrevention = requests.NewInteger(default_password_reuse_prevention)
	request.MaxLoginAttemps = requests.NewInteger(default_max_login_attempts)
	raw, err := client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
		return ramClient.SetPasswordPolicy(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_ram_account_password_policy", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw)
	return nil
}

// below copy/pasta from https://github.com/hashicorp/terraform/blob/master/helper/validation/validation.go
// alicloud vendor contains very old version of Terraform which lacks this functions

// IntBetween returns a SchemaValidateFunc which tests if the provided value
// is of type int and is between min and max (inclusive)
func intBetween(min, max int) schema.SchemaValidateFunc {
	return func(i interface{}, k string) (s []string, es []error) {
		v, ok := i.(int)
		if !ok {
			es = append(es, fmt.Errorf("expected type of %s to be int", k))
			return
		}

		if v < min || v > max {
			es = append(es, fmt.Errorf("expected %s to be in the range (%d - %d), got %d", k, min, max, v))
			return
		}

		return
	}
}

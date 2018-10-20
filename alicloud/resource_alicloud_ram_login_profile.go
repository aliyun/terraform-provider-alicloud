package alicloud

import (
	"fmt"
	"time"

	"github.com/denverdino/aliyungo/ram"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func resourceAlicloudRamLoginProfile() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudRamLoginProfileCreate,
		Read:   resourceAlicloudRamLoginProfileRead,
		Update: resourceAlicloudRamLoginProfileUpdate,
		Delete: resourceAlicloudRamLoginProfileDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"user_name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateRamName,
			},
			"password": &schema.Schema{
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			},
			"password_reset_required": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"mfa_bind_required": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
}

func resourceAlicloudRamLoginProfileCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	args := ram.ProfileRequest{
		UserName:              d.Get("user_name").(string),
		Password:              d.Get("password").(string),
		PasswordResetRequired: d.Get("password_reset_required").(bool),
		MFABindRequired:       d.Get("mfa_bind_required").(bool),
	}

	_, err := client.WithRamClient(func(ramClient ram.RamClientInterface) (interface{}, error) {
		return ramClient.CreateLoginProfile(args)
	})
	if err != nil {
		return fmt.Errorf("CreateLoginProfile got an error: %#v", err)
	}

	d.SetId(args.UserName)
	return resourceAlicloudRamLoginProfileUpdate(d, meta)
}

func resourceAlicloudRamLoginProfileUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	d.Partial(true)

	args := ram.ProfileRequest{
		UserName: d.Id(),
		Password: d.Get("password").(string),
	}
	attributeUpdate := false

	if d.HasChange("password") {
		d.SetPartial("password")
		attributeUpdate = true
	}

	if d.HasChange("password_reset_required") {
		d.SetPartial("password_reset_required")
		args.PasswordResetRequired = d.Get("password_reset_required").(bool)
		attributeUpdate = true
	}

	if d.HasChange("mfa_bind_required") {
		d.SetPartial("mfa_bind_required")
		args.MFABindRequired = d.Get("mfa_bind_required").(bool)
		attributeUpdate = true
	}

	if attributeUpdate && !d.IsNewResource() {
		_, err := client.WithRamClient(func(ramClient ram.RamClientInterface) (interface{}, error) {
			return ramClient.UpdateLoginProfile(args)
		})
		if err != nil {
			return fmt.Errorf("UpdateLoginProfile got an error: %v", err)
		}
	}

	d.Partial(false)
	return resourceAlicloudRamLoginProfileRead(d, meta)
}

func resourceAlicloudRamLoginProfileRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	args := ram.UserQueryRequest{
		UserName: d.Id(),
	}

	raw, err := client.WithRamClient(func(ramClient ram.RamClientInterface) (interface{}, error) {
		return ramClient.GetLoginProfile(args)
	})
	if err != nil {
		if RamEntityNotExist(err) {
			d.SetId("")
		}
		return fmt.Errorf("GetLoginProfile got an error: %#v", err)
	}
	response, _ := raw.(ram.ProfileResponse)
	profile := response.LoginProfile
	d.Set("user_name", profile.UserName)
	d.Set("mfa_bind_required", profile.MFABindRequired)
	d.Set("password_reset_required", profile.PasswordResetRequired)
	return nil
}

func resourceAlicloudRamLoginProfileDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	args := ram.UserQueryRequest{
		UserName: d.Id(),
	}

	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		_, err := client.WithRamClient(func(ramClient ram.RamClientInterface) (interface{}, error) {
			return ramClient.DeleteLoginProfile(args)
		})
		if err != nil {
			if RamEntityNotExist(err) {
				return nil
			}
			return resource.NonRetryableError(fmt.Errorf("Error deleting login profile: %#v", err))
		}

		raw, err := client.WithRamClient(func(ramClient ram.RamClientInterface) (interface{}, error) {
			return ramClient.GetLoginProfile(args)
		})
		if err != nil {
			if RamEntityNotExist(err) {
				return nil
			}

			return resource.NonRetryableError(err)
		}
		response, _ := raw.(ram.ProfileResponse)
		if response.LoginProfile.UserName == args.UserName {
			return resource.RetryableError(fmt.Errorf("Error deleting login profile - trying again while it is deleted."))
		}
		return nil
	})
}

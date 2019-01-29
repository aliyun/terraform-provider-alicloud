package alicloud

import (
	"strconv"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ram"
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
			"user_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateRamName,
			},
			"password": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			},
			"password_reset_required": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"mfa_bind_required": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
}

func resourceAlicloudRamLoginProfileCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := ram.CreateCreateLoginProfileRequest()
	request.UserName = d.Get("user_name").(string)
	request.Password = d.Get("password").(string)
	request.PasswordResetRequired = requests.Boolean(strconv.FormatBool(d.Get("password_reset_required").(bool)))
	request.MFABindRequired = requests.Boolean(strconv.FormatBool(d.Get("mfa_bind_required").(bool)))

	_, err := client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
		return ramClient.CreateLoginProfile(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "ram_login_profile", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	d.SetId(request.UserName)
	return resourceAlicloudRamLoginProfileUpdate(d, meta)
}

func resourceAlicloudRamLoginProfileUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	d.Partial(true)

	request := ram.CreateUpdateLoginProfileRequest()
	request.Password = d.Get("password").(string)
	request.UserName = d.Id()

	attributeUpdate := false

	if d.HasChange("password") {
		d.SetPartial("password")
		attributeUpdate = true
	}

	if d.HasChange("password_reset_required") {
		d.SetPartial("password_reset_required")
		request.PasswordResetRequired = requests.Boolean(strconv.FormatBool(d.Get("password_reset_required").(bool)))
		attributeUpdate = true
	}

	if d.HasChange("mfa_bind_required") {
		d.SetPartial("mfa_bind_required")
		request.MFABindRequired = requests.Boolean(strconv.FormatBool(d.Get("mfa_bind_required").(bool)))
		attributeUpdate = true
	}

	if attributeUpdate && !d.IsNewResource() {
		_, err := client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
			return ramClient.UpdateLoginProfile(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
	}

	d.Partial(false)
	return resourceAlicloudRamLoginProfileRead(d, meta)
}

func resourceAlicloudRamLoginProfileRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := ram.CreateGetLoginProfileRequest()
	request.UserName = d.Id()

	raw, err := client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
		return ramClient.GetLoginProfile(request)
	})
	if err != nil {
		if RamEntityNotExist(err) {
			d.SetId("")
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	response, _ := raw.(*ram.GetLoginProfileResponse)
	profile := response.LoginProfile
	d.Set("user_name", profile.UserName)
	d.Set("mfa_bind_required", profile.MFABindRequired)
	d.Set("password_reset_required", profile.PasswordResetRequired)
	return nil
}

func resourceAlicloudRamLoginProfileDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := ram.CreateDeleteLoginProfileRequest()
	request.UserName = d.Id()

	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		_, err := client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
			return ramClient.DeleteLoginProfile(request)
		})
		if err != nil {
			if RamEntityNotExist(err) {
				return nil
			}
			return resource.NonRetryableError(WrapError(err))
		}

		raw, err := client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
			request := ram.CreateGetLoginProfileRequest()
			request.UserName = d.Id()
			return ramClient.GetLoginProfile(request)
		})
		if err != nil {
			if RamEntityNotExist(err) {
				return nil
			}

			return resource.NonRetryableError(WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR))
		}
		response, _ := raw.(*ram.GetLoginProfileResponse)
		if response.LoginProfile.UserName == request.UserName {
			return resource.RetryableError(WrapErrorf(err, DeleteTimeoutMsg, d.Id(), request.GetActionName(), ProviderERROR))
		}
		return nil
	})
}

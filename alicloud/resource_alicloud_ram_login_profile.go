package alicloud

import (
	"strconv"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ram"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudRamLoginProfile() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudRamLoginProfileCreate,
		Read:   resourceAliCloudRamLoginProfileRead,
		Update: resourceAliCloudRamLoginProfileUpdate,
		Delete: resourceAliCloudRamLoginProfileDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"user_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
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
				Computed: true,
			},
		},
	}
}

func resourceAliCloudRamLoginProfileCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ramSercvice := RamService{client}
	request := ram.CreateCreateLoginProfileRequest()
	request.RegionId = client.RegionId
	request.UserName = d.Get("user_name").(string)
	request.Password = d.Get("password").(string)
	if v, ok := d.GetOk("password_reset_required"); ok {
		request.PasswordResetRequired = requests.Boolean(strconv.FormatBool(v.(bool)))
	}
	if v, ok := d.GetOk("mfa_bind_required"); ok {
		request.MFABindRequired = requests.Boolean(strconv.FormatBool(v.(bool)))
	}

	var raw interface{}
	var err error
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *resource.RetryError {
		raw, err = client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
			return ramClient.CreateLoginProfile(request)
		})
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_ram_login_profile", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	d.SetId(request.UserName)

	err = ramSercvice.WaitForRamLoginProfile(d.Id(), Normal, DefaultTimeout)
	if err != nil {
		return WrapError(err)
	}

	return resourceAliCloudRamLoginProfileRead(d, meta)
}

func resourceAliCloudRamLoginProfileUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := ram.CreateUpdateLoginProfileRequest()
	request.RegionId = client.RegionId
	request.UserName = d.Id()
	request.Password = d.Get("password").(string)

	if d.HasChange("password_reset_required") {
		request.PasswordResetRequired = requests.Boolean(strconv.FormatBool(d.Get("password_reset_required").(bool)))
	}

	if d.HasChange("mfa_bind_required") {
		request.MFABindRequired = requests.Boolean(strconv.FormatBool(d.Get("mfa_bind_required").(bool)))
	}

	var raw interface{}
	var err error
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
		raw, err = client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
			return ramClient.UpdateLoginProfile(request)
		})
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	return resourceAliCloudRamLoginProfileRead(d, meta)
}

func resourceAliCloudRamLoginProfileRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ramService := RamService{client}
	object, err := ramService.DescribeRamLoginProfile(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	profile := object.LoginProfile
	d.Set("user_name", profile.UserName)
	d.Set("mfa_bind_required", profile.MFABindRequired)
	d.Set("password_reset_required", profile.PasswordResetRequired)

	return nil
}

func resourceAliCloudRamLoginProfileDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ramService := RamService{client}
	request := ram.CreateDeleteLoginProfileRequest()
	request.RegionId = client.RegionId
	request.UserName = d.Id()

	var raw interface{}
	var err error
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *resource.RetryError {
		raw, err = client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
			return ramClient.DeleteLoginProfile(request)
		})
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)

	if err != nil {
		if IsExpectedErrors(err, []string{"EntityNotExist.User", "EntityNotExist.User.LoginProfile"}) || NeedRetry(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	return WrapError(ramService.WaitForRamLoginProfile(d.Id(), Deleted, DefaultTimeout))

}

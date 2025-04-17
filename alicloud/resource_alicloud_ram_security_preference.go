package alicloud

import (
	"encoding/json"
	"log"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudRamSecurityPreference() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudRamSecurityPreferenceCreate,
		Read:   resourceAliCloudRamSecurityPreferenceRead,
		Update: resourceAliCloudRamSecurityPreferenceUpdate,
		Delete: resourceAliCloudRamSecurityPreferenceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"allow_user_to_change_password": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"allow_user_to_login_with_passkey": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"allow_user_to_manage_access_keys": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"allow_user_to_manage_mfa_devices": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"allow_user_to_manage_personal_ding_talk": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"enable_save_mfa_ticket": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"login_network_masks": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"login_session_duration": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"mfa_operation_for_login": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"operation_for_risk_login": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"verification_types": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"enforce_mfa_for_login": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func resourceAliCloudRamSecurityPreferenceCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "SetSecurityPreference"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})

	if v, ok := d.GetOkExists("allow_user_to_change_password"); ok {
		request["AllowUserToChangePassword"] = v
	}
	if v, ok := d.GetOkExists("allow_user_to_manage_access_keys"); ok {
		request["AllowUserToManageAccessKeys"] = v
	}
	if v, ok := d.GetOkExists("login_session_duration"); ok {
		request["LoginSessionDuration"] = v
	}
	if v, ok := d.GetOk("login_network_masks"); ok {
		request["LoginNetworkMasks"] = v
	}
	if v, ok := d.GetOk("verification_types"); ok {
		verificationTypesMapsArray := v.(*schema.Set).List()
		verificationTypesMapsJson, err := json.Marshal(verificationTypesMapsArray)
		if err != nil {
			return WrapError(err)
		}
		request["VerificationTypes"] = string(verificationTypesMapsJson)
	}

	if v, ok := d.GetOkExists("allow_user_to_manage_personal_ding_talk"); ok {
		request["AllowUserToManagePersonalDingTalk"] = v
	}
	if v, ok := d.GetOk("operation_for_risk_login"); ok {
		request["OperationForRiskLogin"] = v
	}
	if v, ok := d.GetOkExists("enable_save_mfa_ticket"); ok {
		request["EnableSaveMFATicket"] = v
	}
	if v, ok := d.GetOkExists("allow_user_to_manage_mfa_devices"); ok {
		request["AllowUserToManageMFADevices"] = v
	}
	if v, ok := d.GetOk("mfa_operation_for_login"); ok {
		request["MFAOperationForLogin"] = v
	}
	if v, ok := d.GetOkExists("allow_user_to_login_with_passkey"); ok {
		request["AllowUserToLoginWithPasskey"] = v
	}
	if v, ok := d.GetOkExists("enforce_mfa_for_login"); ok {
		request["EnforceMFAForLogin"] = v
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_ram_security_preference", action, AlibabaCloudSdkGoERROR)
	}

	accountId, err := client.AccountId()
	d.SetId(accountId)

	return resourceAliCloudRamSecurityPreferenceRead(d, meta)
}

func resourceAliCloudRamSecurityPreferenceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ramServiceV2 := RamServiceV2{client}

	objectRaw, err := ramServiceV2.DescribeRamSecurityPreference(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_ram_security_preference DescribeRamSecurityPreference Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	accessKeyPreferenceRawObj, _ := jsonpath.Get("$.AccessKeyPreference", objectRaw)
	accessKeyPreferenceRaw := make(map[string]interface{})
	if accessKeyPreferenceRawObj != nil {
		accessKeyPreferenceRaw = accessKeyPreferenceRawObj.(map[string]interface{})
	}
	d.Set("allow_user_to_manage_access_keys", accessKeyPreferenceRaw["AllowUserToManageAccessKeys"])

	loginProfilePreferenceRawObj, _ := jsonpath.Get("$.LoginProfilePreference", objectRaw)
	loginProfilePreferenceRaw := make(map[string]interface{})
	if loginProfilePreferenceRawObj != nil {
		loginProfilePreferenceRaw = loginProfilePreferenceRawObj.(map[string]interface{})
	}
	d.Set("allow_user_to_change_password", loginProfilePreferenceRaw["AllowUserToChangePassword"])
	d.Set("allow_user_to_login_with_passkey", loginProfilePreferenceRaw["AllowUserToLoginWithPasskey"])
	d.Set("enable_save_mfa_ticket", loginProfilePreferenceRaw["EnableSaveMFATicket"])
	d.Set("login_network_masks", loginProfilePreferenceRaw["LoginNetworkMasks"])
	d.Set("login_session_duration", loginProfilePreferenceRaw["LoginSessionDuration"])
	d.Set("mfa_operation_for_login", loginProfilePreferenceRaw["MFAOperationForLogin"])
	d.Set("operation_for_risk_login", loginProfilePreferenceRaw["OperationForRiskLogin"])
	d.Set("enforce_mfa_for_login", loginProfilePreferenceRaw["EnforceMFAForLogin"])

	mFAPreferenceRawObj, _ := jsonpath.Get("$.MFAPreference", objectRaw)
	mFAPreferenceRaw := make(map[string]interface{})
	if mFAPreferenceRawObj != nil {
		mFAPreferenceRaw = mFAPreferenceRawObj.(map[string]interface{})
	}
	d.Set("allow_user_to_manage_mfa_devices", mFAPreferenceRaw["AllowUserToManageMFADevices"])

	personalInfoPreferenceRawObj, _ := jsonpath.Get("$.PersonalInfoPreference", objectRaw)
	personalInfoPreferenceRaw := make(map[string]interface{})
	if personalInfoPreferenceRawObj != nil {
		personalInfoPreferenceRaw = personalInfoPreferenceRawObj.(map[string]interface{})
	}
	d.Set("allow_user_to_manage_personal_ding_talk", personalInfoPreferenceRaw["AllowUserToManagePersonalDingTalk"])

	verificationTypesRaw, _ := jsonpath.Get("$.VerificationPreference.VerificationTypes", objectRaw)
	d.Set("verification_types", verificationTypesRaw)

	return nil
}

func resourceAliCloudRamSecurityPreferenceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false

	var err error
	action := "SetSecurityPreference"
	request = make(map[string]interface{})
	query = make(map[string]interface{})

	if d.HasChange("allow_user_to_change_password") {
		update = true
	}
	if v, ok := d.GetOkExists("allow_user_to_change_password"); ok {
		request["AllowUserToChangePassword"] = v
	}

	if d.HasChange("allow_user_to_manage_access_keys") {
		update = true
	}
	if v, ok := d.GetOkExists("allow_user_to_manage_access_keys"); ok {
		request["AllowUserToManageAccessKeys"] = v
	}

	if d.HasChange("login_session_duration") {
		update = true
	}
	if v, ok := d.GetOk("login_session_duration"); ok {
		request["LoginSessionDuration"] = v
	}

	if d.HasChange("login_network_masks") {
		update = true
		request["LoginNetworkMasks"] = d.Get("login_network_masks")
	}
	if v, ok := d.GetOk("login_network_masks"); ok {
		request["LoginNetworkMasks"] = v
	}

	if d.HasChange("verification_types") {
		update = true
	}
	if v, ok := d.GetOk("verification_types"); ok || d.HasChange("verification_types") {
		verificationTypesMapsArray := v.(*schema.Set).List()
		verificationTypesMapsJson, err := json.Marshal(verificationTypesMapsArray)
		if err != nil {
			return WrapError(err)
		}
		request["VerificationTypes"] = string(verificationTypesMapsJson)
	}

	if d.HasChange("allow_user_to_manage_personal_ding_talk") {
		update = true
	}
	if v, ok := d.GetOkExists("allow_user_to_manage_personal_ding_talk"); ok {
		request["AllowUserToManagePersonalDingTalk"] = v
	}

	if d.HasChange("operation_for_risk_login") {
		update = true
	}
	if v, ok := d.GetOk("operation_for_risk_login"); ok {
		request["OperationForRiskLogin"] = v
	}

	if d.HasChange("enable_save_mfa_ticket") {
		update = true
	}
	if v, ok := d.GetOkExists("enable_save_mfa_ticket"); ok {
		request["EnableSaveMFATicket"] = v
	}

	if d.HasChange("allow_user_to_manage_mfa_devices") {
		update = true
	}
	if v, ok := d.GetOkExists("allow_user_to_manage_mfa_devices"); ok {
		request["AllowUserToManageMFADevices"] = v
	}

	if d.HasChange("mfa_operation_for_login") {
		update = true
	}
	if v, ok := d.GetOk("mfa_operation_for_login"); ok {
		request["MFAOperationForLogin"] = v
	}

	if d.HasChange("allow_user_to_login_with_passkey") {
		update = true
	}
	if v, ok := d.GetOkExists("allow_user_to_login_with_passkey"); ok {
		request["AllowUserToLoginWithPasskey"] = v
	}

	if d.HasChange("enforce_mfa_for_login") {
		update = true
	}
	if v, ok := d.GetOkExists("enforce_mfa_for_login"); ok {
		request["EnforceMFAForLogin"] = v
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

	return resourceAliCloudRamSecurityPreferenceRead(d, meta)
}

func resourceAliCloudRamSecurityPreferenceDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[WARN] Cannot destroy resource AliCloud Resource Security Preference. Terraform will remove this resource from the state file, however resources may remain.")
	return nil
}

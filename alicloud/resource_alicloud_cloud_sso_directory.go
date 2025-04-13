// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/blues/jsonata-go"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudCloudSSODirectory() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudCloudSSODirectoryCreate,
		Read:   resourceAliCloudCloudSSODirectoryRead,
		Update: resourceAliCloudCloudSSODirectoryUpdate,
		Delete: resourceAliCloudCloudSSODirectoryDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"directory_global_access_status": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: StringInSlice([]string{"Disabled", "Enabled"}, false),
			},
			"directory_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"login_preference": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"login_network_masks": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"allow_user_to_get_credentials": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
			"mfa_authentication_status": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: StringInSlice([]string{"Enabled", "Disabled"}, false),
			},
			"mfa_authentication_setting_info": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"operation_for_risk_login": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"mfa_authentication_advance_settings": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
			"password_policy": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"min_password_different_chars": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
						"hard_expire": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"require_numbers": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"password_not_contain_username": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
						"max_password_age": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
						"password_reuse_prevention": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
						"require_lower_case_chars": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"max_password_length": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"min_password_length": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
						"require_upper_case_chars": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"max_login_attempts": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
						"require_symbols": {
							Type:     schema.TypeBool,
							Computed: true,
						},
					},
				},
			},
			"saml_identity_provider_configuration": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"entity_id": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							ForceNew: true,
						},
						"encoded_metadata_document": {
							Type:      schema.TypeString,
							Optional:  true,
							Sensitive: true,
						},
						"binding_type": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							ForceNew: true,
						},
						"create_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"want_request_signed": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
							ForceNew: true,
						},
						"update_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"certificate_ids": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"sso_status": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: StringInSlice([]string{"Disabled", "Enabled"}, false),
						},
						"login_url": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							ForceNew: true,
						},
					},
				},
			},
			"scim_synchronization_status": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: StringInSlice([]string{"Enabled", "Disabled"}, false),
			},
			"saml_service_provider": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"entity_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"certificate_type": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"authn_sign_algo": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"encoded_metadata_document": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"support_encrypted_assertion": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
						"acs_url": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"user_provisioning_configuration": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"session_duration": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"default_landing_page": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func resourceAliCloudCloudSSODirectoryCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateDirectory"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})

	if v, ok := d.GetOk("directory_name"); ok {
		request["DirectoryName"] = v
	}
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("cloudsso", "2021-05-15", action, query, request, true)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_cloud_sso_directory", action, AlibabaCloudSdkGoERROR)
	}

	id, _ := jsonpath.Get("$.Directory.DirectoryId", response)
	d.SetId(fmt.Sprint(id))

	return resourceAliCloudCloudSSODirectoryUpdate(d, meta)
}

func resourceAliCloudCloudSSODirectoryRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cloudSSOServiceV2 := CloudSSOServiceV2{client}

	objectRaw, err := cloudSSOServiceV2.DescribeCloudSSODirectory(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_cloud_sso_directory DescribeCloudSSODirectory Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	samlServiceProviderMaps := make([]map[string]interface{}, 0)
	samlServiceProviderMap := make(map[string]interface{})

	samlServiceProviderMap["acs_url"] = objectRaw["AcsUrl"]
	samlServiceProviderMap["authn_sign_algo"] = objectRaw["AuthnSignAlgo"]
	samlServiceProviderMap["certificate_type"] = objectRaw["CertificateType"]
	samlServiceProviderMap["encoded_metadata_document"] = objectRaw["EncodedMetadataDocument"]
	samlServiceProviderMap["entity_id"] = objectRaw["EntityId"]
	samlServiceProviderMap["support_encrypted_assertion"] = objectRaw["SupportEncryptedAssertion"]

	samlServiceProviderMaps = append(samlServiceProviderMaps, samlServiceProviderMap)
	if err := d.Set("saml_service_provider", samlServiceProviderMaps); err != nil {
		return err
	}

	objectRaw, err = cloudSSOServiceV2.DescribeDirectoryGetPasswordPolicy(d.Id())
	if err != nil && !NotFoundError(err) {
		return WrapError(err)
	}

	passwordPolicyMaps := make([]map[string]interface{}, 0)
	passwordPolicyMap := make(map[string]interface{})

	passwordPolicyMap["hard_expire"] = objectRaw["HardExpire"]
	passwordPolicyMap["max_login_attempts"] = objectRaw["MaxLoginAttempts"]
	passwordPolicyMap["max_password_age"] = objectRaw["MaxPasswordAge"]
	passwordPolicyMap["max_password_length"] = objectRaw["MaxPasswordLength"]
	passwordPolicyMap["min_password_different_chars"] = objectRaw["MinPasswordDifferentChars"]
	passwordPolicyMap["min_password_length"] = objectRaw["MinPasswordLength"]
	passwordPolicyMap["password_not_contain_username"] = objectRaw["PasswordNotContainUsername"]
	passwordPolicyMap["password_reuse_prevention"] = objectRaw["PasswordReusePrevention"]
	passwordPolicyMap["require_lower_case_chars"] = objectRaw["RequireLowerCaseChars"]
	passwordPolicyMap["require_numbers"] = objectRaw["RequireNumbers"]
	passwordPolicyMap["require_symbols"] = objectRaw["RequireSymbols"]
	passwordPolicyMap["require_upper_case_chars"] = objectRaw["RequireUpperCaseChars"]

	passwordPolicyMaps = append(passwordPolicyMaps, passwordPolicyMap)
	if err := d.Set("password_policy", passwordPolicyMaps); err != nil {
		return err
	}

	objectRaw, err = cloudSSOServiceV2.DescribeDirectoryGetExternalSAMLIdentityProvider(d.Id())
	if err != nil && !NotFoundError(err) {
		return WrapError(err)
	}

	sAMLIdentityProviderConfigurationMaps := make([]map[string]interface{}, 0)
	sAMLIdentityProviderConfigurationMap := make(map[string]interface{})

	sAMLIdentityProviderConfigurationMap["binding_type"] = objectRaw["BindingType"]
	sAMLIdentityProviderConfigurationMap["create_time"] = objectRaw["CreateTime"]
	sAMLIdentityProviderConfigurationMap["encoded_metadata_document"] = objectRaw["EncodedMetadataDocument"]
	sAMLIdentityProviderConfigurationMap["entity_id"] = objectRaw["EntityId"]
	sAMLIdentityProviderConfigurationMap["login_url"] = objectRaw["LoginUrl"]
	sAMLIdentityProviderConfigurationMap["sso_status"] = objectRaw["SSOStatus"]
	sAMLIdentityProviderConfigurationMap["update_time"] = objectRaw["UpdateTime"]
	sAMLIdentityProviderConfigurationMap["want_request_signed"] = objectRaw["WantRequestSigned"]

	certificateIdsRaw, _ := jsonpath.Get("$.CertificateIds", objectRaw)
	sAMLIdentityProviderConfigurationMap["certificate_ids"] = certificateIdsRaw
	sAMLIdentityProviderConfigurationMaps = append(sAMLIdentityProviderConfigurationMaps, sAMLIdentityProviderConfigurationMap)
	if err := d.Set("saml_identity_provider_configuration", sAMLIdentityProviderConfigurationMaps); err != nil {
		return err
	}

	objectRaw, err = cloudSSOServiceV2.DescribeDirectoryGetDirectory(d.Id())
	if err != nil && !NotFoundError(err) {
		return WrapError(err)
	}

	d.Set("create_time", objectRaw["CreateTime"])
	d.Set("directory_name", objectRaw["DirectoryName"])

	objectRaw, err = cloudSSOServiceV2.DescribeDirectoryGetUserProvisioningConfiguration(d.Id())
	if err != nil && !NotFoundError(err) {
		return WrapError(err)
	}

	userProvisioningConfigurationMaps := make([]map[string]interface{}, 0)
	userProvisioningConfigurationMap := make(map[string]interface{})

	userProvisioningConfigurationMap["default_landing_page"] = objectRaw["DefaultLandingPage"]
	userProvisioningConfigurationMap["session_duration"] = objectRaw["SessionDuration"]

	userProvisioningConfigurationMaps = append(userProvisioningConfigurationMaps, userProvisioningConfigurationMap)
	if err := d.Set("user_provisioning_configuration", userProvisioningConfigurationMaps); err != nil {
		return err
	}

	objectRaw, err = cloudSSOServiceV2.DescribeDirectoryGetMFAAuthenticationSettingInfo(d.Id())
	if err != nil && !NotFoundError(err) {
		return WrapError(err)
	}

	mfaAuthenticationSettingInfoMaps := make([]map[string]interface{}, 0)
	mfaAuthenticationSettingInfoMap := make(map[string]interface{})

	mfaAuthenticationSettingInfoMap["mfa_authentication_advance_settings"] = objectRaw["MfaAuthenticationAdvanceSettings"]
	mfaAuthenticationSettingInfoMap["operation_for_risk_login"] = objectRaw["OperationForRiskLogin"]

	mfaAuthenticationSettingInfoMaps = append(mfaAuthenticationSettingInfoMaps, mfaAuthenticationSettingInfoMap)
	if err := d.Set("mfa_authentication_setting_info", mfaAuthenticationSettingInfoMaps); err != nil {
		return err
	}

	objectRaw, err = cloudSSOServiceV2.DescribeDirectoryGetLoginPreference(d.Id())
	if err != nil && !NotFoundError(err) {
		return WrapError(err)
	}

	loginPreferenceMaps := make([]map[string]interface{}, 0)
	loginPreferenceMap := make(map[string]interface{})

	loginPreferenceMap["allow_user_to_get_credentials"] = objectRaw["AllowUserToGetCredentials"]
	loginPreferenceMap["login_network_masks"] = objectRaw["LoginNetworkMasks"]

	loginPreferenceMaps = append(loginPreferenceMaps, loginPreferenceMap)
	if err := d.Set("login_preference", loginPreferenceMaps); err != nil {
		return err
	}

	objectRaw, err = cloudSSOServiceV2.DescribeDirectoryGetSCIMSynchronizationStatus(d.Id())
	if err != nil && !NotFoundError(err) {
		return WrapError(err)
	}

	d.Set("scim_synchronization_status", objectRaw["SCIMSynchronizationStatus"])

	objectRaw, err = cloudSSOServiceV2.DescribeDirectoryGetMFAAuthenticationStatus(d.Id())
	if err != nil && !NotFoundError(err) {
		return WrapError(err)
	}

	d.Set("mfa_authentication_status", objectRaw["MFAAuthenticationStatus"])

	objectRaw, err = cloudSSOServiceV2.DescribeDirectoryGetDirectoryGlobalAccessStatus(d.Id())
	if err != nil && !NotFoundError(err) {
		return WrapError(err)
	}

	e := jsonata.MustCompile("$.GlobalAccessStatus.status")
	evaluation, _ := e.Eval(objectRaw)
	d.Set("directory_global_access_status", evaluation)

	return nil
}

func resourceAliCloudCloudSSODirectoryUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false
	d.Partial(true)

	if d.HasChange("saml_identity_provider_configuration.0.encoded_metadata_document") {
		var err error

		target := d.Get("saml_identity_provider_configuration.0.encoded_metadata_document").(string)
		if target == "" {
			action := "ClearExternalSAMLIdentityProvider"
			request = make(map[string]interface{})
			query = make(map[string]interface{})
			request["DirectoryId"] = d.Id()

			wait := incrementalWait(3*time.Second, 5*time.Second)
			err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
				response, err = client.RpcPost("cloudsso", "2021-05-15", action, query, request, true)
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
	}

	var err error
	action := "UpdateDirectory"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["DirectoryId"] = d.Id()

	if !d.IsNewResource() && d.HasChange("directory_name") {
		update = true
		request["NewDirectoryName"] = d.Get("directory_name")
	}

	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("cloudsso", "2021-05-15", action, query, request, true)
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
	update = false
	action = "SetExternalSAMLIdentityProvider"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["DirectoryId"] = d.Id()

	if d.HasChange("saml_identity_provider_configuration.0.encoded_metadata_document") {
		update = true
		jsonPathResult, err := jsonpath.Get("$[0].encoded_metadata_document", d.Get("saml_identity_provider_configuration"))
		if err == nil {
			if jsonPathResult != "" {
				request["EncodedMetadataDocument"] = jsonPathResult
			}
		}
	}

	if d.HasChange("saml_identity_provider_configuration.0.sso_status") {
		update = true
		jsonPathResult1, err := jsonpath.Get("$[0].sso_status", d.Get("saml_identity_provider_configuration"))
		if err == nil {
			request["SSOStatus"] = jsonPathResult1
		}
	}

	if d.HasChange("saml_identity_provider_configuration.0.entity_id") {
		update = true
		jsonPathResult2, err := jsonpath.Get("$[0].entity_id", d.Get("saml_identity_provider_configuration"))
		if err == nil {
			request["EntityId"] = jsonPathResult2
		}
	}

	if d.HasChange("saml_identity_provider_configuration.0.login_url") {
		update = true
		jsonPathResult3, err := jsonpath.Get("$[0].login_url", d.Get("saml_identity_provider_configuration"))
		if err == nil {
			request["LoginUrl"] = jsonPathResult3
		}
	}

	if d.HasChange("saml_identity_provider_configuration.0.want_request_signed") {
		update = true
		jsonPathResult4, err := jsonpath.Get("$[0].want_request_signed", d.Get("saml_identity_provider_configuration"))
		if err == nil {
			request["WantRequestSigned"] = jsonPathResult4
		}
	}

	if d.HasChange("saml_identity_provider_configuration.0.binding_type") {
		update = true
		jsonPathResult5, err := jsonpath.Get("$[0].binding_type", d.Get("saml_identity_provider_configuration"))
		if err == nil {
			request["BindingType"] = jsonPathResult5
		}
	}

	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("cloudsso", "2021-05-15", action, query, request, true)
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
	update = false
	action = "SetMFAAuthenticationStatus"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["DirectoryId"] = d.Id()

	if d.HasChange("mfa_authentication_status") {
		update = true
		request["MFAAuthenticationStatus"] = d.Get("mfa_authentication_status")
	}

	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("cloudsso", "2021-05-15", action, query, request, true)
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
	update = false
	action = "SetSCIMSynchronizationStatus"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["DirectoryId"] = d.Id()

	if d.HasChange("scim_synchronization_status") {
		update = true
		request["SCIMSynchronizationStatus"] = d.Get("scim_synchronization_status")
	}

	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("cloudsso", "2021-05-15", action, query, request, true)
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
	update = false
	action = "UpdateMFAAuthenticationSettings"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["DirectoryId"] = d.Id()

	if d.HasChange("mfa_authentication_setting_info.0.operation_for_risk_login") {
		update = true
		jsonPathResult, err := jsonpath.Get("$[0].operation_for_risk_login", d.Get("mfa_authentication_setting_info"))
		if err == nil {
			request["OperationForRiskLogin"] = jsonPathResult
		}
	}

	if d.HasChange("mfa_authentication_setting_info.0.mfa_authentication_advance_settings") {
		update = true
	}
	jsonPathResult1, err := jsonpath.Get("$[0].mfa_authentication_advance_settings", d.Get("mfa_authentication_setting_info"))
	if err == nil {
		request["MFAAuthenticationSettings"] = jsonPathResult1
	}

	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("cloudsso", "2021-05-15", action, query, request, true)
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
	update = false
	action = "UpdateUserProvisioningConfiguration"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["DirectoryId"] = d.Id()

	if d.HasChange("user_provisioning_configuration.0.session_duration") {
		update = true
		jsonPathResult, err := jsonpath.Get("$[0].session_duration", d.Get("user_provisioning_configuration"))
		if err == nil {
			request["NewSessionDuration"] = jsonPathResult
		}
	}

	if d.HasChange("user_provisioning_configuration.0.default_landing_page") {
		update = true
		jsonPathResult1, err := jsonpath.Get("$[0].default_landing_page", d.Get("user_provisioning_configuration"))
		if err == nil {
			request["NewDefaultLandingPage"] = jsonPathResult1
		}
	}

	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("cloudsso", "2021-05-15", action, query, request, true)
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
	update = false
	action = "UpdateDirectoryGlobalAccessStatus"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["DirectoryId"] = d.Id()

	if d.HasChange("directory_global_access_status") {
		update = true
	}
	request["GlobalAccessStatus"] = d.Get("directory_global_access_status")
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("cloudsso", "2021-05-15", action, query, request, true)
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
	update = false
	action = "SetPasswordPolicy"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["DirectoryId"] = d.Id()

	if d.HasChange("password_policy.0.min_password_length") {
		update = true
		jsonPathResult, err := jsonpath.Get("$[0].min_password_length", d.Get("password_policy"))
		if err == nil {
			request["MinPasswordLength"] = jsonPathResult
		}
	}

	if d.HasChange("password_policy.0.min_password_different_chars") {
		update = true
		jsonPathResult1, err := jsonpath.Get("$[0].min_password_different_chars", d.Get("password_policy"))
		if err == nil {
			request["MinPasswordDifferentChars"] = jsonPathResult1
		}
	}

	if d.HasChange("password_policy.0.password_not_contain_username") {
		update = true
		jsonPathResult2, err := jsonpath.Get("$[0].password_not_contain_username", d.Get("password_policy"))
		if err == nil {
			request["PasswordNotContainUsername"] = jsonPathResult2
		}
	}

	if d.HasChange("password_policy.0.max_password_age") {
		update = true
		jsonPathResult3, err := jsonpath.Get("$[0].max_password_age", d.Get("password_policy"))
		if err == nil {
			request["MaxPasswordAge"] = jsonPathResult3
		}
	}

	if d.HasChange("password_policy.0.password_reuse_prevention") {
		update = true
		jsonPathResult4, err := jsonpath.Get("$[0].password_reuse_prevention", d.Get("password_policy"))
		if err == nil {
			request["PasswordReusePrevention"] = jsonPathResult4
		}
	}

	if d.HasChange("password_policy.0.max_login_attempts") {
		update = true
		jsonPathResult5, err := jsonpath.Get("$[0].max_login_attempts", d.Get("password_policy"))
		if err == nil {
			request["MaxLoginAttempts"] = jsonPathResult5
		}
	}

	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("cloudsso", "2021-05-15", action, query, request, true)
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
	update = false
	action = "SetLoginPreference"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["DirectoryId"] = d.Id()

	if d.HasChange("login_preference.0.login_network_masks") {
		update = true
		jsonPathResult, err := jsonpath.Get("$[0].login_network_masks", d.Get("login_preference"))
		if err == nil {
			request["LoginNetworkMasks"] = jsonPathResult
		}
	}

	if d.HasChange("login_preference.0.allow_user_to_get_credentials") {
		update = true
	}
	jsonPathResult1, err = jsonpath.Get("$[0].allow_user_to_get_credentials", d.Get("login_preference"))
	if err == nil {
		request["AllowUserToGetCredentials"] = jsonPathResult1
	}

	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("cloudsso", "2021-05-15", action, query, request, true)
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
	update = false
	action = "SetDirectorySAMLServiceProviderInfo"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["DirectoryId"] = d.Id()

	if d.HasChange("saml_service_provider.0.authn_sign_algo") {
		update = true
		jsonPathResult, err := jsonpath.Get("$[0].authn_sign_algo", d.Get("saml_service_provider"))
		if err == nil {
			request["AuthnSignAlgo"] = jsonPathResult
		}
	}

	if d.HasChange("saml_service_provider.0.certificate_type") {
		update = true
		jsonPathResult1, err := jsonpath.Get("$[0].certificate_type", d.Get("saml_service_provider"))
		if err == nil {
			request["CertificateType"] = jsonPathResult1
		}
	}

	if d.HasChange("saml_service_provider.0.support_encrypted_assertion") {
		update = true
		jsonPathResult2, err := jsonpath.Get("$[0].support_encrypted_assertion", d.Get("saml_service_provider"))
		if err == nil {
			request["SupportEncryptedAssertion"] = jsonPathResult2
		}
	}

	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("cloudsso", "2021-05-15", action, query, request, true)
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

	d.Partial(false)
	return resourceAliCloudCloudSSODirectoryRead(d, meta)
}

func resourceAliCloudCloudSSODirectoryDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	action := "DeleteDirectory"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["DirectoryId"] = d.Id()

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("cloudsso", "2021-05-15", action, query, request, true)

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
		if IsExpectedErrors(err, []string{"EntityNotExists.Directory"}) || NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return nil
}

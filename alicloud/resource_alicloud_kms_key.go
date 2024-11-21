package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudKmsKey() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudKmsKeyCreate,
		Read:   resourceAliCloudKmsKeyRead,
		Update: resourceAliCloudKmsKeyUpdate,
		Delete: resourceAliCloudKmsKeyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"key_usage": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Computed:     true,
				ValidateFunc: StringInSlice([]string{"ENCRYPT/DECRYPT", "SIGN/VERIFY"}, false),
			},
			"origin": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Computed:     true,
				ValidateFunc: StringInSlice([]string{"Aliyun_KMS", "EXTERNAL"}, false),
			},
			"key_spec": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Computed:     true,
				ValidateFunc: StringInSlice([]string{"Aliyun_AES_256", "Aliyun_AES_128", "Aliyun_AES_192", "Aliyun_SM4", "RSA_2048", "RSA_3072", "EC_P256", "EC_P256K", "EC_SM2"}, false),
			},
			"dkms_instance_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"protection_level": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      "SOFTWARE",
				ValidateFunc: StringInSlice([]string{"SOFTWARE", "HSM"}, false),
			},
			"automatic_rotation": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: StringInSlice([]string{"Enabled", "Disabled"}, false),
			},
			"rotation_interval": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"policy": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					equal, _ := compareJsonTemplateAreEquivalent(old, new)
					return equal
				},
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"status": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ValidateFunc:  StringInSlice([]string{"Enabled", "Disabled", "PendingDeletion"}, false),
				ConflictsWith: []string{"key_state"},
			},
			"pending_window_in_days": {
				Type:          schema.TypeInt,
				Optional:      true,
				ValidateFunc:  IntBetween(7, 366),
				ConflictsWith: []string{"deletion_window_in_days"},
			},
			"tags": tagsSchema(),
			"arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"primary_key_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"last_rotation_date": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"next_rotation_date": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"material_expire_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"creator": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"creation_date": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"delete_date": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"deletion_window_in_days": {
				Type:          schema.TypeInt,
				Optional:      true,
				ValidateFunc:  IntBetween(7, 366),
				ConflictsWith: []string{"pending_window_in_days"},
				Deprecated:    "Field `deletion_window_in_days` has been deprecated from provider version 1.85.0. New field `pending_window_in_days` instead.",
			},
			"key_state": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ValidateFunc:  StringInSlice([]string{"Enabled", "Disabled", "PendingDeletion"}, false),
				ConflictsWith: []string{"status"},
				Deprecated:    "Field `key_state` has been deprecated from provider version 1.123.1. New field `status` instead.",
			},
			"is_enabled": {
				Type:       schema.TypeBool,
				Optional:   true,
				Computed:   true,
				Deprecated: "Field `is_enabled` has been deprecated from provider version 1.85.0. New field `key_state` instead.",
			},
		},
	}
}

func resourceAliCloudKmsKeyCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateKey"
	request := make(map[string]interface{})
	conn, err := client.NewKmsClient()
	if err != nil {
		return WrapError(err)
	}

	if v, ok := d.GetOk("key_usage"); ok {
		request["KeyUsage"] = v
	}

	if v, ok := d.GetOk("origin"); ok {
		request["Origin"] = v
	}

	if v, ok := d.GetOk("key_spec"); ok {
		request["KeySpec"] = v
	}

	if v, ok := d.GetOk("dkms_instance_id"); ok {
		request["DKMSInstanceId"] = v
	}

	if v, ok := d.GetOk("protection_level"); ok {
		request["ProtectionLevel"] = v
	}

	if v, ok := d.GetOk("automatic_rotation"); ok {
		request["EnableAutomaticRotation"] = convertKmsKeyAutomaticRotationRequest(v.(string))
	}

	if v, ok := d.GetOk("rotation_interval"); ok {
		request["RotationInterval"] = v
	}

	if v, ok := d.GetOk("policy"); ok {
		request["Policy"] = v
	}

	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}

	if v, ok := d.GetOk("tags"); ok {
		tagsMaps := ConvertTagsForKms(v.(map[string]interface{}))
		tagsJson, err := convertArrayObjectToJsonString(tagsMaps)
		if err != nil {
			return WrapError(err)
		}

		request["Tags"] = tagsJson
	}

	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-01-20"), StringPointer("AK"), nil, request, &runtime)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_kms_key", action, AlibabaCloudSdkGoERROR)
	}

	if resp, err := jsonpath.Get("$.KeyMetadata", response); err != nil || resp == nil {
		return WrapErrorf(err, IdMsg, "alicloud_kms_key")
	} else {
		keyId := resp.(map[string]interface{})["KeyId"]
		d.SetId(fmt.Sprint(keyId))
	}

	return resourceAliCloudKmsKeyUpdate(d, meta)
}

func resourceAliCloudKmsKeyRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	kmsService := KmsService{client}

	object, err := kmsService.DescribeKmsKey(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_kms_key kmsService.DescribeKmsKey Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("key_usage", object["KeyUsage"])
	d.Set("origin", object["Origin"])
	d.Set("key_spec", object["KeySpec"])
	d.Set("protection_level", object["ProtectionLevel"])
	d.Set("automatic_rotation", object["AutomaticRotation"])
	d.Set("rotation_interval", object["RotationInterval"])
	d.Set("description", object["Description"])
	d.Set("status", object["KeyState"])
	d.Set("arn", object["Arn"])
	d.Set("primary_key_version", object["PrimaryKeyVersion"])
	d.Set("last_rotation_date", object["LastRotationDate"])
	d.Set("next_rotation_date", object["NextRotationDate"])
	d.Set("material_expire_time", object["MaterialExpireTime"])
	d.Set("creator", object["Creator"])
	d.Set("creation_date", object["CreationDate"])
	d.Set("delete_date", object["DeleteDate"])
	d.Set("key_state", object["KeyState"])
	d.Set("is_enabled", convertKmsKeyIsEnabledResponse(object["KeyState"]))

	if dkmsInstanceId, ok := object["DKMSInstanceId"].(string); ok {
		d.Set("dkms_instance_id", dkmsInstanceId)

		if dkmsInstanceId != "" {
			policy, err := kmsService.DescribeKmsKeyPolicy(d.Id())
			if err != nil {
				return WrapError(err)
			}

			if keyPolicy, ok := policy["Policy"]; ok {
				d.Set("policy", keyPolicy)
			}
		}
	}

	listTagResourcesObject, err := kmsService.ListTagResources(d.Id(), "key")
	if err != nil {
		return WrapError(err)
	}

	d.Set("tags", tagsToMap(listTagResourcesObject))

	return nil
}

func resourceAliCloudKmsKeyUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	kmsService := KmsService{client}
	var response map[string]interface{}
	d.Partial(true)

	update := false
	updateRotationPolicyReq := map[string]interface{}{
		"KeyId": d.Id(),
	}

	if !d.IsNewResource() && d.HasChange("automatic_rotation") {
		update = true
	}
	if v, ok := d.GetOk("automatic_rotation"); ok {
		updateRotationPolicyReq["EnableAutomaticRotation"] = convertKmsKeyAutomaticRotationRequest(v)
	}

	if !d.IsNewResource() && d.HasChange("rotation_interval") {
		update = true
	}
	if v, ok := d.GetOk("rotation_interval"); ok {
		updateRotationPolicyReq["RotationInterval"] = v
	}

	if update {
		action := "UpdateRotationPolicy"
		conn, err := client.NewKmsClient()
		if err != nil {
			return WrapError(err)
		}

		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-01-20"), StringPointer("AK"), nil, updateRotationPolicyReq, &runtime)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, updateRotationPolicyReq)

		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}

		d.SetPartial("automatic_rotation")
		d.SetPartial("rotation_interval")
	}

	update = false
	setKeyPolicyReq := map[string]interface{}{
		"KeyId": d.Id(),
	}

	if !d.IsNewResource() && d.HasChange("policy") {
		update = true
	}
	if v, ok := d.GetOk("policy"); ok {
		setKeyPolicyReq["Policy"] = v
	}

	if update {
		action := "SetKeyPolicy"
		conn, err := client.NewKmsClient()
		if err != nil {
			return WrapError(err)
		}

		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-01-20"), StringPointer("AK"), nil, setKeyPolicyReq, &runtime)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, setKeyPolicyReq)

		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}

		d.SetPartial("policy")
	}

	update = false
	updateKeyDescriptionReq := map[string]interface{}{
		"KeyId": d.Id(),
	}

	if !d.IsNewResource() && d.HasChange("description") {
		update = true
	}
	if v, ok := d.GetOk("description"); ok {
		updateKeyDescriptionReq["Description"] = v
	}

	if update {
		action := "UpdateKeyDescription"
		conn, err := client.NewKmsClient()
		if err != nil {
			return WrapError(err)
		}

		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-01-20"), StringPointer("AK"), nil, updateKeyDescriptionReq, &runtime)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, updateKeyDescriptionReq)

		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}

		d.SetPartial("description")
	}

	update = false
	var statusTarget string

	if d.HasChange("status") {
		update = true

		if v, ok := d.GetOk("status"); ok {
			statusTarget = v.(string)
		}
	}

	if d.HasChange("key_state") {
		update = true

		if v, ok := d.GetOk("key_state"); ok {
			statusTarget = v.(string)
		}
	}

	if d.HasChange("is_enabled") {
		update = true

		if v, ok := d.GetOkExists("is_enabled"); ok {
			statusTarget = convertKmsKeyIsEnabledRequest(v)
		}
	}

	if update {
		object, err := kmsService.DescribeKmsKey(d.Id())
		if err != nil {
			return WrapError(err)
		}

		if object["KeyState"].(string) != statusTarget {
			if statusTarget == "Disabled" {
				request := map[string]interface{}{
					"KeyId": d.Id(),
				}

				action := "DisableKey"
				conn, err := client.NewKmsClient()
				if err != nil {
					return WrapError(err)
				}

				runtime := util.RuntimeOptions{}
				runtime.SetAutoretry(true)
				wait := incrementalWait(3*time.Second, 3*time.Second)
				err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
					response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-01-20"), StringPointer("AK"), nil, request, &runtime)
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

				stateConf := BuildStateConf([]string{}, []string{"Disabled"}, d.Timeout(schema.TimeoutUpdate), 10*time.Second, kmsService.KmsKeyStateRefreshFunc(d.Id(), []string{}))
				if _, err := stateConf.WaitForState(); err != nil {
					return WrapErrorf(err, IdMsg, d.Id())
				}
			}

			if statusTarget == "Enabled" {
				request := map[string]interface{}{
					"KeyId": d.Id(),
				}

				action := "EnableKey"
				conn, err := client.NewKmsClient()
				if err != nil {
					return WrapError(err)
				}

				runtime := util.RuntimeOptions{}
				runtime.SetAutoretry(true)
				wait := incrementalWait(3*time.Second, 3*time.Second)
				err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
					response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-01-20"), StringPointer("AK"), nil, request, &runtime)
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

				stateConf := BuildStateConf([]string{}, []string{"Enabled"}, d.Timeout(schema.TimeoutUpdate), 10*time.Second, kmsService.KmsKeyStateRefreshFunc(d.Id(), []string{}))
				if _, err := stateConf.WaitForState(); err != nil {
					return WrapErrorf(err, IdMsg, d.Id())
				}
			}

			d.SetPartial("status")
			d.SetPartial("key_state")
			d.SetPartial("is_enabled")
		}
	}

	if !d.IsNewResource() && d.HasChange("tags") {
		if err := kmsService.SetResourceTags(d, "key"); err != nil {
			return WrapError(err)
		}

		d.SetPartial("tags")
	}

	d.Partial(false)

	return resourceAliCloudKmsKeyRead(d, meta)
}

func resourceAliCloudKmsKeyDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "ScheduleKeyDeletion"
	var response map[string]interface{}
	conn, err := client.NewKmsClient()
	if err != nil {
		return WrapError(err)
	}

	request := map[string]interface{}{
		"KeyId": d.Id(),
	}

	if v, ok := d.GetOk("pending_window_in_days"); ok {
		request["PendingWindowInDays"] = v
	} else if v, ok := d.GetOk("deletion_window_in_days"); ok {
		request["PendingWindowInDays"] = v
	} else {
		return WrapError(Error(`[ERROR] Argument "pending_window_in_days" or "deletion_window_in_days" must be set one!`))
	}

	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-01-20"), StringPointer("AK"), nil, request, &runtime)
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

	return nil
}

func convertKmsKeyAutomaticRotationRequest(source interface{}) interface{} {
	switch source {
	case "Enabled":
		return true
	case "Disabled":
		return false
	}

	return false
}

func convertKmsKeyIsEnabledRequest(source interface{}) string {
	switch source {
	case true:
		return "Enabled"
	case false:
		return "Disabled"
	}

	return ""
}

func convertKmsKeyIsEnabledResponse(source interface{}) interface{} {
	switch source {
	case "Enabled":
		return true
	case "Disabled":
		return false
	}

	return false
}

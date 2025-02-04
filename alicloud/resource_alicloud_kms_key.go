// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/PaesslerAG/jsonpath"
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
			"automatic_rotation": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: StringInSlice([]string{"Enabled", "Disabled"}, false),
			},
			"deletion_protection": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: StringInSlice([]string{"Enabled", "Disabled"}, false),
			},
			"deletion_protection_description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"dkms_instance_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"key_spec": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"Aliyun_AES_256", "Aliyun_AES_128", "Aliyun_AES_192", "Aliyun_SM4", "RSA_2048", "RSA_3072", "EC_P256", "EC_SM2", "EC_P256K"}, false),
			},
			"key_usage": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"ENCRYPT/DECRYPT", "SIGN/VERIFY"}, false),
			},
			"origin": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"Aliyun_KMS", "EXTERNAL"}, false),
			},
			"pending_window_in_days": {
				Type:          schema.TypeInt,
				Optional:      true,
				ValidateFunc:  IntBetween(7, 366),
				ConflictsWith: []string{"deletion_window_in_days"},
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
			"protection_level": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      "SOFTWARE",
				ValidateFunc: StringInSlice([]string{"SOFTWARE", "HSM"}, false),
			},
			"rotation_interval": {
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

	action := "CreateKey"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})

	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}
	if v, ok := d.GetOk("key_usage"); ok {
		request["KeyUsage"] = v
	}
	if v, ok := d.GetOk("origin"); ok {
		request["Origin"] = v
	}
	if v, ok := d.GetOk("protection_level"); ok {
		request["ProtectionLevel"] = v
	}
	if v, ok := d.GetOk("rotation_interval"); ok {
		request["RotationInterval"] = v
	}
	if v, ok := d.GetOk("key_spec"); ok {
		request["KeySpec"] = v
	}
	if v, ok := d.GetOk("policy"); ok {
		request["Policy"] = v
	}
	if v, ok := d.GetOk("automatic_rotation"); ok {
		request["EnableAutomaticRotation"] = convertKmsKeyAutomaticRotationRequest(v.(string))
	}
	if v, ok := d.GetOk("dkms_instance_id"); ok {
		request["DKMSInstanceId"] = v
	}

	if v, ok := d.GetOk("tags"); ok {
		tagsMaps := ConvertTagsForKms(v.(map[string]interface{}))
		tagsJson, err := convertArrayObjectToJsonString(tagsMaps)
		if err != nil {
			return WrapError(err)
		}

		request["Tags"] = tagsJson
	}
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("Kms", "2016-01-20", action, query, request, true)
		if err != nil {
			if IsExpectedErrors(err, []string{"Forbidden.DKMSInstanceStateInvalid"}) || NeedRetry(err) {
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
	kmsServiceV2 := KmsServiceV2{client}

	objectRaw, err := kmsServiceV2.DescribeKmsKey(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_kms_key DescribeKmsKey Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	if objectRaw["AutomaticRotation"] != nil {
		d.Set("automatic_rotation", objectRaw["AutomaticRotation"])
	}
	if objectRaw["DeletionProtection"] != nil {
		d.Set("deletion_protection", objectRaw["DeletionProtection"])
	}
	d.Set("deletion_protection_description", objectRaw["DeletionProtectionDescription"])
	if objectRaw["Description"] != nil {
		d.Set("description", objectRaw["Description"])
	}

	if dkmsInstanceId, ok := objectRaw["DKMSInstanceId"].(string); ok {
		d.Set("dkms_instance_id", dkmsInstanceId)

		if dkmsInstanceId != "" {
			policy, err := kmsServiceV2.DescribeKmsKeyPolicy(d.Id())
			if err != nil {
				return WrapError(err)
			}

			if keyPolicy, ok := policy["Policy"]; ok {
				d.Set("policy", keyPolicy)
			}
		}
	}

	if objectRaw["KeySpec"] != nil {
		d.Set("key_spec", objectRaw["KeySpec"])
	}
	if objectRaw["KeyUsage"] != nil {
		d.Set("key_usage", objectRaw["KeyUsage"])
	}
	if objectRaw["Origin"] != nil {
		d.Set("origin", objectRaw["Origin"])
	}
	if objectRaw["ProtectionLevel"] != nil {
		d.Set("protection_level", objectRaw["ProtectionLevel"])
	}
	d.Set("rotation_interval", objectRaw["RotationInterval"])
	if objectRaw["KeyState"] != nil {
		d.Set("status", objectRaw["KeyState"])
		d.Set("key_state", objectRaw["KeyState"])
		d.Set("is_enabled", convertKmsKeyIsEnabledResponse(objectRaw["KeyState"]))
	}
	if objectRaw["Arn"] != nil {
		d.Set("arn", objectRaw["Arn"])
	}
	if objectRaw["PrimaryKeyVersion"] != nil {
		d.Set("primary_key_version", objectRaw["PrimaryKeyVersion"])
	}
	if objectRaw["LastRotationDate"] != nil {
		d.Set("last_rotation_date", objectRaw["LastRotationDate"])
	}
	d.Set("next_rotation_date", objectRaw["NextRotationDate"])
	if objectRaw["MaterialExpireTime"] != nil {
		d.Set("material_expire_time", objectRaw["MaterialExpireTime"])
	}
	if objectRaw["Creator"] != nil {
		d.Set("creator", objectRaw["Creator"])
	}
	if objectRaw["CreationDate"] != nil {
		d.Set("creation_date", objectRaw["CreationDate"])
	}
	if objectRaw["DeleteDate"] != nil {
		d.Set("delete_date", objectRaw["DeleteDate"])
	}

	listTagResourcesObject, err := kmsServiceV2.ListTagResources(d.Id(), "key")
	if err != nil {
		return WrapError(err)
	}

	d.Set("tags", tagsToMap(listTagResourcesObject))

	return nil
}

func resourceAliCloudKmsKeyUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	kmsServiceV2 := KmsServiceV2{client}

	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false
	d.Partial(true)

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
		object, err := kmsServiceV2.DescribeKmsKey(d.Id())
		if err != nil {
			return WrapError(err)
		}

		if object["KeyState"].(string) != statusTarget {
			if statusTarget == "Disabled" {
				action := "DisableKey"
				request = make(map[string]interface{})
				query = make(map[string]interface{})
				request["KeyId"] = d.Id()

				wait := incrementalWait(3*time.Second, 3*time.Second)
				err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
					response, err = client.RpcPost("Kms", "2016-01-20", action, query, request, true)
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

				stateConf := BuildStateConf([]string{}, []string{"Disabled"}, d.Timeout(schema.TimeoutUpdate), 10*time.Second, kmsServiceV2.KmsKeyStateRefreshFunc(d.Id(), "KeyState", []string{}))
				if _, err := stateConf.WaitForState(); err != nil {
					return WrapErrorf(err, IdMsg, d.Id())
				}
			}

			if statusTarget == "Enabled" {
				action := "EnableKey"
				request = make(map[string]interface{})
				query = make(map[string]interface{})
				request["KeyId"] = d.Id()

				wait := incrementalWait(3*time.Second, 3*time.Second)
				err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
					response, err = client.RpcPost("Kms", "2016-01-20", action, query, request, true)
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

				stateConf := BuildStateConf([]string{}, []string{"Enabled"}, d.Timeout(schema.TimeoutUpdate), 10*time.Second, kmsServiceV2.KmsKeyStateRefreshFunc(d.Id(), "KeyState", []string{}))
				if _, err := stateConf.WaitForState(); err != nil {
					return WrapErrorf(err, IdMsg, d.Id())
				}
			}
		}
	}
	update = false
	action := "UpdateKeyDescription"
	var err error
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["KeyId"] = d.Id()

	if !d.IsNewResource() && d.HasChange("description") {
		update = true
	}
	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("Kms", "2016-01-20", action, query, request, true)
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
	action = "UpdateRotationPolicy"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["KeyId"] = d.Id()

	if !d.IsNewResource() && d.HasChange("automatic_rotation") {
		update = true
	}
	if v, ok := d.GetOk("automatic_rotation"); ok {
		request["EnableAutomaticRotation"] = convertKmsKeyAutomaticRotationRequest(v)
	}

	if !d.IsNewResource() && d.HasChange("rotation_interval") {
		update = true
	}
	if v, ok := d.GetOk("rotation_interval"); ok {
		request["RotationInterval"] = v
	}

	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("Kms", "2016-01-20", action, query, request, true)
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
	action = "SetDeletionProtection"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["KeyId"] = d.Id()

	if d.HasChange("deletion_protection_description") {
		update = true
	}
	if v, ok := d.GetOk("deletion_protection_description"); ok {
		request["DeletionProtectionDescription"] = v
	}

	if d.HasChange("deletion_protection") {
		update = true
	}
	if v, ok := d.GetOk("deletion_protection"); ok {
		request["EnableDeletionProtection"] = convertKmsKeyDeletionProtectionRequest(v)
	}

	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("Kms", "2016-01-20", action, query, request, true)
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
	action = "SetKeyPolicy"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["KeyId"] = d.Id()

	if !d.IsNewResource() && d.HasChange("policy") {
		update = true
	}
	if v, ok := d.GetOk("policy"); ok {
		request["Policy"] = v
	}
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("Kms", "2016-01-20", action, query, request, true)
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

	if !d.IsNewResource() && d.HasChange("tags") {
		if err := kmsServiceV2.SetResourceTags(d, "key"); err != nil {
			return WrapError(err)
		}
	}

	d.Partial(false)
	return resourceAliCloudKmsKeyRead(d, meta)
}

func resourceAliCloudKmsKeyDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	action := "ScheduleKeyDeletion"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["KeyId"] = d.Id()

	if v, ok := d.GetOk("pending_window_in_days"); ok {
		request["PendingWindowInDays"] = v
	} else if v, ok := d.GetOk("deletion_window_in_days"); ok {
		request["PendingWindowInDays"] = v
	} else {
		return WrapError(Error(`[ERROR] Argument "pending_window_in_days" or "deletion_window_in_days" must be set one!`))
	}

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("Kms", "2016-01-20", action, query, request, true)

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
		if NotFoundError(err) {
			return nil
		}
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

func convertKmsKeyDeletionProtectionRequest(source interface{}) interface{} {
	source = fmt.Sprint(source)
	switch source {
	case "Enabled":
		return true
	case "Disabled":
		return false
	}

	return false
}

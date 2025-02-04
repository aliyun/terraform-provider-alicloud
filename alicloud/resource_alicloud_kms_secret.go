package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudKmsSecret() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudKmsSecretCreate,
		Read:   resourceAliCloudKmsSecretRead,
		Update: resourceAliCloudKmsSecretUpdate,
		Delete: resourceAliCloudKmsSecretDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"secret_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"secret_data": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if v, ok := d.GetOk("secret_type"); ok && v.(string) != "Generic" {
						return d.Id() != ""
					}
					return false
				},
			},
			"version_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"secret_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Computed:     true,
				ValidateFunc: StringInSlice([]string{"Generic", "Rds", "RAMCredentials", "ECS"}, false),
			},
			"secret_data_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "text",
				ValidateFunc: StringInSlice([]string{"text", "binary"}, false),
			},
			"encryption_key_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"dkms_instance_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"extended_config": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if v, ok := d.GetOk("secret_type"); ok && v.(string) != "Generic" {
						return d.Id() != ""
					}
					return false
				},
			},
			"enable_automatic_rotation": {
				Type:     schema.TypeBool,
				Optional: true,
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
			"force_delete_without_recovery": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"recovery_window_in_days": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  30,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return d.Get("force_delete_without_recovery").(bool)
				},
			},
			"version_stages": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"tags": tagsSchema(),
			"arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"planned_delete_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAliCloudKmsSecretCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateSecret"
	request := make(map[string]interface{})
	var err error

	request["SecretName"] = d.Get("secret_name")
	request["SecretData"] = d.Get("secret_data")
	request["VersionId"] = d.Get("version_id")

	if v, ok := d.GetOk("secret_type"); ok {
		request["SecretType"] = v
	}

	if v, ok := d.GetOk("secret_data_type"); ok {
		request["SecretDataType"] = v
	}

	if v, ok := d.GetOk("encryption_key_id"); ok {
		request["EncryptionKeyId"] = v
	}

	if v, ok := d.GetOk("dkms_instance_id"); ok {
		request["DKMSInstanceId"] = v
	}

	if v, ok := d.GetOk("extended_config"); ok {
		request["ExtendedConfig"] = v
	}

	if v, ok := d.GetOkExists("enable_automatic_rotation"); ok {
		request["EnableAutomaticRotation"] = v
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

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *resource.RetryError {
		response, err = client.RpcPost("Kms", "2016-01-20", action, nil, request, true)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_kms_secret", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(request["SecretName"]))

	return resourceAliCloudKmsSecretUpdate(d, meta)
}

func resourceAliCloudKmsSecretRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	kmsService := KmsService{client}

	object, err := kmsService.DescribeKmsSecret(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_kms_secret kmsService.DescribeKmsSecret Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("secret_name", object["SecretName"])
	d.Set("secret_type", object["SecretType"])
	d.Set("encryption_key_id", object["EncryptionKeyId"])
	d.Set("extended_config", object["ExtendedConfig"])
	d.Set("enable_automatic_rotation", convertKmsSecretEnableAutomaticRotationResponse(object["AutomaticRotation"]))
	d.Set("rotation_interval", object["RotationInterval"])
	d.Set("description", object["Description"])
	d.Set("arn", object["Arn"])
	d.Set("create_time", object["CreateTime"])
	d.Set("planned_delete_time", object["PlannedDeleteTime"])

	if dkmsInstanceId, ok := object["DKMSInstanceId"].(string); ok {
		d.Set("dkms_instance_id", dkmsInstanceId)

		if dkmsInstanceId != "" {
			policy, err := kmsService.DescribeKmsSecretPolicy(d.Id())
			if err != nil {
				return WrapError(err)
			}

			if secretPolicy, ok := policy["Policy"]; ok {
				d.Set("policy", secretPolicy)
			}
		}
	}

	secretValue, err := kmsService.GetSecretValue(d.Id())
	if err != nil {
		return WrapError(err)
	}

	d.Set("secret_data", secretValue["SecretData"])
	d.Set("version_id", secretValue["VersionId"])
	d.Set("secret_data_type", secretValue["SecretDataType"])

	if versionStages, ok := secretValue["VersionStages"]; ok {
		if versionStageList, ok := versionStages.(map[string]interface{})["VersionStage"]; ok {
			d.Set("version_stages", versionStageList)
		}
	}

	listTagResourcesObject, err := kmsService.ListTagResources(d.Id(), "secret")
	if err != nil {
		return WrapError(err)
	}

	d.Set("tags", tagsToMap(listTagResourcesObject))

	return nil
}

func resourceAliCloudKmsSecretUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	kmsService := KmsService{client}
	var response map[string]interface{}
	var err error
	d.Partial(true)

	update := false
	putSecretValueReq := map[string]interface{}{
		"SecretName": d.Id(),
	}

	if !d.IsNewResource() && d.HasChange("secret_data") {
		update = true
	}
	putSecretValueReq["SecretData"] = d.Get("secret_data")

	if !d.IsNewResource() && d.HasChange("version_id") {
		update = true
	}
	putSecretValueReq["VersionId"] = d.Get("version_id")

	if !d.IsNewResource() && d.HasChange("secret_data_type") {
		update = true
	}
	if v, ok := d.GetOk("secret_data_type"); ok {
		putSecretValueReq["SecretDataType"] = v
	}

	if d.HasChange("version_stages") {
		update = true
	}
	if v, ok := d.GetOk("version_stages"); ok {
		versionStagesJson, err := convertArrayObjectToJsonString(v.(*schema.Set).List())
		if err != nil {
			return WrapError(err)
		}

		putSecretValueReq["VersionStages"] = versionStagesJson
	}

	if update {
		action := "PutSecretValue"
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			response, err = client.RpcPost("Kms", "2016-01-20", action, nil, putSecretValueReq, false)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, putSecretValueReq)

		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}

		d.SetPartial("secret_data")
		d.SetPartial("version_id")
		d.SetPartial("secret_data_type")
		d.SetPartial("version_stages")
	}

	update = false
	updateSecretRotationPolicyReq := map[string]interface{}{
		"SecretName": d.Id(),
	}

	if !d.IsNewResource() && d.HasChange("enable_automatic_rotation") {
		update = true
	}
	if v, ok := d.GetOkExists("enable_automatic_rotation"); ok {
		updateSecretRotationPolicyReq["EnableAutomaticRotation"] = v
	}

	if !d.IsNewResource() && d.HasChange("rotation_interval") {
		update = true
	}
	if v, ok := d.GetOk("rotation_interval"); ok {
		updateSecretRotationPolicyReq["RotationInterval"] = v
	}

	if update {
		action := "UpdateSecretRotationPolicy"
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			response, err = client.RpcPost("Kms", "2016-01-20", action, nil, updateSecretRotationPolicyReq, false)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, updateSecretRotationPolicyReq)

		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}

		d.SetPartial("enable_automatic_rotation")
		d.SetPartial("rotation_interval")
	}

	update = false
	setSecretPolicyReq := map[string]interface{}{
		"SecretName": d.Id(),
	}

	if !d.IsNewResource() && d.HasChange("policy") {
		update = true
	}
	if v, ok := d.GetOk("policy"); ok {
		setSecretPolicyReq["Policy"] = v
	}

	if update {
		action := "SetSecretPolicy"
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			response, err = client.RpcPost("Kms", "2016-01-20", action, nil, setSecretPolicyReq, false)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, setSecretPolicyReq)

		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}

		d.SetPartial("policy")
	}

	update = false
	updateSecretReq := map[string]interface{}{
		"SecretName": d.Id(),
	}

	if !d.IsNewResource() && d.HasChange("description") {
		update = true
	}
	if v, ok := d.GetOk("description"); ok {
		updateSecretReq["Description"] = v
	}

	if update {
		action := "UpdateSecret"
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			response, err = client.RpcPost("Kms", "2016-01-20", action, nil, updateSecretReq, false)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, updateSecretReq)

		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}

		d.SetPartial("description")
	}

	if !d.IsNewResource() && d.HasChange("tags") {
		if err := kmsService.SetResourceTags(d, "secret"); err != nil {
			return WrapError(err)
		}

		d.SetPartial("tags")
	}

	d.Partial(false)

	return resourceAliCloudKmsSecretRead(d, meta)
}

func resourceAliCloudKmsSecretDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DeleteSecret"
	var response map[string]interface{}
	var err error

	request := map[string]interface{}{
		"SecretName": d.Id(),
	}

	if v, ok := d.GetOkExists("force_delete_without_recovery"); ok {
		request["ForceDeleteWithoutRecovery"] = v
	}

	if v, ok := d.GetOkExists("recovery_window_in_days"); ok {
		request["RecoveryWindowInDays"] = v
	}

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *resource.RetryError {
		response, err = client.RpcPost("Kms", "2016-01-20", action, nil, request, true)
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
		if IsExpectedErrors(err, []string{"Forbidden.ResourceNotFound"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return nil
}

func convertKmsSecretEnableAutomaticRotationResponse(source interface{}) interface{} {
	switch source {
	case "Enabled":
		return true
	case "Disabled":
		return false
	}

	return false
}

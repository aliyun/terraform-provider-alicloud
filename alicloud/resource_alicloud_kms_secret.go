// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"log"
	"time"
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
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"certificate_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"custom_data": {
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
				ForceNew: true,
			},
			"enable_automatic_rotation": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"encryption_key_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"extended_config": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"extended_config_custom_data": {
				Type:     schema.TypeMap,
				Optional: true,
			},
			"key_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"move_to_version": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"policy": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"policy_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"recovery_window_in_days": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"remove_from_version": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"rotation_interval": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"secret_data": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"secret_data_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"secret_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"secret_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"tags": tagsSchema(),
			"version_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAliCloudKmsSecretCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateSecret"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	conn, err := client.NewKmsClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	query["SecretName"] = d.Get("secret_name")

	if v, ok := d.GetOk("secret_data_type"); ok {
		request["SecretDataType"] = v
	}
	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}
	request["VersionId"] = d.Get("version_id")
	if v, ok := d.GetOk("encryption_key_id"); ok {
		request["EncryptionKeyId"] = v
	}
	request["SecretData"] = d.Get("secret_data")
	if v, ok := d.GetOk("secret_type"); ok {
		request["SecretType"] = v
	}
	if v, ok := d.GetOkExists("enable_automatic_rotation"); ok {
		request["EnableAutomaticRotation"] = v
	}
	if v, ok := d.GetOk("rotation_interval"); ok {
		request["RotationInterval"] = v
	}
	if v, ok := d.GetOk("extended_config"); ok {
		request["ExtendedConfig"] = v
	}
	if v, ok := d.GetOk("dkms_instance_id"); ok {
		request["DKMSInstanceId"] = v
	}
	if v, ok := d.GetOk("policy"); ok {
		request["Policy"] = v
	}
	if v, ok := d.GetOk("tags"); ok {
		jsonPathResult11, err := jsonpath.Get("$", v)
		if err == nil && jsonPathResult11 != "" {
			request["Tags"] = jsonPathResult11
		}
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-01-20"), StringPointer("AK"), query, request, &runtime)

		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_kms_secret", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["SecretName"]))

	return resourceAliCloudKmsSecretRead(d, meta)
}

func resourceAliCloudKmsSecretRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	kmsServiceV2 := KmsServiceV2{client}

	objectRaw, err := kmsServiceV2.DescribeKmsSecret(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_kms_secret DescribeKmsSecret Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("create_time", objectRaw["CreateTime"])
	d.Set("description", objectRaw["Description"])
	d.Set("dkms_instance_id", objectRaw["DKMSInstanceId"])
	d.Set("encryption_key_id", objectRaw["EncryptionKeyId"])
	d.Set("extended_config", objectRaw["ExtendedConfig"])
	d.Set("rotation_interval", objectRaw["RotationInterval"])
	d.Set("secret_type", objectRaw["SecretType"])
	d.Set("secret_name", objectRaw["SecretName"])

	tagsMaps, _ := jsonpath.Get("$.Tags.Tag", objectRaw)
	d.Set("tags", tagsToMap(tagsMaps))

	objectRaw, err = kmsServiceV2.DescribeGetSecretValue(d.Id())
	if err != nil {
		return WrapError(err)
	}

	d.Set("create_time", objectRaw["CreateTime"])
	d.Set("extended_config", objectRaw["ExtendedConfig"])
	d.Set("rotation_interval", objectRaw["RotationInterval"])
	d.Set("secret_data", objectRaw["SecretData"])
	d.Set("secret_data_type", objectRaw["SecretDataType"])
	d.Set("secret_type", objectRaw["SecretType"])
	d.Set("version_id", objectRaw["VersionId"])
	d.Set("secret_name", objectRaw["SecretName"])

	objectRaw, err = kmsServiceV2.DescribeGetSecretPolicy(d.Id())
	if err != nil {
		return WrapError(err)
	}

	d.Set("policy", objectRaw["Policy"])

	d.Set("secret_name", d.Id())

	return nil
}

func resourceAliCloudKmsSecretUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false
	d.Partial(true)
	action := "UpdateSecret"
	conn, err := client.NewKmsClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	query["SecretName"] = d.Id()
	if d.HasChange("description") {
		update = true
		request["Description"] = d.Get("description")
	}

	if v, ok := d.GetOk("extended_config_custom_data"); ok {
		request["ExtendedConfig.CustomData"] = convertMapToJsonStringIgnoreError(v.(map[string]interface{}))
	}
	if update {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-01-20"), StringPointer("AK"), query, request, &runtime)

			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(action, response, request)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
	}
	update = false
	action = "SetSecretPolicy"
	conn, err = client.NewKmsClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	query["SecretName"] = d.Id()
	if v, ok := d.GetOk("policy_name"); ok {
		request["PolicyName"] = v
	}
	if d.HasChange("policy") {
		update = true
		request["Policy"] = d.Get("policy")
	}

	if update {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-01-20"), StringPointer("AK"), query, request, &runtime)

			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(action, response, request)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
	}

	d.Partial(false)
	return resourceAliCloudKmsSecretRead(d, meta)
}

func resourceAliCloudKmsSecretDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	action := "DeleteSecret"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	conn, err := client.NewKmsClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	query["SecretName"] = d.Id()

	request["ForceDeleteWithoutRecovery"] = "true"
	if v, ok := d.GetOk("recovery_window_in_days"); ok {
		request["RecoveryWindowInDays"] = v
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-01-20"), StringPointer("AK"), query, request, &runtime)

		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return nil
}

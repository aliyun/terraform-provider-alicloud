// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
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
			"automatic_rotation": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"certificate_id": {
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
			"enable_deletion_protection": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"key_spec": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"key_usage": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"origin": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"pending_window_in_days": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: IntBetween(7, 366),
			},
			"policy": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"policy_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"protection_level": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"rotation_interval": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"secret_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"tags": tagsSchema(),
		},
	}
}

func resourceAliCloudKmsKeyCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateKey"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	conn, err := client.NewKmsClient()
	if err != nil {
		return WrapError(err)
	}
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
	if v, ok := d.GetOk("tags"); ok {
		jsonPathResult6, err := jsonpath.Get("$", v)
		if err == nil && jsonPathResult6 != "" {
			request["Tags"] = jsonPathResult6
		}
	}
	if v, ok := d.GetOk("policy"); ok {
		request["Policy"] = v
	}
	if v, ok := d.GetOkExists("enable_automatic_rotation"); ok {
		request["EnableAutomaticRotation"] = v
	}
	if v, ok := d.GetOk("dkms_instance_id"); ok {
		request["DKMSInstanceId"] = v
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_kms_key", action, AlibabaCloudSdkGoERROR)
	}

	id, _ := jsonpath.Get("$.KeyMetadata.KeyId", response)
	d.SetId(fmt.Sprint(id))

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

	d.Set("automatic_rotation", objectRaw["AutomaticRotation"])
	d.Set("description", objectRaw["Description"])
	d.Set("dkms_instance_id", objectRaw["DKMSInstanceId"])
	d.Set("key_spec", objectRaw["KeySpec"])
	d.Set("key_usage", objectRaw["KeyUsage"])
	d.Set("origin", objectRaw["Origin"])
	d.Set("protection_level", objectRaw["ProtectionLevel"])
	d.Set("rotation_interval", objectRaw["RotationInterval"])
	d.Set("status", objectRaw["KeyState"])

	objectRaw, err = kmsServiceV2.DescribeListResourceTags(d.Id())
	if err != nil {
		return WrapError(err)
	}

	tagsMaps, _ := jsonpath.Get("$.Tags.Tag", objectRaw)
	d.Set("tags", tagsToMap(tagsMaps))

	objectRaw, err = kmsServiceV2.DescribeGetKeyPolicy(d.Id())
	if err != nil {
		return WrapError(err)
	}

	d.Set("policy", objectRaw["Policy"])

	return nil
}

func resourceAliCloudKmsKeyUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false
	d.Partial(true)
	action := "UpdateKeyDescription"
	conn, err := client.NewKmsClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	query["KeyId"] = d.Id()
	if !d.IsNewResource() && d.HasChange("description") {
		update = true
		request["Description"] = d.Get("description")
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
	action = "UpdateRotationPolicy"
	conn, err = client.NewKmsClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	query["KeyId"] = d.Id()
	if d.HasChange("automatic_rotation") {
		update = true
		request["EnableAutomaticRotation"] = convertKmsKeyEnableAutomaticRotationRequest(d.Get("automatic_rotation").(string))
	}

	if !d.IsNewResource() && d.HasChange("rotation_interval") {
		update = true
		request["RotationInterval"] = d.Get("rotation_interval")
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
	action = "SetKeyPolicy"
	conn, err = client.NewKmsClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	query["KeyId"] = d.Id()
	if v, ok := d.GetOk("policy_name"); ok {
		request["PolicyName"] = v
	}
	if !d.IsNewResource() && d.HasChange("policy") {
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

	if d.HasChange("status") {
		client := meta.(*connectivity.AliyunClient)
		kmsServiceV2 := KmsServiceV2{client}
		object, err := kmsServiceV2.DescribeKmsKey(d.Id())
		if err != nil {
			return WrapError(err)
		}

		target := d.Get("status").(string)
		if object["KeyState"].(string) != target {
			if target == "Disabled" {
				action = "DisableKey"
				conn, err = client.NewKmsClient()
				if err != nil {
					return WrapError(err)
				}
				request = make(map[string]interface{})
				query = make(map[string]interface{})
				query["KeyId"] = d.Id()
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
			if target == "Enabled" {
				action = "CancelKeyDeletion"
				conn, err = client.NewKmsClient()
				if err != nil {
					return WrapError(err)
				}
				request = make(map[string]interface{})
				query = make(map[string]interface{})
				query["KeyId"] = d.Id()
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
			if target == "PendingDeletion" {
				action = "ScheduleKeyDeletion"
				conn, err = client.NewKmsClient()
				if err != nil {
					return WrapError(err)
				}
				request = make(map[string]interface{})
				query = make(map[string]interface{})
				query["KeyId"] = d.Id()
				if v, ok := d.GetOk("pending_window_in_days"); ok {
					request["PendingWindowInDays"] = v
				}
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
		}
	}

	d.Partial(false)
	return resourceAliCloudKmsKeyRead(d, meta)
}

func resourceAliCloudKmsKeyDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[WARN] Cannot destroy resource AliCloud Resource Key. Terraform will remove this resource from the state file, however resources may remain.")
	return nil
}

func convertKmsKeyEnableAutomaticRotationRequest(source interface{}) interface{} {
	switch source {
	case "Enabled":
		return "true"
	case "Disabled":
		return "false"
	}
	return source
}

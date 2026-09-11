// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudApigSecret() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudApigSecretCreate,
		Read:   resourceAliCloudApigSecretRead,
		Update: resourceAliCloudApigSecretUpdate,
		Delete: resourceAliCloudApigSecretDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"create_timestamp": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"gateway_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"kms_config": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"kms_instance_id": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"kms_key_id": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
					},
				},
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"reference_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"secret_source": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"secret_data": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"update_timestamp": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func resourceAliCloudApigSecretCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := fmt.Sprintf("/v1/secrets")
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	body := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})

	if v, ok := d.GetOk("kms_config"); ok {
		kmsConfigMap := map[string]interface{}{}
		for _, kmsConfigList := range convertToInterfaceArray(v) {
			kmsConfigArg := kmsConfigList.(map[string]interface{})
			if kmsInstanceId, ok := kmsConfigArg["kms_instance_id"]; ok {
				kmsConfigMap["kmsInstanceId"] = kmsInstanceId
			}

			if kmsKeyId, ok := kmsConfigArg["kms_key_id"]; ok {
				kmsConfigMap["kmsKeyId"] = kmsKeyId
			}
		}

		request["kmsConfig"] = kmsConfigMap
	}

	if v, ok := d.GetOk("description"); ok {
		request["description"] = v
	}
	if v, ok := d.GetOk("secret_source"); ok {
		request["secretSource"] = v
	}
	if v, ok := d.GetOk("gateway_type"); ok {
		request["gatewayType"] = v
	}
	request["name"] = d.Get("name")
	if v, ok := d.GetOk("secret_data"); ok {
		request["secretData"] = v
	}
	body = request
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RoaPost("APIG", "2024-03-27", action, query, nil, body, true)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_apig_secret", action, AlibabaCloudSdkGoERROR)
	}

	id, _ := jsonpath.Get("$.data.secretId", response)
	d.SetId(fmt.Sprint(id))

	return resourceAliCloudApigSecretRead(d, meta)
}

func resourceAliCloudApigSecretRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	apigServiceV2 := ApigServiceV2{client}

	objectRaw, err := apigServiceV2.DescribeApigSecret(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_apig_secret DescribeApigSecret Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("reference_count", objectRaw["referenceCount"])
	d.Set("status", objectRaw["status"])
	d.Set("gateway_type", objectRaw["gatewayType"])
	d.Set("name", objectRaw["name"])
	d.Set("secret_source", objectRaw["secretSource"])
	d.Set("description", objectRaw["description"])

	if kmsConfig, ok := objectRaw["kmsConfig"]; ok {
		kmsConfigsMaps := make([]map[string]interface{}, 0)
		kmsConfigMap := map[string]interface{}{}
		kmsConfigArg := kmsConfig.(map[string]interface{})

		if kmsInstanceId, ok := kmsConfigArg["kmsInstanceId"]; ok {
			kmsConfigMap["kms_instance_id"] = kmsInstanceId
		}

		if kmsKeyId, ok := kmsConfigArg["kmsKeyId"]; ok {
			kmsConfigMap["kms_key_id"] = kmsKeyId
		}

		kmsConfigsMaps = append(kmsConfigsMaps, kmsConfigMap)

		d.Set("kms_config", kmsConfigsMaps)
	}

	createTimestamp, err := strconv.ParseInt(objectRaw["createTimestamp"].(json.Number).String(), 10, 64)
	if err != nil {
		return WrapError(err)
	}
	d.Set("create_timestamp", createTimestamp)

	updateTimestamp, err := strconv.ParseInt(objectRaw["updateTimestamp"].(json.Number).String(), 10, 64)
	if err != nil {
		return WrapError(err)
	}
	d.Set("update_timestamp", updateTimestamp)

	return nil
}

func resourceAliCloudApigSecretUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	var body map[string]interface{}
	update := false

	var err error
	secretId := d.Id()
	action := fmt.Sprintf("/v1/secrets/%s", secretId)
	request = make(map[string]interface{})
	query = make(map[string]*string)
	body = make(map[string]interface{})

	if v, ok := d.GetOk("secret_data"); ok {
		request["secretData"] = v
	}
	body = request
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RoaPut("APIG", "2024-03-27", action, query, nil, body, true)
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

	return resourceAliCloudApigSecretRead(d, meta)
}

func resourceAliCloudApigSecretDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	secretId := d.Id()
	action := fmt.Sprintf("/v1/secrets/%s", secretId)
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	var err error
	request = make(map[string]interface{})

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RoaDelete("APIG", "2024-03-27", action, query, nil, nil, true)
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

// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudKmsClientKey() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudKmsClientKeyCreate,
		Read:   resourceAliCloudKmsClientKeyRead,
		Delete: resourceAliCloudKmsClientKeyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"aap_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"not_after": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: ValidateRFC3339TimeString(true),
			},
			"not_before": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: ValidateRFC3339TimeString(true),
			},
			"password": {
				Type:      schema.TypeString,
				Required:  true,
				ForceNew:  true,
				Sensitive: true,
			},
			"private_key_data_file": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAliCloudKmsClientKeyCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateClientKey"
	var request map[string]interface{}
	var response map[string]interface{}
	var err error
	request = make(map[string]interface{})

	request["AapName"] = d.Get("aap_name")
	if v, ok := d.GetOk("not_before"); ok {
		request["NotBefore"] = v
	}
	if v, ok := d.GetOk("not_after"); ok {
		request["NotAfter"] = v
	}
	request["Password"] = d.Get("password")
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("Kms", "2016-01-20", action, nil, request, false)

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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_kms_client_key", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["ClientKeyId"]))
	privateKeyData := map[string]interface{}{
		"KeyId":          response["ClientKeyId"],
		"PrivateKeyData": response["PrivateKeyData"],
	}
	jsonData, _ := json.Marshal(privateKeyData)
	if output, ok := d.GetOk("private_key_data_file"); ok && output != nil {
		// create a private_key_data_file and write private key to it.
		writeToFile(output.(string), string(jsonData))
	}

	return resourceAliCloudKmsClientKeyRead(d, meta)
}

func resourceAliCloudKmsClientKeyRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	kmsServiceV2 := KmsServiceV2{client}

	objectRaw, err := kmsServiceV2.DescribeKmsClientKey(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_kms_client_key DescribeKmsClientKey Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("aap_name", objectRaw["AapName"])
	d.Set("create_time", objectRaw["CreateTime"])
	d.Set("not_after", objectRaw["NotAfter"])
	d.Set("not_before", objectRaw["NotBefore"])

	return nil
}

func resourceAliCloudKmsClientKeyUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Cannot update resource Alicloud Resource Client Key.")
	return nil
}

func resourceAliCloudKmsClientKeyDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	action := "DeleteClientKey"
	var request map[string]interface{}
	var response map[string]interface{}
	var err error
	request = make(map[string]interface{})
	request["ClientKeyId"] = d.Id()

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("Kms", "2016-01-20", action, nil, request, false)

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
		if NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return nil
}

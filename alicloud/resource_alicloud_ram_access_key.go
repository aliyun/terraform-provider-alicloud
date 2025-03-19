// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/encryption"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudRamAccessKey() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudRamAccessKeyCreate,
		Read:   resourceAliCloudRamAccessKeyRead,
		Update: resourceAliCloudRamAccessKeyUpdate,
		Delete: resourceAliCloudRamAccessKeyDelete,
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
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: StringInSlice([]string{"Active", "Inactive"}, false),
			},
			"user_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"secret_file": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"secret": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
			"pgp_key": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"key_fingerprint": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"encrypted_secret": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAliCloudRamAccessKeyCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateAccessKey"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})

	if v, ok := d.GetOk("user_name"); ok {
		request["UserName"] = v
	}
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("Ram", "2015-05-01", action, query, request, true)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_ram_access_key", action, AlibabaCloudSdkGoERROR)
	}

	resp, err := jsonpath.Get("$.AccessKey", response)
	if err != nil {
		return WrapErrorf(err, FailedGetAttributeMsg, "alicloud_ram_access_key", "$.AccessKey", response)
	}

	accessKeyId := resp.(map[string]interface{})["AccessKeyId"]
	accessKeySecret := resp.(map[string]interface{})["AccessKeySecret"]

	d.SetId(fmt.Sprint(accessKeyId))

	if v, ok := d.GetOk("pgp_key"); ok {
		pgpKey := v.(string)
		encryptionKey, err := encryption.RetrieveGPGKey(pgpKey)
		if err != nil {
			return WrapError(err)
		}

		fingerprint, encrypted, err := encryption.EncryptValue(encryptionKey, fmt.Sprint(accessKeySecret), "Alicloud RAM Access Key Secret")
		if err != nil {
			return WrapError(err)
		}
		d.Set("key_fingerprint", fingerprint)
		d.Set("encrypted_secret", encrypted)
	} else {
		if err := d.Set("secret", fmt.Sprint(accessKeySecret)); err != nil {
			return WrapError(err)
		}
	}

	if output, ok := d.GetOk("secret_file"); ok && output != nil {
		// create a secret_file and write access key to it.
		writeToFile(output.(string), resp)
	}

	return resourceAliCloudRamAccessKeyUpdate(d, meta)
}

func resourceAliCloudRamAccessKeyRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ramServiceV2 := RamServiceV2{client}
	userName := ""
	if v, ok := d.GetOk("user_name"); ok && v.(string) != "" {
		userName = v.(string)
	}

	objectRaw, err := ramServiceV2.DescribeRamAccessKey(d.Id(), userName)
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_ram_access_key DescribeRamAccessKey Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("create_time", objectRaw["CreateDate"])
	d.Set("status", objectRaw["Status"])

	return nil
}

func resourceAliCloudRamAccessKeyUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false

	var err error
	action := "UpdateAccessKey"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["UserAccessKeyId"] = d.Id()

	if d.HasChange("status") {
		update = true
	}
	if v, ok := d.GetOk("status"); ok {
		request["Status"] = v
	}

	if v, ok := d.GetOk("user_name"); ok {
		request["UserName"] = v
	}

	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("Ram", "2015-05-01", action, query, request, true)
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

	return resourceAliCloudRamAccessKeyRead(d, meta)
}

func resourceAliCloudRamAccessKeyDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	action := "DeleteAccessKey"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["UserAccessKeyId"] = d.Id()

	if v, ok := d.GetOk("user_name"); ok {
		request["UserName"] = v
	}
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("Ram", "2015-05-01", action, query, request, true)

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
		if IsExpectedErrors(err, []string{"EntityNotExist.User.AccessKey", "EntityNotExist.User"}) || NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return nil
}

// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package ram

import (
	"fmt"
	"log"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/validators"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceAccessKey() *schema.Resource {
	return &schema.Resource{
		Create: ResourceAccessKeyCreate,
		Read:   ResourceAccessKeyRead,
		Update: ResourceAccessKeyUpdate,
		Delete: ResourceAccessKeyDelete,
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
				ValidateFunc: validators.StringInSlice([]string{"Active", "Inactive"}, false),
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

func ResourceAccessKeyCreate(d *schema.ResourceData, meta interface{}) error {

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
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		response, err = client.RpcPost("Ram", "2015-05-01", action, query, request, true)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)

	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_ram_access_key", action, errs.AlibabaCloudSdkGoERROR)
	}

	resp, err := jsonpath.Get("$.AccessKey", response)
	if err != nil {
		return errs.WrapErrorf(err, errs.FailedGetAttributeMsg, "alicloud_ram_access_key", "$.AccessKey", response)
	}

	accessKeyId := resp.(map[string]interface{})["AccessKeyId"]
	accessKeySecret := resp.(map[string]interface{})["AccessKeySecret"]

	d.SetId(fmt.Sprint(accessKeyId))

	if v, ok := d.GetOk("pgp_key"); ok {
		pgpKey := v.(string)
		encryptionKey, err := helper.RetrieveGPGKey(pgpKey)
		if err != nil {
			return errs.WrapError(err)
		}

		fingerprint, encrypted, err := helper.EncryptValue(encryptionKey, fmt.Sprint(accessKeySecret), "Alicloud RAM Access Key Secret")
		if err != nil {
			return errs.WrapError(err)
		}
		d.Set("key_fingerprint", fingerprint)
		d.Set("encrypted_secret", encrypted)
	} else {
		if err := d.Set("secret", fmt.Sprint(accessKeySecret)); err != nil {
			return errs.WrapError(err)
		}
	}

	if output, ok := d.GetOk("secret_file"); ok && output != nil {
		// create a secret_file and write access key to it.
		helper.WriteToFile(output.(string), resp)
	}

	return ResourceAccessKeyUpdate(d, meta)
}

func ResourceAccessKeyRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ramServiceV2 := RamServiceV2{client}
	userName := ""
	if v, ok := d.GetOk("user_name"); ok && v.(string) != "" {
		userName = v.(string)
	}

	objectRaw, err := ramServiceV2.DescribeRamAccessKey(d.Id(), userName)
	if err != nil {
		if !d.IsNewResource() && errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_ram_access_key DescribeRamAccessKey Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}

	d.Set("create_time", objectRaw["CreateDate"])
	d.Set("status", objectRaw["Status"])

	return nil
}

func ResourceAccessKeyUpdate(d *schema.ResourceData, meta interface{}) error {
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
		wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RpcPost("Ram", "2015-05-01", action, query, request, true)
			if err != nil {
				if helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		logs.AddDebug(action, response, request)
		if err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
		}
	}

	return ResourceAccessKeyRead(d, meta)
}

func ResourceAccessKeyDelete(d *schema.ResourceData, meta interface{}) error {

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
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		response, err = client.RpcPost("Ram", "2015-05-01", action, query, request, true)

		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)

	if err != nil {
		if errs.IsExpectedErrors(err, []string{"EntityNotExist.User.AccessKey", "EntityNotExist.User"}) || errs.NotFoundError(err) {
			return nil
		}
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
	}

	return nil
}

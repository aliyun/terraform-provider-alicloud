// Package cloudphone. This file is generated automatically. Please do not modify it manually, thank you!
package cloudphone

import (
	"fmt"
	"log"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceKeyPair() *schema.Resource {
	return &schema.Resource{
		Create: resourceKeyPairCreate,
		Read:   resourceKeyPairRead,
		Update: resourceKeyPairUpdate,
		Delete: resourceKeyPairDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"key_pair_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"public_key_body": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceKeyPairCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	invalidCreate := false
	if _, ok := d.GetOk("public_key_body"); ok {
		invalidCreate = true
	}
	if !invalidCreate {

		action := "CreateKeyPair"
		var request map[string]interface{}
		var response map[string]interface{}
		query := make(map[string]interface{})
		var err error
		request = make(map[string]interface{})

		request["KeyPairName"] = d.Get("key_pair_name")
		wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
			response, err = client.RpcPost("eds-aic", "2023-09-30", action, query, request, true)
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
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_cloud_phone_key_pair", action, errs.AlibabaCloudSdkGoERROR)
		}

		id, _ := jsonpath.Get("$.Data.KeyPairId", response)
		d.SetId(fmt.Sprint(id))

	}

	if v, ok := d.GetOk("public_key_body"); ok && fmt.Sprint(v) != "" {

		action := "ImportKeyPair"
		var request map[string]interface{}
		var response map[string]interface{}
		query := make(map[string]interface{})
		var err error
		request = make(map[string]interface{})

		request["PublicKeyBody"] = d.Get("public_key_body")
		request["KeyPairName"] = d.Get("key_pair_name")
		wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
			response, err = client.RpcPost("eds-aic", "2023-09-30", action, query, request, true)
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
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_cloud_phone_key_pair", action, errs.AlibabaCloudSdkGoERROR)
		}

		id, _ := jsonpath.Get("$.Data.KeyPairId", response)
		d.SetId(fmt.Sprint(id))

	}

	return resourceKeyPairRead(d, meta)
}

func resourceKeyPairRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cloudPhoneService := NewService(client)

	objectRaw, err := cloudPhoneService.DescribeKeyPair(d.Id())
	if err != nil {
		if !d.IsNewResource() && errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_cloud_phone_key_pair DescribeKeyPair Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}

	d.Set("key_pair_name", objectRaw["KeyPairName"])

	return nil
}

func resourceKeyPairUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false

	var err error
	action := "ModifyKeyPairName"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["KeyPairId"] = d.Id()

	if d.HasChange("key_pair_name") {
		update = true
	}
	request["NewKeyPairName"] = d.Get("key_pair_name")
	if update {
		wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RpcPost("eds-aic", "2023-09-30", action, query, request, true)
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

	return resourceKeyPairRead(d, meta)
}

func resourceKeyPairDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	action := "DeleteKeyPairs"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["KeyPairIds.1"] = d.Id()

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		response, err = client.RpcPost("eds-aic", "2023-09-30", action, query, request, true)

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
		if errs.NotFoundError(err) {
			return nil
		}
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
	}

	return nil
}

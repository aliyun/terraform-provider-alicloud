// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudOssBucketRequestPayment() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudOssBucketRequestPaymentCreate,
		Read:   resourceAliCloudOssBucketRequestPaymentRead,
		Update: resourceAliCloudOssBucketRequestPaymentUpdate,
		Delete: resourceAliCloudOssBucketRequestPaymentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"bucket": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"payer": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: StringInSlice([]string{"Requester", "BucketOwner"}, true),
			},
		},
	}
}

func resourceAliCloudOssBucketRequestPaymentCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := fmt.Sprintf("/?requestPayment")
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	body := make(map[string]interface{})
	hostMap := make(map[string]*string)
	var err error
	request = make(map[string]interface{})
	hostMap["bucket"] = StringPointer(d.Get("bucket").(string))

	objectDataLocalMap := make(map[string]interface{})
	if v := d.Get("payer"); !IsNil(v) {
		objectDataLocalMap["Payer"] = v
		request["RequestPaymentConfiguration"] = objectDataLocalMap
	}

	body = request
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.Do("Oss", xmlParam("PUT", "2019-05-17", "PutBucketRequestPayment", action), query, body, nil, hostMap, false)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_oss_bucket_request_payment", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(*hostMap["bucket"]))

	return resourceAliCloudOssBucketRequestPaymentRead(d, meta)
}

func resourceAliCloudOssBucketRequestPaymentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ossServiceV2 := OssServiceV2{client}

	objectRaw, err := ossServiceV2.DescribeOssBucketRequestPayment(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_oss_bucket_request_payment DescribeOssBucketRequestPayment Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("payer", objectRaw["Payer"])

	d.Set("bucket", d.Id())

	return nil
}

func resourceAliCloudOssBucketRequestPaymentUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	var body map[string]interface{}
	update := false
	action := fmt.Sprintf("/?requestPayment")
	var err error
	request = make(map[string]interface{})
	query = make(map[string]*string)
	body = make(map[string]interface{})
	hostMap := make(map[string]*string)
	hostMap["bucket"] = StringPointer(d.Id())
	if d.HasChange("payer") {
		update = true
	}
	objectDataLocalMap := make(map[string]interface{})
	if v := d.Get("payer"); v != nil {
		objectDataLocalMap["Payer"] = d.Get("payer")
		request["RequestPaymentConfiguration"] = objectDataLocalMap
	}

	body = request
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.Do("Oss", xmlParam("PUT", "2019-05-17", "PutBucketRequestPayment", action), query, body, nil, hostMap, false)
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

	return resourceAliCloudOssBucketRequestPaymentRead(d, meta)
}

func resourceAliCloudOssBucketRequestPaymentDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[WARN] Cannot destroy resource AliCloud Resource Bucket Request Payment. Terraform will remove this resource from the state file, however resources may remain.")
	return nil
}

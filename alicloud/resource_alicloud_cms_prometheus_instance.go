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

func resourceAliCloudCmsPrometheusInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudCmsPrometheusInstanceCreate,
		Read:   resourceAliCloudCmsPrometheusInstanceRead,
		Update: resourceAliCloudCmsPrometheusInstanceUpdate,
		Delete: resourceAliCloudCmsPrometheusInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"archive_duration": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: IntBetween(0, 3650),
			},
			"auth_free_read_policy": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"auth_free_write_policy": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"enable_auth_free_read": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"enable_auth_free_write": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"payment_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"prometheus_instance_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"region_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"storage_duration": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: IntBetween(15, 3650),
			},
			"workspace": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAliCloudCmsPrometheusInstanceCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := fmt.Sprintf("/prometheus-instances")
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	body := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})

	request["prometheusInstanceName"] = d.Get("prometheus_instance_name")
	if v, ok := d.GetOkExists("enable_auth_free_write"); ok {
		request["enableAuthFreeWrite"] = v
	}
	if v, ok := d.GetOkExists("archive_duration"); ok {
		request["archiveDuration"] = v
	}
	if v, ok := d.GetOkExists("enable_auth_free_read"); ok {
		request["enableAuthFreeRead"] = v
	}
	if v, ok := d.GetOk("auth_free_read_policy"); ok {
		request["authFreeReadPolicy"] = v
	}
	if v, ok := d.GetOk("auth_free_write_policy"); ok {
		request["authFreeWritePolicy"] = v
	}
	request["workspace"] = d.Get("workspace")
	if v, ok := d.GetOkExists("storage_duration"); ok && v.(int) > 0 {
		request["storageDuration"] = v
	}
	body = request
	wait := incrementalWait(10*time.Second, 10*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RoaPost("Cms", "2024-03-30", action, query, nil, body, true)
		if err != nil {
			if IsExpectedErrors(err, []string{"500"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_cms_prometheus_instance", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["prometheusInstanceId"]))

	return resourceAliCloudCmsPrometheusInstanceRead(d, meta)
}

func resourceAliCloudCmsPrometheusInstanceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cmsServiceV2 := CmsServiceV2{client}

	objectRaw, err := cmsServiceV2.DescribeCmsPrometheusInstance(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_cms_prometheus_instance DescribeCmsPrometheusInstance Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("archive_duration", objectRaw["archiveDuration"])
	d.Set("auth_free_read_policy", objectRaw["authFreeReadPolicy"])
	d.Set("auth_free_write_policy", objectRaw["authFreeWritePolicy"])
	d.Set("create_time", objectRaw["createTime"])
	d.Set("enable_auth_free_read", objectRaw["enableAuthFreeRead"])
	d.Set("enable_auth_free_write", objectRaw["enableAuthFreeWrite"])
	d.Set("payment_type", objectRaw["paymentType"])
	d.Set("prometheus_instance_name", objectRaw["prometheusInstanceName"])
	d.Set("region_id", objectRaw["regionId"])
	d.Set("storage_duration", objectRaw["storageDuration"])
	d.Set("workspace", objectRaw["workspace"])

	return nil
}

func resourceAliCloudCmsPrometheusInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	var body map[string]interface{}
	update := false

	var err error
	prometheusInstanceId := d.Id()
	action := fmt.Sprintf("/prometheus-instances/%s", prometheusInstanceId)
	request = make(map[string]interface{})
	query = make(map[string]*string)
	body = make(map[string]interface{})

	if d.HasChange("prometheus_instance_name") {
		update = true
	}
	request["prometheusInstanceName"] = d.Get("prometheus_instance_name")
	if d.HasChange("enable_auth_free_write") {
		update = true
	}
	if v, ok := d.GetOkExists("enable_auth_free_write"); ok || d.HasChange("enable_auth_free_write") {
		request["enableAuthFreeWrite"] = v
	}
	if d.HasChange("archive_duration") {
		update = true
	}
	if v, ok := d.GetOkExists("archive_duration"); ok || d.HasChange("archive_duration") {
		request["archiveDuration"] = v
	}
	if d.HasChange("enable_auth_free_read") {
		update = true
	}
	if v, ok := d.GetOkExists("enable_auth_free_read"); ok || d.HasChange("enable_auth_free_read") {
		request["enableAuthFreeRead"] = v
	}
	if d.HasChange("auth_free_read_policy") {
		update = true
	}
	if v, ok := d.GetOk("auth_free_read_policy"); ok || d.HasChange("auth_free_read_policy") {
		request["authFreeReadPolicy"] = v
	}
	if d.HasChange("auth_free_write_policy") {
		update = true
	}
	if v, ok := d.GetOk("auth_free_write_policy"); ok || d.HasChange("auth_free_write_policy") {
		request["authFreeWritePolicy"] = v
	}
	if d.HasChange("storage_duration") {
		update = true
	}
	if v, ok := d.GetOkExists("storage_duration"); (ok || d.HasChange("storage_duration")) && v.(int) > 0 {
		request["storageDuration"] = v
	}
	body = request
	if update {
		wait := incrementalWait(10*time.Second, 10*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RoaPut("Cms", "2024-03-30", action, query, nil, body, true)
			if err != nil {
				if IsExpectedErrors(err, []string{"500"}) || NeedRetry(err) {
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

	return resourceAliCloudCmsPrometheusInstanceRead(d, meta)
}

func resourceAliCloudCmsPrometheusInstanceDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	prometheusInstanceId := d.Id()
	action := fmt.Sprintf("/prometheus-instances/%s", prometheusInstanceId)
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	var err error
	request = make(map[string]interface{})

	wait := incrementalWait(10*time.Second, 10*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RoaDelete("Cms", "2024-03-30", action, query, nil, nil, true)
		if err != nil {
			if IsExpectedErrors(err, []string{"500"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		if IsExpectedErrors(err, []string{"404"}) || NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return nil
}

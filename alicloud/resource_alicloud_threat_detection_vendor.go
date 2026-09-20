package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudThreatDetectionVendor() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudThreatDetectionVendorCreate,
		Read:   resourceAlicloudThreatDetectionVendorRead,
		Update: resourceAlicloudThreatDetectionVendorUpdate,
		Delete: resourceAlicloudThreatDetectionVendorDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(3 * time.Minute),
			Update: schema.DefaultTimeout(3 * time.Minute),
			Delete: schema.DefaultTimeout(3 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"vendor_id": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"vendor_name": {
				Required: true,
				Type:     schema.TypeString,
			},
			"vendor_type": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"create_time": {
				Computed: true,
				Type:     schema.TypeInt,
			},
			"update_time": {
				Computed: true,
				Type:     schema.TypeInt,
			},
			"lang": {
				Optional: true,
				Default:  "en",
				Type:     schema.TypeString,
			},
			"role_for": {
				Optional: true,
				Type:     schema.TypeInt,
			},
		},
	}
}

func resourceAlicloudThreatDetectionVendorCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	request := make(map[string]interface{})

	if v, ok := d.GetOk("lang"); ok {
		request["Lang"] = v
	}
	if v, ok := d.GetOk("role_for"); ok {
		request["RoleFor"] = v
	}
	if v, ok := d.GetOk("vendor_name"); ok {
		request["VendorName"] = v
	}

	var response map[string]interface{}
	action := "CreateVendor"
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err := resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *resource.RetryError {
		resp, err := client.RpcPost("cloud-siem", "2024-12-12", action, nil, request, false)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		response = resp
		addDebug(action, response, request)
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_threat_detection_vendor", action, AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$.VendorId", response)
	if err != nil || v == nil {
		return WrapErrorf(err, IdMsg, "alicloud_threat_detection_vendor")
	}
	d.SetId(fmt.Sprint(v))

	return resourceAlicloudThreatDetectionVendorRead(d, meta)
}

func resourceAlicloudThreatDetectionVendorRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cloudSiemService := CloudSiemService{client}

	object, err := cloudSiemService.DescribeThreatDetectionVendor(d.Id(), d.Get("lang").(string))
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_threat_detection_vendor DescribeThreatDetectionVendor Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("vendor_id", object["VendorId"])
	d.Set("vendor_name", object["VendorName"])
	d.Set("vendor_type", object["VendorType"])
	if createTime, err := jsonpath.Get("$.CreateTime", object); err == nil {
		d.Set("create_time", createTime)
	}
	if updateTime, err := jsonpath.Get("$.UpdateTime", object); err == nil {
		d.Set("update_time", updateTime)
	}
	return nil
}

func resourceAlicloudThreatDetectionVendorUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	request := make(map[string]interface{})
	update := false

	request["VendorId"] = d.Id()
	request["Lang"] = d.Get("lang")

	if d.HasChange("vendor_name") {
		update = true
		request["VendorName"] = d.Get("vendor_name")
	}
	if d.HasChange("lang") {
		update = true
	}
	if d.HasChange("role_for") {
		update = true
		request["RoleFor"] = d.Get("role_for")
	}

	if update {
		action := "UpdateVendor"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err := resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			resp, err := client.RpcPost("cloud-siem", "2024-12-12", action, nil, request, false)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(action, resp, request)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
	}

	return resourceAlicloudThreatDetectionVendorRead(d, meta)
}

func resourceAlicloudThreatDetectionVendorDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := map[string]interface{}{
		"VendorId": d.Id(),
		"Lang":     d.Get("lang"),
	}

	action := "DeleteVendor"
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err := resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *resource.RetryError {
		resp, err := client.RpcPost("cloud-siem", "2024-12-12", action, nil, request, false)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, resp, request)
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

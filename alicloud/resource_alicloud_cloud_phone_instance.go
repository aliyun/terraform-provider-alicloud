// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
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

func resourceAliCloudCloudPhoneInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudCloudPhoneInstanceCreate,
		Read:   resourceAliCloudCloudPhoneInstanceRead,
		Update: resourceAliCloudCloudPhoneInstanceUpdate,
		Delete: resourceAliCloudCloudPhoneInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"android_instance_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"android_instance_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceAliCloudCloudPhoneInstanceCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "DescribeAndroidInstances"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})

	if v, ok := d.GetOk("android_instance_group_id"); ok {
		request["InstanceGroupId"] = v
	}
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("eds-aic", "2023-09-30", action, query, request, true)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_cloud_phone_instance", action, AlibabaCloudSdkGoERROR)
	}

	id, _ := jsonpath.Get("$.InstanceModel[0].AndroidInstanceId", response)
	d.SetId(fmt.Sprint(id))

	return resourceAliCloudCloudPhoneInstanceUpdate(d, meta)
}

func resourceAliCloudCloudPhoneInstanceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cloudPhoneServiceV2 := CloudPhoneServiceV2{client}

	objectRaw, err := cloudPhoneServiceV2.DescribeCloudPhoneInstance(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_cloud_phone_instance DescribeCloudPhoneInstance Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("android_instance_group_id", objectRaw["AndroidInstanceGroupId"])
	d.Set("android_instance_name", objectRaw["AndroidInstanceName"])

	return nil
}

func resourceAliCloudCloudPhoneInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false

	var err error
	action := "ModifyAndroidInstance"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["AndroidInstanceId"] = d.Id()

	if d.HasChange("android_instance_name") {
		update = true
	}
	if v, ok := d.GetOk("android_instance_name"); ok || d.HasChange("android_instance_name") {
		request["NewAndroidInstanceName"] = v
	}
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("eds-aic", "2023-09-30", action, query, request, true)
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

	return resourceAliCloudCloudPhoneInstanceRead(d, meta)
}

func resourceAliCloudCloudPhoneInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[WARN] Cannot destroy resource AliCloud Resource Instance. Terraform will remove this resource from the state file, however resources may remain.")
	return nil
}

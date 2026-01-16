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

func resourceAliCloudAckInspectConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudAckInspectConfigCreate,
		Read:   resourceAliCloudAckInspectConfigRead,
		Update: resourceAliCloudAckInspectConfigUpdate,
		Delete: resourceAliCloudAckInspectConfigDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"disabled_check_items": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"enabled": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"inspect_config_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"recurrence": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceAliCloudAckInspectConfigCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	clusterId := d.Get("inspect_config_id")
	action := fmt.Sprintf("/clusters/%s/inspectConfig", clusterId)
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	body := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})

	if v, ok := d.GetOk("disabled_check_items"); ok {
		disabledCheckItemsMapsArray := convertToInterfaceArray(v)

		request["disabledCheckItems"] = disabledCheckItemsMapsArray
	}

	request["enabled"] = d.Get("enabled")
	request["recurrence"] = d.Get("recurrence")
	body = request
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RoaPost("CS", "2015-12-15", action, query, nil, body, true)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_ack_inspect_config", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(clusterId))

	return resourceAliCloudAckInspectConfigRead(d, meta)
}

func resourceAliCloudAckInspectConfigRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ackServiceV2 := AckServiceV2{client}

	objectRaw, err := ackServiceV2.DescribeAckInspectConfig(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_ack_inspect_config DescribeAckInspectConfig Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("enabled", objectRaw["enabled"])
	d.Set("recurrence", objectRaw["recurrence"])

	disabledCheckItemsRaw := make([]interface{}, 0)
	if objectRaw["disabledCheckItems"] != nil {
		disabledCheckItemsRaw = convertToInterfaceArray(objectRaw["disabledCheckItems"])
	}

	d.Set("disabled_check_items", disabledCheckItemsRaw)

	d.Set("inspect_config_id", d.Id())

	return nil
}

func resourceAliCloudAckInspectConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var header map[string]*string
	var query map[string]*string
	var body map[string]interface{}
	update := false

	var err error
	clusterId := d.Id()
	action := fmt.Sprintf("/clusters/%s/inspectConfig", clusterId)
	request = make(map[string]interface{})
	query = make(map[string]*string)
	body = make(map[string]interface{})

	if d.HasChange("disabled_check_items") {
		update = true
		if v, ok := d.GetOk("disabled_check_items"); ok || d.HasChange("disabled_check_items") {
			disabledCheckItemsMapsArray := convertToInterfaceArray(v)

			request["disabledCheckItems"] = disabledCheckItemsMapsArray
		}
	}

	if d.HasChange("enabled") {
		update = true
	}
	request["enabled"] = d.Get("enabled")
	if d.HasChange("recurrence") {
		update = true
	}
	request["scheduleTime"] = d.Get("recurrence")
	body = request
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RoaPut("CS", "2015-12-15", action, query, header, body, true)
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

	return resourceAliCloudAckInspectConfigRead(d, meta)
}

func resourceAliCloudAckInspectConfigDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	clusterId := d.Id()
	action := fmt.Sprintf("/clusters/%s/inspectConfig", clusterId)
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	var err error
	request = make(map[string]interface{})

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RoaDelete("CS", "2015-12-15", action, query, nil, nil, true)
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

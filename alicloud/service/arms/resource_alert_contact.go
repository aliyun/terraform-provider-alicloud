package arms

import (
	"fmt"
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceAlertContact() *schema.Resource {
	return &schema.Resource{
		Create: ResourceAlertContactCreate,
		Read:   ResourceAlertContactRead,
		Update: ResourceAlertContactUpdate,
		Delete: ResourceAlertContactDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"alert_contact_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ding_robot_webhook_url": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"email": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"phone_num": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"system_noc": {
				Type:     schema.TypeBool,
				Optional: true,
			},
		},
	}
}

func ResourceAlertContactCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateAlertContact"
	request := make(map[string]interface{})
	var err error
	if v, ok := d.GetOk("alert_contact_name"); ok {
		request["ContactName"] = v
	}
	if v, ok := d.GetOk("ding_robot_webhook_url"); ok {
		request["DingRobotWebhookUrl"] = v
	} else if v, ok := d.GetOk("email"); ok && v.(string) == "" {
		if v, ok := d.GetOk("phone_num"); ok && v.(string) == "" {
			return errs.WrapError(fmt.Errorf("attribute '%s' is required when '%s' is %v and '%s' is %v ", "ding_robot_webhook_url", "email", d.Get("email"), "phone_num", d.Get("phone_num")))
		}
	}
	if v, ok := d.GetOk("email"); ok {
		request["Email"] = v
	} else if v, ok := d.GetOk("ding_robot_webhook_url"); ok && v.(string) == "" {
		if v, ok := d.GetOk("phone_num"); ok && v.(string) == "" {
			return errs.WrapError(fmt.Errorf("attribute '%s' is required when '%s' is %v and '%s' is %v ", "email", "ding_robot_webhook_url", d.Get("ding_robot_webhook_url"), "phone_num", d.Get("phone_num")))
		}
	}
	if v, ok := d.GetOk("phone_num"); ok {
		request["PhoneNum"] = v
	} else if v, ok := d.GetOk("ding_robot_webhook_url"); ok && v.(string) == "" {
		if v, ok := d.GetOk("email"); ok && v.(string) == "" {
			return errs.WrapError(fmt.Errorf("attribute '%s' is required when '%s' is %v and '%s' is %v ", "phone_num", "ding_robot_webhook_url", d.Get("ding_robot_webhook_url"), "email", d.Get("email")))
		}
	}
	request["RegionId"] = client.RegionId
	if v, ok := d.GetOkExists("system_noc"); ok {
		request["SystemNoc"] = v
	}
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		response, err = client.RpcPost("ARMS", "2019-08-08", action, nil, request, false)
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
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_arms_alert_contact", action, errs.AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["ContactId"]))

	return ResourceAlertContactRead(d, meta)
}
func ResourceAlertContactRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	armsService := Service{client}
	object, err := armsService.DescribeArmsAlertContact(d.Id())
	if err != nil {
		if errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_arms_alert_contact armsService.DescribeArmsAlertContact Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}
	d.Set("alert_contact_name", object["ContactName"])
	d.Set("ding_robot_webhook_url", object["DingRobot"])
	d.Set("email", object["Email"])
	d.Set("phone_num", object["Phone"])
	d.Set("system_noc", object["SystemNoc"])
	return nil
}
func ResourceAlertContactUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var err error
	var response map[string]interface{}
	update := false
	request := map[string]interface{}{
		"ContactId": d.Id(),
	}
	request["RegionId"] = client.RegionId
	if d.HasChange("alert_contact_name") {
		update = true
	}
	if v, ok := d.GetOk("alert_contact_name"); ok {
		request["ContactName"] = v
	}
	if d.HasChange("ding_robot_webhook_url") {
		update = true
	}
	if v, ok := d.GetOk("ding_robot_webhook_url"); ok {
		request["DingRobotWebhookUrl"] = v
	} else if v, ok := d.GetOk("email"); ok && v.(string) == "" {
		if v, ok := d.GetOk("phone_num"); ok && v.(string) == "" {
			return errs.WrapError(fmt.Errorf("attribute '%s' is required when '%s' is %v and '%s' is %v ", "ding_robot_webhook_url", "email", d.Get("email"), "phone_num", d.Get("phone_num")))
		}
	}
	if d.HasChange("email") {
		update = true
	}
	if v, ok := d.GetOk("email"); ok {
		request["Email"] = v
	} else if v, ok := d.GetOk("ding_robot_webhook_url"); ok && v.(string) == "" {
		if v, ok := d.GetOk("phone_num"); ok && v.(string) == "" {
			return errs.WrapError(fmt.Errorf("attribute '%s' is required when '%s' is %v and '%s' is %v ", "email", "ding_robot_webhook_url", d.Get("ding_robot_webhook_url"), "phone_num", d.Get("phone_num")))
		}
	}
	if d.HasChange("phone_num") {
		update = true
	}
	if v, ok := d.GetOk("phone_num"); ok {
		request["PhoneNum"] = v
	} else if v, ok := d.GetOk("ding_robot_webhook_url"); ok && v.(string) == "" {
		if v, ok := d.GetOk("email"); ok && v.(string) == "" {
			return errs.WrapError(fmt.Errorf("attribute '%s' is required when '%s' is %v and '%s' is %v ", "phone_num", "ding_robot_webhook_url", d.Get("ding_robot_webhook_url"), "email", d.Get("email")))
		}
	}
	if d.HasChange("system_noc") || d.IsNewResource() {
		update = true
	}
	if v, ok := d.GetOkExists("system_noc"); ok {
		request["SystemNoc"] = v
	}
	if update {
		action := "UpdateAlertContact"
		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RpcPost("ARMS", "2019-08-08", action, nil, request, false)
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
	return ResourceAlertContactRead(d, meta)
}
func ResourceAlertContactDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DeleteAlertContact"
	var response map[string]interface{}
	var err error
	request := map[string]interface{}{
		"ContactId": d.Id(),
	}

	request["RegionId"] = client.RegionId
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		response, err = client.RpcPost("ARMS", "2019-08-08", action, nil, request, false)
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
	return nil
}

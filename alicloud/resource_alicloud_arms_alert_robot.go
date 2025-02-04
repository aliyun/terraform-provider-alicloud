package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudArmsAlertRobot() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudArmsAlertRobotCreate,
		Read:   resourceAlicloudArmsAlertRobotRead,
		Update: resourceAlicloudArmsAlertRobotUpdate,
		Delete: resourceAlicloudArmsAlertRobotDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"alert_robot_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"robot_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"wechat", "dingding", "feishu"}, false),
			},
			"robot_addr": {
				Type:     schema.TypeString,
				Required: true,
			},
			"daily_noc": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"daily_noc_time": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceAlicloudArmsAlertRobotCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateOrUpdateIMRobot"
	request := make(map[string]interface{})
	var err error
	request["RobotName"] = d.Get("alert_robot_name")
	request["Type"] = d.Get("robot_type")
	request["RobotAddress"] = d.Get("robot_addr")
	if v, ok := d.GetOk("daily_noc"); ok {
		request["DailyNoc"] = v
		if v, ok := d.GetOk("daily_noc_time"); ok && v.(string) == "" {
			return WrapError(fmt.Errorf("attribute '%s' is required when '%s' is %v", "daily_noc_time", "daily_noc", d.Get("daily_noc")))
		}
	}
	request["DailyNocTime"] = d.Get("daily_noc_time")
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("ARMS", "2019-08-08", action, nil, request, false)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_arms_alert_robot", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["AlertRobot"].(map[string]interface{})["RobotId"]))

	return resourceAlicloudArmsAlertRobotRead(d, meta)
}
func resourceAlicloudArmsAlertRobotRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	armsService := ArmsService{client}
	object, err := armsService.DescribeArmsAlertRobot(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_arms_alert_robot armsService.DescribeArmsAlertRobot Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("alert_robot_name", object["RobotName"])
	d.Set("robot_type", object["Type"])
	d.Set("robot_addr", object["RobotAddr"])
	d.Set("daily_noc", object["DailyNoc"])
	d.Set("daily_noc_time", object["DailyNocTime"])
	return nil
}
func resourceAlicloudArmsAlertRobotUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var err error
	var response map[string]interface{}
	update := false
	request := map[string]interface{}{
		"RobotId": d.Id(),
	}
	if d.HasChange("alert_robot_name") {
		update = true
	}
	if v, ok := d.GetOk("alert_robot_name"); ok {
		request["RobotName"] = v
	}
	if d.HasChange("robot_type") {
		update = true
	}
	if v, ok := d.GetOk("robot_type"); ok {
		request["Type"] = v
	}
	if d.HasChange("robot_addr") {
		update = true
	}
	if v, ok := d.GetOk("robot_addr"); ok {
		request["RobotAddress"] = v
	}
	if d.HasChange("daily_noc") {
		update = true
	}
	if v, ok := d.GetOk("daily_noc"); ok {
		request["DailyNoc"] = v
	}
	if d.HasChange("daily_noc_time") {
		update = true
	}
	if v, ok := d.GetOk("daily_noc_time"); ok {
		if d.Get("daily_noc").(bool) && v.(string) == "" {
			return WrapError(fmt.Errorf("attribute '%s' is required when '%s' is %v", "daily_noc_time", "daily_noc", d.Get("daily_noc")))
		}
		request["DailyNocTime"] = v
	}

	if update {
		action := "CreateOrUpdateIMRobot"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("ARMS", "2019-08-08", action, nil, request, false)
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
	return resourceAlicloudArmsAlertRobotRead(d, meta)
}
func resourceAlicloudArmsAlertRobotDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DeleteIMRobot"
	var response map[string]interface{}
	var err error
	request := map[string]interface{}{
		"RobotId": d.Id(),
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("ARMS", "2019-08-08", action, nil, request, false)
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
	return nil
}

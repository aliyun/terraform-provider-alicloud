package alicloud

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudCmsAlertRobot() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudCmsAlertRobotCreate,
		Read:   resourceAliCloudCmsAlertRobotRead,
		Update: resourceAliCloudCmsAlertRobotUpdate,
		Delete: resourceAliCloudCmsAlertRobotDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"alert_robot_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"alert_robot_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"DING", "WEIXIN", "FEISHU", "SLACK", "TEAMS", "DING_COOL_APP"}, false),
			},
			"lang": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: StringInSlice([]string{"zh_CN", "en_US"}, false),
			},
			"url": {
				Type:     schema.TypeString,
				Required: true,
			},
			"workspace": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"digital_employee_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"robot_sign_key": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
		},
	}
}

func resourceAliCloudCmsAlertRobotCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	body := make(map[string]interface{})
	action := "/robot"

	body["type"] = d.Get("type")
	body["name"] = d.Get("alert_robot_name")
	body["webhook"] = d.Get("url")
	langVal := d.Get("lang").(string)
	if langVal == "" {
		langVal = "zh_CN"
	}
	body["lang"] = langVal
	if v, ok := d.GetOk("workspace"); ok {
		body["workspace"] = v
	}
	if v, ok := d.GetOk("digital_employee_name"); ok {
		body["digitalEmployeeName"] = v
	}
	if v, ok := d.GetOk("robot_sign_key"); ok {
		body["robotSignKey"] = v
	}

	var err error
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RoaPost("Cms", "2024-03-30", action, make(map[string]*string), nil, body, true)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, body)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_cms_alert_robot", action, AlibabaCloudSdkGoERROR)
	}

	// Create returns alertRobotId at the top level of the response (no "data" nesting).
	if v, ok := response["alertRobotId"].(string); ok && v != "" {
		d.SetId(v)
	} else {
		return WrapError(fmt.Errorf("failed to create alert robot: no alertRobotId in response"))
	}

	return resourceAliCloudCmsAlertRobotRead(d, meta)
}

func resourceAliCloudCmsAlertRobotRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cmsServiceV2 := CmsServiceV2{client}

	object, err := cmsServiceV2.DescribeCmsAlertRobot(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_cms_alert_robot describeCmsAlertRobot Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("alert_robot_id", object["robotId"])
	d.Set("alert_robot_name", object["name"])
	d.Set("type", object["type"])
	d.Set("lang", object["lang"])
	// "url" (webhook) is write-only: the List API does not return it, so we
	// preserve existing state to avoid spurious diffs instead of overwriting.
	// "robot_sign_key" (robotSignKey) is returned by the List API when a sign
	// key was configured; when the robot has no sign key the API omits the
	// field. Skipping d.Set in that case preserves the planned null value and
	// avoids the "inconsistent values for sensitive attribute" warning that
	// would otherwise mark the apply dirty and drop data-source state.
	if v, ok := object["robotSignKey"]; ok && v != nil {
		d.Set("robot_sign_key", v)
	}
	// "workspace" and "digital_employee_name" may also be absent for robot
	// types that do not use them; skip d.Set when nil to keep state stable.
	if v, ok := object["workspace"]; ok && v != nil {
		d.Set("workspace", v)
	}
	if v, ok := object["digitalEmployeeName"]; ok && v != nil {
		d.Set("digital_employee_name", v)
	}

	return nil
}

func resourceAliCloudCmsAlertRobotUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := fmt.Sprintf("/robot/%s", d.Id())
	body := make(map[string]interface{})
	update := false

	if d.HasChange("alert_robot_name") {
		update = true
	}
	if v, ok := d.GetOk("alert_robot_name"); ok || d.HasChange("alert_robot_name") {
		body["name"] = v
	}
	if d.HasChange("lang") {
		update = true
	}
	body["lang"] = d.Get("lang")
	if d.HasChange("url") {
		update = true
	}
	if v, ok := d.GetOk("url"); ok || d.HasChange("url") {
		body["webhook"] = v
	}
	if d.HasChange("digital_employee_name") {
		update = true
	}
	if v, ok := d.GetOk("digital_employee_name"); ok || d.HasChange("digital_employee_name") {
		body["digitalEmployeeName"] = v
	}
	if d.HasChange("robot_sign_key") {
		update = true
	}
	if v, ok := d.GetOk("robot_sign_key"); ok || d.HasChange("robot_sign_key") {
		body["robotSignKey"] = v
	}

	if update {
		var err error
		var response map[string]interface{}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RoaPatch("Cms", "2024-03-30", action, make(map[string]*string), nil, body, true)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, body)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
	}

	return resourceAliCloudCmsAlertRobotRead(d, meta)
}

func resourceAliCloudCmsAlertRobotDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "/robots"
	query := make(map[string]*string)

	robotIds, _ := json.Marshal([]string{d.Id()})
	query["robotIds"] = StringPointer(string(robotIds))
	if v, ok := d.GetOk("type"); ok {
		query["type"] = StringPointer(v.(string))
	}

	var err error
	var response map[string]interface{}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RoaDelete("Cms", "2024-03-30", action, query, nil, nil, true)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, query)
	if err != nil {
		if IsExpectedErrors(err, []string{"ResourceNotFound", "InvalidParameter"}) || NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	return nil
}

func (s *CmsServiceV2) DescribeCmsAlertRobot(id string) (map[string]interface{}, error) {
	action := "/robots"
	query := make(map[string]*string)
	robotIds, _ := json.Marshal([]string{id})
	query["robotIds"] = StringPointer(string(robotIds))
	query["pageSize"] = StringPointer("100")
	query["pageNumber"] = StringPointer("1")

	var response map[string]interface{}
	var err error
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = s.client.RoaGet("Cms", "2024-03-30", action, query, nil, nil)
		if err != nil {
			if NeedRetry(err) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, query)
	if err != nil {
		return nil, err
	}

	robots, ok := response["robots"].([]interface{})
	if !ok || len(robots) == 0 {
		return nil, WrapErrorf(NotFoundErr("AlertRobot", id), NotFoundMsg, response)
	}
	for _, r := range robots {
		robot, ok := r.(map[string]interface{})
		if !ok {
			continue
		}
		if robotId, ok := robot["robotId"].(string); ok && robotId == id {
			return robot, nil
		}
	}
	return nil, WrapErrorf(NotFoundErr("AlertRobot", id), NotFoundMsg, response)
}

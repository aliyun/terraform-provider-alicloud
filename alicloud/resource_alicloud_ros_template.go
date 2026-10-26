package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAlicloudRosTemplate() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudRosTemplateCreate,
		Read:   resourceAlicloudRosTemplateRead,
		Update: resourceAlicloudRosTemplateUpdate,
		Delete: resourceAlicloudRosTemplateDelete,
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
		},
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"rotate_strategy": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "None",
				ValidateFunc: validation.StringInSlice([]string{"None", "DeleteOldestNonSharedVersionWhenLimitExceeded"}, false),
			},
			"tags": tagsSchema(),
			"template_body": {
				Type:         schema.TypeString,
				Optional:     true,
				ExactlyOneOf: []string{"template_body", "template_url"},
				ValidateFunc: validation.StringIsJSON,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					equal, _ := compareJsonTemplateAreEquivalent(old, new)
					return equal
				},
			},
			"template_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"template_url": {
				Type:         schema.TypeString,
				Optional:     true,
				ExactlyOneOf: []string{"template_body", "template_url"},
			},
			"validation_options": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringInSlice([]string{"BasicValidation"}, false),
				},
			},
		},
	}
}

func resourceAlicloudRosTemplateCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateTemplate"
	request := make(map[string]interface{})
	var err error
	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}

	if v, ok := d.GetOk("template_body"); ok {
		request["TemplateBody"] = v
	}

	request["TemplateName"] = d.Get("template_name")
	request["RegionId"] = client.RegionId
	if v, ok := d.GetOk("template_url"); ok {
		request["TemplateURL"] = v
	}

	if v, ok := d.GetOk("resource_group_id"); ok {
		request["ResourceGroupId"] = v
	}

	if v, ok := d.GetOk("validation_options"); ok {
		for i, option := range v.([]interface{}) {
			request[fmt.Sprintf("ValidationOptions.%d", i+1)] = option.(string)
		}
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		response, err = client.RpcPost("ROS", "2019-09-10", action, nil, request, false)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_ros_template", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["TemplateId"]))

	return resourceAlicloudRosTemplateUpdate(d, meta)
}
func resourceAlicloudRosTemplateRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	rosService := RosService{client}
	object, err := rosService.DescribeRosTemplate(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_ros_template rosService.DescribeRosTemplate Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("description", object["Description"])
	d.Set("resource_group_id", object["ResourceGroupId"])
	if _, ok := d.GetOk("template_url"); ok {
		// TemplateURL content is not tracked in state; keep the body empty so a
		// URL-only configuration does not drift against the body returned by GetTemplate.
		d.Set("template_body", "")
	} else {
		d.Set("template_body", object["TemplateBody"])
	}
	d.Set("template_name", object["TemplateName"])

	listTagResourcesObject, err := rosService.ListTagResources(d.Id(), "template")
	if err != nil {
		return WrapError(err)
	}
	d.Set("tags", tagsToMap(listTagResourcesObject))
	return nil
}
func resourceAlicloudRosTemplateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	rosService := RosService{client}
	var err error
	var response map[string]interface{}
	d.Partial(true)

	if d.HasChange("tags") {
		if err := rosService.SetResourceTags(d, "template"); err != nil {
			return WrapError(err)
		}
	}
	update := false
	request := map[string]interface{}{
		"TemplateId": d.Id(),
		"RegionId":   client.RegionId,
	}
	if !d.IsNewResource() && d.HasChange("description") {
		update = true
		request["Description"] = d.Get("description")
	}
	if !d.IsNewResource() && d.HasChange("template_body") {
		// TemplateBody and TemplateURL are mutually exclusive; never send an empty body.
		if v, ok := d.GetOk("template_body"); ok {
			update = true
			request["TemplateBody"] = v
		}
	}
	if !d.IsNewResource() && d.HasChange("template_name") {
		update = true
		request["TemplateName"] = d.Get("template_name")
	}
	if d.HasChange("rotate_strategy") && (!d.IsNewResource() || d.Get("rotate_strategy") != "None") {
		update = true
		request["RotateStrategy"] = d.Get("rotate_strategy")
	}
	if !d.IsNewResource() && d.HasChange("template_url") {
		// TemplateURL and TemplateBody are mutually exclusive; never send an empty URL.
		if v, ok := d.GetOk("template_url"); ok {
			update = true
			request["TemplateURL"] = v
		}
	}
	if update {
		// ValidationOptions only affects the content carried by the same call.
		if v, ok := d.GetOk("validation_options"); ok {
			for i, option := range v.([]interface{}) {
				request[fmt.Sprintf("ValidationOptions.%d", i+1)] = option.(string)
			}
		}
		action := "UpdateTemplate"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RpcPost("ROS", "2019-09-10", action, nil, request, false)
			if err != nil {
				if IsExpectedErrors(err, []string{"TemplateBeingProcessed"}) || NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			addDebug(action, response, request)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
	}
	if !d.IsNewResource() && d.HasChange("resource_group_id") {
		moveAction := "MoveResourceGroup"
		moveRequest := map[string]interface{}{
			"RegionId":           client.RegionId,
			"ResourceId":         d.Id(),
			"ResourceType":       "template",
			"NewResourceGroupId": d.Get("resource_group_id"),
		}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RpcPost("ROS", "2019-09-10", moveAction, nil, moveRequest, false)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			addDebug(moveAction, response, moveRequest)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), moveAction, AlibabaCloudSdkGoERROR)
		}
	}
	d.Partial(false)
	return resourceAlicloudRosTemplateRead(d, meta)
}
func resourceAlicloudRosTemplateDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DeleteTemplate"
	var response map[string]interface{}
	var err error
	request := map[string]interface{}{
		"TemplateId": d.Id(),
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		response, err = client.RpcPost("ROS", "2019-09-10", action, nil, request, false)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"ChangeSetNotFound", "StackNotFound", "TemplateNotFound"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	return nil
}

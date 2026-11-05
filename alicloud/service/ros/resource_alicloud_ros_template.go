package ros

import (
	"fmt"
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/tags"
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
			"tags": tags.TagsSchema(),
			"template_body": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsJSON,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					equal, _ := helper.CompareJsonTemplateAreEquivalent(old, new)
					return equal
				},
			},
			"template_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"template_url": {
				Type:     schema.TypeString,
				Optional: true,
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
	if v, ok := d.GetOk("template_url"); ok {
		request["TemplateURL"] = v
	}

	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		response, err = client.RpcPost("ROS", "2019-09-10", action, nil, request, false)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		logs.AddDebug(action, response, request)
		return nil
	})
	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_ros_template", action, errs.AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["TemplateId"]))

	return resourceAlicloudRosTemplateUpdate(d, meta)
}
func resourceAlicloudRosTemplateRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	rosService := RosService{client}
	object, err := rosService.DescribeRosTemplate(d.Id())
	if err != nil {
		if errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_ros_template rosService.DescribeRosTemplate Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}
	d.Set("description", object["Description"])
	d.Set("template_body", object["TemplateBody"])
	d.Set("template_name", object["TemplateName"])

	listTagResourcesObject, err := rosService.ListTagResources(d.Id(), "template")
	if err != nil {
		return errs.WrapError(err)
	}
	d.Set("tags", tags.TagsToMap(listTagResourcesObject))
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
			return errs.WrapError(err)
		}
	}
	update := false
	request := map[string]interface{}{
		"TemplateId": d.Id(),
	}
	if !d.IsNewResource() && d.HasChange("description") {
		update = true
		request["Description"] = d.Get("description")
	}
	if !d.IsNewResource() && d.HasChange("template_body") {
		update = true
		request["TemplateBody"] = d.Get("template_body")
	}
	if !d.IsNewResource() && d.HasChange("template_name") {
		update = true
		request["TemplateName"] = d.Get("template_name")
	}
	if update {
		if _, ok := d.GetOk("template_url"); ok {
			request["TemplateURL"] = d.Get("template_url")
		}
		action := "UpdateTemplate"
		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RpcPost("ROS", "2019-09-10", action, nil, request, false)
			if err != nil {
				if helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			logs.AddDebug(action, response, request)
			return nil
		})
		if err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
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

	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		response, err = client.RpcPost("ROS", "2019-09-10", action, nil, request, false)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		logs.AddDebug(action, response, request)
		return nil
	})
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"ChangeSetNotFound", "StackNotFound", "TemplateNotFound"}) {
			return nil
		}
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
	}
	return nil
}

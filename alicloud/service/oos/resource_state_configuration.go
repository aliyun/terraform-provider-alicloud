package oos

import (
	"fmt"
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/convert"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/tags"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourceStateConfiguration() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudOosStateConfigurationCreate,
		Read:   resourceAlicloudOosStateConfigurationRead,
		Update: resourceAlicloudOosStateConfigurationUpdate,
		Delete: resourceAlicloudOosStateConfigurationDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"configure_mode": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"ApplyAndAutoCorrect", "ApplyAndMonitor", "ApplyOnly"}, false),
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"parameters": {
				Type:     schema.TypeString,
				Optional: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					equal, _ := helper.CompareJsonTemplateAreEquivalent(old, new)
					return equal
				},
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"schedule_expression": {
				Type:     schema.TypeString,
				Required: true,
			},
			"schedule_type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"rate"}, false),
			},
			"tags": tags.TagsSchema(),
			"targets": {
				Type:     schema.TypeString,
				Required: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					equal, _ := helper.CompareJsonTemplateAreEquivalent(old, new)
					return equal
				},
			},
			"template_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"ACS-ECS-InventoryDataCollection"}, false),
			},
			"template_version": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAlicloudOosStateConfigurationCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateStateConfiguration"
	request := make(map[string]interface{})
	var err error
	if v, ok := d.GetOk("configure_mode"); ok {
		request["ConfigureMode"] = v
	}
	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}
	if v, ok := d.GetOk("parameters"); ok {
		request["Parameters"] = v
	}
	request["RegionId"] = client.RegionId
	if v, ok := d.GetOk("resource_group_id"); ok {
		request["ResourceGroupId"] = v
	}
	request["ScheduleExpression"] = d.Get("schedule_expression")
	request["ScheduleType"] = d.Get("schedule_type")
	if v, ok := d.GetOk("tags"); ok {
		respJson, err := convert.ConvertMaptoJsonString(v.(map[string]interface{}))
		if err != nil {
			return errs.WrapError(err)
		}
		request["Tags"] = respJson
	}
	request["Targets"] = d.Get("targets")
	request["TemplateName"] = d.Get("template_name")
	if v, ok := d.GetOk("template_version"); ok {
		request["TemplateVersion"] = v
	}
	request["ClientToken"] = helper.BuildClientToken("CreateStateConfiguration")
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		response, err = client.RpcPost("oos", "2019-06-01", action, nil, request, true)
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
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_oos_state_configuration", action, errs.AlibabaCloudSdkGoERROR)
	}
	responseStateConfiguration := response["StateConfiguration"].(map[string]interface{})
	d.SetId(fmt.Sprint(responseStateConfiguration["StateConfigurationId"]))

	return resourceAlicloudOosStateConfigurationRead(d, meta)
}
func resourceAlicloudOosStateConfigurationRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	oosService := OosService{client}
	object, err := oosService.DescribeOosStateConfiguration(d.Id())
	if err != nil {
		if errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_oos_state_configuration oosService.DescribeOosStateConfiguration Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}
	d.Set("configure_mode", object["ConfigureMode"])
	d.Set("description", object["Description"])
	d.Set("parameters", object["Parameters"])
	d.Set("resource_group_id", object["ResourceGroupId"])
	d.Set("schedule_expression", object["ScheduleExpression"])
	d.Set("schedule_type", object["ScheduleType"])
	d.Set("tags", tags.TagsToMap(object["Tags"]))
	d.Set("targets", object["Targets"])
	d.Set("template_name", object["TemplateName"])
	d.Set("template_version", object["TemplateVersion"])
	return nil
}
func resourceAlicloudOosStateConfigurationUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var err error
	var response map[string]interface{}
	update := false
	request := map[string]interface{}{
		"StateConfigurationId": d.Id(),
	}
	if d.HasChange("configure_mode") {
		update = true
		if v, ok := d.GetOk("configure_mode"); ok {
			request["ConfigureMode"] = v
		}
	}
	if d.HasChange("description") {
		update = true
		if v, ok := d.GetOk("description"); ok {
			request["Description"] = v
		}
	}
	if d.HasChange("parameters") {
		update = true
		if v, ok := d.GetOk("parameters"); ok {
			request["Parameters"] = v
		}
	}
	request["RegionId"] = client.RegionId
	if d.HasChange("resource_group_id") {
		update = true
		if v, ok := d.GetOk("resource_group_id"); ok {
			request["ResourceGroupId"] = v
		}
	}
	if d.HasChange("schedule_expression") {
		update = true
		request["ScheduleExpression"] = d.Get("schedule_expression")
	}
	if d.HasChange("schedule_type") {
		update = true
		request["ScheduleType"] = d.Get("schedule_type")
	}
	if d.HasChange("tags") {
		update = true
		if v, ok := d.GetOk("tags"); ok {
			respJson, err := convert.ConvertMaptoJsonString(v.(map[string]interface{}))
			if err != nil {
				return errs.WrapError(err)
			}
			request["Tags"] = respJson
		}
	}
	if d.HasChange("targets") {
		update = true
		request["Targets"] = d.Get("targets")
	}
	if update {
		action := "UpdateStateConfiguration"
		request["ClientToken"] = helper.BuildClientToken("UpdateStateConfiguration")
		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RpcPost("oos", "2019-06-01", action, nil, request, true)
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
	return resourceAlicloudOosStateConfigurationRead(d, meta)
}
func resourceAlicloudOosStateConfigurationDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DeleteStateConfigurations"
	var response map[string]interface{}
	var err error
	request := map[string]interface{}{
		"StateConfigurationIds": convert.ConvertListToJsonString([]interface{}{d.Id()}),
	}

	request["RegionId"] = client.RegionId
	request["ClientToken"] = helper.BuildClientToken("DeleteStateConfigurations")
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		response, err = client.RpcPost("oos", "2019-06-01", action, nil, request, true)
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

package imp

import (
	"fmt"
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/convert"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourceAppTemplate() *schema.Resource {
	return &schema.Resource{
		Create: resourceAppTemplateCreate,
		Read:   resourceAppTemplateRead,
		Update: resourceAppTemplateUpdate,
		Delete: resourceAppTemplateDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"app_template_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"component_list": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringInSlice([]string{"component.live", "component.liveRecord", "component.liveBeauty", "component.rtc", "component.rtcRecord", "component.im", "component.whiteboard", "component.liveSecurity", "component.chatSecurity"}, false),
				},
				ForceNew: true,
			},
			"config_list": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"value": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
			"integration_mode": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"paasSDK", "standardRoom"}, false),
			},
			"scene": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"business", "classroom"}, false),
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAppTemplateCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	request := make(map[string]interface{})
	var err error
	request["AppTemplateName"] = d.Get("app_template_name")
	request["ComponentList"] = convert.ConvertListToJsonString(d.Get("component_list").([]interface{}))
	if v, ok := d.GetOk("integration_mode"); ok {
		request["IntegrationMode"] = v
	}
	if v, ok := d.GetOk("scene"); ok {
		request["Scene"] = v
	}

	var response map[string]interface{}
	action := "CreateAppTemplate"
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		resp, err := client.RpcPost("imp", "2021-06-30", action, nil, request, false)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		response = resp
		logs.AddDebug(action, resp, request)
		return nil
	})
	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_imp_app_template", action, errs.AlibabaCloudSdkGoERROR)
	}
	responseResult := response["Result"].(map[string]interface{})
	d.SetId(fmt.Sprint(responseResult["AppTemplateId"]))

	return resourceAppTemplateRead(d, meta)
}
func resourceAppTemplateRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	impService := NewService(client)
	object, err := impService.DescribeAppTemplate(d.Id())
	if err != nil {
		if errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_imp_app_template impService.DescribeAppTemplate Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}
	d.Set("app_template_name", object["AppTemplateName"])
	if componentList, ok := object["ComponentList"]; ok && componentList != nil {
		d.Set("component_list", componentList)
	}
	if configList, ok := object["ConfigList"]; ok && configList != nil {
		configListMaps := make([]map[string]interface{}, 0)
		for _, configListItem := range configList.([]interface{}) {
			if configListItemMap, ok := configListItem.(map[string]interface{}); ok {
				configListMap := make(map[string]interface{}, 0)
				configListMap["key"] = configListItemMap["Key"]
				configListMap["value"] = configListItemMap["Value"]
				configListMaps = append(configListMaps, configListMap)
			}
		}
		d.Set("config_list", configListMaps)
	}
	d.Set("integration_mode", object["IntegrationMode"])
	d.Set("scene", object["Scene"])
	d.Set("status", object["Status"])
	return nil
}
func resourceAppTemplateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	d.Partial(true)
	var err error

	update := false
	request := map[string]interface{}{
		"AppTemplateId": d.Id(),
	}
	if d.HasChange("app_template_name") {
		update = true
		request["AppTemplateName"] = d.Get("app_template_name")
	}
	if update {
		action := "UpdateAppTemplate"
		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			resp, err := client.RpcPost("imp", "2021-06-30", action, nil, request, false)
			if err != nil {
				if helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			logs.AddDebug(action, resp, request)
			return nil
		})
		if err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
		}
	}
	update = false
	updateAppTemplateConfigReq := map[string]interface{}{
		"AppTemplateId": d.Id(),
	}
	if d.HasChanges("config_list") {
		update = true
		if v, ok := d.GetOk("config_list"); ok {
			configListMaps := make([]map[string]interface{}, 0)
			for _, configList := range v.(*schema.Set).List() {
				configListArg := configList.(map[string]interface{})
				configListMap := map[string]interface{}{
					"Key":   configListArg["key"],
					"Value": configListArg["value"],
				}
				configListMaps = append(configListMaps, configListMap)
			}
			if v, err := convert.ConvertListMapToJsonString(configListMaps); err != nil {
				return errs.WrapError(err)
			} else {
				updateAppTemplateConfigReq["ConfigList"] = v
			}
		}
	}
	if update {
		action := "UpdateAppTemplateConfig"
		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			resp, err := client.RpcPost("imp", "2021-06-30", action, nil, updateAppTemplateConfigReq, false)
			if err != nil {
				if helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			logs.AddDebug(action, resp, updateAppTemplateConfigReq)
			return nil
		})
		if err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
		}
	}
	d.Partial(false)
	return resourceAppTemplateRead(d, meta)
}
func resourceAppTemplateDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var err error
	request := map[string]interface{}{
		"AppTemplateId": d.Id(),
	}

	action := "DeleteAppTemplate"
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		resp, err := client.RpcPost("imp", "2021-06-30", action, nil, request, false)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		logs.AddDebug(action, resp, request)
		return nil
	})
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"InvalidAppTemplateId.App.NotFound"}) || errs.NotFoundError(err) {
			return nil
		}
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
	}
	return nil
}

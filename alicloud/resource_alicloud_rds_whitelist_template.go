package alicloud

import (
	"fmt"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudRdsWhitelistTemplate() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudWhitelistTemplateCreate,
		Read:   resourceAliCloudWhitelistTemplateRead,
		Update: resourceAliCloudWhitelistTemplateUpdate,
		Delete: resourceAliCloudWhitelistTemplateDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"ip_white_list": {
				Type:     schema.TypeString,
				Required: true,
			},
			"template_name": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceAliCloudWhitelistTemplateCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	rdsService := RdsService{client}
	action := "ModifyWhitelistTemplate"
	request := map[string]interface{}{
		"RegionId":     client.RegionId,
		"IpWhitelist":  Trim(d.Get("ip_white_list").(string)),
		"TemplateName": Trim(d.Get("template_name").(string)),
	}
	if err := resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *resource.RetryError {
		response, err := client.RpcPost("Rds", "2014-08-15", action, nil, request, false)
		if err != nil {
			if NeedRetry(err) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	}); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	//临时处理方案 需要等待接口优化:成功返回TemplateId
	template, templateErr := rdsService.DescribeAllWhitelistTemplate(d.Get("template_name").(string))
	if templateErr != nil {
		if NotFoundError(templateErr) {
			d.SetId("")
			return nil
		}
		return WrapError(templateErr)
	}
	d.SetId(fmt.Sprintf("%s", template["TemplateId"]))

	return resourceAliCloudWhitelistTemplateRead(d, meta)
}

func resourceAliCloudWhitelistTemplateRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	rdsService := RdsService{client}
	template, templateErr := rdsService.DescribeWhitelistTemplate(d.Id())
	if templateErr != nil {
		if NotFoundError(templateErr) {
			d.SetId("")
			return nil
		}
		return WrapError(templateErr)
	}
	d.Set("template_name", template["TemplateName"])
	d.Set("ip_white_list", template["Ips"])
	return nil
}

func resourceAliCloudWhitelistTemplateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	if d.HasChanges("ip_white_list", "template_name") {
		action := "ModifyWhitelistTemplate"
		request := map[string]interface{}{
			"RegionId": client.RegionId,
		}
		request["TemplateId"] = d.Id()
		Update := false
		if v, ok := d.GetOk("ip_white_list"); ok {
			request["IpWhitelist"] = v
			Update = true
		}
		if v, ok := d.GetOk("template_name"); ok {
			request["TemplateName"] = v
			Update = true
		}
		if Update {
			if err := resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
				response, err := client.RpcPost("Rds", "2014-08-15", action, nil, request, false)
				if err != nil {
					if NeedRetry(err) {
						return resource.RetryableError(err)
					}
					return resource.NonRetryableError(err)
				}
				addDebug(action, response, request)
				return nil
			}); err != nil {
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
			}
		}
	}
	return resourceAliCloudWhitelistTemplateRead(d, meta)
}
func resourceAliCloudWhitelistTemplateDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	if d.Id() == "" {
		return nil
	}
	action := "ModifyWhitelistTemplate"
	request := map[string]interface{}{
		"RegionId":     client.RegionId,
		"IpWhitelist":  "",
		"TemplateName": d.Get("template_name"),
		"TemplateId":   d.Id(),
	}
	if err := resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *resource.RetryError {
		response, err := client.RpcPost("Rds", "2014-08-15", action, nil, request, false)
		if err != nil {
			if IsExpectedErrors(err, OperationDeniedDBStatus) || NeedRetry(err) {
				return resource.RetryableError(err)
			}
			if NotFoundError(err) {
				return nil
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	}); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	return nil
}

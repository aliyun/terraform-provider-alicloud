package sls

import (
	"fmt"

	sls "github.com/aliyun/aliyun-log-go-sdk"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlicloudLogtailAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudLogtailAttachmentCreate,
		Read:   resourceAlicloudLogtailAttachmentRead,
		Delete: resourceAlicloudLogtailAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"project": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"logtail_config_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"machine_group_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAlicloudLogtailAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	project := d.Get("project").(string)
	config_name := d.Get("logtail_config_name").(string)
	group_name := d.Get("machine_group_name").(string)
	var requestInfo *sls.Client
	raw, err := client.WithLogClient(func(slsClient *sls.Client) (interface{}, error) {
		requestInfo = slsClient
		return nil, slsClient.ApplyConfigToMachineGroup(project, config_name, group_name)
	})
	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_logtail_attachment", "ApplyConfigToMachineGroup", errs.AliyunLogGoSdkERROR)
	}
	if logs.DebugOn() {
		logs.AddDebug("ApplyConfigToMachineGroup", raw, requestInfo, map[string]string{
			"project":   project,
			"confName":  config_name,
			"groupName": group_name,
		})
	}
	d.SetId(fmt.Sprintf("%s%s%s%s%s", project, COLON_SEPARATED, config_name, COLON_SEPARATED, group_name))
	return resourceAlicloudLogtailAttachmentRead(d, meta)
}

func resourceAlicloudLogtailAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	logService := Service{client}
	parts, err := helper.ParseResourceId(d.Id(), 3)
	if err != nil {
		return errs.WrapError(err)
	}
	object, err := logService.DescribeLogtailAttachment(d.Id())
	if err != nil {
		if errs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}

	d.Set("project", parts[0])
	d.Set("logtail_config_name", parts[1])
	d.Set("machine_group_name", object)

	return nil
}

func resourceAlicloudLogtailAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	logService := Service{client}
	parts, err := helper.ParseResourceId(d.Id(), 3)
	if err != nil {
		return errs.WrapError(err)
	}
	var requestInfo *sls.Client
	raw, err := client.WithLogClient(func(slsClient *sls.Client) (interface{}, error) {
		requestInfo = slsClient
		return nil, slsClient.RemoveConfigFromMachineGroup(parts[0], parts[1], parts[2])
	})
	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), "RemoveConfigFromMachineGroup", errs.AliyunLogGoSdkERROR)
	}
	if logs.DebugOn() {
		logs.AddDebug("RemoveConfigFromMachineGroup", raw, requestInfo, map[string]string{
			"project":   parts[0],
			"confName":  parts[1],
			"groupName": parts[2],
		})
	}
	return errs.WrapError(logService.WaitForLogtailAttachment(d.Id(), Deleted, DefaultTimeout))

}

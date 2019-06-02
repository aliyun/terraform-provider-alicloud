package alicloud

import (
	"fmt"
	"strings"
	"time"

	sls "github.com/aliyun/aliyun-log-go-sdk"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
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
			"project": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"logtail_config_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"machine_group_name": &schema.Schema{
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
	_, err := client.WithLogClient(func(slsClient *sls.Client) (interface{}, error) {
		return nil, slsClient.ApplyConfigToMachineGroup(project, config_name, group_name)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "logtail_attachment", "ApplyConfigToMachineGroup", AliyunLogGoSdkERROR)
	}
	d.SetId(fmt.Sprintf("%s%s%s%s%s", project, COLON_SEPARATED, config_name, COLON_SEPARATED, group_name))
	return resourceAlicloudLogtailAttachmentRead(d, meta)
}

func resourceAlicloudLogtailAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	logService := LogService{client}
	split := strings.Split(d.Id(), COLON_SEPARATED)
	groupNames, err := logService.DescribeLogtailAttachment(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	var groupName string
	for _, group_name := range groupNames {
		if group_name == d.Get("machine_group_name").(string) {
			groupName = group_name
		}
	}
	d.Set("project", split[0])
	d.Set("logtail_config_name", split[1])
	d.Set("machine_group_name", groupName)

	return nil
}

func resourceAlicloudLogtailAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	logService := LogService{client}
	split := strings.Split(d.Id(), COLON_SEPARATED)
	return resource.Retry(3*time.Minute, func() *resource.RetryError {
		_, err := client.WithLogClient(func(slsClient *sls.Client) (interface{}, error) {
			return nil, slsClient.RemoveConfigFromMachineGroup(split[0], split[1], split[2])
		})
		if err != nil {
			return resource.NonRetryableError(WrapErrorf(err, DefaultErrorMsg, d.Id(), "RemoveConfigFromMachineGroup", AliyunLogGoSdkERROR))
		}
		if _, err1 := logService.DescribeLogtailAttachment(d.Id()); err1 != nil {
			if NotFoundError(err1) {
				return nil
			}
			return resource.NonRetryableError(WrapError(err1))
		}
		return resource.RetryableError(WrapErrorf(err, DeleteTimeoutMsg, d.Id(), "RemoveConfigFromMachineGroup", ProviderERROR))
	})
}

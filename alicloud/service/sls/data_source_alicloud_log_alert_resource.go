package sls

import (
	"fmt"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceAlicloudLogAlertResource() *schema.Resource {
	return &schema.Resource{
		Read:               dataSourceAlicloudLogAlertResourceRead,
		DeprecationMessage: "This data source has been deprecated since v1.219.0 and will be removed in the future. Please use 'alicloud_log_alert_resource' resource instead.",
		Schema: map[string]*schema.Schema{
			"type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"user", "project"}, true),
			},
			"lang": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"project": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}

}
func dataSourceAlicloudLogAlertResourceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	resourceType := d.Get("type").(string)
	lang := d.Get("lang").(string)
	project := d.Get("project").(string)
	logService := Service{client}
	if err := retry.Retry(2*time.Minute, func() *retry.RetryError {
		err := logService.InitLogAlertResource(resourceType, lang, project)
		if err != nil {
			if errs.IsExpectedErrors(err, []string{errs.LogClientTimeout}) {
				time.Sleep(5 * time.Second)
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	}); err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_log_alert_resource", "CreateAlertResource", errs.AliyunLogGoSdkERROR)
	}
	d.SetId(fmt.Sprintf("alert_resource_%s", resourceType))
	return nil
}

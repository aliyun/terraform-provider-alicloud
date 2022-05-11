package alicloud

import (
	"errors"
	"fmt"
	"time"

	slsPop "github.com/aliyun/alibaba-cloud-sdk-go/services/sls"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAlicloudLogAlertResource() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudLogAlertResourceRead,
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
	lang, _ := d.Get("lang").(string)
	project, _ := d.Get("project").(string)
	if err := resource.Retry(2*time.Minute, func() *resource.RetryError {
		_, err := client.WithLogPopClient(func(slsPopClient *slsPop.Client) (interface{}, error) {
			switch resourceType {
			case "user":
				request := slsPop.CreateInitUserAlertResourceRequest()
				request.Region = client.RegionId
				request.App = "none"
				request.Language = lang
				return slsPopClient.InitUserAlertResource(request)
			case "project":
				request := slsPop.CreateInitProjectAlertResourceRequest()
				request.Region = client.RegionId
				request.ProjectName = project
				return slsPopClient.InitProjectAlertResource(request)
			default:
				return nil, WrapErrorf(errors.New("type error"), DefaultErrorMsg, "alicloud_log_alert_resource", "CreateAlertResource", AliyunLogGoSdkERROR)
			}
		})
		if err != nil {
			if IsExpectedErrors(err, []string{LogClientTimeout}) {
				time.Sleep(5 * time.Second)
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	}); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_log_alert_resource", "CreateAlertResource", AliyunLogGoSdkERROR)
	}
	d.SetId(fmt.Sprintf("alert_resource_%s", resourceType))
	return nil
}

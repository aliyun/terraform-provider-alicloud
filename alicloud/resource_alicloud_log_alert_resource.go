package alicloud

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	slsPop "github.com/aliyun/alibaba-cloud-sdk-go/services/sls"
	sls "github.com/aliyun/aliyun-log-go-sdk"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

type AlertGlobalConfig struct {
	ConfigId     string `json:"config_id"`
	ConfigName   string `json:"config_name"`
	ConfigDetail struct {
		AlertCenterLog struct {
			Region string `json:"region"`
		} `json:"alert_center_log"`
	} `json:"config_detail"`
}

func resourceAlicloudLogAlertResource() *schema.Resource {
	return &schema.Resource{
		Create: resourcelicloudLogAlertResourceCreate,
		Read:   resourcelicloudLogAlertResourceRead,
		Delete: resourcelicloudLogAlertResourceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"user", "project"}, true),
				ForceNew:     true,
			},
			"lang": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"project": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

func resourcelicloudLogAlertResourceCreate(d *schema.ResourceData, meta interface{}) error {
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
				_, err := client.WithLogClient(func(slsClient *sls.Client) (interface{}, error) {
					return slsClient.GetLogStore(project, "internal-alert-history")
				})
				if err != nil {
					if IsExpectedErrors(err, []string{"LogStoreNotExist"}) {
						request := slsPop.CreateAnalyzeProductLogRequest()
						request.CloudProduct = "sls.alert"
						request.Project = project
						request.Logstore = "internal-alert-history"
						request.Overwrite = "true"
						request.Region = client.RegionId
						return slsPopClient.AnalyzeProductLog(request)
					}
					return nil, err
				}
				return nil, nil
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
	if resourceType == "user" {
		d.SetId(fmt.Sprintf("alert_resource%s%s%s%s", COLON_SEPARATED, resourceType, COLON_SEPARATED, lang))
	} else {
		d.SetId(fmt.Sprintf("alert_resource%s%s%s%s", COLON_SEPARATED, resourceType, COLON_SEPARATED, project))
	}
	return nil
}

func resourcelicloudLogAlertResourceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	parts, err := ParseResourceId(d.Id(), 3)
	if err != nil {
		return WrapError(err)
	}
	resourceType := parts[1]
	if err := resource.Retry(2*time.Minute, func() *resource.RetryError {
		_, err := client.WithLogPopClient(func(slsPopClient *slsPop.Client) (interface{}, error) {
			switch resourceType {
			case "user":
				_, err := client.WithLogClient(func(slsClient *sls.Client) (interface{}, error) {
					record, err := slsClient.GetResourceRecord("sls.alert.global_config", "default_config")
					if err != nil {
						return nil, err
					}
					var alertGlobalConfig AlertGlobalConfig
					err = json.Unmarshal([]byte(record.Value), &alertGlobalConfig)
					if err != nil {
						return nil, err
					}
					region := alertGlobalConfig.ConfigDetail.AlertCenterLog.Region
					accountId, err := client.AccountId()
					if err != nil {
						return nil, err
					}
					projectName := fmt.Sprintf("sls-alert-%s-%s", accountId, region)
					endpoint := slsClient.Endpoint
					slsClient.Endpoint = strings.Replace(endpoint, client.RegionId, region, 1)
					_, err = slsClient.GetProject(projectName)
					if err != nil {
						slsClient.Endpoint = endpoint
						return nil, err
					}
					_, err = slsClient.GetLogStore(projectName, "internal-alert-center-log")
					slsClient.Endpoint = endpoint
					if err != nil {
						return nil, err
					}
					return nil, nil
				})
				if err != nil {
					if IsExpectedErrors(err, []string{"ProjectNotExist"}) || IsExpectedErrors(err, []string{"LogStoreNotExist"}) {
						d.SetId("")
						return nil, nil
					}
					return nil, err
				}
				lang := parts[2]
				d.Set("type", resourceType)
				d.Set("project", nil)
				d.Set("lang", lang)
				return nil, nil
			case "project":
				project := parts[2]
				_, err := client.WithLogClient(func(slsClient *sls.Client) (interface{}, error) {
					return slsClient.GetLogStore(project, "internal-alert-history")
				})
				if err != nil {
					if IsExpectedErrors(err, []string{"LogStoreNotExist"}) {
						d.SetId("")
						return nil, nil
					}
					return nil, err
				}
				d.Set("type", resourceType)
				d.Set("project", project)
				d.Set("lang", nil)
				return nil, nil
			default:
				return nil, WrapErrorf(errors.New("type error"), DefaultErrorMsg, "alicloud_log_alert_resource", "ReadAlertResource", AliyunLogGoSdkERROR)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_log_alert_resource", "ReadAlertResource", AliyunLogGoSdkERROR)
	}
	return nil
}

func resourcelicloudLogAlertResourceDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}

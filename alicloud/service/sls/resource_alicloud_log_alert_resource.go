package sls

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	sls "github.com/aliyun/aliyun-log-go-sdk"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
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
	if resourceType == "user" {
		d.SetId(fmt.Sprintf("alert_resource%s%s%s%s", COLON_SEPARATED, resourceType, COLON_SEPARATED, lang))
	} else {
		d.SetId(fmt.Sprintf("alert_resource%s%s%s%s", COLON_SEPARATED, resourceType, COLON_SEPARATED, project))
	}
	return nil
}

func resourcelicloudLogAlertResourceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	parts, err := helper.ParseResourceId(d.Id(), 3)
	if err != nil {
		return errs.WrapError(err)
	}
	resourceType := parts[1]
	if err := retry.Retry(2*time.Minute, func() *retry.RetryError {
		// The alert-resource probe is SLS native data plane (GetResourceRecord /
		// GetLogStore), so it stays on WithLogClient; the old WithLogPopClient
		// wrapper never used its POP client.
		err := readLogAlertResourceOnce(client, resourceType, parts, d)
		if err != nil {
			if errs.IsExpectedErrors(err, []string{errs.LogClientTimeout}) {
				time.Sleep(5 * time.Second)
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	}); err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_log_alert_resource", "ReadAlertResource", errs.AliyunLogGoSdkERROR)
	}
	return nil
}

func readLogAlertResourceOnce(client *connectivity.AliyunClient, resourceType string, parts []string, d *schema.ResourceData) error {
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
			if errs.IsExpectedErrors(err, []string{"ProjectNotExist"}) || errs.IsExpectedErrors(err, []string{"LogStoreNotExist"}) {
				d.SetId("")
				return nil
			}
			return err
		}
		lang := parts[2]
		d.Set("type", resourceType)
		d.Set("project", nil)
		d.Set("lang", lang)
		return nil
	case "project":
		project := parts[2]
		_, err := client.WithLogClient(func(slsClient *sls.Client) (interface{}, error) {
			return slsClient.GetLogStore(project, "internal-alert-history")
		})
		if err != nil {
			if errs.IsExpectedErrors(err, []string{"LogStoreNotExist"}) {
				d.SetId("")
				return nil
			}
			return err
		}
		d.Set("type", resourceType)
		d.Set("project", project)
		d.Set("lang", nil)
		return nil
	default:
		return errs.WrapErrorf(errors.New("type error"), errs.DefaultErrorMsg, "alicloud_log_alert_resource", "ReadAlertResource", errs.AliyunLogGoSdkERROR)
	}
}

func resourcelicloudLogAlertResourceDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}

package sls

import (
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"time"

	sls "github.com/aliyun/aliyun-log-go-sdk"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/validators"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAlicloudLogIngestion() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudLogIngestionCreate,
		Read:   resourceAlicloudLogIngestionRead,
		Update: resourceAlicloudLogIngestionUpdate,
		Delete: resourceAlicloudLogIngestionDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(1 * time.Minute),
			Delete: schema.DefaultTimeout(1 * time.Minute),
			Update: schema.DefaultTimeout(1 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"project": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"ingestion_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringMatch(regexp.MustCompile(`^[A-Za-z0-9-_]{2,128}$`), "ingestion name can only contain lowercase letters, numbers, dashes `-` and underscores `_`. and the name must be 2 to 128 characters long."),
			},
			"logstore": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"display_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"interval": {
				Type:     schema.TypeString,
				Required: true,
			},
			"run_immediately": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"time_zone": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"source": {
				Type:     schema.TypeString,
				Required: true,
				StateFunc: func(v interface{}) string {
					yaml, _ := validators.NormalizeJsonString(v)
					return yaml
				},
				ValidateFunc: validation.StringIsJSON,
			},
		},
	}
}

func resourceAlicloudLogIngestionCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var requestinfo *sls.Client
	ingestion := getIngestion(d)
	project := d.Get("project").(string)
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err := retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		raw, err := client.WithLogClient(func(slsClient *sls.Client) (interface{}, error) {
			requestinfo = slsClient
			return nil, slsClient.CreateIngestion(project, ingestion)
		})
		if err != nil {
			if errs.IsExpectedErrors(err, []string{"InternalServerError", errs.LogClientTimeout}) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		if logs.DebugOn() {
			logs.AddDebug("CreateIngestion", raw, requestinfo, map[string]interface{}{
				"project":  project,
				"logstore": d.Get("logstore").(string),
			})
		}
		return nil
	})
	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_log_ingestion", "CreateIngestion", errs.AliyunLogGoSdkERROR)
	}
	d.SetId(fmt.Sprintf("%s%s%s%s%s", project, COLON_SEPARATED, d.Get("logstore").(string), COLON_SEPARATED, d.Get("ingestion_name").(string)))
	return resourceAlicloudLogIngestionRead(d, meta)
}

func resourceAlicloudLogIngestionRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	logService := Service{client}
	parts, err := helper.ParseResourceId(d.Id(), 3)
	if err != nil {
		return errs.WrapError(err)
	}
	ingestion, err := logService.DescribeLogIngestion(d.Id())
	if err != nil {
		if errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_log_ingestion LogService.DescribeLogIngestion Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}
	d.Set("project", parts[0])
	d.Set("ingestion_name", ingestion.Name)
	d.Set("display_name", ingestion.DisplayName)
	d.Set("description", ingestion.Description)
	d.Set("interval", ingestion.Schedule.Interval)
	d.Set("run_immediately", ingestion.Schedule.RunImmediately)
	d.Set("time_zone", ingestion.Schedule.TimeZone)
	d.Set("logstore", ingestion.IngestionConfiguration.LogStore)
	sourceBytes, err := json.Marshal(ingestion.IngestionConfiguration.DataSource)
	if err != nil {
		return errs.WrapError(err)
	}
	d.Set("source", string(sourceBytes))
	return nil
}

func resourceAlicloudLogIngestionUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var requestinfo *sls.Client
	ingestion := getIngestion(d)
	project := d.Get("project").(string)
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err := retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		raw, err := client.WithLogClient(func(slsClient *sls.Client) (interface{}, error) {
			requestinfo = slsClient
			return nil, slsClient.UpdateIngestion(project, ingestion)
		})
		if err != nil {
			if errs.IsExpectedErrors(err, []string{"InternalServerError", errs.LogClientTimeout}) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		if logs.DebugOn() {
			logs.AddDebug("UpdateIngestion", raw, requestinfo, map[string]interface{}{
				"project":  project,
				"logstore": d.Get("logstore").(string),
			})
		}
		return nil
	})
	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_log_ingestion", "UpdateIngestion", errs.AliyunLogGoSdkERROR)
	}
	return resourceAlicloudLogIngestionRead(d, meta)
}

func resourceAlicloudLogIngestionDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	logService := Service{client}
	var requestinfo *sls.Client
	parts, err := helper.ParseResourceId(d.Id(), 3)
	if err != nil {
		return errs.WrapError(err)
	}
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		raw, err := client.WithLogClient(func(slsClient *sls.Client) (interface{}, error) {
			requestinfo = slsClient
			return nil, slsClient.DeleteIngestion(parts[0], parts[2])
		})
		if err != nil {
			if errs.IsExpectedErrors(err, []string{"InternalServerError", errs.LogClientTimeout}) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		if logs.DebugOn() {
			logs.AddDebug("DeleteIngestion", raw, requestinfo, map[string]interface{}{
				"project":  parts[0],
				"logstore": parts[1],
			})
		}
		return nil
	})
	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_log_ingestion", "DeleteIngestion", errs.AliyunLogGoSdkERROR)
	}
	return errs.WrapError(logService.WaitForLogIngestion(d.Id(), Deleted, DefaultTimeout))
}

func getIngestion(d *schema.ResourceData) *sls.Ingestion {
	ingest := &sls.Ingestion{
		ScheduledJob: sls.ScheduledJob{
			BaseJob: sls.BaseJob{
				Name:        d.Get("ingestion_name").(string),
				DisplayName: d.Get("display_name").(string),
				Description: d.Get("description").(string),
				Type:        sls.INGESTION_JOB,
			},
			Schedule: &sls.Schedule{
				Type:           "FixedRate",
				Interval:       d.Get("interval").(string),
				RunImmediately: d.Get("run_immediately").(bool),
				TimeZone:       d.Get("time_zone").(string),
			},
		},
		IngestionConfiguration: &sls.IngestionConfiguration{
			LogStore: d.Get("logstore").(string),
		},
	}
	s := d.Get("source").(string)
	m := map[string]interface{}{}
	json.Unmarshal([]byte(s), &m)
	ingest.IngestionConfiguration.DataSource = m
	return ingest
}

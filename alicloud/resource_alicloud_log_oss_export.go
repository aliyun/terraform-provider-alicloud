package alicloud

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	sls "github.com/aliyun/aliyun-log-go-sdk"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudLogOssExport() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudLogOssExportCreate,
		Read:   resourceAlicloudLogOssExportRead,
		Update: resourceAlicloudLogOssExportUpdate,
		Delete: resourceAlicloudLogOssExportDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(1 * time.Minute),
			Delete: schema.DefaultTimeout(1 * time.Minute),
			Update: schema.DefaultTimeout(1 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"project_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"logstore_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"export_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"display_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"from_time": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"bucket": {
				Type:     schema.TypeString,
				Required: true,
			},
			"prefix": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"suffix": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"buffer_interval": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"buffer_size": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"role_arn": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"log_read_role_arn": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"compress_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"none", "zstd", "gzip", "snappy"}, false),
				Computed:     true,
			},
			"path_format": {
				Type:     schema.TypeString,
				Required: true,
			},
			"time_zone": {
				Type:     schema.TypeString,
				Required: true,
			},
			"content_type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"json", "parquet", "csv", "orc"}, false),
			},
			"json_enable_tag": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"csv_config_delimiter": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"csv_config_header": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"csv_config_linefeed": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"csv_config_columns": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"csv_config_null": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"csv_config_quote": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"csv_config_escape": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"config_columns": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"type": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
		},
	}
}

func resourceAlicloudLogOssExportCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var requestInfo *sls.Client
	projectName := d.Get("project_name").(string)
	logstoreName := d.Get("logstore_name").(string)
	exportName := d.Get("export_name").(string)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	if err := resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		raw, err := client.WithLogClient(func(slsClient *sls.Client) (interface{}, error) {
			return nil, slsClient.CreateExport(projectName, buildOSSExport(d))
		})
		if err != nil {
			if IsExpectedErrors(err, []string{LogClientTimeout}) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug("CreateOSSExport", raw, requestInfo, map[string]string{
			"project_name":  projectName,
			"logstore_name": logstoreName,
			"export_name":   exportName,
		})
		return nil
	}); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_log_oss_export", "CreateLogOssExport", AliyunLogGoSdkERROR)
	}
	d.SetId(fmt.Sprintf("%s%s%s%s%s", projectName, COLON_SEPARATED, logstoreName, COLON_SEPARATED, exportName))
	return resourceAlicloudLogOssExportRead(d, meta)
}

func resourceAlicloudLogOssExportRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	logService := LogService{client}
	parts, err := ParseResourceId(d.Id(), 3)
	if err != nil {
		return WrapError(err)
	}
	ossExport, err := logService.DescribeLogOssExport(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_log_oss_export LogService.DescribeLogExport Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	ossDataSink := ossExport.ExportConfiguration.DataSink.(*sls.AliyunOSSSink)
	d.Set("project_name", parts[0])
	d.Set("logstore_name", parts[1])
	d.Set("export_name", parts[2])
	d.Set("display_name", ossExport.DisplayName)
	d.Set("from_time", ossExport.ExportConfiguration.FromTime)
	d.Set("log_read_role_arn", ossExport.ExportConfiguration.RoleArn)
	d.Set("bucket", ossDataSink.Bucket)
	d.Set("prefix", ossDataSink.Prefix)
	d.Set("suffix", ossDataSink.Suffix)
	d.Set("buffer_interval", ossDataSink.BufferInterval)
	d.Set("buffer_size", ossDataSink.BufferSize)
	d.Set("time_zone", ossDataSink.TimeZone)
	d.Set("role_arn", ossDataSink.RoleArn)
	d.Set("compress_type", ossDataSink.CompressionType)
	d.Set("path_format", ossDataSink.PathFormat)
	d.Set("content_type", ossDataSink.ContentType)

	if ossDataSink.ContentType == "json" {
		detail := new(sls.JsonContentDetail)
		contentDetailBytes, _ := json.Marshal(ossDataSink.ContentDetail)
		json.Unmarshal(contentDetailBytes, detail)
		d.Set("json_enable_tag", detail.EnableTag)
	} else if ossDataSink.ContentType == "csv" {
		detail := new(sls.CsvContentDetail)
		contentDetailBytes, _ := json.Marshal(ossDataSink.ContentDetail)
		json.Unmarshal(contentDetailBytes, detail)
		d.Set("csv_config_delimiter", detail.Delimiter)
		d.Set("csv_config_header", detail.Header)
		d.Set("csv_config_linefeed", detail.LineFeed)
		d.Set("csv_config_columns", detail.ColumnNames)
		d.Set("csv_config_null", detail.Null)
		d.Set("csv_config_quote", detail.Quote)
	} else if ossDataSink.ContentType == "parquet" || ossDataSink.ContentType == "orc" {
		var config []map[string]interface{}
		contentDetailBytes, _ := json.Marshal(ossDataSink.ContentDetail)
		detail := new(sls.ParquetContentDetail)
		json.Unmarshal(contentDetailBytes, detail)
		for _, column := range detail.Columns {
			tempMap := map[string]interface{}{
				"name": column.Name,
				"type": column.Type,
			}
			config = append(config, tempMap)
		}
		d.Set("config_columns", config)
	}
	return nil
}

func resourceAlicloudLogOssExportUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	parts, err := ParseResourceId(d.Id(), 3)
	if err != nil {
		return WrapError(err)
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	if err := resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
		_, err := client.WithLogClient(func(slsClient *sls.Client) (interface{}, error) {
			return nil, slsClient.RestartExport(parts[0], buildOSSExport(d))
		})
		if err != nil {
			if IsExpectedErrors(err, []string{LogClientTimeout}) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	}); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), "UpdateLogOssExport", AliyunLogGoSdkERROR)
	}
	return resourceAlicloudLogOssExportRead(d, meta)

}

func resourceAlicloudLogOssExportDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	logService := LogService{client}
	parts, err := ParseResourceId(d.Id(), 3)
	if err != nil {
		return WrapError(err)
	}
	var requestInfo *sls.Client
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		raw, err := client.WithLogClient(func(slsClient *sls.Client) (interface{}, error) {
			return nil, slsClient.DeleteExport(parts[0], parts[2])
		})
		if err != nil {
			if IsExpectedErrors(err, []string{LogClientTimeout}) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		if debugOn() {
			addDebug("DeleteLogOssExport", raw, requestInfo, map[string]interface{}{
				"project_name":  parts[0],
				"logstore_name": parts[1],
				"export_name":   parts[2],
			})
		}
		return nil
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"JobNotExist"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_log_oss_export", "DeleteLogOssExport", AliyunLogGoSdkERROR)
	}
	return WrapError(logService.WaitForLogOssExport(d.Id(), Deleted, DefaultTimeout))

}

func buildOSSExport(d *schema.ResourceData) *sls.Export {
	contentType := d.Get("content_type").(string)
	ossExportConfig := &sls.AliyunOSSSink{
		Type:           sls.DataSinkOSS,
		Bucket:         d.Get("bucket").(string),
		PathFormat:     d.Get("path_format").(string),
		PathFormatType: "time",
		BufferSize:     int64(d.Get("buffer_size").(int)),
		BufferInterval: int64(d.Get("buffer_interval").(int)),
		TimeZone:       d.Get("time_zone").(string),
		ContentType:    sls.OSSContentType(contentType),
	}

	roleArn := ""
	if v, ok := d.GetOk("role_arn"); ok {
		roleArn = v.(string)
	}
	ossExportConfig.RoleArn = roleArn
	if v, ok := d.GetOk("prefix"); ok {
		ossExportConfig.Prefix = v.(string)
	}
	if v, ok := d.GetOk("suffix"); ok {
		ossExportConfig.Suffix = v.(string)
	}
	if v, ok := d.GetOk("compress_type"); ok {
		ossExportConfig.CompressionType = sls.OSSCompressionType(v.(string))
	}

	if contentType == "json" {
		enableTag := false
		if v, ok := d.GetOk("json_enable_tag"); ok {
			enableTag = v.(bool)
		}
		ossExportConfig.ContentDetail = sls.JsonContentDetail{EnableTag: enableTag}
	} else if contentType == "parquet" || contentType == "orc" {
		detail := sls.ParquetContentDetail{}
		if configColumns, ok := d.GetOk("config_columns"); ok {
			for _, f := range configColumns.(*schema.Set).List() {
				v := f.(map[string]interface{})
				config := sls.Column{
					Name: v["name"].(string),
					Type: v["type"].(string),
				}
				detail.Columns = append(detail.Columns, config)
			}
		}
		ossExportConfig.ContentDetail = detail
	} else if contentType == "csv" {
		detail := sls.CsvContentDetail{}
		if v, ok := d.GetOk("csv_config_delimiter"); ok {
			detail.Delimiter = v.(string)
		}
		if v, ok := d.GetOk("csv_config_header"); ok {
			detail.Header = v.(bool)
		}
		if v, ok := d.GetOk("csv_config_linefeed"); ok {
			detail.LineFeed = v.(string)
		}
		if v, ok := d.GetOk("csv_config_null"); ok {
			detail.Null = v.(string)
		}
		if v, ok := d.GetOk("csv_config_quote"); ok {
			detail.Quote = v.(string)
		}
		if v, ok := d.GetOk("csv_config_escape"); ok {
			detail.Escape = v.(string)
		}
		columns := []string{}
		if v, ok := d.GetOk("csv_config_columns"); ok {
			for _, v := range v.([]interface{}) {
				columns = append(columns, v.(string))
			}
		}

		detail.ColumnNames = columns
		ossExportConfig.ContentDetail = detail
	}
	fromTime := int64(0)
	if v, ok := d.GetOk("from_time"); ok {
		fromTime = int64(v.(int))
	}
	logReadRoleArn := roleArn
	if v, ok := d.GetOk("log_read_role_arn"); ok {
		logReadRoleArn = v.(string)
	}

	return &sls.Export{
		ScheduledJob: sls.ScheduledJob{
			BaseJob: sls.BaseJob{
				Name:        d.Get("export_name").(string),
				DisplayName: d.Get("display_name").(string),
				Description: "",
				Type:        sls.EXPORT_JOB,
			},
			Schedule: &sls.Schedule{
				Type: "Resident",
			},
		},

		ExportConfiguration: &sls.ExportConfiguration{
			FromTime:   fromTime,
			LogStore:   d.Get("logstore_name").(string),
			Parameters: map[string]string{},
			RoleArn:    logReadRoleArn,
			Version:    sls.ExportVersion2,
			DataSink:   ossExportConfig,
		},
	}

}

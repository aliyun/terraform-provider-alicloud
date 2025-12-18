// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudConfigReportTemplate() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudConfigReportTemplateCreate,
		Read:   resourceAliCloudConfigReportTemplateRead,
		Update: resourceAliCloudConfigReportTemplateUpdate,
		Delete: resourceAliCloudConfigReportTemplateDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"report_file_formats": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringInSlice([]string{"excel"}, false),
			},
			"report_granularity": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"report_language": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"report_scope": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"match_type": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: StringInSlice([]string{"In"}, false),
						},
						"value": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"key": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"report_template_description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"report_template_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"subscription_frequency": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceAliCloudConfigReportTemplateCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateReportTemplate"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})

	if v, ok := d.GetOk("report_template_description"); ok {
		request["ReportTemplateDescription"] = v
	}
	request["ReportTemplateName"] = d.Get("report_template_name")
	if v, ok := d.GetOk("subscription_frequency"); ok {
		request["SubscriptionFrequency"] = v
	}

	if v, ok := d.GetOk("report_scope"); ok {
		reportScopeMaps := make([]map[string]interface{}, 0)
		for _, reportScope := range v.([]interface{}) {
			reportScopeMap := map[string]interface{}{}
			reportScopeArg := reportScope.(map[string]interface{})

			if MatchType, ok := reportScopeArg["match_type"]; ok {
				reportScopeMap["MatchType"] = MatchType
			}

			if val, ok := reportScopeArg["value"]; ok {
				reportScopeMap["Value"] = val
			}

			if key, ok := reportScopeArg["key"]; ok {
				reportScopeMap["Key"] = key
			}

			reportScopeMaps = append(reportScopeMaps, reportScopeMap)
		}
		request["ReportScope"], _ = convertArrayObjectToJsonString(reportScopeMaps)
	}

	if v, ok := d.GetOk("report_granularity"); ok {
		request["ReportGranularity"] = v
	}
	if v, ok := d.GetOk("report_language"); ok {
		request["ReportLanguage"] = v
	}
	if v, ok := d.GetOk("report_file_formats"); ok {
		request["ReportFileFormats"] = v
	}
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("Config", "2020-09-07", action, query, request, true)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_config_report_template", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["ReportTemplateId"]))

	return resourceAliCloudConfigReportTemplateRead(d, meta)
}

func resourceAliCloudConfigReportTemplateRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	configServiceV2 := ConfigServiceV2{client}

	objectRaw, err := configServiceV2.DescribeConfigReportTemplate(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_config_report_template DescribeConfigReportTemplate Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("report_file_formats", objectRaw["ReportFileFormats"])
	d.Set("report_granularity", objectRaw["ReportGranularity"])
	d.Set("report_language", objectRaw["ReportLanguage"])
	d.Set("report_template_description", objectRaw["ReportTemplateDescription"])
	d.Set("report_template_name", objectRaw["ReportTemplateName"])
	d.Set("subscription_frequency", objectRaw["SubscriptionFrequency"])

	reportScopeMaps := make([]map[string]interface{}, 0)
	if reportScopeList, ok := objectRaw["ReportScope"].([]interface{}); ok {
		for _, reportScopeChildRaw := range reportScopeList {
			reportScopeMap := make(map[string]interface{})
			reportScopeChildRaw := reportScopeChildRaw.(map[string]interface{})
			reportScopeMap["key"] = reportScopeChildRaw["Key"]
			reportScopeMap["match_type"] = reportScopeChildRaw["MatchType"]
			reportScopeMap["value"] = reportScopeChildRaw["Value"]

			reportScopeMaps = append(reportScopeMaps, reportScopeMap)
		}
	}
	d.Set("report_scope", reportScopeMaps)

	return nil
}

func resourceAliCloudConfigReportTemplateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false

	var err error
	action := "UpdateReportTemplate"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["ReportTemplateId"] = d.Id()

	if d.HasChange("report_template_description") {
		update = true
		request["ReportTemplateDescription"] = d.Get("report_template_description")
	}

	if d.HasChange("report_template_name") {
		update = true
	}
	request["ReportTemplateName"] = d.Get("report_template_name")
	if d.HasChange("subscription_frequency") {
		update = true
		request["SubscriptionFrequency"] = d.Get("subscription_frequency")
	}

	if d.HasChange("report_scope") {
		update = true
		if v, ok := d.GetOk("report_scope"); ok || d.HasChange("report_scope") {
			reportScopeMapsArray := make([]interface{}, 0)
			for _, dataLoop := range convertToInterfaceArray(v) {
				dataLoopTmp := dataLoop.(map[string]interface{})
				dataLoopMap := make(map[string]interface{})
				dataLoopMap["Value"] = dataLoopTmp["value"]
				dataLoopMap["MatchType"] = dataLoopTmp["match_type"]
				dataLoopMap["Key"] = dataLoopTmp["key"]
				reportScopeMapsArray = append(reportScopeMapsArray, dataLoopMap)
			}
			reportScopeMapsJson, err := json.Marshal(reportScopeMapsArray)
			if err != nil {
				return WrapError(err)
			}
			request["ReportScope"] = string(reportScopeMapsJson)
		}
	}

	if d.HasChange("report_granularity") {
		update = true
		request["ReportGranularity"] = d.Get("report_granularity")
	}

	if d.HasChange("report_language") {
		update = true
		request["ReportLanguage"] = d.Get("report_language")
	}

	if d.HasChange("report_file_formats") {
		update = true
		request["ReportFileFormats"] = d.Get("report_file_formats")
	}

	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("Config", "2020-09-07", action, query, request, true)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, request)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
	}

	return resourceAliCloudConfigReportTemplateRead(d, meta)
}

func resourceAliCloudConfigReportTemplateDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	action := "DeleteReportTemplate"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["ReportTemplateId"] = d.Id()

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("Config", "2020-09-07", action, query, request, true)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		if NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return nil
}

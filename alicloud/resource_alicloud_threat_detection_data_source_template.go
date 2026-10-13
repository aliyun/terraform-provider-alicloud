package alicloud

import (
	"log"
	"strings"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAlicloudThreatDetectionDataSourceTemplate() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudThreatDetectionDataSourceTemplateCreate,
		Read:   resourceAlicloudThreatDetectionDataSourceTemplateRead,
		Update: resourceAlicloudThreatDetectionDataSourceTemplateUpdate,
		Delete: resourceAlicloudThreatDetectionDataSourceTemplateDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"data_source_template_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"data_source_recognize_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"auto_scan_new": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"enabled", "disabled"}, false),
			},
			"data_source_template_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"log_project_pattern": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"log_store_pattern": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"log_user_ids": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"log_region_ids": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"lang": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"role_for": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"create_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"update_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"data_source_from": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"data_source_recognizer": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"region_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAlicloudThreatDetectionDataSourceTemplateCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	threatDetectionServiceV2 := ThreatDetectionServiceV2{client}

	id := d.Get("data_source_template_id").(string)
	// Verify the template exists before adopting it
	_, err := threatDetectionServiceV2.DescribeThreatDetectionDataSourceTemplate(id)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_threat_detection_data_source_template", "ListDataSourceTemplates", AlibabaCloudSdkGoERROR)
	}

	// Apply the desired configuration via UpdateDataSourceTemplate
	if err := updateThreatDetectionDataSourceTemplate(d, client); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_threat_detection_data_source_template", "UpdateDataSourceTemplate", AlibabaCloudSdkGoERROR)
	}

	d.SetId(id)
	return resourceAlicloudThreatDetectionDataSourceTemplateRead(d, meta)
}

func resourceAlicloudThreatDetectionDataSourceTemplateRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	threatDetectionServiceV2 := ThreatDetectionServiceV2{client}

	object, err := threatDetectionServiceV2.DescribeThreatDetectionDataSourceTemplate(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_threat_detection_data_source_template DescribeThreatDetectionDataSourceTemplate Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("data_source_template_id", object["dataSourceTemplateId"])
	d.Set("data_source_recognize_enabled", object["dataSourceRecognizeEnabled"])
	d.Set("auto_scan_new", object["autoScanNew"])
	d.Set("data_source_template_name", object["dataSourceTemplateName"])
	d.Set("log_project_pattern", object["logProjectPattern"])
	d.Set("log_store_pattern", object["logStorePattern"])
	d.Set("create_time", object["createTime"])
	d.Set("update_time", object["updateTime"])
	d.Set("data_source_from", object["dataSourceFrom"])
	d.Set("data_source_recognizer", object["dataSourceRecognizer"])
	d.Set("region_id", client.RegionId)

	if v, ok := object["logUserIds"].([]interface{}); ok {
		d.Set("log_user_ids", v)
	}
	if v, ok := object["logRegionIds"].([]interface{}); ok {
		d.Set("log_region_ids", v)
	}
	return nil
}

func resourceAlicloudThreatDetectionDataSourceTemplateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	if err := updateThreatDetectionDataSourceTemplate(d, client); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), "UpdateDataSourceTemplate", AlibabaCloudSdkGoERROR)
	}
	return resourceAlicloudThreatDetectionDataSourceTemplateRead(d, meta)
}

func resourceAlicloudThreatDetectionDataSourceTemplateDelete(d *schema.ResourceData, meta interface{}) error {
	// No DeleteDataSourceTemplate API exists; the template is preserved.
	log.Printf("[DEBUG] Resource alicloud_threat_detection_data_source_template %s is preserved (no delete API).", d.Id())
	return nil
}

func updateThreatDetectionDataSourceTemplate(d *schema.ResourceData, client *connectivity.AliyunClient) error {
	action := "UpdateDataSourceTemplate"
	request := make(map[string]interface{})
	query := make(map[string]interface{})
	request["DataSourceTemplateId"] = d.Get("data_source_template_id")

	if v, ok := d.GetOk("data_source_template_name"); ok {
		request["DataSourceTemplateName"] = v
	}
	if v, ok := d.GetOk("log_project_pattern"); ok {
		request["LogProjectPattern"] = v
	}
	if v, ok := d.GetOk("log_store_pattern"); ok {
		request["LogStorePattern"] = v
	}
	if v, ok := d.GetOk("lang"); ok {
		request["Lang"] = v
	}
	if v, ok := d.GetOkExists("role_for"); ok {
		request["RoleFor"] = v
	}
	if v, ok := d.GetOk("auto_scan_new"); ok {
		request["AutoScanNew"] = v
	}
	// DataSourceRecognizeEnabled is accepted via query per the API definition
	if v, ok := d.GetOkExists("data_source_recognize_enabled"); ok {
		query["DataSourceRecognizeEnabled"] = v
	}
	if v, ok := d.GetOk("log_user_ids"); ok && len(v.([]interface{})) > 0 {
		request["LogUserIds"] = convertToInterfaceArray(v)
	}
	if v, ok := d.GetOk("log_region_ids"); ok && len(v.([]interface{})) > 0 {
		request["LogRegionIds"] = strings.Join(toStringArray(v), ",")
	}

	var response map[string]interface{}
	var err error
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
		resp, err := client.RpcPost("cloud-siem", "2024-12-12", action, query, request, true)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		response = resp
		addDebug(action, response, request)
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	return nil
}

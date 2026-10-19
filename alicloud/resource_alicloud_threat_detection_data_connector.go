package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAlicloudThreatDetectionDataConnector() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudThreatDetectionDataConnectorCreate,
		Read:   resourceAlicloudThreatDetectionDataConnectorRead,
		Update: resourceAlicloudThreatDetectionDataConnectorUpdate,
		Delete: resourceAlicloudThreatDetectionDataConnectorDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"data_connector_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"oss", "s3", "kafka"}, false),
			},
			"data_connector_config": {
				Type:     schema.TypeString,
				Required: true,
			},
			"src_data_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"dest_data_source_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"log_project_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"log_store_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"log_region_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"data_connector_status": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"enabled", "disabled"}, false),
			},
			"auth_config_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"auth_config_vendor": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"auth_config_product": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"lang": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "en",
			},
			"role_for": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"region_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"data_connector_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"data_connector_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"sls_ingestion_job_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"sls_ingestion_job_state": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"creation_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"update_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func resourceAlicloudThreatDetectionDataConnectorCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	request := make(map[string]interface{})

	if v, ok := d.GetOk("data_connector_type"); ok {
		request["DataConnectorType"] = v
	}
	if v, ok := d.GetOk("data_connector_config"); ok {
		request["DataConnectorConfig"] = v
	}
	if v, ok := d.GetOk("src_data_type"); ok {
		request["SrcDataType"] = v
	}
	if v, ok := d.GetOk("dest_data_source_id"); ok {
		request["DestDataSourceId"] = v
	}
	if v, ok := d.GetOk("log_project_name"); ok {
		request["LogProjectName"] = v
	}
	if v, ok := d.GetOk("log_store_name"); ok {
		request["LogStoreName"] = v
	}
	if v, ok := d.GetOk("log_region_id"); ok {
		request["LogRegionId"] = v
	}
	if v, ok := d.GetOk("auth_config_id"); ok {
		request["AuthConfigId"] = v
	}
	if v, ok := d.GetOk("auth_config_vendor"); ok {
		request["AuthConfigVendor"] = v
	}
	if v, ok := d.GetOk("auth_config_product"); ok {
		request["AuthConfigProduct"] = v
	}
	if v, ok := d.GetOk("data_connector_status"); ok {
		request["DataConnectorStatus"] = v
	}
	if v, ok := d.GetOk("lang"); ok {
		request["Lang"] = v
	}
	if v, ok := d.GetOk("role_for"); ok {
		request["RoleFor"] = v
	}

	var response map[string]interface{}
	action := "CreateDataConnector"
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err := resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *resource.RetryError {
		resp, err := client.RpcPost("cloud-siem", "2024-12-12", action, nil, request, false)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_threat_detection_data_connector", action, AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$.DataConnectorId", response)
	if err != nil || v == nil {
		return WrapErrorf(err, IdMsg, "alicloud_threat_detection_data_connector")
	}
	d.SetId(fmt.Sprint(v))

	return resourceAlicloudThreatDetectionDataConnectorRead(d, meta)
}

func resourceAlicloudThreatDetectionDataConnectorRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cloudSiemService := CloudSiemService{client}

	object, err := cloudSiemService.DescribeThreatDetectionDataConnector(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_threat_detection_data_connector cloudSiemService.DescribeThreatDetectionDataConnector Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("data_connector_id", object["DataConnectorId"])
	d.Set("data_connector_type", object["DataConnectorType"])
	d.Set("data_connector_config", object["DataConnectorConfig"])
	d.Set("data_connector_status", object["DataConnectorStatus"])
	d.Set("src_data_type", object["SrcDataType"])
	d.Set("dest_data_source_id", object["DestDataSourceId"])
	d.Set("log_project_name", object["LogProjectName"])
	d.Set("log_store_name", object["LogStoreName"])
	d.Set("log_region_id", object["LogRegionId"])
	d.Set("auth_config_id", object["AuthConfigId"])
	d.Set("auth_config_vendor", object["AuthConfigVendor"])
	d.Set("auth_config_product", object["AuthConfigProduct"])
	d.Set("data_connector_name", object["DataConnectorName"])
	d.Set("sls_ingestion_job_name", object["SlsIngestionJobName"])
	d.Set("sls_ingestion_job_state", object["SlsIngestionJobState"])
	d.Set("creation_time", object["CreationTime"])
	d.Set("update_time", object["UpdateTime"])
	return nil
}

func resourceAlicloudThreatDetectionDataConnectorUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	update := false
	request := map[string]interface{}{
		"DataConnectorId": d.Id(),
	}

	if d.HasChange("data_connector_config") {
		update = true
		if v, ok := d.GetOk("data_connector_config"); ok {
			request["DataConnectorConfig"] = v
		}
	}
	if d.HasChange("data_connector_status") {
		update = true
		if v, ok := d.GetOk("data_connector_status"); ok {
			request["DataConnectorStatus"] = v
		}
	}
	if d.HasChange("auth_config_vendor") {
		update = true
		if v, ok := d.GetOk("auth_config_vendor"); ok {
			request["AuthConfigVendor"] = v
		}
	}
	if d.HasChange("lang") {
		update = true
		if v, ok := d.GetOk("lang"); ok {
			request["Lang"] = v
		}
	}

	if update {
		action := "UpdateDataConnector"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err := resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			resp, err := client.RpcPost("cloud-siem", "2024-12-12", action, nil, request, false)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(action, resp, request)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
	}

	return resourceAlicloudThreatDetectionDataConnectorRead(d, meta)
}

func resourceAlicloudThreatDetectionDataConnectorDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	request := map[string]interface{}{
		"DataConnectorId": d.Id(),
	}
	if v, ok := d.GetOk("lang"); ok {
		request["Lang"] = v
	}
	if v, ok := d.GetOk("role_for"); ok {
		request["RoleFor"] = v
	}

	action := "DeleteDataConnector"
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err := resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *resource.RetryError {
		resp, err := client.RpcPost("cloud-siem", "2024-12-12", action, nil, request, false)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, resp, request)
		return nil
	})
	if err != nil {
		if NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	return nil
}

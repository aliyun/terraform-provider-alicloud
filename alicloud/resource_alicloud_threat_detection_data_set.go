package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudThreatDetectionDataSet() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudThreatDetectionDataSetCreate,
		Read:   resourceAliCloudThreatDetectionDataSetRead,
		Update: resourceAliCloudThreatDetectionDataSetUpdate,
		Delete: resourceAliCloudThreatDetectionDataSetDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"data_set_id": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"data_set_name": {
				Required: true,
				Type:     schema.TypeString,
			},
			"data_set_field_key_name": {
				Required: true,
				ForceNew: true,
				Type:     schema.TypeString,
			},
			"data_set_file_name": {
				Required: true,
				Type:     schema.TypeString,
			},
			"data_set_description": {
				Optional: true,
				Type:     schema.TypeString,
			},
			"data_set_type": {
				Optional: true,
				ForceNew: true,
				Type:     schema.TypeString,
			},
			"data_set_status": {
				Optional: true,
				Type:     schema.TypeInt,
			},
			"role_for": {
				Optional: true,
				Computed: true,
				Type:     schema.TypeInt,
			},
			"ip_whitelist_recognizers": {
				Optional: true,
				Type:     schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"auto_recognize_status": {
							Optional: true,
							Type:     schema.TypeString,
						},
						"recognize_scope": {
							Optional: true,
							Type:     schema.TypeString,
						},
						"ip_whitelist_recognizer_type": {
							Optional: true,
							Type:     schema.TypeString,
						},
					},
				},
			},
			"lang": {
				Optional: true,
				Computed: true,
				Type:     schema.TypeString,
			},
			"region_id": {
				Optional: true,
				Computed: true,
				Type:     schema.TypeString,
			},
			"create_time": {
				Computed: true,
				Type:     schema.TypeInt,
			},
		},
	}
}

func resourceAliCloudThreatDetectionDataSetCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	request := make(map[string]interface{})
	var err error

	request["DataSetName"] = d.Get("data_set_name")
	request["DataSetFieldKeyName"] = d.Get("data_set_field_key_name")
	request["DataSetFileName"] = d.Get("data_set_file_name")
	if v, ok := d.GetOk("data_set_description"); ok {
		request["DataSetDescription"] = v
	}
	if v, ok := d.GetOk("data_set_type"); ok {
		request["DataSetType"] = v
	}
	if v, ok := d.GetOk("data_set_status"); ok {
		request["DataSetStatus"] = v
	}
	if v, ok := d.GetOk("role_for"); ok {
		request["RoleFor"] = v
	}
	if v, ok := d.GetOk("lang"); ok {
		request["Lang"] = v
	}
	if v, ok := d.GetOk("region_id"); ok {
		request["RegionId"] = v
	}
	if v, ok := d.GetOk("ip_whitelist_recognizers"); ok {
		request["IpWhitelistRecognizers"] = flattenThreatDetectionDataSetIpWhitelistRecognizers(v)
	}

	var response map[string]interface{}
	action := "CreateDataSet"
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		resp, err := client.RpcPost("cloud-siem", "2024-12-12", action, nil, request, true)
		if err != nil {
			if NeedRetry(err) || IsExpectedErrors(err, []string{"IdempotentParameterMismatch"}) {
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_threat_detection_data_set", action, AlibabaCloudSdkGoERROR)
	}

	if v, err := jsonpath.Get("$.DataSetRecordStatistic.DataSetId", response); err != nil || v == nil {
		return WrapErrorf(err, IdMsg, "alicloud_threat_detection_data_set")
	} else {
		d.SetId(fmt.Sprint(v))
	}

	return resourceAliCloudThreatDetectionDataSetRead(d, meta)
}

func resourceAliCloudThreatDetectionDataSetRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	threatDetectionServiceV2 := ThreatDetectionServiceV2{client}

	objectRaw, err := threatDetectionServiceV2.DescribeThreatDetectionDataSet(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_threat_detection_data_set DescribeThreatDetectionDataSet Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("data_set_name", objectRaw["DataSetName"])
	d.Set("data_set_field_key_name", objectRaw["DataSetFieldKeyName"])
	d.Set("data_set_file_name", objectRaw["DataSetFileName"])
	d.Set("data_set_description", objectRaw["DataSetDescription"])
	d.Set("data_set_type", objectRaw["DataSetType"])
	d.Set("data_set_status", objectRaw["DataSetStatus"])
	d.Set("role_for", objectRaw["RoleFor"])
	d.Set("lang", objectRaw["Lang"])
	d.Set("region_id", objectRaw["RegionId"])
	d.Set("create_time", objectRaw["CreateTime"])
	d.Set("ip_whitelist_recognizers", flattenThreatDetectionDataSetIpWhitelistRecognizersResponse(objectRaw["IpWhitelistRecognizers"]))

	return nil
}

func resourceAliCloudThreatDetectionDataSetUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	request := make(map[string]interface{})
	var err error
	update := false

	request["DataSetId"] = d.Id()
	if d.HasChange("data_set_name") {
		update = true
		request["DataSetName"] = d.Get("data_set_name")
	}
	if d.HasChange("data_set_file_name") {
		update = true
		request["DataSetFileName"] = d.Get("data_set_file_name")
	}
	if d.HasChange("data_set_description") {
		update = true
		if v, ok := d.GetOk("data_set_description"); ok {
			request["DataSetDescription"] = v
		}
	}
	if d.HasChange("data_set_status") {
		update = true
		request["DataSetStatus"] = d.Get("data_set_status")
	}
	if d.HasChange("role_for") {
		update = true
		request["RoleFor"] = d.Get("role_for")
	}
	if d.HasChange("lang") {
		update = true
		if v, ok := d.GetOk("lang"); ok {
			request["Lang"] = v
		}
	}
	if d.HasChange("region_id") {
		update = true
		if v, ok := d.GetOk("region_id"); ok {
			request["RegionId"] = v
		}
	}
	if d.HasChange("ip_whitelist_recognizers") {
		update = true
		if v, ok := d.GetOk("ip_whitelist_recognizers"); ok {
			request["IpWhitelistRecognizers"] = flattenThreatDetectionDataSetIpWhitelistRecognizers(v)
		}
	}

	if update {
		action := "UpdateDataSet"
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			resp, err := client.RpcPost("cloud-siem", "2024-12-12", action, nil, request, true)
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

	return resourceAliCloudThreatDetectionDataSetRead(d, meta)
}

func resourceAliCloudThreatDetectionDataSetDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	request := make(map[string]interface{})
	var err error

	request["DataSetId"] = d.Id()
	if v, ok := d.GetOk("role_for"); ok {
		request["RoleFor"] = v
	}
	if v, ok := d.GetOk("lang"); ok {
		request["Lang"] = v
	}
	if v, ok := d.GetOk("region_id"); ok {
		request["RegionId"] = v
	}

	action := "DeleteDataSet"
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		resp, err := client.RpcPost("cloud-siem", "2024-12-12", action, nil, request, true)
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

func flattenThreatDetectionDataSetIpWhitelistRecognizers(v interface{}) []interface{} {
	result := make([]interface{}, 0)
	if items, ok := v.([]interface{}); ok {
		for _, item := range items {
			if m, ok := item.(map[string]interface{}); ok {
				result = append(result, map[string]interface{}{
					"AutoRecognizeStatus":       m["auto_recognize_status"],
					"RecognizeScope":            m["recognize_scope"],
					"IpWhitelistRecognizerType": m["ip_whitelist_recognizer_type"],
				})
			}
		}
	}
	return result
}

func flattenThreatDetectionDataSetIpWhitelistRecognizersResponse(v interface{}) []map[string]interface{} {
	result := make([]map[string]interface{}, 0)
	if items, ok := v.([]interface{}); ok {
		for _, item := range items {
			if m, ok := item.(map[string]interface{}); ok {
				result = append(result, map[string]interface{}{
					"auto_recognize_status":        m["AutoRecognizeStatus"],
					"recognize_scope":              m["RecognizeScope"],
					"ip_whitelist_recognizer_type": m["IpWhitelistRecognizerType"],
				})
			}
		}
	}
	return result
}

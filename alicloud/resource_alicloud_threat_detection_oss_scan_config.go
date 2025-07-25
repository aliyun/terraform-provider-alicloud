// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudThreatDetectionOssScanConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudThreatDetectionOssScanConfigCreate,
		Read:   resourceAliCloudThreatDetectionOssScanConfigRead,
		Update: resourceAliCloudThreatDetectionOssScanConfigUpdate,
		Delete: resourceAliCloudThreatDetectionOssScanConfigDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"all_key_prefix": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"bucket_name_list": {
				Type:     schema.TypeSet,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"decompress_max_file_count": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"decompress_max_layer": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"decryption_list": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"enable": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"end_time": {
				Type:     schema.TypeString,
				Required: true,
			},
			"key_prefix_list": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"key_suffix_list": {
				Type:     schema.TypeSet,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"last_modified_start_time": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"oss_scan_config_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"scan_day_list": {
				Type:     schema.TypeSet,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeInt},
			},
			"start_time": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceAliCloudThreatDetectionOssScanConfigCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateOssScanConfig"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})

	if v, ok := d.GetOkExists("decompress_max_file_count"); ok {
		request["DecompressMaxFileCount"] = v
	}
	if v, ok := d.GetOk("key_suffix_list"); ok {
		keySuffixListMapsArray := v.(*schema.Set).List()
		request["KeySuffixList"] = keySuffixListMapsArray
	}

	request["StartTime"] = d.Get("start_time")
	if v, ok := d.GetOkExists("all_key_prefix"); ok {
		request["AllKeyPrefix"] = v
	}
	if v, ok := d.GetOkExists("decompress_max_layer"); ok {
		request["DecompressMaxLayer"] = v
	}
	if v, ok := d.GetOk("key_prefix_list"); ok {
		keyPrefixListMapsArray := v.(*schema.Set).List()
		request["KeyPrefixList"] = keyPrefixListMapsArray
	}

	request["Enable"] = d.Get("enable")
	if v, ok := d.GetOkExists("last_modified_start_time"); ok {
		request["LastModifiedStartTime"] = v
	}
	if v, ok := d.GetOk("decryption_list"); ok {
		decryptionListMapsArray := v.([]interface{})
		request["DecryptionList"] = decryptionListMapsArray
	}

	if v, ok := d.GetOk("scan_day_list"); ok {
		scanDayListMapsArray := v.(*schema.Set).List()
		request["ScanDayList"] = scanDayListMapsArray
	}

	if v, ok := d.GetOk("oss_scan_config_name"); ok {
		request["Name"] = v
	}
	if v, ok := d.GetOk("bucket_name_list"); ok {
		bucketNameListMapsArray := v.(*schema.Set).List()
		request["BucketNameList"] = bucketNameListMapsArray
	}

	request["EndTime"] = d.Get("end_time")
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("Sas", "2018-12-03", action, query, request, true)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_threat_detection_oss_scan_config", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["Id"]))

	return resourceAliCloudThreatDetectionOssScanConfigRead(d, meta)
}

func resourceAliCloudThreatDetectionOssScanConfigRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	threatDetectionServiceV2 := ThreatDetectionServiceV2{client}

	objectRaw, err := threatDetectionServiceV2.DescribeThreatDetectionOssScanConfig(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_threat_detection_oss_scan_config DescribeThreatDetectionOssScanConfig Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("all_key_prefix", objectRaw["AllKeyPrefix"])
	d.Set("decompress_max_file_count", objectRaw["DecompressMaxFileCount"])
	d.Set("decompress_max_layer", objectRaw["DecompressMaxLayer"])
	d.Set("enable", objectRaw["Enable"])
	d.Set("end_time", objectRaw["EndTime"])
	d.Set("last_modified_start_time", objectRaw["LastModifiedStartTime"])
	d.Set("oss_scan_config_name", objectRaw["Name"])
	d.Set("start_time", objectRaw["StartTime"])

	bucketNameListRaw := make([]interface{}, 0)
	if objectRaw["BucketNameList"] != nil {
		bucketNameListRaw = objectRaw["BucketNameList"].([]interface{})
	}

	d.Set("bucket_name_list", bucketNameListRaw)
	decryptionListRaw := make([]interface{}, 0)
	if objectRaw["DecryptionList"] != nil {
		decryptionListRaw = objectRaw["DecryptionList"].([]interface{})
	}

	d.Set("decryption_list", decryptionListRaw)
	keyPrefixListRaw := make([]interface{}, 0)
	if objectRaw["KeyPrefixList"] != nil {
		keyPrefixListRaw = objectRaw["KeyPrefixList"].([]interface{})
	}

	d.Set("key_prefix_list", keyPrefixListRaw)
	keySuffixListRaw := make([]interface{}, 0)
	if objectRaw["KeySuffixList"] != nil {
		keySuffixListRaw = objectRaw["KeySuffixList"].([]interface{})
	}

	d.Set("key_suffix_list", keySuffixListRaw)
	scanDayListRaw := make([]interface{}, 0)
	if objectRaw["ScanDayList"] != nil {
		scanDayListRaw = objectRaw["ScanDayList"].([]interface{})
	}

	d.Set("scan_day_list", scanDayListRaw)

	return nil
}

func resourceAliCloudThreatDetectionOssScanConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false

	var err error
	action := "UpdateOssScanConfig"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["Id"] = d.Id()

	if d.HasChange("decompress_max_file_count") {
		update = true
	}
	if v, ok := d.GetOkExists("decompress_max_file_count"); ok {
		request["DecompressMaxFileCount"] = v
	}

	if d.HasChange("key_suffix_list") {
		update = true
	}
	if v, ok := d.GetOk("key_suffix_list"); ok {
		keySuffixListMapsArray := v.(*schema.Set).List()
		request["KeySuffixList"] = keySuffixListMapsArray
	}

	if d.HasChange("start_time") {
		update = true
	}
	request["StartTime"] = d.Get("start_time")

	if d.HasChange("all_key_prefix") {
		update = true

		if v, ok := d.GetOkExists("all_key_prefix"); ok {
			request["AllKeyPrefix"] = v
		}
	}

	if d.HasChange("decompress_max_layer") {
		update = true
	}
	if v, ok := d.GetOkExists("decompress_max_layer"); ok {
		request["DecompressMaxLayer"] = v
	}

	if d.HasChange("key_prefix_list") {
		update = true
	}
	if v, ok := d.GetOk("key_prefix_list"); ok {
		keyPrefixListMaps := v.(*schema.Set).List()
		request["KeyPrefixList"] = keyPrefixListMaps
	}

	if d.HasChange("enable") {
		update = true
	}
	request["Enable"] = d.Get("enable")

	if d.HasChange("last_modified_start_time") {
		update = true
	}
	if v, ok := d.GetOkExists("last_modified_start_time"); ok {
		request["LastModifiedStartTime"] = v
	}

	if d.HasChange("decryption_list") {
		update = true
	}
	if v, ok := d.GetOk("decryption_list"); ok {
		decryptionListMapsArray := v.([]interface{})
		request["DecryptionList"] = decryptionListMapsArray
	}

	if d.HasChange("scan_day_list") {
		update = true
	}
	if v, ok := d.GetOk("scan_day_list"); ok {
		scanDayListMaps := v.(*schema.Set).List()
		request["ScanDayList"] = scanDayListMaps
	}

	if d.HasChange("oss_scan_config_name") {
		update = true
	}
	if v, ok := d.GetOk("oss_scan_config_name"); ok {
		request["Name"] = v
	}

	if d.HasChange("bucket_name_list") {
		update = true
	}
	if v, ok := d.GetOk("bucket_name_list"); ok {
		bucketNameListMaps := v.(*schema.Set).List()
		request["BucketNameList"] = bucketNameListMaps
	}

	if d.HasChange("end_time") {
		update = true
	}
	request["EndTime"] = d.Get("end_time")
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("Sas", "2018-12-03", action, query, request, true)
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

	return resourceAliCloudThreatDetectionOssScanConfigRead(d, meta)
}

func resourceAliCloudThreatDetectionOssScanConfigDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	action := "DeleteOssScanConfig"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["Id"] = d.Id()

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("Sas", "2018-12-03", action, query, request, true)

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

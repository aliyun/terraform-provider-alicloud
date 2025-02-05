package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudThreatDetectionHoneypotPreset() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudThreatDetectionHoneypotPresetCreate,
		Read:   resourceAlicloudThreatDetectionHoneypotPresetRead,
		Update: resourceAlicloudThreatDetectionHoneypotPresetUpdate,
		Delete: resourceAlicloudThreatDetectionHoneypotPresetDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"honeypot_image_name": {
				Required: true,
				ForceNew: true,
				Type:     schema.TypeString,
			},
			"honeypot_preset_id": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"meta": {
				Required: true,
				ForceNew: true,
				Type:     schema.TypeList,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"portrait_option": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"burp": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"close", "open"}, false),
						},
						"trojan_git": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringInSlice([]string{"close", "open"}, false),
						},
					},
				},
			},
			"node_id": {
				Required: true,
				ForceNew: true,
				Type:     schema.TypeString,
			},
			"preset_name": {
				Required: true,
				Type:     schema.TypeString,
			},
		},
	}
}

func resourceAlicloudThreatDetectionHoneypotPresetCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	request := make(map[string]interface{})
	var err error

	if v, ok := d.GetOk("honeypot_image_name"); ok {
		request["HoneypotImageName"] = v
	}
	if v, ok := d.GetOk("meta"); ok {
		request["Meta"], _ = convertArrayObjectToJsonString(v.([]interface{})[0])
	}
	if v, ok := d.GetOk("node_id"); ok {
		request["NodeId"] = v
	}
	if v, ok := d.GetOk("preset_name"); ok {
		request["PresetName"] = v
	}

	var response map[string]interface{}
	action := "CreateHoneypotPreset"
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *resource.RetryError {
		resp, err := client.RpcPost("Sas", "2018-12-03", action, nil, request, false)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_threat_detection_honeypot_preset", action, AlibabaCloudSdkGoERROR)
	}

	if v, err := jsonpath.Get("$.HoneypotPreset.HoneypotPresetId", response); err != nil || v == nil {
		return WrapErrorf(err, IdMsg, "alicloud_threat_detection_honeypot_preset")
	} else {
		d.SetId(fmt.Sprint(v))
	}

	return resourceAlicloudThreatDetectionHoneypotPresetRead(d, meta)
}

func resourceAlicloudThreatDetectionHoneypotPresetRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	threatDetectionService := ThreatDetectionService{client}

	object, err := threatDetectionService.DescribeThreatDetectionHoneypotPreset(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_threat_detection_honeypot_preset threatDetectionService.DescribeThreatDetectionHoneypotPreset Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("honeypot_preset_id", object["HoneypotPresetId"])
	d.Set("honeypot_image_name", object["HoneypotImageName"])
	if v, ok := object["Meta"]; ok {
		metaMap, err := convertJsonStringToMap(v.(string))
		if err != nil {
			return WrapError(err)
		}
		metaList := make([]map[string]interface{}, 0)
		metaList = append(metaList, map[string]interface{}{
			"portrait_option": metaMap["portrait_option"],
			"burp":            metaMap["burp"],
			"trojan_git":      metaMap["trojan_git"],
		})
		d.Set("meta", metaList)
	}

	d.Set("node_id", object["NodeId"])
	d.Set("preset_name", object["PresetName"])

	return nil
}

func resourceAlicloudThreatDetectionHoneypotPresetUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	var err error
	update := false
	request := map[string]interface{}{
		"HoneypotPresetId": d.Id(),
	}
	if v, ok := d.GetOk("honeypot_image_name"); ok {
		request["HoneypotImageName"] = v
	}
	if v, ok := d.GetOk("preset_name"); ok {
		request["PresetName"] = v
	}
	if d.HasChange("preset_name") {
		update = true
	}

	if update {
		action := "UpdateHoneypotPreset"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			resp, err := client.RpcPost("Sas", "2018-12-03", action, nil, request, false)
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

	return resourceAlicloudThreatDetectionHoneypotPresetRead(d, meta)
}

func resourceAlicloudThreatDetectionHoneypotPresetDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var err error

	request := map[string]interface{}{
		"HoneypotPresetId": d.Id(),
	}

	action := "DeleteHoneypotPreset"
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *resource.RetryError {
		resp, err := client.RpcPost("Sas", "2018-12-03", action, nil, request, false)
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

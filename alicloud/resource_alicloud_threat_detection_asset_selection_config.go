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

func resourceAliCloudThreatDetectionAssetSelectionConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudThreatDetectionAssetSelectionConfigCreate,
		Read:   resourceAliCloudThreatDetectionAssetSelectionConfigRead,
		Delete: resourceAliCloudThreatDetectionAssetSelectionConfigDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"business_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"platform": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"all"}, false),
			},
			"target_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAliCloudThreatDetectionAssetSelectionConfigCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateAssetSelectionConfig"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	if v, ok := d.GetOk("business_type"); ok {
		request["BusinessType"] = v
	}

	request["TargetType"] = d.Get("target_type")
	if v, ok := d.GetOk("platform"); ok {
		request["Platform"] = v
	}
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_threat_detection_asset_selection_config", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(request["BusinessType"]))

	return resourceAliCloudThreatDetectionAssetSelectionConfigRead(d, meta)
}

func resourceAliCloudThreatDetectionAssetSelectionConfigRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	threatDetectionServiceV2 := ThreatDetectionServiceV2{client}

	objectRaw, err := threatDetectionServiceV2.DescribeThreatDetectionAssetSelectionConfig(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_threat_detection_asset_selection_config DescribeThreatDetectionAssetSelectionConfig Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("platform", objectRaw["Platform"])
	d.Set("target_type", objectRaw["TargetType"])

	d.Set("business_type", d.Id())

	return nil
}

func resourceAliCloudThreatDetectionAssetSelectionConfigDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[WARN] Cannot destroy resource AliCloud Resource Asset Selection Config. Terraform will remove this resource from the state file, however resources may remain.")
	return nil
}

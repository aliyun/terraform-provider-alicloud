// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"log"
	"time"
)

func resourceAliCloudThreatDetectionAttackPathSensitiveAssetConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudThreatDetectionAttackPathSensitiveAssetConfigCreate,
		Read:   resourceAliCloudThreatDetectionAttackPathSensitiveAssetConfigRead,
		Update: resourceAliCloudThreatDetectionAttackPathSensitiveAssetConfigUpdate,
		Delete: resourceAliCloudThreatDetectionAttackPathSensitiveAssetConfigDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"attack_path_asset_list": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"asset_type": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"instance_id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"vendor": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: IntInSlice([]int{0}),
						},
						"region_id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"asset_sub_type": {
							Type:     schema.TypeInt,
							Required: true,
						},
					},
				},
			},
		},
	}
}

func resourceAliCloudThreatDetectionAttackPathSensitiveAssetConfigCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateAttackPathSensitiveAssetConfig"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})

	if v, ok := d.GetOk("attack_path_asset_list"); ok {
		attackPathAssetListMapsArray := make([]interface{}, 0)
		for _, dataLoop := range v.(*schema.Set).List() {
			dataLoopTmp := dataLoop.(map[string]interface{})
			dataLoopMap := make(map[string]interface{})
			dataLoopMap["RegionId"] = dataLoopTmp["region_id"]
			dataLoopMap["AssetSubType"] = dataLoopTmp["asset_sub_type"]
			dataLoopMap["InstanceId"] = dataLoopTmp["instance_id"]
			dataLoopMap["AssetType"] = dataLoopTmp["asset_type"]
			dataLoopMap["Vendor"] = dataLoopTmp["vendor"]
			attackPathAssetListMapsArray = append(attackPathAssetListMapsArray, dataLoopMap)
		}
		request["AttackPathAssetList"] = attackPathAssetListMapsArray
	}

	request["ConfigType"] = "asset_instance"
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_threat_detection_attack_path_sensitive_asset_config", action, AlibabaCloudSdkGoERROR)
	}

	id, _ := jsonpath.Get("$.AttackPathSensitiveAssetConfig.AttackPathSensitiveAssetConfigId", response)
	d.SetId(fmt.Sprint(id))

	return resourceAliCloudThreatDetectionAttackPathSensitiveAssetConfigRead(d, meta)
}

func resourceAliCloudThreatDetectionAttackPathSensitiveAssetConfigRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	threatDetectionServiceV2 := ThreatDetectionServiceV2{client}

	objectRaw, err := threatDetectionServiceV2.DescribeThreatDetectionAttackPathSensitiveAssetConfig(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_threat_detection_attack_path_sensitive_asset_config DescribeThreatDetectionAttackPathSensitiveAssetConfig Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	attackPathAssetListRaw := objectRaw["AttackPathAssetList"]
	attackPathAssetListMaps := make([]map[string]interface{}, 0)
	if attackPathAssetListRaw != nil {
		for _, attackPathAssetListChildRaw := range attackPathAssetListRaw.([]interface{}) {
			attackPathAssetListMap := make(map[string]interface{})
			attackPathAssetListChildRaw := attackPathAssetListChildRaw.(map[string]interface{})
			attackPathAssetListMap["asset_sub_type"] = attackPathAssetListChildRaw["AssetSubType"]
			attackPathAssetListMap["asset_type"] = attackPathAssetListChildRaw["AssetType"]
			attackPathAssetListMap["instance_id"] = attackPathAssetListChildRaw["InstanceId"]
			attackPathAssetListMap["region_id"] = attackPathAssetListChildRaw["RegionId"]
			attackPathAssetListMap["vendor"] = attackPathAssetListChildRaw["Vendor"]

			attackPathAssetListMaps = append(attackPathAssetListMaps, attackPathAssetListMap)
		}
	}
	if err := d.Set("attack_path_asset_list", attackPathAssetListMaps); err != nil {
		return err
	}

	return nil
}

func resourceAliCloudThreatDetectionAttackPathSensitiveAssetConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false

	var err error
	action := "UpdateAttackPathSensitiveAssetConfig"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["AttackPathSensitiveAssetConfigId"] = d.Id()

	if d.HasChange("attack_path_asset_list") {
		update = true
	}
	if v, ok := d.GetOk("attack_path_asset_list"); ok {
		attackPathAssetListMapsArray := make([]interface{}, 0)
		for _, dataLoop := range v.(*schema.Set).List() {
			dataLoopTmp := dataLoop.(map[string]interface{})
			dataLoopMap := make(map[string]interface{})
			dataLoopMap["RegionId"] = dataLoopTmp["region_id"]
			dataLoopMap["AssetSubType"] = dataLoopTmp["asset_sub_type"]
			dataLoopMap["InstanceId"] = dataLoopTmp["instance_id"]
			dataLoopMap["AssetType"] = dataLoopTmp["asset_type"]
			dataLoopMap["Vendor"] = dataLoopTmp["vendor"]
			attackPathAssetListMapsArray = append(attackPathAssetListMapsArray, dataLoopMap)
		}
		request["AttackPathAssetList"] = attackPathAssetListMapsArray
	}

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

	return resourceAliCloudThreatDetectionAttackPathSensitiveAssetConfigRead(d, meta)
}

func resourceAliCloudThreatDetectionAttackPathSensitiveAssetConfigDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	action := "DeleteAttackPathSensitiveAssetConfig"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["AttackPathSensitiveAssetConfigId"] = d.Id()

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
		if IsExpectedErrors(err, []string{"DataNotExists"}) || NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return nil
}

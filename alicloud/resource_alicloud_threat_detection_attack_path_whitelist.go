// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
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

func resourceAliCloudThreatDetectionAttackPathWhitelist() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudThreatDetectionAttackPathWhitelistCreate,
		Read:   resourceAliCloudThreatDetectionAttackPathWhitelistRead,
		Update: resourceAliCloudThreatDetectionAttackPathWhitelistUpdate,
		Delete: resourceAliCloudThreatDetectionAttackPathWhitelistDelete,
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
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"asset_type": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"instance_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"node_type": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: StringInSlice([]string{"start", "end"}, false),
						},
						"region_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"vendor": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: IntInSlice([]int{0, 1}),
						},
						"asset_sub_type": {
							Type:     schema.TypeInt,
							Optional: true,
						},
					},
				},
			},
			"path_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"path_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"remark": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"whitelist_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"whitelist_type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: StringInSlice([]string{"ALL_ASSET", "PART_ASSET"}, false),
			},
		},
	}
}

func resourceAliCloudThreatDetectionAttackPathWhitelistCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateAttackPathWhitelist"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})

	if v, ok := d.GetOk("attack_path_asset_list"); ok {
		attackPathAssetListMapsArray := make([]interface{}, 0)
		for _, dataLoop := range convertToInterfaceArray(v) {
			dataLoopTmp := dataLoop.(map[string]interface{})
			dataLoopMap := make(map[string]interface{})
			dataLoopMap["RegionId"] = dataLoopTmp["region_id"]
			dataLoopMap["AssetType"] = dataLoopTmp["asset_type"]
			dataLoopMap["NodeType"] = dataLoopTmp["node_type"]
			dataLoopMap["AssetSubType"] = dataLoopTmp["asset_sub_type"]
			dataLoopMap["InstanceId"] = dataLoopTmp["instance_id"]
			dataLoopMap["Vendor"] = dataLoopTmp["vendor"]
			attackPathAssetListMapsArray = append(attackPathAssetListMapsArray, dataLoopMap)
		}
		request["AttackPathAssetList"] = attackPathAssetListMapsArray
	}

	request["WhitelistName"] = d.Get("whitelist_name")
	request["WhitelistType"] = d.Get("whitelist_type")
	request["PathType"] = d.Get("path_type")
	request["PathName"] = d.Get("path_name")
	if v, ok := d.GetOk("remark"); ok {
		request["Remark"] = v
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_threat_detection_attack_path_whitelist", action, AlibabaCloudSdkGoERROR)
	}

	id, _ := jsonpath.Get("$.AttackPathWhitelist.AttackPathWhitelistId", response)
	d.SetId(fmt.Sprint(id))

	return resourceAliCloudThreatDetectionAttackPathWhitelistRead(d, meta)
}

func resourceAliCloudThreatDetectionAttackPathWhitelistRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	threatDetectionServiceV2 := ThreatDetectionServiceV2{client}

	objectRaw, err := threatDetectionServiceV2.DescribeThreatDetectionAttackPathWhitelist(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_threat_detection_attack_path_whitelist DescribeThreatDetectionAttackPathWhitelist Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("path_name", objectRaw["PathName"])
	d.Set("path_type", objectRaw["PathType"])
	d.Set("remark", objectRaw["Remark"])
	d.Set("whitelist_name", objectRaw["WhitelistName"])
	d.Set("whitelist_type", objectRaw["WhitelistType"])

	attackPathAssetListRaw := objectRaw["AttackPathAssetList"]
	attackPathAssetListMaps := make([]map[string]interface{}, 0)
	if attackPathAssetListRaw != nil {
		for _, attackPathAssetListChildRaw := range convertToInterfaceArray(attackPathAssetListRaw) {
			attackPathAssetListMap := make(map[string]interface{})
			attackPathAssetListChildRaw := attackPathAssetListChildRaw.(map[string]interface{})
			attackPathAssetListMap["asset_sub_type"] = attackPathAssetListChildRaw["AssetSubType"]
			attackPathAssetListMap["asset_type"] = attackPathAssetListChildRaw["AssetType"]
			attackPathAssetListMap["instance_id"] = attackPathAssetListChildRaw["InstanceId"]
			attackPathAssetListMap["node_type"] = attackPathAssetListChildRaw["NodeType"]
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

func resourceAliCloudThreatDetectionAttackPathWhitelistUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false

	var err error
	action := "UpdateAttackPathWhitelist"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["AttackPathWhitelistId"] = d.Id()

	if d.HasChange("attack_path_asset_list") {
		update = true
		if v, ok := d.GetOk("attack_path_asset_list"); ok || d.HasChange("attack_path_asset_list") {
			attackPathAssetListMapsArray := make([]interface{}, 0)
			for _, dataLoop := range convertToInterfaceArray(v) {
				dataLoopTmp := dataLoop.(map[string]interface{})
				dataLoopMap := make(map[string]interface{})
				dataLoopMap["RegionId"] = dataLoopTmp["region_id"]
				dataLoopMap["AssetType"] = dataLoopTmp["asset_type"]
				dataLoopMap["NodeType"] = dataLoopTmp["node_type"]
				dataLoopMap["AssetSubType"] = dataLoopTmp["asset_sub_type"]
				dataLoopMap["InstanceId"] = dataLoopTmp["instance_id"]
				dataLoopMap["Vendor"] = dataLoopTmp["vendor"]
				attackPathAssetListMapsArray = append(attackPathAssetListMapsArray, dataLoopMap)
			}
			request["AttackPathAssetList"] = attackPathAssetListMapsArray
		}
	}

	if d.HasChange("whitelist_name") {
		update = true
	}
	request["WhitelistName"] = d.Get("whitelist_name")
	if d.HasChange("whitelist_type") {
		update = true
	}
	request["WhitelistType"] = d.Get("whitelist_type")
	if d.HasChange("path_type") {
		update = true
	}
	request["PathType"] = d.Get("path_type")
	if d.HasChange("path_name") {
		update = true
	}
	request["PathName"] = d.Get("path_name")
	if d.HasChange("remark") {
		update = true
		request["Remark"] = d.Get("remark")
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

	return resourceAliCloudThreatDetectionAttackPathWhitelistRead(d, meta)
}

func resourceAliCloudThreatDetectionAttackPathWhitelistDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	action := "DeleteAttackPathWhitelist"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["AttackPathWhitelistId"] = d.Id()

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

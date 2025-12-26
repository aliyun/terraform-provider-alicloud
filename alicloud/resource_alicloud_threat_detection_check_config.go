// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"encoding/json"
	"log"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudThreatDetectionCheckConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudThreatDetectionCheckConfigCreate,
		Read:   resourceAliCloudThreatDetectionCheckConfigRead,
		Update: resourceAliCloudThreatDetectionCheckConfigUpdate,
		Delete: resourceAliCloudThreatDetectionCheckConfigDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"added_check": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"check_id": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"section_id": {
							Type:     schema.TypeInt,
							Optional: true,
						},
					},
				},
			},
			"config_requirement_ids": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"remove_ids": {
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeInt},
						},
						"add_ids": {
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeInt},
						},
					},
				},
			},
			"config_standard_ids": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"remove_ids": {
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeInt},
						},
						"add_ids": {
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeInt},
						},
					},
				},
			},
			"configure": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"cycle_days": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeInt},
			},
			"enable_add_check": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"enable_auto_check": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"end_time": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"removed_check": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"check_id": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"section_id": {
							Type:     schema.TypeInt,
							Optional: true,
						},
					},
				},
			},
			"standard_ids": {
				Type:       schema.TypeList,
				Optional:   true,
				Deprecated: "Field 'standard_ids' has been deprecated from provider version 1.267.0. An array that consists of the information about the check item.",
				Elem:       &schema.Schema{Type: schema.TypeInt},
			},
			"start_time": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"system_config": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"vendors": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceAliCloudThreatDetectionCheckConfigCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "ChangeCheckConfig"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["RegionId"] = client.RegionId

	if v, ok := d.GetOk("cycle_days"); ok {
		cycleDaysMapsArray := convertToInterfaceArray(v)

		request["CycleDays"] = cycleDaysMapsArray
	}

	if v, ok := d.GetOkExists("enable_auto_check"); ok {
		request["EnableAutoCheck"] = v
	}
	if v, ok := d.GetOk("added_check"); ok {
		addedCheckMapsArray := make([]interface{}, 0)
		for _, dataLoop1 := range convertToInterfaceArray(v) {
			dataLoop1Tmp := dataLoop1.(map[string]interface{})
			dataLoop1Map := make(map[string]interface{})
			dataLoop1Map["CheckId"] = dataLoop1Tmp["check_id"]
			dataLoop1Map["SectionId"] = dataLoop1Tmp["section_id"]
			addedCheckMapsArray = append(addedCheckMapsArray, dataLoop1Map)
		}
		request["AddedCheck"] = addedCheckMapsArray
	}

	if v, ok := d.GetOk("removed_check"); ok {
		removedCheckMapsArray := make([]interface{}, 0)
		for _, dataLoop2 := range convertToInterfaceArray(v) {
			dataLoop2Tmp := dataLoop2.(map[string]interface{})
			dataLoop2Map := make(map[string]interface{})
			dataLoop2Map["CheckId"] = dataLoop2Tmp["check_id"]
			dataLoop2Map["SectionId"] = dataLoop2Tmp["section_id"]
			removedCheckMapsArray = append(removedCheckMapsArray, dataLoop2Map)
		}
		request["RemovedCheck"] = removedCheckMapsArray
	}

	if v, ok := d.GetOk("configure"); ok {
		request["Configure"] = v
	}
	if v, ok := d.GetOkExists("start_time"); ok {
		request["StartTime"] = v
	}
	if v, ok := d.GetOkExists("enable_add_check"); ok {
		request["EnableAddCheck"] = v
	}
	if v, ok := d.GetOkExists("system_config"); ok {
		request["SystemConfig"] = v
	}
	if v, ok := d.GetOk("standard_ids"); ok {
		standardIdsMapsArray := convertToInterfaceArray(v)

		request["StandardIds"] = standardIdsMapsArray
	}

	configStandardIds := make(map[string]interface{})

	if v := d.Get("config_standard_ids"); !IsNil(v) {
		removeIds1, _ := jsonpath.Get("$[0].remove_ids", v)
		if removeIds1 != nil && removeIds1 != "" {
			configStandardIds["RemoveIds"] = removeIds1
		}
		addIds1, _ := jsonpath.Get("$[0].add_ids", v)
		if addIds1 != nil && addIds1 != "" {
			configStandardIds["AddIds"] = addIds1
		}

		configStandardIdsJson, err := json.Marshal(configStandardIds)
		if err != nil {
			return WrapError(err)
		}
		request["ConfigStandardIds"] = string(configStandardIdsJson)
	}

	if v, ok := d.GetOk("vendors"); ok {
		vendorsMapsArray := convertToInterfaceArray(v)

		request["Vendors"] = vendorsMapsArray
	}

	configRequirementIds := make(map[string]interface{})

	if v := d.Get("config_requirement_ids"); !IsNil(v) {
		removeIds3, _ := jsonpath.Get("$[0].remove_ids", v)
		if removeIds3 != nil && removeIds3 != "" {
			configRequirementIds["RemoveIds"] = removeIds3
		}
		addIds3, _ := jsonpath.Get("$[0].add_ids", v)
		if addIds3 != nil && addIds3 != "" {
			configRequirementIds["AddIds"] = addIds3
		}

		configRequirementIdsJson, err := json.Marshal(configRequirementIds)
		if err != nil {
			return WrapError(err)
		}
		request["ConfigRequirementIds"] = string(configRequirementIdsJson)
	}

	if v, ok := d.GetOkExists("end_time"); ok {
		request["EndTime"] = v
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_threat_detection_check_config", action, AlibabaCloudSdkGoERROR)
	}

	accountId, err := client.AccountId()
	d.SetId(accountId)

	return resourceAliCloudThreatDetectionCheckConfigRead(d, meta)
}

func resourceAliCloudThreatDetectionCheckConfigRead(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	threatDetectionServiceV2 := ThreatDetectionServiceV2{client}

	objectRaw, err := threatDetectionServiceV2.DescribeThreatDetectionCheckConfig(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_threat_detection_check_config DescribeThreatDetectionCheckConfig Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("enable_add_check", objectRaw["EnableAddCheck"])
	d.Set("enable_auto_check", objectRaw["EnableAutoCheck"])
	d.Set("end_time", objectRaw["EndTime"])
	d.Set("start_time", objectRaw["StartTime"])

	cycleDaysRaw := make([]interface{}, 0)
	if objectRaw["CycleDays"] != nil {
		cycleDaysRaw = convertToInterfaceArray(objectRaw["CycleDays"])
	}

	d.Set("cycle_days", cycleDaysRaw)

	return nil
}

func resourceAliCloudThreatDetectionCheckConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	return resourceAliCloudThreatDetectionCheckConfigCreate(d, meta)

}

func resourceAliCloudThreatDetectionCheckConfigDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[WARN] Cannot destroy resource AliCloud Resource Check Config. Terraform will remove this resource from the state file, however resources may remain.")
	return nil
}

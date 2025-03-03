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

func resourceAliCloudSlsCollectionPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudSlsCollectionPolicyCreate,
		Read:   resourceAliCloudSlsCollectionPolicyRead,
		Update: resourceAliCloudSlsCollectionPolicyUpdate,
		Delete: resourceAliCloudSlsCollectionPolicyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"centralize_config": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"dest_ttl": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"dest_logstore": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"dest_region": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"dest_project": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"centralize_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"data_code": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"data_config": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"data_region": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"data_project": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"enabled": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"policy_config": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"resource_tags": {
							Type:     schema.TypeMap,
							Optional: true,
						},
						"regions": {
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"instance_ids": {
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"resource_mode": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: StringInSlice([]string{"all", "instanceMode", "attributeMode"}, false),
						},
					},
				},
			},
			"policy_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"product_code": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"resource_directory": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"account_group_type": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"members": {
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
		},
	}
}

func resourceAliCloudSlsCollectionPolicyCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := fmt.Sprintf("/collectionpolicy")
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	body := make(map[string]interface{})
	hostMap := make(map[string]*string)
	var err error
	request = make(map[string]interface{})
	if v, ok := d.GetOk("policy_name"); ok {
		request["policyName"] = v
	}

	if v, ok := d.GetOkExists("centralize_enabled"); ok {
		request["centralizeEnabled"] = v
	}
	request["enabled"] = d.Get("enabled")
	request["productCode"] = d.Get("product_code")
	request["dataCode"] = d.Get("data_code")
	objectDataLocalMap := make(map[string]interface{})

	if v := d.Get("policy_config"); v != nil {
		resourceMode1, _ := jsonpath.Get("$[0].resource_mode", d.Get("policy_config"))
		if resourceMode1 != nil && resourceMode1 != "" {
			objectDataLocalMap["resourceMode"] = resourceMode1
		}
		resourceTags1, _ := jsonpath.Get("$[0].resource_tags", d.Get("policy_config"))
		if resourceTags1 != nil && resourceTags1 != "" {
			objectDataLocalMap["resourceTags"] = resourceTags1
		}
		regions1, _ := jsonpath.Get("$[0].regions", v)
		if regions1 != nil && regions1 != "" {
			objectDataLocalMap["regions"] = regions1
		}
		instanceIds1, _ := jsonpath.Get("$[0].instance_ids", v)
		if instanceIds1 != nil && instanceIds1 != "" {
			objectDataLocalMap["instanceIds"] = instanceIds1
		}

		request["policyConfig"] = objectDataLocalMap
	}

	objectDataLocalMap1 := make(map[string]interface{})

	if v := d.Get("centralize_config"); !IsNil(v) {
		destRegion1, _ := jsonpath.Get("$[0].dest_region", d.Get("centralize_config"))
		if destRegion1 != nil && destRegion1 != "" {
			objectDataLocalMap1["destRegion"] = destRegion1
		}
		destProject1, _ := jsonpath.Get("$[0].dest_project", d.Get("centralize_config"))
		if destProject1 != nil && destProject1 != "" {
			objectDataLocalMap1["destProject"] = destProject1
		}
		destLogstore1, _ := jsonpath.Get("$[0].dest_logstore", d.Get("centralize_config"))
		if destLogstore1 != nil && destLogstore1 != "" {
			objectDataLocalMap1["destLogstore"] = destLogstore1
		}
		destTtl, _ := jsonpath.Get("$[0].dest_ttl", d.Get("centralize_config"))
		if destTtl != nil && destTtl != "" {
			objectDataLocalMap1["destTTL"] = destTtl
		}

		request["centralizeConfig"] = objectDataLocalMap1
	}

	objectDataLocalMap2 := make(map[string]interface{})

	if v := d.Get("data_config"); !IsNil(v) {
		dataRegion1, _ := jsonpath.Get("$[0].data_region", d.Get("data_config"))
		if dataRegion1 != nil && dataRegion1 != "" {
			objectDataLocalMap2["dataRegion"] = dataRegion1
		}

		request["dataConfig"] = objectDataLocalMap2
	}

	objectDataLocalMap3 := make(map[string]interface{})

	if v := d.Get("resource_directory"); !IsNil(v) {
		accountGroupType1, _ := jsonpath.Get("$[0].account_group_type", d.Get("resource_directory"))
		if accountGroupType1 != nil && accountGroupType1 != "" {
			objectDataLocalMap3["accountGroupType"] = accountGroupType1
		}
		members1, _ := jsonpath.Get("$[0].members", v)
		if members1 != nil && members1 != "" {
			objectDataLocalMap3["members"] = members1
		}

		request["resourceDirectory"] = objectDataLocalMap3
	}

	body = request
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.Do("Sls", roaParam("POST", "2020-12-30", "UpsertCollectionPolicy", action), query, body, nil, hostMap, false)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_sls_collection_policy", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(request["policyName"]))

	slsServiceV2 := SlsServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{"#CHECKSET"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, slsServiceV2.SlsCollectionPolicyStateRefreshFunc(d.Id(), "#policyName", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAliCloudSlsCollectionPolicyRead(d, meta)
}

func resourceAliCloudSlsCollectionPolicyRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	slsServiceV2 := SlsServiceV2{client}

	objectRaw, err := slsServiceV2.DescribeSlsCollectionPolicy(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_sls_collection_policy DescribeSlsCollectionPolicy Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	if objectRaw["centralizeEnabled"] != nil {
		d.Set("centralize_enabled", objectRaw["centralizeEnabled"])
	}
	if objectRaw["dataCode"] != nil {
		d.Set("data_code", objectRaw["dataCode"])
	}
	if objectRaw["enabled"] != nil {
		d.Set("enabled", objectRaw["enabled"])
	}
	if objectRaw["productCode"] != nil {
		d.Set("product_code", objectRaw["productCode"])
	}
	if objectRaw["policyName"] != nil {
		d.Set("policy_name", objectRaw["policyName"])
	}

	centralizeConfigMaps := make([]map[string]interface{}, 0)
	centralizeConfigMap := make(map[string]interface{})
	centralizeConfig1Raw := make(map[string]interface{})
	if objectRaw["centralizeConfig"] != nil {
		centralizeConfig1Raw = objectRaw["centralizeConfig"].(map[string]interface{})
	}
	if len(centralizeConfig1Raw) > 0 {
		centralizeConfigMap["dest_logstore"] = centralizeConfig1Raw["destLogstore"]
		centralizeConfigMap["dest_project"] = centralizeConfig1Raw["destProject"]
		centralizeConfigMap["dest_region"] = centralizeConfig1Raw["destRegion"]
		centralizeConfigMap["dest_ttl"] = centralizeConfig1Raw["destTTL"]

		centralizeConfigMaps = append(centralizeConfigMaps, centralizeConfigMap)
	}
	if objectRaw["centralizeConfig"] != nil {
		if err := d.Set("centralize_config", centralizeConfigMaps); err != nil {
			return err
		}
	}
	dataConfigMaps := make([]map[string]interface{}, 0)
	dataConfigMap := make(map[string]interface{})
	dataConfig1Raw := make(map[string]interface{})
	if objectRaw["dataConfig"] != nil {
		dataConfig1Raw = objectRaw["dataConfig"].(map[string]interface{})
	}
	if len(dataConfig1Raw) > 0 {
		dataConfigMap["data_project"] = dataConfig1Raw["dataProject"]
		dataConfigMap["data_region"] = dataConfig1Raw["dataRegion"]

		dataConfigMaps = append(dataConfigMaps, dataConfigMap)
	}
	if objectRaw["dataConfig"] != nil {
		if err := d.Set("data_config", dataConfigMaps); err != nil {
			return err
		}
	}
	policyConfigMaps := make([]map[string]interface{}, 0)
	policyConfigMap := make(map[string]interface{})
	policyConfig1Raw := make(map[string]interface{})
	if objectRaw["policyConfig"] != nil {
		policyConfig1Raw = objectRaw["policyConfig"].(map[string]interface{})
	}
	if len(policyConfig1Raw) > 0 {
		policyConfigMap["resource_mode"] = policyConfig1Raw["resourceMode"]
		policyConfigMap["resource_tags"] = policyConfig1Raw["resourceTags"]

		instanceIds1Raw := make([]interface{}, 0)
		if policyConfig1Raw["instanceIds"] != nil {
			instanceIds1Raw = policyConfig1Raw["instanceIds"].([]interface{})
		}

		policyConfigMap["instance_ids"] = instanceIds1Raw
		regions1Raw := make([]interface{}, 0)
		if policyConfig1Raw["regions"] != nil {
			regions1Raw = policyConfig1Raw["regions"].([]interface{})
		}

		policyConfigMap["regions"] = regions1Raw
		policyConfigMaps = append(policyConfigMaps, policyConfigMap)
	}
	if objectRaw["policyConfig"] != nil {
		if err := d.Set("policy_config", policyConfigMaps); err != nil {
			return err
		}
	}
	resourceDirectoryMaps := make([]map[string]interface{}, 0)
	resourceDirectoryMap := make(map[string]interface{})
	resourceDirectory1Raw := make(map[string]interface{})
	if objectRaw["resourceDirectory"] != nil {
		resourceDirectory1Raw = objectRaw["resourceDirectory"].(map[string]interface{})
	}
	if len(resourceDirectory1Raw) > 0 {
		resourceDirectoryMap["account_group_type"] = resourceDirectory1Raw["accountGroupType"]

		members1Raw := make([]interface{}, 0)
		if resourceDirectory1Raw["members"] != nil {
			members1Raw = resourceDirectory1Raw["members"].([]interface{})
		}

		resourceDirectoryMap["members"] = members1Raw
		resourceDirectoryMaps = append(resourceDirectoryMaps, resourceDirectoryMap)
	}
	if objectRaw["resourceDirectory"] != nil && objectRaw["resourceDirectory"].(map[string]interface{})["accountGroupType"] != "" {
		if err := d.Set("resource_directory", resourceDirectoryMaps); err != nil {
			return err
		}
	}

	d.Set("policy_name", d.Id())

	return nil
}

func resourceAliCloudSlsCollectionPolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	var body map[string]interface{}
	update := false
	action := fmt.Sprintf("/collectionpolicy")
	var err error
	request = make(map[string]interface{})
	query = make(map[string]*string)
	body = make(map[string]interface{})
	hostMap := make(map[string]*string)
	request["policyName"] = d.Id()

	if d.HasChange("centralize_enabled") {
		update = true
	}
	if v, ok := d.GetOkExists("centralize_enabled"); ok || d.HasChange("centralize_enabled") {
		request["centralizeEnabled"] = v
	}
	if d.HasChange("enabled") {
		update = true
	}
	request["enabled"] = d.Get("enabled")
	if d.HasChange("product_code") {
		update = true
	}
	request["productCode"] = d.Get("product_code")
	if d.HasChange("data_code") {
		update = true
	}
	request["dataCode"] = d.Get("data_code")
	if d.HasChange("policy_config") {
		update = true
	}
	objectDataLocalMap := make(map[string]interface{})

	if v := d.Get("policy_config"); v != nil {
		resourceMode1, _ := jsonpath.Get("$[0].resource_mode", v)
		if resourceMode1 != nil && (d.HasChange("policy_config.0.resource_mode") || resourceMode1 != "") {
			objectDataLocalMap["resourceMode"] = resourceMode1
		}
		resourceTags1, _ := jsonpath.Get("$[0].resource_tags", v)
		if resourceTags1 != nil && (d.HasChange("policy_config.0.resource_tags") || resourceTags1 != "") {
			objectDataLocalMap["resourceTags"] = resourceTags1
		}
		regions1, _ := jsonpath.Get("$[0].regions", d.Get("policy_config"))
		if regions1 != nil && (d.HasChange("policy_config.0.regions") || regions1 != "") {
			objectDataLocalMap["regions"] = regions1
		}
		instanceIds1, _ := jsonpath.Get("$[0].instance_ids", d.Get("policy_config"))
		if instanceIds1 != nil && (d.HasChange("policy_config.0.instance_ids") || instanceIds1 != "") {
			objectDataLocalMap["instanceIds"] = instanceIds1
		}

		request["policyConfig"] = objectDataLocalMap
	}

	if d.HasChange("centralize_config") {
		update = true
	}
	objectDataLocalMap1 := make(map[string]interface{})

	if v := d.Get("centralize_config"); !IsNil(v) || d.HasChange("centralize_config") {
		destRegion1, _ := jsonpath.Get("$[0].dest_region", v)
		if destRegion1 != nil && (d.HasChange("centralize_config.0.dest_region") || destRegion1 != "") {
			objectDataLocalMap1["destRegion"] = destRegion1
		}
		destProject1, _ := jsonpath.Get("$[0].dest_project", v)
		if destProject1 != nil && (d.HasChange("centralize_config.0.dest_project") || destProject1 != "") {
			objectDataLocalMap1["destProject"] = destProject1
		}
		destLogstore1, _ := jsonpath.Get("$[0].dest_logstore", v)
		if destLogstore1 != nil && (d.HasChange("centralize_config.0.dest_logstore") || destLogstore1 != "") {
			objectDataLocalMap1["destLogstore"] = destLogstore1
		}
		destTtl, _ := jsonpath.Get("$[0].dest_ttl", v)
		if destTtl != nil && (d.HasChange("centralize_config.0.dest_ttl") || destTtl != "") {
			objectDataLocalMap1["destTTL"] = destTtl
		}

		request["centralizeConfig"] = objectDataLocalMap1
	}

	if d.HasChange("resource_directory") {
		update = true
	}
	objectDataLocalMap2 := make(map[string]interface{})

	if v := d.Get("resource_directory"); !IsNil(v) || d.HasChange("resource_directory") {
		accountGroupType1, _ := jsonpath.Get("$[0].account_group_type", v)
		if accountGroupType1 != nil && (d.HasChange("resource_directory.0.account_group_type") || accountGroupType1 != "") {
			objectDataLocalMap2["accountGroupType"] = accountGroupType1
		}
		members1, _ := jsonpath.Get("$[0].members", d.Get("resource_directory"))
		if members1 != nil && (d.HasChange("resource_directory.0.members") || members1 != "") {
			objectDataLocalMap2["members"] = members1
		}

		request["resourceDirectory"] = objectDataLocalMap2
	}

	if d.HasChange("data_config") {
		update = true
	}
	objectDataLocalMap3 := make(map[string]interface{})

	if v := d.Get("data_config"); !IsNil(v) || d.HasChange("data_config") {
		dataRegion1, _ := jsonpath.Get("$[0].data_region", v)
		if dataRegion1 != nil && (d.HasChange("data_config.0.data_region") || dataRegion1 != "") {
			objectDataLocalMap3["dataRegion"] = dataRegion1
		}

		request["dataConfig"] = objectDataLocalMap3
	}

	body = request
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.Do("Sls", roaParam("POST", "2020-12-30", "UpsertCollectionPolicy", action), query, body, nil, hostMap, false)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(action, response, request)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		slsServiceV2 := SlsServiceV2{client}
		stateConf := BuildStateConf([]string{}, []string{"#CHECKSET"}, d.Timeout(schema.TimeoutUpdate), 10*time.Second, slsServiceV2.SlsCollectionPolicyStateRefreshFunc(d.Id(), "#policyName", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}

	return resourceAliCloudSlsCollectionPolicyRead(d, meta)
}

func resourceAliCloudSlsCollectionPolicyDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	policyName := d.Id()
	action := fmt.Sprintf("/collectionpolicy/%s", policyName)
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	hostMap := make(map[string]*string)
	var err error
	request = make(map[string]interface{})
	request["policyName"] = d.Id()

	if v, ok := d.GetOk("product_code"); ok {
		query["productCode"] = StringPointer(v.(string))
	}

	if v, ok := d.GetOk("data_code"); ok {
		query["dataCode"] = StringPointer(v.(string))
	}

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.Do("Sls", roaParam("DELETE", "2020-12-30", "DeleteCollectionPolicy", action), query, nil, nil, hostMap, false)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})

	if err != nil {
		if IsExpectedErrors(err, []string{"PolicyNotExist"}) || NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	slsServiceV2 := SlsServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 5*time.Second, slsServiceV2.SlsCollectionPolicyStateRefreshFunc(d.Id(), "policyName", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return nil
}

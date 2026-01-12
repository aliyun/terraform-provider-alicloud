package alicloud

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudHbrPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudHbrPolicyCreate,
		Read:   resourceAliCloudHbrPolicyRead,
		Update: resourceAliCloudHbrPolicyUpdate,
		Delete: resourceAliCloudHbrPolicyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"policy_description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"policy_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"policy_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"rules": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"keep_latest_snapshots": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"rule_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"schedule": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"tag_filters": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"operator": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"value": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"key": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
						"backup_type": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: StringInSlice([]string{"COMPLETE", "INCREMENTAL"}, false),
						},
						"archive_days": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
						"rule_type": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: StringInSlice([]string{"BACKUP", "TRANSITION", "REPLICATION", "TAG"}, false),
						},
						"data_source_filters": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"source_type": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
						"retention": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"vault_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"retention_rules": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"advanced_retention_type": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: StringInSlice([]string{"WEEKLY", "MONTHLY", "YEARLY", "DAILY"}, false),
									},
									"retention": {
										Type:     schema.TypeInt,
										Optional: true,
									},
								},
							},
						},
						"replication_region_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
		},
	}
}

func resourceAliCloudHbrPolicyCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreatePolicyV2"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})

	if v, ok := d.GetOk("rules"); ok {
		rulesMapsArray := make([]interface{}, 0)
		for _, dataLoop := range convertToInterfaceArray(v) {
			dataLoopTmp := dataLoop.(map[string]interface{})
			dataLoopMap := make(map[string]interface{})
			dataLoopMap["ArchiveDays"] = dataLoopTmp["archive_days"]
			localMaps := make([]interface{}, 0)
			localData1 := dataLoopTmp["tag_filters"]
			for _, dataLoop1 := range convertToInterfaceArray(localData1) {
				dataLoop1Tmp := dataLoop1.(map[string]interface{})
				dataLoop1Map := make(map[string]interface{})
				dataLoop1Map["Value"] = dataLoop1Tmp["value"]
				dataLoop1Map["Operator"] = dataLoop1Tmp["operator"]
				dataLoop1Map["Key"] = dataLoop1Tmp["key"]
				localMaps = append(localMaps, dataLoop1Map)
			}
			dataLoopMap["TagFilters"] = localMaps
			dataLoopMap["KeepLatestSnapshots"] = dataLoopTmp["keep_latest_snapshots"]
			dataLoopMap["VaultId"] = dataLoopTmp["vault_id"]
			localMaps1 := make([]interface{}, 0)
			localData2 := dataLoopTmp["data_source_filters"]
			for _, dataLoop2 := range convertToInterfaceArray(localData2) {
				dataLoop2Tmp := dataLoop2.(map[string]interface{})
				dataLoop2Map := make(map[string]interface{})
				dataLoop2Map["SourceType"] = dataLoop2Tmp["source_type"]
				localMaps1 = append(localMaps1, dataLoop2Map)
			}
			dataLoopMap["DataSourceFilters"] = localMaps1
			localMaps2 := make([]interface{}, 0)
			localData3 := dataLoopTmp["retention_rules"]
			for _, dataLoop3 := range convertToInterfaceArray(localData3) {
				dataLoop3Tmp := dataLoop3.(map[string]interface{})
				dataLoop3Map := make(map[string]interface{})
				dataLoop3Map["AdvancedRetentionType"] = dataLoop3Tmp["advanced_retention_type"]
				dataLoop3Map["Retention"] = dataLoop3Tmp["retention"]
				localMaps2 = append(localMaps2, dataLoop3Map)
			}
			dataLoopMap["RetentionRules"] = localMaps2
			dataLoopMap["RuleType"] = dataLoopTmp["rule_type"]
			dataLoopMap["Schedule"] = dataLoopTmp["schedule"]
			dataLoopMap["Retention"] = dataLoopTmp["retention"]
			dataLoopMap["ReplicationRegionId"] = dataLoopTmp["replication_region_id"]
			if backupType, ok := dataLoopTmp["backup_type"]; ok && backupType != "" {
				dataLoopMap["BackupType"] = dataLoopTmp["backup_type"]
			}
			rulesMapsArray = append(rulesMapsArray, dataLoopMap)
		}
		rulesMapsJson, err := json.Marshal(rulesMapsArray)
		if err != nil {
			return WrapError(err)
		}
		request["Rules"] = string(rulesMapsJson)
	}

	if v, ok := d.GetOk("policy_description"); ok {
		request["PolicyDescription"] = v
	}
	if v, ok := d.GetOk("policy_type"); ok {
		request["PolicyType"] = v
	}
	if v, ok := d.GetOk("policy_name"); ok {
		request["PolicyName"] = v
	}
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("hbr", "2017-09-08", action, query, request, true)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_hbr_policy", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["PolicyId"]))

	return resourceAliCloudHbrPolicyRead(d, meta)
}

func resourceAliCloudHbrPolicyRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	hbrServiceV2 := HbrServiceV2{client}

	objectRaw, err := hbrServiceV2.DescribeHbrPolicy(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_hbr_policy DescribeHbrPolicy Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("create_time", objectRaw["CreatedTime"])
	d.Set("policy_description", objectRaw["PolicyDescription"])
	d.Set("policy_name", objectRaw["PolicyName"])
	d.Set("policy_type", objectRaw["PolicyType"])

	rulesRaw := objectRaw["Rules"]
	rulesMaps := make([]map[string]interface{}, 0)
	if rulesRaw != nil {
		for _, rulesChildRaw := range convertToInterfaceArray(rulesRaw) {
			rulesMap := make(map[string]interface{})
			rulesChildRaw := rulesChildRaw.(map[string]interface{})
			rulesMap["archive_days"] = rulesChildRaw["ArchiveDays"]
			rulesMap["backup_type"] = rulesChildRaw["BackupType"]
			rulesMap["keep_latest_snapshots"] = rulesChildRaw["KeepLatestSnapshots"]
			rulesMap["replication_region_id"] = rulesChildRaw["ReplicationRegionId"]
			rulesMap["retention"] = rulesChildRaw["Retention"]
			rulesMap["rule_id"] = rulesChildRaw["RuleId"]
			rulesMap["rule_type"] = rulesChildRaw["RuleType"]
			rulesMap["schedule"] = rulesChildRaw["Schedule"]
			rulesMap["vault_id"] = rulesChildRaw["VaultId"]

			dataSourceFiltersRaw := rulesChildRaw["DataSourceFilters"]
			dataSourceFiltersMaps := make([]map[string]interface{}, 0)
			if dataSourceFiltersRaw != nil {
				for _, dataSourceFiltersChildRaw := range convertToInterfaceArray(dataSourceFiltersRaw) {
					dataSourceFiltersMap := make(map[string]interface{})
					dataSourceFiltersChildRaw := dataSourceFiltersChildRaw.(map[string]interface{})
					dataSourceFiltersMap["source_type"] = dataSourceFiltersChildRaw["SourceType"]

					dataSourceFiltersMaps = append(dataSourceFiltersMaps, dataSourceFiltersMap)
				}
			}
			rulesMap["data_source_filters"] = dataSourceFiltersMaps
			retentionRulesRaw := rulesChildRaw["RetentionRules"]
			retentionRulesMaps := make([]map[string]interface{}, 0)
			if retentionRulesRaw != nil {
				for _, retentionRulesChildRaw := range convertToInterfaceArray(retentionRulesRaw) {
					retentionRulesMap := make(map[string]interface{})
					retentionRulesChildRaw := retentionRulesChildRaw.(map[string]interface{})
					retentionRulesMap["advanced_retention_type"] = retentionRulesChildRaw["AdvancedRetentionType"]
					retentionRulesMap["retention"] = retentionRulesChildRaw["Retention"]

					retentionRulesMaps = append(retentionRulesMaps, retentionRulesMap)
				}
			}
			rulesMap["retention_rules"] = retentionRulesMaps
			tagFiltersRaw := rulesChildRaw["TagFilters"]
			tagFiltersMaps := make([]map[string]interface{}, 0)
			if tagFiltersRaw != nil {
				for _, tagFiltersChildRaw := range convertToInterfaceArray(tagFiltersRaw) {
					tagFiltersMap := make(map[string]interface{})
					tagFiltersChildRaw := tagFiltersChildRaw.(map[string]interface{})
					tagFiltersMap["key"] = tagFiltersChildRaw["Key"]
					tagFiltersMap["operator"] = tagFiltersChildRaw["Operator"]
					tagFiltersMap["value"] = tagFiltersChildRaw["Value"]

					tagFiltersMaps = append(tagFiltersMaps, tagFiltersMap)
				}
			}
			rulesMap["tag_filters"] = tagFiltersMaps
			rulesMaps = append(rulesMaps, rulesMap)
		}
	}
	if err := d.Set("rules", rulesMaps); err != nil {
		return err
	}

	return nil
}

func resourceAliCloudHbrPolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false

	var err error
	action := "UpdatePolicyV2"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["PolicyId"] = d.Id()

	if d.HasChange("rules") {
		update = true
		if v, ok := d.GetOk("rules"); ok || d.HasChange("rules") {
			rulesMapsArray := make([]interface{}, 0)
			for _, dataLoop := range convertToInterfaceArray(v) {
				dataLoopTmp := dataLoop.(map[string]interface{})
				dataLoopMap := make(map[string]interface{})
				dataLoopMap["ArchiveDays"] = dataLoopTmp["archive_days"]
				localMaps := make([]interface{}, 0)
				localData1 := dataLoopTmp["tag_filters"]
				for _, dataLoop1 := range convertToInterfaceArray(localData1) {
					dataLoop1Tmp := dataLoop1.(map[string]interface{})
					dataLoop1Map := make(map[string]interface{})
					dataLoop1Map["Value"] = dataLoop1Tmp["value"]
					dataLoop1Map["Operator"] = dataLoop1Tmp["operator"]
					dataLoop1Map["Key"] = dataLoop1Tmp["key"]
					localMaps = append(localMaps, dataLoop1Map)
				}
				dataLoopMap["TagFilters"] = localMaps
				dataLoopMap["KeepLatestSnapshots"] = dataLoopTmp["keep_latest_snapshots"]
				dataLoopMap["VaultId"] = dataLoopTmp["vault_id"]
				localMaps1 := make([]interface{}, 0)
				localData2 := dataLoopTmp["data_source_filters"]
				for _, dataLoop2 := range convertToInterfaceArray(localData2) {
					dataLoop2Tmp := dataLoop2.(map[string]interface{})
					dataLoop2Map := make(map[string]interface{})
					dataLoop2Map["SourceType"] = dataLoop2Tmp["source_type"]
					localMaps1 = append(localMaps1, dataLoop2Map)
				}
				dataLoopMap["DataSourceFilters"] = localMaps1
				localMaps2 := make([]interface{}, 0)
				localData3 := dataLoopTmp["retention_rules"]
				for _, dataLoop3 := range convertToInterfaceArray(localData3) {
					dataLoop3Tmp := dataLoop3.(map[string]interface{})
					dataLoop3Map := make(map[string]interface{})
					dataLoop3Map["AdvancedRetentionType"] = dataLoop3Tmp["advanced_retention_type"]
					dataLoop3Map["Retention"] = dataLoop3Tmp["retention"]
					localMaps2 = append(localMaps2, dataLoop3Map)
				}
				dataLoopMap["RetentionRules"] = localMaps2
				dataLoopMap["RuleType"] = dataLoopTmp["rule_type"]
				dataLoopMap["Schedule"] = dataLoopTmp["schedule"]
				dataLoopMap["Retention"] = dataLoopTmp["retention"]
				dataLoopMap["ReplicationRegionId"] = dataLoopTmp["replication_region_id"]
				if backupType, ok := dataLoopTmp["backup_type"]; ok && backupType != "" {
					dataLoopMap["BackupType"] = dataLoopTmp["backup_type"]
				}
				rulesMapsArray = append(rulesMapsArray, dataLoopMap)
			}
			rulesMapsJson, err := json.Marshal(rulesMapsArray)
			if err != nil {
				return WrapError(err)
			}
			request["Rules"] = string(rulesMapsJson)
		}
	}

	if d.HasChange("policy_description") {
		update = true
		request["PolicyDescription"] = d.Get("policy_description")
	}

	if d.HasChange("policy_name") {
		update = true
		request["PolicyName"] = d.Get("policy_name")
	}

	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("hbr", "2017-09-08", action, query, request, true)
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
	return resourceAliCloudHbrPolicyRead(d, meta)
}

func resourceAliCloudHbrPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DeletePolicyV2"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["PolicyId"] = d.Id()

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("hbr", "2017-09-08", action, query, request, true)
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

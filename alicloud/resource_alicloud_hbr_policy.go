// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"log"
	"time"

	"encoding/json"
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
							ValidateFunc: StringInSlice([]string{"BACKUP", "TRANSITION", "REPLICATION"}, false),
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

	if v, ok := d.GetOk("policy_name"); ok {
		request["PolicyName"] = v
	}
	if v, ok := d.GetOk("policy_description"); ok {
		request["PolicyDescription"] = v
	}
	if v, ok := d.GetOk("rules"); ok {
		rulesMapsArray := make([]interface{}, 0)
		for _, dataLoop := range v.([]interface{}) {
			dataLoopTmp := dataLoop.(map[string]interface{})
			dataLoopMap := make(map[string]interface{})
			dataLoopMap["Retention"] = dataLoopTmp["retention"]
			dataLoopMap["RuleType"] = dataLoopTmp["rule_type"]
			dataLoopMap["Schedule"] = dataLoopTmp["schedule"]
			dataLoopMap["ReplicationRegionId"] = dataLoopTmp["replication_region_id"]
			dataLoopMap["ArchiveDays"] = dataLoopTmp["archive_days"]
			localMaps := make([]interface{}, 0)
			localData1 := dataLoopTmp["retention_rules"]
			for _, dataLoop1 := range localData1.([]interface{}) {
				dataLoop1Tmp := dataLoop1.(map[string]interface{})
				dataLoop1Map := make(map[string]interface{})
				dataLoop1Map["AdvancedRetentionType"] = dataLoop1Tmp["advanced_retention_type"]
				dataLoop1Map["Retention"] = dataLoop1Tmp["retention"]
				localMaps = append(localMaps, dataLoop1Map)
			}
			dataLoopMap["RetentionRules"] = localMaps
			dataLoopMap["VaultId"] = dataLoopTmp["vault_id"]
			dataLoopMap["KeepLatestSnapshots"] = dataLoopTmp["keep_latest_snapshots"]
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

	if v, ok := d.GetOk("policy_type"); ok {
		request["PolicyType"] = v
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
		for _, rulesChildRaw := range rulesRaw.([]interface{}) {
			rulesMap := make(map[string]interface{})
			rulesChildRawArg := rulesChildRaw.(map[string]interface{})
			rulesMap["archive_days"] = rulesChildRawArg["ArchiveDays"]
			rulesMap["backup_type"] = rulesChildRawArg["BackupType"]
			rulesMap["keep_latest_snapshots"] = rulesChildRawArg["KeepLatestSnapshots"]
			rulesMap["replication_region_id"] = rulesChildRawArg["ReplicationRegionId"]
			rulesMap["retention"] = rulesChildRawArg["Retention"]
			rulesMap["rule_id"] = rulesChildRawArg["RuleId"]
			rulesMap["rule_type"] = rulesChildRawArg["RuleType"]
			rulesMap["schedule"] = rulesChildRawArg["Schedule"]
			rulesMap["vault_id"] = rulesChildRawArg["VaultId"]

			retentionRulesRaw := rulesChildRawArg["RetentionRules"]
			retentionRulesMaps := make([]map[string]interface{}, 0)
			if retentionRulesRaw != nil {
				for _, retentionRulesChildRaw := range retentionRulesRaw.([]interface{}) {
					retentionRulesMap := make(map[string]interface{})
					retentionRulesChildRawArg := retentionRulesChildRaw.(map[string]interface{})
					retentionRulesMap["advanced_retention_type"] = retentionRulesChildRawArg["AdvancedRetentionType"]
					retentionRulesMap["retention"] = retentionRulesChildRawArg["Retention"]

					retentionRulesMaps = append(retentionRulesMaps, retentionRulesMap)
				}
			}
			rulesMap["retention_rules"] = retentionRulesMaps
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

	if d.HasChange("policy_name") {
		update = true
		request["PolicyName"] = d.Get("policy_name")
	}

	if d.HasChange("policy_description") {
		update = true
		request["PolicyDescription"] = d.Get("policy_description")
	}

	if d.HasChange("rules") {
		update = true
		if v, ok := d.GetOk("rules"); ok {
			rulesMapsArray := make([]interface{}, 0)
			for _, dataLoop := range v.([]interface{}) {
				dataLoopTmp := dataLoop.(map[string]interface{})
				dataLoopMap := make(map[string]interface{})
				dataLoopMap["RuleType"] = dataLoopTmp["rule_type"]
				dataLoopMap["Schedule"] = dataLoopTmp["schedule"]
				dataLoopMap["Retention"] = dataLoopTmp["retention"]
				dataLoopMap["ReplicationRegionId"] = dataLoopTmp["replication_region_id"]
				dataLoopMap["KeepLatestSnapshots"] = dataLoopTmp["keep_latest_snapshots"]
				dataLoopMap["ArchiveDays"] = dataLoopTmp["archive_days"]
				localMaps := make([]interface{}, 0)
				localData1 := dataLoopTmp["retention_rules"]
				for _, dataLoop1 := range localData1.([]interface{}) {
					dataLoop1Tmp := dataLoop1.(map[string]interface{})
					dataLoop1Map := make(map[string]interface{})
					dataLoop1Map["AdvancedRetentionType"] = dataLoop1Tmp["advanced_retention_type"]
					dataLoop1Map["Retention"] = dataLoop1Tmp["retention"]
					localMaps = append(localMaps, dataLoop1Map)
				}
				dataLoopMap["RetentionRules"] = localMaps
				dataLoopMap["VaultId"] = dataLoopTmp["vault_id"]
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

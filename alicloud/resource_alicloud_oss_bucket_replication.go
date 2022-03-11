package alicloud

import (
	"encoding/xml"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudOssBucketReplication() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudOssBucketReplicationCreate,
		Read:   resourceAlicloudOssBucketReplicationRead,
		Delete: resourceAlicloudOssBucketReplicationDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"bucket": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"rule_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"prefix_set": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"prefixes": {
							Type:     schema.TypeList,
							Required: true,
							ForceNew: true,
							MaxItems: 10,
							Elem:     &schema.Schema{Type: schema.TypeString, ValidateFunc: validation.StringLenBetween(0, 1023)},
						},
					},
				},
			},
			"action": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"destination": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"bucket": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"location": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"transfer_type": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
					},
				},
			},
			"historical_object_replication": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"enabled",
					"disabled",
				}, false),
			},
			"sync_role": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"source_selection_criteria": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"sse_kms_encrypted_objects": {
							Type:     schema.TypeList,
							Optional: true,
							ForceNew: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"status": {
										Type:     schema.TypeString,
										Optional: true,
										ForceNew: true,
										ValidateFunc: validation.StringInSlice([]string{
											"Enabled",
											"Disabled",
										}, false),
									},
								},
							},
						},
					},
				},
			},
			"encryption_configuration": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"replica_kms_key_id": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
					},
				},
			},
			"progress": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"historical_object": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"new_object": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

type PrefixSetType struct {
	Prefixes []string `xml:"Prefix"`
}

type DestinationType struct {
	Bucket       string `xml:"Bucket"`
	Location     string `xml:"Location"`
	TransferType string `xml:"TransferType,omitempty"`
}

type SseKmsEncryptedObjectsType struct {
	Status string `xml:"Status"`
}

type SourceSelectionCriteriaType struct {
	SseKmsEncryptedObjects *SseKmsEncryptedObjectsType `xml:"SseKmsEncryptedObjects,omitempty"`
}

type EncryptionConfigurationType struct {
	ReplicaKmsKeyID string `xml:"ReplicaKmsKeyID,omitempty"`
}

type ReplicationRule struct {
	ID                          string                       `xml:"ID,omitempty"`
	Action                      string                       `xml:"Action,omitempty"`
	PrefixSet                   *PrefixSetType               `xml:"PrefixSet,omitempty"`
	Destination                 *DestinationType             `xml:"Destination"`
	HistoricalObjectReplication string                       `xml:"HistoricalObjectReplication,omitempty"`
	Status                      string                       `xml:"Status,omitempty"`
	SyncRole                    string                       `xml:"SyncRole,omitempty"`
	SourceSelectionCriteria     *SourceSelectionCriteriaType `xml:"SourceSelectionCriteria,omitempty"`
	EncryptionConfiguration     *EncryptionConfigurationType `xml:"EncryptionConfiguration,omitempty"`
}

type ReplicationConfiguration struct {
	XMLName xml.Name          `xml:"ReplicationConfiguration"`
	Rules   []ReplicationRule `xml:"Rule"`
}

type ReplicationProgress struct {
	XMLName          xml.Name `xml:"ReplicationProgress"`
	ID               string   `xml:"Rule>ID"`
	HistoricalObject string   `xml:"Rule>Progress>HistoricalObject"`
	NewObject        string   `xml:"Rule>Progress>NewObject"`
}

func expandReplicationRule(d *schema.ResourceData) ReplicationRule {
	r := d.Get("").(map[string]interface{})
	rule := ReplicationRule{}
	// ID
	if val, ok := r["rule_id"].(string); ok && val != "" {
		rule.ID = val
	}

	// Action
	if val, ok := r["action"].(string); ok && val != "" {
		rule.Action = val
	}

	// Status
	if val, ok := r["status"].(string); ok && val != "" {
		rule.Status = val
	}

	// HistoricalObjectReplication
	if val, ok := r["historical_object_replication"].(string); ok && val != "" {
		rule.HistoricalObjectReplication = val
	}

	// SyncRole
	if val, ok := r["sync_role"].(string); ok && val != "" {
		rule.SyncRole = val
	}

	// PrefixSet
	if val, ok := r["prefix_set"].([]interface{}); ok && len(val) > 0 && val[0] != nil {
		e := val[0].(map[string]interface{})
		i := PrefixSetType{}
		if v, ok := e["prefixes"]; ok {
			var prefixes []string
			for _, prefix := range v.([]interface{}) {
				prefixes = append(prefixes, prefix.(string))
			}
			i.Prefixes = prefixes
		}
		rule.PrefixSet = &i
	}

	// Destination
	if val, ok := r["destination"].([]interface{}); ok && len(val) > 0 && val[0] != nil {
		e := val[0].(map[string]interface{})
		i := DestinationType{}
		// Bucket
		if val, ok := e["bucket"].(string); ok && val != "" {
			i.Bucket = val
		}

		// Location
		if val, ok := e["location"].(string); ok && val != "" {
			i.Location = val
		}

		// TransferType
		if val, ok := e["transfer_type"].(string); ok && val != "" {
			i.TransferType = val
		}
		rule.Destination = &i
	}

	// SourceSelectionCriteria
	if val, ok := r["source_selection_criteria"].([]interface{}); ok && len(val) > 0 && val[0] != nil {
		e := val[0].(map[string]interface{})
		i := SourceSelectionCriteriaType{}
		//SseKmsEncryptedObjects
		if val1, ok := e["sse_kms_encrypted_objects"].([]interface{}); ok && len(val1) > 0 && val1[0] != nil {
			s := val1[0].(map[string]interface{})
			j := SseKmsEncryptedObjectsType{}
			if v, ok := s["status"].(string); ok && v != "" {
				j.Status = v
			}
			i.SseKmsEncryptedObjects = &j
		}
		rule.SourceSelectionCriteria = &i
	}

	// EncryptionConfiguration
	if val, ok := r["encryption_configuration"].([]interface{}); ok && len(val) > 0 && val[0] != nil {
		e := val[0].(map[string]interface{})
		i := EncryptionConfigurationType{}
		// ReplicaKmsKeyID
		if val, ok := e["replica_kms_key_id"].(string); ok && val != "" {
			i.ReplicaKmsKeyID = val
		}
		rule.EncryptionConfiguration = &i
	}

	return rule
}

func flattenReplicationRule(d *schema.ResourceData, rc *ReplicationConfiguration, rp map[string]interface{}, ruleId string) error {
	rule := make(map[string]interface{})
	for _, r := range rc.Rules {
		if r.ID != ruleId {
			continue
		}

		// ID
		if r.ID != "" {
			rule["rule_id"] = r.ID
		}

		// Action
		if r.Action != "" {
			rule["action"] = r.Action
		}

		// Status
		if r.Status != "" {
			rule["status"] = r.Status
		}

		// HistoricalObjectReplication
		if r.HistoricalObjectReplication != "" {
			rule["historical_object_replication"] = r.HistoricalObjectReplication
		}

		// SyncRole
		if r.SyncRole != "" {
			rule["sync_role"] = r.SyncRole
		}

		// PrefixSet
		if r.PrefixSet != nil {
			prefixSet := make(map[string]interface{})
			// prefixes
			if len(r.PrefixSet.Prefixes) != 0 {
				prefixSet["prefixes"] = r.PrefixSet.Prefixes
			}
			rule["prefix_set"] = []interface{}{prefixSet}
		}

		// Destination
		if r.Destination != nil {
			destination := make(map[string]interface{})
			// Bucket
			if r.Destination.Bucket != "" {
				destination["bucket"] = r.Destination.Bucket
			}
			// Location
			if r.Destination.Location != "" {
				destination["location"] = r.Destination.Location
			}
			// TransferType
			if r.Destination.TransferType != "" {
				destination["transfer_type"] = r.Destination.TransferType
			}
			rule["destination"] = []interface{}{destination}
		}

		// SourceSelectionCriteria
		if r.SourceSelectionCriteria != nil {
			sourceSelectionCriteria := make(map[string]interface{})
			if r.SourceSelectionCriteria.SseKmsEncryptedObjects != nil {
				sseKmsEncryptedObjects := make(map[string]interface{})
				if r.SourceSelectionCriteria.SseKmsEncryptedObjects.Status != "" {
					sseKmsEncryptedObjects["status"] = r.SourceSelectionCriteria.SseKmsEncryptedObjects.Status
				}
				sourceSelectionCriteria["sse_kms_encrypted_objects"] = []interface{}{sseKmsEncryptedObjects}
			}
			rule["source_selection_criteria"] = []interface{}{sourceSelectionCriteria}
		}

		// EncryptionConfiguration
		if r.EncryptionConfiguration != nil {
			encryptionConfiguration := make(map[string]interface{})
			if r.EncryptionConfiguration.ReplicaKmsKeyID != "" {
				encryptionConfiguration["replica_kms_key_id"] = r.EncryptionConfiguration.ReplicaKmsKeyID
			}
			rule["encryption_configuration"] = []interface{}{encryptionConfiguration}
		}

		//Progress
		if rp != nil {
			if val, ok := rp[r.ID].(*ReplicationProgress); ok && val != nil {
				progress := make(map[string]interface{})
				if val.HistoricalObject != "" {
					progress["historical_object"] = val.HistoricalObject
				}

				if val.NewObject != "" {
					progress["new_object"] = val.NewObject
				}
				rule["progress"] = []interface{}{progress}
			}
		}
	}

	for k, v := range rule {
		if err := d.Set(k, v); err != nil {
			return err
		}
	}

	return nil
}

func retrieveReplicationRules(client *connectivity.AliyunClient, bucket string) (*ReplicationConfiguration, error) {
	// Read the replication configuration
	var requestInfo *oss.Client
	raw, err := client.WithOssClient(func(ossClient *oss.Client) (interface{}, error) {
		requestInfo = ossClient
		return ossClient.GetBucketReplication(bucket)
	})

	if err != nil {
		return nil, err
	}

	addDebug("GetBucketReplication", raw, requestInfo, map[string]interface{}{
		"bucketName": bucket,
	})

	var rc ReplicationConfiguration
	var body []byte = []byte(raw.(string))
	if err := xml.Unmarshal(body, &rc); err != nil {
		return nil, err
	}
	return &rc, nil
}

func isReplicationRuleExist(client *connectivity.AliyunClient, bucket string, ruleId string) (bool, error) {
	// Read the replication Progress by rule id
	_, err := client.WithOssClient(func(ossClient *oss.Client) (interface{}, error) {
		return ossClient.GetBucketReplicationProgress(bucket, ruleId)
	})

	if err == nil {
		return true, nil
	}

	if IsExpectedErrors(err, []string{"NoSuchReplicationRule", "NoSuchBucket"}) {
		return false, nil
	}

	return false, err
}

func retrieveReplicationRuleProgress(client *connectivity.AliyunClient, bucket string, ruleId string) (ReplicationProgress, error) {
	var rp ReplicationProgress
	var requestInfo *oss.Client

	// Read the replication Progress by rule id
	raw, err := client.WithOssClient(func(ossClient *oss.Client) (interface{}, error) {
		requestInfo = ossClient
		return ossClient.GetBucketReplicationProgress(bucket, ruleId)
	})

	if err != nil {
		return rp, err
	}

	addDebug("GetBucketReplicationProgress", raw, requestInfo, map[string]interface{}{
		"bucketName": bucket,
		"ruleId":     ruleId,
	})

	var body []byte = []byte(raw.(string))
	if err := xml.Unmarshal(body, &rp); err != nil {
		return rp, WrapErrorf(err, DefaultErrorMsg, bucket, "Unmarshal XML", AliyunOssGoSdk)
	}

	return rp, err
}

func hasProgressBlock(d *schema.ResourceData) bool {
	return true
}

func resourceAlicloudOssBucketReplicationCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	bucket := d.Get("bucket").(string)

	rules := make([]ReplicationRule, 0)
	rule := expandReplicationRule(d)
	rc := &ReplicationConfiguration{
		Rules: append(rules, rule),
	}

	bs, err := xml.Marshal(rc)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, bucket, "Marshal to XML", AliyunOssGoSdk)
	}
	var xmlString string = string(bs[:])
	var requestInfo *oss.Client
	var replicationAlreadyExist bool

	replicationAlreadyExist = false
	raw, err := client.WithOssClient(func(ossClient *oss.Client) (interface{}, error) {
		requestInfo = ossClient
		return nil, ossClient.PutBucketReplication(bucket, xmlString)
	})
	if err != nil {
		if !IsExpectedErrors(err, []string{"BucketReplicationAlreadyExist"}) {
			return WrapErrorf(err, DefaultErrorMsg, bucket, "PutBucketReplication", AliyunOssGoSdk)
		}
		replicationAlreadyExist = true
	}
	addDebug("PutBucketReplication", raw, requestInfo, map[string]interface{}{
		"bucketName":               bucket,
		"ReplicationConfiguration": xmlString,
		"ReplicationAlreadyExist":  replicationAlreadyExist,
	})

	//OSS server does not return rule-id and only supports one rule currently, obtains rule id through GetBucketReplication
	rc, err = retrieveReplicationRules(client, bucket)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, bucket, "retrieveReplicationRules", AliyunOssGoSdk)
	}

	d.SetId(fmt.Sprintf("%s%s%s", bucket, COLON_SEPARATED, rc.Rules[0].ID))
	return resourceAlicloudOssBucketReplicationRead(d, meta)
}

func resourceAlicloudOssBucketReplicationRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	bucket := parts[0]
	ruleId := parts[1]

	// Read the replication configuration
	rc, err := retrieveReplicationRules(client, bucket)

	if IsExpectedErrors(err, []string{"NoSuchReplicationConfiguration", "NoSuchBucket"}) {
		log.Printf("[WARN] OSS Bucket Replication Configuration (%s) not found, removing from state", bucket)
		d.SetId("")
		return nil
	}

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, bucket, "retrieveReplicationRules", AliyunOssGoSdk)
	}

	//Read the replication progress by rule id
	rule_progress := make(map[string]interface{})
	if hasProgressBlock(d) {
		if val, err := retrieveReplicationRuleProgress(client, bucket, ruleId); err == nil {
			rule_progress[ruleId] = &val
		}
	}

	// Update to resouce data
	d.Set("bucket", bucket)
	err = flattenReplicationRule(d, rc, rule_progress, ruleId)
	if err != nil {
		return WrapError(err)
	}
	return nil
}

func resourceAlicloudOssBucketReplicationDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	bucket := parts[0]
	ruleId := parts[1]

	var requestInfo *oss.Client
	raw, err := client.WithOssClient(func(ossClient *oss.Client) (interface{}, error) {
		requestInfo = ossClient
		return nil, ossClient.DeleteBucketReplication(bucket, ruleId)
	})

	if IsExpectedErrors(err, []string{"NoSuchReplicationConfiguration", "NoSuchBucket"}) {
		return nil
	}

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, bucket, "DeleteBucketReplication", AliyunOssGoSdk)
	}
	addDebug("DeleteBucketReplication", raw, requestInfo, map[string]interface{}{
		"bucketName": bucket,
		"ruleId":     ruleId,
	})

	// wait until the replication configuration is closed
	_ = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		raw, _ := isReplicationRuleExist(client, bucket, ruleId)
		if raw {
			time.Sleep(time.Duration(10) * time.Second)
			return resource.RetryableError(Error("in closing status"))
		}
		return nil
	})

	return nil
}

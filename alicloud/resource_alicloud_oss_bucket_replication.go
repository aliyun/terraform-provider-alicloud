package alicloud

import (
	"bytes"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/hashcode"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudOssBucketReplication() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudOssBucketReplicationCreate,
		Read:   resourceAlicloudOssBucketReplicationRead,
		Update: resourceAlicloudOssBucketReplicationUpdate,
		Delete: resourceAlicloudOssBucketReplicationDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"bucket": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(3, 63),
				Default:      resource.PrefixedUniqueId("tf-oss-bucket-replication-"),
			},
			"replication_rule": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"prefix_set": {
							Type:     schema.TypeSet,
							Optional: true,
							Set:      prefixSetHash,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"prefixes": {
										Type:     schema.TypeList,
										Optional: true,
										Elem:     &schema.Schema{Type: schema.TypeString, ValidateFunc: validation.StringLenBetween(0, 1023)},
										MaxItems: 10,
									},
								},
							},
						},
						"action": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  oss.ReplicationActionALL,
							ValidateFunc: validation.StringInSlice([]string{
								string(oss.ReplicationActionALL),
								string(oss.ReplicationActionPUT),
								string(oss.ReplicationActionDELETE),
								string(oss.ReplicationActionABORT),
								string(oss.ReplicationActionPUT_ABORT),
								string(oss.ReplicationActionPUT_DELETE),
								string(oss.ReplicationActionDELETE_ABORT),
							}, false),
						},
						"destination": {
							Type:     schema.TypeSet,
							Optional: true,
							Set:      destinationHash,
							ForceNew: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"bucket": {
										Type:     schema.TypeString,
										Optional: true,
										ForceNew: true,
									},
									"location": {
										Type:     schema.TypeString,
										Optional: true,
										ForceNew: true,
									},
									"transfer_type": {
										Type:     schema.TypeString,
										Optional: true,
										ValidateFunc: validation.StringInSlice([]string{
											string(oss.TransferInternal),
											string(oss.TransferOssAcc),
										}, false),
									},
								},
							},
						},
						"historical_object_replication": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  oss.HistoricalEnabled,
							ValidateFunc: validation.StringInSlice([]string{
								string(oss.HistoricalEnabled),
								string(oss.HistoricalDisabled),
							}, false),
						},
						"sync_role": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"source_selection_criteria": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"sse_kms_encrypted_objects": {
										Type:     schema.TypeSet,
										Optional: true,
										Set:      sseKmsEncryptedObjectsHash,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"status": {
													Type:     schema.TypeString,
													Optional: true,
													ValidateFunc: validation.StringInSlice([]string{
														string(oss.SourceSSEKMSEnabled),
														string(oss.SourceSSEKMSDisabled),
													}, false),
												},
											},
										},
									},
								},
							},
						},
						"encryption_configuration": {
							Type:     schema.TypeSet,
							Optional: true,
							Set:      encryptionConfigurationHash,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"replica_kms_key_id": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func resourceAlicloudOssBucketReplicationCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	request := map[string]string{"bucketName": d.Get("bucket").(string)}

	var requestInfo *oss.Client
	raw, err := client.WithOssClient(func(ossClient *oss.Client) (interface{}, error) {
		requestInfo = ossClient
		return ossClient.IsBucketExist(request["bucketName"])
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_oss_bucket_replication", "IsBucketExist", AliyunOssGoSdk)
	}
	addDebug("IsBucketExist", raw, requestInfo, request)

	isExist, _ := raw.(bool)
	if isExist {
		return resourceAlicloudOssBucketReplicationUpdate(d, meta)
		//return WrapError(Error("[ERROR] The specified bucket name: %#v is not available. The bucket namespace is shared by all users of the OSS system. Please select a different name and try again.", request["bucketName"]))
	}

	for {
		raw, err = client.WithOssClient(func(ossClient *oss.Client) (interface{}, error) {
			return nil, ossClient.CreateBucket(request["bucketName"])
		})
		if BucketAlreadyExistsError(err) {
			log.Printf("[DEBUG] Bucket: %s still exists, wait 60s and retry...", d.Id())
			time.Sleep(60 * time.Second)
			continue
		} else if err == nil {
			break
		} else {
			return WrapErrorf(err, DefaultErrorMsg, "alicloud_oss_bucket_replication", "CreateBucket", AliyunOssGoSdk)
		}
	}

	addDebug("CreateBucket", raw, requestInfo, request["bucketName"])
	err = resource.Retry(3*time.Minute, func() *resource.RetryError {
		raw, err = client.WithOssClient(func(ossClient *oss.Client) (interface{}, error) {
			return ossClient.IsBucketExist(request["bucketName"])
		})

		if err != nil {
			return resource.NonRetryableError(err)
		}
		isExist, _ := raw.(bool)
		if !isExist {
			return resource.RetryableError(Error("Trying to ensure new OSS bucket %#v has been created successfully.", request["bucketName"]))
		}
		addDebug("IsBucketExist", raw, requestInfo, request)
		return nil
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_oss_bucket_replication", "IsBucketExist", AliyunOssGoSdk)
	}

	// Assign the bucket name as the resource ID
	d.SetId(request["bucketName"])
	return resourceAlicloudOssBucketReplicationUpdate(d, meta)
}

func resourceAlicloudOssBucketReplicationRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ossService := OssService{client}

	if d.Id() == "" {
		d.SetId("")
		d.Set("bucket", d.Id())
		d.Set("replication_rule", []map[string]interface{}{})
		return nil
	}

	replication, err := ossService.DescribeOssBucketReplication(d.Id())
	if err != nil {
		if NoSuchReplicationConfigurationError(err) {
			fmt.Printf(">>> NoSuchReplicationConfigurationError \n")
			d.SetId(d.Id())
			d.Set("bucket", d.Id())
			return nil
		}
		return WrapError(err)
	}
	d.Set("bucket", d.Id())

	repliRule := make([]map[string]interface{}, 0)

	if replication.Rule.Status != "" && replication.Rule.ID != "" {
		rule := make(map[string]interface{})
		rule["id"] = replication.Rule.ID
		rule["action"] = replication.Rule.Action
		rule["status"] = replication.Rule.Status
		if replication.Rule.HistoricalObjectReplication != "" {
			rule["historical_object_replication"] = replication.Rule.HistoricalObjectReplication
		}
		if replication.Rule.SyncRole != "" {
			rule["sync_role"] = replication.Rule.SyncRole
		}
		if replication.Rule.Destination != nil {
			d := make(map[string]interface{})
			d["bucket"] = replication.Rule.Destination.Bucket
			d["location"] = replication.Rule.Destination.Location
			if replication.Rule.Destination.TransferType == "" {
				replication.Rule.Destination.TransferType = string(oss.TransferInternal)
			}
			d["transfer_type"] = replication.Rule.Destination.TransferType

			rule["destination"] = schema.NewSet(destinationHash, []interface{}{d})
		}
		if replication.Rule.PrefixSet != nil {
			if len(replication.Rule.PrefixSet.Prefix) != 0 {
				p := make(map[string]interface{})
				p["prefixes"] = replication.Rule.PrefixSet.Prefix
				rule["prefix_set"] = schema.NewSet(prefixSetHash, []interface{}{p})
			}
		}
		if replication.Rule.SourceSelectionCriteria != nil { // todo 发现公有云replication功能不反回以下两个参数，待查清情况再继续
			if replication.Rule.SourceSelectionCriteria.SseKmsEncryptedObjects != nil {
				if replication.Rule.SourceSelectionCriteria.SseKmsEncryptedObjects.Status != "" {
					s := make(map[string]interface{})
					x := make(map[string]interface{})
					x["status"] = replication.Rule.SourceSelectionCriteria.SseKmsEncryptedObjects.Status
					s["sse_kms_encrypted_objects"] = schema.NewSet(sseKmsEncryptedObjectsHash, []interface{}{x})
					rule["source_selection_criteria"] = []interface{}{s}
				}
			}
		}
		if replication.Rule.EncryptionConfiguration != nil {
			if replication.Rule.EncryptionConfiguration.ReplicaKmsKeyID != "" {
				e := make(map[string]interface{})
				e["replica_kms_key_id"] = replication.Rule.EncryptionConfiguration.ReplicaKmsKeyID
				rule["encryption_configuration"] = schema.NewSet(encryptionConfigurationHash, []interface{}{e})
			}
		}

		repliRule = append(repliRule, rule)
	}
	if err := d.Set("replication_rule", repliRule); err != nil {
		return WrapError(err)
	}
	return nil
}

func resourceAlicloudOssBucketReplicationUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	d.Partial(true)

	// Request OSS according to Replication_rule in .tf file
	if d.HasChange("replication_rule") {
		_, n := d.GetChange("replication_rule")
		// The Replication interface does not support overwriting operations. Therefore, once the new Replication_rule
		// is not empty (that is, deleting replication), you need to delete the previous replication and then create
		// replication again. And length 16 means empty
		if len(fmt.Sprintf("%#v", n)) > 16 {
			err := getThenDeleteReplicationRule(client, d)
			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), "GetThenDeleteReplicationRule", AliyunOssGoSdk)
			}
		}
		if err := replicationUpdate(client, d); err != nil {
			return WrapError(err)
		}
		d.SetPartial("replication_rule")
	}

	d.Partial(false)

	return resourceAlicloudOssBucketReplicationRead(d, meta)
}

func resourceAlicloudOssBucketReplicationDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var requestInfo *oss.Client
	raw, err := client.WithOssClient(func(ossClient *oss.Client) (interface{}, error) {
		requestInfo = ossClient
		return ossClient.IsBucketExist(d.Id())
	})
	addDebug("IsBucketExist", raw, requestInfo, map[string]string{"bucketName": d.Id()})

	exist, _ := raw.(bool)
	if !exist {
		return nil
	}

	err = getThenDeleteReplicationRule(client, d)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), "GetThenDeleteReplicationRule", AliyunOssGoSdk)
	}

	raw, err = client.WithOssClient(func(ossClient *oss.Client) (interface{}, error) {
		return nil, ossClient.DeleteBucket(d.Id())
	})

	return nil
}

func getThenDeleteReplicationRule(client *connectivity.AliyunClient, d *schema.ResourceData) error {
	// The DELETE Replication interface must provide the ID of the replication rule, you can only get and then delete
	var requestInfo *oss.Client
	var repId string

	raw, err := client.WithOssClient(func(ossClient *oss.Client) (interface{}, error) {
		requestInfo = ossClient
		return ossClient.GetBucketReplication(d.Id())
	})
	if err == nil {
	} else if err != nil && NoSuchReplicationConfigurationError(err) {
		return nil
	} else {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), "GetBucketReplication", AliyunOssGoSdk)
	}
	addDebug("GetBucketReplication", raw, requestInfo, map[string]string{"bucketName": d.Id()})

	repRule, _ := raw.(oss.GetBucketReplicationResult)
	repId = repRule.Rule.ID

	if len(repId) != 36 {
		return nil
	}

	raw, err = client.WithOssClient(func(ossClient *oss.Client) (interface{}, error) {
		requestInfo = ossClient
		return nil, ossClient.DeleteBucketReplication(d.Id(), repId)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), "DeleteBucketReplication", AliyunOssGoSdk)
	}
	addDebug("DeleteBucketReplication", raw, requestInfo, map[string]string{"bucketName": d.Id()})

	for {
		raw, err = client.WithOssClient(func(ossClient *oss.Client) (interface{}, error) {
			return ossClient.GetBucketReplication(d.Id())
		})
		if NoSuchReplicationConfigurationError(err) {
			break
		} else if err == nil {
			log.Printf("[DEBUG] replication_rule of Bucket: %s still exists, wait 60s and retry...", d.Id())
			time.Sleep(60 * time.Second)
			continue
		} else {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), "GetBucketReplication", AliyunOssGoSdk)
		}
	}

	return nil
}

func replicationUpdate(client *connectivity.AliyunClient, d *schema.ResourceData) error {
	// Request the PUT Replication interface based on replication_rule in the .tf file
	replication_rule := d.Get("replication_rule").([]interface{})
	var requestInfo *oss.Client

	if replication_rule == nil || len(replication_rule) == 0 {
		err := getThenDeleteReplicationRule(client, d)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), "GetThenDeleteReplicationRule", AliyunOssGoSdk)
		}
	}

	c := replication_rule[0].(map[string]interface{})
	var repRule oss.ReplicationRule

	// PrefixSet
	prefixSet := d.Get("replication_rule.0.prefix_set").(*schema.Set).List()
	if len(prefixSet) > 0 {
		x := prefixSet[0].(map[string]interface{})
		s := oss.PrefixSet{}

		if v, ok := x["prefixes"]; ok {
			for _, prefix := range v.([]interface{}) {
				s.Prefix = append(s.Prefix, prefix.(string))
			}
			repRule.PrefixSet = &s
		}
	}

	// Action
	if val, ok := c["action"].(string); ok && val != "" {
		repRule.Action = val
	}

	// Destination
	destination := d.Get("replication_rule.0.destination").(*schema.Set).List()
	if len(destination) > 0 {
		e := destination[0].(map[string]interface{})
		i := oss.Destination{}

		if val, ok := e["bucket"].(string); ok && val != "" {
			i.Bucket = val
		}
		if val, ok := e["location"].(string); ok && val != "" {
			i.Location = val
		}
		if val, ok := e["transfer_type"].(string); ok && val != "" {
			i.TransferType = val
		}

		repRule.Destination = &i
	}

	// HistoricalObjectReplication
	if val, ok := c["historical_object_replication"].(string); ok && val != "" {
		repRule.HistoricalObjectReplication = val
	}

	// SyncRole
	if val, ok := c["sync_role"].(string); ok && val != "" {
		repRule.SyncRole = val
	}

	// SourceSelectionCriteria
	var sseKmsStatus string
	sseKmsEncryptedObjects := d.Get("replication_rule.0.source_selection_criteria.0.sse_kms_encrypted_objects").(*schema.Set).List()
	if len(sseKmsEncryptedObjects) > 0 {
		e := sseKmsEncryptedObjects[0].(map[string]interface{})

		if val, ok := e["status"].(string); ok && val != "" {
			repRule.SourceSelectionCriteria = &oss.SourceSelectionCriteria{
				SseKmsEncryptedObjects: &oss.SseKmsEncryptedObjects{
					Status: val,
				},
			}
		}
	}

	// EncryptionConfiguration
	encryptionConfiguration := d.Get("replication_rule.0.encryption_configuration").(*schema.Set).List()
	if len(encryptionConfiguration) > 0 {
		e := encryptionConfiguration[0].(map[string]interface{})

		if val, ok := e["replica_kms_key_id"].(string); ok && val != "" {
			repRule.EncryptionConfiguration = &oss.EncryptionConfiguration{
				ReplicaKmsKeyID: val,
			}
		} else {
			if sseKmsStatus == string(oss.SourceSSEKMSEnabled) {
				return WrapError(Error("if 'SourceSelectionCriteria.SseKmsEncryptedObjects.Status' was defined as 'Enabled', then 'EncryptionConfiguration' and 'ReplicaKmsKeyID' should be defined."))
			}
		}

	} else {
		if sseKmsStatus == string(oss.SourceSSEKMSEnabled) {
			return WrapError(Error("if 'SourceSelectionCriteria.SseKmsEncryptedObjects.Status' was defined as 'Enabled', then 'EncryptionConfiguration' and 'ReplicaKmsKeyID' should be defined."))
		}
	}

	raw, err := client.WithOssClient(func(ossClient *oss.Client) (interface{}, error) {
		requestInfo = ossClient
		return nil, ossClient.SetBucketReplication(d.Id(), repRule)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), "SetBucketReplication", AliyunOssGoSdk)
	}
	addDebug("SetBucketReplication", raw, requestInfo, map[string]interface{}{
		"bucketName":     d.Id(),
		"encryptionRule": repRule,
	})

	return nil
}

func prefixSetHash(v interface{}) int {
	var buf bytes.Buffer
	m := v.(map[string]interface{})

	if v, ok := m["prefixes"]; ok {
		buf.WriteString(fmt.Sprintf("%v-", v))
	}

	return hashcode.String(buf.String())
}

func destinationHash(v interface{}) int {
	var buf bytes.Buffer
	m := v.(map[string]interface{})

	if v, ok := m["bucket"]; ok {
		buf.WriteString(fmt.Sprintf("%s-", v.(string)))
	}
	if v, ok := m["location"]; ok {
		buf.WriteString(fmt.Sprintf("%s-", v.(string)))
	}
	if v, ok := m["transfer_type"]; ok {
		buf.WriteString(fmt.Sprintf("%s-", v.(string)))
	}

	return hashcode.String(buf.String())
}

func sseKmsEncryptedObjectsHash(v interface{}) int {
	var buf bytes.Buffer
	m := v.(map[string]interface{})

	if v, ok := m["status"]; ok {
		buf.WriteString(fmt.Sprintf("%s-", v.(string)))
	}

	return hashcode.String(buf.String())
}

func encryptionConfigurationHash(v interface{}) int {
	var buf bytes.Buffer
	m := v.(map[string]interface{})

	if v, ok := m["replica_kms_key_id"]; ok {
		buf.WriteString(fmt.Sprintf("%s-", v.(string)))
	}

	return hashcode.String(buf.String())
}

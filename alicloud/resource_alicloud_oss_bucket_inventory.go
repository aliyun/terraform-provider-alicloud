// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudOssBucketInventory() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudOssBucketInventoryCreate,
		Read:   resourceAliCloudOssBucketInventoryRead,
		Update: resourceAliCloudOssBucketInventoryUpdate,
		Delete: resourceAliCloudOssBucketInventoryDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"bucket": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"destination": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"oss_bucket_destination": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"format": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"account_id": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"bucket": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"prefix": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"encryption": {
										Type:     schema.TypeList,
										Optional: true,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"ssekms": {
													Type:     schema.TypeList,
													Optional: true,
													MaxItems: 1,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"key_id": {
																Type:     schema.TypeString,
																Optional: true,
															},
														},
													},
												},
												"sseoss": {
													Type:     schema.TypeString,
													Optional: true,
												},
											},
										},
									},
									"role_arn": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
					},
				},
			},
			"filter": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"lower_size_bound": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"storage_class": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"last_modify_begin_time_stamp": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"last_modify_end_time_stamp": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"upper_size_bound": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"prefix": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"included_object_versions": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"incremental_inventory": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"optional_fields": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"field": {
										Type:     schema.TypeList,
										Optional: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
						"is_enabled": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"schedule": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"frequency": {
										Type:     schema.TypeInt,
										Optional: true,
									},
								},
							},
						},
					},
				},
			},
			"inventory_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"is_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"optional_fields": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"field": {
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"schedule": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"frequency": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
		},
	}
}

func resourceAliCloudOssBucketInventoryCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := fmt.Sprintf("/?inventory")
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	body := make(map[string]interface{})
	hostMap := make(map[string]*string)
	var err error
	request = make(map[string]interface{})
	hostMap["bucket"] = StringPointer(d.Get("bucket").(string))
	query["inventoryId"] = StringPointer(d.Get("inventory_id").(string))

	inventoryConfiguration := make(map[string]interface{})
	if v := d.Get("optional_fields"); !IsNil(v) {
		optionalFields := make(map[string]interface{})
		field1, _ := jsonpath.Get("$[0].field", d.Get("optional_fields"))
		if field1 != nil && field1 != "" {
			optionalFields["Field"] = field1
		}

		if len(optionalFields) > 0 {
			inventoryConfiguration["OptionalFields"] = optionalFields
		}
	}
	if v, ok := d.GetOk("is_enabled"); ok {
		inventoryConfiguration["IsEnabled"] = v
	}
	if v := d.Get("filter"); !IsNil(v) {
		filter := make(map[string]interface{})
		storageClass1, _ := jsonpath.Get("$[0].storage_class", d.Get("filter"))
		if storageClass1 != nil && storageClass1 != "" {
			filter["StorageClass"] = storageClass1
		}
		lastModifyEndTimeStamp1, _ := jsonpath.Get("$[0].last_modify_end_time_stamp", d.Get("filter"))
		if lastModifyEndTimeStamp1 != nil && lastModifyEndTimeStamp1 != "" {
			filter["LastModifyEndTimeStamp"] = lastModifyEndTimeStamp1
		}
		upperSizeBound1, _ := jsonpath.Get("$[0].upper_size_bound", d.Get("filter"))
		if upperSizeBound1 != nil && upperSizeBound1 != "" {
			filter["UpperSizeBound"] = upperSizeBound1
		}
		prefix1, _ := jsonpath.Get("$[0].prefix", d.Get("filter"))
		if prefix1 != nil && prefix1 != "" {
			filter["Prefix"] = prefix1
		}
		lastModifyBeginTimeStamp1, _ := jsonpath.Get("$[0].last_modify_begin_time_stamp", d.Get("filter"))
		if lastModifyBeginTimeStamp1 != nil && lastModifyBeginTimeStamp1 != "" {
			filter["LastModifyBeginTimeStamp"] = lastModifyBeginTimeStamp1
		}
		lowerSizeBound1, _ := jsonpath.Get("$[0].lower_size_bound", d.Get("filter"))
		if lowerSizeBound1 != nil && lowerSizeBound1 != "" {
			filter["LowerSizeBound"] = lowerSizeBound1
		}

		if len(filter) > 0 {
			inventoryConfiguration["Filter"] = filter
		}
	}
	if v := d.Get("incremental_inventory"); !IsNil(v) {
		incrementalInventory := make(map[string]interface{})
		isEnabled3, _ := jsonpath.Get("$[0].is_enabled", d.Get("incremental_inventory"))
		if isEnabled3 != nil && isEnabled3 != "" {
			incrementalInventory["IsEnabled"] = isEnabled3
		}
		optionalFields1 := make(map[string]interface{})
		field3, _ := jsonpath.Get("$[0].optional_fields[0].field", d.Get("incremental_inventory"))
		if field3 != nil && field3 != "" {
			optionalFields1["Field"] = field3
		}

		if len(optionalFields1) > 0 {
			incrementalInventory["OptionalFields"] = optionalFields1
		}
		schedule := make(map[string]interface{})
		frequency1, _ := jsonpath.Get("$[0].schedule[0].frequency", d.Get("incremental_inventory"))
		if frequency1 != nil && frequency1 != "" {
			schedule["Frequency"] = frequency1
		}

		if len(schedule) > 0 {
			incrementalInventory["Schedule"] = schedule
		}

		if len(incrementalInventory) > 0 {
			inventoryConfiguration["IncrementalInventory"] = incrementalInventory
		}
	}
	if v, ok := d.GetOk("included_object_versions"); ok {
		inventoryConfiguration["IncludedObjectVersions"] = v
	}
	if v := d.Get("schedule"); !IsNil(v) {
		schedule1 := make(map[string]interface{})
		frequency3, _ := jsonpath.Get("$[0].frequency", d.Get("schedule"))
		if frequency3 != nil && frequency3 != "" {
			schedule1["Frequency"] = frequency3
		}

		if len(schedule1) > 0 {
			inventoryConfiguration["Schedule"] = schedule1
		}
	}
	if v := d.Get("destination"); !IsNil(v) {
		destination := make(map[string]interface{})
		oSSBucketDestination := make(map[string]interface{})
		encryption := make(map[string]interface{})
		sSEOSS1, _ := jsonpath.Get("$[0].oss_bucket_destination[0].encryption[0].sseoss", d.Get("destination"))
		if sSEOSS1 != nil && sSEOSS1 != "" {
			encryption["SSE-OSS"] = sSEOSS1
		}
		sSEKMS := make(map[string]interface{})
		keyId1, _ := jsonpath.Get("$[0].oss_bucket_destination[0].encryption[0].ssekms[0].key_id", d.Get("destination"))
		if keyId1 != nil && keyId1 != "" {
			sSEKMS["KeyId"] = keyId1
		}

		if len(sSEKMS) > 0 {
			encryption["SSE-KMS"] = sSEKMS
		}

		if len(encryption) > 0 {
			oSSBucketDestination["Encryption"] = encryption
		}
		prefix3, _ := jsonpath.Get("$[0].oss_bucket_destination[0].prefix", d.Get("destination"))
		if prefix3 != nil && prefix3 != "" {
			oSSBucketDestination["Prefix"] = prefix3
		}
		accountId1, _ := jsonpath.Get("$[0].oss_bucket_destination[0].account_id", d.Get("destination"))
		if accountId1 != nil && accountId1 != "" {
			oSSBucketDestination["AccountId"] = accountId1
		}
		format1, _ := jsonpath.Get("$[0].oss_bucket_destination[0].format", d.Get("destination"))
		if format1 != nil && format1 != "" {
			oSSBucketDestination["Format"] = format1
		}
		bucket1, _ := jsonpath.Get("$[0].oss_bucket_destination[0].bucket", d.Get("destination"))
		if bucket1 != nil && bucket1 != "" {
			oSSBucketDestination["Bucket"] = bucket1
		}
		roleArn1, _ := jsonpath.Get("$[0].oss_bucket_destination[0].role_arn", d.Get("destination"))
		if roleArn1 != nil && roleArn1 != "" {
			oSSBucketDestination["RoleArn"] = roleArn1
		}

		if len(oSSBucketDestination) > 0 {
			destination["OSSBucketDestination"] = oSSBucketDestination
		}

		if len(destination) > 0 {
			inventoryConfiguration["Destination"] = destination
		}
	}

	inventoryConfiguration["Id"] = d.Get("inventory_id")
	request["InventoryConfiguration"] = inventoryConfiguration

	body = request
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.Do("Oss", xmlParam("PUT", "2019-05-17", "PutBucketInventory", action), query, body, nil, hostMap, false)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_oss_bucket_inventory", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%v:%v", *hostMap["bucket"], *query["inventoryId"]))

	return resourceAliCloudOssBucketInventoryRead(d, meta)
}

func resourceAliCloudOssBucketInventoryRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ossServiceV2 := OssServiceV2{client}

	objectRaw, err := ossServiceV2.DescribeOssBucketInventory(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_oss_bucket_inventory DescribeOssBucketInventory Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("included_object_versions", objectRaw["IncludedObjectVersions"])
	d.Set("is_enabled", objectRaw["IsEnabled"])
	d.Set("inventory_id", objectRaw["Id"])

	destinationMaps := make([]map[string]interface{}, 0)
	destinationMap := make(map[string]interface{})
	destinationRaw := make(map[string]interface{})
	if objectRaw["Destination"] != nil {
		destinationRaw = objectRaw["Destination"].(map[string]interface{})
	}
	if len(destinationRaw) > 0 {

		oSSBucketDestinationMaps := make([]map[string]interface{}, 0)
		oSSBucketDestinationMap := make(map[string]interface{})
		oSSBucketDestinationRaw := make(map[string]interface{})
		if destinationRaw["OSSBucketDestination"] != nil {
			oSSBucketDestinationRaw = destinationRaw["OSSBucketDestination"].(map[string]interface{})
		}
		if len(oSSBucketDestinationRaw) > 0 {
			oSSBucketDestinationMap["account_id"] = oSSBucketDestinationRaw["AccountId"]
			oSSBucketDestinationMap["bucket"] = oSSBucketDestinationRaw["Bucket"]
			oSSBucketDestinationMap["format"] = oSSBucketDestinationRaw["Format"]
			oSSBucketDestinationMap["prefix"] = oSSBucketDestinationRaw["Prefix"]
			oSSBucketDestinationMap["role_arn"] = oSSBucketDestinationRaw["RoleArn"]

			encryptionMaps := make([]map[string]interface{}, 0)
			encryptionMap := make(map[string]interface{})
			encryptionRaw := make(map[string]interface{})
			if oSSBucketDestinationRaw["Encryption"] != nil {
				encryptionRaw = oSSBucketDestinationRaw["Encryption"].(map[string]interface{})
			}
			if len(encryptionRaw) > 0 {
				encryptionMap["sseoss"] = encryptionRaw["SSE-OSS"]

				sSEKMSMaps := make([]map[string]interface{}, 0)
				sSEKMSMap := make(map[string]interface{})
				sSEKMSRaw := make(map[string]interface{})
				if encryptionRaw["SSE-KMS"] != nil {
					sSEKMSRaw = encryptionRaw["SSE-KMS"].(map[string]interface{})
				}
				if len(sSEKMSRaw) > 0 {
					sSEKMSMap["key_id"] = sSEKMSRaw["KeyId"]

					sSEKMSMaps = append(sSEKMSMaps, sSEKMSMap)
				}
				encryptionMap["ssekms"] = sSEKMSMaps
				encryptionMaps = append(encryptionMaps, encryptionMap)
			}
			oSSBucketDestinationMap["encryption"] = encryptionMaps
			oSSBucketDestinationMaps = append(oSSBucketDestinationMaps, oSSBucketDestinationMap)
		}
		destinationMap["oss_bucket_destination"] = oSSBucketDestinationMaps
		destinationMaps = append(destinationMaps, destinationMap)
	}
	if err := d.Set("destination", destinationMaps); err != nil {
		return err
	}
	filterMaps := make([]map[string]interface{}, 0)
	filterMap := make(map[string]interface{})
	filterRaw := make(map[string]interface{})
	if objectRaw["Filter"] != nil {
		filterRaw = objectRaw["Filter"].(map[string]interface{})
	}
	// OSS's GetBucketInventory always echoes back a Filter block populated with
	// zero-value defaults even when the caller did not set one. Treat an
	// all-empty payload as "no filter" to avoid a perpetual plan diff.
	filterHasValue := false
	for _, k := range []string{"LastModifyBeginTimeStamp", "LastModifyEndTimeStamp", "LowerSizeBound", "Prefix", "StorageClass", "UpperSizeBound"} {
		if v, ok := filterRaw[k]; ok && v != nil && fmt.Sprint(v) != "" && fmt.Sprint(v) != "0" {
			filterHasValue = true
			break
		}
	}
	if filterHasValue {
		filterMap["last_modify_begin_time_stamp"] = filterRaw["LastModifyBeginTimeStamp"]
		filterMap["last_modify_end_time_stamp"] = filterRaw["LastModifyEndTimeStamp"]
		filterMap["lower_size_bound"] = filterRaw["LowerSizeBound"]
		filterMap["prefix"] = filterRaw["Prefix"]
		filterMap["storage_class"] = filterRaw["StorageClass"]
		filterMap["upper_size_bound"] = filterRaw["UpperSizeBound"]

		filterMaps = append(filterMaps, filterMap)
	}
	if err := d.Set("filter", filterMaps); err != nil {
		return err
	}
	incrementalInventoryMaps := make([]map[string]interface{}, 0)
	incrementalInventoryMap := make(map[string]interface{})
	incrementalInventoryRaw := make(map[string]interface{})
	if objectRaw["IncrementalInventory"] != nil {
		incrementalInventoryRaw = objectRaw["IncrementalInventory"].(map[string]interface{})
	}
	if len(incrementalInventoryRaw) > 0 {
		incrementalInventoryMap["is_enabled"] = incrementalInventoryRaw["IsEnabled"]

		optionalFieldsMaps := make([]map[string]interface{}, 0)
		optionalFieldsMap := make(map[string]interface{})
		optionalFieldsRaw := make(map[string]interface{})
		if incrementalInventoryRaw["OptionalFields"] != nil {
			optionalFieldsRaw = incrementalInventoryRaw["OptionalFields"].(map[string]interface{})
		}
		if len(optionalFieldsRaw) > 0 {

			fieldRaw := make([]interface{}, 0)
			if optionalFieldsRaw["Field"] != nil {
				fieldRaw = convertToInterfaceArray(optionalFieldsRaw["Field"])
			}

			optionalFieldsMap["field"] = fieldRaw
			optionalFieldsMaps = append(optionalFieldsMaps, optionalFieldsMap)
		}
		incrementalInventoryMap["optional_fields"] = optionalFieldsMaps
		scheduleMaps := make([]map[string]interface{}, 0)
		scheduleMap := make(map[string]interface{})
		scheduleRaw := make(map[string]interface{})
		if incrementalInventoryRaw["Schedule"] != nil {
			scheduleRaw = incrementalInventoryRaw["Schedule"].(map[string]interface{})
		}
		if len(scheduleRaw) > 0 {
			scheduleMap["frequency"] = scheduleRaw["Frequency"]

			scheduleMaps = append(scheduleMaps, scheduleMap)
		}
		incrementalInventoryMap["schedule"] = scheduleMaps
		incrementalInventoryMaps = append(incrementalInventoryMaps, incrementalInventoryMap)
	}
	if err := d.Set("incremental_inventory", incrementalInventoryMaps); err != nil {
		return err
	}
	optionalFieldsMaps := make([]map[string]interface{}, 0)
	optionalFieldsMap := make(map[string]interface{})
	optionalFieldsRaw := make(map[string]interface{})
	if objectRaw["OptionalFields"] != nil {
		optionalFieldsRaw = objectRaw["OptionalFields"].(map[string]interface{})
	}
	if len(optionalFieldsRaw) > 0 {

		fieldRaw := make([]interface{}, 0)
		if optionalFieldsRaw["Field"] != nil {
			fieldRaw = convertToInterfaceArray(optionalFieldsRaw["Field"])
		}

		optionalFieldsMap["field"] = fieldRaw
		optionalFieldsMaps = append(optionalFieldsMaps, optionalFieldsMap)
	}
	if err := d.Set("optional_fields", optionalFieldsMaps); err != nil {
		return err
	}
	scheduleMaps := make([]map[string]interface{}, 0)
	scheduleMap := make(map[string]interface{})
	scheduleRaw := make(map[string]interface{})
	if objectRaw["Schedule"] != nil {
		scheduleRaw = objectRaw["Schedule"].(map[string]interface{})
	}
	if len(scheduleRaw) > 0 {
		scheduleMap["frequency"] = scheduleRaw["Frequency"]

		scheduleMaps = append(scheduleMaps, scheduleMap)
	}
	if err := d.Set("schedule", scheduleMaps); err != nil {
		return err
	}

	parts := strings.Split(d.Id(), ":")
	d.Set("bucket", parts[0])

	return nil
}

func resourceAliCloudOssBucketInventoryUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	var body map[string]interface{}
	update := false

	var err error
	parts := strings.Split(d.Id(), ":")
	action := fmt.Sprintf("/?inventory")
	request = make(map[string]interface{})
	query = make(map[string]*string)
	body = make(map[string]interface{})
	hostMap := make(map[string]*string)
	hostMap["bucket"] = StringPointer(parts[0])
	query["inventoryId"] = StringPointer(parts[1])

	inventoryConfiguration := make(map[string]interface{})
	if d.HasChange("optional_fields") {
		update = true
	}
	if v := d.Get("optional_fields"); !IsNil(v) || d.HasChange("optional_fields") {
		optionalFields := make(map[string]interface{})
		field1, _ := jsonpath.Get("$[0].field", d.Get("optional_fields"))
		if field1 != nil && field1 != "" {
			optionalFields["Field"] = field1
		}

		if len(optionalFields) > 0 {
			inventoryConfiguration["OptionalFields"] = optionalFields
		}
	}
	if d.HasChange("is_enabled") {
		update = true
	}
	if v, ok := d.GetOkExists("is_enabled"); ok {
		inventoryConfiguration["IsEnabled"] = v
	}
	if d.HasChange("filter") {
		update = true
	}
	if v := d.Get("filter"); !IsNil(v) || d.HasChange("filter") {
		filter := make(map[string]interface{})
		storageClass1, _ := jsonpath.Get("$[0].storage_class", d.Get("filter"))
		if storageClass1 != nil && storageClass1 != "" {
			filter["StorageClass"] = storageClass1
		}
		lastModifyEndTimeStamp1, _ := jsonpath.Get("$[0].last_modify_end_time_stamp", d.Get("filter"))
		if lastModifyEndTimeStamp1 != nil && lastModifyEndTimeStamp1 != "" {
			filter["LastModifyEndTimeStamp"] = lastModifyEndTimeStamp1
		}
		upperSizeBound1, _ := jsonpath.Get("$[0].upper_size_bound", d.Get("filter"))
		if upperSizeBound1 != nil && upperSizeBound1 != "" {
			filter["UpperSizeBound"] = upperSizeBound1
		}
		prefix1, _ := jsonpath.Get("$[0].prefix", d.Get("filter"))
		if prefix1 != nil && prefix1 != "" {
			filter["Prefix"] = prefix1
		}
		lastModifyBeginTimeStamp1, _ := jsonpath.Get("$[0].last_modify_begin_time_stamp", d.Get("filter"))
		if lastModifyBeginTimeStamp1 != nil && lastModifyBeginTimeStamp1 != "" {
			filter["LastModifyBeginTimeStamp"] = lastModifyBeginTimeStamp1
		}
		lowerSizeBound1, _ := jsonpath.Get("$[0].lower_size_bound", d.Get("filter"))
		if lowerSizeBound1 != nil && lowerSizeBound1 != "" {
			filter["LowerSizeBound"] = lowerSizeBound1
		}

		if len(filter) > 0 {
			inventoryConfiguration["Filter"] = filter
		}
	}
	if d.HasChange("incremental_inventory") {
		update = true
	}
	if v := d.Get("incremental_inventory"); !IsNil(v) || d.HasChange("incremental_inventory") {
		incrementalInventory := make(map[string]interface{})
		isEnabled3, _ := jsonpath.Get("$[0].is_enabled", d.Get("incremental_inventory"))
		if isEnabled3 != nil && isEnabled3 != "" {
			incrementalInventory["IsEnabled"] = isEnabled3
		}
		optionalFields1 := make(map[string]interface{})
		field3, _ := jsonpath.Get("$[0].optional_fields[0].field", d.Get("incremental_inventory"))
		if field3 != nil && field3 != "" {
			optionalFields1["Field"] = field3
		}

		if len(optionalFields1) > 0 {
			incrementalInventory["OptionalFields"] = optionalFields1
		}
		schedule := make(map[string]interface{})
		frequency1, _ := jsonpath.Get("$[0].schedule[0].frequency", d.Get("incremental_inventory"))
		if frequency1 != nil && frequency1 != "" {
			schedule["Frequency"] = frequency1
		}

		if len(schedule) > 0 {
			incrementalInventory["Schedule"] = schedule
		}

		if len(incrementalInventory) > 0 {
			inventoryConfiguration["IncrementalInventory"] = incrementalInventory
		}
	}
	if d.HasChange("included_object_versions") {
		update = true
	}
	if v, ok := d.GetOk("included_object_versions"); ok {
		inventoryConfiguration["IncludedObjectVersions"] = v
	}
	if d.HasChange("schedule") {
		update = true
	}
	if v := d.Get("schedule"); !IsNil(v) || d.HasChange("schedule") {
		schedule1 := make(map[string]interface{})
		frequency3, _ := jsonpath.Get("$[0].frequency", d.Get("schedule"))
		if frequency3 != nil && frequency3 != "" {
			schedule1["Frequency"] = frequency3
		}

		if len(schedule1) > 0 {
			inventoryConfiguration["Schedule"] = schedule1
		}
	}
	if d.HasChange("destination") {
		update = true
	}
	if v := d.Get("destination"); !IsNil(v) || d.HasChange("destination") {
		destination := make(map[string]interface{})
		oSSBucketDestination := make(map[string]interface{})
		encryption := make(map[string]interface{})
		sSEOSS1, _ := jsonpath.Get("$[0].oss_bucket_destination[0].encryption[0].sseoss", d.Get("destination"))
		if sSEOSS1 != nil && sSEOSS1 != "" {
			encryption["SSE-OSS"] = sSEOSS1
		}
		sSEKMS := make(map[string]interface{})
		keyId1, _ := jsonpath.Get("$[0].oss_bucket_destination[0].encryption[0].ssekms[0].key_id", d.Get("destination"))
		if keyId1 != nil && keyId1 != "" {
			sSEKMS["KeyId"] = keyId1
		}

		if len(sSEKMS) > 0 {
			encryption["SSE-KMS"] = sSEKMS
		}

		if len(encryption) > 0 {
			oSSBucketDestination["Encryption"] = encryption
		}
		prefix3, _ := jsonpath.Get("$[0].oss_bucket_destination[0].prefix", d.Get("destination"))
		if prefix3 != nil && prefix3 != "" {
			oSSBucketDestination["Prefix"] = prefix3
		}
		accountId1, _ := jsonpath.Get("$[0].oss_bucket_destination[0].account_id", d.Get("destination"))
		if accountId1 != nil && accountId1 != "" {
			oSSBucketDestination["AccountId"] = accountId1
		}
		format1, _ := jsonpath.Get("$[0].oss_bucket_destination[0].format", d.Get("destination"))
		if format1 != nil && format1 != "" {
			oSSBucketDestination["Format"] = format1
		}
		bucket1, _ := jsonpath.Get("$[0].oss_bucket_destination[0].bucket", d.Get("destination"))
		if bucket1 != nil && bucket1 != "" {
			oSSBucketDestination["Bucket"] = bucket1
		}
		roleArn1, _ := jsonpath.Get("$[0].oss_bucket_destination[0].role_arn", d.Get("destination"))
		if roleArn1 != nil && roleArn1 != "" {
			oSSBucketDestination["RoleArn"] = roleArn1
		}

		if len(oSSBucketDestination) > 0 {
			destination["OSSBucketDestination"] = oSSBucketDestination
		}

		if len(destination) > 0 {
			inventoryConfiguration["Destination"] = destination
		}
	}

	inventoryConfiguration["Id"] = d.Get("inventory_id")
	request["InventoryConfiguration"] = inventoryConfiguration

	body = request
	if update {
		// OSS PutBucketInventory is create-only; if the rule already exists it
		// returns 409 InventoryConfigurationAlreadyExists. Delete the existing
		// rule first, then re-create with the new configuration.
		deleteAction := fmt.Sprintf("/?inventory&inventoryId=%s", parts[1])
		waitDel := incrementalWait(3*time.Second, 5*time.Second)
		errDel := resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			_, err = client.Do("Oss", xmlParam("DELETE", "2019-05-17", "DeleteBucketInventory", deleteAction), query, nil, nil, hostMap, false)
			if err != nil {
				if NeedRetry(err) {
					waitDel()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		if errDel != nil && !IsExpectedErrors(errDel, []string{"NoSuchInventoryConfiguration"}) && !NotFoundError(errDel) {
			return WrapErrorf(errDel, DefaultErrorMsg, d.Id(), deleteAction, AlibabaCloudSdkGoERROR)
		}

		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.Do("Oss", xmlParam("PUT", "2019-05-17", "PutBucketInventory", action), query, body, nil, hostMap, false)
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

	return resourceAliCloudOssBucketInventoryRead(d, meta)
}

func resourceAliCloudOssBucketInventoryDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	parts := strings.Split(d.Id(), ":")
	action := fmt.Sprintf("/?inventory")
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	hostMap := make(map[string]*string)
	var err error
	request = make(map[string]interface{})
	hostMap["bucket"] = StringPointer(parts[0])
	query["inventoryId"] = StringPointer(parts[1])

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.Do("Oss", xmlParam("DELETE", "2019-05-17", "DeleteBucketInventory", action), query, nil, nil, hostMap, false)
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

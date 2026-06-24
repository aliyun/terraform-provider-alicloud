// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"encoding/xml"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

// xmlEscape returns v rendered as an XML-escaped text node value.
func xmlEscape(v interface{}) string {
	return XmlEscape(fmt.Sprintf("%v", v))
}

// XmlEscape escapes the characters that are not allowed in an XML text node.
func XmlEscape(s string) string {
	var b strings.Builder
	_ = xml.EscapeText(&b, []byte(s))
	return b.String()
}

// buildOssBucketInventoryConfigurationXML assembles the inner elements of the
// InventoryConfiguration body in the exact order required by OSS PutBucketInventory's
// published schema. The <InventoryConfiguration> wrapper is added by the tea-xml
// serializer from the map key, so it is not included here. A map cannot be used for the
// body because Go map iteration order is random and OSS rejects an out-of-order body
// with 400 MalformedXML.
func buildOssBucketInventoryConfigurationXML(d *schema.ResourceData, inventoryId string) string {
	var b strings.Builder
	b.WriteString("<Id>" + xmlEscape(inventoryId) + "</Id>")

	if v, ok := d.GetOkExists("is_enabled"); ok {
		b.WriteString(fmt.Sprintf("<IsEnabled>%v</IsEnabled>", v))
	}

	if v, ok := d.GetOk("filter"); ok && len(v.([]interface{})) > 0 && v.([]interface{})[0] != nil {
		f := v.([]interface{})[0].(map[string]interface{})
		var fb strings.Builder
		if s, _ := f["prefix"].(string); s != "" {
			fb.WriteString("<Prefix>" + xmlEscape(s) + "</Prefix>")
		}
		if n, _ := f["last_modify_begin_time_stamp"].(int); n > 0 {
			fb.WriteString(fmt.Sprintf("<LastModifyBeginTimeStamp>%d</LastModifyBeginTimeStamp>", n))
		}
		if n, _ := f["last_modify_end_time_stamp"].(int); n > 0 {
			fb.WriteString(fmt.Sprintf("<LastModifyEndTimeStamp>%d</LastModifyEndTimeStamp>", n))
		}
		if n, _ := f["lower_size_bound"].(int); n > 0 {
			fb.WriteString(fmt.Sprintf("<LowerSizeBound>%d</LowerSizeBound>", n))
		}
		if n, _ := f["upper_size_bound"].(int); n > 0 {
			fb.WriteString(fmt.Sprintf("<UpperSizeBound>%d</UpperSizeBound>", n))
		}
		if s, _ := f["storage_class"].(string); s != "" {
			fb.WriteString("<StorageClass>" + xmlEscape(s) + "</StorageClass>")
		}
		if fb.Len() > 0 {
			b.WriteString("<Filter>" + fb.String() + "</Filter>")
		}
	}

	if v, ok := d.GetOk("destination"); ok && len(v.([]interface{})) > 0 && v.([]interface{})[0] != nil {
		dest := v.([]interface{})[0].(map[string]interface{})
		if ds, ok := dest["oss_bucket_destination"].([]interface{}); ok && len(ds) > 0 && ds[0] != nil {
			od := ds[0].(map[string]interface{})
			var ob strings.Builder
			if s, _ := od["format"].(string); s != "" {
				ob.WriteString("<Format>" + xmlEscape(s) + "</Format>")
			}
			if s, _ := od["account_id"].(string); s != "" {
				ob.WriteString("<AccountId>" + xmlEscape(s) + "</AccountId>")
			}
			if s, _ := od["role_arn"].(string); s != "" {
				ob.WriteString("<RoleArn>" + xmlEscape(s) + "</RoleArn>")
			}
			if s, _ := od["bucket"].(string); s != "" {
				ob.WriteString("<Bucket>" + xmlEscape(s) + "</Bucket>")
			}
			if s, _ := od["prefix"].(string); s != "" {
				ob.WriteString("<Prefix>" + xmlEscape(s) + "</Prefix>")
			}
			if enc, ok := od["encryption"].([]interface{}); ok && len(enc) > 0 && enc[0] != nil {
				em := enc[0].(map[string]interface{})
				if s, _ := em["key_id"].(string); s != "" {
					ob.WriteString("<Encryption><SSE-KMS><KeyId>" + xmlEscape(s) + "</KeyId></SSE-KMS></Encryption>")
				}
			}
			if ob.Len() > 0 {
				b.WriteString("<Destination><OSSBucketDestination>" + ob.String() + "</OSSBucketDestination></Destination>")
			}
		}
	}

	if v, ok := d.GetOk("schedule"); ok && len(v.([]interface{})) > 0 && v.([]interface{})[0] != nil {
		s := v.([]interface{})[0].(map[string]interface{})
		if f, _ := s["frequency"].(string); f != "" {
			b.WriteString("<Schedule><Frequency>" + xmlEscape(f) + "</Frequency></Schedule>")
		}
	}

	if v, ok := d.GetOk("included_object_versions"); ok {
		b.WriteString("<IncludedObjectVersions>" + xmlEscape(v) + "</IncludedObjectVersions>")
	}

	if v, ok := d.GetOk("optional_fields"); ok && len(v.([]interface{})) > 0 && v.([]interface{})[0] != nil {
		of := v.([]interface{})[0].(map[string]interface{})
		if fields, ok := of["field"].([]interface{}); ok && len(fields) > 0 {
			b.WriteString("<OptionalFields>")
			for _, f := range fields {
				b.WriteString("<Field>" + xmlEscape(f) + "</Field>")
			}
			b.WriteString("</OptionalFields>")
		}
	}

	if v, ok := d.GetOk("incremental_inventory"); ok && len(v.([]interface{})) > 0 && v.([]interface{})[0] != nil {
		ii := v.([]interface{})[0].(map[string]interface{})
		var ib strings.Builder
		if be, ok := ii["is_enabled"].(bool); ok {
			ib.WriteString(fmt.Sprintf("<IsEnabled>%v</IsEnabled>", be))
		}
		if of, ok := ii["optional_fields"].([]interface{}); ok && len(of) > 0 && of[0] != nil {
			ofm := of[0].(map[string]interface{})
			if fields, ok := ofm["field"].([]interface{}); ok && len(fields) > 0 {
				ib.WriteString("<OptionalFields>")
				for _, f := range fields {
					ib.WriteString("<Field>" + xmlEscape(f) + "</Field>")
				}
				ib.WriteString("</OptionalFields>")
			}
		}
		if sc, ok := ii["schedule"].([]interface{}); ok && len(sc) > 0 && sc[0] != nil {
			scm := sc[0].(map[string]interface{})
			if f, _ := scm["frequency"].(string); f != "" {
				ib.WriteString("<Schedule><Frequency>" + xmlEscape(f) + "</Frequency></Schedule>")
			}
		}
		if ib.Len() > 0 {
			b.WriteString("<IncrementalInventory>" + ib.String() + "</IncrementalInventory>")
		}
	}

	return b.String()
}

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
												"key_id": {
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
										Type:     schema.TypeString,
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

	// OSS PutBucketInventory validates the request body against a strict, ordered
	// schema. A Go map serializes in random order, so the configuration is rendered
	// as a deterministic XML string in canonical element order.
	request["InventoryConfiguration"] = buildOssBucketInventoryConfigurationXML(d, d.Get("inventory_id").(string))

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
		addDebug(action, response, request)
		return nil
	})

	if err != nil && !IsExpectedErrors(err, []string{"InventoryConfigurationAlreadyExists"}) {
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

	inventoryConfigurationRawObj, _ := jsonpath.Get("$.InventoryConfiguration", objectRaw)
	inventoryConfigurationRaw := make(map[string]interface{})
	if inventoryConfigurationRawObj != nil {
		inventoryConfigurationRaw = inventoryConfigurationRawObj.(map[string]interface{})
	}
	d.Set("included_object_versions", inventoryConfigurationRaw["IncludedObjectVersions"])
	if v, ok := inventoryConfigurationRaw["IsEnabled"]; ok {
		d.Set("is_enabled", formatBool(v))
	}

	destinationMaps := make([]map[string]interface{}, 0)
	destinationMap := make(map[string]interface{})
	destinationRawObj, _ := jsonpath.Get("$.InventoryConfiguration.Destination", objectRaw)
	destinationRaw := make(map[string]interface{})
	if destinationRawObj != nil {
		destinationRaw = destinationRawObj.(map[string]interface{})
	}
	if len(destinationRaw) > 0 {

		oSSBucketDestinationMaps := make([]map[string]interface{}, 0)
		oSSBucketDestinationMap := make(map[string]interface{})
		oSSBucketDestinationRawObj, _ := jsonpath.Get("$.InventoryConfiguration.Destination.OSSBucketDestination", objectRaw)
		oSSBucketDestinationRaw := make(map[string]interface{})
		if oSSBucketDestinationRawObj != nil {
			oSSBucketDestinationRaw = oSSBucketDestinationRawObj.(map[string]interface{})
		}
		if len(oSSBucketDestinationRaw) > 0 {
			oSSBucketDestinationMap["account_id"] = oSSBucketDestinationRaw["AccountId"]
			oSSBucketDestinationMap["bucket"] = oSSBucketDestinationRaw["Bucket"]
			oSSBucketDestinationMap["format"] = oSSBucketDestinationRaw["Format"]
			oSSBucketDestinationMap["prefix"] = oSSBucketDestinationRaw["Prefix"]
			oSSBucketDestinationMap["role_arn"] = oSSBucketDestinationRaw["RoleArn"]

			encryptionMaps := make([]map[string]interface{}, 0)
			encryptionMap := make(map[string]interface{})
			encryptionRawObj, _ := jsonpath.Get("$.InventoryConfiguration.Destination.OSSBucketDestination.Encryption", objectRaw)
			encryptionRaw := make(map[string]interface{})
			if encryptionRawObj != nil {
				encryptionRaw = encryptionRawObj.(map[string]interface{})
			}
			if len(encryptionRaw) > 0 {
				encryptionMap["key_id"] = encryptionRaw["KeyId"]

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
	filterRawObj, _ := jsonpath.Get("$.InventoryConfiguration.Filter", objectRaw)
	filterRaw := make(map[string]interface{})
	if filterRawObj != nil {
		filterRaw = filterRawObj.(map[string]interface{})
	}
	if len(filterRaw) > 0 {
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
	incrementalInventoryRawObj, _ := jsonpath.Get("$.InventoryConfiguration.IncrementalInventory", objectRaw)
	incrementalInventoryRaw := make(map[string]interface{})
	if incrementalInventoryRawObj != nil {
		incrementalInventoryRaw = incrementalInventoryRawObj.(map[string]interface{})
	}
	if len(incrementalInventoryRaw) > 0 {
		incrementalInventoryMap["is_enabled"] = incrementalInventoryRaw["IsEnabled"]

		optionalFieldsMaps := make([]map[string]interface{}, 0)
		optionalFieldsMap := make(map[string]interface{})
		optionalFieldsRawObj, _ := jsonpath.Get("$.InventoryConfiguration.IncrementalInventory.OptionalFields", objectRaw)
		optionalFieldsRaw := make(map[string]interface{})
		if optionalFieldsRawObj != nil {
			optionalFieldsRaw = optionalFieldsRawObj.(map[string]interface{})
		}
		if len(optionalFieldsRaw) > 0 {

			fieldRaw, _ := jsonpath.Get("$.InventoryConfiguration.IncrementalInventory.OptionalFields.Field", objectRaw)
			optionalFieldsMap["field"] = fieldRaw
			optionalFieldsMaps = append(optionalFieldsMaps, optionalFieldsMap)
		}
		incrementalInventoryMap["optional_fields"] = optionalFieldsMaps
		scheduleMaps := make([]map[string]interface{}, 0)
		scheduleMap := make(map[string]interface{})
		scheduleRawObj, _ := jsonpath.Get("$.InventoryConfiguration.IncrementalInventory.Schedule", objectRaw)
		scheduleRaw := make(map[string]interface{})
		if scheduleRawObj != nil {
			scheduleRaw = scheduleRawObj.(map[string]interface{})
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
	optionalFieldsRawObj, _ := jsonpath.Get("$.InventoryConfiguration.OptionalFields", objectRaw)
	optionalFieldsRaw := make(map[string]interface{})
	if optionalFieldsRawObj != nil {
		optionalFieldsRaw = optionalFieldsRawObj.(map[string]interface{})
	}
	if len(optionalFieldsRaw) > 0 {

		fieldRaw, _ := jsonpath.Get("$.InventoryConfiguration.OptionalFields.Field", objectRaw)
		optionalFieldsMap["field"] = fieldRaw
		optionalFieldsMaps = append(optionalFieldsMaps, optionalFieldsMap)
	}
	if err := d.Set("optional_fields", optionalFieldsMaps); err != nil {
		return err
	}
	scheduleMaps := make([]map[string]interface{}, 0)
	scheduleMap := make(map[string]interface{})
	scheduleRawObj, _ := jsonpath.Get("$.InventoryConfiguration.Schedule", objectRaw)
	scheduleRaw := make(map[string]interface{})
	if scheduleRawObj != nil {
		scheduleRaw = scheduleRawObj.(map[string]interface{})
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
	d.Set("inventory_id", parts[1])

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

	if d.HasChange("optional_fields") || d.HasChange("is_enabled") || d.HasChange("filter") ||
		d.HasChange("incremental_inventory") || d.HasChange("included_object_versions") ||
		d.HasChange("schedule") || d.HasChange("destination") {
		update = true
	}

	// Render the configuration as a deterministic, schema-ordered XML string; a map
	// would serialize in random order and trip OSS's MalformedXML schema validation.
	request["InventoryConfiguration"] = buildOssBucketInventoryConfigurationXML(d, parts[1])

	body = request
	if update {
		// PutBucketInventory rejects an existing configuration with InventoryConfigurationAlreadyExists,
		// so an update must remove the old rule before writing the new one.
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.Do("Oss", xmlParam("DELETE", "2019-05-17", "DeleteBucketInventory", action), query, nil, nil, hostMap, false)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(action, response, nil)
			return nil
		})
		if err != nil && !IsExpectedErrors(err, []string{"NoSuchInventory", "NoSuchBucket"}) {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}

		wait = incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.Do("Oss", xmlParam("PUT", "2019-05-17", "PutBucketInventory", action), query, body, nil, hostMap, false)
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
		if err != nil && !IsExpectedErrors(err, []string{"InventoryConfigurationAlreadyExists"}) {
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
		response, err = client.Do("Oss", xmlParam("DELETE", "2019-05-17", "DeleteBucketInventory", action), query, request, nil, hostMap, false)
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
		if IsExpectedErrors(err, []string{"NoSuchInventory", "NoSuchBucket"}) || NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return nil
}

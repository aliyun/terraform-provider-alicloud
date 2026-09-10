// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAliCloudOssBucketInventories() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAliCloudOssBucketInventoryRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"bucket": {
				Type:     schema.TypeString,
				Required: true,
			},
			"inventories": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"destination": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"oss_bucket_destination": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"format": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"account_id": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"bucket": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"prefix": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"encryption": {
													Type:     schema.TypeList,
													Computed: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"ssekms": {
																Type:     schema.TypeList,
																Computed: true,
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"key_id": {
																			Type:     schema.TypeString,
																			Computed: true,
																		},
																	},
																},
															},
															"sseoss": {
																Type:     schema.TypeString,
																Computed: true,
															},
														},
													},
												},
												"role_arn": {
													Type:     schema.TypeString,
													Computed: true,
												},
											},
										},
									},
								},
							},
						},
						"filter": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"lower_size_bound": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"storage_class": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"last_modify_begin_time_stamp": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"last_modify_end_time_stamp": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"upper_size_bound": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"prefix": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"included_object_versions": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"incremental_inventory": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"optional_fields": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"field": {
													Type:     schema.TypeList,
													Computed: true,
													Elem:     &schema.Schema{Type: schema.TypeString},
												},
											},
										},
									},
									"is_enabled": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"schedule": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"frequency": {
													Type:     schema.TypeInt,
													Computed: true,
												},
											},
										},
									},
								},
							},
						},
						"inventory_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"is_enabled": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"optional_fields": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"field": {
										Type:     schema.TypeList,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
						"schedule": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"frequency": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func dataSourceAliCloudOssBucketInventoryRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	var objects []map[string]interface{}

	idsMap := make(map[string]string)
	if v, ok := d.GetOk("ids"); ok {
		for _, vv := range v.([]interface{}) {
			if vv == nil {
				continue
			}
			idsMap[vv.(string)] = vv.(string)
		}
	}

	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	// ListBucketInventory
	action := fmt.Sprintf("/?inventory")
	var err error
	request = make(map[string]interface{})
	query = make(map[string]*string)
	hostMap := make(map[string]*string)

	hostMap["bucket"] = StringPointer(d.Get("bucket").(string))
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutRead), func() *resource.RetryError {
		response, err = client.Do("Oss", xmlParam("GET", "2019-05-17", "ListBucketInventory", action), query, nil, nil, hostMap, false)

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

	resp, _ := jsonpath.Get("$.ListInventoryConfigurationsResult.InventoryConfiguration[*]", response)

	result, _ := resp.([]interface{})
	for _, v := range result {
		item := v.(map[string]interface{})
		if len(idsMap) > 0 {
			if _, ok := idsMap[fmt.Sprint(*hostMap["bucket"], ":", item["Id"])]; !ok {
				continue
			}
		}
		objects = append(objects, item)
	}

	ids := make([]string, 0)
	s := make([]map[string]interface{}, 0)
	for _, objectRaw := range objects {
		mapping := map[string]interface{}{}

		mapping["id"] = fmt.Sprint(*hostMap["bucket"], ":", objectRaw["Id"])

		mapping["included_object_versions"] = objectRaw["IncludedObjectVersions"]
		mapping["is_enabled"] = objectRaw["IsEnabled"]
		mapping["inventory_id"] = objectRaw["Id"]

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
		mapping["destination"] = destinationMaps
		filterMaps := make([]map[string]interface{}, 0)
		filterMap := make(map[string]interface{})
		filterRaw := make(map[string]interface{})
		if objectRaw["Filter"] != nil {
			filterRaw = objectRaw["Filter"].(map[string]interface{})
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
		mapping["filter"] = filterMaps
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
		mapping["incremental_inventory"] = incrementalInventoryMaps
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
		mapping["optional_fields"] = optionalFieldsMaps
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
		mapping["schedule"] = scheduleMaps

		ids = append(ids, fmt.Sprint(mapping["id"]))
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("inventories", s); err != nil {
		return WrapError(err)
	}

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}
